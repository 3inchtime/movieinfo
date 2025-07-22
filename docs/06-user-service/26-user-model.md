# 5.1 用户服务实现

## 5.1.1 概述

用户服务是MovieInfo项目的核心业务模块之一，负责处理用户注册、登录、个人信息管理、权限控制等功能。作为业务逻辑层的重要组成部分，用户服务需要实现复杂的业务规则、安全控制和数据一致性保证。

## 5.1.2 为什么需要专业的用户服务？

### 5.1.2.1 **安全性要求**
- **身份认证**：确保用户身份的真实性和唯一性
- **密码安全**：安全的密码存储和验证机制
- **会话管理**：安全的用户会话和令牌管理
- **权限控制**：细粒度的用户权限和访问控制

### 5.1.2.2 **业务复杂性**
- **注册流程**：复杂的用户注册和邮箱验证流程
- **登录逻辑**：多种登录方式和安全检查
- **个人资料**：用户个人信息的管理和隐私保护
- **社交功能**：用户关注、收藏等社交功能

### 5.1.2.3 **数据一致性**
- **事务处理**：确保用户相关操作的事务一致性
- **缓存同步**：用户数据在缓存和数据库间的同步
- **状态管理**：用户状态变更的一致性管理
- **关联数据**：用户与其他实体关联数据的一致性

### 5.1.2.4 **性能优化**
- **缓存策略**：用户数据的智能缓存策略
- **查询优化**：高效的用户查询和检索
- **批量操作**：用户批量操作的性能优化
- **并发控制**：高并发场景下的用户操作控制

## 5.1.3 用户服务架构设计

### 5.1.3.1 **服务层次结构**

```
用户服务架构
├── 服务接口层 (Service Interface)
│   ├── 用户服务接口 (UserService)
│   ├── 认证服务接口 (AuthService)
│   ├── 个人资料服务接口 (ProfileService)
│   └── 社交服务接口 (SocialService)
├── 业务逻辑层 (Business Logic)
│   ├── 注册逻辑 (Registration Logic)
│   ├── 登录逻辑 (Authentication Logic)
│   ├── 权限逻辑 (Authorization Logic)
│   └── 社交逻辑 (Social Logic)
├── 数据访问层 (Data Access)
│   ├── 用户Repository (UserRepository)
│   ├── 会话Repository (SessionRepository)
│   ├── 权限Repository (PermissionRepository)
│   └── 社交Repository (SocialRepository)
├── 外部集成层 (External Integration)
│   ├── 邮件服务 (Email Service)
│   ├── 短信服务 (SMS Service)
│   ├── 第三方登录 (OAuth Service)
│   └── 文件存储 (File Storage)
└── 缓存层 (Cache Layer)
    ├── 用户缓存 (User Cache)
    ├── 会话缓存 (Session Cache)
    ├── 权限缓存 (Permission Cache)
    └── 社交缓存 (Social Cache)
```

### 5.1.3.2 **用户服务接口定义**

