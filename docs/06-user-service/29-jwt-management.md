# 第29步：JWT Token 管理

## 📋 概述

JWT (JSON Web Token) 是现代Web应用中广泛使用的身份认证和授权机制。对于MovieInfo项目，我们需要实现一个完整的JWT管理系统，包括Token生成、验证、刷新和撤销等功能。

## 🎯 设计目标

### 1. **安全性**
- 强加密算法保护
- Token签名验证
- 过期时间控制
- 撤销机制支持

### 2. **性能**
- 快速Token验证
- 缓存优化
- 批量操作支持
- 内存高效使用

### 3. **可扩展性**
- 多种Token类型
- 自定义Claims支持
- 密钥轮换机制
- 分布式部署支持

## 🔧 JWT架构设计

### 1. **Token结构设计**

```go
// JWT Claims结构
type JWTClaims struct {
    UserID   string `json:"user_id"`
    Email    string `json:"email"`
    Username string `json:"username"`
    Role     string `json:"role"`
    Type     string `json:"type"` // access, refresh
    DeviceID string `json:"device_id,omitempty"`
    jwt.StandardClaims
}

// Token配置
type TokenConfig struct {
    AccessTokenExpiry  time.Duration `yaml:"access_token_expiry"`
    RefreshTokenExpiry time.Duration `yaml:"refresh_token_expiry"`
    SigningMethod      string        `yaml:"signing_method"`
    SecretKey          string        `yaml:"secret_key"`
    Issuer             string        `yaml:"issuer"`
    EnableRefresh      bool          `yaml:"enable_refresh"`
}

// Token对象
type Token struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token,omitempty"`
    TokenType    string    `json:"token_type"`
    ExpiresIn    int64     `json:"expires_in"`
    IssuedAt     time.Time `json:"issued_at"`
    ExpiresAt    time.Time `json:"expires_at"`
}
```

### 2. **JWT管理器实现**

```go
type JWTManager struct {
    config       *TokenConfig
    signingKey   []byte
    blacklist    TokenBlacklist
    keyRotator   *KeyRotator
    logger       *logrus.Logger
    metrics      *JWTMetrics
}

func NewJWTManager(config *TokenConfig, blacklist TokenBlacklist) *JWTManager {
    return &JWTManager{
        config:     config,
        signingKey: []byte(config.SecretKey),
        blacklist:  blacklist,
        keyRotator: NewKeyRotator(config.SecretKey),
        logger:     logrus.New(),
        metrics:    NewJWTMetrics(),
    }
}

// 生成Token
func (jm *JWTManager) Generate(claims *JWTClaims) (string, error) {
    start := time.Now()
    defer func() {
        jm.metrics.ObserveTokenGeneration(time.Since(start))
    }()

    // 设置标准Claims
    now := time.Now()
    claims.IssuedAt = now.Unix()
    claims.Issuer = jm.config.Issuer
    
    if claims.ExpiresAt == 0 {
        expiry := jm.getExpiryForType(claims.Type)
        claims.ExpiresAt = now.Add(expiry).Unix()
    }

    // 创建Token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    // 添加Key ID用于密钥轮换
    token.Header["kid"] = jm.keyRotator.GetCurrentKeyID()

    // 签名Token
    tokenString, err := token.SignedString(jm.signingKey)
    if err != nil {
        jm.metrics.IncTokenGenerationErrors()
        jm.logger.Errorf("Failed to sign token: %v", err)
        return "", err
    }

    jm.metrics.IncTokenGenerated(claims.Type)
    return tokenString, nil
}

// 验证Token
func (jm *JWTManager) Verify(tokenString string) (*JWTClaims, error) {
    start := time.Now()
    defer func() {
        jm.metrics.ObserveTokenVerification(time.Since(start))
    }()

    // 检查黑名单
    if jm.blacklist.IsBlacklisted(tokenString) {
        jm.metrics.IncTokenVerificationErrors("blacklisted")
        return nil, errors.New("token has been revoked")
    }

    // 解析Token
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        // 验证签名方法
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }

        // 获取密钥（支持密钥轮换）
        keyID, ok := token.Header["kid"].(string)
        if ok {
            if key := jm.keyRotator.GetKey(keyID); key != nil {
                return key, nil
            }
        }

        return jm.signingKey, nil
    })

    if err != nil {
        jm.metrics.IncTokenVerificationErrors("parse_failed")
        jm.logger.Warnf("Token verification failed: %v", err)
        return nil, err
    }

    // 验证Token有效性
    if !token.Valid {
        jm.metrics.IncTokenVerificationErrors("invalid")
        return nil, errors.New("invalid token")
    }

    claims, ok := token.Claims.(*JWTClaims)
    if !ok {
        jm.metrics.IncTokenVerificationErrors("invalid_claims")
        return nil, errors.New("invalid token claims")
    }

    // 验证Issuer
    if claims.Issuer != jm.config.Issuer {
        jm.metrics.IncTokenVerificationErrors("invalid_issuer")
        return nil, errors.New("invalid token issuer")
    }

    jm.metrics.IncTokenVerified(claims.Type)
    return claims, nil
}

