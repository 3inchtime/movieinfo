# 1.4 API 接口设计

## 概述

API接口设计是连接前端和后端的桥梁，它定义了系统的对外服务能力和交互方式。对于MovieInfo项目，我们需要设计一套既符合RESTful规范，又能满足业务需求的API接口体系。

## 为什么API设计如此重要？

### 1. **系统解耦**
- **前后端分离**：API作为契约，使前后端可以独立开发
- **服务边界**：清晰的API定义了服务的职责边界
- **版本管理**：API版本化支持系统平滑升级

### 2. **开发效率**
- **并行开发**：前后端团队可以基于API文档并行开发
- **测试便利**：标准化的API便于自动化测试
- **文档化**：API文档即系统功能说明

### 3. **用户体验**
- **响应速度**：合理的API设计提升响应性能
- **数据精确**：按需返回数据，减少网络传输
- **错误处理**：统一的错误格式提升用户体验

### 4. **系统扩展**
- **第三方集成**：标准化API便于第三方系统集成
- **移动端支持**：统一API支持多端应用
- **微服务演进**：为微服务化提供接口基础

## API设计原则

### 1. **RESTful设计原则**

#### 1.1 资源导向
API应该围绕资源（Resource）设计，而不是动作：
```
✅ 正确：GET /api/movies/123
❌ 错误：GET /api/getMovie?id=123

✅ 正确：POST /api/users
❌ 错误：POST /api/createUser
```

#### 1.2 HTTP方法语义
```
GET    - 获取资源（幂等）
POST   - 创建资源
PUT    - 更新整个资源（幂等）
PATCH  - 部分更新资源
DELETE - 删除资源（幂等）
```

#### 1.3 状态码规范
```
200 OK           - 请求成功
201 Created      - 资源创建成功
204 No Content   - 请求成功但无返回内容
400 Bad Request  - 请求参数错误
401 Unauthorized - 未认证
403 Forbidden    - 无权限
404 Not Found    - 资源不存在
422 Unprocessable Entity - 数据验证失败
500 Internal Server Error - 服务器内部错误
```

### 2. **一致性原则**

#### 2.1 命名规范
```
- 使用小写字母和连字符
- 复数形式表示集合：/api/movies
- 单数形式表示单个资源：/api/movies/123
- 嵌套资源：/api/movies/123/comments
```

#### 2.2 响应格式统一
```json
{
    "code": 200,
    "message": "success",
    "data": {
        // 具体数据
    },
    "meta": {
        "timestamp": "2024-01-01T00:00:00Z",
        "request_id": "req_123456"
    }
}
```

### 3. **性能优化原则**

#### 3.1 分页设计
```
GET /api/movies?page=1&page_size=20&sort=rating&order=desc
```

#### 3.2 字段选择
```
GET /api/movies?fields=id,title,rating,poster_url
```

#### 3.3 缓存策略
```
Cache-Control: public, max-age=3600
ETag: "33a64df551425fcc55e4d42a148795d9f25f89d4"
```

## API架构设计

### 1. 整体架构

```
┌─────────────────┐
│   前端应用      │
│  (Browser/App)  │
└─────────┬───────┘
          │ HTTP/HTTPS
          ▼
┌─────────────────┐
│   API网关       │
│  (Web Service)  │ ← 统一入口，路由分发
└─────────┬───────┘
          │ gRPC
    ┌─────┼─────┐
    │     │     │
    ▼     ▼     ▼
┌─────┐ ┌─────┐ ┌─────┐
│用户 │ │电影 │ │评论 │
│服务 │ │服务 │ │服务 │
└─────┘ └─────┘ └─────┘
```

### 2. API版本策略

#### 2.1 URL版本化
```
/api/v1/movies
/api/v2/movies
```

#### 2.2 Header版本化
```
Accept: application/vnd.movieinfo.v1+json
API-Version: v1
```

#### 2.3 版本兼容性
- **向后兼容**：新版本保持对旧版本的兼容
- **废弃通知**：提前通知API废弃计划
- **平滑迁移**：提供迁移指南和工具

## 核心API设计

### 1. 用户相关API

#### 1.1 用户注册
```http
POST /api/v1/users/register
Content-Type: application/json

{
    "email": "user@example.com",
    "username": "testuser",
    "password": "password123"
}
```

**响应示例**：
```json
{
    "code": 201,
    "message": "用户注册成功",
    "data": {
        "user": {
            "id": 123,
            "email": "user@example.com",
            "username": "testuser",
            "avatar_url": null,
            "created_at": "2024-01-01T00:00:00Z"
        },
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expires_at": "2024-01-02T00:00:00Z"
    }
}
```

**错误响应**：
```json
{
    "code": 422,
    "message": "数据验证失败",
    "errors": {
        "email": ["邮箱格式不正确"],
        "password": ["密码长度至少8位"]
    }
}
```

#### 1.2 用户登录
```http
POST /api/v1/users/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password123"
}
```

#### 1.3 获取用户信息
```http
GET /api/v1/users/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### 1.4 更新用户信息
```http
PATCH /api/v1/users/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "username": "newusername",
    "bio": "电影爱好者"
}
```

#### 1.5 密码重置
```http
POST /api/v1/users/password/reset
Content-Type: application/json

{
    "email": "user@example.com"
}
```

### 2. 电影相关API

#### 2.1 获取电影列表
```http
GET /api/v1/movies?page=1&page_size=20&category=action&sort=rating&order=desc
```

**查询参数**：
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20，最大100）
- `category`: 分类筛选
- `year`: 年份筛选
- `sort`: 排序字段（title, year, rating, created_at）
- `order`: 排序方向（asc, desc）
- `fields`: 返回字段选择

**响应示例**：
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "movies": [
            {
                "id": 1,
                "title": "肖申克的救赎",
                "poster_url": "https://example.com/poster1.jpg",
                "average_rating": 4.8,
                "rating_count": 1000,
                "release_year": 1994,
                "duration": 142
            }
        ],
        "pagination": {
            "current_page": 1,
            "page_size": 20,
            "total_pages": 50,
            "total_count": 1000,
            "has_next": true,
            "has_prev": false
        }
    }
}
```

