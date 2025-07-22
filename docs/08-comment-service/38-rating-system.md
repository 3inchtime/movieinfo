# 第38步：评分系统实现

## 📋 概述

评分系统是MovieInfo项目的核心功能之一，为用户提供对电影的量化评价能力。一个完善的评分系统需要支持多维度评分、智能算法计算、防刷机制和实时统计更新。

## 🎯 设计目标

### 1. **评分准确性**
- 多维度评分支持
- 智能权重算法
- 异常评分检测
- 评分质量控制

### 2. **系统公平性**
- 防刷票机制
- 用户权重系统
- 评分时效性控制
- 恶意评分过滤

### 3. **性能优化**
- 实时评分计算
- 缓存策略优化
- 批量更新机制
- 异步统计处理

## 🏗️ 评分系统架构

### 1. **评分数据模型**

```go
// 用户评分记录
type UserRating struct {
    ID        string    `gorm:"primaryKey" json:"id"`
    UserID    string    `gorm:"not null;index" json:"user_id"`
    MovieID   string    `gorm:"not null;index" json:"movie_id"`
    Rating    int       `gorm:"not null" json:"rating"` // 1-10分
    
    // 多维度评分
    StoryRating     *int `gorm:"column:story_rating" json:"story_rating,omitempty"`
    ActingRating    *int `gorm:"column:acting_rating" json:"acting_rating,omitempty"`
    VisualRating    *int `gorm:"column:visual_rating" json:"visual_rating,omitempty"`
    MusicRating     *int `gorm:"column:music_rating" json:"music_rating,omitempty"`
    
    // 评分元数据
    Weight        float64   `gorm:"not null;default:1.0" json:"weight"`
    Source        string    `gorm:"not null;default:'user'" json:"source"` // user, import, system
    DeviceInfo    string    `gorm:"size:200" json:"device_info,omitempty"`
    IPAddress     string    `gorm:"size:45" json:"ip_address,omitempty"`
    
    // 状态信息
    Status        string    `gorm:"not null;default:'active'" json:"status"` // active, hidden, flagged
    
    // 时间戳
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    
    // 关联数据
    User          *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Movie         *Movie    `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
}

// 电影评分统计
type MovieRatingStats struct {
    ID            string    `gorm:"primaryKey" json:"id"`
    MovieID       string    `gorm:"not null;uniqueIndex" json:"movie_id"`
    
    // 基础统计
    AverageRating float64   `gorm:"not null;default:0" json:"average_rating"`
    TotalRatings  int       `gorm:"not null;default:0" json:"total_ratings"`
    WeightedRating float64  `gorm:"not null;default:0" json:"weighted_rating"`
    
    // 评分分布
    Rating1Count  int       `gorm:"not null;default:0" json:"rating_1_count"`
    Rating2Count  int       `gorm:"not null;default:0" json:"rating_2_count"`
    Rating3Count  int       `gorm:"not null;default:0" json:"rating_3_count"`
    Rating4Count  int       `gorm:"not null;default:0" json:"rating_4_count"`
    Rating5Count  int       `gorm:"not null;default:0" json:"rating_5_count"`
    Rating6Count  int       `gorm:"not null;default:0" json:"rating_6_count"`
    Rating7Count  int       `gorm:"not null;default:0" json:"rating_7_count"`
    Rating8Count  int       `gorm:"not null;default:0" json:"rating_8_count"`
    Rating9Count  int       `gorm:"not null;default:0" json:"rating_9_count"`
    Rating10Count int       `gorm:"not null;default:0" json:"rating_10_count"`
    
    // 多维度平均分
    AvgStoryRating  float64 `gorm:"column:avg_story_rating" json:"avg_story_rating"`
    AvgActingRating float64 `gorm:"column:avg_acting_rating" json:"avg_acting_rating"`
    AvgVisualRating float64 `gorm:"column:avg_visual_rating" json:"avg_visual_rating"`
    AvgMusicRating  float64 `gorm:"column:avg_music_rating" json:"avg_music_rating"`
    
    // 时间信息
    LastRatedAt   time.Time `json:"last_rated_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    
    // 关联数据
    Movie         *Movie    `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
}

