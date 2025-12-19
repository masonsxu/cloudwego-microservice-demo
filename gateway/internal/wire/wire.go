//go:build wireinject

// Package wire Wire依赖注入配置
// 使用Google Wire进行依赖注入管理，实现清晰的分层架构
package wire

import (
	"github.com/google/wire"
)

// AllSet 所有依赖注入集合
// 按照分层架构组织：基础设施层 -> 应用层 -> 领域层 -> 中间件层
var AllSet = wire.NewSet(
	// 基础设施层
	InfrastructureSet,

	// 应用层
	ApplicationSet,

	// 领域服务层
	DomainServiceSet,

	// 中间件层
	MiddlewareSet,

	// 容器
	NewServiceContainer,
	NewAppContainer,
)

// InitializeApp 初始化应用容器
// 统一初始化所有依赖，避免重复创建 Logger 等实例
func InitializeApp() (*AppContainer, error) {
	wire.Build(AllSet)
	return &AppContainer{}, nil
}
