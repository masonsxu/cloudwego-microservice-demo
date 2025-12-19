// Package errors 提供了 Hertz 网关层的统一错误处理机制
// 本文件定义了错误类型、错误码常量和预定义错误变量
package errors

import "fmt"

// APIError API 层业务错误结构体
// 采用私有字段设计，通过方法访问，确保数据安全性
type APIError struct {
	code    int32  // 错误码
	message string // 错误消息
}

// Error 实现 error 接口
func (e APIError) Error() string {
	return fmt.Sprintf("APIError: code=%d, message=%s", e.code, e.message)
}

// NewAPIError 创建新的 API 业务错误
func NewAPIError(code int32, message string) APIError {
	return APIError{
		code:    code,
		message: message,
	}
}

// WithMessage 设置错误消息，返回新的错误实例
func (e APIError) WithMessage(message string) APIError {
	e.message = message
	return e
}

// Code 获取错误码
func (e APIError) Code() int32 {
	return e.code
}

// Message 获取错误消息
func (e APIError) Message() string {
	return e.message
}

// =================================================================
//
//	网关层错误码规范
//
// =================================================================
// 网关层错误码分为两类：
//
//  1. 网关自有错误（100xxx）：
//     处理 HTTP 层面的错误，格式：1-BB-CCC
//     - 10: 通用系统错误（认证、授权、参数验证等）
//     - 11: 网关基础设施错误（超时、服务不可用、限流等）
//
//  2. RPC 业务错误（200xxx）：
//     透传自下游 RPC 服务的业务错误码，由 ProcessRPCError 自动转换
//     格式：2-BB-CCC（按业务领域编码）
//     - 00: 通用业务错误
//     - 01: 用户领域（identity_srv）
//     - 02: 组织领域（identity_srv）
//     - 03: 部门领域（identity_srv）
//     - 04: 数据一致性领域（identity_srv）
//     - 05: 角色定义领域（permission_srv）
//     - 06: 角色分配领域（permission_srv）
//     - 07: 菜单权限领域（permission_srv）
//
// HTTP 状态码映射规则：
//   - 200xxx 业务错误：HTTP 200（错误信息在 response body 中）
//   - 100xxx 系统错误：对应的 HTTP 状态码（401, 403, 500 等）
const (
	// 系统级错误 (100xxx)
	CodeSuccess          = 0      // 成功
	CodeInvalidParams    = 100001 // 无效参数
	CodeUnauthorized     = 100002 // 未授权/认证失败
	CodeForbidden        = 100003 // 权限不足
	CodeNotFound         = 100004 // 资源未找到
	CodeInternalError    = 100005 // 服务内部错误
	CodeMethodNotAllowed = 100006 // 请求方法不被允许

	// JWT认证相关错误 (102xxx)
	CodeJWTTokenMissing    = 102001 // JWT令牌缺失
	CodeJWTTokenInvalid    = 102002 // JWT令牌格式无效
	CodeJWTTokenExpired    = 102003 // JWT令牌已过期
	CodeJWTTokenNotActive  = 102004 // JWT令牌未生效（nbf校验失败）
	CodeJWTTokenMalformed  = 102005 // JWT令牌结构错误
	CodeJWTValidationFail  = 102006 // JWT验证失败（通用验证错误）
	CodeJWTSigningError    = 102007 // JWT签名生成失败
	CodeJWTCreationFail    = 102008 // JWT令牌创建失败
	CodeInvalidCredentials = 102009 // 认证凭据无效（用户名密码错误）

	// 授权和权限相关错误 (103xxx)
	CodeUserNoAvailableRoles = 103001 // 用户无可用角色

	// 网关特有错误 (110xxx)
	CodeGatewayTimeout = 110001 // 网关超时
	CodeServiceDown    = 110002 // 下游服务不可用
	CodeRateLimited    = 110003 // 请求限流

	// =================================================================
	// RPC 业务错误码范围（200xxx）
	// =================================================================
	// 以下错误码由下游 RPC 服务定义，网关层通过 ProcessRPCError 自动透传
	// 不需要在网关层定义常量，仅作为文档说明
	//
	// 通用业务错误（200xxx）：
	//   200100: 参数错误
	//   200101: 操作失败

	// =================================================================
	// RPC 业务错误码（需要映射 HTTP 状态码的特殊错误）
	// =================================================================
	// 这些错误码来自下游 RPC 服务，但需要在网关层映射为特定的 HTTP 状态码
	// 而非默认的 HTTP 200

	// 用户认证相关的 RPC 业务错误 (201xxx - identity_srv)
	CodeRPCUserNotFound       = 201001 // 用户不存在
	CodeRPCUserInactive       = 201012 // 用户未激活
	CodeRPCInvalidCredentials = 201016 // 用户名或密码错误
	CodeRPCUserSuspended      = 201017 // 用户已停用
	CodeRPCMustChangePassword = 201018 // 需要修改密码
	// 角色分配相关的 RPC 业务错误 (207xxx - identity_srv)
	CodeRPCUserNoAvailableRoles = 207016 // 用户没有可用角色
	// 数据源配置相关的 RPC 业务错误 (208xxx - cancer_srv)
	CodeRPCDataSourceInUse = 208001 // 数据源正在被字典配置使用，无法删除
)

