# 第41步：Web服务器搭建

## 📋 概述

Web服务器是MovieInfo项目的前端入口，负责处理用户请求、渲染页面和提供API服务。基于Gin框架构建高性能的Web服务器，集成模板引擎、静态资源管理和后端服务调用。

## 🎯 设计目标

### 1. **高性能**
- 快速的请求响应
- 高效的资源利用
- 优化的并发处理
- 智能的缓存策略

### 2. **可扩展性**
- 模块化的架构设计
- 灵活的路由配置
- 可插拔的中间件
- 微服务集成支持

### 3. **用户体验**
- 快速的页面加载
- 响应式设计支持
- SEO友好的结构
- 优雅的错误处理

## 🏗️ 服务器架构

### 1. **Web服务器结构**

```go
// Web服务器配置
type WebServerConfig struct {
    Port            string        `yaml:"port"`
    Mode            string        `yaml:"mode"` // debug, release, test
    ReadTimeout     time.Duration `yaml:"read_timeout"`
    WriteTimeout    time.Duration `yaml:"write_timeout"`
    MaxHeaderBytes  int           `yaml:"max_header_bytes"`

    // 静态资源配置
    StaticPath      string        `yaml:"static_path"`
    TemplatePath    string        `yaml:"template_path"`

    // 安全配置
    TrustedProxies  []string      `yaml:"trusted_proxies"`
    EnableHTTPS     bool          `yaml:"enable_https"`
    CertFile        string        `yaml:"cert_file"`
    KeyFile         string        `yaml:"key_file"`

    // 性能配置
    EnableGzip      bool          `yaml:"enable_gzip"`
    EnableCache     bool          `yaml:"enable_cache"`
    CacheMaxAge     int           `yaml:"cache_max_age"`
}

// Web服务器主结构
type WebServer struct {
    config          *WebServerConfig
    router          *gin.Engine
    server          *http.Server

    // 服务客户端
    userClient      user.UserServiceClient
    movieClient     movie.MovieServiceClient
    commentClient   comment.CommentServiceClient

    // 模板引擎
    templateEngine  *TemplateEngine

    // 中间件
    authMiddleware  *AuthMiddleware
    corsMiddleware  *CORSMiddleware

    logger          *logrus.Logger
    metrics         *WebServerMetrics
}

func NewWebServer(config *WebServerConfig) *WebServer {
    // 设置Gin模式
    gin.SetMode(config.Mode)

    ws := &WebServer{
        config: config,
        router: gin.New(),
        logger: logrus.New(),
        metrics: NewWebServerMetrics(),
    }

    // 初始化组件
    ws.initializeClients()
    ws.initializeTemplateEngine()
    ws.initializeMiddlewares()
    ws.setupRoutes()

    return ws
}
```

### 2. **Gin框架集成**

```go
// 初始化Web服务器
func (ws *WebServer) Initialize() error {
    // 配置Gin引擎
    ws.setupGinEngine()

    // 设置HTTP服务器
    ws.server = &http.Server{
        Addr:           ":" + ws.config.Port,
        Handler:        ws.router,
        ReadTimeout:    ws.config.ReadTimeout,
        WriteTimeout:   ws.config.WriteTimeout,
        MaxHeaderBytes: ws.config.MaxHeaderBytes,
    }

    return nil
}

// 配置Gin引擎
func (ws *WebServer) setupGinEngine() {
    // 基础中间件
    ws.router.Use(gin.Logger())
    ws.router.Use(gin.Recovery())

    // 自定义中间件
    ws.router.Use(ws.corsMiddleware.Handler())
    ws.router.Use(ws.requestIDMiddleware())
    ws.router.Use(ws.metricsMiddleware())

    // Gzip压缩
    if ws.config.EnableGzip {
        ws.router.Use(gzip.Gzip(gzip.DefaultCompression))
    }

    // 信任代理
    if len(ws.config.TrustedProxies) > 0 {
        ws.router.SetTrustedProxies(ws.config.TrustedProxies)
    }

    // 静态资源
    ws.router.Static("/static", ws.config.StaticPath)
    ws.router.StaticFile("/favicon.ico", ws.config.StaticPath+"/favicon.ico")

    // 模板引擎
    ws.router.HTMLRender = ws.templateEngine.Render()
}

// 启动服务器
func (ws *WebServer) Start() error {
    ws.logger.Infof("Starting web server on port %s", ws.config.Port)

    if ws.config.EnableHTTPS {
        return ws.server.ListenAndServeTLS(ws.config.CertFile, ws.config.KeyFile)
    }

    return ws.server.ListenAndServe()
}

// 优雅关闭
func (ws *WebServer) Shutdown(ctx context.Context) error {
    ws.logger.Info("Shutting down web server...")

    return ws.server.Shutdown(ctx)
}
```

