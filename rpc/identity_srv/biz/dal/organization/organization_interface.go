package organization

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// OrganizationRepository 组织仓储接口（简化版）
// 基于 models.Organization，管理机构的简单2级组织架构
type OrganizationRepository interface {
	// 嵌入基础仓储接口
	base.BaseRepository[models.Organization]

	// ============================================================================
	// 核心查询方法（保留）
	// ============================================================================

	// GetByCode 根据组织代码获取组织
	GetByCode(ctx context.Context, code string) (*models.Organization, error)

	// ============================================================================
	// 简化的关系管理
	// ============================================================================

	// UpdateParent 更新组织的父组织（仅支持2级层级）
	UpdateParent(ctx context.Context, organizationID, newParentID string) error

	// HasChildren 检查组织是否有直接子组织
	HasChildren(ctx context.Context, organizationID string) (bool, error)

	// CountChildren 统计直接子组织数量
	CountChildren(ctx context.Context, parentID string) (int64, error)
	// ExistsByID 检查组织是否存在
	ExistsByID(ctx context.Context, organizationID string) (bool, error)

	// ============================================================================
	// 数据完整性检查（简化）
	// ============================================================================

	// CheckNameConflict 检查组织名称是否冲突
	CheckNameConflict(
		ctx context.Context,
		name string,
		parentID string,
		excludeID ...string,
	) (bool, error)

	// ============================================================================
	// 统一查询方法
	// ============================================================================

	// FindWithConditions 根据组合查询条件查询组织列表
	FindWithConditions(
		ctx context.Context,
		conditions *OrganizationQueryConditions,
	) ([]*models.Organization, *models.PageResult, error)
}

// OrganizationQueryConditions 组织查询条件
// 支持多条件组合查询，提供灵活的查询能力
type OrganizationQueryConditions struct {
	Name         *string            // 组织名称（模糊匹配）
	Code         *string            // 组织代码（精确匹配）
	ParentID     *string            // 父组织ID（nil表示查询根组织）
	FacilityType *string            // 机构类型
	ProvinceCity *string            // 省市
	Page         *base.QueryOptions // 分页、排序、搜索参数
}
