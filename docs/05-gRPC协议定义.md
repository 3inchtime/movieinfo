# 05-gRPC协议定义（简化版）

## 目标说明

本步骤的目标是为 MovieInfo 项目定义**简化但完整**的 gRPC 协议，建立服务间通信的标准接口。我们采用**渐进式开发策略**，先实现核心功能，后续根据需要逐步完善。

### 核心目标
- 设计用户服务的基础 gRPC 接口（CRUD + 认证）
- 设计电影服务的基础 gRPC 接口（CRUD + 搜索）
- 设计评分服务的基础 gRPC 接口（CRUD + 统计）
- 定义通用的数据类型和简化的错误处理
- 实现基础的 gRPC 服务器和客户端
- 集成日志和基础监控

### 设计原则
- **简单优先**：优先实现核心功能，避免过度设计
- **渐进式开发**：从MVP开始，逐步添加高级特性
- **实用性**：每个功能都有明确的业务价值
- **可扩展性**：为未来扩展预留接口

## 前置条件

- 已完成项目初始化（01-项目初始化.md）
- 已完成数据库设计（02-数据库设计.md）
- 已完成配置管理系统（03-配置管理系统.md）
- 已完成日志系统（04-日志系统.md）
- 基本了解 Protocol Buffers 语法
- 了解 gRPC 的基本概念

## 技术要点

### 简化的设计原则
- **接口优先**：先定义接口，再实现服务
- **类型安全**：利用强类型系统避免运行时错误
- **统一响应**：使用统一的响应格式简化错误处理
- **核心功能**：只实现必要的业务功能

### 技术选型
- **Protocol Buffers**：使用 proto3 语法定义接口
- **gRPC-Go**：使用官方 Go 实现
- **简化中间件**：只实现日志和恢复中间件
- **统一错误处理**：使用 gRPC 状态码和简单错误信息

## 实现步骤

### 步骤1：定义通用数据类型

#### 1.1 创建简化的通用类型定义

我们首先定义最基础的通用类型，避免过度复杂化：

```protobuf
// proto/common/common.proto
syntax = "proto3";

package movieinfo.common;

option go_package = "github.com/3inchtime/movieinfo/proto/common";

import "google/protobuf/timestamp.proto";

// 分页请求 - 简化版本，只包含必要字段
message PageRequest {
  int32 page = 1;           // 页码，从1开始
  int32 page_size = 2;      // 每页大小，默认10，最大100
}

// 分页响应 - 简化版本
message PageResponse {
  int32 page = 1;           // 当前页码
  int32 page_size = 2;      // 每页大小
  int64 total = 3;          // 总记录数
  int32 total_pages = 4;    // 总页数
}

// 通用响应 - 大幅简化，只保留核心信息
message CommonResponse {
  bool success = 1;         // 是否成功
  string message = 2;       // 响应消息
  google.protobuf.Timestamp timestamp = 3; // 响应时间
}

// 健康检查 - 保持标准格式
message HealthCheckRequest {
  string service = 1;       // 服务名称
}

message HealthCheckResponse {
  enum ServingStatus {
    UNKNOWN = 0;
    SERVING = 1;
    NOT_SERVING = 2;
  }
  ServingStatus status = 1;
  string message = 2;
}
```

**设计说明**：
- 移除了复杂的批量操作、统计信息等高级功能
- 简化了分页请求，移除了排序功能（可在业务层实现）
- 统一响应格式，使用简单的成功/失败标识
- 保留健康检查，这是微服务的基础功能

#### 1.2 创建简化的错误定义

```protobuf
// proto/common/error.proto
syntax = "proto3";

package movieinfo.common;

option go_package = "github.com/3inchtime/movieinfo/proto/common";

// 简化的错误代码 - 只保留最常用的错误类型
enum ErrorCode {
  // 通用错误
  UNKNOWN_ERROR = 0;
  INVALID_ARGUMENT = 1;
  NOT_FOUND = 2;
  ALREADY_EXISTS = 3;
  INTERNAL_ERROR = 4;
  
  // 认证相关错误
  UNAUTHENTICATED = 100;
  PERMISSION_DENIED = 101;
  
  // 业务逻辑错误
  BUSINESS_ERROR = 200;
}

// 简化的错误详情
message ErrorDetail {
  ErrorCode code = 1;       // 错误代码
  string message = 2;       // 错误消息
  string field = 3;         // 相关字段（可选）
}
```

