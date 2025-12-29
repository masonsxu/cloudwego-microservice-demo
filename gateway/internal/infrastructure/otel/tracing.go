// Package otel provides OpenTelemetry initialization for the Gateway service.
package otel

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/config"
)

// NewIgnorePathsOption 根据配置创建跳过指定路径的追踪选项。
// 支持精确匹配和通配符匹配（路径以 * 结尾表示前缀匹配）。
func NewIgnorePathsOption(cfg *config.Configuration) hertztracing.Option {
	ignorePaths := cfg.Tracing.IgnorePaths
	if len(ignorePaths) == 0 {
		return nil
	}

	// 预处理路径，分离精确匹配和前缀匹配
	exactPaths := make(map[string]struct{})
	prefixPaths := make([]string, 0)

	for _, path := range ignorePaths {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}

		if strings.HasSuffix(path, "*") {
			// 前缀匹配：去掉末尾的 *
			prefixPaths = append(prefixPaths, strings.TrimSuffix(path, "*"))
		} else {
			// 精确匹配
			exactPaths[path] = struct{}{}
		}
	}

	return hertztracing.WithShouldIgnore(func(_ context.Context, c *app.RequestContext) bool {
		reqPath := string(c.Path())

		// 精确匹配
		if _, ok := exactPaths[reqPath]; ok {
			return true
		}

		// 前缀匹配
		for _, prefix := range prefixPaths {
			if strings.HasPrefix(reqPath, prefix) {
				return true
			}
		}

		return false
	})
}

// NewServerTracerOptions 根据配置创建服务端追踪器选项列表。
func NewServerTracerOptions(cfg *config.Configuration) []hertztracing.Option {
	opts := make([]hertztracing.Option, 0)

	if opt := NewIgnorePathsOption(cfg); opt != nil {
		opts = append(opts, opt)
	}

	return opts
}
