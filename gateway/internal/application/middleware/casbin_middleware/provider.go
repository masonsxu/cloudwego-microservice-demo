package casbin_middleware

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"

	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
)

// Config Casbin 配置
type Config struct {
	// ModelPath 模型文件路径
	ModelPath string
	// Enabled 是否启用 Casbin 权限检查
	Enabled bool
	// LogEnabled 是否启用 Casbin 日志
	LogEnabled bool
	// SyncInterval 策略同步间隔（秒）
	SyncInterval int
	// SkipPaths 跳过权限检查的路径列表
	SkipPaths []string
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		ModelPath:    "./config/casbin_model.conf",
		Enabled:      true,
		LogEnabled:   false,
		SyncInterval: 300, // 5分钟同步一次
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/swagger",
			"/api/v1/identity/auth/login",
		},
	}
}

// LoadConfigFromEnv 从环境变量加载配置
func LoadConfigFromEnv() *Config {
	config := DefaultConfig()

	if modelPath := os.Getenv("CASBIN_MODEL_PATH"); modelPath != "" {
		config.ModelPath = modelPath
	}

	if enabled := os.Getenv("CASBIN_ENABLED"); enabled == "false" {
		config.Enabled = false
	}

	if logEnabled := os.Getenv("CASBIN_LOG_ENABLED"); logEnabled == "true" {
		config.LogEnabled = true
	}

	if skipPaths := os.Getenv("CASBIN_SKIP_PATHS"); skipPaths != "" {
		config.SkipPaths = splitAndTrim(skipPaths, ",")
	}

	return config
}

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

// ProvideCasbinConfig 提供 Casbin 配置
func ProvideCasbinConfig() *Config {
	return LoadConfigFromEnv()
}

// ProvideCasbinEnforcer 提供 Casbin Enforcer（使用内存 Adapter）
// 策略从 RPC 服务同步，不依赖数据库
func ProvideCasbinEnforcer(config *Config, logger *zerolog.Logger) (*CasbinEnforcer, error) {
	if !config.Enabled {
		logger.Info().Msg("Casbin is disabled")
		return nil, nil
	}

	// 检查模型文件是否存在
	absPath, err := filepath.Abs(config.ModelPath)
	if err != nil {
		logger.Warn().
			Err(err).
			Str("path", config.ModelPath).
			Msg("Failed to get absolute path for model, using embedded model")
		return NewCasbinEnforcer(logger)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		logger.Warn().
			Str("path", absPath).
			Msg("Casbin model file not found, using embedded model")
		return NewCasbinEnforcer(logger)
	}

	logger.Info().
		Str("model_path", absPath).
		Msg("Creating Casbin enforcer with model file")

	return NewCasbinEnforcerFromFile(absPath, logger)
}

// ProvideCasbinMiddleware 提供 Casbin 中间件
func ProvideCasbinMiddleware(config *Config, logger *zerolog.Logger) *CasbinMiddleware {
	if !config.Enabled {
		logger.Info().Msg("Casbin middleware disabled, using no-op middleware")
		return &CasbinMiddleware{
			enforcer:    nil,
			logger:      logger,
			skipPaths:   []string{"/"}, // 跳过所有路径
			pathMapping: make(map[string]string),
		}
	}

	enforcer, err := ProvideCasbinEnforcer(config, logger)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create Casbin enforcer, using no-op middleware")
		return &CasbinMiddleware{
			enforcer:    nil,
			logger:      logger,
			skipPaths:   []string{"/"},
			pathMapping: make(map[string]string),
		}
	}

	middleware := NewCasbinMiddleware(enforcer, logger)
	if len(config.SkipPaths) > 0 {
		middleware.SetSkipPaths(config.SkipPaths)
	}

	logger.Info().Strs("skip_paths", middleware.skipPaths).Msg("Casbin middleware created successfully")
	return middleware
}

// ProvideNoOpMiddleware 提供无操作的中间件（用于禁用 Casbin 时）
func ProvideNoOpMiddleware(logger *zerolog.Logger) *CasbinMiddleware {
	return &CasbinMiddleware{
		enforcer:    nil,
		logger:      logger,
		skipPaths:   []string{},
		pathMapping: make(map[string]string),
	}
}

// ProvidePolicySyncService 提供策略同步服务
func ProvidePolicySyncService(
	config *Config,
	middleware *CasbinMiddleware,
	identityClient identitycli.IdentityClient,
	logger *zerolog.Logger,
) *PolicySyncService {
	if middleware == nil || middleware.enforcer == nil {
		logger.Warn().Msg("Casbin middleware or enforcer is nil, policy sync service disabled")
		return nil
	}

	return NewPolicySyncService(middleware.enforcer, identityClient, logger, config.SyncInterval)
}
