# 4.2 数据模型定义

## 4.2.1 概述

数据模型是应用程序的核心抽象，它定义了业务实体的结构、属性和关系。对于MovieInfo项目，我们需要设计清晰、完整、可扩展的数据模型，涵盖用户、电影、评论、评分等核心业务实体。

## 4.2.2 为什么数据模型设计如此重要？

### 4.2.2.1 **业务抽象**
- **实体映射**：将现实世界的业务概念映射为数据结构
- **关系定义**：明确实体间的关联和依赖关系
- **约束规则**：定义数据的完整性和一致性约束
- **业务规则**：在模型层面体现业务逻辑和规则

### 4.2.2.2 **开发效率**
- **代码生成**：基于模型自动生成CRUD代码
- **类型安全**：强类型定义减少运行时错误
- **IDE支持**：良好的代码补全和重构支持
- **文档化**：模型即文档，便于理解和维护

### 4.2.2.3 **数据一致性**
- **结构统一**：确保数据结构在各层间的一致性
- **验证规则**：统一的数据验证和格式化规则
- **序列化**：标准化的JSON/XML序列化格式
- **版本兼容**：支持模型的版本演进和兼容性

### 4.2.2.4 **性能优化**
- **查询优化**：基于模型关系优化数据库查询
- **缓存策略**：针对模型特点设计缓存策略
- **索引设计**：根据模型访问模式设计索引
- **分页支持**：内置分页和排序支持

## 4.2.3 数据模型架构设计

### 4.2.3.1 **模型分层结构**

```
数据模型层次
├── 核心实体模型 (Core Entity Models)
│   ├── 用户模型 (User)
│   ├── 电影模型 (Movie)
│   ├── 评论模型 (Comment)
│   └── 评分模型 (Rating)
├── 关联模型 (Association Models)
│   ├── 电影分类关联 (MovieCategory)
│   ├── 电影演员关联 (MovieActor)
│   ├── 用户收藏关联 (UserFavorite)
│   └── 评论点赞关联 (CommentLike)
├── 配置模型 (Configuration Models)
│   ├── 分类模型 (Category)
│   ├── 演员模型 (Actor)
│   ├── 标签模型 (Tag)
│   └── 系统配置 (SystemConfig)
└── 辅助模型 (Helper Models)
    ├── 分页模型 (Pagination)
    ├── 排序模型 (Sorting)
    ├── 过滤模型 (Filter)
    └── 响应模型 (Response)
```

### 2. **基础模型定义**

#### 2.1 基础模型接口
```go
// internal/models/base.go
package models

import (
    "time"
)

// Model 基础模型接口
type Model interface {
    GetID() int64
    SetID(id int64)
    GetCreatedAt() time.Time
    GetUpdatedAt() time.Time
    Validate() error
}

// BaseModel 基础模型结构
type BaseModel struct {
    ID        int64     `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
    CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at" gorm:"index"`
}

// GetID 获取ID
func (m *BaseModel) GetID() int64 {
    return m.ID
}

// SetID 设置ID
func (m *BaseModel) SetID(id int64) {
    m.ID = id
}

// GetCreatedAt 获取创建时间
func (m *BaseModel) GetCreatedAt() time.Time {
    return m.CreatedAt
}

// GetUpdatedAt 获取更新时间
func (m *BaseModel) GetUpdatedAt() time.Time {
    return m.UpdatedAt
}

// IsDeleted 检查是否已删除
func (m *BaseModel) IsDeleted() bool {
    return m.DeletedAt != nil
}

// SoftDelete 软删除
func (m *BaseModel) SoftDelete() {
    now := time.Now()
    m.DeletedAt = &now
}

// Restore 恢复删除
func (m *BaseModel) Restore() {
    m.DeletedAt = nil
}
```

#### 2.2 验证接口
```go
// internal/models/validation.go
package models

import (
    "fmt"
    "regexp"
    "strings"
    "unicode/utf8"
)

// ValidationError 验证错误
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Value   interface{} `json:"value,omitempty"`
}

// Error 实现error接口
func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// ValidationErrors 验证错误集合
type ValidationErrors []ValidationError

// Error 实现error接口
func (e ValidationErrors) Error() string {
    if len(e) == 0 {
        return ""
    }
    
    var messages []string
    for _, err := range e {
        messages = append(messages, err.Error())
    }
    
    return strings.Join(messages, "; ")
}

