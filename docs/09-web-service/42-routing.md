# 第42步：路由设计

## 📋 概述

路由设计是Web服务器的核心组件，负责将用户请求映射到相应的处理器。一个良好的路由设计需要考虑RESTful原则、用户体验、SEO优化和系统可维护性。

## 🎯 设计目标

### 1. **用户友好**
- 直观的URL结构
- 语义化的路径命名
- SEO友好的URL
- 简洁的参数传递

### 2. **系统可维护**
- 清晰的路由分组
- 统一的命名规范
- 灵活的中间件配置
- 易于扩展的结构

### 3. **性能优化**
- 高效的路由匹配
- 合理的缓存策略
- 最小化的重定向
- 优化的静态资源路径

## 🗺️ 路由架构设计

### 1. **路由结构规划**

```go
// 路由配置
type RouteConfig struct {
    Prefix      string
    Middlewares []gin.HandlerFunc
    Routes      []Route
}

type Route struct {
    Method      string
    Path        string
    Handler     gin.HandlerFunc
    Middlewares []gin.HandlerFunc
    Name        string
    Description string
}

// 路由管理器
type RouteManager struct {
    router      *gin.Engine
    authMW      *AuthMiddleware
    rateLimitMW *RateLimitMiddleware
    cacheMW     *CacheMiddleware
    logger      *logrus.Logger
}

func NewRouteManager(router *gin.Engine, authMW *AuthMiddleware) *RouteManager {
    return &RouteManager{
        router:      router,
        authMW:      authMW,
        rateLimitMW: NewRateLimitMiddleware(),
        cacheMW:     NewCacheMiddleware(),
        logger:      logrus.New(),
    }
}

// 设置所有路由
func (rm *RouteManager) SetupRoutes(controllers *Controllers) {
    // 页面路由
    rm.setupPageRoutes(controllers.PageController)
    
    // API路由
    rm.setupAPIRoutes(controllers)
    
    // 管理后台路由
    rm.setupAdminRoutes(controllers.AdminController)
    
    // 健康检查路由
    rm.setupHealthRoutes()
}
```

### 2. **页面路由设计**

```go
// 页面路由配置
func (rm *RouteManager) setupPageRoutes(pc *PageController) {
    // 首页路由
    rm.router.GET("/", rm.cacheMW.Cache(5*time.Minute), pc.Home)
    rm.router.GET("/home", func(c *gin.Context) {
        c.Redirect(http.StatusMovedPermanently, "/")
    })
    
    // 电影相关页面
    movieGroup := rm.router.Group("/movies")
    {
        movieGroup.GET("", rm.cacheMW.Cache(10*time.Minute), pc.MovieList)
        movieGroup.GET("/popular", rm.cacheMW.Cache(15*time.Minute), pc.PopularMovies)
        movieGroup.GET("/latest", rm.cacheMW.Cache(10*time.Minute), pc.LatestMovies)
        movieGroup.GET("/top-rated", rm.cacheMW.Cache(30*time.Minute), pc.TopRatedMovies)
        movieGroup.GET("/genres/:genre", rm.cacheMW.Cache(20*time.Minute), pc.MoviesByGenre)
        movieGroup.GET("/search", pc.SearchMovies)
        movieGroup.GET("/:id", rm.cacheMW.Cache(30*time.Minute), pc.MovieDetail)
        movieGroup.GET("/:id/reviews", rm.cacheMW.Cache(5*time.Minute), pc.MovieReviews)
    }
    
    // 用户相关页面
    userGroup := rm.router.Group("/user")
    {
        userGroup.GET("/login", pc.LoginPage)
        userGroup.GET("/register", pc.RegisterPage)
        userGroup.GET("/forgot-password", pc.ForgotPasswordPage)
        userGroup.GET("/reset-password", pc.ResetPasswordPage)
        userGroup.GET("/logout", rm.authMW.RequireAuth(), pc.Logout)
        
        // 需要认证的用户页面
        authUserGroup := userGroup.Group("", rm.authMW.RequireAuth())
        {
            authUserGroup.GET("/profile", pc.UserProfile)
            authUserGroup.GET("/settings", pc.UserSettings)
            authUserGroup.GET("/favorites", pc.UserFavorites)
            authUserGroup.GET("/watchlist", pc.UserWatchlist)
            authUserGroup.GET("/reviews", pc.UserReviews)
            authUserGroup.GET("/ratings", pc.UserRatings)
        }
    }
    
    // 分类页面
    categoryGroup := rm.router.Group("/categories")
    {
        categoryGroup.GET("", rm.cacheMW.Cache(1*time.Hour), pc.Categories)
        categoryGroup.GET("/:category", rm.cacheMW.Cache(30*time.Minute), pc.CategoryMovies)
    }
    
    // 静态页面
    staticGroup := rm.router.Group("/pages")
    {
        staticGroup.GET("/about", rm.cacheMW.Cache(1*time.Hour), pc.About)
        staticGroup.GET("/contact", rm.cacheMW.Cache(1*time.Hour), pc.Contact)
        staticGroup.GET("/privacy", rm.cacheMW.Cache(1*time.Hour), pc.Privacy)
        staticGroup.GET("/terms", rm.cacheMW.Cache(1*time.Hour), pc.Terms)
        staticGroup.GET("/help", rm.cacheMW.Cache(1*time.Hour), pc.Help)
    }
}
```

