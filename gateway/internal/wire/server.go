// Package wire 服务器层依赖注入提供者
package wire

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/google/wire"
	"github.com/hertz-contrib/etag"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/requestid"

	identityHandler "github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/handler/identity"
	permissionHandler "github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/handler/permission"
	identityService "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/service/identity"
	permissionService "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/service/permission"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/otel"
)

// ServerSet 服务器层依赖注入集合
var ServerSet = wire.NewSet(
	ProvideOtelProvider,
	ProvideTracer,
	ProvideServerFactory,
	ProvideHandlerRegistry,
)

// ProvideOtelProvider 提供 OpenTelemetry Provider
func ProvideOtelProvider(cfg *config.Configuration) (*otel.Provider, func(), error) {
	return otel.NewProvider(cfg)
}

// ProvideTracer 提供 Hertz Server Tracer
func ProvideTracer(cfg *config.Configuration) *otel.Tracer {
	return otel.NewTracer(cfg)
}

// ProvideServerFactory 提供 Hertz Server 工厂
func ProvideServerFactory(cfg *config.Configuration, tracer *otel.Tracer, provider *otel.Provider) *otel.ServerFactory {
	return otel.NewServerFactory(cfg, tracer, provider)
}

// HandlerRegistry Handler 注册器
// 负责将依赖注入到 Handler 层并注册中间件
type HandlerRegistry struct {
	server            *server.Hertz
	tracer            *otel.Tracer
	middlewares       *MiddlewareContainer
	identityService   identityService.Service
	permissionService permissionService.Service
	logger            *hertzZerolog.Logger
}

// ProvideHandlerRegistry 提供 Handler 注册器
func ProvideHandlerRegistry(
	factory *otel.ServerFactory,
	tracer *otel.Tracer,
	middlewares *MiddlewareContainer,
	services *ServiceContainer,
	logger *hertzZerolog.Logger,
) *HandlerRegistry {
	return &HandlerRegistry{
		server:            factory.Server(),
		tracer:            tracer,
		middlewares:       middlewares,
		identityService:   services.IdentityService,
		permissionService: services.PermissionService,
		logger:            logger,
	}
}

// RegisterMiddlewares 注册全局中间件
func (r *HandlerRegistry) RegisterMiddlewares() {
	r.server.Use(
		hertztracing.ServerMiddleware(r.tracer.Config), // 追踪：最先执行，生成/提取追踪信息
		requestid.New(), // RequestID：生成和传递请求ID
		r.middlewares.ResponseHeaderMiddleware.MiddlewareFunc(), // 响应头：添加标准 HTTP Date 头部
		r.middlewares.TraceMiddleware.MiddlewareFunc(),          // 追踪：最先执行，生成/提取追踪信息
		r.middlewares.CORSMiddleware.MiddlewareFunc(),           // 跨域：处理预检，避免被后续中间件拦截
		r.middlewares.ErrorHandlerMiddleware.MiddlewareFunc(),   // 错误处理：后续所有错误均由其捕获
		r.middlewares.JWTMiddleware.MiddlewareFunc(),            // 认证：解析用户身份，存入上下文
		r.middlewares.CasbinMiddleware.MiddlewareFunc(),         // 权限：基于 Casbin RBAC 进行权限校验
		etag.New(), // ETag：计算和验证 ETag
	)

	r.logger.Infof("Global middlewares registered successfully")
}

// RegisterHandlers 注册 Handler 依赖
func (r *HandlerRegistry) RegisterHandlers() {
	// 设置 Handler 层的服务依赖
	identityHandler.SetIdentityService(r.identityService, r.middlewares.JWTMiddleware)
	permissionHandler.SetPermissionService(r.permissionService)

	r.logger.Infof("Handler dependencies registered successfully")
}

// Server 返回 Hertz 服务器实例
func (r *HandlerRegistry) Server() *server.Hertz {
	return r.server
}

// Initialize 初始化所有注册（中间件 + Handler）
func (r *HandlerRegistry) Initialize() {
	r.RegisterMiddlewares()
	r.RegisterHandlers()
}
