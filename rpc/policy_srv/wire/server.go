package wire

import (
	"fmt"
	"net"
	"time"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/rs/zerolog"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/config"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/internal/otel"
)

// ServerOptions Kitex 服务器选项
type ServerOptions struct {
	Options []server.Option
	Addr    *net.TCPAddr
}

// ProvideOtelProvider 提供 OpenTelemetry Provider
func ProvideOtelProvider(cfg *config.Config) (*otel.Provider, func(), error) {
	return otel.NewProvider(cfg)
}

// ProvideEtcdRegistry 提供 etcd 注册中心
func ProvideEtcdRegistry(cfg *config.Config) (registry.Registry, error) {
	r, err := etcd.NewEtcdRegistry([]string{cfg.Etcd.Address})
	if err != nil {
		return nil, fmt.Errorf("create etcd registry: %w", err)
	}

	return r, nil
}

// ProvideServerOptions 提供 Kitex 服务器选项
// 依赖 *otel.Provider 确保 OTel 在 Server 启动前完成初始化
func ProvideServerOptions(
	cfg *config.Config,
	reg registry.Registry,
	provider *otel.Provider,
	logger *zerolog.Logger,
) (*ServerOptions, error) {
	addr, err := net.ResolveTCPAddr("tcp", cfg.Server.Address)
	if err != nil {
		return nil, fmt.Errorf("resolve address: %w", err)
	}

	_ = logger

	opts := []server.Option{
		server.WithServiceAddr(addr),
		server.WithRegistry(reg),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: cfg.Server.Name,
		}),
		server.WithLimit(&limit.Option{
			MaxConnections: 1000,
			MaxQPS:         500,
		}),
		server.WithReadWriteTimeout(30 * time.Second),
		server.WithExitWaitTime(5 * time.Second),
	}

	if provider.IsEnabled() {
		opts = append(opts, server.WithSuite(tracing.NewServerSuite()))
	}

	return &ServerOptions{
		Options: opts,
		Addr:    addr,
	}, nil
}
