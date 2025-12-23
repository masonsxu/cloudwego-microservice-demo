package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/config"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/internal/middleware"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/internal/otel"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv/identityservice"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/wire"
)

// dbForHealthCheck 用于健康检查的数据库连接
// 在 main 函数中通过 Wire 依赖注入初始化后赋值
var dbForHealthCheck *sql.DB

// runHealthCheckServer 启动独立的 HTTP 健康检查服务器
func runHealthCheckServer(port int) {
	mux := http.NewServeMux()

	// /live 端点用于存活探测，确认进程正在运行
	mux.HandleFunc("/live", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	// /ready 端点用于就绪探测，确认依赖项是否健康
	mux.HandleFunc("/ready", func(w http.ResponseWriter, _ *http.Request) {
		// 检查依赖项
		err := checkDependencies()
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	log.Printf("Health check server starting on port %d", port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not start health check server: %v", err)
	}
}

// checkDependencies 运行所有依赖项检查
func checkDependencies() error {
	// 检查数据库连接
	if err := checkDatabase(dbForHealthCheck); err != nil {
		return fmt.Errorf("数据库检查失败: %w", err)
	}

	return nil
}

// checkDatabase 测试数据库连接以确保其可达
func checkDatabase(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("健康检查的数据库连接未初始化")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		klog.Fatalf("failed to load config: %v", err)
	}

	// 在独立的 goroutine 中启动健康检查服务器
	// 使用不同的端口进行健康检查是一个最佳实践
	go runHealthCheckServer(cfg.HealthCheck.Port)

	// 2. 初始化 OpenTelemetry Provider
	shutdown, err := otel.NewProvider(cfg)
	if err != nil {
		klog.Fatalf("failed to init OpenTelemetry provider: %v", err)
	}

	defer func() {
		if err := shutdown(context.Background()); err != nil {
			klog.Fatalf("failed to shutdown OpenTelemetry provider: %v", err)
		}
	}()

	// 3. 创建 handler 实例并获取数据库连接
	serviceImpl, serviceWithDB, err := NewIdentityServiceImplWithDB()
	if err != nil {
		klog.Fatalf("failed to create service impl: %v", err)
	}

	// 从 GORM 获取底层的 *sql.DB 用于健康检查
	sqlDB, err := serviceWithDB.DB.DB()
	if err != nil {
		klog.Fatalf("failed to get sql.DB from gorm: %v", err)
	}

	dbForHealthCheck = sqlDB

	// 4. 配置并启动服务器
	// 解析监听地址
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		klog.Fatalf("failed to resolve server address: %v", err)
	}

	// 服务注册名称（用于服务发现）
	serviceName := cfg.Server.Name

	// 构建 Etcd 注册中心实例（用于服务注册与发现）
	r, err := etcd.NewEtcdRegistry([]string{cfg.Etcd.Address})
	if err != nil {
		klog.Fatal(err)
	}

	// 初始化 logger
	logger, err := wire.InitializeLogger()
	if err != nil {
		klog.Fatalf("failed to initialize logger: %v", err)
	}
	// 初始化 Kitex logger（用于 Kitex 框架日志）
	kitexLogger, err := config.CreateKitexLogger(cfg)
	if err != nil {
		klog.Fatalf("failed to create Kitex logger: %v", err)
	}

	// 设置 Kitex 全局 logger
	klog.SetLogger(kitexLogger)

	// 根据配置设置日志级别
	switch cfg.Log.Level {
	case "debug":
		klog.SetLevel(klog.LevelDebug)
	case "info":
		klog.SetLevel(klog.LevelInfo)
	case "warn":
		klog.SetLevel(klog.LevelWarn)
	case "error":
		klog.SetLevel(klog.LevelError)
	default:
		klog.SetLevel(klog.LevelInfo)
	}

	// 创建MetaInfo中间件
	metaMiddleware := middleware.NewMetaInfoMiddleware(logger)

	// 构建 Server Options
	serverOpts := []server.Option{
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
		server.WithRegistry(r),
		server.WithServiceAddr(addr),
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler),
	}

	// 6. 如果启用追踪，添加 OpenTelemetry Suite（必须在 MetaInfoMiddleware 之前）
	// 这样 OTel 会先创建 span，MetaInfoMiddleware 才能提取 trace_id/span_id
	if cfg.Tracing.Enabled {
		serverOpts = append(serverOpts, server.WithSuite(tracing.NewServerSuite()))
		klog.Infof("OpenTelemetry tracing enabled, endpoint: %s", cfg.Tracing.Endpoint)
	}

	// 7. 添加 MetaInfoMiddleware（在 OTel Suite 之后，确保 span 已创建）
	serverOpts = append(serverOpts, server.WithMiddleware(metaMiddleware.ServerMiddleware()))

	// 创建并配置 Kitex Server
	svr := identityservice.NewServer(serviceImpl, serverOpts...)

	klog.Infof("Identity service starting on %s", addr.String())

	if err := svr.Run(); err != nil {
		klog.Fatalf("server stopped with error: %v", err)
	}
}
