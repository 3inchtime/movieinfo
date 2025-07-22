# 第24步：服务注册与发现

## 📋 概述

服务注册与发现是微服务架构的核心组件，它解决了服务间如何找到彼此的问题。对于MovieInfo项目，我们需要实现一个简单而有效的服务发现机制，为未来的微服务化做好准备。

## 🎯 设计目标

### 1. **服务自动发现**
- 服务启动时自动注册
- 服务停止时自动注销
- 动态发现可用服务

### 2. **健康检查**
- 定期检查服务健康状态
- 自动剔除不健康服务
- 服务恢复后自动重新加入

### 3. **负载均衡**
- 支持多种负载均衡策略
- 动态调整服务权重
- 故障转移机制

### 4. **配置简单**
- 最小化配置要求
- 支持多种部署环境
- 易于集成和使用

## 🏗️ 服务发现架构

### 1. **整体架构**

```
┌─────────────────────────────────────────────────────────────┐
│                    服务发现架构                              │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │  主页服务    │  │  用户服务    │  │  电影服务    │          │
│  │             │  │             │  │             │          │
│  │ • 服务消费者 │  │ • 服务提供者 │  │ • 服务提供者 │          │
│  │ • 负载均衡   │  │ • 健康检查   │  │ • 健康检查   │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
│         │                 │                 │               │
│         │                 │                 │               │
│         ▼                 ▼                 ▼               │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │              服务注册中心                                │ │
│  │                                                         │ │
│  │ • 服务注册    • 服务发现    • 健康检查    • 负载均衡     │ │
│  │ • 配置管理    • 事件通知    • 故障转移    • 监控统计     │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 2. **服务注册流程**

```
服务启动 → 读取配置 → 注册到注册中心 → 启动健康检查 → 开始提供服务
    ↓           ↓           ↓             ↓             ↓
  初始化     获取地址     发送注册请求    定期心跳      处理请求
```

### 3. **服务发现流程**

```
客户端启动 → 查询注册中心 → 获取服务列表 → 选择服务实例 → 发起调用
     ↓            ↓            ↓            ↓            ↓
   初始化      发送查询请求   缓存服务信息   负载均衡     建立连接
```

## 🔧 核心组件实现

### 1. **服务注册中心**

```go
// 服务实例信息
type ServiceInstance struct {
    ID       string            `json:"id"`
    Name     string            `json:"name"`
    Address  string            `json:"address"`
    Port     int               `json:"port"`
    Tags     []string          `json:"tags"`
    Meta     map[string]string `json:"meta"`
    Health   HealthStatus      `json:"health"`
    LastSeen time.Time         `json:"last_seen"`
}

type HealthStatus string

const (
    HealthStatusHealthy   HealthStatus = "healthy"
    HealthStatusUnhealthy HealthStatus = "unhealthy"
    HealthStatusUnknown   HealthStatus = "unknown"
)

// 服务注册中心接口
type ServiceRegistry interface {
    Register(instance *ServiceInstance) error
    Deregister(serviceID string) error
    Discover(serviceName string) ([]*ServiceInstance, error)
    Watch(serviceName string) (<-chan []*ServiceInstance, error)
    HealthCheck(serviceID string) error
}

// 内存实现的服务注册中心
type MemoryServiceRegistry struct {
    services map[string]map[string]*ServiceInstance // serviceName -> serviceID -> instance
    watchers map[string][]chan []*ServiceInstance   // serviceName -> watchers
    mutex    sync.RWMutex
    logger   *logrus.Logger
}

func NewMemoryServiceRegistry() *MemoryServiceRegistry {
    registry := &MemoryServiceRegistry{
        services: make(map[string]map[string]*ServiceInstance),
        watchers: make(map[string][]chan []*ServiceInstance),
        logger:   logrus.New(),
    }
    
    // 启动健康检查
    go registry.startHealthCheck()
    
    return registry
}

func (msr *MemoryServiceRegistry) Register(instance *ServiceInstance) error {
    msr.mutex.Lock()
    defer msr.mutex.Unlock()
    
    if msr.services[instance.Name] == nil {
        msr.services[instance.Name] = make(map[string]*ServiceInstance)
    }
    
    instance.LastSeen = time.Now()
    instance.Health = HealthStatusHealthy
    msr.services[instance.Name][instance.ID] = instance
    
    msr.logger.Infof("Service registered: %s (%s)", instance.Name, instance.ID)
    
    // 通知观察者
    msr.notifyWatchers(instance.Name)
    
    return nil
}