#### 5.1.3.2.1 核心用户服务接口
```go
// internal/services/user/interface.go
package user

import (
    "context"
    
    "github.com/yourname/movieinfo/internal/models"
    "github.com/yourname/movieinfo/internal/repository"
)

// UserService 用户服务接口
type UserService interface {
    // 用户注册和认证
    Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
    Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
    Logout(ctx context.Context, userID int64, token string) error
    RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
    
    // 用户信息管理
    GetUser(ctx context.Context, userID int64) (*models.User, error)
    GetUserByEmail(ctx context.Context, email string) (*models.User, error)
    GetUserByUsername(ctx context.Context, username string) (*models.User, error)
    UpdateUser(ctx context.Context, userID int64, req *UpdateUserRequest) (*models.User, error)
    DeleteUser(ctx context.Context, userID int64) error
    
    // 密码管理
    ChangePassword(ctx context.Context, userID int64, req *ChangePasswordRequest) error
    ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
    ForgotPassword(ctx context.Context, email string) error
    
    // 邮箱验证
    SendVerificationEmail(ctx context.Context, userID int64) error
    VerifyEmail(ctx context.Context, token string) error
    
    // 用户状态管理
    ActivateUser(ctx context.Context, userID int64) error
    DeactivateUser(ctx context.Context, userID int64) error
    BanUser(ctx context.Context, userID int64, reason string) error
    UnbanUser(ctx context.Context, userID int64) error
    
    // 用户查询
    ListUsers(ctx context.Context, opts *ListUsersOptions) (*repository.PaginationResult[models.User], error)
    SearchUsers(ctx context.Context, query string, opts *SearchUsersOptions) ([]*models.User, error)
    GetUserStats(ctx context.Context, userID int64) (*UserStats, error)
}

// AuthService 认证服务接口
type AuthService interface {
    // Token管理
    GenerateTokens(ctx context.Context, user *models.User) (*TokenResponse, error)
    ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
    RevokeToken(ctx context.Context, token string) error
    
    // 会话管理
    CreateSession(ctx context.Context, userID int64, deviceInfo *DeviceInfo) (*Session, error)
    GetSession(ctx context.Context, sessionID string) (*Session, error)
    UpdateSession(ctx context.Context, sessionID string, req *UpdateSessionRequest) error
    DeleteSession(ctx context.Context, sessionID string) error
    DeleteAllSessions(ctx context.Context, userID int64) error
    
    // 权限验证
    HasPermission(ctx context.Context, userID int64, permission string) (bool, error)
    GetUserPermissions(ctx context.Context, userID int64) ([]string, error)
    GrantPermission(ctx context.Context, userID int64, permission string) error
    RevokePermission(ctx context.Context, userID int64, permission string) error
}
```

#### 5.1.3.2.2 请求和响应结构
```go
// internal/services/user/types.go
package user

import (
    "time"
    
    "github.com/yourname/movieinfo/internal/models"
)

// RegisterRequest 注册请求
type RegisterRequest struct {
    Email           string `json:"email" validate:"required,email"`
    Username        string `json:"username" validate:"required,min=3,max=20,alphanum"`
    Password        string `json:"password" validate:"required,min=8,max=128"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
    Nickname        string `json:"nickname,omitempty" validate:"omitempty,min=1,max=50"`
    AcceptTerms     bool   `json:"accept_terms" validate:"required,eq=true"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
    User         *models.User  `json:"user"`
    AccessToken  string        `json:"access_token"`
    RefreshToken string        `json:"refresh_token"`
    ExpiresIn    int64         `json:"expires_in"`
    TokenType    string        `json:"token_type"`
    Message      string        `json:"message"`
}

// LoginRequest 登录请求
type LoginRequest struct {
    EmailOrUsername string `json:"email_or_username" validate:"required"`
    Password        string `json:"password" validate:"required"`
    RememberMe      bool   `json:"remember_me"`
    DeviceInfo      *DeviceInfo `json:"device_info,omitempty"`
}

// LoginResponse 登录响应
type LoginResponse struct {
    User         *models.User  `json:"user"`
    AccessToken  string        `json:"access_token"`
    RefreshToken string        `json:"refresh_token"`
    ExpiresIn    int64         `json:"expires_in"`
    TokenType    string        `json:"token_type"`
    SessionID    string        `json:"session_id"`
}

// TokenResponse Token响应
type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
    TokenType    string `json:"token_type"`
}

// TokenClaims Token声明
type TokenClaims struct {
    UserID    int64  `json:"user_id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    SessionID string `json:"session_id"`
    IssuedAt  int64  `json:"iat"`
    ExpiresAt int64  `json:"exp"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
    Nickname *string           `json:"nickname,omitempty" validate:"omitempty,min=1,max=50"`
    Bio      *string           `json:"bio,omitempty" validate:"omitempty,max=500"`
    Gender   *models.Gender    `json:"gender,omitempty"`
    Birthday *time.Time        `json:"birthday,omitempty"`
    Location *string           `json:"location,omitempty" validate:"omitempty,max=100"`
    Website  *string           `json:"website,omitempty" validate:"omitempty,url,max=255"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
    CurrentPassword string `json:"current_password" validate:"required"`
    NewPassword     string `json:"new_password" validate:"required,min=8,max=128"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
    Token           string `json:"token" validate:"required"`
    NewPassword     string `json:"new_password" validate:"required,min=8,max=128"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// DeviceInfo 设备信息
type DeviceInfo struct {
    UserAgent string `json:"user_agent"`
    IP        string `json:"ip"`
    Platform  string `json:"platform"`
    Browser   string `json:"browser"`
    OS        string `json:"os"`
}

