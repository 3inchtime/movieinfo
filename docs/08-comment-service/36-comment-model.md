# 5.3 评论服务实现

## 5.3.1 概述

评论服务是MovieInfo项目的重要互动模块，负责处理用户评论、评分、点赞、举报等社交功能。作为用户参与度的核心组件，评论服务需要实现内容审核、垃圾评论过滤、情感分析和社区管理等复杂功能。

## 5.3.2 为什么需要专业的评论服务？

### 5.3.2.1 **内容管理**
- **内容审核**：自动和人工审核机制防止不当内容
- **垃圾过滤**：识别和过滤垃圾评论和恶意内容
- **敏感词检测**：检测和处理敏感词汇和违规内容
- **质量控制**：评论质量评估和优质内容推荐

### 5.3.2.2 **用户体验**
- **实时互动**：支持实时评论和即时反馈
- **个性化展示**：基于用户偏好的评论排序和展示
- **社交功能**：点赞、回复、关注等社交互动
- **通知系统**：评论相关的实时通知和提醒

### 5.3.2.3 **数据分析**
- **情感分析**：分析评论情感倾向和用户态度
- **热点发现**：识别热门话题和讨论焦点
- **用户画像**：基于评论行为构建用户画像
- **内容洞察**：为内容推荐提供数据支持

### 5.3.2.4 **系统稳定性**
- **高并发处理**：支持大量用户同时评论和互动
- **数据一致性**：确保评论数据的一致性和完整性
- **性能优化**：评论列表的高效加载和分页
- **缓存策略**：热门评论和用户数据的智能缓存

## 5.3.3 评论服务架构设计

### 5.3.3.1 **服务层次结构**

```
评论服务架构
├── 服务接口层 (Service Interface)
│   ├── 评论服务接口 (CommentService)
│   ├── 评分服务接口 (RatingService)
│   ├── 点赞服务接口 (LikeService)
│   └── 举报服务接口 (ReportService)
├── 业务逻辑层 (Business Logic)
│   ├── 评论管理逻辑 (Comment Management)
│   ├── 内容审核逻辑 (Content Moderation)
│   ├── 社交互动逻辑 (Social Interaction)
│   └── 数据分析逻辑 (Analytics Logic)
├── 数据访问层 (Data Access)
│   ├── 评论Repository (CommentRepository)
│   ├── 评分Repository (RatingRepository)
│   ├── 点赞Repository (LikeRepository)
│   └── 举报Repository (ReportRepository)
├── 外部集成层 (External Integration)
│   ├── 内容审核API (Moderation API)
│   ├── 情感分析API (Sentiment API)
│   ├── 通知服务 (Notification Service)
│   └── 消息队列 (Message Queue)
├── 缓存层 (Cache Layer)
│   ├── 评论缓存 (Comment Cache)
│   ├── 评分缓存 (Rating Cache)
│   ├── 用户缓存 (User Cache)
│   └── 统计缓存 (Stats Cache)
└── 分析层 (Analytics Layer)
    ├── 情感分析 (Sentiment Analysis)
    ├── 热点检测 (Trending Detection)
    ├── 质量评估 (Quality Assessment)
    └── 用户行为分析 (Behavior Analysis)
```

### 5.3.3.2 **评论服务接口定义**

