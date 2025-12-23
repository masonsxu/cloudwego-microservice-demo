package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
)

// CORSMiddleware CORS中间件实现
type CORSMiddleware struct {
	config *config.CORSConfig
	logger *hertzZerolog.Logger
}

// NewCORSMiddleware 创建新的CORS中间件实例
func NewCORSMiddleware(
	config *config.CORSConfig,
	logger *hertzZerolog.Logger,
) CORSMiddlewareService {
	if logger == nil {
		logger = hertzZerolog.New()
	}

	return &CORSMiddleware{
		config: config,
		logger: logger,
	}
}

// MiddlewareFunc 返回Hertz中间件处理函数
func (cm *CORSMiddleware) MiddlewareFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 如果中间件被禁用，直接跳过
		if !cm.config.Enabled {
			c.Next(ctx)
			return
		}

		// 设置CORS响应头
		cm.setAccessControlHeaders(c)

		// 如果是OPTIONS请求，直接返回200
		if string(c.Method()) == http.MethodOptions {
			cm.logger.Debugf("Handling CORS preflight request: origin=%s, method=%s",
				string(c.GetHeader("Origin")), string(c.GetHeader("Access-Control-Request-Method")))
			c.AbortWithStatus(http.StatusOK)

			return
		}

		c.Next(ctx)
	}
}

// setAccessControlHeaders 设置Access-Control相关头部
func (cm *CORSMiddleware) setAccessControlHeaders(c *app.RequestContext) {
	// 设置 Allow-Origin
	if len(cm.config.AllowOrigins) > 0 {
		origin := string(c.GetHeader("Origin"))
		if cm.isOriginAllowed(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		} else if cm.hasWildcardOrigin() {
			c.Header("Access-Control-Allow-Origin", "*")
		}
	} else {
		// 默认允许所有来源
		c.Header("Access-Control-Allow-Origin", "*")
	}

	// 设置 Allow-Methods
	if len(cm.config.AllowMethods) > 0 {
		c.Header("Access-Control-Allow-Methods", strings.Join(cm.config.AllowMethods, ", "))
	} else {
		// 默认允许的方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	}

	// 设置 Allow-Headers
	if len(cm.config.AllowHeaders) > 0 {
		c.Header("Access-Control-Allow-Headers", strings.Join(cm.config.AllowHeaders, ", "))
	} else {
		// 默认允许的头部
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
	}

	// 设置 Allow-Credentials
	if cm.config.AllowCredentials {
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	// 设置 Max-Age
	c.Header("Access-Control-Max-Age", "86400") // 24小时
}

// isOriginAllowed 检查来源是否被允许
func (cm *CORSMiddleware) isOriginAllowed(origin string) bool {
	for _, allowedOrigin := range cm.config.AllowOrigins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
	}

	return false
}

// hasWildcardOrigin 检查是否包含通配符来源
func (cm *CORSMiddleware) hasWildcardOrigin() bool {
	for _, origin := range cm.config.AllowOrigins {
		if origin == "*" {
			return true
		}
	}

	return false
}
