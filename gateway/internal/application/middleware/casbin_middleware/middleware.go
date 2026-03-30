package casbin_middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/rs/zerolog"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/context/auth_context"
	tracelog "github.com/masonsxu/cloudwego-microservice-demo/gateway/pkg/log"
)

// CasbinMiddleware Casbin权限中间件实现
type CasbinMiddleware struct {
	enforcer                  *CasbinEnforcer
	logger                    *zerolog.Logger
	skipPaths                 []string
	pathMapping               map[string]string // API路径到权限编码的映射
	superAdminBypassEnabled   bool
	superAdminSubjectAllowSet map[string]struct{}
}

// NewCasbinMiddleware 创建新的 Casbin 中间件
func NewCasbinMiddleware(enforcer *CasbinEnforcer, logger *zerolog.Logger) *CasbinMiddleware {
	return &CasbinMiddleware{
		enforcer:                  enforcer,
		logger:                    logger,
		skipPaths:                 defaultSkipPaths(),
		pathMapping:               defaultPathMapping(),
		superAdminBypassEnabled:   true,
		superAdminSubjectAllowSet: toLookupSet(defaultSuperAdminSubjects()),
	}
}

// defaultSkipPaths 默认跳过权限检查的路径
func defaultSkipPaths() []string {
	return []string{
		"/health",
		"/metrics",
		"/swagger",
		"/api/v1/identity/auth/login",
	}
}

// defaultPathMapping 默认 API 路径到 Casbin 资源映射
func defaultPathMapping() map[string]string {
	return map[string]string{
		"/api/v1/permission/roles":                    "menu:role_permissions",
		"/api/v1/permission/roles/*":                  "menu:role_permissions",
		"/api/v1/identity/users":                      "menu:account_management",
		"/api/v1/identity/users/*":                    "menu:account_management",
		"/api/v1/identity/organizations":              "menu:organization_management",
		"/api/v1/identity/organizations/*":            "menu:organization_management",
		"/api/v1/identity/departments":                "menu:organization_management",
		"/api/v1/identity/departments/*":              "menu:organization_management",
		"/api/v1/identity/organization-logos":         "menu:organization_management",
		"/api/v1/identity/organization-logos/*":       "menu:organization_management",
		"/api/v1/identity/users/*/memberships":        "menu:organization_management",
		"/api/v1/identity/users/*/primary-membership": "menu:organization_management",
		"/api/v1/identity/audit-logs":                 "menu:audit_logs",
		"/api/v1/identity/audit-logs/*":               "menu:audit_logs",
	}
}

func defaultSuperAdminSubjects() []string {
	return []string{
		"role:superadmin",
		"username:superadmin",
	}
}

// SetSkipPaths 设置跳过权限检查的路径
func (m *CasbinMiddleware) SetSkipPaths(paths []string) {
	m.skipPaths = paths
}

// AddPathMapping 添加路径到权限编码的映射
func (m *CasbinMiddleware) AddPathMapping(path, permCode string) {
	m.pathMapping[path] = permCode
}

// SetSuperAdminBypassConfig 设置超级管理员兜底放行配置
func (m *CasbinMiddleware) SetSuperAdminBypassConfig(enabled bool, subjects []string) {
	m.superAdminBypassEnabled = enabled

	if len(subjects) == 0 {
		subjects = defaultSuperAdminSubjects()
	}

	m.superAdminSubjectAllowSet = toLookupSet(subjects)
}

