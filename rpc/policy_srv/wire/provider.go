package wire

import (
	"database/sql"

	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/biz/logic"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/config"
)

// ProvideDB 提供数据库连接
func ProvideDB(cfg *config.Config, logger *zerolog.Logger) (*gorm.DB, error) {
	return config.InitDB(cfg, logger)
}

// ProvideLogger 提供 logger
func ProvideLogger(cfg *config.Config) (*zerolog.Logger, error) {
	return config.CreateLogger(cfg)
}

// ProvideEnforcerService 提供 Casbin enforcer 服务
func ProvideEnforcerService(db *gorm.DB, logger *zerolog.Logger) (*logic.EnforcerService, error) {
	return logic.NewEnforcerService(db, logger)
}

// ProvideDecisionService 提供决策服务
func ProvideDecisionService(enforcer *logic.EnforcerService) *logic.DecisionService {
	return logic.NewDecisionService(enforcer)
}

// ProvideSQLDB 从 GORM DB 提取 *sql.DB（健康检查用）
func ProvideSQLDB(db *gorm.DB) (*sql.DB, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	return sqlDB, nil
}
