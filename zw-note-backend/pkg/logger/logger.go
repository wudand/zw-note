package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"zw-note-backend/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var global *zap.Logger

// Init creates a zap logger based on the given configuration and stores it
// as the package-level singleton accessible via L().
func Init(cfg config.LogConfig) (*zap.Logger, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if cfg.Format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	var writeSyncer zapcore.WriteSyncer
	if cfg.Output == "file" && cfg.FilePath != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0o755); err != nil {
			return nil, fmt.Errorf("create log dir: %w", err)
		}
		f, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, fmt.Errorf("open log file: %w", err)
		}
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(f))
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	global = logger
	return logger, nil
}

// L returns the package-level logger, initializing a default development
// logger if Init has not been called yet.
func L() *zap.Logger {
	if global == nil {
		global, _ = zap.NewDevelopment()
	}
	return global
}

// Sync flushes any buffered log entries.
func Sync() {
	if global != nil {
		_ = global.Sync()
	}
}

func parseLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info", "":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unknown log level: %s", level)
	}
}
