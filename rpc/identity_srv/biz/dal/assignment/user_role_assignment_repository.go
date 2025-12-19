package assignment

import (
	"context"

	"github.com/google/uuid"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"gorm.io/gorm"
)

// UserRoleAssignmentRepositoryImpl 用户角色分配仓储实现
// 提供完整的RBAC角色分配管理功能，支持多维度查询和业务约束验证
type UserRoleAssignmentRepositoryImpl struct {
	base.BaseRepository[models.UserRoleAssignment]
	db *gorm.DB
}

// NewUserRoleAssignmentRepository 创建用户角色分配仓储实例
func NewUserRoleAssignmentRepository(db *gorm.DB) UserRoleAssignmentRepository {
	return &UserRoleAssignmentRepositoryImpl{
		BaseRepository: base.NewBaseRepository[models.UserRoleAssignment](db),
		db:             db,
	}
}

// ============================================================================
// 核心查询方法
// ============================================================================

// FindByUserID 查询指定用户的所有角色分配记录
func (r *UserRoleAssignmentRepositoryImpl) FindByUserID(
	ctx context.Context,
	userID string,
	page *base.QueryOptions,
) ([]*models.UserRoleAssignment, *models.PageResult, error) {
	if page == nil {
		page = base.NewQueryOptions()
	}

	opts := base.NewQueryOptions().
		WithFilter("user_id", userID).
		WithPage(page.Page, page.PageSize).
		WithOrder(page.OrderBy, page.OrderDesc)

	return r.FindAll(ctx, opts)
}

// FindByRoleID 查询指定角色的所有分配记录
func (r *UserRoleAssignmentRepositoryImpl) FindByRoleID(
	ctx context.Context,
	roleID string,
	page *base.QueryOptions,
) ([]*models.UserRoleAssignment, *models.PageResult, error) {
	if page == nil {
		page = base.NewQueryOptions()
	}

	opts := base.NewQueryOptions().
		WithFilter("role_id", roleID).
		WithPage(page.Page, page.PageSize).
		WithOrder(page.OrderBy, page.OrderDesc)

	return r.FindAll(ctx, opts)
}

// FindByUserAndRole 查询指定用户和角色的分配记录
func (r *UserRoleAssignmentRepositoryImpl) FindByUserAndRole(
	ctx context.Context,
	userID, roleID string,
) (*models.UserRoleAssignment, error) {
	var assignment models.UserRoleAssignment

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		First(&assignment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &assignment, nil
}

// GetLastUserRoleAssignment 获取用户最后一次的角色分配信息
func (r *UserRoleAssignmentRepositoryImpl) GetLastUserRoleAssignment(
	ctx context.Context,
	userID string,
) (*models.UserRoleAssignment, error) {
	var assignment models.UserRoleAssignment

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&assignment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &assignment, nil
}

// GetActiveRolesByUserID 获取用户所有活跃的角色ID列表（未删除）
func (r *UserRoleAssignmentRepositoryImpl) GetActiveRolesByUserID(
	ctx context.Context,
	userID string,
) ([]string, error) {
	var assignments []models.UserRoleAssignment

	// 查询用户所有活跃的角色分配
	err := r.db.WithContext(ctx).
		Select("role_id").
		Where("user_id = ?", userID).
		Find(&assignments).Error
	if err != nil {
		return nil, err
	}

	// 提取角色ID列表
	roleIDs := make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		roleIDs = append(roleIDs, assignment.RoleID.String())
	}

	return roleIDs, nil
}

// GetActiveRoleIDsWithStatus 获取用户所有处于指定状态的角色ID列表
// 此方法会联表查询 user_role_assignments 和 role_definitions
// 只返回角色状态匹配的角色ID（通常用于获取 Active 状态的角色）
// 专门用于登录等需要验证角色可用性的场景
func (r *UserRoleAssignmentRepositoryImpl) GetActiveRoleIDsWithStatus(
	ctx context.Context,
	userID string,
	status models.RoleStatus,
) ([]string, error) {
	var roleIDs []string

	// 联表查询：user_role_assignments JOIN role_definitions
	// 只返回匹配指定状态的角色ID
	err := r.db.WithContext(ctx).
		Model(&models.UserRoleAssignment{}).
		Select("user_role_assignments.role_id").
		Joins("JOIN role_definitions ON role_definitions.id = user_role_assignments.role_id").
		Where("user_role_assignments.user_id = ?", userID).
		Where("role_definitions.status = ?", status).
		Where("role_definitions.deleted_at IS NULL").
		Pluck("user_role_assignments.role_id", &roleIDs).Error
	if err != nil {
		return nil, err
	}

	return roleIDs, nil
}

// FindWithConditions 根据组合查询条件查询角色分配记录
func (r *UserRoleAssignmentRepositoryImpl) FindWithConditions(
	ctx context.Context,
	conditions *UserRoleAssignmentQueryConditions,
) ([]*models.UserRoleAssignment, *models.PageResult, error) {
	opts := base.NewQueryOptions()

	if conditions != nil {
		if conditions.UserID != nil {
			opts = opts.WithFilter("user_id", *conditions.UserID)
		}

		if conditions.RoleID != nil {
			opts = opts.WithFilter("role_id", *conditions.RoleID)
		}

		if conditions.Page != nil {
			opts = opts.WithPage(conditions.Page.Page, conditions.Page.PageSize).
				WithOrder(conditions.Page.OrderBy, conditions.Page.OrderDesc)
		}
	}

	return r.FindAll(ctx, opts)
}

