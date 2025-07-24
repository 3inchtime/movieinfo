package main

// gRPC服务器示例
// 演示简化版gRPC协议定义的基本用法

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 模拟的用户服务实现
type UserServiceServer struct {
	// 这里应该嵌入生成的UnimplementedUserServiceServer
	// UnimplementedUserServiceServer
}

// 模拟的电影服务实现
type MovieServiceServer struct {
	// 这里应该嵌入生成的UnimplementedMovieServiceServer
	// UnimplementedMovieServiceServer
}

// 模拟的评分服务实现
type RatingServiceServer struct {
	// 这里应该嵌入生成的UnimplementedRatingServiceServer
	// UnimplementedRatingServiceServer
}

// 日志中间件
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	log.Printf("[INFO] gRPC call: %s started", info.FullMethod)
	
	resp, err := handler(ctx, req)
	
	duration := time.Since(start)
	if err != nil {
		log.Printf("[ERROR] gRPC call: %s failed in %v: %v", info.FullMethod, duration, err)
	} else {
		log.Printf("[INFO] gRPC call: %s completed in %v", info.FullMethod, duration)
	}
	
	return resp, err
}

// 错误处理中间件
func errorHandlingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	
	if err != nil {
		// 将内部错误转换为gRPC状态码
		switch err.Error() {
		case "user not found":
			return nil, status.Error(codes.NotFound, "用户不存在")
		case "invalid input":
			return nil, status.Error(codes.InvalidArgument, "输入参数无效")
		default:
			return nil, status.Error(codes.Internal, "内部服务器错误")
		}
	}
	
	return resp, err
}

func main() {
	// 创建gRPC服务器
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingInterceptor,
			errorHandlingInterceptor,
		),
		grpc.MaxRecvMsgSize(4*1024*1024), // 4MB
		grpc.MaxSendMsgSize(4*1024*1024), // 4MB
	)

	// 注册服务
	// 注意：这里需要等proto文件生成后才能正常编译
	// pb.RegisterUserServiceServer(server, &UserServiceServer{})
	// pb.RegisterMovieServiceServer(server, &MovieServiceServer{})
	// pb.RegisterRatingServiceServer(server, &RatingServiceServer{})

	// 启动服务器
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("gRPC服务器启动成功，监听端口 :9090")
	fmt.Println("服务列表:")
	fmt.Println("  - UserService    (用户服务)")
	fmt.Println("  - MovieService   (电影服务)")
	fmt.Println("  - RatingService  (评分服务)")
	fmt.Println("")
	fmt.Println("中间件:")
	fmt.Println("  - 日志记录")
	fmt.Println("  - 错误处理")
	fmt.Println("")
	fmt.Println("按 Ctrl+C 停止服务器")

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}