// Package otel provides OpenTelemetry initialization for the Identity RPC service.
// It configures tracing and metrics exporters using OTLP protocol.
package otel

import (
	"context"

	"github.com/kitex-contrib/obs-opentelemetry/provider"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/config"
)

// Provider 封装 OpenTelemetry Provider，支持 Wire 依赖注入
type Provider struct {
	shutdown func(context.Context) error
	enabled  bool
}

// NewProvider creates a new OpenTelemetry provider based on configuration.
// Returns Provider wrapper and cleanup function for Wire.
func NewProvider(cfg *config.Config) (*Provider, func(), error) {
	if !cfg.Tracing.Enabled {
		return &Provider{enabled: false}, func() {}, nil
	}

	opts := []provider.Option{
		provider.WithEnableMetrics(false), // Jaeger 不支持 OTLP Metrics，需要 Prometheus 后端
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

// IsEnabled returns whether tracing is enabled.
func (p *Provider) IsEnabled() bool {
	return p.enabled
}
