// Package logger 提供日志相关的辅助工具
package logger

import (
	"context"
	"log/slog"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/internal/middleware"
	"github.com/rs/zerolog"
)

// WithRequestID 为 slog.Logger 添加 request_id 字段
// 用于业务逻辑层快速创建带追踪信息的日志记录器
//
// 使用示例：
//
//	logger := logger.WithRequestID(ctx, baseLogger)
//	logger.Info("user created", "user_id", userID)
func WithRequestID(ctx context.Context, logger *slog.Logger) *slog.Logger {
	requestID := middleware.GetRequestID(ctx)
	if requestID == "" {
		return logger
	}

	return logger.With("request_id", requestID)
}

// WithRequestIDZerolog 为 zerolog.Logger 添加 request_id 字段
// 用于使用 zerolog 的组件（如 casbin_manager）
//
// 使用示例：
//
//	logger := logger.WithRequestIDZerolog(ctx, cm.logger)
//	logger.Info().Str("user_id", userID).Msg("operation completed")
func WithRequestIDZerolog(ctx context.Context, logger *zerolog.Logger) *zerolog.Logger {
	requestID := middleware.GetRequestID(ctx)
	if requestID == "" {
		return logger
	}

	contextLogger := logger.With().Str("request_id", requestID).Logger()

	return &contextLogger
}
