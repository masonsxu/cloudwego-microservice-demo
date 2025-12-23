// Package errors 提供了 Hertz 网关层的统一错误处理机制
// 本文件专门处理错误响应的生成和发送
package errors

import (
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/requestid"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
)

// 错误码到 HTTP 状态码的映射
// 支持网关层错误码（100xxx）和 RPC 业务层认证相关错误码（201xxx）
var httpStatusMap = map[int32]int{
	// 系统级错误
	CodeSuccess:          http.StatusOK,
	CodeInvalidParams:    http.StatusBadRequest,
	CodeUnauthorized:     http.StatusUnauthorized,
	CodeForbidden:        http.StatusForbidden,
	CodeNotFound:         http.StatusNotFound,
	CodeInternalError:    http.StatusInternalServerError,
	CodeMethodNotAllowed: http.StatusMethodNotAllowed,

	// JWT认证相关错误
	CodeJWTTokenMissing:    http.StatusUnauthorized,
	CodeJWTTokenInvalid:    http.StatusUnauthorized,
	CodeJWTTokenExpired:    http.StatusUnauthorized,
	CodeJWTTokenNotActive:  http.StatusUnauthorized,
	CodeJWTTokenMalformed:  http.StatusBadRequest,
	CodeJWTValidationFail:  http.StatusUnauthorized,
	CodeJWTSigningError:    http.StatusInternalServerError,
	CodeJWTCreationFail:    http.StatusInternalServerError,
	CodeInvalidCredentials: http.StatusUnauthorized,

	// 网关特有错误
	CodeGatewayTimeout: http.StatusGatewayTimeout,
	CodeServiceDown:    http.StatusServiceUnavailable,
	CodeRateLimited:    http.StatusTooManyRequests,

	// RPC 业务层认证相关错误 (201xxx - identity_srv)
	// 这些错误来自下游 RPC 服务，需要在网关层映射为正确的 HTTP 状态码
	CodeRPCUserInactive:         http.StatusForbidden,    // 用户未激活
	CodeRPCUserNotFound:         http.StatusUnauthorized, // 用户不存在
	CodeRPCInvalidCredentials:   http.StatusUnauthorized, // 用户名或密码错误
	CodeRPCUserSuspended:        http.StatusForbidden,    // 用户已停用
	CodeRPCMustChangePassword:   http.StatusForbidden,    // 需要修改密码
	CodeRPCUserNoAvailableRoles: http.StatusForbidden,    // 用户没有可用角色
}

// AbortWithError 中断请求并返回错误响应
// 与成功响应保持一致的结构，便于前端统一处理
// RequestID 会通过 HTTP Header (X-Request-ID) 传递，由 requestid 中间件自动处理
// Date 响应头由 ResponseHeaderMiddleware 自动添加
func AbortWithError(c *app.RequestContext, err APIError) {
	httpStatus := GetHTTPStatus(err.Code())

	// 使用 OperationStatusResponseDTO 结构
	response := &http_base.OperationStatusResponseDTO{
		BaseResp: &http_base.BaseResponseDTO{
			Code:    err.Code(),
			Message: err.Message(),
		},
	}

	c.JSON(httpStatus, response)
	c.Abort()
}

// AbortWithErrorMessage 中断请求并返回自定义错误消息
// 与成功响应保持一致的结构，便于前端统一处理
// RequestID 会通过 HTTP Header (X-Request-ID) 传递，由 requestid 中间件自动处理
// Date 响应头由 ResponseHeaderMiddleware 自动添加
func AbortWithErrorMessage(c *app.RequestContext, err APIError, message string) {
	httpStatus := GetHTTPStatus(err.Code())

	// 使用 OperationStatusResponseDTO 结构
	response := &http_base.OperationStatusResponseDTO{
		BaseResp: &http_base.BaseResponseDTO{
			Code:    err.Code(),
			Message: message,
		},
	}

	c.JSON(httpStatus, response)
	c.Abort()
}

// GenerateRequestID 从请求上下文获取 RequestID
// 使用 hertz-contrib/requestid 中间件提供的 Get 函数
func GenerateRequestID(c *app.RequestContext) string {
	return requestid.Get(c)
}

// GetHTTPStatus 根据业务错误码获取对应的 HTTP 状态码
func GetHTTPStatus(code int32) int {
	if status, exists := httpStatusMap[code]; exists {
		return status
	}

	// 根据错误码规范判断
	switch {
	case code == 0:
		return http.StatusOK
	case code >= 20000 && code < 30000: // RPC 业务错误，统一返回 200
		return http.StatusOK
	default:
		// 系统级错误默认返回 500
		return http.StatusInternalServerError
	}
}
