package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// mapToViper 通用环境变量映射函数
func mapToViper(v *viper.Viper, envVar, configKey string, transformer func(string) interface{}) {
	if value := os.Getenv(envVar); value != "" {
		if transformer != nil {
			v.Set(configKey, transformer(value))
		} else {
			v.Set(configKey, value)
		}
	}
}

// parseDurationWithDefault 解析持续时间，支持多种格式
func parseDurationWithDefault(value string, defaultValue time.Duration) time.Duration {
	// 如果是纯数字，视为秒数
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}

	// 尝试解析为 time.Duration
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}

	// 解析失败，返回默认值
	return defaultValue
}

// mapEnvVarsToConfig 映射环境变量到配置结构
func mapEnvVarsToConfig(v *viper.Viper) {
	// 数据库配置映射
	mapDatabaseEnvVars(v)

	// 服务器配置映射
	mapServerEnvVars(v)

	// 健康检查配置映射
	mapHealthCheckEnvVars(v)

	// etcd配置映射
	mapEtcdEnvVars(v)

	// 日志配置映射
	mapLogEnvVars(v)

	// 链路追踪配置映射
	mapTracingEnvVars(v)

	// 管理员配置映射
	mapAdminEnvVars(v)

	// Logo存储配置映射
	mapLogoStorageEnvVars(v)
}

// mapDatabaseEnvVars 映射数据库相关环境变量
func mapDatabaseEnvVars(v *viper.Viper) {
	mapToViper(v, "DB_DRIVER", "database.driver", nil)
	mapToViper(v, "DB_HOST", "database.host", nil)
	mapToViper(v, "DB_PORT", "database.port", nil)
	mapToViper(v, "DB_USERNAME", "database.username", nil)
	mapToViper(v, "DB_PASSWORD", "database.password", nil)
	mapToViper(v, "DB_NAME", "database.dbname", nil)
	mapToViper(v, "DB_SSLMODE", "database.sslmode", nil)
	mapToViper(v, "DB_TIMEZONE", "database.timezone", nil)

	// 连接池配置
	mapToViper(v, "DB_MAX_IDLE_CONNS", "database.max_idle_conns", nil)
	mapToViper(v, "DB_MAX_OPEN_CONNS", "database.max_open_conns", nil)
	mapToViper(
		v,
		"DB_CONN_MAX_LIFETIME",
		"database.conn_max_lifetime",
		func(value string) interface{} {
			return parseDurationWithDefault(value, 60*time.Minute)
		},
	)
	mapToViper(
		v,
		"DB_CONN_MAX_IDLE_TIME",
		"database.conn_max_idle_time",
		func(value string) interface{} {
			return parseDurationWithDefault(value, 5*time.Minute)
		},
	)
}

// mapServerEnvVars 映射服务器相关环境变量
func mapServerEnvVars(v *viper.Viper) {
	mapToViper(v, "SERVER_ADDRESS", "server.address", nil)
	mapToViper(v, "SERVER_NAME", "server.name", nil)
	mapToViper(v, "SERVER_HOST", "server.host", nil)
	mapToViper(v, "SERVER_PORT", "server.port", func(value string) interface{} {
		if val, err := strconv.Atoi(value); err == nil {
			return val
		}

		return 8891
	})
	mapToViper(v, "SERVER_DEBUG", "server.debug", func(value string) interface{} {
		return value == "true"
	})
}

// mapHealthCheckEnvVars 映射健康检查相关环境变量
func mapHealthCheckEnvVars(v *viper.Viper) {
	mapToViper(v, "HEALTH_CHECK_PORT", "health_check.port", func(value string) interface{} {
		if val, err := strconv.Atoi(value); err == nil {
			return val
		}

		return 10000
	})
}

// mapEtcdEnvVars 映射etcd相关环境变量
func mapEtcdEnvVars(v *viper.Viper) {
	mapToViper(v, "ETCD_ADDRESS", "etcd.address", nil)
	mapToViper(v, "ETCD_USERNAME", "etcd.username", nil)
	mapToViper(v, "ETCD_PASSWORD", "etcd.password", nil)
	mapToViper(v, "ETCD_TIMEOUT", "etcd.timeout", func(value string) interface{} {
		return parseDurationWithDefault(value, 5*time.Second)
	})
}

