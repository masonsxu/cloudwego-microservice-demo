package config

import (
	"os"

	"github.com/spf13/viper"
)

func mapEnvVarsToConfig(v *viper.Viper) {
	// Database
	_ = v.BindEnv("database.host", "DB_HOST")
	_ = v.BindEnv("database.port", "DB_PORT")
	_ = v.BindEnv("database.user", "DB_USER")
	_ = v.BindEnv("database.password", "DB_PASSWORD")
	_ = v.BindEnv("database.dbname", "DB_NAME")
	_ = v.BindEnv("database.sslmode", "DB_SSLMODE")

	// Server
	_ = v.BindEnv("server.port", "SERVER_PORT")

	// etcd
	_ = v.BindEnv("etcd.address", "ETCD_ADDRESS")

	// Tracing
	_ = v.BindEnv("tracing.enabled", "TRACING_ENABLED")
	_ = v.BindEnv("tracing.endpoint", "TRACING_ENDPOINT")
}

func loadDotEnvFile(v *viper.Viper) {
	searchPaths := []string{".", "../..", "../../../"}
	for _, p := range searchPaths {
		envPath := p + "/.env"
		if _, err := os.Stat(envPath); err == nil {
			v.SetConfigFile(envPath)
			_ = v.ReadInConfig()
			return
		}
	}
}

// parseDurationWithDefault 和 parseBool 预留：
// policy_srv 当前通过 viper 自动 Unmarshal 处理类型转换，无需额外解析
