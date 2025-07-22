# 3.3 配置文件设计

## 概述

配置文件是应用程序的重要组成部分，它决定了应用在不同环境下的行为。对于MovieInfo项目，我们需要设计一个灵活、安全、易维护的配置系统，支持多环境部署和动态配置管理。

## 为什么配置文件设计如此重要？

### 1. **环境适配**
- **多环境支持**：开发、测试、生产环境的不同配置
- **部署灵活性**：无需修改代码即可适配不同环境
- **配置隔离**：敏感信息与代码分离
- **动态调整**：运行时配置的动态更新

### 2. **安全性**
- **敏感信息保护**：数据库密码、API密钥等敏感信息
- **权限控制**：不同环境的访问权限控制
- **加密存储**：敏感配置的加密存储
- **审计追踪**：配置变更的审计日志

### 3. **可维护性**
- **集中管理**：统一的配置管理入口
- **版本控制**：配置变更的版本管理
- **文档化**：清晰的配置说明和示例
- **验证机制**：配置有效性的自动验证

### 4. **运维友好**
- **热更新**：无需重启的配置更新
- **监控告警**：配置异常的监控和告警
- **回滚机制**：配置错误的快速回滚
- **批量管理**：多实例的批量配置管理

## 配置系统架构设计

### 1. **配置层次结构**

```
配置系统
├── 默认配置 (Default)     # 内置默认值
├── 文件配置 (File)        # 配置文件
├── 环境变量 (Environment) # 环境变量
└── 命令行参数 (CLI)       # 命令行参数
```

**优先级顺序**：命令行参数 > 环境变量 > 配置文件 > 默认配置

### 2. **配置文件结构**

#### 2.1 主配置文件 (config.yaml)
```yaml
# MovieInfo 配置文件
app:
  name: "MovieInfo"
  version: "1.0.0"
  environment: "development"  # development, testing, production
  debug: true
  timezone: "Asia/Shanghai"

# 服务配置
server:
  web:
    host: "0.0.0.0"
    port: 8080
    read_timeout: 30s
    write_timeout: 30s
    idle_timeout: 60s
    max_header_bytes: 1048576
  user:
    http_port: 8081
    grpc_port: 9081
  movie:
    http_port: 8082
    grpc_port: 9082
  comment:
    http_port: 8083
    grpc_port: 9083

# 数据库配置
database:
  mysql:
    host: "localhost"
    port: 3306
    username: "movieinfo"
    password: "movieinfo123"
    dbname: "movieinfo"
    charset: "utf8mb4"
    parse_time: true
    loc: "Local"
    max_open_conns: 100
    max_idle_conns: 10
    conn_max_lifetime: 3600s
    conn_max_idle_time: 300s
    
# Redis配置
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10
  min_idle_conns: 5
  dial_timeout: 5s
  read_timeout: 3s
  write_timeout: 3s
  pool_timeout: 4s
  idle_timeout: 300s

# JWT配置
jwt:
  secret: "your_jwt_secret_key_change_in_production"
  expire_hours: 24
  refresh_expire_hours: 168  # 7 days
  issuer: "movieinfo"

# 日志配置
log:
  level: "info"  # debug, info, warn, error
  format: "json"  # json, text
  output: "stdout"  # stdout, stderr, file
  file_path: "logs/app.log"
  max_size: 100  # MB
  max_backups: 10
  max_age: 30  # days
  compress: true

# 缓存配置
cache:
  default_expire: 3600s  # 1 hour
  cleanup_interval: 600s  # 10 minutes
  max_size: 1000000  # 1M items
  
  # 具体业务缓存配置
  user:
    expire: 1800s  # 30 minutes
  movie:
    expire: 86400s  # 24 hours
  search:
    expire: 300s   # 5 minutes
  hot_movies:
    expire: 3600s  # 1 hour

# 文件上传配置
upload:
  max_size: 10485760  # 10MB
  allowed_types: ["jpg", "jpeg", "png", "gif"]
  upload_path: "uploads"
  url_prefix: "/static/uploads"

# 邮件配置
email:
  smtp_host: "smtp.gmail.com"
  smtp_port: 587
  username: "your-email@gmail.com"
  password: "your-app-password"
  from_name: "MovieInfo"
  from_email: "noreply@movieinfo.com"

# 第三方服务配置
external:
  tmdb:
    api_key: "your_tmdb_api_key"
    base_url: "https://api.themoviedb.org/3"
    timeout: 10s
  
  imdb:
    api_key: "your_imdb_api_key"
    base_url: "https://imdb-api.com"
    timeout: 10s

# 安全配置
security:
  cors:
    allowed_origins: ["http://localhost:3000", "https://movieinfo.com"]
    allowed_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allowed_headers: ["Origin", "Content-Type", "Authorization"]
    expose_headers: ["Content-Length"]
    allow_credentials: true
    max_age: 86400
  
  rate_limit:
    enabled: true
    requests_per_minute: 100
    burst: 200
    
  csrf:
    enabled: true
    secret: "csrf_secret_key"
    
# 监控配置
monitoring:
  metrics:
    enabled: true
    path: "/metrics"
    port: 9090
  
  health:
    enabled: true
    path: "/health"
  
  pprof:
    enabled: false  # 仅在开发环境启用
    path: "/debug/pprof"

# 功能开关
features:
  user_registration: true
  email_verification: true
  social_login: false
  movie_recommendation: true
  comment_moderation: true
  search_suggestion: true
```