### 3. **API路由设计**

```go
// API路由配置
func (rm *RouteManager) setupAPIRoutes(controllers *Controllers) {
    apiGroup := rm.router.Group("/api/v1")
    apiGroup.Use(rm.rateLimitMW.Limit(100, time.Minute)) // API限流
    
    // 认证API
    authGroup := apiGroup.Group("/auth")
    {
        authGroup.POST("/login", controllers.AuthController.Login)
        authGroup.POST("/register", controllers.AuthController.Register)
        authGroup.POST("/logout", rm.authMW.RequireAuth(), controllers.AuthController.Logout)
        authGroup.POST("/refresh", controllers.AuthController.RefreshToken)
        authGroup.POST("/forgot-password", controllers.AuthController.ForgotPassword)
        authGroup.POST("/reset-password", controllers.AuthController.ResetPassword)
        authGroup.GET("/verify-email", controllers.AuthController.VerifyEmail)
    }
    
    // 电影API
    movieGroup := apiGroup.Group("/movies")
    {
        movieGroup.GET("", rm.cacheMW.Cache(10*time.Minute), controllers.MovieController.GetMovieList)
        movieGroup.GET("/popular", rm.cacheMW.Cache(15*time.Minute), controllers.MovieController.GetPopularMovies)
        movieGroup.GET("/search", controllers.MovieController.SearchMovies)
        movieGroup.GET("/:id", rm.cacheMW.Cache(30*time.Minute), controllers.MovieController.GetMovieDetail)
        movieGroup.GET("/:id/similar", rm.cacheMW.Cache(1*time.Hour), controllers.MovieController.GetSimilarMovies)
        movieGroup.GET("/:id/recommendations", rm.authMW.OptionalAuth(), controllers.MovieController.GetRecommendations)
        
        // 需要认证的电影API
        authMovieGroup := movieGroup.Group("", rm.authMW.RequireAuth())
        {
            authMovieGroup.POST("/:id/favorite", controllers.MovieController.AddToFavorites)
            authMovieGroup.DELETE("/:id/favorite", controllers.MovieController.RemoveFromFavorites)
            authMovieGroup.POST("/:id/watchlist", controllers.MovieController.AddToWatchlist)
            authMovieGroup.DELETE("/:id/watchlist", controllers.MovieController.RemoveFromWatchlist)
        }
    }
    
    // 评论API
    commentGroup := apiGroup.Group("/comments")
    {
        commentGroup.GET("/movie/:movie_id", controllers.CommentController.GetComments)
        
        // 需要认证的评论API
        authCommentGroup := commentGroup.Group("", rm.authMW.RequireAuth())
        {
            authCommentGroup.POST("", controllers.CommentController.CreateComment)
            authCommentGroup.PUT("/:id", controllers.CommentController.UpdateComment)
            authCommentGroup.DELETE("/:id", controllers.CommentController.DeleteComment)
            authCommentGroup.POST("/:id/like", controllers.CommentController.LikeComment)
            authCommentGroup.DELETE("/:id/like", controllers.CommentController.UnlikeComment)
            authCommentGroup.POST("/:id/report", controllers.CommentController.ReportComment)
        }
    }
    
    // 评分API
    ratingGroup := apiGroup.Group("/ratings")
    {
        ratingGroup.GET("/movie/:movie_id", rm.cacheMW.Cache(10*time.Minute), controllers.RatingController.GetMovieRatings)
        
        // 需要认证的评分API
        authRatingGroup := ratingGroup.Group("", rm.authMW.RequireAuth())
        {
            authRatingGroup.POST("", controllers.RatingController.SubmitRating)
            authRatingGroup.GET("/user", controllers.RatingController.GetUserRatings)
        }
    }
    
    // 用户API
    userGroup := apiGroup.Group("/users", rm.authMW.RequireAuth())
    {
        userGroup.GET("/profile", controllers.UserController.GetProfile)
        userGroup.PUT("/profile", controllers.UserController.UpdateProfile)
        userGroup.GET("/favorites", controllers.UserController.GetFavorites)
        userGroup.GET("/watchlist", controllers.UserController.GetWatchlist)
        userGroup.GET("/reviews", controllers.UserController.GetUserReviews)
        userGroup.GET("/ratings", controllers.UserController.GetUserRatings)
        userGroup.POST("/change-password", controllers.UserController.ChangePassword)
        userGroup.DELETE("/account", controllers.UserController.DeleteAccount)
    }
    
    // 搜索API
    searchGroup := apiGroup.Group("/search")
    {
        searchGroup.GET("/movies", controllers.SearchController.SearchMovies)
        searchGroup.GET("/suggestions", controllers.SearchController.GetSuggestions)
        searchGroup.GET("/hot", rm.cacheMW.Cache(1*time.Hour), controllers.SearchController.GetHotSearches)
    }
}
```

