# 05-gRPC协议定义

## 目标说明

本步骤的目标是为 MovieInfo 项目定义完整的 gRPC 协议，建立服务间通信的标准接口：
- 设计用户服务的 gRPC 接口
- 设计电影服务的 gRPC 接口
- 设计评分服务的 gRPC 接口
- 定义通用的数据类型和错误处理
- 实现 gRPC 拦截器和中间件
- 集成日志、监控和链路追踪
- 提供 gRPC 客户端和服务端的通用工具

完成本步骤后，将拥有一套完整的 gRPC 协议定义，为后续的微服务开发提供标准的通信接口。

## 前置条件

- 已完成项目初始化（01-项目初始化.md）
- 已完成数据库设计（02-数据库设计.md）
- 已完成配置管理系统（03-配置管理系统.md）
- 已完成日志系统（04-日志系统.md）
- 熟悉 Protocol Buffers 语法
- 了解 gRPC 的基本概念和最佳实践
- 理解微服务架构的设计原则

## 技术要点

### gRPC 设计原则
- **接口优先**：先定义接口，再实现服务
- **向后兼容**：确保协议的向后兼容性
- **类型安全**：利用强类型系统避免运行时错误
- **性能优化**：使用二进制序列化提高性能

### 技术选型
- **Protocol Buffers**：使用 proto3 语法定义接口
- **gRPC-Go**：使用官方 Go 实现
- **gRPC 中间件**：集成认证、日志、监控等功能
- **错误处理**：使用 gRPC 状态码和详细错误信息

### 服务设计模式
- **领域驱动设计**：按业务领域划分服务
- **CQRS 模式**：分离命令和查询操作
- **分页查询**：支持大数据量的分页查询
- **批量操作**：提供批量操作接口提高效率

## 实现步骤

### 步骤1：定义通用数据类型

#### 1.1 创建通用类型定义

```protobuf
// api/proto/common/common.proto
syntax = "proto3";

package movieinfo.common;

option go_package = "movieinfo/api/proto/common";

import "google/protobuf/timestamp.proto";
import "google/protobuf/empty.proto";

// 分页请求
message PageRequest {
  int32 page = 1;           // 页码，从1开始
  int32 page_size = 2;      // 每页大小，默认10，最大100
  string sort_by = 3;       // 排序字段
  string sort_order = 4;    // 排序方向：asc, desc
}

// 分页响应
message PageResponse {
  int32 page = 1;           // 当前页码
  int32 page_size = 2;      // 每页大小
  int64 total = 3;          // 总记录数
  int32 total_pages = 4;    // 总页数
  bool has_next = 5;        // 是否有下一页
  bool has_prev = 6;        // 是否有上一页
}

// 通用响应状态
enum ResponseStatus {
  SUCCESS = 0;              // 成功
  FAILED = 1;               // 失败
  PARTIAL_SUCCESS = 2;      // 部分成功
}

// 通用响应
message CommonResponse {
  ResponseStatus status = 1; // 响应状态
  string message = 2;        // 响应消息
  string request_id = 3;     // 请求ID
  google.protobuf.Timestamp timestamp = 4; // 响应时间
}

// 批量操作请求
message BatchRequest {
  repeated string ids = 1;   // ID列表
  map<string, string> options = 2; // 操作选项
}

// 批量操作响应
message BatchResponse {
  int32 success_count = 1;   // 成功数量
  int32 failed_count = 2;    // 失败数量
  repeated BatchResult results = 3; // 详细结果
}

// 批量操作结果
message BatchResult {
  string id = 1;             // 操作的ID
  bool success = 2;          // 是否成功
  string error = 3;          // 错误信息
}

// 健康检查请求
message HealthCheckRequest {
  string service = 1;        // 服务名称
}

// 健康检查响应
message HealthCheckResponse {
  enum ServingStatus {
    UNKNOWN = 0;
    SERVING = 1;
    NOT_SERVING = 2;
    SERVICE_UNKNOWN = 3;
  }
  ServingStatus status = 1;
  string message = 2;
  map<string, string> details = 3;
}

// 版本信息
message VersionInfo {
  string version = 1;        // 版本号
  string build_time = 2;     // 构建时间
  string git_commit = 3;     // Git提交哈希
  string go_version = 4;     // Go版本
}

// 统计信息
message Statistics {
  int64 total_requests = 1;  // 总请求数
  int64 success_requests = 2; // 成功请求数
  int64 failed_requests = 3;  // 失败请求数
  double avg_response_time = 4; // 平均响应时间(ms)
  google.protobuf.Timestamp last_request_time = 5; // 最后请求时间
}
```

#### 1.2 创建错误定义

```protobuf
// api/proto/common/error.proto
syntax = "proto3";

package movieinfo.common;

option go_package = "movieinfo/api/proto/common";

// 错误代码
enum ErrorCode {
  // 通用错误
  UNKNOWN_ERROR = 0;
  INVALID_ARGUMENT = 1;
  PERMISSION_DENIED = 2;
  NOT_FOUND = 3;
  ALREADY_EXISTS = 4;
  RESOURCE_EXHAUSTED = 5;
  INTERNAL_ERROR = 6;
  UNAVAILABLE = 7;
  DEADLINE_EXCEEDED = 8;
  
  // 认证相关错误
  UNAUTHENTICATED = 100;
  TOKEN_EXPIRED = 101;
  TOKEN_INVALID = 102;
  
  // 用户相关错误
  USER_NOT_FOUND = 200;
  USER_ALREADY_EXISTS = 201;
  USER_DISABLED = 202;
  INVALID_CREDENTIALS = 203;
  
  // 电影相关错误
  MOVIE_NOT_FOUND = 300;
  MOVIE_ALREADY_EXISTS = 301;
  INVALID_MOVIE_DATA = 302;
  
  // 评分相关错误
  RATING_NOT_FOUND = 400;
  RATING_ALREADY_EXISTS = 401;
  INVALID_RATING_VALUE = 402;
  
  // 业务逻辑错误
  BUSINESS_RULE_VIOLATION = 500;
  DATA_CONSISTENCY_ERROR = 501;
  OPERATION_NOT_ALLOWED = 502;
}

// 错误详情
message ErrorDetail {
  ErrorCode code = 1;        // 错误代码
  string message = 2;        // 错误消息
  string field = 3;          // 相关字段
  map<string, string> metadata = 4; // 额外元数据
}

// 验证错误
message ValidationError {
  string field = 1;          // 字段名
  string message = 2;        // 错误消息
  string value = 3;          // 错误值
}

// 业务错误
message BusinessError {
  string code = 1;           // 业务错误代码
  string message = 2;        // 错误消息
  map<string, string> context = 3; // 错误上下文
}
```

### 步骤2：定义用户服务接口

#### 2.1 创建用户数据类型

