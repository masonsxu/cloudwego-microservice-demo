// Package wire 基础设施层依赖注入提供者
package wire

import (
	"log"

	"github.com/google/wire"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/redis"
)

// InfrastructureSet 基础设施层依赖注入集合
var InfrastructureSet = wire.NewSet(
	ProvideConfig,
	ProvideLogger,
	ProvideIdentityClient,
	ProvideJWTConfig,
	ProvideRedisConfig,
	ProvideRedisClient,
	ProvideTokenCache,
	// ProvideCasbinManager,
)

// ProvideConfig 提供配置服务
func ProvideConfig() *config.Configuration {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	return config.Config
}

// ProvideLogger 提供日志服务
// 根据配置创建结构化日志实例
func ProvideLogger(cfg *config.Configuration) *hertzZerolog.Logger {
	logger, err := config.CreateLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	return logger
}

// ProvideIdentityClient 提供身份服务客户端
// 创建与身份认证RPC服务的客户端连接
func ProvideIdentityClient(logger *hertzZerolog.Logger) identitycli.IdentityClient {
	client, err := identitycli.NewIdentityClient()
	if err != nil {
		logger.Errorf("Failed to create identity client: %v", err)
		panic(err)
	}

	logger.Infof("Identity client created successfully")

	return client
}

// ProvideJWTConfig 提供JWT配置
// 从主配置中提取JWT相关配置
func ProvideJWTConfig(cfg *config.Configuration) *config.JWTConfig {
	return &cfg.Middleware.JWT
}

// ProvideDataLakeConfig 提供数据湖配置
// 从主配置中提取DataLake相关配置
func ProvideDataLakeConfig(cfg *config.Configuration) config.DataLakeConfig {
	return cfg.DataLake
}

// ProvideRedisConfig 提供Redis配置
// 从主配置中提取Redis相关配置
func ProvideRedisConfig(cfg *config.Configuration) *config.RedisConfig {
	return &cfg.Redis
}

// ProvideRedisClient 提供Redis客户端
// 创建并配置Redis客户端连接
func ProvideRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
	return redis.NewClient(cfg)
}

// ProvideTokenCache 提供Token缓存服务
// 创建Token缓存服务实例
func ProvideTokenCache(client *redis.Client, logger *hertzZerolog.Logger) redis.TokenCacheService {
	return redis.NewTokenCache(client, logger)
}
