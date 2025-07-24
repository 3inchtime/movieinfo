-- scripts/database/01_create_database.sql
-- 创建数据库和用户

-- 创建数据库
CREATE DATABASE IF NOT EXISTS movieinfo
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

-- 创建用户（根据实际需要修改用户名和密码）
CREATE USER IF NOT EXISTS 'movieinfo_user'@'localhost' IDENTIFIED BY 'movieinfo_password';

-- 授权
GRANT ALL PRIVILEGES ON movieinfo.* TO 'movieinfo_user'@'localhost';

-- 刷新权限
FLUSH PRIVILEGES;

-- 使用数据库
USE movieinfo;