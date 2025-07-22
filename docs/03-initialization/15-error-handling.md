# 3.5 错误处理机制

## 概述

错误处理是应用程序健壮性的重要保障，它决定了系统在异常情况下的表现和用户体验。对于MovieInfo项目，我们需要建立一个统一、清晰、用户友好的错误处理机制，包括错误定义、错误码管理、错误恢复和错误监控。

## 为什么需要统一的错误处理？

### 1. **用户体验**
- **友好提示**：为用户提供清晰、有用的错误信息
- **一致性**：统一的错误格式和处理方式
- **国际化**：支持多语言的错误消息
- **引导性**：提供解决问题的建议和指导

### 2. **开发效率**
- **标准化**：统一的错误处理模式和规范
- **可维护性**：集中的错误定义和管理
- **调试便利**：详细的错误信息和堆栈跟踪
- **代码复用**：通用的错误处理逻辑

### 3. **系统稳定性**
- **优雅降级**：在错误发生时保持系统稳定
- **错误隔离**：防止错误传播和级联失败
- **自动恢复**：某些错误的自动恢复机制
- **监控告警**：错误的实时监控和告警

### 4. **运维支持**
- **问题定位**：快速定位错误原因和位置
- **趋势分析**：错误趋势的分析和预警
- **性能影响**：了解错误对系统性能的影响
- **容量规划**：基于错误数据进行容量规划

## 错误处理架构设计

### 1. **错误分类体系**

```
错误分类
├── 系统错误 (System Errors)
│   ├── 数据库错误 (Database)
│   ├── 网络错误 (Network)
│   ├── 文件系统错误 (FileSystem)
│   └── 第三方服务错误 (External)
├── 业务错误 (Business Errors)
│   ├── 验证错误 (Validation)
│   ├── 权限错误 (Permission)
│   ├── 资源错误 (Resource)
│   └── 状态错误 (State)
├── 客户端错误 (Client Errors)
│   ├── 请求格式错误 (Format)
│   ├── 参数错误 (Parameter)
│   ├── 认证错误 (Authentication)
│   └── 授权错误 (Authorization)
└── 未知错误 (Unknown Errors)
    ├── 恐慌错误 (Panic)
    ├── 运行时错误 (Runtime)
    └── 其他错误 (Others)
```

### 2. **错误码设计**

#### 2.1 错误码结构
```
错误码格式: XXYYZZ
├── XX: 服务代码 (01-99)
│   ├── 01: 通用服务
│   ├── 10: 用户服务
│   ├── 20: 电影服务
│   └── 30: 评论服务
├── YY: 错误类型 (01-99)
│   ├── 01: 系统错误
│   ├── 02: 业务错误
│   ├── 03: 客户端错误
│   └── 04: 第三方错误
└── ZZ: 具体错误 (01-99)
    └── 按错误类型递增编号
```

#### 2.2 错误码定义
```go
// pkg/errors/codes.go
package errors

// ErrorCode 错误码类型
type ErrorCode int

const (
    // 通用错误 (01xxxx)
    CodeSuccess           ErrorCode = 0      // 成功
    CodeInternalError     ErrorCode = 10101  // 内部服务器错误
    CodeDatabaseError     ErrorCode = 10102  // 数据库错误
    CodeCacheError        ErrorCode = 10103  // 缓存错误
    CodeNetworkError      ErrorCode = 10104  // 网络错误
    CodeConfigError       ErrorCode = 10105  // 配置错误
    
    // 客户端错误 (01xxxx)
    CodeBadRequest        ErrorCode = 10301  // 请求格式错误
    CodeUnauthorized      ErrorCode = 10302  // 未认证
    CodeForbidden         ErrorCode = 10303  // 无权限
    CodeNotFound          ErrorCode = 10304  // 资源不存在
    CodeMethodNotAllowed  ErrorCode = 10305  // 方法不允许
    CodeTooManyRequests   ErrorCode = 10306  // 请求过多
    CodeValidationFailed  ErrorCode = 10307  // 参数验证失败
    
    // 用户服务错误 (10xxxx)
    CodeUserNotFound      ErrorCode = 100201 // 用户不存在
    CodeUserExists        ErrorCode = 100202 // 用户已存在
    CodeInvalidPassword   ErrorCode = 100203 // 密码错误
    CodeUserDisabled      ErrorCode = 100204 // 用户已禁用
    CodeEmailNotVerified  ErrorCode = 100205 // 邮箱未验证
    CodeTokenExpired      ErrorCode = 100206 // Token已过期
    CodeTokenInvalid      ErrorCode = 100207 // Token无效
    
    // 电影服务错误 (20xxxx)
    CodeMovieNotFound     ErrorCode = 200201 // 电影不存在
    CodeMovieExists       ErrorCode = 200202 // 电影已存在
    CodeCategoryNotFound  ErrorCode = 200203 // 分类不存在
    CodeInvalidRating     ErrorCode = 200204 // 无效评分
    
    // 评论服务错误 (30xxxx)
    CodeCommentNotFound   ErrorCode = 300201 // 评论不存在
    CodeCommentExists     ErrorCode = 300202 // 评论已存在
    CodeCommentTooLong    ErrorCode = 300203 // 评论过长
    CodeCommentBlocked    ErrorCode = 300204 // 评论被屏蔽
    CodeDuplicateRating   ErrorCode = 300205 // 重复评分
)

// String 返回错误码字符串
func (c ErrorCode) String() string {
    return fmt.Sprintf("%06d", int(c))
}

// HTTPStatus 返回对应的HTTP状态码
func (c ErrorCode) HTTPStatus() int {
    switch {
    case c == CodeSuccess:
        return 200
    case c >= 10301 && c <= 10399:
        return 400 // Bad Request
    case c == CodeUnauthorized || c == CodeTokenExpired || c == CodeTokenInvalid:
        return 401 // Unauthorized
    case c == CodeForbidden:
        return 403 // Forbidden
    case c == CodeNotFound || c >= 100201 && c <= 100299:
        return 404 // Not Found
    case c == CodeMethodNotAllowed:
        return 405 // Method Not Allowed
    case c == CodeTooManyRequests:
        return 429 // Too Many Requests
    case c >= 10101 && c <= 10199:
        return 500 // Internal Server Error
    default:
        return 500
    }
}
```

