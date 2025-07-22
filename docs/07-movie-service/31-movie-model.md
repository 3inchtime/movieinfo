# 5.2 电影服务实现

## 5.2.1 概述

电影服务是MovieInfo项目的核心业务模块，负责电影数据管理、搜索、分类、推荐等功能。作为内容管理的核心，电影服务需要处理复杂的数据关系、高效的搜索算法和智能的推荐系统。

## 5.2.2 为什么需要专业的电影服务？

### 5.2.2.1 **数据复杂性**
- **多维度数据**：电影包含标题、导演、演员、分类等多维度信息
- **关联关系**：电影与演员、分类、评分等实体的复杂关联
- **外部数据源**：集成TMDB、IMDB等外部电影数据源
- **数据同步**：保持内部数据与外部数据源的同步

### 5.2.2.2 **搜索需求**
- **全文搜索**：支持电影标题、导演、演员的全文搜索
- **多条件筛选**：按分类、年份、评分等条件筛选
- **智能排序**：按相关度、评分、热度等多种方式排序
- **搜索优化**：搜索性能优化和结果缓存

### 5.2.2.3 **推荐算法**
- **协同过滤**：基于用户行为的协同过滤推荐
- **内容推荐**：基于电影内容特征的推荐
- **热门推荐**：基于热度和趋势的推荐
- **个性化推荐**：基于用户偏好的个性化推荐

### 5.2.2.4 **性能要求**
- **高并发访问**：支持大量用户同时访问电影数据
- **快速响应**：毫秒级的电影信息查询响应
- **缓存策略**：多层缓存提升访问性能
- **数据预加载**：热门数据的预加载和预热

## 5.2.3 电影服务架构设计

### 5.2.3.1 **服务层次结构**

```
电影服务架构
├── 服务接口层 (Service Interface)
│   ├── 电影服务接口 (MovieService)
│   ├── 搜索服务接口 (SearchService)
│   ├── 推荐服务接口 (RecommendationService)
│   └── 分类服务接口 (CategoryService)
├── 业务逻辑层 (Business Logic)
│   ├── 电影管理逻辑 (Movie Management)
│   ├── 搜索算法逻辑 (Search Algorithm)
│   ├── 推荐算法逻辑 (Recommendation Algorithm)
│   └── 数据同步逻辑 (Data Sync Logic)
├── 数据访问层 (Data Access)
│   ├── 电影Repository (MovieRepository)
│   ├── 分类Repository (CategoryRepository)
│   ├── 演员Repository (ActorRepository)
│   └── 搜索Repository (SearchRepository)
├── 外部集成层 (External Integration)
│   ├── TMDB API集成 (TMDB Integration)
│   ├── IMDB API集成 (IMDB Integration)
│   ├── 图片存储服务 (Image Storage)
│   └── 搜索引擎集成 (Search Engine)
├── 缓存层 (Cache Layer)
│   ├── 电影缓存 (Movie Cache)
│   ├── 搜索缓存 (Search Cache)
│   ├── 推荐缓存 (Recommendation Cache)
│   └── 热门缓存 (Hot Cache)
└── 算法层 (Algorithm Layer)
    ├── 搜索算法 (Search Algorithm)
    ├── 推荐算法 (Recommendation Algorithm)
    ├── 排序算法 (Ranking Algorithm)
    └── 相似度算法 (Similarity Algorithm)
```

### 5.2.3.2 **电影服务接口定义**

