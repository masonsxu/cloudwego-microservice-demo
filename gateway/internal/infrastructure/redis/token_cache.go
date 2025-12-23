package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	hertzZerolog "github.com/hertz-contrib/logger/zerolog"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
)

// TokenCacheService Token缓存服务接口
type TokenCacheService interface {
	// CacheToken 缓存Token
	CacheToken(ctx context.Context, token string, userID string, expiration time.Duration) error

	// IsTokenCached 检查Token是否已缓存（是否存在）
	IsTokenCached(ctx context.Context, token string) (bool, error)

	// RemoveToken 移除Token缓存
	RemoveToken(ctx context.Context, token string) error

	// RefreshTokenExpiration 刷新Token过期时间
	RefreshTokenExpiration(ctx context.Context, token string, expiration time.Duration) error

	// RemoveUserTokens 移除用户所有Token（用户主动登出所有设备）
	RemoveUserTokens(ctx context.Context, userID string) error

	// GetUserTokens 获取用户所有活跃Token
	GetUserTokens(ctx context.Context, userID string) ([]string, error)

	// IsTokenRevoked 检查token是否被吊销
	IsTokenRevoked(ctx context.Context, token string) (bool, error)

	// RevokeToken 吊销token，expiration为token剩余有效期
	RevokeToken(ctx context.Context, token string, expiration time.Duration) error
}

// TokenCache Token缓存服务实现
type TokenCache struct {
	client *Client
	logger *hertzZerolog.Logger
}

// NewTokenCache 创建Token缓存服务
func NewTokenCache(client *Client, logger *hertzZerolog.Logger) TokenCacheService {
	return &TokenCache{
		client: client,
		logger: logger,
	}
}

// hashToken 对Token进行SHA256哈希处理，避免明文存储
func (tc *TokenCache) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// getTokenKey 获取Token存储的Redis Key
func (tc *TokenCache) getTokenKey(tokenHash string) string {
	return fmt.Sprintf("gateway:jwt:%s", tokenHash)
}

// getUserTokensKey 获取用户Token集合的Redis Key
func (tc *TokenCache) getUserTokensKey(userID string) string {
	return fmt.Sprintf("gateway:user:%s:tokens", userID)
}

// CacheToken 缓存Token
func (tc *TokenCache) CacheToken(
	ctx context.Context,
	token string,
	userID string,
	expiration time.Duration,
) error {
	tokenHash := tc.hashToken(token)
	tokenKey := tc.getTokenKey(tokenHash)
	userTokensKey := tc.getUserTokensKey(userID)

	// 使用管道操作提高性能
	pipe := tc.client.GetClient().Pipeline()

	// 缓存Token信息，存储用户ID和创建时间
	tokenData := fmt.Sprintf(`{"userID":"%s","createdAt":%d}`, userID, time.Now().Unix())
	pipe.Set(ctx, tokenKey, tokenData, expiration)

	// 将Token添加到用户的Token集合中
	pipe.SAdd(ctx, userTokensKey, tokenHash)
	pipe.Expire(ctx, userTokensKey, expiration)

	_, err := pipe.Exec(ctx)
	if err != nil {
		tc.logger.Errorf("Failed to cache token: error=%v, userID=%s", err, userID)
		return fmt.Errorf("缓存Token失败: %w", err)
	}

	tc.logger.Infof("Token cached successfully: userID=%s, expiration=%v", userID, expiration)

	return nil
}

// IsTokenCached 检查Token是否已缓存
func (tc *TokenCache) IsTokenCached(ctx context.Context, token string) (bool, error) {
	tokenHash := tc.hashToken(token)
	tokenKey := tc.getTokenKey(tokenHash)

	exists, err := tc.client.Exists(ctx, tokenKey)
	if err != nil {
		return false, fmt.Errorf("检查Token缓存失败: %w", err)
	}

	return exists, nil
}

// RemoveToken 移除Token缓存
func (tc *TokenCache) RemoveToken(ctx context.Context, token string) error {
	tokenHash := tc.hashToken(token)
	tokenKey := tc.getTokenKey(tokenHash)

	// 首先获取Token信息以获取用户ID
	tokenData, err := tc.client.Get(ctx, tokenKey)
	if err != nil && err.Error() != "redis: nil" {
		return fmt.Errorf("获取Token信息失败: %w", err)
	}

	// 使用管道操作保证原子性
	pipe := tc.client.GetClient().Pipeline()

	// 删除Token Key
	pipe.Del(ctx, tokenKey)

	// 如果成功获取到Token数据，从用户Token集合中移除
	if tokenData != "" {
		// 解析Token数据获取用户ID
		var userID string
		// 简化解析，实际项目中可能需要更复杂的JSON解析
		if strings.Contains(tokenData, `"userID":"`) {
			start := strings.Index(tokenData, `"userID":"`) + 9

			end := strings.Index(tokenData[start:], `"`)
			if start > 8 && end > 0 {
				userID = tokenData[start : start+end]
			}
		}

		if userID != "" {
			userTokensKey := tc.getUserTokensKey(userID)
			pipe.SRem(ctx, userTokensKey, tokenHash)
		}
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		tc.logger.Errorf("Failed to remove token: error=%v", err)
		return fmt.Errorf("删除Token缓存失败: %w", err)
	}

	tc.logger.Infof("Token removed successfully")

	return nil
}

