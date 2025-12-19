package membership

import (
	"context"
	"fmt"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"gorm.io/gorm"
)

// MembershipStatistics 成员关系统计信息
type MembershipStatistics struct {
	TotalMembers    int64 `json:"total_members"`    // 总成员数
	ActiveMembers   int64 `json:"active_members"`   // 活跃成员数
	InactiveMembers int64 `json:"inactive_members"` // 非活跃成员数
}

// UserMembershipRepositoryImpl 用户成员关系仓储实现
type UserMembershipRepositoryImpl struct {
	base.BaseRepository[models.UserMembership]
	db *gorm.DB
}

// NewUserMembershipRepository 创建用户成员关系仓储实例
func NewUserMembershipRepository(db *gorm.DB) UserMembershipRepository {
	return &UserMembershipRepositoryImpl{
		BaseRepository: base.NewBaseRepository[models.UserMembership](db),
		db:             db,
	}
}

// ============================================================================
// 核心成员关系查询实现
// ============================================================================

// FindAll 覆盖基础仓库的查询方法，实现成员关系特有的搜索逻辑
func (r *UserMembershipRepositoryImpl) FindAll(
	ctx context.Context,
	opts *base.QueryOptions,
) ([]*models.UserMembership, *models.PageResult, error) {
	if opts == nil {
		opts = base.NewQueryOptions()
	}

	// 构建查询
	query := r.buildMembershipQuery(ctx, opts)

	// 计算总数
	var total int64

	countQuery := r.buildMembershipCountQuery(ctx, opts)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("查询成员关系总数失败: %w", err)
	}

	// 分页查询
	var memberships []*models.UserMembership

	offset := (opts.Page - 1) * opts.PageSize

	orderClause := r.buildOrderClause(opts)
	if err := query.Order(orderClause).
		Offset(int(offset)).
		Limit(int(opts.PageSize)).
		Find(&memberships).Error; err != nil {
		return nil, nil, fmt.Errorf("查询成员关系列表失败: %w", err)
	}

	pageResult := models.NewPageResult(int32(total), opts.Page, opts.PageSize)

	return memberships, pageResult, nil
}

// buildMembershipQuery 构建成员关系查询语句
func (r *UserMembershipRepositoryImpl) buildMembershipQuery(
	ctx context.Context,
	opts *base.QueryOptions,
) *gorm.DB {
	query := r.db.WithContext(ctx).Model(&models.UserMembership{})

	// 处理过滤条件
	for field, value := range opts.Filters {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	return query
}

// buildMembershipCountQuery 构建成员关系计数查询语句
func (r *UserMembershipRepositoryImpl) buildMembershipCountQuery(
	ctx context.Context,
	opts *base.QueryOptions,
) *gorm.DB {
	query := r.db.WithContext(ctx).Model(&models.UserMembership{})

	// 处理过滤条件
	for field, value := range opts.Filters {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	return query
}

// GetByUserAndOrganization 获取用户在指定组织的成员关系
func (r *UserMembershipRepositoryImpl) GetByUserAndOrganization(
	ctx context.Context,
	userID, organizationID string,
) (*models.UserMembership, error) {
	var membership models.UserMembership

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND organization_id = ? AND status = ?", userID, organizationID, models.MembershipStatusActive).
		First(&membership).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, fmt.Errorf("查询用户组织成员关系失败: %w", err)
	}

	return &membership, nil
}

// GetByUserAndDepartment 获取用户在指定部门的成员关系
func (r *UserMembershipRepositoryImpl) GetByUserAndDepartment(
	ctx context.Context,
	userID, departmentID string,
) (*models.UserMembership, error) {
	var membership models.UserMembership

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND department_id = ? AND status = ?", userID, departmentID, models.MembershipStatusActive).
		First(&membership).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, fmt.Errorf("查询用户部门成员关系失败: %w", err)
	}

	return &membership, nil
}

// GetPrimaryMembership 获取用户的主要成员关系
func (r *UserMembershipRepositoryImpl) GetPrimaryMembership(
	ctx context.Context,
	userID string,
) (*models.UserMembership, error) {
	var membership models.UserMembership

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_primary = ? AND status = ?", userID, true, models.MembershipStatusActive).
		First(&membership).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, fmt.Errorf("查询用户主要成员关系失败: %w", err)
	}

	return &membership, nil
}