#### 5.2.3.2.1 核心电影服务接口
```go
// internal/services/movie/interface.go
package movie

import (
    "context"
    
    "github.com/yourname/movieinfo/internal/models"
    "github.com/yourname/movieinfo/internal/repository"
)

// MovieService 电影服务接口
type MovieService interface {
    // 电影基础操作
    CreateMovie(ctx context.Context, req *CreateMovieRequest) (*models.Movie, error)
    GetMovie(ctx context.Context, movieID int64) (*models.Movie, error)
    GetMovieByIMDBID(ctx context.Context, imdbID string) (*models.Movie, error)
    UpdateMovie(ctx context.Context, movieID int64, req *UpdateMovieRequest) (*models.Movie, error)
    DeleteMovie(ctx context.Context, movieID int64) error
    
    // 电影列表和查询
    ListMovies(ctx context.Context, opts *ListMoviesOptions) (*repository.PaginationResult[models.Movie], error)
    GetMoviesByCategory(ctx context.Context, categoryID int64, opts *ListMoviesOptions) (*repository.PaginationResult[models.Movie], error)
    GetMoviesByActor(ctx context.Context, actorID int64, opts *ListMoviesOptions) (*repository.PaginationResult[models.Movie], error)
    GetMoviesByDirector(ctx context.Context, director string, opts *ListMoviesOptions) (*repository.PaginationResult[models.Movie], error)
    
    // 电影统计和排行
    GetTopRatedMovies(ctx context.Context, limit int) ([]*models.Movie, error)
    GetMostPopularMovies(ctx context.Context, limit int) ([]*models.Movie, error)
    GetLatestMovies(ctx context.Context, limit int) ([]*models.Movie, error)
    GetTrendingMovies(ctx context.Context, timeWindow string, limit int) ([]*models.Movie, error)
    
    // 电影详情和关联
    GetMovieWithDetails(ctx context.Context, movieID int64) (*MovieDetails, error)
    GetMovieCategories(ctx context.Context, movieID int64) ([]*models.Category, error)
    GetMovieActors(ctx context.Context, movieID int64) ([]*models.Actor, error)
    GetSimilarMovies(ctx context.Context, movieID int64, limit int) ([]*models.Movie, error)
    
    // 电影操作
    IncrementViewCount(ctx context.Context, movieID int64) error
    UpdateMovieRating(ctx context.Context, movieID int64, rating float64, userID int64) error
    AddMovieToFavorites(ctx context.Context, movieID int64, userID int64) error
    RemoveMovieFromFavorites(ctx context.Context, movieID int64, userID int64) error
    
    // 数据同步
    SyncMovieFromTMDB(ctx context.Context, tmdbID int) (*models.Movie, error)
    SyncMovieFromIMDB(ctx context.Context, imdbID string) (*models.Movie, error)
    BatchSyncMovies(ctx context.Context, req *BatchSyncRequest) (*BatchSyncResponse, error)
}

// SearchService 搜索服务接口
type SearchService interface {
    // 电影搜索
    SearchMovies(ctx context.Context, query string, opts *SearchOptions) (*SearchResult, error)
    SearchMoviesByTitle(ctx context.Context, title string, opts *SearchOptions) (*SearchResult, error)
    SearchMoviesByKeyword(ctx context.Context, keyword string, opts *SearchOptions) (*SearchResult, error)
    
    // 高级搜索
    AdvancedSearch(ctx context.Context, req *AdvancedSearchRequest) (*SearchResult, error)
    FilterMovies(ctx context.Context, filters *MovieFilters, opts *SearchOptions) (*SearchResult, error)
    
    // 搜索建议
    GetSearchSuggestions(ctx context.Context, query string, limit int) ([]string, error)
    GetPopularSearches(ctx context.Context, limit int) ([]string, error)
    
    // 搜索统计
    RecordSearch(ctx context.Context, query string, userID int64, resultCount int) error
    GetSearchStats(ctx context.Context, timeRange string) (*SearchStats, error)
}

// RecommendationService 推荐服务接口
type RecommendationService interface {
    // 个性化推荐
    GetPersonalizedRecommendations(ctx context.Context, userID int64, limit int) ([]*models.Movie, error)
    GetRecommendationsBasedOnHistory(ctx context.Context, userID int64, limit int) ([]*models.Movie, error)
    GetRecommendationsBasedOnRatings(ctx context.Context, userID int64, limit int) ([]*models.Movie, error)
    
    // 内容推荐
    GetSimilarMovies(ctx context.Context, movieID int64, limit int) ([]*models.Movie, error)
    GetMoviesByGenre(ctx context.Context, userID int64, genreID int64, limit int) ([]*models.Movie, error)
    GetMoviesByDirector(ctx context.Context, userID int64, director string, limit int) ([]*models.Movie, error)
    
    // 热门推荐
    GetTrendingRecommendations(ctx context.Context, userID int64, limit int) ([]*models.Movie, error)
    GetNewReleaseRecommendations(ctx context.Context, userID int64, limit int) ([]*models.Movie, error)
    
    // 推荐反馈
    RecordRecommendationClick(ctx context.Context, userID int64, movieID int64, recommendationType string) error
    RecordRecommendationFeedback(ctx context.Context, userID int64, movieID int64, feedback string) error
    
    // 推荐模型
    TrainRecommendationModel(ctx context.Context) error
    UpdateUserProfile(ctx context.Context, userID int64) error
}
```

