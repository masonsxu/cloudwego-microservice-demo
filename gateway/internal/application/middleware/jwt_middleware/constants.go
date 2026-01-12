// Package middleware 提供认证中间件实现
// 基于 github.com/hertz-contrib/jwt 实现高性能JWT认证
package middleware

// JWT claims 中的键名定义
const (
	// IdentityKey 表示用户ID (改为使用新的字段名)
	IdentityKey = "userProfileID"

	// OrganizationID 表示组织ID
	OrganizationID = "organizationID"

	// DepartmentIDs 表示部门ID列表（多部门模式）
	DepartmentIDs = "departmentIDs"

	// Username 表示用户名
	Username = "username"

	// Status 表示用户状态
	Status = "status"

	// RoleIDs 表示角色ID列表（多角色模式）
	RoleIDs = "roleIDs"

	// CorePermission 表示核心权限
	CorePermission = "corePermission"

	// DataScope 表示数据范围（self/dept/org）
	DataScope = "dataScope"
)

// Context中存储登录用户信息的键名
const (
	// LoginUserContextKey 在 Context 中存储登录用户信息的键名
	LoginUserContextKey = "login_user_info"
)
