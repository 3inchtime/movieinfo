# 第40步：统计分析功能

## 📋 概述

统计分析功能为MovieInfo项目提供全面的数据洞察和运营支持。通过收集、处理和分析用户行为、内容质量、平台健康度等多维度数据，为产品优化和运营决策提供科学依据。

## 🎯 设计目标

### 1. **数据完整性**
- 全面的数据收集
- 多维度指标体系
- 实时数据更新
- 历史数据追踪

### 2. **分析深度**
- 用户行为分析
- 内容质量分析
- 趋势预测分析
- 异常检测分析

### 3. **可视化展示**
- 直观的图表展示
- 交互式数据探索
- 自定义报表生成
- 移动端适配

## 🔧 统计分析架构

### 1. **数据模型设计**

```go
// 用户行为统计
type UserBehaviorStats struct {
    ID              string    `gorm:"primaryKey" json:"id"`
    UserID          string    `gorm:"not null;index" json:"user_id"`
    Date            time.Time `gorm:"not null;index" json:"date"`
    
    // 活跃度指标
    LoginCount      int       `gorm:"not null;default:0" json:"login_count"`
    PageViews       int       `gorm:"not null;default:0" json:"page_views"`
    SessionDuration int       `gorm:"not null;default:0" json:"session_duration"` // 秒
    
    // 内容互动
    CommentsPosted  int       `gorm:"not null;default:0" json:"comments_posted"`
    RatingsGiven    int       `gorm:"not null;default:0" json:"ratings_given"`
    LikesGiven      int       `gorm:"not null;default:0" json:"likes_given"`
    SharesCount     int       `gorm:"not null;default:0" json:"shares_count"`
    
    // 搜索行为
    SearchQueries   int       `gorm:"not null;default:0" json:"search_queries"`
    MoviesViewed    int       `gorm:"not null;default:0" json:"movies_viewed"`
    
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

// 内容统计
type ContentStats struct {
    ID              string    `gorm:"primaryKey" json:"id"`
    ContentID       string    `gorm:"not null;index" json:"content_id"`
    ContentType     string    `gorm:"not null;size:50" json:"content_type"` // movie, comment, rating
    Date            time.Time `gorm:"not null;index" json:"date"`
    
    // 浏览统计
    ViewCount       int       `gorm:"not null;default:0" json:"view_count"`
    UniqueViews     int       `gorm:"not null;default:0" json:"unique_views"`
    
    // 互动统计
    LikeCount       int       `gorm:"not null;default:0" json:"like_count"`
    CommentCount    int       `gorm:"not null;default:0" json:"comment_count"`
    ShareCount      int       `gorm:"not null;default:0" json:"share_count"`
    
    // 评分统计
    RatingCount     int       `gorm:"not null;default:0" json:"rating_count"`
    AverageRating   float64   `gorm:"not null;default:0" json:"average_rating"`
    
    // 质量指标
    QualityScore    float64   `gorm:"not null;default:0" json:"quality_score"`
    EngagementRate  float64   `gorm:"not null;default:0" json:"engagement_rate"`
    
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

// 平台统计
type PlatformStats struct {
    ID              string    `gorm:"primaryKey" json:"id"`
    Date            time.Time `gorm:"not null;uniqueIndex" json:"date"`
    
    // 用户指标
    TotalUsers      int       `gorm:"not null;default:0" json:"total_users"`
    ActiveUsers     int       `gorm:"not null;default:0" json:"active_users"`
    NewUsers        int       `gorm:"not null;default:0" json:"new_users"`
    RetentionRate   float64   `gorm:"not null;default:0" json:"retention_rate"`
    
    // 内容指标
    TotalMovies     int       `gorm:"not null;default:0" json:"total_movies"`
    TotalComments   int       `gorm:"not null;default:0" json:"total_comments"`
    TotalRatings    int       `gorm:"not null;default:0" json:"total_ratings"`
    
    // 活跃度指标
    PageViews       int       `gorm:"not null;default:0" json:"page_views"`
    Sessions        int       `gorm:"not null;default:0" json:"sessions"`
    AvgSessionTime  int       `gorm:"not null;default:0" json:"avg_session_time"`
    
    // 质量指标
    ContentQuality  float64   `gorm:"not null;default:0" json:"content_quality"`
    UserSatisfaction float64  `gorm:"not null;default:0" json:"user_satisfaction"`
    
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

// 统计查询请求
type AnalyticsRequest struct {
    MetricType  string    `json:"metric_type" binding:"required"` // user, content, platform
    StartDate   time.Time `json:"start_date" binding:"required"`
    EndDate     time.Time `json:"end_date" binding:"required"`
    Granularity string    `json:"granularity"` // hour, day, week, month
    Filters     map[string]interface{} `json:"filters,omitempty"`
    GroupBy     []string  `json:"group_by,omitempty"`
}

// 统计响应
type AnalyticsResponse struct {
    Success   bool                   `json:"success"`
    Message   string                 `json:"message,omitempty"`
    Data      []AnalyticsDataPoint   `json:"data,omitempty"`
    Summary   *AnalyticsSummary      `json:"summary,omitempty"`
    Metadata  *AnalyticsMetadata     `json:"metadata,omitempty"`
}

type AnalyticsDataPoint struct {
    Timestamp time.Time              `json:"timestamp"`
    Metrics   map[string]interface{} `json:"metrics"`
    Dimensions map[string]string     `json:"dimensions,omitempty"`
}

type AnalyticsSummary struct {
    TotalRecords int                    `json:"total_records"`
    Aggregations map[string]interface{} `json:"aggregations"`
    Trends       map[string]float64     `json:"trends"`
}

type AnalyticsMetadata struct {
    QueryTime    time.Duration `json:"query_time_ms"`
    CacheHit     bool          `json:"cache_hit"`
    DataSource   string        `json:"data_source"`
    LastUpdated  time.Time     `json:"last_updated"`
}
```

