package middleware

import (
	"bytes"
	"context"
	"testing"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetaInfoMiddleware(t *testing.T) {
	t.Run("with custom logger", func(t *testing.T) {
		logger := zerolog.Nop()
		middleware := NewMetaInfoMiddleware(&logger)

		assert.NotNil(t, middleware)
		assert.Equal(t, &logger, middleware.logger)
	})

	t.Run("with nil logger uses default", func(t *testing.T) {
		middleware := NewMetaInfoMiddleware(nil)

		assert.NotNil(t, middleware)
		assert.NotNil(t, middleware.logger)
	})
}

func TestMetaInfoMiddleware_ServerMiddleware(t *testing.T) {
	tests := []struct {
		name            string
		existingMeta    map[string]string
		expectGenerated bool
		validateFunc    func(*testing.T, context.Context, *bytes.Buffer)
	}{
		{
			name: "request_id exists - no generation needed",
			existingMeta: map[string]string{
				"request_id": "existing-request-id",
			},
			expectGenerated: false,
			validateFunc: func(t *testing.T, ctx context.Context, logBuf *bytes.Buffer) {
				assert.Equal(t, "existing-request-id", GetRequestID(ctx))

				// Should not contain warning about generation
				logOutput := logBuf.String()
				assert.NotContains(t, logOutput, "Generated missing request_id")
			},
		},
		{
			name:            "no request_id - generate one",
			existingMeta:    map[string]string{},
			expectGenerated: true,
			validateFunc: func(t *testing.T, ctx context.Context, logBuf *bytes.Buffer) {
				requestID := GetRequestID(ctx)

				assert.NotEmpty(t, requestID)

				// Should contain warning about generation
				logOutput := logBuf.String()
				assert.Contains(t, logOutput, "Generated missing request_id")
				assert.Contains(t, logOutput, requestID)
			},
		},
		{
			name: "empty string request_id - should regenerate",
			existingMeta: map[string]string{
				"request_id": "",
			},
			expectGenerated: true,
			validateFunc: func(t *testing.T, ctx context.Context, logBuf *bytes.Buffer) {
				requestID := GetRequestID(ctx)

				assert.NotEmpty(t, requestID)

				// Should contain warning about generation
				logOutput := logBuf.String()
				assert.Contains(t, logOutput, "Generated missing request_id")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create logger with buffer to capture logs
			var logBuf bytes.Buffer

			logger := zerolog.New(&logBuf).With().Timestamp().Logger()

			middleware := NewMetaInfoMiddleware(&logger)
			serverMiddleware := middleware.ServerMiddleware()

			// Create context with existing meta info
			ctx := createContextWithMeta(tt.existingMeta)

			// Mock endpoint to capture final context
			var finalCtx context.Context

			mockEndpoint := func(ctx context.Context, req, resp interface{}) error {
				finalCtx = ctx
				return nil
			}

			// Apply middleware
			wrappedEndpoint := serverMiddleware(mockEndpoint)
			err := wrappedEndpoint(ctx, nil, nil)

			// Verify no error
			require.NoError(t, err)
			require.NotNil(t, finalCtx)

			// Run validation
			tt.validateFunc(t, finalCtx, &logBuf)
		})
	}
}

func TestMetaInfoMiddleware_ensureRequestID(t *testing.T) {
	logger := zerolog.Nop()
	middleware := NewMetaInfoMiddleware(&logger)

	t.Run("generates valid UUID", func(t *testing.T) {
		ctx := context.Background()

		resultCtx := middleware.ensureRequestID(ctx)

		requestID := GetRequestID(resultCtx)

		// Verify UUID is generated and valid format
		assert.NotEmpty(t, requestID)

		// Basic UUID format check (36 characters with hyphens)
		assert.Len(t, requestID, 36)
		assert.Contains(t, requestID, "-")
	})

	t.Run("preserves existing request_id", func(t *testing.T) {
		ctx := createContextWithMeta(map[string]string{
			"request_id": "preserve-me",
		})

		resultCtx := middleware.ensureRequestID(ctx)

		assert.Equal(t, "preserve-me", GetRequestID(resultCtx))
	})
}

func TestMetaInfoMiddleware_logTraceInfo(t *testing.T) {
	var logBuf bytes.Buffer

	logger := zerolog.New(&logBuf).With().Timestamp().Logger()

	middleware := NewMetaInfoMiddleware(&logger)

	t.Run("logs trace info when request_id present", func(t *testing.T) {
		ctx := createContextWithMeta(map[string]string{
			"request_id": "test-request-id",
		})

		middleware.logTraceInfo(ctx)

		logOutput := logBuf.String()
		assert.Contains(t, logOutput, "RPC request started")
		assert.Contains(t, logOutput, "test-request-id")
		assert.Contains(t, logOutput, `"service":"identity_srv"`)
		assert.Contains(t, logOutput, `"method":"unknown"`)
	})

	t.Run("handles empty context gracefully", func(t *testing.T) {
		logBuf.Reset()

		ctx := context.Background()

		middleware.logTraceInfo(ctx)

		// Should not log anything when no trace info is available
		logOutput := logBuf.String()
		assert.Empty(t, logOutput)
	})
}

func TestGetRequestID(t *testing.T) {
	t.Run("returns ID when present", func(t *testing.T) {
		ctx := createContextWithMeta(map[string]string{
			"request_id": "test-request-id",
		})

		result := GetRequestID(ctx)
		assert.Equal(t, "test-request-id", result)
	})

	t.Run("returns empty string when not present", func(t *testing.T) {
		ctx := context.Background()

		result := GetRequestID(ctx)
		assert.Equal(t, "", result)
	})

	t.Run("returns empty string for empty value", func(t *testing.T) {
		ctx := createContextWithMeta(map[string]string{
			"request_id": "",
		})

		result := GetRequestID(ctx)
		assert.Equal(t, "", result)
	})
}

func TestLoggingAttrs(t *testing.T) {
	t.Run("returns attributes for request_id", func(t *testing.T) {
		ctx := createContextWithMeta(map[string]string{
			"request_id": "test-request-id",
		})

		attrs := LoggingAttrs(ctx)

		assert.Len(t, attrs, 1)
		assert.Equal(t, "test-request-id", attrs["request_id"])
	})

	t.Run("returns empty map when no request_id present", func(t *testing.T) {
		ctx := context.Background()

		attrs := LoggingAttrs(ctx)

		assert.Len(t, attrs, 0)
	})
}

func TestMetaInfoMiddleware_Integration(t *testing.T) {
	t.Run("complete middleware flow with ID generation", func(t *testing.T) {
		var logBuf bytes.Buffer

		logger := zerolog.New(&logBuf).With().Timestamp().Logger()

		middleware := NewMetaInfoMiddleware(&logger)
		serverMiddleware := middleware.ServerMiddleware()

		// Start with empty context
		ctx := context.Background()

		// Business logic that uses the request_id
		businessLogic := func(ctx context.Context, req, resp interface{}) error {
			// Verify request_id is available in business logic
			requestID := GetRequestID(ctx)

			assert.NotEmpty(t, requestID)

			// Verify logging attributes work
			attrs := LoggingAttrs(ctx)
			assert.Len(t, attrs, 1)
			assert.Equal(t, requestID, attrs["request_id"])

			return nil
		}

		// Apply middleware and execute
		wrappedEndpoint := serverMiddleware(businessLogic)
		err := wrappedEndpoint(ctx, nil, nil)

		require.NoError(t, err)

		// Verify logging occurred
		logOutput := logBuf.String()
		assert.Contains(t, logOutput, "Generated missing request_id")
		assert.Contains(t, logOutput, "RPC request started")
		assert.Contains(t, logOutput, "RPC request completed")
		assert.Contains(t, logOutput, "identity_srv")
		assert.Contains(t, logOutput, `"status":"success"`)
	})

	t.Run("middleware handles metainfo correctly", func(t *testing.T) {
		logger := zerolog.Nop()
		middleware := NewMetaInfoMiddleware(&logger)
		serverMiddleware := middleware.ServerMiddleware()

		// Create context with metainfo values
		ctx := createContextWithMeta(map[string]string{
			"request_id": "existing-id",
		})

		businessLogic := func(ctx context.Context, req, resp interface{}) error {
			// Verify request_id is preserved
			assert.Equal(t, "existing-id", GetRequestID(ctx))

			return nil
		}

		wrappedEndpoint := serverMiddleware(businessLogic)
		err := wrappedEndpoint(ctx, nil, nil)

		require.NoError(t, err)
	})
}

// Helper functions

// createContextWithMeta creates a context with metainfo values
func createContextWithMeta(metaInfo map[string]string) context.Context {
	ctx := context.Background()

	for key, value := range metaInfo {
		ctx = metainfo.WithPersistentValue(ctx, key, value)
	}

	return ctx
}

// Benchmark tests for performance validation

func BenchmarkMetaInfoMiddleware_WithExistingID(b *testing.B) {
	logger := zerolog.Nop()
	middleware := NewMetaInfoMiddleware(&logger)
	serverMiddleware := middleware.ServerMiddleware()

	ctx := createContextWithMeta(map[string]string{
		"request_id": "benchmark-request-id",
	})

	mockEndpoint := func(ctx context.Context, req, resp interface{}) error {
		return nil
	}

	wrappedEndpoint := serverMiddleware(mockEndpoint)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = wrappedEndpoint(ctx, nil, nil)
	}
}

func BenchmarkMetaInfoMiddleware_WithGeneration(b *testing.B) {
	logger := zerolog.Nop()
	middleware := NewMetaInfoMiddleware(&logger)
	serverMiddleware := middleware.ServerMiddleware()

	ctx := context.Background()

	mockEndpoint := func(ctx context.Context, req, resp interface{}) error {
		return nil
	}

	wrappedEndpoint := serverMiddleware(mockEndpoint)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = wrappedEndpoint(ctx, nil, nil)
	}
}

func BenchmarkGetRequestID(b *testing.B) {
	ctx := createContextWithMeta(map[string]string{
		"request_id": "benchmark-request-id",
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = GetRequestID(ctx)
	}
}

func BenchmarkLoggingAttrs(b *testing.B) {
	ctx := createContextWithMeta(map[string]string{
		"request_id": "benchmark-request-id",
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = LoggingAttrs(ctx)
	}
}
