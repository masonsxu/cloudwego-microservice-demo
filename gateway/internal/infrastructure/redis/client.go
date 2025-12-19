package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
	"github.com/redis/go-redis/v9"
)

// Client Redis客户端封装
type Client struct {
	rdb    *redis.Client
	config *config.RedisConfig
}

// NewClient 创建Redis客户端
func NewClient(cfg *config.RedisConfig) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Address,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolTimeout:  cfg.PoolTimeout,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis连接失败: %w", err)
	}

	return &Client{
		rdb:    rdb,
		config: cfg,
	}, nil
}

// GetClient 获取原生Redis客户端
func (c *Client) GetClient() *redis.Client {
	return c.rdb
}

// Close 关闭Redis连接
func (c *Client) Close() error {
	return c.rdb.Close()
}

// Set 设置键值对
func (c *Client) Set(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) error {
	return c.rdb.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

// Del 删除键
func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}

// Expire 设置键的过期时间
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.rdb.Expire(ctx, key, expiration).Err()
}

// SAdd 向集合添加成员
func (c *Client) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return c.rdb.SAdd(ctx, key, members...).Err()
}

// SRem 从集合移除成员
func (c *Client) SRem(ctx context.Context, key string, members ...interface{}) error {
	return c.rdb.SRem(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func (c *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.rdb.SMembers(ctx, key).Result()
}

// SIsMember 检查成员是否在集合中
func (c *Client) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return c.rdb.SIsMember(ctx, key, member).Result()
}