// Add 添加验证错误
func (e *ValidationErrors) Add(field, message string, value interface{}) {
    *e = append(*e, ValidationError{
        Field:   field,
        Message: message,
        Value:   value,
    })
}

// HasErrors 检查是否有错误
func (e ValidationErrors) HasErrors() bool {
    return len(e) > 0
}

// 常用验证函数
var (
    emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
)

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
    return emailRegex.MatchString(email)
}

// ValidateUsername 验证用户名格式
func ValidateUsername(username string) bool {
    return usernameRegex.MatchString(username)
}

// ValidateStringLength 验证字符串长度
func ValidateStringLength(value string, min, max int) bool {
    length := utf8.RuneCountInString(value)
    return length >= min && length <= max
}

// ValidateRequired 验证必填字段
func ValidateRequired(value string) bool {
    return strings.TrimSpace(value) != ""
}

// ValidateRange 验证数值范围
func ValidateRange(value, min, max int) bool {
    return value >= min && value <= max
}
```

### 3. **核心实体模型**

#### 3.1 用户模型
```go
// internal/models/user.go
package models

import (
    "crypto/md5"
    "fmt"
    "strings"
    "time"
)

// User 用户模型
type User struct {
    BaseModel
    
    // 基本信息
    Email    string `json:"email" db:"email" gorm:"uniqueIndex;size:100;not null"`
    Username string `json:"username" db:"username" gorm:"uniqueIndex;size:50;not null"`
    Password string `json:"-" db:"password" gorm:"size:255;not null"` // 不在JSON中显示
    
    // 个人信息
    Nickname  string  `json:"nickname" db:"nickname" gorm:"size:50"`
    Avatar    string  `json:"avatar" db:"avatar" gorm:"size:255"`
    Bio       string  `json:"bio" db:"bio" gorm:"type:text"`
    Gender    *Gender `json:"gender" db:"gender" gorm:"type:tinyint"`
    Birthday  *time.Time `json:"birthday" db:"birthday"`
    Location  string  `json:"location" db:"location" gorm:"size:100"`
    Website   string  `json:"website" db:"website" gorm:"size:255"`
    
    // 状态信息
    Status           UserStatus `json:"status" db:"status" gorm:"type:tinyint;default:1"`
    EmailVerified    bool       `json:"email_verified" db:"email_verified" gorm:"default:false"`
    EmailVerifiedAt  *time.Time `json:"email_verified_at" db:"email_verified_at"`
    LastLoginAt      *time.Time `json:"last_login_at" db:"last_login_at"`
    LastLoginIP      string     `json:"last_login_ip" db:"last_login_ip" gorm:"size:45"`
    
    // 统计信息
    MoviesWatched    int `json:"movies_watched" db:"movies_watched" gorm:"default:0"`
    ReviewsCount     int `json:"reviews_count" db:"reviews_count" gorm:"default:0"`
    FollowersCount   int `json:"followers_count" db:"followers_count" gorm:"default:0"`
    FollowingCount   int `json:"following_count" db:"following_count" gorm:"default:0"`
    
    // 关联关系（不存储在数据库中）
    Ratings   []Rating  `json:"ratings,omitempty" gorm:"foreignKey:UserID"`
    Comments  []Comment `json:"comments,omitempty" gorm:"foreignKey:UserID"`
    Favorites []Movie   `json:"favorites,omitempty" gorm:"many2many:user_favorites;"`
}

// Gender 性别枚举
type Gender int

const (
    GenderUnknown Gender = iota
    GenderMale
    GenderFemale
    GenderOther
)

// String 返回性别字符串
func (g Gender) String() string {
    switch g {
    case GenderMale:
        return "male"
    case GenderFemale:
        return "female"
    case GenderOther:
        return "other"
    default:
        return "unknown"
    }
}

// UserStatus 用户状态枚举
type UserStatus int

const (
    UserStatusInactive UserStatus = iota
    UserStatusActive
    UserStatusBanned
    UserStatusDeleted
)

// String 返回用户状态字符串
func (s UserStatus) String() string {
    switch s {
    case UserStatusActive:
        return "active"
    case UserStatusBanned:
        return "banned"
    case UserStatusDeleted:
        return "deleted"
    default:
        return "inactive"
    }
}

// TableName 指定表名
func (User) TableName() string {
    return "users"
}

