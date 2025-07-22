# 1.3 数据库设计

## 概述

数据库设计是系统架构的核心组成部分，它决定了数据的存储结构、访问效率和系统的可扩展性。对于MovieInfo项目，我们需要设计一个既能满足当前业务需求，又能支持未来扩展的数据库架构。

## 为什么数据库设计如此重要？

### 1. **数据完整性保证**
- **约束机制**：通过主键、外键、唯一约束等保证数据一致性
- **事务支持**：ACID特性确保数据操作的可靠性
- **数据验证**：在数据库层面进行数据格式和范围验证

### 2. **查询性能优化**
- **索引设计**：合理的索引提升查询效率
- **表结构优化**：规范化和反规范化的平衡
- **查询路径优化**：减少表连接和数据扫描

### 3. **系统扩展性**
- **水平扩展**：支持分库分表的扩展策略
- **垂直扩展**：支持读写分离和主从复制
- **业务扩展**：预留字段和表结构支持功能扩展

### 4. **维护便利性**
- **清晰的命名规范**：便于理解和维护
- **文档化设计**：完整的字段说明和关系描述
- **版本控制**：数据库变更的版本管理

## 数据库设计原则

### 1. **规范化原则**

#### 1.1 第一范式（1NF）
**定义**：每个字段都是原子性的，不可再分
**应用**：
```sql
-- 错误示例：演员字段包含多个值
actors VARCHAR(500) -- "张三,李四,王五"

-- 正确示例：使用关联表
CREATE TABLE movie_actors (
    movie_id INT,
    actor_name VARCHAR(100)
);
```

#### 1.2 第二范式（2NF）
**定义**：满足1NF，且非主键字段完全依赖于主键
**应用**：避免部分依赖，将相关字段分离到独立表

#### 1.3 第三范式（3NF）
**定义**：满足2NF，且非主键字段不依赖于其他非主键字段
**应用**：消除传递依赖，减少数据冗余

### 2. **性能优化原则**

#### 2.1 适度反规范化
在某些场景下，为了查询性能，可以适当冗余数据：
```sql
-- 在电影表中冗余评分信息，避免每次都计算
CREATE TABLE movies (
    id INT PRIMARY KEY,
    title VARCHAR(200),
    average_rating DECIMAL(3,2), -- 冗余字段
    rating_count INT             -- 冗余字段
);
```

#### 2.2 索引设计策略
- **主键索引**：每个表必须有主键
- **唯一索引**：保证数据唯一性
- **复合索引**：多字段查询优化
- **覆盖索引**：减少回表查询

### 3. **扩展性原则**

#### 3.1 预留扩展字段
```sql
CREATE TABLE users (
    id INT PRIMARY KEY,
    -- 核心字段
    email VARCHAR(100),
    username VARCHAR(50),
    -- 预留扩展字段
    ext_field1 VARCHAR(255),
    ext_field2 TEXT,
    metadata JSON
);
```

#### 3.2 软删除设计
```sql
-- 使用deleted_at字段实现软删除
deleted_at TIMESTAMP NULL DEFAULT NULL
```

## 核心表设计

### 1. 用户相关表

#### 1.1 用户主表（users）
```sql
CREATE TABLE users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱地址',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    avatar_url VARCHAR(255) DEFAULT NULL COMMENT '头像URL',
    bio TEXT DEFAULT NULL COMMENT '个人简介',
    status ENUM('active', 'inactive', 'banned') DEFAULT 'active' COMMENT '用户状态',
    email_verified_at TIMESTAMP NULL DEFAULT NULL COMMENT '邮箱验证时间',
    last_login_at TIMESTAMP NULL DEFAULT NULL COMMENT '最后登录时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',
    
    INDEX idx_email (email),
    INDEX idx_username (username),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
```

**设计说明**：
- **主键设计**：使用自增INT作为主键，性能好且便于分库分表
- **唯一约束**：email和username都设置唯一约束，支持两种登录方式
- **密码安全**：存储密码哈希而非明文，使用bcrypt算法
- **状态管理**：支持用户状态管理，便于后台管理
- **软删除**：使用deleted_at实现软删除，保留数据完整性

