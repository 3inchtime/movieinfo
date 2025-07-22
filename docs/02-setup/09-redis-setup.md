# 2.4 Redis 缓存配置

## 概述

Redis作为MovieInfo项目的缓存层，承担着提升系统性能、减少数据库压力的重要职责。合理的Redis配置和缓存策略设计，能够显著提升用户体验和系统吞吐量。本文档将详细介绍Redis的高级配置、缓存策略设计和性能优化。

## 为什么需要专门的缓存配置？

### 1. **性能提升**
- **响应时间**：内存访问比磁盘访问快1000倍以上
- **并发处理**：减少数据库连接压力，提升并发能力
- **热点数据**：将频繁访问的数据放入缓存
- **计算缓存**：缓存复杂计算结果，避免重复计算

### 2. **用户体验**
- **页面加载**：缓存页面数据，快速响应用户请求
- **搜索优化**：缓存搜索结果，提升搜索体验
- **个性化推荐**：缓存用户推荐数据
- **实时更新**：缓存实时统计数据

### 3. **系统稳定性**
- **降低数据库压力**：减少数据库查询次数
- **故障隔离**：缓存失效不影响核心功能
- **流量削峰**：在高并发时起到缓冲作用
- **容错机制**：提供降级方案

## Redis高级配置

### 1. **内存管理配置**

#### 1.1 内存策略配置
```conf
# redis.conf

# 最大内存限制
maxmemory 4gb

# 内存淘汰策略
maxmemory-policy allkeys-lru

# 内存淘汰策略说明：
# noeviction: 不淘汰，内存满时返回错误
# allkeys-lru: 在所有键中使用LRU算法淘汰
# allkeys-lfu: 在所有键中使用LFU算法淘汰
# volatile-lru: 在设置了过期时间的键中使用LRU淘汰
# volatile-lfu: 在设置了过期时间的键中使用LFU淘汰
# volatile-random: 在设置了过期时间的键中随机淘汰
# volatile-ttl: 淘汰即将过期的键

# LRU样本数量
maxmemory-samples 5

# 内存使用报告
memory-usage-threshold 80
```

#### 1.2 持久化配置
```conf
# RDB持久化配置
save 900 1      # 900秒内至少1个键变化时保存
save 300 10     # 300秒内至少10个键变化时保存
save 60 10000   # 60秒内至少10000个键变化时保存

# RDB文件配置
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
dir /data

# 后台保存出错时停止写入
stop-writes-on-bgsave-error yes

# AOF持久化配置
appendonly yes
appendfilename "appendonly.aof"

# AOF同步策略
appendfsync everysec  # 每秒同步一次（推荐）
# appendfsync always  # 每次写入都同步（最安全但性能差）
# appendfsync no      # 由操作系统决定（性能最好但可能丢失数据）

# AOF重写配置
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb

# AOF文件损坏时的处理
aof-load-truncated yes
aof-use-rdb-preamble yes
```

### 2. **网络和连接配置**

#### 2.1 连接优化
```conf
# 最大客户端连接数
maxclients 10000

# 客户端空闲超时时间（秒）
timeout 300

# TCP keepalive
tcp-keepalive 300

# TCP监听队列长度
tcp-backlog 511

# 禁用Nagle算法
tcp-nodelay yes
```

#### 2.2 安全配置
```conf
# 绑定地址
bind 127.0.0.1 192.168.1.100

# 保护模式
protected-mode yes

# 密码认证
requirepass movieinfo_redis_secure_password_2024

# 重命名危险命令
rename-command FLUSHDB ""
rename-command FLUSHALL ""
rename-command KEYS ""
rename-command CONFIG "CONFIG_a1b2c3d4"
rename-command SHUTDOWN "SHUTDOWN_e5f6g7h8"
rename-command DEBUG ""
rename-command EVAL "EVAL_i9j0k1l2"

# 禁用某些命令
rename-command DEL ""
```

### 3. **性能优化配置**

#### 3.1 慢查询配置
```conf
# 慢查询阈值（微秒）
slowlog-log-slower-than 10000

# 慢查询日志最大长度
slowlog-max-len 128
```

