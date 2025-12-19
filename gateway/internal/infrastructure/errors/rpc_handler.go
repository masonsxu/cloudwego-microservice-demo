// Package errors 提供了 Hertz 网关层的统一错误处理机制
// 本文件专门处理 RPC 调用相关的错误
package errors

import (
	"github.com/cloudwego/kitex/pkg/kerrors"
)

// ProcessRPCError 统一处理 RPC 调用错误，返回适合的 APIError
//
// 处理逻辑：
//  1. 如果是 RPC BizStatusError（业务错误），直接透传错误码和消息
//     - RPC 业务错误码范围：200xxx（按业务领域编码）
//     - 错误信息保持原样，不做转换
//  2. 如果是 RPC 框架错误（网络超时、连接失败等），返回网关内部错误（100005）
//     - 使用 fallbackMessage 作为用户友好的错误提示
//
// 返回的 APIError 将通过 GetHTTPStatus 映射为对应的 HTTP 状态码：
//   - RPC 业务错误（200xxx）-> HTTP 200（错误信息在 response body 中）
//   - 网关系统错误（100xxx）-> 对应 HTTP 状态码（500）
//
// 使用示例：
//
//	resp, err := s.identityClient.GetUser(ctx, req)
//	if err != nil {
//	    return nil, errors.ProcessRPCError(err, "获取用户信息失败")
//	}
func ProcessRPCError(err error, fallbackMessage string) error {
	if err == nil {
		return nil
	}

	// 使用 kerrors.FromBizStatusError 提取业务异常
	if bizErr, isBizErr := kerrors.FromBizStatusError(err); isBizErr {
		// 直接将 RPC 业务错误转换为 API 错误，保持原始错误码和消息
		return NewAPIError(bizErr.BizStatusCode(), bizErr.BizMessage())
	}

	// 其他类型错误按系统错误处理
	return ErrInternal.WithMessage(fallbackMessage)
}
