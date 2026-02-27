// Package log 提供统一的日志工具，支持 OpenTelemetry 追踪信息注入
//
// 本包解决日志与链路追踪对齐的问题，确保每条日志都包含：
// - trace_id: OpenTelemetry 链路追踪 ID
// - span_id: OpenTelemetry Span ID
// - request_id: 业务层请求 ID (从 metainfo 传播)
//
// 追踪字段由 OTelHook 在日志写入时动态注入（非静态固化），需通过 Ctx(ctx) 或
// event.Ctx(ctx) 将 context 绑定到 logger/event，OTelHook 才能读取 span 信息。
//
// 使用示例:
//
//	// 在业务代码中获取带追踪信息的 logger（TraceMiddleware 已绑定到 ctx）
//	logger := log.Ctx(ctx)
//	logger.Info().Str("user_id", "123").Msg("User created")
//
//	// 链式调用场景
//	log.Event(ctx, logger.Error()).Err(err).Msg("操作失败")
package log

import (
	"context"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
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
// 优先从 context 获取已绑定的 logger（由 TraceMiddleware 通过 BindToContext 绑定），
// 回退时使用 WithTrace 将 ctx 绑定到基础 logger，由 OTelHook 在写入时动态注入追踪字段。
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

	return nil
}

// WithTrace 为 logger 绑定 context，追踪字段由 OTelHook 在日志写入时动态注入
//
// 使用方式:
//
//	logger := log.WithTrace(ctx, existingLogger)
//	logger.Info().Msg("hello")
func WithTrace(ctx context.Context, logger zerolog.Logger) zerolog.Logger {
	return logger.With().Ctx(ctx).Logger()
}

// WithTraceAndService 为 logger 绑定 context 和服务信息
// 追踪字段由 OTelHook 在日志写入时动态注入（而非静态固化到 logger 中）
//
// Deprecated: 使用 BindToContext 替代
func WithTraceAndService(
	ctx context.Context,
	logger zerolog.Logger,
	service, method string,
) zerolog.Logger {
	return logger.With().
		Str(FieldService, service).
		Str(FieldMethod, method).
		Ctx(ctx).
		Logger()
}

// Event 为 zerolog.Event 绑定 context，追踪字段由 OTelHook 在日志写入时动态注入
// 这是让 OTelHook 能够调用 span.RecordError()/span.SetStatus() 的关键——
// OTelHook.Run() 通过 e.GetCtx() 获取 span，必须先调用 event.Ctx(ctx)。
//
// 使用方式:
//
//	log.Event(ctx, logger.Error()).
//	    Err(err).
//	    Str("component", "user_handler").
//	    Msg("操作失败")
func Event(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	return event.Ctx(ctx).Str(FieldService, "gateway")
}

// BindToContext 将带服务信息的 logger 绑定到 context
// 追踪字段由 OTelHook 在日志写入时动态注入（而非静态固化到 logger 中）
// 后续代码可通过 zerolog.Ctx(ctx) 或 log.Ctx(ctx) 获取
func BindToContext(
	ctx context.Context,
	logger zerolog.Logger,
	service, method string,
) context.Context {
	ctxLogger := logger.With().
		Str(FieldService, service).
		Str(FieldMethod, method).
		Ctx(ctx).
		Logger()
	return ctxLogger.WithContext(ctx)
}

// RecordSpanHTTPError 将 HTTP 错误的详细信息作为 span attributes 直接写入 OTel Span
//
// OTelHook 只将日志消息（msg）传给 span.RecordError()/span.AddEvent()，
// zerolog 的结构化字段不会自动成为 span attributes。
// 本函数补充这一缺口，将 status_code、path、method、error_body 等附加到 span，
// 使其在 Jaeger span 详情中可见。
//
// 4xx 使用 span.AddEvent + attributes（不标红），5xx 使用 span.SetAttributes + RecordError（标红）。
func RecordSpanHTTPError(
	ctx context.Context,
	statusCode int,
	method, path, errorBody string,
) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	attrs := []attribute.KeyValue{
		attribute.Int("http.status_code", statusCode),
		attribute.String("http.method", method),
		attribute.String("http.path", path),
	}
	if errorBody != "" {
		attrs = append(attrs, attribute.String("error.body", errorBody))
	}

	if statusCode >= 500 {
		// 5xx：将属性附加到 span 本身，配合 OTelHook 的 RecordError 一起在 Jaeger 中显示
		span.SetAttributes(attrs...)
	} else {
		// 4xx：以具名 span event 的形式记录，不覆盖 OTelHook 已添加的事件
		span.AddEvent("http.client_error", trace.WithAttributes(attrs...))
	}
}
