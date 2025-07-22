# 3.4 日志系统搭建

## 概述

日志系统是应用程序的重要基础设施，它为问题诊断、性能监控、安全审计提供了关键信息。对于MovieInfo项目，我们需要构建一个统一、高效、可扩展的日志系统，支持结构化日志、多级别输出和集中化管理。

## 为什么需要专业的日志系统？

### 1. **问题诊断**
- **错误追踪**：快速定位和分析系统错误
- **性能分析**：识别性能瓶颈和优化点
- **用户行为**：了解用户操作路径和习惯
- **系统状态**：监控系统运行状态和健康度

### 2. **运维支持**
- **故障排查**：提供详细的故障信息和上下文
- **容量规划**：基于日志数据进行容量规划
- **趋势分析**：分析系统使用趋势和模式
- **告警触发**：基于日志内容触发告警

### 3. **安全审计**
- **访问记录**：记录用户访问和操作行为
- **安全事件**：记录安全相关的事件和异常
- **合规要求**：满足法规和合规性要求
- **取证支持**：为安全事件调查提供证据

### 4. **业务分析**
- **用户画像**：基于日志数据分析用户特征
- **功能使用**：统计功能使用频率和效果
- **业务指标**：计算关键业务指标和KPI
- **A/B测试**：支持功能测试和效果评估

## 日志系统架构设计

### 1. **日志层次结构**

```
日志系统架构
├── 应用层日志 (Application)
│   ├── 业务日志 (Business)
│   ├── 访问日志 (Access)
│   └── 错误日志 (Error)
├── 中间件日志 (Middleware)
│   ├── 认证日志 (Auth)
│   ├── 限流日志 (RateLimit)
│   └── CORS日志 (CORS)
├── 基础设施日志 (Infrastructure)
│   ├── 数据库日志 (Database)
│   ├── 缓存日志 (Cache)
│   └── gRPC日志 (gRPC)
└── 系统日志 (System)
    ├── 启动日志 (Startup)
    ├── 配置日志 (Config)
    └── 健康检查日志 (Health)
```

### 2. **日志级别定义**

```go
// pkg/logger/level.go
package logger

// Level 日志级别
type Level int

const (
    DebugLevel Level = iota // 调试信息
    InfoLevel              // 一般信息
    WarnLevel              // 警告信息
    ErrorLevel             // 错误信息
    FatalLevel             // 致命错误
    PanicLevel             // 恐慌级别
)

// String 返回日志级别字符串
func (l Level) String() string {
    switch l {
    case DebugLevel:
        return "debug"
    case InfoLevel:
        return "info"
    case WarnLevel:
        return "warn"
    case ErrorLevel:
        return "error"
    case FatalLevel:
        return "fatal"
    case PanicLevel:
        return "panic"
    default:
        return "unknown"
    }
}
```

### 3. **日志接口设计**

#### 3.1 统一日志接口
```go
// pkg/logger/interface.go
package logger

import (
    "context"
)

// Logger 日志接口
type Logger interface {
    // 基础日志方法
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Fatal(msg string, fields ...Field)
    Panic(msg string, fields ...Field)
    
    // 格式化日志方法
    Debugf(format string, args ...interface{})
    Infof(format string, args ...interface{})
    Warnf(format string, args ...interface{})
    Errorf(format string, args ...interface{})
    Fatalf(format string, args ...interface{})
    Panicf(format string, args ...interface{})
    
    // 上下文日志方法
    WithContext(ctx context.Context) Logger
    WithFields(fields ...Field) Logger
    WithError(err error) Logger
    
    // 配置方法
    SetLevel(level Level)
    GetLevel() Level
}

// Field 日志字段
type Field struct {
    Key   string
    Value interface{}
}

// 便捷的字段创建函数
func String(key, value string) Field {
    return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
    return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
    return Field{Key: key, Value: value}
}

func Float64(key string, value float64) Field {
    return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field {
    return Field{Key: key, Value: value}
}

func Any(key string, value interface{}) Field {
    return Field{Key: key, Value: value}
}

func Error(err error) Field {
    return Field{Key: "error", Value: err}
}

func Duration(key string, value time.Duration) Field {
    return Field{Key: key, Value: value}
}
```

### 4. **Zap日志实现**

