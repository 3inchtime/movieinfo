# Movieinfo Project Makefile
# 简化版gRPC开发管理

.PHONY: help proto-gen proto-clean build test clean

# 默认目标
help:
	@echo "Available targets:"
	@echo "  proto-gen    - Generate gRPC code from proto files"
	@echo "  proto-clean  - Clean generated proto files"
	@echo "  build        - Build all services"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"

# 生成gRPC代码（手动方式，不依赖protoc）
proto-gen:
	@echo "Generating gRPC code..."
	@mkdir -p proto/gen/common proto/gen/user proto/gen/movie proto/gen/rating
	@echo "Note: This requires protoc to be installed."
	@echo "If protoc is not available, the proto files are ready for manual generation."
	@echo "Proto files location: proto/"
	@echo "Target generation directory: proto/gen/"

# 清理生成的代码
proto-clean:
	@echo "Cleaning generated proto files..."
	@rm -rf proto/gen

# 构建所有服务
build:
	@echo "Building services..."
	go build -o bin/user-service ./cmd/user
	go build -o bin/movie-service ./cmd/movie
	go build -o bin/rating-service ./cmd/rating
	go build -o bin/web-service ./cmd/web

# 运行测试
test:
	@echo "Running tests..."
	go test ./...

# 清理构建产物
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf proto/gen/

# 初始化开发环境
init:
	@echo "Initializing development environment..."
	go mod tidy
	go mod download
	@echo "Development environment ready!"