// Session 会话信息
type Session struct {
    ID         string     `json:"id"`
    UserID     int64      `json:"user_id"`
    DeviceInfo DeviceInfo `json:"device_info"`
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
    ExpiresAt  time.Time  `json:"expires_at"`
    IsActive   bool       `json:"is_active"`
}

// UserStats 用户统计
type UserStats struct {
    MoviesWatched   int `json:"movies_watched"`
    ReviewsCount    int `json:"reviews_count"`
    RatingsCount    int `json:"ratings_count"`
    FavoritesCount  int `json:"favorites_count"`
    FollowersCount  int `json:"followers_count"`
    FollowingCount  int `json:"following_count"`
    JoinedDaysAgo   int `json:"joined_days_ago"`
    LastActiveAgo   int `json:"last_active_ago"`
}

// ListUsersOptions 用户列表选项
type ListUsersOptions struct {
    Page     int                    `json:"page" validate:"min=1"`
    PageSize int                    `json:"page_size" validate:"min=1,max=100"`
    OrderBy  string                 `json:"order_by" validate:"omitempty,oneof=id username email created_at"`
    Order    string                 `json:"order" validate:"omitempty,oneof=asc desc"`
    Status   *models.UserStatus     `json:"status,omitempty"`
    Filters  map[string]interface{} `json:"filters,omitempty"`
}

// SearchUsersOptions 用户搜索选项
type SearchUsersOptions struct {
    Limit   int      `json:"limit" validate:"min=1,max=50"`
    Fields  []string `json:"fields,omitempty"` // username, email, nickname
    Filters map[string]interface{} `json:"filters,omitempty"`
}
```

### 5.1.3.3 **用户服务实现**

#### 5.1.3.3.1 核心用户服务实现
```go
// internal/services/user/service.go
package user

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "time"
    
    "golang.org/x/crypto/bcrypt"
    
    "github.com/yourname/movieinfo/internal/models"
    "github.com/yourname/movieinfo/internal/repository"
    "github.com/yourname/movieinfo/pkg/cache"
    "github.com/yourname/movieinfo/pkg/errors"
    "github.com/yourname/movieinfo/pkg/logger"
    "github.com/yourname/movieinfo/pkg/validator"
)

// ServiceImpl 用户服务实现
type ServiceImpl struct {
    userRepo    repository.UserRepository
    authService AuthService
    cache       cache.Cache
    validator   validator.Validator
    logger      logger.Logger
    
    // 配置
    config *Config
}

// Config 用户服务配置
type Config struct {
    PasswordCost        int           `yaml:"password_cost"`
    TokenExpiry         time.Duration `yaml:"token_expiry"`
    RefreshTokenExpiry  time.Duration `yaml:"refresh_token_expiry"`
    VerificationExpiry  time.Duration `yaml:"verification_expiry"`
    MaxLoginAttempts    int           `yaml:"max_login_attempts"`
    LoginAttemptWindow  time.Duration `yaml:"login_attempt_window"`
    RequireEmailVerify  bool          `yaml:"require_email_verify"`
}

// NewService 创建用户服务
func NewService(
    userRepo repository.UserRepository,
    authService AuthService,
    cache cache.Cache,
    validator validator.Validator,
    config *Config,
) UserService {
    return &ServiceImpl{
        userRepo:    userRepo,
        authService: authService,
        cache:       cache,
        validator:   validator,
        logger:      logger.GetGlobalLogger(),
        config:      config,
    }
}

