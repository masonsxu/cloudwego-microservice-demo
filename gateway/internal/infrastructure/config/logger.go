package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"

	tracelog "github.com/masonsxu/cloudwego-microservice-demo/gateway/pkg/log"
)

// CreateLogger 根据配置创建zerolog.Logger实例
// 支持标准输出和文件输出，文件输出时自动启用日志轮转
//
// 日志输出策略：
//   - SERVER_DEBUG=false（生产环境）：只输出到文件（JSON 格式）
//   - SERVER_DEBUG=true（开发环境）：同时输出到终端（美化格式）和文件（JSON 格式）
func CreateLogger(cfg *Configuration) (*hertzZerolog.Logger, error) {
	// 设置 zerolog 全局时间格式（影响 JSON 输出）
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"

	// 解析日志级别
	var level hlog.Level

	logLevel := cfg.Log.Level

	// 如果开启了调试模式且未显式设置日志级别，自动使用 debug 级别
	if cfg.Server.Debug && logLevel == "" {
		logLevel = "debug"
	}

	switch logLevel {
	case "debug":
		level = hlog.LevelDebug
	case "info":
		level = hlog.LevelInfo
	case "warn":
		level = hlog.LevelWarn
	case "error":
		level = hlog.LevelError
	default:
		level = hlog.LevelInfo
	}

	// 构建输出目标
	var outputWriter io.Writer

	// 创建文件 writer（如果配置了文件输出）
	var fileWriter io.Writer

	if cfg.Log.Output == "file" && cfg.Log.FilePath != "" {
		fw, err := createLogWriter(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create log writer: %w", err)
		}

		fileWriter = fw
	}

	// 根据调试模式决定输出目标
	outputWriter = buildOutputWriter(cfg.Server.Debug, fileWriter)

	// 创建 zerolog logger 选项
	opts := []hertzZerolog.Opt{
		hertzZerolog.WithLevel(level),
		hertzZerolog.WithOutput(outputWriter),
		hertzZerolog.WithTimestamp(),
		// 使用 CallerSkipFrameCount 来跳过 hertz-zerolog 库内部的调用栈层数
		// 这样 caller 显示的是业务代码位置而不是库内部位置
		// 值为 4 可以跳过：hertz-zerolog 内部封装层 -> zerolog Event 层
		hertzZerolog.WithCallerSkipFrameCount(4),
		// 安装 OTel Hook：自动注入 trace_id/span_id/request_id，
		// 并将 Error/Warn 级别日志同步到 Jaeger Span
		hertzZerolog.WithHook(tracelog.NewOTelHook()),
	}

	// 创建 logger
	logger := hertzZerolog.New(opts...)

	// 记录日志初始化信息（使用 hlog 接口）
	logger.Infof(
		"Logger initialized: level=%s, format=%s, output=%s, file_path=%s, "+
			"debug_mode=%v, max_size_mb=%d, max_age_days=%d, max_backups=%d",
		logLevel,
		cfg.Log.Format,
		cfg.Log.Output,
		cfg.Log.FilePath,
		cfg.Server.Debug,
		cfg.Log.MaxSize,
		cfg.Log.MaxAge,
		cfg.Log.MaxBackups,
	)

	return logger, nil
}

// buildOutputWriter 根据调试模式和文件writer构建输出目标
func buildOutputWriter(debugMode bool, fileWriter io.Writer) io.Writer {
	if debugMode {
		// 调试模式：同时输出到终端（美化格式）和文件（如果配置了）
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05",
			NoColor:    false,
		}

		if fileWriter != nil {
			return io.MultiWriter(consoleWriter, fileWriter)
		}

		return consoleWriter
	}

	// 生产模式：只输出到文件（如果配置了），否则输出到标准输出
	if fileWriter != nil {
		return fileWriter
	}

	return os.Stdout
}

// createLogWriter 创建支持轮转的日志writer
// 使用 lumberjack 实现日志轮转，自动处理文件的创建和追加
func createLogWriter(cfg *Configuration) (*lumberjack.Logger, error) {
	// 确保日志目录存在
	logDir := filepath.Dir(cfg.Log.FilePath)

	// 检查路径是否存在以及它的类型
	if stat, err := os.Stat(logDir); err != nil {
		// 路径不存在，创建目录
		if err := os.MkdirAll(logDir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create log directory %s: %w", logDir, err)
		}
	} else if !stat.IsDir() {
		// 路径存在但不是目录，这是一个错误状态
		return nil, fmt.Errorf("log directory path %s exists but is not a directory", logDir)
	}

	// 如果日志文件已存在，检查权限
	if _, err := os.Stat(cfg.Log.FilePath); err == nil {
		// 文件存在，尝试以追加模式打开以验证权限
		file, err := os.OpenFile(cfg.Log.FilePath, os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot write to existing log file %s: %w",
				cfg.Log.FilePath,
				err,
			)
		}

		file.Close()
	}

	// 创建 lumberjack logger
	// lumberjack 会自动以追加模式打开现有文件
	writer := &lumberjack.Logger{
		Filename:   cfg.Log.FilePath,
		MaxSize:    cfg.Log.MaxSize,    // 单个文件最大尺寸（MB）
		MaxAge:     cfg.Log.MaxAge,     // 文件最大保存天数
		MaxBackups: cfg.Log.MaxBackups, // 最多保留文件数
		LocalTime:  true,               // 使用本地时间命名备份文件
		Compress:   true,               // 压缩旧日志文件
	}

	return writer, nil
}
