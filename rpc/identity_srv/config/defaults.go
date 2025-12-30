package config

import (
	"time"

	"github.com/spf13/viper"
)

// setDefaults 设置默认配置值
func setDefaults(v *viper.Viper) {
	// 服务器配置默认值
	v.SetDefault("server.name", "identity-service")
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8891)
	v.SetDefault("server.debug", false) // 默认关闭，开发环境通过 .env 设置为 true

	// 健康检查配置默认值
	v.SetDefault("health_check.port", 10000)

	// 数据库默认配置
	v.SetDefault("database.driver", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.username", "postgres")
	v.SetDefault("database.password", "")
	v.SetDefault("database.dbname", "postgres")
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.max_open_conns", 20)
	v.SetDefault("database.conn_max_lifetime", 60*time.Minute)
	v.SetDefault("database.conn_max_idle_time", 5*time.Minute)
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.timezone", "UTC")

	// etcd配置默认值
	v.SetDefault("etcd.address", "localhost:2379")
	v.SetDefault("etcd.timeout", 5)

	// 日志配置默认值
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "file")
	v.SetDefault("log.file_path", "./logs/identity-srv.log")
	v.SetDefault("log.max_size", 100)
	v.SetDefault("log.max_age", 30)
	v.SetDefault("log.max_backups", 10)

	// 链路追踪配置默认值 (OpenTelemetry OTLP gRPC)
	v.SetDefault("tracing.enabled", false)
	v.SetDefault("tracing.endpoint", "localhost:4317")
	v.SetDefault("tracing.sampler_ratio", 0.1)

	// 组织Logo存储配置默认值
	v.SetDefault("logo_storage.s3_endpoint", "http://localhost:9000")
	v.SetDefault("logo_storage.s3_public_endpoint", "http://localhost:9000") // 默认与 s3_endpoint 相同
	v.SetDefault("logo_storage.s3_region", "us-east-1")
	v.SetDefault("logo_storage.s3_use_ssl", false)
	v.SetDefault("logo_storage.use_path_style", true) // 默认使用 Path Style 模式（适合开发环境）
	v.SetDefault("logo_storage.access_key", "")
	v.SetDefault("logo_storage.secret_key", "")
	v.SetDefault("logo_storage.s3_log_enabled", false)
	v.SetDefault("logo_storage.timeout", 5)
	v.SetDefault("logo_storage.max_file_size", 10485760) // 10MB
	// 允许的图片类型
	v.SetDefault("logo_storage.allowed_file_types", []string{
		"image/jpeg",    // JPEG 图片
		"image/png",     // PNG 图片
		"image/gif",     // GIF 图片
		"image/webp",    // WebP 图片
		"image/svg+xml", // SVG 图片
	})

	// 超级管理员配置默认值
	v.SetDefault("super_admin.role_names", []string{"superadmin"})
}