// GetPrimaryMembershipsByUserIDs 批量获取多个用户的主成员关系
func (r *UserMembershipRepositoryImpl) GetPrimaryMembershipsByUserIDs(
	ctx context.Context,
	userIDs []string,
) (map[string]*models.UserMembership, error) {
	// 空列表直接返回空map
	if len(userIDs) == 0 {
		return make(map[string]*models.UserMembership), nil
	}

	// 批量查询所有用户的主成员关系
	var memberships []*models.UserMembership

	err := r.db.WithContext(ctx).
		Where("user_id IN ? AND is_primary = ? AND status = ?",
			userIDs, true, models.MembershipStatusActive).
		Find(&memberships).Error
	if err != nil {
		return nil, fmt.Errorf("批量查询用户主要成员关系失败: %w", err)
	}

	// 转换为 map[userID]*UserMembership
	result := make(map[string]*models.UserMembership, len(memberships))
	for _, membership := range memberships {
		result[membership.UserID.String()] = membership
	}

	return result, nil
}

// ============================================================================
// 成员关系列表查询实现
// ============================================================================

// CountByDepartmentID 统计部门的成员数量
func (r *UserMembershipRepositoryImpl) CountByDepartmentID(
	ctx context.Context,
	departmentID string,
) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.UserMembership{}).
		Where("department_id = ? AND status = ?", departmentID, models.MembershipStatusActive).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("统计部门成员数量失败: %w", err)
	}

	return count, nil
}

// ============================================================================
// 组织层级查询实现
// ============================================================================

// FindUsersByOrganizationHierarchy 查询组织及其子组织的所有用户成员关系
func (r *UserMembershipRepositoryImpl) FindUsersByOrganizationHierarchy(
	ctx context.Context,
	rootOrganizationID string,
	opts *base.QueryOptions,
) ([]*models.UserMembership, *models.PageResult, error) {
	if opts == nil {
		opts = base.NewQueryOptions()
	}

	// 构建递归查询：获取组织及其所有子组织的ID
	// 这里简化实现，仅查询当前组织，实际项目中可使用WITH RECURSIVE或应用层递归
	subQuery := r.db.Model(&models.Organization{}).
		Select("id").
		Where("id = ? OR parent_id = ?", rootOrganizationID, rootOrganizationID)

	// 构建主查询
	query := r.db.WithContext(ctx).Model(&models.UserMembership{}).
		Where("organization_id IN (?) AND status = ?", subQuery, models.MembershipStatusActive)

	// 处理其他过滤条件
	for field, value := range opts.Filters {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("查询组织层级成员关系总数失败: %w", err)
	}

	// 分页查询
	var memberships []*models.UserMembership

	offset := (opts.Page - 1) * opts.PageSize
	orderClause := r.buildOrderClause(opts)

	if err := query.Order(orderClause).Offset(int(offset)).Limit(int(opts.PageSize)).Find(&memberships).Error; err != nil {
		return nil, nil, fmt.Errorf("查询组织层级成员关系失败: %w", err)
	}

	// 构建分页结果
	pageResult := models.NewPageResult(int32(total), opts.Page, opts.PageSize)

	return memberships, pageResult, nil
}

// GetUserOrganizations 获取用户所属的所有组织ID列表
func (r *UserMembershipRepositoryImpl) GetUserOrganizations(
	ctx context.Context,
	userID string,
) ([]string, error) {
	var organizationIDs []string

	err := r.db.WithContext(ctx).
		Model(&models.UserMembership{}).
		Select("organization_id").
		Where("user_id = ? AND status = ?", userID, models.MembershipStatusActive).
		Find(&organizationIDs).Error
	if err != nil {
		return nil, fmt.Errorf("获取用户组织列表失败: %w", err)
	}

	return organizationIDs, nil
}

// GetOrganizationUsers 获取组织的所有用户ID列表
func (r *UserMembershipRepositoryImpl) GetOrganizationUsers(
	ctx context.Context,
	organizationID string,
) ([]string, error) {
	var userIDs []string

	err := r.db.WithContext(ctx).
		Model(&models.UserMembership{}).
		Select("user_id").
		Where("organization_id = ? AND status = ?", organizationID, models.MembershipStatusActive).
		Find(&userIDs).Error
	if err != nil {
		return nil, fmt.Errorf("获取组织用户列表失败: %w", err)
	}

	return userIDs, nil
}

// ============================================================================
// 成员关系管理操作实现
// ============================================================================

// UpdateStatus 更新成员关系状态
func (r *UserMembershipRepositoryImpl) UpdateStatus(
	ctx context.Context,
	membershipID string,
	status models.MembershipStatus,
) error {
	result := r.db.WithContext(ctx).
		Model(&models.UserMembership{}).
		Where("id = ?", membershipID).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("更新成员关系状态失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("成员关系不存在: %s", membershipID)
	}

	return nil
}

