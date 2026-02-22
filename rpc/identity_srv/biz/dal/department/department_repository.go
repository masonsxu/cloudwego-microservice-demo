package department

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// DepartmentRepositoryImpl 部门仓储实现
type DepartmentRepositoryImpl struct {
	base.BaseRepository[models.Department]
	db *gorm.DB
}

// NewDepartmentRepository 创建部门仓储实例
func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &DepartmentRepositoryImpl{
		BaseRepository: base.NewBaseRepository[models.Department](db),
		db:             db,
	}
}

// ============================================================================
// 部门查询实现
// ============================================================================

// ExistsByID 检查部门是否存在
func (r *DepartmentRepositoryImpl) ExistsByID(
	ctx context.Context,
	departmentID string,
) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.Department{}).
		Where("id = ?", departmentID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查部门是否存在失败: %w", err)
	}

	return count > 0, nil
}

// GetByName 根据名称获取部门（在指定组织内）
func (r *DepartmentRepositoryImpl) GetByName(
	ctx context.Context,
	name string,
	organizationID string,
) (*models.Department, error) {
	var dept models.Department

	err := r.db.WithContext(ctx).
		Where("name = ? AND organization_id = ?", name, organizationID).
		First(&dept).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, fmt.Errorf("根据名称查询部门失败: %w", err)
	}

	return &dept, nil
}

// ============================================================================
// 设备管理相关实现
// ============================================================================

// GetDepartmentEquipment 获取部门的所有可用设备ID列表
func (r *DepartmentRepositoryImpl) GetDepartmentEquipment(
	ctx context.Context,
	departmentID string,
) ([]string, error) {
	var dept models.Department

	err := r.db.WithContext(ctx).
		Select("available_equipment").
		Where("id = ?", departmentID).
		First(&dept).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("部门不存在: %s", departmentID)
		}

		return nil, fmt.Errorf("获取部门设备失败: %w", err)
	}

	// 解析JSON字段（这里简化处理，实际应使用JSON解析）
	// 假设AvailableEquipment已经是字符串切片
	equipmentIDs := make([]string, 0)
	// TODO: 实现JSON解析逻辑

	return equipmentIDs, nil
}

// AddEquipment 为部门添加设备
func (r *DepartmentRepositoryImpl) AddEquipment(
	ctx context.Context,
	departmentID string,
	equipmentIDs []string,
) error {
	// 获取当前设备列表
	currentEquipment, err := r.GetDepartmentEquipment(ctx, departmentID)
	if err != nil {
		return fmt.Errorf("获取当前设备列表失败: %w", err)
	}

	// 合并设备列表并去重
	equipmentMap := make(map[string]bool)
	for _, id := range currentEquipment {
		equipmentMap[id] = true
	}

	for _, id := range equipmentIDs {
		equipmentMap[id] = true
	}

	// 构建新的设备列表
	newEquipment := make([]string, 0, len(equipmentMap))
	for id := range equipmentMap {
		newEquipment = append(newEquipment, id)
	}

	return r.UpdateEquipment(ctx, departmentID, newEquipment)
}

// RemoveEquipment 从部门移除设备
func (r *DepartmentRepositoryImpl) RemoveEquipment(
	ctx context.Context,
	departmentID string,
	equipmentIDs []string,
) error {
	// 获取当前设备列表
	currentEquipment, err := r.GetDepartmentEquipment(ctx, departmentID)
	if err != nil {
		return fmt.Errorf("获取当前设备列表失败: %w", err)
	}

	// 创建要移除的设备ID映射
	removeMap := make(map[string]bool)
	for _, id := range equipmentIDs {
		removeMap[id] = true
	}

	// 过滤设备列表
	newEquipment := make([]string, 0)

	for _, id := range currentEquipment {
		if !removeMap[id] {
			newEquipment = append(newEquipment, id)
		}
	}

	return r.UpdateEquipment(ctx, departmentID, newEquipment)
}

