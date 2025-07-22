# 第33步：电影详情接口

## 📋 概述

电影详情接口是MovieInfo项目的核心功能之一，为用户提供完整、详细的电影信息展示。一个优秀的电影详情接口需要整合多维度的电影数据，提供丰富的内容展示和良好的性能表现。

## 🎯 设计目标

### 1. **信息完整性**
- 电影基础信息展示
- 演职人员详细信息
- 媒体资源集成
- 相关推荐内容

### 2. **性能优化**
- 快速数据加载
- 智能缓存策略
- 图片懒加载
- 数据预取机制

### 3. **用户体验**
- 丰富的视觉展示
- 交互式内容浏览
- 响应式设计
- 无缝导航体验

## 🔧 接口设计

### 1. **电影详情响应结构**

```go
// 电影详情响应
type MovieDetailResponse struct {
    Success bool              `json:"success"`
    Message string            `json:"message,omitempty"`
    Data    *MovieDetailData  `json:"data,omitempty"`
}

// 电影详情数据
type MovieDetailData struct {
    // 基础信息
    ID               string           `json:"id"`
    Title            string           `json:"title"`
    OriginalTitle    string           `json:"original_title"`
    Overview         string           `json:"overview"`
    Tagline          string           `json:"tagline,omitempty"`
    
    // 评分信息
    Rating           float64          `json:"rating"`
    VoteCount        int              `json:"vote_count"`
    Popularity       float64          `json:"popularity"`
    
    // 发布信息
    ReleaseDate      string           `json:"release_date"`
    Status           string           `json:"status"`
    Runtime          int              `json:"runtime"`
    
    // 分类信息
    Genres           []Genre          `json:"genres"`
    Countries        []Country        `json:"countries"`
    Languages        []Language       `json:"languages"`
    
    // 媒体资源
    PosterURL        string           `json:"poster_url"`
    BackdropURL      string           `json:"backdrop_url"`
    TrailerURL       string           `json:"trailer_url,omitempty"`
    Images           *MovieImages     `json:"images,omitempty"`
    Videos           []MovieVideo     `json:"videos,omitempty"`
    
    // 制作信息
    Budget           int64            `json:"budget,omitempty"`
    Revenue          int64            `json:"revenue,omitempty"`
    ProductionCompanies []Company     `json:"production_companies,omitempty"`
    
    // 演职人员
    Cast             []CastMember     `json:"cast,omitempty"`
    Crew             []CrewMember     `json:"crew,omitempty"`
    Director         *Person          `json:"director,omitempty"`
    
    // 相关内容
    Similar          []MovieListItem  `json:"similar,omitempty"`
    Recommendations  []MovieListItem  `json:"recommendations,omitempty"`
    
    // 统计信息
    ViewCount        int64            `json:"view_count"`
    FavoriteCount    int64            `json:"favorite_count"`
    WatchlistCount   int64            `json:"watchlist_count"`
    
    // 用户相关（需要登录）
    UserRating       *float64         `json:"user_rating,omitempty"`
    IsFavorite       bool             `json:"is_favorite"`
    IsInWatchlist    bool             `json:"is_in_watchlist"`
    
    // 元数据
    CreatedAt        time.Time        `json:"created_at"`
    UpdatedAt        time.Time        `json:"updated_at"`
}

// 电影图片
type MovieImages struct {
    Posters    []ImageInfo `json:"posters,omitempty"`
    Backdrops  []ImageInfo `json:"backdrops,omitempty"`
    Logos      []ImageInfo `json:"logos,omitempty"`
}

// 图片信息
type ImageInfo struct {
    URL        string  `json:"url"`
    Width      int     `json:"width"`
    Height     int     `json:"height"`
    Language   string  `json:"language,omitempty"`
    VoteAverage float64 `json:"vote_average,omitempty"`
}

// 电影视频
type MovieVideo struct {
    ID          string `json:"id"`
    Key         string `json:"key"`
    Name        string `json:"name"`
    Site        string `json:"site"`
    Type        string `json:"type"`
    Size        int    `json:"size"`
    Official    bool   `json:"official"`
    PublishedAt string `json:"published_at"`
}

// 演员信息
type CastMember struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Character   string `json:"character"`
    ProfileURL  string `json:"profile_url,omitempty"`
    Order       int    `json:"order"`
    Gender      int    `json:"gender,omitempty"`
}

// 工作人员信息
type CrewMember struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Job         string `json:"job"`
    Department  string `json:"department"`
    ProfileURL  string `json:"profile_url,omitempty"`
    Gender      int    `json:"gender,omitempty"`
}
```