```protobuf
// api/proto/user/user.proto
syntax = "proto3";

package movieinfo.user;

option go_package = "movieinfo/api/proto/user";

import "google/protobuf/timestamp.proto";
import "google/protobuf/empty.proto";
import "common/common.proto";
import "common/error.proto";

// 用户信息
message User {
  int64 id = 1;              // 用户ID
  string username = 2;       // 用户名
  string email = 3;          // 邮箱
  string phone = 4;          // 手机号
  string nickname = 5;       // 昵称
  string avatar = 6;         // 头像URL
  UserStatus status = 7;     // 用户状态
  google.protobuf.Timestamp created_at = 8;  // 创建时间
  google.protobuf.Timestamp updated_at = 9;  // 更新时间
  google.protobuf.Timestamp last_login_at = 10; // 最后登录时间
  UserProfile profile = 11;  // 用户资料
}

// 用户状态
enum UserStatus {
  USER_STATUS_UNKNOWN = 0;
  USER_STATUS_ACTIVE = 1;    // 活跃
  USER_STATUS_INACTIVE = 2;  // 非活跃
  USER_STATUS_SUSPENDED = 3; // 暂停
  USER_STATUS_DELETED = 4;   // 已删除
}

// 用户资料
message UserProfile {
  string bio = 1;            // 个人简介
  string location = 2;       // 所在地
  string website = 3;        // 个人网站
  google.protobuf.Timestamp birthday = 4; // 生日
  Gender gender = 5;         // 性别
  map<string, string> preferences = 6; // 用户偏好
}

// 性别
enum Gender {
  GENDER_UNKNOWN = 0;
  GENDER_MALE = 1;
  GENDER_FEMALE = 2;
  GENDER_OTHER = 3;
}

// 创建用户请求
message CreateUserRequest {
  string username = 1;       // 用户名
  string email = 2;          // 邮箱
  string password = 3;       // 密码
  string phone = 4;          // 手机号
  string nickname = 5;       // 昵称
  UserProfile profile = 6;   // 用户资料
}

// 创建用户响应
message CreateUserResponse {
  movieinfo.common.CommonResponse common = 1;
  User user = 2;             // 创建的用户信息
}

// 获取用户请求
message GetUserRequest {
  oneof identifier {
    int64 id = 1;            // 用户ID
    string username = 2;     // 用户名
    string email = 3;        // 邮箱
  }
  bool include_profile = 4;  // 是否包含详细资料
}

// 获取用户响应
message GetUserResponse {
  movieinfo.common.CommonResponse common = 1;
  User user = 2;             // 用户信息
}

// 更新用户请求
message UpdateUserRequest {
  int64 id = 1;              // 用户ID
  string nickname = 2;       // 昵称
  string avatar = 3;         // 头像URL
  UserProfile profile = 4;   // 用户资料
  repeated string update_mask = 5; // 更新字段掩码
}

// 更新用户响应
message UpdateUserResponse {
  movieinfo.common.CommonResponse common = 1;
  User user = 2;             // 更新后的用户信息
}

// 删除用户请求
message DeleteUserRequest {
  int64 id = 1;              // 用户ID
  bool soft_delete = 2;      // 是否软删除
}

// 删除用户响应
message DeleteUserResponse {
  movieinfo.common.CommonResponse common = 1;
}

// 列出用户请求
message ListUsersRequest {
  movieinfo.common.PageRequest page = 1; // 分页参数
  UserStatus status = 2;     // 用户状态过滤
  string search = 3;         // 搜索关键词
  google.protobuf.Timestamp created_after = 4;  // 创建时间过滤
  google.protobuf.Timestamp created_before = 5; // 创建时间过滤
}

// 列出用户响应
message ListUsersResponse {
  movieinfo.common.CommonResponse common = 1;
  repeated User users = 2;   // 用户列表
  movieinfo.common.PageResponse page = 3; // 分页信息
}

// 用户登录请求
message LoginRequest {
  oneof credential {
    string username = 1;     // 用户名
    string email = 2;        // 邮箱
    string phone = 3;        // 手机号
  }
  string password = 4;       // 密码
  string client_ip = 5;      // 客户端IP
  string user_agent = 6;     // 用户代理
}

// 用户登录响应
message LoginResponse {
  movieinfo.common.CommonResponse common = 1;
  string access_token = 2;   // 访问令牌
  string refresh_token = 3;  // 刷新令牌
  int64 expires_in = 4;      // 过期时间（秒）
  User user = 5;             // 用户信息
}

// 刷新令牌请求
message RefreshTokenRequest {
  string refresh_token = 1;  // 刷新令牌
}

// 刷新令牌响应
message RefreshTokenResponse {
  movieinfo.common.CommonResponse common = 1;
  string access_token = 2;   // 新的访问令牌
  string refresh_token = 3;  // 新的刷新令牌
  int64 expires_in = 4;      // 过期时间（秒）
}

// 用户登出请求
message LogoutRequest {
  string access_token = 1;   // 访问令牌
}

// 用户登出响应
message LogoutResponse {
  movieinfo.common.CommonResponse common = 1;
}

// 修改密码请求
message ChangePasswordRequest {
  int64 user_id = 1;         // 用户ID
  string old_password = 2;   // 旧密码
  string new_password = 3;   // 新密码
}

// 修改密码响应
message ChangePasswordResponse {
  movieinfo.common.CommonResponse common = 1;
}

// 重置密码请求
message ResetPasswordRequest {
  string email = 1;          // 邮箱
  string reset_code = 2;     // 重置码
  string new_password = 3;   // 新密码
}

// 重置密码响应
message ResetPasswordResponse {
  movieinfo.common.CommonResponse common = 1;
}

// 发送重置码请求
message SendResetCodeRequest {
  string email = 1;          // 邮箱
}

// 发送重置码响应
message SendResetCodeResponse {
  movieinfo.common.CommonResponse common = 1;
  string message = 2;        // 提示信息
}
```

#### 2.2 创建用户服务接口

```protobuf
// api/proto/user/user_service.proto
syntax = "proto3";

package movieinfo.user;

option go_package = "movieinfo/api/proto/user";

import "user/user.proto";
import "common/common.proto";
import "google/protobuf/empty.proto";

// 用户服务
service UserService {
  // 用户管理
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  
  // 批量操作
  rpc BatchGetUsers(movieinfo.common.BatchRequest) returns (stream GetUserResponse);
  rpc BatchDeleteUsers(movieinfo.common.BatchRequest) returns (movieinfo.common.BatchResponse);
  
  // 认证相关
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc Logout(LogoutRequest) returns (LogoutResponse);
  rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse);
  
  // 密码管理
  rpc ChangePassword(ChangePasswordRequest) returns (ChangePasswordResponse);
  rpc ResetPassword(ResetPasswordRequest) returns (ResetPasswordResponse);
  rpc SendResetCode(SendResetCodeRequest) returns (SendResetCodeResponse);
  
  // 健康检查和统计
  rpc HealthCheck(movieinfo.common.HealthCheckRequest) returns (movieinfo.common.HealthCheckResponse);
  rpc GetStatistics(google.protobuf.Empty) returns (movieinfo.common.Statistics);
  rpc GetVersion(google.protobuf.Empty) returns (movieinfo.common.VersionInfo);
}
```

### 步骤3：定义电影服务接口

#### 3.1 创建电影数据类型

