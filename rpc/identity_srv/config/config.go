package config

// LoadConfig 加载配置 - 对外暴露的统一入口
func LoadConfig() (*Config, error) {
	return loadConfig()
}
