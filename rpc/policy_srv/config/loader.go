package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig 加载应用配置
func LoadConfig() (*Config, error) {
	v := viper.New()

	loadDotEnvFile(v)
	setDefaults(v)
	mapEnvVarsToConfig(v)

	// 自动绑定所有环境变量（作为 fallback）
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// 后处理：组装地址
	cfg.Server.Address = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	// 传播 database 环境变量给子模块（GORM 可能需要）
	if cfg.Database.Password == "" {
		cfg.Database.Password = os.Getenv("DB_PASSWORD")
	}

	return &cfg, nil
}