// MiddlewareFunc 返回权限校验中间件
func (m *CasbinMiddleware) MiddlewareFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Request.URI().Path())
		method := string(c.Request.Method())

		// 检查是否跳过权限检查
		shouldSkip := m.shouldSkip(path)
		tracelog.Event(ctx, m.logger.Debug()).
			Str("component", "casbin_middleware").
			Str("path", path).
			Strs("skip_paths", m.skipPaths).
			Bool("should_skip", shouldSkip).
			Msg("Casbin skip decision")
		if shouldSkip {
			c.Next(ctx)
			return
		}

		// 如果 enforcer 为空，跳过权限检查
		if m.enforcer == nil {
			c.Next(ctx)
			return
		}

		// 获取认证上下文
		authCtx, exists := auth_context.GetAuthContext(c)
		if !exists || authCtx == nil {
			tracelog.Event(ctx, m.logger.Warn()).
				Str("path", path).
				Str("method", method).
				Msg("No auth context found")
			abortWithUnauthorized(c, "Authentication required")
			return
		}

		tracelog.Event(ctx, m.logger.Debug()).
			Str("path", path).
			Str("method", method).
			Msg("Auth context found")

		// 获取用户信息
		userID, ok := auth_context.GetCurrentUserProfileID(c)
		if !ok {
			tracelog.Event(ctx, m.logger.Warn()).
				Str("path", path).
				Str("method", method).
				Msg("No user ID found in auth context")
			abortWithUnauthorized(c, "User not authenticated")
			return
		}

		// 获取部门ID列表
		deptIDs := auth_context.GetCurrentDepartmentIDs(c)
		roleIDs := auth_context.GetCurrentRoleIDs(c)
		username, _ := auth_context.GetCurrentUsername(c)
		tracelog.Event(ctx, m.logger.Debug()).
			Str("path", path).
			Str("method", method).
			Str("user_id", userID).
			Str("username", username).
			Strs("role_ids", roleIDs).
			Strs("department_ids", deptIDs).
			Msg("Casbin auth context extracted")

		if m.shouldBypassForSuperAdmin(ctx, userID, username, roleIDs, deptIDs) {
			auth_context.SetDataScope(c, "org")
			c.Next(ctx)
			return
		}

		// 获取角色ID列表（多角色模式）
		if len(roleIDs) == 0 {
			tracelog.Event(ctx, m.logger.Warn()).
				Str("path", path).
				Str("method", method).
				Str("user_id", userID).
				Msg("No roles found for user")
			abortWithPermissionDenied(c, "No roles assigned")
			return
		}

		// 将HTTP方法转换为操作类型
		action := methodToAction(method)

		// 获取权限资源（优先使用路径映射，其次使用API路径）
		resource := m.getResource(path)

		tracelog.Event(ctx, m.logger.Debug()).
			Str("path", path).
			Str("method", method).
			Str("user_id", userID).
			Strs("role_ids", roleIDs).
			Strs("department_ids", deptIDs).
			Str("resource", resource).
			Str("action", action).
			Msg("Starting Casbin permission evaluation")

		// 执行权限检查（多角色取并集）
		allowed, dataScope, err := m.checkMultiRolePermission(ctx, userID, roleIDs, deptIDs, resource, action)
		if err != nil {
			tracelog.Event(ctx, m.logger.Error()).Err(err).
				Str("path", path).
				Str("user_id", userID).
				Strs("role_ids", roleIDs).
				Msg("Permission check error")
			abortWithInternalError(c, "Permission check failed")
			return
		}

		if !allowed {
			tracelog.Event(ctx, m.logger.Info()).
				Str("path", path).
				Str("user_id", userID).
				Strs("role_ids", roleIDs).
				Str("resource", resource).
				Str("action", action).
				Msg("Permission denied")
			abortWithPermissionDenied(c, "Access denied to resource")
			return
		}

		// 将数据范围注入到上下文中
		if dataScope != "" {
			auth_context.SetDataScope(c, dataScope)
		}

		tracelog.Event(ctx, m.logger.Debug()).
			Str("path", path).
			Str("user_id", userID).
			Str("resource", resource).
			Str("action", action).
			Str("data_scope", dataScope).
			Msg("Permission granted")

		c.Next(ctx)
	}
}

// abortWithPermissionDenied 权限拒绝响应
func abortWithPermissionDenied(c *app.RequestContext, message string) {
	c.AbortWithStatusJSON(403, map[string]interface{}{
		"code":    403,
		"message": message,
	})
}

// abortWithInternalError 内部错误响应
func abortWithInternalError(c *app.RequestContext, message string) {
	c.AbortWithStatusJSON(500, map[string]interface{}{
		"code":    500,
		"message": message,
	})
}

// abortWithUnauthorized 未授权响应
func abortWithUnauthorized(c *app.RequestContext, message string) {
	c.AbortWithStatusJSON(401, map[string]interface{}{
		"code":    401,
		"message": message,
	})
}

// shouldSkip 检查是否应该跳过权限检查
func (m *CasbinMiddleware) shouldSkip(path string) bool {
	for _, skipPath := range m.skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// getResource 获取权限资源标识
func (m *CasbinMiddleware) getResource(path string) string {
	// 优先使用路径映射
	if permCode, ok := m.pathMapping[path]; ok {
		return permCode
	}

	// 尝试模式匹配
	for pattern, permCode := range m.pathMapping {
		if matchPath(pattern, path) {
			return permCode
		}
	}

	// 默认使用API路径
	return path
}

// matchPath 简单的路径模式匹配
func matchPath(pattern, path string) bool {
	// 支持简单的通配符匹配，如 /api/v1/users/* 匹配 /api/v1/users/123
	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		return strings.HasPrefix(path, prefix)
	}
	return pattern == path
}