// 刷新Token
func (jm *JWTManager) Refresh(refreshToken string) (*Token, error) {
    // 验证刷新Token
    claims, err := jm.Verify(refreshToken)
    if err != nil {
        return nil, err
    }

    if claims.Type != "refresh" {
        jm.metrics.IncTokenRefreshErrors("wrong_type")
        return nil, errors.New("invalid token type for refresh")
    }

    // 生成新的访问Token
    newAccessClaims := &JWTClaims{
        UserID:   claims.UserID,
        Email:    claims.Email,
        Username: claims.Username,
        Role:     claims.Role,
        Type:     "access",
        DeviceID: claims.DeviceID,
    }

    accessToken, err := jm.Generate(newAccessClaims)
    if err != nil {
        jm.metrics.IncTokenRefreshErrors("generation_failed")
        return nil, err
    }

    // 可选：生成新的刷新Token
    var newRefreshToken string
    if jm.config.EnableRefresh {
        newRefreshClaims := &JWTClaims{
            UserID:   claims.UserID,
            Email:    claims.Email,
            Type:     "refresh",
            DeviceID: claims.DeviceID,
        }

        newRefreshToken, err = jm.Generate(newRefreshClaims)
        if err != nil {
            jm.metrics.IncTokenRefreshErrors("refresh_generation_failed")
            return nil, err
        }

        // 将旧的刷新Token加入黑名单
        jm.blacklist.Add(refreshToken, time.Unix(claims.ExpiresAt, 0))
    }

    jm.metrics.IncTokenRefreshed()

    return &Token{
        AccessToken:  accessToken,
        RefreshToken: newRefreshToken,
        TokenType:    "Bearer",
        ExpiresIn:    int64(jm.config.AccessTokenExpiry.Seconds()),
        IssuedAt:     time.Now(),
        ExpiresAt:    time.Now().Add(jm.config.AccessTokenExpiry),
    }, nil
}

// 撤销Token
func (jm *JWTManager) Revoke(tokenString string) error {
    claims, err := jm.Verify(tokenString)
    if err != nil {
        return err
    }

    // 添加到黑名单
    expiry := time.Unix(claims.ExpiresAt, 0)
    if err := jm.blacklist.Add(tokenString, expiry); err != nil {
        jm.metrics.IncTokenRevocationErrors()
        return err
    }

    jm.metrics.IncTokenRevoked(claims.Type)
    jm.logger.Infof("Token revoked for user: %s", claims.UserID)
    return nil
}

func (jm *JWTManager) getExpiryForType(tokenType string) time.Duration {
    switch tokenType {
    case "access":
        return jm.config.AccessTokenExpiry
    case "refresh":
        return jm.config.RefreshTokenExpiry
    default:
        return jm.config.AccessTokenExpiry
    }
}
```

### 3. **Token黑名单实现**

```go
type TokenBlacklist interface {
    Add(token string, expiry time.Time) error
    IsBlacklisted(token string) bool
    Remove(token string) error
    Cleanup() error
}

type RedisTokenBlacklist struct {
    redis  *redis.Client
    prefix string
    logger *logrus.Logger
}

func NewRedisTokenBlacklist(redis *redis.Client) *RedisTokenBlacklist {
    blacklist := &RedisTokenBlacklist{
        redis:  redis,
        prefix: "blacklist:",
        logger: logrus.New(),
    }

    // 启动清理任务
    go blacklist.startCleanupTask()

    return blacklist
}

func (rtb *RedisTokenBlacklist) Add(token string, expiry time.Time) error {
    ctx := context.Background()
    key := rtb.prefix + rtb.hashToken(token)
    
    // 计算TTL
    ttl := expiry.Sub(time.Now())
    if ttl <= 0 {
        return nil // Token已过期，无需加入黑名单
    }

    return rtb.redis.Set(ctx, key, "1", ttl).Err()
}

func (rtb *RedisTokenBlacklist) IsBlacklisted(token string) bool {
    ctx := context.Background()
    key := rtb.prefix + rtb.hashToken(token)
    
    exists, err := rtb.redis.Exists(ctx, key).Result()
    if err != nil {
        rtb.logger.Errorf("Failed to check blacklist: %v", err)
        return false
    }

    return exists > 0
}

