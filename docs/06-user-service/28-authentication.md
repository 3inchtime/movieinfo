# 第28步：登录认证系统

## 📋 概述

登录认证系统是用户服务的核心功能，负责验证用户身份并建立安全的会话。一个完善的认证系统需要平衡安全性和用户体验，提供多种认证方式和完善的安全防护机制。

## 🎯 设计目标

### 1. **安全性**
- 密码安全验证
- 防暴力破解
- 会话安全管理
- 多因素认证支持

### 2. **用户体验**
- 快速登录响应
- 记住登录状态
- 友好的错误提示
- 多端登录支持

### 3. **可扩展性**
- 支持多种认证方式
- 第三方登录集成
- SSO单点登录
- 移动端适配

## 🔐 认证流程设计

### 1. **登录流程图**

```
用户访问登录页面 → 输入凭据 → 前端验证 → 提交登录请求
        ↓              ↓        ↓           ↓
    显示登录表单    填写用户信息  客户端验证   发送HTTP请求
        ↓              ↓        ↓           ↓
服务端验证凭据 → 检查用户状态 → 生成Token → 返回登录结果
        ↓              ↓          ↓           ↓
    密码验证        账户状态检查   JWT生成     响应客户端
        ↓              ↓          ↓           ↓
记录登录日志 ← 更新登录时间 ← 缓存会话 ← 设置Cookie
```

### 2. **认证架构**

```go
// 认证请求结构
type LoginRequest struct {
    Email      string `json:"email" binding:"required,email"`
    Password   string `json:"password" binding:"required"`
    RememberMe bool   `json:"remember_me"`
    DeviceInfo string `json:"device_info"`
}

// 认证响应结构
type LoginResponse struct {
    Success      bool      `json:"success"`
    Message      string    `json:"message"`
    AccessToken  string    `json:"access_token,omitempty"`
    RefreshToken string    `json:"refresh_token,omitempty"`
    ExpiresIn    int64     `json:"expires_in,omitempty"`
    User         *UserInfo `json:"user,omitempty"`
}

// 用户信息结构
type UserInfo struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    Username  string    `json:"username"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Avatar    string    `json:"avatar"`
    Role      string    `json:"role"`
    LastLogin time.Time `json:"last_login"`
}
```

## 🔧 核心组件实现

### 1. **认证服务**

```go
type AuthenticationService struct {
    userRepo      UserRepository
    passwordHasher PasswordHasher
    jwtManager    *JWTManager
    loginAttempts *LoginAttemptTracker
    sessionStore  SessionStore
    logger        *logrus.Logger
    metrics       *AuthMetrics
}

func NewAuthenticationService(
    userRepo UserRepository,
    passwordHasher PasswordHasher,
    jwtManager *JWTManager,
    loginAttempts *LoginAttemptTracker,
    sessionStore SessionStore,
) *AuthenticationService {
    return &AuthenticationService{
        userRepo:      userRepo,
        passwordHasher: passwordHasher,
        jwtManager:    jwtManager,
        loginAttempts: loginAttempts,
        sessionStore:  sessionStore,
        logger:        logrus.New(),
        metrics:       NewAuthMetrics(),
    }
}

