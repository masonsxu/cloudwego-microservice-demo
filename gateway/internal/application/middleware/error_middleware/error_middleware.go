package middleware

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	"github.com/rs/zerolog"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/context/auth_context"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/errors"
	tracelog "github.com/masonsxu/cloudwego-microservice-demo/gateway/pkg/log"
)

// ErrorHandlerMiddlewareImpl 错误处理中间件实现
type ErrorHandlerMiddlewareImpl struct {
	config *config.ErrorHandlerConfig
	logger *zerolog.Logger
}

// NewErrorHandlerMiddleware 创建新的错误处理中间件实例
func NewErrorHandlerMiddleware(
	config *config.ErrorHandlerConfig,
	logger *hertzZerolog.Logger,
) ErrorHandlerMiddlewareService {
	var zlog *zerolog.Logger

	if logger != nil {
		// 从 hertzZerolog.Logger 获取底层 zerolog.Logger
		unwrapped := logger.Unwrap()
		zlog = &unwrapped
	} else {
		// 如果 logger 为 nil，创建一个默认的 nop logger
		nopLogger := zerolog.Nop()
		zlog = &nopLogger
	}

	return &ErrorHandlerMiddlewareImpl{
		config: config,
		logger: zlog,
	}
}

// MiddlewareFunc 返回Hertz中间件处理函数
func (ehm *ErrorHandlerMiddlewareImpl) MiddlewareFunc() app.HandlerFunc {
	return app.HandlerFunc(func(ctx context.Context, c *app.RequestContext) {
		// 如果中间件被禁用，直接跳过
		if !ehm.config.Enabled {
			c.Next(ctx)
			return
		}

		// 记录请求开始时间
		startTime := time.Now()

		// 如果启用panic恢复，设置recover
		if ehm.config.EnablePanicRecovery {
			defer func() {
				if r := recover(); r != nil {
					ehm.handlePanicError(ctx, c, r)
				}
			}()
		}

		// 记录请求信息
		if ehm.config.EnableRequestLogging {
			ehm.logRequestInfo(ctx, c)
		}

		// 执行下一个处理器
		c.Next(ctx)

		statusCode := c.Response.StatusCode()

		if statusCode == 500 && len(c.Response.Body()) == 0 {
			tracelog.Event(ctx, ehm.logger.Warn()).
				Str("component", "error_middleware").
				Str("method", string(c.Method())).
				Str("path", string(c.Request.Path())).
				Msg("Permission denied")
			c.SetStatusCode(403)
			c.JSON(403, map[string]string{"msg": "权限不足"})

			return
		}

		// 检查响应状态码并记录错误日志
		// 注意：不能依赖 c.IsAborted() 判断，因为 AbortWithError 会设置 Abort 标志，
		// 但我们仍然需要对所有 4xx/5xx 响应进行日志和 span 上报
		if statusCode >= 400 {
			ehm.logHTTPError(ctx, c, statusCode)
			return
		}

		// 记录响应信息
		if ehm.config.EnableResponseLogging {
			ehm.logResponseInfo(ctx, c, startTime)
		}
	})
}

// handlePanicError 处理panic错误
func (ehm *ErrorHandlerMiddlewareImpl) handlePanicError(
	ctx context.Context,
	c *app.RequestContext,
	r any,
) {
	// 获取堆栈信息
	bufSize := ehm.config.MaxStackTraceSize
	if bufSize <= 0 {
		bufSize = 4096 // 默认值
	}

	buf := make([]byte, bufSize)
	n := runtime.Stack(buf, false)
	stackTrace := string(buf[:n])

	// 使用结构化日志记录 panic 详情
	tracelog.Event(ctx, ehm.logger.Error()).
		Str("component", "error_middleware").
		Str("path", string(c.Request.Path())).
		Str("method", string(c.Method())).
		Str("user_agent", string(c.UserAgent())).
		Str("remote_addr", c.RemoteAddr().String()).
		Interface("panic", r).
		Str("stack_trace", stackTrace).
		Msg("Panic recovered in error handler")

	// 构建错误响应
	bizErr := errors.ErrInternal
	if ehm.config.EnableDetailedErrors {
		bizErr = bizErr.WithMessage(fmt.Sprintf("内部服务器错误：%v", r))
	}

	// 使用统一的错误处理函数
	errors.AbortWithError(c, bizErr)
}