func (rtb *RedisTokenBlacklist) Remove(token string) error {
    ctx := context.Background()
    key := rtb.prefix + rtb.hashToken(token)
    return rtb.redis.Del(ctx, key).Err()
}

func (rtb *RedisTokenBlacklist) hashToken(token string) string {
    hash := sha256.Sum256([]byte(token))
    return hex.EncodeToString(hash[:])
}

func (rtb *RedisTokenBlacklist) startCleanupTask() {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        rtb.Cleanup()
    }
}

func (rtb *RedisTokenBlacklist) Cleanup() error {
    // Redis会自动清理过期的键，这里可以添加额外的清理逻辑
    rtb.logger.Info("Token blacklist cleanup completed")
    return nil
}
```

### 4. **密钥轮换机制**

```go
type KeyRotator struct {
    currentKey   []byte
    currentKeyID string
    keys         map[string][]byte
    rotationInterval time.Duration
    mutex        sync.RWMutex
    logger       *logrus.Logger
}

func NewKeyRotator(initialKey string) *KeyRotator {
    kr := &KeyRotator{
        currentKey:   []byte(initialKey),
        currentKeyID: "key-1",
        keys:         make(map[string][]byte),
        rotationInterval: 24 * time.Hour, // 24小时轮换一次
        logger:       logrus.New(),
    }

    kr.keys[kr.currentKeyID] = kr.currentKey

    // 启动密钥轮换任务
    go kr.startRotationTask()

    return kr
}

func (kr *KeyRotator) GetCurrentKeyID() string {
    kr.mutex.RLock()
    defer kr.mutex.RUnlock()
    return kr.currentKeyID
}

func (kr *KeyRotator) GetKey(keyID string) []byte {
    kr.mutex.RLock()
    defer kr.mutex.RUnlock()
    return kr.keys[keyID]
}

func (kr *KeyRotator) RotateKey() {
    kr.mutex.Lock()
    defer kr.mutex.Unlock()

    // 生成新密钥
    newKey := make([]byte, 32)
    if _, err := rand.Read(newKey); err != nil {
        kr.logger.Errorf("Failed to generate new key: %v", err)
        return
    }

    // 更新密钥
    oldKeyID := kr.currentKeyID
    kr.currentKeyID = fmt.Sprintf("key-%d", time.Now().Unix())
    kr.currentKey = newKey
    kr.keys[kr.currentKeyID] = newKey

    kr.logger.Infof("Key rotated from %s to %s", oldKeyID, kr.currentKeyID)

    // 清理旧密钥（保留一段时间以支持现有Token）
    go func() {
        time.Sleep(48 * time.Hour) // 48小时后清理
        kr.mutex.Lock()
        delete(kr.keys, oldKeyID)
        kr.mutex.Unlock()
        kr.logger.Infof("Old key %s cleaned up", oldKeyID)
    }()
}

func (kr *KeyRotator) startRotationTask() {
    ticker := time.NewTicker(kr.rotationInterval)
    defer ticker.Stop()

    for range ticker.C {
        kr.RotateKey()
    }
}
```

### 5. **Token工具函数**

```go
type TokenUtils struct {
    jwtManager *JWTManager
}

func NewTokenUtils(jwtManager *JWTManager) *TokenUtils {
    return &TokenUtils{
        jwtManager: jwtManager,
    }
}

// 从HTTP请求中提取Token
func (tu *TokenUtils) ExtractTokenFromRequest(r *http.Request) string {
    // 从Authorization头提取
    authHeader := r.Header.Get("Authorization")
    if authHeader != "" {
        parts := strings.Split(authHeader, " ")
        if len(parts) == 2 && parts[0] == "Bearer" {
            return parts[1]
        }
    }

    // 从Cookie提取
    if cookie, err := r.Cookie("access_token"); err == nil {
        return cookie.Value
    }

    // 从查询参数提取
    return r.URL.Query().Get("token")
}

// 验证Token并返回用户信息
func (tu *TokenUtils) ValidateAndGetUser(tokenString string) (*JWTClaims, error) {
    if tokenString == "" {
        return nil, errors.New("missing token")
    }

    claims, err := tu.jwtManager.Verify(tokenString)
    if err != nil {
        return nil, err
    }

    if claims.Type != "access" {
        return nil, errors.New("invalid token type")
    }

    return claims, nil
}

