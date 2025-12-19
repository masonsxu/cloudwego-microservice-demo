package middleware

import (
	casbinmw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/casbin_middleware"
)

// 全局中间件实例，由 main.go 初始化后设置
var (
	globalCasbinMiddleware casbinmw.CasbinMiddleware
)

// SetGlobalCasbinMiddleware 设置全局 Casbin 中间件实例
// 在 main.go 中初始化中间件后调用
func SetGlobalCasbinMiddleware(casbin casbinmw.CasbinMiddleware) {
	globalCasbinMiddleware = casbin
}

// GetGlobalCasbinMiddleware 获取全局 Casbin 中间件实例
// 用于在 biz/router/**/middleware.go 中使用
func GetGlobalCasbinMiddleware() casbinmw.CasbinMiddleware {
	return globalCasbinMiddleware
}