```protobuf
// api/proto/movie/movie.proto
syntax = "proto3";

package movieinfo.movie;

option go_package = "movieinfo/api/proto/movie";

import "google/protobuf/timestamp.proto";
import "google/protobuf/empty.proto";
import "common/common.proto";
import "common/error.proto";

// 电影信息
message Movie {
  int64 id = 1;              // 电影ID
  string title = 2;          // 电影标题
  string original_title = 3; // 原始标题
  string description = 4;    // 电影描述
  string poster_url = 5;     // 海报URL
  string trailer_url = 6;    // 预告片URL
  int32 duration = 7;        // 时长（分钟）
  google.protobuf.Timestamp release_date = 8; // 上映日期
  string language = 9;       // 语言
  string country = 10;       // 国家
  MovieStatus status = 11;   // 电影状态
  repeated string genres = 12; // 类型列表
  repeated string directors = 13; // 导演列表
  repeated string actors = 14;    // 演员列表
  MovieRating rating = 15;   // 评分信息
  google.protobuf.Timestamp created_at = 16; // 创建时间
  google.protobuf.Timestamp updated_at = 17; // 更新时间
  map<string, string> metadata = 18; // 额外元数据
}

// 电影状态
enum MovieStatus {
  MOVIE_STATUS_UNKNOWN = 0;
  MOVIE_STATUS_DRAFT = 1;    // 草稿
  MOVIE_STATUS_PUBLISHED = 2; // 已发布
  MOVIE_STATUS_ARCHIVED = 3;  // 已归档
  MOVIE_STATUS_DELETED = 4;   // 已删除
}

// 电影评分信息
message MovieRating {
  double average_rating = 1; // 平均评分
  int64 rating_count = 2;    // 评分数量
  repeated RatingDistribution distribution = 3; // 评分分布
}

// 评分分布
message RatingDistribution {
  int32 rating = 1;          // 评分值（1-5）
  int64 count = 2;           // 该评分的数量
  double percentage = 3;     // 百分比
}

// 电影类别
message Category {
  int64 id = 1;              // 类别ID
  string name = 2;           // 类别名称
  string description = 3;    // 类别描述
  string slug = 4;           // URL友好的标识符
  int32 sort_order = 5;      // 排序顺序
  bool is_active = 6;        // 是否激活
  google.protobuf.Timestamp created_at = 7; // 创建时间
  google.protobuf.Timestamp updated_at = 8; // 更新时间
}

// 创建电影请求
message CreateMovieRequest {
  string title = 1;          // 电影标题
  string original_title = 2; // 原始标题
  string description = 3;    // 电影描述
  string poster_url = 4;     // 海报URL
  string trailer_url = 5;    // 预告片URL
  int32 duration = 6;        // 时长（分钟）
  google.protobuf.Timestamp release_date = 7; // 上映日期
  string language = 8;       // 语言
  string country = 9;        // 国家
  repeated int64 category_ids = 10; // 类别ID列表
  repeated string directors = 11;   // 导演列表
  repeated string actors = 12;      // 演员列表
  map<string, string> metadata = 13; // 额外元数据
}

// 创建电影响应
message CreateMovieResponse {
  movieinfo.common.CommonResponse common = 1;
  Movie movie = 2;           // 创建的电影信息
}

// 获取电影请求
message GetMovieRequest {
  int64 id = 1;              // 电影ID
  bool include_rating = 2;   // 是否包含评分信息
}

// 获取电影响应
message GetMovieResponse {
  movieinfo.common.CommonResponse common = 1;
  Movie movie = 2;           // 电影信息
}

// 更新电影请求
message UpdateMovieRequest {
  int64 id = 1;              // 电影ID
  string title = 2;          // 电影标题
  string description = 3;    // 电影描述
  string poster_url = 4;     // 海报URL
  string trailer_url = 5;    // 预告片URL
  int32 duration = 6;        // 时长（分钟）
  google.protobuf.Timestamp release_date = 7; // 上映日期
  string language = 8;       // 语言
  string country = 9;        // 国家
  MovieStatus status = 10;   // 电影状态
  repeated int64 category_ids = 11; // 类别ID列表
  repeated string directors = 12;   // 导演列表
  repeated string actors = 13;      // 演员列表
  map<string, string> metadata = 14; // 额外元数据
  repeated string update_mask = 15;  // 更新字段掩码
}

// 更新电影响应
message UpdateMovieResponse {
  movieinfo.common.CommonResponse common = 1;
  Movie movie = 2;           // 更新后的电影信息
}

// 删除电影请求
message DeleteMovieRequest {
  int64 id = 1;              // 电影ID
  bool soft_delete = 2;      // 是否软删除
}

// 删除电影响应
message DeleteMovieResponse {
  movieinfo.common.CommonResponse common = 1;
}

// 列出电影请求
message ListMoviesRequest {
  movieinfo.common.PageRequest page = 1; // 分页参数
  repeated int64 category_ids = 2;       // 类别过滤
  MovieStatus status = 3;                // 状态过滤
  string search = 4;                     // 搜索关键词
  string language = 5;                   // 语言过滤
  string country = 6;                    // 国家过滤
  google.protobuf.Timestamp release_after = 7;  // 上映时间过滤
  google.protobuf.Timestamp release_before = 8; // 上映时间过滤
  double min_rating = 9;                 // 最低评分
  double max_rating = 10;                // 最高评分
  bool include_rating = 11;              // 是否包含评分信息
}

// 列出电影响应
message ListMoviesResponse {
  movieinfo.common.CommonResponse common = 1;
  repeated Movie movies = 2;  // 电影列表
  movieinfo.common.PageResponse page = 3; // 分页信息
}

// 搜索电影请求
message SearchMoviesRequest {
  string query = 1;          // 搜索查询
  movieinfo.common.PageRequest page = 2; // 分页参数
  repeated string fields = 3; // 搜索字段
  map<string, string> filters = 4; // 过滤条件
}

// 搜索电影响应
message SearchMoviesResponse {
  movieinfo.common.CommonResponse common = 1;
  repeated Movie movies = 2;  // 电影列表
  movieinfo.common.PageResponse page = 3; // 分页信息
  repeated string suggestions = 4; // 搜索建议
}

// 获取热门电影请求
message GetPopularMoviesRequest {
  movieinfo.common.PageRequest page = 1; // 分页参数
  string time_range = 2;     // 时间范围：day, week, month, year
  repeated int64 category_ids = 3; // 类别过滤
}

// 获取热门电影响应
message GetPopularMoviesResponse {
  movieinfo.common.CommonResponse common = 1;
  repeated Movie movies = 2;  // 电影列表
  movieinfo.common.PageResponse page = 3; // 分页信息
}

// 获取推荐电影请求
message GetRecommendedMoviesRequest {
  int64 user_id = 1;         // 用户ID
  movieinfo.common.PageRequest page = 2; // 分页参数
  string algorithm = 3;      // 推荐算法
}

// 获取推荐电影响应
message GetRecommendedMoviesResponse {
  movieinfo.common.CommonResponse common = 1;
  repeated Movie movies = 2;  // 电影列表
  movieinfo.common.PageResponse page = 3; // 分页信息
  string algorithm_used = 4; // 使用的算法
}
```

#### 3.2 创建电影服务接口

```protobuf
// api/proto/movie/movie_service.proto
syntax = "proto3";

package movieinfo.movie;

option go_package = "movieinfo/api/proto/movie";

import "movie/movie.proto";
import "common/common.proto";
import "google/protobuf/empty.proto";

// 电影服务
service MovieService {
  // 电影管理
  rpc CreateMovie(CreateMovieRequest) returns (CreateMovieResponse);
  rpc GetMovie(GetMovieRequest) returns (GetMovieResponse);
  rpc UpdateMovie(UpdateMovieRequest) returns (UpdateMovieResponse);
  rpc DeleteMovie(DeleteMovieRequest) returns (DeleteMovieResponse);
  rpc ListMovies(ListMoviesRequest) returns (ListMoviesResponse);
  
  // 批量操作
  rpc BatchGetMovies(movieinfo.common.BatchRequest) returns (stream GetMovieResponse);
  rpc BatchDeleteMovies(movieinfo.common.BatchRequest) returns (movieinfo.common.BatchResponse);
  
  // 搜索和发现
  rpc SearchMovies(SearchMoviesRequest) returns (SearchMoviesResponse);
  rpc GetPopularMovies(GetPopularMoviesRequest) returns (GetPopularMoviesResponse);
  rpc GetRecommendedMovies(GetRecommendedMoviesRequest) returns (GetRecommendedMoviesResponse);
  
  // 类别管理
  rpc CreateCategory(CreateCategoryRequest) returns (CreateCategoryResponse);
  rpc GetCategory(GetCategoryRequest) returns (GetCategoryResponse);
  rpc UpdateCategory(UpdateCategoryRequest) returns (UpdateCategoryResponse);
  rpc DeleteCategory(DeleteCategoryRequest) returns (DeleteCategoryResponse);
  rpc ListCategories(ListCategoriesRequest) returns (ListCategoriesResponse);
  
  // 健康检查和统计
  rpc HealthCheck(movieinfo.common.HealthCheckRequest) returns (movieinfo.common.HealthCheckResponse);
  rpc GetStatistics(google.protobuf.Empty) returns (movieinfo.common.Statistics);
  rpc GetVersion(google.protobuf.Empty) returns (movieinfo.common.VersionInfo);
}

// 类别相关的消息定义
message CreateCategoryRequest {
  string name = 1;
  string description = 2;
  string slug = 3;
  int32 sort_order = 4;
}

message CreateCategoryResponse {
  movieinfo.common.CommonResponse common = 1;
  Category category = 2;
}

message GetCategoryRequest {
  oneof identifier {
    int64 id = 1;
    string slug = 2;
  }
}

message GetCategoryResponse {
  movieinfo.common.CommonResponse common = 1;
  Category category = 2;
}

message UpdateCategoryRequest {
  int64 id = 1;
  string name = 2;
  string description = 3;
  string slug = 4;
  int32 sort_order = 5;
  bool is_active = 6;
  repeated string update_mask = 7;
}

message UpdateCategoryResponse {
  movieinfo.common.CommonResponse common = 1;
  Category category = 2;
}

message DeleteCategoryRequest {
  int64 id = 1;
}

message DeleteCategoryResponse {
  movieinfo.common.CommonResponse common = 1;
}

message ListCategoriesRequest {
  movieinfo.common.PageRequest page = 1;
  bool active_only = 2;
  string search = 3;
}

message ListCategoriesResponse {
  movieinfo.common.CommonResponse common = 1;
  repeated Category categories = 2;
  movieinfo.common.PageResponse page = 3;
}
```