// 生成Token对
func (tu *TokenUtils) GenerateTokenPair(userID, email, username, role, deviceID string) (*Token, error) {
    // 生成访问Token
    accessClaims := &JWTClaims{
        UserID:   userID,
        Email:    email,
        Username: username,
        Role:     role,
        Type:     "access",
        DeviceID: deviceID,
    }

    accessToken, err := tu.jwtManager.Generate(accessClaims)
    if err != nil {
        return nil, err
    }

    // 生成刷新Token
    refreshClaims := &JWTClaims{
        UserID:   userID,
        Email:    email,
        Type:     "refresh",
        DeviceID: deviceID,
    }

    refreshToken, err := tu.jwtManager.Generate(refreshClaims)
    if err != nil {
        return nil, err
    }

    return &Token{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        TokenType:    "Bearer",
        ExpiresIn:    int64(tu.jwtManager.config.AccessTokenExpiry.Seconds()),
        IssuedAt:     time.Now(),
        ExpiresAt:    time.Now().Add(tu.jwtManager.config.AccessTokenExpiry),
    }, nil
}
```

## 📊 监控指标

### 1. **JWT指标收集**

```go
type JWTMetrics struct {
    tokenGenerated     *prometheus.CounterVec
    tokenVerified      *prometheus.CounterVec
    tokenRefreshed     prometheus.Counter
    tokenRevoked       *prometheus.CounterVec
    generationDuration prometheus.Histogram
    verificationDuration prometheus.Histogram
    generationErrors   prometheus.Counter
    verificationErrors *prometheus.CounterVec
    refreshErrors      *prometheus.CounterVec
    revocationErrors   prometheus.Counter
    activeTokens       prometheus.Gauge
}

func NewJWTMetrics() *JWTMetrics {
    return &JWTMetrics{
        tokenGenerated: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "jwt_tokens_generated_total",
                Help: "Total number of JWT tokens generated",
            },
            []string{"type"},
        ),
        tokenVerified: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "jwt_tokens_verified_total",
                Help: "Total number of JWT tokens verified",
            },
            []string{"type"},
        ),
        tokenRefreshed: prometheus.NewCounter(
            prometheus.CounterOpts{
                Name: "jwt_tokens_refreshed_total",
                Help: "Total number of JWT tokens refreshed",
            },
        ),
        tokenRevoked: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "jwt_tokens_revoked_total",
                Help: "Total number of JWT tokens revoked",
            },
            []string{"type"},
        ),
        generationDuration: prometheus.NewHistogram(
            prometheus.HistogramOpts{
                Name: "jwt_token_generation_duration_seconds",
                Help: "Duration of JWT token generation",
            },
        ),
        verificationDuration: prometheus.NewHistogram(
            prometheus.HistogramOpts{
                Name: "jwt_token_verification_duration_seconds",
                Help: "Duration of JWT token verification",
            },
        ),
        generationErrors: prometheus.NewCounter(
            prometheus.CounterOpts{
                Name: "jwt_token_generation_errors_total",
                Help: "Total number of JWT token generation errors",
            },
        ),
        verificationErrors: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "jwt_token_verification_errors_total",
                Help: "Total number of JWT token verification errors",
            },
            []string{"reason"},
        ),
        refreshErrors: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "jwt_token_refresh_errors_total",
                Help: "Total number of JWT token refresh errors",
            },
            []string{"reason"},
        ),
        revocationErrors: prometheus.NewCounter(
            prometheus.CounterOpts{
                Name: "jwt_token_revocation_errors_total",
                Help: "Total number of JWT token revocation errors",
            },
        ),
        activeTokens: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "jwt_active_tokens",
                Help: "Number of active JWT tokens",
            },
        ),
    }
}
```

## 🔧 配置示例

### 1. **JWT配置文件**

```yaml
jwt:
  access_token_expiry: 1h
  refresh_token_expiry: 720h  # 30 days
  signing_method: HS256
  secret_key: ${JWT_SECRET_KEY}
  issuer: movieinfo
  enable_refresh: true

# Token黑名单配置
token_blacklist:
  redis_prefix: "blacklist:"
  cleanup_interval: 1h

# 密钥轮换配置
key_rotation:
  enabled: true
  interval: 24h
  keep_old_keys_duration: 48h
```

## 📝 总结

JWT Token管理系统为MovieInfo项目提供了完整的令牌生命周期管理：

**核心功能**：
1. **Token生成**：安全的JWT生成机制
2. **Token验证**：快速可靠的验证流程
3. **Token刷新**：无缝的令牌刷新机制
4. **Token撤销**：黑名单管理和撤销功能

**安全特性**：
- 强加密算法保护
- 密钥轮换机制
- Token黑名单管理
- 过期时间控制

**性能优化**：
- Redis缓存黑名单
- 快速Token验证
- 批量操作支持
- 监控指标收集

下一步，我们将实现密码重置功能，为用户提供安全的密码恢复机制。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第30步：密码重置功能](30-password-reset.md)