#### 5.2.3.2.2 请求和响应结构
```go
// internal/services/movie/types.go
package movie

import (
    "time"
    
    "github.com/yourname/movieinfo/internal/models"
)

// CreateMovieRequest 创建电影请求
type CreateMovieRequest struct {
    Title         string   `json:"title" validate:"required,min=1,max=200"`
    OriginalTitle string   `json:"original_title,omitempty" validate:"omitempty,max=200"`
    Director      string   `json:"director,omitempty" validate:"omitempty,max=100"`
    ReleaseYear   *int     `json:"release_year,omitempty" validate:"omitempty,min=1888,max=2030"`
    Duration      *int     `json:"duration,omitempty" validate:"omitempty,min=1,max=1000"`
    Country       string   `json:"country,omitempty" validate:"omitempty,max=100"`
    Language      string   `json:"language,omitempty" validate:"omitempty,max=50"`
    Plot          string   `json:"plot,omitempty" validate:"omitempty,max=2000"`
    Tagline       string   `json:"tagline,omitempty" validate:"omitempty,max=255"`
    PosterURL     string   `json:"poster_url,omitempty" validate:"omitempty,url"`
    BackdropURL   string   `json:"backdrop_url,omitempty" validate:"omitempty,url"`
    TrailerURL    string   `json:"trailer_url,omitempty" validate:"omitempty,url"`
    IMDBID        string   `json:"imdb_id,omitempty" validate:"omitempty,len=9"`
    TMDBID        *int     `json:"tmdb_id,omitempty"`
    CategoryIDs   []int64  `json:"category_ids,omitempty"`
    ActorIDs      []int64  `json:"actor_ids,omitempty"`
    Keywords      []string `json:"keywords,omitempty"`
}

// UpdateMovieRequest 更新电影请求
type UpdateMovieRequest struct {
    Title         *string  `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
    OriginalTitle *string  `json:"original_title,omitempty" validate:"omitempty,max=200"`
    Director      *string  `json:"director,omitempty" validate:"omitempty,max=100"`
    ReleaseYear   *int     `json:"release_year,omitempty" validate:"omitempty,min=1888,max=2030"`
    Duration      *int     `json:"duration,omitempty" validate:"omitempty,min=1,max=1000"`
    Country       *string  `json:"country,omitempty" validate:"omitempty,max=100"`
    Language      *string  `json:"language,omitempty" validate:"omitempty,max=50"`
    Plot          *string  `json:"plot,omitempty" validate:"omitempty,max=2000"`
    Tagline       *string  `json:"tagline,omitempty" validate:"omitempty,max=255"`
    PosterURL     *string  `json:"poster_url,omitempty" validate:"omitempty,url"`
    BackdropURL   *string  `json:"backdrop_url,omitempty" validate:"omitempty,url"`
    TrailerURL    *string  `json:"trailer_url,omitempty" validate:"omitempty,url"`
    CategoryIDs   []int64  `json:"category_ids,omitempty"`
    ActorIDs      []int64  `json:"actor_ids,omitempty"`
    Keywords      []string `json:"keywords,omitempty"`
    Status        *models.MovieStatus `json:"status,omitempty"`
}

// MovieDetails 电影详情
type MovieDetails struct {
    *models.Movie
    Categories    []*models.Category `json:"categories"`
    Actors        []*models.Actor    `json:"actors"`
    Ratings       []*models.Rating   `json:"ratings,omitempty"`
    Comments      []*models.Comment  `json:"comments,omitempty"`
    SimilarMovies []*models.Movie    `json:"similar_movies,omitempty"`
    UserRating    *float64           `json:"user_rating,omitempty"`
    IsFavorite    bool               `json:"is_favorite"`
}

