package casbin_middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// CasbinMiddlewareService Casbin权限中间件服务接口
type CasbinMiddlewareService interface {
	// MiddlewareFunc 返回权限校验中间件
	MiddlewareFunc() app.HandlerFunc

	// CheckPermission 检查用户是否有权限访问指定资源
	// userID: 用户ID
	// roleIDs: 用户角色ID列表（多角色模式）
	// deptIDs: 用户科室ID列表
	// resource: 资源对象（菜单ID或API路径）
	// action: 操作类型（read/write/manage/full）
	CheckPermission(ctx context.Context, userID string, roleIDs []string, deptIDs []string, resource, action string) (bool, error)

	// CheckPermissionWithDataScope 检查权限并返回数据范围
	// 返回: (是否有权限, 数据范围, 错误)
	CheckPermissionWithDataScope(ctx context.Context, userID string, roleIDs []string, deptIDs []string, resource, action string) (bool, string, error)

	// GetUserPermissions 获取用户的所有权限列表
	GetUserPermissions(ctx context.Context, userID string, roleIDs []string) ([]PermissionInfo, error)

	// RefreshUserPermissions 刷新用户权限缓存
	RefreshUserPermissions(ctx context.Context, userID string) error

	// RefreshAllPermissions 刷新所有权限策略（从数据库重新加载）
	RefreshAllPermissions(ctx context.Context) error
}

// PermissionInfo 权限信息
type PermissionInfo struct {
	Resource  string `json:"resource"`   // 资源对象
	Action    string `json:"action"`     // 操作类型
	DataScope string `json:"data_scope"` // 数据范围
	Domain    string `json:"domain"`     // 域（科室）
}

// EnforcerService Casbin Enforcer 服务接口
type EnforcerService interface {
	// Enforce 执行权限校验
	Enforce(sub, dom, obj, act string) (bool, error)

	// EnforceWithDataScope 执行权限校验并返回数据范围
	EnforceWithDataScope(sub, dom, obj, act string) (bool, string, error)

	// AddPolicy 添加策略
	AddPolicy(sub, dom, obj, act, dataScope string) (bool, error)

	// RemovePolicy 移除策略
	RemovePolicy(sub, dom, obj, act, dataScope string) (bool, error)

	// AddRoleForUserInDomain 添加用户-角色-域绑定
	AddRoleForUserInDomain(user, role, domain string) (bool, error)

	// RemoveRoleForUserInDomain 移除用户-角色-域绑定
	RemoveRoleForUserInDomain(user, role, domain string) (bool, error)

	// AddRoleInheritance 添加角色继承关系
	AddRoleInheritance(childRole, parentRole string) (bool, error)

	// RemoveRoleInheritance 移除角色继承关系
	RemoveRoleInheritance(childRole, parentRole string) (bool, error)

	// GetRolesForUserInDomain 获取用户在指定域下的角色
	GetRolesForUserInDomain(user, domain string) []string

	// GetPermissionsForUser 获取用户的所有权限
	GetPermissionsForUser(user string) [][]string

	// LoadPolicy 从数据库加载策略
	LoadPolicy() error

	// SavePolicy 保存策略到数据库
	SavePolicy() error
}
