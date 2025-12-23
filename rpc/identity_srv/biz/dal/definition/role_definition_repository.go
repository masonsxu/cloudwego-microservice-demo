package definition

import (
	"context"

	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// RoleDefinitionRepositoryImpl 角色定义仓储实现
// 提供完整的角色定义数据访问功能，包含权限查询和系统角色管理
type RoleDefinitionRepositoryImpl struct {
	base.BaseRepository[models.RoleDefinition]
	db *gorm.DB
}

// NewRoleDefinitionRepository 创建角色定义仓储实例
func NewRoleDefinitionRepository(db *gorm.DB) RoleDefinitionRepository {
	return &RoleDefinitionRepositoryImpl{
		BaseRepository: base.NewBaseRepository[models.RoleDefinition](db),
		db:             db,
	}
}

// ============================================================================
// 核心查询方法
// ============================================================================

// FindByName 根据角色名称查询角色定义（精确匹配）
func (r *RoleDefinitionRepositoryImpl) FindByName(
	ctx context.Context,
	name string,
) (*models.RoleDefinition, error) {
	var role models.RoleDefinition

	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &role, nil
}

// CheckNameExists 检查指定角色名称是否已存在
func (r *RoleDefinitionRepositoryImpl) CheckNameExists(
	ctx context.Context,
	name string,
) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.RoleDefinition{}).
		Where("name = ?", name).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// FindByStatus 根据角色状态查询角色定义列表
func (r *RoleDefinitionRepositoryImpl) FindByStatus(
	ctx context.Context,
	status models.RoleStatus,
	page *base.QueryOptions,
) ([]*models.RoleDefinition, *models.PageResult, error) {
	if page == nil {
		page = base.NewQueryOptions()
	}

	opts := base.NewQueryOptions().
		WithFilter("status", status).
		WithPage(page.Page, page.PageSize).
		WithOrder(page.OrderBy, page.OrderDesc).
		WithFetchAll(page.FetchAll)

	return r.FindAll(ctx, opts)
}

// FindBySystemRole 根据系统角色标识查询角色定义列表
func (r *RoleDefinitionRepositoryImpl) FindBySystemRole(
	ctx context.Context,
	isSystemRole bool,
	page *base.QueryOptions,
) ([]*models.RoleDefinition, *models.PageResult, error) {
	if page == nil {
		page = base.NewQueryOptions()
	}

	opts := base.NewQueryOptions().
		WithFilter("is_system_role", isSystemRole).
		WithPage(page.Page, page.PageSize).
		WithOrder(page.OrderBy, page.OrderDesc).
		WithFetchAll(page.FetchAll)

	return r.FindAll(ctx, opts)
}

// FindWithConditions 根据组合查询条件查询角色定义列表
func (r *RoleDefinitionRepositoryImpl) FindWithConditions(
	ctx context.Context,
	conditions *RoleDefinitionQueryConditions,
) ([]*models.RoleDefinition, *models.PageResult, error) {
	opts := base.NewQueryOptions()

	if conditions != nil {
		if conditions.Name != nil {
			opts = opts.WithFilter("name", *conditions.Name)
		}

		if conditions.Status != nil {
			opts = opts.WithFilter("status", *conditions.Status)
		}

		if conditions.IsSystemRole != nil {
			opts = opts.WithFilter("is_system_role", *conditions.IsSystemRole)
		}

		if conditions.Page != nil {
			opts = opts.WithPage(conditions.Page.Page, conditions.Page.PageSize).
				WithOrder(conditions.Page.OrderBy, conditions.Page.OrderDesc).
				WithFetchAll(conditions.Page.FetchAll)
		}
	}

	return r.FindAll(ctx, opts)
}

// ============================================================================
// 业务统计方法
// ============================================================================

// CountByStatus 统计指定状态的角色定义数量
func (r *RoleDefinitionRepositoryImpl) CountByStatus(
	ctx context.Context,
	status models.RoleStatus,
) (int64, error) {
	opts := base.NewQueryOptions().WithFilter("status", status)
	return r.Count(ctx, opts)
}

// ListActiveRoles 列出所有活跃状态的角色定义
func (r *RoleDefinitionRepositoryImpl) ListActiveRoles(
	ctx context.Context,
	page *base.QueryOptions,
) ([]*models.RoleDefinition, *models.PageResult, error) {
	if page == nil {
		page = base.NewQueryOptions()
	}

	opts := base.NewQueryOptions().
		WithFilter("status", models.RoleStatusActive).
		WithPage(page.Page, page.PageSize).
		WithOrder(page.OrderBy, page.OrderDesc).
		WithFetchAll(page.FetchAll)

	return r.FindAll(ctx, opts)
}
