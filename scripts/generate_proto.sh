#!/bin/bash
# gRPC代码生成脚本 - Linux/Mac版本
# 根据简化版gRPC协议定义生成Go代码

set -e

echo "开始生成gRPC代码..."

# 设置项目根目录
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$PROJECT_ROOT"

# 检查protoc是否安装
if ! command -v protoc &> /dev/null; then
    echo "错误: protoc未安装或不在PATH中"
    echo "请安装Protocol Buffers编译器"
    exit 1
fi

# 检查protoc-gen-go是否安装
if ! command -v protoc-gen-go &> /dev/null; then
    echo "错误: protoc-gen-go未安装"
    echo "请运行: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    exit 1
fi

# 检查protoc-gen-go-grpc是否安装
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "错误: protoc-gen-go-grpc未安装"
    echo "请运行: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

# 创建输出目录
mkdir -p proto/gen/{common,user,movie,rating}

echo "生成通用模块代码..."
protoc --proto_path=proto \
    --go_out=proto/gen \
    --go_opt=paths=source_relative \
    --go-grpc_out=proto/gen \
    --go-grpc_opt=paths=source_relative \
    common/common.proto common/error.proto

echo "生成用户服务代码..."
protoc --proto_path=proto \
    --go_out=proto/gen \
    --go_opt=paths=source_relative \
    --go-grpc_out=proto/gen \
    --go-grpc_opt=paths=source_relative \
    user/user.proto user/user_service.proto

echo "生成电影服务代码..."
protoc --proto_path=proto \
    --go_out=proto/gen \
    --go_opt=paths=source_relative \
    --go-grpc_out=proto/gen \
    --go-grpc_opt=paths=source_relative \
    movie/movie.proto movie/movie_service.proto

echo "生成评分服务代码..."
protoc --proto_path=proto \
    --go_out=proto/gen \
    --go_opt=paths=source_relative \
    --go-grpc_out=proto/gen \
    --go-grpc_opt=paths=source_relative \
    rating/rating.proto rating/rating_service.proto

echo
echo "gRPC代码生成完成！"
echo "生成的文件位于: proto/gen/"
echo
echo "文件结构:"
echo "  proto/gen/common/    - 通用数据类型和错误定义"
echo "  proto/gen/user/      - 用户服务相关代码"
echo "  proto/gen/movie/     - 电影服务相关代码"
echo "  proto/gen/rating/    - 评分服务相关代码"
echo