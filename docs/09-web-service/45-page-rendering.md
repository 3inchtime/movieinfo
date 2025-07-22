# 第45步：页面渲染逻辑

## 📋 概述

页面渲染逻辑是Web应用的核心组件，负责将后端数据与前端模板结合，生成用户最终看到的HTML页面。MovieInfo项目采用服务端渲染(SSR)策略，提供快速的首屏加载和良好的SEO支持。

## 🎯 设计目标

### 1. **渲染性能**
- 快速的页面生成
- 高效的数据获取
- 智能的缓存策略
- 异步数据加载

### 2. **用户体验**
- 快速的首屏渲染
- 渐进式页面增强
- 响应式设计支持
- 无障碍访问支持

### 3. **SEO优化**
- 完整的HTML结构
- 语义化的标签
- 结构化数据支持
- 社交媒体优化

## 🔧 页面渲染架构

### 1. **页面控制器结构**

```go
// 页面数据结构
type PageData struct {
    // 基础信息
    Title       string                 `json:"title"`
    Description string                 `json:"description"`
    Keywords    string                 `json:"keywords"`
    CanonicalURL string                `json:"canonical_url"`
    
    // 用户信息
    User        *User                  `json:"user,omitempty"`
    IsLoggedIn  bool                   `json:"is_logged_in"`
    
    // 页面特定数据
    Data        interface{}            `json:"data"`
    
    // 元数据
    Meta        map[string]interface{} `json:"meta"`
    
    // 导航和面包屑
    Navigation  *Navigation            `json:"navigation"`
    Breadcrumbs []Breadcrumb           `json:"breadcrumbs"`
    
    // 性能数据
    RenderTime  time.Duration          `json:"render_time"`
    RequestID   string                 `json:"request_id"`
}

// 页面控制器
type PageController struct {
    serviceClient   *ServiceClient
    templateEngine  *TemplateEngine
    cacheManager    *CacheManager
    seoManager      *SEOManager
    logger          *logrus.Logger
    metrics         *PageMetrics
}

func NewPageController(
    serviceClient *ServiceClient,
    templateEngine *TemplateEngine,
    cacheManager *CacheManager,
    seoManager *SEOManager,
) *PageController {
    return &PageController{
        serviceClient:  serviceClient,
        templateEngine: templateEngine,
        cacheManager:   cacheManager,
        seoManager:     seoManager,
        logger:         logrus.New(),
        metrics:        NewPageMetrics(),
    }
}

// 渲染页面基础方法
func (pc *PageController) renderPage(c *gin.Context, template string, data interface{}) {
    start := time.Now()
    defer func() {
        pc.metrics.ObserveRenderTime(template, time.Since(start))
    }()
    
    // 构建页面数据
    pageData := pc.buildPageData(c, data)
    
    // 设置响应头
    pc.setResponseHeaders(c, template)
    
    // 渲染模板
    c.HTML(http.StatusOK, template, pageData)
    
    pc.metrics.IncPageViews(template)
    pc.logger.Debugf("Rendered page: %s in %v", template, time.Since(start))
}

// 构建页面数据
func (pc *PageController) buildPageData(c *gin.Context, data interface{}) *PageData {
    pageData := &PageData{
        Data:        data,
        Meta:        make(map[string]interface{}),
        RequestID:   c.GetString("request_id"),
        RenderTime:  0, // 将在渲染完成后设置
    }
    
    // 设置用户信息
    if user, exists := c.Get("user"); exists {
        pageData.User = user.(*User)
        pageData.IsLoggedIn = true
    }
    
    // 设置导航信息
    pageData.Navigation = pc.buildNavigation(c)
    
    // 设置面包屑
    pageData.Breadcrumbs = pc.buildBreadcrumbs(c)
    
    return pageData
}
```

### 2. **首页渲染**

