// Package middleware 提供请求追踪中间件
// 负责将 requestid 中间件生成的 RequestID 和 OpenTelemetry 追踪信息注入到 RPC 调用链中
package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/requestid"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/errors"
)

// TraceMiddlewareImpl 追踪中间件实现
type TraceMiddlewareImpl struct{}

// NewTraceMiddleware 创建追踪中间件实例
func NewTraceMiddleware() TraceMiddlewareService {
	return &TraceMiddlewareImpl{}
}

// MiddlewareFunc 返回追踪中间件函数
// 此中间件执行以下操作：
// 1. 从 requestid 中间件获取 request_id（使用 requestid.Get 函数）
// 2. 将 RequestID 注入到 Go context (metainfo) 供 RPC 调用传播
// 3. 将 OpenTelemetry TraceID/SpanID 注入到 metainfo 供 RPC 调用传播
//
// 注意：此中间件应在 requestid 中间件和 OTel tracing 中间件之后执行
func (m *TraceMiddlewareImpl) MiddlewareFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从 requestid 中间件获取 RequestID
		// requestid 中间件会自动从 X-Request-ID header 获取或生成 RequestID
		requestID := requestid.Get(c)
		if requestID == "" {
			// 如果 requestid 中间件还没有设置，说明可能有问题，跳过
			// requestid 中间件应该已经处理了这种情况
			c.Next(ctx)
			return
		}

		// 将 RequestID 注入到 Go context (metainfo) 供 RPC 调用传播
		ctx = errors.InjectRequestIDToContext(ctx, requestID)

		// 将 OpenTelemetry TraceID/SpanID 注入到 metainfo 供 RPC 调用传播
		// 这确保日志和链路追踪可以对齐
		ctx = errors.InjectTraceToContext(ctx)

		// 继续处理请求
		c.Next(ctx)
	}
}
