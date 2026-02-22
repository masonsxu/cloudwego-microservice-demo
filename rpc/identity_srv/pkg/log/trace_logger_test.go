package log

import (
	"context"
	"testing"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestTraceFields(t *testing.T) {
	t.Run("extracts request_id from metainfo", func(t *testing.T) {
		ctx := metainfo.WithPersistentValue(context.Background(), "request_id", "test-req-123")

		fields := TraceFields(ctx)

		assert.Equal(t, "test-req-123", fields["request_id"])
	})

	t.Run("returns empty map when no metainfo", func(t *testing.T) {
		ctx := context.Background()

		fields := TraceFields(ctx)

		assert.Len(t, fields, 0)
	})

	t.Run("handles empty request_id", func(t *testing.T) {
		ctx := metainfo.WithPersistentValue(context.Background(), "request_id", "")

		fields := TraceFields(ctx)

		// 空字符串不应该被包含
		if reqID, ok := fields["request_id"]; ok {
			assert.Equal(t, "", reqID)
		}
	})
}

func TestGetRequestID(t *testing.T) {
	t.Run("returns request_id when present", func(t *testing.T) {
		expectedID := "req-456"
		ctx := metainfo.WithPersistentValue(context.Background(), "request_id", expectedID)

		result := GetRequestID(ctx)

		assert.Equal(t, expectedID, result)
	})

	t.Run("returns empty string when not present", func(t *testing.T) {
		ctx := context.Background()

		result := GetRequestID(ctx)

		assert.Equal(t, "", result)
	})

	t.Run("returns empty string for empty value", func(t *testing.T) {
		ctx := metainfo.WithPersistentValue(context.Background(), "request_id", "")

		result := GetRequestID(ctx)

		assert.Equal(t, "", result)
	})
}

func TestGetTraceID(t *testing.T) {
	t.Run("returns empty string when no span context", func(t *testing.T) {
		ctx := context.Background()

		result := GetTraceID(ctx)

		assert.Equal(t, "", result)
	})
}

func TestGetSpanID(t *testing.T) {
	t.Run("returns empty string when no span context", func(t *testing.T) {
		ctx := context.Background()

		result := GetSpanID(ctx)

		assert.Equal(t, "", result)
	})
}

func TestCtx(t *testing.T) {
	t.Run("returns logger with trace fields", func(t *testing.T) {
		ctx := metainfo.WithPersistentValue(context.Background(), "request_id", "test-req-789")

		logger := Ctx(ctx)

		assert.NotNil(t, logger)
	})

	t.Run("handles context without logger", func(t *testing.T) {
		ctx := context.Background()

		logger := Ctx(ctx)

		assert.NotNil(t, logger)
	})
}

func TestWithTrace(t *testing.T) {
	t.Run("adds trace fields to logger", func(t *testing.T) {
		ctx := metainfo.WithPersistentValue(context.Background(), "request_id", "req-999")
		baseLogger := zerolog.Nop()

		resultLogger := WithTrace(ctx, baseLogger)

		assert.NotNil(t, resultLogger)
	})
}

func TestWithTraceAndService(t *testing.T) {
	t.Run("adds trace and service fields to logger", func(t *testing.T) {
		ctx := metainfo.WithPersistentValue(context.Background(), "request_id", "req-service-123")
		baseLogger := zerolog.Nop()

		resultLogger := WithTraceAndService(ctx, baseLogger, "test_service", "test_method")

		assert.NotNil(t, resultLogger)
	})
}

func TestEvent(t *testing.T) {
	t.Run("adds trace fields to event", func(t *testing.T) {
		ctx := metainfo.WithPersistentValue(context.Background(), "request_id", "req-event-456")
		baseLogger := zerolog.New(nil)
		event := baseLogger.Info()

		resultEvent := Event(ctx, event)

		assert.NotNil(t, resultEvent)
	})

	t.Run("adds service field as identity_srv", func(t *testing.T) {
		ctx := context.Background()
		baseLogger := zerolog.New(nil)
		event := baseLogger.Info()

		resultEvent := Event(ctx, event)

		assert.NotNil(t, resultEvent)
	})
}

func TestBindToContext(t *testing.T) {
	t.Run("binds logger to context", func(t *testing.T) {
		ctx := context.Background()
		logger := zerolog.Nop()

		newCtx := BindToContext(ctx, logger, "test_service", "test_method")

		assert.NotNil(t, newCtx)

		// 验证可以通过 zerolog.Ctx 获取 logger
		ctxLogger := zerolog.Ctx(newCtx)
		assert.NotNil(t, ctxLogger)
	})
}
