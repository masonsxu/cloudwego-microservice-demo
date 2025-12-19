package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/middleware/common"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/errors"
)

// JWTMiddlewareImpl JWT中间件实现
type JWTMiddlewareImpl struct {
	jwtConfig      *config.JWTConfig
	mw             *jwt.HertzJWTMiddleware
	tokenCache     TokenCacheService
	tokenExtractor TokenExtractor
	logger         *hertzZerolog.Logger
}

// MiddlewareFunc 返回JWT认证中间件函数
func (m *JWTMiddlewareImpl) MiddlewareFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 检查是否需要跳过认证
		if common.ShouldSkip(c, m.jwtConfig.SkipPaths) {
			c.Next(ctx)
			return
		}

		// 检查Token是否被吊销
		tokenString := m.tokenExtractor.ExtractToken(c)
		if tokenString != "" {
			if isRevoked, err := m.tokenCache.IsTokenRevoked(ctx, tokenString); err == nil &&
				isRevoked {
				m.logger.Warnf("Access denied: token has been revoked")
				c.JSON(http.StatusUnauthorized, &http_base.OperationStatusResponseDTO{
					BaseResp: &http_base.BaseResponseDTO{
						Code:    errors.ErrJWTTokenExpired.Code(),
						Message: errors.ErrJWTTokenExpired.Message(),
					},
				})
				c.Abort()

				return
			}
		}

		// 调用底层JWT中间件
		m.mw.MiddlewareFunc()(ctx, c)
	}
}

// LoginHandler 处理登录请求
func (m *JWTMiddlewareImpl) LoginHandler(ctx context.Context, c *app.RequestContext) {
	m.mw.LoginHandler(ctx, c)
}

// LogoutHandler 处理登出请求
// 从请求中提取Token并吊销，即使Token无效也返回成功响应
func (m *JWTMiddlewareImpl) LogoutHandler(ctx context.Context, c *app.RequestContext) {
	// 从请求中提取token字符串
	tokenString := m.tokenExtractor.ExtractToken(c)
	if tokenString == "" {
		m.logger.Warnf("No token found in request during logout")
		logoutResponseHandler(ctx, c, http.StatusOK)

		return
	}

	// 使用 hertz-contrib/jwt 提供的 ExtractClaims 获取 claims
	claims := jwt.ExtractClaims(ctx, c)
	if claims == nil {
		m.logger.Warnf("Failed to extract claims from token")
		logoutResponseHandler(ctx, c, http.StatusOK)

		return
	}

	m.logger.Debugf("JWT claims extracted: claims=%v", claims)

	// 获取token过期时间
	expClaim, exists := claims["exp"]
	if !exists {
		m.logger.Warnf("Token does not have expiration claim")
		logoutResponseHandler(ctx, c, http.StatusOK)

		return
	}

	// 转换过期时间
	var expTime float64

	switch v := expClaim.(type) {
	case float64:
		expTime = v
	case int64:
		expTime = float64(v)
	case int:
		expTime = float64(v)
	default:
		m.logger.Warnf("Invalid expiration claim type")
		logoutResponseHandler(ctx, c, http.StatusOK)

		return
	}

	// 计算剩余有效期
	exp := time.Unix(int64(expTime), 0)

	now := time.Now()
	if exp.Before(now) {
		// token已经过期，无需吊销
		m.logger.Debugf("Token already expired, no need to revoke")
		logoutResponseHandler(ctx, c, http.StatusOK)

		return
	}

	expiration := exp.Sub(now)

	// 吊销token
	if err := m.tokenCache.RevokeToken(ctx, tokenString, expiration); err != nil {
		m.logger.Errorf("Failed to revoke token during logout: %v", err)
		// 即使吊销失败，也继续返回登出成功
	} else {
		m.logger.Infof("Token successfully revoked during logout")
	}

	// 返回登出响应
	logoutResponseHandler(ctx, c, http.StatusOK)
}

// RefreshHandler 处理刷新Token请求
func (m *JWTMiddlewareImpl) RefreshHandler(ctx context.Context, c *app.RequestContext) {
	m.mw.RefreshHandler(ctx, c)
}