// methodToAction 将HTTP方法转换为操作类型
func methodToAction(method string) string {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return "read"
	case http.MethodPost:
		return "write"
	case http.MethodPut, http.MethodPatch:
		return "write"
	case http.MethodDelete:
		return "manage"
	default:
		return "read"
	}
}

func toLookupSet(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}

		result[trimmed] = struct{}{}
	}

	return result
}

func (m *CasbinMiddleware) shouldBypassForSuperAdmin(
	ctx context.Context,
	userID string,
	username string,
	roleIDs []string,
	deptIDs []string,
) bool {
	if !m.superAdminBypassEnabled {
		return false
	}

	subjects := m.buildBypassSubjects(userID, username, roleIDs, deptIDs)
	for subject := range subjects {
		if _, ok := m.superAdminSubjectAllowSet[subject]; ok {
			tracelog.Event(ctx, m.logger.Debug()).
				Str("user_id", userID).
				Str("username", username).
				Str("subject", subject).
				Msg("Casbin bypass: superadmin subject matched")
			return true
		}
	}

	return false
}

func (m *CasbinMiddleware) buildBypassSubjects(
	userID string,
	username string,
	roleIDs []string,
	deptIDs []string,
) map[string]struct{} {
	subjects := m.buildEffectiveRoleSubjects(userID, roleIDs, deptIDs)

	if userID != "" {
		subjects["user:"+userID] = struct{}{}
	}

	trimmedName := strings.TrimSpace(username)
	if trimmedName != "" {
		subjects["username:"+trimmedName] = struct{}{}
		subjects["username:"+strings.ToLower(trimmedName)] = struct{}{}
	}

	return subjects
}

func (m *CasbinMiddleware) buildEffectiveRoleSubjects(
	userID string,
	roleIDs []string,
	deptIDs []string,
) map[string]struct{} {
	effectiveRoleSubjects := make(map[string]struct{}, len(roleIDs)+4)
	for _, roleID := range roleIDs {
		effectiveRoleSubjects["role:"+roleID] = struct{}{}
	}

	if m.enforcer == nil {
		return effectiveRoleSubjects
	}

	userSub := "user:" + userID
	for _, resolvedRole := range m.enforcer.GetRolesForUserInDomain(userSub, "*") {
		effectiveRoleSubjects[resolvedRole] = struct{}{}
	}

	for _, deptID := range deptIDs {
		domain := "dept:" + deptID
		for _, resolvedRole := range m.enforcer.GetRolesForUserInDomain(userSub, domain) {
			effectiveRoleSubjects[resolvedRole] = struct{}{}
		}
	}

	return effectiveRoleSubjects
}