```go
// 首页控制器
func (pc *PageController) Home(c *gin.Context) {
    // 检查缓存
    cacheKey := "page:home"
    if cached := pc.cacheManager.Get(cacheKey); cached != nil {
        c.HTML(http.StatusOK, "pages/home", cached)
        return
    }
    
    // 并发获取首页数据
    var (
        popularMovies    []*Movie
        latestMovies     []*Movie
        topRatedMovies   []*Movie
        featuredMovies   []*Movie
        categories       []*Category
        wg               sync.WaitGroup
        mu               sync.Mutex
        errors           []error
    )
    
    // 获取热门电影
    wg.Add(1)
    go func() {
        defer wg.Done()
        if movies, err := pc.serviceClient.GetPopularMovies(c.Request.Context(), 12); err != nil {
            mu.Lock()
            errors = append(errors, err)
            mu.Unlock()
        } else {
            popularMovies = movies
        }
    }()
    
    // 获取最新电影
    wg.Add(1)
    go func() {
        defer wg.Done()
        if movies, err := pc.serviceClient.GetLatestMovies(c.Request.Context(), 12); err != nil {
            mu.Lock()
            errors = append(errors, err)
            mu.Unlock()
        } else {
            latestMovies = movies
        }
    }()
    
    // 获取高分电影
    wg.Add(1)
    go func() {
        defer wg.Done()
        if movies, err := pc.serviceClient.GetTopRatedMovies(c.Request.Context(), 12); err != nil {
            mu.Lock()
            errors = append(errors, err)
            mu.Unlock()
        } else {
            topRatedMovies = movies
        }
    }()
    
    // 获取精选电影
    wg.Add(1)
    go func() {
        defer wg.Done()
        if movies, err := pc.serviceClient.GetFeaturedMovies(c.Request.Context(), 6); err != nil {
            mu.Lock()
            errors = append(errors, err)
            mu.Unlock()
        } else {
            featuredMovies = movies
        }
    }()
    
    // 获取分类
    wg.Add(1)
    go func() {
        defer wg.Done()
        if cats, err := pc.serviceClient.GetCategories(c.Request.Context()); err != nil {
            mu.Lock()
            errors = append(errors, err)
            mu.Unlock()
        } else {
            categories = cats
        }
    }()
    
    wg.Wait()
    
    // 检查是否有严重错误
    if len(errors) > 3 {
        pc.logger.Errorf("Too many errors loading home page: %v", errors)
        c.HTML(http.StatusInternalServerError, "error", gin.H{
            "title":   "服务器错误",
            "message": "页面加载失败，请稍后重试",
        })
        return
    }
    
    // 构建首页数据
    homeData := &HomePageData{
        PopularMovies:  popularMovies,
        LatestMovies:   latestMovies,
        TopRatedMovies: topRatedMovies,
        FeaturedMovies: featuredMovies,
        Categories:     categories,
        Stats:          pc.getHomePageStats(c.Request.Context()),
    }
    
    // 设置SEO信息
    pageData := pc.buildPageData(c, homeData)
    pageData.Title = "MovieInfo - 专业的电影信息平台"
    pageData.Description = "MovieInfo提供最新的电影信息、影评、评分和推荐，帮您发现好电影"
    pageData.Keywords = "电影,影评,评分,推荐,MovieInfo"
    pageData.CanonicalURL = pc.buildCanonicalURL(c, "/")
    
    // 缓存页面数据
    pc.cacheManager.Set(cacheKey, pageData, 10*time.Minute)
    
    pc.renderPage(c, "pages/home", pageData)
}

type HomePageData struct {
    PopularMovies  []*Movie    `json:"popular_movies"`
    LatestMovies   []*Movie    `json:"latest_movies"`
    TopRatedMovies []*Movie    `json:"top_rated_movies"`
    FeaturedMovies []*Movie    `json:"featured_movies"`
    Categories     []*Category `json:"categories"`
    Stats          *HomeStats  `json:"stats"`
}

type HomeStats struct {
    TotalMovies   int `json:"total_movies"`
    TotalUsers    int `json:"total_users"`
    TotalReviews  int `json:"total_reviews"`
    TotalRatings  int `json:"total_ratings"`
}
```

### 3. **电影详情页渲染**

