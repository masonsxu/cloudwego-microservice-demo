package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
)

// PolicyCacheService 策略缓存服务接口
// 用于缓存 Casbin 权限检查结果，提高权限校验性能
type PolicyCacheService interface {
	// CachePermissionResult 缓存权限检查结果
	CachePermissionResult(ctx context.Context, key string, allowed bool, dataScope string, expiration time.Duration) error

	// GetPermissionResult 获取缓存的权限检查结果
	GetPermissionResult(ctx context.Context, key string) (allowed bool, dataScope string, exists bool, err error)

	// InvalidateUserPermissions 使用户权限缓存失效
	InvalidateUserPermissions(ctx context.Context, userID string) error

	// InvalidateRolePermissions 使角色权限缓存失效
	InvalidateRolePermissions(ctx context.Context, roleID string) error

	// InvalidateAllPermissions 使所有权限缓存失效
	InvalidateAllPermissions(ctx context.Context) error

	// CacheUserRoles 缓存用户角色列表
	CacheUserRoles(ctx context.Context, userID string, roleIDs []string, expiration time.Duration) error

	// GetUserRoles 获取缓存的用户角色列表
	GetUserRoles(ctx context.Context, userID string) ([]string, bool, error)

	// InvalidateUserRoles 使用户角色缓存失效
	InvalidateUserRoles(ctx context.Context, userID string) error
}

// PermissionCacheResult 权限缓存结果
type PermissionCacheResult struct {
	Allowed   bool   `json:"allowed"`
	DataScope string `json:"data_scope"`
	CachedAt  int64  `json:"cached_at"`
}

// PolicyCache 策略缓存服务实现
type PolicyCache struct {
	client *Client
	logger *hertzZerolog.Logger
	prefix string // 缓存键前缀
}

// NewPolicyCache 创建策略缓存服务
func NewPolicyCache(client *Client, logger *hertzZerolog.Logger) PolicyCacheService {
	return &PolicyCache{
		client: client,
		logger: logger,
		prefix: "gateway:permission:",
	}
}

// getPermissionKey 生成权限缓存键
func (pc *PolicyCache) getPermissionKey(key string) string {
	return fmt.Sprintf("%s%s", pc.prefix, key)
}

// getUserPermissionsPattern 获取用户权限缓存模式
func (pc *PolicyCache) getUserPermissionsPattern(userID string) string {
	return fmt.Sprintf("%suser:%s:*", pc.prefix, userID)
}

// getRolePermissionsPattern 获取角色权限缓存模式
func (pc *PolicyCache) getRolePermissionsPattern(roleID string) string {
	return fmt.Sprintf("%srole:%s:*", pc.prefix, roleID)
}

// getUserRolesKey 获取用户角色列表缓存键
func (pc *PolicyCache) getUserRolesKey(userID string) string {
	return fmt.Sprintf("%suser_roles:%s", pc.prefix, userID)
}

// CachePermissionResult 缓存权限检查结果
func (pc *PolicyCache) CachePermissionResult(
	ctx context.Context,
	key string,
	allowed bool,
	dataScope string,
	expiration time.Duration,
) error {
	cacheKey := pc.getPermissionKey(key)

	result := PermissionCacheResult{
		Allowed:   allowed,
		DataScope: dataScope,
		CachedAt:  time.Now().Unix(),
	}

	data, err := json.Marshal(result)
	if err != nil {
		pc.logger.Errorf("Failed to marshal permission result: error=%v", err)
		return fmt.Errorf("序列化权限结果失败: %w", err)
	}

	err = pc.client.Set(ctx, cacheKey, string(data), expiration)
	if err != nil {
		pc.logger.Errorf("Failed to cache permission result: error=%v, key=%s", err, key)
		return fmt.Errorf("缓存权限结果失败: %w", err)
	}

	pc.logger.Debugf("Permission result cached: key=%s, allowed=%v, dataScope=%s", key, allowed, dataScope)
	return nil
}

// GetPermissionResult 获取缓存的权限检查结果
func (pc *PolicyCache) GetPermissionResult(
	ctx context.Context,
	key string,
) (allowed bool, dataScope string, exists bool, err error) {
	cacheKey := pc.getPermissionKey(key)

	data, err := pc.client.Get(ctx, cacheKey)
	if err != nil {
		// Redis nil 表示键不存在
		if err.Error() == "redis: nil" {
			return false, "", false, nil
		}
		pc.logger.Errorf("Failed to get permission result: error=%v, key=%s", err, key)
		return false, "", false, fmt.Errorf("获取权限缓存失败: %w", err)
	}

	var result PermissionCacheResult
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		pc.logger.Errorf("Failed to unmarshal permission result: error=%v, key=%s", err, key)
		return false, "", false, fmt.Errorf("解析权限缓存失败: %w", err)
	}

	pc.logger.Debugf("Permission result cache hit: key=%s, allowed=%v", key, result.Allowed)
	return result.Allowed, result.DataScope, true, nil
}

