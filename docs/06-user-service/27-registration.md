# 第27步：注册功能实现

## 📋 概述

用户注册是MovieInfo项目的核心功能之一，它为新用户提供了加入平台的入口。一个完善的注册系统不仅要确保用户信息的安全性，还要提供良好的用户体验和完整的数据验证机制。

## 🎯 设计目标

### 1. **安全性**
- 密码安全存储
- 邮箱验证机制
- 防止重复注册
- 输入数据验证

### 2. **用户体验**
- 简洁的注册流程
- 清晰的错误提示
- 快速的响应时间
- 友好的界面设计

### 3. **数据完整性**
- 用户信息验证
- 数据格式检查
- 业务规则校验
- 事务一致性保证

### 4. **可扩展性**
- 支持多种注册方式
- 可配置的验证规则
- 灵活的用户属性
- 第三方登录预留

## 🏗️ 注册流程设计

### 1. **注册流程图**

```
用户访问注册页面 → 填写注册信息 → 前端验证 → 提交注册请求
        ↓                ↓            ↓           ↓
    显示注册表单      输入用户信息    客户端验证   发送HTTP请求
        ↓                ↓            ↓           ↓
服务端接收请求 → 数据验证 → 检查重复 → 创建用户 → 发送验证邮件
        ↓            ↓        ↓        ↓           ↓
    解析请求参数    格式验证   查询数据库  插入记录    邮件服务
        ↓            ↓        ↓        ↓           ↓
返回注册结果 ← 生成响应 ← 事务提交 ← 密码加密 ← 生成验证码
```

### 2. **数据流设计**

```
┌─────────────────┐    HTTP POST    ┌─────────────────┐
│   前端页面       │ ──────────────► │   主页服务       │
└─────────────────┘                 └─────────────────┘
                                             │
                                             │ gRPC
                                             ▼
                                    ┌─────────────────┐
                                    │   用户服务       │
                                    └─────────────────┘
                                             │
                                             ▼
                    ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
                    │   数据验证   │  │  业务逻辑    │  │  数据存储    │
                    └─────────────┘  └─────────────┘  └─────────────┘
                             │              │              │
                             ▼              ▼              ▼
                    ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
                    │  格式检查    │  │  重复检查    │  │   MySQL     │
                    │  长度验证    │  │  密码加密    │  │   Redis     │
                    │  规则校验    │  │  用户创建    │  │  邮件队列    │
                    └─────────────┘  └─────────────┘  └─────────────┘
```

## 🔧 核心组件实现

### 1. **注册请求模型**

```go
// 注册请求结构
type RegisterRequest struct {
    Email           string `json:"email" binding:"required,email"`
    Password        string `json:"password" binding:"required,min=8,max=128"`
    ConfirmPassword string `json:"confirm_password" binding:"required"`
    Username        string `json:"username" binding:"required,min=3,max=50"`
    FirstName       string `json:"first_name" binding:"max=50"`
    LastName        string `json:"last_name" binding:"max=50"`
    AcceptTerms     bool   `json:"accept_terms" binding:"required"`
}

// 注册响应结构
type RegisterResponse struct {
    Success     bool   `json:"success"`
    Message     string `json:"message"`
    UserID      string `json:"user_id,omitempty"`
    NeedVerify  bool   `json:"need_verify"`
    VerifyEmail string `json:"verify_email,omitempty"`
}

// 验证错误结构
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Code    string `json:"code"`
}

type RegisterErrorResponse struct {
    Success bool              `json:"success"`
    Message string            `json:"message"`
    Errors  []ValidationError `json:"errors,omitempty"`
}
```

### 2. **数据验证器**

