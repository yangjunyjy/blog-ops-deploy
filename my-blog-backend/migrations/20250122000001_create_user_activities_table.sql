-- 创建用户活动记录表
CREATE TABLE IF NOT EXISTS `user_activities` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `article_id` bigint unsigned NOT NULL COMMENT '文章ID',
  `type` varchar(20) NOT NULL COMMENT '活动类型：like(点赞)、comment(评论)、share(分享)、favorite(收藏)',
  `content` text COMMENT '评论内容（仅评论类型）',
  `platform` varchar(50) DEFAULT NULL COMMENT '分享平台（仅分享类型）',
  `folder_name` varchar(100) DEFAULT NULL COMMENT '收藏文件夹名称（仅收藏类型）',
  `article_title` varchar(200) NOT NULL COMMENT '文章标题（冗余字段，便于查询）',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_article_id` (`article_id`),
  KEY `idx_type` (`type`),
  KEY `idx_user_type` (`user_id`, `type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户活动记录表';