### 步骤4：定义评分服务接口

#### 4.1 创建评分数据类型

```protobuf
// api/proto/rating/rating.proto
syntax = "proto3";

package movieinfo.rating;

option go_package = "movieinfo/api/proto/rating";

import "google/protobuf/timestamp.proto";
import "google/protobuf/empty.proto";
import "common/common.proto";
import "common/error.proto";

// 用户评分
message UserRating {
  int64 id = 1;              // 评分ID
  int64 user_id = 2;         // 用户ID
  int64 movie_id = 3;        // 电影ID
  int32 rating = 4;          // 评分值（1-5）
  string comment = 5;        // 评论内容
  RatingStatus status = 6;   // 评分状态
  google.protobuf.Timestamp created_at = 7; // 创建时间
  google.protobuf.Timestamp updated_at = 8; // 更新时间
  UserInfo user_info = 9;    // 用户信息（可选）
  MovieInfo movie_info = 10; // 电影信息（可选）
  int32 helpful_count = 11;  // 有用数
  int32 unhelpful_count = 12; // 无用数
}

// 评分状态
enum RatingStatus {
  RATING_STATUS_UNKNOWN = 0;
  RATING_STATUS_ACTIVE = 1;   // 活跃
  RATING_STATUS_HIDDEN = 2;   // 隐藏
  RATING_STATUS_DELETED = 3;  // 已删除
  RATING_STATUS_REPORTED = 4; // 被举报
}

// 用户信息（简化版）
message UserInfo {
  int64 id = 1;
  string username = 2;
  string nickname = 3;
  string avatar = 4;
}

// 电影信息（简化版）
message MovieInfo {
  int64 id = 1;
  string title = 2;
  string poster_url = 3;
  google.protobuf.Timestamp release_date = 4;
}

// 评分统计
message RatingStatistics {
  int64 movie_id = 1;        // 电影ID
  double average_rating = 2; // 平均评分
  int64 total_ratings = 3;   // 总评分数
  repeated RatingCount rating_counts = 4; // 各评分数量
  google.protobuf.Timestamp last_updated = 5; // 最后更新时间
}

// 评分数量统计
message RatingCount {
  int32 rating = 1;          // 评分值
  int64 count = 2;           // 数量
  double percentage = 3;     // 百分比
}

// 创建评分请求
message CreateRatingRequest {
  int64 user_id = 1;         // 用户ID
  int64 movie_id = 2;        // 电影ID
  int32 rating = 3;          // 评分值（1-5）
  string comment = 4;        // 评论内容
}

// 创建评分响应
message CreateRatingResponse {
  movieinfo.common.CommonResponse common = 1;
  UserRating rating = 2;     // 创建的评分
}

// 获取评分请求
message GetRatingRequest {
  oneof identifier {
    int64 id = 1;            // 评分ID
    UserMovieRating user_movie = 2; // 用户-电影组合
  }
  bool include_user_info = 3;  // 是否包含用户信息
  bool include_movie_info = 4; // 是否包含电影信息
}

// 用户-电影评分标识
message UserMovieRating {
  int64 user_id = 1;
  int64 movie_id = 2;
}

// 获取评分响应
message GetRatingResponse {
  movieinfo.common.CommonResponse common = 1;
  UserRating rating = 2;     // 评分信息
}

// 更新评分请求
message UpdateRatingRequest {
  int64 id = 1;              // 评分ID
  int32 rating = 2;          // 评分值（1-5）
  string comment = 3;        // 评论内容
  RatingStatus status = 4;   // 评分状态
  repeated string update_mask = 5; // 更新字段掩码
}

// 更新评分响应
message UpdateRatingResponse {
  movieinfo.common.CommonResponse common = 1;
  UserRating rating = 2;     // 更新后的评分
}

// 删除评分请求
message DeleteRatingRequest {
  int64 id = 1;              // 评分ID
  bool soft_delete = 2;      // 是否软删除
}

// 删除评分响应
message DeleteRatingResponse {
  movieinfo.common.CommonResponse common = 1;
}

// 列出评分请求
message ListRatingsRequest {
  movieinfo.common.PageRequest page = 1; // 分页参数
  int64 user_id = 2;         // 用户ID过滤
  int64 movie_id = 3;        // 电影ID过滤
  int32 min_rating = 4;      // 最低评分
  int32 max_rating = 5;      // 最高评分
  RatingStatus status = 6;   // 状态过滤
  google.protobuf.Timestamp created_after = 7;  // 创建时间过滤
  google.protobuf.Timestamp created_before = 8; // 创建时间过滤
  bool include_user_info = 9;  // 是否包含用户信息
  bool include_movie_info = 10; // 是否包含电影信息
  string sort_by = 11;       // 排序字段：created_at, rating, helpful_count
}

// 列出评分响应
message ListRatingsResponse {
  movieinfo.common.CommonResponse common = 1;
  repeated UserRating ratings = 2; // 评分列表
  movieinfo.common.PageResponse page = 3; // 分页信息
}

// 获取用户评分请求
message GetUserRatingsRequest {
  int64 user_id = 1;         // 用户ID
  movieinfo.common.PageRequest page = 2; // 分页参数
  bool include_movie_info = 3; // 是否包含电影信息
}

// 获取用户评分响应
message GetUserRatingsResponse {
  movieinfo.common.CommonResponse common = 1;
  repeated UserRating ratings = 2; // 评分列表
  movieinfo.common.PageResponse page = 3; // 分页信息
}

// 获取电影评分请求
message GetMovieRatingsRequest {
  int64 movie_id = 1;        // 电影ID
  movieinfo.common.PageRequest page = 2; // 分页参数
  bool include_user_info = 3; // 是否包含用户信息
  string sort_by = 4;        // 排序方式：newest, oldest, highest, lowest, helpful
}

// 获取电影评分响应
message GetMovieRatingsResponse {
  movieinfo.common.CommonResponse common = 1;
  repeated UserRating ratings = 2; // 评分列表
  movieinfo.common.PageResponse page = 3; // 分页信息
  RatingStatistics statistics = 4; // 评分统计
}

// 获取评分统计请求
message GetRatingStatisticsRequest {
  int64 movie_id = 1;        // 电影ID
}

// 获取评分统计响应
message GetRatingStatisticsResponse {
  movieinfo.common.CommonResponse common = 1;
  RatingStatistics statistics = 2; // 评分统计
}

// 标记评分有用请求
message MarkRatingHelpfulRequest {
  int64 rating_id = 1;       // 评分ID
  int64 user_id = 2;         // 用户ID
  bool helpful = 3;          // 是否有用
}

// 标记评分有用响应
message MarkRatingHelpfulResponse {
  movieinfo.common.CommonResponse common = 1;
  int32 helpful_count = 2;   // 有用数
  int32 unhelpful_count = 3; // 无用数
}

// 举报评分请求
message ReportRatingRequest {
  int64 rating_id = 1;       // 评分ID
  int64 reporter_id = 2;     // 举报者ID
  string reason = 3;         // 举报原因
  string description = 4;    // 详细描述
}

// 举报评分响应
message ReportRatingResponse {
  movieinfo.common.CommonResponse common = 1;
}
```

#### 4.2 创建评分服务接口

```protobuf
// api/proto/rating/rating_service.proto
syntax = "proto3";

package movieinfo.rating;

option go_package = "movieinfo/api/proto/rating";

import "rating/rating.proto";
import "common/common.proto";
import "google/protobuf/empty.proto";

// 评分服务
service RatingService {
  // 评分管理
  rpc CreateRating(CreateRatingRequest) returns (CreateRatingResponse);
  rpc GetRating(GetRatingRequest) returns (GetRatingResponse);
  rpc UpdateRating(UpdateRatingRequest) returns (UpdateRatingResponse);
  rpc DeleteRating(DeleteRatingRequest) returns (DeleteRatingResponse);
  rpc ListRatings(ListRatingsRequest) returns (ListRatingsResponse);
  
  // 批量操作
  rpc BatchGetRatings(movieinfo.common.BatchRequest) returns (stream GetRatingResponse);
  rpc BatchDeleteRatings(movieinfo.common.BatchRequest) returns (movieinfo.common.BatchResponse);
  
  // 用户相关
  rpc GetUserRatings(GetUserRatingsRequest) returns (GetUserRatingsResponse);
  
  // 电影相关
  rpc GetMovieRatings(GetMovieRatingsRequest) returns (GetMovieRatingsResponse);
  rpc GetRatingStatistics(GetRatingStatisticsRequest) returns (GetRatingStatisticsResponse);
  
  // 交互功能
  rpc MarkRatingHelpful(MarkRatingHelpfulRequest) returns (MarkRatingHelpfulResponse);
  rpc ReportRating(ReportRatingRequest) returns (ReportRatingResponse);
  
  // 健康检查和统计
  rpc HealthCheck(movieinfo.common.HealthCheckRequest) returns (movieinfo.common.HealthCheckResponse);
  rpc GetStatistics(google.protobuf.Empty) returns (movieinfo.common.Statistics);
  rpc GetVersion(google.protobuf.Empty) returns (movieinfo.common.VersionInfo);
}
```

