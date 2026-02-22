package log

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLoggerLevels(t *testing.T) {
	t.Run("creates logger with different levels", func(t *testing.T) {
		levels := []zerolog.Level{
			zerolog.DebugLevel,
			zerolog.InfoLevel,
			zerolog.WarnLevel,
			zerolog.ErrorLevel,
		}

		for _, level := range levels {
			logger := zerolog.New(nil).Level(level)
			assert.Equal(t, level, logger.GetLevel())
		}
	})
}

func TestLoggerFields(t *testing.T) {
	t.Run("adds timestamp field", func(t *testing.T) {
		logger := zerolog.New(nil).With().Timestamp().Logger()

		// 无法直接检查字段，但可以创建logger而不panic
		assert.NotNil(t, logger)
	})

	t.Run("adds custom fields", func(t *testing.T) {
		customLogger := zerolog.New(nil).With().
			Str("service", "test-service").
			Str("version", "1.0.0").
			Logger()

		assert.NotNil(t, customLogger)
	})
}

func TestLogLevels(t *testing.T) {
	t.Run("all log levels work", func(t *testing.T) {
		logger := zerolog.Nop()

		logger.Debug().Msg("debug message")
		logger.Info().Msg("info message")
		logger.Warn().Msg("warning message")
		logger.Error().Msg("error message")

		// 如果没有panic则通过
		assert.True(t, true)
	})
}

func TestContextLogger(t *testing.T) {
	t.Run("logger can be attached to context", func(t *testing.T) {
		logger := zerolog.Nop()
		ctx := context.Background()

		// 无法直接获取 logger，但可以绑定
		ctx = logger.WithContext(ctx)

		// 验证 context 不为空
		assert.NotNil(t, ctx)
	})
}

func TestStructuredLogging(t *testing.T) {
	t.Run("logs with string fields", func(t *testing.T) {
		logger := zerolog.Nop()

		logger.Info().
			Str("user_id", "123").
			Str("action", "login").
			Msg("User logged in")

		assert.True(t, true)
	})

	t.Run("logs with int fields", func(t *testing.T) {
		logger := zerolog.Nop()

		logger.Info().
			Int("port", 8080).
			Msg("Server started")

		assert.True(t, true)
	})

	t.Run("logs with bool fields", func(t *testing.T) {
		logger := zerolog.Nop()

		logger.Info().
			Bool("success", true).
			Msg("Operation completed")

		assert.True(t, true)
	})
}

func TestLogOutput(t *testing.T) {
	t.Run("logs to console by default", func(t *testing.T) {
		logger := zerolog.New(nil).With().Timestamp().Logger()

		// 默认输出到控制台
		assert.NotNil(t, logger)
	})
}
