package config

import "time"

// Configuration 总配置结构
// 包含服务器、服务发现、客户端、中间件、日志、Tracing、Metrics 等配置段。
// 加载顺序：默认值 -> 配置文件(config.yaml) -> .env -> 环境变量(同名覆盖)。
type Configuration struct {
	Server     ServerConfig     `mapstructure:"server"`
	Etcd       EtcdConfig       `mapstructure:"etcd"`
	Client     ClientConfig     `mapstructure:"client"`
	Middleware MiddlewareConfig `mapstructure:"middleware"`
	Log        LogConfig        `mapstructure:"log"`
	Tracing    TracingConfig    `mapstructure:"tracing"`
	DataLake   DataLakeConfig   `mapstructure:"data_lake"`
	Redis      RedisConfig      `mapstructure:"redis"`
}

// ServerConfig 服务器配置
// 相关环境变量：SERVER_NAME, SERVER_HOST, SERVER_PORT, SERVER_DEBUG,
// SERVER_READ_TIMEOUT, SERVER_WRITE_TIMEOUT, SERVER_IDLE_TIMEOUT
type ServerConfig struct {
	Name         string        `mapstructure:"name"`  // 服务名称（用于服务发现、RPC调用标识）
	Host         string        `mapstructure:"host"`  // 服务监听主机
	Port         int           `mapstructure:"port"`  // 服务监听端口
	Debug        bool          `mapstructure:"debug"` // 调试模式开关（控制日志详细度、错误堆栈、SQL日志等）
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// EtcdConfig etcd配置
// 相关环境变量：ETCD_ADDRESS, ETCD_USERNAME, ETCD_PASSWORD, ETCD_TIMEOUT
// 用于服务发现和配置中心
type EtcdConfig struct {
	Address  string `mapstructure:"address"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Timeout  int    `mapstructure:"timeout"`
}

// ClientConfig 客户端配置
// 相关环境变量：CLIENT_CONNECTION_TIMEOUT, CLIENT_REQUEST_TIMEOUT, CLIENT_POOL_MAX_IDLE_PER_ADDRESS,
// CLIENT_POOL_MAX_IDLE_GLOBAL, CLIENT_POOL_MAX_IDLE_TIMEOUT
// 用于配置 RPC 客户端的连接和请求超时、连接池等参数
type ClientConfig struct {
	ConnectionTimeout time.Duration            `mapstructure:"connection_timeout"`
	RequestTimeout    time.Duration            `mapstructure:"request_timeout"`
	Pool              ConnectionPoolConfig     `mapstructure:"pool"`
	Services          map[string]ServiceConfig `mapstructure:"services"`
}

// ConnectionPoolConfig 连接池配置
// 相关环境变量：CLIENT_POOL_MAX_IDLE_PER_ADDRESS, CLIENT_POOL_MAX_IDLE_GLOBAL, CLIENT_POOL_MAX_IDLE_TIMEOUT
type ConnectionPoolConfig struct {
	MaxIdlePerAddress int           `mapstructure:"max_idle_per_address"`
	MaxIdleGlobal     int           `mapstructure:"max_idle_global"`
	MaxIdleTimeout    time.Duration `mapstructure:"max_idle_timeout"`
}

// ServiceConfig 服务配置
// 相关环境变量：IDENTITY_SRV_NAME（示例服务名称）
type ServiceConfig struct {
	Name string `mapstructure:"name"`
}

// MiddlewareConfig 中间件配置
// 包含 CORS、限流、JWT、错误处理、Casbin 等中间件配置
type MiddlewareConfig struct {
	CORS         CORSConfig         `mapstructure:"cors"`
	RateLimit    RateLimitConfig    `mapstructure:"rate_limit"`
	JWT          JWTConfig          `mapstructure:"jwt"`
	ErrorHandler ErrorHandlerConfig `mapstructure:"error_handler"`
	Casbin       CasbinConfig       `mapstructure:"casbin"`
}

// CORSConfig CORS配置
// 相关环境变量：CORS_ENABLED, CORS_ALLOW_ORIGINS, CORS_ALLOW_CREDENTIALS
// 用于配置跨域资源共享策略
type CORSConfig struct {
	Enabled          bool     `mapstructure:"enabled"`
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

// RateLimitConfig 限流配置
// 相关环境变量：RATE_LIMIT_ENABLED, RATE_LIMIT_RPS, RATE_LIMIT_BURST
// 用于控制请求速率，防止服务过载
type RateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	RequestsPerSecond int  `mapstructure:"requests_per_second"`
	Burst             int  `mapstructure:"burst"`
}

// JWTConfig 身份验证配置
// 相关环境变量：JWT_ENABLED, JWT_SIGNING_KEY, JWT_TIMEOUT, JWT_MAX_REFRESH, JWT_IDENTITY_KEY,
// JWT_REALM, JWT_TOKEN_LOOKUP, JWT_TOKEN_HEAD_NAME, JWT_SEND_AUTHORIZATION, JWT_SKIP_PATHS,
// JWT_COOKIE_SEND_COOKIE, JWT_COOKIE_COOKIE_NAME, JWT_COOKIE_COOKIE_DOMAIN, JWT_COOKIE_COOKIE_PATH,
// JWT_COOKIE_COOKIE_MAX_AGE, JWT_COOKIE_COOKIE_SAME_SITE, JWT_COOKIE_SECURE_COOKIE, JWT_COOKIE_HTTP_ONLY
// 用于配置 JWT 认证和 Cookie 相关设置
type JWTConfig struct {
	Realm             string        `mapstructure:"realm"`              // 认证领域
	SigningKey        string        `mapstructure:"signing_key"`        // HS265 密钥
	Timeout           time.Duration `mapstructure:"timeout"`            // access-token 有效期(秒)
	MaxRefresh        time.Duration `mapstructure:"max_refresh"`        // refresh-token 有效期(秒)
	IdentityKey       string        `mapstructure:"identity_key"`       // JWT中存储用户标识的键
	SkipPaths         []string      `mapstructure:"skip_paths"`         // 跳过认证的路径列表
	TokenLookup       string        `mapstructure:"token_lookup"`       // 获取token的lookup方式
	TokenHeadName     string        `mapstructure:"token_head_name"`    // token头前缀
	SendAuthorization bool          `mapstructure:"send_authorization"` // 是否在响应中返回 Authorization header

	// Cookie配置（前后端分离架构）
	Cookie CookieConfig `mapstructure:"cookie"` // Cookie配置
}

// CookieConfig 前后端分离Cookie配置
// 用于配置 JWT token 在 Cookie 中的存储和传输方式
type CookieConfig struct {
	SendCookie     bool          `mapstructure:"send_cookie"`      // 是否同时写回 cookie
	CookieName     string        `mapstructure:"cookie_name"`      // Cookie 名称
	CookieMaxAge   time.Duration `mapstructure:"cookie_max_age"`   // Cookie 最大有效期(秒)
	CookieDomain   string        `mapstructure:"cookie_domain"`    // Cookie 域
	CookiePath     string        `mapstructure:"cookie_path"`      // Cookie 路径（默认"/"）
	CookieSameSite string        `mapstructure:"cookie_same_site"` // Cookie SameSite 策略
	SecureCookie   bool          `mapstructure:"secure_cookie"`    // 是否仅通过 HTTPS 发送 Cookie
	CookieHTTPOnly bool          `mapstructure:"cookie_http_only"` // 防止 JavaScript 访问（防 XSS 攻击）
}

// LogConfig 日志配置
// Level: info/debug/warn/error
// Format: json/text
// Output: stdout/stderr/file；当为 file 时使用 FilePath 并遵循滚动策略(MaxSize/MaxAge/MaxBackups)
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// TracingConfig 链路追踪配置
// 当 Enabled=true 时，向 Endpoint 上报链路数据；SamplerRatio 控制采样率[0.0,1.0]
type TracingConfig struct {
	Enabled      bool     `mapstructure:"enabled"`
	Endpoint     string   `mapstructure:"endpoint"`
	SamplerRatio float64  `mapstructure:"sampler_ratio"`
	IgnorePaths  []string `mapstructure:"ignore_paths"`
}

// ErrorHandlerConfig 错误处理中间件配置
// 相关环境变量：ERROR_HANDLER_ENABLED, ERROR_HANDLER_ENABLE_DETAILED_ERRORS,
// ERROR_HANDLER_ENABLE_REQUEST_LOGGING, ERROR_HANDLER_ENABLE_RESPONSE_LOGGING, ERROR_HANDLER_ENABLE_PANIC_RECOVERY
type ErrorHandlerConfig struct {
	Enabled               bool `mapstructure:"enabled"`                 // 是否启用错误处理中间件
	EnableDetailedErrors  bool `mapstructure:"enable_detailed_errors"`  // 是否启用详细错误信息（生产环境建议关闭）
	EnableRequestLogging  bool `mapstructure:"enable_request_logging"`  // 是否记录请求信息到日志
	EnableResponseLogging bool `mapstructure:"enable_response_logging"` // 是否记录响应信息到日志
	EnablePanicRecovery   bool `mapstructure:"enable_panic_recovery"`   // 是否启用panic恢复
	MaxStackTraceSize     int  `mapstructure:"max_stack_trace_size"`    // 最大堆栈跟踪大小（字节）
	EnableErrorMetrics    bool `mapstructure:"enable_error_metrics"`    // 是否启用错误指标收集
	ErrorResponseTimeout  int  `mapstructure:"error_response_timeout"`  // 错误响应超时时间（毫秒）
}

// CasbinConfig Casbin 权限控制配置
// 相关环境变量：CASBIN_ENABLED, CASBIN_MODEL_PATH, CASBIN_SKIP_PATHS, CASBIN_ENABLE_LOG
type CasbinConfig struct {
	Enabled   bool     `mapstructure:"enabled"`    // 是否启用Casbin权限校验
	SkipPaths []string `mapstructure:"skip_paths"` // 跳过权限校验的路径列表
}

// DataLakeConfig DataLake 配置
// 相关环境变量：DATALAKE_URL
// 用于配置 DataLake 服务地址
type DataLakeConfig struct {
	DataLakeURL string `mapstructure:"data_lake_url"` // DataLake 服务地址
}

// RedisConfig Redis配置
// 相关环境变量：REDIS_ADDRESS, REDIS_PASSWORD, REDIS_DB, REDIS_POOL_SIZE, REDIS_MIN_IDLE_CONNS,
// REDIS_MAX_RETRIES, REDIS_DIAL_TIMEOUT, REDIS_READ_TIMEOUT, REDIS_WRITE_TIMEOUT,
// REDIS_POOL_TIMEOUT, REDIS_IDLE_TIMEOUT, REDIS_IDLE_CHECK_FREQ
// 用于配置 Redis 连接池和超时参数
type RedisConfig struct {
	Address       string        `mapstructure:"address"`
	Password      string        `mapstructure:"password"`
	DB            int           `mapstructure:"db"`
	PoolSize      int           `mapstructure:"pool_size"`
	MinIdleConns  int           `mapstructure:"min_idle_conns"`
	MaxRetries    int           `mapstructure:"max_retries"`
	DialTimeout   time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout   time.Duration `mapstructure:"read_timeout"`
	WriteTimeout  time.Duration `mapstructure:"write_timeout"`
	PoolTimeout   time.Duration `mapstructure:"pool_timeout"`
	IdleTimeout   time.Duration `mapstructure:"idle_timeout"`
	IdleCheckFreq time.Duration `mapstructure:"idle_check_freq"`
}