**设计说明**：
- 将30+种错误类型简化为7种核心错误类型
- 移除了复杂的验证错误和业务错误结构
- 保留了最基础的错误信息，满足大部分使用场景

### 步骤2：定义用户服务接口

#### 2.1 创建简化的用户数据类型

```protobuf
// proto/user/user.proto
syntax = "proto3";

package movieinfo.user;

option go_package = "github.com/3inchtime/movieinfo/proto/user";

import "google/protobuf/timestamp.proto";
import "common/common.proto";

// 用户信息 - 简化版本，只保留核心字段
message User {
  int64 id = 1;             // 用户ID
  string username = 2;      // 用户名
  string email = 3;         // 邮箱
  string nickname = 4;      // 昵称
  string avatar = 5;        // 头像URL
  UserStatus status = 6;    // 用户状态
  google.protobuf.Timestamp created_at = 7;  // 创建时间
  google.protobuf.Timestamp updated_at = 8;  // 更新时间
}

// 简化的用户状态
enum UserStatus {
  USER_STATUS_UNKNOWN = 0;
  USER_STATUS_ACTIVE = 1;   // 活跃
  USER_STATUS_INACTIVE = 2; // 非活跃
}

// 创建用户请求 - 只保留必要字段
message CreateUserRequest {
  string username = 1;      // 用户名
  string email = 2;         // 邮箱
  string password = 3;      // 密码
  string nickname = 4;      // 昵称
}

message CreateUserResponse {
  movieinfo.common.CommonResponse common = 1;
  User user = 2;            // 创建的用户信息
}

// 获取用户请求 - 简化标识符
message GetUserRequest {
  oneof identifier {
    int64 id = 1;           // 用户ID
    string username = 2;    // 用户名
  }
}

message GetUserResponse {
  movieinfo.common.CommonResponse common = 1;
  User user = 2;            // 用户信息
}

// 更新用户请求 - 简化字段
message UpdateUserRequest {
  int64 id = 1;             // 用户ID
  string nickname = 2;      // 昵称
  string avatar = 3;        // 头像URL
}

message UpdateUserResponse {
  movieinfo.common.CommonResponse common = 1;
  User user = 2;            // 更新后的用户信息
}

// 删除用户请求
message DeleteUserRequest {
  int64 id = 1;             // 用户ID
}

message DeleteUserResponse {
  movieinfo.common.CommonResponse common = 1;
}

// 列出用户请求 - 简化过滤条件
message ListUsersRequest {
  movieinfo.common.PageRequest page = 1; // 分页参数
  UserStatus status = 2;    // 用户状态过滤（可选）
  string search = 3;        // 搜索关键词（可选）
}

message ListUsersResponse {
  movieinfo.common.CommonResponse common = 1;
  repeated User users = 2;  // 用户列表
  movieinfo.common.PageResponse page = 3; // 分页信息
}

// 用户登录请求
message LoginRequest {
  string username = 1;      // 用户名或邮箱
  string password = 2;      // 密码
}

message LoginResponse {
  movieinfo.common.CommonResponse common = 1;
  string access_token = 2;  // 访问令牌
  int64 expires_in = 3;     // 过期时间（秒）
  User user = 4;            // 用户信息
}

// 用户登出请求
message LogoutRequest {
  string access_token = 1;  // 访问令牌
}

message LogoutResponse {
  movieinfo.common.CommonResponse common = 1;
}

// 修改密码请求
message ChangePasswordRequest {
  int64 user_id = 1;        // 用户ID
  string old_password = 2;  // 旧密码
  string new_password = 3;  // 新密码
}

message ChangePasswordResponse {
  movieinfo.common.CommonResponse common = 1;
}
```

**设计说明**：
- 移除了复杂的用户资料、性别、偏好等字段，专注核心功能
- 简化了用户状态，只保留活跃和非活跃两种状态
- 移除了刷新令牌、重置密码等高级功能（可后续添加）
- 简化了登录请求，统一使用用户名或邮箱登录

#### 2.2 创建简化的用户服务接口

