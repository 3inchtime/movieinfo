package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// SimpleLogger 简单日志器实现
type SimpleLogger struct {
	logger *slog.Logger
	level  Level
	fields map[string]interface{}
}

// NewLogger 创建新的日志器
func NewLogger(config *Config) (Logger, error) {
	// 设置输出
	var writer io.Writer
	switch config.Output {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	case "file":
		writer = &lumberjack.Logger{
			Filename:   config.File.Path,
			MaxSize:    config.File.MaxSize,
			MaxBackups: config.File.MaxBackups,
			MaxAge:     config.File.MaxAge,
			Compress:   config.File.Compress,
		}
	default:
		writer = os.Stdout
	}

	// 设置处理器
	var handler slog.Handler
	if config.Format == "json" {
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level: convertLevel(ParseLevel(config.Level)),
		})
	} else {
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{
			Level: convertLevel(ParseLevel(config.Level)),
		})
	}

	return &SimpleLogger{
		logger: slog.New(handler),
		level:  ParseLevel(config.Level),
		fields: make(map[string]interface{}),
	}, nil
}

// convertLevel 转换日志级别
func convertLevel(level Level) slog.Level {
	switch level {
	case DebugLevel:
		return slog.LevelDebug
	case InfoLevel:
		return slog.LevelInfo
	case WarnLevel:
		return slog.LevelWarn
	case ErrorLevel:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// 实现Logger接口
func (l *SimpleLogger) Debug(msg string) {
	l.log(DebugLevel, msg)
}

func (l *SimpleLogger) Info(msg string) {
	l.log(InfoLevel, msg)
}

func (l *SimpleLogger) Warn(msg string) {
	l.log(WarnLevel, msg)
}

func (l *SimpleLogger) Error(msg string) {
	l.log(ErrorLevel, msg)
}

func (l *SimpleLogger) Debugf(format string, args ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(format, args...))
}

func (l *SimpleLogger) Infof(format string, args ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(format, args...))
}

func (l *SimpleLogger) Warnf(format string, args ...interface{}) {
	l.log(WarnLevel, fmt.Sprintf(format, args...))
}

func (l *SimpleLogger) Errorf(format string, args ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(format, args...))
}

func (l *SimpleLogger) WithField(key string, value interface{}) Logger {
	newLogger := &SimpleLogger{
		logger: l.logger,
		level:  l.level,
		fields: make(map[string]interface{}),
	}
	// 复制现有字段
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	// 添加新字段
	newLogger.fields[key] = value
	return newLogger
}

func (l *SimpleLogger) WithError(err error) Logger {
	return l.WithField("error", err.Error())
}

// log 内部日志方法
func (l *SimpleLogger) log(level Level, msg string) {
	if level < l.level {
		return
	}

	// 构建属性
	args := make([]interface{}, 0, len(l.fields)*2)
	for k, v := range l.fields {
		args = append(args, k, v)
	}

	// 记录日志
	switch level {
	case DebugLevel:
		l.logger.Debug(msg, args...)
	case InfoLevel:
		l.logger.Info(msg, args...)
	case WarnLevel:
		l.logger.Warn(msg, args...)
	case ErrorLevel:
		l.logger.Error(msg, args...)
	}
}
