# 4.1 数据库连接池

## 概述

数据库连接池是现代应用程序与数据库交互的核心组件，它管理着数据库连接的创建、复用和销毁。对于MovieInfo项目，我们需要实现一个高效、稳定、可配置的数据库连接池，支持连接复用、健康检查和性能监控。

## 为什么需要数据库连接池？

### 1. **性能优化**
- **连接复用**：避免频繁创建和销毁数据库连接的开销
- **并发处理**：支持多个并发请求同时访问数据库
- **资源管理**：合理分配和管理数据库连接资源
- **响应时间**：减少连接建立时间，提升响应速度

### 2. **资源控制**
- **连接限制**：防止过多连接导致数据库压力过大
- **内存管理**：控制连接池占用的内存资源
- **超时管理**：自动清理长时间空闲的连接
- **负载均衡**：在多个数据库实例间分配连接

### 3. **稳定性保障**
- **故障恢复**：自动检测和恢复失效的连接
- **健康检查**：定期检查连接的健康状态
- **优雅降级**：在连接不足时的降级策略
- **错误隔离**：防止单个连接错误影响整个应用

### 4. **监控和调试**
- **连接统计**：实时监控连接池的使用情况
- **性能指标**：收集连接池的性能数据
- **问题诊断**：提供连接问题的诊断信息
- **容量规划**：基于使用数据进行容量规划

## 连接池架构设计

### 1. **连接池组件结构**

```
数据库连接池
├── 连接管理器 (Connection Manager)
│   ├── 连接创建 (Connection Factory)
│   ├── 连接验证 (Connection Validator)
│   ├── 连接回收 (Connection Recycler)
│   └── 连接监控 (Connection Monitor)
├── 连接池 (Connection Pool)
│   ├── 活跃连接 (Active Connections)
│   ├── 空闲连接 (Idle Connections)
│   ├── 等待队列 (Waiting Queue)
│   └── 连接缓存 (Connection Cache)
├── 配置管理 (Configuration)
│   ├── 池大小配置 (Pool Size Config)
│   ├── 超时配置 (Timeout Config)
│   ├── 健康检查配置 (Health Check Config)
│   └── 重试配置 (Retry Config)
└── 监控统计 (Metrics)
    ├── 连接统计 (Connection Stats)
    ├── 性能指标 (Performance Metrics)
    ├── 错误统计 (Error Stats)
    └── 告警机制 (Alerting)
```

### 2. **连接池配置结构**