#### 1.2 用户会话表（user_sessions）
```sql
CREATE TABLE user_sessions (
    id VARCHAR(128) PRIMARY KEY COMMENT '会话ID',
    user_id INT UNSIGNED NOT NULL COMMENT '用户ID',
    ip_address VARCHAR(45) DEFAULT NULL COMMENT 'IP地址',
    user_agent TEXT DEFAULT NULL COMMENT '用户代理',
    payload TEXT NOT NULL COMMENT '会话数据',
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '最后活动时间',
    expires_at TIMESTAMP NOT NULL COMMENT '过期时间',
    
    INDEX idx_user_id (user_id),
    INDEX idx_last_activity (last_activity),
    INDEX idx_expires_at (expires_at),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户会话表';
```

#### 1.3 密码重置表（password_resets）
```sql
CREATE TABLE password_resets (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(100) NOT NULL COMMENT '邮箱地址',
    token VARCHAR(255) NOT NULL COMMENT '重置令牌',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    expires_at TIMESTAMP NOT NULL COMMENT '过期时间',
    used_at TIMESTAMP NULL DEFAULT NULL COMMENT '使用时间',
    
    INDEX idx_email (email),
    INDEX idx_token (token),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='密码重置表';
```

### 2. 电影相关表

#### 2.1 电影主表（movies）
```sql
CREATE TABLE movies (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '电影ID',
    title VARCHAR(200) NOT NULL COMMENT '电影标题',
    original_title VARCHAR(200) DEFAULT NULL COMMENT '原始标题',
    director VARCHAR(100) DEFAULT NULL COMMENT '导演',
    release_year YEAR DEFAULT NULL COMMENT '上映年份',
    duration INT DEFAULT NULL COMMENT '时长(分钟)',
    plot TEXT DEFAULT NULL COMMENT '剧情简介',
    poster_url VARCHAR(255) DEFAULT NULL COMMENT '海报URL',
    trailer_url VARCHAR(255) DEFAULT NULL COMMENT '预告片URL',
    imdb_id VARCHAR(20) DEFAULT NULL COMMENT 'IMDB ID',
    tmdb_id INT DEFAULT NULL COMMENT 'TMDB ID',
    
    -- 评分相关字段（冗余设计，提升查询性能）
    average_rating DECIMAL(3,2) DEFAULT 0.00 COMMENT '平均评分',
    rating_count INT UNSIGNED DEFAULT 0 COMMENT '评分人数',
    
    -- 状态和元数据
    status ENUM('draft', 'published', 'archived') DEFAULT 'published' COMMENT '状态',
    view_count INT UNSIGNED DEFAULT 0 COMMENT '浏览次数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',
    
    INDEX idx_title (title),
    INDEX idx_director (director),
    INDEX idx_release_year (release_year),
    INDEX idx_average_rating (average_rating),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    UNIQUE KEY uk_imdb_id (imdb_id),
    UNIQUE KEY uk_tmdb_id (tmdb_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='电影表';
```

**设计说明**：
- **标题字段**：支持中文标题和原始标题
- **外部ID**：支持IMDB和TMDB的ID，便于数据同步
- **评分冗余**：在电影表中冗余评分信息，避免频繁计算
- **状态管理**：支持草稿、发布、归档等状态
- **性能优化**：为常用查询字段建立索引

#### 2.2 电影分类表（categories）
```sql
CREATE TABLE categories (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '分类ID',
    name VARCHAR(50) NOT NULL COMMENT '分类名称',
    slug VARCHAR(50) NOT NULL UNIQUE COMMENT '分类别名',
    description TEXT DEFAULT NULL COMMENT '分类描述',
    sort_order INT DEFAULT 0 COMMENT '排序权重',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_slug (slug),
    INDEX idx_sort_order (sort_order),
    INDEX idx_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='电影分类表';
```

#### 2.3 电影分类关联表（movie_categories）
```sql
CREATE TABLE movie_categories (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    movie_id INT UNSIGNED NOT NULL COMMENT '电影ID',
    category_id INT UNSIGNED NOT NULL COMMENT '分类ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    UNIQUE KEY uk_movie_category (movie_id, category_id),
    INDEX idx_movie_id (movie_id),
    INDEX idx_category_id (category_id),
    
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='电影分类关联表';
```

