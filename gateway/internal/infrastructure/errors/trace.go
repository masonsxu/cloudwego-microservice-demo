// Package errors 提供了 Hertz 网关层的统一错误处理机制
// 本文件提供 RequestID 在 RPC 调用链中的传播功能
package errors

import (
	"context"

	"github.com/bytedance/gopkg/cloud/metainfo"
)

// InjectRequestIDToContext 将 RequestID 注入到 context 中（用于 RPC 调用）
// 使用 metainfo.WithPersistentValue 确保 RequestID 通过 TTHeader 传递到 RPC 服务
func InjectRequestIDToContext(ctx context.Context, requestID string) context.Context {
	if requestID == "" {
		return ctx
	}

	// 注入 request_id 到 metainfo，确保通过 TTHeader 传递到 RPC 服务
	ctx = metainfo.WithPersistentValue(ctx, "request_id", requestID)

	return ctx
}

// GetRequestIDFromContext 从 context 中获取 request_id（辅助函数）
// 优先从 metainfo 读取，如果没有则从普通 context 读取
func GetRequestIDFromContext(ctx context.Context) string {
	// 优先从 metainfo 提取
	if requestID, ok := metainfo.GetPersistentValue(ctx, "request_id"); ok && requestID != "" {
		return requestID
	}

	// 如果 metainfo 中没有，尝试从普通 context 提取
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}

	return ""
}
