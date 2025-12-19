package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/hertz-contrib/jwt"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	authservice "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/service/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
)

// validateJWTConfig 验证JWT配置的合理性
func validateJWTConfig(cfg *config.JWTConfig) error {
	// 验证 SigningKey 格式（Base64）
	if cfg.SigningKey == "" {
		return fmt.Errorf("JWT signing key cannot be empty")
	}

	if _, err := base64.StdEncoding.DecodeString(cfg.SigningKey); err != nil {
		return fmt.Errorf("JWT signing key must be valid Base64 encoded: %w", err)
	}

	// 验证 Timeout 和 MaxRefresh 的合理性
	if cfg.Timeout <= 0 {
		return fmt.Errorf("JWT timeout must be greater than 0")
	}

	if cfg.MaxRefresh <= 0 {
		return fmt.Errorf("JWT max refresh must be greater than 0")
	}

	if cfg.MaxRefresh < cfg.Timeout {
		return fmt.Errorf(
			"JWT max refresh (%v) must be greater than or equal to timeout (%v)",
			cfg.MaxRefresh,
			cfg.Timeout,
		)
	}

	// 验证 SkipPaths 不为空（至少包含登录路径）
	if len(cfg.SkipPaths) == 0 {
		return fmt.Errorf("JWT skip paths cannot be empty, at least login path should be included")
	}

	// 验证 Cookie 配置的合理性（如果启用）
	if cfg.Cookie.SendCookie {
		if cfg.Cookie.CookieName == "" {
			return fmt.Errorf("JWT cookie name cannot be empty when cookie is enabled")
		}

		if cfg.Cookie.CookieMaxAge <= 0 {
			return fmt.Errorf("JWT cookie max age must be greater than 0 when cookie is enabled")
		}
	}

	return nil
}

// parseSameSite 解析SameSite设置
func parseSameSite(sameSite string) int {
	switch strings.ToLower(sameSite) {
	case "lax":
		return int(http.SameSiteLaxMode)
	case "strict":
		return int(http.SameSiteStrictMode)
	case "none":
		return int(http.SameSiteNoneMode)
	default:
		return int(http.SameSiteDefaultMode)
	}
}

// JWTMiddlewareProvider 创建JWT中间件实例
// 这是依赖注入的入口函数，负责创建和配置JWT中间件
func JWTMiddlewareProvider(
	authService authservice.AuthService,
	jwtConfig *config.JWTConfig,
	tokenCache TokenCacheService,
	logger *hertzZerolog.Logger,
) (JWTMiddlewareService, error) {
	// 验证配置
	if err := validateJWTConfig(jwtConfig); err != nil {
		return nil, fmt.Errorf("JWT配置验证失败: %w", err)
	}

	// Base64 解码签名密钥
	signingKey, err := base64.StdEncoding.DecodeString(jwtConfig.SigningKey)
	if err != nil {
		return nil, fmt.Errorf("签名密钥解码失败: %w", err)
	}

	// 创建 Token 提取器
	tokenExtractor := NewDefaultTokenExtractor(jwtConfig)

	// 创建 HTTP 状态消息处理函数（适配 hertz-contrib/jwt 的接口）
	// hertz-contrib/jwt 的 HTTPStatusMessageFunc 签名是：
	// func(e error, ctx context.Context, c *app.RequestContext) string
	// 我们创建一个闭包来传递 logger
	httpStatusMessageFunc := func(e error, ctx context.Context, c *app.RequestContext) string {
		return customHTTPStatusMessageFunc(e, ctx, c, logger)
	}

	// 创建 hertz-contrib/jwt 中间件
	mw, err := jwt.New(&jwt.HertzJWTMiddleware{
		Realm:            jwtConfig.Realm,
		SigningAlgorithm: "HS256",
		Key:              signingKey,
		Timeout:          jwtConfig.Timeout,
		MaxRefresh:       jwtConfig.MaxRefresh,
		IdentityKey:      jwtConfig.IdentityKey,

		// Token查找策略（支持Header、Query、Cookie）
		TokenLookup:   jwtConfig.TokenLookup,
		TokenHeadName: jwtConfig.TokenHeadName,

		// Authorization Header 配置
		SendAuthorization: jwtConfig.SendAuthorization,

		// Cookie配置（前后端分离架构）
		SendCookie:   jwtConfig.Cookie.SendCookie,
		CookieName:   jwtConfig.Cookie.CookieName,
		CookieDomain: jwtConfig.Cookie.CookieDomain,
		// 注意：CookiePath 在当前版本的 hertz-contrib/jwt (v1.0.4) 中不支持，默认为 "/"
		SecureCookie:   jwtConfig.Cookie.SecureCookie,
		CookieHTTPOnly: jwtConfig.Cookie.CookieHTTPOnly,
		CookieSameSite: protocol.CookieSameSite(parseSameSite(jwtConfig.Cookie.CookieSameSite)),
		CookieMaxAge:   jwtConfig.Cookie.CookieMaxAge,

		// 时间函数：方便单测 mock
		TimeFunc: time.Now,

		// 核心处理函数
		PayloadFunc:     payloadFunc,
		IdentityHandler: identityHandler,
		Authenticator:   authenticatorWithoutAbort(authService),
		Authorizator:    authorizator,

		// 关键：使用自定义的HTTP状态消息函数
		HTTPStatusMessageFunc: httpStatusMessageFunc,

		// 未认证处理
		Unauthorized: unauthorizedHandler,
		// 登录响应处理
		LoginResponse: loginResponseHandler,
		// 登出响应处理
		LogoutResponse: logoutResponseHandler,
		// 刷新Token响应处理
		RefreshResponse: refreshResponseHandler,
	})
	if err != nil {
		return nil, fmt.Errorf("创建JWT中间件失败: %w", err)
	}

	// 初始化中间件
	if err := mw.MiddlewareInit(); err != nil {
		return nil, fmt.Errorf("初始化JWT中间件失败: %w", err)
	}

	// 创建中间件实现实例
	return &JWTMiddlewareImpl{
		jwtConfig:      jwtConfig,
		mw:             mw,
		tokenCache:     tokenCache,
		tokenExtractor: tokenExtractor,
		logger:         logger,
	}, nil
}
