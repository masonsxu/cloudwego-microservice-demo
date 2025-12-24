// Package wire 健康检查服务依赖注入提供者
package wire

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/config"
)

// HealthCheckSet 健康检查服务 Provider 集合
var HealthCheckSet = wire.NewSet(
	ProvideSQLDB,
	ProvideHealthCheckServer,
)

// ProvideSQLDB 从 GORM 提取底层的 *sql.DB
// 用于健康检查探测
func ProvideSQLDB(db *gorm.DB) (*sql.DB, error) {
	return db.DB()
}

// HealthCheckServer 健康检查 HTTP 服务器
type HealthCheckServer struct {
	server *http.Server
	db     *sql.DB
	port   int
}

// ProvideHealthCheckServer 提供健康检查服务器
func ProvideHealthCheckServer(cfg *config.Config, db *sql.DB) *HealthCheckServer {
	return &HealthCheckServer{
		db:   db,
		port: cfg.HealthCheck.Port,
	}
}

// Start 启动健康检查服务器（在独立的 goroutine 中运行）
func (h *HealthCheckServer) Start() {
	mux := http.NewServeMux()

	// /live 端点用于存活探测，确认进程正在运行
	mux.HandleFunc("/live", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	// /ready 端点用于就绪探测，确认依赖项是否健康
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := h.checkDependencies(); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	h.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", h.port),
		Handler: mux,
	}

	go func() {
		log.Printf("Health check server starting on port %d", h.port)

		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not start health check server: %v", err)
		}
	}()
}

// Stop 停止健康检查服务器
func (h *HealthCheckServer) Stop() error {
	if h.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return h.server.Shutdown(ctx)
	}

	return nil
}

// checkDependencies 运行所有依赖项检查
func (h *HealthCheckServer) checkDependencies() error {
	if err := h.checkDatabase(); err != nil {
		return fmt.Errorf("数据库检查失败: %w", err)
	}

	return nil
}

// checkDatabase 测试数据库连接以确保其可达
func (h *HealthCheckServer) checkDatabase() error {
	if h.db == nil {
		return fmt.Errorf("健康检查的数据库连接未初始化")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return h.db.PingContext(ctx)
}
