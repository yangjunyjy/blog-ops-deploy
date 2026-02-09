-- ============================================
-- 博客系统数据库初始化脚本
-- 支持数据库: MySQL 5.7+, PostgreSQL 9.6+, SQLite 3.x
-- ============================================

-- ============================================
-- 1. 用户表
-- ============================================
CREATE TABLE IF NOT EXISTS `user` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `username` VARCHAR(50) NOT NULL UNIQUE,
    `email` VARCHAR(100) NOT NULL UNIQUE,
    `password` VARCHAR(255) NOT NULL,
    `nickname` VARCHAR(50),
    `avatar` VARCHAR(500),
    `bio` TEXT,
    `website` VARCHAR(200),
    `github` VARCHAR(200),
    `role` INTEGER NOT NULL DEFAULT 0,
    `status` INTEGER NOT NULL DEFAULT 1,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` DATETIME
);

CREATE INDEX IF NOT EXISTS idx_user_deleted_at ON `user`(`deleted_at`);

-- ============================================
-- 2. 角色表
-- ============================================
CREATE TABLE IF NOT EXISTS `roles` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` VARCHAR(50) NOT NULL UNIQUE,
    `code` VARCHAR(50) NOT NULL UNIQUE,
    `status` INTEGER DEFAULT 1,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- 3. 分类表
-- ============================================
CREATE TABLE IF NOT EXISTS `categorie` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` VARCHAR(50) NOT NULL,
    `slug` VARCHAR(100) NOT NULL UNIQUE,
    `description` VARCHAR(200),
    `icon` VARCHAR(50),
    `sort_order` INTEGER NOT NULL DEFAULT 0,
    `status` INTEGER NOT NULL DEFAULT 1,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- 4. 标签表
-- ============================================
CREATE TABLE IF NOT EXISTS `tag` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` VARCHAR(50) NOT NULL UNIQUE,
    `slug` VARCHAR(100) NOT NULL UNIQUE,
    `description` VARCHAR(200),
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- 5. 文章表
-- ============================================
CREATE TABLE IF NOT EXISTS `articles` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `title` VARCHAR(200) NOT NULL,
    `slug` VARCHAR(200) NOT NULL UNIQUE,
    `summary` TEXT,
    `content` TEXT NOT NULL,
    `cover` VARCHAR(500),
    `category_id` INTEGER,
    `author_id` INTEGER NOT NULL,
    `views` INTEGER NOT NULL DEFAULT 0,
    `likes` INTEGER NOT NULL DEFAULT 0,
    `favorites` INTEGER NOT NULL DEFAULT 0,
    `comment_count` INTEGER NOT NULL DEFAULT 0,
    `status` INTEGER NOT NULL DEFAULT 1,
    `is_top` INTEGER NOT NULL DEFAULT 0,
    `sort_order` INTEGER NOT NULL DEFAULT 0,
    `published_at` DATETIME,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_articles_category_id ON `articles`(`category_id`);
CREATE INDEX IF NOT EXISTS idx_articles_author_id ON `articles`(`author_id`);

-- ============================================
-- 6. 评论表
-- ============================================
CREATE TABLE IF NOT EXISTS `comment` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `article_id` INTEGER NOT NULL,
    `user_id` INTEGER NOT NULL,
    `parent_id` INTEGER,
    `content` TEXT NOT NULL,
    `status` INTEGER NOT NULL DEFAULT 1,
    `ip_address` VARCHAR(50),
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`article_id`) REFERENCES `articles`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`parent_id`) REFERENCES `comment`(`id`) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_comment_article_id ON `comment`(`article_id`);
CREATE INDEX IF NOT EXISTS idx_comment_user_id ON `comment`(`user_id`);
CREATE INDEX IF NOT EXISTS idx_comment_parent_id ON `comment`(`parent_id`);

-- ============================================
-- 7. 文章标签关联表
-- ============================================
CREATE TABLE IF NOT EXISTS `article_tag` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `article_id` INTEGER NOT NULL,
    `tag_id` INTEGER NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`article_id`) REFERENCES `articles`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`tag_id`) REFERENCES `tag`(`id`) ON DELETE CASCADE,
    UNIQUE(`article_id`, `tag_id`)
);

CREATE INDEX IF NOT EXISTS idx_article_tag_article_id ON `article_tag`(`article_id`);
CREATE INDEX IF NOT EXISTS idx_article_tag_tag_id ON `article_tag`(`tag_id`);

-- ============================================
-- 8. 文章点赞表
-- ============================================
CREATE TABLE IF NOT EXISTS `article_likes` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `article_id` INTEGER NOT NULL,
    `user_id` INTEGER NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`article_id`) REFERENCES `articles`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    UNIQUE(`article_id`, `user_id`)
);

