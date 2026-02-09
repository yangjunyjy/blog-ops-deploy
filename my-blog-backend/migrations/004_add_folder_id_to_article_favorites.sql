-- 添加 folder_id 字段到 article_favorites 表
-- Date: 2025-01-18

-- 添加 folder_id 字段（必需字段）
ALTER TABLE `article_favorites`
ADD COLUMN `folder_id` BIGINT UNSIGNED NOT NULL COMMENT '收藏文件夹ID' AFTER `user_id`;

-- 添加索引以加速按文件夹查询
CREATE INDEX IF NOT EXISTS `idx_article_favorites_folder_id` ON `article_favorites`(`folder_id`);

-- 添加外键约束到 favorite_folders 表（如果该表存在）
-- ALTER TABLE `article_favorites`
-- ADD CONSTRAINT `fk_article_favorites_folder`
-- FOREIGN KEY (`folder_id`) REFERENCES `favorite_folders` (`id`)
-- ON DELETE CASCADE;
