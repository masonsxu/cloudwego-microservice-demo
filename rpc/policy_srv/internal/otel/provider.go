// Package otel 提供 policy_srv 的 OpenTelemetry 初始化。
package otel

import (
	"context"

	"github.com/kitex-contrib/obs-opentelemetry/provider"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/config"
)

// Provider 封装 OpenTelemetry Provider，支持 Wire 依赖注入
type Provider struct {
	shutdown func(context.Context) error
	enabled  bool
}

// NewProvider 根据配置创建 OpenTelemetry provider，返回 Provider 与 cleanup
func NewProvider(cfg *config.Config) (*Provider, func(), error) {
	if !cfg.Tracing.Enabled {
		return &Provider{enabled: false}, func() {}, nil
	}

	opts := []provider.Option{
		provider.WithEnableMetrics(false),
		provider.WithServiceName(cfg.Server.Name),
		provider.WithExportEndpoint(cfg.Tracing.Endpoint),
		provider.WithInsecure(),
	}

	p := provider.NewOpenTelemetryProvider(opts...)

	otelProvider := &Provider{
		shutdown: p.Shutdown,
		enabled:  true,
	}

	cleanup := func() {
		if otelProvider.shutdown != nil {
			_ = otelProvider.shutdown(context.Background())
		}
	}

	return otelProvider, cleanup, nil
}

// IsEnabled 返回 tracing 是否启用
func (p *Provider) IsEnabled() bool {
	return p.enabled
}
