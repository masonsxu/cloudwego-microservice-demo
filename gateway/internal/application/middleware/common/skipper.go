package common

import (
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

// ShouldSkip 检查是否跳过权限校验
func ShouldSkip(c *app.RequestContext, skipPaths []string) bool {
	path := string(c.Request.URI().Path())
	method := string(c.Request.Method())

	for _, pattern := range skipPaths {
		if matchSkip(pattern, method, path) {
			return true // 继续下一个 handler，不再走鉴权
		}
	}

	return false
}

func matchSkip(pattern, method, path string) bool {
	if strings.Contains(pattern, ":") {
		parts := strings.SplitN(pattern, ":", 2)
		return parts[0] == method && parts[1] == path
	}

	if strings.HasSuffix(pattern, "/*") {
		return strings.HasPrefix(path, strings.TrimSuffix(pattern, "/*"))
	}

	return pattern == path
}
