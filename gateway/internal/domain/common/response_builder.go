package common

import (
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/errors"
)

// ResponseBuilder 统一响应构建器
// 提供标准化的响应构建方法，减少重复代码并确保响应格式一致
type ResponseBuilder struct{}

// NewResponseBuilder 创建响应构建器
func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{}
}

// BuildSuccessResponse 构建成功响应的基础部分
// 注意：
// - request_id 通过 HTTP Header (X-Request-ID) 传递，由 requestid 中间件自动处理
// - timestamp 已从响应体移除，改用标准 HTTP Date 响应头
// - Date 响应头由 ResponseHeaderMiddleware 自动添加
func (rb *ResponseBuilder) BuildSuccessResponse() *http_base.BaseResponseDTO {
	return &http_base.BaseResponseDTO{
		Code:    errors.ErrSuccess.Code(),
		Message: errors.ErrSuccess.Message(),
	}
}

// BuildOperationStatusResponse 构建操作状态响应
func (rb *ResponseBuilder) BuildOperationStatusResponse() *http_base.OperationStatusResponseDTO {
	return &http_base.OperationStatusResponseDTO{
		BaseResp: rb.BuildSuccessResponse(),
	}
}

// WithData 为响应添加数据（通用方法）
func (rb *ResponseBuilder) WithData(data interface{}) interface{} {
	// 这是一个通用的数据包装方法
	// 实际使用时应该根据具体的响应类型来调用
	return data
}

// BuildDataResponse 构建带数据的响应（泛型模拟）
func BuildDataResponse[T any](data T) T {
	// Go 1.18+ 泛型写法，如果项目使用较老版本需要调整
	return data
}

// ============================================================================
// 业务特定响应构建器（可以根据需要扩展）
// ============================================================================

// BuildDepartmentResponse 构建部门响应
func (rb *ResponseBuilder) BuildDepartmentResponse(department interface{}) interface{} {
	return map[string]interface{}{
		"base_resp":  rb.BuildSuccessResponse(),
		"department": department,
	}
}

// BuildListResponse 构建列表响应
func (rb *ResponseBuilder) BuildListResponse(items interface{}, pagination interface{}) interface{} {
	return map[string]interface{}{
		"base_resp":  rb.BuildSuccessResponse(),
		"items":      items,
		"pagination": pagination,
	}
}