#### 2.2 环境特定配置

**开发环境 (config.development.yaml)**:
```yaml
app:
  debug: true
  environment: "development"

database:
  mysql:
    host: "localhost"
    dbname: "movieinfo_dev"

log:
  level: "debug"
  format: "text"

monitoring:
  pprof:
    enabled: true

features:
  email_verification: false
```

**生产环境 (config.production.yaml)**:
```yaml
app:
  debug: false
  environment: "production"

server:
  web:
    read_timeout: 10s
    write_timeout: 10s

log:
  level: "info"
  format: "json"
  output: "file"

security:
  rate_limit:
    requests_per_minute: 1000
    
monitoring:
  pprof:
    enabled: false
```

### 3. **配置结构定义**

#### 3.1 Go配置结构体
```go
// internal/config/config.go
package config

import (
    "time"
)

// Config 应用配置结构
type Config struct {
    App        AppConfig        `mapstructure:"app" yaml:"app"`
    Server     ServerConfig     `mapstructure:"server" yaml:"server"`
    Database   DatabaseConfig   `mapstructure:"database" yaml:"database"`
    Redis      RedisConfig      `mapstructure:"redis" yaml:"redis"`
    JWT        JWTConfig        `mapstructure:"jwt" yaml:"jwt"`
    Log        LogConfig        `mapstructure:"log" yaml:"log"`
    Cache      CacheConfig      `mapstructure:"cache" yaml:"cache"`
    Upload     UploadConfig     `mapstructure:"upload" yaml:"upload"`
    Email      EmailConfig      `mapstructure:"email" yaml:"email"`
    External   ExternalConfig   `mapstructure:"external" yaml:"external"`
    Security   SecurityConfig   `mapstructure:"security" yaml:"security"`
    Monitoring MonitoringConfig `mapstructure:"monitoring" yaml:"monitoring"`
    Features   FeaturesConfig   `mapstructure:"features" yaml:"features"`
}

// AppConfig 应用基础配置
type AppConfig struct {
    Name        string `mapstructure:"name" yaml:"name"`
    Version     string `mapstructure:"version" yaml:"version"`
    Environment string `mapstructure:"environment" yaml:"environment"`
    Debug       bool   `mapstructure:"debug" yaml:"debug"`
    Timezone    string `mapstructure:"timezone" yaml:"timezone"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
    Web     HTTPServerConfig `mapstructure:"web" yaml:"web"`
    User    ServiceConfig    `mapstructure:"user" yaml:"user"`
    Movie   ServiceConfig    `mapstructure:"movie" yaml:"movie"`
    Comment ServiceConfig    `mapstructure:"comment" yaml:"comment"`
}

