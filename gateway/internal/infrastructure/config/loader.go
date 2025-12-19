package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// loadConfig 内部配置加载函数
func loadConfig() (*Configuration, error) {
	v := viper.New()

	// 设置环境变量支持
	v.SetEnvPrefix("") // 不设置前缀，直接使用环境变量名
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	// 尝试加载 .env（可选，找不到不报错）
	loadDotEnvFirst([]string{".", "../..", "../../../"})

	// 设置默认值
	setDefaults(v)

	// 从环境变量映射到配置结构
	mapEnvVarsToConfig(v)

	// 解析配置
	var config Configuration
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("配置解析失败: %w", err)
	}

	// 后处理配置
	postProcessConfig(&config)

	return &config, nil
}

// postProcessConfig 配置后处理
func postProcessConfig(config *Configuration) {
	// 根据 SERVER_DEBUG 自动设置日志级别和错误详情开关
	// 这允许在生产环境临时开启 DEBUG 模式来排查问题
	if config.Server.Debug {
		// 调试模式下，如果日志级别未设置或为 info，自动调整为 debug
		if config.Log.Level == "" || config.Log.Level == "info" {
			config.Log.Level = "debug"
		}
		// 调试模式下，默认启用详细错误信息
		if !config.Middleware.ErrorHandler.EnableDetailedErrors {
			config.Middleware.ErrorHandler.EnableDetailedErrors = true
		}
	}

	// 确保日志目录存在
	if config.Log.Output == "file" && config.Log.FilePath != "" {
		// 使用 filepath.Dir 获取目录路径，而不是完整文件路径
		logDir := filepath.Dir(config.Log.FilePath)
		if err := os.MkdirAll(logDir, 0o755); err != nil {
			// 日志目录创建失败，降级到标准输出
			config.Log.Output = "stdout"
		}
	}

	// 初始化服务配置映射
	if config.Client.Services == nil {
		config.Client.Services = make(map[string]ServiceConfig)
	}
}
