// Package middleware 提供RPC服务端中间件
package middleware

import (
	"context"
	"time"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// MetaInfoMiddleware RPC服务端追踪中间件
// 职责：
// 1. 从 metainfo 提取 request_id
// 2. 如果不存在，自动生成并注入到 metainfo
// 3. 记录追踪信息日志
//
// 设计原则：
// - 直接使用 metainfo，不使用 context.WithValue（避免重复存储）
// - 确保每个请求都有完整的追踪 ID（自动生成缺失的）
// - 聚焦核心功能（只处理 request_id）
type MetaInfoMiddleware struct {
	logger *zerolog.Logger
}

// NewMetaInfoMiddleware 创建新的MetaInfo中间件实例
func NewMetaInfoMiddleware(logger *zerolog.Logger) *MetaInfoMiddleware {
	if logger == nil {
		defaultLogger := zerolog.Nop()
		logger = &defaultLogger
	}

	return &MetaInfoMiddleware{
		logger: logger,
	}
}

// ServerMiddleware 返回Kitex服务端中间件
func (m *MetaInfoMiddleware) ServerMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) error {
			// 记录开始时间
			start := time.Now()

			// 确保追踪 ID 存在，不存在则自动生成
			ctx = m.ensureRequestID(ctx)

			// 记录请求开始
			m.logTraceInfo(ctx)

			// 执行业务逻辑
			err := next(ctx, req, resp)

			// 记录请求结束和耗时
			duration := time.Since(start)
			m.logRequestComplete(ctx, duration, err)

			return err
		}
	}
}

// ensureRequestID 确保追踪 ID 存在，不存在则自动生成
// 关键设计：
// - 直接操作 metainfo，不使用 context.WithValue
// - 缺失的 ID 会自动生成并注入到 metainfo
func (m *MetaInfoMiddleware) ensureRequestID(ctx context.Context) context.Context {
	// 检查 request_id 是否已存在
	if id, ok := metainfo.GetPersistentValue(ctx, "request_id"); ok && id != "" {
		return ctx // 已存在，直接返回
	}

	// 生成新的 request_id
	requestID := uuid.New().String()
	ctx = metainfo.WithPersistentValue(ctx, "request_id", requestID)

	m.logger.Warn().
		Str("request_id", requestID).
		Str("service", "identity_srv").
		Msg("Generated missing request_id")

	return ctx
}

// getMethodName 从 Kitex RPCInfo 提取方法名
func (m *MetaInfoMiddleware) getMethodName(ctx context.Context) string {
	ri := rpcinfo.GetRPCInfo(ctx)
	if ri != nil {
		return ri.Invocation().MethodName()
	}
	return "unknown"
}

// logTraceInfo 记录追踪信息日志（请求开始）
func (m *MetaInfoMiddleware) logTraceInfo(ctx context.Context) {
	requestID := GetRequestID(ctx)
	if requestID == "" {
		return
	}

	methodName := m.getMethodName(ctx)

	m.logger.Info().
		Str("request_id", requestID).
		Str("service", "identity_srv").
		Str("method", methodName).
		Msg("RPC request started")
}

// logRequestComplete 记录请求完成和性能指标
func (m *MetaInfoMiddleware) logRequestComplete(ctx context.Context, duration time.Duration, err error) {
	requestID := GetRequestID(ctx)
	methodName := m.getMethodName(ctx)

	event := m.logger.Info()
	if duration > 100*time.Millisecond {
		event = m.logger.Warn() // 慢调用使用 Warn 级别
	}

	event = event.
		Str("request_id", requestID).
		Str("service", "identity_srv").
		Str("method", methodName).
		Dur("duration_ms", duration)

	if err != nil {
		event = event.Err(err).Str("status", "error")
	} else {
		event = event.Str("status", "success")
	}

	event.Msg("RPC request completed")
}

// Context 访问辅助函数，供业务逻辑使用

// GetRequestID 从 RPC 上下文获取 RequestID
// 直接从 metainfo 读取，不使用 context.Value
func GetRequestID(ctx context.Context) string {
	if id, ok := metainfo.GetPersistentValue(ctx, "request_id"); ok {
		return id
	}

	return ""
}

// LoggingAttrs 返回用于结构化日志的属性
// 返回 map[string]interface{} 用于 zerolog
func LoggingAttrs(ctx context.Context) map[string]interface{} {
	attrs := make(map[string]interface{})

	if requestID := GetRequestID(ctx); requestID != "" {
		attrs["request_id"] = requestID
	}

	return attrs
}