### 2. **统计分析服务**

```go
type AnalyticsService struct {
    userStatsRepo     UserBehaviorStatsRepository
    contentStatsRepo  ContentStatsRepository
    platformStatsRepo PlatformStatsRepository
    cacheStore        CacheStore
    timeSeriesDB      TimeSeriesDB
    logger            *logrus.Logger
    metrics           *AnalyticsMetrics
}

func NewAnalyticsService(
    userStatsRepo UserBehaviorStatsRepository,
    contentStatsRepo ContentStatsRepository,
    platformStatsRepo PlatformStatsRepository,
    cacheStore CacheStore,
    timeSeriesDB TimeSeriesDB,
) *AnalyticsService {
    return &AnalyticsService{
        userStatsRepo:     userStatsRepo,
        contentStatsRepo:  contentStatsRepo,
        platformStatsRepo: platformStatsRepo,
        cacheStore:        cacheStore,
        timeSeriesDB:      timeSeriesDB,
        logger:            logrus.New(),
        metrics:           NewAnalyticsMetrics(),
    }
}

// 获取分析数据
func (as *AnalyticsService) GetAnalytics(ctx context.Context, req *AnalyticsRequest) (*AnalyticsResponse, error) {
    start := time.Now()
    defer func() {
        as.metrics.ObserveQueryDuration(req.MetricType, time.Since(start))
    }()

    // 验证请求参数
    if err := as.validateRequest(req); err != nil {
        as.metrics.IncInvalidRequests()
        return &AnalyticsResponse{
            Success: false,
            Message: err.Error(),
        }, nil
    }

    // 生成缓存键
    cacheKey := as.generateCacheKey(req)

    // 尝试从缓存获取
    if cachedResult, err := as.getFromCache(ctx, cacheKey); err == nil {
        as.metrics.IncCacheHits()
        return cachedResult, nil
    }
    as.metrics.IncCacheMisses()

    // 根据指标类型执行查询
    var data []AnalyticsDataPoint
    var err error

    switch req.MetricType {
    case "user":
        data, err = as.getUserAnalytics(ctx, req)
    case "content":
        data, err = as.getContentAnalytics(ctx, req)
    case "platform":
        data, err = as.getPlatformAnalytics(ctx, req)
    default:
        as.metrics.IncInvalidRequests()
        return &AnalyticsResponse{
            Success: false,
            Message: "不支持的指标类型",
        }, nil
    }

    if err != nil {
        as.logger.Errorf("Failed to get analytics data: %v", err)
        as.metrics.IncQueryErrors()
        return nil, errors.New("获取分析数据失败")
    }

    // 计算汇总信息
    summary := as.calculateSummary(data, req)

    // 构建响应
    response := &AnalyticsResponse{
        Success: true,
        Data:    data,
        Summary: summary,
        Metadata: &AnalyticsMetadata{
            QueryTime:   time.Since(start),
            CacheHit:    false,
            DataSource:  "database",
            LastUpdated: time.Now(),
        },
    }

    // 异步缓存结果
    go func() {
        if err := as.cacheResult(context.Background(), cacheKey, response); err != nil {
            as.logger.Errorf("Failed to cache analytics result: %v", err)
        }
    }()

    as.metrics.IncSuccessfulQueries()
    return response, nil
}

// 获取用户行为分析
func (as *AnalyticsService) getUserAnalytics(ctx context.Context, req *AnalyticsRequest) ([]AnalyticsDataPoint, error) {
    // 构建查询条件
    queryBuilder := as.buildUserStatsQuery(req)

    // 执行查询
    stats, err := as.userStatsRepo.FindWithQuery(ctx, queryBuilder)
    if err != nil {
        return nil, err
    }

    // 转换为数据点
    dataPoints := make([]AnalyticsDataPoint, 0)
    
    // 按时间分组聚合
    groupedData := as.groupUserStatsByTime(stats, req.Granularity)
    
    for timestamp, groupStats := range groupedData {
        metrics := map[string]interface{}{
            "active_users":     len(groupStats),
            "total_page_views": as.sumUserMetric(groupStats, "page_views"),
            "total_comments":   as.sumUserMetric(groupStats, "comments_posted"),
            "total_ratings":    as.sumUserMetric(groupStats, "ratings_given"),
            "avg_session_time": as.avgUserMetric(groupStats, "session_duration"),
        }

        dataPoints = append(dataPoints, AnalyticsDataPoint{
            Timestamp: timestamp,
            Metrics:   metrics,
        })
    }

    // 按时间排序
    sort.Slice(dataPoints, func(i, j int) bool {
        return dataPoints[i].Timestamp.Before(dataPoints[j].Timestamp)
    })

    return dataPoints, nil
}

// 获取内容分析
func (as *AnalyticsService) getContentAnalytics(ctx context.Context, req *AnalyticsRequest) ([]AnalyticsDataPoint, error) {
    queryBuilder := as.buildContentStatsQuery(req)

    stats, err := as.contentStatsRepo.FindWithQuery(ctx, queryBuilder)
    if err != nil {
        return nil, err
    }

    dataPoints := make([]AnalyticsDataPoint, 0)
    groupedData := as.groupContentStatsByTime(stats, req.Granularity)

    for timestamp, groupStats := range groupedData {
        metrics := map[string]interface{}{
            "total_views":      as.sumContentMetric(groupStats, "view_count"),
            "unique_views":     as.sumContentMetric(groupStats, "unique_views"),
            "total_likes":      as.sumContentMetric(groupStats, "like_count"),
            "total_comments":   as.sumContentMetric(groupStats, "comment_count"),
            "avg_rating":       as.avgContentMetric(groupStats, "average_rating"),
            "engagement_rate":  as.avgContentMetric(groupStats, "engagement_rate"),
        }

        dataPoints = append(dataPoints, AnalyticsDataPoint{
            Timestamp: timestamp,
            Metrics:   metrics,
        })
    }

    sort.Slice(dataPoints, func(i, j int) bool {
        return dataPoints[i].Timestamp.Before(dataPoints[j].Timestamp)
    })

    return dataPoints, nil
}

// 获取平台分析
func (as *AnalyticsService) getPlatformAnalytics(ctx context.Context, req *AnalyticsRequest) ([]AnalyticsDataPoint, error) {
    stats, err := as.platformStatsRepo.FindByDateRange(ctx, req.StartDate, req.EndDate)
    if err != nil {
        return nil, err
    }

    dataPoints := make([]AnalyticsDataPoint, len(stats))
    for i, stat := range stats {
        dataPoints[i] = AnalyticsDataPoint{
            Timestamp: stat.Date,
            Metrics: map[string]interface{}{
                "total_users":       stat.TotalUsers,
                "active_users":      stat.ActiveUsers,
                "new_users":         stat.NewUsers,
                "retention_rate":    stat.RetentionRate,
                "total_movies":      stat.TotalMovies,
                "total_comments":    stat.TotalComments,
                "total_ratings":     stat.TotalRatings,
                "page_views":        stat.PageViews,
                "sessions":          stat.Sessions,
                "avg_session_time":  stat.AvgSessionTime,
                "content_quality":   stat.ContentQuality,
                "user_satisfaction": stat.UserSatisfaction,
            },
        }
    }

    return dataPoints, nil
}

// 计算汇总信息
func (as *AnalyticsService) calculateSummary(data []AnalyticsDataPoint, req *AnalyticsRequest) *AnalyticsSummary {
    if len(data) == 0 {
        return &AnalyticsSummary{
            TotalRecords: 0,
            Aggregations: make(map[string]interface{}),
            Trends:       make(map[string]float64),
        }
    }

    aggregations := make(map[string]interface{})
    trends := make(map[string]float64)

    // 计算聚合值
    for metricName := range data[0].Metrics {
        values := make([]float64, len(data))
        for i, point := range data {
            if val, ok := point.Metrics[metricName].(float64); ok {
                values[i] = val
            } else if val, ok := point.Metrics[metricName].(int); ok {
                values[i] = float64(val)
            }
        }

        aggregations[metricName+"_sum"] = as.sum(values)
        aggregations[metricName+"_avg"] = as.average(values)
        aggregations[metricName+"_max"] = as.max(values)
        aggregations[metricName+"_min"] = as.min(values)

        // 计算趋势（简单线性回归斜率）
        trends[metricName] = as.calculateTrend(values)
    }

    return &AnalyticsSummary{
        TotalRecords: len(data),
        Aggregations: aggregations,
        Trends:       trends,
    }
}

// 实时数据收集
func (as *AnalyticsService) RecordUserBehavior(ctx context.Context, userID string, action string, metadata map[string]interface{}) error {
    // 获取或创建今日统计记录
    today := time.Now().Truncate(24 * time.Hour)
    stats, err := as.userStatsRepo.FindByUserAndDate(ctx, userID, today)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return err
    }

    if stats == nil {
        stats = &UserBehaviorStats{
            ID:     uuid.New().String(),
            UserID: userID,
            Date:   today,
        }
    }

    // 根据行为类型更新统计
    switch action {
    case "login":
        stats.LoginCount++
    case "page_view":
        stats.PageViews++
        if duration, ok := metadata["session_duration"].(int); ok {
            stats.SessionDuration += duration
        }
    case "comment_posted":
        stats.CommentsPosted++
    case "rating_given":
        stats.RatingsGiven++
    case "like_given":
        stats.LikesGiven++
    case "share":
        stats.SharesCount++
    case "search":
        stats.SearchQueries++
    case "movie_viewed":
        stats.MoviesViewed++
    }

    stats.UpdatedAt = time.Now()

    // 保存或更新统计
    if stats.CreatedAt.IsZero() {
        stats.CreatedAt = time.Now()
        return as.userStatsRepo.Create(ctx, stats)
    } else {
        return as.userStatsRepo.Update(ctx, stats)
    }
}

// 批量统计更新
func (as *AnalyticsService) UpdateDailyStats(ctx context.Context) error {
    today := time.Now().Truncate(24 * time.Hour)
    
    // 更新平台统计
    if err := as.updatePlatformStats(ctx, today); err != nil {
        as.logger.Errorf("Failed to update platform stats: %v", err)
        return err
    }

    // 更新内容统计
    if err := as.updateContentStats(ctx, today); err != nil {
        as.logger.Errorf("Failed to update content stats: %v", err)
        return err
    }

    as.logger.Infof("Daily stats updated for %s", today.Format("2006-01-02"))
    return nil
}

// 数据导出
func (as *AnalyticsService) ExportData(ctx context.Context, req *AnalyticsRequest, format string) ([]byte, error) {
    // 获取数据
    response, err := as.GetAnalytics(ctx, req)
    if err != nil {
        return nil, err
    }

    switch format {
    case "csv":
        return as.exportToCSV(response.Data)
    case "json":
        return json.Marshal(response)
    case "excel":
        return as.exportToExcel(response.Data)
    default:
        return nil, errors.New("不支持的导出格式")
    }
}
```

