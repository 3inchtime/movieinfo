# 3.1 项目结构创建

## 概述

项目结构是软件项目的骨架，它决定了代码的组织方式、模块的划分和团队的协作效率。对于MovieInfo项目，我们需要创建一个清晰、可扩展、符合Go语言最佳实践的项目结构。

## 为什么项目结构如此重要？

### 1. **代码组织**
- **模块化设计**：清晰的目录结构体现了系统的模块划分
- **职责分离**：不同类型的代码放在不同的目录中
- **依赖管理**：合理的结构有助于管理模块间的依赖关系
- **代码复用**：公共代码的合理放置便于复用

### 2. **团队协作**
- **统一认知**：团队成员对项目结构有一致的理解
- **并行开发**：不同模块可以由不同开发者并行开发
- **代码审查**：清晰的结构便于代码审查和维护
- **新人上手**：新团队成员能快速理解项目结构

### 3. **项目维护**
- **易于定位**：快速找到需要修改的代码
- **影响分析**：了解修改对其他模块的影响
- **重构支持**：结构化的代码便于重构
- **文档生成**：自动生成项目文档

### 4. **扩展性**
- **功能扩展**：新功能的添加有明确的位置
- **服务拆分**：为微服务化提供清晰的边界
- **技术演进**：支持技术栈的升级和演进
- **性能优化**：便于识别性能瓶颈和优化点

## Go项目结构标准

### 1. **Go项目布局标准**

