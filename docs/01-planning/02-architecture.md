# 1.2 系统架构设计

## 概述

系统架构设计是将需求转化为技术实现的关键环节。它定义了系统的整体结构、组件划分、交互方式和技术选型。对于MovieInfo项目，我们需要设计一个既能满足当前需求，又能为未来微服务化做准备的架构。

## 为什么需要架构设计？

### 1. **指导开发方向**
- **统一认知**：为开发团队提供清晰的技术蓝图
- **降低复杂度**：将复杂问题分解为可管理的模块
- **提高效率**：避免开发过程中的架构争议和返工

### 2. **保证系统质量**
- **性能保证**：合理的架构设计确保系统性能
- **可维护性**：良好的模块划分便于后期维护
- **可扩展性**：为系统未来发展预留空间

### 3. **控制技术风险**
- **技术选型**：基于需求选择合适的技术栈
- **风险评估**：识别和规避潜在的技术风险
- **成本控制**：平衡功能需求和开发成本

## 架构设计原则

### 1. **单一职责原则**
每个服务只负责一个业务领域，职责清晰：
- **用户服务**：专注用户管理和认证
- **电影服务**：专注电影数据管理
- **评论服务**：专注评论和评分功能
- **主页服务**：专注前端展示和API网关

**原理**：单一职责使得服务边界清晰，降低耦合度，便于独立开发和维护。

### 2. **高内聚低耦合**
- **高内聚**：相关功能聚集在同一个服务内
- **低耦合**：服务间通过明确的接口通信，减少依赖

**实现方式**：
- 使用gRPC定义清晰的服务接口
- 避免服务间直接数据库访问
- 通过事件或消息进行异步通信

### 3. **可扩展性设计**
- **水平扩展**：支持通过增加服务实例来提升性能
- **垂直扩展**：支持通过升级硬件来提升性能
- **功能扩展**：预留接口和扩展点

### 4. **容错性设计**
- **故障隔离**：一个服务的故障不影响其他服务
- **优雅降级**：在部分功能不可用时，核心功能仍可正常工作
- **快速恢复**：故障后能够快速恢复服务

## 整体架构设计

### 1. 架构风格选择

#### 1.1 服务化架构（为微服务做准备）
我们选择服务化架构而非传统的单体架构，原因如下：

**优势**：
- **业务隔离**：不同业务逻辑分离，便于团队协作
- **技术栈统一**：当前阶段使用统一技术栈，降低复杂度
- **平滑演进**：为未来微服务化奠定基础
- **独立部署**：每个服务可以独立启动和部署

**与微服务的区别**：
- **共享数据库**：当前阶段使用共享数据库，简化数据一致性
- **单一代码库**：所有服务在同一个代码仓库，便于管理
- **简化通信**：使用gRPC但不引入服务发现等复杂机制

#### 1.2 分层架构
每个服务内部采用分层架构：

```
┌─────────────────────────────────────┐
│           表示层 (Handlers)          │  ← HTTP/gRPC接口
├─────────────────────────────────────┤
│           业务层 (Services)          │  ← 业务逻辑处理
├─────────────────────────────────────┤
│           数据层 (Models/DAO)        │  ← 数据访问对象
├─────────────────────────────────────┤
│           基础设施层 (Database)       │  ← 数据库/缓存
└─────────────────────────────────────┘
```

**各层职责**：
- **表示层**：处理HTTP请求和gRPC调用，参数验证，响应格式化
- **业务层**：实现业务逻辑，调用数据层，处理业务规则
- **数据层**：数据访问抽象，CRUD操作，数据模型定义
- **基础设施层**：数据库连接，缓存，外部服务调用

### 2. 服务架构图