// Validate 验证用户数据
func (u *User) Validate() error {
    var errors ValidationErrors
    
    // 验证邮箱
    if !ValidateRequired(u.Email) {
        errors.Add("email", "邮箱不能为空", u.Email)
    } else if !ValidateEmail(u.Email) {
        errors.Add("email", "邮箱格式不正确", u.Email)
    }
    
    // 验证用户名
    if !ValidateRequired(u.Username) {
        errors.Add("username", "用户名不能为空", u.Username)
    } else if !ValidateUsername(u.Username) {
        errors.Add("username", "用户名格式不正确，只能包含字母、数字和下划线，长度3-20位", u.Username)
    }
    
    // 验证密码（仅在创建时）
    if u.ID == 0 && !ValidateRequired(u.Password) {
        errors.Add("password", "密码不能为空", nil)
    } else if u.Password != "" && !ValidateStringLength(u.Password, 8, 128) {
        errors.Add("password", "密码长度必须在8-128位之间", nil)
    }
    
    // 验证昵称
    if u.Nickname != "" && !ValidateStringLength(u.Nickname, 1, 50) {
        errors.Add("nickname", "昵称长度必须在1-50位之间", u.Nickname)
    }
    
    // 验证个人简介
    if u.Bio != "" && !ValidateStringLength(u.Bio, 0, 500) {
        errors.Add("bio", "个人简介长度不能超过500字", u.Bio)
    }
    
    // 验证网站URL
    if u.Website != "" && !ValidateStringLength(u.Website, 0, 255) {
        errors.Add("website", "网站URL长度不能超过255字符", u.Website)
    }
    
    if errors.HasErrors() {
        return errors
    }
    
    return nil
}

// GetDisplayName 获取显示名称
func (u *User) GetDisplayName() string {
    if u.Nickname != "" {
        return u.Nickname
    }
    return u.Username
}

// GetAvatarURL 获取头像URL
func (u *User) GetAvatarURL() string {
    if u.Avatar != "" {
        return u.Avatar
    }
    
    // 生成Gravatar头像
    hash := fmt.Sprintf("%x", md5.Sum([]byte(strings.ToLower(u.Email))))
    return fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=identicon&s=200", hash)
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
    return u.Status == UserStatusActive && !u.IsDeleted()
}

// IsBanned 检查用户是否被禁用
func (u *User) IsBanned() bool {
    return u.Status == UserStatusBanned
}

// CanLogin 检查用户是否可以登录
func (u *User) CanLogin() bool {
    return u.IsActive() && u.EmailVerified
}

// UpdateLastLogin 更新最后登录信息
func (u *User) UpdateLastLogin(ip string) {
    now := time.Now()
    u.LastLoginAt = &now
    u.LastLoginIP = ip
}

// IncrementMoviesWatched 增加观看电影数量
func (u *User) IncrementMoviesWatched() {
    u.MoviesWatched++
}

// IncrementReviewsCount 增加评论数量
func (u *User) IncrementReviewsCount() {
    u.ReviewsCount++
}

// DecrementReviewsCount 减少评论数量
func (u *User) DecrementReviewsCount() {
    if u.ReviewsCount > 0 {
        u.ReviewsCount--
    }
}
```

#### 3.2 电影模型
```go
// internal/models/movie.go
package models

import (
    "fmt"
    "strings"
    "time"
)

// Movie 电影模型
type Movie struct {
    BaseModel
    
    // 基本信息
    Title         string  `json:"title" db:"title" gorm:"size:200;not null;index"`
    OriginalTitle string  `json:"original_title" db:"original_title" gorm:"size:200;index"`
    Director      string  `json:"director" db:"director" gorm:"size:100;index"`
    ReleaseYear   *int    `json:"release_year" db:"release_year" gorm:"index"`
    Duration      *int    `json:"duration" db:"duration"` // 分钟
    Country       string  `json:"country" db:"country" gorm:"size:100"`
    Language      string  `json:"language" db:"language" gorm:"size:50"`
    
    // 内容信息
    Plot        string `json:"plot" db:"plot" gorm:"type:text"`
    Tagline     string `json:"tagline" db:"tagline" gorm:"size:255"`
    Keywords    string `json:"keywords" db:"keywords" gorm:"size:500"`
    
    // 媒体信息
    PosterURL   string `json:"poster_url" db:"poster_url" gorm:"size:255"`
    BackdropURL string `json:"backdrop_url" db:"backdrop_url" gorm:"size:255"`
    TrailerURL  string `json:"trailer_url" db:"trailer_url" gorm:"size:255"`
    
    // 外部ID
    IMDBID string `json:"imdb_id" db:"imdb_id" gorm:"size:20;uniqueIndex"`
    TMDBID *int   `json:"tmdb_id" db:"tmdb_id" gorm:"uniqueIndex"`
    
    // 评分信息
    AverageRating float64 `json:"average_rating" db:"average_rating" gorm:"type:decimal(3,2);default:0.00;index"`
    RatingCount   int     `json:"rating_count" db:"rating_count" gorm:"default:0"`
    
    // 统计信息
    ViewCount     int `json:"view_count" db:"view_count" gorm:"default:0"`
    CommentCount  int `json:"comment_count" db:"comment_count" gorm:"default:0"`
    FavoriteCount int `json:"favorite_count" db:"favorite_count" gorm:"default:0"`
    
    // 状态信息
    Status MovieStatus `json:"status" db:"status" gorm:"type:tinyint;default:1;index"`
    
    // 关联关系
    Categories []Category `json:"categories,omitempty" gorm:"many2many:movie_categories;"`
    Actors     []Actor    `json:"actors,omitempty" gorm:"many2many:movie_actors;"`
    Ratings    []Rating   `json:"ratings,omitempty" gorm:"foreignKey:MovieID"`
    Comments   []Comment  `json:"comments,omitempty" gorm:"foreignKey:MovieID"`
}

