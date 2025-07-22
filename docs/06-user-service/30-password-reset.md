# 第30步：密码重置功能

## 📋 概述

密码重置功能是用户服务的重要安全功能，为忘记密码的用户提供安全的密码恢复机制。一个完善的密码重置系统需要平衡安全性和用户体验，防止恶意攻击的同时提供便捷的重置流程。

## 🎯 设计目标

### 1. **安全性**
- 防止暴力破解
- 重置链接时效性
- 一次性使用机制
- 身份验证保护

### 2. **用户体验**
- 简单的重置流程
- 清晰的操作指引
- 快速的邮件发送
- 友好的错误提示

### 3. **可靠性**
- 邮件发送保障
- 重试机制支持
- 状态跟踪管理
- 异常处理完善

## 🔄 密码重置流程

### 1. **重置流程图**

```
用户忘记密码 → 访问重置页面 → 输入邮箱 → 提交重置请求
        ↓              ↓           ↓           ↓
    点击忘记密码    显示重置表单   填写邮箱地址   发送HTTP请求
        ↓              ↓           ↓           ↓
验证邮箱存在 → 生成重置Token → 发送重置邮件 → 返回成功响应
        ↓              ↓             ↓           ↓
    查询用户数据    创建重置记录    邮件服务发送   提示查收邮件
        ↓              ↓             ↓           ↓
用户收到邮件 → 点击重置链接 → 访问重置页面 → 输入新密码
        ↓              ↓             ↓           ↓
    查看邮件内容    验证Token有效性  显示密码表单   提交新密码
        ↓              ↓             ↓           ↓
验证新密码 → 更新用户密码 → 清理重置记录 → 重置完成
        ↓          ↓             ↓           ↓
    密码强度检查   加密存储新密码   删除Token记录   跳转登录页面
```

### 2. **安全机制**

```
┌─────────────────────────────────────────────────────────────┐
│                    密码重置安全机制                          │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │  频率限制    │  │  Token验证   │  │  时效控制    │          │
│  │             │  │             │  │             │          │
│  │ • IP限制    │  │ • 唯一性    │  │ • 过期时间   │          │
│  │ • 邮箱限制   │  │ • 一次性    │  │ • 自动清理   │          │
│  │ • 时间窗口   │  │ • 签名验证   │  │ • 状态跟踪   │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
│                              │                              │
│                              ▼                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                  审计日志记录                            │ │
│  │                                                         │ │
│  │ • 重置请求记录  • 邮件发送记录  • Token使用记录          │ │
│  │ • 成功重置记录  • 失败尝试记录  • 异常行为记录          │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 🔧 核心组件实现

### 1. **密码重置服务**

```go
type PasswordResetService struct {
    userRepo       UserRepository
    resetRepo      PasswordResetRepository
    emailSender    EmailSender
    validator      *PasswordValidator
    rateLimiter    *ResetRateLimiter
    tokenGenerator TokenGenerator
    logger         *logrus.Logger
    metrics        *PasswordResetMetrics
}

func NewPasswordResetService(
    userRepo UserRepository,
    resetRepo PasswordResetRepository,
    emailSender EmailSender,
    validator *PasswordValidator,
    rateLimiter *ResetRateLimiter,
) *PasswordResetService {
    return &PasswordResetService{
        userRepo:       userRepo,
        resetRepo:      resetRepo,
        emailSender:    emailSender,
        validator:      validator,
        rateLimiter:    rateLimiter,
        tokenGenerator: NewSecureTokenGenerator(),
        logger:         logrus.New(),
        metrics:        NewPasswordResetMetrics(),
    }
}

