# 第32步：电影列表接口

## 📋 概述

电影列表接口是MovieInfo项目的核心展示功能，为用户提供丰富的电影浏览体验。一个高效的电影列表接口需要支持多种排序方式、筛选条件、分页机制，并具备优秀的性能表现。

## 🎯 设计目标

### 1. **功能完整性**
- 多种排序方式
- 灵活的筛选条件
- 高效的分页机制
- 丰富的数据展示

### 2. **性能优化**
- 快速响应时间
- 智能缓存策略
- 数据库查询优化
- 并发处理能力

### 3. **用户体验**
- 流畅的浏览体验
- 直观的数据展示
- 响应式设计支持
- 加载状态反馈

## 🔧 接口设计

### 1. **电影列表请求结构**

```go
// 电影列表请求参数
type MovieListRequest struct {
    // 分页参数
    Page     int `json:"page" form:"page" binding:"min=1"`
    PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
    
    // 排序参数
    SortBy    string `json:"sort_by" form:"sort_by"`     // rating, release_date, popularity, title
    SortOrder string `json:"sort_order" form:"sort_order"` // asc, desc
    
    // 筛选参数
    Genre       []string `json:"genre" form:"genre"`           // 类型筛选
    Year        int      `json:"year" form:"year"`             // 年份筛选
    YearFrom    int      `json:"year_from" form:"year_from"`   // 年份范围开始
    YearTo      int      `json:"year_to" form:"year_to"`       // 年份范围结束
    Rating      float64  `json:"rating" form:"rating"`         // 最低评分
    Country     string   `json:"country" form:"country"`       // 国家/地区
    Language    string   `json:"language" form:"language"`     // 语言
    
    // 搜索参数
    Keyword string `json:"keyword" form:"keyword"` // 关键词搜索
    
    // 其他参数
    IncludeAdult bool `json:"include_adult" form:"include_adult"` // 是否包含成人内容
}

// 电影列表响应结构
type MovieListResponse struct {
    Success    bool              `json:"success"`
    Message    string            `json:"message,omitempty"`
    Data       *MovieListData    `json:"data,omitempty"`
    Pagination *PaginationInfo   `json:"pagination,omitempty"`
    Filters    *FilterInfo       `json:"filters,omitempty"`
}

// 电影列表数据
type MovieListData struct {
    Movies []MovieListItem `json:"movies"`
    Total  int64           `json:"total"`
}

// 电影列表项
type MovieListItem struct {
    ID           string    `json:"id"`
    Title        string    `json:"title"`
    OriginalTitle string   `json:"original_title,omitempty"`
    PosterURL    string    `json:"poster_url"`
    BackdropURL  string    `json:"backdrop_url,omitempty"`
    Overview     string    `json:"overview"`
    ReleaseDate  string    `json:"release_date"`
    Rating       float64   `json:"rating"`
    VoteCount    int       `json:"vote_count"`
    Popularity   float64   `json:"popularity"`
    Genres       []Genre   `json:"genres"`
    Runtime      int       `json:"runtime,omitempty"`
    Country      string    `json:"country,omitempty"`
    Language     string    `json:"language,omitempty"`
    Adult        bool      `json:"adult"`
}

// 分页信息
type PaginationInfo struct {
    CurrentPage  int   `json:"current_page"`
    PageSize     int   `json:"page_size"`
    TotalPages   int   `json:"total_pages"`
    TotalItems   int64 `json:"total_items"`
    HasNext      bool  `json:"has_next"`
    HasPrevious  bool  `json:"has_previous"`
}

// 筛选信息
type FilterInfo struct {
    AvailableGenres    []Genre           `json:"available_genres"`
    AvailableCountries []Country         `json:"available_countries"`
    AvailableLanguages []Language        `json:"available_languages"`
    YearRange          *YearRange        `json:"year_range"`
    RatingRange        *RatingRange      `json:"rating_range"`
}
```

### 2. **电影列表服务实现**

