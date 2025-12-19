package auth_context

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/core"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
)

// AuthContext 统一的认证上下文管理器
// 提供用户认证信息的存储和提取功能，支持角色和权限管理
type AuthContext struct {
	claims      *http_base.JWTClaimsDTO
	roles       []string // 用户当前拥有的角色列表
	permissions []string // 用户当前拥有的权限列表
}

// AuthContextKey 认证上下文在 context 中的键
const AuthContextKey = "auth_context"

// NewAuthContext 创建新的认证上下文
func NewAuthContext(claims *http_base.JWTClaimsDTO) *AuthContext {
	return &AuthContext{
		claims:      claims,
		roles:       make([]string, 0),
		permissions: make([]string, 0),
	}
}

// NewAuthContextWithRoles 创建带有角色信息的认证上下文
func NewAuthContextWithRoles(
	claims *http_base.JWTClaimsDTO,
	roles []string,
	permissions []string,
) *AuthContext {
	return &AuthContext{
		claims:      claims,
		roles:       roles,
		permissions: permissions,
	}
}

// SetAuthContext 将认证上下文设置到请求上下文中
func SetAuthContext(c *app.RequestContext, authCtx *AuthContext) {
	c.Set(AuthContextKey, authCtx)
}

// GetAuthContext 从请求上下文中获取认证上下文
func GetAuthContext(c *app.RequestContext) (*AuthContext, bool) {
	if value, exists := c.Get(AuthContextKey); exists {
		if authCtx, ok := value.(*AuthContext); ok {
			return authCtx, true
		}
	}

	return nil, false
}

// GetUserProfileID 获取用户ID
func (ac *AuthContext) GetUserProfileID() (string, bool) {
	if ac == nil || ac.claims == nil || ac.claims.UserProfileID == nil {
		return "", false
	}

	return *ac.claims.UserProfileID, true
}

// GetUsername 获取用户名
func (ac *AuthContext) GetUsername() (string, bool) {
	if ac == nil || ac.claims == nil || ac.claims.Username == nil {
		return "", false
	}

	return *ac.claims.Username, true
}

// GetOrganizationID 获取组织ID
func (ac *AuthContext) GetOrganizationID() (string, bool) {
	if ac == nil || ac.claims == nil || ac.claims.OrganizationID == nil {
		return "", false
	}

	return *ac.claims.OrganizationID, true
}

// GetUserStatus 获取用户状态
func (ac *AuthContext) GetUserStatus() (core.UserStatus, bool) {
	if ac == nil || ac.claims == nil || ac.claims.Status == nil {
		return 0, false
	}

	return *ac.claims.Status, true
}

// 注意：AccountType 字段在当前 JWT Claims 中不存在，使用角色系统替代
// GetRoleID 获取角色ID（如果需要类似AccountType的功能，可以通过角色系统实现）
func (ac *AuthContext) GetRoleID() (string, bool) {
	if ac == nil || ac.claims == nil || ac.claims.RoleID == nil {
		return "", false
	}

	return *ac.claims.RoleID, true
}

// GetDepartmentID 获取部门ID
func (ac *AuthContext) GetDepartmentID() (string, bool) {
	if ac == nil || ac.claims == nil || ac.claims.DepartmentID == nil {
		return "", false
	}

	return *ac.claims.DepartmentID, true
}

// GetPermission 获取权限
func (ac *AuthContext) GetPermission() (string, bool) {
	if ac == nil || ac.claims == nil || ac.claims.Permission == nil {
		return "", false
	}

	return *ac.claims.Permission, true
}

// 便利函数：直接从 RequestContext 获取认证信息

// GetCurrentUserProfileID 直接从请求上下文获取当前用户ID
func GetCurrentUserProfileID(c *app.RequestContext) (string, bool) {
	if authCtx, exists := GetAuthContext(c); exists {
		return authCtx.GetUserProfileID()
	}

	return "", false
}

// GetCurrentUsername 直接从请求上下文获取当前用户名
func GetCurrentUsername(c *app.RequestContext) (string, bool) {
	if authCtx, exists := GetAuthContext(c); exists {
		return authCtx.GetUsername()
	}

	return "", false
}

// GetCurrentOrganizationID 直接从请求上下文获取当前组织ID
func GetCurrentOrganizationID(c *app.RequestContext) (string, bool) {
	if authCtx, exists := GetAuthContext(c); exists {
		return authCtx.GetOrganizationID()
	}

	return "", false
}

// GetCurrentUserStatus 直接从请求上下文获取当前用户状态
func GetCurrentUserStatus(c *app.RequestContext) (core.UserStatus, bool) {
	if authCtx, exists := GetAuthContext(c); exists {
		return authCtx.GetUserStatus()
	}

	return 0, false
}

// GetCurrentRoleID 直接从请求上下文获取当前角色ID
func GetCurrentRoleID(c *app.RequestContext) (string, bool) {
	if authCtx, exists := GetAuthContext(c); exists {
		return authCtx.GetRoleID()
	}

	return "", false
}

// GetCurrentPermission 直接从请求上下文获取当前权限
func GetCurrentPermission(c *app.RequestContext) (string, bool) {
	if authCtx, exists := GetAuthContext(c); exists {
		return authCtx.GetPermission()
	}

	return "", false
}