#### 2.1 配置定义
```go
// pkg/database/config.go
package database

import (
    "time"
)

// Config 数据库配置
type Config struct {
    // 连接配置
    Host     string `mapstructure:"host" yaml:"host"`
    Port     int    `mapstructure:"port" yaml:"port"`
    Username string `mapstructure:"username" yaml:"username"`
    Password string `mapstructure:"password" yaml:"password"`
    DBName   string `mapstructure:"dbname" yaml:"dbname"`
    Charset  string `mapstructure:"charset" yaml:"charset"`
    
    // 连接池配置
    MaxOpenConns    int           `mapstructure:"max_open_conns" yaml:"max_open_conns"`       // 最大打开连接数
    MaxIdleConns    int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`       // 最大空闲连接数
    ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"` // 连接最大生存时间
    ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" yaml:"conn_max_idle_time"` // 连接最大空闲时间
    
    // 超时配置
    ConnectTimeout time.Duration `mapstructure:"connect_timeout" yaml:"connect_timeout"` // 连接超时
    ReadTimeout    time.Duration `mapstructure:"read_timeout" yaml:"read_timeout"`       // 读取超时
    WriteTimeout   time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`     // 写入超时
    
    // 健康检查配置
    HealthCheckInterval time.Duration `mapstructure:"health_check_interval" yaml:"health_check_interval"` // 健康检查间隔
    HealthCheckTimeout  time.Duration `mapstructure:"health_check_timeout" yaml:"health_check_timeout"`   // 健康检查超时
    
    // 重试配置
    MaxRetries    int           `mapstructure:"max_retries" yaml:"max_retries"`       // 最大重试次数
    RetryInterval time.Duration `mapstructure:"retry_interval" yaml:"retry_interval"` // 重试间隔
    
    // SSL配置
    SSLMode     string `mapstructure:"ssl_mode" yaml:"ssl_mode"`         // SSL模式
    SSLCert     string `mapstructure:"ssl_cert" yaml:"ssl_cert"`         // SSL证书
    SSLKey      string `mapstructure:"ssl_key" yaml:"ssl_key"`           // SSL密钥
    SSLRootCert string `mapstructure:"ssl_root_cert" yaml:"ssl_root_cert"` // SSL根证书
}

// DefaultConfig 默认配置
func DefaultConfig() Config {
    return Config{
        Host:     "localhost",
        Port:     3306,
        Charset:  "utf8mb4",
        
        MaxOpenConns:    100,
        MaxIdleConns:    10,
        ConnMaxLifetime: time.Hour,
        ConnMaxIdleTime: time.Minute * 30,
        
        ConnectTimeout: time.Second * 10,
        ReadTimeout:    time.Second * 30,
        WriteTimeout:   time.Second * 30,
        
        HealthCheckInterval: time.Minute * 5,
        HealthCheckTimeout:  time.Second * 5,
        
        MaxRetries:    3,
        RetryInterval: time.Second * 2,
        
        SSLMode: "disable",
    }
}

// DSN 生成数据源名称
func (c Config) DSN() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
        c.Username, c.Password, c.Host, c.Port, c.DBName, c.Charset)
}

// Validate 验证配置
func (c Config) Validate() error {
    if c.Host == "" {
        return fmt.Errorf("database host is required")
    }
    
    if c.Port <= 0 || c.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", c.Port)
    }
    
    if c.Username == "" {
        return fmt.Errorf("database username is required")
    }
    
    if c.DBName == "" {
        return fmt.Errorf("database name is required")
    }
    
    if c.MaxOpenConns <= 0 {
        return fmt.Errorf("max_open_conns must be greater than 0")
    }
    
    if c.MaxIdleConns < 0 {
        return fmt.Errorf("max_idle_conns must be non-negative")
    }
    
    if c.MaxIdleConns > c.MaxOpenConns {
        return fmt.Errorf("max_idle_conns cannot be greater than max_open_conns")
    }
    
    return nil
}
```

### 3. **连接池实现**