// UpdateEquipment 更新部门的设备列表
func (r *DepartmentRepositoryImpl) UpdateEquipment(
	ctx context.Context,
	departmentID string,
	equipmentIDs []string,
) error {
	// TODO: 实现JSON序列化逻辑
	// 这里简化处理，实际应将equipmentIDs序列化为JSON字符串
	equipmentJSON := ""

	result := r.db.WithContext(ctx).
		Model(&models.Department{}).
		Where("id = ?", departmentID).
		Update("available_equipment", equipmentJSON)

	if result.Error != nil {
		return fmt.Errorf("更新部门设备失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("部门不存在: %s", departmentID)
	}

	return nil
}

// ============================================================================
// 成员关系查询实现
// ============================================================================

// GetDepartmentMembers 获取部门成员ID列表
func (r *DepartmentRepositoryImpl) GetDepartmentMembers(
	ctx context.Context,
	departmentID string,
) ([]string, error) {
	var userIDs []string

	err := r.db.WithContext(ctx).
		Model(&models.UserMembership{}).
		Select("user_id").
		Where("department_id = ? AND status = ?", departmentID, models.MembershipStatusActive).
		Find(&userIDs).Error
	if err != nil {
		return nil, fmt.Errorf("获取部门成员列表失败: %w", err)
	}

	return userIDs, nil
}

// HasMembers 检查部门是否有成员
func (r *DepartmentRepositoryImpl) HasMembers(
	ctx context.Context,
	departmentID string,
) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.UserMembership{}).
		Where("department_id = ? AND status = ?", departmentID, models.MembershipStatusActive).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查部门成员存在性失败: %w", err)
	}

	return count > 0, nil
}

// CountMembers 统计部门成员数量
func (r *DepartmentRepositoryImpl) CountMembers(
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
// 数据完整性检查实现
// ============================================================================

// CheckNameExists 检查部门名称是否已存在（在指定组织内）
func (r *DepartmentRepositoryImpl) CheckNameExists(
	ctx context.Context,
	name string,
	organizationID string,
	excludeID ...string,
) (bool, error) {
	query := r.db.WithContext(ctx).Model(&models.Department{}).
		Where("name = ? AND organization_id = ?", name, organizationID)

	if len(excludeID) > 0 && excludeID[0] != "" {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("检查部门名称是否存在失败: %w", err)
	}

	return count > 0, nil
}

// ValidateOrganizationExists 验证组织是否存在
func (r *DepartmentRepositoryImpl) ValidateOrganizationExists(
	ctx context.Context,
	organizationID string,
) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.Organization{}).
		Where("id = ?", organizationID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("验证组织存在性失败: %w", err)
	}

	return count > 0, nil
}

// ============================================================================
// 批量操作实现
// ============================================================================

// BatchCreateDepartments 批量创建部门
func (r *DepartmentRepositoryImpl) BatchCreateDepartments(
	ctx context.Context,
	departments []*models.Department,
) error {
	if len(departments) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).CreateInBatches(departments, 100).Error; err != nil {
		return fmt.Errorf("批量创建部门失败: %w", err)
	}

	return nil
}

// BatchUpdateOrganization 批量更新部门的组织归属
func (r *DepartmentRepositoryImpl) BatchUpdateOrganization(
	ctx context.Context,
	departmentIDs []string,
	newOrganizationID string,
) error {
	if len(departmentIDs) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).
		Model(&models.Department{}).
		Where("id IN ?", departmentIDs).
		Update("organization_id", newOrganizationID)

	if result.Error != nil {
		return fmt.Errorf("批量更新部门组织归属失败: %w", result.Error)
	}

	return nil
}

// BatchDeleteByOrganization 批量删除指定组织的所有部门
func (r *DepartmentRepositoryImpl) BatchDeleteByOrganization(
	ctx context.Context,
	organizationID string,
) error {
	if err := r.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Delete(&models.Department{}).Error; err != nil {
		return fmt.Errorf("批量删除组织部门失败: %w", err)
	}

	return nil
}

