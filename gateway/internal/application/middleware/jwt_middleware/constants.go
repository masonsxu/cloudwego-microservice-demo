// Package middleware 提供认证中间件实现
// 基于 github.com/hertz-contrib/jwt 实现高性能JWT认证
package middleware

// JWT claims 中的键名定义
const (
	// IdentityKey 表示用户ID (改为使用新的字段名)
	IdentityKey = "userProfileID"

	// OrganizationID 表示组织ID
	OrganizationID = "organizationID"

	// DepartmentID 表示部门ID
	DepartmentID = "departmentID"

	// Username 表示用户名
	Username = "username"

	// Status 表示用户状态
	Status = "status"

	// RoleID 表示角色ID（简化为单角色）
	RoleID = "roleID"

	// CorePermission 表示核心权限
	CorePermission = "corePermission"
)

// Context中存储登录用户信息的键名
const (
	// LoginUserContextKey 在 Context 中存储登录用户信息的键名
	LoginUserContextKey = "login_user_info"
)