#### 4.1 Zap日志器实现
```go
// pkg/logger/zap.go
package logger

import (
    "context"
    "os"
    "time"
    
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
)

// ZapLogger Zap日志器实现
type ZapLogger struct {
    logger *zap.Logger
    level  Level
}

// NewZapLogger 创建Zap日志器
func NewZapLogger(config LogConfig) (Logger, error) {
    // 创建编码器配置
    encoderConfig := zapcore.EncoderConfig{
        TimeKey:        "timestamp",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        MessageKey:     "message",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.LowercaseLevelEncoder,
        EncodeTime:     zapcore.ISO8601TimeEncoder,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    }
    
    // 创建编码器
    var encoder zapcore.Encoder
    if config.Format == "json" {
        encoder = zapcore.NewJSONEncoder(encoderConfig)
    } else {
        encoder = zapcore.NewConsoleEncoder(encoderConfig)
    }
    
    // 创建输出
    var writeSyncer zapcore.WriteSyncer
    switch config.Output {
    case "stdout":
        writeSyncer = zapcore.AddSync(os.Stdout)
    case "stderr":
        writeSyncer = zapcore.AddSync(os.Stderr)
    case "file":
        writeSyncer = zapcore.AddSync(&lumberjack.Logger{
            Filename:   config.FilePath,
            MaxSize:    config.MaxSize,
            MaxBackups: config.MaxBackups,
            MaxAge:     config.MaxAge,
            Compress:   config.Compress,
        })
    default:
        writeSyncer = zapcore.AddSync(os.Stdout)
    }
    
    // 设置日志级别
    level := parseLevel(config.Level)
    
    // 创建核心
    core := zapcore.NewCore(encoder, writeSyncer, level)
    
    // 创建日志器
    logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
    
    return &ZapLogger{
        logger: logger,
        level:  parseLogLevel(config.Level),
    }, nil
}

// Debug 调试日志
func (l *ZapLogger) Debug(msg string, fields ...Field) {
    l.logger.Debug(msg, l.convertFields(fields)...)
}

// Info 信息日志
func (l *ZapLogger) Info(msg string, fields ...Field) {
    l.logger.Info(msg, l.convertFields(fields)...)
}

// Warn 警告日志
func (l *ZapLogger) Warn(msg string, fields ...Field) {
    l.logger.Warn(msg, l.convertFields(fields)...)
}

// Error 错误日志
func (l *ZapLogger) Error(msg string, fields ...Field) {
    l.logger.Error(msg, l.convertFields(fields)...)
}

// Fatal 致命错误日志
func (l *ZapLogger) Fatal(msg string, fields ...Field) {
    l.logger.Fatal(msg, l.convertFields(fields)...)
}

// Panic 恐慌日志
func (l *ZapLogger) Panic(msg string, fields ...Field) {
    l.logger.Panic(msg, l.convertFields(fields)...)
}

// Debugf 格式化调试日志
func (l *ZapLogger) Debugf(format string, args ...interface{}) {
    l.logger.Sugar().Debugf(format, args...)
}

// Infof 格式化信息日志
func (l *ZapLogger) Infof(format string, args ...interface{}) {
    l.logger.Sugar().Infof(format, args...)
}

// Warnf 格式化警告日志
func (l *ZapLogger) Warnf(format string, args ...interface{}) {
    l.logger.Sugar().Warnf(format, args...)
}

// Errorf 格式化错误日志
func (l *ZapLogger) Errorf(format string, args ...interface{}) {
    l.logger.Sugar().Errorf(format, args...)
}

// Fatalf 格式化致命错误日志
func (l *ZapLogger) Fatalf(format string, args ...interface{}) {
    l.logger.Sugar().Fatalf(format, args...)
}

// Panicf 格式化恐慌日志
func (l *ZapLogger) Panicf(format string, args ...interface{}) {
    l.logger.Sugar().Panicf(format, args...)
}

// WithContext 添加上下文
func (l *ZapLogger) WithContext(ctx context.Context) Logger {
    // 从上下文中提取字段
    fields := extractFieldsFromContext(ctx)
    return &ZapLogger{
        logger: l.logger.With(l.convertFields(fields)...),
        level:  l.level,
    }
}

// WithFields 添加字段
func (l *ZapLogger) WithFields(fields ...Field) Logger {
    return &ZapLogger{
        logger: l.logger.With(l.convertFields(fields)...),
        level:  l.level,
    }
}

// WithError 添加错误
func (l *ZapLogger) WithError(err error) Logger {
    return l.WithFields(Error(err))
}

// SetLevel 设置日志级别
func (l *ZapLogger) SetLevel(level Level) {
    l.level = level
}

// GetLevel 获取日志级别
func (l *ZapLogger) GetLevel() Level {
    return l.level
}

// convertFields 转换字段
func (l *ZapLogger) convertFields(fields []Field) []zap.Field {
    zapFields := make([]zap.Field, len(fields))
    for i, field := range fields {
        zapFields[i] = zap.Any(field.Key, field.Value)
    }
    return zapFields
}

// parseLevel 解析日志级别
func parseLevel(level string) zapcore.Level {
    switch level {
    case "debug":
        return zapcore.DebugLevel
    case "info":
        return zapcore.InfoLevel
    case "warn":
        return zapcore.WarnLevel
    case "error":
        return zapcore.ErrorLevel
    case "fatal":
        return zapcore.FatalLevel
    case "panic":
        return zapcore.PanicLevel
    default:
        return zapcore.InfoLevel
    }
}

// parseLogLevel 解析自定义日志级别
func parseLogLevel(level string) Level {
    switch level {
    case "debug":
        return DebugLevel
    case "info":
        return InfoLevel
    case "warn":
        return WarnLevel
    case "error":
        return ErrorLevel
    case "fatal":
        return FatalLevel
    case "panic":
        return PanicLevel
    default:
        return InfoLevel
    }
}

// extractFieldsFromContext 从上下文提取字段
func extractFieldsFromContext(ctx context.Context) []Field {
    var fields []Field
    
    // 提取请求ID
    if requestID := ctx.Value("request_id"); requestID != nil {
        fields = append(fields, String("request_id", requestID.(string)))
    }
    
    // 提取用户ID
    if userID := ctx.Value("user_id"); userID != nil {
        fields = append(fields, Any("user_id", userID))
    }
    
    // 提取跟踪ID
    if traceID := ctx.Value("trace_id"); traceID != nil {
        fields = append(fields, String("trace_id", traceID.(string)))
    }
    
    return fields
}
```