// ============================================================================
// 统计分析实现
// ============================================================================

// CountByOrganization 统计指定组织的部门数量
func (r *DepartmentRepositoryImpl) CountByOrganization(
	ctx context.Context,
	organizationID string,
) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.Department{}).
		Where("organization_id = ?", organizationID).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("统计组织部门数量失败: %w", err)
	}

	return count, nil
}

// CountByDepartmentType 统计指定类型的部门数量
func (r *DepartmentRepositoryImpl) CountByDepartmentType(
	ctx context.Context,
	departmentType string,
) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.Department{}).
		Where("department_type = ?", departmentType).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("统计部门类型数量失败: %w", err)
	}

	return count, nil
}

// GetDepartmentStatistics 获取部门统计信息
func (r *DepartmentRepositoryImpl) GetDepartmentStatistics(
	ctx context.Context,
	departmentID string,
) (*DepartmentStatistics, error) {
	stats := &DepartmentStatistics{}

	// 统计成员数量
	var err error

	stats.MembersCount, err = r.CountMembers(ctx, departmentID)
	if err != nil {
		return nil, fmt.Errorf("统计部门成员数量失败: %w", err)
	}

	// 统计活跃成员数量（与总成员数相同，因为查询条件已过滤）
	stats.ActiveMembers = stats.MembersCount

	// 统计设备数量
	equipmentIDs, err := r.GetDepartmentEquipment(ctx, departmentID)
	if err != nil {
		return nil, fmt.Errorf("获取部门设备失败: %w", err)
	}

	stats.EquipmentCount = int64(len(equipmentIDs))

	return stats, nil
}

// GetOrganizationDepartmentStatistics 获取组织的部门统计信息
func (r *DepartmentRepositoryImpl) GetOrganizationDepartmentStatistics(
	ctx context.Context,
	organizationID string,
) (*OrganizationDepartmentStatistics, error) {
	stats := &OrganizationDepartmentStatistics{
		TypeDistribution: make(map[string]int64),
	}

	// 统计总部门数
	var err error

	stats.TotalDepartments, err = r.CountByOrganization(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("统计组织部门数量失败: %w", err)
	}

	// 统计部门类型分布
	var typeStats []struct {
		DepartmentType string `gorm:"column:department_type"`
		Count          int64  `gorm:"column:count"`
	}

	err = r.db.WithContext(ctx).
		Model(&models.Department{}).
		Select("department_type, COUNT(*) as count").
		Where("organization_id = ? AND department_type != ''", organizationID).
		Group("department_type").
		Find(&typeStats).Error
	if err != nil {
		return nil, fmt.Errorf("统计部门类型分布失败: %w", err)
	}

	for _, stat := range typeStats {
		if stat.DepartmentType != "" {
			stats.TypeDistribution[stat.DepartmentType] = stat.Count
		}
	}

	// 统计总成员数
	err = r.db.WithContext(ctx).
		Table("user_memberships um").
		Joins("JOIN departments d ON um.department_id = d.id").
		Where("d.organization_id = ? AND um.status = ?", organizationID, models.MembershipStatusActive).
		Count(&stats.TotalMembers).Error
	if err != nil {
		return nil, fmt.Errorf("统计组织总成员数失败: %w", err)
	}

	// 计算平均值
	if stats.TotalDepartments > 0 {
		stats.AverageMembers = float64(stats.TotalMembers) / float64(stats.TotalDepartments)
		// TODO: 计算平均设备数（需要解析JSON字段）
		stats.AverageEquipment = 0
	}

	return stats, nil
}

// ============================================================================
// 统一查询方法实现
// ============================================================================

