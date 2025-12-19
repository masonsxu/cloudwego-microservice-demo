// Package identitycli 提供与身份服务交互的客户端实现
package identitycli

import (
	"log/slog"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	kitextracing "github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	conf "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv/identityservice"
)

const (
	defaultIdentityServiceName = "identity-service"
)

func configureTimeouts(opts []client.Option) []client.Option {
	if conf.Config.Client.ConnectionTimeout > 0 {
		opts = append(opts, client.WithConnectTimeout(conf.Config.Client.ConnectionTimeout))
		slog.Debug(
			"Set connection timeout",
			"timeout_seconds",
			conf.Config.Client.ConnectionTimeout,
		)
	}

	if conf.Config.Client.RequestTimeout > 0 {
		opts = append(opts, client.WithRPCTimeout(conf.Config.Client.RequestTimeout))
		slog.Debug("Set request timeout", "timeout_seconds", conf.Config.Client.RequestTimeout)
	}

	return opts
}

// NewIdentityClient 创建聚合的用户客户端，使用单一的 Kitex 客户端实例
func NewIdentityClient() (IdentityClient, error) {
	// 创建单一的 Kitex 客户端实例
	r, err := etcd.NewEtcdResolver([]string{conf.Config.Etcd.Address})
	if err != nil {
		slog.Error("Failed to create etcd resolver", "error", err)
		return nil, err
	}

	identityServiceName := defaultIdentityServiceName
	if service, exists := conf.Config.Client.Services["identity"]; exists &&
		service.Name != "" {
		identityServiceName = service.Name
	}

	slog.Info("Creating identity client", "service_name", identityServiceName)
	// 连接池配置
	idleConfig := connpool.IdleConfig{
		// MaxIdlePerAddress 建议值 = QPS_per_dest_host * avg_response_time_sec:cite[4]
		// 例如：QPS=100, 平均响应时间0.1秒, 则建议设置为 10 (100 * 0.1)
		MaxIdlePerAddress: conf.Config.Client.Pool.MaxIdlePerAddress,
		MaxIdleGlobal:     conf.Config.Client.Pool.MaxIdleGlobal,
		MaxIdleTimeout:    conf.Config.Client.Pool.MaxIdleTimeout,
	}
	// 构建客户端选项
	opts := []client.Option{
		client.WithResolver(r),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: conf.Config.Server.Name,
		}),
		// 配置TTHeader传输协议
		client.WithTransportProtocol(transport.TTHeader),
		// 修正：统一使用 TTHeader MetaHandler 以支持BizStatusError的正确传递:cite[1]
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		// 配置长连接池:cite[4]
		client.WithLongConnection(idleConfig),
		// 配置熔断器中间件:cite[6]
		// client.WithMiddleware(cbSuite.ServiceCBMW()),
		// 配置实例级别的熔断器中间件（支持自动重试）:cite[6]
		// client.WithInstanceMW(cbSuite.InstanceCBMW()),
		// 可选：配置熔断器中间件（使用 kitex-contrib/cbreaker）
		// client.WithMiddleware(cbreaker.NewCircuitBreakerMiddleware()),
		// 配置负载均衡策略:cite[2]
		// client.WithLoadBalance(loadbalance.NewWeightedRoundBalancer()),
	}

	// 添加 OpenTelemetry tracing Suite (如果启用)
	if conf.Config.Tracing.Enabled {
		opts = append(opts, client.WithSuite(kitextracing.NewClientSuite()))

		slog.Debug("OpenTelemetry tracing enabled for identity client")
	}

	// 配置超时
	opts = configureTimeouts(opts)

	cli, err := identityservice.NewClient(
		identityServiceName,
		opts...,
	)
	if err != nil {
		slog.Error("Failed to create identity client", "error", err)
		return nil, err
	}

	slog.Info("Successfully created identity client")

	return cli, nil
}
