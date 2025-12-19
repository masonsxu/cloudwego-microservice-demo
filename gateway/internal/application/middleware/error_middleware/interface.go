package middleware

import "github.com/cloudwego/hertz/pkg/app"

// ErrorHandlerMiddlewareService 定义错误处理中间件服务接口
type ErrorHandlerMiddlewareService interface {
	MiddlewareFunc() app.HandlerFunc
}