### 2. **电影详情服务实现**

```go
type MovieDetailService struct {
    movieRepo      MovieRepository
    castRepo       CastRepository
    userRepo       UserRepository
    cacheStore     CacheStore
    imageService   ImageService
    recommendationService RecommendationService
    logger         *logrus.Logger
    metrics        *MovieDetailMetrics
}

func NewMovieDetailService(
    movieRepo MovieRepository,
    castRepo CastRepository,
    userRepo UserRepository,
    cacheStore CacheStore,
    imageService ImageService,
    recommendationService RecommendationService,
) *MovieDetailService {
    return &MovieDetailService{
        movieRepo:      movieRepo,
        castRepo:       castRepo,
        userRepo:       userRepo,
        cacheStore:     cacheStore,
        imageService:   imageService,
        recommendationService: recommendationService,
        logger:         logrus.New(),
        metrics:        NewMovieDetailMetrics(),
    }
}

// 获取电影详情
func (mds *MovieDetailService) GetMovieDetail(ctx context.Context, movieID string, userID string) (*MovieDetailResponse, error) {
    start := time.Now()
    defer func() {
        mds.metrics.ObserveDetailRequestDuration(time.Since(start))
    }()

    // 验证电影ID
    if movieID == "" {
        mds.metrics.IncInvalidRequests()
        return &MovieDetailResponse{
            Success: false,
            Message: "电影ID不能为空",
        }, nil
    }

    // 生成缓存键
    cacheKey := mds.generateCacheKey(movieID, userID)

    // 尝试从缓存获取
    if cachedResult, err := mds.getFromCache(ctx, cacheKey); err == nil {
        mds.metrics.IncCacheHits()
        // 异步更新访问统计
        go mds.updateViewCount(context.Background(), movieID)
        return cachedResult, nil
    }
    mds.metrics.IncCacheMisses()

    // 获取电影基础信息
    movie, err := mds.movieRepo.FindByIDWithDetails(ctx, movieID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            mds.metrics.IncNotFoundErrors()
            return &MovieDetailResponse{
                Success: false,
                Message: "电影不存在",
            }, nil
        }
        mds.logger.Errorf("Failed to get movie: %v", err)
        mds.metrics.IncQueryErrors()
        return nil, errors.New("获取电影信息失败")
    }

    // 并发获取相关数据
    var (
        cast            []CastMember
        crew            []CrewMember
        images          *MovieImages
        videos          []MovieVideo
        similar         []MovieListItem
        recommendations []MovieListItem
        userInteraction *UserMovieInteraction
        wg              sync.WaitGroup
        mu              sync.Mutex
        errs            []error
    )

    // 获取演职人员信息
    wg.Add(1)
    go func() {
        defer wg.Done()
        if c, cr, err := mds.getCastAndCrew(ctx, movieID); err != nil {
            mu.Lock()
            errs = append(errs, err)
            mu.Unlock()
        } else {
            cast, crew = c, cr
        }
    }()

    // 获取媒体资源
    wg.Add(1)
    go func() {
        defer wg.Done()
        if img, vid, err := mds.getMediaResources(ctx, movieID); err != nil {
            mu.Lock()
            errs = append(errs, err)
            mu.Unlock()
        } else {
            images, videos = img, vid
        }
    }()

    // 获取相似电影
    wg.Add(1)
    go func() {
        defer wg.Done()
        if sim, err := mds.getSimilarMovies(ctx, movieID); err != nil {
            mu.Lock()
            errs = append(errs, err)
            mu.Unlock()
        } else {
            similar = sim
        }
    }()

    // 获取推荐电影
    wg.Add(1)
    go func() {
        defer wg.Done()
        if rec, err := mds.getRecommendations(ctx, movieID, userID); err != nil {
            mu.Lock()
            errs = append(errs, err)
            mu.Unlock()
        } else {
            recommendations = rec
        }
    }()

    // 获取用户交互信息（如果已登录）
    if userID != "" {
        wg.Add(1)
        go func() {
            defer wg.Done()
            if ui, err := mds.getUserInteraction(ctx, userID, movieID); err != nil {
                mu.Lock()
                errs = append(errs, err)
                mu.Unlock()
            } else {
                userInteraction = ui
            }
        }()
    }

    // 等待所有并发操作完成
    wg.Wait()

    // 检查是否有严重错误
    if len(errs) > 0 {
        mds.logger.Errorf("Errors occurred while fetching movie details: %v", errs)
        // 非关键错误不影响主要功能
    }

    // 构建响应数据
    detailData := mds.buildDetailData(movie, cast, crew, images, videos, similar, recommendations, userInteraction)

    response := &MovieDetailResponse{
        Success: true,
        Data:    detailData,
    }

    // 异步缓存结果
    go func() {
        if err := mds.cacheResult(context.Background(), cacheKey, response); err != nil {
            mds.logger.Errorf("Failed to cache result: %v", err)
        }
    }()

    // 异步更新访问统计
    go mds.updateViewCount(context.Background(), movieID)

    mds.metrics.IncSuccessfulRequests()
    return response, nil
}

// 获取演职人员信息
func (mds *MovieDetailService) getCastAndCrew(ctx context.Context, movieID string) ([]CastMember, []CrewMember, error) {
    // 获取演员信息
    castMembers, err := mds.castRepo.FindCastByMovieID(ctx, movieID)
    if err != nil {
        return nil, nil, err
    }

    // 获取工作人员信息
    crewMembers, err := mds.castRepo.FindCrewByMovieID(ctx, movieID)
    if err != nil {
        return nil, nil, err
    }

    // 转换为响应格式
    cast := make([]CastMember, len(castMembers))
    for i, member := range castMembers {
        cast[i] = CastMember{
            ID:         member.PersonID,
            Name:       member.Person.Name,
            Character:  member.Character,
            ProfileURL: member.Person.ProfileURL,
            Order:      member.Order,
            Gender:     member.Person.Gender,
        }
    }

    crew := make([]CrewMember, len(crewMembers))
    for i, member := range crewMembers {
        crew[i] = CrewMember{
            ID:         member.PersonID,
            Name:       member.Person.Name,
            Job:        member.Job,
            Department: member.Department,
            ProfileURL: member.Person.ProfileURL,
            Gender:     member.Person.Gender,
        }
    }

    return cast, crew, nil
}

// 获取媒体资源
func (mds *MovieDetailService) getMediaResources(ctx context.Context, movieID string) (*MovieImages, []MovieVideo, error) {
    // 获取图片资源
    images, err := mds.imageService.GetMovieImages(ctx, movieID)
    if err != nil {
        mds.logger.Errorf("Failed to get movie images: %v", err)
        images = &MovieImages{} // 使用空结构体而不是返回错误
    }

    // 获取视频资源
    videos, err := mds.movieRepo.FindVideosByMovieID(ctx, movieID)
    if err != nil {
        mds.logger.Errorf("Failed to get movie videos: %v", err)
        videos = []MovieVideo{} // 使用空切片而不是返回错误
    }

    return images, videos, nil
}

// 获取相似电影
func (mds *MovieDetailService) getSimilarMovies(ctx context.Context, movieID string) ([]MovieListItem, error) {
    similar, err := mds.recommendationService.GetSimilarMovies(ctx, movieID, 10)
    if err != nil {
        mds.logger.Errorf("Failed to get similar movies: %v", err)
        return []MovieListItem{}, nil // 返回空切片而不是错误
    }

    return mds.convertToListItems(similar), nil
}

// 获取推荐电影
func (mds *MovieDetailService) getRecommendations(ctx context.Context, movieID, userID string) ([]MovieListItem, error) {
    var recommendations []*Movie
    var err error

    if userID != "" {
        // 个性化推荐
        recommendations, err = mds.recommendationService.GetPersonalizedRecommendations(ctx, userID, movieID, 10)
    } else {
        // 通用推荐
        recommendations, err = mds.recommendationService.GetGeneralRecommendations(ctx, movieID, 10)
    }

    if err != nil {
        mds.logger.Errorf("Failed to get recommendations: %v", err)
        return []MovieListItem{}, nil
    }

    return mds.convertToListItems(recommendations), nil
}

// 获取用户交互信息
func (mds *MovieDetailService) getUserInteraction(ctx context.Context, userID, movieID string) (*UserMovieInteraction, error) {
    interaction, err := mds.userRepo.FindUserMovieInteraction(ctx, userID, movieID)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, err
    }

    if interaction == nil {
        // 创建默认交互信息
        interaction = &UserMovieInteraction{
            UserID:        userID,
            MovieID:       movieID,
            IsFavorite:    false,
            IsInWatchlist: false,
        }
    }

    return interaction, nil
}

// 更新访问统计
func (mds *MovieDetailService) updateViewCount(ctx context.Context, movieID string) {
    if err := mds.movieRepo.IncrementViewCount(ctx, movieID); err != nil {
        mds.logger.Errorf("Failed to update view count: %v", err)
    }
}
```