// SetPrimaryMembership 设置主要成员关系
func (r *UserMembershipRepositoryImpl) SetPrimaryMembership(
	ctx context.Context,
	userID, membershipID string,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 将用户的所有成员关系的primary标志置为false
		if err := tx.Model(&models.UserMembership{}).
			Where("user_id = ?", userID).
			Update("is_primary", false).Error; err != nil {
			return fmt.Errorf("重置主要成员关系标志失败: %w", err)
		}

		// 2. 将指定成员关系的primary标志置为true
		result := tx.Model(&models.UserMembership{}).
			Where("id = ? AND user_id = ?", membershipID, userID).
			Update("is_primary", true)

		if result.Error != nil {
			return fmt.Errorf("设置主要成员关系失败: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("成员关系不存在或不属于指定用户: %s", membershipID)
		}

		return nil
	})
}

// ============================================================================
// 批量操作实现
// ============================================================================

// BatchCreateMemberships 批量创建成员关系
func (r *UserMembershipRepositoryImpl) BatchCreateMemberships(
	ctx context.Context,
	memberships []*models.UserMembership,
) error {
	if len(memberships) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).CreateInBatches(memberships, 100).Error; err != nil {
		return fmt.Errorf("批量创建成员关系失败: %w", err)
	}

	return nil
}

// BatchUpdateStatus 批量更新成员关系状态
func (r *UserMembershipRepositoryImpl) BatchUpdateStatus(
	ctx context.Context,
	membershipIDs []string,
	status models.MembershipStatus,
) error {
	if len(membershipIDs) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).
		Model(&models.UserMembership{}).
		Where("id IN ?", membershipIDs).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("批量更新成员关系状态失败: %w", result.Error)
	}

	return nil
}

// BatchDeleteByUser 批量删除用户的所有成员关系
func (r *UserMembershipRepositoryImpl) BatchDeleteByUser(ctx context.Context, userID string) error {
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.UserMembership{}).Error; err != nil {
		return fmt.Errorf("批量删除用户成员关系失败: %w", err)
	}

	return nil
}

// BatchDeleteByOrganization 批量删除组织的所有成员关系
func (r *UserMembershipRepositoryImpl) BatchDeleteByOrganization(
	ctx context.Context,
	organizationID string,
) error {
	if err := r.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Delete(&models.UserMembership{}).Error; err != nil {
		return fmt.Errorf("批量删除组织成员关系失败: %w", err)
	}

	return nil
}

// ============================================================================
// 数据完整性检查实现
// ============================================================================

// CheckMembershipExists 检查成员关系是否已存在
func (r *UserMembershipRepositoryImpl) CheckMembershipExists(
	ctx context.Context,
	userID, organizationID, departmentID string,
) (bool, error) {
	query := r.db.WithContext(ctx).Model(&models.UserMembership{}).
		Where("user_id = ? AND organization_id = ? AND status = ?", userID, organizationID, models.MembershipStatusActive)

	if departmentID != "" {
		query = query.Where("department_id = ?", departmentID)
	} else {
		query = query.Where("department_id IS NULL")
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("检查成员关系是否存在失败: %w", err)
	}

	return count > 0, nil
}

// ============================================================================
// 统计分析实现
// ============================================================================

// CountByOrganization 统计组织的成员数量
func (r *UserMembershipRepositoryImpl) CountByOrganization(
	ctx context.Context,
	organizationID string,
) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.UserMembership{}).
		Where("organization_id = ? AND status = ?", organizationID, models.MembershipStatusActive).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("统计组织成员数量失败: %w", err)
	}

	return count, nil
}

// GetMembershipStatistics 获取成员关系统计信息
func (r *UserMembershipRepositoryImpl) GetMembershipStatistics(
	ctx context.Context,
	organizationID string,
) (*MembershipStatistics, error) {
	stats := &MembershipStatistics{}

	// 统计总成员数
	var err error

	err = r.db.WithContext(ctx).
		Model(&models.UserMembership{}).
		Where("organization_id = ?", organizationID).
		Count(&stats.TotalMembers).
		Error
	if err != nil {
		return nil, fmt.Errorf("统计总成员数失败: %w", err)
	}

	// 统计活跃成员数
	err = r.db.WithContext(ctx).
		Model(&models.UserMembership{}).
		Where("organization_id = ? AND status = ?", organizationID, models.MembershipStatusActive).
		Count(&stats.ActiveMembers).Error
	if err != nil {
		return nil, fmt.Errorf("统计活跃成员数失败: %w", err)
	}

	// 计算非活跃成员数
	stats.InactiveMembers = stats.TotalMembers - stats.ActiveMembers

	return stats, nil
}

