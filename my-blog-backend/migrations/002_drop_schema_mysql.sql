-- ============================================
-- 博客系统数据库回滚脚本 (MySQL版本)
-- ============================================

SET FOREIGN_KEY_CHECKS = 0;

-- 删除顺序很重要，因为有外键约束

-- 1. 删除子章节文章关联表
DROP TABLE IF EXISTS `subchapter_articles`;

-- 2. 删除专题子章节表
DROP TABLE IF EXISTS `series_subchapters`;

-- 3. 删除专题章节表
DROP TABLE IF EXISTS `series_sections`;

-- 4. 删除专题系列表
DROP TABLE IF EXISTS `series`;

-- 5. 删除文章浏览记录表
DROP TABLE IF EXISTS `article_views`;

-- 6. 删除文章收藏表
DROP TABLE IF EXISTS `article_favorites`;

-- 7. 删除文章点赞表
DROP TABLE IF EXISTS `article_likes`;

-- 8. 删除文章标签关联表
DROP TABLE IF EXISTS `article_tag`;

-- 9. 删除评论表
DROP TABLE IF EXISTS `comment`;

-- 10. 删除文章表
DROP TABLE IF EXISTS `articles`;

-- 11. 删除标签表
DROP TABLE IF EXISTS `tag`;

-- 12. 删除分类表
DROP TABLE IF EXISTS `categorie`;

-- 13. 删除角色表
DROP TABLE IF EXISTS `roles`;

-- 14. 删除用户表
DROP TABLE IF EXISTS `user`;

SET FOREIGN_KEY_CHECKS = 1;
