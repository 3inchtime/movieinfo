# scripts/database/reset_database.sh
#!/bin/bash

# 数据库重置脚本

set -e

echo "警告: 此操作将删除所有数据！"
read -p "确定要重置数据库吗？(y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "操作已取消"
    exit 1
fi

# 数据库连接参数
DB_HOST="localhost"
DB_PORT="3306"
DB_ROOT_USER="root"
DB_ROOT_PASSWORD="your_root_password"

echo "删除现有数据库..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_ROOT_USER" -p"$DB_ROOT_PASSWORD" -e "DROP DATABASE IF EXISTS movieinfo;"

echo "重新初始化数据库..."
./init_database.sh

echo "数据库重置完成！"