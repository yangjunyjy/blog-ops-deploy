# 博客系统数据库表结构文档

## 数据库概述

- **数据库类型**: SQLite / MySQL / PostgreSQL
- **字符集**: utf8mb4
- **排序规则**: utf8mb4_unicode_ci

## 表结构

### 1. 用户表 (users)

| 字段名 | 类型 | 长度 | 允许NULL | 默认值 | 说明 |
|--------|------|------|----------|--------|------|
| id | INTEGER | - | NO | - | 主键, 自增 |
| username | VARCHAR | 50 | NO | - | 用户名, 唯一 |
| email | VARCHAR | 100 | YES | - | 邮箱, 唯一 |
| password | VARCHAR | 255 | NO | - | 密码 (加密存储) |
| nickname | VARCHAR | 50 | YES | - | 昵称 |
| avatar | VARCHAR | 255 | YES | - | 头像URL |
| bio | TEXT | - | YES | - | 个人简介 |
| website | VARCHAR | 255 | YES | - | 个人网站 |
| github | VARCHAR | 255 | YES | - | GitHub账号 |
| role | INTEGER | - | NO | 0 | 角色 (0-普通用户, 1-管理员) |
| status | INTEGER | - | NO | 1 | 状态 (0-禁用, 1-启用) |
| created_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 更新时间 |
| deleted_at | DATETIME | - | YES | - | 软删除时间 |

**索引**:
- PRIMARY KEY (id)
- UNIQUE INDEX idx_username (username)
- UNIQUE INDEX idx_email (email)
- INDEX idx_role (role)
- INDEX idx_status (status)

---

### 2. 角色表 (roles)

| 字段名 | 类型 | 长度 | 允许NULL | 默认值 | 说明 |
|--------|------|------|----------|--------|------|
| id | INTEGER | - | NO | - | 主键, 自增 |
| name | VARCHAR | 50 | NO | - | 角色名称 |
| code | VARCHAR | 50 | NO | - | 角色代码, 唯一 |
| status | INTEGER | - | NO | 1 | 状态 (0-禁用, 1-启用) |
| created_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 更新时间 |

**索引**:
- PRIMARY KEY (id)
- UNIQUE INDEX idx_code (code)

---

### 3. 分类表 (categories)

| 字段名 | 类型 | 长度 | 允许NULL | 默认值 | 说明 |
|--------|------|------|----------|--------|------|
| id | INTEGER | - | NO | - | 主键, 自增 |
| name | VARCHAR | 50 | NO | - | 分类名称 |
| slug | VARCHAR | 100 | NO | - | URL别名, 唯一 |
| description | TEXT | - | YES | - | 分类描述 |
| icon | VARCHAR | 100 | YES | - | 分类图标 |
| sort_order | INTEGER | - | NO | 0 | 排序顺序 |
| status | INTEGER | - | NO | 1 | 状态 (0-禁用, 1-启用) |
| created_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 更新时间 |

**索引**:
- PRIMARY KEY (id)
- UNIQUE INDEX idx_slug (slug)
- INDEX idx_status (status)

---

### 4. 标签表 (tags)

| 字段名 | 类型 | 长度 | 允许NULL | 默认值 | 说明 |
|--------|------|------|----------|--------|------|
| id | INTEGER | - | NO | - | 主键, 自增 |
| name | VARCHAR | 50 | NO | - | 标签名称 |
| slug | VARCHAR | 100 | NO | - | URL别名, 唯一 |
| description | TEXT | - | YES | - | 标签描述 |
| created_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 创建时间 |

**索引**:
- PRIMARY KEY (id)
- UNIQUE INDEX idx_slug (slug)
- UNIQUE INDEX idx_name (name)

---

### 5. 文章表 (articles)