// 请求密码重置
func (prs *PasswordResetService) RequestReset(ctx context.Context, email, clientIP string) error {
    start := time.Now()
    defer func() {
        prs.metrics.ObserveResetRequestDuration(time.Since(start))
    }()

    // 检查频率限制
    if blocked, remaining := prs.rateLimiter.IsBlocked(email, clientIP); blocked {
        prs.metrics.IncBlockedResetRequests()
        return fmt.Errorf("请求过于频繁，请%d分钟后重试", remaining/60)
    }

    // 验证邮箱格式
    if err := prs.validator.ValidateEmail(email); err != nil {
        prs.metrics.IncInvalidResetRequests("invalid_email")
        return err
    }

    // 查找用户
    user, err := prs.userRepo.FindByEmail(ctx, email)
    if err != nil {
        // 为了安全，即使用户不存在也返回成功
        prs.logger.Warnf("Password reset requested for non-existent email: %s", email)
        prs.metrics.IncInvalidResetRequests("user_not_found")
        return nil
    }

    // 检查用户状态
    if user.Status != UserStatusActive {
        prs.logger.Warnf("Password reset requested for inactive user: %s", email)
        prs.metrics.IncInvalidResetRequests("user_inactive")
        return nil
    }

    // 生成重置Token
    resetToken := prs.tokenGenerator.Generate()
    expiresAt := time.Now().Add(1 * time.Hour) // 1小时有效期

    // 创建重置记录
    resetRecord := &PasswordReset{
        ID:        uuid.New().String(),
        UserID:    user.ID,
        Email:     email,
        Token:     resetToken,
        ClientIP:  clientIP,
        ExpiresAt: expiresAt,
        CreatedAt: time.Now(),
        Used:      false,
    }

    // 清理旧的重置记录
    if err := prs.resetRepo.DeleteByUserID(ctx, user.ID); err != nil {
        prs.logger.Errorf("Failed to clean old reset records: %v", err)
    }

    // 保存重置记录
    if err := prs.resetRepo.Create(ctx, resetRecord); err != nil {
        prs.logger.Errorf("Failed to create reset record: %v", err)
        prs.metrics.IncResetRequestErrors("database_error")
        return errors.New("密码重置请求失败，请稍后重试")
    }

    // 发送重置邮件
    if err := prs.sendResetEmail(user, resetToken); err != nil {
        prs.logger.Errorf("Failed to send reset email: %v", err)
        prs.metrics.IncResetRequestErrors("email_error")
        // 邮件发送失败不影响重置记录的创建
    }

    // 记录频率限制
    prs.rateLimiter.RecordRequest(email, clientIP)

    prs.metrics.IncSuccessfulResetRequests()
    prs.logger.Infof("Password reset requested for user: %s", email)

    return nil
}

// 验证重置Token
func (prs *PasswordResetService) ValidateResetToken(ctx context.Context, token string) (*PasswordReset, error) {
    if token == "" {
        prs.metrics.IncInvalidTokenValidations("empty_token")
        return nil, errors.New("重置令牌不能为空")
    }

    // 查找重置记录
    resetRecord, err := prs.resetRepo.FindByToken(ctx, token)
    if err != nil {
        prs.metrics.IncInvalidTokenValidations("token_not_found")
        return nil, errors.New("无效的重置令牌")
    }

    // 检查是否已使用
    if resetRecord.Used {
        prs.metrics.IncInvalidTokenValidations("token_used")
        return nil, errors.New("重置令牌已被使用")
    }

    // 检查是否过期
    if time.Now().After(resetRecord.ExpiresAt) {
        prs.metrics.IncInvalidTokenValidations("token_expired")
        return nil, errors.New("重置令牌已过期")
    }

    prs.metrics.IncValidTokenValidations()
    return resetRecord, nil
}

