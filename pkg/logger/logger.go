package logger

// Level 日志级别
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// String 返回日志级别的字符串表示
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return "INFO"
	}
}

// ParseLevel 解析日志级别字符串
func ParseLevel(level string) Level {
	switch level {
	case "debug", "DEBUG":
		return DebugLevel
	case "info", "INFO":
		return InfoLevel
	case "warn", "WARN", "warning", "WARNING":
		return WarnLevel
	case "error", "ERROR":
		return ErrorLevel
	default:
		return InfoLevel
	}
}

// Logger 简单日志接口
type Logger interface {
	// 基础日志方法
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)

	// 格式化日志方法
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})

	// 带字段的日志方法
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger
}
