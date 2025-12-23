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

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/log"
)

// MetaInfoMiddleware RPC服务端追踪中间件
// 职责：
// 1. 从 metainfo 提取 request_id
// 2. 如果不存在，自动生成并注入到 metainfo
// 3. 记录追踪信息日志（包含 trace_id, span_id）
//
// 设计原则：
// - 直接使用 metainfo，不使用 context.WithValue（避免重复存储）
// - 确保每个请求都有完整的追踪 ID（自动生成缺失的）
// - 聚焦核心功能（处理 request_id 和 OTel trace 信息）
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
			// 将带 request_id 的 logger 绑定到 ctx，供业务逻辑使用
			ctx = m.bindLoggerToContext(ctx)

			// 记录请求开始（此时 OpenTelemetry Span 可能还未创建）
			m.logTraceInfo(ctx)

			// 执行业务逻辑（OpenTelemetry Suite 会在这里创建 Span）
			err := next(ctx, req, resp)

			// 记录请求结束和耗时（此时应该能获取到 trace_id）
			duration := time.Since(start)
			m.logRequestComplete(ctx, duration, err)

			return err
		}
	}
}

// bindLoggerToContext 将带追踪信息的 logger 绑定到 ctx
// 包含 request_id, trace_id, span_id, service, method 等字段
// 业务代码可通过 zerolog.Ctx(ctx) 或 log.Ctx(ctx) 获取带追踪字段的 logger
func (m *MetaInfoMiddleware) bindLoggerToContext(ctx context.Context) context.Context {
	methodName := m.getMethodName(ctx)

	// 使用 log 包绑定完整追踪信息
	return log.BindToContext(ctx, *m.logger, "identity_srv", methodName)
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

	// 使用 log.Event 动态提取追踪信息（包括 OpenTelemetry trace_id/span_id）
	log.Event(ctx, m.logger.Info()).
		Str("method", methodName).
		Msg("RPC request started")
}

// logRequestComplete 记录请求完成和性能指标
func (m *MetaInfoMiddleware) logRequestComplete(
	ctx context.Context,
	duration time.Duration,
	err error,
) {
	methodName := m.getMethodName(ctx)

	event := m.logger.Info()
	if duration > 100*time.Millisecond {
		event = m.logger.Warn() // 慢调用使用 Warn 级别
	}

	// 使用 log.Event 添加追踪字段
	event = log.Event(ctx, event).
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
// 包含 request_id, trace_id, span_id
// 返回 map[string]interface{} 用于 zerolog
func LoggingAttrs(ctx context.Context) map[string]interface{} {
	// 使用 log 包获取完整追踪字段
	return log.TraceFields(ctx)
}

// GetTraceID 从 context 获取 OpenTelemetry TraceID
func GetTraceID(ctx context.Context) string {
	return log.GetTraceID(ctx)
}

// GetSpanID 从 context 获取 OpenTelemetry SpanID
func GetSpanID(ctx context.Context) string {
	return log.GetSpanID(ctx)
}