// 重置密码
func (prs *PasswordResetService) ResetPassword(ctx context.Context, token, newPassword string) error {
    start := time.Now()
    defer func() {
        prs.metrics.ObservePasswordResetDuration(time.Since(start))
    }()

    // 验证Token
    resetRecord, err := prs.ValidateResetToken(ctx, token)
    if err != nil {
        return err
    }

    // 验证新密码
    if err := prs.validator.ValidatePassword(newPassword); err != nil {
        prs.metrics.IncPasswordResetErrors("invalid_password")
        return err
    }

    // 查找用户
    user, err := prs.userRepo.FindByID(ctx, resetRecord.UserID)
    if err != nil {
        prs.logger.Errorf("Failed to find user for reset: %v", err)
        prs.metrics.IncPasswordResetErrors("user_not_found")
        return errors.New("用户不存在")
    }

    // 检查新密码是否与当前密码相同
    if prs.validator.CheckPassword(newPassword, user.Password) {
        prs.metrics.IncPasswordResetErrors("same_password")
        return errors.New("新密码不能与当前密码相同")
    }

    // 加密新密码
    hashedPassword, err := prs.validator.HashPassword(newPassword)
    if err != nil {
        prs.logger.Errorf("Failed to hash new password: %v", err)
        prs.metrics.IncPasswordResetErrors("hash_error")
        return errors.New("密码处理失败")
    }

    // 开始事务
    tx, err := prs.userRepo.BeginTx(ctx)
    if err != nil {
        prs.logger.Errorf("Failed to begin transaction: %v", err)
        return errors.New("密码重置失败")
    }
    defer tx.Rollback()

    // 更新用户密码
    user.Password = hashedPassword
    user.PasswordChangedAt = time.Now()
    user.UpdatedAt = time.Now()

    if err := prs.userRepo.UpdateWithTx(ctx, tx, user); err != nil {
        prs.logger.Errorf("Failed to update user password: %v", err)
        prs.metrics.IncPasswordResetErrors("update_error")
        return errors.New("密码更新失败")
    }

    // 标记重置记录为已使用
    resetRecord.Used = true
    resetRecord.UsedAt = time.Now()

    if err := prs.resetRepo.UpdateWithTx(ctx, tx, resetRecord); err != nil {
        prs.logger.Errorf("Failed to update reset record: %v", err)
        prs.metrics.IncPasswordResetErrors("record_update_error")
        return errors.New("重置记录更新失败")
    }

    // 提交事务
    if err := tx.Commit(); err != nil {
        prs.logger.Errorf("Failed to commit transaction: %v", err)
        return errors.New("密码重置失败")
    }

    // 发送密码重置成功通知邮件
    go func() {
        if err := prs.sendResetSuccessEmail(user); err != nil {
            prs.logger.Errorf("Failed to send reset success email: %v", err)
        }
    }()

    // 撤销用户的所有活跃会话（可选）
    go func() {
        if err := prs.revokeUserSessions(ctx, user.ID); err != nil {
            prs.logger.Errorf("Failed to revoke user sessions: %v", err)
        }
    }()

    prs.metrics.IncSuccessfulPasswordResets()
    prs.logger.Infof("Password reset successfully for user: %s", user.Email)

    return nil
}

func (prs *PasswordResetService) sendResetEmail(user *User, token string) error {
    resetURL := fmt.Sprintf("https://movieinfo.com/reset-password?token=%s", token)
    
    emailData := PasswordResetEmailData{
        Username:  user.Username,
        Email:     user.Email,
        ResetURL:  resetURL,
        ExpiresIn: "1小时",
        Timestamp: time.Now(),
    }

    return prs.emailSender.SendPasswordResetEmail(user.Email, emailData)
}

func (prs *PasswordResetService) sendResetSuccessEmail(user *User) error {
    emailData := PasswordResetSuccessEmailData{
        Username:  user.Username,
        Email:     user.Email,
        Timestamp: time.Now(),
        ClientIP:  "", // 可以从上下文获取
    }

    return prs.emailSender.SendPasswordResetSuccessEmail(user.Email, emailData)
}
```

### 2. **频率限制器**

```go
type ResetRateLimiter struct {
    redis         *redis.Client
    emailLimit    int
    ipLimit       int
    timeWindow    time.Duration
    blockDuration time.Duration
    logger        *logrus.Logger
}

func NewResetRateLimiter(redis *redis.Client) *ResetRateLimiter {
    return &ResetRateLimiter{
        redis:         redis,
        emailLimit:    3,  // 每个邮箱每小时最多3次
        ipLimit:       10, // 每个IP每小时最多10次
        timeWindow:    1 * time.Hour,
        blockDuration: 1 * time.Hour,
        logger:        logrus.New(),
    }
}

func (rrl *ResetRateLimiter) RecordRequest(email, clientIP string) {
    ctx := context.Background()
    now := time.Now()

    // 记录邮箱请求
    emailKey := fmt.Sprintf("reset_rate:email:%s", email)
    rrl.redis.ZAdd(ctx, emailKey, &redis.Z{
        Score:  float64(now.Unix()),
        Member: now.UnixNano(),
    })
    rrl.redis.Expire(ctx, emailKey, rrl.timeWindow)

    // 记录IP请求
    ipKey := fmt.Sprintf("reset_rate:ip:%s", clientIP)
    rrl.redis.ZAdd(ctx, ipKey, &redis.Z{
        Score:  float64(now.Unix()),
        Member: now.UnixNano(),
    })
    rrl.redis.Expire(ctx, ipKey, rrl.timeWindow)
}

