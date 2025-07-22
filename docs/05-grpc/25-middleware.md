# 第25步：gRPC 中间件

## 📋 概述

gRPC中间件是处理横切关注点的重要机制，它允许我们在不修改业务逻辑的情况下，为服务调用添加认证、日志、监控、限流等功能。对于MovieInfo项目，我们需要实现一套完整的中间件体系。

## 🎯 设计目标

### 1. **横切关注点处理**
- 认证授权
- 日志记录
- 性能监控
- 错误处理

### 2. **可组合性**
- 中间件可以自由组合
- 支持条件执行
- 易于扩展和维护

### 3. **性能优化**
- 最小化性能开销
- 支持异步处理
- 资源高效利用

### 4. **可观测性**
- 完整的调用链追踪
- 详细的性能指标
- 结构化日志输出

## 🏗️ 中间件架构

### 1. **中间件执行流程**

```
客户端请求 → 认证中间件 → 日志中间件 → 监控中间件 → 限流中间件 → 业务处理
     ↓            ↓           ↓           ↓           ↓           ↓
   发起调用     验证Token    记录请求     收集指标     检查限制    执行逻辑
     ↑            ↑           ↑           ↑           ↑           ↑
客户端响应 ← 认证响应 ← 日志响应 ← 监控响应 ← 限流响应 ← 业务响应
```

### 2. **中间件类型**

```go
// 一元RPC中间件
type UnaryServerInterceptor func(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error)

// 流式RPC中间件
type StreamServerInterceptor func(
    srv interface{},
    ss grpc.ServerStream,
    info *grpc.StreamServerInfo,
    handler grpc.StreamHandler,
) error

// 客户端中间件
type UnaryClientInterceptor func(
    ctx context.Context,
    method string,
    req, reply interface{},
    cc *grpc.ClientConn,
    invoker grpc.UnaryInvoker,
    opts ...grpc.CallOption,
) error
```

## 🔐 认证中间件

### 1. **JWT认证中间件**

```go
type AuthInterceptor struct {
    jwtManager *JWTManager
    logger     *logrus.Logger
}

func NewAuthInterceptor(jwtManager *JWTManager) *AuthInterceptor {
    return &AuthInterceptor{
        jwtManager: jwtManager,
        logger:     logrus.New(),
    }
}

func (ai *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        // 检查是否需要认证
        if ai.isPublicMethod(info.FullMethod) {
            return handler(ctx, req)
        }
        
        // 提取Token
        token, err := ai.extractToken(ctx)
        if err != nil {
            ai.logger.Warnf("Failed to extract token: %v", err)
            return nil, status.Errorf(codes.Unauthenticated, "missing or invalid token")
        }
        
        // 验证Token
        claims, err := ai.jwtManager.Verify(token)
        if err != nil {
            ai.logger.Warnf("Token verification failed: %v", err)
            return nil, status.Errorf(codes.Unauthenticated, "invalid token")
        }
        
        // 将用户信息添加到上下文
        ctx = ai.addUserToContext(ctx, claims)
        
        // 调用下一个处理器
        return handler(ctx, req)
    }
}

func (ai *AuthInterceptor) extractToken(ctx context.Context) (string, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return "", errors.New("missing metadata")
    }
    
    values := md["authorization"]
    if len(values) == 0 {
        return "", errors.New("missing authorization header")
    }
    
    authHeader := values[0]
    if !strings.HasPrefix(authHeader, "Bearer ") {
        return "", errors.New("invalid authorization header format")
    }
    
    return strings.TrimPrefix(authHeader, "Bearer "), nil
}

func (ai *AuthInterceptor) isPublicMethod(method string) bool {
    publicMethods := []string{
        "/user.UserService/Register",
        "/user.UserService/Login",
        "/movie.MovieService/GetMovieList",
        "/movie.MovieService/GetMovieDetail",
        "/movie.MovieService/SearchMovies",
    }
    
    for _, publicMethod := range publicMethods {
        if method == publicMethod {
            return true
        }
    }
    
    return false
}

func (ai *AuthInterceptor) addUserToContext(ctx context.Context, claims *JWTClaims) context.Context {
    return context.WithValue(ctx, "user", &User{
        ID:    claims.UserID,
        Email: claims.Email,
        Role:  claims.Role,
    })
}
```

### 2. **权限检查中间件**

