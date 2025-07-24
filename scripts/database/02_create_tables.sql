-- 创建所有数据表

USE movieinfo;

-- 创建用户表
CREATE TABLE users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    email VARCHAR(100) NOT NULL COMMENT '邮箱地址',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希值',
    nickname VARCHAR(50) DEFAULT NULL COMMENT '昵称',
    avatar_url VARCHAR(255) DEFAULT NULL COMMENT '头像URL',
    status TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '用户状态：0-禁用，1-正常',
    email_verified BOOLEAN NOT NULL DEFAULT FALSE COMMENT '邮箱是否已验证',
    last_login_at TIMESTAMP NULL DEFAULT NULL COMMENT '最后登录时间',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_username (username),
    UNIQUE KEY uk_email (email),
    KEY idx_status (status),
    KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 创建电影分类表
CREATE TABLE categories (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '分类ID',
    name VARCHAR(50) NOT NULL COMMENT '分类名称',
    description TEXT DEFAULT NULL COMMENT '分类描述',
    sort_order INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序顺序',
    status TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_name (name),
    KEY idx_sort_order (sort_order),
    KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='电影分类表';

-- 创建电影表
CREATE TABLE movies (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '电影ID',
    title VARCHAR(200) NOT NULL COMMENT '电影标题',
    original_title VARCHAR(200) DEFAULT NULL COMMENT '原始标题',
    description TEXT DEFAULT NULL COMMENT '电影描述',
    director VARCHAR(100) DEFAULT NULL COMMENT '导演',
    actors TEXT DEFAULT NULL COMMENT '主要演员，JSON格式存储',
    release_date DATE DEFAULT NULL COMMENT '上映日期',
    duration INT UNSIGNED DEFAULT NULL COMMENT '时长（分钟）',
    country VARCHAR(100) DEFAULT NULL COMMENT '制片国家/地区',
    language VARCHAR(50) DEFAULT NULL COMMENT '语言',
    poster_url VARCHAR(255) DEFAULT NULL COMMENT '海报URL',
    trailer_url VARCHAR(255) DEFAULT NULL COMMENT '预告片URL',
    imdb_id VARCHAR(20) DEFAULT NULL COMMENT 'IMDB ID',
    category_id INT UNSIGNED DEFAULT NULL COMMENT '主要分类ID',
    rating_average DECIMAL(3,1) UNSIGNED DEFAULT 0.0 COMMENT '平均评分',
    rating_count INT UNSIGNED DEFAULT 0 COMMENT '评分人数',
    view_count INT UNSIGNED DEFAULT 0 COMMENT '查看次数',
    status TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：0-下架，1-上架',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_title (title),
    KEY idx_director (director),
    KEY idx_release_date (release_date),
    KEY idx_category_id (category_id),
    KEY idx_rating_average (rating_average),
    KEY idx_status (status),
    KEY idx_created_at (created_at),
    UNIQUE KEY uk_imdb_id (imdb_id),
    CONSTRAINT fk_movies_category FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='电影表';

-- 创建电影分类关联表
CREATE TABLE movie_categories (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关联ID',
    movie_id BIGINT UNSIGNED NOT NULL COMMENT '电影ID',
    category_id INT UNSIGNED NOT NULL COMMENT '分类ID',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_movie_category (movie_id, category_id),
    KEY idx_movie_id (movie_id),
    KEY idx_category_id (category_id),
    CONSTRAINT fk_movie_categories_movie FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE CASCADE,
    CONSTRAINT fk_movie_categories_category FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='电影分类关联表';

-- 创建用户评分表
CREATE TABLE user_ratings (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评分ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    movie_id BIGINT UNSIGNED NOT NULL COMMENT '电影ID',
    rating TINYINT UNSIGNED NOT NULL COMMENT '评分：1-10分',
    comment TEXT DEFAULT NULL COMMENT '评价内容',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_user_movie (user_id, movie_id),
    KEY idx_movie_id (movie_id),
    KEY idx_rating (rating),
    KEY idx_created_at (created_at),
    CONSTRAINT fk_user_ratings_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_ratings_movie FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE CASCADE,
    CONSTRAINT chk_rating_range CHECK (rating >= 1 AND rating <= 10)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户评分表';