#### 5.3.3.2.1 核心评论服务接口
```go
// internal/services/comment/interface.go
package comment

import (
    "context"
    
    "github.com/yourname/movieinfo/internal/models"
    "github.com/yourname/movieinfo/internal/repository"
)

// CommentService 评论服务接口
type CommentService interface {
    // 评论基础操作
    CreateComment(ctx context.Context, req *CreateCommentRequest) (*models.Comment, error)
    GetComment(ctx context.Context, commentID int64) (*models.Comment, error)
    UpdateComment(ctx context.Context, commentID int64, req *UpdateCommentRequest) (*models.Comment, error)
    DeleteComment(ctx context.Context, commentID int64, userID int64) error
    
    // 评论列表和查询
    GetMovieComments(ctx context.Context, movieID int64, opts *CommentListOptions) (*repository.PaginationResult[models.Comment], error)
    GetUserComments(ctx context.Context, userID int64, opts *CommentListOptions) (*repository.PaginationResult[models.Comment], error)
    GetCommentReplies(ctx context.Context, parentID int64, opts *CommentListOptions) (*repository.PaginationResult[models.Comment], error)
    
    // 评论互动
    LikeComment(ctx context.Context, commentID int64, userID int64) error
    UnlikeComment(ctx context.Context, commentID int64, userID int64) error
    ReplyToComment(ctx context.Context, parentID int64, req *CreateCommentRequest) (*models.Comment, error)
    
    // 评论管理
    ReportComment(ctx context.Context, commentID int64, req *ReportCommentRequest) error
    ModerateComment(ctx context.Context, commentID int64, req *ModerateCommentRequest) error
    GetReportedComments(ctx context.Context, opts *CommentListOptions) (*repository.PaginationResult[models.Comment], error)
    
    // 评论统计
    GetCommentStats(ctx context.Context, movieID int64) (*CommentStats, error)
    GetUserCommentStats(ctx context.Context, userID int64) (*UserCommentStats, error)
    GetTrendingComments(ctx context.Context, timeWindow string, limit int) ([]*models.Comment, error)
}

// RatingService 评分服务接口
type RatingService interface {
    // 评分操作
    RateMovie(ctx context.Context, req *RateMovieRequest) (*models.Rating, error)
    UpdateRating(ctx context.Context, ratingID int64, req *UpdateRatingRequest) (*models.Rating, error)
    DeleteRating(ctx context.Context, ratingID int64, userID int64) error
    
    // 评分查询
    GetUserRating(ctx context.Context, movieID int64, userID int64) (*models.Rating, error)
    GetMovieRatings(ctx context.Context, movieID int64, opts *RatingListOptions) (*repository.PaginationResult[models.Rating], error)
    GetUserRatings(ctx context.Context, userID int64, opts *RatingListOptions) (*repository.PaginationResult[models.Rating], error)
    
    // 评分统计
    GetMovieRatingStats(ctx context.Context, movieID int64) (*RatingStats, error)
    GetRatingDistribution(ctx context.Context, movieID int64) (*RatingDistribution, error)
    GetUserRatingStats(ctx context.Context, userID int64) (*UserRatingStats, error)
    
    // 评分分析
    GetSimilarUsers(ctx context.Context, userID int64, limit int) ([]*SimilarUser, error)
    GetRatingTrends(ctx context.Context, movieID int64, timeWindow string) (*RatingTrends, error)
}

// LikeService 点赞服务接口
type LikeService interface {
    // 点赞操作
    LikeContent(ctx context.Context, req *LikeContentRequest) error
    UnlikeContent(ctx context.Context, req *UnlikeContentRequest) error
    
    // 点赞查询
    IsLiked(ctx context.Context, contentType string, contentID int64, userID int64) (bool, error)
    GetContentLikes(ctx context.Context, contentType string, contentID int64, opts *LikeListOptions) (*repository.PaginationResult[models.Like], error)
    GetUserLikes(ctx context.Context, userID int64, opts *LikeListOptions) (*repository.PaginationResult[models.Like], error)
    
    // 点赞统计
    GetLikeCount(ctx context.Context, contentType string, contentID int64) (int64, error)
    GetLikeStats(ctx context.Context, contentType string, contentID int64) (*LikeStats, error)
}
```

