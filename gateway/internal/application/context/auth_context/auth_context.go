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

	// 多角色模式扩展字段
	roleIDs       []string // 用户的角色ID列表
	departmentIDs []string // 用户的部门ID列表
	dataScope     string   // 当前请求的数据范围（由 Casbin 中间件设置）
}

// AuthContextKey 认证上下文在 context 中的键
const AuthContextKey = "auth_context"

// DataScopeKey 数据范围在 context 中的键
const DataScopeKey = "data_scope"

// NewAuthContext 创建新的认证上下文
// 自动从 JWT claims 中提取多角色和多部门信息
func NewAuthContext(claims *http_base.JWTClaimsDTO) *AuthContext {
	ctx := &AuthContext{
		claims:      claims,
		roles:       make([]string, 0),
		permissions: make([]string, 0),
	}

	// 从 claims 中提取多角色信息
	if claims != nil {
		if len(claims.RoleIDs) > 0 {
			ctx.roleIDs = claims.RoleIDs
		}
		if len(claims.DepartmentIDs) > 0 {
			ctx.departmentIDs = claims.DepartmentIDs
		}
		if claims.DataScope != nil {
			ctx.dataScope = *claims.DataScope
		}
	}

	return ctx
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

// NewAuthContextWithMultiRoles 创建带有多角色信息的认证上下文（多角色模式）
func NewAuthContextWithMultiRoles(
	claims *http_base.JWTClaimsDTO,
	roles []string,
	permissions []string,
	roleIDs []string,
	departmentIDs []string,
) *AuthContext {
	return &AuthContext{
		claims:        claims,
		roles:         roles,
		permissions:   permissions,
		roleIDs:       roleIDs,
		departmentIDs: departmentIDs,
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
// GetRoleID 获取角色ID（兼容单角色模式，返回第一个角色）
// 推荐使用 GetRoleIDs() 获取所有角色
func (ac *AuthContext) GetRoleID() (string, bool) {
	if ac == nil {
		return "", false
	}
	// 优先从多角色列表获取
	roleIDs := ac.GetRoleIDs()
	if len(roleIDs) > 0 {
		return roleIDs[0], true
	}
	return "", false
}

// GetDepartmentID 获取部门ID（兼容单部门模式，返回第一个部门）
// 推荐使用 GetDepartmentIDs() 获取所有部门
func (ac *AuthContext) GetDepartmentID() (string, bool) {
	if ac == nil {
		return "", false
	}
	// 优先从多部门列表获取
	deptIDs := ac.GetDepartmentIDs()
	if len(deptIDs) > 0 {
		return deptIDs[0], true
	}
	return "", false
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

// ========== 多角色模式扩展方法 ==========

// GetRoleIDs 获取角色ID列表（多角色模式）
func (ac *AuthContext) GetRoleIDs() []string {
	if ac == nil {
		return nil
	}
	// 如果设置了 roleIDs，优先返回
	if len(ac.roleIDs) > 0 {
		return ac.roleIDs
	}
	// 从 claims 中获取角色ID列表
	if ac.claims != nil && len(ac.claims.RoleIDs) > 0 {
		return ac.claims.RoleIDs
	}
	return nil
}

// SetRoleIDs 设置角色ID列表
func (ac *AuthContext) SetRoleIDs(roleIDs []string) {
	if ac != nil {
		ac.roleIDs = roleIDs
	}
}

// GetDepartmentIDs 获取部门ID列表（多部门模式）
func (ac *AuthContext) GetDepartmentIDs() []string {
	if ac == nil {
		return nil
	}
	// 如果设置了 departmentIDs，优先返回
	if len(ac.departmentIDs) > 0 {
		return ac.departmentIDs
	}
	// 从 claims 中获取部门ID列表
	if ac.claims != nil && len(ac.claims.DepartmentIDs) > 0 {
		return ac.claims.DepartmentIDs
	}
	return nil
}

// SetDepartmentIDs 设置部门ID列表
func (ac *AuthContext) SetDepartmentIDs(departmentIDs []string) {
	if ac != nil {
		ac.departmentIDs = departmentIDs
	}
}

// GetDataScope 获取数据范围
func (ac *AuthContext) GetDataScope() string {
	if ac == nil {
		return ""
	}
	return ac.dataScope
}

// SetDataScopeInternal 设置数据范围（内部方法）
func (ac *AuthContext) SetDataScopeInternal(dataScope string) {
	if ac != nil {
		ac.dataScope = dataScope
	}
}

// GetCurrentRoleIDs 直接从请求上下文获取当前角色ID列表
func GetCurrentRoleIDs(c *app.RequestContext) []string {
	if authCtx, exists := GetAuthContext(c); exists {
		return authCtx.GetRoleIDs()
	}
	return nil
}

// GetCurrentDepartmentIDs 直接从请求上下文获取当前部门ID列表
func GetCurrentDepartmentIDs(c *app.RequestContext) []string {
	if authCtx, exists := GetAuthContext(c); exists {
		return authCtx.GetDepartmentIDs()
	}
	return nil
}

// GetCurrentDataScope 直接从请求上下文获取当前数据范围
func GetCurrentDataScope(c *app.RequestContext) string {
	// 优先从 AuthContext 获取
	if authCtx, exists := GetAuthContext(c); exists {
		if dataScope := authCtx.GetDataScope(); dataScope != "" {
			return dataScope
		}
	}
	// 其次从 RequestContext 直接获取
	if value, exists := c.Get(DataScopeKey); exists {
		if dataScope, ok := value.(string); ok {
			return dataScope
		}
	}
	return ""
}

// SetDataScope 设置数据范围到请求上下文
func SetDataScope(c *app.RequestContext, dataScope string) {
	c.Set(DataScopeKey, dataScope)
	// 同时更新 AuthContext
	if authCtx, exists := GetAuthContext(c); exists {
		authCtx.SetDataScopeInternal(dataScope)
	}
}

// SetRoleIDsToContext 设置角色ID列表到请求上下文
func SetRoleIDsToContext(c *app.RequestContext, roleIDs []string) {
	if authCtx, exists := GetAuthContext(c); exists {
		authCtx.SetRoleIDs(roleIDs)
	}
}

// SetDepartmentIDsToContext 设置部门ID列表到请求上下文
func SetDepartmentIDsToContext(c *app.RequestContext, departmentIDs []string) {
	if authCtx, exists := GetAuthContext(c); exists {
		authCtx.SetDepartmentIDs(departmentIDs)
	}
}
