package rolemenu

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// roleMenuPermissionRepository 实现 RoleMenuPermissionRepository 接口
type roleMenuPermissionRepository struct {
	db *gorm.DB
}

// NewRoleMenuPermissionRepository 创建角色菜单权限仓储实例
func NewRoleMenuPermissionRepository(db *gorm.DB) RoleMenuPermissionRepository {
	return &roleMenuPermissionRepository{db: db}
}

// Create 创建角色菜单权限
func (r *roleMenuPermissionRepository) Create(ctx context.Context, permission *models.RoleMenuPermission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

// BatchCreate 批量创建角色菜单权限
func (r *roleMenuPermissionRepository) BatchCreate(
	ctx context.Context,
	permissions []*models.RoleMenuPermission,
) error {
	if len(permissions) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).CreateInBatches(permissions, 100).Error
}

// Update 更新角色菜单权限
func (r *roleMenuPermissionRepository) Update(ctx context.Context, permission *models.RoleMenuPermission) error {
	return r.db.WithContext(ctx).Save(permission).Error
}

// Delete 删除角色菜单权限
func (r *roleMenuPermissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.RoleMenuPermission{}, "id = ?", id).Error
}

// DeleteByRoleID 删除指定角色的所有菜单权限
func (r *roleMenuPermissionRepository) DeleteByRoleID(ctx context.Context, roleID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.RoleMenuPermission{}, "role_id = ?", roleID).Error
}

// GetByID 根据ID获取角色菜单权限
func (r *roleMenuPermissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.RoleMenuPermission, error) {
	var permission models.RoleMenuPermission

	err := r.db.WithContext(ctx).First(&permission, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &permission, nil
}

// GetByRoleID 获取指定角色的所有菜单权限
func (r *roleMenuPermissionRepository) GetByRoleID(
	ctx context.Context,
	roleID uuid.UUID,
) ([]*models.RoleMenuPermission, error) {
	var permissions []*models.RoleMenuPermission

	err := r.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetByRoleIDs 批量获取多个角色的菜单权限
func (r *roleMenuPermissionRepository) GetByRoleIDs(
	ctx context.Context,
	roleIDs []uuid.UUID,
) ([]*models.RoleMenuPermission, error) {
	if len(roleIDs) == 0 {
		return []*models.RoleMenuPermission{}, nil
	}

	var permissions []*models.RoleMenuPermission

	err := r.db.WithContext(ctx).
		Where("role_id IN ?", roleIDs).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetByRoleAndMenu 获取指定角色和菜单的权限
func (r *roleMenuPermissionRepository) GetByRoleAndMenu(
	ctx context.Context,
	roleID uuid.UUID,
	menuID string,
) (*models.RoleMenuPermission, error) {
	var permission models.RoleMenuPermission

	err := r.db.WithContext(ctx).
		Where("role_id = ? AND menu_id = ?", roleID, menuID).
		First(&permission).Error
	if err != nil {
		return nil, err
	}

	return &permission, nil
}

// SyncRoleMenus 同步角色菜单权限（先删除旧的，再创建新的）
func (r *roleMenuPermissionRepository) SyncRoleMenus(
	ctx context.Context,
	roleID uuid.UUID,
	permissions []*models.RoleMenuPermission,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 删除角色的所有旧权限（使用硬删除，避免唯一索引冲突）
		if err := tx.Unscoped().Delete(&models.RoleMenuPermission{}, "role_id = ?", roleID).Error; err != nil {
			return err
		}

		// 2. 创建新权限
		if len(permissions) > 0 {
			// 确保所有权限都属于该角色
			for _, perm := range permissions {
				perm.RoleID = roleID
			}

			if err := tx.CreateInBatches(permissions, 100).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetMenuIDsByRoleID 获取指定角色有权限的菜单ID列表
func (r *roleMenuPermissionRepository) GetMenuIDsByRoleID(ctx context.Context, roleID uuid.UUID) ([]string, error) {
	var menuIDs []string

	err := r.db.WithContext(ctx).
		Model(&models.RoleMenuPermission{}).
		Where("role_id = ?", roleID).
		Pluck("menu_id", &menuIDs).Error
	if err != nil {
		return nil, err
	}

	return menuIDs, nil
}

// GetMenuIDsByRoleIDs 批量获取多个角色有权限的菜单ID列表（去重）
func (r *roleMenuPermissionRepository) GetMenuIDsByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]string, error) {
	if len(roleIDs) == 0 {
		return []string{}, nil
	}

	var menuIDs []string

	err := r.db.WithContext(ctx).
		Model(&models.RoleMenuPermission{}).
		Where("role_id IN ?", roleIDs).
		Distinct("menu_id").
		Pluck("menu_id", &menuIDs).Error
	if err != nil {
		return nil, err
	}

	return menuIDs, nil
}

// HasPermission 检查角色是否具有指定菜单的指定权限
func (r *roleMenuPermissionRepository) HasPermission(
	ctx context.Context,
	roleID uuid.UUID,
	menuID string,
	permissionType models.MenuPermissionType,
) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.RoleMenuPermission{}).
		Where("role_id = ? AND menu_id = ? AND permission_type >= ?", roleID, menuID, permissionType).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetMergedPermissions 获取多个角色的合并权限（每个菜单取最高权限）
func (r *roleMenuPermissionRepository) GetMergedPermissions(
	ctx context.Context,
	roleIDs []uuid.UUID,
) ([]models.MenuPermissionInfo, error) {
	if len(roleIDs) == 0 {
		return []models.MenuPermissionInfo{}, nil
	}

	// 查询所有权限
	permissions, err := r.GetByRoleIDs(ctx, roleIDs)
	if err != nil {
		return nil, err
	}

	// 按菜单ID分组并合并权限
	permMap := make(map[string]*models.MenuPermissionInfo)
	for _, perm := range permissions {
		if existing, ok := permMap[perm.MenuID]; ok {
			existing.PermissionType = models.MergePermissionTypes(existing.PermissionType, perm.PermissionType)
			existing.DataScope = models.MergeDataScopes(existing.DataScope, perm.DataScope)
		} else {
			permMap[perm.MenuID] = &models.MenuPermissionInfo{
				MenuID:         perm.MenuID,
				PermissionType: perm.PermissionType,
				DataScope:      perm.DataScope,
			}
		}
	}

	// 转换为切片
	result := make([]models.MenuPermissionInfo, 0, len(permMap))
	for _, perm := range permMap {
		result = append(result, *perm)
	}

	return result, nil
}
