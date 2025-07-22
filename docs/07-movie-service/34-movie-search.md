# 第34步：电影搜索功能

## 📋 概述

电影搜索功能是MovieInfo项目的核心发现机制，为用户提供快速、准确、智能的电影查找体验。一个优秀的搜索系统需要支持多种搜索方式、智能排序算法和高性能的查询处理。

## 🎯 设计目标

### 1. **搜索准确性**
- 全文搜索支持
- 模糊匹配算法
- 智能纠错功能
- 多语言搜索支持

### 2. **搜索性能**
- 毫秒级响应时间
- 高并发查询支持
- 智能缓存策略
- 索引优化设计

### 3. **用户体验**
- 实时搜索建议
- 搜索历史记录
- 热门搜索推荐
- 高级筛选功能

## 🔍 搜索架构设计

### 1. **搜索系统架构**

```
┌─────────────────────────────────────────────────────────────┐
│                    电影搜索系统架构                          │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │  搜索接口    │  │  搜索引擎    │  │  索引管理    │          │
│  │             │  │             │  │             │          │
│  │ • 关键词搜索 │  │ • 全文搜索   │  │ • 倒排索引   │          │
│  │ • 高级搜索   │  │ • 模糊匹配   │  │ • 分词索引   │          │
│  │ • 搜索建议   │  │ • 相关性排序 │  │ • 实时更新   │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
│                              │                              │
│                              ▼                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                  搜索优化层                              │ │
│  │                                                         │ │
│  │ • 查询缓存    • 结果缓存    • 热词缓存    • 统计分析     │ │
│  │ • 搜索日志    • 性能监控    • 用户行为    • 智能推荐     │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 2. **搜索请求结构**

```go
// 搜索请求参数
type MovieSearchRequest struct {
    // 基础搜索参数
    Query    string `json:"query" form:"query" binding:"required"`
    Page     int    `json:"page" form:"page" binding:"min=1"`
    PageSize int    `json:"page_size" form:"page_size" binding:"min=1,max=50"`
    
    // 搜索类型
    SearchType string `json:"search_type" form:"search_type"` // all, title, person, company
    
    // 排序参数
    SortBy    string `json:"sort_by" form:"sort_by"`       // relevance, rating, release_date, popularity
    SortOrder string `json:"sort_order" form:"sort_order"` // asc, desc
    
    // 筛选参数
    Genres      []string `json:"genres" form:"genres"`
    YearFrom    int      `json:"year_from" form:"year_from"`
    YearTo      int      `json:"year_to" form:"year_to"`
    RatingFrom  float64  `json:"rating_from" form:"rating_from"`
    RatingTo    float64  `json:"rating_to" form:"rating_to"`
    Countries   []string `json:"countries" form:"countries"`
    Languages   []string `json:"languages" form:"languages"`
    
    // 搜索选项
    IncludeAdult bool `json:"include_adult" form:"include_adult"`
    ExactMatch   bool `json:"exact_match" form:"exact_match"`
    
    // 用户信息（用于个性化）
    UserID string `json:"-"` // 从认证中获取
}

// 搜索响应结构
type MovieSearchResponse struct {
    Success      bool                `json:"success"`
    Message      string              `json:"message,omitempty"`
    Data         *MovieSearchData    `json:"data,omitempty"`
    Pagination   *PaginationInfo     `json:"pagination,omitempty"`
    SearchInfo   *SearchInfo         `json:"search_info,omitempty"`
    Suggestions  []SearchSuggestion  `json:"suggestions,omitempty"`
}

// 搜索数据
type MovieSearchData struct {
    Movies      []MovieSearchResult `json:"movies"`
    Total       int64               `json:"total"`
    SearchTime  float64             `json:"search_time_ms"`
}

// 搜索结果项
type MovieSearchResult struct {
    ID            string    `json:"id"`
    Title         string    `json:"title"`
    OriginalTitle string    `json:"original_title,omitempty"`
    PosterURL     string    `json:"poster_url"`
    Overview      string    `json:"overview"`
    ReleaseDate   string    `json:"release_date"`
    Rating        float64   `json:"rating"`
    VoteCount     int       `json:"vote_count"`
    Popularity    float64   `json:"popularity"`
    Genres        []Genre   `json:"genres"`
    
    // 搜索相关
    Relevance     float64   `json:"relevance"`
    MatchedFields []string  `json:"matched_fields"`
    Highlights    []string  `json:"highlights,omitempty"`
}