```go
// 电影详情页
func (pc *PageController) MovieDetail(c *gin.Context) {
    movieID := c.Param("id")
    if movieID == "" {
        c.HTML(http.StatusNotFound, "error", gin.H{
            "title":   "电影未找到",
            "message": "您访问的电影不存在",
        })
        return
    }
    
    // 检查缓存
    userID := pc.getUserID(c)
    cacheKey := fmt.Sprintf("page:movie:%s:user:%s", movieID, userID)
    if cached := pc.cacheManager.Get(cacheKey); cached != nil {
        c.HTML(http.StatusOK, "pages/movie-detail", cached)
        return
    }
    
    // 并发获取电影相关数据
    var (
        movie           *Movie
        comments        []*Comment
        similarMovies   []*Movie
        recommendations []*Movie
        userRating      *UserRating
        userInteraction *UserInteraction
        wg              sync.WaitGroup
        mu              sync.Mutex
        errors          []error
    )
    
    // 获取电影详情
    wg.Add(1)
    go func() {
        defer wg.Done()
        if m, err := pc.serviceClient.GetMovieDetail(c.Request.Context(), movieID); err != nil {
            mu.Lock()
            errors = append(errors, err)
            mu.Unlock()
        } else {
            movie = m
        }
    }()
    
    // 获取评论
    wg.Add(1)
    go func() {
        defer wg.Done()
        if cs, err := pc.serviceClient.GetMovieComments(c.Request.Context(), movieID, 1, 10); err != nil {
            mu.Lock()
            errors = append(errors, err)
            mu.Unlock()
        } else {
            comments = cs
        }
    }()
    
    // 获取相似电影
    wg.Add(1)
    go func() {
        defer wg.Done()
        if movies, err := pc.serviceClient.GetSimilarMovies(c.Request.Context(), movieID, 12); err != nil {
            mu.Lock()
            errors = append(errors, err)
            mu.Unlock()
        } else {
            similarMovies = movies
        }
    }()
    
    // 获取推荐电影（如果用户已登录）
    if userID != "" {
        wg.Add(1)
        go func() {
            defer wg.Done()
            if movies, err := pc.serviceClient.GetRecommendations(c.Request.Context(), userID, movieID, 12); err != nil {
                mu.Lock()
                errors = append(errors, err)
                mu.Unlock()
            } else {
                recommendations = movies
            }
        }()
        
        // 获取用户评分
        wg.Add(1)
        go func() {
            defer wg.Done()
            if rating, err := pc.serviceClient.GetUserRating(c.Request.Context(), userID, movieID); err == nil {
                userRating = rating
            }
        }()
        
        // 获取用户交互信息
        wg.Add(1)
        go func() {
            defer wg.Done()
            if interaction, err := pc.serviceClient.GetUserInteraction(c.Request.Context(), userID, movieID); err == nil {
                userInteraction = interaction
            }
        }()
    }
    
    wg.Wait()
    
    // 检查电影是否存在
    if movie == nil {
        c.HTML(http.StatusNotFound, "error", gin.H{
            "title":   "电影未找到",
            "message": "您访问的电影不存在",
        })
        return
    }
    
    // 构建电影详情数据
    movieDetailData := &MovieDetailPageData{
        Movie:           movie,
        Comments:        comments,
        SimilarMovies:   similarMovies,
        Recommendations: recommendations,
        UserRating:      userRating,
        UserInteraction: userInteraction,
        RelatedData:     pc.getMovieRelatedData(c.Request.Context(), movie),
    }
    
    // 设置SEO信息
    pageData := pc.buildPageData(c, movieDetailData)
    pageData.Title = fmt.Sprintf("%s - MovieInfo", movie.Title)
    pageData.Description = pc.truncateString(movie.Overview, 160)
    pageData.Keywords = fmt.Sprintf("%s,电影,%s", movie.Title, strings.Join(movie.GenreNames, ","))
    pageData.CanonicalURL = pc.buildCanonicalURL(c, fmt.Sprintf("/movies/%s", movieID))
    
    // 添加结构化数据
    pageData.Meta["structured_data"] = pc.seoManager.GenerateMovieStructuredData(movie)
    pageData.Meta["og_data"] = pc.seoManager.GenerateOpenGraphData(movie)
    pageData.Meta["twitter_data"] = pc.seoManager.GenerateTwitterCardData(movie)
    
    // 缓存页面数据（较短时间，因为包含用户特定数据）
    if userID == "" {
        pc.cacheManager.Set(cacheKey, pageData, 30*time.Minute)
    } else {
        pc.cacheManager.Set(cacheKey, pageData, 5*time.Minute)
    }
    
    pc.renderPage(c, "pages/movie-detail", pageData)
}

type MovieDetailPageData struct {
    Movie           *Movie           `json:"movie"`
    Comments        []*Comment       `json:"comments"`
    SimilarMovies   []*Movie         `json:"similar_movies"`
    Recommendations []*Movie         `json:"recommendations"`
    UserRating      *UserRating      `json:"user_rating,omitempty"`
    UserInteraction *UserInteraction `json:"user_interaction,omitempty"`
    RelatedData     *MovieRelatedData `json:"related_data"`
}

type MovieRelatedData struct {
    Director        *Person   `json:"director"`
    MainCast        []*Person `json:"main_cast"`
    ProductionInfo  *ProductionInfo `json:"production_info"`
    Awards          []*Award  `json:"awards"`
}
```

