package wire

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/biz/logic"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/config"
)

// AppContainer 应用依赖容器
type AppContainer struct {
	Config        *config.Config
	Logger        *zerolog.Logger
	DB            *gorm.DB
	Decision      *logic.DecisionService
	Enforcer      *logic.EnforcerService
	ServerOptions *ServerOptions

	healthSrv *http.Server
}

// NewAppContainer 创建应用容器
func NewAppContainer(
	cfg *config.Config,
	logger *zerolog.Logger,
	db *gorm.DB,
	decision *logic.DecisionService,
	enforcer *logic.EnforcerService,
	serverOpts *ServerOptions,
	sqlDB *sql.DB,
) (*AppContainer, error) {
	healthSrv := createHealthServer(cfg, logger, sqlDB)

	return &AppContainer{
		Config:        cfg,
		Logger:        logger,
		DB:            db,
		Decision:      decision,
		Enforcer:      enforcer,
		ServerOptions: serverOpts,
		healthSrv:     healthSrv,
	}, nil
}

// StartHealthCheck 启动健康检查 HTTP 服务
func (c *AppContainer) StartHealthCheck() {
	go func() {
		c.Logger.Info().Str("addr", c.healthSrv.Addr).Msg("健康检查服务启动")
		if err := c.healthSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			c.Logger.Error().Err(err).Msg("健康检查服务错误")
		}
	}()
}

func createHealthServer(cfg *config.Config, logger *zerolog.Logger, sqlDB *sql.DB) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := sqlDB.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("DB ping failed"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	addr := fmt.Sprintf("%s:%d", cfg.HealthCheck.Host, cfg.HealthCheck.Port)
	return &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}
