-- 添加评论点赞数统计字段
ALTER TABLE `comment` ADD COLUMN `likes` INT UNSIGNED NOT NULL DEFAULT 0 AFTER `status`;

-- 创建评论点赞表
CREATE TABLE IF NOT EXISTS `comment_like` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `comment_id` INT UNSIGNED NOT NULL COMMENT '评论ID',
  `user_id` INT UNSIGNED NOT NULL COMMENT '用户ID',
  `created_at` DATETIME NOT NULL COMMENT '点赞时间',
  INDEX `idx_comment_id` (`comment_id`),
  INDEX `idx_user_id` (`user_id`),
  UNIQUE KEY `idx_comment_user` (`comment_id`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论点赞表';

-- 确保现有评论的点赞数默认为0
UPDATE `comment` SET `likes` = 0 WHERE `likes` IS NULL;