#### 2.2 获取电影详情
```http
GET /api/v1/movies/123
```

**响应示例**：
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "movie": {
            "id": 123,
            "title": "肖申克的救赎",
            "original_title": "The Shawshank Redemption",
            "director": "弗兰克·德拉邦特",
            "actors": [
                {
                    "id": 1,
                    "name": "蒂姆·罗宾斯",
                    "role": "安迪·杜弗雷恩"
                }
            ],
            "categories": ["剧情", "犯罪"],
            "release_year": 1994,
            "duration": 142,
            "plot": "一个关于希望和友谊的故事...",
            "poster_url": "https://example.com/poster.jpg",
            "trailer_url": "https://example.com/trailer.mp4",
            "average_rating": 4.8,
            "rating_count": 1000,
            "user_rating": 5,
            "created_at": "2024-01-01T00:00:00Z"
        }
    }
}
```

#### 2.3 搜索电影
```http
GET /api/v1/movies/search?q=肖申克&category=drama&year=1994
```

**查询参数**：
- `q`: 搜索关键词（必需）
- `category`: 分类筛选
- `year`: 年份筛选
- `min_rating`: 最低评分
- `page`: 页码
- `page_size`: 每页数量

#### 2.4 获取电影分类
```http
GET /api/v1/categories
```

#### 2.5 获取推荐电影
```http
GET /api/v1/movies/recommendations?limit=10
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 3. 评论评分相关API

#### 3.1 获取电影评论
```http
GET /api/v1/movies/123/comments?page=1&page_size=20&sort=created_at&order=desc
```

#### 3.2 发表评论
```http
POST /api/v1/movies/123/comments
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "content": "这是一部非常优秀的电影！",
    "rating": 5
}
```

#### 3.3 更新评论
```http
PATCH /api/v1/comments/456
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "content": "更新后的评论内容"
}
```

#### 3.4 删除评论
```http
DELETE /api/v1/comments/456
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### 3.5 评论点赞
```http
POST /api/v1/comments/456/like
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### 3.6 电影评分
```http
POST /api/v1/movies/123/rating
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "score": 5
}
```

## 错误处理设计

### 1. 错误响应格式
```json
{
    "code": 400,
    "message": "请求参数错误",
    "error_code": "INVALID_PARAMETER",
    "errors": {
        "field_name": ["错误描述"]
    },
    "meta": {
        "timestamp": "2024-01-01T00:00:00Z",
        "request_id": "req_123456"
    }
}
```

### 2. 错误码设计
```
# 用户相关错误
USER_NOT_FOUND          - 用户不存在
USER_ALREADY_EXISTS     - 用户已存在
INVALID_CREDENTIALS     - 登录凭证无效
EMAIL_NOT_VERIFIED      - 邮箱未验证

# 电影相关错误
MOVIE_NOT_FOUND         - 电影不存在
INVALID_CATEGORY        - 无效的分类

# 评论相关错误
COMMENT_NOT_FOUND       - 评论不存在
COMMENT_PERMISSION_DENIED - 无权限操作评论
DUPLICATE_RATING        - 重复评分

# 系统错误
INTERNAL_ERROR          - 内部服务器错误
SERVICE_UNAVAILABLE     - 服务不可用
RATE_LIMIT_EXCEEDED     - 请求频率超限
```

### 3. 参数验证
```go
type MovieListRequest struct {
    Page     int    `form:"page" binding:"min=1"`
    PageSize int    `form:"page_size" binding:"min=1,max=100"`
    Category string `form:"category"`
    Year     int    `form:"year" binding:"min=1900,max=2030"`
    Sort     string `form:"sort" binding:"oneof=title year rating created_at"`
    Order    string `form:"order" binding:"oneof=asc desc"`
}
```

## 认证和授权

### 1. JWT认证
```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**JWT Payload**：
```json
{
    "user_id": 123,
    "username": "testuser",
    "email": "user@example.com",
    "exp": 1640995200,
    "iat": 1640908800
}
```

### 2. 权限控制
```go
// 权限级别
const (
    PermissionGuest  = "guest"    // 游客
    PermissionUser   = "user"     // 普通用户
    PermissionAdmin  = "admin"    // 管理员
)

// API权限要求
GET    /api/v1/movies          - guest
POST   /api/v1/movies/*/rating - user
DELETE /api/v1/comments/*      - user (own) / admin
```

### 3. 限流策略
```go
// 限流配置
type RateLimit struct {
    Guest int // 游客：100 req/min
    User  int // 用户：1000 req/min
    Admin int // 管理员：无限制
}
```

## 总结

API接口设计为MovieInfo项目提供了清晰的服务边界和交互规范。通过RESTful设计原则，我们确保了API的一致性和易用性；通过完善的错误处理和认证机制，我们保证了系统的安全性和可靠性。

**关键设计特点**：
1. **RESTful规范**：遵循REST设计原则，提供直观的API
2. **统一响应格式**：标准化的响应结构便于前端处理
3. **完善的错误处理**：详细的错误信息帮助快速定位问题
4. **安全认证**：JWT认证和权限控制保证系统安全
5. **性能优化**：分页、字段选择、缓存等优化策略

**下一步**：基于这个API设计，我们将进行UI/UX设计，确保前端界面能够很好地展示和操作这些API提供的数据。