#### 3.1 连接池结构
```go
// pkg/database/pool.go
package database

import (
    "context"
    "database/sql"
    "sync"
    "sync/atomic"
    "time"
    
    _ "github.com/go-sql-driver/mysql"
    "github.com/yourname/movieinfo/pkg/logger"
)

// Pool 数据库连接池
type Pool struct {
    db     *sql.DB
    config Config
    logger logger.Logger
    
    // 统计信息
    stats    *Stats
    statsLock sync.RWMutex
    
    // 健康检查
    healthChecker *HealthChecker
    
    // 监控
    monitor *Monitor
    
    // 状态
    closed int32
}

// Stats 连接池统计信息
type Stats struct {
    OpenConnections     int64     // 当前打开的连接数
    InUseConnections    int64     // 当前使用中的连接数
    IdleConnections     int64     // 当前空闲的连接数
    WaitCount           int64     // 等待连接的请求数
    WaitDuration        time.Duration // 平均等待时间
    MaxOpenConnections  int64     // 最大打开连接数
    MaxIdleConnections  int64     // 最大空闲连接数
    MaxLifetime         time.Duration // 连接最大生存时间
    MaxIdleTime         time.Duration // 连接最大空闲时间
    
    // 累计统计
    TotalConnections    int64     // 累计创建的连接数
    TotalQueries        int64     // 累计查询数
    TotalErrors         int64     // 累计错误数
    TotalRetries        int64     // 累计重试数
    
    // 时间统计
    LastActivity        time.Time // 最后活动时间
    StartTime           time.Time // 启动时间
}

// NewPool 创建数据库连接池
func NewPool(config Config) (*Pool, error) {
    // 验证配置
    if err := config.Validate(); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    
    // 创建数据库连接
    db, err := sql.Open("mysql", config.DSN())
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // 配置连接池
    db.SetMaxOpenConns(config.MaxOpenConns)
    db.SetMaxIdleConns(config.MaxIdleConns)
    db.SetConnMaxLifetime(config.ConnMaxLifetime)
    db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
    
    // 测试连接
    ctx, cancel := context.WithTimeout(context.Background(), config.ConnectTimeout)
    defer cancel()
    
    if err := db.PingContext(ctx); err != nil {
        db.Close()
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    pool := &Pool{
        db:     db,
        config: config,
        logger: logger.GetGlobalLogger(),
        stats: &Stats{
            MaxOpenConnections: int64(config.MaxOpenConns),
            MaxIdleConnections: int64(config.MaxIdleConns),
            MaxLifetime:        config.ConnMaxLifetime,
            MaxIdleTime:        config.ConnMaxIdleTime,
            StartTime:          time.Now(),
            LastActivity:       time.Now(),
        },
    }
    
    // 启动健康检查
    pool.healthChecker = NewHealthChecker(pool)
    pool.healthChecker.Start()
    
    // 启动监控
    pool.monitor = NewMonitor(pool)
    pool.monitor.Start()
    
    pool.logger.Info("Database connection pool created",
        logger.String("host", config.Host),
        logger.Int("port", config.Port),
        logger.String("database", config.DBName),
        logger.Int("max_open_conns", config.MaxOpenConns),
        logger.Int("max_idle_conns", config.MaxIdleConns),
    )
    
    return pool, nil
}

// GetDB 获取数据库连接
func (p *Pool) GetDB() *sql.DB {
    return p.db
}

// Ping 检查数据库连接
func (p *Pool) Ping(ctx context.Context) error {
    if atomic.LoadInt32(&p.closed) == 1 {
        return fmt.Errorf("connection pool is closed")
    }
    
    start := time.Now()
    err := p.db.PingContext(ctx)
    duration := time.Since(start)
    
    // 更新统计
    p.updateStats(func(stats *Stats) {
        stats.LastActivity = time.Now()
        if err != nil {
            stats.TotalErrors++
        }
    })
    
    p.logger.Debug("Database ping",
        logger.Duration("duration", duration),
        logger.Error(err),
    )
    
    return err
}

// Query 执行查询
func (p *Pool) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    if atomic.LoadInt32(&p.closed) == 1 {
        return nil, fmt.Errorf("connection pool is closed")
    }
    
    start := time.Now()
    rows, err := p.db.QueryContext(ctx, query, args...)
    duration := time.Since(start)
    
    // 更新统计
    p.updateStats(func(stats *Stats) {
        stats.TotalQueries++
        stats.LastActivity = time.Now()
        if err != nil {
            stats.TotalErrors++
        }
    })
    
    p.logger.Debug("Database query executed",
        logger.String("query", query),
        logger.Duration("duration", duration),
        logger.Error(err),
    )
    
    return rows, err
}

// QueryRow 执行单行查询
func (p *Pool) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
    if atomic.LoadInt32(&p.closed) == 1 {
        return nil
    }
    
    start := time.Now()
    row := p.db.QueryRowContext(ctx, query, args...)
    duration := time.Since(start)
    
    // 更新统计
    p.updateStats(func(stats *Stats) {
        stats.TotalQueries++
        stats.LastActivity = time.Now()
    })
    
    p.logger.Debug("Database query row executed",
        logger.String("query", query),
        logger.Duration("duration", duration),
    )
    
    return row
}

// Exec 执行SQL语句
func (p *Pool) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
    if atomic.LoadInt32(&p.closed) == 1 {
        return nil, fmt.Errorf("connection pool is closed")
    }
    
    start := time.Now()
    result, err := p.db.ExecContext(ctx, query, args...)
    duration := time.Since(start)
    
    // 更新统计
    p.updateStats(func(stats *Stats) {
        stats.TotalQueries++
        stats.LastActivity = time.Now()
        if err != nil {
            stats.TotalErrors++
        }
    })
    
    p.logger.Debug("Database exec executed",
        logger.String("query", query),
        logger.Duration("duration", duration),
        logger.Error(err),
    )
    
    return result, err
}

// Begin 开始事务
func (p *Pool) Begin(ctx context.Context) (*sql.Tx, error) {
    if atomic.LoadInt32(&p.closed) == 1 {
        return nil, fmt.Errorf("connection pool is closed")
    }
    
    start := time.Now()
    tx, err := p.db.BeginTx(ctx, nil)
    duration := time.Since(start)
    
    // 更新统计
    p.updateStats(func(stats *Stats) {
        stats.LastActivity = time.Now()
        if err != nil {
            stats.TotalErrors++
        }
    })
    
    p.logger.Debug("Database transaction started",
        logger.Duration("duration", duration),
        logger.Error(err),
    )
    
    return tx, err
}

// GetStats 获取连接池统计信息
func (p *Pool) GetStats() Stats {
    p.statsLock.RLock()
    defer p.statsLock.RUnlock()
    
    // 获取数据库统计
    dbStats := p.db.Stats()
    
    stats := *p.stats
    stats.OpenConnections = int64(dbStats.OpenConnections)
    stats.InUseConnections = int64(dbStats.InUse)
    stats.IdleConnections = int64(dbStats.Idle)
    stats.WaitCount = dbStats.WaitCount
    stats.WaitDuration = dbStats.WaitDuration
    
    return stats
}

// updateStats 更新统计信息
func (p *Pool) updateStats(fn func(*Stats)) {
    p.statsLock.Lock()
    defer p.statsLock.Unlock()
    fn(p.stats)
}

// Close 关闭连接池
func (p *Pool) Close() error {
    if !atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
        return fmt.Errorf("connection pool already closed")
    }
    
    // 停止健康检查
    if p.healthChecker != nil {
        p.healthChecker.Stop()
    }
    
    // 停止监控
    if p.monitor != nil {
        p.monitor.Stop()
    }
    
    // 关闭数据库连接
    err := p.db.Close()
    
    p.logger.Info("Database connection pool closed")
    
    return err
}

// IsClosed 检查连接池是否已关闭
func (p *Pool) IsClosed() bool {
    return atomic.LoadInt32(&p.closed) == 1
}
```

