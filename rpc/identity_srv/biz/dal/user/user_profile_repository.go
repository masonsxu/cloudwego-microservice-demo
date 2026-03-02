package user

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/errno"
)

// UserProfileRepositoryImpl 用户档案仓储实现
type UserProfileRepositoryImpl struct {
	base.BaseRepository[models.UserProfile]
	db *gorm.DB
}

// NewUserProfileRepository 创建用户档案仓储实例
func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return &UserProfileRepositoryImpl{
		BaseRepository: base.NewBaseRepository[models.UserProfile](db),
		db:             db,
	}
}

// ============================================================================
// 唯一性查询实现
// ============================================================================

// GetByUsername 根据用户名获取用户档案
func (r *UserProfileRepositoryImpl) GetByUsername(
	ctx context.Context,
	username string,
) (*models.UserProfile, error) {
	var user models.UserProfile

	err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, fmt.Errorf("根据用户名查询用户失败: %w", err)
	}

	return &user, nil
}

// GetByEmail 根据邮箱获取用户档案
func (r *UserProfileRepositoryImpl) GetByEmail(
	ctx context.Context,
	email string,
) (*models.UserProfile, error) {
	if email == "" {
		return nil, nil
	}

	var user models.UserProfile

	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ErrUserNotFound
		}

		return nil, fmt.Errorf("根据邮箱查询用户失败: %w", err)
	}

	return &user, nil
}

// GetByPhone 根据手机号获取用户档案
func (r *UserProfileRepositoryImpl) GetByPhone(
	ctx context.Context,
	phone string,
) (*models.UserProfile, error) {
	if phone == "" {
		return nil, nil
	}

	var user models.UserProfile

	err := r.db.WithContext(ctx).
		Where("phone = ?", phone).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ErrUserNotFound
		}

		return nil, fmt.Errorf("根据手机号查询用户失败: %w", err)
	}

	return &user, nil
}

// ExistsByID 检查用户ID是否存在
func (r *UserProfileRepositoryImpl) ExistsByID(
	ctx context.Context,
	userID string,
) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.UserProfile{}).
		Where("id = ?", userID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查用户ID是否存在失败: %w", err)
	}

	return count > 0, nil
}

// ============================================================================
// 状态管理操作实现
// ============================================================================

// UpdateStatus 更新用户状态
func (r *UserProfileRepositoryImpl) UpdateStatus(
	ctx context.Context,
	userID string,
	status models.UserStatus,
) error {
	result := r.db.WithContext(ctx).
		Model(&models.UserProfile{}).
		Where("id = ?", userID).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("更新用户状态失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在或已删除: %s", userID)
	}

	return nil
}

// UpdatePassword 更新用户密码
func (r *UserProfileRepositoryImpl) UpdatePassword(
	ctx context.Context,
	userID string,
	passwordHash string,
) error {
	result := r.db.WithContext(ctx).
		Model(&models.UserProfile{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password_hash":        passwordHash,
			"must_change_password": false, // Reset the must_change_password flag
			"login_attempts":       0,     // Also reset login attempts on password change
		})

	if result.Error != nil {
		return fmt.Errorf("更新用户密码失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在或已删除: %s", userID)
	}

	return nil
}

// UpdateLoginAttempts 更新登录尝试次数
func (r *UserProfileRepositoryImpl) UpdateLoginAttempts(
	ctx context.Context,
	userID string,
	attempts int32,
) error {
	result := r.db.WithContext(ctx).
		Model(&models.UserProfile{}).
		Where("id = ?", userID).
		Update("login_attempts", attempts)

	if result.Error != nil {
		return fmt.Errorf("更新登录尝试次数失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在或已删除: %s", userID)
	}

	return nil
}

// SetMustChangePassword 设置强制修改密码标志
func (r *UserProfileRepositoryImpl) SetMustChangePassword(
	ctx context.Context,
	userID string,
	mustChange bool,
) error {
	result := r.db.WithContext(ctx).
		Model(&models.UserProfile{}).
		Where("id = ?", userID).
		Update("must_change_password", mustChange)

	if result.Error != nil {
		return fmt.Errorf("设置强制修改密码标志失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在或已删除: %s", userID)
	}

	return nil
}

// ============================================================================
// 数据完整性检查实现
// ============================================================================

