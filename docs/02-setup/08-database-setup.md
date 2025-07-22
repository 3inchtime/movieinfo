# 2.3 数据库环境搭建

## 概述

数据库是MovieInfo项目的数据存储核心，承载着用户信息、电影数据、评论评分等关键业务数据。我们需要搭建一个稳定、高性能、易维护的数据库环境，包括MySQL主数据库和Redis缓存数据库。

## 为什么选择MySQL + Redis架构？

### 1. **MySQL作为主数据库**
- **ACID特性**：保证数据的一致性和可靠性
- **成熟稳定**：经过长期验证的关系型数据库
- **丰富功能**：支持复杂查询、事务、索引优化
- **生态完善**：工具链丰富，社区支持好
- **扩展性**：支持主从复制、分库分表

### 2. **Redis作为缓存数据库**
- **高性能**：内存存储，微秒级响应时间
- **数据结构丰富**：支持字符串、哈希、列表、集合等
- **持久化**：支持RDB和AOF持久化
- **高可用**：支持主从复制和哨兵模式
- **原子操作**：支持复杂的原子操作

### 3. **架构优势**
- **读写分离**：MySQL负责持久化存储，Redis负责缓存
- **性能优化**：热点数据缓存，减少数据库压力
- **扩展性**：独立扩展存储和缓存层
- **容错性**：缓存失效不影响核心功能

## MySQL环境搭建

### 1. **MySQL安装**

#### 1.1 Docker方式安装（推荐）
```bash
# 拉取MySQL镜像
docker pull mysql:8.0

# 创建数据目录
mkdir -p ~/docker/mysql/{data,conf,logs}

# 创建MySQL配置文件
cat > ~/docker/mysql/conf/my.cnf << EOF
[mysqld]
# 基础配置
port = 3306
bind-address = 0.0.0.0
server-id = 1

# 字符集配置
character-set-server = utf8mb4
collation-server = utf8mb4_unicode_ci
init-connect = 'SET NAMES utf8mb4'

# InnoDB配置
default-storage-engine = INNODB
innodb_buffer_pool_size = 1G
innodb_log_file_size = 256M
innodb_file_per_table = 1
innodb_flush_log_at_trx_commit = 2

# 连接配置
max_connections = 1000
max_connect_errors = 1000
wait_timeout = 28800
interactive_timeout = 28800

# 查询缓存
query_cache_type = 1
query_cache_size = 128M

# 慢查询日志
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 2

# 二进制日志
log-bin = mysql-bin
binlog_format = ROW
expire_logs_days = 7

# 安全配置
sql_mode = STRICT_TRANS_TABLES,NO_ZERO_DATE,NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO

[mysql]
default-character-set = utf8mb4

[client]
default-character-set = utf8mb4
EOF

# 启动MySQL容器
docker run -d \
  --name movieinfo-mysql \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=movieinfo123 \
  -e MYSQL_DATABASE=movieinfo \
  -e MYSQL_USER=movieinfo \
  -e MYSQL_PASSWORD=movieinfo123 \
  -v ~/docker/mysql/data:/var/lib/mysql \
  -v ~/docker/mysql/conf/my.cnf:/etc/mysql/conf.d/my.cnf \
  -v ~/docker/mysql/logs:/var/log/mysql \
  --restart=unless-stopped \
  mysql:8.0
```

#### 1.2 系统包管理器安装

**Ubuntu/Debian**：
```bash
# 更新包列表
sudo apt update

# 安装MySQL服务器
sudo apt install mysql-server

# 安全配置
sudo mysql_secure_installation

# 启动服务
sudo systemctl start mysql
sudo systemctl enable mysql
```

**CentOS/RHEL**：
```bash
# 安装MySQL仓库
sudo yum install mysql-server

# 启动服务
sudo systemctl start mysqld
sudo systemctl enable mysqld

# 获取临时密码
sudo grep 'temporary password' /var/log/mysqld.log

# 安全配置
sudo mysql_secure_installation
```

**macOS**：
```bash
# 使用Homebrew安装
brew install mysql

# 启动服务
brew services start mysql

# 安全配置
mysql_secure_installation
```

### 2. **MySQL配置优化**