```
                    ┌─────────────────┐
                    │   用户访问      │
                    │   (Browser)     │
                    └─────────┬───────┘
                              │ HTTP
                              ▼
                    ┌─────────────────┐
                    │   主页服务      │
                    │  (Web Gateway)  │ ← Gin + Template
                    │   Port: 8080    │
                    └─────────┬───────┘
                              │ gRPC
                    ┌─────────┼───────┐
                    │         │       │
                    ▼         ▼       ▼
          ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
          │  用户服务   │ │  电影服务   │ │ 评论服务   │
          │(User Service)│ │(Movie Service)│ │(Comment Service)│
          │Port:8081/9081│ │Port:8082/9082│ │Port:8083/9083│
          └─────────────┘ └─────────────┘ └─────────────┘
                    │         │       │
                    └─────────┼───────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │   数据存储层    │
                    │  MySQL + Redis  │
                    └─────────────────┘
```

### 3. 技术栈选择

#### 3.1 后端技术栈

**Go语言 + Gin框架**
- **选择原因**：
  - 高性能：编译型语言，执行效率高
  - 并发优势：goroutine轻量级并发模型
  - 简洁语法：学习成本低，开发效率高
  - 生态丰富：丰富的第三方库支持

**Gin框架特点**：
- 轻量级：核心功能精简，性能优秀
- 中间件支持：灵活的中间件机制
- 路由功能：强大的路由匹配能力
- JSON支持：原生JSON处理能力

#### 3.2 数据存储技术栈

**MySQL数据库**
- **选择原因**：
  - 成熟稳定：经过长期验证的关系型数据库
  - 事务支持：ACID特性保证数据一致性
  - 性能优秀：经过优化的查询引擎
  - 社区支持：丰富的文档和社区资源

**Redis缓存**
- **选择原因**：
  - 高性能：内存存储，访问速度快
  - 数据结构丰富：支持多种数据类型
  - 持久化：支持数据持久化存储
  - 集群支持：支持主从复制和集群模式

#### 3.3 服务通信技术栈

**gRPC**
- **选择原因**：
  - 高性能：基于HTTP/2，支持多路复用
  - 类型安全：Protocol Buffers提供强类型定义
  - 跨语言：支持多种编程语言
  - 流式处理：支持流式数据传输

**Protocol Buffers**
- **优势**：
  - 序列化效率高：比JSON更紧凑
  - 向后兼容：支持协议演进
  - 代码生成：自动生成客户端和服务端代码
  - 文档化：proto文件即文档

## 详细架构设计

### 1. 主页服务架构

#### 1.1 职责定义
- **API网关**：统一对外提供服务入口
- **前端渲染**：使用Gin Template渲染HTML页面
- **请求路由**：将请求路由到相应的后端服务
- **静态资源**：提供CSS、JS、图片等静态资源

#### 1.2 技术组件
```go
// 主要组件
- Gin Router        // HTTP路由
- Template Engine   // 模板引擎
- Static Server     // 静态文件服务
- gRPC Clients      // 后端服务客户端
- Middleware        // 中间件（认证、日志等）
```

#### 1.3 数据流
```
用户请求 → Gin Router → 中间件 → Handler → gRPC调用 → 模板渲染 → 响应
```

### 2. 用户服务架构

#### 2.1 职责定义
- **用户管理**：用户注册、登录、信息管理
- **身份认证**：JWT Token生成和验证
- **密码管理**：密码加密、重置功能
- **会话管理**：用户会话状态管理

#### 2.2 核心组件
```go
// 数据模型
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Email    string `gorm:"unique;not null"`
    Password string `gorm:"not null"`
    Username string `gorm:"unique;not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

