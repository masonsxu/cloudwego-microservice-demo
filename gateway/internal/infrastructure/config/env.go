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

// splitAndTrim 分隔字符串并去除每个元素的首尾空格
// 这样可以避免配置中的空格导致匹配失败
// 例如: "/health, /metrics" -> ["/health", "/metrics"]
func splitAndTrim(value string, sep string) []string {
	parts := strings.Split(value, sep)

	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// mapEnvVarsToConfig 映射环境变量到配置结构
func mapEnvVarsToConfig(v *viper.Viper) {
	// 服务器配置映射
	mapServerEnvVars(v)

	// etcd配置映射
	mapEtcdEnvVars(v)

	// 客户端配置映射
	mapClientEnvVars(v)

	// 中间件配置映射
	mapMiddlewareEnvVars(v)

	// 日志配置映射
	mapLogEnvVars(v)

	// 链路追踪配置映射
	mapTracingEnvVars(v)

	// DataLake配置映射
	mapDataLakeEnvVars(v)

	// Redis配置映射
	mapRedisEnvVars(v)
}

// mapServerEnvVars 映射服务器相关环境变量
func mapServerEnvVars(v *viper.Viper) {
	mapToViper(v, "SERVER_NAME", "server.name", nil)
	mapToViper(v, "SERVER_HOST", "server.host", nil)
	mapToViper(v, "SERVER_PORT", "server.port", func(value string) interface{} {
		if val, err := strconv.Atoi(value); err == nil {
			return val
		}

		return 8888
	})
	mapToViper(v, "SERVER_DEBUG", "server.debug", func(value string) interface{} {
		return value == "true"
	})
	mapToViper(v, "SERVER_READ_TIMEOUT", "server.read_timeout", func(value string) interface{} {
		return parseDurationWithDefault(value, 30*time.Second)
	})
	mapToViper(v, "SERVER_WRITE_TIMEOUT", "server.write_timeout", func(value string) interface{} {
		return parseDurationWithDefault(value, 30*time.Second)
	})
	mapToViper(v, "SERVER_IDLE_TIMEOUT", "server.idle_timeout", func(value string) interface{} {
		return parseDurationWithDefault(value, 120*time.Second)
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

// mapClientEnvVars 映射客户端相关环境变量
func mapClientEnvVars(v *viper.Viper) {
	mapToViper(
		v,
		"CLIENT_CONNECTION_TIMEOUT",
		"client.connection_timeout",
		func(value string) interface{} {
			return parseDurationWithDefault(value, 2*time.Second)
		},
	)
	mapToViper(
		v,
		"CLIENT_REQUEST_TIMEOUT",
		"client.request_timeout",
		func(value string) interface{} {
			return parseDurationWithDefault(value, 60*time.Second)
		},
	)

	// 连接池配置映射
	mapToViper(
		v,
		"CLIENT_POOL_MAX_IDLE_PER_ADDRESS",
		"client.pool.max_idle_per_address",
		func(value string) interface{} {
			if val, err := strconv.Atoi(value); err == nil {
				return val
			}

			return 10
		},
	)
	mapToViper(
		v,
		"CLIENT_POOL_MAX_IDLE_GLOBAL",
		"client.pool.max_idle_global",
		func(value string) interface{} {
			if val, err := strconv.Atoi(value); err == nil {
				return val
			}

			return 100
		},
	)
	mapToViper(
		v,
		"CLIENT_POOL_MAX_IDLE_TIMEOUT",
		"client.pool.max_idle_timeout",
		func(value string) interface{} {
			return parseDurationWithDefault(value, 5*time.Minute)
		},
	)

	// 服务配置映射
	mapToViper(v, "IDENTITY_SRV_NAME", "client.services.identity.name", nil)
}

// mapMiddlewareEnvVars 映射中间件相关环境变量
func mapMiddlewareEnvVars(v *viper.Viper) {
	// CORS配置映射
	mapCORSEnvVars(v)

	// 限流配置映射
	mapRateLimitEnvVars(v)

	// 身份验证配置映射
	mapJWTEnvVars(v)

	// 错误处理中间件配置映射
	mapErrorHandlerEnvVars(v)
}

// mapCORSEnvVars 映射CORS相关环境变量
func mapCORSEnvVars(v *viper.Viper) {
	mapToViper(v, "CORS_ENABLED", "middleware.cors.enabled", func(value string) interface{} {
		return value == "true"
	})
	mapToViper(
		v,
		"CORS_ALLOW_ORIGINS",
		"middleware.cors.allow_origins",
		func(value string) interface{} {
			return splitAndTrim(value, ",")
		},
	)
	mapToViper(
		v,
		"CORS_ALLOW_CREDENTIALS",
		"middleware.cors.allow_credentials",
		func(value string) interface{} {
			return value == "true"
		},
	)
}

// mapRateLimitEnvVars 映射限流相关环境变量
func mapRateLimitEnvVars(v *viper.Viper) {
	mapToViper(
		v,
		"RATE_LIMIT_ENABLED",
		"middleware.rate_limit.enabled",
		func(value string) interface{} {
			return value == "true"
		},
	)
	mapToViper(
		v,
		"RATE_LIMIT_RPS",
		"middleware.rate_limit.requests_per_second",
		func(value string) interface{} {
			if val, err := strconv.Atoi(value); err == nil {
				return val
			}

			return 1000
		},
	)
	mapToViper(
		v,
		"RATE_LIMIT_BURST",
		"middleware.rate_limit.burst",
		func(value string) interface{} {
			if val, err := strconv.Atoi(value); err == nil {
				return val
			}

			return 2000
		},
	)
}

// mapJWTEnvVars 映射身份验证相关环境变量
func mapJWTEnvVars(v *viper.Viper) {
	// JWT配置映射
	mapToViper(v, "JWT_ENABLED", "middleware.jwt.enabled", func(value string) interface{} {
		return value == "true"
	})
	mapToViper(v, "JWT_SIGNING_KEY", "middleware.jwt.signing_key", nil)
	mapToViper(v, "JWT_TIMEOUT", "middleware.jwt.timeout", func(value string) interface{} {
		return parseDurationWithDefault(value, 30*time.Minute)
	})
	mapToViper(v, "JWT_MAX_REFRESH", "middleware.jwt.max_refresh", func(value string) interface{} {
		return parseDurationWithDefault(value, 7*24*time.Hour)
	})
	mapToViper(v, "JWT_IDENTITY_KEY", "middleware.jwt.identity_key", nil)
	mapToViper(v, "JWT_REALM", "middleware.jwt.realm", nil)
	mapToViper(v, "JWT_TOKEN_LOOKUP", "middleware.jwt.token_lookup", nil)
	mapToViper(v, "JWT_TOKEN_HEAD_NAME", "middleware.jwt.token_head_name", nil)
	mapToViper(
		v,
		"JWT_SEND_AUTHORIZATION",
		"middleware.jwt.send_authorization",
		func(value string) interface{} {
			return value == "false" // 默认不发送 Authorization header
		},
	)
	mapToViper(v, "JWT_SKIP_PATHS", "middleware.jwt.skip_paths", func(value string) interface{} {
		return splitAndTrim(value, ",")
	})

	// Cookie配置映射
	mapCookieEnvVars(v)
}

// mapCookieEnvVars 映射Cookie相关环境变量
func mapCookieEnvVars(v *viper.Viper) {
	mapToViper(
		v,
		"JWT_COOKIE_SEND_COOKIE",
		"middleware.jwt.cookie.send_cookie",
		func(value string) interface{} {
			return value == "true"
		},
	)
	mapToViper(v, "JWT_COOKIE_COOKIE_NAME", "middleware.jwt.cookie.cookie_name", nil)
	mapToViper(v, "JWT_COOKIE_COOKIE_DOMAIN", "middleware.jwt.cookie.cookie_domain", nil)
	mapToViper(v, "JWT_COOKIE_COOKIE_PATH", "middleware.jwt.cookie.cookie_path", nil)
	mapToViper(
		v,
		"JWT_COOKIE_COOKIE_MAX_AGE",
		"middleware.jwt.cookie.cookie_max_age",
		func(value string) interface{} {
			return parseDurationWithDefault(value, 7*24*time.Hour)
		},
	)
	mapToViper(v, "JWT_COOKIE_COOKIE_SAME_SITE", "middleware.jwt.cookie.cookie_same_site", nil)
	mapToViper(
		v,
		"JWT_COOKIE_SECURE_COOKIE",
		"middleware.jwt.cookie.secure_cookie",
		func(value string) interface{} {
			return value == "true"
		},
	)
	mapToViper(
		v,
		"JWT_COOKIE_HTTP_ONLY",
		"middleware.jwt.cookie.cookie_http_only",
		func(value string) interface{} {
			return value == "true"
		},
	)
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

		return 100 // 默认100MB
	})
	mapToViper(v, "LOG_MAX_AGE", "log.max_age", func(value string) interface{} {
		if val, err := strconv.Atoi(value); err == nil {
			return val
		}

		return 30 // 默认30天
	})
	mapToViper(v, "LOG_MAX_BACKUPS", "log.max_backups", func(value string) interface{} {
		if val, err := strconv.Atoi(value); err == nil {
			return val
		}

		return 10 // 默认10个备份
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

		return 0.1 // 默认采样率0.1
	})
	mapToViper(v, "TRACING_IGNORE_PATHS", "tracing.ignore_paths", func(value string) interface{} {
		return splitAndTrim(value, ",")
	})
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

// mapErrorHandlerEnvVars 映射错误处理中间件相关环境变量
func mapErrorHandlerEnvVars(v *viper.Viper) {
	mapToViper(
		v,
		"ERROR_HANDLER_ENABLED",
		"middleware.error_handler.enabled",
		func(value string) interface{} { return value == "true" },
	)
	mapToViper(
		v,
		"ERROR_HANDLER_ENABLE_DETAILED_ERRORS",
		"middleware.error_handler.enable_detailed_errors",
		func(value string) interface{} { return value == "true" },
	)
	mapToViper(
		v,
		"ERROR_HANDLER_ENABLE_REQUEST_LOGGING",
		"middleware.error_handler.enable_request_logging",
		func(value string) interface{} { return value == "true" },
	)
	mapToViper(
		v,
		"ERROR_HANDLER_ENABLE_RESPONSE_LOGGING",
		"middleware.error_handler.enable_response_logging",
		func(value string) interface{} { return value == "true" },
	)
	mapToViper(
		v,
		"ERROR_HANDLER_ENABLE_PANIC_RECOVERY",
		"middleware.error_handler.enable_panic_recovery",
		func(value string) interface{} { return value == "true" },
	)
	mapToViper(
		v,
		"ERROR_HANDLER_MAX_STACK_TRACE_SIZE",
		"middleware.error_handler.max_stack_trace_size",
		func(value string) interface{} {
			if val, err := strconv.Atoi(value); err == nil {
				return val
			}

			return 4096 // 默认4KB
		},
	)
	mapToViper(
		v,
		"ERROR_HANDLER_ENABLE_ERROR_METRICS",
		"middleware.error_handler.enable_error_metrics",
		func(value string) interface{} { return value == "true" },
	)
	mapToViper(
		v,
		"ERROR_HANDLER_ERROR_RESPONSE_TIMEOUT",
		"middleware.error_handler.error_response_timeout",
		func(value string) interface{} {
			if val, err := strconv.Atoi(value); err == nil {
				return val
			}

			return 5000 // 默认5000毫秒
		},
	)
}

// mapDataLakeEnvVars 映射 DataLake 相关环境变量
func mapDataLakeEnvVars(v *viper.Viper) {
	mapToViper(v, "DATALAKE_URL", "data_lake.data_lake_url", nil)
}

// mapRedisEnvVars 映射 Redis 相关环境变量
func mapRedisEnvVars(v *viper.Viper) {
	mapToViper(v, "REDIS_ADDRESS", "redis.address", nil)
	mapToViper(v, "REDIS_PASSWORD", "redis.password", nil)
	mapToViper(v, "REDIS_DB", "redis.db", func(value string) interface{} {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}

		return 0
	})
	mapToViper(v, "REDIS_POOL_SIZE", "redis.pool_size", func(value string) interface{} {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}

		return 10
	})
	mapToViper(v, "REDIS_MIN_IDLE_CONNS", "redis.min_idle_conns", func(value string) interface{} {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}

		return 5
	})
	mapToViper(v, "REDIS_MAX_RETRIES", "redis.max_retries", func(value string) interface{} {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}

		return 3
	})
	mapToViper(v, "REDIS_DIAL_TIMEOUT", "redis.dial_timeout", func(value string) interface{} {
		return parseDurationWithDefault(value, 5*time.Second)
	})
	mapToViper(v, "REDIS_READ_TIMEOUT", "redis.read_timeout", func(value string) interface{} {
		return parseDurationWithDefault(value, 3*time.Second)
	})
	mapToViper(v, "REDIS_WRITE_TIMEOUT", "redis.write_timeout", func(value string) interface{} {
		return parseDurationWithDefault(value, 3*time.Second)
	})
	mapToViper(v, "REDIS_POOL_TIMEOUT", "redis.pool_timeout", func(value string) interface{} {
		return parseDurationWithDefault(value, 4*time.Second)
	})
	mapToViper(v, "REDIS_IDLE_TIMEOUT", "redis.idle_timeout", func(value string) interface{} {
		return parseDurationWithDefault(value, 5*time.Minute)
	})
	mapToViper(v, "REDIS_IDLE_CHECK_FREQ", "redis.idle_check_freq", func(value string) interface{} {
		return parseDurationWithDefault(value, 1*time.Minute)
	})
}
