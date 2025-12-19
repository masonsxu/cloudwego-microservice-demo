package middleware

import "github.com/cloudwego/hertz/pkg/app"

// ResponseHeaderMiddlewareService 响应头中间件接口
// 负责为所有响应自动添加标准 HTTP 响应头
type ResponseHeaderMiddlewareService interface {
	// MiddlewareFunc 返回中间件函数
	MiddlewareFunc() app.HandlerFunc
}