| 字段名 | 类型 | 长度 | 允许NULL | 默认值 | 说明 |
|--------|------|------|----------|--------|------|
| id | INTEGER | - | NO | - | 主键, 自增 |
| title | VARCHAR | 200 | NO | - | 文章标题 |
| slug | VARCHAR | 200 | NO | - | URL别名, 唯一 |
| summary | TEXT | - | YES | - | 文章摘要 |
| content | TEXT | - | NO | - | 文章内容 |
| cover | VARCHAR | 255 | YES | - | 封面图片URL |
| category_id | INTEGER | - | YES | - | 分类ID, 外键关联categories.id |
| author_id | INTEGER | - | NO | - | 作者ID, 外键关联users.id |
| views | INTEGER | - | NO | 0 | 浏览量 |
| likes | INTEGER | - | NO | 0 | 点赞数 |
| favorites | INTEGER | - | NO | 0 | 收藏数 |
| comment_count | INTEGER | - | NO | 0 | 评论数 |
| status | INTEGER | - | NO | 0 | 状态 (0-草稿, 1-已发布) |
| is_top | BOOLEAN | - | NO | false | 是否置顶 |
| sort_order | INTEGER | - | NO | 0 | 排序顺序 |
| published_at | DATETIME | - | YES | - | 发布时间 |
| created_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 更新时间 |
| deleted_at | DATETIME | - | YES | - | 软删除时间 |

**索引**:
- PRIMARY KEY (id)
- UNIQUE INDEX idx_slug (slug)
- INDEX idx_category_id (category_id)
- INDEX idx_author_id (author_id)
- INDEX idx_status (status)
- INDEX idx_published_at (published_at)
- INDEX idx_is_top (is_top)

**外键**:
- FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
- FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE

---

### 6. 文章标签关联表 (article_tags)

| 字段名 | 类型 | 长度 | 允许NULL | 默认值 | 说明 |
|--------|------|------|----------|--------|------|
| id | INTEGER | - | NO | - | 主键, 自增 |
| article_id | INTEGER | - | NO | - | 文章ID, 外键关联articles.id |
| tag_id | INTEGER | - | NO | - | 标签ID, 外键关联tags.id |
| created_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 创建时间 |

**索引**:
- PRIMARY KEY (id)
- UNIQUE INDEX idx_article_tag (article_id, tag_id)
- INDEX idx_article_id (article_id)
- INDEX idx_tag_id (tag_id)

**外键**:
- FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
- FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE

---

### 7. 文章点赞表 (article_likes)

| 字段名 | 类型 | 长度 | 允许NULL | 默认值 | 说明 |
|--------|------|------|----------|--------|------|
| id | INTEGER | - | NO | - | 主键, 自增 |
| article_id | INTEGER | - | NO | - | 文章ID, 外键关联articles.id |
| user_id | INTEGER | - | NO | - | 用户ID, 外键关联users.id |
| created_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 创建时间 |

**索引**:
- PRIMARY KEY (id)
- UNIQUE INDEX idx_article_user (article_id, user_id)
- INDEX idx_article_id (article_id)
- INDEX idx_user_id (user_id)

**外键**:
- FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

---

### 8. 文章收藏表 (article_favorites)

| 字段名 | 类型 | 长度 | 允许NULL | 默认值 | 说明 |
|--------|------|------|----------|--------|------|
| id | INTEGER | - | NO | - | 主键, 自增 |
| article_id | INTEGER | - | NO | - | 文章ID, 外键关联articles.id |
| user_id | INTEGER | - | NO | - | 用户ID, 外键关联users.id |
| created_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 创建时间 |

**索引**:
- PRIMARY KEY (id)
- UNIQUE INDEX idx_article_user (article_id, user_id)
- INDEX idx_article_id (article_id)
- INDEX idx_user_id (user_id)

**外键**:
- FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

---

### 9. 文章浏览记录表 (article_views)

| 字段名 | 类型 | 长度 | 允许NULL | 默认值 | 说明 |
|--------|------|------|----------|--------|------|
| id | INTEGER | - | NO | - | 主键, 自增 |
| article_id | INTEGER | - | NO | - | 文章ID, 外键关联articles.id |
| user_id | INTEGER | - | YES | - | 用户ID (可为null, 记录游客浏览), 外键关联users.id |
| ip_address | VARCHAR | 45 | YES | - | IP地址 (IPv4/IPv6) |
| user_agent | VARCHAR | 500 | YES | - | 用户代理字符串 |
| viewed_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 浏览时间 |

