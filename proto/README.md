# gRPC协议定义 - 简化版本

本目录包含了MovieInfo项目的gRPC协议定义文件，采用简化设计策略，专注于核心功能。

## 目录结构

```
proto/
├── common/                 # 通用定义
│   ├── common.proto        # 通用数据类型（分页、响应等）
│   └── error.proto         # 错误定义
├── user/                   # 用户服务
│   ├── user.proto          # 用户数据结构
│   └── user_service.proto  # 用户服务接口
├── movie/                  # 电影服务
│   ├── movie.proto         # 电影数据结构
│   └── movie_service.proto # 电影服务接口
├── rating/                 # 评分服务
│   ├── rating.proto        # 评分数据结构
│   └── rating_service.proto# 评分服务接口
└── gen/                    # 生成的Go代码（自动生成）
```

## 简化设计原则

### 1. 核心功能优先
- 只保留最基础的CRUD操作
- 移除复杂的批量操作和高级功能
- 专注于MVP（最小可行产品）需求

### 2. 配置简化
- 减少配置层次和选项
- 使用合理的默认值
- 避免过度工程化

### 3. 错误处理简化
- 只定义最常用的错误类型
- 使用标准的gRPC状态码
- 简化错误消息结构

### 4. 中间件最小化
- 只实现必要的日志和错误处理中间件
- 避免复杂的认证和授权逻辑（在MVP阶段）

## 服务定义

### 用户服务 (UserService)
- `CreateUser` - 创建用户
- `GetUser` - 获取用户信息
- `UpdateUser` - 更新用户信息
- `DeleteUser` - 删除用户
- `ListUsers` - 列出用户（分页）
- `Login` - 用户登录
- `Logout` - 用户登出
- `ChangePassword` - 修改密码
- `HealthCheck` - 健康检查

### 电影服务 (MovieService)
- `CreateMovie` - 创建电影
- `GetMovie` - 获取电影信息
- `UpdateMovie` - 更新电影信息
- `DeleteMovie` - 删除电影
- `ListMovies` - 列出电影（分页）
- `SearchMovies` - 搜索电影
- `HealthCheck` - 健康检查

### 评分服务 (RatingService)
- `CreateRating` - 创建评分
- `GetRating` - 获取评分信息
- `UpdateRating` - 更新评分
- `DeleteRating` - 删除评分
- `ListRatings` - 列出评分（分页）
- `GetMovieAverageRating` - 获取电影平均评分
- `HealthCheck` - 健康检查

## 代码生成

### 前置要求

1. 安装Protocol Buffers编译器 (protoc)
2. 安装Go插件：
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

### 生成代码

#### Windows
```bash
scripts\generate_proto.bat
```

#### Linux/Mac
```bash
chmod +x scripts/generate_proto.sh
./scripts/generate_proto.sh
```

#### 使用Makefile
```bash
make proto-gen
```

### 生成的文件

生成的Go代码将位于 `proto/gen/` 目录下：

```
proto/gen/
├── common/
│   ├── common.pb.go        # 通用数据类型
│   └── error.pb.go         # 错误定义
├── user/
│   ├── user.pb.go          # 用户数据结构
│   └── user_service_grpc.pb.go # 用户服务接口
├── movie/
│   ├── movie.pb.go         # 电影数据结构
│   └── movie_service_grpc.pb.go # 电影服务接口
└── rating/
    ├── rating.pb.go        # 评分数据结构
    └── rating_service_grpc.pb.go # 评分服务接口
```

## 使用示例

### 服务器端

```go
package main

import (
    "net"
    "google.golang.org/grpc"
    pb "github.com/3inchtime/movieinfo/proto/gen/user"
)

type userServiceServer struct {
    pb.UnimplementedUserServiceServer
}

func main() {
    lis, _ := net.Listen("tcp", ":9090")
    server := grpc.NewServer()
    
    pb.RegisterUserServiceServer(server, &userServiceServer{})
    
    server.Serve(lis)
}
```

### 客户端

```go
package main

import (
    "context"
    "google.golang.org/grpc"
    pb "github.com/3inchtime/movieinfo/proto/gen/user"
)

func main() {
    conn, _ := grpc.Dial("localhost:9090", grpc.WithInsecure())
    defer conn.Close()
    
    client := pb.NewUserServiceClient(conn)
    
    resp, _ := client.GetUser(context.Background(), &pb.GetUserRequest{
        Id: 1,
    })
    
    // 处理响应
}
```

## 配置文件

gRPC相关配置位于 `configs/grpc.yaml`，包含：
- 服务器监听配置
- 连接参数
- 中间件设置
- 超时配置

## 开发建议

1. **渐进式开发**：先实现基础功能，后续根据需要添加高级特性
2. **测试优先**：为每个RPC方法编写单元测试
3. **文档同步**：及时更新proto文件的注释
4. **版本管理**：使用语义化版本管理proto文件变更

## 扩展计划

在MVP版本稳定后，可以考虑添加：
- 批量操作接口
- 流式RPC
- 高级搜索功能
- 缓存策略
- 认证和授权
- 服务发现
- 负载均衡

## 故障排除

### 常见问题

1. **protoc未找到**
   - 确保已安装Protocol Buffers编译器
   - 检查PATH环境变量

2. **插件未找到**
   - 运行 `go install` 命令安装必要插件
   - 确保 `$GOPATH/bin` 在PATH中

3. **导入路径错误**
   - 检查 `go_package` 选项是否正确
   - 确保模块路径与项目结构匹配

### 获取帮助

如果遇到问题，请检查：
1. 项目文档
2. gRPC官方文档
3. Protocol Buffers文档