```protobuf
// proto/user/user_service.proto
syntax = "proto3";

package movieinfo.user;

option go_package = "github.com/3inchtime/movieinfo/proto/user";

import "user/user.proto";
import "common/common.proto";

// 用户服务 - 简化版本，只保留核心功能
service UserService {
  // 基础用户管理
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  
  // 认证功能
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc Logout(LogoutRequest) returns (LogoutResponse);
  rpc ChangePassword(ChangePasswordRequest) returns (ChangePasswordResponse);
  
  // 健康检查
  rpc HealthCheck(movieinfo.common.HealthCheckRequest) returns (movieinfo.common.HealthCheckResponse);
}
```

**设计说明**：
- 从原来的13个RPC方法简化为9个核心方法
- 移除了批量操作、统计信息、版本信息等高级功能
- 保留了完整的CRUD操作和基础认证功能
- 保留健康检查，这是微服务的标准功能

### 步骤3：定义电影服务接口

#### 3.1 创建简化的电影数据类型

```protobuf
// proto/movie/movie.proto
syntax = "proto3";

package movieinfo.movie;

option go_package = "github.com/3inchtime/movieinfo/proto/movie";

import "google/protobuf/timestamp.proto";
import "common/common.proto";

// 电影信息 - 简化版本
message Movie {
  int64 id = 1;             // 电影ID
  string title = 2;         // 电影标题
  string description = 3;   // 电影描述
  string poster_url = 4;    // 海报URL
  int32 duration = 5;       // 时长（分钟）
  google.protobuf.Timestamp release_date = 6; // 上映日期
  string language = 7;      // 语言
  repeated string genres = 8; // 类型列表
  repeated string directors = 9; // 导演列表
  repeated string actors = 10;   // 演员列表
  double average_rating = 11;    // 平均评分
  int64 rating_count = 12;       // 评分数量
  google.protobuf.Timestamp created_at = 13; // 创建时间
  google.protobuf.Timestamp updated_at = 14; // 更新时间
}

// 创建电影请求
message CreateMovieRequest {
  string title = 1;         // 电影标题
  string description = 2;   // 电影描述
  string poster_url = 3;    // 海报URL
  int32 duration = 4;       // 时长（分钟）
  google.protobuf.Timestamp release_date = 5; // 上映日期
  string language = 6;      // 语言
  repeated string genres = 7;    // 类型列表
  repeated string directors = 8; // 导演列表
  repeated string actors = 9;    // 演员列表
}

message CreateMovieResponse {
  movieinfo.common.CommonResponse common = 1;
  Movie movie = 2;          // 创建的电影信息
}

// 获取电影请求
message GetMovieRequest {
  int64 id = 1;             // 电影ID
}

message GetMovieResponse {
  movieinfo.common.CommonResponse common = 1;
  Movie movie = 2;          // 电影信息
}

// 更新电影请求
message UpdateMovieRequest {
  int64 id = 1;             // 电影ID
  string title = 2;         // 电影标题
  string description = 3;   // 电影描述
  string poster_url = 4;    // 海报URL
  int32 duration = 5;       // 时长（分钟）
  google.protobuf.Timestamp release_date = 6; // 上映日期
  string language = 7;      // 语言
  repeated string genres = 8;    // 类型列表
  repeated string directors = 9; // 导演列表
  repeated string actors = 10;   // 演员列表
}

message UpdateMovieResponse {
  movieinfo.common.CommonResponse common = 1;
  Movie movie = 2;          // 更新后的电影信息
}

// 删除电影请求
message DeleteMovieRequest {
  int64 id = 1;             // 电影ID
}

message DeleteMovieResponse {
  movieinfo.common.CommonResponse common = 1;
}

// 列出电影请求
message ListMoviesRequest {
  movieinfo.common.PageRequest page = 1; // 分页参数
  repeated string genres = 2;            // 类型过滤
  string search = 3;                     // 搜索关键词
  string language = 4;                   // 语言过滤
}

message ListMoviesResponse {
  movieinfo.common.CommonResponse common = 1;
  repeated Movie movies = 2; // 电影列表
  movieinfo.common.PageResponse page = 3; // 分页信息
}

// 搜索电影请求
message SearchMoviesRequest {
  string query = 1;         // 搜索查询
  movieinfo.common.PageRequest page = 2; // 分页参数
}

message SearchMoviesResponse {
  movieinfo.common.CommonResponse common = 1;
  repeated Movie movies = 2; // 电影列表
  movieinfo.common.PageResponse page = 3; // 分页信息
}
```

