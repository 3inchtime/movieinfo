# MovieInfo - 电影信息网站

[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)
[![Gin Framework](https://img.shields.io/badge/Gin-1.9+-green.svg)](https://gin-gonic.com)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## 项目简介

MovieInfo 是一个基于 Go 语言和 Gin 框架开发的电影信息网站，采用服务化架构设计。项目使用 Gin Template 作为前端模板引擎，提供用户友好的电影浏览、搜索、评论和评分功能。

### 核心特性

- 🎬 **电影信息管理** - 完整的电影详情、列表展示
- 👤 **用户系统** - 注册、登录、密码重置
- 💬 **评论系统** - 用户评论和评分功能
- 🏠 **主页服务** - 统一的用户访问入口
- 🔄 **服务化架构** - 四个独立服务，为后期微服务化做准备
- 🚀 **高性能** - gRPC 服务间通信，Redis 缓存加速

### 技术栈

- **后端框架**: Go + Gin
- **前端模板**: Gin Template
- **服务通信**: gRPC
- **数据库**: MySQL
- **缓存**: Redis
- **容器化**: Docker & Docker Compose

### 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   用户访问      │───▶│   主页服务      │───▶│   用户服务      │
│   (Browser)     │    │  (Web Gateway)  │    │ (User Service)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │                        │
                              │                        │
                              ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   电影服务      │    │  评论打分服务   │
                       │ (Movie Service) │    │(Comment Service)│
                       └─────────────────┘    └─────────────────┘
                              │                        │
                              └────────┬───────────────┘
                                       │
                                       ▼
                              ┌─────────────────┐
                              │   数据存储层    │
                              │  MySQL + Redis  │
                              └─────────────────┘
```



### 服务说明

| 服务名称 | 端口 | 功能描述 | 技术栈 |
|---------|------|----------|--------|
| **主页服务** | 8080 | Web网关，用户访问入口，页面渲染 | Gin + Gin Template |
| **用户服务** | 8081 | 用户注册、登录、密码管理 | Gin + gRPC |
| **电影服务** | 8082 | 电影信息管理、列表查询 | Gin + gRPC |
| **评论打分服务** | 8083 | 电影评论、评分功能 | Gin + gRPC |

## 快速开始

### 环境要求

- Go 1.19+
- MySQL 8.0+
- Redis 6.0+
- Docker & Docker Compose (可选)

### 项目结构

```
movieinfo/
├── cmd/                    # 应用程序入口
│   ├── web/               # 主页服务
│   ├── user/              # 用户服务
│   ├── movie/             # 电影服务
│   └── comment/           # 评论打分服务
├── internal/              # 内部包
│   ├── config/            # 配置管理
│   ├── models/            # 数据模型
│   ├── handlers/          # HTTP处理器
│   ├── services/          # 业务逻辑
│   └── middleware/        # 中间件
├── pkg/                   # 公共包
│   ├── database/          # 数据库连接
│   ├── redis/             # Redis连接
│   └── grpc/              # gRPC客户端
├── proto/                 # Protocol Buffers定义
├── templates/             # HTML模板
├── static/                # 静态资源
├── configs/               # 配置文件
├── docs/                  # 项目文档
├── go.mod
├── go.sum
└── README.md
```

### 本地开发

```bash
# 1. 克隆项目
git clone https://github.com/your-username/movieinfo.git
cd movieinfo

# 2. 初始化Go模块
go mod init movieinfo
go mod tidy

# 3. 启动数据库和缓存（需要先安装Docker）
# 创建docker-compose.yml文件后执行：
docker-compose up -d mysql redis

# 4. 配置数据库连接
# 复制配置文件模板并修改数据库连接信息
cp configs/config.example.yaml configs/config.yaml

# 5. 运行数据库迁移（需要先实现迁移脚本）
go run cmd/migrate/main.go

# 6. 启动各个服务（按顺序启动）
# 终端1: 启动用户服务
go run cmd/user/main.go

# 终端2: 启动电影服务
go run cmd/movie/main.go

# 终端3: 启动评论服务
go run cmd/comment/main.go

# 终端4: 启动主页服务
go run cmd/web/main.go
```

### 配置文件示例

创建 `configs/config.yaml` 文件：

```yaml
# 数据库配置
database:
  host: localhost
  port: 3306
  username: root
  password: your_password
  dbname: movieinfo
  charset: utf8mb4

# Redis配置
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

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
```

### 访问应用

启动成功后，可以通过以下地址访问：

- **主页服务**: http://localhost:8080
- **用户服务API**: http://localhost:8081 (gRPC: 9081)
- **电影服务API**: http://localhost:8082 (gRPC: 9082)
- **评论服务API**: http://localhost:8083 (gRPC: 9083)

## 开发文档目录

### 📋 1. 项目规划与设计
- [1.1 第1步：需求分析](docs/01-planning/01-requirements.md)
- [1.2 第2步：系统架构设计](docs/01-planning/02-architecture.md)
- [1.3 第3步：数据库设计](docs/01-planning/03-database-design.md)
- [1.4 第4步：API 接口设计](docs/01-planning/04-api-design.md)
- [1.5 第5步：UI/UX 设计](docs/01-planning/05-ui-design.md)

### 🛠️ 2. 开发环境搭建
- [2.1 第6步：开发环境准备](docs/02-setup/06-development-environment.md)
- [2.2 第7步：Go 环境配置](docs/02-setup/07-go-setup.md)
- [2.3 第8步：数据库环境搭建](docs/02-setup/08-database-setup.md)
- [2.4 第9步：Redis 缓存配置](docs/02-setup/09-redis-setup.md)
- [2.5 第10步：IDE 配置与插件](docs/02-setup/10-ide-configuration.md)

### 🏗️ 3. 项目初始化
- [3.1 第11步：项目结构创建](docs/03-initialization/11-project-structure.md)
- [3.2 第12步：Go Modules 初始化](docs/03-initialization/12-go-modules.md)
- [3.3 第13步：配置文件设计](docs/03-initialization/13-configuration.md)
- [3.4 第14步：日志系统搭建](docs/03-initialization/14-logging-system.md)
- [3.5 第15步：错误处理机制](docs/03-initialization/15-error-handling.md)

### 🗄️ 4. 数据层开发
- [4.1 第16步：数据库连接池](docs/04-data-layer/16-database-connection.md)
- [4.2 第17步：数据模型定义](docs/04-data-layer/17-data-models.md)
- [4.3 第18步：数据库迁移](docs/04-data-layer/18-database-migration.md)
- [4.4 第19步：CRUD 操作封装](docs/04-data-layer/19-crud-operations.md)
- [4.5 第20步：Redis 缓存集成](docs/04-data-layer/20-redis-integration.md)

### 🔧 5. gRPC 服务开发
- [5.1 第21步：Protocol Buffers 定义](docs/05-grpc/21-protobuf-definition.md)
- [5.2 第22步：gRPC 服务端实现](docs/05-grpc/22-server-implementation.md)
- [5.3 第23步：gRPC 客户端封装](docs/05-grpc/23-client-implementation.md)
- [5.4 第24步：服务注册与发现](docs/05-grpc/24-service-discovery.md)
- [5.5 第25步：gRPC 中间件](docs/05-grpc/25-middleware.md)

### 👤 6. 用户服务开发
- [6.1 第26步：用户模型设计](docs/06-user-service/26-user-model.md)
- [6.2 第27步：注册功能实现](docs/06-user-service/27-registration.md)
- [6.3 第28步：登录认证系统](docs/06-user-service/28-authentication.md)
- [6.4 第29步：JWT Token 管理](docs/06-user-service/29-jwt-management.md)
- [6.5 第30步：密码重置功能](docs/06-user-service/30-password-reset.md)

### 🎬 7. 电影服务开发
- [7.1 第31步：电影数据模型](docs/07-movie-service/31-movie-model.md)
- [7.2 第32步：电影列表接口](docs/07-movie-service/32-movie-list.md)
- [7.3 第33步：电影详情接口](docs/07-movie-service/33-movie-details.md)
- [7.4 第34步：电影搜索功能](docs/07-movie-service/34-movie-search.md)
- [7.5 第35步：电影分类管理](docs/07-movie-service/35-movie-categories.md)

### 💬 8. 评论打分服务开发
- [8.1 第36步：评论数据模型](docs/08-comment-service/36-comment-model.md)
- [8.2 第37步：评论 CRUD 接口](docs/08-comment-service/37-comment-crud.md)
- [8.3 第38步：评分系统实现](docs/08-comment-service/38-rating-system.md)
- [8.4 第39步：评论审核机制](docs/08-comment-service/39-comment-moderation.md)
- [8.5 第40步：统计分析功能](docs/08-comment-service/40-analytics.md)

### 🏠 9. 主页服务开发
- [9.1 第41步：Web 服务器搭建](docs/09-web-service/41-web-server.md)
- [9.2 第42步：路由设计](docs/09-web-service/42-routing.md)
- [9.3 第43步：模板引擎集成](docs/09-web-service/43-template-engine.md)
- [9.4 第44步：静态资源管理](docs/09-web-service/44-static-assets.md)
- [9.5 第45步：页面渲染逻辑](docs/09-web-service/45-page-rendering.md)

### 🎨 10. 前端页面开发
- [10.1 第46步：页面布局设计](docs/10-frontend/46-layout-design.md)
- [10.2 第47步：首页开发](docs/10-frontend/47-homepage.md)
- [10.3 第48步：电影列表页](docs/10-frontend/48-movie-list-page.md)
- [10.4 第49步：电影详情页](docs/10-frontend/49movie-detail-page.md)
- [10.5 第50步：用户相关页面](docs/10-frontend/50-user-pages.md)

### 🔒 11. 安全与认证
- [11.1 第51步：身份认证机制](docs/11-security/51-authentication.md)
- [11.2 第52步：权限控制系统](docs/11-security/52-authorization.md)
- [11.3 第53步：数据验证](docs/11-security/53-data-validation.md)
- [11.4 第54步：安全中间件](docs/11-security/54-security-middleware.md)
- [11.5 第55步：HTTPS 配置](docs/11-security/55-https-setup.md)

### 🧪 12. 项目完善与扩展
- [12.1 第56步：功能测试与验证](docs/12-enhancement/56-testing-validation.md)
- [12.2 第57步：性能调优](docs/12-enhancement/57-performance-tuning.md)
- [12.3 第58步：安全加固](docs/12-enhancement/58-security-hardening.md)
- [12.4 第59步：部署准备](docs/12-enhancement/59-deployment-preparation.md)
- [12.5 第60步：项目总结](docs/12-enhancement/60-project-summary.md)

> **注意**: 当前项目采用服务化架构，使用gRPC进行服务间通信，为后期转向微服务架构做好准备。测试开发、部署运维、监控日志、性能优化等高级主题将在后期版本中添加。

## 开发进度

当前项目处于核心功能开发阶段，按照以下顺序进行开发：

1. ✅ **项目规划与设计** - 完成需求分析和架构设计
2. 🔄 **开发环境搭建** - 正在进行中
3. ⏳ **项目初始化** - 待开始
4. ⏳ **数据层开发** - 待开始
5. ⏳ **gRPC服务开发** - 待开始
6. ⏳ **各服务开发** - 待开始
7. ⏳ **前端页面开发** - 待开始
8. ⏳ **安全与认证** - 待开始
9. ⏳ **项目完善与扩展** - 待开始

## 贡献指南

我们欢迎所有形式的贡献！请按照以下步骤参与项目开发：

### 开发流程

1. Fork 项目到你的GitHub账户
2. 创建功能分支 (`git checkout -b feature/功能名称`)
3. 按照开发文档目录逐步实现功能
4. 提交更改 (`git commit -m '添加某某功能'`)
5. 推送到分支 (`git push origin feature/功能名称`)
6. 创建 Pull Request

### 代码规范

- 遵循Go语言官方代码规范
- 使用有意义的变量和函数命名
- 添加必要的注释说明
- 确保代码可读性和可维护性

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

- 项目维护者: [Your Name](mailto:your.email@example.com)
- 项目地址: [https://github.com/your-username/movieinfo](https://github.com/your-username/movieinfo)
- 问题反馈: [Issues](https://github.com/your-username/movieinfo/issues)

## 致谢

感谢所有为这个项目做出贡献的开发者们！

---

⭐ 如果这个项目对你有帮助，请给我们一个 Star！