### 4. **管理后台路由**

```go
// 管理后台路由
func (rm *RouteManager) setupAdminRoutes(ac *AdminController) {
    adminGroup := rm.router.Group("/admin")
    adminGroup.Use(rm.authMW.RequireAuth())
    adminGroup.Use(rm.authMW.RequireRole("admin", "moderator"))
    
    // 管理后台首页
    adminGroup.GET("", ac.Dashboard)
    adminGroup.GET("/dashboard", ac.Dashboard)
    
    // 用户管理
    userMgmtGroup := adminGroup.Group("/users")
    {
        userMgmtGroup.GET("", ac.UserList)
        userMgmtGroup.GET("/:id", ac.UserDetail)
        userMgmtGroup.PUT("/:id/status", ac.UpdateUserStatus)
        userMgmtGroup.DELETE("/:id", ac.DeleteUser)
        userMgmtGroup.GET("/:id/activities", ac.UserActivities)
    }
    
    // 电影管理
    movieMgmtGroup := adminGroup.Group("/movies")
    {
        movieMgmtGroup.GET("", ac.MovieList)
        movieMgmtGroup.GET("/:id", ac.MovieDetail)
        movieMgmtGroup.POST("", ac.CreateMovie)
        movieMgmtGroup.PUT("/:id", ac.UpdateMovie)
        movieMgmtGroup.DELETE("/:id", ac.DeleteMovie)
        movieMgmtGroup.POST("/:id/approve", ac.ApproveMovie)
    }
    
    // 评论审核
    moderationGroup := adminGroup.Group("/moderation")
    {
        moderationGroup.GET("/comments", ac.PendingComments)
        moderationGroup.PUT("/comments/:id/approve", ac.ApproveComment)
        moderationGroup.PUT("/comments/:id/reject", ac.RejectComment)
        moderationGroup.GET("/reports", ac.ReportList)
        moderationGroup.PUT("/reports/:id/resolve", ac.ResolveReport)
    }
    
    // 系统统计
    analyticsGroup := adminGroup.Group("/analytics")
    {
        analyticsGroup.GET("/overview", ac.AnalyticsOverview)
        analyticsGroup.GET("/users", ac.UserAnalytics)
        analyticsGroup.GET("/content", ac.ContentAnalytics)
        analyticsGroup.GET("/performance", ac.PerformanceAnalytics)
    }
    
    // 系统设置
    settingsGroup := adminGroup.Group("/settings")
    {
        settingsGroup.GET("", ac.SystemSettings)
        settingsGroup.PUT("", ac.UpdateSettings)
        settingsGroup.GET("/cache", ac.CacheStatus)
        settingsGroup.POST("/cache/clear", ac.ClearCache)
    }
}
```

### 5. **路由中间件配置**