```go
type MovieListService struct {
    movieRepo   MovieRepository
    cacheStore  CacheStore
    logger      *logrus.Logger
    metrics     *MovieListMetrics
}

func NewMovieListService(
    movieRepo MovieRepository,
    cacheStore CacheStore,
) *MovieListService {
    return &MovieListService{
        movieRepo:  movieRepo,
        cacheStore: cacheStore,
        logger:     logrus.New(),
        metrics:    NewMovieListMetrics(),
    }
}

// 获取电影列表
func (mls *MovieListService) GetMovieList(ctx context.Context, req *MovieListRequest) (*MovieListResponse, error) {
    start := time.Now()
    defer func() {
        mls.metrics.ObserveListRequestDuration(time.Since(start))
    }()

    // 验证请求参数
    if err := mls.validateRequest(req); err != nil {
        mls.metrics.IncInvalidRequests()
        return &MovieListResponse{
            Success: false,
            Message: err.Error(),
        }, nil
    }

    // 设置默认值
    mls.setDefaults(req)

    // 生成缓存键
    cacheKey := mls.generateCacheKey(req)

    // 尝试从缓存获取
    if cachedResult, err := mls.getFromCache(ctx, cacheKey); err == nil {
        mls.metrics.IncCacheHits()
        return cachedResult, nil
    }
    mls.metrics.IncCacheMisses()

    // 构建查询条件
    queryBuilder := mls.buildQuery(req)

    // 获取总数
    total, err := mls.movieRepo.CountWithQuery(ctx, queryBuilder)
    if err != nil {
        mls.logger.Errorf("Failed to count movies: %v", err)
        mls.metrics.IncQueryErrors()
        return nil, errors.New("获取电影数量失败")
    }

    // 获取电影列表
    movies, err := mls.movieRepo.FindWithQuery(ctx, queryBuilder)
    if err != nil {
        mls.logger.Errorf("Failed to get movies: %v", err)
        mls.metrics.IncQueryErrors()
        return nil, errors.New("获取电影列表失败")
    }

    // 转换为响应格式
    movieItems := mls.convertToListItems(movies)

    // 构建分页信息
    pagination := mls.buildPaginationInfo(req, total)

    // 获取筛选信息
    filters, err := mls.getFilterInfo(ctx)
    if err != nil {
        mls.logger.Errorf("Failed to get filter info: %v", err)
        // 筛选信息获取失败不影响主要功能
    }

    // 构建响应
    response := &MovieListResponse{
        Success: true,
        Data: &MovieListData{
            Movies: movieItems,
            Total:  total,
        },
        Pagination: pagination,
        Filters:    filters,
    }

    // 缓存结果
    go func() {
        if err := mls.cacheResult(context.Background(), cacheKey, response); err != nil {
            mls.logger.Errorf("Failed to cache result: %v", err)
        }
    }()

    mls.metrics.IncSuccessfulRequests()
    return response, nil
}

// 验证请求参数
func (mls *MovieListService) validateRequest(req *MovieListRequest) error {
    if req.Page < 1 {
        req.Page = 1
    }
    
    if req.PageSize < 1 || req.PageSize > 100 {
        req.PageSize = 20
    }

    // 验证排序字段
    validSortFields := map[string]bool{
        "rating":       true,
        "release_date": true,
        "popularity":   true,
        "title":        true,
        "created_at":   true,
    }
    
    if req.SortBy != "" && !validSortFields[req.SortBy] {
        return errors.New("无效的排序字段")
    }

    // 验证排序顺序
    if req.SortOrder != "" && req.SortOrder != "asc" && req.SortOrder != "desc" {
        return errors.New("无效的排序顺序")
    }

    // 验证年份范围
    if req.YearFrom > 0 && req.YearTo > 0 && req.YearFrom > req.YearTo {
        return errors.New("年份范围无效")
    }

    // 验证评分范围
    if req.Rating < 0 || req.Rating > 10 {
        return errors.New("评分范围无效")
    }

    return nil
}

// 设置默认值
func (mls *MovieListService) setDefaults(req *MovieListRequest) {
    if req.Page == 0 {
        req.Page = 1
    }
    
    if req.PageSize == 0 {
        req.PageSize = 20
    }
    
    if req.SortBy == "" {
        req.SortBy = "popularity"
    }
    
    if req.SortOrder == "" {
        req.SortOrder = "desc"
    }
}

// 构建查询条件
func (mls *MovieListService) buildQuery(req *MovieListRequest) *QueryBuilder {
    qb := NewQueryBuilder()

    // 基础条件
    qb.Where("status = ?", "published")

    // 成人内容筛选
    if !req.IncludeAdult {
        qb.Where("adult = ?", false)
    }

    // 类型筛选
    if len(req.Genre) > 0 {
        qb.WhereIn("genres.name", req.Genre).
           Joins("JOIN movie_genres ON movies.id = movie_genres.movie_id").
           Joins("JOIN genres ON movie_genres.genre_id = genres.id")
    }

    // 年份筛选
    if req.Year > 0 {
        qb.Where("YEAR(release_date) = ?", req.Year)
    } else {
        if req.YearFrom > 0 {
            qb.Where("YEAR(release_date) >= ?", req.YearFrom)
        }
        if req.YearTo > 0 {
            qb.Where("YEAR(release_date) <= ?", req.YearTo)
        }
    }

    // 评分筛选
    if req.Rating > 0 {
        qb.Where("rating >= ?", req.Rating)
    }

    // 国家筛选
    if req.Country != "" {
        qb.Where("country = ?", req.Country)
    }

    // 语言筛选
    if req.Language != "" {
        qb.Where("original_language = ?", req.Language)
    }

    // 关键词搜索
    if req.Keyword != "" {
        keyword := "%" + req.Keyword + "%"
        qb.Where("(title LIKE ? OR original_title LIKE ? OR overview LIKE ?)", 
                 keyword, keyword, keyword)
    }

    // 排序
    orderClause := fmt.Sprintf("%s %s", req.SortBy, strings.ToUpper(req.SortOrder))
    qb.Order(orderClause)

    // 分页
    offset := (req.Page - 1) * req.PageSize
    qb.Limit(req.PageSize).Offset(offset)

    // 预加载关联数据
    qb.Preload("Genres").Preload("Countries").Preload("Languages")

    return qb
}

// 转换为列表项
func (mls *MovieListService) convertToListItems(movies []*Movie) []MovieListItem {
    items := make([]MovieListItem, len(movies))
    
    for i, movie := range movies {
        items[i] = MovieListItem{
            ID:            movie.ID,
            Title:         movie.Title,
            OriginalTitle: movie.OriginalTitle,
            PosterURL:     movie.PosterURL,
            BackdropURL:   movie.BackdropURL,
            Overview:      mls.truncateOverview(movie.Overview, 200),
            ReleaseDate:   movie.ReleaseDate.Format("2006-01-02"),
            Rating:        movie.Rating,
            VoteCount:     movie.VoteCount,
            Popularity:    movie.Popularity,
            Genres:        mls.convertGenres(movie.Genres),
            Runtime:       movie.Runtime,
            Country:       movie.Country,
            Language:      movie.OriginalLanguage,
            Adult:         movie.Adult,
        }
    }
    
    return items
}

// 构建分页信息
func (mls *MovieListService) buildPaginationInfo(req *MovieListRequest, total int64) *PaginationInfo {
    totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))
    
    return &PaginationInfo{
        CurrentPage:  req.Page,
        PageSize:     req.PageSize,
        TotalPages:   totalPages,
        TotalItems:   total,
        HasNext:      req.Page < totalPages,
        HasPrevious:  req.Page > 1,
    }
}
```