### 3. **错误结构定义**

#### 3.1 基础错误结构
```go
// pkg/errors/error.go
package errors

import (
    "fmt"
    "runtime"
    "time"
)

// Error 自定义错误结构
type Error struct {
    Code      ErrorCode              `json:"code"`                // 错误码
    Message   string                 `json:"message"`             // 错误消息
    Details   string                 `json:"details,omitempty"`   // 详细信息
    Timestamp time.Time              `json:"timestamp"`           // 时间戳
    RequestID string                 `json:"request_id,omitempty"` // 请求ID
    Stack     string                 `json:"stack,omitempty"`     // 堆栈信息
    Metadata  map[string]interface{} `json:"metadata,omitempty"`  // 元数据
    Cause     error                  `json:"cause,omitempty"`     // 原始错误
}

// Error 实现error接口
func (e *Error) Error() string {
    if e.Details != "" {
        return fmt.Sprintf("[%s] %s: %s", e.Code.String(), e.Message, e.Details)
    }
    return fmt.Sprintf("[%s] %s", e.Code.String(), e.Message)
}

// Unwrap 返回原始错误
func (e *Error) Unwrap() error {
    return e.Cause
}

// WithDetails 添加详细信息
func (e *Error) WithDetails(details string) *Error {
    e.Details = details
    return e
}

// WithRequestID 添加请求ID
func (e *Error) WithRequestID(requestID string) *Error {
    e.RequestID = requestID
    return e
}

// WithMetadata 添加元数据
func (e *Error) WithMetadata(key string, value interface{}) *Error {
    if e.Metadata == nil {
        e.Metadata = make(map[string]interface{})
    }
    e.Metadata[key] = value
    return e
}

// WithCause 添加原始错误
func (e *Error) WithCause(cause error) *Error {
    e.Cause = cause
    return e
}

// WithStack 添加堆栈信息
func (e *Error) WithStack() *Error {
    buf := make([]byte, 2048)
    n := runtime.Stack(buf, false)
    e.Stack = string(buf[:n])
    return e
}
```