func (as *AuthenticationService) Login(ctx context.Context, req *LoginRequest, clientIP string) (*LoginResponse, error) {
    start := time.Now()
    defer func() {
        as.metrics.ObserveLoginDuration(time.Since(start))
    }()

    // 检查登录尝试次数
    if blocked, remaining := as.loginAttempts.IsBlocked(req.Email, clientIP); blocked {
        as.metrics.IncBlockedLogins()
        return &LoginResponse{
            Success: false,
            Message: fmt.Sprintf("登录尝试过多，请%d分钟后重试", remaining/60),
        }, nil
    }

    // 查找用户
    user, err := as.userRepo.FindByEmail(ctx, req.Email)
    if err != nil {
        as.recordFailedLogin(req.Email, clientIP, "user_not_found")
        return &LoginResponse{
            Success: false,
            Message: "邮箱或密码错误",
        }, nil
    }

    // 检查用户状态
    if user.Status != UserStatusActive {
        as.recordFailedLogin(req.Email, clientIP, "account_inactive")
        return &LoginResponse{
            Success: false,
            Message: as.getUserStatusMessage(user.Status),
        }, nil
    }

    // 验证密码
    if !as.passwordHasher.CheckPassword(req.Password, user.Password) {
        as.recordFailedLogin(req.Email, clientIP, "invalid_password")
        return &LoginResponse{
            Success: false,
            Message: "邮箱或密码错误",
        }, nil
    }

    // 生成Token
    accessToken, refreshToken, err := as.generateTokens(user, req.RememberMe)
    if err != nil {
        as.logger.Errorf("Failed to generate tokens: %v", err)
        return nil, errors.New("登录失败，请稍后重试")
    }

    // 创建会话
    session := &UserSession{
        UserID:       user.ID,
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        DeviceInfo:   req.DeviceInfo,
        ClientIP:     clientIP,
        CreatedAt:    time.Now(),
        ExpiresAt:    time.Now().Add(as.getTokenExpiry(req.RememberMe)),
    }

    if err := as.sessionStore.Create(ctx, session); err != nil {
        as.logger.Errorf("Failed to create session: %v", err)
        // 会话创建失败不影响登录成功
    }

    // 更新用户登录信息
    as.updateUserLoginInfo(ctx, user, clientIP)

    // 清除登录尝试记录
    as.loginAttempts.ClearAttempts(req.Email, clientIP)

    // 记录成功登录
    as.recordSuccessfulLogin(user, clientIP)

    as.metrics.IncSuccessfulLogins()
    as.logger.Infof("User logged in successfully: %s", user.Email)

    return &LoginResponse{
        Success:      true,
        Message:      "登录成功",
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    as.getTokenExpiry(req.RememberMe).Milliseconds(),
        User:         as.buildUserInfo(user),
    }, nil
}

func (as *AuthenticationService) recordFailedLogin(email, clientIP, reason string) {
    as.loginAttempts.RecordAttempt(email, clientIP)
    as.metrics.IncFailedLogins(reason)
    as.logger.Warnf("Failed login attempt: email=%s, ip=%s, reason=%s", email, clientIP, reason)
}

func (as *AuthenticationService) recordSuccessfulLogin(user *User, clientIP string) {
    loginLog := &LoginLog{
        UserID:    user.ID,
        Email:     user.Email,
        ClientIP:  clientIP,
        Success:   true,
        Timestamp: time.Now(),
    }
    
    // 异步记录登录日志
    go func() {
        if err := as.userRepo.CreateLoginLog(context.Background(), loginLog); err != nil {
            as.logger.Errorf("Failed to create login log: %v", err)
        }
    }()
}

func (as *AuthenticationService) generateTokens(user *User, rememberMe bool) (string, string, error) {
    // 生成访问令牌
    accessClaims := &JWTClaims{
        UserID:   user.ID,
        Email:    user.Email,
        Username: user.Username,
        Role:     user.Role,
        Type:     "access",
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(as.getAccessTokenExpiry()).Unix(),
            IssuedAt:  time.Now().Unix(),
            Issuer:    "movieinfo",
        },
    }

    accessToken, err := as.jwtManager.Generate(accessClaims)
    if err != nil {
        return "", "", err
    }

    // 生成刷新令牌
    refreshClaims := &JWTClaims{
        UserID: user.ID,
        Email:  user.Email,
        Type:   "refresh",
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(as.getRefreshTokenExpiry(rememberMe)).Unix(),
            IssuedAt:  time.Now().Unix(),
            Issuer:    "movieinfo",
        },
    }

    refreshToken, err := as.jwtManager.Generate(refreshClaims)
    if err != nil {
        return "", "", err
    }

    return accessToken, refreshToken, nil
}
```

### 2. **登录尝试追踪器**

```go
type LoginAttemptTracker struct {
    redis      *redis.Client
    maxAttempts int
    blockDuration time.Duration
    logger     *logrus.Logger
}

func NewLoginAttemptTracker(redis *redis.Client) *LoginAttemptTracker {
    return &LoginAttemptTracker{
        redis:         redis,
        maxAttempts:   5,
        blockDuration: 15 * time.Minute,
        logger:        logrus.New(),
    }
}

