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

		if c.Response.StatusCode() == 500 && len(c.Response.Body()) == 0 {
			tracelog.Event(ctx, ehm.logger.Warn()).
				Str("component", "error_middleware").
				Str("method", string(c.Method())).
				Str("path", string(c.Request.Path())).
				Msg("Permission denied")
			c.SetStatusCode(403)
			c.JSON(403, map[string]string{"msg": "权限不足"})

			return
		}

		// 处理完成后检查是否有错误
		if c.IsAborted() {
			// 请求已被中断，可能是认证失败或其他错误
			return
		}

		// 检查响应状态码，处理HTTP状态错误
		if statusCode := c.Response.StatusCode(); statusCode >= 400 {
			ehm.handleHTTPStatusError(ctx, c, statusCode)
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

// handleHTTPStatusError 处理HTTP状态码错误
func (ehm *ErrorHandlerMiddlewareImpl) handleHTTPStatusError(
	ctx context.Context,
	c *app.RequestContext,
	statusCode int,
) {
	// 使用结构化日志记录 HTTP 状态错误
	event := tracelog.Event(ctx, ehm.logger.Warn()).
		Str("component", "error_middleware").
		Int("status_code", statusCode).
		Str("path", string(c.Request.Path())).
		Str("method", string(c.Method()))

	// 如果有用户上下文，记录用户信息
	if userID, ok := auth_context.GetCurrentUserProfileID(c); ok && userID != "" {
		event = event.Str("user_id", userID)
	}

	if orgID, ok := auth_context.GetCurrentOrganizationID(c); ok && orgID != "" {
		event = event.Str("org_id", orgID)
	}

	// 如果启用了详细错误信息，添加更多上下文
	if ehm.config.EnableDetailedErrors {
		event = event.
			Str("remote_addr", c.RemoteAddr().String()).
			Str("user_agent", string(c.UserAgent()))
	}

	event.Msg("HTTP status error")

	errors.HandleHTTPStatusError(c, statusCode)
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