// 服务接口
type UserService interface {
    Register(req *RegisterRequest) (*User, error)
    Login(req *LoginRequest) (*LoginResponse, error)
    GetUser(userID uint) (*User, error)
    UpdateUser(userID uint, req *UpdateRequest) error
    ResetPassword(email string) error
}
```

#### 2.3 安全设计
- **密码加密**：使用bcrypt算法加密存储
- **JWT认证**：无状态的Token认证机制
- **输入验证**：严格的参数验证和过滤
- **防暴力破解**：登录失败次数限制

### 3. 电影服务架构

#### 3.1 职责定义
- **电影数据管理**：电影信息的CRUD操作
- **搜索功能**：电影搜索和筛选
- **分类管理**：电影分类和标签管理
- **推荐算法**：基础的电影推荐功能

#### 3.2 数据模型设计
```go
type Movie struct {
    ID          uint      `gorm:"primaryKey"`
    Title       string    `gorm:"not null;index"`
    Director    string    `gorm:"index"`
    Actors      string    // JSON存储演员列表
    Genre       string    `gorm:"index"`
    ReleaseYear int       `gorm:"index"`
    Duration    int       // 分钟
    Plot        text      // 剧情简介
    PosterURL   string    // 海报URL
    Rating      float64   `gorm:"default:0"`
    RatingCount int       `gorm:"default:0"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Category struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"unique;not null"`
    Slug string `gorm:"unique;not null"`
}
```

#### 3.3 搜索架构
```go
// 搜索服务接口
type SearchService interface {
    SearchMovies(query string, filters SearchFilters) (*SearchResult, error)
    GetMoviesByCategory(categoryID uint, page int) (*MovieList, error)
    GetRecommendations(userID uint, limit int) (*MovieList, error)
}

