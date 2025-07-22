# 第37步：评论 CRUD 接口

## 📋 概述

评论CRUD接口是评论服务的核心功能，为用户提供完整的评论管理能力。一个完善的评论系统需要支持评论的创建、读取、更新、删除操作，同时确保数据安全性和用户体验。

## 🎯 设计目标

### 1. **功能完整性**
- 完整的CRUD操作
- 评论层级管理
- 批量操作支持
- 软删除机制

### 2. **性能优化**
- 高效的分页查询
- 智能缓存策略
- 数据库查询优化
- 异步处理支持

### 3. **安全保障**
- 权限验证机制
- 内容安全检查
- 操作审计日志
- 防刷机制

## 🔧 接口设计

### 1. **评论数据结构**

```go
// 评论请求结构
type CommentCreateRequest struct {
    MovieID    string `json:"movie_id" binding:"required"`
    Content    string `json:"content" binding:"required,min=1,max=2000"`
    ParentID   string `json:"parent_id,omitempty"`
    Rating     *int   `json:"rating,omitempty" binding:"omitempty,min=1,max=10"`
    Spoiler    bool   `json:"spoiler"`
    Anonymous  bool   `json:"anonymous"`
}

type CommentUpdateRequest struct {
    Content   string `json:"content" binding:"required,min=1,max=2000"`
    Spoiler   bool   `json:"spoiler"`
    Anonymous bool   `json:"anonymous"`
}

// 评论响应结构
type CommentResponse struct {
    Success bool     `json:"success"`
    Message string   `json:"message,omitempty"`
    Data    *Comment `json:"data,omitempty"`
}

type CommentListResponse struct {
    Success    bool            `json:"success"`
    Message    string          `json:"message,omitempty"`
    Data       []Comment       `json:"data,omitempty"`
    Pagination *PaginationInfo `json:"pagination,omitempty"`
    Statistics *CommentStats   `json:"statistics,omitempty"`
}

// 评论详细信息
type Comment struct {
    ID          string    `json:"id"`
    MovieID     string    `json:"movie_id"`
    UserID      string    `json:"user_id"`
    ParentID    *string   `json:"parent_id,omitempty"`
    Content     string    `json:"content"`
    Rating      *int      `json:"rating,omitempty"`
    Spoiler     bool      `json:"spoiler"`
    Anonymous   bool      `json:"anonymous"`
    
    // 用户信息
    User        *CommentUser `json:"user,omitempty"`
    
    // 统计信息
    LikeCount   int       `json:"like_count"`
    ReplyCount  int       `json:"reply_count"`
    ReportCount int       `json:"report_count"`
    
    // 用户交互状态
    IsLiked     bool      `json:"is_liked"`
    IsReported  bool      `json:"is_reported"`
    
    // 审核状态
    Status      string    `json:"status"` // pending, approved, rejected, hidden
    
    // 子评论
    Replies     []Comment `json:"replies,omitempty"`
    
    // 时间信息
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type CommentUser struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Avatar   string `json:"avatar,omitempty"`
    Level    int    `json:"level,omitempty"`
}

type CommentStats struct {
    TotalComments    int                    `json:"total_comments"`
    TotalReplies     int                    `json:"total_replies"`
    AverageRating    float64                `json:"average_rating"`
    RatingDistribution map[string]int       `json:"rating_distribution"`
    StatusDistribution map[string]int       `json:"status_distribution"`
}
```

### 2. **评论服务实现**

