// Package wire 应用容器定义
package wire

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/config"
)

// AppContainer 应用容器
// 统一管理所有依赖，由 Wire 自动构建
type AppContainer struct {
	Config            *config.Config
	Logger            *zerolog.Logger
	DB                *gorm.DB
	Logic             logic.Logic
	ServerOptions     *ServerOptions
	HealthCheckServer *HealthCheckServer
}

// NewAppContainer 创建应用容器
func NewAppContainer(
	cfg *config.Config,
	logger *zerolog.Logger,
	db *gorm.DB,
	logicImpl logic.Logic,
	serverOpts *ServerOptions,
	healthServer *HealthCheckServer,
) *AppContainer {
	klog.Infof("Application container initialized successfully")

	return &AppContainer{
		Config:            cfg,
		Logger:            logger,
		DB:                db,
		Logic:             logicImpl,
		ServerOptions:     serverOpts,
		HealthCheckServer: healthServer,
	}
}

// StartHealthCheck 启动健康检查服务器
func (c *AppContainer) StartHealthCheck() {
	c.HealthCheckServer.Start()
}