```go
type AuthorizationInterceptor struct {
    rbac   *RBACManager
    logger *logrus.Logger
}

func NewAuthorizationInterceptor(rbac *RBACManager) *AuthorizationInterceptor {
    return &AuthorizationInterceptor{
        rbac:   rbac,
        logger: logrus.New(),
    }
}

func (azi *AuthorizationInterceptor) Unary() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        // 获取用户信息
        user, ok := ctx.Value("user").(*User)
        if !ok {
            return nil, status.Errorf(codes.Unauthenticated, "user not found in context")
        }
        
        // 检查权限
        if !azi.rbac.HasPermission(user.Role, info.FullMethod) {
            azi.logger.Warnf("Access denied for user %s to method %s", user.ID, info.FullMethod)
            return nil, status.Errorf(codes.PermissionDenied, "insufficient permissions")
        }
        
        return handler(ctx, req)
    }
}
```

## 📝 日志中间件

### 1. **请求日志中间件**

```go
type LoggingInterceptor struct {
    logger *logrus.Logger
}

func NewLoggingInterceptor() *LoggingInterceptor {
    logger := logrus.New()
    logger.SetFormatter(&logrus.JSONFormatter{})
    
    return &LoggingInterceptor{
        logger: logger,
    }
}

func (li *LoggingInterceptor) Unary() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        start := time.Now()
        
        // 生成请求ID
        requestID := li.generateRequestID()
        ctx = context.WithValue(ctx, "request_id", requestID)
        
        // 记录请求开始
        li.logRequest(ctx, info.FullMethod, req, requestID)
        
        // 执行处理器
        resp, err := handler(ctx, req)
        
        // 记录请求结束
        duration := time.Since(start)
        li.logResponse(ctx, info.FullMethod, resp, err, duration, requestID)
        
        return resp, err
    }
}

func (li *LoggingInterceptor) logRequest(ctx context.Context, method string, req interface{}, requestID string) {
    fields := logrus.Fields{
        "request_id": requestID,
        "method":     method,
        "type":       "request",
    }
    
    // 添加用户信息（如果存在）
    if user, ok := ctx.Value("user").(*User); ok {
        fields["user_id"] = user.ID
    }
    
    // 添加客户端信息
    if peer, ok := peer.FromContext(ctx); ok {
        fields["client_ip"] = peer.Addr.String()
    }
    
    li.logger.WithFields(fields).Info("gRPC request started")
}

func (li *LoggingInterceptor) logResponse(ctx context.Context, method string, resp interface{}, err error, duration time.Duration, requestID string) {
    fields := logrus.Fields{
        "request_id": requestID,
        "method":     method,
        "duration":   duration.Milliseconds(),
        "type":       "response",
    }
    
    if err != nil {
        fields["error"] = err.Error()
        if st, ok := status.FromError(err); ok {
            fields["status_code"] = st.Code().String()
        }
        li.logger.WithFields(fields).Error("gRPC request failed")
    } else {
        fields["status"] = "success"
        li.logger.WithFields(fields).Info("gRPC request completed")
    }
}

func (li *LoggingInterceptor) generateRequestID() string {
    return fmt.Sprintf("%d-%s", time.Now().UnixNano(), uuid.New().String()[:8])
}
```

## 📊 监控中间件

### 1. **性能监控中间件**