#### 2.1 性能配置
```ini
# /etc/mysql/mysql.conf.d/mysqld.cnf 或 my.cnf

[mysqld]
# 内存配置
innodb_buffer_pool_size = 2G          # 设置为系统内存的70-80%
innodb_log_buffer_size = 64M
key_buffer_size = 256M
sort_buffer_size = 4M
read_buffer_size = 2M
read_rnd_buffer_size = 8M
myisam_sort_buffer_size = 128M

# 连接配置
max_connections = 1000
max_user_connections = 950
max_connect_errors = 1000
connect_timeout = 10
wait_timeout = 28800
interactive_timeout = 28800

# InnoDB配置
innodb_file_per_table = 1
innodb_flush_log_at_trx_commit = 2     # 性能优化，可接受少量数据丢失
innodb_log_file_size = 512M
innodb_log_files_in_group = 2
innodb_max_dirty_pages_pct = 75
innodb_lock_wait_timeout = 120

# 查询缓存
query_cache_type = 1
query_cache_size = 256M
query_cache_limit = 2M

# 临时表配置
tmp_table_size = 256M
max_heap_table_size = 256M

# 慢查询配置
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 1
log_queries_not_using_indexes = 1

# 二进制日志
log-bin = mysql-bin
binlog_format = ROW
max_binlog_size = 1G
expire_logs_days = 7
sync_binlog = 1000                     # 性能优化

# 字符集
character-set-server = utf8mb4
collation-server = utf8mb4_unicode_ci
```

#### 2.2 安全配置
```sql
-- 创建应用数据库用户
CREATE USER 'movieinfo'@'%' IDENTIFIED BY 'movieinfo_secure_password';
CREATE DATABASE movieinfo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON movieinfo.* TO 'movieinfo'@'%';

-- 创建只读用户（用于读写分离）
CREATE USER 'movieinfo_readonly'@'%' IDENTIFIED BY 'readonly_password';
GRANT SELECT ON movieinfo.* TO 'movieinfo_readonly'@'%';

-- 刷新权限
FLUSH PRIVILEGES;

-- 删除匿名用户
DELETE FROM mysql.user WHERE User='';

-- 删除test数据库
DROP DATABASE IF EXISTS test;

-- 禁用远程root登录
DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1');
```

### 3. **MySQL监控和维护**

#### 3.1 性能监控
```sql
-- 查看连接状态
SHOW PROCESSLIST;
SHOW STATUS LIKE 'Threads_connected';
SHOW STATUS LIKE 'Max_used_connections';

-- 查看缓存命中率
SHOW STATUS LIKE 'Qcache_hits';
SHOW STATUS LIKE 'Qcache_inserts';

-- 查看InnoDB状态
SHOW ENGINE INNODB STATUS;

-- 查看慢查询
SELECT * FROM mysql.slow_log ORDER BY start_time DESC LIMIT 10;
```

#### 3.2 维护脚本
```bash
#!/bin/bash
# mysql-maintenance.sh

# 数据库连接信息
DB_HOST="localhost"
DB_USER="root"
DB_PASS="password"
DB_NAME="movieinfo"

# 备份目录
BACKUP_DIR="/backup/mysql"
DATE=$(date +%Y%m%d_%H%M%S)

# 创建备份目录
mkdir -p $BACKUP_DIR

# 全量备份
echo "Starting full backup..."
mysqldump -h$DB_HOST -u$DB_USER -p$DB_PASS \
  --single-transaction \
  --routines \
  --triggers \
  --all-databases \
  --master-data=2 \
  | gzip > $BACKUP_DIR/full_backup_$DATE.sql.gz

# 表优化
echo "Optimizing tables..."
mysql -h$DB_HOST -u$DB_USER -p$DB_PASS -e "
USE $DB_NAME;
OPTIMIZE TABLE users;
OPTIMIZE TABLE movies;
OPTIMIZE TABLE comments;
OPTIMIZE TABLE ratings;
"

# 清理旧的二进制日志
echo "Purging old binary logs..."
mysql -h$DB_HOST -u$DB_USER -p$DB_PASS -e "PURGE BINARY LOGS BEFORE DATE_SUB(NOW(), INTERVAL 7 DAY);"

# 删除7天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete

echo "Maintenance completed!"
```

## Redis环境搭建

### 1. **Redis安装**

#### 1.1 Docker方式安装（推荐）
```bash
# 拉取Redis镜像
docker pull redis:7-alpine

# 创建配置目录
mkdir -p ~/docker/redis/{data,conf}

# 创建Redis配置文件
cat > ~/docker/redis/conf/redis.conf << EOF
# 网络配置
bind 0.0.0.0
port 6379
protected-mode yes
requirepass movieinfo_redis_password

# 内存配置
maxmemory 2gb
maxmemory-policy allkeys-lru

# 持久化配置
save 900 1
save 300 10
save 60 10000
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
dir /data

# AOF配置
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb

# 日志配置
loglevel notice
logfile /data/redis.log

# 客户端配置
timeout 300
tcp-keepalive 300
tcp-backlog 511

# 慢查询配置
slowlog-log-slower-than 10000
slowlog-max-len 128

# 键空间通知
notify-keyspace-events Ex
EOF

# 启动Redis容器
docker run -d \
  --name movieinfo-redis \
  -p 6379:6379 \
  -v ~/docker/redis/data:/data \
  -v ~/docker/redis/conf/redis.conf:/usr/local/etc/redis/redis.conf \
  --restart=unless-stopped \
  redis:7-alpine redis-server /usr/local/etc/redis/redis.conf
```