**索引**:
- PRIMARY KEY (id)
- INDEX idx_article_id (article_id)
- INDEX idx_user_id (user_id)
- INDEX idx_viewed_at (viewed_at)

**外键**:
- FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL

---

### 10. 评论表 (comments)

| 字段名 | 类型 | 长度 | 允许NULL | 默认值 | 说明 |
|--------|------|------|----------|--------|------|
| id | INTEGER | - | NO | - | 主键, 自增 |
| article_id | INTEGER | - | NO | - | 文章ID, 外键关联articles.id |
| user_id | INTEGER | - | NO | - | 用户ID, 外键关联users.id |
| parent_id | INTEGER | - | YES | - | 父评论ID, 支持多级评论, 外键关联comments.id |
| content | TEXT | - | NO | - | 评论内容 |
| status | INTEGER | - | NO | 0 | 状态 (0-待审核, 1-已通过, 2-已拒绝) |
| ip_address | VARCHAR | 45 | YES | - | IP地址 |
| user_agent | VARCHAR | 500 | YES | - | 用户代理字符串 |
| created_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 更新时间 |

**索引**:
- PRIMARY KEY (id)
- INDEX idx_article_id (article_id)
- INDEX idx_user_id (user_id)
- INDEX idx_parent_id (parent_id)
- INDEX idx_status (status)
- INDEX idx_created_at (created_at)

**外键**:
- FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
- FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE

---

### 11. 专题系列表 (series)

| 字段名 | 类型 | 长度 | 允许NULL | 默认值 | 说明 |
|--------|------|------|----------|--------|------|
| id | INTEGER | - | NO | - | 主键, 自增 |
| name | VARCHAR | 100 | NO | - | 系列名称 |
| slug | VARCHAR | 100 | NO | - | URL别名, 唯一 |
| icon | VARCHAR | 100 | YES | - | 系列图标 |
| description | TEXT | - | YES | - | 系列描述 |
| cover | VARCHAR | 255 | YES | - | 封面图片URL |
| sort_order | INTEGER | - | NO | 0 | 排序顺序 |
| status | INTEGER | - | NO | 1 | 状态 (0-禁用, 1-启用) |
| created_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 更新时间 |

**索引**:
- PRIMARY KEY (id)
- UNIQUE INDEX idx_slug (slug)
- INDEX idx_status (status)

---

### 12. 系列章节表 (series_sections)

| 字段名 | 类型 | 长度 | 允许NULL | 默认值 | 说明 |
|--------|------|------|----------|--------|------|
| id | INTEGER | - | NO | - | 主键, 自增 |
| series_id | INTEGER | - | NO | - | 系列ID, 外键关联series.id |
| name | VARCHAR | 100 | NO | - | 章节名称 |
| description | TEXT | - | YES | - | 章节描述 |
| sort_order | INTEGER | - | NO | 0 | 排序顺序 |
| created_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 更新时间 |

**索引**:
- PRIMARY KEY (id)
- INDEX idx_series_id (series_id)
- INDEX idx_sort_order (sort_order)

**外键**:
- FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE

---

### 13. 系列子章节表 (series_subchapters)

| 字段名 | 类型 | 长度 | 允许NULL | 默认值 | 说明 |
|--------|------|------|----------|--------|------|
| id | INTEGER | - | NO | - | 主键, 自增 |
| section_id | INTEGER | - | NO | - | 章节ID, 外键关联series_sections.id |
| name | VARCHAR | 100 | NO | - | 子章节名称 |
| description | TEXT | - | YES | - | 子章节描述 |
| sort_order | INTEGER | - | NO | 0 | 排序顺序 |
| created_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 更新时间 |

**索引**:
- PRIMARY KEY (id)
- INDEX idx_section_id (section_id)
- INDEX idx_sort_order (sort_order)