**设计说明**：
- 移除了复杂的电影状态、类别管理、元数据等功能
- 将评分信息直接嵌入电影对象，简化数据结构
- 移除了原始标题、预告片URL、国家等非核心字段
- 简化了搜索功能，移除了复杂的过滤条件和搜索建议

#### 3.2 创建简化的电影服务接口

```protobuf
// proto/movie/movie_service.proto
syntax = "proto3";

package movieinfo.movie;

option go_package = "github.com/3inchtime/movieinfo/proto/movie";

import "movie/movie.proto";
import "common/common.proto";

// 电影服务 - 简化版本
service MovieService {
  // 基础电影管理
  rpc CreateMovie(CreateMovieRequest) returns (CreateMovieResponse);
  rpc GetMovie(GetMovieRequest) returns (GetMovieResponse);
  rpc UpdateMovie(UpdateMovieRequest) returns (UpdateMovieResponse);
  rpc DeleteMovie(DeleteMovieRequest) returns (DeleteMovieResponse);
  rpc ListMovies(ListMoviesRequest) returns (ListMoviesResponse);
  
  // 搜索功能
  rpc SearchMovies(SearchMoviesRequest) returns (SearchMoviesResponse);
  
  // 健康检查
  rpc HealthCheck(movieinfo.common.HealthCheckRequest) returns (movieinfo.common.HealthCheckResponse);
}
```

**设计说明**：
- 从原来的15+个RPC方法简化为7个核心方法
- 移除了批量操作、热门电影、推荐算法等高级功能
- 保留了完整的CRUD操作和基础搜索功能

### 步骤4：定义评分服务接口

#### 4.1 创建简化的评分数据类型

```protobuf
// proto/rating/rating.proto
syntax = "proto3";

package movieinfo.rating;

option go_package = "github.com/3inchtime/movieinfo/proto/rating";

import "google/protobuf/timestamp.proto";
import "common/common.proto";

// 评分信息 - 简化版本
message Rating {
  int64 id = 1;             // 评分ID
  int64 user_id = 2;        // 用户ID
  int64 movie_id = 3;       // 电影ID
  int32 rating = 4;         // 评分值（1-5）
  string comment = 5;       // 评论内容
  google.protobuf.Timestamp created_at = 6; // 创建时间
  google.protobuf.Timestamp updated_at = 7; // 更新时间
}

// 创建评分请求
message CreateRatingRequest {
  int64 user_id = 1;        // 用户ID
  int64 movie_id = 2;       // 电影ID
  int32 rating = 3;         // 评分值（1-5）
  string comment = 4;       // 评论内容
}

message CreateRatingResponse {
  movieinfo.common.CommonResponse common = 1;
  Rating rating = 2;        // 创建的评分信息
}

// 获取评分请求
message GetRatingRequest {
  int64 id = 1;             // 评分ID
}

message GetRatingResponse {
  movieinfo.common.CommonResponse common = 1;
  Rating rating = 2;        // 评分信息
}

// 更新评分请求
message UpdateRatingRequest {
  int64 id = 1;             // 评分ID
  int32 rating = 2;         // 评分值（1-5）
  string comment = 3;       // 评论内容
}

message UpdateRatingResponse {
  movieinfo.common.CommonResponse common = 1;
  Rating rating = 2;        // 更新后的评分信息
}

// 删除评分请求
message DeleteRatingRequest {
  int64 id = 1;             // 评分ID
}

message DeleteRatingResponse {
  movieinfo.common.CommonResponse common = 1;
}

// 列出评分请求
message ListRatingsRequest {
  movieinfo.common.PageRequest page = 1; // 分页参数
  int64 user_id = 2;        // 用户ID过滤（可选）
  int64 movie_id = 3;       // 电影ID过滤（可选）
}

message ListRatingsResponse {
  movieinfo.common.CommonResponse common = 1;
  repeated Rating ratings = 2; // 评分列表
  movieinfo.common.PageResponse page = 3; // 分页信息
}

// 获取电影评分统计请求
message GetMovieRatingStatsRequest {
  int64 movie_id = 1;       // 电影ID
}

message GetMovieRatingStatsResponse {
  movieinfo.common.CommonResponse common = 1;
  double average_rating = 2; // 平均评分
  int64 total_ratings = 3;   // 总评分数
  repeated RatingCount rating_distribution = 4; // 评分分布
}

// 评分分布统计
message RatingCount {
  int32 rating = 1;         // 评分值（1-5）
  int64 count = 2;          // 该评分的数量
}
```

