package wire

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/iamclient"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/config"
)

// =============================================================================
// Provider Functions - 具体的依赖提供者实现
// =============================================================================

// ProvideDB 提供数据库连接实例
// Wire 依赖注入提供者，委托给 config 层处理所有初始化逻辑
func ProvideDB(cfg *config.Config, logger *zerolog.Logger) (*gorm.DB, error) {
	return config.InitDB(cfg, logger)
}

// ProvideLogger 提供结构化日志实例
// 根据配置提供不同环境的日志配置
func ProvideLogger(cfg *config.Config) (*zerolog.Logger, error) {
	logger, err := config.CreateLogger(cfg)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// =============================================================================
// Provider Options - 高级配置选项
// =============================================================================

// DBOption 数据库配置选项
type DBOption func(*gorm.DB) error

// ProvideDBWithOptions 提供带自定义选项的数据库连接
func ProvideDBWithOptions(
	cfg *config.Config,
	logger *zerolog.Logger,
	opts ...DBOption,
) (*gorm.DB, error) {
	db, err := ProvideDB(cfg, logger)
	if err != nil {
		return nil, err
	}

	// 应用自定义选项
	for _, opt := range opts {
		if err := opt(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

// WithDBDebugMode 启用数据库调试模式
func WithDBDebugMode() DBOption {
	return func(db *gorm.DB) error {
		return nil // db.Debug() 返回的是新实例，这里需要根据实际情况调整
	}
}

// WithDBMigration 执行数据库迁移
func WithDBMigration(models ...interface{}) DBOption {
	return func(db *gorm.DB) error {
		return db.AutoMigrate(models...)
	}
}

// ProvideLoggerWithOptions 提供带自定义选项的日志器
// 注意：zerolog 使用不同的配置方式，此函数保留以保持兼容性
// 实际配置通过 config.CreateLogger 处理
func ProvideLoggerWithOptions(cfg *config.Config) (*zerolog.Logger, error) {
	return ProvideLogger(cfg)
}

// ProvideIAMClient 提供 IAM 客户端（PDP 决策入口）
//
// identity_srv 自身作为业务系统接入 iamclient，吃自家狗粮：
// 管理动作 RPC 入口通过 SubjectFromContext + MustCheck 调 policy_srv 决策。
//
// 失败时返回错误，让进程在启动期暴露问题（etcd 不可达、policy_srv 注册名错配等）。
func ProvideIAMClient(cfg *config.Config) (*iamclient.Client, func(), error) {
	cli, err := iamclient.New(iamclient.Config{
		EtcdEndpoints: []string{cfg.Etcd.Address},
		CallerService: cfg.Server.Name,
	})
	if err != nil {
		return nil, nil, err
	}

	return cli, func() { _ = cli.Close() }, nil
}
