//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/casbin"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/logic"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/config"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
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

// CasbinSet Casbin 权限管理 Provider 集合
var CasbinSet = wire.NewSet(
	ProvideCasbinConfig,
	casbin.NewCasbinManager,
	casbin.NewMenuPermissionLogic,
)

// DALSet 数据访问层 Provider 集合
var DALSet = wire.NewSet(
	dal.NewDALImpl,
)

// LogicSet 业务逻辑层 Provider 集合
var LogicSet = wire.NewSet(
	logic.NewLogicImpl,
)

// ApplicationSet 完整应用 Provider 集合
var ApplicationSet = wire.NewSet(
	InfrastructureSet,
	ConverterSet,
	CasbinSet,
	DALSet,
	LogicSet,
)

// TestSet 测试环境专用 Provider 集合
var TestSet = wire.NewSet(
	InfrastructureSet,
	ConverterSet,
	CasbinSet,
	DALSet,
	LogicSet,
)

// =============================================================================
// Injector Functions - 依赖注入函数
// =============================================================================

// InitializeService 初始化服务，返回业务逻辑层实例
func InitializeService() (logic.Logic, error) {
	wire.Build(ApplicationSet)
	return nil, nil
}

// InitializeLogger 仅初始化日志器
func InitializeLogger() (*zerolog.Logger, error) {
	wire.Build(InfrastructureSet)
	return nil, nil
}

// InitializeTestService 初始化测试环境服务
func InitializeTestService() (logic.Logic, error) {
	wire.Build(TestSet)
	return nil, nil
}

// ServiceWithDB 包含服务逻辑和数据库连接的包装结构
// 用于需要同时访问服务和数据库的场景（如健康检查）
type ServiceWithDB struct {
	Service logic.Logic
	DB      *gorm.DB
}

// ProvideServiceWithDB 提供包含服务和数据库的包装结构
func ProvideServiceWithDB(service logic.Logic, db *gorm.DB) *ServiceWithDB {
	return &ServiceWithDB{
		Service: service,
		DB:      db,
	}
}

// InitializeServiceWithDB 初始化服务并返回 DB 连接（用于健康检查）
func InitializeServiceWithDB() (*ServiceWithDB, error) {
	wire.Build(
		ApplicationSet,
		ProvideServiceWithDB,
	)
	return nil, nil
}