### 3. **缓存策略实现**

```go
// 缓存管理器
type MovieListCacheManager struct {
    redis      *redis.Client
    localCache *cache.Cache
    logger     *logrus.Logger
}

func NewMovieListCacheManager(redis *redis.Client) *MovieListCacheManager {
    return &MovieListCacheManager{
        redis:      redis,
        localCache: cache.New(5*time.Minute, 10*time.Minute),
        logger:     logrus.New(),
    }
}

// 生成缓存键
func (mls *MovieListService) generateCacheKey(req *MovieListRequest) string {
    h := sha256.New()
    
    // 将请求参数序列化为字符串
    params := fmt.Sprintf("page:%d,size:%d,sort:%s,%s,genre:%v,year:%d-%d,rating:%.1f,country:%s,lang:%s,keyword:%s,adult:%t",
        req.Page, req.PageSize, req.SortBy, req.SortOrder,
        req.Genre, req.YearFrom, req.YearTo, req.Rating,
        req.Country, req.Language, req.Keyword, req.IncludeAdult)
    
    h.Write([]byte(params))
    return fmt.Sprintf("movie_list:%x", h.Sum(nil))
}

// 从缓存获取
func (mls *MovieListService) getFromCache(ctx context.Context, key string) (*MovieListResponse, error) {
    // 先尝试本地缓存
    if cached, found := mls.cacheStore.Get(key); found {
        if response, ok := cached.(*MovieListResponse); ok {
            return response, nil
        }
    }

    // 再尝试Redis缓存
    data, err := mls.cacheStore.GetFromRedis(ctx, key)
    if err != nil {
        return nil, err
    }

    var response MovieListResponse
    if err := json.Unmarshal([]byte(data), &response); err != nil {
        return nil, err
    }

    // 回写到本地缓存
    mls.cacheStore.SetLocal(key, &response, 5*time.Minute)

    return &response, nil
}

// 缓存结果
func (mls *MovieListService) cacheResult(ctx context.Context, key string, response *MovieListResponse) error {
    // 序列化响应
    data, err := json.Marshal(response)
    if err != nil {
        return err
    }

    // 存储到Redis（15分钟过期）
    if err := mls.cacheStore.SetRedis(ctx, key, string(data), 15*time.Minute); err != nil {
        mls.logger.Errorf("Failed to cache to Redis: %v", err)
    }

    // 存储到本地缓存（5分钟过期）
    mls.cacheStore.SetLocal(key, response, 5*time.Minute)

    return nil
}
```

