// Package wire 中间件层依赖注入提供者
package wire

import (
	"github.com/google/wire"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"

	accesslogmw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/access_log_middleware"
	authzmw "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/authz_middleware"
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
	ProvideAuthZRules,
	ProvideAuthZMiddleware,
	ProvideAccessLogMiddleware,
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
	AuthZMiddleware          authzmw.AuthZMiddlewareService
	AccessLogMiddleware      accesslogmw.AccessLogMiddlewareService
}

// NewMiddlewareContainer 创建中间件容器
func NewMiddlewareContainer(
	traceMiddleware tracemdw.TraceMiddlewareService,
	corsMiddleware corsmdw.CORSMiddlewareService,
	errorHandlerMiddleware errormw.ErrorHandlerMiddlewareService,
	jwtMiddleware jwtmdw.JWTMiddlewareService,
	responseHeaderMiddleware respmw.ResponseHeaderMiddlewareService,
	authzMiddleware authzmw.AuthZMiddlewareService,
	accessLogMiddleware accesslogmw.AccessLogMiddlewareService,
) *MiddlewareContainer {
	return &MiddlewareContainer{
		TraceMiddleware:          traceMiddleware,
		CORSMiddleware:           corsMiddleware,
		ErrorHandlerMiddleware:   errorHandlerMiddleware,
		JWTMiddleware:            jwtMiddleware,
		ResponseHeaderMiddleware: responseHeaderMiddleware,
		AuthZMiddleware:          authzMiddleware,
		AccessLogMiddleware:      accessLogMiddleware,
	}
}

// ProvideTraceMiddleware 提供追踪中间件
// 自动生成和传播请求追踪信息，并将带追踪上下文的 logger 绑定到每个请求的 context
func ProvideTraceMiddleware(logger *hertzZerolog.Logger) tracemdw.TraceMiddlewareService {
	zl := logger.Unwrap()
	return tracemdw.NewTraceMiddleware(&zl)
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

// ProvideAuthZRules 加载路由级 ACL 规则
//
// 失败时（文件不存在 / YAML 语法错误）panic，让进程在启动期就暴露问题，
// 避免运行时静默放行所有请求。
func ProvideAuthZRules(
	cfg *config.Configuration,
	logger *hertzZerolog.Logger,
) *authzmw.Rules {
	zl := logger.Unwrap()

	rulesFile := cfg.Middleware.AuthZ.RulesFile
	if rulesFile == "" {
		zl.Error().Msg("AUTHZ_RULES_FILE 未配置")
		panic("authz rules_file is empty")
	}

	rules, err := authzmw.LoadRulesFromFile(rulesFile)
	if err != nil {
		zl.Error().Err(err).Str("rules_file", rulesFile).Msg("Failed to load authz rules")
		panic(err)
	}

	zl.Info().
		Str("rules_file", rulesFile).
		Int("public", len(rules.Public)).
		Int("authenticated", len(rules.Authenticated)).
		Int("roles", len(rules.Roles)).
		Str("default", string(rules.Default)).
		Msg("AuthZ rules loaded")

	return rules
}

// ProvideAuthZMiddleware 提供路由级 ACL 中间件
func ProvideAuthZMiddleware(
	rules *authzmw.Rules,
	logger *hertzZerolog.Logger,
) authzmw.AuthZMiddlewareService {
	zl := logger.Unwrap()
	mw := authzmw.NewAuthZMiddleware(rules, &zl)

	zl.Info().Msg("AuthZ middleware created successfully")

	return mw
}

// ProvideAccessLogMiddleware 提供访问日志中间件
func ProvideAccessLogMiddleware(
	logger *hertzZerolog.Logger,
) accesslogmw.AccessLogMiddlewareService {
	zl := logger.Unwrap()
	mw := accesslogmw.NewAccessLogMiddleware(&zl)

	zl.Info().Msg("Access log middleware created successfully")

	return mw
}
