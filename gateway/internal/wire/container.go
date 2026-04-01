// Package wire 服务容器定义
package wire

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"

	identityService "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/service/identity"
	permissionService "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/service/permission"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
)

// AppContainer 应用容器
// 统一管理所有依赖实例，避免重复初始化
type AppContainer struct {
	Config          *config.Configuration
	Logger          *hertzZerolog.Logger
	Services        *ServiceContainer
	Middlewares     *MiddlewareContainer
	HandlerRegistry *HandlerRegistry
}

// NewAppContainer 创建应用容器
func NewAppContainer(
	cfg *config.Configuration,
	logger *hertzZerolog.Logger,
	services *ServiceContainer,
	middlewares *MiddlewareContainer,
	handlerRegistry *HandlerRegistry,
) *AppContainer {
	// 设置 hlog 使用统一的 logger
	hlog.SetLogger(logger)

	return &AppContainer{
		Config:          cfg,
		Logger:          logger,
		Services:        services,
		Middlewares:     middlewares,
		HandlerRegistry: handlerRegistry,
	}
}

// ServiceContainer 服务容器
// 统一管理所有业务服务实例
type ServiceContainer struct {
	IdentityService   identityService.Service
	PermissionService permissionService.Service
}

// NewServiceContainer 创建服务容器
func NewServiceContainer(
	identityService identityService.Service,
	permissionService permissionService.Service,
) *ServiceContainer {
	return &ServiceContainer{
		IdentityService:   identityService,
		PermissionService: permissionService,
	}
}
