// Package otel provides OpenTelemetry initialization for the Gateway service.
// It configures tracing and metrics exporters using OTLP protocol.
package otel

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"

	"github.com/hertz-contrib/obs-opentelemetry/provider"

	appconfig "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
)

// Provider 封装 OpenTelemetry Provider，支持 Wire 依赖注入
type Provider struct {
	shutdown func(context.Context) error
	enabled  bool
}

// NewProvider creates a new OpenTelemetry provider based on configuration.
// Returns Provider wrapper and cleanup function for Wire.
func NewProvider(cfg *appconfig.Configuration) (*Provider, func(), error) {
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

// Tracer 封装 Hertz Server Tracer，支持 Wire 依赖注入
type Tracer struct {
	option config.Option
	Config *hertztracing.Config
}

// NewTracer creates a new Hertz server tracer based on configuration.
func NewTracer(cfg *appconfig.Configuration) *Tracer {
	tracerOpts := NewServerTracerOptions(cfg)
	tracer, tracerCfg := hertztracing.NewServerTracer(tracerOpts...)

	return &Tracer{
		option: tracer,
		Config: tracerCfg,
	}
}

// ServerFactory 封装 Hertz Server 创建逻辑
type ServerFactory struct {
	server *server.Hertz
}

// NewServerFactory creates a new Hertz server factory.
// 依赖 Provider 确保 OpenTelemetry 在服务器之前初始化
func NewServerFactory(cfg *appconfig.Configuration, tracer *Tracer, _ *Provider) *ServerFactory {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	serverOpts := []config.Option{
		server.WithHostPorts(addr),
		server.WithMaxRequestBodySize(100 * 1024 * 1024),
	}

	// 添加 tracer 选项
	if tracer != nil {
		serverOpts = append(serverOpts, tracer.option)
	}

	h := server.New(serverOpts...)

	return &ServerFactory{
		server: h,
	}
}

// Server returns the Hertz server instance.
func (f *ServerFactory) Server() *server.Hertz {
	return f.server
}
