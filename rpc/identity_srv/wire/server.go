// Package wire 服务器层依赖注入提供者
package wire

import (
	"fmt"
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/google/wire"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/rs/zerolog"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/config"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/internal/middleware"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/internal/otel"
)

// ServerSet 服务器层依赖注入集合
var ServerSet = wire.NewSet(
	ProvideOtelProvider,
	ProvideEtcdRegistry,
	ProvideKitexLogger,
	ProvideMetaMiddleware,
	ProvideServerOptions,
)

// ProvideOtelProvider 提供 OpenTelemetry Provider
func ProvideOtelProvider(cfg *config.Config) (*otel.Provider, func(), error) {
	return otel.NewProvider(cfg)
}

// ProvideEtcdRegistry 提供 etcd 注册中心
func ProvideEtcdRegistry(cfg *config.Config) (registry.Registry, error) {
	return etcd.NewEtcdRegistry([]string{cfg.Etcd.Address})
}

// ProvideKitexLogger 提供 Kitex 日志配置
// 初始化 Kitex 全局日志器
func ProvideKitexLogger(cfg *config.Config) error {
	kitexLogger, err := config.CreateKitexLogger(cfg)
	if err != nil {
		return fmt.Errorf("failed to create Kitex logger: %w", err)
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

	return nil
}

// ProvideMetaMiddleware 提供 MetaInfo 中间件
func ProvideMetaMiddleware(logger *zerolog.Logger) *middleware.MetaInfoMiddleware {
	return middleware.NewMetaInfoMiddleware(logger)
}

// ServerOptions 封装 Kitex Server 配置选项
type ServerOptions struct {
	Options []server.Option
	Addr    *net.TCPAddr
}

// ProvideServerOptions 提供 Kitex Server 配置选项
// 依赖 Provider 确保 OpenTelemetry 在服务器之前初始化
func ProvideServerOptions(
	cfg *config.Config,
	provider *otel.Provider,
	reg registry.Registry,
	metaMiddleware *middleware.MetaInfoMiddleware,
	_ error, // 依赖 ProvideKitexLogger 确保日志器已初始化
) (*ServerOptions, error) {
	// 解析监听地址
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve server address: %w", err)
	}

	// 服务注册名称（用于服务发现）
	serviceName := cfg.Server.Name

	// 构建 Server Options
	opts := []server.Option{
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
		server.WithRegistry(reg),
		server.WithServiceAddr(addr),
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler),
	}

	// 如果启用追踪，添加 OpenTelemetry Suite（必须在 MetaInfoMiddleware 之前）
	// 这样 OTel 会先创建 span，MetaInfoMiddleware 才能提取 trace_id/span_id
	if provider.IsEnabled() {
		opts = append(opts, server.WithSuite(tracing.NewServerSuite()))
		klog.Infof("OpenTelemetry tracing enabled, endpoint: %s", cfg.Tracing.Endpoint)
	}

	// 添加 MetaInfoMiddleware（在 OTel Suite 之后，确保 span 已创建）
	opts = append(opts, server.WithMiddleware(metaMiddleware.ServerMiddleware()))

	return &ServerOptions{
		Options: opts,
		Addr:    addr,
	}, nil
}