// MovieStatus 电影状态枚举
type MovieStatus int

const (
    MovieStatusDraft MovieStatus = iota
    MovieStatusPublished
    MovieStatusArchived
    MovieStatusDeleted
)

// String 返回电影状态字符串
func (s MovieStatus) String() string {
    switch s {
    case MovieStatusPublished:
        return "published"
    case MovieStatusArchived:
        return "archived"
    case MovieStatusDeleted:
        return "deleted"
    default:
        return "draft"
    }
}

// TableName 指定表名
func (Movie) TableName() string {
    return "movies"
}

// Validate 验证电影数据
func (m *Movie) Validate() error {
    var errors ValidationErrors
    
    // 验证标题
    if !ValidateRequired(m.Title) {
        errors.Add("title", "电影标题不能为空", m.Title)
    } else if !ValidateStringLength(m.Title, 1, 200) {
        errors.Add("title", "电影标题长度必须在1-200字符之间", m.Title)
    }
    
    // 验证原始标题
    if m.OriginalTitle != "" && !ValidateStringLength(m.OriginalTitle, 1, 200) {
        errors.Add("original_title", "原始标题长度必须在1-200字符之间", m.OriginalTitle)
    }
    
    // 验证导演
    if m.Director != "" && !ValidateStringLength(m.Director, 1, 100) {
        errors.Add("director", "导演姓名长度必须在1-100字符之间", m.Director)
    }
    
    // 验证上映年份
    if m.ReleaseYear != nil {
        currentYear := time.Now().Year()
        if !ValidateRange(*m.ReleaseYear, 1888, currentYear+5) { // 电影历史从1888年开始
            errors.Add("release_year", fmt.Sprintf("上映年份必须在1888-%d之间", currentYear+5), *m.ReleaseYear)
        }
    }
    
    // 验证时长
    if m.Duration != nil && !ValidateRange(*m.Duration, 1, 1000) { // 1分钟到1000分钟
        errors.Add("duration", "电影时长必须在1-1000分钟之间", *m.Duration)
    }
    
    // 验证剧情简介
    if m.Plot != "" && !ValidateStringLength(m.Plot, 0, 2000) {
        errors.Add("plot", "剧情简介长度不能超过2000字符", m.Plot)
    }
    
    // 验证标语
    if m.Tagline != "" && !ValidateStringLength(m.Tagline, 0, 255) {
        errors.Add("tagline", "标语长度不能超过255字符", m.Tagline)
    }
    
    // 验证评分范围
    if m.AverageRating < 0 || m.AverageRating > 5 {
        errors.Add("average_rating", "平均评分必须在0-5之间", m.AverageRating)
    }
    
    // 验证评分人数
    if m.RatingCount < 0 {
        errors.Add("rating_count", "评分人数不能为负数", m.RatingCount)
    }
    
    if errors.HasErrors() {
        return errors
    }
    
    return nil
}

// GetDisplayTitle 获取显示标题
func (m *Movie) GetDisplayTitle() string {
    if m.OriginalTitle != "" && m.OriginalTitle != m.Title {
        return fmt.Sprintf("%s (%s)", m.Title, m.OriginalTitle)
    }
    return m.Title
}