### 3. **实时数据收集**

```go
type EventCollector struct {
    analyticsService *AnalyticsService
    eventQueue       chan *Event
    batchSize        int
    flushInterval    time.Duration
    logger           *logrus.Logger
}

type Event struct {
    UserID    string                 `json:"user_id"`
    Action    string                 `json:"action"`
    Timestamp time.Time              `json:"timestamp"`
    Metadata  map[string]interface{} `json:"metadata"`
}

func NewEventCollector(analyticsService *AnalyticsService) *EventCollector {
    ec := &EventCollector{
        analyticsService: analyticsService,
        eventQueue:       make(chan *Event, 10000),
        batchSize:        100,
        flushInterval:    10 * time.Second,
        logger:           logrus.New(),
    }

    // 启动事件处理器
    go ec.startEventProcessor()

    return ec
}

// 收集事件
func (ec *EventCollector) CollectEvent(userID, action string, metadata map[string]interface{}) {
    event := &Event{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now(),
        Metadata:  metadata,
    }

    select {
    case ec.eventQueue <- event:
        // 事件已入队
    default:
        // 队列满，丢弃事件
        ec.logger.Warn("Event queue full, dropping event")
    }
}

// 事件处理器
func (ec *EventCollector) startEventProcessor() {
    ticker := time.NewTicker(ec.flushInterval)
    defer ticker.Stop()

    var batch []*Event

    for {
        select {
        case event := <-ec.eventQueue:
            batch = append(batch, event)
            if len(batch) >= ec.batchSize {
                ec.processBatch(batch)
                batch = batch[:0]
            }

        case <-ticker.C:
            if len(batch) > 0 {
                ec.processBatch(batch)
                batch = batch[:0]
            }
        }
    }
}

// 批量处理事件
func (ec *EventCollector) processBatch(events []*Event) {
    ctx := context.Background()

    for _, event := range events {
        if err := ec.analyticsService.RecordUserBehavior(ctx, event.UserID, event.Action, event.Metadata); err != nil {
            ec.logger.Errorf("Failed to record user behavior: %v", err)
        }
    }

    ec.logger.Debugf("Processed %d events", len(events))
}
```

