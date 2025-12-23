// Package otel provides OpenTelemetry initialization for the Gateway service.
// It configures tracing and metrics exporters using OTLP protocol.
package otel

import (
	"context"

	"github.com/hertz-contrib/obs-opentelemetry/provider"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
)

// NewProvider creates a new OpenTelemetry provider based on configuration.
// If tracing is disabled, it returns a no-op cleanup function.
func NewProvider(cfg *config.Configuration) (func(context.Context) error, error) {
	if !cfg.Tracing.Enabled {
		return nil, nil
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