### 4. **电影列表页渲染**

```go
// 电影列表页
func (pc *PageController) MovieList(c *gin.Context) {
    // 解析查询参数
    params := pc.parseListParams(c)
    
    // 检查缓存
    cacheKey := pc.generateListCacheKey("movies", params)
    if cached := pc.cacheManager.Get(cacheKey); cached != nil {
        c.HTML(http.StatusOK, "pages/movie-list", cached)
        return
    }
    
    // 获取电影列表
    movies, total, err := pc.serviceClient.GetMovieList(c.Request.Context(), params)
    if err != nil {
        pc.logger.Errorf("Failed to get movie list: %v", err)
        c.HTML(http.StatusInternalServerError, "error", gin.H{
            "title":   "加载失败",
            "message": "电影列表加载失败，请稍后重试",
        })
        return
    }
    
    // 获取筛选选项
    filterOptions := pc.getFilterOptions(c.Request.Context())
    
    // 构建分页信息
    pagination := pc.buildPagination(params.Page, params.PageSize, total)
    
    // 构建列表数据
    listData := &MovieListPageData{
        Movies:        movies,
        Total:         total,
        Pagination:    pagination,
        FilterOptions: filterOptions,
        CurrentFilter: params,
        SortOptions:   pc.getSortOptions(),
    }
    
    // 设置SEO信息
    pageData := pc.buildPageData(c, listData)
    pageData.Title = pc.buildListTitle(params)
    pageData.Description = pc.buildListDescription(params)
    pageData.Keywords = "电影列表,电影推荐,热门电影"
    pageData.CanonicalURL = pc.buildCanonicalURL(c, c.Request.URL.Path)
    
    // 设置面包屑
    pageData.Breadcrumbs = pc.buildListBreadcrumbs(params)
    
    // 缓存页面数据
    pc.cacheManager.Set(cacheKey, pageData, 15*time.Minute)
    
    pc.renderPage(c, "pages/movie-list", pageData)
}

type MovieListPageData struct {
    Movies        []*Movie        `json:"movies"`
    Total         int64           `json:"total"`
    Pagination    *Pagination     `json:"pagination"`
    FilterOptions *FilterOptions  `json:"filter_options"`
    CurrentFilter *ListParams     `json:"current_filter"`
    SortOptions   []*SortOption   `json:"sort_options"`
}

type ListParams struct {
    Page      int      `json:"page"`
    PageSize  int      `json:"page_size"`
    Genre     string   `json:"genre"`
    Year      int      `json:"year"`
    Country   string   `json:"country"`
    SortBy    string   `json:"sort_by"`
    SortOrder string   `json:"sort_order"`
    Keywords  string   `json:"keywords"`
}

// 解析列表参数
func (pc *PageController) parseListParams(c *gin.Context) *ListParams {
    params := &ListParams{
        Page:      1,
        PageSize:  20,
        SortBy:    "popularity",
        SortOrder: "desc",
    }
    
    if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && page > 0 {
        params.Page = page
    }
    
    if pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20")); err == nil && pageSize > 0 && pageSize <= 100 {
        params.PageSize = pageSize
    }
    
    params.Genre = c.Query("genre")
    if year, err := strconv.Atoi(c.Query("year")); err == nil {
        params.Year = year
    }
    params.Country = c.Query("country")
    params.SortBy = c.DefaultQuery("sort_by", "popularity")
    params.SortOrder = c.DefaultQuery("sort_order", "desc")
    params.Keywords = c.Query("q")
    
    return params
}
```

### 5. **SEO优化管理**