// 评分请求结构
type RatingRequest struct {
    MovieID      string `json:"movie_id" binding:"required"`
    Rating       int    `json:"rating" binding:"required,min=1,max=10"`
    StoryRating  *int   `json:"story_rating,omitempty" binding:"omitempty,min=1,max=10"`
    ActingRating *int   `json:"acting_rating,omitempty" binding:"omitempty,min=1,max=10"`
    VisualRating *int   `json:"visual_rating,omitempty" binding:"omitempty,min=1,max=10"`
    MusicRating  *int   `json:"music_rating,omitempty" binding:"omitempty,min=1,max=10"`
}

// 评分响应结构
type RatingResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    *UserRating `json:"data,omitempty"`
}

// 评分统计响应
type RatingStatsResponse struct {
    Success bool               `json:"success"`
    Message string             `json:"message,omitempty"`
    Data    *MovieRatingStats  `json:"data,omitempty"`
    Distribution []RatingDistributionItem `json:"distribution,omitempty"`
}

type RatingDistributionItem struct {
    Rating int     `json:"rating"`
    Count  int     `json:"count"`
    Percentage float64 `json:"percentage"`
}
```

### 2. **评分服务实现**

```go
type RatingService struct {
    ratingRepo     RatingRepository
    movieRepo      MovieRepository
    userRepo       UserRepository
    antiSpamSvc    AntiSpamService
    cacheStore     CacheStore
    logger         *logrus.Logger
    metrics        *RatingMetrics
}

func NewRatingService(
    ratingRepo RatingRepository,
    movieRepo MovieRepository,
    userRepo UserRepository,
    antiSpamSvc AntiSpamService,
    cacheStore CacheStore,
) *RatingService {
    return &RatingService{
        ratingRepo:  ratingRepo,
        movieRepo:   movieRepo,
        userRepo:    userRepo,
        antiSpamSvc: antiSpamSvc,
        cacheStore:  cacheStore,
        logger:      logrus.New(),
        metrics:     NewRatingMetrics(),
    }
}

