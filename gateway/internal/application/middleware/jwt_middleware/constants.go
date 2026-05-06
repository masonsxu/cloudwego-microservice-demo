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

// JWT 验证通过后注入下游的 HTTP Header（业务侧契约，提案 §5.2）
//
// 业务系统（含网关内 authz/access_log 中间件）统一从 header 读身份，
// 禁止再次解析 JWT。
const (
	// HeaderUserID 用户唯一标识（claims.sub）
	HeaderUserID = "X-User-Id"

	// HeaderUserName 用户名（claims.username，展示用）
	HeaderUserName = "X-User-Name"

	// HeaderTenantID 租户标识（claims.tenant，= OrganizationID）
	HeaderTenantID = "X-Tenant-Id"

	// HeaderUserRoles 角色 code 列表（claims.roles），多个值以英文逗号 "," 拼接
	HeaderUserRoles = "X-User-Roles"
)