func (msr *MemoryServiceRegistry) Deregister(serviceID string) error {
    msr.mutex.Lock()
    defer msr.mutex.Unlock()
    
    for serviceName, instances := range msr.services {
        if instance, exists := instances[serviceID]; exists {
            delete(instances, serviceID)
            msr.logger.Infof("Service deregistered: %s (%s)", serviceName, serviceID)
            
            // 通知观察者
            msr.notifyWatchers(serviceName)
            return nil
        }
    }
    
    return fmt.Errorf("service not found: %s", serviceID)
}

func (msr *MemoryServiceRegistry) Discover(serviceName string) ([]*ServiceInstance, error) {
    msr.mutex.RLock()
    defer msr.mutex.RUnlock()
    
    instances := msr.services[serviceName]
    if instances == nil {
        return nil, fmt.Errorf("service not found: %s", serviceName)
    }
    
    var healthyInstances []*ServiceInstance
    for _, instance := range instances {
        if instance.Health == HealthStatusHealthy {
            healthyInstances = append(healthyInstances, instance)
        }
    }
    
    return healthyInstances, nil
}

func (msr *MemoryServiceRegistry) Watch(serviceName string) (<-chan []*ServiceInstance, error) {
    msr.mutex.Lock()
    defer msr.mutex.Unlock()
    
    watcher := make(chan []*ServiceInstance, 1)
    msr.watchers[serviceName] = append(msr.watchers[serviceName], watcher)
    
    // 发送当前状态
    if instances, err := msr.Discover(serviceName); err == nil {
        select {
        case watcher <- instances:
        default:
        }
    }
    
    return watcher, nil
}

func (msr *MemoryServiceRegistry) notifyWatchers(serviceName string) {
    if watchers, exists := msr.watchers[serviceName]; exists {
        if instances, err := msr.Discover(serviceName); err == nil {
            for _, watcher := range watchers {
                select {
                case watcher <- instances:
                default:
                    // 非阻塞发送
                }
            }
        }
    }
}

func (msr *MemoryServiceRegistry) startHealthCheck() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        msr.performHealthCheck()
    }
}

func (msr *MemoryServiceRegistry) performHealthCheck() {
    msr.mutex.Lock()
    defer msr.mutex.Unlock()
    
    now := time.Now()
    for serviceName, instances := range msr.services {
        for serviceID, instance := range instances {
            // 检查最后心跳时间
            if now.Sub(instance.LastSeen) > time.Minute*2 {
                instance.Health = HealthStatusUnhealthy
                msr.logger.Warnf("Service unhealthy: %s (%s)", serviceName, serviceID)
                msr.notifyWatchers(serviceName)
            }
        }
    }
}
```

### 2. **服务注册客户端**

```go
// 服务注册客户端
type ServiceRegistrar struct {
    registry     ServiceRegistry
    instance     *ServiceInstance
    stopCh       chan struct{}
    healthCheck  HealthChecker
    logger       *logrus.Logger
}

type HealthChecker interface {
    Check() error
}

func NewServiceRegistrar(registry ServiceRegistry, instance *ServiceInstance) *ServiceRegistrar {
    return &ServiceRegistrar{
        registry: registry,
        instance: instance,
        stopCh:   make(chan struct{}),
        logger:   logrus.New(),
    }
}

func (sr *ServiceRegistrar) Start() error {
    // 注册服务
    if err := sr.registry.Register(sr.instance); err != nil {
        return fmt.Errorf("failed to register service: %v", err)
    }
    
    // 启动心跳
    go sr.startHeartbeat()
    
    sr.logger.Infof("Service registrar started for %s", sr.instance.Name)
    return nil
}

func (sr *ServiceRegistrar) Stop() error {
    close(sr.stopCh)
    
    // 注销服务
    if err := sr.registry.Deregister(sr.instance.ID); err != nil {
        sr.logger.Errorf("Failed to deregister service: %v", err)
        return err
    }
    
    sr.logger.Infof("Service registrar stopped for %s", sr.instance.Name)
    return nil
}

func (sr *ServiceRegistrar) startHeartbeat() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := sr.sendHeartbeat(); err != nil {
                sr.logger.Errorf("Failed to send heartbeat: %v", err)
            }
        case <-sr.stopCh:
            return
        }
    }
}

func (sr *ServiceRegistrar) sendHeartbeat() error {
    // 执行健康检查
    if sr.healthCheck != nil {
        if err := sr.healthCheck.Check(); err != nil {
            sr.instance.Health = HealthStatusUnhealthy
            sr.logger.Warnf("Health check failed: %v", err)
        } else {
            sr.instance.Health = HealthStatusHealthy
        }
    }
    
    // 更新最后心跳时间
    sr.instance.LastSeen = time.Now()
    
    // 重新注册（更新信息）
    return sr.registry.Register(sr.instance)
}
```

### 3. **服务发现客户端**

```go
// 服务发现客户端
type ServiceDiscoverer struct {
    registry    ServiceRegistry
    serviceName string
    instances   []*ServiceInstance
    balancer    LoadBalancer
    mutex       sync.RWMutex
    logger      *logrus.Logger
}

