package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/hertz-contrib/jwt"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"

	authservice "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/service/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
)

// validateJWTConfig 验证 JWT 配置的合理性
func validateJWTConfig(cfg *config.JWTConfig) error {
	if cfg.PrivKeyPath == "" {
		return fmt.Errorf("JWT private key path cannot be empty")
	}

	if _, err := os.Stat(cfg.PrivKeyPath); err != nil {
		return fmt.Errorf("JWT private key file not found: %s", cfg.PrivKeyPath)
	}

	if cfg.PubKeyPath == "" {
		return fmt.Errorf("JWT public key path cannot be empty")
	}

	if _, err := os.Stat(cfg.PubKeyPath); err != nil {
		return fmt.Errorf("JWT public key file not found: %s", cfg.PubKeyPath)
	}

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

	if len(cfg.SkipPaths) == 0 {
		return fmt.Errorf("JWT skip paths cannot be empty")
	}

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

// parseSameSite 解析 SameSite 设置
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

// JWTMiddlewareProvider 创建 JWT 中间件实例
func JWTMiddlewareProvider(
	authService authservice.AuthService,
	jwtConfig *config.JWTConfig,
	tokenCache TokenCacheService,
	logger *hertzZerolog.Logger,
) (JWTMiddlewareService, error) {
	if err := validateJWTConfig(jwtConfig); err != nil {
		return nil, fmt.Errorf("JWT配置验证失败: %w", err)
	}

	tokenExtractor := NewDefaultTokenExtractor(jwtConfig)

	unwrapped := logger.Unwrap()
	zlogger := &unwrapped

	httpStatusMessageFunc := func(e error, ctx context.Context, c *app.RequestContext) string {
		return customHTTPStatusMessageFunc(e, ctx, c, zlogger)
	}

	// RS256: hertz-contrib/jwt 强制要求同时提供私钥（签发）和公钥（验签）
	// readKeys() 会依次调用 privateKey() + publicKey()，缺一不可
	mw, err := jwt.New(&jwt.HertzJWTMiddleware{
		Realm:            jwtConfig.Realm,
		SigningAlgorithm: "RS256",
		PrivKeyFile:      jwtConfig.PrivKeyPath,
		PubKeyFile:       jwtConfig.PubKeyPath,
		Timeout:          jwtConfig.Timeout,
		MaxRefresh:       jwtConfig.MaxRefresh,
		IdentityKey:      jwtConfig.IdentityKey,

		TokenLookup:   jwtConfig.TokenLookup,
		TokenHeadName: jwtConfig.TokenHeadName,

		SendAuthorization: jwtConfig.SendAuthorization,

		SendCookie:   jwtConfig.Cookie.SendCookie,
		CookieName:   jwtConfig.Cookie.CookieName,
		CookieDomain: jwtConfig.Cookie.CookieDomain,
		// CookiePath 在 hertz-contrib/jwt v1.0.4 中不支持，默认为 "/"
		SecureCookie:   jwtConfig.Cookie.SecureCookie,
		CookieHTTPOnly: jwtConfig.Cookie.CookieHTTPOnly,
		CookieSameSite: protocol.CookieSameSite(parseSameSite(jwtConfig.Cookie.CookieSameSite)),
		CookieMaxAge:   jwtConfig.Cookie.CookieMaxAge,

		TimeFunc: time.Now,

		PayloadFunc:     payloadFunc,
		IdentityHandler: identityHandler,
		Authenticator:   authenticatorWithoutAbort(authService),
		Authorizator:    authorizator,

		HTTPStatusMessageFunc: httpStatusMessageFunc,

		Unauthorized:    unauthorizedHandler,
		LoginResponse:   loginResponseHandler,
		LogoutResponse:  logoutResponseHandler,
		RefreshResponse: refreshResponseHandler,
	})
	if err != nil {
		return nil, fmt.Errorf("创建JWT中间件失败: %w", err)
	}

	if err := mw.MiddlewareInit(); err != nil {
		return nil, fmt.Errorf("初始化JWT中间件失败: %w", err)
	}

	return &JWTMiddlewareImpl{
		jwtConfig:      jwtConfig,
		mw:             mw,
		tokenCache:     tokenCache,
		tokenExtractor: tokenExtractor,
		logger:         zlogger,
	}, nil
}
