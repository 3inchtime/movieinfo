# 第23步：gRPC 客户端封装

## 📋 概述

gRPC客户端封装是服务间通信的关键组件，它为主页服务提供了调用后端服务的能力。良好的客户端封装不仅简化了服务调用，还提供了连接管理、错误处理、重试机制等高级功能。

## 🎯 设计目标

### 1. **简化调用**
- 提供简洁的API接口
- 隐藏gRPC连接细节
- 自动处理序列化/反序列化

### 2. **连接管理**
- 连接池管理
- 自动重连机制
- 健康检查

### 3. **错误处理**
- 统一错误处理
- 重试机制
- 熔断保护

### 4. **性能优化**
- 连接复用
- 请求超时控制
- 负载均衡

## 🔧 客户端架构设计

### 1. **客户端层次结构**

```
┌─────────────────────────────────────────────────────────────┐
│                    业务层                                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │  用户控制器  │  │  电影控制器  │  │  评论控制器  │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  客户端封装层                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │用户服务客户端│  │电影服务客户端│  │评论服务客户端│          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  gRPC传输层                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │  连接池管理  │  │  负载均衡    │  │  错误处理    │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

### 2. **客户端接口设计**

```go
// 客户端接口定义
type ServiceClient interface {
    Connect() error
    Close() error
    IsHealthy() bool
}

// 用户服务客户端接口
type UserServiceClient interface {
    ServiceClient
    Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
    Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
    GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error)
    UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoResponse, error)
    ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error)
}

// 电影服务客户端接口
type MovieServiceClient interface {
    ServiceClient
    GetMovieList(ctx context.Context, req *pb.GetMovieListRequest) (*pb.GetMovieListResponse, error)
    GetMovieDetail(ctx context.Context, req *pb.GetMovieDetailRequest) (*pb.GetMovieDetailResponse, error)
    SearchMovies(ctx context.Context, req *pb.SearchMoviesRequest) (*pb.SearchMoviesResponse, error)
    GetMoviesByCategory(ctx context.Context, req *pb.GetMoviesByCategoryRequest) (*pb.GetMoviesByCategoryResponse, error)
}

// 评论服务客户端接口
type CommentServiceClient interface {
    ServiceClient
    CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error)
    GetComments(ctx context.Context, req *pb.GetCommentsRequest) (*pb.GetCommentsResponse, error)
    UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error)
    DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error)
    RateMovie(ctx context.Context, req *pb.RateMovieRequest) (*pb.RateMovieResponse, error)
    GetMovieRating(ctx context.Context, req *pb.GetMovieRatingRequest) (*pb.GetMovieRatingResponse, error)
}
```

## 🏗️ 客户端实现

### 1. **基础客户端结构**

```go
// 基础客户端配置
type ClientConfig struct {
    Address         string        `yaml:"address"`
    Timeout         time.Duration `yaml:"timeout"`
    MaxRetries      int           `yaml:"max_retries"`
    RetryInterval   time.Duration `yaml:"retry_interval"`
    PoolSize        int           `yaml:"pool_size"`
    HealthCheckInterval time.Duration `yaml:"health_check_interval"`
}

// 基础客户端实现
type BaseClient struct {
    config     *ClientConfig
    conn       *grpc.ClientConn
    connPool   *ConnectionPool
    logger     *logrus.Logger
    metrics    *ClientMetrics
    mutex      sync.RWMutex
    healthy    bool
}

func NewBaseClient(config *ClientConfig) *BaseClient {
    return &BaseClient{
        config:  config,
        logger:  logrus.New(),
        metrics: NewClientMetrics(),
        healthy: false,
    }
}