// ListMoviesOptions 电影列表选项
type ListMoviesOptions struct {
    Page       int                    `json:"page" validate:"min=1"`
    PageSize   int                    `json:"page_size" validate:"min=1,max=100"`
    OrderBy    string                 `json:"order_by" validate:"omitempty,oneof=id title release_year average_rating view_count created_at"`
    Order      string                 `json:"order" validate:"omitempty,oneof=asc desc"`
    Status     *models.MovieStatus    `json:"status,omitempty"`
    CategoryID *int64                 `json:"category_id,omitempty"`
    Year       *int                   `json:"year,omitempty"`
    MinRating  *float64               `json:"min_rating,omitempty"`
    MaxRating  *float64               `json:"max_rating,omitempty"`
    Filters    map[string]interface{} `json:"filters,omitempty"`
}

// SearchOptions 搜索选项
type SearchOptions struct {
    Page      int      `json:"page" validate:"min=1"`
    PageSize  int      `json:"page_size" validate:"min=1,max=100"`
    OrderBy   string   `json:"order_by" validate:"omitempty,oneof=relevance rating popularity release_date"`
    Fields    []string `json:"fields,omitempty"` // title, director, actor, plot
    Highlight bool     `json:"highlight"`
}

// SearchResult 搜索结果
type SearchResult struct {
    Movies     []*models.Movie `json:"movies"`
    Total      int64           `json:"total"`
    Page       int             `json:"page"`
    PageSize   int             `json:"page_size"`
    Query      string          `json:"query"`
    Took       int64           `json:"took"` // 搜索耗时(毫秒)
    Highlights map[int64]map[string][]string `json:"highlights,omitempty"`
}

// AdvancedSearchRequest 高级搜索请求
type AdvancedSearchRequest struct {
    Title       string    `json:"title,omitempty"`
    Director    string    `json:"director,omitempty"`
    Actor       string    `json:"actor,omitempty"`
    Genre       string    `json:"genre,omitempty"`
    YearFrom    *int      `json:"year_from,omitempty"`
    YearTo      *int      `json:"year_to,omitempty"`
    RatingFrom  *float64  `json:"rating_from,omitempty"`
    RatingTo    *float64  `json:"rating_to,omitempty"`
    Country     string    `json:"country,omitempty"`
    Language    string    `json:"language,omitempty"`
    Keywords    []string  `json:"keywords,omitempty"`
    Options     *SearchOptions `json:"options,omitempty"`
}

// MovieFilters 电影过滤器
type MovieFilters struct {
    CategoryIDs []int64   `json:"category_ids,omitempty"`
    ActorIDs    []int64   `json:"actor_ids,omitempty"`
    Directors   []string  `json:"directors,omitempty"`
    Countries   []string  `json:"countries,omitempty"`
    Languages   []string  `json:"languages,omitempty"`
    YearRange   *YearRange `json:"year_range,omitempty"`
    RatingRange *RatingRange `json:"rating_range,omitempty"`
    Keywords    []string  `json:"keywords,omitempty"`
}

// YearRange 年份范围
type YearRange struct {
    From *int `json:"from,omitempty"`
    To   *int `json:"to,omitempty"`
}

// RatingRange 评分范围
type RatingRange struct {
    From *float64 `json:"from,omitempty"`
    To   *float64 `json:"to,omitempty"`
}

// SearchStats 搜索统计
type SearchStats struct {
    TotalSearches    int64             `json:"total_searches"`
    UniqueSearches   int64             `json:"unique_searches"`
    AvgResultCount   float64           `json:"avg_result_count"`
    TopQueries       []QueryStat       `json:"top_queries"`
    SearchTrends     []TrendPoint      `json:"search_trends"`
    NoResultQueries  []string          `json:"no_result_queries"`
}

// QueryStat 查询统计
type QueryStat struct {
    Query string `json:"query"`
    Count int64  `json:"count"`
}

// TrendPoint 趋势点
type TrendPoint struct {
    Date  time.Time `json:"date"`
    Count int64     `json:"count"`
}

// BatchSyncRequest 批量同步请求
type BatchSyncRequest struct {
    Source   string   `json:"source" validate:"required,oneof=tmdb imdb"`
    IDs      []string `json:"ids" validate:"required,min=1,max=100"`
    ForceUpdate bool  `json:"force_update"`
}