#### 3.2 键空间通知
```conf
# 键空间通知配置
notify-keyspace-events "Ex"

# 通知类型说明：
# K: 键空间通知，所有通知以__keyspace@<db>__为前缀
# E: 键事件通知，所有通知以__keyevent@<db>__为前缀
# g: DEL、EXPIRE、RENAME等类型无关的通用命令的通知
# $: 字符串命令的通知
# l: 列表命令的通知
# s: 集合命令的通知
# h: 哈希命令的通知
# z: 有序集合命令的通知
# x: 过期事件：每当有过期键被删除时发送
# e: 驱逐(evict)事件：每当有键因为maxmemory政策而被删除时发送
# A: 参数g$lshzxe的别名
```

## 缓存策略设计

### 1. **缓存键设计规范**

#### 1.1 键命名规范
```go
// 键命名规范
const (
    // 用户相关缓存
    UserCacheKey        = "user:%d"                    // user:123
    UserSessionKey      = "session:%s"                 // session:token_abc
    UserProfileKey      = "profile:%d"                 // profile:123
    
    // 电影相关缓存
    MovieCacheKey       = "movie:%d"                   // movie:456
    MovieListKey        = "movies:list:%s:%d"          // movies:list:action:1
    MovieSearchKey      = "search:%s"                  // search:hash_of_query
    MovieHotKey         = "movies:hot"                 // 热门电影列表
    MovieLatestKey      = "movies:latest"              // 最新电影列表
    
    // 评论相关缓存
    CommentListKey      = "comments:movie:%d:%d"       // comments:movie:456:1
    CommentCountKey     = "comments:count:movie:%d"    // comments:count:movie:456
    
    // 评分相关缓存
    RatingKey           = "rating:movie:%d"            // rating:movie:456
    RatingStatsKey      = "rating:stats:movie:%d"      // rating:stats:movie:456
    
    // 统计相关缓存
    StatsViewKey        = "stats:view:movie:%d"        // stats:view:movie:456
    StatsUserKey        = "stats:user:%d"              // stats:user:123
    
    // 配置相关缓存
    ConfigKey           = "config:%s"                  // config:categories
    
    // 临时数据缓存
    TempDataKey         = "temp:%s:%d"                 // temp:reset_token:123
)
```

#### 1.2 键过期策略
```go
// 缓存过期时间配置
const (
    // 短期缓存（5分钟）
    CacheExpireShort    = 5 * time.Minute
    
    // 中期缓存（1小时）
    CacheExpireMedium   = 1 * time.Hour
    
    // 长期缓存（24小时）
    CacheExpireLong     = 24 * time.Hour
    
    // 超长期缓存（7天）
    CacheExpireVeryLong = 7 * 24 * time.Hour
    
    // 临时缓存（10分钟）
    CacheExpireTemp     = 10 * time.Minute
)

// 具体业务的过期时间
var CacheExpireMap = map[string]time.Duration{
    "user":           CacheExpireMedium,    // 用户信息1小时
    "movie":          CacheExpireLong,      // 电影信息24小时
    "movie_list":     CacheExpireShort,     // 电影列表5分钟
    "search":         CacheExpireShort,     // 搜索结果5分钟
    "comments":       CacheExpireShort,     // 评论列表5分钟
    "rating":         CacheExpireMedium,    // 评分信息1小时
    "config":         CacheExpireVeryLong,  // 配置信息7天
    "temp":           CacheExpireTemp,      // 临时数据10分钟
}
```

### 2. **缓存模式实现**

#### 2.1 Cache-Aside模式
```go
// Cache-Aside模式实现
type CacheAsideService struct {
    cache  *redis.Client
    db     *sql.DB
}

// 读取数据
func (s *CacheAsideService) GetMovie(movieID int) (*Movie, error) {
    // 1. 先从缓存读取
    cacheKey := fmt.Sprintf(MovieCacheKey, movieID)
    cached, err := s.cache.Get(context.Background(), cacheKey).Result()
    if err == nil {
        // 缓存命中，反序列化返回
        var movie Movie
        if err := json.Unmarshal([]byte(cached), &movie); err == nil {
            return &movie, nil
        }
    }
    
    // 2. 缓存未命中，从数据库读取
    movie, err := s.getMovieFromDB(movieID)
    if err != nil {
        return nil, err
    }
    
    // 3. 写入缓存
    movieJSON, _ := json.Marshal(movie)
    s.cache.Set(context.Background(), cacheKey, movieJSON, CacheExpireMap["movie"])
    
    return movie, nil
}

// 更新数据
func (s *CacheAsideService) UpdateMovie(movie *Movie) error {
    // 1. 更新数据库
    err := s.updateMovieInDB(movie)
    if err != nil {
        return err
    }
    
    // 2. 删除缓存（让下次读取时重新加载）
    cacheKey := fmt.Sprintf(MovieCacheKey, movie.ID)
    s.cache.Del(context.Background(), cacheKey)
    
    return nil
}
```