#### 2.4 演员表（actors）
```sql
CREATE TABLE actors (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '演员ID',
    name VARCHAR(100) NOT NULL COMMENT '演员姓名',
    original_name VARCHAR(100) DEFAULT NULL COMMENT '原始姓名',
    gender ENUM('male', 'female', 'other') DEFAULT NULL COMMENT '性别',
    birth_date DATE DEFAULT NULL COMMENT '出生日期',
    nationality VARCHAR(50) DEFAULT NULL COMMENT '国籍',
    biography TEXT DEFAULT NULL COMMENT '个人简介',
    photo_url VARCHAR(255) DEFAULT NULL COMMENT '照片URL',
    imdb_id VARCHAR(20) DEFAULT NULL COMMENT 'IMDB ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_name (name),
    INDEX idx_nationality (nationality),
    UNIQUE KEY uk_imdb_id (imdb_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='演员表';
```

#### 2.5 电影演员关联表（movie_actors）
```sql
CREATE TABLE movie_actors (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    movie_id INT UNSIGNED NOT NULL COMMENT '电影ID',
    actor_id INT UNSIGNED NOT NULL COMMENT '演员ID',
    role VARCHAR(100) DEFAULT NULL COMMENT '角色名称',
    role_type ENUM('lead', 'supporting', 'cameo') DEFAULT 'supporting' COMMENT '角色类型',
    sort_order INT DEFAULT 0 COMMENT '排序权重',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    UNIQUE KEY uk_movie_actor_role (movie_id, actor_id, role),
    INDEX idx_movie_id (movie_id),
    INDEX idx_actor_id (actor_id),
    INDEX idx_role_type (role_type),
    
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE,
    FOREIGN KEY (actor_id) REFERENCES actors(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='电影演员关联表';
```

### 3. 评论评分相关表

#### 3.1 评分表（ratings）
```sql
CREATE TABLE ratings (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '评分ID',
    user_id INT UNSIGNED NOT NULL COMMENT '用户ID',
    movie_id INT UNSIGNED NOT NULL COMMENT '电影ID',
    score TINYINT UNSIGNED NOT NULL COMMENT '评分(1-5)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    UNIQUE KEY uk_user_movie (user_id, movie_id),
    INDEX idx_movie_id (movie_id),
    INDEX idx_score (score),
    INDEX idx_created_at (created_at),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE,
    
    CHECK (score >= 1 AND score <= 5)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='电影评分表';
```

#### 3.2 评论表（comments）
```sql
CREATE TABLE comments (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '评论ID',
    user_id INT UNSIGNED NOT NULL COMMENT '用户ID',
    movie_id INT UNSIGNED NOT NULL COMMENT '电影ID',
    parent_id INT UNSIGNED DEFAULT NULL COMMENT '父评论ID',
    content TEXT NOT NULL COMMENT '评论内容',
    status ENUM('pending', 'approved', 'rejected', 'hidden') DEFAULT 'pending' COMMENT '审核状态',
    like_count INT UNSIGNED DEFAULT 0 COMMENT '点赞数',
    dislike_count INT UNSIGNED DEFAULT 0 COMMENT '踩数',
    report_count INT UNSIGNED DEFAULT 0 COMMENT '举报次数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',
    
    INDEX idx_user_id (user_id),
    INDEX idx_movie_id (movie_id),
    INDEX idx_parent_id (parent_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='电影评论表';
```

#### 3.3 评论点赞表（comment_likes）
```sql
CREATE TABLE comment_likes (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNSIGNED NOT NULL COMMENT '用户ID',
    comment_id INT UNSIGNED NOT NULL COMMENT '评论ID',
    type ENUM('like', 'dislike') NOT NULL COMMENT '类型',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    UNIQUE KEY uk_user_comment (user_id, comment_id),
    INDEX idx_comment_id (comment_id),
    INDEX idx_type (type),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论点赞表';
```

## 索引设计策略

### 1. 主键索引
每个表都使用自增INT作为主键，提供最佳的插入性能和存储效率。

### 2. 唯一索引
```sql
-- 用户表
UNIQUE KEY uk_email (email)
UNIQUE KEY uk_username (username)

-- 电影表
UNIQUE KEY uk_imdb_id (imdb_id)
UNIQUE KEY uk_tmdb_id (tmdb_id)

-- 评分表
UNIQUE KEY uk_user_movie (user_id, movie_id)
```