### 3. **缓存策略**

```go
// 生成缓存键
func (mds *MovieDetailService) generateCacheKey(movieID, userID string) string {
    if userID != "" {
        return fmt.Sprintf("movie_detail:%s:user:%s", movieID, userID)
    }
    return fmt.Sprintf("movie_detail:%s", movieID)
}

// 从缓存获取
func (mds *MovieDetailService) getFromCache(ctx context.Context, key string) (*MovieDetailResponse, error) {
    data, err := mds.cacheStore.Get(ctx, key)
    if err != nil {
        return nil, err
    }

    var response MovieDetailResponse
    if err := json.Unmarshal([]byte(data), &response); err != nil {
        return nil, err
    }

    return &response, nil
}

// 缓存结果
func (mds *MovieDetailService) cacheResult(ctx context.Context, key string, response *MovieDetailResponse) error {
    data, err := json.Marshal(response)
    if err != nil {
        return err
    }

    // 根据是否包含用户信息设置不同的过期时间
    var expiry time.Duration
    if strings.Contains(key, ":user:") {
        expiry = 5 * time.Minute // 用户相关数据缓存时间较短
    } else {
        expiry = 30 * time.Minute // 公共数据缓存时间较长
    }

    return mds.cacheStore.Set(ctx, key, string(data), expiry)
}
```