#### 2.2 Write-Through模式
```go
// Write-Through模式实现
func (s *CacheService) SetMovieRating(movieID int, rating *Rating) error {
    // 1. 同时写入数据库和缓存
    err := s.setRatingInDB(rating)
    if err != nil {
        return err
    }
    
    // 2. 更新缓存中的评分统计
    ratingKey := fmt.Sprintf(RatingStatsKey, movieID)
    stats, err := s.calculateRatingStats(movieID)
    if err != nil {
        return err
    }
    
    statsJSON, _ := json.Marshal(stats)
    s.cache.Set(context.Background(), ratingKey, statsJSON, CacheExpireMap["rating"])
    
    return nil
}
```

#### 2.3 Write-Behind模式
```go
// Write-Behind模式实现（异步写入）
type WriteBehindCache struct {
    cache     *redis.Client
    db        *sql.DB
    writeQueue chan WriteOperation
}

type WriteOperation struct {
    Type   string
    Key    string
    Data   interface{}
    Timestamp time.Time
}

func (w *WriteBehindCache) IncrementViewCount(movieID int) error {
    // 1. 立即更新缓存
    viewKey := fmt.Sprintf(StatsViewKey, movieID)
    w.cache.Incr(context.Background(), viewKey)
    
    // 2. 异步写入数据库
    w.writeQueue <- WriteOperation{
        Type: "view_count",
        Key:  viewKey,
        Data: movieID,
        Timestamp: time.Now(),
    }
    
    return nil
}

// 后台写入协程
func (w *WriteBehindCache) backgroundWriter() {
    ticker := time.NewTicker(30 * time.Second) // 每30秒批量写入
    defer ticker.Stop()
    
    batch := make([]WriteOperation, 0, 100)
    
    for {
        select {
        case op := <-w.writeQueue:
            batch = append(batch, op)
            if len(batch) >= 100 {
                w.flushBatch(batch)
                batch = batch[:0]
            }
            
        case <-ticker.C:
            if len(batch) > 0 {
                w.flushBatch(batch)
                batch = batch[:0]
            }
        }
    }
}
```

### 3. **缓存预热策略**

#### 3.1 应用启动时预热
```go
// 缓存预热服务
type CacheWarmupService struct {
    cache *redis.Client
    db    *sql.DB
}

func (s *CacheWarmupService) WarmupOnStartup() error {
    log.Println("Starting cache warmup...")
    
    // 1. 预热热门电影
    if err := s.warmupHotMovies(); err != nil {
        log.Printf("Failed to warmup hot movies: %v", err)
    }
    
    // 2. 预热电影分类
    if err := s.warmupCategories(); err != nil {
        log.Printf("Failed to warmup categories: %v", err)
    }
    
    // 3. 预热系统配置
    if err := s.warmupConfigs(); err != nil {
        log.Printf("Failed to warmup configs: %v", err)
    }
    
    log.Println("Cache warmup completed")
    return nil
}

func (s *CacheWarmupService) warmupHotMovies() error {
    // 查询热门电影
    movies, err := s.getHotMoviesFromDB(50) // 获取前50部热门电影
    if err != nil {
        return err
    }
    
    // 批量写入缓存
    pipe := s.cache.Pipeline()
    for _, movie := range movies {
        movieJSON, _ := json.Marshal(movie)
        cacheKey := fmt.Sprintf(MovieCacheKey, movie.ID)
        pipe.Set(context.Background(), cacheKey, movieJSON, CacheExpireMap["movie"])
    }
    
    // 缓存热门电影列表
    hotListJSON, _ := json.Marshal(movies)
    pipe.Set(context.Background(), MovieHotKey, hotListJSON, CacheExpireMap["movie_list"])
    
    _, err = pipe.Exec(context.Background())
    return err
}
```