// InvalidateUserPermissions 使用户权限缓存失效
func (pc *PolicyCache) InvalidateUserPermissions(ctx context.Context, userID string) error {
	pattern := pc.getUserPermissionsPattern(userID)
	return pc.invalidateByPattern(ctx, pattern)
}

// InvalidateRolePermissions 使角色权限缓存失效
func (pc *PolicyCache) InvalidateRolePermissions(ctx context.Context, roleID string) error {
	pattern := pc.getRolePermissionsPattern(roleID)
	return pc.invalidateByPattern(ctx, pattern)
}

// InvalidateAllPermissions 使所有权限缓存失效
func (pc *PolicyCache) InvalidateAllPermissions(ctx context.Context) error {
	pattern := pc.prefix + "*"
	return pc.invalidateByPattern(ctx, pattern)
}

// invalidateByPattern 根据模式使缓存失效
func (pc *PolicyCache) invalidateByPattern(ctx context.Context, pattern string) error {
	rdb := pc.client.GetClient()

	// 使用 SCAN 命令查找匹配的键
	var cursor uint64
	var keysDeleted int

	for {
		keys, nextCursor, err := rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			pc.logger.Errorf("Failed to scan keys: error=%v, pattern=%s", err, pattern)
			return fmt.Errorf("扫描缓存键失败: %w", err)
		}

		if len(keys) > 0 {
			if err := rdb.Del(ctx, keys...).Err(); err != nil {
				pc.logger.Errorf("Failed to delete keys: error=%v, count=%d", err, len(keys))
				return fmt.Errorf("删除缓存键失败: %w", err)
			}
			keysDeleted += len(keys)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	pc.logger.Infof("Permission cache invalidated: pattern=%s, keysDeleted=%d", pattern, keysDeleted)
	return nil
}

// CacheUserRoles 缓存用户角色列表
func (pc *PolicyCache) CacheUserRoles(
	ctx context.Context,
	userID string,
	roleIDs []string,
	expiration time.Duration,
) error {
	cacheKey := pc.getUserRolesKey(userID)

	data, err := json.Marshal(roleIDs)
	if err != nil {
		pc.logger.Errorf("Failed to marshal user roles: error=%v, userID=%s", err, userID)
		return fmt.Errorf("序列化用户角色列表失败: %w", err)
	}

	err = pc.client.Set(ctx, cacheKey, string(data), expiration)
	if err != nil {
		pc.logger.Errorf("Failed to cache user roles: error=%v, userID=%s", err, userID)
		return fmt.Errorf("缓存用户角色列表失败: %w", err)
	}

	pc.logger.Debugf("User roles cached: userID=%s, roleCount=%d", userID, len(roleIDs))
	return nil
}

// GetUserRoles 获取缓存的用户角色列表
func (pc *PolicyCache) GetUserRoles(ctx context.Context, userID string) ([]string, bool, error) {
	cacheKey := pc.getUserRolesKey(userID)

	data, err := pc.client.Get(ctx, cacheKey)
	if err != nil {
		// Redis nil 表示键不存在
		if err.Error() == "redis: nil" {
			return nil, false, nil
		}
		pc.logger.Errorf("Failed to get user roles: error=%v, userID=%s", err, userID)
		return nil, false, fmt.Errorf("获取用户角色缓存失败: %w", err)
	}

	var roleIDs []string
	if err := json.Unmarshal([]byte(data), &roleIDs); err != nil {
		pc.logger.Errorf("Failed to unmarshal user roles: error=%v, userID=%s", err, userID)
		return nil, false, fmt.Errorf("解析用户角色缓存失败: %w", err)
	}

	pc.logger.Debugf("User roles cache hit: userID=%s, roleCount=%d", userID, len(roleIDs))
	return roleIDs, true, nil
}

// InvalidateUserRoles 使用户角色缓存失效
func (pc *PolicyCache) InvalidateUserRoles(ctx context.Context, userID string) error {
	cacheKey := pc.getUserRolesKey(userID)

	err := pc.client.Del(ctx, cacheKey)
	if err != nil {
		pc.logger.Errorf("Failed to invalidate user roles: error=%v, userID=%s", err, userID)
		return fmt.Errorf("使用户角色缓存失效失败: %w", err)
	}

	pc.logger.Debugf("User roles cache invalidated: userID=%s", userID)
	return nil
}

// GeneratePermissionCacheKey 生成权限缓存键
// 格式: {userID}:{roleIDs}:{resource}:{action}
func GeneratePermissionCacheKey(userID string, roleIDs []string, resource, action string) string {
	rolesKey := ""
	for i, roleID := range roleIDs {
		if i > 0 {
			rolesKey += ","
		}
		rolesKey += roleID
	}
	return fmt.Sprintf("%s:[%s]:%s:%s", userID, rolesKey, resource, action)
}

// ProvidePolicyCache 提供策略缓存服务
func ProvidePolicyCache(client *Client, logger *hertzZerolog.Logger) PolicyCacheService {
	return NewPolicyCache(client, logger)
}