```go
type UserValidator struct {
    db     *gorm.DB
    redis  *redis.Client
    logger *logrus.Logger
}

func NewUserValidator(db *gorm.DB, redis *redis.Client) *UserValidator {
    return &UserValidator{
        db:     db,
        redis:  redis,
        logger: logrus.New(),
    }
}

// 验证注册请求
func (uv *UserValidator) ValidateRegisterRequest(req *RegisterRequest) []ValidationError {
    var errors []ValidationError
    
    // 验证邮箱格式
    if err := uv.validateEmail(req.Email); err != nil {
        errors = append(errors, ValidationError{
            Field:   "email",
            Message: err.Error(),
            Code:    "INVALID_EMAIL",
        })
    }
    
    // 验证密码强度
    if err := uv.validatePassword(req.Password); err != nil {
        errors = append(errors, ValidationError{
            Field:   "password",
            Message: err.Error(),
            Code:    "WEAK_PASSWORD",
        })
    }
    
    // 验证密码确认
    if req.Password != req.ConfirmPassword {
        errors = append(errors, ValidationError{
            Field:   "confirm_password",
            Message: "密码确认不匹配",
            Code:    "PASSWORD_MISMATCH",
        })
    }
    
    // 验证用户名
    if err := uv.validateUsername(req.Username); err != nil {
        errors = append(errors, ValidationError{
            Field:   "username",
            Message: err.Error(),
            Code:    "INVALID_USERNAME",
        })
    }
    
    // 验证服务条款
    if !req.AcceptTerms {
        errors = append(errors, ValidationError{
            Field:   "accept_terms",
            Message: "必须接受服务条款",
            Code:    "TERMS_NOT_ACCEPTED",
        })
    }
    
    return errors
}

func (uv *UserValidator) validateEmail(email string) error {
    // 基本格式验证
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(email) {
        return errors.New("邮箱格式不正确")
    }
    
    // 检查邮箱是否已存在
    var count int64
    if err := uv.db.Model(&User{}).Where("email = ?", email).Count(&count).Error; err != nil {
        uv.logger.Errorf("Failed to check email existence: %v", err)
        return errors.New("邮箱验证失败")
    }
    
    if count > 0 {
        return errors.New("邮箱已被注册")
    }
    
    // 检查邮箱域名黑名单
    if uv.isEmailDomainBlocked(email) {
        return errors.New("不支持该邮箱域名")
    }
    
    return nil
}

func (uv *UserValidator) validatePassword(password string) error {
    // 长度检查
    if len(password) < 8 {
        return errors.New("密码长度至少8位")
    }
    
    if len(password) > 128 {
        return errors.New("密码长度不能超过128位")
    }
    
    // 复杂度检查
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
    
    complexity := 0
    if hasLower {
        complexity++
    }
    if hasUpper {
        complexity++
    }
    if hasNumber {
        complexity++
    }
    if hasSpecial {
        complexity++
    }
    
    if complexity < 3 {
        return errors.New("密码必须包含大写字母、小写字母、数字、特殊字符中的至少3种")
    }
    
    // 常见密码检查
    if uv.isCommonPassword(password) {
        return errors.New("密码过于简单，请使用更复杂的密码")
    }
    
    return nil
}

func (uv *UserValidator) validateUsername(username string) error {
    // 长度检查
    if len(username) < 3 {
        return errors.New("用户名长度至少3位")
    }
    
    if len(username) > 50 {
        return errors.New("用户名长度不能超过50位")
    }
    
    // 格式检查
    usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
    if !usernameRegex.MatchString(username) {
        return errors.New("用户名只能包含字母、数字、下划线和连字符")
    }
    
    // 检查用户名是否已存在
    var count int64
    if err := uv.db.Model(&User{}).Where("username = ?", username).Count(&count).Error; err != nil {
        uv.logger.Errorf("Failed to check username existence: %v", err)
        return errors.New("用户名验证失败")
    }
    
    if count > 0 {
        return errors.New("用户名已被使用")
    }
    
    // 检查保留用户名
    if uv.isReservedUsername(username) {
        return errors.New("该用户名为系统保留，请选择其他用户名")
    }
    
    return nil
}

func (uv *UserValidator) isEmailDomainBlocked(email string) bool {
    domain := strings.Split(email, "@")[1]
    blockedDomains := []string{
        "tempmail.com",
        "10minutemail.com",
        "guerrillamail.com",
    }
    
    for _, blocked := range blockedDomains {
        if domain == blocked {
            return true
        }
    }
    
    return false
}

func (uv *UserValidator) isCommonPassword(password string) bool {
    commonPasswords := []string{
        "password", "123456", "123456789", "qwerty",
        "abc123", "password123", "admin", "letmein",
    }
    
    lowerPassword := strings.ToLower(password)
    for _, common := range commonPasswords {
        if lowerPassword == common {
            return true
        }
    }
    
    return false
}

func (uv *UserValidator) isReservedUsername(username string) bool {
    reserved := []string{
        "admin", "administrator", "root", "system",
        "api", "www", "mail", "ftp", "support",
        "help", "info", "contact", "about",
    }
    
    lowerUsername := strings.ToLower(username)
    for _, res := range reserved {
        if lowerUsername == res {
            return true
        }
    }
    
    return false
}
```

### 3. **注册服务实现**

