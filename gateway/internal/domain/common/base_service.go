package common

import (
	"context"
	"fmt"

	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	"github.com/rs/zerolog"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/errors"
	tracelog "github.com/masonsxu/cloudwego-microservice-demo/gateway/pkg/log"
)

// BaseService 基础服务结构
// 提供通用的服务方法模板，减少样板代码
type BaseService struct {
	logger          *hertzZerolog.Logger
	responseBuilder *ResponseBuilder
}

// NewBaseService 创建基础服务
func NewBaseService(logger *hertzZerolog.Logger) *BaseService {
	return &BaseService{
		logger:          logger,
		responseBuilder: NewResponseBuilder(),
	}
}

// Logger 获取日志记录器
// Deprecated: 使用 Log(ctx) 替代，Logger() 返回的 hertzZerolog 不支持通过 OTelHook 自动上报 span
func (bs *BaseService) Logger() *hertzZerolog.Logger {
	return bs.logger
}

// Log 获取带追踪上下文的 *zerolog.Logger
// 优先从 context 获取已绑定的 logger（由 TraceMiddleware 通过 BindToContext 注入），
// 回退时将 ctx 绑定到基础 logger，由 OTelHook 在写入时动态注入追踪字段。
// 用法：s.Log(ctx).Error().Err(err).Msg("失败")
func (bs *BaseService) Log(ctx context.Context) *zerolog.Logger {
	return bs.getTracedLogger(ctx)
}

// ResponseBuilder 获取响应构建器
func (bs *BaseService) ResponseBuilder() *ResponseBuilder {
	return bs.responseBuilder
}

// getTracedLogger 获取带追踪信息的 logger（私有方法）
func (bs *BaseService) getTracedLogger(ctx context.Context) *zerolog.Logger {
	// 优先尝试从 context 获取已绑定的 logger（由 TraceMiddleware 绑定）
	logger := tracelog.Ctx(ctx)
	if logger != nil {
		return logger
	}

	// 回退：使用 WithTrace 将 ctx 绑定到基础 logger，OTelHook 动态注入追踪字段
	base := bs.logger.Unwrap()
	traced := tracelog.WithTrace(ctx, base)

	return &traced
}

// ProcessRPCCall 处理RPC调用的通用模板
// 自动处理日志记录、错误处理等重复逻辑
func (bs *BaseService) ProcessRPCCall(
	ctx context.Context,
	logMsg string,
	rpcCall func(ctx context.Context) (interface{}, error),
	logFields ...interface{},
) (interface{}, error) {
	// 记录调用日志
	bs.logWithFields(ctx, zerolog.InfoLevel, logMsg, logFields...)

	// 执行RPC调用
	result, err := rpcCall(ctx)
	if err != nil {
		// 统一错误处理
		bs.logWithFields(
			ctx,
			zerolog.ErrorLevel,
			logMsg+"失败",
			append([]interface{}{"error", err}, logFields...)...)

		return nil, errors.ProcessRPCError(err, logMsg+"失败")
	}

	return result, nil
}

// ProcessRPCVoidCall 处理无返回值的RPC调用
func (bs *BaseService) ProcessRPCVoidCall(
	ctx context.Context,
	logMsg string,
	rpcCall func(ctx context.Context) error,
	logFields ...interface{},
) error {
	// 记录调用日志
	bs.logWithFields(ctx, zerolog.InfoLevel, logMsg, logFields...)

	// 执行RPC调用
	err := rpcCall(ctx)
	if err != nil {
		// 统一错误处理
		bs.logWithFields(
			ctx,
			zerolog.ErrorLevel,
			logMsg+"失败",
			append([]interface{}{"error", err}, logFields...)...)

		return errors.ProcessRPCError(err, logMsg+"失败")
	}

	return nil
}

// LogInfo 便捷的信息日志方法
func (bs *BaseService) LogInfo(ctx context.Context, msg string, fields ...interface{}) {
	bs.logWithFields(ctx, zerolog.InfoLevel, msg, fields...)
}

// LogError 便捷的错误日志方法
func (bs *BaseService) LogError(ctx context.Context, msg string, err error, fields ...interface{}) {
	allFields := append([]interface{}{"error", err}, fields...)
	bs.logWithFields(ctx, zerolog.ErrorLevel, msg, allFields...)
}

// logWithFields 辅助方法：将字段转换为 zerolog 的链式调用
func (bs *BaseService) logWithFields(
	ctx context.Context,
	level zerolog.Level,
	msg string,
	fields ...interface{},
) {
	logger := bs.getTracedLogger(ctx)

	// 先根据级别创建 event，ctx 已通过 getTracedLogger 绑定到 logger.ctx，
	// zerolog 创建 event 时会自动将 logger.ctx 传播到 event，OTelHook 因此能拿到 span。
	var event *zerolog.Event

	switch level {
	case zerolog.ErrorLevel:
		event = logger.Error()
	case zerolog.WarnLevel:
		event = logger.Warn()
	case zerolog.DebugLevel:
		event = logger.Debug()
	default:
		event = logger.Info()
	}

	// 处理字段对 (key, value)
	for i := 0; i < len(fields)-1; i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			key = fmt.Sprintf("field_%d", i)
		}

		value := fields[i+1]

		// 根据值的类型选择合适的方法
		switch v := value.(type) {
		case string:
			event = event.Str(key, v)
		case int:
			event = event.Int(key, v)
		case int64:
			event = event.Int64(key, v)
		case bool:
			event = event.Bool(key, v)
		case error:
			event = event.Err(v)
		default:
			event = event.Interface(key, v)
		}
	}

	event.Msg(msg)
}