#### 3.2 错误创建函数
```go
// pkg/errors/builder.go
package errors

import (
    "context"
    "time"
)

// New 创建新错误
func New(code ErrorCode, message string) *Error {
    return &Error{
        Code:      code,
        Message:   message,
        Timestamp: time.Now(),
    }
}

// Newf 创建格式化错误
func Newf(code ErrorCode, format string, args ...interface{}) *Error {
    return New(code, fmt.Sprintf(format, args...))
}

// Wrap 包装错误
func Wrap(err error, code ErrorCode, message string) *Error {
    if err == nil {
        return nil
    }
    
    return &Error{
        Code:      code,
        Message:   message,
        Timestamp: time.Now(),
        Cause:     err,
    }
}

// Wrapf 包装格式化错误
func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *Error {
    return Wrap(err, code, fmt.Sprintf(format, args...))
}

// FromContext 从上下文创建错误
func FromContext(ctx context.Context, code ErrorCode, message string) *Error {
    err := New(code, message)
    
    // 从上下文提取请求ID
    if requestID := ctx.Value("request_id"); requestID != nil {
        err.RequestID = requestID.(string)
    }
    
    return err
}

// 预定义错误创建函数
func BadRequest(message string) *Error {
    return New(CodeBadRequest, message)
}

func Unauthorized(message string) *Error {
    return New(CodeUnauthorized, message)
}

func Forbidden(message string) *Error {
    return New(CodeForbidden, message)
}

func NotFound(message string) *Error {
    return New(CodeNotFound, message)
}

func InternalError(message string) *Error {
    return New(CodeInternalError, message)
}

func ValidationFailed(message string) *Error {
    return New(CodeValidationFailed, message)
}

// 业务错误创建函数
func UserNotFound() *Error {
    return New(CodeUserNotFound, "用户不存在")
}

func UserExists() *Error {
    return New(CodeUserExists, "用户已存在")
}

func InvalidPassword() *Error {
    return New(CodeInvalidPassword, "密码错误")
}

func MovieNotFound() *Error {
    return New(CodeMovieNotFound, "电影不存在")
}

func CommentNotFound() *Error {
    return New(CodeCommentNotFound, "评论不存在")
}
```

### 4. **错误处理中间件**

#### 4.1 Gin错误处理中间件
```go
// internal/middleware/error.go
package middleware

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/yourname/movieinfo/pkg/errors"
    "github.com/yourname/movieinfo/pkg/logger"
)

// ErrorHandlerMiddleware 错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // 处理panic
                logger.Error("Panic recovered",
                    logger.Any("error", err),
                    logger.String("path", c.Request.URL.Path),
                    logger.String("method", c.Request.Method),
                )
                
                appErr := errors.InternalError("服务器内部错误").WithStack()
                handleError(c, appErr)
                c.Abort()
            }
        }()
        
        c.Next()
        
        // 处理错误
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            handleError(c, err)
        }
    }
}

// handleError 处理错误
func handleError(c *gin.Context, err error) {
    var appErr *errors.Error
    
    // 类型断言
    if e, ok := err.(*errors.Error); ok {
        appErr = e
    } else {
        // 包装普通错误
        appErr = errors.Wrap(err, errors.CodeInternalError, "服务器内部错误")
    }
    
    // 添加请求ID
    if requestID := c.GetString("request_id"); requestID != "" {
        appErr.WithRequestID(requestID)
    }
    
    // 记录错误日志
    logError(c, appErr)
    
    // 返回错误响应
    c.JSON(appErr.Code.HTTPStatus(), gin.H{
        "code":       appErr.Code,
        "message":    appErr.Message,
        "details":    appErr.Details,
        "timestamp":  appErr.Timestamp,
        "request_id": appErr.RequestID,
    })
}

// logError 记录错误日志
func logError(c *gin.Context, err *errors.Error) {
    fields := []logger.Field{
        logger.String("error_code", err.Code.String()),
        logger.String("error_message", err.Message),
        logger.String("path", c.Request.URL.Path),
        logger.String("method", c.Request.Method),
        logger.String("client_ip", c.ClientIP()),
        logger.String("user_agent", c.Request.UserAgent()),
    }
    
    if err.RequestID != "" {
        fields = append(fields, logger.String("request_id", err.RequestID))
    }
    
    if err.Details != "" {
        fields = append(fields, logger.String("details", err.Details))
    }
    
    if err.Cause != nil {
        fields = append(fields, logger.String("cause", err.Cause.Error()))
    }
    
    // 根据错误类型选择日志级别
    switch {
    case err.Code >= 10101 && err.Code <= 10199: // 系统错误
        logger.Error("System error occurred", fields...)
    case err.Code >= 10301 && err.Code <= 10399: // 客户端错误
        logger.Warn("Client error occurred", fields...)
    default:
        logger.Info("Business error occurred", fields...)
    }
}
```

### 5. **业务层错误处理**