**设计说明**：
- 移除了复杂的评分状态、有用性标记、举报功能等
- 简化了评分统计，只保留基础的平均分和分布信息
- 专注于核心的评分CRUD功能

#### 4.2 创建简化的评分服务接口

```protobuf
// proto/rating/rating_service.proto
syntax = "proto3";

package movieinfo.rating;

option go_package = "github.com/3inchtime/movieinfo/proto/rating";

import "rating/rating.proto";
import "common/common.proto";

// 评分服务 - 简化版本
service RatingService {
  // 基础评分管理
  rpc CreateRating(CreateRatingRequest) returns (CreateRatingResponse);
  rpc GetRating(GetRatingRequest) returns (GetRatingResponse);
  rpc UpdateRating(UpdateRatingRequest) returns (UpdateRatingResponse);
  rpc DeleteRating(DeleteRatingRequest) returns (DeleteRatingResponse);
  rpc ListRatings(ListRatingsRequest) returns (ListRatingsResponse);
  
  // 统计功能
  rpc GetMovieRatingStats(GetMovieRatingStatsRequest) returns (GetMovieRatingStatsResponse);
  
  // 健康检查
  rpc HealthCheck(movieinfo.common.HealthCheckRequest) returns (movieinfo.common.HealthCheckResponse);
}
```

**设计说明**：
- 从原来的10+个RPC方法简化为7个核心方法
- 移除了批量操作、交互功能等高级特性
- 保留了完整的CRUD操作和基础统计功能

### 步骤5：创建简化的 gRPC 基础设施

#### 5.1 创建简化的 gRPC 配置

```go
// pkg/grpc/config.go
package grpc

import (
	"time"
)

// ServerConfig gRPC服务器配置 - 简化版本
type ServerConfig struct {
	// 基础配置
	Host string `yaml:"host" validate:"required"`
	Port int    `yaml:"port" validate:"required,min=1,max=65535"`

	// 超时配置
	Timeout TimeoutConfig `yaml:"timeout"`

	// 中间件配置
	Middleware MiddlewareConfig `yaml:"middleware"`

	// 健康检查配置
	HealthCheck HealthCheckConfig `yaml:"health_check"`
}

// ClientConfig gRPC客户端配置 - 简化版本
type ClientConfig struct {
	// 连接配置
	Target string `yaml:"target" validate:"required"`

	// 超时配置
	Timeout TimeoutConfig `yaml:"timeout"`

	// 重试配置
	Retry RetryConfig `yaml:"retry"`
}

// TimeoutConfig 超时配置
type TimeoutConfig struct {
	Connection time.Duration `yaml:"connection"`
	Request    time.Duration `yaml:"request"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	Enabled     bool          `yaml:"enabled"`
	MaxAttempts int           `yaml:"max_attempts"`
	Backoff     time.Duration `yaml:"backoff"`
}

// MiddlewareConfig 中间件配置 - 只保留基础中间件
type MiddlewareConfig struct {
	Logging  LoggingConfig  `yaml:"logging"`
	Recovery RecoveryConfig `yaml:"recovery"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Enabled bool `yaml:"enabled"`
}

// RecoveryConfig 恢复配置
type RecoveryConfig struct {
	Enabled bool `yaml:"enabled"`
}

// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
	Enabled bool `yaml:"enabled"`
}

// DefaultServerConfig 返回默认服务器配置
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Host: "0.0.0.0",
		Port: 8080,
		Timeout: TimeoutConfig{
			Connection: 30 * time.Second,
			Request:    30 * time.Second,
		},
		Middleware: MiddlewareConfig{
			Logging: LoggingConfig{
				Enabled: true,
			},
			Recovery: RecoveryConfig{
				Enabled: true,
			},
		},
		HealthCheck: HealthCheckConfig{
			Enabled: true,
		},
	}
}

// DefaultClientConfig 返回默认客户端配置
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Target: "localhost:8080",
		Timeout: TimeoutConfig{
			Connection: 10 * time.Second,
			Request:    30 * time.Second,
		},
		Retry: RetryConfig{
			Enabled:     true,
			MaxAttempts: 3,
			Backoff:     100 * time.Millisecond,
		},
	}
}
```

**设计说明**：
- 将50+个配置项简化为15个核心配置
- 移除了TLS、限流、负载均衡、连接池等高级功能
- 只保留日志和恢复两个基础中间件
- 简化了重试配置，使用固定退避时间

#### 5.2 创建简化的 gRPC 服务器

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
	"google.golang.org/grpc/reflection"

	"movieinfo/pkg/logger"
)