func NewServiceDiscoverer(registry ServiceRegistry, serviceName string) *ServiceDiscoverer {
    discoverer := &ServiceDiscoverer{
        registry:    registry,
        serviceName: serviceName,
        balancer:    NewRoundRobinBalancer(),
        logger:      logrus.New(),
    }
    
    // 启动服务监听
    go discoverer.startWatching()
    
    return discoverer
}

func (sd *ServiceDiscoverer) GetInstance() (*ServiceInstance, error) {
    sd.mutex.RLock()
    defer sd.mutex.RUnlock()
    
    if len(sd.instances) == 0 {
        return nil, errors.New("no healthy instances available")
    }
    
    return sd.balancer.Select(sd.instances)
}

func (sd *ServiceDiscoverer) GetAllInstances() []*ServiceInstance {
    sd.mutex.RLock()
    defer sd.mutex.RUnlock()
    
    // 返回副本
    instances := make([]*ServiceInstance, len(sd.instances))
    copy(instances, sd.instances)
    return instances
}

func (sd *ServiceDiscoverer) startWatching() {
    watcher, err := sd.registry.Watch(sd.serviceName)
    if err != nil {
        sd.logger.Errorf("Failed to watch service %s: %v", sd.serviceName, err)
        return
    }
    
    for instances := range watcher {
        sd.mutex.Lock()
        sd.instances = instances
        sd.mutex.Unlock()
        
        sd.logger.Infof("Service instances updated for %s: %d instances", 
            sd.serviceName, len(instances))
    }
}
```

## ⚖️ 负载均衡策略

### 1. **负载均衡接口**

```go
type LoadBalancer interface {
    Select(instances []*ServiceInstance) (*ServiceInstance, error)
}

// 轮询负载均衡
type RoundRobinBalancer struct {
    counter uint64
}

func NewRoundRobinBalancer() *RoundRobinBalancer {
    return &RoundRobinBalancer{}
}

func (rrb *RoundRobinBalancer) Select(instances []*ServiceInstance) (*ServiceInstance, error) {
    if len(instances) == 0 {
        return nil, errors.New("no instances available")
    }
    
    index := atomic.AddUint64(&rrb.counter, 1) % uint64(len(instances))
    return instances[index], nil
}

// 随机负载均衡
type RandomBalancer struct {
    rand *rand.Rand
    mutex sync.Mutex
}

func NewRandomBalancer() *RandomBalancer {
    return &RandomBalancer{
        rand: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

func (rb *RandomBalancer) Select(instances []*ServiceInstance) (*ServiceInstance, error) {
    if len(instances) == 0 {
        return nil, errors.New("no instances available")
    }
    
    rb.mutex.Lock()
    index := rb.rand.Intn(len(instances))
    rb.mutex.Unlock()
    
    return instances[index], nil
}

// 加权轮询负载均衡
type WeightedRoundRobinBalancer struct {
    weights map[string]int
    current map[string]int
    mutex   sync.Mutex
}

func NewWeightedRoundRobinBalancer(weights map[string]int) *WeightedRoundRobinBalancer {
    return &WeightedRoundRobinBalancer{
        weights: weights,
        current: make(map[string]int),
    }
}

func (wrrb *WeightedRoundRobinBalancer) Select(instances []*ServiceInstance) (*ServiceInstance, error) {
    if len(instances) == 0 {
        return nil, errors.New("no instances available")
    }
    
    wrrb.mutex.Lock()
    defer wrrb.mutex.Unlock()
    
    var selected *ServiceInstance
    maxWeight := -1
    
    for _, instance := range instances {
        weight := wrrb.weights[instance.ID]
        if weight == 0 {
            weight = 1 // 默认权重
        }
        
        wrrb.current[instance.ID] += weight
        
        if wrrb.current[instance.ID] > maxWeight {
            maxWeight = wrrb.current[instance.ID]
            selected = instance
        }
    }
    
    if selected != nil {
        wrrb.current[selected.ID] -= wrrb.getTotalWeight(instances)
    }
    
    return selected, nil
}

func (wrrb *WeightedRoundRobinBalancer) getTotalWeight(instances []*ServiceInstance) int {
    total := 0
    for _, instance := range instances {
        weight := wrrb.weights[instance.ID]
        if weight == 0 {
            weight = 1
        }
        total += weight
    }
    return total
}
```

## 🔧 集成示例

### 1. **服务端集成**

```go
// 在服务启动时注册
func StartUserService(config *Config) error {
    // 创建服务注册中心客户端
    registry := NewMemoryServiceRegistry()
    
    // 创建服务实例信息
    instance := &ServiceInstance{
        ID:      fmt.Sprintf("user-service-%s", uuid.New().String()),
        Name:    "user-service",
        Address: config.Server.Host,
        Port:    config.Server.GRPCPort,
        Tags:    []string{"user", "auth"},
        Meta: map[string]string{
            "version": "1.0.0",
            "region":  "local",
        },
    }
    
    // 创建服务注册器
    registrar := NewServiceRegistrar(registry, instance)
    
    // 启动服务注册
    if err := registrar.Start(); err != nil {
        return fmt.Errorf("failed to start service registrar: %v", err)
    }
    
    // 启动gRPC服务器
    server := grpc.NewServer()
    pb.RegisterUserServiceServer(server, &userServiceImpl{})
    
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Server.GRPCPort))
    if err != nil {
        return fmt.Errorf("failed to listen: %v", err)
    }
    
    // 优雅关闭处理
    go func() {
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
        <-sigCh
        
        registrar.Stop()
        server.GracefulStop()
    }()
    
    return server.Serve(listener)
}
```

### 2. **客户端集成**

```go
// 在客户端使用服务发现
type UserServiceClient struct {
    discoverer *ServiceDiscoverer
    connPool   map[string]*grpc.ClientConn
    mutex      sync.RWMutex
}