#### 1.2 系统包管理器安装

**Ubuntu/Debian**：
```bash
# 安装Redis
sudo apt update
sudo apt install redis-server

# 配置Redis
sudo nano /etc/redis/redis.conf

# 启动服务
sudo systemctl start redis-server
sudo systemctl enable redis-server
```

**CentOS/RHEL**：
```bash
# 安装EPEL仓库
sudo yum install epel-release

# 安装Redis
sudo yum install redis

# 启动服务
sudo systemctl start redis
sudo systemctl enable redis
```

**macOS**：
```bash
# 使用Homebrew安装
brew install redis

# 启动服务
brew services start redis
```

### 2. **Redis配置优化**

#### 2.1 性能配置
```conf
# /etc/redis/redis.conf

# 内存优化
maxmemory 4gb
maxmemory-policy allkeys-lru
maxmemory-samples 5

# 网络优化
tcp-backlog 511
tcp-keepalive 300
timeout 0

# 持久化优化
save 900 1
save 300 10
save 60 10000
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes

# AOF优化
appendonly yes
appendfsync everysec
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
aof-load-truncated yes

# 客户端优化
maxclients 10000

# 慢查询优化
slowlog-log-slower-than 10000
slowlog-max-len 128

# 键过期优化
hz 10
```

#### 2.2 安全配置
```conf
# 密码认证
requirepass your_strong_password

# 绑定地址
bind 127.0.0.1 192.168.1.100

# 保护模式
protected-mode yes

# 禁用危险命令
rename-command FLUSHDB ""
rename-command FLUSHALL ""
rename-command KEYS ""
rename-command CONFIG "CONFIG_9a8b7c6d"
rename-command SHUTDOWN "SHUTDOWN_1a2b3c4d"
rename-command DEBUG ""
rename-command EVAL ""
```

### 3. **Redis监控和维护**

#### 3.1 性能监控
```bash
# 连接Redis
redis-cli -h localhost -p 6379 -a password

# 查看信息
INFO memory
INFO stats
INFO replication
INFO persistence

# 监控命令
MONITOR

# 慢查询日志
SLOWLOG GET 10

# 客户端连接
CLIENT LIST
```

#### 3.2 维护脚本
```bash
#!/bin/bash
# redis-maintenance.sh

REDIS_HOST="localhost"
REDIS_PORT="6379"
REDIS_PASSWORD="password"
BACKUP_DIR="/backup/redis"
DATE=$(date +%Y%m%d_%H%M%S)

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份Redis数据
echo "Starting Redis backup..."
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASSWORD --rdb $BACKUP_DIR/dump_$DATE.rdb

# 清理过期键
echo "Cleaning expired keys..."
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASSWORD EVAL "
local keys = redis.call('keys', ARGV[1])
for i=1,#keys do
    redis.call('expire', keys[i], 0)
end
return #keys
" 0 "*:expired:*"

# 内存碎片整理
echo "Defragmenting memory..."
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASSWORD MEMORY PURGE

# 删除7天前的备份
find $BACKUP_DIR -name "dump_*.rdb" -mtime +7 -delete

echo "Redis maintenance completed!"
```

## 数据库连接测试

### 1. **MySQL连接测试**
```go
// test-mysql.go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // 数据库连接字符串
    dsn := "movieinfo:movieinfo123@tcp(localhost:3306)/movieinfo?charset=utf8mb4&parseTime=True&loc=Local"
    
    // 连接数据库
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    // 测试连接
    err = db.Ping()
    if err != nil {
        log.Fatal("Failed to ping database:", err)
    }
    
    // 查询版本
    var version string
    err = db.QueryRow("SELECT VERSION()").Scan(&version)
    if err != nil {
        log.Fatal("Failed to query version:", err)
    }
    
    fmt.Printf("✅ MySQL connected successfully!\n")
    fmt.Printf("Version: %s\n", version)
    
    // 测试创建表
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS test_connection (
            id INT AUTO_INCREMENT PRIMARY KEY,
            message VARCHAR(255),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        log.Fatal("Failed to create test table:", err)
    }
    
    // 插入测试数据
    _, err = db.Exec("INSERT INTO test_connection (message) VALUES (?)", "Connection test successful")
    if err != nil {
        log.Fatal("Failed to insert test data:", err)
    }
    
    // 查询测试数据
    var message string
    err = db.QueryRow("SELECT message FROM test_connection ORDER BY id DESC LIMIT 1").Scan(&message)
    if err != nil {
        log.Fatal("Failed to query test data:", err)
    }
    
    fmt.Printf("Test message: %s\n", message)
    
    // 清理测试表
    _, err = db.Exec("DROP TABLE test_connection")
    if err != nil {
        log.Printf("Warning: Failed to drop test table: %v", err)
    }
}
```