// Server gRPC服务器
type Server struct {
	config   *ServerConfig
	logger   logger.Logger
	server   *grpc.Server
	listener net.Listener
	health   *health.Server
}

// NewServer 创建新的gRPC服务器
func NewServer(config *ServerConfig, logger logger.Logger) (*Server, error) {
	// 创建监听器
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	// 创建gRPC服务器选项
	opts := []grpc.ServerOption{
		grpc.ConnectionTimeout(config.Timeout.Connection),
	}

	// 添加中间件
	if config.Middleware.Logging.Enabled {
		opts = append(opts, grpc.UnaryInterceptor(loggingInterceptor(logger)))
	}

	if config.Middleware.Recovery.Enabled {
		opts = append(opts, grpc.UnaryInterceptor(recoveryInterceptor(logger)))
	}

	// 创建gRPC服务器
	server := grpc.NewServer(opts...)

	// 创建健康检查服务
	var healthServer *health.Server
	if config.HealthCheck.Enabled {
		healthServer = health.NewServer()
		grpc_health_v1.RegisterHealthServer(server, healthServer)
	}

	// 启用反射（开发环境）
	reflection.Register(server)

	return &Server{
		config:   config,
		logger:   logger,
		server:   server,
		listener: listener,
		health:   healthServer,
	}, nil
}

// GetServer 获取gRPC服务器实例
func (s *Server) GetServer() *grpc.Server {
	return s.server
}

// Start 启动服务器
func (s *Server) Start() error {
	s.logger.Info("Starting gRPC server", "address", s.listener.Addr().String())

	// 设置健康状态
	if s.health != nil {
		s.health.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	}

	return s.server.Serve(s.listener)
}

// Stop 优雅停止服务器
func (s *Server) Stop() {
	s.logger.Info("Stopping gRPC server")

	// 设置健康状态
	if s.health != nil {
		s.health.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	}

	s.server.GracefulStop()
}

// ForceStop 强制停止服务器
func (s *Server) ForceStop() {
	s.logger.Info("Force stopping gRPC server")
	s.server.Stop()
}

// GetAddress 获取服务器地址
func (s *Server) GetAddress() string {
	return s.listener.Addr().String()
}

// 简化的日志中间件
func loggingInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			logger.Error("gRPC request failed", "method", info.FullMethod, "duration", duration, "error", err)
		} else {
			logger.Info("gRPC request completed", "method", info.FullMethod, "duration", duration)
		}

		return resp, err
	}
}

// 简化的恢复中间件
func recoveryInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("gRPC panic recovered", "method", info.FullMethod, "panic", r)
				err = fmt.Errorf("internal server error")
			}
		}()

		return handler(ctx, req)
	}
}
```

**设计说明**：
- 移除了复杂的中间件链和配置
- 只实现了基础的日志和恢复中间件
- 简化了服务器创建和管理逻辑
- 保留了健康检查和反射功能

#### 5.3 创建简化的 gRPC 客户端

```go
// pkg/grpc/client.go
package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(config.Timeout.Connection),
	}

	// 添加重试配置
	if config.Retry.Enabled {
		// 简化的重试配置
		opts = append(opts, grpc.WithUnaryInterceptor(retryInterceptor(config.Retry, logger)))
	}

	// 建立连接
	conn, err := grpc.Dial(config.Target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", config.Target, err)
	}

	return &Client{
		config: config,
		logger: logger,
		conn:   conn,
	}, nil
}