// HTTPServerConfig HTTP服务器配置
type HTTPServerConfig struct {
    Host           string        `mapstructure:"host" yaml:"host"`
    Port           int           `mapstructure:"port" yaml:"port"`
    ReadTimeout    time.Duration `mapstructure:"read_timeout" yaml:"read_timeout"`
    WriteTimeout   time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`
    IdleTimeout    time.Duration `mapstructure:"idle_timeout" yaml:"idle_timeout"`
    MaxHeaderBytes int           `mapstructure:"max_header_bytes" yaml:"max_header_bytes"`
}

// ServiceConfig 微服务配置
type ServiceConfig struct {
    HTTPPort int `mapstructure:"http_port" yaml:"http_port"`
    GRPCPort int `mapstructure:"grpc_port" yaml:"grpc_port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
    MySQL MySQLConfig `mapstructure:"mysql" yaml:"mysql"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
    Host              string        `mapstructure:"host" yaml:"host"`
    Port              int           `mapstructure:"port" yaml:"port"`
    Username          string        `mapstructure:"username" yaml:"username"`
    Password          string        `mapstructure:"password" yaml:"password"`
    DBName            string        `mapstructure:"dbname" yaml:"dbname"`
    Charset           string        `mapstructure:"charset" yaml:"charset"`
    ParseTime         bool          `mapstructure:"parse_time" yaml:"parse_time"`
    Loc               string        `mapstructure:"loc" yaml:"loc"`
    MaxOpenConns      int           `mapstructure:"max_open_conns" yaml:"max_open_conns"`
    MaxIdleConns      int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
    ConnMaxLifetime   time.Duration `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"`
    ConnMaxIdleTime   time.Duration `mapstructure:"conn_max_idle_time" yaml:"conn_max_idle_time"`
}

// RedisConfig Redis配置
type RedisConfig struct {
    Host         string        `mapstructure:"host" yaml:"host"`
    Port         int           `mapstructure:"port" yaml:"port"`
    Password     string        `mapstructure:"password" yaml:"password"`
    DB           int           `mapstructure:"db" yaml:"db"`
    PoolSize     int           `mapstructure:"pool_size" yaml:"pool_size"`
    MinIdleConns int           `mapstructure:"min_idle_conns" yaml:"min_idle_conns"`
    DialTimeout  time.Duration `mapstructure:"dial_timeout" yaml:"dial_timeout"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout" yaml:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`
    PoolTimeout  time.Duration `mapstructure:"pool_timeout" yaml:"pool_timeout"`
    IdleTimeout  time.Duration `mapstructure:"idle_timeout" yaml:"idle_timeout"`
}

// JWTConfig JWT配置
type JWTConfig struct {
    Secret             string `mapstructure:"secret" yaml:"secret"`
    ExpireHours        int    `mapstructure:"expire_hours" yaml:"expire_hours"`
    RefreshExpireHours int    `mapstructure:"refresh_expire_hours" yaml:"refresh_expire_hours"`
    Issuer             string `mapstructure:"issuer" yaml:"issuer"`
}

// LogConfig 日志配置
type LogConfig struct {
    Level      string `mapstructure:"level" yaml:"level"`
    Format     string `mapstructure:"format" yaml:"format"`
    Output     string `mapstructure:"output" yaml:"output"`
    FilePath   string `mapstructure:"file_path" yaml:"file_path"`
    MaxSize    int    `mapstructure:"max_size" yaml:"max_size"`
    MaxBackups int    `mapstructure:"max_backups" yaml:"max_backups"`
    MaxAge     int    `mapstructure:"max_age" yaml:"max_age"`
    Compress   bool   `mapstructure:"compress" yaml:"compress"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
    DefaultExpire   time.Duration            `mapstructure:"default_expire" yaml:"default_expire"`
    CleanupInterval time.Duration            `mapstructure:"cleanup_interval" yaml:"cleanup_interval"`
    MaxSize         int                      `mapstructure:"max_size" yaml:"max_size"`
    User            CacheItemConfig          `mapstructure:"user" yaml:"user"`
    Movie           CacheItemConfig          `mapstructure:"movie" yaml:"movie"`
    Search          CacheItemConfig          `mapstructure:"search" yaml:"search"`
    HotMovies       CacheItemConfig          `mapstructure:"hot_movies" yaml:"hot_movies"`
}

// CacheItemConfig 缓存项配置
type CacheItemConfig struct {
    Expire time.Duration `mapstructure:"expire" yaml:"expire"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
    CORS      CORSConfig      `mapstructure:"cors" yaml:"cors"`
    RateLimit RateLimitConfig `mapstructure:"rate_limit" yaml:"rate_limit"`
    CSRF      CSRFConfig      `mapstructure:"csrf" yaml:"csrf"`
}

// CORSConfig CORS配置
type CORSConfig struct {
    AllowedOrigins   []string `mapstructure:"allowed_origins" yaml:"allowed_origins"`
    AllowedMethods   []string `mapstructure:"allowed_methods" yaml:"allowed_methods"`
    AllowedHeaders   []string `mapstructure:"allowed_headers" yaml:"allowed_headers"`
    ExposedHeaders   []string `mapstructure:"exposed_headers" yaml:"exposed_headers"`
    AllowCredentials bool     `mapstructure:"allow_credentials" yaml:"allow_credentials"`
    MaxAge           int      `mapstructure:"max_age" yaml:"max_age"`
}