// 预定义 API 错误变量
// 提供常用错误的预定义实例，避免重复创建
var (
	// 系统级错误
	ErrSuccess          = NewAPIError(CodeSuccess, "success")
	ErrInvalidParams    = NewAPIError(CodeInvalidParams, "参数错误")
	ErrUnauthorized     = NewAPIError(CodeUnauthorized, "未授权/认证失败")
	ErrForbidden        = NewAPIError(CodeForbidden, "权限不足")
	ErrNotFound         = NewAPIError(CodeNotFound, "资源不存在")
	ErrInternal         = NewAPIError(CodeInternalError, "系统繁忙，请稍后重试")
	ErrMethodNotAllowed = NewAPIError(CodeMethodNotAllowed, "请求方法不被允许")

	// JWT认证相关错误
	ErrJWTTokenMissing    = NewAPIError(CodeJWTTokenMissing, "令牌缺失")
	ErrJWTTokenInvalid    = NewAPIError(CodeJWTTokenInvalid, "令牌格式无效")
	ErrJWTTokenExpired    = NewAPIError(CodeJWTTokenExpired, "令牌已过期")
	ErrJWTTokenNotActive  = NewAPIError(CodeJWTTokenNotActive, "令牌未生效")
	ErrJWTTokenMalformed  = NewAPIError(CodeJWTTokenMalformed, "令牌结构错误")
	ErrJWTValidationFail  = NewAPIError(CodeJWTValidationFail, "令牌验证失败")
	ErrJWTSigningError    = NewAPIError(CodeJWTSigningError, "令牌签名生成失败")
	ErrJWTCreationFail    = NewAPIError(CodeJWTCreationFail, "令牌创建失败")
	ErrInvalidCredentials = NewAPIError(CodeInvalidCredentials, "用户名或密码错误")

	// 授权和权限相关错误
	ErrUserNoAvailableRoles = NewAPIError(CodeUserNoAvailableRoles, "用户无可用角色，无法登录")
	// 数据源配置相关的 RPC 业务错误 (208xxx - cancer_srv)
	ErrDataSourceInUse = NewAPIError(CodeRPCDataSourceInUse, "数据源正在被字典配置使用，无法删除")

	// 网关特有错误
	ErrGatewayTimeout = NewAPIError(CodeGatewayTimeout, "请求超时")
	ErrServiceDown    = NewAPIError(CodeServiceDown, "服务暂不可用")
	ErrRateLimited    = NewAPIError(CodeRateLimited, "请求过于频繁")
)