// 搜索信息
type SearchInfo struct {
    Query           string  `json:"query"`
    ProcessedQuery  string  `json:"processed_query"`
    SearchTime      float64 `json:"search_time_ms"`
    TotalResults    int64   `json:"total_results"`
    CorrectedQuery  string  `json:"corrected_query,omitempty"`
    DidYouMean      string  `json:"did_you_mean,omitempty"`
}

// 搜索建议
type SearchSuggestion struct {
    Text        string  `json:"text"`
    Type        string  `json:"type"` // movie, person, genre
    Count       int     `json:"count"`
    Popularity  float64 `json:"popularity"`
}
```

## 🔧 搜索引擎实现

### 1. **核心搜索服务**

```go
type MovieSearchService struct {
    searchEngine   SearchEngine
    movieRepo      MovieRepository
    cacheStore     CacheStore
    suggestionService SuggestionService
    logger         *logrus.Logger
    metrics        *SearchMetrics
}

func NewMovieSearchService(
    searchEngine SearchEngine,
    movieRepo MovieRepository,
    cacheStore CacheStore,
    suggestionService SuggestionService,
) *MovieSearchService {
    return &MovieSearchService{
        searchEngine:      searchEngine,
        movieRepo:         movieRepo,
        cacheStore:        cacheStore,
        suggestionService: suggestionService,
        logger:            logrus.New(),
        metrics:           NewSearchMetrics(),
    }
}

// 搜索电影
func (mss *MovieSearchService) SearchMovies(ctx context.Context, req *MovieSearchRequest) (*MovieSearchResponse, error) {
    start := time.Now()
    defer func() {
        mss.metrics.ObserveSearchDuration(time.Since(start))
    }()

    // 验证和预处理请求
    if err := mss.validateAndPreprocessRequest(req); err != nil {
        mss.metrics.IncInvalidRequests()
        return &MovieSearchResponse{
            Success: false,
            Message: err.Error(),
        }, nil
    }

    // 生成缓存键
    cacheKey := mss.generateCacheKey(req)

    // 尝试从缓存获取结果
    if cachedResult, err := mss.getFromCache(ctx, cacheKey); err == nil {
        mss.metrics.IncCacheHits()
        return cachedResult, nil
    }
    mss.metrics.IncCacheMisses()

    // 执行搜索
    searchResult, err := mss.executeSearch(ctx, req)
    if err != nil {
        mss.logger.Errorf("Search execution failed: %v", err)
        mss.metrics.IncSearchErrors()
        return nil, errors.New("搜索执行失败")
    }

    // 获取搜索建议
    suggestions, err := mss.getSuggestions(ctx, req.Query)
    if err != nil {
        mss.logger.Errorf("Failed to get suggestions: %v", err)
        // 建议获取失败不影响主要搜索功能
    }

    // 构建响应
    response := &MovieSearchResponse{
        Success:     true,
        Data:        searchResult.Data,
        Pagination:  searchResult.Pagination,
        SearchInfo:  searchResult.SearchInfo,
        Suggestions: suggestions,
    }

    // 异步缓存结果
    go func() {
        if err := mss.cacheResult(context.Background(), cacheKey, response); err != nil {
            mss.logger.Errorf("Failed to cache search result: %v", err)
        }
    }()

    // 异步记录搜索日志
    go mss.logSearchQuery(context.Background(), req, searchResult.Data.Total)

    mss.metrics.IncSuccessfulSearches()
    return response, nil
}

// 执行搜索
func (mss *MovieSearchService) executeSearch(ctx context.Context, req *MovieSearchRequest) (*SearchResult, error) {
    searchStart := time.Now()

    // 构建搜索查询
    searchQuery := mss.buildSearchQuery(req)

    // 执行搜索引擎查询
    engineResult, err := mss.searchEngine.Search(ctx, searchQuery)
    if err != nil {
        return nil, err
    }

    searchTime := time.Since(searchStart).Seconds() * 1000 // 转换为毫秒

    // 转换搜索结果
    movies := mss.convertSearchResults(engineResult.Hits)

    // 构建分页信息
    pagination := &PaginationInfo{
        CurrentPage:  req.Page,
        PageSize:     req.PageSize,
        TotalPages:   int(math.Ceil(float64(engineResult.Total) / float64(req.PageSize))),
        TotalItems:   engineResult.Total,
        HasNext:      req.Page < int(math.Ceil(float64(engineResult.Total)/float64(req.PageSize))),
        HasPrevious:  req.Page > 1,
    }

    // 构建搜索信息
    searchInfo := &SearchInfo{
        Query:          req.Query,
        ProcessedQuery: searchQuery.ProcessedQuery,
        SearchTime:     searchTime,
        TotalResults:   engineResult.Total,
        CorrectedQuery: engineResult.CorrectedQuery,
        DidYouMean:     engineResult.DidYouMean,
    }

    return &SearchResult{
        Data: &MovieSearchData{
            Movies:     movies,
            Total:      engineResult.Total,
            SearchTime: searchTime,
        },
        Pagination: pagination,
        SearchInfo: searchInfo,
    }, nil
}

