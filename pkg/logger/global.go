package logger

import (
	"sync"
)

var (
	globalLogger Logger
	mu           sync.RWMutex
)

// InitGlobalLogger 初始化全局日志器
func InitGlobalLogger(config *Config) error {
	logger, err := NewLogger(config)
	if err != nil {
		return err
	}

	mu.Lock()
	globalLogger = logger
	mu.Unlock()

	return nil
}

// GetGlobalLogger 获取全局日志器
func GetGlobalLogger() Logger {
	mu.RLock()
	defer mu.RUnlock()
	return globalLogger
}

// 全局日志方法
func Debug(msg string) {
	if globalLogger != nil {
		globalLogger.Debug(msg)
	}
}

func Info(msg string) {
	if globalLogger != nil {
		globalLogger.Info(msg)
	}
}

func Warn(msg string) {
	if globalLogger != nil {
		globalLogger.Warn(msg)
	}
}

func Error(msg string) {
	if globalLogger != nil {
		globalLogger.Error(msg)
	}
}

func Debugf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Debugf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Infof(format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Warnf(format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Errorf(format, args...)
	}
}

func WithField(key string, value interface{}) Logger {
	if globalLogger != nil {
		return globalLogger.WithField(key, value)
	}
	return nil
}

func WithError(err error) Logger {
	if globalLogger != nil {
		return globalLogger.WithError(err)
	}
	return nil
}