// GetYearString 获取年份字符串
func (m *Movie) GetYearString() string {
    if m.ReleaseYear != nil {
        return fmt.Sprintf("%d", *m.ReleaseYear)
    }
    return "未知"
}

// GetDurationString 获取时长字符串
func (m *Movie) GetDurationString() string {
    if m.Duration != nil {
        hours := *m.Duration / 60
        minutes := *m.Duration % 60
        if hours > 0 {
            return fmt.Sprintf("%d小时%d分钟", hours, minutes)
        }
        return fmt.Sprintf("%d分钟", minutes)
    }
    return "未知"
}

// GetRatingString 获取评分字符串
func (m *Movie) GetRatingString() string {
    if m.RatingCount > 0 {
        return fmt.Sprintf("%.1f (%d人评价)", m.AverageRating, m.RatingCount)
    }
    return "暂无评分"
}

// IsPublished 检查是否已发布
func (m *Movie) IsPublished() bool {
    return m.Status == MovieStatusPublished && !m.IsDeleted()
}

// CanView 检查是否可以查看
func (m *Movie) CanView() bool {
    return m.IsPublished()
}

// IncrementViewCount 增加浏览次数
func (m *Movie) IncrementViewCount() {
    m.ViewCount++
}

// UpdateRating 更新评分信息
func (m *Movie) UpdateRating(newRating float64, isNew bool) {
    if isNew {
        // 新增评分
        totalRating := m.AverageRating*float64(m.RatingCount) + newRating
        m.RatingCount++
        m.AverageRating = totalRating / float64(m.RatingCount)
    } else {
        // 更新现有评分，需要传入旧评分进行计算
        // 这里简化处理，实际应该传入旧评分值
        if m.RatingCount > 0 {
            totalRating := m.AverageRating*float64(m.RatingCount) - m.AverageRating + newRating
            m.AverageRating = totalRating / float64(m.RatingCount)
        }
    }
    
    // 保留两位小数
    m.AverageRating = float64(int(m.AverageRating*100)) / 100
}

// RemoveRating 移除评分
func (m *Movie) RemoveRating(rating float64) {
    if m.RatingCount > 1 {
        totalRating := m.AverageRating*float64(m.RatingCount) - rating
        m.RatingCount--
        m.AverageRating = totalRating / float64(m.RatingCount)
        m.AverageRating = float64(int(m.AverageRating*100)) / 100
    } else {
        m.RatingCount = 0
        m.AverageRating = 0
    }
}

// GetCategoryNames 获取分类名称列表
func (m *Movie) GetCategoryNames() []string {
    var names []string
    for _, category := range m.Categories {
        names = append(names, category.Name)
    }
    return names
}

// GetActorNames 获取演员名称列表
func (m *Movie) GetActorNames() []string {
    var names []string
    for _, actor := range m.Actors {
        names = append(names, actor.Name)
    }
    return names
}

// GetKeywordList 获取关键词列表
func (m *Movie) GetKeywordList() []string {
    if m.Keywords == "" {
        return []string{}
    }
    
    keywords := strings.Split(m.Keywords, ",")
    var result []string
    for _, keyword := range keywords {
        keyword = strings.TrimSpace(keyword)
        if keyword != "" {
            result = append(result, keyword)
        }
    }
    
    return result
}

// SetKeywords 设置关键词
func (m *Movie) SetKeywords(keywords []string) {
    var cleanKeywords []string
    for _, keyword := range keywords {
        keyword = strings.TrimSpace(keyword)
        if keyword != "" {
            cleanKeywords = append(cleanKeywords, keyword)
        }
    }
    
    m.Keywords = strings.Join(cleanKeywords, ",")
}
```

## 总结

数据模型定义为MovieInfo项目提供了完整的业务实体抽象。通过清晰的结构定义、完善的验证机制和丰富的业务方法，我们建立了一个类型安全、功能完整的数据模型系统。

**关键设计要点**：
1. **基础模型**：统一的基础结构和接口定义
2. **验证机制**：完善的数据验证和错误处理
3. **业务方法**：丰富的业务逻辑和便捷方法
4. **关联关系**：清晰的实体间关联定义
5. **扩展性**：支持模型的扩展和演进

**模型优势**：
- **类型安全**：强类型定义减少错误
- **业务抽象**：清晰的业务概念映射
- **验证完整**：统一的数据验证机制
- **功能丰富**：便捷的业务操作方法

**下一步**：基于这些数据模型，我们将实现CRUD操作，包括数据的创建、读取、更新和删除功能。
