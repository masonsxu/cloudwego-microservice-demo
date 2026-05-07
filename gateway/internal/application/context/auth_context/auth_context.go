package auth_context

import (
	"github.com/cloudwego/hertz/pkg/app"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
)

// AuthContext 统一的认证上下文管理器（目标态：仅稳定身份字段）
//
// 字段对照提案 §5.1 token claims 最小集：sub / username / tenant / roles。
// 任何业务字段（部门、数据范围、权限码等）一律不再持有，决策应通过 PDP
// （policy_srv）在业务系统侧完成。
type AuthContext struct {
	claims *http_base.JWTClaimsDTO
}

// AuthContextKey 认证上下文在 context 中的键
const AuthContextKey = "auth_context"

// NewAuthContext 创建新的认证上下文
func NewAuthContext(claims *http_base.JWTClaimsDTO) *AuthContext {
	return &AuthContext{claims: claims}
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

// GetUserProfileID 获取用户ID（sub）
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

// GetTenant 获取租户ID（Tenant claim）
func (ac *AuthContext) GetTenant() (string, bool) {
	if ac == nil || ac.claims == nil || ac.claims.OrganizationID == nil {
		return "", false
	}

	return *ac.claims.OrganizationID, true
}

// GetRoles 获取角色 code 列表（Roles claim）
func (ac *AuthContext) GetRoles() []string {
	if ac == nil || ac.claims == nil {
		return nil
	}

	return ac.claims.RoleIDs
}

// GetOrganizationID 获取组织ID（兼容别名，= Tenant）
func (ac *AuthContext) GetOrganizationID() (string, bool) {
	return ac.GetTenant()
}

// ---- 便利函数：直接从 RequestContext 获取认证信息 ----

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

// GetCurrentTenant 直接从请求上下文获取当前租户ID
func GetCurrentTenant(c *app.RequestContext) (string, bool) {
	if authCtx, exists := GetAuthContext(c); exists {
		return authCtx.GetTenant()
	}

	return "", false
}