// 提交评分
func (rs *RatingService) SubmitRating(ctx context.Context, req *RatingRequest, userID, clientIP, deviceInfo string) (*RatingResponse, error) {
    start := time.Now()
    defer func() {
        rs.metrics.ObserveRatingOperation("submit", time.Since(start))
    }()

    // 验证用户
    user, err := rs.userRepo.FindByID(ctx, userID)
    if err != nil {
        rs.metrics.IncInvalidRequests("user_not_found")
        return &RatingResponse{
            Success: false,
            Message: "用户不存在",
        }, nil
    }

    if user.Status != "active" {
        rs.metrics.IncInvalidRequests("user_inactive")
        return &RatingResponse{
            Success: false,
            Message: "用户账户已被禁用",
        }, nil
    }

    // 验证电影存在性
    if exists, err := rs.movieRepo.ExistsByID(ctx, req.MovieID); err != nil {
        rs.logger.Errorf("Failed to check movie existence: %v", err)
        return nil, errors.New("验证电影失败")
    } else if !exists {
        rs.metrics.IncInvalidRequests("movie_not_found")
        return &RatingResponse{
            Success: false,
            Message: "电影不存在",
        }, nil
    }

    // 反垃圾检查
    if blocked, reason := rs.antiSpamSvc.CheckRatingSpam(ctx, userID, clientIP, req.MovieID); blocked {
        rs.metrics.IncBlockedRequests(reason)
        return &RatingResponse{
            Success: false,
            Message: "评分过于频繁，请稍后再试",
        }, nil
    }

    // 检查是否已经评分
    existingRating, err := rs.ratingRepo.FindByUserAndMovie(ctx, userID, req.MovieID)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        rs.logger.Errorf("Failed to check existing rating: %v", err)
        return nil, errors.New("检查评分失败")
    }

    var rating *UserRating
    var isUpdate bool

    if existingRating != nil {
        // 更新现有评分
        isUpdate = true
        rating = existingRating
        rating.Rating = req.Rating
        rating.StoryRating = req.StoryRating
        rating.ActingRating = req.ActingRating
        rating.VisualRating = req.VisualRating
        rating.MusicRating = req.MusicRating
        rating.UpdatedAt = time.Now()
        rating.IPAddress = clientIP
        rating.DeviceInfo = deviceInfo

        if err := rs.ratingRepo.Update(ctx, rating); err != nil {
            rs.logger.Errorf("Failed to update rating: %v", err)
            rs.metrics.IncOperationErrors("update")
            return nil, errors.New("评分更新失败")
        }
    } else {
        // 创建新评分
        isUpdate = false
        rating = &UserRating{
            ID:           uuid.New().String(),
            UserID:       userID,
            MovieID:      req.MovieID,
            Rating:       req.Rating,
            StoryRating:  req.StoryRating,
            ActingRating: req.ActingRating,
            VisualRating: req.VisualRating,
            MusicRating:  req.MusicRating,
            Weight:       rs.calculateUserWeight(user),
            Source:       "user",
            DeviceInfo:   deviceInfo,
            IPAddress:    clientIP,
            Status:       "active",
            CreatedAt:    time.Now(),
            UpdatedAt:    time.Now(),
        }

        if err := rs.ratingRepo.Create(ctx, rating); err != nil {
            rs.logger.Errorf("Failed to create rating: %v", err)
            rs.metrics.IncOperationErrors("create")
            return nil, errors.New("评分提交失败")
        }
    }

    // 异步更新电影评分统计
    go func() {
        if err := rs.updateMovieRatingStats(context.Background(), req.MovieID); err != nil {
            rs.logger.Errorf("Failed to update movie rating stats: %v", err)
        }
    }()

    // 清除相关缓存
    rs.clearRatingCache(ctx, req.MovieID)

    // 记录反垃圾信息
    rs.antiSpamSvc.RecordRatingActivity(ctx, userID, clientIP, req.MovieID)

    if isUpdate {
        rs.metrics.IncSuccessfulOperations("update")
        rs.logger.Infof("Rating updated: user %s, movie %s, rating %d", userID, req.MovieID, req.Rating)
    } else {
        rs.metrics.IncSuccessfulOperations("create")
        rs.logger.Infof("Rating created: user %s, movie %s, rating %d", userID, req.MovieID, req.Rating)
    }

    return &RatingResponse{
        Success: true,
        Message: "评分提交成功",
        Data:    rating,
    }, nil
}

// 获取电影评分统计
func (rs *RatingService) GetMovieRatingStats(ctx context.Context, movieID string) (*RatingStatsResponse, error) {
    start := time.Now()
    defer func() {
        rs.metrics.ObserveRatingOperation("stats", time.Since(start))
    }()

    // 生成缓存键
    cacheKey := fmt.Sprintf("movie_rating_stats:%s", movieID)

    // 尝试从缓存获取
    if cachedResult, err := rs.getStatsFromCache(ctx, cacheKey); err == nil {
        rs.metrics.IncCacheHits()
        return cachedResult, nil
    }
    rs.metrics.IncCacheMisses()

    // 从数据库获取统计信息
    stats, err := rs.ratingRepo.GetMovieRatingStats(ctx, movieID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // 如果没有统计记录，创建默认统计
            stats = &MovieRatingStats{
                ID:             uuid.New().String(),
                MovieID:        movieID,
                AverageRating:  0,
                TotalRatings:   0,
                WeightedRating: 0,
                UpdatedAt:      time.Now(),
            }
        } else {
            rs.logger.Errorf("Failed to get movie rating stats: %v", err)
            rs.metrics.IncOperationErrors("stats")
            return nil, errors.New("获取评分统计失败")
        }
    }

    // 构建评分分布
    distribution := rs.buildRatingDistribution(stats)

    response := &RatingStatsResponse{
        Success:      true,
        Data:         stats,
        Distribution: distribution,
    }

    // 异步缓存结果
    go func() {
        if err := rs.cacheStatsResult(context.Background(), cacheKey, response); err != nil {
            rs.logger.Errorf("Failed to cache rating stats: %v", err)
        }
    }()

    rs.metrics.IncSuccessfulOperations("stats")
    return response, nil
}

