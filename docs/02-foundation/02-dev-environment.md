# 开发环境搭建

## 1. 概述

本文档为MovieInfo项目的开发者提供了一份完整的开发环境搭建指南，旨在确保所有开发者都能在一个一致、高效的环境中工作。

## 2. 系统要求

- **操作系统**: Windows 10/11, macOS, or Linux
- **硬件**: 8GB RAM, 4-core CPU, 50GB 可用磁盘空间

## 3. 基础软件安装

### 3.1. Go 语言环境

- **版本**: 1.18 或更高
- **安装**: 从 [Go官方网站](https://golang.org/dl/) 下载并安装。
- **验证**: `go version`

### 3.2. MySQL

- **版本**: 8.0 或更高
- **安装**: 从 [MySQL官方网站](https://dev.mysql.com/downloads/mysql/) 下载并安装，或使用Docker。
- **验证**: `mysql --version`

### 3.3. Redis

- **版本**: 6.0 或更高
- **安装**: 从 [Redis官方网站](https://redis.io/download) 下载并安装，或使用Docker。
- **验证**: `redis-server --version`

## 4. 容器化环境 (推荐)

强烈建议使用Docker和Docker Compose来管理本地的数据库和缓存服务，以简化环境配置。

- **安装**: 参考 [Docker官方文档](https://docs.docker.com/get-docker/) 进行安装。
- **启动依赖服务**: 在项目根目录运行 `docker-compose up -d mysql redis` 来后台启动MySQL和Redis容器。

## 5. 开发工具

### 5.1. Git

- **安装**: 从 [Git官网](https://git-scm.com/downloads) 下载并安装。

### 5.2. IDE/编辑器

- **推荐**: Visual Studio Code
- **VS Code 推荐扩展**:
    - `Go`: Go语言官方支持扩展
    - `Docker`: Docker集成
    - `Remote - Containers`: 用于容器化开发
    - `gRPC`: gRPC/Protobuf支持

## 6. 项目初始化

### 6.1. 克隆项目

```bash
git clone https://github.com/your-username/movieinfo.git
cd movieinfo
```

### 6.2. 安装依赖

```bash
go mod tidy
```

### 6.3. 环境配置

创建并编辑 `config/local.yaml` 文件，配置本地数据库和Redis连接。

```bash
cp config/config.example.yaml config/local.yaml
```

修改 `config/local.yaml`:

```yaml
mysql:
  host: "127.0.0.1"
  port: 3306
  user: "root"
  password: "your_password"
  dbname: "movieinfo"

redis:
  host: "127.0.0.1"
  port: 6379
```

### 6.4. 数据库初始化

```sql
CREATE DATABASE movieinfo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## 7. 开发工作流

1. **拉取最新代码**: `git pull origin main`
2. **启动/重启服务**: 
    - 使用Docker Compose: `docker-compose up --build`
    - 手动运行: `go run cmd/your-service/main.go`
3. **进行代码开发**
4. **运行测试**: `go test ./...`
5. **提交代码**: `git commit`, `git push`

---