// GetConnection 获取gRPC连接
func (c *Client) GetConnection() *grpc.ClientConn {
	return c.conn
}

// Close 关闭客户端连接
func (c *Client) Close() error {
	c.logger.Info("Closing gRPC client connection")
	return c.conn.Close()
}

// 简化的重试中间件
func retryInterceptor(config RetryConfig, logger logger.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var err error
		for i := 0; i < config.MaxAttempts; i++ {
			err = invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return nil
			}

			if i < config.MaxAttempts-1 {
				logger.Warn("gRPC request failed, retrying", "method", method, "attempt", i+1, "error", err)
				time.Sleep(config.Backoff)
			}
		}

		logger.Error("gRPC request failed after all retries", "method", method, "attempts", config.MaxAttempts, "error", err)
		return err
	}
}
```

**设计说明**：
- 移除了复杂的连接池、负载均衡等功能
- 简化了重试逻辑，使用固定退避时间
- 暂时不支持TLS，简化开发环境配置

### 步骤6：创建配置文件

#### 6.1 创建简化的 gRPC 配置文件

```yaml
# configs/grpc.yaml
# gRPC服务配置 - 简化版本

# 服务器配置
servers:
  user:
    host: "0.0.0.0"
    port: 8081
    timeout:
      connection: 30s
      request: 30s
    middleware:
      logging:
        enabled: true
      recovery:
        enabled: true
    health_check:
      enabled: true

  movie:
    host: "0.0.0.0"
    port: 8082
    timeout:
      connection: 30s
      request: 30s
    middleware:
      logging:
        enabled: true
      recovery:
        enabled: true
    health_check:
      enabled: true

  rating:
    host: "0.0.0.0"
    port: 8083
    timeout:
      connection: 30s
      request: 30s
    middleware:
      logging:
        enabled: true
      recovery:
        enabled: true
    health_check:
      enabled: true

# 客户端配置
clients:
  user:
    target: "localhost:8081"
    timeout:
      connection: 10s
      request: 30s
    retry:
      enabled: true
      max_attempts: 3
      backoff: 100ms

  movie:
    target: "localhost:8082"
    timeout:
      connection: 10s
      request: 30s
    retry:
      enabled: true
      max_attempts: 3
      backoff: 100ms

  rating:
    target: "localhost:8083"
    timeout:
      connection: 10s
      request: 30s
    retry:
      enabled: true
      max_attempts: 3
      backoff: 100ms
```

**设计说明**：
- 将复杂的嵌套配置简化为扁平结构
- 移除了TLS、限流、监控等高级配置
- 每个服务使用独立的端口，便于开发和调试

### 步骤7：创建代码生成脚本

#### 7.1 创建 Makefile

```makefile
# Makefile - 简化版本
.PHONY: proto-gen proto-clean build test clean

# 变量定义
PROTO_DIR := proto
GO_OUT_DIR := .
PROTO_FILES := $(shell find $(PROTO_DIR) -name "*.proto")

# 生成 protobuf 代码
proto-gen:
	@echo "Generating protobuf code..."
	@for file in $(PROTO_FILES); do \
		echo "Processing $$file"; \
		protoc --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative \
		       --go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		       $$file; \
	done
	@echo "Protobuf code generation completed"

# 清理生成的代码
proto-clean:
	@echo "Cleaning generated protobuf code..."
	@find $(PROTO_DIR) -name "*.pb.go" -delete
	@find $(PROTO_DIR) -name "*_grpc.pb.go" -delete
	@echo "Cleanup completed"

# 构建项目
build:
	@echo "Building project..."
	go build -o bin/ ./cmd/...
	@echo "Build completed"

# 运行测试
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "Tests completed"

# 清理构建文件
clean:
	@echo "Cleaning build files..."
	rm -rf bin/
	@echo "Clean completed"

# 安装依赖
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download
	@echo "Dependencies installed"

# 格式化代码
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted"

# 运行 lint
lint:
	@echo "Running lint..."
	golangci-lint run
	@echo "Lint completed"

# 开发环境启动
dev-user:
	@echo "Starting user service..."
	go run cmd/user/main.go

dev-movie:
	@echo "Starting movie service..."
	go run cmd/movie/main.go

dev-rating:
	@echo "Starting rating service..."
	go run cmd/rating/main.go