func (bc *BaseClient) Connect() error {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()
    
    if bc.conn != nil {
        return nil
    }
    
    // 创建连接选项
    opts := []grpc.DialOption{
        grpc.WithInsecure(),
        grpc.WithTimeout(bc.config.Timeout),
        grpc.WithKeepaliveParams(keepalive.ClientParameters{
            Time:                10 * time.Second,
            Timeout:             3 * time.Second,
            PermitWithoutStream: true,
        }),
        grpc.WithUnaryInterceptor(bc.unaryInterceptor),
    }
    
    // 建立连接
    conn, err := grpc.Dial(bc.config.Address, opts...)
    if err != nil {
        bc.logger.Errorf("Failed to connect to %s: %v", bc.config.Address, err)
        return err
    }
    
    bc.conn = conn
    bc.healthy = true
    
    // 启动健康检查
    go bc.healthCheck()
    
    bc.logger.Infof("Connected to %s", bc.config.Address)
    return nil
}

func (bc *BaseClient) Close() error {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()
    
    if bc.conn != nil {
        err := bc.conn.Close()
        bc.conn = nil
        bc.healthy = false
        return err
    }
    
    return nil
}

func (bc *BaseClient) IsHealthy() bool {
    bc.mutex.RLock()
    defer bc.mutex.RUnlock()
    return bc.healthy
}
```

### 2. **用户服务客户端实现**

```go
type userServiceClient struct {
    *BaseClient
    client pb.UserServiceClient
}

func NewUserServiceClient(config *ClientConfig) UserServiceClient {
    baseClient := NewBaseClient(config)
    return &userServiceClient{
        BaseClient: baseClient,
    }
}

func (usc *userServiceClient) Connect() error {
    if err := usc.BaseClient.Connect(); err != nil {
        return err
    }
    
    usc.client = pb.NewUserServiceClient(usc.conn)
    return nil
}

func (usc *userServiceClient) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
    if !usc.IsHealthy() {
        return nil, errors.New("user service client is not healthy")
    }
    
    // 添加超时控制
    ctx, cancel := context.WithTimeout(ctx, usc.config.Timeout)
    defer cancel()
    
    // 调用服务
    resp, err := usc.client.Register(ctx, req)
    if err != nil {
        usc.metrics.IncErrorCount("Register")
        usc.logger.Errorf("Register failed: %v", err)
        return nil, err
    }
    
    usc.metrics.IncSuccessCount("Register")
    return resp, nil
}

func (usc *userServiceClient) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    if !usc.IsHealthy() {
        return nil, errors.New("user service client is not healthy")
    }
    
    ctx, cancel := context.WithTimeout(ctx, usc.config.Timeout)
    defer cancel()
    
    resp, err := usc.client.Login(ctx, req)
    if err != nil {
        usc.metrics.IncErrorCount("Login")
        usc.logger.Errorf("Login failed: %v", err)
        return nil, err
    }
    
    usc.metrics.IncSuccessCount("Login")
    return resp, nil
}

// 其他方法实现类似...
```

## 🔄 连接池管理

### 1. **连接池设计**

```go
type ConnectionPool struct {
    address    string
    maxSize    int
    currentSize int
    connections chan *grpc.ClientConn
    factory    func() (*grpc.ClientConn, error)
    mutex      sync.Mutex
}

func NewConnectionPool(address string, maxSize int) *ConnectionPool {
    pool := &ConnectionPool{
        address:     address,
        maxSize:     maxSize,
        connections: make(chan *grpc.ClientConn, maxSize),
        factory: func() (*grpc.ClientConn, error) {
            return grpc.Dial(address, grpc.WithInsecure())
        },
    }
    
    // 预创建连接
    for i := 0; i < maxSize/2; i++ {
        if conn, err := pool.factory(); err == nil {
            pool.connections <- conn
            pool.currentSize++
        }
    }
    
    return pool
}

func (cp *ConnectionPool) Get() (*grpc.ClientConn, error) {
    select {
    case conn := <-cp.connections:
        // 检查连接状态
        if conn.GetState() == connectivity.Ready {
            return conn, nil
        }
        // 连接不可用，创建新连接
        conn.Close()
        cp.mutex.Lock()
        cp.currentSize--
        cp.mutex.Unlock()
        fallthrough
    default:
        // 创建新连接
        return cp.factory()
    }
}