**外键**:
- FOREIGN KEY (section_id) REFERENCES series_sections(id) ON DELETE CASCADE

---

### 14. 子章节文章关联表 (subchapter_articles)

| 字段名 | 类型 | 长度 | 允许NULL | 默认值 | 说明 |
|--------|------|------|----------|--------|------|
| id | INTEGER | - | NO | - | 主键, 自增 |
| subchapter_id | INTEGER | - | NO | - | 子章节ID, 外键关联series_subchapters.id |
| article_id | INTEGER | - | NO | - | 文章ID, 外键关联articles.id |
| sort_order | INTEGER | - | NO | 0 | 排序顺序 |
| created_at | DATETIME | - | NO | CURRENT_TIMESTAMP | 创建时间 |

**索引**:
- PRIMARY KEY (id)
- UNIQUE INDEX idx_subchapter_article (subchapter_id, article_id)
- INDEX idx_subchapter_id (subchapter_id)
- INDEX idx_article_id (article_id)

**外键**:
- FOREIGN KEY (subchapter_id) REFERENCES series_subchapters(id) ON DELETE CASCADE
- FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE

---

## 表关系图

```
┌─────────────┐
│    users    │
└──────┬──────┘
       │
       ├──┬──────────────────────────────────────────┐
       ││                                           │
┌──────▼─────┐    ┌───────────┐    ┌───────────┐    │
│  articles  │    │  roles    │    │ comments  │    │
└──────┬──────┘    └─────┬─────┘    └─────┬─────┘    │
       │                │                │           │
       │         ┌──────▼─────┐          │           │
       │         │   users   │          │           │
       │         └───────────┘          │           │
       │                               │           │
┌──────▼─────┐    ┌──────────────┐    │           │
│categories  │    │ article_tags │    │           │
└────────────┘    └──────┬───────┘    │           │
                        │            │           │
                  ┌─────▼─────┐  ┌───▼─────┐  ┌──▼─────┐
                  │   tags    │  │comments │  │comments│
                  └───────────┘  │(parent) │  │(child) │
                                 └─────────┘  └────────┘

┌─────────────┐
│   series    │
└──────┬──────┘
       │
┌──────▼──────────┐
│series_sections  │
└──────┬──────────┘
       │
┌──────▼──────────────┐
│series_subchapters   │
└──────┬──────────────┘
       │
┌──────▼──────────────┐
│subchapter_articles  │
└──────┬──────────────┘
       │
┌──────▼─────┐
│  articles  │
└────────────┘
```

---

## 枚举值说明

### 用户角色 (role)
- `0`: 普通用户
- `1`: 管理员

### 用户状态 (status)
- `0`: 禁用
- `1`: 启用

### 文章状态 (status)
- `0`: 草稿
- `1`: 已发布

### 评论状态 (status)
- `0`: 待审核
- `1`: 已通过
- `2`: 已拒绝

---

## 数据迁移

数据库迁移文件位于 `migrations/` 目录：

- `001_init_schema.sql`: 初始化表结构
- `002_seed_data.sql`: 测试数据
- `003_seed_data.sql`: 补充测试数据

执行迁移命令：

```bash
# SQLite
sqlite3 myblog.db < migrations/001_init_schema.sql

# MySQL
mysql -u root -p myblog < migrations/001_init_schema.sql

# PostgreSQL
psql -U postgres -d myblog -f migrations/001_init_schema.sql
```

---

## 备份与恢复

### SQLite 备份
```bash
# 备份
cp myblog.db myblog_backup_$(date +%Y%m%d).db

# 恢复
cp myblog_backup_20240101.db myblog.db
```

### MySQL 备份
```bash
# 备份
mysqldump -u root -p myblog > myblog_backup_$(date +%Y%m%d).sql

# 恢复
mysql -u root -p myblog < myblog_backup_20240101.sql
```

### PostgreSQL 备份
```bash
# 备份
pg_dump -U postgres myblog > myblog_backup_$(date +%Y%m%d).sql

# 恢复
psql -U postgres -d myblog < myblog_backup_20240101.sql
```