### 4. **图片服务集成**

```go
type ImageService struct {
    imageRepo  ImageRepository
    cdnService CDNService
    logger     *logrus.Logger
}

func (is *ImageService) GetMovieImages(ctx context.Context, movieID string) (*MovieImages, error) {
    // 获取电影图片记录
    imageRecords, err := is.imageRepo.FindByMovieID(ctx, movieID)
    if err != nil {
        return nil, err
    }

    images := &MovieImages{
        Posters:   make([]ImageInfo, 0),
        Backdrops: make([]ImageInfo, 0),
        Logos:     make([]ImageInfo, 0),
    }

    for _, record := range imageRecords {
        imageInfo := ImageInfo{
            URL:         is.cdnService.GetImageURL(record.FilePath),
            Width:       record.Width,
            Height:      record.Height,
            Language:    record.Language,
            VoteAverage: record.VoteAverage,
        }

        switch record.Type {
        case "poster":
            images.Posters = append(images.Posters, imageInfo)
        case "backdrop":
            images.Backdrops = append(images.Backdrops, imageInfo)
        case "logo":
            images.Logos = append(images.Logos, imageInfo)
        }
    }

    return images, nil
}
```

## 📊 性能监控

### 1. **电影详情指标**

```go
type MovieDetailMetrics struct {
    requestCount       *prometheus.CounterVec
    requestDuration    prometheus.Histogram
    cacheHitRate       *prometheus.CounterVec
    componentDuration  *prometheus.HistogramVec
    errorCount         *prometheus.CounterVec
}

func NewMovieDetailMetrics() *MovieDetailMetrics {
    return &MovieDetailMetrics{
        requestCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "movie_detail_requests_total",
                Help: "Total number of movie detail requests",
            },
            []string{"status"},
        ),
        requestDuration: prometheus.NewHistogram(
            prometheus.HistogramOpts{
                Name: "movie_detail_request_duration_seconds",
                Help: "Duration of movie detail requests",
            },
        ),
        cacheHitRate: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "movie_detail_cache_operations_total",
                Help: "Total number of cache operations",
            },
            []string{"type"},
        ),
        componentDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "movie_detail_component_duration_seconds",
                Help: "Duration of different components loading",
            },
            []string{"component"},
        ),
        errorCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "movie_detail_errors_total",
                Help: "Total number of movie detail errors",
            },
            []string{"type"},
        ),
    }
}
```