// 搜索过滤器
type SearchFilters struct {
    Genre       string
    Year        int
    MinRating   float64
    SortBy      string // title, year, rating
    SortOrder   string // asc, desc
}
```

### 4. 评论服务架构

#### 4.1 职责定义
- **评论管理**：用户评论的发布、编辑、删除
- **评分系统**：电影评分功能
- **内容审核**：评论内容的审核和过滤
- **统计分析**：评分统计和趋势分析

#### 4.2 数据模型
```go
type Comment struct {
    ID       uint      `gorm:"primaryKey"`
    UserID   uint      `gorm:"not null;index"`
    MovieID  uint      `gorm:"not null;index"`
    Content  string    `gorm:"type:text"`
    Rating   int       `gorm:"check:rating >= 1 AND rating <= 5"`
    Status   string    `gorm:"default:'pending'"` // pending, approved, rejected
    LikeCount int      `gorm:"default:0"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Rating struct {
    ID      uint `gorm:"primaryKey"`
    UserID  uint `gorm:"not null"`
    MovieID uint `gorm:"not null"`
    Score   int  `gorm:"check:score >= 1 AND score <= 5"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

#### 4.3 审核机制
```go
type ModerationService interface {
    CheckContent(content string) (*ModerationResult, error)
    ApproveComment(commentID uint) error
    RejectComment(commentID uint, reason string) error
    GetPendingComments(page int) (*CommentList, error)
}

type ModerationResult struct {
    IsApproved bool
    Confidence float64
    Reasons    []string
}
```

## 数据架构设计

### 1. 数据库设计原则

#### 1.1 规范化设计
- **第三范式**：消除数据冗余，保证数据一致性
- **适度反规范化**：为性能考虑，适当冗余常用字段
- **索引优化**：为查询频繁的字段建立索引

#### 1.2 数据分离策略
虽然当前使用共享数据库，但在逻辑上按服务分离：
- **用户数据**：users, user_sessions, user_profiles
- **电影数据**：movies, categories, movie_categories
- **评论数据**：comments, ratings, comment_likes

### 2. 缓存架构设计

#### 2.1 缓存策略
```go
// 缓存层次
L1: 应用内存缓存 (sync.Map)
L2: Redis缓存
L3: 数据库

// 缓存模式
- Cache-Aside: 应用控制缓存
- Write-Through: 写入时同步更新缓存
- Write-Behind: 异步写入数据库
```

#### 2.2 缓存键设计
```
用户缓存: user:{userID}
电影缓存: movie:{movieID}
电影列表: movies:page:{page}:category:{categoryID}
搜索结果: search:{hash(query+filters)}
评分统计: rating:movie:{movieID}
```

#### 2.3 缓存失效策略
- **TTL过期**：设置合理的过期时间
- **主动失效**：数据更新时主动清除相关缓存
- **版本控制**：使用版本号控制缓存一致性

## 接口设计

### 1. gRPC接口设计

#### 1.1 用户服务接口
```protobuf
service UserService {
    rpc Register(RegisterRequest) returns (RegisterResponse);
    rpc Login(LoginRequest) returns (LoginResponse);
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
    rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
    rpc ResetPassword(ResetPasswordRequest) returns (ResetPasswordResponse);
}

message RegisterRequest {
    string email = 1;
    string password = 2;
    string username = 3;
}

message LoginResponse {
    string token = 1;
    User user = 2;
    int64 expires_at = 3;
}
```

#### 1.2 电影服务接口
```protobuf
service MovieService {
    rpc GetMovie(GetMovieRequest) returns (GetMovieResponse);
    rpc ListMovies(ListMoviesRequest) returns (ListMoviesResponse);
    rpc SearchMovies(SearchMoviesRequest) returns (SearchMoviesResponse);
    rpc GetRecommendations(GetRecommendationsRequest) returns (GetRecommendationsResponse);
}

message ListMoviesRequest {
    int32 page = 1;
    int32 page_size = 2;
    string category = 3;
    string sort_by = 4;
}
```

### 2. HTTP接口设计

#### 2.1 RESTful API设计
```
GET    /api/movies              # 获取电影列表
GET    /api/movies/{id}         # 获取电影详情
GET    /api/movies/search       # 搜索电影
POST   /api/users/register      # 用户注册
POST   /api/users/login         # 用户登录
POST   /api/comments            # 发表评论
PUT    /api/comments/{id}       # 更新评论
DELETE /api/comments/{id}       # 删除评论
```

#### 2.2 响应格式标准化
```json
{
    "code": 200,
    "message": "success",
    "data": {
        // 具体数据
    },
    "timestamp": "2024-01-01T00:00:00Z"
}
```

## 部署架构

### 1. 单机部署架构
```
┌─────────────────────────────────────────────────────────┐
│                    Docker Host                          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │
│  │ Web Service │ │User Service │ │Movie Service│       │
│  │   :8080     │ │   :8081     │ │   :8082     │       │
│  └─────────────┘ └─────────────┘ └─────────────┘       │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │
│  │Comment Svc  │ │   MySQL     │ │   Redis     │       │
│  │   :8083     │ │   :3306     │ │   :6379     │       │
│  └─────────────┘ └─────────────┘ └─────────────┘       │
└─────────────────────────────────────────────────────────┘
```

### 2. 容器编排
```yaml
# docker-compose.yml
version: '3.8'
services:
  web:
    build: ./cmd/web
    ports: ["8080:8080"]
    depends_on: [user, movie, comment]
  
  user:
    build: ./cmd/user
    ports: ["8081:8081", "9081:9081"]
    depends_on: [mysql, redis]
  
  movie:
    build: ./cmd/movie
    ports: ["8082:8082", "9082:9082"]
    depends_on: [mysql, redis]
  
  comment:
    build: ./cmd/comment
    ports: ["8083:8083", "9083:9083"]
    depends_on: [mysql, redis]
```

## 总结

系统架构设计为MovieInfo项目提供了清晰的技术蓝图。通过服务化架构，我们实现了业务逻辑的分离和技术栈的统一，为未来的微服务化奠定了基础。

**关键设计决策**：
1. **服务化架构**：平衡了复杂度和扩展性
2. **技术栈统一**：降低了学习和维护成本
3. **gRPC通信**：提供了高性能的服务间通信
4. **分层设计**：确保了代码的可维护性

**架构优势**：
- **可扩展**：支持独立扩展各个服务
- **可维护**：清晰的模块划分便于维护
- **高性能**：优化的技术栈保证性能
- **易演进**：为微服务化预留了空间

下一步，我们将基于这个架构设计进行详细的数据库设计，确保数据模型能够支撑业务需求。