// Register 用户注册
func (s *ServiceImpl) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
    // 验证请求参数
    if err := s.validator.Validate(req); err != nil {
        return nil, errors.ValidationFailed(err.Error())
    }
    
    s.logger.WithContext(ctx).Info("User registration started",
        logger.String("email", req.Email),
        logger.String("username", req.Username),
    )
    
    // 检查邮箱是否已存在
    existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err != nil && !errors.IsNotFound(err) {
        return nil, errors.Wrap(err, errors.CodeDatabaseError, "检查邮箱失败")
    }
    if existingUser != nil {
        return nil, errors.UserExists().WithDetails("邮箱已被注册")
    }
    
    // 检查用户名是否已存在
    existingUser, err = s.userRepo.GetByUsername(ctx, req.Username)
    if err != nil && !errors.IsNotFound(err) {
        return nil, errors.Wrap(err, errors.CodeDatabaseError, "检查用户名失败")
    }
    if existingUser != nil {
        return nil, errors.UserExists().WithDetails("用户名已被使用")
    }
    
    // 加密密码
    hashedPassword, err := s.hashPassword(req.Password)
    if err != nil {
        return nil, errors.InternalError("密码加密失败")
    }
    
    // 创建用户
    user := &models.User{
        Email:    req.Email,
        Username: req.Username,
        Password: hashedPassword,
        Nickname: req.Nickname,
        Status:   models.UserStatusActive,
    }
    
    if s.config.RequireEmailVerify {
        user.Status = models.UserStatusInactive
    }
    
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, errors.Wrap(err, errors.CodeDatabaseError, "创建用户失败")
    }
    
    // 生成Token
    tokens, err := s.authService.GenerateTokens(ctx, user)
    if err != nil {
        return nil, errors.Wrap(err, errors.CodeInternalError, "生成Token失败")
    }
    
    // 发送验证邮件
    if s.config.RequireEmailVerify {
        if err := s.SendVerificationEmail(ctx, user.ID); err != nil {
            s.logger.WithContext(ctx).Error("Failed to send verification email",
                logger.Int64("user_id", user.ID),
                logger.Error(err),
            )
            // 不阻断注册流程
        }
    }
    
    // 清除密码字段
    user.Password = ""
    
    s.logger.WithContext(ctx).Info("User registration completed",
        logger.Int64("user_id", user.ID),
        logger.String("email", user.Email),
    )
    
    message := "注册成功"
    if s.config.RequireEmailVerify {
        message = "注册成功，请查收验证邮件"
    }
    
    return &RegisterResponse{
        User:         user,
        AccessToken:  tokens.AccessToken,
        RefreshToken: tokens.RefreshToken,
        ExpiresIn:    tokens.ExpiresIn,
        TokenType:    tokens.TokenType,
        Message:      message,
    }, nil
}

// Login 用户登录
func (s *ServiceImpl) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
    // 验证请求参数
    if err := s.validator.Validate(req); err != nil {
        return nil, errors.ValidationFailed(err.Error())
    }
    
    s.logger.WithContext(ctx).Info("User login attempt",
        logger.String("email_or_username", req.EmailOrUsername),
    )
    
    // 检查登录尝试次数
    if err := s.checkLoginAttempts(ctx, req.EmailOrUsername); err != nil {
        return nil, err
    }
    
    // 获取用户
    user, err := s.userRepo.GetByEmailOrUsername(ctx, req.EmailOrUsername)
    if err != nil {
        if errors.IsNotFound(err) {
            s.recordFailedLogin(ctx, req.EmailOrUsername)
            return nil, errors.Unauthorized("用户名或密码错误")
        }
        return nil, errors.Wrap(err, errors.CodeDatabaseError, "查询用户失败")
    }
    
    // 验证密码
    if !s.verifyPassword(req.Password, user.Password) {
        s.recordFailedLogin(ctx, req.EmailOrUsername)
        return nil, errors.Unauthorized("用户名或密码错误")
    }
    
    // 检查用户状态
    if !user.CanLogin() {
        return nil, s.getUserStatusError(user)
    }
    
    // 生成Token
    tokens, err := s.authService.GenerateTokens(ctx, user)
    if err != nil {
        return nil, errors.Wrap(err, errors.CodeInternalError, "生成Token失败")
    }
    
    // 创建会话
    session, err := s.authService.CreateSession(ctx, user.ID, req.DeviceInfo)
    if err != nil {
        s.logger.WithContext(ctx).Error("Failed to create session",
            logger.Int64("user_id", user.ID),
            logger.Error(err),
        )
        // 不阻断登录流程
    }
    
    // 更新最后登录信息
    if req.DeviceInfo != nil {
        user.UpdateLastLogin(req.DeviceInfo.IP)
        if err := s.userRepo.UpdateLastLogin(ctx, user.ID, req.DeviceInfo.IP); err != nil {
            s.logger.WithContext(ctx).Error("Failed to update last login",
                logger.Int64("user_id", user.ID),
                logger.Error(err),
            )
        }
    }
    
    // 清除登录失败记录
    s.clearFailedLogins(ctx, req.EmailOrUsername)
    
    // 清除密码字段
    user.Password = ""
    
    s.logger.WithContext(ctx).Info("User login successful",
        logger.Int64("user_id", user.ID),
        logger.String("email", user.Email),
    )
    
    sessionID := ""
    if session != nil {
        sessionID = session.ID
    }
    
    return &LoginResponse{
        User:         user,
        AccessToken:  tokens.AccessToken,
        RefreshToken: tokens.RefreshToken,
        ExpiresIn:    tokens.ExpiresIn,
        TokenType:    tokens.TokenType,
        SessionID:    sessionID,
    }, nil
}

