package middleware

import (
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
)

// TokenExtractor 定义 Token 提取接口
// 用于从请求中提取 JWT Token，支持多种提取策略（Header、Query、Cookie）
type TokenExtractor interface {
	// ExtractToken 从请求上下文中提取 Token 字符串
	ExtractToken(c *app.RequestContext) string
}

// DefaultTokenExtractor 默认的 Token 提取器实现
// 按照优先级顺序从以下位置提取 Token：
// 1. Authorization Header (Bearer token)
// 2. Query 参数 (token=xxx)
// 3. Cookie (如果启用)
type DefaultTokenExtractor struct {
	jwtConfig *config.JWTConfig
}

// NewDefaultTokenExtractor 创建默认 Token 提取器
func NewDefaultTokenExtractor(jwtConfig *config.JWTConfig) *DefaultTokenExtractor {
	return &DefaultTokenExtractor{
		jwtConfig: jwtConfig,
	}
}

// ExtractToken 从请求中提取 token 字符串
// 按照优先级顺序：Header > Query > Cookie
func (e *DefaultTokenExtractor) ExtractToken(c *app.RequestContext) string {
	// 1. 从 Authorization Header 提取
	token := string(c.GetHeader("Authorization"))
	if token != "" {
		// 移除 Bearer 前缀
		token = strings.TrimPrefix(token, "Bearer ")
		token = strings.TrimPrefix(token, e.jwtConfig.TokenHeadName+" ")

		return strings.TrimSpace(token)
	}

	// 2. 从 Query 参数提取
	token = c.Query("token")
	if token != "" {
		return token
	}

	// 3. 从 Cookie 提取（如果启用）
	if e.jwtConfig.Cookie.SendCookie {
		token = string(c.Cookie(e.jwtConfig.Cookie.CookieName))
		if token != "" {
			return token
		}
	}

	return ""
}
