package casbin

import (
	"context"

	"github.com/casbin/casbin/v3"
)

// PolicySyncRepository Casbin 策略同步仓储接口
type PolicySyncRepository interface {
	// GetEnforcer 获取 Casbin Enforcer 实例
	GetEnforcer() *casbin.SyncedEnforcer

	// SyncRolePolicies 同步角色权限策略
	// 从 role_menu_permissions 表生成 Casbin 策略
	SyncRolePolicies(ctx context.Context) error

	// SyncUserRoleBindings 同步用户角色绑定
	// 从 user_role_assignments 表生成 Casbin g 规则
	SyncUserRoleBindings(ctx context.Context) error

	// SyncRoleInheritance 同步角色继承关系
	// 从 role_definitions.parent_role_id 生成 Casbin g2 规则
	SyncRoleInheritance(ctx context.Context) error

	// SyncAll 同步所有策略数据
	SyncAll(ctx context.Context) error

	// AddRolePolicy 添加角色策略
	AddRolePolicy(ctx context.Context, roleCode, domain, resource, action, dataScope string) error

	// RemoveRolePolicy 移除角色策略
	RemoveRolePolicy(ctx context.Context, roleCode, domain, resource, action, dataScope string) error

	// AddUserRoleBinding 添加用户角色绑定
	AddUserRoleBinding(ctx context.Context, userID, roleCode, domain string) error

	// RemoveUserRoleBinding 移除用户角色绑定
	RemoveUserRoleBinding(ctx context.Context, userID, roleCode, domain string) error

	// AddRoleInheritance 添加角色继承
	AddRoleInheritance(ctx context.Context, childRole, parentRole string) error

	// RemoveRoleInheritance 移除角色继承
	RemoveRoleInheritance(ctx context.Context, childRole, parentRole string) error

	// CheckPermission 检查权限
	CheckPermission(ctx context.Context, userID, domain, resource, action string) (bool, error)

	// GetUserDataScope 获取用户在指定资源上的数据范围
	GetUserDataScope(ctx context.Context, userID, domain, resource, action string) (string, error)
}
