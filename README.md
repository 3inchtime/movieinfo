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
┌─────────────────┐    ┌─────────────────┐
│   用户访问      │───▶│   主页服务      │
│   (Browser)     │    │  (Web Gateway)  │
└─────────────────┘    └─────────────────┘
                              │
                    ┌─────────┼─────────┐
                    │         │         │
                    ▼         ▼         ▼
            ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
            │   用户服务  │ │   电影服务  │ │ 评论打分服务│
            │(User Service)│ │(Movie Service)│ │(Comment Service)│
            └─────────────┘ └─────────────┘ └─────────────┘
                    │         │         │
                    └─────────┼─────────┘
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

### 功能说明

#### 主页服务 (Web Gateway)
- **首页展示**: 展示热门电影、最新电影、推荐电影
- **电影搜索**: 支持按电影名称、类型、年份等条件搜索
- **电影详情**: 展示电影详细信息、评分、评论列表
- **用户界面**: 提供用户登录、注册、个人中心页面
- **页面路由**: 统一管理所有前端页面路由和模板渲染

#### 用户服务 (User Service)
- **用户注册**: 支持邮箱注册，密码加密存储
- **用户登录**: JWT token认证，支持记住登录状态
- **密码管理**: 密码重置、修改密码功能
- **用户信息**: 个人资料管理、头像上传
- **权限控制**: 基于角色的访问控制

#### 电影服务 (Movie Service)
- **电影管理**: 电影信息的增删改查
- **分类管理**: 电影类型、年份、地区分类
- **搜索功能**: 支持多条件组合搜索和模糊搜索
- **数据缓存**: 热门电影数据Redis缓存
- **分页查询**: 支持电影列表分页展示

#### 评论打分服务 (Comment Service)
- **评论功能**: 用户对电影发表评论
- **评分系统**: 1-10分评分，支持半分
- **评分统计**: 计算电影平均分、评分分布
- **评论管理**: 评论的审核、删除功能
- **互动功能**: 评论点赞、回复功能

### 访问应用

启动成功后，可以通过以下地址访问：

- **主页服务**: http://localhost:8080
- **用户服务API**: http://localhost:8081 (gRPC: 9081)
- **电影服务API**: http://localhost:8082 (gRPC: 9082)
- **评论服务API**: http://localhost:8083 (gRPC: 9083)

## 开发文档

### 📋 1. 项目设计
- [项目需求分析](docs/01-design/01-requirements.md)
- [系统架构设计](docs/01-design/02-architecture.md)
- [技术选型说明](docs/01-design/03-tech-stack.md)
- [数据库设计](docs/01-design/04-database-design.md)
- [API接口设计](docs/01-design/05-api-design.md)
- [缓存策略设计](docs/01-design/06-cache-strategy.md)

### 🏗️ 2. 基础设施
- [项目结构设计](docs/02-foundation/01-project-structure.md)
- [开发环境搭建](docs/02-foundation/02-dev-environment.md)
- [代码规范](docs/02-foundation/03-coding-standards.md)
- [数据库连接池](docs/02-foundation/04-database-connection.md)
- [GORM模型定义](docs/02-foundation/05-gorm-models.md)
- [gRPC框架配置](docs/02-foundation/06-grpc-setup.md)

### 👤 3. 用户服务
- [用户注册API](docs/03-user-service/01-registration.md)
- [用户认证系统](docs/03-user-service/02-authentication.md)
- [密码管理功能](docs/03-user-service/03-password-management.md)
- [用户资料管理](docs/03-user-service/04-profile-management.md)
- [权限控制系统](docs/03-user-service/05-access-control.md)

### 🎬 4. 电影服务
- [电影CRUD操作](docs/04-movie-service/01-movie-crud.md)
- [电影分类管理](docs/04-movie-service/02-category-management.md)
- [搜索功能实现](docs/04-movie-service/03-search-functionality.md)
- [Redis缓存策略](docs/04-movie-service/04-caching-strategy.md)
- [分页查询实现](docs/04-movie-service/05-pagination.md)

### 💬 5. 评论服务
- [评论功能开发](docs/05-comment-service/01-comment-system.md)
- [评分系统实现](docs/05-comment-service/02-rating-system.md)
- [评分统计计算](docs/05-comment-service/03-rating-statistics.md)
- [评论审核管理](docs/05-comment-service/04-comment-moderation.md)
- [互动功能开发](docs/05-comment-service/05-interaction-features.md)

### 🏠 6. 主页服务
- [Gin服务器配置](docs/06-web-service/01-server-config.md)
- [首页功能实现](docs/06-web-service/02-homepage.md)
- [搜索界面开发](docs/06-web-service/03-search-interface.md)
- [电影详情页面](docs/06-web-service/04-movie-detail.md)
- [用户界面开发](docs/06-web-service/05-user-interface.md)
- [模板引擎使用](docs/06-web-service/06-template-engine.md)
- [gRPC客户端集成](docs/06-web-service/07-grpc-client.md)

### 🔒 7. 安全与测试
- [安全机制配置](docs/07-security/01-security-config.md)
- [CORS跨域配置](docs/07-security/02-cors-setup.md)
- [输入验证机制](docs/07-security/03-input-validation.md)
- [单元测试编写](docs/07-security/04-unit-testing.md)
- [集成测试实现](docs/07-security/05-integration-testing.md)

### 🚀 8. 部署运维
- [Docker容器化](docs/08-deployment/01-docker-setup.md)
- [Docker Compose配置](docs/08-deployment/02-docker-compose.md)
- [生产环境部署](docs/08-deployment/03-production-deployment.md)
- [监控与日志](docs/08-deployment/04-monitoring-logging.md)
- [性能优化指南](docs/08-deployment/05-performance-optimization.md)
  

## 开发进度

当前项目处于核心功能开发阶段，按照以下顺序进行开发：

1. ✅ **项目设计与环境搭建** - 完成需求分析和架构设计
2. ⏳ **项目初始化与基础设施** - 待开始
3. ⏳ **用户服务开发** - 待开始
4. ⏳ **电影服务开发** - 待开始
5. ⏳ **评论打分服务开发** - 待开始
6. ⏳ **主页服务与前端开发** - 待开始
7. ⏳ **安全与部署** - 待开始

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