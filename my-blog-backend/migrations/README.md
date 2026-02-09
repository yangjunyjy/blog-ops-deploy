# 数据库迁移文档

## 概述

本目录包含博客系统的所有数据库迁移脚本。

## 迁移脚本列表

### MySQL版本

#### 1. 001_init_schema_mysql.sql

初始化数据库结构和基础数据。

**包含内容：**
- 创建所有数据表（14张）
- 创建索引
- 设置外键约束
- 插入默认管理员用户
- 插入默认角色
- 插入默认分类和标签

**使用方法：**
```bash
mysql -u root -p myblog < migrations/001_init_schema_mysql.sql
```

#### 2. 002_drop_schema_mysql.sql

删除所有数据表（用于回滚）。

**注意：**
- 此脚本会永久删除所有数据
- 执行前请务必备份数据
- 删除顺序遵循外键依赖关系

**使用方法：**
```bash
mysql -u root -p myblog < migrations/002_drop_schema_mysql.sql
```

#### 3. 003_seed_data_mysql.sql

插入测试数据（仅用于开发环境）。

**包含内容：**
- 测试用户（user1, user2）
- 更多分类（前端、后端、DevOps等）
- 更多标签（JavaScript、Docker等）
- 专题系列和章节
- 测试文章（5篇）
- 文章标签关联
- 测试评论
- 模拟浏览记录
- 点赞和收藏记录

**使用方法：**
```bash
mysql -u root -p myblog < migrations/003_seed_data_mysql.sql
```

### SQLite版本

#### 1. 001_init_schema.sql

初始化数据库结构和基础数据。

**使用方法：**
```bash
sqlite3 myblog.db < migrations/001_init_schema.sql
```

#### 2. 002_drop_schema.sql

删除所有数据表（用于回滚）。

**使用方法：**
```bash
sqlite3 myblog.db < migrations/002_drop_schema.sql
```

#### 3. 003_seed_data.sql

插入测试数据（仅用于开发环境）。

**使用方法：**
```bash
sqlite3 myblog.db < migrations/003_seed_data.sql
```

## 数据表说明

### 用户相关表

#### user（用户表）
- 存储用户基本信息
- 包含软删除支持
- 支持角色和状态管理

#### roles（角色表）
- 存储系统角色信息
- 用于RBAC权限控制

### 内容管理表

#### categorie（分类表）
- 文章分类管理
- 支持排序和状态

#### tag（标签表）
- 文章标签管理
- 支持多对多关联

#### articles（文章表）
- 文章主表
- 支持草稿和发布状态
- 包含浏览、点赞、收藏统计

#### comment（评论表）
- 评论内容
- 支持回复功能
- 包含审核状态

### 文章交互表

#### article_tag（文章标签关联表）
- 文章和标签多对多关联
- 唯一约束防止重复关联

#### article_likes（文章点赞表）
- 记录用户点赞
- 唯一约束防止重复点赞

#### article_favorites（文章收藏表）
- 记录用户收藏
- 唯一约束防止重复收藏

#### article_views（文章浏览记录表）
- 记录文章浏览
- 支持IP和UserAgent记录

### 专题系列表

#### series（专题系列表）
- 管理文章系列
- 支持封面和图标

#### series_sections（专题章节表）
- 系列章节管理
- 支持排序

#### series_subchapters（专题子章节表）
- 子章节管理
- 可关联多篇文章

#### subchapter_articles（子章节文章关联表）
- 子章节和文章关联
- 支持排序

## 执行顺序

1. 首次部署：执行 `001_init_schema.sql`
2. 开发环境：执行 `003_seed_data.sql`（可选）
3. 需要重置：先执行 `002_drop_schema.sql`，再执行 `001_init_schema.sql`

## 默认数据说明

### 默认管理员

- 用户名：`admin`
- 密码：`admin123`
- 邮箱：`admin@example.com`
- 角色：管理员

### 默认分类

1. 技术（tech）
2. 生活（life）
3. 读书（reading）

### 默认标签

1. Go
2. Vue
3. 数据库
4. 后端

## 注意事项

1. **密码安全**：默认密码仅用于开发，生产环境必须修改
2. **数据备份**：执行删除脚本前务必备份
3. **外键约束**：数据库启用了外键约束，删除数据时注意顺序
4. **索引优化**：已为常用查询字段创建索引
5. **软删除**：用户表支持软删除，其他表使用级联删除

## 扩展指南

### 添加新表

1. 创建新的迁移脚本 `004_add_XXX_table.sql`
2. 在文档中说明表结构和用途
3. 更新本README文档

### 修改表结构

1. 创建新的迁移脚本（如 `005_alter_XXX_table.sql`）
2. 使用 `ALTER TABLE` 语句
3. 考虑数据迁移兼容性

## 常见问题

### Q: 如何重置数据库？

```bash
# 删除数据库
rm myblog.db

# 重新初始化
sqlite3 myblog.db < migrations/001_init_schema.sql
```

### Q: 如何只初始化表结构不插入数据？

只执行 `001_init_schema.sql`，不执行 `003_seed_data.sql`。

### Q: 如何修改默认管理员密码？

修改 `001_init_schema.sql` 中的密码哈希值，或在应用启动后通过API修改。

### Q: 如何在生产环境使用？

1. 执行 `001_init_schema.sql`
2. 不要执行 `003_seed_data.sql`
3. 修改默认管理员密码
4. 根据需要调整数据库配置