func (lat *LoginAttemptTracker) RecordAttempt(email, clientIP string) {
    ctx := context.Background()
    
    // 记录邮箱尝试次数
    emailKey := fmt.Sprintf("login_attempts:email:%s", email)
    lat.redis.Incr(ctx, emailKey)
    lat.redis.Expire(ctx, emailKey, lat.blockDuration)
    
    // 记录IP尝试次数
    ipKey := fmt.Sprintf("login_attempts:ip:%s", clientIP)
    lat.redis.Incr(ctx, ipKey)
    lat.redis.Expire(ctx, ipKey, lat.blockDuration)
}

func (lat *LoginAttemptTracker) IsBlocked(email, clientIP string) (bool, int64) {
    ctx := context.Background()
    
    // 检查邮箱尝试次数
    emailKey := fmt.Sprintf("login_attempts:email:%s", email)
    emailAttempts, _ := lat.redis.Get(ctx, emailKey).Int()
    
    // 检查IP尝试次数
    ipKey := fmt.Sprintf("login_attempts:ip:%s", clientIP)
    ipAttempts, _ := lat.redis.Get(ctx, ipKey).Int()
    
    if emailAttempts >= lat.maxAttempts || ipAttempts >= lat.maxAttempts {
        // 获取剩余阻塞时间
        var remaining int64
        if emailAttempts >= lat.maxAttempts {
            remaining, _ = lat.redis.TTL(ctx, emailKey).Result().Seconds()
        } else {
            remaining, _ = lat.redis.TTL(ctx, ipKey).Result().Seconds()
        }
        return true, remaining
    }
    
    return false, 0
}

func (lat *LoginAttemptTracker) ClearAttempts(email, clientIP string) {
    ctx := context.Background()
    
    emailKey := fmt.Sprintf("login_attempts:email:%s", email)
    ipKey := fmt.Sprintf("login_attempts:ip:%s", clientIP)
    
    lat.redis.Del(ctx, emailKey, ipKey)
}
```

### 3. **会话管理**

```go
type SessionStore interface {
    Create(ctx context.Context, session *UserSession) error
    Get(ctx context.Context, token string) (*UserSession, error)
    Update(ctx context.Context, session *UserSession) error
    Delete(ctx context.Context, token string) error
    DeleteUserSessions(ctx context.Context, userID string) error
}

type RedisSessionStore struct {
    redis  *redis.Client
    prefix string
    logger *logrus.Logger
}

func NewRedisSessionStore(redis *redis.Client) *RedisSessionStore {
    return &RedisSessionStore{
        redis:  redis,
        prefix: "session:",
        logger: logrus.New(),
    }
}

func (rss *RedisSessionStore) Create(ctx context.Context, session *UserSession) error {
    key := rss.prefix + session.AccessToken
    
    data, err := json.Marshal(session)
    if err != nil {
        return err
    }
    
    expiry := session.ExpiresAt.Sub(time.Now())
    return rss.redis.Set(ctx, key, data, expiry).Err()
}

func (rss *RedisSessionStore) Get(ctx context.Context, token string) (*UserSession, error) {
    key := rss.prefix + token
    
    data, err := rss.redis.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return nil, errors.New("session not found")
        }
        return nil, err
    }
    
    var session UserSession
    if err := json.Unmarshal([]byte(data), &session); err != nil {
        return nil, err
    }
    
    return &session, nil
}