### 5. **日志配置结构**

#### 5.1 日志配置定义
```go
// pkg/logger/config.go
package logger

// LogConfig 日志配置
type LogConfig struct {
    Level      string `mapstructure:"level" yaml:"level"`           // 日志级别
    Format     string `mapstructure:"format" yaml:"format"`         // 日志格式 (json/text)
    Output     string `mapstructure:"output" yaml:"output"`         // 输出方式 (stdout/stderr/file)
    FilePath   string `mapstructure:"file_path" yaml:"file_path"`   // 文件路径
    MaxSize    int    `mapstructure:"max_size" yaml:"max_size"`     // 最大文件大小(MB)
    MaxBackups int    `mapstructure:"max_backups" yaml:"max_backups"` // 最大备份数
    MaxAge     int    `mapstructure:"max_age" yaml:"max_age"`       // 最大保存天数
    Compress   bool   `mapstructure:"compress" yaml:"compress"`     // 是否压缩
}

// DefaultConfig 默认配置
func DefaultConfig() LogConfig {
    return LogConfig{
        Level:      "info",
        Format:     "json",
        Output:     "stdout",
        FilePath:   "logs/app.log",
        MaxSize:    100,
        MaxBackups: 10,
        MaxAge:     30,
        Compress:   true,
    }
}
```

### 6. **全局日志器管理**

#### 6.1 全局日志器
```go
// pkg/logger/global.go
package logger

import (
    "sync"
)

var (
    globalLogger Logger
    once         sync.Once
)

// InitGlobalLogger 初始化全局日志器
func InitGlobalLogger(config LogConfig) error {
    var err error
    once.Do(func() {
        globalLogger, err = NewZapLogger(config)
    })
    return err
}

// GetGlobalLogger 获取全局日志器
func GetGlobalLogger() Logger {
    if globalLogger == nil {
        // 使用默认配置初始化
        globalLogger, _ = NewZapLogger(DefaultConfig())
    }
    return globalLogger
}

// 全局便捷方法
func Debug(msg string, fields ...Field) {
    GetGlobalLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
    GetGlobalLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
    GetGlobalLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
    GetGlobalLogger().Error(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
    GetGlobalLogger().Fatal(msg, fields...)
}

func Panic(msg string, fields ...Field) {
    GetGlobalLogger().Panic(msg, fields...)
}

func Debugf(format string, args ...interface{}) {
    GetGlobalLogger().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
    GetGlobalLogger().Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
    GetGlobalLogger().Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
    GetGlobalLogger().Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
    GetGlobalLogger().Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
    GetGlobalLogger().Panicf(format, args...)
}
```