### 3. 复合索引
```sql
-- 电影查询优化
INDEX idx_status_rating (status, average_rating)
INDEX idx_category_year (category_id, release_year)

-- 评论查询优化
INDEX idx_movie_status_created (movie_id, status, created_at)
```

### 4. 覆盖索引
```sql
-- 电影列表查询覆盖索引
INDEX idx_movie_list (status, average_rating, id, title, poster_url)
```

## 数据完整性约束

### 1. 外键约束
```sql
-- 确保数据引用完整性
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE
```

### 2. 检查约束
```sql
-- 评分范围约束
CHECK (score >= 1 AND score <= 5)

-- 状态值约束
status ENUM('active', 'inactive', 'banned')
```

### 3. 非空约束
```sql
-- 关键字段不允许为空
email VARCHAR(100) NOT NULL
title VARCHAR(200) NOT NULL
content TEXT NOT NULL
```

## 性能优化策略

### 1. 分区策略
```sql
-- 按时间分区评论表
CREATE TABLE comments (
    ...
) PARTITION BY RANGE (YEAR(created_at)) (
    PARTITION p2023 VALUES LESS THAN (2024),
    PARTITION p2024 VALUES LESS THAN (2025),
    PARTITION p_future VALUES LESS THAN MAXVALUE
);
```

### 2. 读写分离准备
```sql
-- 为读写分离预留配置
-- 主库：写操作
-- 从库：读操作
```

### 3. 缓存友好设计
```sql
-- 设计缓存键友好的查询
SELECT id, title, poster_url, average_rating 
FROM movies 
WHERE status = 'published' 
ORDER BY average_rating DESC 
LIMIT 20;
```

## 数据迁移策略

### 1. 版本控制
```sql
CREATE TABLE schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 2. 迁移脚本示例
```sql
-- V001_create_users_table.sql
CREATE TABLE users (...);

-- V002_create_movies_table.sql  
CREATE TABLE movies (...);

