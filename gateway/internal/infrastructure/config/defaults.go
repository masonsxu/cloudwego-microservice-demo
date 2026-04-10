package config

import (
	"time"

	"github.com/spf13/viper"
)

// setDefaults 设置默认配置值
func setDefaults(v *viper.Viper) {
	// 服务器默认值
	v.SetDefault("server.name", "api-gateway")
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.debug", false) // 默认关闭，开发环境通过 .env 设置为 true
	v.SetDefault("server.read_timeout", 30*time.Second)
	v.SetDefault("server.write_timeout", 30*time.Second)
	v.SetDefault("server.idle_timeout", 120*time.Second)

	// etcd默认值
	v.SetDefault("etcd.address", "localhost:2379")
	v.SetDefault("etcd.timeout", 5*time.Second)

	// 客户端默认值
	v.SetDefault("client.connection_timeout", 500*time.Millisecond)
	v.SetDefault("client.request_timeout", 3*time.Second)

	// 连接池默认值
	// MaxIdlePerAddress 估算: QPS_per_host * avg_response_time_sec
	v.SetDefault("client.pool.max_idle_per_address", 10)
	v.SetDefault("client.pool.min_idle_per_address", 2)
	v.SetDefault("client.pool.max_idle_global", 1000)
	v.SetDefault("client.pool.max_idle_timeout", time.Minute)

	// 默认服务配置
	v.SetDefault("client.services.identity.name", "identity-service")

	// 中间件默认值
	v.SetDefault("middleware.cors.enabled", true)
	// CORS 允许的来源、方法和头部
	// 注意：使用 withCredentials 时不能设置为 *，必须指定具体的 origin
	v.SetDefault("middleware.cors.allow_origins", []string{
		"http://localhost:5173",
	})
	v.SetDefault(
		"middleware.cors.allow_methods",
		[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	)
	v.SetDefault("middleware.cors.allow_headers", []string{
		"Content-Type",
		"Authorization",
		"X-Requested-With",
	})
	// 启用 credentials 支持跨域 Cookie
	v.SetDefault("middleware.cors.allow_credentials", true)

	v.SetDefault("middleware.rate_limit.enabled", false)
	v.SetDefault("middleware.jwt.enabled", true)
	v.SetDefault("middleware.oidc.enabled", true)
	v.SetDefault("middleware.oidc.issuer", "http://localhost:8080")
	v.SetDefault("middleware.oidc.access_token_lifespan", 30*time.Minute)
	v.SetDefault("middleware.oidc.refresh_token_lifespan", 7*24*time.Hour)
	v.SetDefault("middleware.oidc.auth_code_lifespan", 10*time.Minute)
	v.SetDefault("middleware.oidc.id_token_lifespan", 30*time.Minute)
	v.SetDefault("middleware.oidc.enforce_pkce", true)
	v.SetDefault("middleware.oidc.consent_page_url", "")
	v.SetDefault("middleware.casbin.enabled", true)
	v.SetDefault("middleware.casbin.model_path", "./config/casbin_model.conf")
	v.SetDefault("middleware.casbin.log_enabled", false)
	v.SetDefault("middleware.casbin.sync_interval", 300)
	v.SetDefault("middleware.casbin.skip_extra_paths", []string{})
	v.SetDefault("middleware.casbin.superadmin_bypass_enabled", true)
	v.SetDefault("middleware.casbin.superadmin_subjects", []string{"role:superadmin", "username:superadmin"})
	v.SetDefault("middleware.casbin.menu_mapping_file", "menu.yaml")
	v.SetDefault("middleware.jwt.signing_key", "")
	v.SetDefault("middleware.jwt.timeout", 30*time.Minute)
	v.SetDefault("middleware.jwt.max_refresh", 7*24*time.Hour)
	v.SetDefault("middleware.jwt.identity_key", "identity")
	v.SetDefault("middleware.jwt.realm", "API Gateway")
	v.SetDefault(
		"middleware.jwt.token_lookup",
		"header:Authorization,cookie:auth_token,query:token",
	)
	v.SetDefault("middleware.jwt.token_head_name", "Bearer")
	v.SetDefault("middleware.jwt.send_authorization", false)
	// JWT 跳过认证的路径列表（默认跳过健康检查、指标、OIDC 公开端点）
	v.SetDefault("middleware.jwt.skip_paths", []string{
		"/health",
		"/metrics",
		"/ping",
		"/.well-known/openid-configuration",
		"/keys",
		"/oauth/token",
		"/authorize",
		"/authorize/callback",
		"/login",
		"/revoke",
		"/oauth/introspect",
		"/userinfo",
	})

	// Cookie默认值
	v.SetDefault("middleware.jwt.cookie.send_cookie", true)
	v.SetDefault("middleware.jwt.cookie.cookie_name", "auth_token")
	v.SetDefault("middleware.jwt.cookie.cookie_max_age", 7*24*time.Hour)
	v.SetDefault("middleware.jwt.cookie.cookie_domain", "")
	v.SetDefault("middleware.jwt.cookie.cookie_path", "/")
	v.SetDefault("middleware.jwt.cookie.cookie_same_site", "lax")
	v.SetDefault("middleware.jwt.cookie.secure_cookie", false)
	v.SetDefault("middleware.jwt.cookie.cookie_http_only", true) // 生产环境必须启用，防止XSS攻击

	// 日志默认值
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "file")
	v.SetDefault("log.file_path", "./logs/api_gateway.log")
	v.SetDefault("log.max_size", 100)
	v.SetDefault("log.max_age", 30)
	v.SetDefault("log.max_backups", 10)

	// 链路追踪默认值 (OpenTelemetry OTLP gRPC)
	v.SetDefault("tracing.enabled", false)
	v.SetDefault("tracing.endpoint", "jaeger:4317")
	v.SetDefault("tracing.sampler_ratio", 0.1)
	v.SetDefault("tracing.ignore_paths", []string{"/health", "/metrics", "/ping", "/swagger/*"})

	// ErrorHandler 中间件默认配置
	v.SetDefault("middleware.error_handler.enabled", true)
	v.SetDefault("middleware.error_handler.enable_detailed_errors", false) // 生产环境建议关闭
	v.SetDefault("middleware.error_handler.enable_request_logging", true)
	v.SetDefault("middleware.error_handler.enable_response_logging", true)
	v.SetDefault("middleware.error_handler.enable_panic_recovery", true)
	v.SetDefault("middleware.error_handler.max_stack_trace_size", 4096)
	v.SetDefault("middleware.error_handler.enable_error_metrics", false)
	v.SetDefault("middleware.error_handler.error_response_timeout", 5000)

	// Redis 默认值
	v.SetDefault("redis.address", "localhost:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.min_idle_conns", 5)
	v.SetDefault("redis.max_retries", 3)
	v.SetDefault("redis.dial_timeout", 5*time.Second)
	v.SetDefault("redis.read_timeout", 3*time.Second)
	v.SetDefault("redis.write_timeout", 3*time.Second)
	v.SetDefault("redis.pool_timeout", 4*time.Second)
	v.SetDefault("redis.idle_timeout", 5*time.Minute)
	v.SetDefault("redis.idle_check_freq", 1*time.Minute)
}

// DefaultErrorHandlerConfig 返回默认的错误处理中间件配置
// 主要用于测试和快速初始化
func DefaultErrorHandlerConfig() ErrorHandlerConfig {
	return ErrorHandlerConfig{
		Enabled:               true,
		EnableDetailedErrors:  false, // 生产环境建议关闭
		EnableRequestLogging:  true,
		EnableResponseLogging: true,
		EnablePanicRecovery:   true,
		MaxStackTraceSize:     4096,
		EnableErrorMetrics:    false,
		ErrorResponseTimeout:  5000, // 5秒
	}
}

// DefaultDataLakeConfig 返回默认的 DataLake 配置
// 主要用于测试和快速初始化
func DefaultDataLakeConfig() DataLakeConfig {
	return DataLakeConfig{
		DataLakeURL: "http://localhost:8080", // 默认本地开发地址
	}
}