func NewUserServiceClient(registry ServiceRegistry) *UserServiceClient {
    return &UserServiceClient{
        discoverer: NewServiceDiscoverer(registry, "user-service"),
        connPool:   make(map[string]*grpc.ClientConn),
    }
}

func (usc *UserServiceClient) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    // 获取服务实例
    instance, err := usc.discoverer.GetInstance()
    if err != nil {
        return nil, fmt.Errorf("failed to get service instance: %v", err)
    }
    
    // 获取连接
    conn, err := usc.getConnection(instance)
    if err != nil {
        return nil, fmt.Errorf("failed to get connection: %v", err)
    }
    
    // 创建客户端并调用
    client := pb.NewUserServiceClient(conn)
    return client.Login(ctx, req)
}

func (usc *UserServiceClient) getConnection(instance *ServiceInstance) (*grpc.ClientConn, error) {
    address := fmt.Sprintf("%s:%d", instance.Address, instance.Port)
    
    usc.mutex.RLock()
    if conn, exists := usc.connPool[address]; exists {
        usc.mutex.RUnlock()
        return conn, nil
    }
    usc.mutex.RUnlock()
    
    usc.mutex.Lock()
    defer usc.mutex.Unlock()
    
    // 双重检查
    if conn, exists := usc.connPool[address]; exists {
        return conn, nil
    }
    
    // 创建新连接
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    
    usc.connPool[address] = conn
    return conn, nil
}
```

## 📊 监控与指标

### 1. **服务发现指标**

```go
type ServiceDiscoveryMetrics struct {
    registeredServices prometheus.Gauge
    healthyServices    prometheus.Gauge
    discoveryRequests  *prometheus.CounterVec
    registrationTime   *prometheus.HistogramVec
}

func NewServiceDiscoveryMetrics() *ServiceDiscoveryMetrics {
    return &ServiceDiscoveryMetrics{
        registeredServices: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "service_discovery_registered_services",
                Help: "Number of registered services",
            },
        ),
        healthyServices: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "service_discovery_healthy_services",
                Help: "Number of healthy services",
            },
        ),
        discoveryRequests: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "service_discovery_requests_total",
                Help: "Total number of service discovery requests",
            },
            []string{"service_name", "status"},
        ),
        registrationTime: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "service_discovery_registration_duration_seconds",
                Help: "Time taken to register a service",
            },
            []string{"service_name"},
        ),
    }
}
```

## 📝 总结

服务注册与发现为MovieInfo项目提供了强大的服务管理能力：

**核心功能**：
1. **自动注册**：服务启动时自动注册到注册中心
2. **健康检查**：定期检查服务健康状态
3. **动态发现**：客户端动态发现可用服务
4. **负载均衡**：支持多种负载均衡策略

**技术特点**：
- 简单易用的API接口
- 支持多种负载均衡策略
- 完整的健康检查机制
- 实时的服务状态更新

**扩展能力**：
- 支持分布式部署
- 可集成外部注册中心（如Consul、Etcd）
- 支持服务网格演进

下一步，我们将实现gRPC中间件，为服务调用添加认证、日志、监控等横切关注点。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第25步：gRPC中间件](25-middleware.md)