// checkMultiRolePermission 多角色权限检查（取并集）
func (m *CasbinMiddleware) checkMultiRolePermission(
	ctx context.Context,
	userID string,
	roleIDs []string,
	deptIDs []string,
	resource string,
	action string,
) (bool, string, error) {
	maxDataScope := ""

	// 构建有效角色主体集合：JWT 角色ID + Casbin g 关系解析出的角色编码
	effectiveRoleSubjects := m.buildEffectiveRoleSubjects(userID, roleIDs, deptIDs)
	userSub := "user:" + userID

	// 为每个角色和部门组合检查权限
	for sub := range effectiveRoleSubjects {
		// 先检查全局权限（域为 *）
		allowed, dataScope, err := m.enforcer.EnforceWithDataScope(sub, "*", resource, action)
		if err != nil {
			return false, "", err
		}

		tracelog.Event(ctx, m.logger.Debug()).
			Str("subject", sub).
			Str("domain", "*").
			Str("resource", resource).
			Str("action", action).
			Bool("allowed", allowed).
			Str("data_scope", dataScope).
			Msg("Casbin role global-domain evaluation")

		if allowed {
			if compareDataScope(dataScope, maxDataScope) > 0 {
				maxDataScope = dataScope
			}
		}

		// 再检查每个部门的权限
		for _, deptID := range deptIDs {
			domain := "dept:" + deptID
			allowed, dataScope, err := m.enforcer.EnforceWithDataScope(sub, domain, resource, action)
			if err != nil {
				return false, "", err
			}

			tracelog.Event(ctx, m.logger.Debug()).
				Str("subject", sub).
				Str("domain", domain).
				Str("resource", resource).
				Str("action", action).
				Bool("allowed", allowed).
				Str("data_scope", dataScope).
				Msg("Casbin role dept-domain evaluation")

			if allowed {
				if compareDataScope(dataScope, maxDataScope) > 0 {
					maxDataScope = dataScope
				}
			}
		}
	}

	// 检查用户主体权限（可通过 g 关系间接命中角色策略）
	// 先检查全局域
	allowed, dataScope, err := m.enforcer.EnforceWithDataScope(userSub, "*", resource, action)
	if err != nil {
		return false, "", err
	}

	tracelog.Event(ctx, m.logger.Debug()).
		Str("subject", userSub).
		Str("domain", "*").
		Str("resource", resource).
		Str("action", action).
		Bool("allowed", allowed).
		Str("data_scope", dataScope).
		Msg("Casbin user global-domain evaluation")

	if allowed {
		if compareDataScope(dataScope, maxDataScope) > 0 {
			maxDataScope = dataScope
		}
	}

	// 再检查部门域
	for _, deptID := range deptIDs {
		domain := "dept:" + deptID
		allowed, dataScope, err := m.enforcer.EnforceWithDataScope(userSub, domain, resource, action)
		if err != nil {
			return false, "", err
		}

		tracelog.Event(ctx, m.logger.Debug()).
			Str("subject", userSub).
			Str("domain", domain).
			Str("resource", resource).
			Str("action", action).
			Bool("allowed", allowed).
			Str("data_scope", dataScope).
			Msg("Casbin direct-user evaluation")

		if allowed {
			if compareDataScope(dataScope, maxDataScope) > 0 {
				maxDataScope = dataScope
			}
		}
	}

	if maxDataScope == "" {
		tracelog.Event(ctx, m.logger.Warn()).
			Str("user_id", userID).
			Strs("role_ids", roleIDs).
			Strs("department_ids", deptIDs).
			Str("resource", resource).
			Str("action", action).
			Msg("Casbin denied: no matching policy found")
	}

	return maxDataScope != "", maxDataScope, nil
}

// CheckPermission 检查用户是否有权限访问指定资源
func (m *CasbinMiddleware) CheckPermission(
	ctx context.Context,
	userID string,
	roleIDs []string,
	deptIDs []string,
	resource string,
	action string,
) (bool, error) {
	allowed, _, err := m.checkMultiRolePermission(ctx, userID, roleIDs, deptIDs, resource, action)
	return allowed, err
}

// CheckPermissionWithDataScope 检查权限并返回数据范围
func (m *CasbinMiddleware) CheckPermissionWithDataScope(
	ctx context.Context,
	userID string,
	roleIDs []string,
	deptIDs []string,
	resource string,
	action string,
) (bool, string, error) {
	return m.checkMultiRolePermission(ctx, userID, roleIDs, deptIDs, resource, action)
}

// GetUserPermissions 获取用户的所有权限列表
func (m *CasbinMiddleware) GetUserPermissions(
	ctx context.Context,
	userID string,
	roleIDs []string,
) ([]PermissionInfo, error) {
	var permissions []PermissionInfo
	seen := make(map[string]bool)

	// 获取每个角色的权限
	for _, roleID := range roleIDs {
		perms := m.enforcer.GetPermissionsForUser("role:" + roleID)
		for _, perm := range perms {
			if len(perm) >= 4 {
				key := perm[0] + "|" + perm[1] + "|" + perm[2] + "|" + perm[3]
				if !seen[key] {
					seen[key] = true
					info := PermissionInfo{
						Domain:   perm[1],
						Resource: perm[2],
						Action:   perm[3],
					}
					if len(perm) >= 5 {
						info.DataScope = perm[4]
					}
					permissions = append(permissions, info)
				}
			}
		}
	}

	return permissions, nil
}

// RefreshUserPermissions 刷新用户权限缓存
// 在内存模式下，策略通过 PolicySyncService 从 RPC 同步
func (m *CasbinMiddleware) RefreshUserPermissions(ctx context.Context, userID string) error {
	// 内存模式下，策略通过 PolicySyncService 同步，此方法无需操作
	m.logger.Debug().Str("user_id", userID).Msg("RefreshUserPermissions called (no-op in memory mode)")
	return nil
}

// RefreshAllPermissions 刷新所有权限策略
// 在内存模式下，策略通过 PolicySyncService 从 RPC 同步
func (m *CasbinMiddleware) RefreshAllPermissions(ctx context.Context) error {
	// 内存模式下，策略通过 PolicySyncService 同步，此方法无需操作
	m.logger.Debug().Msg("RefreshAllPermissions called (no-op in memory mode)")
	return nil
}