### 3. **中间件系统**

```go
// 请求ID中间件
func (ws *WebServer) requestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }

        c.Header("X-Request-ID", requestID)
        c.Set("request_id", requestID)

        c.Next()
    }
}

// 指标收集中间件
func (ws *WebServer) metricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        method := c.Request.Method

        c.Next()

        duration := time.Since(start)
        status := c.Writer.Status()

        ws.metrics.RecordRequest(method, path, status, duration)
    }
}

// 认证中间件
type AuthMiddleware struct {
    jwtManager *JWTManager
    userClient user.UserServiceClient
    logger     *logrus.Logger
}

func NewAuthMiddleware(jwtManager *JWTManager, userClient user.UserServiceClient) *AuthMiddleware {
    return &AuthMiddleware{
        jwtManager: jwtManager,
        userClient: userClient,
        logger:     logrus.New(),
    }
}

func (am *AuthMiddleware) RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := am.extractToken(c)
        if token == "" {
            c.Redirect(http.StatusFound, "/login")
            c.Abort()
            return
        }

        claims, err := am.jwtManager.Verify(token)
        if err != nil {
            am.logger.Warnf("Invalid token: %v", err)
            c.Redirect(http.StatusFound, "/login")
            c.Abort()
            return
        }

        // 验证用户状态
        userResp, err := am.userClient.GetUser(c.Request.Context(), &user.GetUserRequest{
            UserId: claims.UserID,
        })
        if err != nil || userResp.User.Status != "active" {
            c.Redirect(http.StatusFound, "/login")
            c.Abort()
            return
        }

        // 设置用户信息到上下文
        c.Set("user", userResp.User)
        c.Set("user_id", claims.UserID)

        c.Next()
    }
}

func (am *AuthMiddleware) extractToken(c *gin.Context) string {
    // 从Cookie获取
    if token, err := c.Cookie("access_token"); err == nil && token != "" {
        return token
    }

    // 从Authorization头获取
    authHeader := c.GetHeader("Authorization")
    if strings.HasPrefix(authHeader, "Bearer ") {
        return strings.TrimPrefix(authHeader, "Bearer ")
    }

    return ""
}

// CORS中间件
type CORSMiddleware struct {
    config *CORSConfig
}

type CORSConfig struct {
    AllowOrigins     []string
    AllowMethods     []string
    AllowHeaders     []string
    ExposeHeaders    []string
    AllowCredentials bool
    MaxAge           time.Duration
}

func NewCORSMiddleware(config *CORSConfig) *CORSMiddleware {
    return &CORSMiddleware{
        config: config,
    }
}

func (cm *CORSMiddleware) Handler() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     cm.config.AllowOrigins,
        AllowMethods:     cm.config.AllowMethods,
        AllowHeaders:     cm.config.AllowHeaders,
        ExposeHeaders:    cm.config.ExposeHeaders,
        AllowCredentials: cm.config.AllowCredentials,
        MaxAge:           cm.config.MaxAge,
    })
}
```

### 4. **gRPC客户端集成**