// FeaturesConfig 功能开关配置
type FeaturesConfig struct {
    UserRegistration    bool `mapstructure:"user_registration" yaml:"user_registration"`
    EmailVerification   bool `mapstructure:"email_verification" yaml:"email_verification"`
    SocialLogin         bool `mapstructure:"social_login" yaml:"social_login"`
    MovieRecommendation bool `mapstructure:"movie_recommendation" yaml:"movie_recommendation"`
    CommentModeration   bool `mapstructure:"comment_moderation" yaml:"comment_moderation"`
    SearchSuggestion    bool `mapstructure:"search_suggestion" yaml:"search_suggestion"`
}
```

### 4. **配置加载实现**

#### 4.1 配置加载器
```go
// internal/config/loader.go
package config

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    
    "github.com/spf13/viper"
)

var (
    globalConfig *Config
)

// Load 加载配置
func Load(configPath string) (*Config, error) {
    v := viper.New()
    
    // 设置配置文件路径和名称
    if configPath != "" {
        v.SetConfigFile(configPath)
    } else {
        v.SetConfigName("config")
        v.SetConfigType("yaml")
        v.AddConfigPath("./configs")
        v.AddConfigPath("../configs")
        v.AddConfigPath("../../configs")
    }
    
    // 设置环境变量前缀
    v.SetEnvPrefix("MOVIEINFO")
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    v.AutomaticEnv()
    
    // 设置默认值
    setDefaults(v)
    
    // 读取配置文件
    if err := v.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, fmt.Errorf("failed to read config file: %w", err)
        }
    }
    
    // 根据环境加载特定配置
    env := v.GetString("app.environment")
    if env != "" {
        envConfigPath := fmt.Sprintf("config.%s.yaml", env)
        if _, err := os.Stat(filepath.Join("configs", envConfigPath)); err == nil {
            envViper := viper.New()
            envViper.SetConfigFile(filepath.Join("configs", envConfigPath))
            if err := envViper.ReadInConfig(); err == nil {
                // 合并环境特定配置
                for key, value := range envViper.AllSettings() {
                    v.Set(key, value)
                }
            }
        }
    }
    
    // 解析配置到结构体
    var config Config
    if err := v.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }
    
    // 验证配置
    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }
    
    globalConfig = &config
    return &config, nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
    return globalConfig
}

// setDefaults 设置默认值
func setDefaults(v *viper.Viper) {
    // 应用默认配置
    v.SetDefault("app.name", "MovieInfo")
    v.SetDefault("app.version", "1.0.0")
    v.SetDefault("app.environment", "development")
    v.SetDefault("app.debug", true)
    v.SetDefault("app.timezone", "UTC")
    
    // 服务器默认配置
    v.SetDefault("server.web.host", "0.0.0.0")
    v.SetDefault("server.web.port", 8080)
    v.SetDefault("server.web.read_timeout", "30s")
    v.SetDefault("server.web.write_timeout", "30s")
    v.SetDefault("server.web.idle_timeout", "60s")
    v.SetDefault("server.web.max_header_bytes", 1048576)
    
    // 数据库默认配置
    v.SetDefault("database.mysql.host", "localhost")
    v.SetDefault("database.mysql.port", 3306)
    v.SetDefault("database.mysql.charset", "utf8mb4")
    v.SetDefault("database.mysql.parse_time", true)
    v.SetDefault("database.mysql.loc", "Local")
    v.SetDefault("database.mysql.max_open_conns", 100)
    v.SetDefault("database.mysql.max_idle_conns", 10)
    v.SetDefault("database.mysql.conn_max_lifetime", "3600s")
    
    // Redis默认配置
    v.SetDefault("redis.host", "localhost")
    v.SetDefault("redis.port", 6379)
    v.SetDefault("redis.db", 0)
    v.SetDefault("redis.pool_size", 10)
    v.SetDefault("redis.min_idle_conns", 5)
    
    // JWT默认配置
    v.SetDefault("jwt.expire_hours", 24)
    v.SetDefault("jwt.refresh_expire_hours", 168)
    v.SetDefault("jwt.issuer", "movieinfo")
    
    // 日志默认配置
    v.SetDefault("log.level", "info")
    v.SetDefault("log.format", "json")
    v.SetDefault("log.output", "stdout")
    
    // 缓存默认配置
    v.SetDefault("cache.default_expire", "3600s")
    v.SetDefault("cache.cleanup_interval", "600s")
    v.SetDefault("cache.max_size", 1000000)
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
    // 验证必需字段
    if config.Database.MySQL.Username == "" {
        return fmt.Errorf("database username is required")
    }
    
    if config.Database.MySQL.DBName == "" {
        return fmt.Errorf("database name is required")
    }
    
    if config.JWT.Secret == "" || config.JWT.Secret == "your_jwt_secret_key_change_in_production" {
        if config.App.Environment == "production" {
            return fmt.Errorf("JWT secret must be set in production")
        }
    }
    
    // 验证端口范围
    if config.Server.Web.Port < 1 || config.Server.Web.Port > 65535 {
        return fmt.Errorf("invalid web server port: %d", config.Server.Web.Port)
    }
    
    // 验证日志级别
    validLogLevels := []string{"debug", "info", "warn", "error"}
    isValidLevel := false
    for _, level := range validLogLevels {
        if config.Log.Level == level {
            isValidLevel = true
            break
        }
    }
    if !isValidLevel {
        return fmt.Errorf("invalid log level: %s", config.Log.Level)
    }
    
    return nil
}
```

### 5. **配置热更新**

#### 5.1 配置监听器
```go
// internal/config/watcher.go
package config