### 4. **查询优化器**

```go
type QueryBuilder struct {
    db         *gorm.DB
    conditions []string
    args       []interface{}
    joins      []string
    preloads   []string
    orderBy    string
    limitNum   int
    offsetNum  int
}

func NewQueryBuilder() *QueryBuilder {
    return &QueryBuilder{
        conditions: make([]string, 0),
        args:       make([]interface{}, 0),
        joins:      make([]string, 0),
        preloads:   make([]string, 0),
    }
}

func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
    qb.conditions = append(qb.conditions, condition)
    qb.args = append(qb.args, args...)
    return qb
}

func (qb *QueryBuilder) WhereIn(column string, values []string) *QueryBuilder {
    if len(values) > 0 {
        placeholders := make([]string, len(values))
        for i := range values {
            placeholders[i] = "?"
            qb.args = append(qb.args, values[i])
        }
        condition := fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ","))
        qb.conditions = append(qb.conditions, condition)
    }
    return qb
}

func (qb *QueryBuilder) Joins(join string) *QueryBuilder {
    qb.joins = append(qb.joins, join)
    return qb
}

func (qb *QueryBuilder) Preload(association string) *QueryBuilder {
    qb.preloads = append(qb.preloads, association)
    return qb
}

func (qb *QueryBuilder) Order(order string) *QueryBuilder {
    qb.orderBy = order
    return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
    qb.limitNum = limit
    return qb
}

func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
    qb.offsetNum = offset
    return qb
}

func (qb *QueryBuilder) Build(db *gorm.DB) *gorm.DB {
    query := db

    // 添加JOIN
    for _, join := range qb.joins {
        query = query.Joins(join)
    }

    // 添加WHERE条件
    if len(qb.conditions) > 0 {
        whereClause := strings.Join(qb.conditions, " AND ")
        query = query.Where(whereClause, qb.args...)
    }

    // 添加预加载
    for _, preload := range qb.preloads {
        query = query.Preload(preload)
    }

    // 添加排序
    if qb.orderBy != "" {
        query = query.Order(qb.orderBy)
    }

    // 添加分页
    if qb.limitNum > 0 {
        query = query.Limit(qb.limitNum)
    }
    if qb.offsetNum > 0 {
        query = query.Offset(qb.offsetNum)
    }

    return query
}
```