#### 5.1 服务层错误处理示例
```go
// internal/services/user/service.go
package user

import (
    "context"
    "database/sql"
    
    "github.com/yourname/movieinfo/pkg/errors"
    "github.com/yourname/movieinfo/pkg/logger"
)

// UserService 用户服务
type UserService struct {
    repo   UserRepository
    logger logger.Logger
}

// GetUser 获取用户
func (s *UserService) GetUser(ctx context.Context, userID int64) (*User, error) {
    user, err := s.repo.GetByID(ctx, userID)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.FromContext(ctx, errors.CodeUserNotFound, "用户不存在").
                WithMetadata("user_id", userID)
        }
        return nil, errors.Wrap(err, errors.CodeDatabaseError, "查询用户失败").
            WithMetadata("user_id", userID)
    }
    
    return user, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // 验证请求参数
    if err := s.validateCreateUserRequest(req); err != nil {
        return nil, err
    }
    
    // 检查用户是否已存在
    existingUser, err := s.repo.GetByEmail(ctx, req.Email)
    if err != nil && err != sql.ErrNoRows {
        return nil, errors.Wrap(err, errors.CodeDatabaseError, "检查用户是否存在失败")
    }
    
    if existingUser != nil {
        return nil, errors.FromContext(ctx, errors.CodeUserExists, "用户已存在").
            WithMetadata("email", req.Email)
    }
    
    // 创建用户
    user := &User{
        Email:    req.Email,
        Username: req.Username,
        Password: hashPassword(req.Password),
    }
    
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, errors.Wrap(err, errors.CodeDatabaseError, "创建用户失败").
            WithMetadata("email", req.Email)
    }
    
    s.logger.WithContext(ctx).Info("User created successfully",
        logger.Int64("user_id", user.ID),
        logger.String("email", user.Email),
    )
    
    return user, nil
}

// validateCreateUserRequest 验证创建用户请求
func (s *UserService) validateCreateUserRequest(req *CreateUserRequest) error {
    if req.Email == "" {
        return errors.ValidationFailed("邮箱不能为空")
    }
    
    if !isValidEmail(req.Email) {
        return errors.ValidationFailed("邮箱格式不正确")
    }
    
    if req.Username == "" {
        return errors.ValidationFailed("用户名不能为空")
    }
    
    if len(req.Username) < 3 || len(req.Username) > 20 {
        return errors.ValidationFailed("用户名长度必须在3-20个字符之间")
    }
    
    if req.Password == "" {
        return errors.ValidationFailed("密码不能为空")
    }
    
    if len(req.Password) < 8 {
        return errors.ValidationFailed("密码长度不能少于8个字符")
    }
    
    return nil
}
```

### 6. **错误监控和告警**

#### 6.1 错误统计
```go
// pkg/errors/metrics.go
package errors

import (
    "sync"
    "time"
)

// ErrorMetrics 错误统计
type ErrorMetrics struct {
    mutex      sync.RWMutex
    counters   map[ErrorCode]int64
    lastErrors map[ErrorCode]time.Time
}

// NewErrorMetrics 创建错误统计
func NewErrorMetrics() *ErrorMetrics {
    return &ErrorMetrics{
        counters:   make(map[ErrorCode]int64),
        lastErrors: make(map[ErrorCode]time.Time),
    }
}

// Record 记录错误
func (m *ErrorMetrics) Record(code ErrorCode) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    m.counters[code]++
    m.lastErrors[code] = time.Now()
}

// GetCount 获取错误次数
func (m *ErrorMetrics) GetCount(code ErrorCode) int64 {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    return m.counters[code]
}

// GetAllCounts 获取所有错误统计
func (m *ErrorMetrics) GetAllCounts() map[ErrorCode]int64 {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    result := make(map[ErrorCode]int64)
    for code, count := range m.counters {
        result[code] = count
    }
    
    return result
}

// Reset 重置统计
func (m *ErrorMetrics) Reset() {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    m.counters = make(map[ErrorCode]int64)
    m.lastErrors = make(map[ErrorCode]time.Time)
}

// 全局错误统计实例
var globalMetrics = NewErrorMetrics()

// RecordError 记录错误到全局统计
func RecordError(code ErrorCode) {
    globalMetrics.Record(code)
}

// GetErrorMetrics 获取全局错误统计
func GetErrorMetrics() *ErrorMetrics {
    return globalMetrics
}
```

## 总结

错误处理机制为MovieInfo项目提供了完整的错误管理解决方案。通过统一的错误定义、标准化的错误码和完善的处理流程，我们建立了一个健壮、用户友好的错误处理系统。

**关键设计要点**：
1. **统一错误结构**：标准化的错误定义和格式
2. **分层错误码**：清晰的错误分类和编码规则
3. **上下文传递**：错误信息的上下文关联
4. **日志集成**：错误与日志系统的深度集成
5. **监控支持**：错误统计和监控机制

**错误处理优势**：
- **用户友好**：清晰、有用的错误提示
- **开发效率**：统一的错误处理模式
- **系统稳定**：优雅的错误处理和恢复
- **运维支持**：完善的错误监控和分析

**下一步**：基于完整的项目初始化基础，我们将开始数据层开发，包括数据库连接池、数据模型定义和CRUD操作实现。
