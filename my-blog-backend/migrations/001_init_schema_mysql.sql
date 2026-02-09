-- ============================================
-- 博客系统数据库初始化脚本 (MySQL版本)
-- 支持数据库: MySQL 5.7+
-- ============================================

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================
-- 1. 用户表
-- ============================================
CREATE TABLE IF NOT EXISTS `user` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(50) NOT NULL COMMENT '用户名',
    `email` VARCHAR(100) NOT NULL COMMENT '邮箱',
    `password` VARCHAR(255) NOT NULL COMMENT '密码',
    `nickname` VARCHAR(50) DEFAULT NULL COMMENT '昵称',
    `avatar` VARCHAR(500) DEFAULT NULL COMMENT '头像',
    `bio` TEXT COMMENT '个人简介',
    `website` VARCHAR(200) DEFAULT NULL COMMENT '个人网站',
    `github` VARCHAR(200) DEFAULT NULL COMMENT 'GitHub',
    `role` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '角色：0-普通用户，1-管理员',
    `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    UNIQUE KEY `uk_email` (`email`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================
-- 2. 角色表
-- ============================================
CREATE TABLE IF NOT EXISTS `roles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(50) NOT NULL COMMENT '角色名称',
    `code` VARCHAR(50) NOT NULL COMMENT '角色代码',
    `status` INT DEFAULT 1 COMMENT '状态',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name` (`name`),
    UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- ============================================
-- 3. 分类表
-- ============================================
CREATE TABLE IF NOT EXISTS `categorie` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
    `slug` VARCHAR(100) NOT NULL COMMENT '分类别名',
    `description` VARCHAR(200) DEFAULT NULL COMMENT '分类描述',
    `icon` VARCHAR(50) DEFAULT NULL COMMENT '分类图标',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类表';

-- ============================================
-- 4. 标签表
-- ============================================
CREATE TABLE IF NOT EXISTS `tag` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(50) NOT NULL COMMENT '标签名称',
    `slug` VARCHAR(100) NOT NULL COMMENT '标签别名',
    `description` VARCHAR(200) DEFAULT NULL COMMENT '标签描述',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name` (`name`),
    UNIQUE KEY `uk_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='标签表';

-- ============================================
-- 5. 文章表
-- ============================================
CREATE TABLE IF NOT EXISTS `articles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `title` VARCHAR(200) NOT NULL COMMENT '文章标题',
    `slug` VARCHAR(200) NOT NULL COMMENT '文章别名',
    `summary` TEXT COMMENT '文章摘要',
    `content` LONGTEXT NOT NULL COMMENT '文章内容',
    `cover` VARCHAR(500) DEFAULT NULL COMMENT '封面图',
    `category_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '分类ID',
    `author_id` BIGINT UNSIGNED NOT NULL COMMENT '作者ID',
    `views` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '浏览次数',
    `likes` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数',
    `favorites` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '收藏数',
    `comment_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '评论数',
    `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：0-草稿，1-已发布',
    `is_top` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否置顶',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `published_at` DATETIME DEFAULT NULL COMMENT '发布时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_slug` (`slug`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_author_id` (`author_id`),
    KEY `idx_status` (`status`),
    KEY `idx_is_top` (`is_top`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章表';

-- ============================================
-- 6. 评论表
-- ============================================
CREATE TABLE IF NOT EXISTS `comment` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `parent_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '父评论ID',
    `content` TEXT NOT NULL COMMENT '评论内容',
    `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：0-待审核，1-已通过，2-已拒绝',
    `ip_address` VARCHAR(50) DEFAULT NULL COMMENT 'IP地址',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_article_id` (`article_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_comment_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_comment_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_comment_parent` FOREIGN KEY (`parent_id`) REFERENCES `comment` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- ============================================
-- 7. 文章标签关联表
-- ============================================
CREATE TABLE IF NOT EXISTS `article_tag` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章ID',
    `tag_id` BIGINT UNSIGNED NOT NULL COMMENT '标签ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_article_tag` (`article_id`, `tag_id`),
    KEY `idx_article_id` (`article_id`),
    KEY `idx_tag_id` (`tag_id`),
    CONSTRAINT `fk_article_tag_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_article_tag_tag` FOREIGN KEY (`tag_id`) REFERENCES `tag` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章标签关联表';

-- ============================================
-- 8. 文章点赞表
-- ============================================
CREATE TABLE IF NOT EXISTS `article_likes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_article_user` (`article_id`, `user_id`),
    KEY `idx_article_id` (`article_id`),
    KEY `idx_user_id` (`user_id`),
    CONSTRAINT `fk_article_likes_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_article_likes_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章点赞表';

-- ============================================
-- 9. 文章收藏表
-- ============================================
CREATE TABLE IF NOT EXISTS `article_favorites` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_article_user` (`article_id`, `user_id`),
    KEY `idx_article_id` (`article_id`),
    KEY `idx_user_id` (`user_id`),
    CONSTRAINT `fk_article_favorites_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_article_favorites_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章收藏表';

-- ============================================
-- 10. 文章浏览记录表
-- ============================================
CREATE TABLE IF NOT EXISTS `article_views` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章ID',
    `ip_address` VARCHAR(50) DEFAULT NULL COMMENT 'IP地址',
    `user_agent` VARCHAR(500) DEFAULT NULL COMMENT 'User Agent',
    `user_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_article_id` (`article_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `fk_article_views_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_article_views_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章浏览记录表';

-- ============================================
-- 11. 专题系列表
-- ============================================
CREATE TABLE IF NOT EXISTS `series` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(100) NOT NULL COMMENT '专题名称',
    `slug` VARCHAR(100) NOT NULL COMMENT '专题别名',
    `icon` VARCHAR(50) DEFAULT NULL COMMENT '图标',
    `description` TEXT COMMENT '专题描述',
    `cover` VARCHAR(500) DEFAULT NULL COMMENT '封面',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='专题系列表';

-- ============================================
-- 12. 专题章节表
-- ============================================
CREATE TABLE IF NOT EXISTS `series_sections` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `series_id` BIGINT UNSIGNED NOT NULL COMMENT '专题ID',
    `name` VARCHAR(100) NOT NULL COMMENT '章节名称',
    `description` TEXT COMMENT '章节描述',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_series_id` (`series_id`),
    CONSTRAINT `fk_series_sections_series` FOREIGN KEY (`series_id`) REFERENCES `series` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='专题章节表';

-- ============================================
-- 13. 专题子章节表
-- ============================================
CREATE TABLE IF NOT EXISTS `series_subchapters` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `section_id` BIGINT UNSIGNED NOT NULL COMMENT '章节ID',
    `name` VARCHAR(100) NOT NULL COMMENT '子章节名称',
    `description` TEXT COMMENT '子章节描述',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_section_id` (`section_id`),
    CONSTRAINT `fk_series_subchapters_section` FOREIGN KEY (`section_id`) REFERENCES `series_sections` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='专题子章节表';

-- ============================================
-- 14. 子章节文章关联表
-- ============================================
CREATE TABLE IF NOT EXISTS `subchapter_articles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `subchapter_id` BIGINT UNSIGNED NOT NULL COMMENT '子章节ID',
    `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章ID',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_subchapter_article` (`subchapter_id`, `article_id`),
    KEY `idx_subchapter_id` (`subchapter_id`),
    KEY `idx_article_id` (`article_id`),
    CONSTRAINT `fk_subchapter_articles_subchapter` FOREIGN KEY (`subchapter_id`) REFERENCES `series_subchapters` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_subchapter_articles_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='子章节文章关联表';

SET FOREIGN_KEY_CHECKS = 1;

-- ============================================
-- 插入初始数据
-- ============================================

-- 插入默认管理员用户 (密码: admin123)
-- 密码哈希: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
INSERT IGNORE INTO `user` (`username`, `email`, `password`, `nickname`, `role`, `status`)
VALUES ('admin', 'admin@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Administrator', 1, 1);

-- 插入默认角色
INSERT IGNORE INTO `roles` (`name`, `code`, `status`) VALUES
('普通用户', 'user', 1),
('管理员', 'admin', 1);

-- 插入默认分类
INSERT IGNORE INTO `categorie` (`name`, `slug`, `description`, `sort_order`, `status`) VALUES
('技术', 'tech', '技术相关文章', 1, 1),
('生活', 'life', '生活感悟', 2, 1),
('读书', 'reading', '读书笔记', 3, 1);

-- 插入默认标签
INSERT IGNORE INTO `tag` (`name`, `slug`, `description`) VALUES
('Go', 'go', 'Go语言相关'),
('Vue', 'vue', 'Vue框架相关'),
('数据库', 'database', '数据库相关'),
('后端', 'backend', '后端开发');