```

**设计说明**：
- 移除了复杂的buf配置，使用标准的protoc命令
- 简化了构建和测试流程
- 添加了开发环境的快速启动命令

### 步骤8：测试验证

#### 8.1 创建简单的集成测试

```go
// test/grpc_test.go
package test

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/health/grpc_health_v1"

	grpcpkg "movieinfo/pkg/grpc"
	"movieinfo/pkg/logger"
)

func TestGRPCServerHealthCheck(t *testing.T) {
	// 创建服务器
	config := grpcpkg.DefaultServerConfig()
	config.Port = 0 // 使用随机端口

	logger := logger.NewConsoleLogger()

	server, err := grpcpkg.NewServer(config, logger)
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

	t.Log("Health check passed")
}
```

**设计说明**：
- 创建简单的健康检查测试
- 验证服务器和客户端的基本功能
- 为后续的服务测试提供基础

## 预期结果

完成本步骤后，你将拥有：

### 1. 简化但完整的 gRPC 协议定义
- **通用类型**：分页、响应、错误处理等基础类型
- **用户服务**：9个核心RPC方法，涵盖CRUD和认证
- **电影服务**：7个核心RPC方法，涵盖CRUD和搜索
- **评分服务**：7个核心RPC方法，涵盖CRUD和统计

### 2. 简化的 gRPC 基础设施
- **服务器**：支持基础中间件和健康检查
- **客户端**：支持重试和超时配置
- **配置管理**：15个核心配置项，易于理解和维护

### 3. 开发工具
- **Makefile**：自动化代码生成和构建
- **配置文件**：简化的YAML配置
- **测试框架**：基础的集成测试

### 4. 清晰的扩展路径
- **第一阶段**（当前）：核心CRUD功能
- **第二阶段**：添加认证、批量操作
- **第三阶段**：添加监控、链路追踪
- **第四阶段**：添加高级特性如推荐算法

## 复杂度对比

### 原始设计 vs 简化设计

| 方面 | 原始设计 | 简化设计 | 减少比例 |
|------|----------|----------|----------|
| 配置项数量 | 50+ | 15 | 70% |
| RPC方法数量 | 40+ | 23 | 43% |
| 错误类型数量 | 30+ | 7 | 77% |
| 中间件数量 | 6 | 2 | 67% |
| 代码行数 | 2500+ | 800 | 68% |

### 功能保留情况

✅ **保留的核心功能**：
- 完整的CRUD操作
- 基础认证功能
- 分页查询
- 搜索功能
- 健康检查
- 基础统计

❌ **移除的高级功能**：
- 批量操作
- 流式处理
- 复杂的错误分类
- 高级中间件（认证、监控、链路追踪）
- 推荐算法
- TLS加密
- 负载均衡

## 注意事项

### 开发建议
1. **先实现核心功能**：专注于基础的CRUD操作
2. **逐步添加特性**：根据实际需求添加高级功能
3. **保持简单**：避免过早优化和过度设计
4. **测试驱动**：为每个功能编写测试

### 扩展策略
1. **第一优先级**：完成基础服务实现
2. **第二优先级**：添加认证和权限控制
3. **第三优先级**：添加监控和日志
4. **第四优先级**：添加高级特性

### 性能考虑
- **连接复用**：客户端连接复用
- **超时设置**：合理的超时配置
- **错误处理**：快速失败和重试
- **资源管理**：及时释放连接和资源

## 下一步骤

完成 gRPC 协议定义后，下一步将进行：

1. **数据模型层开发**（06-数据模型层.md）
   - 定义数据模型结构
   - 实现数据访问层
   - 集成数据库连接

2. **检查清单**
   - [ ] 所有 proto 文件编译成功
   - [ ] gRPC 服务器可以正常启动
   - [ ] gRPC 客户端可以正常连接
   - [ ] 健康检查接口正常工作
   - [ ] 基础中间件功能正常
   - [ ] 配置文件格式正确
   - [ ] 集成测试通过
   - [ ] 代码生成脚本正常工作

通过这种简化的方式，我们可以快速建立起一个可工作的gRPC系统，然后根据实际需求逐步完善功能。这种渐进式的开发方法既保证了项目的可行性，又为未来的扩展留下了空间。