package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"gorm.io/gorm"
)

// OrganizationRepositoryImpl 组织仓储实现（简化版）
// 仅支持2级层级结构：根组织→子组织
type OrganizationRepositoryImpl struct {
	base.BaseRepository[models.Organization]
	db *gorm.DB
}

// NewOrganizationRepository 创建组织仓储实例
func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &OrganizationRepositoryImpl{
		BaseRepository: base.NewBaseRepository[models.Organization](db),
		db:             db,
	}
}

// ============================================================================
// 核心查询方法实现
// ============================================================================

// FindAll 覆盖基础仓库的查询方法，实现组织特有的搜索逻辑
func (r *OrganizationRepositoryImpl) FindAll(
	ctx context.Context,
	opts *base.QueryOptions,
) ([]*models.Organization, *models.PageResult, error) {
	if opts == nil {
		opts = base.NewQueryOptions()
	}

	// 获取 BaseRepositoryImpl 实例以访问 QueryBuilder
	baseRepo := r.BaseRepository.(*base.BaseRepositoryImpl[models.Organization])

	// 使用 QueryBuilder 构建查询
	qb := baseRepo.NewQueryBuilder(ctx).
		WithSoftDelete(opts).                             // 应用软删除过滤
		WithSearch(opts.Search, r.applySearchConditions). // 自定义搜索逻辑
		WithPreload(opts.Preloads...).                    // 预加载关联
		WithOrder(opts)                                   // 排序

	// 应用过滤条件
	for field, value := range opts.Filters {
		qb = qb.WhereEqual(field, value)
	}

	// 执行分页查询
	return qb.FindWithPagination(opts)
}

// GetByCode 根据组织代码获取组织
func (r *OrganizationRepositoryImpl) GetByCode(
	ctx context.Context,
	code string,
) (*models.Organization, error) {
	var org models.Organization

	err := r.db.WithContext(ctx).Where("code = ?", code).First(&org).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, fmt.Errorf("根据代码查询组织失败: %w", err)
	}

	return &org, nil
}

// ============================================================================
// 简化的关系管理实现
// ============================================================================

// UpdateParent 更新组织的父组织（仅支持2级层级）
func (r *OrganizationRepositoryImpl) UpdateParent(
	ctx context.Context,
	organizationID, newParentID string,
) error {
	// 如果设置了新的父组织，需要验证2级层级限制
	if newParentID != "" {
		// 检查新父组织是否存在
		var parent models.Organization

		err := r.db.WithContext(ctx).Where("id = ?", newParentID).First(&parent).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("新父组织不存在: %s", newParentID)
			}

			return fmt.Errorf("查询新父组织失败: %w", err)
		}

		// 验证2级层级限制：新父组织不能有父组织
		if parent.ParentID != uuid.Nil {
			return fmt.Errorf("不支持超过2级的组织层级，新父组织不能为子组织")
		}
	}

	// 更新父组织
	result := r.db.WithContext(ctx).Model(&models.Organization{}).
		Where("id = ?", organizationID).
		Update("parent_id", newParentID)

	if result.Error != nil {
		return fmt.Errorf("更新组织父级失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("组织不存在: %s", organizationID)
	}

	return nil
}

// HasChildren 检查组织是否有直接子组织
func (r *OrganizationRepositoryImpl) HasChildren(
	ctx context.Context,
	organizationID string,
) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.Organization{}).
		Where("parent_id = ?", organizationID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查子组织失败: %w", err)
	}

	return count > 0, nil
}

// CountChildren 统计直接子组织数量
func (r *OrganizationRepositoryImpl) CountChildren(
	ctx context.Context,
	parentID string,
) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.Organization{}).
		Where("parent_id = ?", parentID).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("统计子组织数量失败: %w", err)
	}

	return count, nil
}

// ExistsByID 检查组织是否存在
func (r *OrganizationRepositoryImpl) ExistsByID(
	ctx context.Context,
	organizationID string,
) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.Organization{}).
		Where("id = ?", organizationID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查组织是否存在失败: %w", err)
	}

	return count > 0, nil
}

// ============================================================================
// 数据完整性检查实现
// ============================================================================

// CheckNameConflict 检查组织名称是否冲突
func (r *OrganizationRepositoryImpl) CheckNameConflict(
	ctx context.Context,
	name string,
	parentID string,
	excludeID ...string,
) (bool, error) {
	query := r.db.WithContext(ctx).Model(&models.Organization{}).
		Where("name = ? AND parent_id = ?", name, parentID)

	// 排除指定ID
	if len(excludeID) > 0 && excludeID[0] != "" {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64

	err := query.Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查名称冲突失败: %w", err)
	}

	return count > 0, nil
}

// ============================================================================
// 统一查询方法实现
// ============================================================================

// FindWithConditions 根据组合查询条件查询组织列表
// 使用 QueryBuilder 模式重构，消除代码重复，提高可维护性
func (r *OrganizationRepositoryImpl) FindWithConditions(
	ctx context.Context,
	conditions *OrganizationQueryConditions,
) ([]*models.Organization, *models.PageResult, error) {
	// 初始化分页参数
	opts := base.NewQueryOptions()
	if conditions != nil && conditions.Page != nil {
		opts = conditions.Page
	}

	// 获取 BaseRepositoryImpl 实例以访问 QueryBuilder
	baseRepo := r.BaseRepository.(*base.BaseRepositoryImpl[models.Organization])

	// 使用 QueryBuilder 构建查询
	qb := baseRepo.NewQueryBuilder(ctx).
		WithSoftDelete(opts) // 应用软删除过滤（使用动态列名）

	// 应用精确匹配条件（WhereEqual 会自动跳过 nil 值）
	if conditions != nil {
		qb = qb.WhereEqual("code", conditions.Code).
			WhereEqual("facility_type", conditions.FacilityType).
			WhereEqual("province_city", conditions.ProvinceCity)

		// ParentID 特殊处理：空字符串表示查询根组织
		if conditions.ParentID != nil {
			if *conditions.ParentID == "" {
				// 空字符串表示查询根组织（parent_id IS NULL 或 uuid.Nil）
				qb = qb.WhereCustom(func(db *gorm.DB) *gorm.DB {
					return db.Where("parent_id IS NULL OR parent_id = ?", uuid.Nil)
				})
			} else {
				qb = qb.WhereEqual("parent_id", conditions.ParentID)
			}
		}

		// Name 模糊匹配
		if conditions.Name != nil && *conditions.Name != "" {
			qb = qb.WhereCustom(func(db *gorm.DB) *gorm.DB {
				return db.Where("name LIKE ?", "%"+*conditions.Name+"%")
			})
		}
	}

	// 应用搜索、预加载和排序
	qb = qb.WithSearch(opts.Search, r.applySearchConditions). // 自定义搜索逻辑
									WithPreload(opts.Preloads...). // 预加载关联
									WithOrder(opts)                // 排序

	// 执行分页查询（自动处理计数和分页）
	return qb.FindWithPagination(opts)
}

// applySearchConditions 应用搜索条件
// 定义组织特有的多字段搜索逻辑
func (r *OrganizationRepositoryImpl) applySearchConditions(
	query *gorm.DB,
	searchTerm string,
) *gorm.DB {
	searchPattern := "%" + searchTerm + "%"

	// ⚠️ 重要：使用括号包裹 OR 条件，避免与其他 AND 条件产生优先级问题
	return query.Where(
		"(name LIKE ? OR code LIKE ? OR facility_type LIKE ?)",
		searchPattern, searchPattern, searchPattern,
	)
}
