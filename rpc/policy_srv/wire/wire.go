//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/biz/logic"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/config"
)

// InfrastructureSet 基础设施
var InfrastructureSet = wire.NewSet(
	config.LoadConfig,
	ProvideLogger,
	ProvideDB,
	ProvideSQLDB,
)

// ApplicationSet 应用层
var ApplicationSet = wire.NewSet(
	ProvideEnforcerService,
	ProvideDecisionService,
)

// ServerSet 服务器配置
var ServerSet = wire.NewSet(
	ProvideEtcdRegistry,
	ProvideServerOptions,
)

// AllSet 全部
var AllSet = wire.NewSet(
	InfrastructureSet,
	ApplicationSet,
	ServerSet,
	NewAppContainer,
)

// InitializeApp Wire 注入入口
func InitializeApp() (*AppContainer, func(), error) {
	wire.Build(
		AllSet,
	)
	return nil, nil, nil
}

// 确保类型被导入
var (
	_ *zerolog.Logger
	_ *gorm.DB
	_ *logic.EnforcerService
	_ *logic.DecisionService
)
