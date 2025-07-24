# scripts/database/init_database.sh
#!/bin/bash

# 数据库初始化脚本

set -e

echo "开始初始化数据库..."

# 数据库连接参数
DB_HOST="localhost"
DB_PORT="3306"
DB_ROOT_USER="root"
DB_ROOT_PASSWORD="your_root_password"

# 检查 MySQL 是否运行
echo "检查 MySQL 服务状态..."
if ! mysqladmin ping -h"$DB_HOST" -P"$DB_PORT" --silent; then
    echo "错误: MySQL 服务未运行或无法连接"
    exit 1
fi

echo "MySQL 服务正常运行"

# 执行数据库创建脚本
echo "创建数据库和用户..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_ROOT_USER" -p"$DB_ROOT_PASSWORD" < 01_create_database.sql

echo "创建数据表..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_ROOT_USER" -p"$DB_ROOT_PASSWORD" < 02_create_tables.sql

echo "插入初始数据..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_ROOT_USER" -p"$DB_ROOT_PASSWORD" < 03_insert_initial_data.sql

echo "数据库初始化完成！"