## 📊 监控指标

### 1. **分析系统指标**

```go
type AnalyticsMetrics struct {
    queryCount     *prometheus.CounterVec
    queryDuration  *prometheus.HistogramVec
    cacheHitRate   *prometheus.CounterVec
    dataPoints     prometheus.Gauge
    eventCount     *prometheus.CounterVec
}

func NewAnalyticsMetrics() *AnalyticsMetrics {
    return &AnalyticsMetrics{
        queryCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "analytics_queries_total",
                Help: "Total number of analytics queries",
            },
            []string{"metric_type", "status"},
        ),
        queryDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "analytics_query_duration_seconds",
                Help: "Duration of analytics queries",
            },
            []string{"metric_type"},
        ),
        cacheHitRate: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "analytics_cache_operations_total",
                Help: "Total number of analytics cache operations",
            },
            []string{"type"},
        ),
        dataPoints: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "analytics_data_points_total",
                Help: "Total number of data points in analytics",
            },
        ),
        eventCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "analytics_events_total",
                Help: "Total number of collected events",
            },
            []string{"action"},
        ),
    }
}
```

## 🔧 HTTP处理器

### 1. **分析API端点**

```go
func (ac *AnalyticsController) GetAnalytics(c *gin.Context) {
    var req AnalyticsRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": "请求参数错误",
            "error":   err.Error(),
        })
        return
    }

    // 验证权限
    userID := ac.getUserIDFromContext(c)
    if !ac.hasAnalyticsPermission(userID) {
        c.JSON(403, gin.H{
            "success": false,
            "message": "无权限访问分析数据",
        })
        return
    }

    response, err := ac.analyticsService.GetAnalytics(c.Request.Context(), &req)
    if err != nil {
        ac.logger.Errorf("Failed to get analytics: %v", err)
        c.JSON(500, gin.H{
            "success": false,
            "message": "获取分析数据失败",
        })
        return
    }

    c.Header("Cache-Control", "private, max-age=300") // 5分钟缓存
    c.JSON(200, response)
}

func (ac *AnalyticsController) ExportData(c *gin.Context) {
    var req AnalyticsRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": "请求参数错误",
        })
        return
    }

    format := c.DefaultQuery("format", "csv")
    
    data, err := ac.analyticsService.ExportData(c.Request.Context(), &req, format)
    if err != nil {
        ac.logger.Errorf("Failed to export data: %v", err)
        c.JSON(500, gin.H{
            "success": false,
            "message": "数据导出失败",
        })
        return
    }

    filename := fmt.Sprintf("analytics_%s.%s", time.Now().Format("20060102"), format)
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
    
    switch format {
    case "csv":
        c.Header("Content-Type", "text/csv")
    case "json":
        c.Header("Content-Type", "application/json")
    case "excel":
        c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    }

    c.Data(200, c.GetHeader("Content-Type"), data)
}
```

## 📝 总结

统计分析功能为MovieInfo项目提供了全面的数据洞察能力：

**核心功能**：
1. **多维分析**：用户行为、内容质量、平台健康度分析
2. **实时收集**：事件驱动的实时数据收集机制
3. **趋势分析**：历史数据趋势和预测分析
4. **数据导出**：支持多种格式的数据导出

**技术特性**：
- 高性能的时序数据处理
- 智能的数据聚合算法
- 灵活的查询和过滤机制
- 完善的缓存优化策略

**业务价值**：
- 用户行为洞察
- 内容质量监控
- 运营决策支持
- 产品优化指导

至此，评论打分服务的核心功能已经完成。下一步，我们将继续完成主页服务的开发文档。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第41步：Web服务器搭建](../09-web-service/41-web-server.md)
