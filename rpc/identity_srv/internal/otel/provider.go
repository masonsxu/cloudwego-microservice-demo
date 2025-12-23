// Package otel provides OpenTelemetry initialization for the Identity RPC service.
// It configures tracing and metrics exporters using OTLP protocol.
package otel

import (
	"context"

	"github.com/kitex-contrib/obs-opentelemetry/provider"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/config"
)

// NewProvider creates a new OpenTelemetry provider based on configuration.
// If tracing is disabled, it returns a no-op provider.
func NewProvider(cfg *config.Config) (func(context.Context) error, error) {
	if !cfg.Tracing.Enabled {
		return func(ctx context.Context) error {
			return nil
		}, nil
	}

	opts := []provider.Option{
		provider.WithEnableMetrics(false), // Jaeger 不支持 OTLP Metrics，需要 Prometheus 后端
		provider.WithServiceName(cfg.Server.Name),
		provider.WithExportEndpoint(cfg.Tracing.Endpoint),
		provider.WithInsecure(),
	}

	p := provider.NewOpenTelemetryProvider(opts...)

	return p.Shutdown, nil
}