```go
type CommentService struct {
    commentRepo    CommentRepository
    userRepo       UserRepository
    movieRepo      MovieRepository
    moderationSvc  ModerationService
    notificationSvc NotificationService
    cacheStore     CacheStore
    logger         *logrus.Logger
    metrics        *CommentMetrics
}

func NewCommentService(
    commentRepo CommentRepository,
    userRepo UserRepository,
    movieRepo MovieRepository,
    moderationSvc ModerationService,
    notificationSvc NotificationService,
    cacheStore CacheStore,
) *CommentService {
    return &CommentService{
        commentRepo:     commentRepo,
        userRepo:        userRepo,
        movieRepo:       movieRepo,
        moderationSvc:   moderationSvc,
        notificationSvc: notificationSvc,
        cacheStore:      cacheStore,
        logger:          logrus.New(),
        metrics:         NewCommentMetrics(),
    }
}

// 创建评论
func (cs *CommentService) CreateComment(ctx context.Context, req *CommentCreateRequest, userID string) (*CommentResponse, error) {
    start := time.Now()
    defer func() {
        cs.metrics.ObserveCommentOperation("create", time.Since(start))
    }()

    // 验证用户权限
    user, err := cs.userRepo.FindByID(ctx, userID)
    if err != nil {
        cs.metrics.IncInvalidRequests("user_not_found")
        return &CommentResponse{
            Success: false,
            Message: "用户不存在",
        }, nil
    }

    if user.Status != "active" {
        cs.metrics.IncInvalidRequests("user_inactive")
        return &CommentResponse{
            Success: false,
            Message: "用户账户已被禁用",
        }, nil
    }

    // 验证电影存在性
    if exists, err := cs.movieRepo.ExistsByID(ctx, req.MovieID); err != nil {
        cs.logger.Errorf("Failed to check movie existence: %v", err)
        return nil, errors.New("验证电影失败")
    } else if !exists {
        cs.metrics.IncInvalidRequests("movie_not_found")
        return &CommentResponse{
            Success: false,
            Message: "电影不存在",
        }, nil
    }

    // 验证父评论（如果是回复）
    if req.ParentID != "" {
        parentComment, err := cs.commentRepo.FindByID(ctx, req.ParentID)
        if err != nil {
            cs.metrics.IncInvalidRequests("parent_not_found")
            return &CommentResponse{
                Success: false,
                Message: "父评论不存在",
            }, nil
        }

        // 检查父评论是否属于同一电影
        if parentComment.MovieID != req.MovieID {
            cs.metrics.IncInvalidRequests("parent_movie_mismatch")
            return &CommentResponse{
                Success: false,
                Message: "回复评论必须属于同一电影",
            }, nil
        }

        // 限制回复层级（最多2层）
        if parentComment.ParentID != nil {
            cs.metrics.IncInvalidRequests("reply_depth_exceeded")
            return &CommentResponse{
                Success: false,
                Message: "回复层级过深",
            }, nil
        }
    }

    // 检查用户是否已经评论过该电影（如果不是回复）
    if req.ParentID == "" {
        if exists, err := cs.commentRepo.UserHasCommented(ctx, userID, req.MovieID); err != nil {
            cs.logger.Errorf("Failed to check user comment: %v", err)
            return nil, errors.New("检查用户评论失败")
        } else if exists {
            cs.metrics.IncInvalidRequests("already_commented")
            return &CommentResponse{
                Success: false,
                Message: "您已经评论过该电影",
            }, nil
        }
    }

    // 内容审核
    moderationResult, err := cs.moderationSvc.ModerateContent(ctx, req.Content)
    if err != nil {
        cs.logger.Errorf("Content moderation failed: %v", err)
        return nil, errors.New("内容审核失败")
    }

    // 创建评论对象
    comment := &Comment{
        ID:        uuid.New().String(),
        MovieID:   req.MovieID,
        UserID:    userID,
        Content:   req.Content,
        Rating:    req.Rating,
        Spoiler:   req.Spoiler,
        Anonymous: req.Anonymous,
        Status:    cs.determineCommentStatus(moderationResult),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if req.ParentID != "" {
        comment.ParentID = &req.ParentID
    }

    // 保存评论
    if err := cs.commentRepo.Create(ctx, comment); err != nil {
        cs.logger.Errorf("Failed to create comment: %v", err)
        cs.metrics.IncOperationErrors("create")
        return nil, errors.New("评论创建失败")
    }

    // 更新统计信息
    go func() {
        cs.updateCommentStatistics(context.Background(), req.MovieID)
        if req.ParentID != "" {
            cs.updateReplyCount(context.Background(), req.ParentID)
        }
    }()

    // 发送通知
    go func() {
        cs.sendCommentNotifications(context.Background(), comment)
    }()

    // 清除相关缓存
    cs.clearCommentCache(ctx, req.MovieID)

    cs.metrics.IncSuccessfulOperations("create")
    cs.logger.Infof("Comment created: %s by user %s", comment.ID, userID)

    // 加载完整的评论信息
    fullComment, err := cs.loadCommentWithDetails(ctx, comment.ID, userID)
    if err != nil {
        cs.logger.Errorf("Failed to load comment details: %v", err)
        fullComment = comment
    }

    return &CommentResponse{
        Success: true,
        Message: "评论发表成功",
        Data:    fullComment,
    }, nil
}

// 获取评论列表
func (cs *CommentService) GetComments(ctx context.Context, movieID string, page, pageSize int, sortBy string, userID string) (*CommentListResponse, error) {
    start := time.Now()
    defer func() {
        cs.metrics.ObserveCommentOperation("list", time.Since(start))
    }()

    // 验证参数
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }
    if sortBy == "" {
        sortBy = "created_at"
    }

    // 生成缓存键
    cacheKey := cs.generateListCacheKey(movieID, page, pageSize, sortBy, userID)

    // 尝试从缓存获取
    if cachedResult, err := cs.getListFromCache(ctx, cacheKey); err == nil {
        cs.metrics.IncCacheHits()
        return cachedResult, nil
    }
    cs.metrics.IncCacheMisses()

    // 构建查询条件
    queryOptions := &CommentQueryOptions{
        MovieID:  movieID,
        Status:   []string{"approved"},
        ParentID: nil, // 只获取顶级评论
        Page:     page,
        PageSize: pageSize,
        SortBy:   sortBy,
        SortOrder: "desc",
    }

    // 获取评论列表
    comments, total, err := cs.commentRepo.FindWithOptions(ctx, queryOptions)
    if err != nil {
        cs.logger.Errorf("Failed to get comments: %v", err)
        cs.metrics.IncOperationErrors("list")
        return nil, errors.New("获取评论列表失败")
    }

    // 加载评论详细信息
    enrichedComments := make([]Comment, len(comments))
    for i, comment := range comments {
        enrichedComment, err := cs.enrichComment(ctx, &comment, userID)
        if err != nil {
            cs.logger.Errorf("Failed to enrich comment %s: %v", comment.ID, err)
            enrichedComment = &comment
        }
        enrichedComments[i] = *enrichedComment
    }

    // 加载回复（限制层级）
    for i := range enrichedComments {
        replies, err := cs.loadReplies(ctx, enrichedComments[i].ID, userID, 2) // 最多2层回复
        if err != nil {
            cs.logger.Errorf("Failed to load replies for comment %s: %v", enrichedComments[i].ID, err)
        } else {
            enrichedComments[i].Replies = replies
        }
    }

    // 构建分页信息
    pagination := &PaginationInfo{
        CurrentPage:  page,
        PageSize:     pageSize,
        TotalPages:   int(math.Ceil(float64(total) / float64(pageSize))),
        TotalItems:   total,
        HasNext:      page < int(math.Ceil(float64(total)/float64(pageSize))),
        HasPrevious:  page > 1,
    }

    // 获取统计信息
    stats, err := cs.getCommentStatistics(ctx, movieID)
    if err != nil {
        cs.logger.Errorf("Failed to get comment statistics: %v", err)
        // 统计信息获取失败不影响主要功能
    }

    response := &CommentListResponse{
        Success:    true,
        Data:       enrichedComments,
        Pagination: pagination,
        Statistics: stats,
    }

    // 异步缓存结果
    go func() {
        if err := cs.cacheListResult(context.Background(), cacheKey, response); err != nil {
            cs.logger.Errorf("Failed to cache comment list: %v", err)
        }
    }()

    cs.metrics.IncSuccessfulOperations("list")
    return response, nil
}

// 更新评论
func (cs *CommentService) UpdateComment(ctx context.Context, commentID string, req *CommentUpdateRequest, userID string) (*CommentResponse, error) {
    start := time.Now()
    defer func() {
        cs.metrics.ObserveCommentOperation("update", time.Since(start))
    }()

    // 获取现有评论
    comment, err := cs.commentRepo.FindByID(ctx, commentID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            cs.metrics.IncInvalidRequests("comment_not_found")
            return &CommentResponse{
                Success: false,
                Message: "评论不存在",
            }, nil
        }
        cs.logger.Errorf("Failed to get comment: %v", err)
        return nil, errors.New("获取评论失败")
    }

    // 验证权限
    if comment.UserID != userID {
        cs.metrics.IncInvalidRequests("permission_denied")
        return &CommentResponse{
            Success: false,
            Message: "无权限修改此评论",
        }, nil
    }

    // 检查评论状态
    if comment.Status == "deleted" {
        cs.metrics.IncInvalidRequests("comment_deleted")
        return &CommentResponse{
            Success: false,
            Message: "评论已被删除",
        }, nil
    }

    // 检查修改时间限制（24小时内可修改）
    if time.Since(comment.CreatedAt) > 24*time.Hour {
        cs.metrics.IncInvalidRequests("edit_time_expired")
        return &CommentResponse{
            Success: false,
            Message: "评论发表超过24小时，无法修改",
        }, nil
    }

    // 内容审核
    moderationResult, err := cs.moderationSvc.ModerateContent(ctx, req.Content)
    if err != nil {
        cs.logger.Errorf("Content moderation failed: %v", err)
        return nil, errors.New("内容审核失败")
    }

    // 更新评论字段
    comment.Content = req.Content
    comment.Spoiler = req.Spoiler
    comment.Anonymous = req.Anonymous
    comment.Status = cs.determineCommentStatus(moderationResult)
    comment.UpdatedAt = time.Now()

    // 保存更新
    if err := cs.commentRepo.Update(ctx, comment); err != nil {
        cs.logger.Errorf("Failed to update comment: %v", err)
        cs.metrics.IncOperationErrors("update")
        return nil, errors.New("评论更新失败")
    }

    // 清除相关缓存
    cs.clearCommentCache(ctx, comment.MovieID)

    cs.metrics.IncSuccessfulOperations("update")
    cs.logger.Infof("Comment updated: %s by user %s", comment.ID, userID)

    // 加载完整的评论信息
    fullComment, err := cs.loadCommentWithDetails(ctx, comment.ID, userID)
    if err != nil {
        cs.logger.Errorf("Failed to load comment details: %v", err)
        fullComment = comment
    }

    return &CommentResponse{
        Success: true,
        Message: "评论更新成功",
        Data:    fullComment,
    }, nil
}

// 删除评论
func (cs *CommentService) DeleteComment(ctx context.Context, commentID string, userID string, isAdmin bool) (*CommentResponse, error) {
    start := time.Now()
    defer func() {
        cs.metrics.ObserveCommentOperation("delete", time.Since(start))
    }()

    // 获取评论
    comment, err := cs.commentRepo.FindByID(ctx, commentID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            cs.metrics.IncInvalidRequests("comment_not_found")
            return &CommentResponse{
                Success: false,
                Message: "评论不存在",
            }, nil
        }
        return nil, errors.New("获取评论失败")
    }

    // 验证权限
    if !isAdmin && comment.UserID != userID {
        cs.metrics.IncInvalidRequests("permission_denied")
        return &CommentResponse{
            Success: false,
            Message: "无权限删除此评论",
        }, nil
    }

    // 检查是否已删除
    if comment.Status == "deleted" {
        cs.metrics.IncInvalidRequests("already_deleted")
        return &CommentResponse{
            Success: false,
            Message: "评论已被删除",
        }, nil
    }

    // 软删除评论
    comment.Status = "deleted"
    comment.UpdatedAt = time.Now()

    if err := cs.commentRepo.Update(ctx, comment); err != nil {
        cs.logger.Errorf("Failed to delete comment: %v", err)
        cs.metrics.IncOperationErrors("delete")
        return nil, errors.New("评论删除失败")
    }

    // 处理子评论
    if err := cs.handleChildComments(ctx, commentID); err != nil {
        cs.logger.Errorf("Failed to handle child comments: %v", err)
        // 子评论处理失败不影响主评论删除
    }

    // 更新统计信息
    go func() {
        cs.updateCommentStatistics(context.Background(), comment.MovieID)
        if comment.ParentID != nil {
            cs.updateReplyCount(context.Background(), *comment.ParentID)
        }
    }()

    // 清除相关缓存
    cs.clearCommentCache(ctx, comment.MovieID)

    cs.metrics.IncSuccessfulOperations("delete")
    cs.logger.Infof("Comment deleted: %s by user %s (admin: %t)", comment.ID, userID, isAdmin)

    return &CommentResponse{
        Success: true,
        Message: "评论删除成功",
    }, nil
}

// 丰富评论信息
func (cs *CommentService) enrichComment(ctx context.Context, comment *Comment, currentUserID string) (*Comment, error) {
    enriched := *comment

    // 加载用户信息
    if !comment.Anonymous {
        user, err := cs.userRepo.FindByID(ctx, comment.UserID)
        if err != nil {
            cs.logger.Errorf("Failed to load user for comment %s: %v", comment.ID, err)
        } else {
            enriched.User = &CommentUser{
                ID:       user.ID,
                Username: user.Username,
                Avatar:   user.Avatar,
                Level:    user.Level,
            }
        }
    }

    // 加载用户交互状态（如果已登录）
    if currentUserID != "" {
        isLiked, err := cs.commentRepo.IsLikedByUser(ctx, comment.ID, currentUserID)
        if err != nil {
            cs.logger.Errorf("Failed to check like status: %v", err)
        } else {
            enriched.IsLiked = isLiked
        }

        isReported, err := cs.commentRepo.IsReportedByUser(ctx, comment.ID, currentUserID)
        if err != nil {
            cs.logger.Errorf("Failed to check report status: %v", err)
        } else {
            enriched.IsReported = isReported
        }
    }

    return &enriched, nil
}
```

