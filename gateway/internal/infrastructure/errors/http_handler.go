// Package errors 提供了 Hertz 网关层的统一错误处理机制
// 本文件专门处理 HTTP 状态码相关的错误
package errors

import (
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

// HandleHTTPStatusError 处理 HTTP 状态码错误
// 将 HTTP 状态码错误转换为统一的业务错误响应
func HandleHTTPStatusError(c *app.RequestContext, statusCode int) {
	var apiErr APIError

	switch statusCode {
	case http.StatusNotFound:
		apiErr = ErrNotFound
	case http.StatusMethodNotAllowed:
		apiErr = ErrMethodNotAllowed
	case http.StatusInternalServerError:
		apiErr = ErrInternal
	case http.StatusBadRequest:
		apiErr = ErrInvalidParams
	case http.StatusUnauthorized:
		apiErr = ErrUnauthorized
	case http.StatusForbidden:
		apiErr = ErrForbidden
	case http.StatusGatewayTimeout:
		apiErr = ErrGatewayTimeout
	case http.StatusServiceUnavailable:
		apiErr = ErrServiceDown
	case http.StatusTooManyRequests:
		apiErr = ErrRateLimited
	default:
		// 其他 HTTP 错误码统一返回内部错误
		apiErr = ErrInternal
	}

	AbortWithError(c, apiErr)
}