// GetUser 获取用户信息
func (s *ServiceImpl) GetUser(ctx context.Context, userID int64) (*models.User, error) {
    // 先从缓存获取
    cacheKey := fmt.Sprintf("user:%d", userID)
    if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
        if user, ok := cached.(*models.User); ok {
            return user, nil
        }
    }
    
    // 从数据库获取
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        if errors.IsNotFound(err) {
            return nil, errors.UserNotFound()
        }
        return nil, errors.Wrap(err, errors.CodeDatabaseError, "查询用户失败")
    }
    
    // 清除密码字段
    user.Password = ""
    
    // 缓存用户信息
    s.cache.Set(ctx, cacheKey, user, time.Hour)
    
    return user, nil
}

// hashPassword 加密密码
func (s *ServiceImpl) hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.config.PasswordCost)
    return string(bytes), err
}

// verifyPassword 验证密码
func (s *ServiceImpl) verifyPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// checkLoginAttempts 检查登录尝试次数
func (s *ServiceImpl) checkLoginAttempts(ctx context.Context, identifier string) error {
    cacheKey := fmt.Sprintf("login_attempts:%s", identifier)
    attempts, _ := s.cache.Get(ctx, cacheKey)
    
    if count, ok := attempts.(int); ok && count >= s.config.MaxLoginAttempts {
        return errors.TooManyRequests("登录尝试次数过多，请稍后再试")
    }
    
    return nil
}

// recordFailedLogin 记录登录失败
func (s *ServiceImpl) recordFailedLogin(ctx context.Context, identifier string) {
    cacheKey := fmt.Sprintf("login_attempts:%s", identifier)
    attempts, _ := s.cache.Get(ctx, cacheKey)
    
    count := 1
    if existingCount, ok := attempts.(int); ok {
        count = existingCount + 1
    }
    
    s.cache.Set(ctx, cacheKey, count, s.config.LoginAttemptWindow)
}

// clearFailedLogins 清除登录失败记录
func (s *ServiceImpl) clearFailedLogins(ctx context.Context, identifier string) {
    cacheKey := fmt.Sprintf("login_attempts:%s", identifier)
    s.cache.Delete(ctx, cacheKey)
}

// getUserStatusError 获取用户状态错误
func (s *ServiceImpl) getUserStatusError(user *models.User) error {
    switch user.Status {
    case models.UserStatusInactive:
        if !user.EmailVerified {
            return errors.Unauthorized("请先验证邮箱")
        }
        return errors.Unauthorized("账户未激活")
    case models.UserStatusBanned:
        return errors.Forbidden("账户已被禁用")
    case models.UserStatusDeleted:
        return errors.NotFound("账户不存在")
    default:
        return errors.Unauthorized("账户状态异常")
    }
}
```

## 5.1.4 总结

用户服务实现为MovieInfo项目提供了完整的用户管理解决方案。通过分层架构、安全设计和性能优化，我们建立了一个功能完整、安全可靠的用户服务系统。

**关键设计要点**：
1. **安全性**：密码加密、登录限制、会话管理
2. **业务完整性**：注册、登录、权限、状态管理
3. **性能优化**：缓存策略、查询优化、批量操作
4. **错误处理**：统一的错误处理和用户友好提示
5. **扩展性**：支持第三方登录、社交功能扩展

**服务优势**：
- **安全可靠**：多层安全防护和验证机制
- **功能完整**：覆盖用户管理的各个方面
- **性能优化**：缓存和查询优化提升性能
- **易于扩展**：清晰的接口设计便于功能扩展

**下一步**：基于用户服务基础，我们将实现电影服务，处理电影数据管理、搜索、推荐等核心业务功能。
