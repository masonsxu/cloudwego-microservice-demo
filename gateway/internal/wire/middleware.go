// Package wire 中间件层依赖注入提供者
package wire

import (
	"github.com/google/wire"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"

	corsmdw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/cors_middleware"
	errormw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/error_middleware"
	jwtmdw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/jwt_middleware"
	responsemw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/response_middleware"
	tracemdw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/trace_middleware"
	identityService "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/service/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/redis"
)

// MiddlewareSet 中间件层依赖注入集合
var MiddlewareSet = wire.NewSet(
	ProvideTraceMiddleware,
	ProvideCORSMiddleware,
	ProvideErrorHandlerMiddleware,
	ProvideJWTMiddleware,
	ProvideResponseHeaderMiddleware,
	NewMiddlewareContainer,
)

// MiddlewareContainer 中间件容器
// 统一管理所有中间件实例
type MiddlewareContainer struct {
	TraceMiddleware          tracemdw.TraceMiddlewareService
	CORSMiddleware           corsmdw.CORSMiddlewareService
	ErrorHandlerMiddleware   errormw.ErrorHandlerMiddlewareService
	JWTMiddleware            jwtmdw.JWTMiddlewareService
	ResponseHeaderMiddleware responsemw.ResponseHeaderMiddlewareService
}

// NewMiddlewareContainer 创建中间件容器
func NewMiddlewareContainer(
	traceMiddleware tracemdw.TraceMiddlewareService,
	corsMiddleware corsmdw.CORSMiddlewareService,
	errorHandlerMiddleware errormw.ErrorHandlerMiddlewareService,
	jwtMiddleware jwtmdw.JWTMiddlewareService,
	responseHeaderMiddleware responsemw.ResponseHeaderMiddlewareService,
) *MiddlewareContainer {
	return &MiddlewareContainer{
		TraceMiddleware:          traceMiddleware,
		CORSMiddleware:           corsMiddleware,
		ErrorHandlerMiddleware:   errorHandlerMiddleware,
		JWTMiddleware:            jwtMiddleware,
		ResponseHeaderMiddleware: responseHeaderMiddleware,
	}
}

// ProvideTraceMiddleware 提供追踪中间件
// 自动生成和传播请求追踪信息
func ProvideTraceMiddleware(logger *hertzZerolog.Logger) tracemdw.TraceMiddlewareService {
	middleware := tracemdw.NewTraceMiddleware()

	logger.Infof("Trace middleware created successfully")

	return middleware
}

// ProvideJWTMiddleware 提供JWT中间件
// 配置JWT认证中间件，用于API权限控制
func ProvideJWTMiddleware(
	identityService identityService.Service,
	jwtConfig *config.JWTConfig,
	tokenCache redis.TokenCacheService,
	logger *hertzZerolog.Logger,
) jwtmdw.JWTMiddlewareService {
	middleware, err := jwtmdw.JWTMiddlewareProvider(identityService, jwtConfig, tokenCache, logger)
	if err != nil {
		logger.Errorf("Failed to create JWT middleware: %v", err)
		panic(err)
	}

	logger.Infof("JWT middleware created successfully")

	return middleware
}

// ProvideCORSMiddleware 提供跨域中间件
// 处理跨域资源共享(CORS)配置
func ProvideCORSMiddleware(
	cfg *config.Configuration,
	logger *hertzZerolog.Logger,
) corsmdw.CORSMiddlewareService {
	middleware := corsmdw.NewCORSMiddleware(&cfg.Middleware.CORS, logger)
	logger.Infof("CORS middleware created successfully")

	return middleware
}

// ProvideErrorHandlerMiddleware 提供错误处理中间件
// 统一处理请求中的错误响应
func ProvideErrorHandlerMiddleware(
	cfg *config.Configuration,
	logger *hertzZerolog.Logger,
) errormw.ErrorHandlerMiddlewareService {
	middleware := errormw.NewErrorHandlerMiddleware(&cfg.Middleware.ErrorHandler, logger)
	logger.Infof("Error Handler middleware created successfully")

	return middleware
}

// ProvideResponseHeaderMiddleware 提供响应头中间件
// 自动为所有响应添加标准 HTTP Date 响应头
func ProvideResponseHeaderMiddleware() responsemw.ResponseHeaderMiddlewareService {
	return responsemw.NewResponseHeaderMiddleware()
}