func (cp *ConnectionPool) Put(conn *grpc.ClientConn) {
    if conn == nil {
        return
    }
    
    select {
    case cp.connections <- conn:
        // 成功放回连接池
    default:
        // 连接池已满，关闭连接
        conn.Close()
        cp.mutex.Lock()
        cp.currentSize--
        cp.mutex.Unlock()
    }
}

func (cp *ConnectionPool) Close() {
    close(cp.connections)
    for conn := range cp.connections {
        conn.Close()
    }
}
```

## 🛡️ 错误处理与重试

### 1. **重试机制**

```go
type RetryConfig struct {
    MaxRetries    int
    InitialDelay  time.Duration
    MaxDelay      time.Duration
    Multiplier    float64
    RetryableErrors []codes.Code
}

func (bc *BaseClient) withRetry(ctx context.Context, operation func() error) error {
    var lastErr error
    delay := bc.config.RetryInterval
    
    for attempt := 0; attempt <= bc.config.MaxRetries; attempt++ {
        if attempt > 0 {
            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-time.After(delay):
                delay = time.Duration(float64(delay) * 1.5) // 指数退避
                if delay > time.Second*30 {
                    delay = time.Second * 30
                }
            }
        }
        
        if err := operation(); err != nil {
            lastErr = err
            
            // 检查是否为可重试错误
            if !bc.isRetryableError(err) {
                return err
            }
            
            bc.logger.Warnf("Attempt %d failed: %v", attempt+1, err)
            continue
        }
        
        return nil
    }
    
    return fmt.Errorf("operation failed after %d attempts: %v", bc.config.MaxRetries+1, lastErr)
}

func (bc *BaseClient) isRetryableError(err error) bool {
    if err == nil {
        return false
    }
    
    // gRPC状态码检查
    if st, ok := status.FromError(err); ok {
        switch st.Code() {
        case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted:
            return true
        }
    }
    
    // 网络错误检查
    if strings.Contains(err.Error(), "connection refused") ||
       strings.Contains(err.Error(), "timeout") {
        return true
    }
    
    return false
}
```

### 2. **熔断器实现**

```go
type CircuitBreaker struct {
    maxFailures  int
    resetTimeout time.Duration
    state        CircuitState
    failures     int
    lastFailTime time.Time
    mutex        sync.RWMutex
}

type CircuitState int

const (
    StateClosed CircuitState = iota
    StateOpen
    StateHalfOpen
)

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        maxFailures:  maxFailures,
        resetTimeout: resetTimeout,
        state:        StateClosed,
    }
}

func (cb *CircuitBreaker) Call(operation func() error) error {
    if !cb.allowRequest() {
        return errors.New("circuit breaker is open")
    }
    
    err := operation()
    cb.recordResult(err)
    return err
}

func (cb *CircuitBreaker) allowRequest() bool {
    cb.mutex.RLock()
    defer cb.mutex.RUnlock()
    
    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        return time.Since(cb.lastFailTime) > cb.resetTimeout
    case StateHalfOpen:
        return true
    default:
        return false
    }
}

func (cb *CircuitBreaker) recordResult(err error) {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    if err != nil {
        cb.failures++
        cb.lastFailTime = time.Now()
        
        if cb.failures >= cb.maxFailures {
            cb.state = StateOpen
        }
    } else {
        cb.failures = 0
        cb.state = StateClosed
    }
}
```

## 📊 监控与指标

### 1. **客户端指标**

```go
type ClientMetrics struct {
    requestCount    *prometheus.CounterVec
    requestDuration *prometheus.HistogramVec
    errorCount      *prometheus.CounterVec
    connectionCount prometheus.Gauge
}

