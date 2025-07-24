package main

// gRPC客户端示例
// 演示如何连接和使用简化版gRPC服务

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// gRPC客户端配置
type GRPCClientConfig struct {
	Address     string
	Timeout     time.Duration
	MaxRecvSize int
	MaxSendSize int
}

// 默认客户端配置
func DefaultClientConfig() *GRPCClientConfig {
	return &GRPCClientConfig{
		Address:     "localhost:9090",
		Timeout:     5 * time.Second,
		MaxRecvSize: 4 * 1024 * 1024, // 4MB
		MaxSendSize: 4 * 1024 * 1024, // 4MB
	}
}

// 创建gRPC连接
func createGRPCConnection(config *GRPCClientConfig) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, config.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(config.MaxRecvSize),
			grpc.MaxCallSendMsgSize(config.MaxSendSize),
		),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return conn, nil
}

// 用户服务客户端示例
func demonstrateUserService(conn *grpc.ClientConn) {
	fmt.Println("=== 用户服务示例 ===")
	
	// 注意：这里需要等proto文件生成后才能正常编译
	// client := pb.NewUserServiceClient(conn)
	// 
	// // 创建用户
	// createReq := &pb.CreateUserRequest{
	// 	Username: "testuser",
	// 	Email:    "test@example.com",
	// 	Password: "password123",
	// }
	// 
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// 
	// createResp, err := client.CreateUser(ctx, createReq)
	// if err != nil {
	// 	log.Printf("创建用户失败: %v", err)
	// 	return
	// }
	// 
	// fmt.Printf("用户创建成功: ID=%d\n", createResp.User.Id)
	
	fmt.Println("用户服务功能:")
	fmt.Println("  - CreateUser: 创建新用户")
	fmt.Println("  - GetUser: 获取用户信息")
	fmt.Println("  - UpdateUser: 更新用户信息")
	fmt.Println("  - DeleteUser: 删除用户")
	fmt.Println("  - ListUsers: 列出用户")
	fmt.Println("  - Login: 用户登录")
	fmt.Println("  - Logout: 用户登出")
	fmt.Println("  - ChangePassword: 修改密码")
	fmt.Println()
}

// 电影服务客户端示例
func demonstrateMovieService(conn *grpc.ClientConn) {
	fmt.Println("=== 电影服务示例 ===")
	
	fmt.Println("电影服务功能:")
	fmt.Println("  - CreateMovie: 创建新电影")
	fmt.Println("  - GetMovie: 获取电影信息")
	fmt.Println("  - UpdateMovie: 更新电影信息")
	fmt.Println("  - DeleteMovie: 删除电影")
	fmt.Println("  - ListMovies: 列出电影")
	fmt.Println("  - SearchMovies: 搜索电影")
	fmt.Println()
}

// 评分服务客户端示例
func demonstrateRatingService(conn *grpc.ClientConn) {
	fmt.Println("=== 评分服务示例 ===")
	
	fmt.Println("评分服务功能:")
	fmt.Println("  - CreateRating: 创建新评分")
	fmt.Println("  - GetRating: 获取评分信息")
	fmt.Println("  - UpdateRating: 更新评分")
	fmt.Println("  - DeleteRating: 删除评分")
	fmt.Println("  - ListRatings: 列出评分")
	fmt.Println("  - GetMovieAverageRating: 获取电影平均评分")
	fmt.Println()
}

// 健康检查示例
func demonstrateHealthCheck(conn *grpc.ClientConn) {
	fmt.Println("=== 健康检查示例 ===")
	
	// 注意：这里需要等proto文件生成后才能正常编译
	// client := pb.NewUserServiceClient(conn) // 任意服务都有健康检查
	// 
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()
	// 
	// healthReq := &pb.HealthCheckRequest{
	// 	Service: "user",
	// }
	// 
	// healthResp, err := client.HealthCheck(ctx, healthReq)
	// if err != nil {
	// 	log.Printf("健康检查失败: %v", err)
	// 	return
	// }
	// 
	// fmt.Printf("服务状态: %s\n", healthResp.Status)
	
	fmt.Println("健康检查功能:")
	fmt.Println("  - 检查服务是否正常运行")
	fmt.Println("  - 返回服务状态信息")
	fmt.Println()
}

func main() {
	fmt.Println("gRPC客户端示例")
	fmt.Println("===============")
	fmt.Println()

	// 创建客户端配置
	config := DefaultClientConfig()
	fmt.Printf("连接地址: %s\n", config.Address)
	fmt.Printf("连接超时: %v\n", config.Timeout)
	fmt.Println()

	// 创建gRPC连接
	fmt.Println("正在连接gRPC服务器...")
	conn, err := createGRPCConnection(config)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	fmt.Println("连接成功！")
	fmt.Println()

	// 演示各个服务
	demonstrateUserService(conn)
	demonstrateMovieService(conn)
	demonstrateRatingService(conn)
	demonstrateHealthCheck(conn)

	fmt.Println("注意: 以上示例需要proto文件生成后才能正常编译运行")
	fmt.Println("请先运行代码生成脚本生成gRPC代码")
}