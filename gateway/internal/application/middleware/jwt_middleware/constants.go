// Package middleware 提供认证中间件实现
// 基于 github.com/hertz-contrib/jwt 实现 RS256 JWT 认证
package middleware

// JWT claims 中的键名定义（目标态 schema，仅稳定身份字段）
const (
	// IdentityKey 用户唯一标识，映射 OIDC sub
	IdentityKey = "sub"

	// Username 用户名，展示用
	Username = "username"

	// Tenant 租户标识，= OrganizationID（单值）
	Tenant = "tenant"

	// Roles 角色 code 列表（不是 role ID）
	Roles = "roles"
)

// Context中存储登录用户信息的键名
const (
	// LoginUserContextKey 在 Context 中存储登录用户信息的键名
	LoginUserContextKey = "login_user_info"
)