// 构建搜索查询
func (mss *MovieSearchService) buildSearchQuery(req *MovieSearchRequest) *SearchQuery {
    query := &SearchQuery{
        Query:      req.Query,
        SearchType: req.SearchType,
        Page:       req.Page,
        PageSize:   req.PageSize,
        SortBy:     req.SortBy,
        SortOrder:  req.SortOrder,
        Filters:    make(map[string]interface{}),
    }

    // 处理查询文本
    query.ProcessedQuery = mss.preprocessQuery(req.Query)

    // 添加筛选条件
    if len(req.Genres) > 0 {
        query.Filters["genres"] = req.Genres
    }
    
    if req.YearFrom > 0 || req.YearTo > 0 {
        yearFilter := make(map[string]int)
        if req.YearFrom > 0 {
            yearFilter["gte"] = req.YearFrom
        }
        if req.YearTo > 0 {
            yearFilter["lte"] = req.YearTo
        }
        query.Filters["release_year"] = yearFilter
    }

    if req.RatingFrom > 0 || req.RatingTo > 0 {
        ratingFilter := make(map[string]float64)
        if req.RatingFrom > 0 {
            ratingFilter["gte"] = req.RatingFrom
        }
        if req.RatingTo > 0 {
            ratingFilter["lte"] = req.RatingTo
        }
        query.Filters["rating"] = ratingFilter
    }

    if len(req.Countries) > 0 {
        query.Filters["countries"] = req.Countries
    }

    if len(req.Languages) > 0 {
        query.Filters["languages"] = req.Languages
    }

    if !req.IncludeAdult {
        query.Filters["adult"] = false
    }

    return query
}

// 预处理查询文本
func (mss *MovieSearchService) preprocessQuery(query string) string {
    // 去除多余空格
    query = strings.TrimSpace(query)
    query = regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")

    // 转换为小写
    query = strings.ToLower(query)

    // 移除特殊字符（保留基本标点）
    query = regexp.MustCompile(`[^\w\s\-\.\,\!\?]`).ReplaceAllString(query, "")

    return query
}
```

### 2. **搜索引擎接口**

```go
type SearchEngine interface {
    Search(ctx context.Context, query *SearchQuery) (*SearchEngineResult, error)
    Index(ctx context.Context, document *SearchDocument) error
    Update(ctx context.Context, id string, document *SearchDocument) error
    Delete(ctx context.Context, id string) error
    Suggest(ctx context.Context, query string, limit int) ([]SearchSuggestion, error)
}

// Elasticsearch实现
type ElasticsearchEngine struct {
    client *elasticsearch.Client
    index  string
    logger *logrus.Logger
}

func NewElasticsearchEngine(client *elasticsearch.Client, index string) *ElasticsearchEngine {
    return &ElasticsearchEngine{
        client: client,
        index:  index,
        logger: logrus.New(),
    }
}

func (es *ElasticsearchEngine) Search(ctx context.Context, query *SearchQuery) (*SearchEngineResult, error) {
    // 构建Elasticsearch查询
    esQuery := es.buildElasticsearchQuery(query)

    // 执行搜索
    res, err := es.client.Search(
        es.client.Search.WithContext(ctx),
        es.client.Search.WithIndex(es.index),
        es.client.Search.WithBody(esQuery),
        es.client.Search.WithTrackTotalHits(true),
    )
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    // 解析响应
    var response ElasticsearchResponse
    if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
        return nil, err
    }

    // 转换为通用结果格式
    result := &SearchEngineResult{
        Total: response.Hits.Total.Value,
        Hits:  make([]SearchHit, len(response.Hits.Hits)),
    }

    for i, hit := range response.Hits.Hits {
        result.Hits[i] = SearchHit{
            ID:        hit.ID,
            Score:     hit.Score,
            Source:    hit.Source,
            Highlight: hit.Highlight,
        }
    }

    return result, nil
}