func (rrl *ResetRateLimiter) IsBlocked(email, clientIP string) (bool, int64) {
    ctx := context.Background()
    now := time.Now()
    windowStart := now.Add(-rrl.timeWindow)

    // 检查邮箱频率
    emailKey := fmt.Sprintf("reset_rate:email:%s", email)
    emailCount, _ := rrl.redis.ZCount(ctx, emailKey, 
        fmt.Sprintf("%d", windowStart.Unix()), 
        fmt.Sprintf("%d", now.Unix())).Result()

    // 检查IP频率
    ipKey := fmt.Sprintf("reset_rate:ip:%s", clientIP)
    ipCount, _ := rrl.redis.ZCount(ctx, ipKey,
        fmt.Sprintf("%d", windowStart.Unix()),
        fmt.Sprintf("%d", now.Unix())).Result()

    if emailCount >= int64(rrl.emailLimit) || ipCount >= int64(rrl.ipLimit) {
        // 计算剩余阻塞时间
        var oldestTime int64
        if emailCount >= int64(rrl.emailLimit) {
            oldest, _ := rrl.redis.ZRange(ctx, emailKey, 0, 0).Result()
            if len(oldest) > 0 {
                oldestTime, _ = strconv.ParseInt(oldest[0], 10, 64)
            }
        } else {
            oldest, _ := rrl.redis.ZRange(ctx, ipKey, 0, 0).Result()
            if len(oldest) > 0 {
                oldestTime, _ = strconv.ParseInt(oldest[0], 10, 64)
            }
        }

        if oldestTime > 0 {
            remaining := rrl.timeWindow.Seconds() - float64(now.Unix()-oldestTime)
            if remaining > 0 {
                return true, int64(remaining)
            }
        }
    }

    return false, 0
}
```

### 3. **安全Token生成器**

```go
type TokenGenerator interface {
    Generate() string
}

type SecureTokenGenerator struct {
    length int
}

func NewSecureTokenGenerator() *SecureTokenGenerator {
    return &SecureTokenGenerator{
        length: 32,
    }
}

func (stg *SecureTokenGenerator) Generate() string {
    // 生成随机字节
    bytes := make([]byte, stg.length)
    if _, err := rand.Read(bytes); err != nil {
        // 如果随机数生成失败，使用时间戳作为后备
        return fmt.Sprintf("%d-%s", time.Now().UnixNano(), uuid.New().String())
    }

    // 转换为URL安全的Base64编码
    return base64.URLEncoding.EncodeToString(bytes)
}

// 验证Token格式
func (stg *SecureTokenGenerator) ValidateFormat(token string) bool {
    if len(token) == 0 {
        return false
    }

    // 检查是否为有效的Base64编码
    if _, err := base64.URLEncoding.DecodeString(token); err != nil {
        return false
    }

    return true
}
```

### 4. **密码重置数据模型**

```go
type PasswordReset struct {
    ID        string    `gorm:"primaryKey" json:"id"`
    UserID    string    `gorm:"not null;index" json:"user_id"`
    Email     string    `gorm:"not null;index" json:"email"`
    Token     string    `gorm:"not null;uniqueIndex" json:"token"`
    ClientIP  string    `gorm:"not null" json:"client_ip"`
    ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
    CreatedAt time.Time `gorm:"not null" json:"created_at"`
    Used      bool      `gorm:"not null;default:false" json:"used"`
    UsedAt    time.Time `json:"used_at"`
}

type PasswordResetRepository interface {
    Create(ctx context.Context, reset *PasswordReset) error
    FindByToken(ctx context.Context, token string) (*PasswordReset, error)
    FindByUserID(ctx context.Context, userID string) ([]*PasswordReset, error)
    Update(ctx context.Context, reset *PasswordReset) error
    UpdateWithTx(ctx context.Context, tx *gorm.DB, reset *PasswordReset) error
    DeleteByUserID(ctx context.Context, userID string) error
    DeleteExpired(ctx context.Context) error
}

type passwordResetRepository struct {
    db     *gorm.DB
    logger *logrus.Logger
}

func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
    return &passwordResetRepository{
        db:     db,
        logger: logrus.New(),
    }
}

func (prr *passwordResetRepository) Create(ctx context.Context, reset *PasswordReset) error {
    return prr.db.WithContext(ctx).Create(reset).Error
}

func (prr *passwordResetRepository) FindByToken(ctx context.Context, token string) (*PasswordReset, error) {
    var reset PasswordReset
    err := prr.db.WithContext(ctx).Where("token = ?", token).First(&reset).Error
    if err != nil {
        return nil, err
    }
    return &reset, nil
}

