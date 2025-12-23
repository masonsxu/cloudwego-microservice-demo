package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	kitexzerolog "github.com/kitex-contrib/obs-opentelemetry/logging/zerolog"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// CreateLogger 根据配置创建zerolog.Logger实例
// 支持标准输出和文件输出，文件输出时自动启用日志轮转
//
// 日志输出策略：
//   - SERVER_DEBUG=false（生产环境）：只输出到文件（JSON 格式）
//   - SERVER_DEBUG=true（开发环境）：同时输出到终端（美化格式）和文件（JSON 格式）
func CreateLogger(cfg *Config) (*zerolog.Logger, error) {
	// 解析日志级别
	logLevel := cfg.Log.Level

	// 如果开启了调试模式且未显式设置日志级别，自动使用 debug 级别
	if cfg.Server.Debug && logLevel == "" {
		logLevel = "debug"
	}

	var level zerolog.Level

	switch logLevel {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	default:
		level = zerolog.InfoLevel
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
	if cfg.Server.Debug {
		// 调试模式：同时输出到终端（美化格式）和文件（如果配置了）
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}

		if fileWriter != nil {
			// 同时输出到终端和文件
			outputWriter = io.MultiWriter(consoleWriter, fileWriter)
		} else {
			// 只输出到终端
			outputWriter = consoleWriter
		}
	} else {
		// 生产模式：只输出到文件（如果配置了），否则输出到标准输出
		if fileWriter != nil {
			outputWriter = fileWriter
		} else {
			outputWriter = os.Stdout
		}
	}

	// 创建 zerolog logger
	logger := zerolog.New(outputWriter).With().
		Timestamp().
		Caller(). // 添加调用者信息，便于定位日志来源
		Logger().
		Level(level)

	// 记录日志初始化信息
	logger.Info().
		Str("log_level", logLevel).
		Str("format", cfg.Log.Format).
		Str("output", cfg.Log.Output).
		Str("file_path", cfg.Log.FilePath).
		Bool("debug_mode", cfg.Server.Debug).
		Int("max_size_mb", cfg.Log.MaxSize).
		Int("max_age_days", cfg.Log.MaxAge).
		Int("max_backups", cfg.Log.MaxBackups).
		Msg("Logger initialized")

	return &logger, nil
}

// createLogWriter 创建支持轮转的日志writer
// 使用 lumberjack 实现日志轮转，自动处理文件的创建和追加
func createLogWriter(cfg *Config) (*lumberjack.Logger, error) {
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

	// 如果日志文件已存在，检查是否为常规文件并验证权限
	if fileInfo, err := os.Stat(cfg.Log.FilePath); err == nil {
		// 检查是否为目录
		if fileInfo.IsDir() {
			return nil, fmt.Errorf("log path %s is a directory, not a file", cfg.Log.FilePath)
		}

		// 文件存在且是常规文件，尝试以追加模式打开以验证权限
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

// CreateKitexLogger 根据配置创建 Kitex logger 实例
// 使用 kitex-contrib/obs-opentelemetry/logging/zerolog 集成
// 返回的 logger 可以直接通过 klog.SetLogger() 设置到 Kitex
func CreateKitexLogger(cfg *Config) (*kitexzerolog.Logger, error) {
	// 先创建基础的 zerolog.Logger
	zerologLoggerPtr, err := CreateLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create zerolog logger: %w", err)
	}

	// 使用 kitexzerolog.NewLogger 包装，传入自定义的 zerolog.Logger
	// WithLogger 需要 *zerolog.Logger 类型
	kitexLogger := kitexzerolog.NewLogger(
		kitexzerolog.WithLogger(zerologLoggerPtr),
	)

	return kitexLogger, nil
}
