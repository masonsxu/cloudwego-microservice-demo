package middleware

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/context/auth_context"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/errors"
)

// ErrorHandlerMiddlewareImpl 错误处理中间件实现
type ErrorHandlerMiddlewareImpl struct {
	config *config.ErrorHandlerConfig
	logger *hertzZerolog.Logger
}

// NewErrorHandlerMiddleware 创建新的错误处理中间件实例
func NewErrorHandlerMiddleware(
	config *config.ErrorHandlerConfig,
	logger *hertzZerolog.Logger,
) ErrorHandlerMiddlewareService {
	if logger == nil {
		// 如果 logger 为 nil，创建一个默认的 logger
		logger = hertzZerolog.New()
	}

	return &ErrorHandlerMiddlewareImpl{
		config: config,
		logger: logger,
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
			ehm.logger.Warnf("Casbin permission denied: method=%s, path=%s",
				string(c.Method()), string(c.Request.Path()))
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

	// 记录详细的panic日志
	ehm.logger.Errorf(
		"Panic recovered in error handler: panic=%v, path=%s, method=%s, user_agent=%s, remote_addr=%s, stack_trace=%s",
		r,
		string(c.Request.Path()),
		string(c.Method()),
		string(c.UserAgent()),
		c.RemoteAddr(),
		stackTrace,
	)

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
	// 记录用户信息（如果有的话）
	logMsg := fmt.Sprintf("HTTP status error: status_code=%d, path=%s, method=%s",
		statusCode, string(c.Request.Path()), string(c.Method()))

	// 如果有用户上下文，记录用户信息
	if userID, ok := auth_context.GetCurrentUserProfileID(c); ok && userID != "" {
		logMsg += fmt.Sprintf(", user_id=%s", userID)
	}

	if orgID, ok := auth_context.GetCurrentOrganizationID(c); ok && orgID != "" {
		logMsg += fmt.Sprintf(", org_id=%s", orgID)
	}

	ehm.logger.Warnf("%s", logMsg)

	// 如果启用了详细错误信息，可以在这里添加更多上下文
	if ehm.config.EnableDetailedErrors {
		// 可以在这里添加更详细的错误信息处理
		ehm.logger.Debugf("Detailed error context: remote_addr=%s, user_agent=%s",
			c.RemoteAddr(), string(c.UserAgent()))
	}

	errors.HandleHTTPStatusError(c, statusCode)
}

// logRequestInfo 记录请求信息
func (ehm *ErrorHandlerMiddlewareImpl) logRequestInfo(ctx context.Context, c *app.RequestContext) {
	UserID, hasUserID := auth_context.GetCurrentUserProfileID(c)
	orgID, hasOrg := auth_context.GetCurrentOrganizationID(c)

	logMsg := fmt.Sprintf(
		"Request started: method=%s, path=%s, query=%s, remote_addr=%s",
		string(
			c.Method(),
		),
		string(c.Request.Path()),
		string(c.Request.QueryString()),
		c.RemoteAddr(),
	)

	if hasUserID && UserID != "" {
		logMsg += fmt.Sprintf(", user_id=%s", UserID)
	}

	if hasOrg && orgID != "" {
		logMsg += fmt.Sprintf(", org_id=%s", orgID)
	}

	ehm.logger.Infof("%s", logMsg)
}

// logResponseInfo 记录响应信息
func (ehm *ErrorHandlerMiddlewareImpl) logResponseInfo(
	ctx context.Context,
	c *app.RequestContext,
	startTime time.Time,
) {
	duration := time.Since(startTime)
	statusCode := c.Response.StatusCode()

	logMsg := fmt.Sprintf(
		"Request completed: method=%s, path=%s, status_code=%d, duration_ms=%d, response_size=%d",
		string(
			c.Method(),
		),
		string(c.Request.Path()),
		statusCode,
		duration.Milliseconds(),
		len(c.Response.Body()),
	)

	if userID, ok := auth_context.GetCurrentUserProfileID(c); ok && userID != "" {
		logMsg += fmt.Sprintf(", user_id=%s", userID)
	}

	ehm.logger.Infof("%s", logMsg)
}