// BatchSyncResponse 批量同步响应
type BatchSyncResponse struct {
    Total     int                `json:"total"`
    Success   int                `json:"success"`
    Failed    int                `json:"failed"`
    Results   []SyncResult       `json:"results"`
    Errors    []SyncError        `json:"errors,omitempty"`
}

// SyncResult 同步结果
type SyncResult struct {
    ID       string        `json:"id"`
    MovieID  int64         `json:"movie_id"`
    Status   string        `json:"status"` // created, updated, skipped
    Movie    *models.Movie `json:"movie,omitempty"`
}

// SyncError 同步错误
type SyncError struct {
    ID    string `json:"id"`
    Error string `json:"error"`
}
```

### 5.2.3.3 **电影服务实现**

#### 5.2.3.3.1 核心电影服务实现
```go
// internal/services/movie/service.go
package movie

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

// ServiceImpl 电影服务实现
type ServiceImpl struct {
    movieRepo     repository.MovieRepository
    categoryRepo  repository.CategoryRepository
    actorRepo     repository.ActorRepository
    searchService SearchService
    recomService  RecommendationService
    cache         cache.Cache
    validator     validator.Validator
    logger        logger.Logger
    
    // 外部服务
    tmdbClient TMDBClient
    imdbClient IMDBClient
    
    // 配置
    config *Config
}

// Config 电影服务配置
type Config struct {
    CacheExpiry        time.Duration `yaml:"cache_expiry"`
    SearchCacheExpiry  time.Duration `yaml:"search_cache_expiry"`
    HotMoviesLimit     int           `yaml:"hot_movies_limit"`
    SimilarMoviesLimit int           `yaml:"similar_movies_limit"`
    EnableAutoSync     bool          `yaml:"enable_auto_sync"`
    SyncBatchSize      int           `yaml:"sync_batch_size"`
}

// NewService 创建电影服务
func NewService(
    movieRepo repository.MovieRepository,
    categoryRepo repository.CategoryRepository,
    actorRepo repository.ActorRepository,
    searchService SearchService,
    recomService RecommendationService,
    cache cache.Cache,
    validator validator.Validator,
    tmdbClient TMDBClient,
    imdbClient IMDBClient,
    config *Config,
) MovieService {
    return &ServiceImpl{
        movieRepo:     movieRepo,
        categoryRepo:  categoryRepo,
        actorRepo:     actorRepo,
        searchService: searchService,
        recomService:  recomService,
        cache:         cache,
        validator:     validator,
        logger:        logger.GetGlobalLogger(),
        tmdbClient:    tmdbClient,
        imdbClient:    imdbClient,
        config:        config,
    }
}

// CreateMovie 创建电影
func (s *ServiceImpl) CreateMovie(ctx context.Context, req *CreateMovieRequest) (*models.Movie, error) {
    // 验证请求参数
    if err := s.validator.Validate(req); err != nil {
        return nil, errors.ValidationFailed(err.Error())
    }
    
    s.logger.WithContext(ctx).Info("Creating movie",
        logger.String("title", req.Title),
        logger.String("director", req.Director),
    )
    
    // 检查电影是否已存在
    if req.IMDBID != "" {
        existingMovie, err := s.movieRepo.GetByIMDBID(ctx, req.IMDBID)
        if err != nil && !errors.IsNotFound(err) {
            return nil, errors.Wrap(err, errors.CodeDatabaseError, "检查IMDB ID失败")
        }
        if existingMovie != nil {
            return nil, errors.MovieExists().WithDetails("IMDB ID已存在")
        }
    }
    
    // 创建电影实体
    movie := &models.Movie{
        Title:         req.Title,
        OriginalTitle: req.OriginalTitle,
        Director:      req.Director,
        ReleaseYear:   req.ReleaseYear,
        Duration:      req.Duration,
        Country:       req.Country,
        Language:      req.Language,
        Plot:          req.Plot,
        Tagline:       req.Tagline,
        PosterURL:     req.PosterURL,
        BackdropURL:   req.BackdropURL,
        TrailerURL:    req.TrailerURL,
        IMDBID:        req.IMDBID,
        TMDBID:        req.TMDBID,
        Status:        models.MovieStatusPublished,
    }
    
    // 设置关键词
    if len(req.Keywords) > 0 {
        movie.SetKeywords(req.Keywords)
    }
    
    // 创建电影
    if err := s.movieRepo.Create(ctx, movie); err != nil {
        return nil, errors.Wrap(err, errors.CodeDatabaseError, "创建电影失败")
    }
    
    // 关联分类
    if len(req.CategoryIDs) > 0 {
        if err := s.associateCategories(ctx, movie.ID, req.CategoryIDs); err != nil {
            s.logger.WithContext(ctx).Error("Failed to associate categories",
                logger.Int64("movie_id", movie.ID),
                logger.Error(err),
            )
        }
    }
    
    // 关联演员
    if len(req.ActorIDs) > 0 {
        if err := s.associateActors(ctx, movie.ID, req.ActorIDs); err != nil {
            s.logger.WithContext(ctx).Error("Failed to associate actors",
                logger.Int64("movie_id", movie.ID),
                logger.Error(err),
            )
        }
    }
    
    // 清除相关缓存
    s.clearMovieCache(ctx, movie.ID)
    
    s.logger.WithContext(ctx).Info("Movie created successfully",
        logger.Int64("movie_id", movie.ID),
        logger.String("title", movie.Title),
    )
    
    return movie, nil
}

