package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

// ResponseHeaderMiddleware 响应头中间件实现
type ResponseHeaderMiddleware struct{}

// NewResponseHeaderMiddleware 创建响应头中间件实例
func NewResponseHeaderMiddleware() ResponseHeaderMiddlewareService {
	return &ResponseHeaderMiddleware{}
}

// MiddlewareFunc 返回中间件处理函数
// 自动为所有响应添加标准 HTTP Date 响应头（符合 RFC 7231 规范）
func (m *ResponseHeaderMiddleware) MiddlewareFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 设置 Date 响应头（RFC 7231 标准格式）
		// 注意：使用 UTC 时间和 RFC1123 格式
		c.Header("Date", time.Now().UTC().Format(time.RFC1123))

		// 继续处理请求
		c.Next(ctx)
	}
}