// ============================================================================
// 私有辅助方法
// ============================================================================

// buildOrderClause 构建排序子句
func (r *UserMembershipRepositoryImpl) buildOrderClause(opts *base.QueryOptions) string {
	orderBy := "created_at"
	if opts.OrderBy != "" {
		orderBy = opts.OrderBy
	}

	if opts.OrderDesc {
		return orderBy + " DESC"
	}

	return orderBy + " ASC"
}

// ============================================================================
// Logic层需要的简化查询方法实现
// ============================================================================

// UnsetPrimaryByUserID 取消用户的所有主要成员关系标志
func (r *UserMembershipRepositoryImpl) UnsetPrimaryByUserID(
	ctx context.Context,
	userID string,
) error {
	result := r.db.WithContext(ctx).
		Session(&gorm.Session{SkipHooks: true}).
		Model(&models.UserMembership{}).
		Where("user_id = ? AND is_primary = ?", userID, true).
		Update("is_primary", false)

	if result.Error != nil {
		return fmt.Errorf("取消用户主要成员关系标志失败: %w", result.Error)
	}

	return nil
}

// ExistsByID 检查成员关系是否存在
func (r *UserMembershipRepositoryImpl) ExistsByID(
	ctx context.Context,
	membershipID string,
) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.UserMembership{}).
		Where("id = ?", membershipID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查成员关系是否存在失败: %w", err)
	}

	return count > 0, nil
}

// ============================================================================
// 统一查询方法实现
// ============================================================================

// FindWithConditions 根据组合查询条件查询成员关系列表
// 使用 QueryBuilder 模式重构，消除代码重复，提高可维护性
func (r *UserMembershipRepositoryImpl) FindWithConditions(
	ctx context.Context,
	conditions *UserMembershipQueryConditions,
) ([]*models.UserMembership, *models.PageResult, error) {
	// 初始化分页参数
	opts := base.NewQueryOptions()
	if conditions != nil && conditions.Page != nil {
		opts = conditions.Page
	}

	// 获取 BaseRepositoryImpl 实例以访问 QueryBuilder
	baseRepo := r.BaseRepository.(*base.BaseRepositoryImpl[models.UserMembership])

	// 使用 QueryBuilder 构建查询
	qb := baseRepo.NewQueryBuilder(ctx).
		WithSoftDelete(opts) // 应用软删除过滤（使用动态列名）

	// 应用精确匹配条件（WhereEqual 会自动跳过 nil 值）
	if conditions != nil {
		qb = qb.WhereEqual("user_id", conditions.UserID).
			WhereEqual("organization_id", conditions.OrganizationID).
			WhereEqual("department_id", conditions.DepartmentID).
			WhereEqual("status", conditions.Status).
			WhereEqual("is_primary", conditions.IsPrimary)
	}

	// 应用搜索、预加载和排序
	qb = qb.WithSearch(opts.Search, r.applySearchConditions). // 自定义搜索逻辑
									WithPreload(opts.Preloads...). // 预加载关联
									WithOrder(opts)                // 排序

	// 执行分页查询（自动处理计数和分页）
	return qb.FindWithPagination(opts)
}

// applySearchConditions 应用搜索条件
func (r *UserMembershipRepositoryImpl) applySearchConditions(
	query *gorm.DB,
	searchTerm string,
) *gorm.DB {
	// 成员关系搜索逻辑可以根据需要扩展
	// 目前可以基于 user_id, organization_id 等字段进行搜索
	searchPattern := "%" + searchTerm + "%"

	// ⚠️ 重要：使用括号包裹 OR 条件，避免与其他 AND 条件产生优先级问题
	return query.Where(
		"(user_id LIKE ? OR organization_id LIKE ? OR department_id LIKE ?)",
		searchPattern, searchPattern, searchPattern,
	)
}

// WithTx 使用指定事务创建新的仓储实例
func (r *UserMembershipRepositoryImpl) WithTx(
	tx *gorm.DB,
) base.BaseRepository[models.UserMembership] {
	return &UserMembershipRepositoryImpl{
		BaseRepository: base.NewBaseRepository[models.UserMembership](tx),
		db:             tx,
	}
}
