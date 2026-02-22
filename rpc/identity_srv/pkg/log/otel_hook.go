// Package log 提供 OTel zerolog Hook，自动将追踪信息注入日志并同步到 Jaeger
//
// OTelHook 实现 zerolog.Hook 接口，在每条日志事件触发时：
//   - 从 OTel Span 动态注入 trace_id/span_id/trace_flags
//   - 从 metainfo 注入 request_id
//   - Error 级别日志：调用 span.RecordError() + span.SetStatus()，Jaeger Span 标红
//   - Warn 级别日志：调用 span.AddEvent()，Jaeger Span 详情中可见
package log

import (
	"errors"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// OTelHook 是一个 zerolog.Hook 实现，负责：
// 1. 在日志事件中自动注入 OTel 追踪字段（trace_id, span_id, trace_flags）
// 2. 在日志事件中自动注入 request_id（从 metainfo 获取）
// 3. 将 Error/Warn 级别日志同步到 OTel Span（Jaeger 可见）
type OTelHook struct {
	// errorSpanLevel 及以上级别的日志会调用 span.RecordError() + span.SetStatus(Error)
	errorSpanLevel zerolog.Level
	// spanEventLevel 及以上且低于 errorSpanLevel 的日志会调用 span.AddEvent()
	spanEventLevel zerolog.Level
	// recordStackTrace 是否在 RecordError 时记录堆栈
	recordStackTrace bool
}

// OTelHookOption 配置 OTelHook 的选项函数
type OTelHookOption func(*OTelHook)

// WithErrorSpanLevel 设置触发 span.RecordError() 的最低日志级别
// 默认为 zerolog.ErrorLevel
func WithErrorSpanLevel(level zerolog.Level) OTelHookOption {
	return func(h *OTelHook) {
		h.errorSpanLevel = level
	}
}

// WithSpanEventLevel 设置触发 span.AddEvent() 的最低日志级别
// 默认为 zerolog.WarnLevel
func WithSpanEventLevel(level zerolog.Level) OTelHookOption {
	return func(h *OTelHook) {
		h.spanEventLevel = level
	}
}

// WithRecordStackTrace 设置是否在 RecordError 时记录堆栈信息
// 默认为 true
func WithRecordStackTrace(enabled bool) OTelHookOption {
	return func(h *OTelHook) {
		h.recordStackTrace = enabled
	}
}

// NewOTelHook 创建 OTelHook 实例
func NewOTelHook(opts ...OTelHookOption) *OTelHook {
	h := &OTelHook{
		errorSpanLevel:   zerolog.ErrorLevel,
		spanEventLevel:   zerolog.WarnLevel,
		recordStackTrace: true,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// Run 实现 zerolog.Hook 接口
// 在每次日志事件写入时被调用
func (h *OTelHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	if ctx == nil {
		return
	}

	// 1. 从 OTel Span 注入 trace_id / span_id / trace_flags
	span := trace.SpanFromContext(ctx)

	spanCtx := span.SpanContext()
	if spanCtx.IsValid() {
		e.Str(FieldTraceID, spanCtx.TraceID().String())
		e.Str(FieldSpanID, spanCtx.SpanID().String())
		e.Str("trace_flags", spanCtx.TraceFlags().String())
	}

	// 2. 从 metainfo 注入 request_id（官方 Hook 不处理）
	if id, ok := metainfo.GetPersistentValue(ctx, FieldRequestID); ok && id != "" {
		e.Str(FieldRequestID, id)
	}

	// 3. 只有 Span 正在记录时才写入 Span Event
	if !span.IsRecording() {
		return
	}

	// 4. Error 及以上级别 → span.RecordError() + span.SetStatus(Error)
	//    让 Jaeger 中 Span 标红，并在详情中显示错误事件
	if level >= h.errorSpanLevel {
		span.SetStatus(codes.Error, msg)
		span.RecordError(
			errors.New(msg),
			trace.WithStackTrace(h.recordStackTrace),
		)

		return
	}

	// 5. Warn 级别（>= spanEventLevel 且 < errorSpanLevel）→ span.AddEvent()
	//    Jaeger Span 详情中可见，但不标记为错误
	if level >= h.spanEventLevel {
		span.AddEvent(msg, trace.WithAttributes(
			attribute.String("log.level", level.String()),
		))
	}
}
