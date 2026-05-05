package config

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// CreateLogger 创建 zerolog logger
func CreateLogger(cfg *Config) (*zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	var writers []io.Writer

	if cfg.Log.Output == "file" || cfg.Log.Output == "both" {
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.Log.FilePath,
			MaxSize:    cfg.Log.MaxSize,
			MaxAge:     cfg.Log.MaxAge,
			MaxBackups: cfg.Log.MaxBackups,
		}
		writers = append(writers, fileWriter)
	}

	if cfg.Log.Output == "console" || cfg.Log.Output == "both" {
		writers = append(writers, os.Stdout)
	}

	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	multiWriter := io.MultiWriter(writers...)
	logger := zerolog.New(multiWriter).Level(level).With().Timestamp().Caller().Logger()

	return &logger, nil
}