-- V003_add_rating_index.sql
CREATE INDEX idx_average_rating ON movies(average_rating);
```

## 数据库连接和配置

### 1. 连接池配置
```go
// 数据库连接池配置
type DatabaseConfig struct {
    Host            string `yaml:"host"`
    Port            int    `yaml:"port"`
    Username        string `yaml:"username"`
    Password        string `yaml:"password"`
    Database        string `yaml:"database"`
    Charset         string `yaml:"charset"`
    MaxOpenConns    int    `yaml:"max_open_conns"`
    MaxIdleConns    int    `yaml:"max_idle_conns"`
    ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

// 推荐配置值
MaxOpenConns:    100  // 最大打开连接数
MaxIdleConns:    10   // 最大空闲连接数
ConnMaxLifetime: 3600 // 连接最大生存时间(秒)
```

### 2. 字符集和排序规则
```sql
-- 数据库级别
CREATE DATABASE movieinfo
CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;

-- 表级别
DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
```

**选择原因**：
- **utf8mb4**：支持完整的UTF-8字符集，包括emoji
- **unicode_ci**：不区分大小写的Unicode排序规则

## 数据安全和备份策略

### 1. 数据加密
```sql
-- 敏感字段加密存储
password_hash VARCHAR(255) NOT NULL  -- bcrypt哈希
email_token VARCHAR(255)             -- 邮箱验证令牌
reset_token VARCHAR(255)             -- 密码重置令牌
```

### 2. 备份策略
```bash
# 每日全量备份
mysqldump --single-transaction --routines --triggers \
  --all-databases > backup_$(date +%Y%m%d).sql

# 增量备份（二进制日志）
mysqlbinlog --start-datetime="2024-01-01 00:00:00" \
  --stop-datetime="2024-01-01 23:59:59" \
  mysql-bin.000001 > incremental_backup.sql
```

### 3. 数据恢复
```bash
# 全量恢复
mysql < backup_20240101.sql

# 增量恢复
mysql < incremental_backup.sql
```

## 监控和维护

### 1. 性能监控指标
```sql
-- 查询性能监控
SHOW PROCESSLIST;
SHOW ENGINE INNODB STATUS;

-- 索引使用情况
SELECT * FROM sys.schema_unused_indexes;
SELECT * FROM sys.schema_redundant_indexes;

-- 表大小监控
SELECT
    table_name,
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size (MB)'
FROM information_schema.tables
WHERE table_schema = 'movieinfo'
ORDER BY (data_length + index_length) DESC;
```

### 2. 慢查询优化
```sql
-- 开启慢查询日志
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1;

-- 分析慢查询
SELECT * FROM mysql.slow_log
WHERE start_time > DATE_SUB(NOW(), INTERVAL 1 DAY);
```

### 3. 定期维护任务
```sql
-- 表优化
OPTIMIZE TABLE movies;
OPTIMIZE TABLE comments;

-- 索引重建
ALTER TABLE movies DROP INDEX idx_title, ADD INDEX idx_title (title);

-- 统计信息更新
ANALYZE TABLE movies;
```

## 扩展性考虑

### 1. 分库分表准备
```sql
-- 用户表按用户ID分表
users_0, users_1, users_2, users_3

-- 评论表按电影ID分表
comments_movie_0, comments_movie_1, comments_movie_2

-- 分片键选择原则
- 用户相关：按user_id分片
- 电影相关：按movie_id分片
- 时间相关：按created_at分片
```

### 2. 读写分离准备
```go
// 数据库路由配置
type DatabaseRouter struct {
    Master DatabaseConfig  // 主库配置
    Slaves []DatabaseConfig // 从库配置列表
}

// 读写分离逻辑
func (r *DatabaseRouter) GetDB(operation string) *sql.DB {
    if operation == "write" {
        return r.Master.DB
    }
    // 负载均衡选择从库
    return r.selectSlave()
}
```

### 3. 缓存层设计
```go
// 缓存键设计
type CacheKeys struct {
    User   string // "user:{id}"
    Movie  string // "movie:{id}"
    Rating string // "rating:movie:{movie_id}"
    Search string // "search:{hash}"
}

// 缓存更新策略
func UpdateMovieCache(movieID int) {
    // 删除相关缓存
    cache.Delete(fmt.Sprintf("movie:%d", movieID))
    cache.Delete("movies:hot")
    cache.Delete("movies:latest")
}
```

## 数据初始化

### 1. 基础数据
```sql
-- 电影分类初始数据
INSERT INTO categories (name, slug, sort_order) VALUES
('动作', 'action', 1),
('喜剧', 'comedy', 2),
('剧情', 'drama', 3),
('科幻', 'sci-fi', 4),
('恐怖', 'horror', 5),
('爱情', 'romance', 6),
('悬疑', 'thriller', 7),
('动画', 'animation', 8);

-- 管理员用户
INSERT INTO users (email, username, password_hash, status) VALUES
('admin@movieinfo.com', 'admin', '$2a$10$...', 'active');
```

### 2. 测试数据
```sql
-- 测试电影数据
INSERT INTO movies (title, director, release_year, plot, average_rating, rating_count) VALUES
('肖申克的救赎', '弗兰克·德拉邦特', 1994, '一个关于希望和友谊的故事...', 4.8, 1000),
('霸王别姬', '陈凯歌', 1993, '一部关于京剧演员的史诗...', 4.7, 800),
('阿甘正传', '罗伯特·泽米吉斯', 1994, '一个智商不高但心地善良的男子...', 4.6, 900);
```

## 总结

数据库设计为MovieInfo项目提供了坚实的数据基础。通过规范化设计保证了数据完整性，通过性能优化确保了查询效率，通过扩展性设计为未来发展预留了空间。

**关键设计特点**：
1. **规范化与性能平衡**：在保证数据完整性的同时优化查询性能
2. **索引策略完善**：为常用查询建立合适的索引
3. **约束机制健全**：通过各种约束保证数据质量
4. **扩展性良好**：支持未来的功能扩展和性能优化
5. **安全性考虑**：敏感数据加密，完善的备份策略
6. **监控维护**：提供完整的监控和维护方案

**设计优势**：
- **高性能**：合理的索引和查询优化
- **高可用**：完善的备份和恢复机制
- **可扩展**：支持分库分表和读写分离
- **易维护**：清晰的表结构和命名规范

**下一步**：基于这个数据库设计，我们将进行API接口设计，确保接口能够高效地操作数据库，为前端提供稳定可靠的数据服务。