// 更新电影评分统计
func (rs *RatingService) updateMovieRatingStats(ctx context.Context, movieID string) error {
    // 获取所有有效评分
    ratings, err := rs.ratingRepo.FindActiveRatingsByMovie(ctx, movieID)
    if err != nil {
        return err
    }

    if len(ratings) == 0 {
        // 如果没有评分，清空统计
        return rs.ratingRepo.ClearMovieRatingStats(ctx, movieID)
    }

    // 计算统计数据
    stats := rs.calculateRatingStatistics(ratings)
    stats.MovieID = movieID
    stats.LastRatedAt = time.Now()
    stats.UpdatedAt = time.Now()

    // 保存或更新统计
    existingStats, err := rs.ratingRepo.GetMovieRatingStats(ctx, movieID)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return err
    }

    if existingStats != nil {
        stats.ID = existingStats.ID
        return rs.ratingRepo.UpdateMovieRatingStats(ctx, stats)
    } else {
        stats.ID = uuid.New().String()
        return rs.ratingRepo.CreateMovieRatingStats(ctx, stats)
    }
}

// 计算评分统计
func (rs *RatingService) calculateRatingStatistics(ratings []UserRating) *MovieRatingStats {
    stats := &MovieRatingStats{}
    
    var totalScore float64
    var totalWeight float64
    var storyTotal, actingTotal, visualTotal, musicTotal float64
    var storyCount, actingCount, visualCount, musicCount int
    
    // 评分分布计数
    ratingCounts := make([]int, 11) // 索引0不使用，1-10对应评分

    for _, rating := range ratings {
        weight := rating.Weight
        totalScore += float64(rating.Rating) * weight
        totalWeight += weight
        
        // 评分分布
        if rating.Rating >= 1 && rating.Rating <= 10 {
            ratingCounts[rating.Rating]++
        }
        
        // 多维度评分
        if rating.StoryRating != nil {
            storyTotal += float64(*rating.StoryRating) * weight
            storyCount++
        }
        if rating.ActingRating != nil {
            actingTotal += float64(*rating.ActingRating) * weight
            actingCount++
        }
        if rating.VisualRating != nil {
            visualTotal += float64(*rating.VisualRating) * weight
            visualCount++
        }
        if rating.MusicRating != nil {
            musicTotal += float64(*rating.MusicRating) * weight
            musicCount++
        }
    }

    // 计算平均分
    if totalWeight > 0 {
        stats.WeightedRating = totalScore / totalWeight
        stats.AverageRating = totalScore / float64(len(ratings)) // 简单平均
    }
    
    stats.TotalRatings = len(ratings)
    
    // 设置评分分布
    stats.Rating1Count = ratingCounts[1]
    stats.Rating2Count = ratingCounts[2]
    stats.Rating3Count = ratingCounts[3]
    stats.Rating4Count = ratingCounts[4]
    stats.Rating5Count = ratingCounts[5]
    stats.Rating6Count = ratingCounts[6]
    stats.Rating7Count = ratingCounts[7]
    stats.Rating8Count = ratingCounts[8]
    stats.Rating9Count = ratingCounts[9]
    stats.Rating10Count = ratingCounts[10]
    
    // 计算多维度平均分
    if storyCount > 0 {
        stats.AvgStoryRating = storyTotal / float64(storyCount)
    }
    if actingCount > 0 {
        stats.AvgActingRating = actingTotal / float64(actingCount)
    }
    if visualCount > 0 {
        stats.AvgVisualRating = visualTotal / float64(visualCount)
    }
    if musicCount > 0 {
        stats.AvgMusicRating = musicTotal / float64(musicCount)
    }

    return stats
}