func (rss *RedisSessionStore) Delete(ctx context.Context, token string) error {
    key := rss.prefix + token
    return rss.redis.Del(ctx, key).Err()
}
```

### 4. **Token刷新机制**

```go
func (as *AuthenticationService) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
    // 验证刷新令牌
    claims, err := as.jwtManager.Verify(refreshToken)
    if err != nil {
        as.metrics.IncTokenRefreshErrors("invalid_token")
        return &LoginResponse{
            Success: false,
            Message: "无效的刷新令牌",
        }, nil
    }
    
    if claims.Type != "refresh" {
        as.metrics.IncTokenRefreshErrors("wrong_token_type")
        return &LoginResponse{
            Success: false,
            Message: "令牌类型错误",
        }, nil
    }
    
    // 查找用户
    user, err := as.userRepo.FindByID(ctx, claims.UserID)
    if err != nil {
        as.metrics.IncTokenRefreshErrors("user_not_found")
        return &LoginResponse{
            Success: false,
            Message: "用户不存在",
        }, nil
    }
    
    // 检查用户状态
    if user.Status != UserStatusActive {
        as.metrics.IncTokenRefreshErrors("account_inactive")
        return &LoginResponse{
            Success: false,
            Message: "账户已被禁用",
        }, nil
    }
    
    // 生成新的访问令牌
    newAccessClaims := &JWTClaims{
        UserID:   user.ID,
        Email:    user.Email,
        Username: user.Username,
        Role:     user.Role,
        Type:     "access",
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(as.getAccessTokenExpiry()).Unix(),
            IssuedAt:  time.Now().Unix(),
            Issuer:    "movieinfo",
        },
    }
    
    newAccessToken, err := as.jwtManager.Generate(newAccessClaims)
    if err != nil {
        as.logger.Errorf("Failed to generate new access token: %v", err)
        return nil, errors.New("令牌刷新失败")
    }
    
    as.metrics.IncSuccessfulTokenRefresh()
    
    return &LoginResponse{
        Success:     true,
        Message:     "令牌刷新成功",
        AccessToken: newAccessToken,
        ExpiresIn:   as.getAccessTokenExpiry().Milliseconds(),
        User:        as.buildUserInfo(user),
    }, nil
}
```

## 🔒 安全增强

### 1. **多因素认证**

```go
type MFAService struct {
    userRepo   UserRepository
    totpSecret string
    logger     *logrus.Logger
}

func (mfa *MFAService) EnableTOTP(ctx context.Context, userID string) (*TOTPSetupInfo, error) {
    user, err := mfa.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 生成TOTP密钥
    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "MovieInfo",
        AccountName: user.Email,
    })
    if err != nil {
        return nil, err
    }
    
    // 保存密钥到用户记录
    user.TOTPSecret = key.Secret()
    user.TOTPEnabled = false // 需要验证后才启用
    
    if err := mfa.userRepo.Update(ctx, user); err != nil {
        return nil, err
    }
    
    return &TOTPSetupInfo{
        Secret:  key.Secret(),
        QRCode:  key.URL(),
        BackupCodes: mfa.generateBackupCodes(),
    }, nil
}

func (mfa *MFAService) VerifyTOTP(ctx context.Context, userID, code string) error {
    user, err := mfa.userRepo.FindByID(ctx, userID)
    if err != nil {
        return err
    }
    
    if !totp.Validate(code, user.TOTPSecret) {
        return errors.New("验证码错误")
    }
    
    // 启用TOTP
    user.TOTPEnabled = true
    return mfa.userRepo.Update(ctx, user)
}
```

### 2. **设备管理**

```go
type DeviceManager struct {
    redis  *redis.Client
    logger *logrus.Logger
}

func (dm *DeviceManager) RegisterDevice(ctx context.Context, userID, deviceInfo, clientIP string) string {
    deviceID := dm.generateDeviceID(deviceInfo, clientIP)
    
    device := &UserDevice{
        ID:         deviceID,
        UserID:     userID,
        DeviceInfo: deviceInfo,
        ClientIP:   clientIP,
        FirstSeen:  time.Now(),
        LastSeen:   time.Now(),
        Trusted:    false,
    }
    
    key := fmt.Sprintf("device:%s:%s", userID, deviceID)
    data, _ := json.Marshal(device)
    dm.redis.Set(ctx, key, data, 30*24*time.Hour) // 30天过期
    
    return deviceID
}

func (dm *DeviceManager) IsDeviceTrusted(ctx context.Context, userID, deviceID string) bool {
    key := fmt.Sprintf("device:%s:%s", userID, deviceID)
    
    data, err := dm.redis.Get(ctx, key).Result()
    if err != nil {
        return false
    }
    
    var device UserDevice
    if err := json.Unmarshal([]byte(data), &device); err != nil {
        return false
    }
    
    return device.Trusted
}
```

## 📊 监控指标

### 1. **认证指标**

```go
type AuthMetrics struct {
    loginAttempts     *prometheus.CounterVec
    successfulLogins  prometheus.Counter
    failedLogins      *prometheus.CounterVec
    blockedLogins     prometheus.Counter
    loginDuration     prometheus.Histogram
    tokenRefresh      *prometheus.CounterVec
    activeSessions    prometheus.Gauge
}