### 步骤5：创建 gRPC 工具和中间件

#### 5.1 创建 gRPC 配置

```go
// pkg/grpc/config.go
package grpc

import (
	"time"
)

// ServerConfig gRPC服务器配置
type ServerConfig struct {
	// 基础配置
	Host string `yaml:"host" validate:"required"`
	Port int    `yaml:"port" validate:"required,min=1,max=65535"`

	// TLS配置
	TLS TLSConfig `yaml:"tls"`

	// 超时配置
	Timeout TimeoutConfig `yaml:"timeout"`

	// 限流配置
	RateLimit RateLimitConfig `yaml:"rate_limit"`

	// 中间件配置
	Middleware MiddlewareConfig `yaml:"middleware"`

	// 健康检查配置
	HealthCheck HealthCheckConfig `yaml:"health_check"`

	// 反射配置
	Reflection bool `yaml:"reflection"`
}

// ClientConfig gRPC客户端配置
type ClientConfig struct {
	// 连接配置
	Target string `yaml:"target" validate:"required"`

	// TLS配置
	TLS TLSConfig `yaml:"tls"`

	// 超时配置
	Timeout TimeoutConfig `yaml:"timeout"`

	// 重试配置
	Retry RetryConfig `yaml:"retry"`

	// 负载均衡配置
	LoadBalancer LoadBalancerConfig `yaml:"load_balancer"`

	// 连接池配置
	ConnectionPool ConnectionPoolConfig `yaml:"connection_pool"`

	// 中间件配置
	Middleware MiddlewareConfig `yaml:"middleware"`
}

// TLSConfig TLS配置
type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
	CAFile   string `yaml:"ca_file"`
	Insecure bool   `yaml:"insecure"`
}

// TimeoutConfig 超时配置
type TimeoutConfig struct {
	Connection time.Duration `yaml:"connection"`
	Request    time.Duration `yaml:"request"`
	Idle       time.Duration `yaml:"idle"`
	Keepalive  time.Duration `yaml:"keepalive"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled bool    `yaml:"enabled"`
	RPS     float64 `yaml:"rps"`
	Burst   int     `yaml:"burst"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	Enabled     bool          `yaml:"enabled"`
	MaxAttempts int           `yaml:"max_attempts"`
	InitialBackoff time.Duration `yaml:"initial_backoff"`
	MaxBackoff  time.Duration `yaml:"max_backoff"`
	Multiplier  float64       `yaml:"multiplier"`
}

// LoadBalancerConfig 负载均衡配置
type LoadBalancerConfig struct {
	Policy string `yaml:"policy"` // round_robin, pick_first, grpclb
}

// ConnectionPoolConfig 连接池配置
type ConnectionPoolConfig struct {
	MaxConnections int           `yaml:"max_connections"`
	MaxIdle        int           `yaml:"max_idle"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	Auth       AuthConfig       `yaml:"auth"`
	Logging    LoggingConfig    `yaml:"logging"`
	Metrics    MetricsConfig    `yaml:"metrics"`
	Tracing    TracingConfig    `yaml:"tracing"`
	Recovery   RecoveryConfig   `yaml:"recovery"`
	Validation ValidationConfig `yaml:"validation"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	Enabled   bool     `yaml:"enabled"`
	JWTSecret string   `yaml:"jwt_secret"`
	SkipPaths []string `yaml:"skip_paths"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Enabled       bool     `yaml:"enabled"`
	LogRequests   bool     `yaml:"log_requests"`
	LogResponses  bool     `yaml:"log_responses"`
	LogPayloads   bool     `yaml:"log_payloads"`
	SkipPaths     []string `yaml:"skip_paths"`
	MaxPayloadSize int     `yaml:"max_payload_size"`
}

// MetricsConfig 指标配置
type MetricsConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Namespace string `yaml:"namespace"`
	Subsystem string `yaml:"subsystem"`
}

// TracingConfig 链路追踪配置
type TracingConfig struct {
	Enabled     bool    `yaml:"enabled"`
	ServiceName string  `yaml:"service_name"`
	SampleRate  float64 `yaml:"sample_rate"`
	Endpoint    string  `yaml:"endpoint"`
}

// RecoveryConfig 恢复配置
type RecoveryConfig struct {
	Enabled bool `yaml:"enabled"`
}

// ValidationConfig 验证配置
type ValidationConfig struct {
	Enabled bool `yaml:"enabled"`
}

// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
}

// DefaultServerConfig 返回默认服务器配置
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Host: "0.0.0.0",
		Port: 8080,
		TLS: TLSConfig{
			Enabled: false,
		},
		Timeout: TimeoutConfig{
			Connection: 30 * time.Second,
			Request:    30 * time.Second,
			Idle:       60 * time.Second,
			Keepalive:  30 * time.Second,
		},
		RateLimit: RateLimitConfig{
			Enabled: false,
			RPS:     1000,
			Burst:   100,
		},
		Middleware: MiddlewareConfig{
			Auth: AuthConfig{
				Enabled: true,
			},
			Logging: LoggingConfig{
				Enabled:     true,
				LogRequests: true,
			},
			Metrics: MetricsConfig{
				Enabled: true,
			},
			Tracing: TracingConfig{
				Enabled:    true,
				SampleRate: 0.1,
			},
			Recovery: RecoveryConfig{
				Enabled: true,
			},
			Validation: ValidationConfig{
				Enabled: true,
			},
		},
		HealthCheck: HealthCheckConfig{
			Enabled:  true,
			Interval: 30 * time.Second,
			Timeout:  5 * time.Second,
		},
		Reflection: false,
	}
}

// DefaultClientConfig 返回默认客户端配置
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Target: "localhost:8080",
		TLS: TLSConfig{
			Enabled: false,
		},
		Timeout: TimeoutConfig{
			Connection: 10 * time.Second,
			Request:    30 * time.Second,
			Idle:       60 * time.Second,
			Keepalive:  30 * time.Second,
		},
		Retry: RetryConfig{
			Enabled:        true,
			MaxAttempts:    3,
			InitialBackoff: 100 * time.Millisecond,
			MaxBackoff:     5 * time.Second,
			Multiplier:     2.0,
		},
		LoadBalancer: LoadBalancerConfig{
			Policy: "round_robin",
		},
		ConnectionPool: ConnectionPoolConfig{
			MaxConnections: 100,
			MaxIdle:        10,
			IdleTimeout:    60 * time.Second,
		},
		Middleware: MiddlewareConfig{
			Auth: AuthConfig{
				Enabled: true,
			},
			Logging: LoggingConfig{
				Enabled: true,
			},
			Metrics: MetricsConfig{
				Enabled: true,
			},
			Tracing: TracingConfig{
				Enabled:    true,
				SampleRate: 0.1,
			},
			Recovery: RecoveryConfig{
				Enabled: true,
			},
		},
	}
}
```

#### 5.2 创建 gRPC 中间件

```go
// pkg/grpc/middleware/auth.go
package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"movieinfo/pkg/auth"
	"movieinfo/pkg/logger"
)

// AuthInterceptor 认证拦截器
type AuthInterceptor struct {
	auth      auth.TokenManager
	logger    logger.Logger
	skipPaths []string
}

// NewAuthInterceptor 创建认证拦截器
func NewAuthInterceptor(auth auth.TokenManager, logger logger.Logger, skipPaths []string) *AuthInterceptor {
	return &AuthInterceptor{
		auth:      auth,
		logger:    logger,
		skipPaths: skipPaths,
	}
}