### 7. **中间件集成**

#### 7.1 Gin日志中间件
```go
// internal/middleware/logger.go
package middleware

import (
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/yourname/movieinfo/pkg/logger"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        // 记录访问日志
        logger.Info("HTTP Request",
            logger.String("method", param.Method),
            logger.String("path", param.Path),
            logger.String("protocol", param.Request.Proto),
            logger.Int("status_code", param.StatusCode),
            logger.String("client_ip", param.ClientIP),
            logger.String("user_agent", param.Request.UserAgent()),
            logger.Duration("latency", param.Latency),
            logger.Int("body_size", param.BodySize),
        )
        return ""
    })
}

// RequestIDMiddleware 请求ID中间件
func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := generateRequestID()
        c.Set("request_id", requestID)
        c.Header("X-Request-ID", requestID)
        
        // 添加到日志上下文
        ctx := context.WithValue(c.Request.Context(), "request_id", requestID)
        c.Request = c.Request.WithContext(ctx)
        
        c.Next()
    }
}

// generateRequestID 生成请求ID
func generateRequestID() string {
    return fmt.Sprintf("%d-%s", time.Now().UnixNano(), randomString(8))
}
```

### 8. **业务日志实践**

#### 8.1 业务日志示例
```go
// internal/services/user/service.go
package user

import (
    "context"
    
    "github.com/yourname/movieinfo/pkg/logger"
)

// UserService 用户服务
type UserService struct {
    logger logger.Logger
}

// NewUserService 创建用户服务
func NewUserService() *UserService {
    return &UserService{
        logger: logger.GetGlobalLogger(),
    }
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
    // 记录业务开始日志
    s.logger.WithContext(ctx).Info("User registration started",
        logger.String("email", req.Email),
        logger.String("username", req.Username),
    )
    
    // 业务逻辑处理
    user, err := s.createUser(ctx, req)
    if err != nil {
        // 记录错误日志
        s.logger.WithContext(ctx).Error("User registration failed",
            logger.String("email", req.Email),
            logger.Error(err),
        )
        return nil, err
    }
    
    // 记录成功日志
    s.logger.WithContext(ctx).Info("User registration completed",
        logger.Int64("user_id", user.ID),
        logger.String("email", user.Email),
    )
    
    return user, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
    // 记录登录尝试
    s.logger.WithContext(ctx).Info("User login attempt",
        logger.String("email", req.Email),
        logger.String("ip", getClientIP(ctx)),
    )
    
    user, err := s.authenticateUser(ctx, req)
    if err != nil {
        // 记录登录失败
        s.logger.WithContext(ctx).Warn("User login failed",
            logger.String("email", req.Email),
            logger.String("reason", err.Error()),
        )
        return nil, err
    }
    
    // 记录登录成功
    s.logger.WithContext(ctx).Info("User login successful",
        logger.Int64("user_id", user.ID),
        logger.String("email", user.Email),
    )
    
    return &LoginResponse{
        Token: generateToken(user),
        User:  user,
    }, nil
}
```

## 总结

日志系统搭建为MovieInfo项目提供了完整的日志解决方案。通过统一的日志接口、结构化日志格式和灵活的配置管理，我们建立了一个生产级的日志系统。

**关键设计要点**：
1. **统一接口**：标准化的日志接口和使用方式
2. **结构化日志**：便于分析和处理的结构化格式
3. **多级别支持**：灵活的日志级别控制
4. **上下文传递**：请求上下文的自动传递
5. **性能优化**：高性能的日志输出和轮转

**日志优势**：
- **问题诊断**：快速定位和分析问题
- **性能监控**：实时监控系统性能
- **安全审计**：完整的操作记录和审计
- **业务分析**：支持业务数据分析

**下一步**：基于这个日志系统，我们将实现错误处理机制，包括统一的错误定义、错误码管理和错误恢复策略。