```go
// 初始化gRPC客户端
func (ws *WebServer) initializeClients() error {
    // 用户服务客户端
    userConn, err := grpc.Dial("user-service:50051", grpc.WithInsecure())
    if err != nil {
        return fmt.Errorf("failed to connect to user service: %v", err)
    }
    ws.userClient = user.NewUserServiceClient(userConn)

    // 电影服务客户端
    movieConn, err := grpc.Dial("movie-service:50052", grpc.WithInsecure())
    if err != nil {
        return fmt.Errorf("failed to connect to movie service: %v", err)
    }
    ws.movieClient = movie.NewMovieServiceClient(movieConn)

    // 评论服务客户端
    commentConn, err := grpc.Dial("comment-service:50053", grpc.WithInsecure())
    if err != nil {
        return fmt.Errorf("failed to connect to comment service: %v", err)
    }
    ws.commentClient = comment.NewCommentServiceClient(commentConn)

    return nil
}

// 服务调用封装
type ServiceClient struct {
    userClient    user.UserServiceClient
    movieClient   movie.MovieServiceClient
    commentClient comment.CommentServiceClient
    logger        *logrus.Logger
}

func NewServiceClient(userClient user.UserServiceClient, movieClient movie.MovieServiceClient, commentClient comment.CommentServiceClient) *ServiceClient {
    return &ServiceClient{
        userClient:    userClient,
        movieClient:   movieClient,
        commentClient: commentClient,
        logger:        logrus.New(),
    }
}

// 获取电影列表
func (sc *ServiceClient) GetMovieList(ctx context.Context, req *movie.GetMovieListRequest) (*movie.GetMovieListResponse, error) {
    resp, err := sc.movieClient.GetMovieList(ctx, req)
    if err != nil {
        sc.logger.Errorf("Failed to get movie list: %v", err)
        return nil, err
    }
    return resp, nil
}

// 获取电影详情
func (sc *ServiceClient) GetMovieDetail(ctx context.Context, movieID string) (*movie.GetMovieDetailResponse, error) {
    resp, err := sc.movieClient.GetMovieDetail(ctx, &movie.GetMovieDetailRequest{
        MovieId: movieID,
    })
    if err != nil {
        sc.logger.Errorf("Failed to get movie detail: %v", err)
        return nil, err
    }
    return resp, nil
}

// 用户认证
func (sc *ServiceClient) AuthenticateUser(ctx context.Context, email, password string) (*user.LoginResponse, error) {
    resp, err := sc.userClient.Login(ctx, &user.LoginRequest{
        Email:    email,
        Password: password,
    })
    if err != nil {
        sc.logger.Errorf("Failed to authenticate user: %v", err)
        return nil, err
    }
    return resp, nil
}
```

### 5. **错误处理机制**

```go
// 错误处理器
type ErrorHandler struct {
    logger *logrus.Logger
}

func NewErrorHandler() *ErrorHandler {
    return &ErrorHandler{
        logger: logrus.New(),
    }
}

// 处理HTTP错误
func (eh *ErrorHandler) HandleError(c *gin.Context, err error, statusCode int) {
    requestID := c.GetString("request_id")

    eh.logger.WithFields(logrus.Fields{
        "request_id": requestID,
        "path":       c.Request.URL.Path,
        "method":     c.Request.Method,
        "error":      err.Error(),
    }).Error("HTTP request error")

    // 根据请求类型返回不同格式的错误
    if c.GetHeader("Accept") == "application/json" || strings.HasPrefix(c.Request.URL.Path, "/api/") {
        c.JSON(statusCode, gin.H{
            "success":    false,
            "message":    eh.getErrorMessage(statusCode),
            "request_id": requestID,
        })
    } else {
        // 渲染错误页面
        c.HTML(statusCode, "error.html", gin.H{
            "title":      eh.getErrorTitle(statusCode),
            "message":    eh.getErrorMessage(statusCode),
            "status":     statusCode,
            "request_id": requestID,
        })
    }
}

func (eh *ErrorHandler) getErrorTitle(statusCode int) string {
    switch statusCode {
    case 404:
        return "页面未找到"
    case 500:
        return "服务器内部错误"
    case 403:
        return "访问被拒绝"
    case 401:
        return "未授权访问"
    default:
        return "发生错误"
    }
}

func (eh *ErrorHandler) getErrorMessage(statusCode int) string {
    switch statusCode {
    case 404:
        return "您访问的页面不存在"
    case 500:
        return "服务器遇到了一个错误，请稍后重试"
    case 403:
        return "您没有权限访问此资源"
    case 401:
        return "请先登录后再访问"
    default:
        return "发生了未知错误"
    }
}

// 全局错误处理中间件
func (ws *WebServer) errorHandlerMiddleware() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        if err, ok := recovered.(string); ok {
            ws.errorHandler.HandleError(c, errors.New(err), 500)
        } else if err, ok := recovered.(error); ok {
            ws.errorHandler.HandleError(c, err, 500)
        } else {
            ws.errorHandler.HandleError(c, errors.New("unknown error"), 500)
        }
        c.Abort()
    })
}
```