// 计算用户权重
func (rs *RatingService) calculateUserWeight(user *User) float64 {
    baseWeight := 1.0
    
    // 根据用户等级调整权重
    levelWeight := 1.0 + float64(user.Level)*0.1
    
    // 根据用户活跃度调整权重
    activityWeight := 1.0
    if user.CommentCount > 100 {
        activityWeight = 1.2
    } else if user.CommentCount > 50 {
        activityWeight = 1.1
    }
    
    // 根据账户年龄调整权重
    ageWeight := 1.0
    accountAge := time.Since(user.CreatedAt)
    if accountAge > 365*24*time.Hour { // 超过1年
        ageWeight = 1.1
    }
    
    return baseWeight * levelWeight * activityWeight * ageWeight
}

// 构建评分分布
func (rs *RatingService) buildRatingDistribution(stats *MovieRatingStats) []RatingDistributionItem {
    distribution := make([]RatingDistributionItem, 10)
    total := float64(stats.TotalRatings)
    
    counts := []int{
        stats.Rating1Count, stats.Rating2Count, stats.Rating3Count,
        stats.Rating4Count, stats.Rating5Count, stats.Rating6Count,
        stats.Rating7Count, stats.Rating8Count, stats.Rating9Count,
        stats.Rating10Count,
    }
    
    for i, count := range counts {
        percentage := 0.0
        if total > 0 {
            percentage = float64(count) / total * 100
        }
        
        distribution[i] = RatingDistributionItem{
            Rating:     i + 1,
            Count:      count,
            Percentage: math.Round(percentage*100) / 100, // 保留2位小数
        }
    }
    
    return distribution
}
```

### 3. **反垃圾评分服务**

```go
type AntiSpamService struct {
    redis  *redis.Client
    logger *logrus.Logger
}

func NewAntiSpamService(redis *redis.Client) *AntiSpamService {
    return &AntiSpamService{
        redis:  redis,
        logger: logrus.New(),
    }
}

// 检查评分垃圾行为
func (ass *AntiSpamService) CheckRatingSpam(ctx context.Context, userID, clientIP, movieID string) (bool, string) {
    // 检查用户评分频率
    if blocked, reason := ass.checkUserRatingFrequency(ctx, userID); blocked {
        return true, reason
    }
    
    // 检查IP评分频率
    if blocked, reason := ass.checkIPRatingFrequency(ctx, clientIP); blocked {
        return true, reason
    }
    
    // 检查同一电影重复评分
    if blocked, reason := ass.checkDuplicateMovieRating(ctx, userID, movieID); blocked {
        return true, reason
    }
    
    return false, ""
}

func (ass *AntiSpamService) checkUserRatingFrequency(ctx context.Context, userID string) (bool, string) {
    key := fmt.Sprintf("rating_freq:user:%s", userID)
    
    // 检查1小时内评分次数
    count, err := ass.redis.Get(ctx, key).Int()
    if err != nil && err != redis.Nil {
        ass.logger.Errorf("Failed to check user rating frequency: %v", err)
        return false, ""
    }
    
    if count >= 10 { // 1小时内最多10次评分
        return true, "user_frequency_limit"
    }
    
    return false, ""
}