// UnaryServerInterceptor 一元服务器拦截器
func (a *AuthInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 检查是否需要跳过认证
		if a.shouldSkip(info.FullMethod) {
			return handler(ctx, req)
		}

		// 从元数据中获取token
		token, err := a.extractToken(ctx)
		if err != nil {
			a.logger.Warn("Failed to extract token", "error", err, "method", info.FullMethod)
			return nil, status.Error(codes.Unauthenticated, "missing or invalid token")
		}

		// 验证token
		claims, err := a.auth.ValidateToken(token)
		if err != nil {
			a.logger.Warn("Token validation failed", "error", err, "method", info.FullMethod)
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// 将用户信息添加到上下文
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)

		return handler(ctx, req)
	}
}

// StreamServerInterceptor 流服务器拦截器
func (a *AuthInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 检查是否需要跳过认证
		if a.shouldSkip(info.FullMethod) {
			return handler(srv, stream)
		}

		// 从元数据中获取token
		token, err := a.extractToken(stream.Context())
		if err != nil {
			a.logger.Warn("Failed to extract token", "error", err, "method", info.FullMethod)
			return status.Error(codes.Unauthenticated, "missing or invalid token")
		}

		// 验证token
		claims, err := a.auth.ValidateToken(token)
		if err != nil {
			a.logger.Warn("Token validation failed", "error", err, "method", info.FullMethod)
			return status.Error(codes.Unauthenticated, "invalid token")
		}

		// 创建包装的流
		wrappedStream := &wrappedServerStream{
			ServerStream: stream,
			ctx: context.WithValue(stream.Context(), "user_id", claims.UserID),
		}

		return handler(srv, wrappedStream)
	}
}

// shouldSkip 检查是否应该跳过认证
func (a *AuthInterceptor) shouldSkip(method string) bool {
	for _, path := range a.skipPaths {
		if strings.Contains(method, path) {
			return true
		}
	}
	return false
}

// extractToken 从上下文中提取token
func (a *AuthInterceptor) extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	authorization := md.Get("authorization")
	if len(authorization) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization header")
	}

	token := authorization[0]
	if !strings.HasPrefix(token, "Bearer ") {
		return "", status.Error(codes.Unauthenticated, "invalid authorization format")
	}

	return strings.TrimPrefix(token, "Bearer "), nil
}

// wrappedServerStream 包装的服务器流
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
```

```go
// pkg/grpc/middleware/logging.go
package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"movieinfo/pkg/logger"
)

// LoggingInterceptor 日志拦截器
type LoggingInterceptor struct {
	logger         logger.Logger
	logRequests    bool
	logResponses   bool
	logPayloads    bool
	maxPayloadSize int
	skipPaths      []string
}

// NewLoggingInterceptor 创建日志拦截器
func NewLoggingInterceptor(logger logger.Logger, config LoggingConfig) *LoggingInterceptor {
	return &LoggingInterceptor{
		logger:         logger,
		logRequests:    config.LogRequests,
		logResponses:   config.LogResponses,
		logPayloads:    config.LogPayloads,
		maxPayloadSize: config.MaxPayloadSize,
		skipPaths:      config.SkipPaths,
	}
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	LogRequests    bool
	LogResponses   bool
	LogPayloads    bool
	MaxPayloadSize int
	SkipPaths      []string
}

// UnaryServerInterceptor 一元服务器拦截器
func (l *LoggingInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// 记录请求
		if l.logRequests && l.shouldLog(info.FullMethod) {
			fields := []interface{}{
				"method", info.FullMethod,
				"start_time", start,
			}

			if l.logPayloads {
				fields = append(fields, "request", l.truncatePayload(req))
			}

			l.logger.Info("gRPC request started", fields...)
		}

		// 执行处理器
		resp, err := handler(ctx, req)

		// 计算耗时
		duration := time.Since(start)

		// 记录响应
		if l.logResponses && l.shouldLog(info.FullMethod) {
			fields := []interface{}{
				"method", info.FullMethod,
				"duration", duration,
				"success", err == nil,
			}

			if err != nil {
				st := status.Convert(err)
				fields = append(fields, "error_code", st.Code(), "error_message", st.Message())
			}

			if l.logPayloads && resp != nil {
				fields = append(fields, "response", l.truncatePayload(resp))
			}

			if err != nil {
				l.logger.Error("gRPC request completed with error", fields...)
			} else {
				l.logger.Info("gRPC request completed successfully", fields...)
			}
		}

		return resp, err
	}
}

// StreamServerInterceptor 流服务器拦截器
func (l *LoggingInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		// 记录流开始
		if l.shouldLog(info.FullMethod) {
			l.logger.Info("gRPC stream started", "method", info.FullMethod, "start_time", start)
		}

		// 执行处理器
		err := handler(srv, stream)

		// 计算耗时
		duration := time.Since(start)

		// 记录流结束
		if l.shouldLog(info.FullMethod) {
			fields := []interface{}{
				"method", info.FullMethod,
				"duration", duration,
				"success", err == nil,
			}

			if err != nil {
				st := status.Convert(err)
				fields = append(fields, "error_code", st.Code(), "error_message", st.Message())
				l.logger.Error("gRPC stream completed with error", fields...)
			} else {
				l.logger.Info("gRPC stream completed successfully", fields...)
			}
		}

		return err
	}
}

// shouldLog 检查是否应该记录日志
func (l *LoggingInterceptor) shouldLog(method string) bool {
	for _, path := range l.skipPaths {
		if strings.Contains(method, path) {
			return false
		}
	}
	return true
}

// truncatePayload 截断负载
func (l *LoggingInterceptor) truncatePayload(payload interface{}) interface{} {
	if !l.logPayloads {
		return nil
	}

	// 这里可以实现更复杂的截断逻辑
	// 例如：序列化为JSON并检查大小
	return payload
}
```

```go
// pkg/grpc/middleware/recovery.go
package middleware

import (
	"context"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"movieinfo/pkg/logger"
)

// RecoveryInterceptor 恢复拦截器
type RecoveryInterceptor struct {
	logger logger.Logger
}

// NewRecoveryInterceptor 创建恢复拦截器
func NewRecoveryInterceptor(logger logger.Logger) *RecoveryInterceptor {
	return &RecoveryInterceptor{
		logger: logger,
	}
}

// UnaryServerInterceptor 一元服务器拦截器
func (r *RecoveryInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if p := recover(); p != nil {
				r.logger.Error("gRPC panic recovered",
					"method", info.FullMethod,
					"panic", p,
					"stack", string(debug.Stack()),
				)
				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

// StreamServerInterceptor 流服务器拦截器
func (r *RecoveryInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if p := recover(); p != nil {
				r.logger.Error("gRPC stream panic recovered",
					"method", info.FullMethod,
					"panic", p,
					"stack", string(debug.Stack()),
				)
				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		return handler(srv, stream)
	}
}
```

#### 5.3 创建 gRPC 服务器

```go
// pkg/grpc/server.go
package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"movieinfo/pkg/grpc/middleware"
	"movieinfo/pkg/logger"
)

// Server gRPC服务器
type Server struct {
	config     *ServerConfig
	logger     logger.Logger
	server     *grpc.Server
	listener   net.Listener
	healthSrv  *health.Server
	interceptors *Interceptors
}

// Interceptors 拦截器集合
type Interceptors struct {
	Auth     *middleware.AuthInterceptor
	Logging  *middleware.LoggingInterceptor
	Recovery *middleware.RecoveryInterceptor
	// 可以添加更多拦截器
}

// NewServer 创建新的gRPC服务器
func NewServer(config *ServerConfig, logger logger.Logger) (*Server, error) {
	// 创建监听器
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	// 创建服务器选项
	opts := []grpc.ServerOption{
		// Keepalive配置
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    config.Timeout.Keepalive,
			Timeout: config.Timeout.Connection,
		}),
		// 连接超时
		grpc.ConnectionTimeout(config.Timeout.Connection),
	}

	// 创建拦截器
	interceptors := &Interceptors{
		Recovery: middleware.NewRecoveryInterceptor(logger),
		Logging:  middleware.NewLoggingInterceptor(logger, middleware.LoggingConfig{
			LogRequests:    config.Middleware.Logging.LogRequests,
			LogResponses:   config.Middleware.Logging.LogResponses,
			LogPayloads:    config.Middleware.Logging.LogPayloads,
			MaxPayloadSize: config.Middleware.Logging.MaxPayloadSize,
			SkipPaths:      config.Middleware.Logging.SkipPaths,
		}),
	}

	// 添加一元拦截器
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		interceptors.Recovery.UnaryServerInterceptor(),
		interceptors.Logging.UnaryServerInterceptor(),
	}

	// 添加流拦截器
	streamInterceptors := []grpc.StreamServerInterceptor{
		interceptors.Recovery.StreamServerInterceptor(),
		interceptors.Logging.StreamServerInterceptor(),
	}

	// 如果启用认证，添加认证拦截器
	if config.Middleware.Auth.Enabled {
		// 这里需要实际的认证管理器实现
		// authManager := auth.NewTokenManager(config.Middleware.Auth.JWTSecret)
		// interceptors.Auth = middleware.NewAuthInterceptor(authManager, logger, config.Middleware.Auth.SkipPaths)
		// unaryInterceptors = append(unaryInterceptors, interceptors.Auth.UnaryServerInterceptor())
		// streamInterceptors = append(streamInterceptors, interceptors.Auth.StreamServerInterceptor())
	}

	// 添加拦截器到选项
	opts = append(opts,
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	// 创建gRPC服务器
	grpcServer := grpc.NewServer(opts...)

	// 创建健康检查服务
	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)

	// 如果启用反射，注册反射服务
	if config.Reflection {
		reflection.Register(grpcServer)
	}

	return &Server{
		config:       config,
		logger:       logger,
		server:       grpcServer,
		listener:     listener,
		healthSrv:    healthSrv,
		interceptors: interceptors,
	}, nil
}

