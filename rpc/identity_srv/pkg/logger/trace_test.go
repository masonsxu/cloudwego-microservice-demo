package logger

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestWithRequestID(t *testing.T) {
	t.Run("adds request_id when present", func(t *testing.T) {
		var buf bytes.Buffer

		baseLogger := slog.New(slog.NewJSONHandler(&buf, nil))

		ctx := metainfo.WithPersistentValue(context.Background(), "request_id", "test-request-id")

		logger := WithRequestID(ctx, baseLogger)
		logger.Info("test message")

		output := buf.String()
		assert.Contains(t, output, "test-request-id")
		assert.Contains(t, output, "test message")
	})

	t.Run("returns original logger when request_id is empty", func(t *testing.T) {
		var buf bytes.Buffer

		baseLogger := slog.New(slog.NewJSONHandler(&buf, nil))

		ctx := context.Background()

		logger := WithRequestID(ctx, baseLogger)
		logger.Info("test message")

		output := buf.String()
		assert.NotContains(t, output, "request_id")
		assert.Contains(t, output, "test message")
	})
}

func TestWithRequestIDZerolog(t *testing.T) {
	t.Run("adds request_id when present", func(t *testing.T) {
		var buf bytes.Buffer

		baseLogger := zerolog.New(&buf).With().Timestamp().Logger()

		ctx := metainfo.WithPersistentValue(context.Background(), "request_id", "test-request-id")

		logger := WithRequestIDZerolog(ctx, &baseLogger)
		logger.Info().Msg("test message")

		output := buf.String()
		assert.Contains(t, output, "test-request-id")
		assert.Contains(t, output, "test message")
	})

	t.Run("returns original logger when request_id is empty", func(t *testing.T) {
		var buf bytes.Buffer

		baseLogger := zerolog.New(&buf).With().Timestamp().Logger()

		ctx := context.Background()

		logger := WithRequestIDZerolog(ctx, &baseLogger)
		logger.Info().Msg("test message")

		output := buf.String()
		assert.NotContains(t, output, "request_id")
		assert.Contains(t, output, "test message")
	})

	t.Run("preserves other logger fields", func(t *testing.T) {
		var buf bytes.Buffer

		baseLogger := zerolog.New(&buf).With().
			Str("service", "test-service").
			Timestamp().
			Logger()

		ctx := metainfo.WithPersistentValue(context.Background(), "request_id", "test-request-id")

		logger := WithRequestIDZerolog(ctx, &baseLogger)
		logger.Info().Str("user_id", "123").Msg("test message")

		output := buf.String()
		assert.Contains(t, output, "test-request-id")
		assert.Contains(t, output, "test-service")
		assert.Contains(t, output, "123")
		assert.Contains(t, output, "test message")
	})
}