// RefreshTokenExpiration 刷新Token过期时间
func (tc *TokenCache) RefreshTokenExpiration(
	ctx context.Context,
	token string,
	expiration time.Duration,
) error {
	tokenHash := tc.hashToken(token)
	tokenKey := tc.getTokenKey(tokenHash)

	err := tc.client.Expire(ctx, tokenKey, expiration)
	if err != nil {
		tc.logger.Errorf("Failed to refresh token expiration: error=%v", err)
		return fmt.Errorf("刷新Token过期时间失败: %w", err)
	}

	tc.logger.Infof("Token expiration refreshed: expiration=%v", expiration)

	return nil
}

// RemoveUserTokens 移除用户所有Token
func (tc *TokenCache) RemoveUserTokens(ctx context.Context, userID string) error {
	userTokensKey := tc.getUserTokensKey(userID)

	// 获取用户所有Token
	tokenHashes, err := tc.client.SMembers(ctx, userTokensKey)
	if err != nil {
		tc.logger.Errorf("Failed to get user tokens: error=%v, userID=%s", err, userID)
		return fmt.Errorf("获取用户Token列表失败: %w", err)
	}

	// 批量删除Token
	if len(tokenHashes) > 0 {
		tokenKeys := make([]string, len(tokenHashes))
		for i, tokenHash := range tokenHashes {
			tokenKeys[i] = tc.getTokenKey(tokenHash)
		}

		// 使用管道批量删除
		pipe := tc.client.GetClient().Pipeline()
		pipe.Del(ctx, tokenKeys...)
		pipe.Del(ctx, userTokensKey)

		_, err = pipe.Exec(ctx)
		if err != nil {
			tc.logger.Errorf(
				"Failed to batch remove user tokens: error=%v, userID=%s, tokenCount=%d",
				err,
				userID,
				len(tokenHashes),
			)

			return fmt.Errorf("批量删除用户Token失败: %w", err)
		}

		tc.logger.Infof("User tokens removed successfully: userID=%s, tokenCount=%d",
			userID, len(tokenHashes))
	} else {
		tc.logger.Infof("No tokens found for user: userID=%s", userID)
	}

	return nil
}

// GetUserTokens 获取用户所有活跃Token
func (tc *TokenCache) GetUserTokens(ctx context.Context, userID string) ([]string, error) {
	userTokensKey := tc.getUserTokensKey(userID)

	tokenHashes, err := tc.client.SMembers(ctx, userTokensKey)
	if err != nil {
		tc.logger.Errorf("Failed to get user tokens: error=%v, userID=%s", err, userID)
		return nil, fmt.Errorf("获取用户Token列表失败: %w", err)
	}

	tc.logger.Debugf("Retrieved user tokens: userID=%s, tokenCount=%d", userID, len(tokenHashes))

	return tokenHashes, nil
}

// getRevokedTokenKey 获取吊销Token的Redis Key
func (tc *TokenCache) getRevokedTokenKey(tokenHash string) string {
	return fmt.Sprintf("gateway:jwt:revoked:%s", tokenHash)
}

// RevokeToken 吊销Token
func (tc *TokenCache) RevokeToken(
	ctx context.Context,
	token string,
	expiration time.Duration,
) error {
	tokenHash := tc.hashToken(token)
	tokenKey := tc.getRevokedTokenKey(tokenHash)

	// 存储吊销标记，值为吊销时间
	revokedAt := time.Now().Unix()

	err := tc.client.GetClient().Set(ctx, tokenKey, revokedAt, expiration).Err()
	if err != nil {
		tc.logger.Errorf("Failed to revoke token: error=%v", err)
		return fmt.Errorf("吊销Token失败: %w", err)
	}

	tc.logger.Infof("Token revoked successfully: expiration=%v", expiration)

	return nil
}

// IsTokenRevoked 检查Token是否被吊销
func (tc *TokenCache) IsTokenRevoked(ctx context.Context, token string) (bool, error) {
	tokenHash := tc.hashToken(token)
	tokenKey := tc.getRevokedTokenKey(tokenHash)

	exists, err := tc.client.Exists(ctx, tokenKey)
	if err != nil {
		tc.logger.Errorf("Failed to check if token is revoked: error=%v", err)
		return false, fmt.Errorf("检查Token吊销状态失败: %w", err)
	}

	return exists, nil
}

// ProvideRedisClient 提供Redis客户端
func ProvideRedisClient(cfg *config.RedisConfig) (*Client, error) {
	return NewClient(cfg)
}

// ProvideTokenCache 提供Token缓存服务
func ProvideTokenCache(client *Client, logger *hertzZerolog.Logger) TokenCacheService {
	return NewTokenCache(client, logger)
}
