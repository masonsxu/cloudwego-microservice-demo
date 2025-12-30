package rolemenu

import (
	"context"

	"github.com/google/uuid"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// RoleMenuPermissionRepository 角色菜单权限数据访问接口
type RoleMenuPermissionRepository interface {
	// Create 创建角色菜单权限
	Create(ctx context.Context, permission *models.RoleMenuPermission) error

	// BatchCreate 批量创建角色菜单权限
	BatchCreate(ctx context.Context, permissions []*models.RoleMenuPermission) error

	// Update 更新角色菜单权限
	Update(ctx context.Context, permission *models.RoleMenuPermission) error

	// Delete 删除角色菜单权限
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByRoleID 删除指定角色的所有菜单权限
	DeleteByRoleID(ctx context.Context, roleID uuid.UUID) error

	// GetByID 根据ID获取角色菜单权限
	GetByID(ctx context.Context, id uuid.UUID) (*models.RoleMenuPermission, error)

	// GetByRoleID 获取指定角色的所有菜单权限
	GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]*models.RoleMenuPermission, error)

	// GetByRoleIDs 批量获取多个角色的菜单权限
	GetByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]*models.RoleMenuPermission, error)

	// GetByRoleAndMenu 获取指定角色和菜单的权限
	GetByRoleAndMenu(ctx context.Context, roleID uuid.UUID, menuID string) (*models.RoleMenuPermission, error)

	// SyncRoleMenus 同步角色菜单权限（先删除旧的，再创建新的）
	SyncRoleMenus(ctx context.Context, roleID uuid.UUID, permissions []*models.RoleMenuPermission) error

	// GetMenuIDsByRoleID 获取指定角色有权限的菜单ID列表
	GetMenuIDsByRoleID(ctx context.Context, roleID uuid.UUID) ([]string, error)

	// GetMenuIDsByRoleIDs 批量获取多个角色有权限的菜单ID列表（去重）
	GetMenuIDsByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]string, error)

	// HasPermission 检查角色是否具有指定菜单的指定权限
	HasPermission(ctx context.Context, roleID uuid.UUID, menuID string, permissionType models.MenuPermissionType) (bool, error)

	// GetMergedPermissions 获取多个角色的合并权限（每个菜单取最高权限）
	GetMergedPermissions(ctx context.Context, roleIDs []uuid.UUID) ([]models.MenuPermissionInfo, error)
}