```go
type RegistrationService struct {
    userRepo    UserRepository
    validator   *UserValidator
    hasher      PasswordHasher
    emailSender EmailSender
    logger      *logrus.Logger
    metrics     *RegistrationMetrics
}

func NewRegistrationService(
    userRepo UserRepository,
    validator *UserValidator,
    hasher PasswordHasher,
    emailSender EmailSender,
) *RegistrationService {
    return &RegistrationService{
        userRepo:    userRepo,
        validator:   validator,
        hasher:      hasher,
        emailSender: emailSender,
        logger:      logrus.New(),
        metrics:     NewRegistrationMetrics(),
    }
}

func (rs *RegistrationService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
    start := time.Now()
    defer func() {
        rs.metrics.ObserveRegistrationDuration(time.Since(start))
    }()
    
    // 数据验证
    if errors := rs.validator.ValidateRegisterRequest(req); len(errors) > 0 {
        rs.metrics.IncRegistrationErrors("validation_failed")
        return &RegisterResponse{
            Success: false,
            Message: "注册信息验证失败",
        }, &ValidationErrors{Errors: errors}
    }
    
    // 密码加密
    hashedPassword, err := rs.hasher.HashPassword(req.Password)
    if err != nil {
        rs.logger.Errorf("Failed to hash password: %v", err)
        rs.metrics.IncRegistrationErrors("hash_failed")
        return nil, errors.New("密码处理失败")
    }
    
    // 创建用户对象
    user := &User{
        ID:        uuid.New().String(),
        Email:     req.Email,
        Username:  req.Username,
        Password:  hashedPassword,
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Status:    UserStatusPending,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // 生成邮箱验证码
    verificationCode := rs.generateVerificationCode()
    user.EmailVerificationCode = verificationCode
    user.EmailVerificationExpiry = time.Now().Add(24 * time.Hour)
    
    // 保存用户
    if err := rs.userRepo.Create(ctx, user); err != nil {
        rs.logger.Errorf("Failed to create user: %v", err)
        rs.metrics.IncRegistrationErrors("create_failed")
        
        if strings.Contains(err.Error(), "duplicate") {
            return &RegisterResponse{
                Success: false,
                Message: "邮箱或用户名已被注册",
            }, nil
        }
        
        return nil, errors.New("用户创建失败")
    }
    
    // 发送验证邮件
    if err := rs.sendVerificationEmail(user); err != nil {
        rs.logger.Errorf("Failed to send verification email: %v", err)
        // 邮件发送失败不影响注册成功
    }
    
    rs.metrics.IncSuccessfulRegistrations()
    rs.logger.Infof("User registered successfully: %s", user.Email)
    
    return &RegisterResponse{
        Success:     true,
        Message:     "注册成功，请查收验证邮件",
        UserID:      user.ID,
        NeedVerify:  true,
        VerifyEmail: user.Email,
    }, nil
}

func (rs *RegistrationService) generateVerificationCode() string {
    return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (rs *RegistrationService) sendVerificationEmail(user *User) error {
    emailData := EmailVerificationData{
        Username:         user.Username,
        Email:           user.Email,
        VerificationCode: user.EmailVerificationCode,
        ExpiryTime:      user.EmailVerificationExpiry,
        VerificationURL: fmt.Sprintf("https://movieinfo.com/verify?code=%s&email=%s", 
            user.EmailVerificationCode, user.Email),
    }
    
    return rs.emailSender.SendVerificationEmail(user.Email, emailData)
}
```

### 4. **邮箱验证实现**

```go
type EmailVerificationService struct {
    userRepo UserRepository
    logger   *logrus.Logger
}

func NewEmailVerificationService(userRepo UserRepository) *EmailVerificationService {
    return &EmailVerificationService{
        userRepo: userRepo,
        logger:   logrus.New(),
    }
}

func (evs *EmailVerificationService) VerifyEmail(ctx context.Context, email, code string) error {
    // 查找用户
    user, err := evs.userRepo.FindByEmail(ctx, email)
    if err != nil {
        evs.logger.Errorf("Failed to find user by email: %v", err)
        return errors.New("用户不存在")
    }
    
    // 检查验证码
    if user.EmailVerificationCode != code {
        evs.logger.Warnf("Invalid verification code for user: %s", email)
        return errors.New("验证码不正确")
    }
    
    // 检查验证码是否过期
    if time.Now().After(user.EmailVerificationExpiry) {
        evs.logger.Warnf("Verification code expired for user: %s", email)
        return errors.New("验证码已过期")
    }
    
    // 更新用户状态
    user.Status = UserStatusActive
    user.EmailVerified = true
    user.EmailVerificationCode = ""
    user.EmailVerificationExpiry = time.Time{}
    user.UpdatedAt = time.Now()
    
    if err := evs.userRepo.Update(ctx, user); err != nil {
        evs.logger.Errorf("Failed to update user status: %v", err)
        return errors.New("邮箱验证失败")
    }
    
    evs.logger.Infof("Email verified successfully for user: %s", email)
    return nil
}

func (evs *EmailVerificationService) ResendVerificationCode(ctx context.Context, email string) error {
    user, err := evs.userRepo.FindByEmail(ctx, email)
    if err != nil {
        return errors.New("用户不存在")
    }
    
    if user.EmailVerified {
        return errors.New("邮箱已验证")
    }
    
    // 生成新的验证码
    user.EmailVerificationCode = fmt.Sprintf("%06d", rand.Intn(1000000))
    user.EmailVerificationExpiry = time.Now().Add(24 * time.Hour)
    user.UpdatedAt = time.Now()
    
    if err := evs.userRepo.Update(ctx, user); err != nil {
        return errors.New("验证码生成失败")
    }
    
    // 发送验证邮件
    // ... 邮件发送逻辑
    
    return nil
}
```