#### 3.2 定时预热任务
```go
// 定时预热任务
func (s *CacheWarmupService) StartScheduledWarmup() {
    // 每小时更新热门电影缓存
    hotMoviesTicker := time.NewTicker(1 * time.Hour)
    go func() {
        for range hotMoviesTicker.C {
            s.warmupHotMovies()
        }
    }()
    
    // 每天更新电影统计缓存
    statsTicker := time.NewTicker(24 * time.Hour)
    go func() {
        for range statsTicker.C {
            s.warmupMovieStats()
        }
    }()
}

func (s *CacheWarmupService) warmupMovieStats() error {
    // 获取需要预热的电影ID列表
    movieIDs, err := s.getPopularMovieIDs(100)
    if err != nil {
        return err
    }
    
    // 并发预热评分统计
    semaphore := make(chan struct{}, 10) // 限制并发数
    var wg sync.WaitGroup
    
    for _, movieID := range movieIDs {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            s.warmupMovieRatingStats(id)
        }(movieID)
    }
    
    wg.Wait()
    return nil
}
```

## 缓存监控和维护

### 1. **性能监控**

#### 1.1 监控指标收集
```go
// Redis监控服务
type RedisMonitor struct {
    client *redis.Client
}

type RedisMetrics struct {
    ConnectedClients    int64   `json:"connected_clients"`
    UsedMemory         int64   `json:"used_memory"`
    UsedMemoryPeak     int64   `json:"used_memory_peak"`
    UsedMemoryRss      int64   `json:"used_memory_rss"`
    MemoryFragmentation float64 `json:"memory_fragmentation_ratio"`
    KeyspaceHits       int64   `json:"keyspace_hits"`
    KeyspaceMisses     int64   `json:"keyspace_misses"`
    HitRate            float64 `json:"hit_rate"`
    EvictedKeys        int64   `json:"evicted_keys"`
    ExpiredKeys        int64   `json:"expired_keys"`
    TotalCommands      int64   `json:"total_commands_processed"`
    OpsPerSecond       float64 `json:"instantaneous_ops_per_sec"`
}

func (m *RedisMonitor) GetMetrics() (*RedisMetrics, error) {
    info, err := m.client.Info(context.Background(), "memory", "stats", "clients").Result()
    if err != nil {
        return nil, err
    }
    
    metrics := &RedisMetrics{}
    
    // 解析INFO命令输出
    lines := strings.Split(info, "\r\n")
    for _, line := range lines {
        if strings.Contains(line, ":") {
            parts := strings.Split(line, ":")
            if len(parts) != 2 {
                continue
            }
            
            key, value := parts[0], parts[1]
            switch key {
            case "connected_clients":
                metrics.ConnectedClients, _ = strconv.ParseInt(value, 10, 64)
            case "used_memory":
                metrics.UsedMemory, _ = strconv.ParseInt(value, 10, 64)
            case "used_memory_peak":
                metrics.UsedMemoryPeak, _ = strconv.ParseInt(value, 10, 64)
            case "used_memory_rss":
                metrics.UsedMemoryRss, _ = strconv.ParseInt(value, 10, 64)
            case "mem_fragmentation_ratio":
                metrics.MemoryFragmentation, _ = strconv.ParseFloat(value, 64)
            case "keyspace_hits":
                metrics.KeyspaceHits, _ = strconv.ParseInt(value, 10, 64)
            case "keyspace_misses":
                metrics.KeyspaceMisses, _ = strconv.ParseInt(value, 10, 64)
            case "evicted_keys":
                metrics.EvictedKeys, _ = strconv.ParseInt(value, 10, 64)
            case "expired_keys":
                metrics.ExpiredKeys, _ = strconv.ParseInt(value, 10, 64)
            case "total_commands_processed":
                metrics.TotalCommands, _ = strconv.ParseInt(value, 10, 64)
            case "instantaneous_ops_per_sec":
                metrics.OpsPerSecond, _ = strconv.ParseFloat(value, 64)
            }
        }
    }
    
    // 计算命中率
    if metrics.KeyspaceHits+metrics.KeyspaceMisses > 0 {
        metrics.HitRate = float64(metrics.KeyspaceHits) / float64(metrics.KeyspaceHits+metrics.KeyspaceMisses) * 100
    }
    
    return metrics, nil
}
```

