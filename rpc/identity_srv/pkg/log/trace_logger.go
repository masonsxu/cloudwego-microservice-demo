// Package log 提供统一的日志工具，支持 OpenTelemetry 追踪信息注入
//
// 本包解决日志与链路追踪对齐的问题，确保每条日志都包含：
// - trace_id: OpenTelemetry 链路追踪 ID
// - span_id: OpenTelemetry Span ID
// - request_id: 业务层请求 ID (从 metainfo 传播)
//
// 使用示例:
//
//	// 在业务代码中获取带追踪信息的 logger
//	logger := log.Ctx(ctx)
//	logger.Info().Str("user_id", "123").Msg("User created")
//
//	// 或者使用 WithTrace 为现有 logger 添加追踪字段
//	logger := log.WithTrace(ctx, existingLogger)
package log

import (
	"context"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// 追踪字段常量
const (
	// FieldTraceID OpenTelemetry TraceID 字段名
	FieldTraceID = "trace_id"
	// FieldSpanID OpenTelemetry SpanID 字段名
	FieldSpanID = "span_id"
	// FieldRequestID 业务层 RequestID 字段名
	FieldRequestID = "request_id"
	// FieldService 服务名字段
	FieldService = "service"
	// FieldMethod 方法名字段
	FieldMethod = "method"
	// FieldComponent 组件层级字段
	FieldComponent = "component"
)

// TraceFields 从 context 提取所有追踪字段
// 返回包含 trace_id, span_id, request_id 的 map
func TraceFields(ctx context.Context) map[string]interface{} {
	fields := make(map[string]interface{})

	// 1. 业务层 request_id (从 metainfo 传播)
	if id, ok := metainfo.GetPersistentValue(ctx, FieldRequestID); ok && id != "" {
		fields[FieldRequestID] = id
	}

	// 2. OpenTelemetry trace_id 和 span_id
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		fields[FieldTraceID] = span.SpanContext().TraceID().String()
		fields[FieldSpanID] = span.SpanContext().SpanID().String()
	}

	return fields
}

// GetTraceID 从 context 获取 OpenTelemetry TraceID
func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().TraceID().String()
	}

	return ""
}

// GetSpanID 从 context 获取 OpenTelemetry SpanID
func GetSpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().SpanID().String()
	}

	return ""
}

// GetRequestID 从 context 获取业务层 RequestID
func GetRequestID(ctx context.Context) string {
	if id, ok := metainfo.GetPersistentValue(ctx, FieldRequestID); ok {
		return id
	}

	return ""
}

// Ctx 获取带追踪信息的 logger
// 如果 context 中已有 logger (通过 zerolog.Ctx)，则返回该 logger
// 否则返回带有追踪字段的新 logger
//
// 使用方式:
//
//	log.Ctx(ctx).Info().Msg("hello")
func Ctx(ctx context.Context) *zerolog.Logger {
	logger := zerolog.Ctx(ctx)
	// zerolog.Ctx 返回的 logger 如果未设置，会返回一个 disabled logger
	// 通过检查 GetLevel 来判断是否有效
	if logger != nil && logger.GetLevel() != zerolog.Disabled {
		return logger
	}

	// 创建默认 logger 并附加追踪字段
	l := zerolog.New(nil).With().
		Timestamp().
		Fields(TraceFields(ctx)).
		Logger()

	return &l
}

// WithTrace 为 logger 添加追踪字段
// 返回包含 trace_id, span_id, request_id 的新 logger
//
// 使用方式:
//
//	logger := log.WithTrace(ctx, existingLogger)
//	logger.Info().Msg("hello")
func WithTrace(ctx context.Context, logger zerolog.Logger) zerolog.Logger {
	return logger.With().Ctx(ctx).Logger()
}

// WithTraceAndService 为 logger 添加追踪字段和服务信息
// 返回包含完整追踪上下文的 logger
func WithTraceAndService(ctx context.Context, logger zerolog.Logger, service, method string) zerolog.Logger {
	return logger.With().
		Str(FieldService, service).
		Str(FieldMethod, method).
		Ctx(ctx).
		Logger()
}

// Event 为 zerolog.Event 添加追踪字段（通过 context 传递给 Hook）
// 适用于链式调用场景
//
// 使用方式:
//
//	log.Event(ctx, logger.Info()).
//	    Str("component", "user_handler").
//	    Msg("hello")
func Event(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	return event.Ctx(ctx)
}

// BindToContext 将带追踪信息的 logger 绑定到 context
// 后续代码可通过 zerolog.Ctx(ctx) 或 log.Ctx(ctx) 获取
func BindToContext(ctx context.Context, logger zerolog.Logger, service, method string) context.Context {
	ctxLogger := logger.With().
		Str(FieldService, service).
		Str(FieldMethod, method).
		Ctx(ctx).
		Logger()

	return ctxLogger.WithContext(ctx)
}