func NewClientMetrics() *ClientMetrics {
    return &ClientMetrics{
        requestCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "grpc_client_requests_total",
                Help: "Total number of gRPC client requests",
            },
            []string{"service", "method", "status"},
        ),
        requestDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "grpc_client_request_duration_seconds",
                Help: "Duration of gRPC client requests",
            },
            []string{"service", "method"},
        ),
        errorCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "grpc_client_errors_total",
                Help: "Total number of gRPC client errors",
            },
            []string{"service", "method", "error_code"},
        ),
        connectionCount: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "grpc_client_connections",
                Help: "Number of active gRPC client connections",
            },
        ),
    }
}

func (cm *ClientMetrics) IncRequestCount(service, method, status string) {
    cm.requestCount.WithLabelValues(service, method, status).Inc()
}

func (cm *ClientMetrics) ObserveRequestDuration(service, method string, duration time.Duration) {
    cm.requestDuration.WithLabelValues(service, method).Observe(duration.Seconds())
}

func (cm *ClientMetrics) IncErrorCount(service, method, errorCode string) {
    cm.errorCount.WithLabelValues(service, method, errorCode).Inc()
}
```

## 🔧 客户端配置

### 1. **配置文件示例**

```yaml
# gRPC客户端配置
grpc_clients:
  user_service:
    address: "user:9081"
    timeout: 5s
    max_retries: 3
    retry_interval: 1s
    pool_size: 10
    health_check_interval: 30s
    
  movie_service:
    address: "movie:9082"
    timeout: 5s
    max_retries: 3
    retry_interval: 1s
    pool_size: 10
    health_check_interval: 30s
    
  comment_service:
    address: "comment:9083"
    timeout: 5s
    max_retries: 3
    retry_interval: 1s
    pool_size: 10
    health_check_interval: 30s

# 熔断器配置
circuit_breaker:
  max_failures: 5
  reset_timeout: 60s
```

## 📋 使用示例

### 1. **客户端初始化**

```go
func InitGRPCClients(config *Config) (*GRPCClients, error) {
    clients := &GRPCClients{}
    
    // 初始化用户服务客户端
    userClient := NewUserServiceClient(&config.GRPCClients.UserService)
    if err := userClient.Connect(); err != nil {
        return nil, fmt.Errorf("failed to connect user service: %v", err)
    }
    clients.UserService = userClient
    
    // 初始化电影服务客户端
    movieClient := NewMovieServiceClient(&config.GRPCClients.MovieService)
    if err := movieClient.Connect(); err != nil {
        return nil, fmt.Errorf("failed to connect movie service: %v", err)
    }
    clients.MovieService = movieClient
    
    // 初始化评论服务客户端
    commentClient := NewCommentServiceClient(&config.GRPCClients.CommentService)
    if err := commentClient.Connect(); err != nil {
        return nil, fmt.Errorf("failed to connect comment service: %v", err)
    }
    clients.CommentService = commentClient
    
    return clients, nil
}
```

### 2. **在控制器中使用**

```go
func (uc *UserController) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 调用用户服务
    resp, err := uc.clients.UserService.Login(c.Request.Context(), &pb.LoginRequest{
        Email:    req.Email,
        Password: req.Password,
    })
    
    if err != nil {
        c.JSON(500, gin.H{"error": "Login failed"})
        return
    }
    
    c.JSON(200, gin.H{
        "token": resp.Token,
        "user":  resp.User,
    })
}
```

## 📝 总结

gRPC客户端封装为主页服务提供了强大的服务调用能力：

**关键特性**：
1. **连接管理**：自动连接、重连和健康检查
2. **错误处理**：重试机制和熔断保护
3. **性能优化**：连接池和负载均衡
4. **监控支持**：完整的指标收集

**最佳实践**：
- 使用连接池提升性能
- 实现重试和熔断机制
- 添加完整的监控指标
- 统一错误处理和日志记录

下一步，我们将实现服务注册与发现机制，为分布式部署做准备。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第24步：服务注册与发现](24-service-discovery.md)