## 📊 监控与指标

### 1. **注册指标收集**

```go
type RegistrationMetrics struct {
    registrationAttempts  *prometheus.CounterVec
    successfulRegistrations prometheus.Counter
    registrationErrors    *prometheus.CounterVec
    registrationDuration  prometheus.Histogram
    verificationAttempts  *prometheus.CounterVec
}

func NewRegistrationMetrics() *RegistrationMetrics {
    return &RegistrationMetrics{
        registrationAttempts: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "user_registration_attempts_total",
                Help: "Total number of user registration attempts",
            },
            []string{"status"},
        ),
        successfulRegistrations: prometheus.NewCounter(
            prometheus.CounterOpts{
                Name: "user_registrations_successful_total",
                Help: "Total number of successful user registrations",
            },
        ),
        registrationErrors: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "user_registration_errors_total",
                Help: "Total number of user registration errors",
            },
            []string{"error_type"},
        ),
        registrationDuration: prometheus.NewHistogram(
            prometheus.HistogramOpts{
                Name: "user_registration_duration_seconds",
                Help: "Duration of user registration process",
            },
        ),
        verificationAttempts: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "email_verification_attempts_total",
                Help: "Total number of email verification attempts",
            },
            []string{"status"},
        ),
    }
}
```

## 🔧 HTTP处理器

### 1. **注册API端点**

```go
type UserController struct {
    registrationService *RegistrationService
    verificationService *EmailVerificationService
    logger             *logrus.Logger
}

func (uc *UserController) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, RegisterErrorResponse{
            Success: false,
            Message: "请求参数错误",
            Errors:  []ValidationError{{Field: "request", Message: err.Error(), Code: "INVALID_REQUEST"}},
        })
        return
    }
    
    resp, err := uc.registrationService.Register(c.Request.Context(), &req)
    if err != nil {
        if validationErr, ok := err.(*ValidationErrors); ok {
            c.JSON(400, RegisterErrorResponse{
                Success: false,
                Message: "注册信息验证失败",
                Errors:  validationErr.Errors,
            })
            return
        }
        
        uc.logger.Errorf("Registration failed: %v", err)
        c.JSON(500, RegisterErrorResponse{
            Success: false,
            Message: "注册失败，请稍后重试",
        })
        return
    }
    
    c.JSON(200, resp)
}

func (uc *UserController) VerifyEmail(c *gin.Context) {
    email := c.Query("email")
    code := c.Query("code")
    
    if email == "" || code == "" {
        c.JSON(400, gin.H{
            "success": false,
            "message": "缺少必要参数",
        })
        return
    }
    
    if err := uc.verificationService.VerifyEmail(c.Request.Context(), email, code); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": err.Error(),
        })
        return
    }
    
    c.JSON(200, gin.H{
        "success": true,
        "message": "邮箱验证成功",
    })
}
```

## 📝 总结

用户注册功能为MovieInfo项目提供了安全可靠的用户入口：

**核心特性**：
1. **完整验证**：邮箱、密码、用户名等全面验证
2. **安全存储**：密码加密存储，防止泄露
3. **邮箱验证**：确保邮箱真实性和用户身份
4. **错误处理**：详细的错误信息和用户友好提示

**安全措施**：
- 密码强度验证
- 防重复注册检查
- 邮箱域名黑名单
- 验证码过期机制

**性能优化**：
- 异步邮件发送
- 数据库索引优化
- 缓存验证结果
- 监控指标收集

下一步，我们将实现登录认证系统，为用户提供安全的身份验证机制。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第28步：登录认证系统](28-authentication.md)