## 📊 性能监控

### 1. **Web服务器指标**

```go
type WebServerMetrics struct {
    requestCount     *prometheus.CounterVec
    requestDuration  *prometheus.HistogramVec
    responseSize     *prometheus.HistogramVec
    activeConnections prometheus.Gauge
    errorCount       *prometheus.CounterVec
}

func NewWebServerMetrics() *WebServerMetrics {
    return &WebServerMetrics{
        requestCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "http_requests_total",
                Help: "Total number of HTTP requests",
            },
            []string{"method", "path", "status"},
        ),
        requestDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "http_request_duration_seconds",
                Help: "Duration of HTTP requests",
            },
            []string{"method", "path"},
        ),
        responseSize: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "http_response_size_bytes",
                Help: "Size of HTTP responses",
            },
            []string{"method", "path"},
        ),
        activeConnections: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "http_active_connections",
                Help: "Number of active HTTP connections",
            },
        ),
        errorCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "http_errors_total",
                Help: "Total number of HTTP errors",
            },
            []string{"method", "path", "error_type"},
        ),
    }
}

func (wsm *WebServerMetrics) RecordRequest(method, path string, status int, duration time.Duration) {
    statusStr := strconv.Itoa(status)
    wsm.requestCount.WithLabelValues(method, path, statusStr).Inc()
    wsm.requestDuration.WithLabelValues(method, path).Observe(duration.Seconds())

    if status >= 400 {
        errorType := "client_error"
        if status >= 500 {
            errorType = "server_error"
        }
        wsm.errorCount.WithLabelValues(method, path, errorType).Inc()
    }
}
```

## 🔧 配置文件

### 1. **Web服务器配置**

```yaml
# config/web.yaml
web_server:
  port: "8080"
  mode: "release"  # debug, release, test
  read_timeout: 30s
  write_timeout: 30s
  max_header_bytes: 1048576  # 1MB

  # 静态资源配置
  static_path: "./web/static"
  template_path: "./web/templates"

  # 安全配置
  trusted_proxies:
    - "127.0.0.1"
    - "10.0.0.0/8"
  enable_https: false
  cert_file: ""
  key_file: ""

  # 性能配置
  enable_gzip: true
  enable_cache: true
  cache_max_age: 3600

# CORS配置
cors:
  allow_origins:
    - "http://localhost:3000"
    - "https://movieinfo.com"
  allow_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allow_headers:
    - "Origin"
    - "Content-Type"
    - "Authorization"
    - "X-Requested-With"
  expose_headers:
    - "X-Request-ID"
  allow_credentials: true
  max_age: 86400

# gRPC服务配置
grpc_services:
  user_service:
    address: "user-service:50051"
    timeout: 10s
    retry_count: 3
  movie_service:
    address: "movie-service:50052"
    timeout: 10s
    retry_count: 3
  comment_service:
    address: "comment-service:50053"
    timeout: 10s
    retry_count: 3
```

## 📝 总结

Web服务器搭建为MovieInfo项目提供了完整的前端服务基础：

**核心功能**：
1. **高性能Web服务器**：基于Gin框架的高性能HTTP服务器
2. **完整的中间件系统**：认证、CORS、指标收集等中间件
3. **gRPC服务集成**：与后端微服务的无缝集成
4. **优雅的错误处理**：统一的错误处理和用户友好的错误页面

**技术特性**：
- 模块化的架构设计
- 灵活的配置管理
- 完善的监控指标
- 优雅的服务关闭

**性能优化**：
- Gzip压缩支持
- 静态资源缓存
- 连接池管理
- 请求指标监控

下一步，我们将设计详细的路由系统，为不同的页面和API提供清晰的访问路径。

---

**文档状态**: ✅ 已完成
**最后更新**: 2025-07-22
**下一步**: [第42步：路由设计](42-routing.md)
```