```go
type MetricsInterceptor struct {
    requestCount    *prometheus.CounterVec
    requestDuration *prometheus.HistogramVec
    requestSize     *prometheus.HistogramVec
    responseSize    *prometheus.HistogramVec
}

func NewMetricsInterceptor() *MetricsInterceptor {
    return &MetricsInterceptor{
        requestCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "grpc_server_requests_total",
                Help: "Total number of gRPC requests",
            },
            []string{"service", "method", "status"},
        ),
        requestDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "grpc_server_request_duration_seconds",
                Help: "Duration of gRPC requests",
                Buckets: prometheus.DefBuckets,
            },
            []string{"service", "method"},
        ),
        requestSize: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "grpc_server_request_size_bytes",
                Help: "Size of gRPC requests",
                Buckets: prometheus.ExponentialBuckets(1, 2, 20),
            },
            []string{"service", "method"},
        ),
        responseSize: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "grpc_server_response_size_bytes",
                Help: "Size of gRPC responses",
                Buckets: prometheus.ExponentialBuckets(1, 2, 20),
            },
            []string{"service", "method"},
        ),
    }
}

func (mi *MetricsInterceptor) Unary() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        start := time.Now()
        
        // 解析服务和方法名
        service, method := mi.parseMethod(info.FullMethod)
        
        // 记录请求大小
        if reqSize := mi.getMessageSize(req); reqSize > 0 {
            mi.requestSize.WithLabelValues(service, method).Observe(float64(reqSize))
        }
        
        // 执行处理器
        resp, err := handler(ctx, req)
        
        // 记录指标
        duration := time.Since(start)
        status := "success"
        if err != nil {
            status = "error"
        }
        
        mi.requestCount.WithLabelValues(service, method, status).Inc()
        mi.requestDuration.WithLabelValues(service, method).Observe(duration.Seconds())
        
        // 记录响应大小
        if resp != nil {
            if respSize := mi.getMessageSize(resp); respSize > 0 {
                mi.responseSize.WithLabelValues(service, method).Observe(float64(respSize))
            }
        }
        
        return resp, err
    }
}

func (mi *MetricsInterceptor) parseMethod(fullMethod string) (string, string) {
    parts := strings.Split(fullMethod, "/")
    if len(parts) >= 3 {
        serviceParts := strings.Split(parts[1], ".")
        service := serviceParts[len(serviceParts)-1]
        method := parts[2]
        return service, method
    }
    return "unknown", "unknown"
}

func (mi *MetricsInterceptor) getMessageSize(msg interface{}) int {
    if msg == nil {
        return 0
    }
    
    if protoMsg, ok := msg.(proto.Message); ok {
        return proto.Size(protoMsg)
    }
    
    return 0
}
```

## 🚦 限流中间件

### 1. **令牌桶限流中间件**

```go
type RateLimitInterceptor struct {
    limiters map[string]*rate.Limiter
    mutex    sync.RWMutex
    logger   *logrus.Logger
}

func NewRateLimitInterceptor() *RateLimitInterceptor {
    return &RateLimitInterceptor{
        limiters: make(map[string]*rate.Limiter),
        logger:   logrus.New(),
    }
}

func (rli *RateLimitInterceptor) Unary() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        // 获取限流器
        limiter := rli.getLimiter(info.FullMethod)
        
        // 检查是否允许请求
        if !limiter.Allow() {
            rli.logger.Warnf("Rate limit exceeded for method: %s", info.FullMethod)
            return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
        }
        
        return handler(ctx, req)
    }
}

func (rli *RateLimitInterceptor) getLimiter(method string) *rate.Limiter {
    rli.mutex.RLock()
    limiter, exists := rli.limiters[method]
    rli.mutex.RUnlock()
    
    if exists {
        return limiter
    }
    
    rli.mutex.Lock()
    defer rli.mutex.Unlock()
    
    // 双重检查
    if limiter, exists := rli.limiters[method]; exists {
        return limiter
    }
    
    // 创建新的限流器（每秒100个请求，突发200个）
    limiter = rate.NewLimiter(100, 200)
    rli.limiters[method] = limiter
    
    return limiter
}
```

## 🔗 链路追踪中间件

### 1. **分布式追踪中间件**

```go
type TracingInterceptor struct {
    tracer opentracing.Tracer
    logger *logrus.Logger
}

func NewTracingInterceptor(tracer opentracing.Tracer) *TracingInterceptor {
    return &TracingInterceptor{
        tracer: tracer,
        logger: logrus.New(),
    }
}

func (ti *TracingInterceptor) Unary() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        // 从metadata中提取span上下文
        md, _ := metadata.FromIncomingContext(ctx)
        spanContext, err := ti.tracer.Extract(
            opentracing.TextMap,
            MetadataTextMap(md),
        )
        
        // 创建新的span
        span := ti.tracer.StartSpan(
            info.FullMethod,
            ext.RPCServerOption(spanContext),
        )
        defer span.Finish()
        
        // 设置span标签
        ext.Component.Set(span, "gRPC")
        ext.SpanKind.Set(span, ext.SpanKindRPCServerEnum)
        
        // 将span添加到上下文
        ctx = opentracing.ContextWithSpan(ctx, span)
        
        // 执行处理器
        resp, err := handler(ctx, req)
        
        // 记录错误信息
        if err != nil {
            ext.Error.Set(span, true)
            span.LogFields(
                log.String("error.kind", "grpc_error"),
                log.String("error.object", err.Error()),
            )
        }
        
        return resp, err
    }
}

// MetadataTextMap 实现 opentracing.TextMapReader 和 opentracing.TextMapWriter
type MetadataTextMap metadata.MD

func (m MetadataTextMap) ForeachKey(handler func(key, val string) error) error {
    for k, vs := range m {
        for _, v := range vs {
            if err := handler(k, v); err != nil {
                return err
            }
        }
    }
    return nil
}

func (m MetadataTextMap) Set(key, val string) {
    m[key] = append(m[key], val)
}
```