#### 5.3.3.2.2 请求和响应结构
```go
// internal/services/comment/types.go
package comment

import (
    "time"
    
    "github.com/yourname/movieinfo/internal/models"
)

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
    MovieID  int64  `json:"movie_id" validate:"required"`
    UserID   int64  `json:"user_id" validate:"required"`
    Content  string `json:"content" validate:"required,min=1,max=2000"`
    ParentID *int64 `json:"parent_id,omitempty"`
    Rating   *float64 `json:"rating,omitempty" validate:"omitempty,min=0,max=5"`
}

// UpdateCommentRequest 更新评论请求
type UpdateCommentRequest struct {
    Content *string  `json:"content,omitempty" validate:"omitempty,min=1,max=2000"`
    Rating  *float64 `json:"rating,omitempty" validate:"omitempty,min=0,max=5"`
}

// CommentListOptions 评论列表选项
type CommentListOptions struct {
    Page     int    `json:"page" validate:"min=1"`
    PageSize int    `json:"page_size" validate:"min=1,max=100"`
    OrderBy  string `json:"order_by" validate:"omitempty,oneof=created_at updated_at likes_count rating"`
    Order    string `json:"order" validate:"omitempty,oneof=asc desc"`
    Status   *models.CommentStatus `json:"status,omitempty"`
    HasRating bool  `json:"has_rating,omitempty"`
    MinRating *float64 `json:"min_rating,omitempty"`
    MaxRating *float64 `json:"max_rating,omitempty"`
}

// ReportCommentRequest 举报评论请求
type ReportCommentRequest struct {
    UserID   int64  `json:"user_id" validate:"required"`
    Reason   string `json:"reason" validate:"required,oneof=spam inappropriate offensive copyright other"`
    Details  string `json:"details,omitempty" validate:"omitempty,max=500"`
}

// ModerateCommentRequest 审核评论请求
type ModerateCommentRequest struct {
    ModeratorID int64  `json:"moderator_id" validate:"required"`
    Action      string `json:"action" validate:"required,oneof=approve reject hide"`
    Reason      string `json:"reason,omitempty" validate:"omitempty,max=500"`
}

// RateMovieRequest 评分电影请求
type RateMovieRequest struct {
    MovieID int64   `json:"movie_id" validate:"required"`
    UserID  int64   `json:"user_id" validate:"required"`
    Rating  float64 `json:"rating" validate:"required,min=0,max=5"`
    Review  string  `json:"review,omitempty" validate:"omitempty,max=2000"`
}

// UpdateRatingRequest 更新评分请求
type UpdateRatingRequest struct {
    Rating *float64 `json:"rating,omitempty" validate:"omitempty,min=0,max=5"`
    Review *string  `json:"review,omitempty" validate:"omitempty,max=2000"`
}

// RatingListOptions 评分列表选项
type RatingListOptions struct {
    Page      int      `json:"page" validate:"min=1"`
    PageSize  int      `json:"page_size" validate:"min=1,max=100"`
    OrderBy   string   `json:"order_by" validate:"omitempty,oneof=created_at rating helpful_count"`
    Order     string   `json:"order" validate:"omitempty,oneof=asc desc"`
    MinRating *float64 `json:"min_rating,omitempty"`
    MaxRating *float64 `json:"max_rating,omitempty"`
    HasReview bool     `json:"has_review,omitempty"`
}

// LikeContentRequest 点赞内容请求
type LikeContentRequest struct {
    UserID      int64  `json:"user_id" validate:"required"`
    ContentType string `json:"content_type" validate:"required,oneof=comment rating movie"`
    ContentID   int64  `json:"content_id" validate:"required"`
}

// UnlikeContentRequest 取消点赞请求
type UnlikeContentRequest struct {
    UserID      int64  `json:"user_id" validate:"required"`
    ContentType string `json:"content_type" validate:"required,oneof=comment rating movie"`
    ContentID   int64  `json:"content_id" validate:"required"`
}

// LikeListOptions 点赞列表选项
type LikeListOptions struct {
    Page        int    `json:"page" validate:"min=1"`
    PageSize    int    `json:"page_size" validate:"min=1,max=100"`
    ContentType string `json:"content_type,omitempty" validate:"omitempty,oneof=comment rating movie"`
    OrderBy     string `json:"order_by" validate:"omitempty,oneof=created_at"`
    Order       string `json:"order" validate:"omitempty,oneof=asc desc"`
}

// CommentStats 评论统计
type CommentStats struct {
    TotalComments    int64   `json:"total_comments"`
    TotalReplies     int64   `json:"total_replies"`
    AverageRating    float64 `json:"average_rating"`
    RatingCount      int64   `json:"rating_count"`
    TotalLikes       int64   `json:"total_likes"`
    ActiveUsers      int64   `json:"active_users"`
    RecentComments   int64   `json:"recent_comments"`   // 最近24小时
    TrendingScore    float64 `json:"trending_score"`
}

// UserCommentStats 用户评论统计
type UserCommentStats struct {
    TotalComments     int64   `json:"total_comments"`
    TotalRatings      int64   `json:"total_ratings"`
    AverageRating     float64 `json:"average_rating"`
    TotalLikesReceived int64  `json:"total_likes_received"`
    TotalLikesGiven   int64   `json:"total_likes_given"`
    CommentsThisMonth int64   `json:"comments_this_month"`
    RatingsThisMonth  int64   `json:"ratings_this_month"`
    PopularComments   int64   `json:"popular_comments"`   // 获得10+点赞的评论数
}

// RatingStats 评分统计
type RatingStats struct {
    TotalRatings    int64   `json:"total_ratings"`
    AverageRating   float64 `json:"average_rating"`
    MedianRating    float64 `json:"median_rating"`
    StandardDev     float64 `json:"standard_deviation"`
    RatingSum       float64 `json:"rating_sum"`
    WeightedRating  float64 `json:"weighted_rating"`
    BayesianRating  float64 `json:"bayesian_rating"`
}

// RatingDistribution 评分分布
type RatingDistribution struct {
    OneStar    int64   `json:"one_star"`
    TwoStar    int64   `json:"two_star"`
    ThreeStar  int64   `json:"three_star"`
    FourStar   int64   `json:"four_star"`
    FiveStar   int64   `json:"five_star"`
    Percentages map[string]float64 `json:"percentages"`
}

// UserRatingStats 用户评分统计
type UserRatingStats struct {
    TotalRatings      int64   `json:"total_ratings"`
    AverageRating     float64 `json:"average_rating"`
    HighestRating     float64 `json:"highest_rating"`
    LowestRating      float64 `json:"lowest_rating"`
    RatingVariance    float64 `json:"rating_variance"`
    FavoriteGenres    []GenreRating `json:"favorite_genres"`
    RatingHistory     []MonthlyRating `json:"rating_history"`
    SimilarityScore   float64 `json:"similarity_score"`
}

// SimilarUser 相似用户
type SimilarUser struct {
    UserID         int64   `json:"user_id"`
    Username       string  `json:"username"`
    SimilarityScore float64 `json:"similarity_score"`
    CommonRatings  int64   `json:"common_ratings"`
    CorrelationCoeff float64 `json:"correlation_coefficient"`
}

// RatingTrends 评分趋势
type RatingTrends struct {
    MovieID      int64         `json:"movie_id"`
    TimeWindow   string        `json:"time_window"`
    TrendPoints  []TrendPoint  `json:"trend_points"`
    OverallTrend string        `json:"overall_trend"` // increasing, decreasing, stable
    TrendScore   float64       `json:"trend_score"`
}

// TrendPoint 趋势点
type TrendPoint struct {
    Date          time.Time `json:"date"`
    AverageRating float64   `json:"average_rating"`
    RatingCount   int64     `json:"rating_count"`
    CumulativeAvg float64   `json:"cumulative_average"`
}

// LikeStats 点赞统计
type LikeStats struct {
    TotalLikes    int64            `json:"total_likes"`
    RecentLikes   int64            `json:"recent_likes"`   // 最近24小时
    LikeVelocity  float64          `json:"like_velocity"`  // 点赞速度
    TopLikers     []TopLiker       `json:"top_likers"`
    LikeHistory   []DailyLikeCount `json:"like_history"`
}

// TopLiker 热门点赞用户
type TopLiker struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    LikeCount int64 `json:"like_count"`
}

// DailyLikeCount 每日点赞数
type DailyLikeCount struct {
    Date  time.Time `json:"date"`
    Count int64     `json:"count"`
}

// GenreRating 分类评分
type GenreRating struct {
    Genre         string  `json:"genre"`
    AverageRating float64 `json:"average_rating"`
    RatingCount   int64   `json:"rating_count"`
}

// MonthlyRating 月度评分
type MonthlyRating struct {
    Month         time.Time `json:"month"`
    AverageRating float64   `json:"average_rating"`
    RatingCount   int64     `json:"rating_count"`
}
```