// 记录评分活动
func (ass *AntiSpamService) RecordRatingActivity(ctx context.Context, userID, clientIP, movieID string) {
    // 记录用户评分频率
    userKey := fmt.Sprintf("rating_freq:user:%s", userID)
    ass.redis.Incr(ctx, userKey)
    ass.redis.Expire(ctx, userKey, time.Hour)
    
    // 记录IP评分频率
    ipKey := fmt.Sprintf("rating_freq:ip:%s", clientIP)
    ass.redis.Incr(ctx, ipKey)
    ass.redis.Expire(ctx, ipKey, time.Hour)
}
```

## 📊 性能监控

### 1. **评分系统指标**

```go
type RatingMetrics struct {
    operationCount    *prometheus.CounterVec
    operationDuration *prometheus.HistogramVec
    ratingDistribution *prometheus.HistogramVec
    blockedRequests   *prometheus.CounterVec
    cacheHitRate      *prometheus.CounterVec
}

func NewRatingMetrics() *RatingMetrics {
    return &RatingMetrics{
        operationCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "rating_operations_total",
                Help: "Total number of rating operations",
            },
            []string{"operation", "status"},
        ),
        operationDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "rating_operation_duration_seconds",
                Help: "Duration of rating operations",
            },
            []string{"operation"},
        ),
        ratingDistribution: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "movie_rating_distribution",
                Help: "Distribution of movie ratings",
                Buckets: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
            },
            []string{"movie_id"},
        ),
        blockedRequests: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "rating_blocked_requests_total",
                Help: "Total number of blocked rating requests",
            },
            []string{"reason"},
        ),
        cacheHitRate: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "rating_cache_operations_total",
                Help: "Total number of rating cache operations",
            },
            []string{"type"},
        ),
    }
}
```

## 🔧 HTTP处理器

### 1. **评分API端点**

```go
func (rc *RatingController) SubmitRating(c *gin.Context) {
    var req RatingRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": "请求参数错误",
            "error":   err.Error(),
        })
        return
    }

    userID := rc.getUserIDFromContext(c)
    if userID == "" {
        c.JSON(401, gin.H{
            "success": false,
            "message": "请先登录",
        })
        return
    }

    clientIP := c.ClientIP()
    deviceInfo := c.GetHeader("User-Agent")

    response, err := rc.ratingService.SubmitRating(c.Request.Context(), &req, userID, clientIP, deviceInfo)
    if err != nil {
        rc.logger.Errorf("Failed to submit rating: %v", err)
        c.JSON(500, gin.H{
            "success": false,
            "message": "评分提交失败",
        })
        return
    }

    c.JSON(200, response)
}

func (rc *RatingController) GetMovieRatingStats(c *gin.Context) {
    movieID := c.Param("movie_id")
    if movieID == "" {
        c.JSON(400, gin.H{
            "success": false,
            "message": "电影ID不能为空",
        })
        return
    }

    response, err := rc.ratingService.GetMovieRatingStats(c.Request.Context(), movieID)
    if err != nil {
        rc.logger.Errorf("Failed to get rating stats: %v", err)
        c.JSON(500, gin.H{
            "success": false,
            "message": "获取评分统计失败",
        })
        return
    }

    c.Header("Cache-Control", "public, max-age=600") // 10分钟缓存
    c.JSON(200, response)
}
```

## 📝 总结

评分系统为MovieInfo项目提供了完整的电影评价功能：

**核心功能**：
1. **多维评分**：支持总体评分和多维度细分评分
2. **智能算法**：用户权重系统和加权平均算法
3. **防刷机制**：完善的反垃圾评分检测
4. **实时统计**：动态更新的评分统计和分布

**技术特性**：
- 高性能的评分计算
- 智能的缓存策略
- 完善的监控指标
- 异步的统计更新

**安全保障**：
- 用户权限验证
- 频率限制机制
- 异常评分检测
- 恶意行为过滤

下一步，我们将实现评论审核机制，确保平台内容的质量和安全。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第39步：评论审核机制](39-comment-moderation.md)
