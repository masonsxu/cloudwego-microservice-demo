package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/redis"
)

// JWTMiddlewareService 定义JWT中间件服务接口
type JWTMiddlewareService interface {
	// 认证中间件
	MiddlewareFunc() app.HandlerFunc

	// 处理器方法
	LoginHandler(ctx context.Context, c *app.RequestContext)
	LogoutHandler(ctx context.Context, c *app.RequestContext)
	RefreshHandler(ctx context.Context, c *app.RequestContext)
}

// TokenCacheService Token缓存服务接口（直接使用redis包的接口）
type TokenCacheService = redis.TokenCacheService