// logHTTPError 记录 HTTP 错误日志并上报到 OTel Span
// 4xx 客户端错误使用 Warn 级别（span.AddEvent），5xx 服务端错误使用 Error 级别（span.RecordError + span 标红）
// 不修改响应体——响应已由 AbortWithError 写入具体错误信息
func (ehm *ErrorHandlerMiddlewareImpl) logHTTPError(
	ctx context.Context,
	c *app.RequestContext,
	statusCode int,
) {
	// 从已写入的响应体中提取错误信息（避免信息丢失）
	responseBody := string(c.Response.Body())

	var event *zerolog.Event
	if statusCode >= 500 {
		// 5xx 服务端错误：Error 级别 → OTelHook 调用 span.RecordError() + span.SetStatus(Error)，Jaeger 标红
		event = ehm.logger.Error()
	} else {
		// 4xx 客户端错误：Warn 级别 → OTelHook 调用 span.AddEvent()，Jaeger span 详情中可见
		event = ehm.logger.Warn()
	}

	event = tracelog.Event(ctx, event).
		Str("component", "error_middleware").
		Int("status_code", statusCode).
		Str("path", string(c.Request.Path())).
		Str("method", string(c.Method())).
		Str("response_body", responseBody)

	if userID, ok := auth_context.GetCurrentUserProfileID(c); ok && userID != "" {
		event = event.Str("user_id", userID)
	}

	if orgID, ok := auth_context.GetCurrentOrganizationID(c); ok && orgID != "" {
		event = event.Str("org_id", orgID)
	}

	if ehm.config.EnableDetailedErrors {
		event = event.
			Str("remote_addr", c.RemoteAddr().String()).
			Str("user_agent", string(c.UserAgent()))
	}

	event.Msg("HTTP error response")

	// OTelHook 只将日志消息传给 span，结构化字段不会成为 span attributes。
	// 直接写入 span attributes，使 response_body 等信息在 Jaeger 中可见。
	tracelog.RecordSpanHTTPError(ctx, statusCode, string(c.Method()), string(c.Request.Path()), responseBody)
}

// handleHTTPStatusError 处理HTTP状态码错误
// Deprecated: 改用 logHTTPError，此方法仅保留以兼容 panic 恢复路径
func (ehm *ErrorHandlerMiddlewareImpl) handleHTTPStatusError(
	ctx context.Context,
	c *app.RequestContext,
	statusCode int,
) {
	ehm.logHTTPError(ctx, c, statusCode)
}

// logRequestInfo 记录请求信息
func (ehm *ErrorHandlerMiddlewareImpl) logRequestInfo(ctx context.Context, c *app.RequestContext) {
	userID, hasUserID := auth_context.GetCurrentUserProfileID(c)
	orgID, hasOrg := auth_context.GetCurrentOrganizationID(c)

	// 使用结构化日志记录请求开始
	event := tracelog.Event(ctx, ehm.logger.Info()).
		Str("component", "error_middleware").
		Str("method", string(c.Method())).
		Str("path", string(c.Request.Path())).
		Str("query", string(c.Request.QueryString())).
		Str("remote_addr", c.RemoteAddr().String())

	if hasUserID && userID != "" {
		event = event.Str("user_id", userID)
	}

	if hasOrg && orgID != "" {
		event = event.Str("org_id", orgID)
	}

	event.Msg("Request started")
}

// logResponseInfo 记录响应信息
func (ehm *ErrorHandlerMiddlewareImpl) logResponseInfo(
	ctx context.Context,
	c *app.RequestContext,
	startTime time.Time,
) {
	duration := time.Since(startTime)
	statusCode := c.Response.StatusCode()

	// 选择日志级别：慢请求使用 Warn
	var event *zerolog.Event
	if duration > 100*time.Millisecond {
		event = ehm.logger.Warn()
	} else {
		event = ehm.logger.Info()
	}

	// 使用结构化日志记录请求完成
	event = tracelog.Event(ctx, event).
		Str("component", "error_middleware").
		Str("method", string(c.Method())).
		Str("path", string(c.Request.Path())).
		Int("status_code", statusCode).
		Dur("duration_ms", duration).
		Int("response_size", len(c.Response.Body()))

	if userID, ok := auth_context.GetCurrentUserProfileID(c); ok && userID != "" {
		event = event.Str("user_id", userID)
	}

	event.Msg("Request completed")
}