func (prr *passwordResetRepository) DeleteExpired(ctx context.Context) error {
    return prr.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&PasswordReset{}).Error
}
```

## 📊 监控指标

### 1. **密码重置指标**

```go
type PasswordResetMetrics struct {
    resetRequests        *prometheus.CounterVec
    resetCompletions     prometheus.Counter
    resetRequestDuration prometheus.Histogram
    resetDuration        prometheus.Histogram
    tokenValidations     *prometheus.CounterVec
    blockedRequests      prometheus.Counter
    emailSendErrors      prometheus.Counter
}

func NewPasswordResetMetrics() *PasswordResetMetrics {
    return &PasswordResetMetrics{
        resetRequests: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "password_reset_requests_total",
                Help: "Total number of password reset requests",
            },
            []string{"status"},
        ),
        resetCompletions: prometheus.NewCounter(
            prometheus.CounterOpts{
                Name: "password_reset_completions_total",
                Help: "Total number of completed password resets",
            },
        ),
        resetRequestDuration: prometheus.NewHistogram(
            prometheus.HistogramOpts{
                Name: "password_reset_request_duration_seconds",
                Help: "Duration of password reset request processing",
            },
        ),
        resetDuration: prometheus.NewHistogram(
            prometheus.HistogramOpts{
                Name: "password_reset_duration_seconds",
                Help: "Duration of password reset completion",
            },
        ),
        tokenValidations: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "password_reset_token_validations_total",
                Help: "Total number of password reset token validations",
            },
            []string{"status"},
        ),
        blockedRequests: prometheus.NewCounter(
            prometheus.CounterOpts{
                Name: "password_reset_blocked_requests_total",
                Help: "Total number of blocked password reset requests",
            },
        ),
        emailSendErrors: prometheus.NewCounter(
            prometheus.CounterOpts{
                Name: "password_reset_email_errors_total",
                Help: "Total number of password reset email send errors",
            },
        ),
    }
}
```

## 🔧 HTTP处理器

### 1. **密码重置API端点**

```go
func (uc *UserController) RequestPasswordReset(c *gin.Context) {
    var req struct {
        Email string `json:"email" binding:"required,email"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": "邮箱格式不正确",
        })
        return
    }

    clientIP := c.ClientIP()
    if err := uc.passwordResetService.RequestReset(c.Request.Context(), req.Email, clientIP); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": err.Error(),
        })
        return
    }

    c.JSON(200, gin.H{
        "success": true,
        "message": "如果该邮箱已注册，您将收到密码重置邮件",
    })
}

func (uc *UserController) ValidateResetToken(c *gin.Context) {
    token := c.Query("token")
    if token == "" {
        c.JSON(400, gin.H{
            "success": false,
            "message": "缺少重置令牌",
        })
        return
    }

    _, err := uc.passwordResetService.ValidateResetToken(c.Request.Context(), token)
    if err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": err.Error(),
        })
        return
    }

    c.JSON(200, gin.H{
        "success": true,
        "message": "重置令牌有效",
    })
}

func (uc *UserController) ResetPassword(c *gin.Context) {
    var req struct {
        Token       string `json:"token" binding:"required"`
        NewPassword string `json:"new_password" binding:"required,min=8"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": "请求参数错误",
        })
        return
    }

    if err := uc.passwordResetService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": err.Error(),
        })
        return
    }

    c.JSON(200, gin.H{
        "success": true,
        "message": "密码重置成功，请使用新密码登录",
    })
}
```

## 📝 总结

密码重置功能为MovieInfo项目提供了安全可靠的密码恢复机制：

**核心功能**：
1. **安全重置**：Token验证、时效控制、一次性使用
2. **频率限制**：防止暴力破解和恶意请求
3. **邮件通知**：重置链接发送和成功通知
4. **审计日志**：完整的操作记录和监控

**安全措施**：
- 重置Token时效性控制
- 频率限制和IP封禁
- 一次性使用机制
- 会话撤销保护

**用户体验**：
- 简单的重置流程
- 清晰的操作指引
- 友好的错误提示
- 快速的邮件响应

至此，用户服务的核心功能已经完成。下一步，我们将继续完成电影服务、评论服务等其他模块的开发文档。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第31步：电影数据模型](../07-movie-service/31-movie-model.md)