// RegisterService 注册服务
func (s *Server) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.server.RegisterService(desc, impl)
}

// Start 启动服务器
func (s *Server) Start() error {
	s.logger.Info("Starting gRPC server",
		"address", s.listener.Addr().String(),
		"reflection", s.config.Reflection,
	)

	// 设置健康检查状态
	s.healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	return s.server.Serve(s.listener)
}

// Stop 停止服务器
func (s *Server) Stop() {
	s.logger.Info("Stopping gRPC server")

	// 设置健康检查状态
	s.healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	// 优雅停止
	s.server.GracefulStop()
}

// ForceStop 强制停止服务器
func (s *Server) ForceStop() {
	s.logger.Info("Force stopping gRPC server")
	s.server.Stop()
}

// GetServer 获取gRPC服务器实例
func (s *Server) GetServer() *grpc.Server {
	return s.server
}

// GetAddress 获取服务器地址
func (s *Server) GetAddress() string {
	return s.listener.Addr().String()
}

// SetHealthStatus 设置健康检查状态
func (s *Server) SetHealthStatus(service string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	s.healthSrv.SetServingStatus(service, status)
}
```

#### 5.4 创建 gRPC 客户端

```go
// pkg/grpc/client.go
package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"movieinfo/pkg/logger"
)

// Client gRPC客户端
type Client struct {
	config *ClientConfig
	logger logger.Logger
	conn   *grpc.ClientConn
}

// NewClient 创建新的gRPC客户端
func NewClient(config *ClientConfig, logger logger.Logger) (*Client, error) {
	// 创建连接选项
	opts := []grpc.DialOption{
		// 使用不安全连接（开发环境）
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// Keepalive配置
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                config.Timeout.Keepalive,
			Timeout:             config.Timeout.Connection,
			PermitWithoutStream: true,
		}),
		// 连接超时
		grpc.WithTimeout(config.Timeout.Connection),
		// 负载均衡
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, config.LoadBalancer.Policy)),
	}

	// 创建连接
	conn, err := grpc.Dial(config.Target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return &Client{
		config: config,
		logger: logger,
		conn:   conn,
	}, nil
}

// GetConnection 获取连接
func (c *Client) GetConnection() *grpc.ClientConn {
	return c.conn
}

// Close 关闭连接
func (c *Client) Close() error {
	c.logger.Info("Closing gRPC client connection")
	return c.conn.Close()
}

// WithTimeout 创建带超时的上下文
func (c *Client) WithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// WithRequestTimeout 创建带请求超时的上下文
func (c *Client) WithRequestTimeout() (context.Context, context.CancelFunc) {
	return c.WithTimeout(c.config.Timeout.Request)
}

// IsConnected 检查连接状态
func (c *Client) IsConnected() bool {
	state := c.conn.GetState()
	return state == grpc.Ready
}

// WaitForReady 等待连接就绪
func (c *Client) WaitForReady(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return c.conn.WaitForStateChange(ctx, c.conn.GetState())
}
```

### 步骤6：创建 Makefile 和脚本

#### 6.1 创建 Makefile

```makefile
# Makefile
.PHONY: proto proto-clean proto-gen proto-lint proto-format

# Protocol Buffers相关
PROTO_DIR := api/proto
PROTO_OUT_DIR := api/proto
PROTO_FILES := $(shell find $(PROTO_DIR) -name "*.proto")

# 工具
PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc
BUF := buf

# 安装依赖
install-deps:
	@echo "Installing dependencies..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/bufbuild/buf/cmd/buf@latest

# 生成Protocol Buffers代码
proto-gen:
	@echo "Generating Protocol Buffers code..."
	@for file in $(PROTO_FILES); do \
		echo "Generating $$file..."; \
		$(PROTOC) \
			--proto_path=$(PROTO_DIR) \
			--go_out=$(PROTO_OUT_DIR) \
			--go_opt=paths=source_relative \
			--go-grpc_out=$(PROTO_OUT_DIR) \
			--go-grpc_opt=paths=source_relative \
			$$file; \
	done

# 清理生成的文件
proto-clean:
	@echo "Cleaning generated Protocol Buffers files..."
	find $(PROTO_OUT_DIR) -name "*.pb.go" -delete
	find $(PROTO_OUT_DIR) -name "*_grpc.pb.go" -delete

# 格式化proto文件
proto-format:
	@echo "Formatting Protocol Buffers files..."
	$(BUF) format -w

# 检查proto文件
proto-lint:
	@echo "Linting Protocol Buffers files..."
	$(BUF) lint

# 生成并检查
proto: proto-clean proto-format proto-lint proto-gen
	@echo "Protocol Buffers generation completed!"

# 帮助
help:
	@echo "Available targets:"
	@echo "  install-deps  - Install required dependencies"
	@echo "  proto-gen     - Generate Protocol Buffers code"
	@echo "  proto-clean   - Clean generated files"
	@echo "  proto-format  - Format proto files"
	@echo "  proto-lint    - Lint proto files"
	@echo "  proto         - Full proto workflow (clean, format, lint, generate)"
	@echo "  help          - Show this help message"
```

#### 6.2 创建 buf 配置

```yaml
# buf.yaml
version: v1
deps:
  - buf.build/googleapis/googleapis
lint:
  use:
    - DEFAULT
  except:
    - UNARY_RPC
    - COMMENT_FIELD
    - SERVICE_SUFFIX
  ignore_only:
    PACKAGE_DIRECTORY_MATCH:
      - api/proto/common
breaking:
  use:
    - FILE
  except:
    - EXTENSION_NO_DELETE
    - FIELD_SAME_DEFAULT
```

```yaml
# buf.gen.yaml
version: v1
plugins:
  - plugin: buf.build/protocolbuffers/go
    out: .
    opt:
      - paths=source_relative
  - plugin: buf.build/grpc/go
    out: .
    opt:
      - paths=source_relative
