package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/requestid"
	"github.com/rs/zerolog"

	jwtmw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/jwt_middleware"
	tracelog "github.com/masonsxu/cloudwego-microservice-demo/gateway/pkg/log"
)

// AccessLogMiddlewareImpl 访问日志中间件实现
type AccessLogMiddlewareImpl struct {
	logger *zerolog.Logger
}

// NewAccessLogMiddleware 创建访问日志中间件实例
func NewAccessLogMiddleware(logger *zerolog.Logger) *AccessLogMiddlewareImpl {
	return &AccessLogMiddlewareImpl{logger: logger}
}

// MiddlewareFunc 返回中间件函数
//
// 在请求 handler chain 末尾输出一条 Info 级别日志，包含：
// method、path、status、duration、user_id（如有）、request_id（如有）。
//
// 不替代 RequestID 中间件本身——后者负责生成/传播 ID，本中间件仅消费。
func (m *AccessLogMiddlewareImpl) MiddlewareFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		c.Next(ctx)
		duration := time.Since(start)

		evt := tracelog.Event(ctx, m.logger.Info()).
			Str("component", "access_log").
			Str("method", string(c.Request.Method())).
			Str("path", string(c.Request.URI().Path())).
			Int("status", c.Response.StatusCode()).
			Dur("duration", duration)

		if reqID := requestid.Get(c); reqID != "" {
			evt = evt.Str("request_id", reqID)
		}

		if uid := string(c.Request.Header.Get(jwtmw.HeaderUserID)); uid != "" {
			evt = evt.Str("user_id", uid)
		}

		evt.Msg("access")
	}
}