// FindWithConditions 根据组合查询条件查询部门列表
func (r *DepartmentRepositoryImpl) FindWithConditions(
	ctx context.Context,
	conditions *DepartmentQueryConditions,
) ([]*models.Department, *models.PageResult, error) {
	// 初始化分页参数
	opts := base.NewQueryOptions()
	if conditions != nil && conditions.Page != nil {
		opts = conditions.Page
	}

	// 构建基础查询并应用过滤条件
	query := r.db.WithContext(ctx).Model(&models.Department{})
	query = r.applyFilterConditions(query, conditions, opts)

	// 处理搜索条件
	if opts.Search != "" {
		query = r.applySearchConditions(query, opts.Search)
	}

	// 处理预加载
	for _, preload := range opts.Preloads {
		query = query.Preload(preload)
	}

	// 计算总数
	var total int64

	countQuery := r.db.WithContext(ctx).Model(&models.Department{})
	countQuery = r.applyFilterConditions(countQuery, conditions, opts)

	if opts.Search != "" {
		countQuery = r.applySearchConditions(countQuery, opts.Search)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("统计部门总数失败: %w", err)
	}

	// 处理排序
	orderClause := r.buildOrderClause(opts)
	query = query.Order(orderClause)

	// 分页查询
	var departments []*models.Department

	offset := (opts.Page - 1) * opts.PageSize
	if err := query.Offset(int(offset)).Limit(int(opts.PageSize)).Find(&departments).Error; err != nil {
		return nil, nil, fmt.Errorf("查询部门列表失败: %w", err)
	}

	// 构建分页结果
	pageResult := models.NewPageResult(int32(total), opts.Page, opts.PageSize)

	return departments, pageResult, nil
}

// applyFilterConditions 应用过滤条件到查询
func (r *DepartmentRepositoryImpl) applyFilterConditions(
	query *gorm.DB,
	conditions *DepartmentQueryConditions,
	opts *base.QueryOptions,
) *gorm.DB {
	// 处理软删除
	if !opts.IncludeDeleted && !opts.DeletedOnly {
		// 默认：仅查询未删除记录
	} else if opts.DeletedOnly {
		query = query.Unscoped().Where("deleted_at IS NOT NULL")
	} else if opts.IncludeDeleted {
		query = query.Unscoped()
	}

	// 处理设备过滤（需要 JOIN department_equipment 表）
	if conditions != nil && conditions.EquipmentID != nil && *conditions.EquipmentID != "" {
		// 使用 JSONB 查询（PostgreSQL）
		query = query.Where("equipment_ids::jsonb ? ?", *conditions.EquipmentID)
	}

	// 处理精确匹配条件
	if conditions != nil {
		if conditions.OrganizationID != nil {
			query = query.Where("organization_id = ?", *conditions.OrganizationID)
		}

		if conditions.DepartmentType != nil {
			query = query.Where("department_type = ?", *conditions.DepartmentType)
		}

		// 名称模糊匹配
		if conditions.Name != nil && *conditions.Name != "" {
			query = query.Where("name LIKE ?", "%"+*conditions.Name+"%")
		}
	}

	return query
}

// ============================================================================
// 私有辅助方法
// ============================================================================

// applySearchConditions 应用搜索条件
func (r *DepartmentRepositoryImpl) applySearchConditions(
	query *gorm.DB,
	searchTerm string,
) *gorm.DB {
	searchPattern := "%" + searchTerm + "%"

	// ⚠️ 重要：使用括号包裹 OR 条件，避免与其他 AND 条件产生优先级问题
	return query.Where(
		"(name LIKE ? OR department_type LIKE ?)",
		searchPattern, searchPattern,
	)
}

// buildOrderClause 构建排序子句
func (r *DepartmentRepositoryImpl) buildOrderClause(opts *base.QueryOptions) string {
	orderBy := "created_at"
	if opts.OrderBy != "" {
		orderBy = opts.OrderBy
	}

	if opts.OrderDesc {
		return orderBy + " DESC"
	}

	return orderBy + " ASC"
}

// WithTx 使用指定事务创建新的仓储实例
func (r *DepartmentRepositoryImpl) WithTx(tx *gorm.DB) base.BaseRepository[models.Department] {
	return &DepartmentRepositoryImpl{
		BaseRepository: base.NewBaseRepository[models.Department](tx),
		db:             tx,
	}
}
