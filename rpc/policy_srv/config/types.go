package config

import "time"

// Config 应用配置根结构
type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	HealthCheck HealthCheckConfig `mapstructure:"health_check"`
	Etcd        EtcdConfig        `mapstructure:"etcd"`
	Log         LogConfig         `mapstructure:"log"`
	Tracing     TracingConfig     `mapstructure:"tracing"`
	Database    DatabaseConfig    `mapstructure:"database"`
}

// ServerConfig Kitex RPC 服务器配置
type ServerConfig struct {
	Name         string `mapstructure:"name"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Address      string // Host:Port，由 LoadConfig 后处理填充
}

// HealthCheckConfig 健康检查 HTTP 服务器配置
type HealthCheckConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
}

// EtcdConfig etcd 服务发现配置
type EtcdConfig struct {
	Address string        `mapstructure:"address"`
	Timeout time.Duration `mapstructure:"timeout"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// TracingConfig OpenTelemetry 配置
type TracingConfig struct {
	Enabled      bool    `mapstructure:"enabled"`
	Endpoint     string  `mapstructure:"endpoint"`
	SamplerRatio float64 `mapstructure:"sampler_ratio"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}
