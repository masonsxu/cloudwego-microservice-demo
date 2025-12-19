package casbin_middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// CasbinMiddleware Casbin权限控制中间件接口
// 基于hertz-contrib/casbin的权限中间件，提供类型安全的权限控制
type CasbinMiddleware interface {
	// RequiresPermissions 要求特定权限的中间件
	// 权限格式: resource:action (如: "users:read", "departments:create")
	// 域(domain)通过认证上下文自动提取，支持多租户权限隔离
	RequiresPermissions(permission string) []app.HandlerFunc

	// RequiresRoles 要求特定角色的中间件
	// 支持单个角色或多个角色（逗号分隔）
	RequiresRoles(roles string) []app.HandlerFunc

	// RequiresAnyPermissions 要求任一权限的中间件
	// 用户拥有其中任意一个权限即可通过验证
	RequiresAnyPermissions(permissions ...string) []app.HandlerFunc

	// RequiresAllPermissions 要求所有权限的中间件
	// 用户必须拥有所有指定权限才能通过验证
	RequiresAllPermissions(permissions ...string) []app.HandlerFunc

	// HasPermission 检查用户是否拥有特定权限（工具方法）
	// 不返回中间件，仅用于业务逻辑中的权限检查
	HasPermission(ctx context.Context, userID, domain, resource, action string) (bool, error)

	// HasRole 检查用户是否拥有特定角色（工具方法）
	HasRole(ctx context.Context, userID, domain, role string) (bool, error)

	// GetUserPermissions 获取用户所有权限（工具方法）
	GetUserPermissions(ctx context.Context, userID, domain string) ([]string, error)

	// GetUserRoles 获取用户所有角色（工具方法）
	GetUserRoles(ctx context.Context, userID, domain string) ([]string, error)
}

// PermissionConfig 权限中间件配置
type PermissionConfig struct {
	// SkipPaths 跳过权限检查的路径列表
	SkipPaths []string

	// EnableCache 是否启用权限缓存
	EnableCache bool

	// CacheTimeout 权限缓存超时时间（秒）
	CacheTimeout int

	// LogLevel 日志级别 (DEBUG, INFO, WARN, ERROR)
	LogLevel string

	// UnauthorizedHandler 未认证处理器
	UnauthorizedHandler app.HandlerFunc

	// ForbiddenHandler 无权限处理器
	ForbiddenHandler app.HandlerFunc
}

// DefaultPermissionConfig 默认权限配置
func DefaultPermissionConfig() *PermissionConfig {
	return &PermissionConfig{
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/favicon.ico",
		},
		EnableCache:  true,
		CacheTimeout: 300, // 5分钟
		LogLevel:     "INFO",
	}
}