### 5.3.3.3 **评论服务实现**

#### 5.3.3.3.1 核心评论服务实现
```go
// internal/services/comment/service.go
package comment

import (
    "context"
    "fmt"
    "strings"
    "time"
    
    "github.com/yourname/movieinfo/internal/models"
    "github.com/yourname/movieinfo/internal/repository"
    "github.com/yourname/movieinfo/pkg/cache"
    "github.com/yourname/movieinfo/pkg/errors"
    "github.com/yourname/movieinfo/pkg/logger"
    "github.com/yourname/movieinfo/pkg/validator"
)

// ServiceImpl 评论服务实现
type ServiceImpl struct {
    commentRepo   repository.CommentRepository
    ratingRepo    repository.RatingRepository
    likeRepo      repository.LikeRepository
    userRepo      repository.UserRepository
    movieRepo     repository.MovieRepository
    cache         cache.Cache
    validator     validator.Validator
    logger        logger.Logger
    
    // 外部服务
    moderationService ModerationService
    notificationService NotificationService
    analyticsService  AnalyticsService
    
    // 配置
    config *Config
}

// Config 评论服务配置
type Config struct {
    CacheExpiry           time.Duration `yaml:"cache_expiry"`
    MaxCommentLength      int           `yaml:"max_comment_length"`
    MaxCommentsPerUser    int           `yaml:"max_comments_per_user"`
    CommentCooldown       time.Duration `yaml:"comment_cooldown"`
    EnableModeration      bool          `yaml:"enable_moderation"`
    EnableSentimentAnalysis bool        `yaml:"enable_sentiment_analysis"`
    AutoHideThreshold     float64       `yaml:"auto_hide_threshold"`
    TrendingTimeWindow    time.Duration `yaml:"trending_time_window"`
}

// NewService 创建评论服务
func NewService(
    commentRepo repository.CommentRepository,
    ratingRepo repository.RatingRepository,
    likeRepo repository.LikeRepository,
    userRepo repository.UserRepository,
    movieRepo repository.MovieRepository,
    cache cache.Cache,
    validator validator.Validator,
    moderationService ModerationService,
    notificationService NotificationService,
    analyticsService AnalyticsService,
    config *Config,
) CommentService {
    return &ServiceImpl{
        commentRepo:         commentRepo,
        ratingRepo:          ratingRepo,
        likeRepo:            likeRepo,
        userRepo:            userRepo,
        movieRepo:           movieRepo,
        cache:               cache,
        validator:           validator,
        logger:              logger.GetGlobalLogger(),
        moderationService:   moderationService,
        notificationService: notificationService,
        analyticsService:    analyticsService,
        config:              config,
    }
}

// CreateComment 创建评论
func (s *ServiceImpl) CreateComment(ctx context.Context, req *CreateCommentRequest) (*models.Comment, error) {
    // 验证请求参数
    if err := s.validator.Validate(req); err != nil {
        return nil, errors.ValidationFailed(err.Error())
    }
    
    s.logger.WithContext(ctx).Info("Creating comment",
        logger.Int64("user_id", req.UserID),
        logger.Int64("movie_id", req.MovieID),
        logger.Int("content_length", len(req.Content)),
    )
    
    // 检查用户是否存在
    user, err := s.userRepo.GetByID(ctx, req.UserID)
    if err != nil {
        if errors.IsNotFound(err) {
            return nil, errors.UserNotFound()
        }
        return nil, errors.Wrap(err, errors.CodeDatabaseError, "查询用户失败")
    }
    
    // 检查用户状态
    if !user.IsActive() {
        return nil, errors.Forbidden("用户状态异常，无法发表评论")
    }
    
    // 检查电影是否存在
    movie, err := s.movieRepo.GetByID(ctx, req.MovieID)
    if err != nil {
        if errors.IsNotFound(err) {
            return nil, errors.MovieNotFound()
        }
        return nil, errors.Wrap(err, errors.CodeDatabaseError, "查询电影失败")
    }
    
    // 检查评论冷却时间
    if err := s.checkCommentCooldown(ctx, req.UserID); err != nil {
        return nil, err
    }
    
    // 检查用户评论数量限制
    if err := s.checkUserCommentLimit(ctx, req.UserID); err != nil {
        return nil, err
    }
    
    // 内容审核
    moderationResult, err := s.moderateContent(ctx, req.Content)
    if err != nil {
        s.logger.WithContext(ctx).Error("Content moderation failed",
            logger.Error(err),
        )
        // 不阻断评论流程，使用默认状态
        moderationResult = &ModerationResult{
            IsApproved: true,
            Confidence: 0.5,
        }
    }
    
    // 创建评论
    comment := &models.Comment{
        MovieID:  req.MovieID,
        UserID:   req.UserID,
        Content:  req.Content,
        ParentID: req.ParentID,
        Status:   models.CommentStatusPending,
    }
    
    // 根据审核结果设置状态
    if moderationResult.IsApproved {
        comment.Status = models.CommentStatusApproved
    } else if moderationResult.Confidence > s.config.AutoHideThreshold {
        comment.Status = models.CommentStatusHidden
    }
    
    // 设置情感分析结果
    if s.config.EnableSentimentAnalysis && moderationResult.Sentiment != nil {
        comment.SentimentScore = &moderationResult.Sentiment.Score
        comment.SentimentLabel = &moderationResult.Sentiment.Label
    }
    
    // 保存评论
    if err := s.commentRepo.Create(ctx, comment); err != nil {
        return nil, errors.Wrap(err, errors.CodeDatabaseError, "创建评论失败")
    }
    
    // 如果包含评分，创建评分记录
    if req.Rating != nil {
        ratingReq := &RateMovieRequest{
            MovieID: req.MovieID,
            UserID:  req.UserID,
            Rating:  *req.Rating,
        }
        
        if _, err := s.createRating(ctx, ratingReq); err != nil {
            s.logger.WithContext(ctx).Error("Failed to create rating",
                logger.Error(err),
            )
            // 不阻断评论流程
        }
    }
    
    // 更新电影评论数
    go func() {
        if err := s.movieRepo.IncrementCommentCount(ctx, req.MovieID); err != nil {
            s.logger.Error("Failed to increment movie comment count",
                logger.Int64("movie_id", req.MovieID),
                logger.Error(err),
            )
        }
    }()
    
    // 更新用户评论数
    go func() {
        if err := s.userRepo.IncrementReviewsCount(ctx, req.UserID); err != nil {
            s.logger.Error("Failed to increment user reviews count",
                logger.Int64("user_id", req.UserID),
                logger.Error(err),
            )
        }
    }()
    
    // 发送通知
    go func() {
        s.sendCommentNotifications(ctx, comment, user, movie)
    }()
    
    // 清除相关缓存
    s.clearCommentCache(ctx, req.MovieID, req.UserID)
    
    s.logger.WithContext(ctx).Info("Comment created successfully",
        logger.Int64("comment_id", comment.ID),
        logger.String("status", comment.Status.String()),
    )
    
    return comment, nil
}

// GetMovieComments 获取电影评论
func (s *ServiceImpl) GetMovieComments(ctx context.Context, movieID int64, opts *CommentListOptions) (*repository.PaginationResult[models.Comment], error) {
    // 构建缓存键
    cacheKey := fmt.Sprintf("movie:comments:%d:%d:%d:%s:%s", 
        movieID, opts.Page, opts.PageSize, opts.OrderBy, opts.Order)
    
    // 先从缓存获取
    if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
        if result, ok := cached.(*repository.PaginationResult[models.Comment]); ok {
            return result, nil
        }
    }
    
    // 设置默认排序
    if opts.OrderBy == "" {
        opts.OrderBy = "created_at"
        opts.Order = "desc"
    }
    
    // 只显示已审核的评论
    if opts.Status == nil {
        status := models.CommentStatusApproved
        opts.Status = &status
    }
    
    // 从数据库获取
    listOpts := &repository.ListOptions{
        Offset:  (opts.Page - 1) * opts.PageSize,
        Limit:   opts.PageSize,
        OrderBy: opts.OrderBy,
        Order:   opts.Order,
        Filters: map[string]interface{}{
            "movie_id": movieID,
            "status":   *opts.Status,
        },
    }
    
    // 添加其他过滤条件
    if opts.HasRating {
        listOpts.Filters["has_rating"] = true
    }
    if opts.MinRating != nil {
        listOpts.Filters["min_rating"] = *opts.MinRating
    }
    if opts.MaxRating != nil {
        listOpts.Filters["max_rating"] = *opts.MaxRating
    }
    
    comments, err := s.commentRepo.List(ctx, listOpts)
    if err != nil {
        return nil, errors.Wrap(err, errors.CodeDatabaseError, "获取评论列表失败")
    }
    
    // 获取总数
    countOpts := &repository.CountOptions{
        Filters: listOpts.Filters,
    }
    total, err := s.commentRepo.Count(ctx, countOpts)
    if err != nil {
        return nil, errors.Wrap(err, errors.CodeDatabaseError, "获取评论总数失败")
    }
    
    // 构建分页结果
    result := repository.NewPaginationResult(comments, total, opts.Page, opts.PageSize)
    
    // 缓存结果
    s.cache.Set(ctx, cacheKey, result, s.config.CacheExpiry)
    
    return result, nil
}

// LikeComment 点赞评论
func (s *ServiceImpl) LikeComment(ctx context.Context, commentID int64, userID int64) error {
    // 检查评论是否存在
    comment, err := s.commentRepo.GetByID(ctx, commentID)
    if err != nil {
        if errors.IsNotFound(err) {
            return errors.CommentNotFound()
        }
        return errors.Wrap(err, errors.CodeDatabaseError, "查询评论失败")
    }
    
    // 检查是否已经点赞
    exists, err := s.likeRepo.Exists(ctx, "comment", commentID, userID)
    if err != nil {
        return errors.Wrap(err, errors.CodeDatabaseError, "检查点赞状态失败")
    }
    if exists {
        return errors.BadRequest("已经点赞过该评论")
    }
    
    // 创建点赞记录
    like := &models.Like{
        UserID:      userID,
        ContentType: "comment",
        ContentID:   commentID,
    }
    
    if err := s.likeRepo.Create(ctx, like); err != nil {
        return errors.Wrap(err, errors.CodeDatabaseError, "创建点赞记录失败")
    }
    
    // 更新评论点赞数
    if err := s.commentRepo.IncrementLikeCount(ctx, commentID); err != nil {
        s.logger.WithContext(ctx).Error("Failed to increment comment like count",
            logger.Int64("comment_id", commentID),
            logger.Error(err),
        )
    }
    
    // 发送通知
    go func() {
        s.sendLikeNotification(ctx, comment, userID)
    }()
    
    // 清除缓存
    s.clearLikeCache(ctx, commentID, userID)
    
    return nil
}

// checkCommentCooldown 检查评论冷却时间
func (s *ServiceImpl) checkCommentCooldown(ctx context.Context, userID int64) error {
    cacheKey := fmt.Sprintf("user:comment_cooldown:%d", userID)
    
    if exists, _ := s.cache.Exists(ctx, cacheKey); exists {
        return errors.TooManyRequests("评论过于频繁，请稍后再试")
    }
    
    // 设置冷却时间
    s.cache.Set(ctx, cacheKey, true, s.config.CommentCooldown)
    
    return nil
}

// checkUserCommentLimit 检查用户评论数量限制
func (s *ServiceImpl) checkUserCommentLimit(ctx context.Context, userID int64) error {
    // 检查今日评论数量
    cacheKey := fmt.Sprintf("user:daily_comments:%d:%s", userID, time.Now().Format("2006-01-02"))
    
    count, _ := s.cache.Get(ctx, cacheKey)
    if dailyCount, ok := count.(int); ok && dailyCount >= s.config.MaxCommentsPerUser {
        return errors.TooManyRequests("今日评论数量已达上限")
    }
    
    // 增加计数
    newCount := 1
    if count != nil {
        if dailyCount, ok := count.(int); ok {
            newCount = dailyCount + 1
        }
    }
    
    s.cache.Set(ctx, cacheKey, newCount, 24*time.Hour)
    
    return nil
}

// moderateContent 内容审核
func (s *ServiceImpl) moderateContent(ctx context.Context, content string) (*ModerationResult, error) {
    if !s.config.EnableModeration {
        return &ModerationResult{
            IsApproved: true,
            Confidence: 1.0,
        }, nil
    }
    
    return s.moderationService.ModerateText(ctx, content)
}

// clearCommentCache 清除评论相关缓存
func (s *ServiceImpl) clearCommentCache(ctx context.Context, movieID, userID int64) {
    patterns := []string{
        fmt.Sprintf("movie:comments:%d:*", movieID),
        fmt.Sprintf("user:comments:%d:*", userID),
        "comments:trending:*",
        "comments:stats:*",
    }
    
    for _, pattern := range patterns {
        s.cache.DeletePattern(ctx, pattern)
    }
}
```

## 5.3.4 总结

评论服务实现为MovieInfo项目提供了完整的用户互动解决方案。通过内容审核、社交功能和数据分析，我们建立了一个安全、活跃、智能的评论系统。

**关键设计要点**：
1. **内容安全**：自动审核、敏感词过滤、垃圾评论检测
2. **用户体验**：实时互动、个性化排序、智能推荐
3. **社交功能**：点赞、回复、关注、通知系统
4. **数据分析**：情感分析、热点检测、用户画像
5. **性能优化**：多层缓存、异步处理、批量操作

**服务优势**：
- **内容质量**：智能审核保证内容质量
- **用户参与**：丰富的互动功能提升参与度
- **数据洞察**：深度分析提供业务洞察
- **系统稳定**：高并发处理和性能优化

**下一步**：基于完整的业务逻辑层，我们将开始API层开发，实现HTTP路由设计和请求处理器。
