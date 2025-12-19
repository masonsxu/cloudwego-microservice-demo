// Package otel provides OpenTelemetry initialization for the Identity RPC service.
// It configures tracing and metrics exporters using OTLP protocol.
package otel

import (
	"context"

	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/config"
)

// Provider wraps the OpenTelemetry provider for dependency injection.
type Provider struct {
	provider provider.OtelProvider
	enabled  bool
}

// NewProvider creates a new OpenTelemetry provider based on configuration.
// If tracing is disabled, it returns a no-op provider.
func NewProvider(cfg *config.Config) (*Provider, error) {
	if !cfg.Tracing.Enabled {
		return &Provider{enabled: false}, nil
	}

	opts := []provider.Option{
		provider.WithEnableMetrics(false), // Jaeger 不支持 OTLP Metrics，需要 Prometheus 后端
		provider.WithServiceName(cfg.Server.Name),
		provider.WithExportEndpoint(cfg.Tracing.Endpoint),
		provider.WithInsecure(),
	}

	p := provider.NewOpenTelemetryProvider(opts...)

	return &Provider{
		provider: p,
		enabled:  true,
	}, nil
}

// Shutdown gracefully shuts down the OpenTelemetry provider.
func (p *Provider) Shutdown(ctx context.Context) error {
	if !p.enabled || p.provider == nil {
		return nil
	}

	return p.provider.Shutdown(ctx)
}

// IsEnabled returns whether OpenTelemetry is enabled.
func (p *Provider) IsEnabled() bool {
	return p.enabled
}