### 4. **健康检查机制**

#### 4.1 健康检查器
```go
// pkg/database/health_checker.go
package database

import (
    "context"
    "sync"
    "time"
    
    "github.com/yourname/movieinfo/pkg/logger"
)

// HealthChecker 健康检查器
type HealthChecker struct {
    pool   *Pool
    logger logger.Logger
    
    ticker *time.Ticker
    done   chan struct{}
    wg     sync.WaitGroup
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(pool *Pool) *HealthChecker {
    return &HealthChecker{
        pool:   pool,
        logger: logger.GetGlobalLogger(),
        done:   make(chan struct{}),
    }
}

// Start 启动健康检查
func (h *HealthChecker) Start() {
    h.ticker = time.NewTicker(h.pool.config.HealthCheckInterval)
    
    h.wg.Add(1)
    go h.run()
    
    h.logger.Info("Database health checker started",
        logger.Duration("interval", h.pool.config.HealthCheckInterval),
    )
}

// Stop 停止健康检查
func (h *HealthChecker) Stop() {
    if h.ticker != nil {
        h.ticker.Stop()
    }
    
    close(h.done)
    h.wg.Wait()
    
    h.logger.Info("Database health checker stopped")
}

// run 运行健康检查
func (h *HealthChecker) run() {
    defer h.wg.Done()
    
    for {
        select {
        case <-h.ticker.C:
            h.checkHealth()
        case <-h.done:
            return
        }
    }
}

// checkHealth 执行健康检查
func (h *HealthChecker) checkHealth() {
    ctx, cancel := context.WithTimeout(context.Background(), h.pool.config.HealthCheckTimeout)
    defer cancel()
    
    start := time.Now()
    err := h.pool.Ping(ctx)
    duration := time.Since(start)
    
    if err != nil {
        h.logger.Error("Database health check failed",
            logger.Duration("duration", duration),
            logger.Error(err),
        )
        
        // 触发告警
        h.triggerAlert("health_check_failed", err)
    } else {
        h.logger.Debug("Database health check passed",
            logger.Duration("duration", duration),
        )
    }
    
    // 检查连接池统计
    h.checkPoolStats()
}

// checkPoolStats 检查连接池统计
func (h *HealthChecker) checkPoolStats() {
    stats := h.pool.GetStats()
    
    // 检查连接使用率
    if stats.OpenConnections > 0 {
        usageRate := float64(stats.InUseConnections) / float64(stats.OpenConnections)
        if usageRate > 0.9 { // 使用率超过90%
            h.logger.Warn("High database connection usage",
                logger.Float64("usage_rate", usageRate),
                logger.Int64("in_use", stats.InUseConnections),
                logger.Int64("open", stats.OpenConnections),
            )
            
            h.triggerAlert("high_connection_usage", nil)
        }
    }
    
    // 检查等待时间
    if stats.WaitDuration > time.Second*5 { // 等待时间超过5秒
        h.logger.Warn("High database connection wait time",
            logger.Duration("wait_duration", stats.WaitDuration),
            logger.Int64("wait_count", stats.WaitCount),
        )
        
        h.triggerAlert("high_wait_time", nil)
    }
    
    // 检查错误率
    if stats.TotalQueries > 0 {
        errorRate := float64(stats.TotalErrors) / float64(stats.TotalQueries)
        if errorRate > 0.1 { // 错误率超过10%
            h.logger.Warn("High database error rate",
                logger.Float64("error_rate", errorRate),
                logger.Int64("total_errors", stats.TotalErrors),
                logger.Int64("total_queries", stats.TotalQueries),
            )
            
            h.triggerAlert("high_error_rate", nil)
        }
    }
}

// triggerAlert 触发告警
func (h *HealthChecker) triggerAlert(alertType string, err error) {
    // 这里可以集成告警系统，如发送邮件、短信、钉钉等
    h.logger.Error("Database alert triggered",
        logger.String("alert_type", alertType),
        logger.Error(err),
    )
    
    // TODO: 集成具体的告警系统
}
```

## 总结

数据库连接池为MovieInfo项目提供了高效、稳定的数据库访问基础设施。通过合理的连接管理、健康检查和性能监控，我们建立了一个生产级的数据库连接解决方案。

**关键设计要点**：
1. **连接复用**：高效的连接池管理和复用机制
2. **配置灵活**：丰富的配置选项和默认值
3. **健康监控**：完善的健康检查和监控机制
4. **错误处理**：优雅的错误处理和恢复策略
5. **性能优化**：连接池参数的性能优化

**连接池优势**：
- **高性能**：连接复用减少创建开销
- **高可用**：健康检查保证连接可用性
- **可监控**：详细的统计和监控信息
- **可配置**：灵活的配置和调优选项

**下一步**：基于这个连接池基础，我们将定义数据模型，包括用户、电影、评论等核心业务实体的结构定义和关系映射。
