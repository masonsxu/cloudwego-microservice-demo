package definition

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// RoleDefinitionQueryConditions 角色定义查询条件
// 对应 IDL 中的 RoleDefinitionQueryRequest，用于仓储层查询
type RoleDefinitionQueryConditions struct {
	// Name 角色名称，支持精确匹配
	Name *string `json:"name,omitempty"`

	// Status 角色状态过滤
	Status *models.RoleStatus `json:"status,omitempty"`

	// IsSystemRole 是否为系统内置角色
	IsSystemRole *bool `json:"is_system_role,omitempty"`

	// Page 分页查询选项
	Page *base.QueryOptions `json:"page,omitempty"`
}

// RoleDefinitionRepository 角色定义仓储接口
// 基于 models.RoleDefinition，管理系统中的角色定义和权限配置
type RoleDefinitionRepository interface {
	// 嵌入基础仓储接口
	base.BaseRepository[models.RoleDefinition]

	// ============================================================================
	// 核心查询方法
	// ============================================================================

	// FindByName 根据角色名称查询角色定义（精确匹配）
	// 由于角色名称具有唯一索引，此方法返回单个结果或 nil
	FindByName(ctx context.Context, name string) (*models.RoleDefinition, error)

	// CheckNameExists 检查指定角色名称是否已存在
	// 用于创建角色前的唯一性验证，避免数据库约束冲突
	CheckNameExists(ctx context.Context, name string) (bool, error)

	// FindByStatus 根据角色状态查询角色定义列表
	// 支持分页，常用于管理界面的状态筛选
	FindByStatus(
		ctx context.Context,
		status models.RoleStatus,
		page *base.QueryOptions,
	) ([]*models.RoleDefinition, *models.PageResult, error)

	// FindBySystemRole 根据系统角色标识查询角色定义列表
	// isSystemRole: true=系统内置角色, false=用户自定义角色
	FindBySystemRole(
		ctx context.Context,
		isSystemRole bool,
		page *base.QueryOptions,
	) ([]*models.RoleDefinition, *models.PageResult, error)

	// FindWithConditions 根据组合查询条件查询角色定义列表
	// 这是最灵活的查询方法，对应 IDL 中的 RoleDefinitionQueryRequest
	FindWithConditions(
		ctx context.Context,
		conditions *RoleDefinitionQueryConditions,
	) ([]*models.RoleDefinition, *models.PageResult, error)

	// ============================================================================
	// 业务统计方法
	// ============================================================================

	// CountByStatus 统计指定状态的角色定义数量
	// 用于仪表板统计和数据分析
	CountByStatus(ctx context.Context, status models.RoleStatus) (int64, error)

	// ListActiveRoles 列出所有活跃状态的角色定义
	// 常用于角色分配场景，只显示可用角色
	ListActiveRoles(
		ctx context.Context,
		page *base.QueryOptions,
	) ([]*models.RoleDefinition, *models.PageResult, error)
}
