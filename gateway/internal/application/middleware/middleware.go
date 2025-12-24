// Package middleware 提供中间件相关功能
//
// Deprecated: 此文件已废弃，中间件注册已迁移至 wire/server.go 中的 HandlerRegistry。
// 保留此文件仅供参考，将在未来版本中删除。
package middleware

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/etag"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/requestid"

	corsmw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/cors_middleware"
	errormw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/error_middleware"
	jwtmw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/jwt_middleware"
	responsemw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/response_middleware"
	tracemdw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/trace_middleware"
)

// DefaultMiddleware 默认中间件注册
//
// Deprecated: 请使用 wire.HandlerRegistry.RegisterMiddlewares() 代替。
// 此函数将在未来版本中删除。
func DefaultMiddleware(
	h *server.Hertz,
	tracerCfg *hertztracing.Config,
	traceMiddleware tracemdw.TraceMiddlewareService,
	corsMiddleware corsmw.CORSMiddlewareService,
	errorMiddleware errormw.ErrorHandlerMiddlewareService,
	jwtMiddleware jwtmw.JWTMiddlewareService,
	responseHeaderMiddleware responsemw.ResponseHeaderMiddlewareService,
) {
	h.Use(
		hertztracing.ServerMiddleware(tracerCfg),  // 追踪：最先执行，生成/提取追踪信息
		requestid.New(),                           // RequestID：生成和传递请求ID
		responseHeaderMiddleware.MiddlewareFunc(), // 响应头：添加标准 HTTP Date 头部
		traceMiddleware.MiddlewareFunc(),          // 追踪：最先执行，生成/提取追踪信息
		corsMiddleware.MiddlewareFunc(),           // 跨域：处理预检，避免被后续中间件拦截
		errorMiddleware.MiddlewareFunc(),          // 错误处理：后续所有错误均由其捕获
		jwtMiddleware.MiddlewareFunc(),            // 认证：解析用户身份，存入上下文
		etag.New(),                                // ETag：计算和验证 ETag
	)
}