```go
// 可选认证中间件
func (am *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := am.extractToken(c)
        if token != "" {
            if claims, err := am.jwtManager.Verify(token); err == nil {
                // 验证用户状态
                if userResp, err := am.userClient.GetUser(c.Request.Context(), &user.GetUserRequest{
                    UserId: claims.UserID,
                }); err == nil && userResp.User.Status == "active" {
                    c.Set("user", userResp.User)
                    c.Set("user_id", claims.UserID)
                }
            }
        }
        c.Next()
    }
}

// 角色验证中间件
func (am *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user, exists := c.Get("user")
        if !exists {
            c.JSON(403, gin.H{
                "success": false,
                "message": "需要登录",
            })
            c.Abort()
            return
        }
        
        userInfo := user.(*user.User)
        for _, role := range roles {
            if userInfo.Role == role {
                c.Next()
                return
            }
        }
        
        c.JSON(403, gin.H{
            "success": false,
            "message": "权限不足",
        })
        c.Abort()
    }
}

// 限流中间件
type RateLimitMiddleware struct {
    limiter *rate.Limiter
    store   map[string]*rate.Limiter
    mutex   sync.RWMutex
}

func NewRateLimitMiddleware() *RateLimitMiddleware {
    return &RateLimitMiddleware{
        store: make(map[string]*rate.Limiter),
    }
}

func (rlm *RateLimitMiddleware) Limit(requests int, duration time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        clientIP := c.ClientIP()
        
        rlm.mutex.RLock()
        limiter, exists := rlm.store[clientIP]
        rlm.mutex.RUnlock()
        
        if !exists {
            rlm.mutex.Lock()
            limiter = rate.NewLimiter(rate.Every(duration/time.Duration(requests)), requests)
            rlm.store[clientIP] = limiter
            rlm.mutex.Unlock()
        }
        
        if !limiter.Allow() {
            c.JSON(429, gin.H{
                "success": false,
                "message": "请求过于频繁，请稍后重试",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// 缓存中间件
type CacheMiddleware struct {
    cache *cache.Cache
}

func NewCacheMiddleware() *CacheMiddleware {
    return &CacheMiddleware{
        cache: cache.New(5*time.Minute, 10*time.Minute),
    }
}

func (cm *CacheMiddleware) Cache(duration time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 只缓存GET请求
        if c.Request.Method != "GET" {
            c.Next()
            return
        }
        
        cacheKey := cm.generateCacheKey(c)
        
        // 尝试从缓存获取
        if cached, found := cm.cache.Get(cacheKey); found {
            if response, ok := cached.(CachedResponse); ok {
                // 设置缓存头
                c.Header("X-Cache", "HIT")
                c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", int(duration.Seconds())))
                
                // 返回缓存的响应
                for key, value := range response.Headers {
                    c.Header(key, value)
                }
                c.Data(response.StatusCode, response.ContentType, response.Body)
                c.Abort()
                return
            }
        }
        
        // 缓存未命中，继续处理请求
        c.Header("X-Cache", "MISS")
        
        // 创建响应写入器来捕获响应
        writer := &CacheResponseWriter{
            ResponseWriter: c.Writer,
            body:          bytes.NewBuffer([]byte{}),
            headers:       make(map[string]string),
        }
        c.Writer = writer
        
        c.Next()
        
        // 缓存响应
        if writer.statusCode == 200 {
            response := CachedResponse{
                StatusCode:  writer.statusCode,
                ContentType: writer.Header().Get("Content-Type"),
                Headers:     writer.headers,
                Body:        writer.body.Bytes(),
            }
            cm.cache.Set(cacheKey, response, duration)
        }
    }
}

type CachedResponse struct {
    StatusCode  int
    ContentType string
    Headers     map[string]string
    Body        []byte
}

type CacheResponseWriter struct {
    gin.ResponseWriter
    body       *bytes.Buffer
    headers    map[string]string
    statusCode int
}

func (crw *CacheResponseWriter) Write(data []byte) (int, error) {
    crw.body.Write(data)
    return crw.ResponseWriter.Write(data)
}

func (crw *CacheResponseWriter) WriteHeader(statusCode int) {
    crw.statusCode = statusCode
    crw.ResponseWriter.WriteHeader(statusCode)
}

func (cm *CacheMiddleware) generateCacheKey(c *gin.Context) string {
    return fmt.Sprintf("%s:%s:%s", c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery)
}
```

## 📝 总结

路由设计为MovieInfo项目提供了完整的URL结构和请求处理机制：

**核心特性**：
1. **清晰的路由结构**：页面路由、API路由、管理后台路由分离
2. **灵活的中间件系统**：认证、限流、缓存等中间件支持
3. **RESTful API设计**：符合REST原则的API接口设计
4. **SEO友好的URL**：语义化的URL结构，有利于搜索引擎优化

**技术特性**：
- 模块化的路由管理
- 高效的路由匹配
- 智能的缓存策略
- 完善的权限控制

**用户体验**：
- 直观的URL结构
- 快速的页面响应
- 合理的重定向处理
- 友好的错误页面

下一步，我们将集成模板引擎，为页面渲染提供强大的模板支持。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第43步：模板引擎集成](43-template-engine.md)