## 🔧 HTTP处理器

### 1. **电影详情API端点**

```go
func (mc *MovieController) GetMovieDetail(c *gin.Context) {
    movieID := c.Param("id")
    if movieID == "" {
        c.JSON(400, gin.H{
            "success": false,
            "message": "电影ID不能为空",
        })
        return
    }

    // 获取用户ID（如果已登录）
    userID := mc.getUserIDFromContext(c)

    // 调用服务
    response, err := mc.movieDetailService.GetMovieDetail(c.Request.Context(), movieID, userID)
    if err != nil {
        mc.logger.Errorf("Failed to get movie detail: %v", err)
        c.JSON(500, gin.H{
            "success": false,
            "message": "获取电影详情失败",
        })
        return
    }

    // 设置缓存头
    if userID == "" {
        c.Header("Cache-Control", "public, max-age=1800") // 30分钟缓存
    } else {
        c.Header("Cache-Control", "private, max-age=300") // 5分钟缓存
    }

    c.JSON(200, response)
}

func (mc *MovieController) getUserIDFromContext(c *gin.Context) string {
    if user, exists := c.Get("user"); exists {
        if userInfo, ok := user.(*User); ok {
            return userInfo.ID
        }
    }
    return ""
}
```

## 📝 总结

电影详情接口为MovieInfo项目提供了完整的电影信息展示：

**核心功能**：
1. **完整信息**：电影基础信息、演职人员、媒体资源
2. **个性化内容**：用户交互信息、个性化推荐
3. **相关推荐**：相似电影、智能推荐
4. **统计信息**：访问量、收藏数等统计数据

**性能优化**：
- 并发数据获取
- 智能缓存策略
- 图片CDN加速
- 数据预取机制

**用户体验**：
- 丰富的视觉内容
- 快速加载响应
- 个性化推荐
- 无缝导航体验

下一步，我们将实现电影搜索功能，为用户提供强大的电影发现能力。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第34步：电影搜索功能](34-movie-search.md)
