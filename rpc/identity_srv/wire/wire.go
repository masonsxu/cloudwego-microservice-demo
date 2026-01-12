//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/config"
	"github.com/rs/zerolog"
)

// =============================================================================
// Provider Sets - 分层组织依赖注入
// =============================================================================

// InfrastructureSet 基础设施层 Provider 集合
// 包含配置、数据库、日志等基础组件
var InfrastructureSet = wire.NewSet(
	config.LoadConfig,
	ProvideDB,
	ProvideLogger,
)

// ConverterSet 转换器 Provider 集合
var ConverterSet = wire.NewSet(
	converter.NewConverter,
)

// DALSet 数据访问层 Provider 集合
var DALSet = wire.NewSet(
	dal.NewDALImpl,
)

// LogicSet 业务逻辑层 Provider 集合
var LogicSet = wire.NewSet(
	logic.NewLogicImpl,
)

// CasbinSet Casbin 权限管理 Provider 集合
var CasbinSet = wire.NewSet(
	ProvideCasbinRepository,
	ProvideCasbinService,
)

// ApplicationSet 完整应用 Provider 集合
// 包含业务逻辑相关的所有依赖
var ApplicationSet = wire.NewSet(
	InfrastructureSet,
	ConverterSet,
	DALSet,
	CasbinSet,
	LogicSet,
)

// AllSet 所有依赖注入集合
// 按照分层架构组织：基础设施层 -> 业务层 -> 服务器层 -> 健康检查层
var AllSet = wire.NewSet(
	ApplicationSet,
	ServerSet,
	HealthCheckSet,
	NewAppContainer,
)

// =============================================================================
// Injector Functions - 依赖注入函数
// =============================================================================

// InitializeApp 初始化应用容器
// 统一初始化所有依赖，避免重复创建实例
// 返回 cleanup 函数用于资源清理（如 OpenTelemetry shutdown）
func InitializeApp() (*AppContainer, func(), error) {
	wire.Build(AllSet)
	return nil, nil, nil
}

// InitializeLogger 仅初始化日志器
func InitializeLogger() (*zerolog.Logger, error) {
	wire.Build(InfrastructureSet)
	return nil, nil
}