func (es *ElasticsearchEngine) buildElasticsearchQuery(query *SearchQuery) io.Reader {
    esQuery := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []map[string]interface{}{
                    {
                        "multi_match": map[string]interface{}{
                            "query":  query.ProcessedQuery,
                            "fields": []string{"title^3", "original_title^2", "overview", "genres.name", "cast.name", "crew.name"},
                            "type":   "best_fields",
                            "fuzziness": "AUTO",
                        },
                    },
                },
                "filter": es.buildFilters(query.Filters),
            },
        },
        "highlight": map[string]interface{}{
            "fields": map[string]interface{}{
                "title":          map[string]interface{}{},
                "original_title": map[string]interface{}{},
                "overview":       map[string]interface{}{},
            },
        },
        "sort": es.buildSort(query.SortBy, query.SortOrder),
        "from": (query.Page - 1) * query.PageSize,
        "size": query.PageSize,
    }

    jsonQuery, _ := json.Marshal(esQuery)
    return bytes.NewReader(jsonQuery)
}
```

### 3. **搜索建议服务**

```go
type SuggestionService struct {
    redis      *redis.Client
    movieRepo  MovieRepository
    logger     *logrus.Logger
}

func NewSuggestionService(redis *redis.Client, movieRepo MovieRepository) *SuggestionService {
    return &SuggestionService{
        redis:     redis,
        movieRepo: movieRepo,
        logger:    logrus.New(),
    }
}

// 获取搜索建议
func (ss *SuggestionService) GetSuggestions(ctx context.Context, query string, limit int) ([]SearchSuggestion, error) {
    if len(query) < 2 {
        return []SearchSuggestion{}, nil
    }

    var suggestions []SearchSuggestion

    // 获取电影标题建议
    movieSuggestions, err := ss.getMovieSuggestions(ctx, query, limit/2)
    if err != nil {
        ss.logger.Errorf("Failed to get movie suggestions: %v", err)
    } else {
        suggestions = append(suggestions, movieSuggestions...)
    }

    // 获取人物建议
    personSuggestions, err := ss.getPersonSuggestions(ctx, query, limit/4)
    if err != nil {
        ss.logger.Errorf("Failed to get person suggestions: %v", err)
    } else {
        suggestions = append(suggestions, personSuggestions...)
    }

    // 获取类型建议
    genreSuggestions, err := ss.getGenreSuggestions(ctx, query, limit/4)
    if err != nil {
        ss.logger.Errorf("Failed to get genre suggestions: %v", err)
    } else {
        suggestions = append(suggestions, genreSuggestions...)
    }

    // 按相关性排序
    sort.Slice(suggestions, func(i, j int) bool {
        return suggestions[i].Popularity > suggestions[j].Popularity
    })

    if len(suggestions) > limit {
        suggestions = suggestions[:limit]
    }

    return suggestions, nil
}

// 获取热门搜索
func (ss *SuggestionService) GetHotSearches(ctx context.Context, limit int) ([]SearchSuggestion, error) {
    key := "hot_searches"
    
    // 从Redis获取热门搜索
    results, err := ss.redis.ZRevRange(ctx, key, 0, int64(limit-1)).Result()
    if err != nil {
        return nil, err
    }

    suggestions := make([]SearchSuggestion, len(results))
    for i, result := range results {
        score, _ := ss.redis.ZScore(ctx, key, result).Result()
        suggestions[i] = SearchSuggestion{
            Text:       result,
            Type:       "hot",
            Popularity: score,
        }
    }

    return suggestions, nil
}

// 记录搜索查询
func (ss *SuggestionService) RecordSearchQuery(ctx context.Context, query string) {
    if len(query) < 2 {
        return
    }

    key := "hot_searches"
    
    // 增加搜索次数
    ss.redis.ZIncrBy(ctx, key, 1, query)
    
    // 设置过期时间（7天）
    ss.redis.Expire(ctx, key, 7*24*time.Hour)
}
```

## 📊 性能监控

### 1. **搜索指标收集**

```go
type SearchMetrics struct {
    searchRequests     *prometheus.CounterVec
    searchDuration     prometheus.Histogram
    cacheHitRate       *prometheus.CounterVec
    searchErrors       *prometheus.CounterVec
    queryComplexity    prometheus.Histogram
    resultCount        prometheus.Histogram
}

