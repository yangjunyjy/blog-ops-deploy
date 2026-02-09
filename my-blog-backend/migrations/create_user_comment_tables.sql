-- 用户表结构
CREATE TABLE IF NOT EXISTS `user` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户唯一标识ID',
  `username` VARCHAR(50) NOT NULL COMMENT '用户登录用户名',
  `email` VARCHAR(100) NOT NULL COMMENT '用户邮箱地址',
  `password` VARCHAR(255) NOT NULL COMMENT '用户密码哈希值',
  `nickname` VARCHAR(50) DEFAULT NULL COMMENT '用户昵称',
  `avatar` VARCHAR(500) DEFAULT NULL COMMENT '用户头像URL地址',
  `bio` TEXT COMMENT '用户个人简介',
  `role` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户角色权限：0-普通用户，1-管理员',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '账户状态：0-禁用，1-启用',
  `created_at` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '账户创建时间',
  `updated_at` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '账户最后更新时间',
  `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_email` (`email`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 评论表结构
CREATE TABLE IF NOT EXISTS `comment` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评论唯一标识ID',
  `article_id` BIGINT UNSIGNED NOT NULL COMMENT '被评论文章的ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '发表评论用户的ID',
  `parent_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '父评论ID，用于回复评论',
  `content` TEXT NOT NULL COMMENT '评论内容',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '评论审核状态：0-待审核，1-已通过，2-已拒绝',
  `ip_address` VARCHAR(50) DEFAULT NULL COMMENT '发表评论的IP地址',
  `created_at` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '评论发表时间',
  `updated_at` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '评论最后更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_article_id` (`article_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- 用户评论中间表（用于查询用户参与评论的所有文章）
CREATE TABLE IF NOT EXISTS `user_comment_article` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '记录唯一标识ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章ID',
  `comment_id` BIGINT UNSIGNED NOT NULL COMMENT '评论ID',
  `is_author` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否为文章作者：0-否，1-是',
  `created_at` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_article_comment` (`user_id`, `article_id`, `comment_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_article_id` (`article_id`),
  KEY `idx_comment_id` (`comment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户评论文章关联表';

-- 触发器：在评论表插入记录时，同步插入到用户评论中间表
DELIMITER $$

CREATE TRIGGER `trigger_comment_insert`
AFTER INSERT ON `comment`
FOR EACH ROW
BEGIN
  -- 获取文章作者ID
  DECLARE author_id BIGINT UNSIGNED;

  SELECT author_id INTO author_id FROM articles WHERE id = NEW.article_id;

  -- 插入评论者记录
  INSERT IGNORE INTO `user_comment_article` (`user_id`, `article_id`, `comment_id`, `is_author`)
  VALUES (NEW.user_id, NEW.article_id, NEW.id, 0);

  -- 如果评论者不是作者，插入作者记录
  IF NEW.user_id != author_id THEN
    INSERT IGNORE INTO `user_comment_article` (`user_id`, `article_id`, `comment_id`, `is_author`)
    VALUES (author_id, NEW.article_id, NEW.id, 1);
  END IF;
END$$

DELIMITER ;

-- 触发器：在评论表删除记录时，同步删除用户评论中间表记录
DELIMITER $$

CREATE TRIGGER `trigger_comment_delete`
AFTER DELETE ON `comment`
FOR EACH ROW
BEGIN
  DELETE FROM `user_comment_article` WHERE comment_id = OLD.id;
END$$

DELIMITER ;