```go
// SEO管理器
type SEOManager struct {
    config *SEOConfig
    logger *logrus.Logger
}

type SEOConfig struct {
    SiteName    string `yaml:"site_name"`
    SiteURL     string `yaml:"site_url"`
    DefaultMeta *DefaultMeta `yaml:"default_meta"`
}

type DefaultMeta struct {
    Title       string `yaml:"title"`
    Description string `yaml:"description"`
    Keywords    string `yaml:"keywords"`
    Author      string `yaml:"author"`
}

func NewSEOManager(config *SEOConfig) *SEOManager {
    return &SEOManager{
        config: config,
        logger: logrus.New(),
    }
}

// 生成电影结构化数据
func (sm *SEOManager) GenerateMovieStructuredData(movie *Movie) map[string]interface{} {
    return map[string]interface{}{
        "@context": "https://schema.org",
        "@type":    "Movie",
        "name":     movie.Title,
        "description": movie.Overview,
        "image":    movie.PosterURL,
        "datePublished": movie.ReleaseDate,
        "genre":    movie.GenreNames,
        "director": map[string]interface{}{
            "@type": "Person",
            "name":  movie.Director,
        },
        "aggregateRating": map[string]interface{}{
            "@type":       "AggregateRating",
            "ratingValue": movie.Rating,
            "ratingCount": movie.VoteCount,
            "bestRating":  10,
            "worstRating": 1,
        },
        "url": fmt.Sprintf("%s/movies/%s", sm.config.SiteURL, movie.ID),
    }
}

// 生成Open Graph数据
func (sm *SEOManager) GenerateOpenGraphData(movie *Movie) map[string]string {
    return map[string]string{
        "og:type":        "video.movie",
        "og:title":       movie.Title,
        "og:description": sm.truncateString(movie.Overview, 200),
        "og:image":       movie.PosterURL,
        "og:url":         fmt.Sprintf("%s/movies/%s", sm.config.SiteURL, movie.ID),
        "og:site_name":   sm.config.SiteName,
    }
}

// 生成Twitter Card数据
func (sm *SEOManager) GenerateTwitterCardData(movie *Movie) map[string]string {
    return map[string]string{
        "twitter:card":        "summary_large_image",
        "twitter:title":       movie.Title,
        "twitter:description": sm.truncateString(movie.Overview, 200),
        "twitter:image":       movie.PosterURL,
    }
}

func (sm *SEOManager) truncateString(s string, length int) string {
    if len(s) <= length {
        return s
    }
    return s[:length-3] + "..."
}
```

## 📊 性能监控

### 1. **页面渲染指标**

```go
type PageMetrics struct {
    pageViews    *prometheus.CounterVec
    renderTime   *prometheus.HistogramVec
    cacheHitRate *prometheus.CounterVec
    errorRate    *prometheus.CounterVec
}

func NewPageMetrics() *PageMetrics {
    return &PageMetrics{
        pageViews: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "page_views_total",
                Help: "Total number of page views",
            },
            []string{"page", "status"},
        ),
        renderTime: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "page_render_time_seconds",
                Help: "Page render time",
            },
            []string{"page"},
        ),
        cacheHitRate: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "page_cache_operations_total",
                Help: "Total number of page cache operations",
            },
            []string{"type"},
        ),
        errorRate: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "page_errors_total",
                Help: "Total number of page errors",
            },
            []string{"page", "error_type"},
        ),
    }
}

func (pm *PageMetrics) IncPageViews(page string) {
    pm.pageViews.WithLabelValues(page, "success").Inc()
}

func (pm *PageMetrics) ObserveRenderTime(page string, duration time.Duration) {
    pm.renderTime.WithLabelValues(page).Observe(duration.Seconds())
}

func (pm *PageMetrics) IncCacheHits() {
    pm.cacheHitRate.WithLabelValues("hit").Inc()
}

func (pm *PageMetrics) IncCacheMisses() {
    pm.cacheHitRate.WithLabelValues("miss").Inc()
}
```

## 📝 总结

页面渲染逻辑为MovieInfo项目提供了完整的服务端渲染能力：

**核心功能**：
1. **高效渲染**：并发数据获取和智能缓存策略
2. **SEO优化**：完整的元数据和结构化数据支持
3. **用户体验**：快速的首屏加载和响应式设计
4. **错误处理**：优雅的错误页面和降级策略

**技术特性**：
- 模块化的页面控制器
- 灵活的数据组装机制
- 智能的缓存管理
- 完善的性能监控

**SEO支持**：
- 结构化数据标记
- Open Graph协议
- Twitter Card支持
- 语义化HTML结构

至此，主页服务的核心功能已经完成。下一步，我们将继续完成前端页面开发的相关文档。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第46步：页面布局设计](../10-frontend-pages/46-layout-design.md)
