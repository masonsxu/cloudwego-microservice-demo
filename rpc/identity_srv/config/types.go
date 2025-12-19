package config

import "time"

// Config 总配置结构
// 包含服务、数据库、日志、Tracing、Metrics 等配置段。
// 加载顺序：默认值 -> .env -> 环境变量(同名覆盖)。
type Config struct {
	Database    DatabaseConfig    `mapstructure:"database"`
	Server      ServerConfig      `mapstructure:"server"`
	HealthCheck HealthCheckConfig `mapstructure:"health_check"`
	Etcd        EtcdConfig        `mapstructure:"etcd"`
	Log         LogConfig         `mapstructure:"log"`
	Tracing     TracingConfig     `mapstructure:"tracing"`
	LogoStorage LogoStorageConfig `mapstructure:"logo_storage"`
	Casbin      CasbinConfig      `mapstructure:"casbin"`
	SuperAdmin  SuperAdminConfig  `mapstructure:"super_admin"`
}

// DatabaseConfig 数据库配置
// 相关环境变量：DB_HOST, DB_PORT, DB_USERNAME, DB_PASSWORD, DB_NAME,
// DB_SSLMODE, DB_TIMEZONE, DB_DRIVER, DB_MAX_IDLE_CONNS, DB_MAX_OPEN_CONNS, DB_CONN_MAX_LIFETIME
type DatabaseConfig struct {
	// 基础配置
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	Timezone string `mapstructure:"timezone"`

	// 连接池配置
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`     // 最大空闲连接数
	MaxOpenConns    int           `mapstructure:"max_open_conns"`     // 最大打开连接数
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`  // 连接最大生命周期(分钟)
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"` // 连接最大空闲时间(分钟)
}

// ServerConfig 服务器配置
// 相关环境变量：SERVER_NAME, SERVER_HOST, SERVER_PORT, SERVER_ADDRESS, SERVER_DEBUG
// Address 可直接提供监听地址（如 ":8891" 或 "0.0.0.0:8891"），
// 若未设置 Address，则由 Host + Port 组合生成。
type ServerConfig struct {
	Name    string `mapstructure:"name"`    // 服务名称（用于服务发现、RPC调用标识）
	Host    string `mapstructure:"host"`    // 服务监听主机
	Port    int    `mapstructure:"port"`    // 服务监听端口
	Address string `mapstructure:"address"` // 兼容旧配置，可直接提供完整地址
	Debug   bool   `mapstructure:"debug"`   // 调试模式开关（控制日志详细度、GORM SQL日志等）
}

// HealthCheckConfig 健康检查配置
// 相关环境变量：HEALTH_CHECK_PORT
// 健康检查服务器运行在独立的 HTTP 端口上，提供 /live 和 /ready 端点
type HealthCheckConfig struct {
	Port int `mapstructure:"port"` // 健康检查服务器端口
}

// EtcdConfig etcd配置
// 相关环境变量：ETCD_ADDRESS, ETCD_USERNAME, ETCD_PASSWORD, ETCD_TIMEOUT
type EtcdConfig struct {
	Address  string `mapstructure:"address"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Timeout  int    `mapstructure:"timeout"`
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
	Enabled      bool    `mapstructure:"enabled"`
	Endpoint     string  `mapstructure:"endpoint"`
	SamplerRatio float64 `mapstructure:"sampler_ratio"`
}

// LogoStorageConfig 组织Logo存储配置
// 用于组织Logo的上传、存储和访问管理，基于S3兼容存储（MinIO/RustFS）
// 相关环境变量：LOGO_STORAGE_S3_ENDPOINT, LOGO_STORAGE_S3_PUBLIC_ENDPOINT, LOGO_STORAGE_S3_REGION,
// LOGO_STORAGE_S3_USE_SSL, LOGO_STORAGE_USE_PATH_STYLE, LOGO_STORAGE_ACCESS_KEY, LOGO_STORAGE_SECRET_KEY,
// LOGO_STORAGE_MAX_FILE_SIZE, LOGO_STORAGE_ALLOWED_FILE_TYPES
type LogoStorageConfig struct {
	// S3兼容存储配置
	S3Endpoint       string `mapstructure:"s3_endpoint"`        // S3服务内部端点地址（容器间通信）
	S3PublicEndpoint string `mapstructure:"s3_public_endpoint"` // S3服务公共端点地址（生成预签名URL，浏览器可访问）
	S3Region         string `mapstructure:"s3_region"`          // S3区域
	S3UseSSL         bool   `mapstructure:"s3_use_ssl"`         // 是否使用SSL连接
	UsePathStyle     bool   `mapstructure:"use_path_style"`     // 使用Path Style模式（true）或Virtual Host Style（false）
	AccessKey        string `mapstructure:"access_key"`         // S3访问密钥
	SecretKey        string `mapstructure:"secret_key"`         // S3私钥
	S3LogEnabled     bool   `mapstructure:"s3_log_enabled"`     // S3操作日志启用选项
	Timeout          int    `mapstructure:"timeout"`            // 超时时间（秒）

	// 文件管理配置
	MaxFileSize      int64    `mapstructure:"max_file_size"`      // 最大文件大小（字节，默认10MB）
	AllowedFileTypes []string `mapstructure:"allowed_file_types"` // 允许的图片类型（如 image/png, image/jpeg）
}

// CasbinConfig Casbin 配置
// 相关环境变量：CASBIN_MODEL_PATH, CASBIN_ENABLE_LOG
type CasbinConfig struct {
	ModelPath string `mapstructure:"model_path"`
	EnableLog bool   `mapstructure:"enable_log"`
}

// SuperAdminConfig 超级管理员配置
// 相关环境变量：SUPER_ADMIN_ROLE_NAMES
type SuperAdminConfig struct {
	// RoleNames 超管角色名称列表，这些角色将拥有所有菜单的完整权限
	// 支持多个角色名称，例如：["super_admin", "system_admin"]
	RoleNames []string `mapstructure:"role_names"`
}