Go社区推荐的项目布局遵循以下原则：
- **cmd/**: 应用程序的主要入口点
- **internal/**: 私有应用程序和库代码
- **pkg/**: 外部应用程序可以使用的库代码
- **api/**: OpenAPI/Swagger规范，JSON模式文件，协议定义文件
- **web/**: Web应用程序特定的组件
- **configs/**: 配置文件模板或默认配置
- **init/**: 系统初始化（systemd，upstart，sysv）和进程管理器（runit，supervisord）配置
- **scripts/**: 执行各种构建，安装，分析等操作的脚本
- **build/**: 打包和持续集成
- **deployments/**: IaaS，PaaS，系统和容器编排部署配置和模板
- **test/**: 额外的外部测试应用程序和测试数据

### 2. **MovieInfo项目结构设计**

```
movieinfo/
├── cmd/                    # 应用程序入口点
│   ├── web/               # 主页服务
│   │   └── main.go
│   ├── user/              # 用户服务
│   │   └── main.go
│   ├── movie/             # 电影服务
│   │   └── main.go
│   ├── comment/           # 评论服务
│   │   └── main.go
│   └── migrate/           # 数据库迁移工具
│       └── main.go
├── internal/              # 私有应用程序代码
│   ├── config/            # 配置管理
│   │   ├── config.go
│   │   └── database.go
│   ├── models/            # 数据模型
│   │   ├── user.go
│   │   ├── movie.go
│   │   ├── comment.go
│   │   └── rating.go
│   ├── handlers/          # HTTP处理器
│   │   ├── web/           # Web页面处理器
│   │   ├── user/          # 用户API处理器
│   │   ├── movie/         # 电影API处理器
│   │   └── comment/       # 评论API处理器
│   ├── services/          # 业务逻辑服务
│   │   ├── user/
│   │   ├── movie/
│   │   └── comment/
│   ├── repository/        # 数据访问层
│   │   ├── user/
│   │   ├── movie/
│   │   └── comment/
│   └── middleware/        # 中间件
│       ├── auth.go
│       ├── cors.go
│       ├── logger.go
│       └── recovery.go
├── pkg/                   # 公共库代码
│   ├── database/          # 数据库连接
│   │   ├── mysql.go
│   │   └── connection.go
│   ├── redis/             # Redis连接
│   │   ├── client.go
│   │   └── cache.go
│   ├── grpc/              # gRPC客户端/服务端
│   │   ├── client/
│   │   ├── server/
│   │   └── interceptor/
│   ├── logger/            # 日志工具
│   │   ├── logger.go
│   │   └── formatter.go
│   ├── utils/             # 工具函数
│   │   ├── hash.go
│   │   ├── jwt.go
│   │   ├── validator.go
│   │   └── response.go
│   └── errors/            # 错误定义
│       ├── codes.go
│       └── errors.go
├── api/                   # API定义
│   ├── openapi/           # OpenAPI规范
│   │   └── movieinfo.yaml
│   └── proto/             # Protocol Buffers定义
│       ├── user/
│       │   └── user.proto
│       ├── movie/
│       │   └── movie.proto
│       └── comment/
│           └── comment.proto
├── web/                   # Web资源
│   ├── templates/         # HTML模板
│   │   ├── layouts/
│   │   ├── pages/
│   │   └── partials/
│   └── static/            # 静态资源
│       ├── css/
│       ├── js/
│       ├── images/
│       └── fonts/
├── configs/               # 配置文件
│   ├── config.yaml
│   ├── config.example.yaml
│   └── database.yaml
├── scripts/               # 脚本文件
│   ├── build.sh
│   ├── test.sh
│   ├── deploy.sh
│   └── migrate.sh
├── build/                 # 构建相关
│   ├── Dockerfile.web
│   ├── Dockerfile.user
│   ├── Dockerfile.movie
│   ├── Dockerfile.comment
│   └── docker/
│       ├── mysql/
│       └── redis/
├── deployments/           # 部署配置
│   ├── docker-compose.yml
│   ├── docker-compose.prod.yml
│   └── k8s/
├── test/                  # 测试相关
│   ├── integration/       # 集成测试
│   ├── fixtures/          # 测试数据
│   └── mocks/             # Mock对象
├── docs/                  # 项目文档
│   ├── api/               # API文档
│   ├── architecture/      # 架构文档
│   └── development/       # 开发文档
├── tools/                 # 工具和实用程序
│   └── codegen/           # 代码生成工具
├── vendor/                # 依赖包（可选）
├── .gitignore
├── .golangci.yml
├── .air.toml
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── LICENSE
```

## 创建项目结构

### 1. **自动化创建脚本**

#### 1.1 项目初始化脚本
```bash
#!/bin/bash
# create-project-structure.sh

PROJECT_NAME="movieinfo"
PROJECT_ROOT=$(pwd)/$PROJECT_NAME

echo "Creating MovieInfo project structure..."

# 创建根目录
mkdir -p $PROJECT_ROOT
cd $PROJECT_ROOT

# 创建主要目录结构
echo "Creating directory structure..."

# 应用程序入口点
mkdir -p cmd/{web,user,movie,comment,migrate}

# 私有应用程序代码
mkdir -p internal/{config,models,handlers/{web,user,movie,comment},services/{user,movie,comment},repository/{user,movie,comment},middleware}

# 公共库代码
mkdir -p pkg/{database,redis,grpc/{client,server,interceptor},logger,utils,errors}

# API定义
mkdir -p api/{openapi,proto/{user,movie,comment}}

# Web资源
mkdir -p web/{templates/{layouts,pages,partials},static/{css,js,images,fonts}}

# 配置文件
mkdir -p configs

# 脚本文件
mkdir -p scripts

# 构建相关
mkdir -p build/{docker/{mysql,redis}}

# 部署配置
mkdir -p deployments/k8s

# 测试相关
mkdir -p test/{integration,fixtures,mocks}

# 项目文档
mkdir -p docs/{api,architecture,development}

# 工具
mkdir -p tools/codegen

echo "Directory structure created successfully!"

# 创建基础文件
echo "Creating basic files..."

# 创建主要入口文件
cat > cmd/web/main.go << 'EOF'
package main

import (
    "log"
    "net/http"
    
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    
    r.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "MovieInfo Web Service",
            "version": "1.0.0",
        })
    })
    
    log.Println("Web service starting on :8080")
    log.Fatal(r.Run(":8080"))
}
EOF

cat > cmd/user/main.go << 'EOF'
package main

import (
    "log"
    "net/http"
    
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "service": "user",
            "status": "healthy",
        })
    })
    
    log.Println("User service starting on :8081")
    log.Fatal(r.Run(":8081"))
}
EOF

cat > cmd/movie/main.go << 'EOF'
package main

import (
    "log"
    "net/http"
    
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "service": "movie",
            "status": "healthy",
        })
    })
    
    log.Println("Movie service starting on :8082")
    log.Fatal(r.Run(":8082"))
}
EOF

cat > cmd/comment/main.go << 'EOF'
package main

import (
    "log"
    "net/http"
    
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "service": "comment",
            "status": "healthy",
        })
    })
    
    log.Println("Comment service starting on :8083")
    log.Fatal(r.Run(":8083"))
}
EOF

# 创建配置文件示例
cat > configs/config.example.yaml << 'EOF'
# 数据库配置
database:
  host: localhost
  port: 3306
  username: movieinfo
  password: movieinfo123
  dbname: movieinfo
  charset: utf8mb4
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600

# Redis配置
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
  pool_size: 10

# 服务端口配置
services:
  web:
    port: 8080
  user:
    http_port: 8081
    grpc_port: 9081
  movie:
    http_port: 8082
    grpc_port: 9082
  comment:
    http_port: 8083
    grpc_port: 9083

# JWT配置
jwt:
  secret: your_jwt_secret_key
  expire_hours: 24

# 日志配置
log:
  level: info
  format: json
  output: stdout
EOF

# 创建.gitignore文件
cat > .gitignore << 'EOF'
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
vendor/

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Project specific
bin/
tmp/
logs/
*.log
coverage.html
coverage.out

# Configuration files (keep examples)
configs/config.yaml
configs/database.yaml
!configs/*.example.*

# Build artifacts
build/bin/
build/tmp/

# Docker volumes
docker/mysql/data/
docker/redis/data/

# Environment files
.env
.env.local
.env.*.local
EOF

# 创建README.md
cat > README.md << 'EOF'
# MovieInfo

MovieInfo 是一个基于 Go 语言和 Gin 框架开发的电影信息网站，采用服务化架构设计。

## 项目结构

```
movieinfo/
├── cmd/                    # 应用程序入口点
├── internal/              # 私有应用程序代码
├── pkg/                   # 公共库代码
├── api/                   # API定义
├── web/                   # Web资源
├── configs/               # 配置文件
├── scripts/               # 脚本文件
├── build/                 # 构建相关
├── deployments/           # 部署配置
├── test/                  # 测试相关
├── docs/                  # 项目文档
└── tools/                 # 工具和实用程序
```

## 快速开始

1. 克隆项目
```bash
git clone <repository-url>
cd movieinfo
```

2. 安装依赖
```bash
go mod tidy
```

3. 配置环境
```bash
cp configs/config.example.yaml configs/config.yaml
# 编辑 configs/config.yaml 配置数据库连接等信息
```

4. 启动服务
```bash
# 启动所有服务
make run-all

# 或者单独启动服务
go run cmd/web/main.go
go run cmd/user/main.go
go run cmd/movie/main.go
go run cmd/comment/main.go
```

## 开发指南

详细的开发指南请参考 [docs/development/](docs/development/) 目录。

## API 文档

API 文档请参考 [docs/api/](docs/api/) 目录。

## 许可证

MIT License
EOF

echo "Basic files created successfully!"
echo "Project structure for MovieInfo has been created at: $PROJECT_ROOT"
echo ""
echo "Next steps:"
echo "1. cd $PROJECT_ROOT"
echo "2. go mod init github.com/yourname/movieinfo"
echo "3. go mod tidy"
echo "4. Copy configs/config.example.yaml to configs/config.yaml and configure"
echo "5. Start development!"
```

#### 1.2 运行创建脚本
```bash
# 使脚本可执行
chmod +x create-project-structure.sh

# 运行脚本
./create-project-structure.sh

# 进入项目目录
cd movieinfo
```

### 2. **目录详细说明**

#### 2.1 cmd/ - 应用程序入口
```
cmd/
├── web/           # 主页服务 - Web网关和前端页面
├── user/          # 用户服务 - 用户管理和认证
├── movie/         # 电影服务 - 电影数据管理
├── comment/       # 评论服务 - 评论和评分
└── migrate/       # 数据库迁移工具
```

**设计原则**：
- 每个服务一个独立的main.go
- 保持main.go简洁，主要逻辑在internal/
- 支持命令行参数和环境变量配置

#### 2.2 internal/ - 私有代码
```
internal/
├── config/        # 配置管理 - 统一的配置加载和管理
├── models/        # 数据模型 - 数据库实体和业务对象
├── handlers/      # HTTP处理器 - 按服务分组的API处理器
├── services/      # 业务逻辑 - 核心业务逻辑实现
├── repository/    # 数据访问 - 数据库操作抽象层
└── middleware/    # 中间件 - HTTP中间件实现
```

**设计原则**：
- 按功能模块组织代码
- 清晰的分层架构
- 避免循环依赖

#### 2.3 pkg/ - 公共库
```
pkg/
├── database/      # 数据库连接 - MySQL连接池管理
├── redis/         # Redis客户端 - 缓存操作封装
├── grpc/          # gRPC支持 - 客户端和服务端工具
├── logger/        # 日志工具 - 统一的日志接口
├── utils/         # 工具函数 - 通用工具函数
└── errors/        # 错误定义 - 统一的错误码和错误处理
```

**设计原则**：
- 可被外部项目使用的代码
- 无业务逻辑，纯技术实现
- 良好的接口设计

#### 2.4 api/ - API定义
```
api/
├── openapi/       # OpenAPI规范 - REST API文档
└── proto/         # Protocol Buffers - gRPC接口定义
    ├── user/
    ├── movie/
    └── comment/
```

**设计原则**：
- API优先设计
- 版本化管理
- 自动生成代码

### 3. **项目结构验证**

#### 3.1 结构验证脚本
```bash
#!/bin/bash
# verify-structure.sh

echo "Verifying MovieInfo project structure..."

# 检查必要目录
REQUIRED_DIRS=(
    "cmd/web"
    "cmd/user"
    "cmd/movie"
    "cmd/comment"
    "internal/config"
    "internal/models"
    "internal/handlers"
    "internal/services"
    "internal/repository"
    "internal/middleware"
    "pkg/database"
    "pkg/redis"
    "pkg/grpc"
    "pkg/logger"
    "pkg/utils"
    "pkg/errors"
    "api/openapi"
    "api/proto"
    "web/templates"
    "web/static"
    "configs"
    "scripts"
    "build"
    "deployments"
    "test"
    "docs"
    "tools"
)

MISSING_DIRS=()

for dir in "${REQUIRED_DIRS[@]}"; do
    if [ ! -d "$dir" ]; then
        MISSING_DIRS+=("$dir")
    fi
done

if [ ${#MISSING_DIRS[@]} -eq 0 ]; then
    echo "✅ All required directories exist"
else
    echo "❌ Missing directories:"
    for dir in "${MISSING_DIRS[@]}"; do
        echo "  - $dir"
    done
    exit 1
fi

# 检查必要文件
REQUIRED_FILES=(
    "go.mod"
    "README.md"
    ".gitignore"
    "configs/config.example.yaml"
    "cmd/web/main.go"
    "cmd/user/main.go"
    "cmd/movie/main.go"
    "cmd/comment/main.go"
)

MISSING_FILES=()

for file in "${REQUIRED_FILES[@]}"; do
    if [ ! -f "$file" ]; then
        MISSING_FILES+=("$file")
    fi
done

if [ ${#MISSING_FILES[@]} -eq 0 ]; then
    echo "✅ All required files exist"
else
    echo "❌ Missing files:"
    for file in "${MISSING_FILES[@]}"; do
        echo "  - $file"
    done
    exit 1
fi

echo "✅ Project structure verification completed successfully!"
```

#### 3.2 运行验证
```bash
# 使脚本可执行
chmod +x verify-structure.sh

# 运行验证
./verify-structure.sh
```

## 项目结构最佳实践

### 1. **命名规范**
- **目录名**：使用小写字母和下划线
- **文件名**：使用小写字母和下划线
- **包名**：简短、小写、单数形式
- **接口名**：以er结尾（如Reader、Writer）

### 2. **依赖管理**
- **internal包**：不能被外部导入
- **pkg包**：可以被外部导入
- **避免循环依赖**：合理设计包的依赖关系
- **接口隔离**：使用接口减少包间耦合

### 3. **文件组织**
- **单一职责**：每个文件只负责一个功能
- **合理大小**：单个文件不超过500行
- **相关性**：相关的类型和函数放在同一个文件
- **测试文件**：与源文件放在同一个包中

### 4. **扩展性考虑**
- **预留扩展点**：为未来功能扩展预留目录
- **插件机制**：支持插件式的功能扩展
- **配置驱动**：通过配置控制功能开关
- **版本兼容**：考虑API版本兼容性

## 总结

项目结构创建为MovieInfo项目奠定了坚实的基础。通过合理的目录组织和文件布局，我们建立了一个清晰、可维护、可扩展的项目架构。

**关键设计要点**：
1. **标准化**：遵循Go社区的项目布局标准
2. **模块化**：清晰的模块划分和职责分离
3. **可扩展**：为未来功能扩展和微服务化预留空间
4. **自动化**：提供脚本自动创建和验证项目结构

**结构优势**：
- **开发效率**：清晰的结构提升开发效率
- **团队协作**：统一的结构便于团队协作
- **代码质量**：良好的组织有助于代码质量
- **项目维护**：结构化的代码易于维护和扩展

**下一步**：基于这个项目结构，我们将初始化Go模块，配置依赖管理和版本控制。