// GetRolesByUserIDs 批量查询多个用户的角色分配
func (r *UserRoleAssignmentRepositoryImpl) GetRolesByUserIDs(
	ctx context.Context,
	userIDs []string,
) (map[string][]string, error) {
	// 空列表直接返回空map
	if len(userIDs) == 0 {
		return make(map[string][]string), nil
	}

	// 批量查询所有用户的角色分配
	var assignments []models.UserRoleAssignment

	err := r.db.WithContext(ctx).
		Select("user_id, role_id").
		Where("user_id IN ?", userIDs).
		Find(&assignments).Error
	if err != nil {
		return nil, err
	}

	// 转换为 map[userID][]roleID
	result := make(map[string][]string)

	for _, assignment := range assignments {
		userID := assignment.UserID.String()
		roleID := assignment.RoleID.String()
		result[userID] = append(result[userID], roleID)
	}

	return result, nil
}

// ============================================================================
// 业务验证方法
// ============================================================================

// CheckUserRoleExists 检查用户是否已分配指定角色
func (r *UserRoleAssignmentRepositoryImpl) CheckUserRoleExists(
	ctx context.Context,
	userID, roleID string,
) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.UserRoleAssignment{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CountByUserID 统计指定用户的角色分配数量
func (r *UserRoleAssignmentRepositoryImpl) CountByUserID(
	ctx context.Context,
	userID string,
) (int64, error) {
	opts := base.NewQueryOptions().WithFilter("user_id", userID)
	return r.Count(ctx, opts)
}

// CountByRoleID 统计指定角色的分配数量
func (r *UserRoleAssignmentRepositoryImpl) CountByRoleID(
	ctx context.Context,
	roleID string,
) (int64, error) {
	opts := base.NewQueryOptions().WithFilter("role_id", roleID)
	return r.Count(ctx, opts)
}

// ============================================================================
// 批量操作方法
// ============================================================================

// BatchAssignRoleToUsers 批量为多个用户分配同一角色
func (r *UserRoleAssignmentRepositoryImpl) BatchAssignRoleToUsers(
	ctx context.Context,
	userIDs []string,
	roleID string,
	assignedBy string,
) error {
	if len(userIDs) == 0 {
		return nil
	}

	assignments := make([]*models.UserRoleAssignment, 0, len(userIDs))
	for _, userID := range userIDs {
		assignment := &models.UserRoleAssignment{
			UserID: uuid.MustParse(userID),
			RoleID: uuid.MustParse(roleID),
		}
		if assignedBy != "" {
			createdByUUID := uuid.MustParse(assignedBy)
			assignment.CreatedBy = &createdByUUID
		}

		assignments = append(assignments, assignment)
	}

	return r.BatchCreate(ctx, assignments)
}

// BatchRevokeUserRoles 批量撤销用户的多个角色分配
func (r *UserRoleAssignmentRepositoryImpl) BatchRevokeUserRoles(
	ctx context.Context,
	assignmentIDs []string,
) error {
	if len(assignmentIDs) == 0 {
		return nil
	}

	return r.BatchDelete(ctx, assignmentIDs)
}

// GetAllUserIDsByRoleID 获取指定角色下所有用户ID（不分页）
func (r *UserRoleAssignmentRepositoryImpl) GetAllUserIDsByRoleID(
	ctx context.Context,
	roleID string,
) ([]string, error) {
	var assignments []models.UserRoleAssignment

	err := r.db.WithContext(ctx).
		Select("user_id").
		Where("role_id = ?", roleID).
		Find(&assignments).Error
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		userIDs = append(userIDs, assignment.UserID.String())
	}

	return userIDs, nil
}

// ReplaceRoleUsers 批量替换角色的用户绑定（事务操作）
func (r *UserRoleAssignmentRepositoryImpl) ReplaceRoleUsers(
	ctx context.Context,
	roleID string,
	userIDs []string,
	operatorID string,
) error {
	// 使用事务确保数据一致性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 删除该角色下所有旧的用户绑定
		if err := tx.Where("role_id = ?", roleID).Delete(&models.UserRoleAssignment{}).Error; err != nil {
			return err
		}

		// 2. 如果没有新的用户列表，直接返回（清空该角色的所有用户）
		if len(userIDs) == 0 {
			return nil
		}

		// 3. 批量创建新的用户绑定
		assignments := make([]*models.UserRoleAssignment, 0, len(userIDs))
		for _, userID := range userIDs {
			assignment := &models.UserRoleAssignment{
				UserID: uuid.MustParse(userID),
				RoleID: uuid.MustParse(roleID),
			}
			if operatorID != "" {
				createdByUUID := uuid.MustParse(operatorID)
				assignment.CreatedBy = &createdByUUID
			}

			assignments = append(assignments, assignment)
		}

		// 使用 Create 批量插入
		if err := tx.Create(&assignments).Error; err != nil {
			return err
		}

		return nil
	})
}