```

### 步骤7：创建配置文件

#### 7.1 创建 gRPC 配置文件

```yaml
# configs/grpc.yaml
# gRPC服务配置
services:
  user:
    server:
      host: "0.0.0.0"
      port: 8081
      tls:
        enabled: false
      timeout:
        connection: 30s
        request: 30s
        idle: 60s
        keepalive: 30s
      middleware:
        auth:
          enabled: true
          jwt_secret: "your-jwt-secret-key"
          skip_paths:
            - "/movieinfo.user.UserService/Login"
            - "/movieinfo.user.UserService/RefreshToken"
            - "/grpc.health.v1.Health/Check"
        logging:
          enabled: true
          log_requests: true
          log_responses: true
          log_payloads: false
          max_payload_size: 1024
          skip_paths:
            - "/grpc.health.v1.Health/Check"
        metrics:
          enabled: true
          namespace: "movieinfo"
          subsystem: "user_service"
        tracing:
          enabled: true
          service_name: "user-service"
          sample_rate: 0.1
          endpoint: "http://localhost:14268/api/traces"
        recovery:
          enabled: true
        validation:
          enabled: true
      health_check:
        enabled: true
        interval: 30s
        timeout: 5s
      reflection: false

  movie:
    server:
      host: "0.0.0.0"
      port: 8082
      tls:
        enabled: false
      timeout:
        connection: 30s
        request: 30s
        idle: 60s
        keepalive: 30s
      middleware:
        auth:
          enabled: true
          jwt_secret: "your-jwt-secret-key"
          skip_paths:
            - "/movieinfo.movie.MovieService/ListMovies"
            - "/movieinfo.movie.MovieService/GetMovie"
            - "/movieinfo.movie.MovieService/SearchMovies"
            - "/movieinfo.movie.MovieService/GetPopularMovies"
            - "/grpc.health.v1.Health/Check"
        logging:
          enabled: true
          log_requests: true
          log_responses: true
          log_payloads: false
          max_payload_size: 1024
          skip_paths:
            - "/grpc.health.v1.Health/Check"
        metrics:
          enabled: true
          namespace: "movieinfo"
          subsystem: "movie_service"
        tracing:
          enabled: true
          service_name: "movie-service"
          sample_rate: 0.1
          endpoint: "http://localhost:14268/api/traces"
        recovery:
          enabled: true
        validation:
          enabled: true
      health_check:
        enabled: true
        interval: 30s
        timeout: 5s
      reflection: false

  rating:
    server:
      host: "0.0.0.0"
      port: 8083
      tls:
        enabled: false
      timeout:
        connection: 30s
        request: 30s
        idle: 60s
        keepalive: 30s
      middleware:
        auth:
          enabled: true
          jwt_secret: "your-jwt-secret-key"
          skip_paths:
            - "/movieinfo.rating.RatingService/GetMovieRatings"
            - "/movieinfo.rating.RatingService/GetRatingStatistics"
            - "/grpc.health.v1.Health/Check"
        logging:
          enabled: true
          log_requests: true
          log_responses: true
          log_payloads: false
          max_payload_size: 1024
          skip_paths:
            - "/grpc.health.v1.Health/Check"
        metrics:
          enabled: true
          namespace: "movieinfo"
          subsystem: "rating_service"
        tracing:
          enabled: true
          service_name: "rating-service"
          sample_rate: 0.1
          endpoint: "http://localhost:14268/api/traces"
        recovery:
          enabled: true
        validation:
          enabled: true
      health_check:
        enabled: true
        interval: 30s
        timeout: 5s
      reflection: false

# 客户端配置
clients:
  user:
    target: "localhost:8081"
    tls:
      enabled: false
    timeout:
      connection: 10s
      request: 30s
      idle: 60s
      keepalive: 30s
    retry:
      enabled: true
      max_attempts: 3
      initial_backoff: 100ms
      max_backoff: 5s
      multiplier: 2.0
    load_balancer:
      policy: "round_robin"
    connection_pool:
      max_connections: 100
      max_idle: 10
      idle_timeout: 60s

  movie:
    target: "localhost:8082"
    tls:
      enabled: false
    timeout:
      connection: 10s
      request: 30s
      idle: 60s
      keepalive: 30s
    retry:
      enabled: true
      max_attempts: 3
      initial_backoff: 100ms
      max_backoff: 5s
      multiplier: 2.0
    load_balancer:
      policy: "round_robin"
    connection_pool:
      max_connections: 100
      max_idle: 10
      idle_timeout: 60s

  rating:
    target: "localhost:8083"
    tls:
      enabled: false
    timeout:
      connection: 10s
      request: 30s
      idle: 60s
      keepalive: 30s
    retry:
      enabled: true
      max_attempts: 3
      initial_backoff: 100ms
      max_backoff: 5s
      multiplier: 2.0
    load_balancer:
      policy: "round_robin"
    connection_pool:
      max_connections: 100
      max_idle: 10
      idle_timeout: 60s
```

## 测试验证

### 单元测试

```go
// pkg/grpc/server_test.go
package grpc

import (
	"testing"

	"movieinfo/pkg/logger"
)

func TestNewServer(t *testing.T) {
	config := DefaultServerConfig()
	config.Port = 0 // 使用随机端口
	
	logger := logger.NewConsoleLogger()
	
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	
	if server == nil {
		t.Fatal("Server is nil")
	}
	
	// 测试获取地址
	addr := server.GetAddress()
	if addr == "" {
		t.Fatal("Server address is empty")
	}
	
	// 清理
	server.ForceStop()
}
```

### 集成测试

```go
// test/grpc_integration_test.go
package test

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	grpcpkg "movieinfo/pkg/grpc"
	"movieinfo/pkg/logger"
)

func TestGRPCServerHealthCheck(t *testing.T) {
	// 创建服务器
	config := grpcpkg.DefaultServerConfig()
	config.Port = 0 // 使用随机端口
	
	logger := logger.NewConsoleLogger()
	
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	
	// 启动服务器
	go func() {
		if err := server.Start(); err != nil {
			t.Errorf("Server failed to start: %v", err)
		}
	}()
	
	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)
	
	defer server.ForceStop()
	
	// 创建客户端
	clientConfig := grpcpkg.DefaultClientConfig()
	clientConfig.Target = server.GetAddress()
	
	client, err := grpcpkg.NewClient(clientConfig, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	
	// 测试健康检查
	healthClient := grpc_health_v1.NewHealthClient(client.GetConnection())
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
	}
	
	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Fatalf("Expected SERVING status, got %v", resp.Status)
	}
}
```

## 预期结果

完成本步骤后，你将拥有：

1. **完整的 gRPC 协议定义**
   - 通用数据类型和错误处理
   - 用户服务的完整接口定义
   - 电影服务的完整接口定义
   - 评分服务的完整接口定义

2. **gRPC 基础设施**
   - 可配置的 gRPC 服务器
   - 可配置的 gRPC 客户端
   - 完整的中间件系统
   - 健康检查和监控支持

3. **开发工具**
   - Makefile 自动化脚本
   - buf 配置和代码生成
   - 完整的配置文件

4. **测试覆盖**
   - 单元测试
   - 集成测试
   - 健康检查测试

## 注意事项

### 安全性
- **认证和授权**：确保所有敏感接口都有适当的认证
- **输入验证**：对所有输入进行严格验证
- **错误处理**：不要在错误信息中泄露敏感信息
- **TLS 加密**：生产环境必须启用 TLS

### 性能
- **连接池**：合理配置客户端连接池
- **超时设置**：设置合适的超时时间
- **负载均衡**：使用适当的负载均衡策略
- **监控指标**：收集关键性能指标

### 可维护性
- **版本兼容**：确保协议的向后兼容性
- **文档更新**：及时更新接口文档
- **代码生成**：使用自动化工具生成代码
- **测试覆盖**：保持高测试覆盖率

### 扩展性
- **服务发现**：考虑集成服务发现机制
- **配置中心**：支持动态配置更新
- **链路追踪**：集成分布式链路追踪
- **指标监控**：集成 Prometheus 等监控系统

## 下一步骤

完成 gRPC 协议定义后，下一步将进行：

1. **数据模型层开发**（06-数据模型层.md）
   - 定义数据模型结构
   - 实现数据访问层
   - 集成数据库连接
   - 实现缓存机制

2. **检查清单**
   - [ ] 所有 proto 文件编译成功
   - [ ] gRPC 服务器可以正常启动
   - [ ] gRPC 客户端可以正常连接
   - [ ] 健康检查接口正常工作
   - [ ] 中间件功能正常
   - [ ] 配置文件格式正确
   - [ ] 单元测试通过
   - [ ] 集成测试通过
   - [ ] 代码生成脚本正常工作
   - [ ] 文档更新完成

请确保所有检查项都已完成，然后继续下一个开发步骤。
```