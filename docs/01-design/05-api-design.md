# API接口设计

## 1. 设计原则

MovieInfo项目的API设计遵循RESTful架构风格，旨在提供一套直观、易用且可预测的接口。

- **无状态**：每个请求都包含所有必要信息，服务器不保存客户端会话状态。
- **统一接口**：使用标准的HTTP方法 (GET, POST, PUT, DELETE) 对资源进行操作。
- **资源导向**：URL代表资源，使用名词复数形式 (e.g., `/api/users`, `/api/movies`)。
- **JSON格式**：请求体和响应体均使用JSON格式。
- **版本控制**：API通过URL进行版本控制 (e.g., `/api/v1/...`)。

## 2. API通用规范

### 2.1. Base URL

所有API的根路径为 `/api/v1`。

### 2.2. 认证与授权

- **认证方式**：使用JWT (JSON Web Token) 进行认证。
- **流程**：
    1. 用户通过 `POST /api/v1/users/login` 登录，成功后获得一个JWT。
    2. 访问受保护的接口时，需在HTTP请求头中携带 `Authorization: Bearer <token>`。
- **公开接口**：登录、注册、电影列表、电影详情等接口无需认证。
- **受保护接口**：发表评论、修改个人信息等需要认证。

### 2.3. 响应格式

所有API响应都遵循统一的JSON结构：

```json
{
  "code": 0,          // 业务状态码，0表示成功，非0表示失败
  "message": "Success", // 提示信息
  "data": {}          // 响应数据，可以是对象或数组
}
```

### 2.4. 错误码

| Code | Message | Description |
|---|---|---|
| 0 | Success | 操作成功 |
| 40001 | Invalid Parameter | 请求参数无效 |
| 40101 | Unauthorized | 未认证或Token无效 |
| 40301 | Forbidden | 无权限访问 |
| 40401 | Not Found | 资源不存在 |
| 50001 | Internal Server Error | 服务器内部错误 |

## 3. API接口详述

### 3.1. 用户服务 (User Service)

#### `POST /api/v1/users/register`
- **描述**: 用户注册
- **请求体**: `{"email": "user@example.com", "password": "password123"}`
- **响应**: `{"code": 0, "message": "Success", "data": {"id": 1, "email": "user@example.com"}}`

#### `POST /api/v1/users/login`
- **描述**: 用户登录
- **请求体**: `{"email": "user@example.com", "password": "password123"}`
- **响应**: `{"code": 0, "message": "Success", "data": {"token": "jwt-token-string"}}`

#### `GET /api/v1/users/me`
- **描述**: 获取当前用户信息 (需认证)
- **响应**: `{"code": 0, "message": "Success", "data": {"id": 1, "email": "user@example.com", "nickname": "test"}}`

### 3.2. 电影服务 (Movie Service)

#### `GET /api/v1/movies`
- **描述**: 获取电影列表 (支持分页和分类过滤)
- **查询参数**: `?page=1&limit=10&category_id=1`
- **响应**: `{"code": 0, "message": "Success", "data": {"movies": [...], "total": 100}}`

#### `GET /api/v1/movies/{id}`
- **描述**: 获取单部电影详情
- **响应**: `{"code": 0, "message": "Success", "data": {"id": 1, "title": "Inception", ...}}`

### 3.3. 评论与评分服务 (Comment & Rating Service)

#### `GET /api/v1/movies/{movieId}/comments`
- **描述**: 获取某部电影的评论列表 (支持分页)
- **查询参数**: `?page=1&limit=10`
- **响应**: `{"code": 0, "message": "Success", "data": {"comments": [...], "total": 50}}`

#### `POST /api/v1/movies/{movieId}/comments`
- **描述**: 发表评论 (需认证)
- **请求体**: `{"content": "This is a great movie!", "parent_id": null}`
- **响应**: `{"code": 0, "message": "Success", "data": {"id": 101, ...}}`

#### `POST /api/v1/movies/{movieId}/ratings`
- **描述**: 对电影进行评分 (需认证)
- **请求体**: `{"score": 9}`
- **响应**: `{"code": 0, "message": "Success", "data": {}}`

---