## 🔧 中间件组合

### 1. **中间件链构建器**

```go
type InterceptorChain struct {
    interceptors []grpc.UnaryServerInterceptor
}

func NewInterceptorChain() *InterceptorChain {
    return &InterceptorChain{
        interceptors: make([]grpc.UnaryServerInterceptor, 0),
    }
}

func (ic *InterceptorChain) Add(interceptor grpc.UnaryServerInterceptor) *InterceptorChain {
    ic.interceptors = append(ic.interceptors, interceptor)
    return ic
}

func (ic *InterceptorChain) Build() grpc.UnaryServerInterceptor {
    return grpc_middleware.ChainUnaryServer(ic.interceptors...)
}

// 使用示例
func setupInterceptors() grpc.UnaryServerInterceptor {
    // 创建各种中间件
    authInterceptor := NewAuthInterceptor(jwtManager)
    loggingInterceptor := NewLoggingInterceptor()
    metricsInterceptor := NewMetricsInterceptor()
    rateLimitInterceptor := NewRateLimitInterceptor()
    tracingInterceptor := NewTracingInterceptor(tracer)
    
    // 构建中间件链
    return NewInterceptorChain().
        Add(tracingInterceptor.Unary()).     // 链路追踪（最外层）
        Add(loggingInterceptor.Unary()).     // 日志记录
        Add(metricsInterceptor.Unary()).     // 性能监控
        Add(rateLimitInterceptor.Unary()).   // 限流控制
        Add(authInterceptor.Unary()).        // 认证授权（最内层）
        Build()
}
```

### 2. **服务器配置**

```go
func NewGRPCServer(config *Config) *grpc.Server {
    // 设置服务器选项
    opts := []grpc.ServerOption{
        grpc.UnaryInterceptor(setupInterceptors()),
        grpc.StreamInterceptor(setupStreamInterceptors()),
        grpc.KeepaliveParams(keepalive.ServerParameters{
            Time:    60 * time.Second,
            Timeout: 10 * time.Second,
        }),
        grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
            MinTime:             30 * time.Second,
            PermitWithoutStream: true,
        }),
    }
    
    return grpc.NewServer(opts...)
}
```

## 📊 中间件监控

### 1. **中间件性能监控**

```go
type MiddlewareMetrics struct {
    middlewareCount    *prometheus.CounterVec
    middlewareDuration *prometheus.HistogramVec
}

func NewMiddlewareMetrics() *MiddlewareMetrics {
    return &MiddlewareMetrics{
        middlewareCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "grpc_middleware_executions_total",
                Help: "Total number of middleware executions",
            },
            []string{"middleware", "status"},
        ),
        middlewareDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "grpc_middleware_duration_seconds",
                Help: "Duration of middleware execution",
            },
            []string{"middleware"},
        ),
    }
}

func (mm *MiddlewareMetrics) WrapInterceptor(name string, interceptor grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        start := time.Now()
        
        resp, err := interceptor(ctx, req, info, handler)
        
        duration := time.Since(start)
        status := "success"
        if err != nil {
            status = "error"
        }
        
        mm.middlewareCount.WithLabelValues(name, status).Inc()
        mm.middlewareDuration.WithLabelValues(name).Observe(duration.Seconds())
        
        return resp, err
    }
}
```

## 📝 总结

gRPC中间件为MovieInfo项目提供了强大的横切关注点处理能力：

**核心功能**：
1. **认证授权**：JWT认证和RBAC权限控制
2. **日志记录**：结构化请求响应日志
3. **性能监控**：详细的性能指标收集
4. **限流保护**：防止服务过载
5. **链路追踪**：分布式调用链追踪

**设计特点**：
- 可组合的中间件架构
- 最小化性能开销
- 完整的可观测性支持
- 易于扩展和维护

**最佳实践**：
- 合理安排中间件执行顺序
- 使用结构化日志和指标
- 实现优雅的错误处理
- 添加完整的监控和告警

下一步，我们将开始用户服务的具体实现，包括用户模型设计和核心功能开发。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第26步：用户模型设计](../06-user-service/26-user-model.md)