func NewSearchMetrics() *SearchMetrics {
    return &SearchMetrics{
        searchRequests: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "movie_search_requests_total",
                Help: "Total number of movie search requests",
            },
            []string{"search_type", "status"},
        ),
        searchDuration: prometheus.NewHistogram(
            prometheus.HistogramOpts{
                Name: "movie_search_duration_seconds",
                Help: "Duration of movie search requests",
            },
        ),
        cacheHitRate: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "movie_search_cache_operations_total",
                Help: "Total number of search cache operations",
            },
            []string{"type"},
        ),
        searchErrors: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "movie_search_errors_total",
                Help: "Total number of search errors",
            },
            []string{"error_type"},
        ),
        queryComplexity: prometheus.NewHistogram(
            prometheus.HistogramOpts{
                Name: "movie_search_query_complexity",
                Help: "Complexity of search queries",
            },
        ),
        resultCount: prometheus.NewHistogram(
            prometheus.HistogramOpts{
                Name: "movie_search_result_count",
                Help: "Number of search results returned",
            },
        ),
    }
}
```

## 🔧 HTTP处理器

### 1. **搜索API端点**

```go
func (mc *MovieController) SearchMovies(c *gin.Context) {
    var req MovieSearchRequest
    
    // 绑定查询参数
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": "请求参数错误",
            "error":   err.Error(),
        })
        return
    }

    // 获取用户ID（用于个性化）
    req.UserID = mc.getUserIDFromContext(c)

    // 调用搜索服务
    response, err := mc.movieSearchService.SearchMovies(c.Request.Context(), &req)
    if err != nil {
        mc.logger.Errorf("Failed to search movies: %v", err)
        c.JSON(500, gin.H{
            "success": false,
            "message": "搜索失败，请稍后重试",
        })
        return
    }

    // 设置缓存头
    c.Header("Cache-Control", "public, max-age=300") // 5分钟缓存

    c.JSON(200, response)
}

// 获取搜索建议
func (mc *MovieController) GetSearchSuggestions(c *gin.Context) {
    query := c.Query("q")
    if len(query) < 2 {
        c.JSON(200, gin.H{
            "success":     true,
            "suggestions": []SearchSuggestion{},
        })
        return
    }

    suggestions, err := mc.suggestionService.GetSuggestions(c.Request.Context(), query, 10)
    if err != nil {
        mc.logger.Errorf("Failed to get suggestions: %v", err)
        c.JSON(500, gin.H{
            "success": false,
            "message": "获取建议失败",
        })
        return
    }

    c.Header("Cache-Control", "public, max-age=60") // 1分钟缓存
    c.JSON(200, gin.H{
        "success":     true,
        "suggestions": suggestions,
    })
}

// 获取热门搜索
func (mc *MovieController) GetHotSearches(c *gin.Context) {
    hotSearches, err := mc.suggestionService.GetHotSearches(c.Request.Context(), 20)
    if err != nil {
        mc.logger.Errorf("Failed to get hot searches: %v", err)
        c.JSON(500, gin.H{
            "success": false,
            "message": "获取热门搜索失败",
        })
        return
    }

    c.Header("Cache-Control", "public, max-age=3600") // 1小时缓存
    c.JSON(200, gin.H{
        "success": true,
        "data":    hotSearches,
    })
}
```

## 📝 总结

电影搜索功能为MovieInfo项目提供了强大的电影发现能力：

**核心功能**：
1. **全文搜索**：支持电影标题、演员、导演等多字段搜索
2. **智能匹配**：模糊搜索、拼写纠错、相关性排序
3. **高级筛选**：多维度筛选条件组合
4. **搜索建议**：实时建议、热门搜索、搜索历史

**性能优化**：
- Elasticsearch全文搜索引擎
- 多层缓存策略
- 查询优化和索引设计
- 异步处理和预加载

**用户体验**：
- 毫秒级搜索响应
- 智能搜索建议
- 丰富的筛选选项
- 个性化搜索结果

下一步，我们将实现电影分类管理功能，为电影内容提供有序的组织结构。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第35步：电影分类管理](35-category-management.md)
