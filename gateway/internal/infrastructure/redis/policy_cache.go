package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	"github.com/rs/zerolog"

	tracelog "github.com/masonsxu/cloudwego-microservice-demo/gateway/pkg/log"
)

// PolicyCacheService 策略缓存服务接口
// 用于缓存 Casbin 权限检查结果，提高权限校验性能
type PolicyCacheService interface {
	// CachePermissionResult 缓存权限检查结果
	CachePermissionResult(
		ctx context.Context,
		key string,
		allowed bool,
		dataScope string,
		expiration time.Duration,
	) error

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
	logger *zerolog.Logger
	prefix string // 缓存键前缀
}

// NewPolicyCache 创建策略缓存服务
func NewPolicyCache(client *Client, logger *hertzZerolog.Logger) PolicyCacheService {
	var zlogger *zerolog.Logger

	if logger != nil {
		unwrapped := logger.Unwrap()
		zlogger = &unwrapped
	} else {
		nop := zerolog.Nop()
		zlogger = &nop
	}

	return &PolicyCache{
		client: client,
		logger: zlogger,
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
		tracelog.Event(ctx, pc.logger.Error()).Err(err).Msg("Failed to marshal permission result")
		return fmt.Errorf("序列化权限结果失败: %w", err)
	}

	err = pc.client.Set(ctx, cacheKey, string(data), expiration)
	if err != nil {
		tracelog.Event(ctx, pc.logger.Error()).Err(err).Str("key", key).Msg("Failed to cache permission result")
		return fmt.Errorf("缓存权限结果失败: %w", err)
	}

	tracelog.Event(ctx, pc.logger.Debug()).
		Str("key", key).
		Bool("allowed", allowed).
		Str("data_scope", dataScope).
		Msg("Permission result cached")

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

		tracelog.Event(ctx, pc.logger.Error()).Err(err).Str("key", key).Msg("Failed to get permission result")

		return false, "", false, fmt.Errorf("获取权限缓存失败: %w", err)
	}

	var result PermissionCacheResult
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		tracelog.Event(ctx, pc.logger.Error()).Err(err).Str("key", key).Msg("Failed to unmarshal permission result")
		return false, "", false, fmt.Errorf("解析权限缓存失败: %w", err)
	}

	tracelog.Event(ctx, pc.logger.Debug()).
		Str("key", key).
		Bool("allowed", result.Allowed).
		Msg("Permission result cache hit")

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
	var (
		cursor      uint64
		keysDeleted int
	)

	for {
		keys, nextCursor, err := rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			tracelog.Event(ctx, pc.logger.Error()).Err(err).Str("pattern", pattern).Msg("Failed to scan keys")
			return fmt.Errorf("扫描缓存键失败: %w", err)
		}

		if len(keys) > 0 {
			if err := rdb.Del(ctx, keys...).Err(); err != nil {
				tracelog.Event(ctx, pc.logger.Error()).Err(err).Int("count", len(keys)).Msg("Failed to delete keys")
				return fmt.Errorf("删除缓存键失败: %w", err)
			}

			keysDeleted += len(keys)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	tracelog.Event(ctx, pc.logger.Info()).
		Str("pattern", pattern).
		Int("keys_deleted", keysDeleted).
		Msg("Permission cache invalidated")

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
		tracelog.Event(ctx, pc.logger.Error()).Err(err).Str("user_id", userID).Msg("Failed to marshal user roles")
		return fmt.Errorf("序列化用户角色列表失败: %w", err)
	}

	err = pc.client.Set(ctx, cacheKey, string(data), expiration)
	if err != nil {
		tracelog.Event(ctx, pc.logger.Error()).Err(err).Str("user_id", userID).Msg("Failed to cache user roles")
		return fmt.Errorf("缓存用户角色列表失败: %w", err)
	}

	tracelog.Event(ctx, pc.logger.Debug()).
		Str("user_id", userID).
		Int("role_count", len(roleIDs)).
		Msg("User roles cached")

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

		tracelog.Event(ctx, pc.logger.Error()).Err(err).Str("user_id", userID).Msg("Failed to get user roles")

		return nil, false, fmt.Errorf("获取用户角色缓存失败: %w", err)
	}

	var roleIDs []string
	if err := json.Unmarshal([]byte(data), &roleIDs); err != nil {
		tracelog.Event(ctx, pc.logger.Error()).Err(err).Str("user_id", userID).Msg("Failed to unmarshal user roles")
		return nil, false, fmt.Errorf("解析用户角色缓存失败: %w", err)
	}

	tracelog.Event(ctx, pc.logger.Debug()).
		Str("user_id", userID).
		Int("role_count", len(roleIDs)).
		Msg("User roles cache hit")

	return roleIDs, true, nil
}

// InvalidateUserRoles 使用户角色缓存失效
func (pc *PolicyCache) InvalidateUserRoles(ctx context.Context, userID string) error {
	cacheKey := pc.getUserRolesKey(userID)

	err := pc.client.Del(ctx, cacheKey)
	if err != nil {
		tracelog.Event(ctx, pc.logger.Error()).Err(err).Str("user_id", userID).Msg("Failed to invalidate user roles")
		return fmt.Errorf("使用户角色缓存失效失败: %w", err)
	}

	tracelog.Event(ctx, pc.logger.Debug()).Str("user_id", userID).Msg("User roles cache invalidated")

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