// GetMovie 获取电影信息
func (s *ServiceImpl) GetMovie(ctx context.Context, movieID int64) (*models.Movie, error) {
    // 先从缓存获取
    cacheKey := fmt.Sprintf("movie:%d", movieID)
    if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
        if movie, ok := cached.(*models.Movie); ok {
            return movie, nil
        }
    }
    
    // 从数据库获取
    movie, err := s.movieRepo.GetByID(ctx, movieID)
    if err != nil {
        if errors.IsNotFound(err) {
            return nil, errors.MovieNotFound()
        }
        return nil, errors.Wrap(err, errors.CodeDatabaseError, "查询电影失败")
    }
    
    // 检查电影状态
    if !movie.CanView() {
        return nil, errors.MovieNotFound()
    }
    
    // 缓存电影信息
    s.cache.Set(ctx, cacheKey, movie, s.config.CacheExpiry)
    
    return movie, nil
}

// GetMovieWithDetails 获取电影详情
func (s *ServiceImpl) GetMovieWithDetails(ctx context.Context, movieID int64) (*MovieDetails, error) {
    // 获取基础电影信息
    movie, err := s.GetMovie(ctx, movieID)
    if err != nil {
        return nil, err
    }
    
    // 并发获取关联数据
    type result struct {
        categories []*models.Category
        actors     []*models.Actor
        similar    []*models.Movie
        err        error
    }
    
    resultChan := make(chan result, 1)
    
    go func() {
        var res result
        
        // 获取分类
        res.categories, res.err = s.GetMovieCategories(ctx, movieID)
        if res.err != nil {
            resultChan <- res
            return
        }
        
        // 获取演员
        res.actors, res.err = s.GetMovieActors(ctx, movieID)
        if res.err != nil {
            resultChan <- res
            return
        }
        
        // 获取相似电影
        res.similar, res.err = s.GetSimilarMovies(ctx, movieID, s.config.SimilarMoviesLimit)
        if res.err != nil {
            s.logger.WithContext(ctx).Error("Failed to get similar movies",
                logger.Int64("movie_id", movieID),
                logger.Error(res.err),
            )
            res.err = nil // 不阻断主流程
        }
        
        resultChan <- res
    }()
    
    // 增加浏览次数
    go func() {
        if err := s.IncrementViewCount(ctx, movieID); err != nil {
            s.logger.WithContext(ctx).Error("Failed to increment view count",
                logger.Int64("movie_id", movieID),
                logger.Error(err),
            )
        }
    }()
    
    // 等待关联数据
    res := <-resultChan
    if res.err != nil {
        return nil, res.err
    }
    
    details := &MovieDetails{
        Movie:         movie,
        Categories:    res.categories,
        Actors:        res.actors,
        SimilarMovies: res.similar,
    }
    
    return details, nil
}