### 2. **Redis连接测试**
```go
// test-redis.go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/redis/go-redis/v9"
)

func main() {
    // Redis客户端配置
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "movieinfo_redis_password",
        DB:       0,
    })
    defer rdb.Close()
    
    ctx := context.Background()
    
    // 测试连接
    pong, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatal("Failed to connect to Redis:", err)
    }
    
    fmt.Printf("✅ Redis connected successfully!\n")
    fmt.Printf("Ping response: %s\n", pong)
    
    // 测试基本操作
    // 设置键值
    err = rdb.Set(ctx, "test:connection", "Redis connection test", time.Minute).Err()
    if err != nil {
        log.Fatal("Failed to set key:", err)
    }
    
    // 获取值
    val, err := rdb.Get(ctx, "test:connection").Result()
    if err != nil {
        log.Fatal("Failed to get key:", err)
    }
    fmt.Printf("Test value: %s\n", val)
    
    // 测试哈希操作
    err = rdb.HSet(ctx, "test:hash", map[string]interface{}{
        "field1": "value1",
        "field2": "value2",
    }).Err()
    if err != nil {
        log.Fatal("Failed to set hash:", err)
    }
    
    // 获取哈希
    hash, err := rdb.HGetAll(ctx, "test:hash").Result()
    if err != nil {
        log.Fatal("Failed to get hash:", err)
    }
    fmt.Printf("Test hash: %v\n", hash)
    
    // 清理测试数据
    rdb.Del(ctx, "test:connection", "test:hash")
    
    // 获取Redis信息
    info, err := rdb.Info(ctx, "server").Result()
    if err != nil {
        log.Printf("Warning: Failed to get Redis info: %v", err)
    } else {
        fmt.Printf("Redis server info retrieved successfully\n")
    }
}
```

### 3. **运行连接测试**
```bash
# 初始化Go模块（如果还没有）
go mod init movieinfo-test

# 添加依赖
go get github.com/go-sql-driver/mysql
go get github.com/redis/go-redis/v9

# 运行MySQL测试
go run test-mysql.go

# 运行Redis测试
go run test-redis.go
```

## Docker Compose集成

### 1. **完整的docker-compose.yml**
```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: movieinfo-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: movieinfo123
      MYSQL_DATABASE: movieinfo
      MYSQL_USER: movieinfo
      MYSQL_PASSWORD: movieinfo123
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./docker/mysql/conf/my.cnf:/etc/mysql/conf.d/my.cnf
      - ./docker/mysql/init:/docker-entrypoint-initdb.d
    command: --default-authentication-plugin=mysql_native_password
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      timeout: 20s
      retries: 10

  redis:
    image: redis:7-alpine
    container_name: movieinfo-redis
    restart: unless-stopped
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
      - ./docker/redis/conf/redis.conf:/usr/local/etc/redis/redis.conf
    command: redis-server /usr/local/etc/redis/redis.conf
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      timeout: 10s
      retries: 5

volumes:
  mysql_data:
  redis_data:

networks:
  default:
    name: movieinfo-network
```

### 2. **启动和管理**
```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs mysql
docker-compose logs redis

# 停止服务
docker-compose down

# 重启服务
docker-compose restart mysql redis
```

## 总结

数据库环境搭建为MovieInfo项目提供了稳定可靠的数据存储基础。通过MySQL和Redis的组合，我们实现了高性能的数据存储和缓存解决方案。

**关键配置要点**：
1. **MySQL优化**：内存配置、连接池、慢查询监控
2. **Redis优化**：内存策略、持久化配置、安全设置
3. **监控维护**：性能监控、备份策略、维护脚本
4. **容器化部署**：Docker Compose统一管理

**环境优势**：
- **高性能**：优化的配置保证数据库性能
- **高可用**：完善的备份和恢复机制
- **易维护**：自动化的监控和维护脚本
- **可扩展**：为读写分离和集群部署做好准备

**下一步**：基于这个数据库环境，我们将配置Redis缓存系统，包括缓存策略、键设计和性能优化。
