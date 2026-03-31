package casbin_middleware

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"

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
	// PathMapping API 路径到权限编码映射
	PathMapping map[string]string
	// SuperAdminBypassEnabled 是否启用超管兜底放行
	SuperAdminBypassEnabled bool
	// SuperAdminSubjects 超管主体列表（匹配 Casbin subject，例如 role:superadmin）
	SuperAdminSubjects []string
	// MenuMappingFile 菜单映射文件路径
	MenuMappingFile string
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

	if len(config.PathMapping) > 0 {
		for path, permCode := range config.PathMapping {
			middleware.AddPathMapping(path, permCode)
		}
	} else {
		mapping, err := LoadPathMappingFromMenuFile(config.MenuMappingFile)
		if err != nil {
			logger.Warn().Err(err).Str("menu_mapping_file", config.MenuMappingFile).
				Msg("Failed to load Casbin path mapping from menu file")
		} else {
			for path, permCode := range mapping {
				middleware.AddPathMapping(path, permCode)
			}
		}
	}

	middleware.SetSuperAdminBypassConfig(config.SuperAdminBypassEnabled, config.SuperAdminSubjects)

	logger.Info().
		Strs("skip_paths", middleware.skipPaths).
		Int("path_mapping_count", len(middleware.pathMapping)).
		Bool("superadmin_bypass_enabled", config.SuperAdminBypassEnabled).
		Strs("superadmin_subjects", config.SuperAdminSubjects).
		Msg("Casbin middleware created successfully")
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

type menuConfigFile struct {
	Menu []menuNode `yaml:"menu"`
}

type menuNode struct {
	ID       string     `yaml:"id"`
	APIPaths []string   `yaml:"api_paths"`
	Children []menuNode `yaml:"children"`
}

func LoadPathMappingFromMenuFile(menuFilePath string) (map[string]string, error) {
	if strings.TrimSpace(menuFilePath) == "" {
		return map[string]string{}, nil
	}

	resolvedPath := resolveMenuMappingPath(menuFilePath)

	content, err := os.ReadFile(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("read menu mapping file failed: %w", err)
	}

	var cfg menuConfigFile
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal menu mapping file failed: %w", err)
	}

	mapping := make(map[string]string)
	for _, node := range cfg.Menu {
		collectMenuMapping(node, mapping)
	}

	return mapping, nil
}

func resolveMenuMappingPath(configPath string) string {
	candidates := []string{configPath}
	if !filepath.IsAbs(configPath) {
		candidates = append(candidates,
			filepath.Join("..", configPath),
			filepath.Join("..", "..", configPath),
		)
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return configPath
}

func collectMenuMapping(node menuNode, mapping map[string]string) {
	permCode := ""
	if strings.TrimSpace(node.ID) != "" {
		permCode = "menu:" + strings.TrimSpace(node.ID)
	}

	if permCode != "" {
		for _, apiPath := range node.APIPaths {
			trimmedPath := strings.TrimSpace(apiPath)
			if trimmedPath == "" {
				continue
			}

			mapping[trimmedPath] = permCode
		}
	}

	for _, child := range node.Children {
		collectMenuMapping(child, mapping)
	}
}