CREATE INDEX IF NOT EXISTS idx_article_likes_article_id ON `article_likes`(`article_id`);
CREATE INDEX IF NOT EXISTS idx_article_likes_user_id ON `article_likes`(`user_id`);

-- ============================================
-- 9. 文章收藏表
-- ============================================
CREATE TABLE IF NOT EXISTS `article_favorites` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `article_id` INTEGER NOT NULL,
    `user_id` INTEGER NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`article_id`) REFERENCES `articles`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    UNIQUE(`article_id`, `user_id`)
);

CREATE INDEX IF NOT EXISTS idx_article_favorites_article_id ON `article_favorites`(`article_id`);
CREATE INDEX IF NOT EXISTS idx_article_favorites_user_id ON `article_favorites`(`user_id`);

-- ============================================
-- 10. 文章浏览记录表
-- ============================================
CREATE TABLE IF NOT EXISTS `article_views` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `article_id` INTEGER NOT NULL,
    `ip_address` VARCHAR(50),
    `user_agent` VARCHAR(500),
    `user_id` INTEGER,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`article_id`) REFERENCES `articles`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_article_views_article_id ON `article_views`(`article_id`);
CREATE INDEX IF NOT EXISTS idx_article_views_user_id ON `article_views`(`user_id`);

-- ============================================
-- 11. 专题系列表
-- ============================================
CREATE TABLE IF NOT EXISTS `series` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` VARCHAR(100) NOT NULL,
    `slug` VARCHAR(100) NOT NULL UNIQUE,
    `icon` VARCHAR(50),
    `description` TEXT,
    `cover` VARCHAR(500),
    `sort_order` INTEGER NOT NULL DEFAULT 0,
    `status` INTEGER NOT NULL DEFAULT 1,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- 12. 专题章节表
-- ============================================
CREATE TABLE IF NOT EXISTS `series_sections` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `series_id` INTEGER NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `description` TEXT,
    `sort_order` INTEGER NOT NULL DEFAULT 0,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`series_id`) REFERENCES `series`(`id`) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_series_sections_series_id ON `series_sections`(`series_id`);

-- ============================================
-- 13. 专题子章节表
-- ============================================
CREATE TABLE IF NOT EXISTS `series_subchapters` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `section_id` INTEGER NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `description` TEXT,
    `sort_order` INTEGER NOT NULL DEFAULT 0,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`section_id`) REFERENCES `series_sections`(`id`) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_series_subchapters_section_id ON `series_subchapters`(`section_id`);

-- ============================================
-- 14. 子章节文章关联表
-- ============================================
CREATE TABLE IF NOT EXISTS `subchapter_articles` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `subchapter_id` INTEGER NOT NULL,
    `article_id` INTEGER NOT NULL,
    `sort_order` INTEGER NOT NULL DEFAULT 0,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`subchapter_id`) REFERENCES `series_subchapters`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`article_id`) REFERENCES `articles`(`id`) ON DELETE CASCADE,
    UNIQUE(`subchapter_id`, `article_id`)
);

CREATE INDEX IF NOT EXISTS idx_subchapter_articles_subchapter_id ON `subchapter_articles`(`subchapter_id`);
CREATE INDEX IF NOT EXISTS idx_subchapter_articles_article_id ON `subchapter_articles`(`article_id`);

-- ============================================
-- 插入初始数据
-- ============================================

-- 插入默认管理员用户 (密码: admin123, 需要在应用层加密)
-- 注意: 实际应用中应该使用bcrypt加密后的密码
INSERT INTO `user` (`username`, `email`, `password`, `nickname`, `role`, `status`)
VALUES ('admin', 'admin@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Administrator', 1, 1)
ON CONFLICT(username) DO NOTHING;

-- 插入默认角色
INSERT INTO `roles` (`name`, `code`, `status`) VALUES
('普通用户', 'user', 1),
('管理员', 'admin', 1)
ON CONFLICT(code) DO NOTHING;

-- 插入默认分类
INSERT INTO `categorie` (`name`, `slug`, `description`, `sort_order`, `status`) VALUES
('技术', 'tech', '技术相关文章', 1, 1),
('生活', 'life', '生活感悟', 2, 1),
('读书', 'reading', '读书笔记', 3, 1)
ON CONFLICT(slug) DO NOTHING;

-- 插入默认标签
INSERT INTO `tag` (`name`, `slug`, `description`) VALUES
('Go', 'go', 'Go语言相关'),
('Vue', 'vue', 'Vue框架相关'),
('数据库', 'database', '数据库相关'),
('后端', 'backend', '后端开发')
ON CONFLICT(slug) DO NOTHING;