#### 1.2 告警机制
```go
// 告警配置
type AlertConfig struct {
    MemoryUsageThreshold    float64 // 内存使用率阈值
    HitRateThreshold       float64 // 命中率阈值
    ConnectionThreshold    int64   // 连接数阈值
    FragmentationThreshold float64 // 内存碎片率阈值
}

func (m *RedisMonitor) CheckAlerts(config *AlertConfig) []string {
    metrics, err := m.GetMetrics()
    if err != nil {
        return []string{fmt.Sprintf("Failed to get Redis metrics: %v", err)}
    }
    
    var alerts []string
    
    // 检查内存使用率
    memoryUsage := float64(metrics.UsedMemory) / float64(4*1024*1024*1024) * 100 // 假设最大4GB
    if memoryUsage > config.MemoryUsageThreshold {
        alerts = append(alerts, fmt.Sprintf("High memory usage: %.2f%%", memoryUsage))
    }
    
    // 检查命中率
    if metrics.HitRate < config.HitRateThreshold {
        alerts = append(alerts, fmt.Sprintf("Low hit rate: %.2f%%", metrics.HitRate))
    }
    
    // 检查连接数
    if metrics.ConnectedClients > config.ConnectionThreshold {
        alerts = append(alerts, fmt.Sprintf("High connection count: %d", metrics.ConnectedClients))
    }
    
    // 检查内存碎片率
    if metrics.MemoryFragmentation > config.FragmentationThreshold {
        alerts = append(alerts, fmt.Sprintf("High memory fragmentation: %.2f", metrics.MemoryFragmentation))
    }
    
    return alerts
}
```

### 2. **缓存维护**

#### 2.1 过期键清理
```go
// 过期键清理服务
type CacheCleanupService struct {
    client *redis.Client
}

func (s *CacheCleanupService) CleanupExpiredKeys() error {
    // 获取所有数据库
    for db := 0; db < 16; db++ {
        err := s.cleanupDatabase(db)
        if err != nil {
            log.Printf("Failed to cleanup database %d: %v", db, err)
        }
    }
    return nil
}

func (s *CacheCleanupService) cleanupDatabase(db int) error {
    // 切换到指定数据库
    s.client.Do(context.Background(), "SELECT", db)
    
    // 扫描过期键模式
    patterns := []string{
        "temp:*",      // 临时数据
        "session:*",   // 会话数据
        "search:*",    // 搜索结果
    }
    
    for _, pattern := range patterns {
        s.cleanupPattern(pattern)
    }
    
    return nil
}

func (s *CacheCleanupService) cleanupPattern(pattern string) {
    var cursor uint64
    for {
        keys, nextCursor, err := s.client.Scan(context.Background(), cursor, pattern, 100).Result()
        if err != nil {
            break
        }
        
        // 检查键是否过期
        for _, key := range keys {
            ttl := s.client.TTL(context.Background(), key).Val()
            if ttl == -1 { // 没有设置过期时间的键
                // 根据键类型设置合适的过期时间
                s.setDefaultExpire(key)
            }
        }
        
        cursor = nextCursor
        if cursor == 0 {
            break
        }
    }
}
```

#### 2.2 内存优化
```go
// 内存优化服务
func (s *CacheCleanupService) OptimizeMemory() error {
    // 1. 内存碎片整理
    result := s.client.Do(context.Background(), "MEMORY", "PURGE")
    log.Printf("Memory purge result: %v", result)
    
    // 2. 清理大键
    s.cleanupLargeKeys()
    
    // 3. 优化数据结构
    s.optimizeDataStructures()
    
    return nil
}

func (s *CacheCleanupService) cleanupLargeKeys() {
    // 查找大键
    var cursor uint64
    for {
        keys, nextCursor, err := s.client.Scan(context.Background(), cursor, "*", 1000).Result()
        if err != nil {
            break
        }
        
        for _, key := range keys {
            // 检查键的内存使用
            memUsage := s.client.MemoryUsage(context.Background(), key).Val()
            if memUsage > 1024*1024 { // 大于1MB的键
                log.Printf("Large key found: %s, size: %d bytes", key, memUsage)
                
                // 检查是否可以优化
                s.optimizeKey(key)
            }
        }
        
        cursor = nextCursor
        if cursor == 0 {
            break
        }
    }
}
```

## 总结

Redis缓存配置为MovieInfo项目提供了高性能的缓存解决方案。通过合理的配置优化、缓存策略设计和监控维护，我们建立了一个稳定、高效的缓存系统。

**关键配置要点**：
1. **内存管理**：合理的内存策略和淘汰机制
2. **持久化配置**：平衡性能和数据安全
3. **缓存策略**：多种缓存模式的合理应用
4. **监控维护**：完善的监控和自动化维护

**缓存优势**：
- **高性能**：微秒级响应时间
- **高可用**：完善的故障处理机制
- **易扩展**：支持集群和分片
- **智能化**：自动化的监控和维护

**下一步**：基于这个Redis配置，我们将进行IDE配置与插件安装，为开发人员提供更好的开发体验。
