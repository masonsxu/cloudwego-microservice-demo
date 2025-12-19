package errors

import (
	"github.com/cloudwego/hertz/pkg/app"
)

// handleServiceError 统一处理服务层错误的辅助函数
// 遵循项目的错误处理规范，减少代码重复
func HandleServiceError(c *app.RequestContext, err error, defaultMessage string) {
	if err == nil {
		return
	}

	// 检查是否为 APIError 类型
	if apiErr, ok := err.(APIError); ok {
		// 业务层返回的是 APIError，使用统一错误处理
		AbortWithError(c, apiErr)
		return
	}
	// 其他未知错误，转换为系统内部错误
	AbortWithError(c, ErrInternal.WithMessage(defaultMessage))
}