// CheckUsernameExists 检查用户名是否已存在
func (r *UserProfileRepositoryImpl) CheckUsernameExists(
	ctx context.Context,
	username string,
	excludeID ...string,
) (bool, error) {
	query := r.db.WithContext(ctx).Model(&models.UserProfile{}).
		Where("username = ?", username)

	if len(excludeID) > 0 && excludeID[0] != "" {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("检查用户名是否存在失败: %w", err)
	}

	return count > 0, nil
}

// CheckEmailExists 检查邮箱是否已存在
func (r *UserProfileRepositoryImpl) CheckEmailExists(
	ctx context.Context,
	email string,
	excludeID ...string,
) (bool, error) {
	if email == "" {
		return false, nil
	}

	query := r.db.WithContext(ctx).Model(&models.UserProfile{}).
		Where("email = ?", email)

	if len(excludeID) > 0 && excludeID[0] != "" {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("检查邮箱是否存在失败: %w", err)
	}

	return count > 0, nil
}

// CheckPhoneExists 检查手机号是否已存在
func (r *UserProfileRepositoryImpl) CheckPhoneExists(
	ctx context.Context,
	phone string,
	excludeID ...string,
) (bool, error) {
	if phone == "" {
		return false, nil
	}

	query := r.db.WithContext(ctx).Model(&models.UserProfile{}).
		Where("phone = ?", phone)

	if len(excludeID) > 0 && excludeID[0] != "" {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("检查手机号是否存在失败: %w", err)
	}

	return count > 0, nil
}

// ============================================================================
// 专业信息查询实现
// ============================================================================

// FindByMedicalLicense 根据执照号查询用户档案
func (r *UserProfileRepositoryImpl) FindByMedicalLicense(
	ctx context.Context,
	licenseNumber string,
) (*models.UserProfile, error) {
	if licenseNumber == "" {
		return nil, nil
	}

	var user models.UserProfile

	err := r.db.WithContext(ctx).
		Where("license_number = ?", licenseNumber).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ErrUserNotFound
		}

		return nil, fmt.Errorf("根据执照号查询用户失败: %w", err)
	}

	return &user, nil
}