func NewAuthMetrics() *AuthMetrics {
    return &AuthMetrics{
        loginAttempts: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "auth_login_attempts_total",
                Help: "Total number of login attempts",
            },
            []string{"status"},
        ),
        successfulLogins: prometheus.NewCounter(
            prometheus.CounterOpts{
                Name: "auth_successful_logins_total",
                Help: "Total number of successful logins",
            },
        ),
        failedLogins: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "auth_failed_logins_total",
                Help: "Total number of failed logins",
            },
            []string{"reason"},
        ),
        blockedLogins: prometheus.NewCounter(
            prometheus.CounterOpts{
                Name: "auth_blocked_logins_total",
                Help: "Total number of blocked login attempts",
            },
        ),
        loginDuration: prometheus.NewHistogram(
            prometheus.HistogramOpts{
                Name: "auth_login_duration_seconds",
                Help: "Duration of login process",
            },
        ),
        tokenRefresh: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "auth_token_refresh_total",
                Help: "Total number of token refresh attempts",
            },
            []string{"status"},
        ),
        activeSessions: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "auth_active_sessions",
                Help: "Number of active user sessions",
            },
        ),
    }
}
```

## 🔧 HTTP处理器

### 1. **登录API端点**

```go
func (uc *UserController) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": "请求参数错误",
        })
        return
    }
    
    clientIP := c.ClientIP()
    resp, err := uc.authService.Login(c.Request.Context(), &req, clientIP)
    if err != nil {
        uc.logger.Errorf("Login failed: %v", err)
        c.JSON(500, gin.H{
            "success": false,
            "message": "登录失败，请稍后重试",
        })
        return
    }
    
    // 设置Cookie（如果登录成功）
    if resp.Success && resp.AccessToken != "" {
        maxAge := 3600 // 1小时
        if req.RememberMe {
            maxAge = 30 * 24 * 3600 // 30天
        }
        
        c.SetCookie("access_token", resp.AccessToken, maxAge, "/", "", false, true)
        c.SetCookie("refresh_token", resp.RefreshToken, maxAge, "/", "", false, true)
    }
    
    c.JSON(200, resp)
}

func (uc *UserController) Logout(c *gin.Context) {
    token := uc.extractToken(c)
    if token != "" {
        // 删除会话
        uc.authService.Logout(c.Request.Context(), token)
    }
    
    // 清除Cookie
    c.SetCookie("access_token", "", -1, "/", "", false, true)
    c.SetCookie("refresh_token", "", -1, "/", "", false, true)
    
    c.JSON(200, gin.H{
        "success": true,
        "message": "退出登录成功",
    })
}

func (uc *UserController) RefreshToken(c *gin.Context) {
    refreshToken := c.GetHeader("Refresh-Token")
    if refreshToken == "" {
        refreshToken, _ = c.Cookie("refresh_token")
    }
    
    if refreshToken == "" {
        c.JSON(401, gin.H{
            "success": false,
            "message": "缺少刷新令牌",
        })
        return
    }
    
    resp, err := uc.authService.RefreshToken(c.Request.Context(), refreshToken)
    if err != nil {
        c.JSON(500, gin.H{
            "success": false,
            "message": "令牌刷新失败",
        })
        return
    }
    
    c.JSON(200, resp)
}
```

## 📝 总结

登录认证系统为MovieInfo项目提供了安全可靠的身份验证：

**核心特性**：
1. **安全认证**：密码验证、防暴力破解、会话管理
2. **Token管理**：JWT访问令牌和刷新令牌机制
3. **多因素认证**：TOTP支持，增强安全性
4. **设备管理**：设备识别和信任机制

**安全措施**：
- 登录尝试限制和IP封禁
- 安全的会话管理
- Token过期和刷新机制
- 设备指纹识别

**性能优化**：
- Redis缓存会话信息
- 异步日志记录
- 连接池优化
- 监控指标收集

下一步，我们将实现JWT Token管理系统，提供完整的令牌生命周期管理。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第29步：JWT Token管理](29-jwt-management.md)
