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
	enforcer    *CasbinEnforcer
	logger      *zerolog.Logger
	skipPaths   []string
	pathMapping map[string]string // API路径到权限编码的映射
}

// NewCasbinMiddleware 创建新的 Casbin 中间件
func NewCasbinMiddleware(enforcer *CasbinEnforcer, logger *zerolog.Logger) *CasbinMiddleware {
	return &CasbinMiddleware{
		enforcer:    enforcer,
		logger:      logger,
		skipPaths:   defaultSkipPaths(),
		pathMapping: make(map[string]string),
	}
}

// defaultSkipPaths 默认跳过权限检查的路径
func defaultSkipPaths() []string {
	return []string{
		"/login",
		"/logout",
		"/refresh",
		"/health",
		"/metrics",
		"/swagger",
		"/api/v1/auth/login",
		"/api/v1/auth/logout",
		"/api/v1/auth/refresh",
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

// MiddlewareFunc 返回权限校验中间件
func (m *CasbinMiddleware) MiddlewareFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Request.URI().Path())
		method := string(c.Request.Method())

		// 检查是否跳过权限检查
		if m.shouldSkip(path) {
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
			tracelog.Event(ctx, m.logger.Warn()).Str("path", path).Msg("No auth context found")
			abortWithUnauthorized(c, "Authentication required")
			return
		}

		// 获取用户信息
		userID, ok := auth_context.GetCurrentUserProfileID(c)
		if !ok {
			tracelog.Event(ctx, m.logger.Warn()).Str("path", path).Msg("No user ID found in auth context")
			abortWithUnauthorized(c, "User not authenticated")
			return
		}

		// 获取角色ID列表（多角色模式）
		roleIDs := auth_context.GetCurrentRoleIDs(c)
		if len(roleIDs) == 0 {
			tracelog.Event(ctx, m.logger.Warn()).Str("path", path).Str("user_id", userID).Msg("No roles found for user")
			abortWithPermissionDenied(c, "No roles assigned")
			return
		}

		// 获取部门ID列表
		deptIDs := auth_context.GetCurrentDepartmentIDs(c)

		// 将HTTP方法转换为操作类型
		action := methodToAction(method)

		// 获取权限资源（优先使用路径映射，其次使用API路径）
		resource := m.getResource(path)

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

	// 为每个角色和部门组合检查权限
	for _, roleID := range roleIDs {
		// 先检查全局权限（域为 *）
		allowed, dataScope, err := m.enforcer.EnforceWithDataScope("role:"+roleID, "*", resource, action)
		if err != nil {
			return false, "", err
		}
		if allowed {
			if compareDataScope(dataScope, maxDataScope) > 0 {
				maxDataScope = dataScope
			}
		}

		// 再检查每个部门的权限
		for _, deptID := range deptIDs {
			allowed, dataScope, err := m.enforcer.EnforceWithDataScope("role:"+roleID, "dept:"+deptID, resource, action)
			if err != nil {
				return false, "", err
			}
			if allowed {
				if compareDataScope(dataScope, maxDataScope) > 0 {
					maxDataScope = dataScope
				}
			}
		}
	}

	// 也检查用户直接授予的权限
	for _, deptID := range deptIDs {
		allowed, dataScope, err := m.enforcer.EnforceWithDataScope("user:"+userID, "dept:"+deptID, resource, action)
		if err != nil {
			return false, "", err
		}
		if allowed {
			if compareDataScope(dataScope, maxDataScope) > 0 {
				maxDataScope = dataScope
			}
		}
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