// mapLogEnvVars 映射日志相关环境变量
func mapLogEnvVars(v *viper.Viper) {
	mapToViper(v, "LOG_LEVEL", "log.level", nil)
	mapToViper(v, "LOG_FORMAT", "log.format", nil)
	mapToViper(v, "LOG_OUTPUT", "log.output", nil)
	mapToViper(v, "LOG_FILE_PATH", "log.file_path", nil)
	mapToViper(v, "LOG_MAX_SIZE", "log.max_size", func(value string) interface{} {
		if val, err := strconv.Atoi(value); err == nil {
			return val
		}

		return 100
	})
	mapToViper(v, "LOG_MAX_AGE", "log.max_age", func(value string) interface{} {
		if val, err := strconv.Atoi(value); err == nil {
			return val
		}

		return 30
	})
	mapToViper(v, "LOG_MAX_BACKUPS", "log.max_backups", func(value string) interface{} {
		if val, err := strconv.Atoi(value); err == nil {
			return val
		}

		return 10
	})
}

// mapTracingEnvVars 映射链路追踪相关环境变量
func mapTracingEnvVars(v *viper.Viper) {
	mapToViper(v, "TRACING_ENABLED", "tracing.enabled", func(value string) interface{} {
		return value == "true"
	})
	mapToViper(v, "TRACING_ENDPOINT", "tracing.endpoint", nil)
	mapToViper(v, "TRACING_SAMPLER_RATIO", "tracing.sampler_ratio", func(value string) interface{} {
		if val, err := strconv.ParseFloat(value, 64); err == nil {
			return val
		}

		return 0.1
	})
}

// mapAdminEnvVars 映射管理员相关环境变量
func mapAdminEnvVars(v *viper.Viper) {
	mapToViper(v, "ADMIN_USERNAME", "admin.username", nil)
	mapToViper(v, "ADMIN_PASSWORD", "admin.password", nil)
}

// mapLogoStorageEnvVars 映射组织Logo存储相关环境变量
func mapLogoStorageEnvVars(v *viper.Viper) {
	mapToViper(v, "LOGO_STORAGE_S3_ENDPOINT", "logo_storage.s3_endpoint", nil)
	mapToViper(v, "LOGO_STORAGE_S3_PUBLIC_ENDPOINT", "logo_storage.s3_public_endpoint", nil)
	mapToViper(v, "LOGO_STORAGE_S3_REGION", "logo_storage.s3_region", nil)
	mapToViper(
		v,
		"LOGO_STORAGE_S3_USE_SSL",
		"logo_storage.s3_use_ssl",
		func(value string) interface{} {
			return value == "true"
		},
	)
	mapToViper(
		v,
		"LOGO_STORAGE_USE_PATH_STYLE",
		"logo_storage.use_path_style",
		func(value string) interface{} {
			return value == "true"
		},
	)
	mapToViper(v, "LOGO_STORAGE_ACCESS_KEY", "logo_storage.access_key", nil)
	mapToViper(v, "LOGO_STORAGE_SECRET_KEY", "logo_storage.secret_key", nil)
	mapToViper(
		v,
		"LOGO_STORAGE_S3_LOG_ENABLED",
		"logo_storage.s3_log_enabled",
		func(value string) interface{} {
			return value == "true"
		},
	)
	mapToViper(v, "LOGO_STORAGE_TIMEOUT", "logo_storage.timeout", func(value string) interface{} {
		if val, err := strconv.Atoi(value); err == nil {
			return val
		}

		return 5
	})
	mapToViper(
		v,
		"LOGO_STORAGE_MAX_FILE_SIZE",
		"logo_storage.max_file_size",
		func(value string) interface{} {
			if val, err := strconv.Atoi(value); err == nil {
				return val
			}

			return 10485760 // 10MB
		},
	)
	mapToViper(
		v,
		"LOGO_STORAGE_ALLOWED_FILE_TYPES",
		"logo_storage.allowed_file_types",
		func(value string) interface{} {
			// 分割 MIME 类型列表，并去除每个类型的空格
			types := strings.Split(value, ",")

			result := make([]string, 0, len(types))
			for _, t := range types {
				trimmed := strings.TrimSpace(t)
				if trimmed != "" {
					result = append(result, trimmed)
				}
			}

			return result
		},
	)
}

// loadDotEnvFirst 在给定路径列表中查找首个 .env 并加载到环境变量（若未找到则忽略）。
func loadDotEnvFirst(paths []string) {
	for _, p := range paths {
		fp := filepath.Join(p, ".env")
		if _, err := os.Stat(fp); err == nil {
			_ = loadDotEnvFile(fp)
			return
		}
	}
}

// loadDotEnvFile 读取 .env 文件并将未设置的键注入到进程环境变量。
func loadDotEnvFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])

		value := strings.TrimSpace(parts[1])
		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}

	return nil
}
