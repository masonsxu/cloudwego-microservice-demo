// Package middleware 提供请求追踪中间件
// 负责生成/提取请求追踪信息并传播到后续调用链
package middleware

import (
	"github.com/cloudwego/hertz/pkg/app"
)

// TraceMiddlewareService 追踪中间件服务接口
type TraceMiddlewareService interface {
	// MiddlewareFunc 返回追踪中间件函数
	MiddlewareFunc() app.HandlerFunc
}
