package middleware

import "github.com/cloudwego/hertz/pkg/app"

// CORSMiddlewareService 定义CORS中间件服务接口
type CORSMiddlewareService interface {
	MiddlewareFunc() app.HandlerFunc
}