## 📊 性能监控

### 1. **电影列表指标**

```go
type MovieListMetrics struct {
    requestCount       *prometheus.CounterVec
    requestDuration    prometheus.Histogram
    cacheHitRate       *prometheus.CounterVec
    queryDuration      prometheus.Histogram
    errorCount         *prometheus.CounterVec
}

func NewMovieListMetrics() *MovieListMetrics {
    return &MovieListMetrics{
        requestCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "movie_list_requests_total",
                Help: "Total number of movie list requests",
            },
            []string{"status"},
        ),
        requestDuration: prometheus.NewHistogram(
            prometheus.HistogramOpts{
                Name: "movie_list_request_duration_seconds",
                Help: "Duration of movie list requests",
            },
        ),
        cacheHitRate: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "movie_list_cache_operations_total",
                Help: "Total number of cache operations",
            },
            []string{"type"}, // hit, miss
        ),
        queryDuration: prometheus.NewHistogram(
            prometheus.HistogramOpts{
                Name: "movie_list_query_duration_seconds",
                Help: "Duration of database queries",
            },
        ),
        errorCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "movie_list_errors_total",
                Help: "Total number of movie list errors",
            },
            []string{"type"},
        ),
    }
}
```

## 🔧 HTTP处理器

### 1. **电影列表API端点**

```go
func (mc *MovieController) GetMovieList(c *gin.Context) {
    var req MovieListRequest
    
    // 绑定查询参数
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": "请求参数错误",
            "error":   err.Error(),
        })
        return
    }

    // 调用服务
    response, err := mc.movieListService.GetMovieList(c.Request.Context(), &req)
    if err != nil {
        mc.logger.Errorf("Failed to get movie list: %v", err)
        c.JSON(500, gin.H{
            "success": false,
            "message": "获取电影列表失败",
        })
        return
    }

    // 设置缓存头
    c.Header("Cache-Control", "public, max-age=300") // 5分钟缓存

    c.JSON(200, response)
}

// 获取热门电影
func (mc *MovieController) GetPopularMovies(c *gin.Context) {
    req := MovieListRequest{
        Page:     1,
        PageSize: 20,
        SortBy:   "popularity",
        SortOrder: "desc",
    }

    response, err := mc.movieListService.GetMovieList(c.Request.Context(), &req)
    if err != nil {
        c.JSON(500, gin.H{
            "success": false,
            "message": "获取热门电影失败",
        })
        return
    }

    c.Header("Cache-Control", "public, max-age=600") // 10分钟缓存
    c.JSON(200, response)
}

// 获取最新电影
func (mc *MovieController) GetLatestMovies(c *gin.Context) {
    req := MovieListRequest{
        Page:     1,
        PageSize: 20,
        SortBy:   "release_date",
        SortOrder: "desc",
    }

    response, err := mc.movieListService.GetMovieList(c.Request.Context(), &req)
    if err != nil {
        c.JSON(500, gin.H{
            "success": false,
            "message": "获取最新电影失败",
        })
        return
    }

    c.Header("Cache-Control", "public, max-age=600")
    c.JSON(200, response)
}
```

## 📝 总结

电影列表接口为MovieInfo项目提供了强大的电影浏览功能：

**核心功能**：
1. **多维筛选**：类型、年份、评分、国家等多维度筛选
2. **灵活排序**：评分、热度、发布时间等多种排序方式
3. **高效分页**：支持大数据量的分页浏览
4. **智能缓存**：多层缓存提升响应性能

**性能优化**：
- 查询优化和索引设计
- 多层缓存策略
- 数据预加载机制
- 监控指标收集

**用户体验**：
- 快速响应时间
- 丰富的筛选选项
- 直观的数据展示
- 流畅的分页体验

下一步，我们将实现电影详情接口，为用户提供完整的电影信息展示。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第33步：电影详情接口](33-movie-detail.md)