// GetTopRatedMovies 获取高评分电影
func (s *ServiceImpl) GetTopRatedMovies(ctx context.Context, limit int) ([]*models.Movie, error) {
    cacheKey := fmt.Sprintf("movies:top_rated:%d", limit)
    
    // 先从缓存获取
    if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
        if movies, ok := cached.([]*models.Movie); ok {
            return movies, nil
        }
    }
    
    // 从数据库获取
    movies, err := s.movieRepo.GetTopRated(ctx, limit)
    if err != nil {
        return nil, errors.Wrap(err, errors.CodeDatabaseError, "获取高评分电影失败")
    }
    
    // 缓存结果
    s.cache.Set(ctx, cacheKey, movies, s.config.CacheExpiry)
    
    return movies, nil
}

// IncrementViewCount 增加浏览次数
func (s *ServiceImpl) IncrementViewCount(ctx context.Context, movieID int64) error {
    // 使用缓存计数器，定期同步到数据库
    cacheKey := fmt.Sprintf("movie:view_count:%d", movieID)
    
    // 增加缓存计数
    count, err := s.cache.Increment(ctx, cacheKey, 1)
    if err != nil {
        // 缓存失败，直接更新数据库
        return s.movieRepo.UpdateViewCount(ctx, movieID)
    }
    
    // 每100次浏览同步一次数据库
    if count%100 == 0 {
        go func() {
            if err := s.movieRepo.UpdateViewCount(ctx, movieID); err != nil {
                s.logger.Error("Failed to sync view count to database",
                    logger.Int64("movie_id", movieID),
                    logger.Int64("count", count),
                    logger.Error(err),
                )
            }
        }()
    }
    
    return nil
}

// associateCategories 关联分类
func (s *ServiceImpl) associateCategories(ctx context.Context, movieID int64, categoryIDs []int64) error {
    // 验证分类是否存在
    for _, categoryID := range categoryIDs {
        exists, err := s.categoryRepo.Exists(ctx, categoryID)
        if err != nil {
            return err
        }
        if !exists {
            return errors.NotFound(fmt.Sprintf("分类 %d 不存在", categoryID))
        }
    }
    
    // 关联分类（这里需要实现多对多关联的Repository方法）
    return s.movieRepo.AssociateCategories(ctx, movieID, categoryIDs)
}

// associateActors 关联演员
func (s *ServiceImpl) associateActors(ctx context.Context, movieID int64, actorIDs []int64) error {
    // 验证演员是否存在
    for _, actorID := range actorIDs {
        exists, err := s.actorRepo.Exists(ctx, actorID)
        if err != nil {
            return err
        }
        if !exists {
            return errors.NotFound(fmt.Sprintf("演员 %d 不存在", actorID))
        }
    }
    
    // 关联演员
    return s.movieRepo.AssociateActors(ctx, movieID, actorIDs)
}

// clearMovieCache 清除电影相关缓存
func (s *ServiceImpl) clearMovieCache(ctx context.Context, movieID int64) {
    cacheKeys := []string{
        fmt.Sprintf("movie:%d", movieID),
        fmt.Sprintf("movie:details:%d", movieID),
        "movies:top_rated:*",
        "movies:popular:*",
        "movies:latest:*",
    }
    
    for _, key := range cacheKeys {
        if strings.Contains(key, "*") {
            s.cache.DeletePattern(ctx, key)
        } else {
            s.cache.Delete(ctx, key)
        }
    }
}
```

## 5.2.4 总结

电影服务实现为MovieInfo项目提供了完整的电影管理解决方案。通过分层架构、缓存优化和外部集成，我们建立了一个功能丰富、性能优异的电影服务系统。

**关键设计要点**：
1. **数据管理**：完整的电影数据CRUD操作和关联管理
2. **搜索功能**：高效的全文搜索和多条件筛选
3. **推荐算法**：智能的电影推荐和相似度计算
4. **缓存策略**：多层缓存提升访问性能
5. **外部集成**：与TMDB、IMDB等外部数据源集成

**服务优势**：
- **功能完整**：覆盖电影管理的各个方面
- **性能优化**：缓存和查询优化提升响应速度
- **智能推荐**：基于多种算法的推荐系统
- **数据丰富**：集成外部数据源丰富电影信息

**下一步**：基于电影服务基础，我们将实现评论服务，处理用户评论、评分、点赞等互动功能。
