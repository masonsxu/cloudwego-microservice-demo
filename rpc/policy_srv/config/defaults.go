package config

import (
	"time"

	"github.com/spf13/viper"
)

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.name", "policy-service")
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8892)

	v.SetDefault("health_check.enabled", true)
	v.SetDefault("health_check.host", "0.0.0.0")
	v.SetDefault("health_check.port", 10001)

	v.SetDefault("etcd.address", "localhost:2379")
	v.SetDefault("etcd.timeout", 5*time.Second)

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "file")
	v.SetDefault("log.file_path", "./logs/policy-srv.log")
	v.SetDefault("log.max_size", 100)
	v.SetDefault("log.max_age", 30)
	v.SetDefault("log.max_backups", 10)

	v.SetDefault("tracing.enabled", false)
	v.SetDefault("tracing.endpoint", "jaeger:4317")
	v.SetDefault("tracing.sampler_ratio", 0.1)

	v.SetDefault("database.driver", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "Admin")
	v.SetDefault("database.password", "")
	v.SetDefault("database.dbname", "identity_srv")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", 5*time.Minute)
}
