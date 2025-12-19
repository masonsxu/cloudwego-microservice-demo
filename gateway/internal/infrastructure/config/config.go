package config

var Config *Configuration

// LoadConfig 加载配置 - 对外暴露的统一入口
func LoadConfig() (*Configuration, error) {
	return loadConfig()
}

// Init 初始化配置（保持向后兼容）
func Init() error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	Config = config

	return nil
}
