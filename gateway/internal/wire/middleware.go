// Package wire 中间件层依赖注入提供者
package wire

import (
	"github.com/google/wire"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	"github.com/rs/zerolog"

	casbnmw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/casbin_middleware"
	corsmdw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/cors_middleware"
	errormw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/error_middleware"
	jwtmdw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/jwt_middleware"
	respmw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/response_middleware"
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
	ProvideCasbinConfig,
	ProvideCasbinMiddleware,
	NewMiddlewareContainer,
)

// MiddlewareContainer 中间件容器
// 统一管理所有中间件实例
type MiddlewareContainer struct {
	TraceMiddleware          tracemdw.TraceMiddlewareService
	CORSMiddleware           corsmdw.CORSMiddlewareService
	ErrorHandlerMiddleware   errormw.ErrorHandlerMiddlewareService
	JWTMiddleware            jwtmdw.JWTMiddlewareService
	ResponseHeaderMiddleware respmw.ResponseHeaderMiddlewareService
	CasbinMiddleware         *casbnmw.CasbinMiddleware
}

// NewMiddlewareContainer 创建中间件容器
func NewMiddlewareContainer(
	traceMiddleware tracemdw.TraceMiddlewareService,
	corsMiddleware corsmdw.CORSMiddlewareService,
	errorHandlerMiddleware errormw.ErrorHandlerMiddlewareService,
	jwtMiddleware jwtmdw.JWTMiddlewareService,
	responseHeaderMiddleware respmw.ResponseHeaderMiddlewareService,
	casbinMiddleware *casbnmw.CasbinMiddleware,
) *MiddlewareContainer {
	return &MiddlewareContainer{
		TraceMiddleware:          traceMiddleware,
		CORSMiddleware:           corsMiddleware,
		ErrorHandlerMiddleware:   errorHandlerMiddleware,
		JWTMiddleware:            jwtMiddleware,
		ResponseHeaderMiddleware: responseHeaderMiddleware,
		CasbinMiddleware:         casbinMiddleware,
	}
}

// ProvideTraceMiddleware 提供追踪中间件
// 自动生成和传播请求追踪信息
func ProvideTraceMiddleware(logger *hertzZerolog.Logger) tracemdw.TraceMiddlewareService {
	middleware := tracemdw.NewTraceMiddleware()

	zl := logger.Unwrap()
	zl.Info().Msg("Trace middleware created successfully")

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
		zl := logger.Unwrap()
		zl.Error().Err(err).Msg("Failed to create JWT middleware")
		panic(err)
	}

	zl := logger.Unwrap()
	zl.Info().Msg("JWT middleware created successfully")

	return middleware
}

// ProvideCORSMiddleware 提供跨域中间件
// 处理跨域资源共享(CORS)配置
func ProvideCORSMiddleware(
	cfg *config.Configuration,
	logger *hertzZerolog.Logger,
) corsmdw.CORSMiddlewareService {
	middleware := corsmdw.NewCORSMiddleware(&cfg.Middleware.CORS, logger)

	zl := logger.Unwrap()
	zl.Info().Msg("CORS middleware created successfully")

	return middleware
}

// ProvideErrorHandlerMiddleware 提供错误处理中间件
// 统一处理请求中的错误响应
func ProvideErrorHandlerMiddleware(
	cfg *config.Configuration,
	logger *hertzZerolog.Logger,
) errormw.ErrorHandlerMiddlewareService {
	middleware := errormw.NewErrorHandlerMiddleware(&cfg.Middleware.ErrorHandler, logger)

	zl := logger.Unwrap()
	zl.Info().Msg("Error Handler middleware created successfully")

	return middleware
}

// ProvideResponseHeaderMiddleware 提供响应头中间件
// 自动为所有响应添加标准 HTTP Date 响应头
func ProvideResponseHeaderMiddleware() respmw.ResponseHeaderMiddlewareService {
	return respmw.NewResponseHeaderMiddleware()
}

// ProvideCasbinConfig 提供 Casbin 配置
func ProvideCasbinConfig() *casbnmw.Config {
	return casbnmw.LoadConfigFromEnv()
}

// ProvideCasbinMiddleware 提供 Casbin 权限中间件
// 使用内存 Adapter，策略通过 RPC 从 Identity Service 同步
func ProvideCasbinMiddleware(
	casbinConfig *casbnmw.Config,
	logger *hertzZerolog.Logger,
) *casbnmw.CasbinMiddleware {
	// 创建一个标准输出的 zerolog.Logger
	zlogger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// 使用新的 ProvideCasbinMiddleware，不再需要数据库连接
	middleware := casbnmw.ProvideCasbinMiddleware(casbinConfig, &zlogger)

	zl := logger.Unwrap()
	zl.Info().Msg("Casbin middleware created successfully (memory adapter)")

	return middleware
}