// FindBySpecialty 根据专业领域查询用户档案列表
func (r *UserProfileRepositoryImpl) FindBySpecialty(
	ctx context.Context,
	specialty string,
	opts *base.QueryOptions,
) ([]*models.UserProfile, *models.PageResult, error) {
	if opts == nil {
		opts = base.NewQueryOptions()
	}

	// 构建查询条件：在JSON字段中搜索专业领域
	query := r.db.WithContext(ctx).Model(&models.UserProfile{}).
		Where("specialties LIKE ?", "%"+specialty+"%")

	// 处理其他过滤条件
	for field, value := range opts.Filters {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	// 处理搜索
	if opts.Search != "" {
		query = r.applySearchConditions(query, opts.Search)
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("查询专业领域用户总数失败: %w", err)
	}

	// 分页查询
	var users []*models.UserProfile

	offset := (opts.Page - 1) * opts.PageSize
	orderClause := r.buildOrderClause(opts)

	if err := query.Order(orderClause).Offset(int(offset)).Limit(int(opts.PageSize)).Find(&users).Error; err != nil {
		return nil, nil, fmt.Errorf("查询专业领域用户列表失败: %w", err)
	}

	// 构建分页结果
	pageResult := models.NewPageResult(int32(total), opts.Page, opts.PageSize)

	return users, pageResult, nil
}

// ============================================================================
// 统一查询方法实现
// ============================================================================

// FindWithConditions 根据组合查询条件查询用户档案列表
// 使用 QueryBuilder 模式重构，消除代码重复，提高可维护性
func (r *UserProfileRepositoryImpl) FindWithConditions(
	ctx context.Context,
	conditions *UserProfileQueryConditions,
) ([]*models.UserProfile, *models.PageResult, error) {
	// 初始化分页参数
	opts := base.NewQueryOptions()
	if conditions != nil && conditions.Page != nil {
		opts = conditions.Page
	}

	// 获取 BaseRepositoryImpl 实例以访问 QueryBuilder
	baseRepo := r.BaseRepository.(*base.BaseRepositoryImpl[models.UserProfile])

	// 使用 QueryBuilder 构建查询
	qb := baseRepo.NewQueryBuilder(ctx).
		WithSoftDelete(opts) // 应用软删除过滤（使用动态列名）

	// 处理组织过滤（需要 JOIN user_memberships 表）
	if conditions != nil && conditions.OrgID != nil && *conditions.OrgID != "" {
		qb = qb.WhereCustom(func(db *gorm.DB) *gorm.DB {
			return db.Joins("JOIN user_memberships um ON user_profiles.id = um.user_id").
				Where("um.organization_id = ? AND um.status = ?",
					*conditions.OrgID, models.MembershipStatusActive)
		})
	}

	// 应用精确匹配条件（需要表名前缀，因为可能有 JOIN）
	if conditions != nil {
		qb = r.applyUserProfileConditions(qb, conditions)
	}

	// 应用搜索、预加载和排序
	qb = qb.WithSearch(opts.Search, r.applySearchConditions). // 自定义搜索逻辑
									WithPreload(opts.Preloads...). // 预加载关联
									WithOrder(opts)                // 排序

	// 执行分页查询（自动处理计数和分页）
	return qb.FindWithPagination(opts)
}

// ============================================================================
// 私有辅助方法
// ============================================================================

// applyUserProfileConditions 应用用户档案查询条件
func (r *UserProfileRepositoryImpl) applyUserProfileConditions(
	qb *base.QueryBuilder[models.UserProfile],
	conditions *UserProfileQueryConditions,
) *base.QueryBuilder[models.UserProfile] {
	if conditions.Username != nil {
		qb = qb.WhereCustom(func(db *gorm.DB) *gorm.DB {
			return db.Where("user_profiles.username = ?", *conditions.Username)
		})
	}

	if conditions.Email != nil {
		qb = qb.WhereCustom(func(db *gorm.DB) *gorm.DB {
			return db.Where("user_profiles.email = ?", *conditions.Email)
		})
	}

	if conditions.Phone != nil {
		qb = qb.WhereCustom(func(db *gorm.DB) *gorm.DB {
			return db.Where("user_profiles.phone = ?", *conditions.Phone)
		})
	}

	if conditions.Status != nil {
		qb = qb.WhereCustom(func(db *gorm.DB) *gorm.DB {
			return db.Where("user_profiles.status = ?", *conditions.Status)
		})
	}

	if conditions.MedicalLicense != nil {
		qb = qb.WhereCustom(func(db *gorm.DB) *gorm.DB {
			return db.Where(
				"user_profiles.license_number = ?",
				*conditions.MedicalLicense,
			)
		})
	}

	// Specialty 模糊匹配
	if conditions.Specialty != nil {
		qb = qb.WhereCustom(func(db *gorm.DB) *gorm.DB {
			return db.Where("user_profiles.specialties LIKE ?", "%"+*conditions.Specialty+"%")
		})
	}

	return qb
}

// applySearchConditions 应用搜索条件
func (r *UserProfileRepositoryImpl) applySearchConditions(
	query *gorm.DB,
	searchTerm string,
) *gorm.DB {
	searchPattern := "%" + searchTerm + "%"

	// 在所有相关文本字段上应用搜索条件，并明确指定表名
	// ⚠️ 重要：使用括号包裹 OR 条件，避免与其他 AND 条件产生优先级问题
	return query.Where(
		`(user_profiles.username LIKE ? OR
		  user_profiles.email LIKE ? OR
		  user_profiles.first_name LIKE ? OR
		  user_profiles.last_name LIKE ? OR
		  user_profiles.real_name LIKE ? OR
		  user_profiles.medical_license_number LIKE ?)`,
		searchPattern,
		searchPattern,
		searchPattern,
		searchPattern,
		searchPattern,
		searchPattern,
	)
}

// buildOrderClause 构建排序子句
func (r *UserProfileRepositoryImpl) buildOrderClause(opts *base.QueryOptions) string {
	orderBy := "user_profiles.created_at" // 默认使用主表的全限定名

	if opts.OrderBy != "" {
		// 简单的安全检查，防止SQL注入
		// 在实际应用中，这里应该有一个更健壮的允许列表
		if strings.Contains(opts.OrderBy, ";") || strings.Contains(opts.OrderBy, " ") {
			orderBy = "user_profiles.created_at"
		} else {
			orderBy = "user_profiles." + opts.OrderBy
		}
	}

	if opts.OrderDesc {
		return orderBy + " DESC"
	}

	return orderBy + " ASC"
}

// WithTx 使用指定事务创建新的仓储实例
func (r *UserProfileRepositoryImpl) WithTx(tx *gorm.DB) base.BaseRepository[models.UserProfile] {
	return &UserProfileRepositoryImpl{
		BaseRepository: base.NewBaseRepository[models.UserProfile](tx),
		db:             tx,
	}
}

// ============================================================================
// 登录相关状态更新方法
// ============================================================================

// IncrementLoginAttempts 增加登录失败次数
func (r *UserProfileRepositoryImpl) IncrementLoginAttempts(
	ctx context.Context,
	userID string,
) error {
	result := r.db.WithContext(ctx).
		Model(&models.UserProfile{}).
		Where("id = ?", userID).
		Update("login_attempts", gorm.Expr("login_attempts + 1"))

	if result.Error != nil {
		return fmt.Errorf("增加用户登录失败次数失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在或已删除")
	}

	return nil
}

// ResetLoginAttempts 重置登录尝试次数
func (r *UserProfileRepositoryImpl) ResetLoginAttempts(
	ctx context.Context,
	userID string,
) error {
	result := r.db.WithContext(ctx).
		Model(&models.UserProfile{}).
		Where("id = ?", userID).
		Update("login_attempts", 0)

	if result.Error != nil {
		return fmt.Errorf("重置用户登录失败次数失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在或已删除")
	}

	return nil
}

// UpdateLastLoginTime 更新最后登录时间
func (r *UserProfileRepositoryImpl) UpdateLastLoginTime(
	ctx context.Context,
	userID string,
) error {
	// 使用当前时间戳
	currentTime := models.GetCurrentTimestamp()

	result := r.db.WithContext(ctx).
		Model(&models.UserProfile{}).
		Where("id = ?", userID).
		Update("last_login_time", currentTime)

	if result.Error != nil {
		return fmt.Errorf("更新用户最后登录时间失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在或已删除")
	}

	return nil
}

// ============================================================================
// 系统用户保护机制（重写基类方法）
// ============================================================================

// SoftDelete 软删除用户（增强保护）
// 重写基类方法，添加系统用户保护
func (r *UserProfileRepositoryImpl) SoftDelete(ctx context.Context, id string) error {
	// 1. 获取用户信息
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. 检查是否为系统用户
	if user.IsSystemUser {
		return errno.ErrSystemUserCannotDelete
	}

	// 3. 调用基类方法执行软删除
	return r.BaseRepository.SoftDelete(ctx, id)
}

// HardDelete 物理删除用户（增强保护）
// 重写基类方法，添加系统用户保护
func (r *UserProfileRepositoryImpl) HardDelete(ctx context.Context, id string) error {
	// 系统用户绝对禁止物理删除
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if user.IsSystemUser {
		return errno.ErrSystemUserCannotDelete
	}

	return r.BaseRepository.HardDelete(ctx, id)
}

// Update 更新用户（增强保护）
// 重写基类方法，添加系统用户保护
func (r *UserProfileRepositoryImpl) Update(ctx context.Context, user *models.UserProfile) error {
	if user == nil || user.ID.String() == "" {
		return errno.ErrInvalidParams.WithMessage("用户ID不能为空")
	}

	// 1. 获取原始用户信息
	original, err := r.GetByID(ctx, user.ID.String())
	if err != nil {
		return err
	}

	// 2. 如果是系统用户，检查关键属性是否被修改
	if original.IsSystemUser {
		// 禁止修改用户名
		if original.Username != user.Username {
			return errno.ErrSystemUserCannotModifyKey.WithMessage("系统用户的用户名不能被修改")
		}

		// 禁止修改 is_system_user 标识
		if !user.IsSystemUser {
			return errno.ErrSystemUserCannotModifyKey.WithMessage("系统用户的 is_system_user 标识不能被修改")
		}
	}

	// 3. 调用基类方法执行更新
	return r.BaseRepository.Update(ctx, user)
}

// ============================================================================
// 系统用户查询方法
// ============================================================================

// FindSystemUsers 查询所有系统用户
func (r *UserProfileRepositoryImpl) FindSystemUsers(
	ctx context.Context,
) ([]*models.UserProfile, error) {
	var users []*models.UserProfile

	err := r.db.WithContext(ctx).
		Where("is_system_user = ?", true).
		Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("查询系统用户失败: %w", err)
	}

	return users, nil
}

// IsSystemUser 判断用户是否为系统用户
func (r *UserProfileRepositoryImpl) IsSystemUser(ctx context.Context, userID string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.UserProfile{}).
		Where("id = ? AND is_system_user = ?", userID, true).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