import (
    "context"
    "log"
    "sync"
    "time"
    
    "github.com/fsnotify/fsnotify"
    "github.com/spf13/viper"
)

// Watcher 配置监听器
type Watcher struct {
    viper    *viper.Viper
    config   *Config
    mutex    sync.RWMutex
    handlers []ReloadHandler
}

// ReloadHandler 配置重载处理器
type ReloadHandler func(oldConfig, newConfig *Config) error

// NewWatcher 创建配置监听器
func NewWatcher(configPath string) (*Watcher, error) {
    config, err := Load(configPath)
    if err != nil {
        return nil, err
    }
    
    v := viper.New()
    v.SetConfigFile(configPath)
    
    watcher := &Watcher{
        viper:    v,
        config:   config,
        handlers: make([]ReloadHandler, 0),
    }
    
    return watcher, nil
}

// AddReloadHandler 添加重载处理器
func (w *Watcher) AddReloadHandler(handler ReloadHandler) {
    w.handlers = append(w.handlers, handler)
}

// GetConfig 获取当前配置
func (w *Watcher) GetConfig() *Config {
    w.mutex.RLock()
    defer w.mutex.RUnlock()
    return w.config
}

// Watch 开始监听配置变化
func (w *Watcher) Watch(ctx context.Context) error {
    w.viper.WatchConfig()
    
    w.viper.OnConfigChange(func(e fsnotify.Event) {
        log.Printf("Config file changed: %s", e.Name)
        
        // 重新加载配置
        var newConfig Config
        if err := w.viper.Unmarshal(&newConfig); err != nil {
            log.Printf("Failed to unmarshal new config: %v", err)
            return
        }
        
        // 验证新配置
        if err := validateConfig(&newConfig); err != nil {
            log.Printf("New config validation failed: %v", err)
            return
        }
        
        w.mutex.Lock()
        oldConfig := w.config
        w.config = &newConfig
        w.mutex.Unlock()
        
        // 执行重载处理器
        for _, handler := range w.handlers {
            if err := handler(oldConfig, &newConfig); err != nil {
                log.Printf("Config reload handler failed: %v", err)
            }
        }
        
        log.Println("Config reloaded successfully")
    })
    
    <-ctx.Done()
    return nil
}
```

## 总结

配置文件设计为MovieInfo项目提供了灵活、安全、可维护的配置管理系统。通过分层配置、环境适配和热更新机制，我们建立了一个生产级的配置解决方案。

**关键设计要点**：
1. **分层配置**：支持默认值、文件配置、环境变量、命令行参数
2. **环境适配**：不同环境的配置隔离和覆盖
3. **类型安全**：强类型的配置结构定义
4. **验证机制**：配置有效性的自动验证
5. **热更新**：运行时配置的动态更新

**配置优势**：
- **灵活性**：支持多种配置来源和格式
- **安全性**：敏感信息的保护和验证
- **可维护性**：清晰的配置结构和文档
- **运维友好**：热更新和监控支持

**下一步**：基于这个配置系统，我们将实现统一的日志系统，支持结构化日志、日志级别控制和日志轮转。
