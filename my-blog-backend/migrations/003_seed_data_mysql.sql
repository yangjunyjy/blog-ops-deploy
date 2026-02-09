-- ============================================
-- 博客系统测试数据脚本 (MySQL版本)
-- 仅用于开发和测试环境
-- ============================================

-- ============================================
-- 插入测试用户
-- ============================================

INSERT IGNORE INTO `user` (`username`, `email`, `password`, `nickname`, `bio`, `role`, `status`)
VALUES
('user1', 'user1@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '测试用户1', '这是一个测试用户', 0, 1),
('user2', 'user2@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '测试用户2', '另一个测试用户', 0, 1);

-- ============================================
-- 插入更多分类
-- ============================================

INSERT IGNORE INTO `categorie` (`name`, `slug`, `description`, `sort_order`, `status`) VALUES
('前端', 'frontend', '前端开发相关', 4, 1),
('后端', 'backend', '后端开发相关', 5, 1),
('DevOps', 'devops', '运维部署相关', 6, 1),
('算法', 'algorithm', '算法和数据结构', 7, 1);

-- ============================================
-- 插入更多标签
-- ============================================

INSERT IGNORE INTO `tag` (`name`, `slug`, `description`) VALUES
('JavaScript', 'javascript', 'JavaScript编程语言'),
('TypeScript', 'typescript', 'TypeScript编程语言'),
('React', 'react', 'React框架'),
('Node.js', 'nodejs', 'Node.js运行环境'),
('Docker', 'docker', 'Docker容器技术'),
('Kubernetes', 'kubernetes', 'Kubernetes容器编排'),
('Git', 'git', 'Git版本控制'),
('Linux', 'linux', 'Linux操作系统'),
('Python', 'python', 'Python编程语言'),
('Java', 'java', 'Java编程语言'),
('Spring', 'spring', 'Spring框架'),
('MyBatis', 'mybatis', 'MyBatis持久层框架'),
('Redis', 'redis', 'Redis缓存数据库'),
('MongoDB', 'mongodb', 'MongoDB文档数据库'),
('Nginx', 'nginx', 'Nginx服务器');

-- ============================================
-- 插入专题系列
-- ============================================

INSERT IGNORE INTO `series` (`name`, `slug`, `description`, `sort_order`, `status`) VALUES
('Go语言入门', 'go-intro', 'Go语言从入门到精通', 1, 1),
('Vue3实战', 'vue3-practice', 'Vue3项目实战开发', 2, 1),
('微服务架构', 'microservices', '微服务架构设计与实践', 3, 1);

-- ============================================
-- 插入专题章节
-- ============================================

-- Go语言入门
INSERT IGNORE INTO `series_sections` (`series_id`, `name`, `description`, `sort_order`) VALUES
((SELECT id FROM `series` WHERE slug = 'go-intro'), '基础语法', 'Go语言基础语法', 1),
((SELECT id FROM `series` WHERE slug = 'go-intro'), '并发编程', 'Go语言并发编程', 2),
((SELECT id FROM `series` WHERE slug = 'go-intro'), '网络编程', 'Go语言网络编程', 3);

-- Vue3实战
INSERT IGNORE INTO `series_sections` (`series_id`, `name`, `description`, `sort_order`) VALUES
((SELECT id FROM `series` WHERE slug = 'vue3-practice'), '基础组件', 'Vue3基础组件开发', 1),
((SELECT id FROM `series` WHERE slug = 'vue3-practice'), '状态管理', 'Vuex和Pinia状态管理', 2),
((SELECT id FROM `series` WHERE slug = 'vue3-practice'), '路由管理', 'Vue Router路由管理', 3);

-- 微服务架构
INSERT IGNORE INTO `series_sections` (`series_id`, `name`, `description`, `sort_order`) VALUES
((SELECT id FROM `series` WHERE slug = 'microservices'), '架构设计', '微服务架构设计原则', 1),
((SELECT id FROM `series` WHERE slug = 'microservices'), '服务治理', '微服务服务治理', 2),
((SELECT id FROM `series` WHERE slug = 'microservices'), '容器化部署', 'Docker和Kubernetes部署', 3);

-- ============================================
-- 插入专题子章节
-- ============================================

-- Go语言入门 - 基础语法
INSERT IGNORE INTO `series_subchapters` (`section_id`, `name`, `description`, `sort_order`)
VALUES
((SELECT id FROM `series_sections` WHERE name = '基础语法' AND series_id = (SELECT id FROM `series` WHERE slug = 'go-intro')), '变量与常量', 'Go语言变量和常量的定义', 1),
((SELECT id FROM `series_sections` WHERE name = '基础语法' AND series_id = (SELECT id FROM `series` WHERE slug = 'go-intro')), '数据类型', 'Go语言的基本数据类型', 2),
((SELECT id FROM `series_sections` WHERE name = '基础语法' AND series_id = (SELECT id FROM `series` WHERE slug = 'go-intro')), '控制流', 'Go语言的条件和循环', 3);

-- Go语言入门 - 并发编程
INSERT IGNORE INTO `series_subchapters` (`section_id`, `name`, `description`, `sort_order`)
VALUES
((SELECT id FROM `series_sections` WHERE name = '并发编程' AND series_id = (SELECT id FROM `series` WHERE slug = 'go-intro')), 'Goroutine', 'Go语言协程', 1),
((SELECT id FROM `series_sections` WHERE name = '并发编程' AND series_id = (SELECT id FROM `series` WHERE slug = 'go-intro')), 'Channel', 'Go语言通道', 2),
((SELECT id FROM `series_sections` WHERE name = '并发编程' AND series_id = (SELECT id FROM `series` WHERE slug = 'go-intro')), 'Sync包', 'Go语言同步原语', 3);

-- ============================================
-- 插入测试文章
-- ============================================

INSERT IGNORE INTO `articles` (`title`, `slug`, `summary`, `content`, `category_id`, `author_id`, `status`, `is_top`, `published_at`)
VALUES
('Go语言入门指南', 'go-introduction', '本文介绍Go语言的基础知识和开发环境搭建', 'Go语言是Google开发的一种静态强类型、编译型语言...', (SELECT id FROM `categorie` WHERE slug = 'tech'), (SELECT id FROM `user` WHERE username = 'admin'), 1, 1, NOW()),
('Vue3 Composition API详解', 'vue3-composition-api', '详解Vue3的Composition API使用方法', 'Composition API是Vue3的重要特性...', (SELECT id FROM `categorie` WHERE slug = 'frontend'), (SELECT id FROM `user` WHERE username = 'admin'), 1, 0, NOW()),
('微服务架构设计原则', 'microservices-principles', '介绍微服务架构的核心设计原则', '微服务架构是一种将单一应用程序开发为一组小型服务的方法...', (SELECT id FROM `categorie` WHERE slug = 'backend'), (SELECT id FROM `user` WHERE username = 'admin'), 1, 1, NOW()),
('Docker容器化部署实践', 'docker-deployment', '分享Docker容器化部署的实际经验', 'Docker是一种开源的容器化平台...', (SELECT id FROM `categorie` WHERE slug = 'devops'), (SELECT id FROM `user` WHERE username = 'admin'), 1, 0, NOW()),
('Redis缓存最佳实践', 'redis-best-practices', '总结Redis缓存使用的最佳实践', 'Redis是一个高性能的键值对数据库...', (SELECT id FROM `categorie` WHERE slug = 'tech'), (SELECT id FROM `user` WHERE username = 'admin'), 1, 0, NOW());

-- ============================================
-- 插入文章标签关联
-- ============================================

-- Go语言入门指南
INSERT IGNORE INTO `article_tag` (`article_id`, `tag_id`)
VALUES
((SELECT id FROM `articles` WHERE slug = 'go-introduction'), (SELECT id FROM `tag` WHERE slug = 'go')),
((SELECT id FROM `articles` WHERE slug = 'go-introduction'), (SELECT id FROM `tag` WHERE slug = 'backend'));

-- Vue3 Composition API详解
INSERT IGNORE INTO `article_tag` (`article_id`, `tag_id`)
VALUES
((SELECT id FROM `articles` WHERE slug = 'vue3-composition-api'), (SELECT id FROM `tag` WHERE slug = 'vue')),
((SELECT id FROM `articles` WHERE slug = 'vue3-composition-api'), (SELECT id FROM `tag` WHERE slug = 'frontend'));

-- 微服务架构设计原则
INSERT IGNORE INTO `article_tag` (`article_id`, `tag_id`)
VALUES
((SELECT id FROM `articles` WHERE slug = 'microservices-principles'), (SELECT id FROM `tag` WHERE slug = 'backend')),
((SELECT id FROM `articles` WHERE slug = 'microservices-principles'), (SELECT id FROM `tag` WHERE slug = 'devops'));

-- Docker容器化部署实践
INSERT IGNORE INTO `article_tag` (`article_id`, `tag_id`)
VALUES
((SELECT id FROM `articles` WHERE slug = 'docker-deployment'), (SELECT id FROM `tag` WHERE slug = 'docker')),
((SELECT id FROM `articles` WHERE slug = 'docker-deployment'), (SELECT id FROM `tag` WHERE slug = 'devops'));

-- Redis缓存最佳实践
INSERT IGNORE INTO `article_tag` (`article_id`, `tag_id`)
VALUES
((SELECT id FROM `articles` WHERE slug = 'redis-best-practices'), (SELECT id FROM `tag` WHERE slug = 'redis')),
((SELECT id FROM `articles` WHERE slug = 'redis-best-practices'), (SELECT id FROM `tag` WHERE slug = 'database'));

-- ============================================
-- 插入测试评论
-- ============================================

INSERT IGNORE INTO `comment` (`article_id`, `user_id`, `content`, `status`)
VALUES
((SELECT id FROM `articles` WHERE slug = 'go-introduction'), (SELECT id FROM `user` WHERE username = 'user1'), '非常棒的入门教程，感谢分享！', 1),
((SELECT id FROM `articles` WHERE slug = 'go-introduction'), (SELECT id FROM `user` WHERE username = 'user2'), '学到了很多，继续加油', 1),
((SELECT id FROM `articles` WHERE slug = 'vue3-composition-api'), (SELECT id FROM `user` WHERE username = 'user1'), 'Composition API确实比Options API更好用', 1);

-- ============================================
-- 插入文章浏览记录（模拟）
-- ============================================

INSERT IGNORE INTO `article_views` (`article_id`, `ip_address`, `user_agent`)
SELECT id, '127.0.0.1', 'Mozilla/5.0' FROM `articles` LIMIT 3;

-- ============================================
-- 插入文章点赞记录
-- ============================================

INSERT IGNORE INTO `article_likes` (`article_id`, `user_id`)
VALUES
((SELECT id FROM `articles` WHERE slug = 'go-introduction'), (SELECT id FROM `user` WHERE username = 'user1')),
((SELECT id FROM `articles` WHERE slug = 'vue3-composition-api'), (SELECT id FROM `user` WHERE username = 'user2'));

-- ============================================
-- 插入文章收藏记录
-- ============================================

INSERT IGNORE INTO `article_favorites` (`article_id`, `user_id`)
VALUES
((SELECT id FROM `articles` WHERE slug = 'go-introduction'), (SELECT id FROM `user` WHERE username = 'user1')),
((SELECT id FROM `articles` WHERE slug = 'redis-best-practices'), (SELECT id FROM `user` WHERE username = 'user2'));