## 📊 性能监控

### 1. **评论操作指标**

```go
type CommentMetrics struct {
    operationCount    *prometheus.CounterVec
    operationDuration *prometheus.HistogramVec
    cacheHitRate      *prometheus.CounterVec
    commentCount      prometheus.Gauge
    errorCount        *prometheus.CounterVec
}

func NewCommentMetrics() *CommentMetrics {
    return &CommentMetrics{
        operationCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "comment_operations_total",
                Help: "Total number of comment operations",
            },
            []string{"operation", "status"},
        ),
        operationDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "comment_operation_duration_seconds",
                Help: "Duration of comment operations",
            },
            []string{"operation"},
        ),
        cacheHitRate: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "comment_cache_operations_total",
                Help: "Total number of comment cache operations",
            },
            []string{"type"},
        ),
        commentCount: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "comments_total",
                Help: "Total number of comments",
            },
        ),
        errorCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "comment_errors_total",
                Help: "Total number of comment errors",
            },
            []string{"operation", "error_type"},
        ),
    }
}
```

## 🔧 HTTP处理器

### 1. **评论API端点**

```go
func (cc *CommentController) CreateComment(c *gin.Context) {
    var req CommentCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": "请求参数错误",
            "error":   err.Error(),
        })
        return
    }

    userID := cc.getUserIDFromContext(c)
    if userID == "" {
        c.JSON(401, gin.H{
            "success": false,
            "message": "请先登录",
        })
        return
    }

    response, err := cc.commentService.CreateComment(c.Request.Context(), &req, userID)
    if err != nil {
        cc.logger.Errorf("Failed to create comment: %v", err)
        c.JSON(500, gin.H{
            "success": false,
            "message": "评论创建失败",
        })
        return
    }

    c.JSON(200, response)
}

func (cc *CommentController) GetComments(c *gin.Context) {
    movieID := c.Param("movie_id")
    if movieID == "" {
        c.JSON(400, gin.H{
            "success": false,
            "message": "电影ID不能为空",
        })
        return
    }

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    sortBy := c.DefaultQuery("sort_by", "created_at")
    userID := cc.getUserIDFromContext(c)

    response, err := cc.commentService.GetComments(c.Request.Context(), movieID, page, pageSize, sortBy, userID)
    if err != nil {
        cc.logger.Errorf("Failed to get comments: %v", err)
        c.JSON(500, gin.H{
            "success": false,
            "message": "获取评论失败",
        })
        return
    }

    c.Header("Cache-Control", "public, max-age=300") // 5分钟缓存
    c.JSON(200, response)
}
```

## 📝 总结

评论CRUD接口为MovieInfo项目提供了完整的评论管理功能：

**核心功能**：
1. **完整CRUD**：评论的创建、读取、更新、删除操作
2. **层级管理**：支持评论回复的层级结构
3. **权限控制**：严格的用户权限验证机制
4. **内容审核**：自动化的内容安全检查

**性能优化**：
- 高效的分页查询
- 智能缓存策略
- 异步统计更新
- 数据库查询优化

**安全保障**：
- 用户权限验证
- 内容安全审核
- 操作时间限制
- 防刷机制

下一步，我们将实现评分系统，为用户提供电影评分功能。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第38步：评分系统实现](38-rating-system.md)
