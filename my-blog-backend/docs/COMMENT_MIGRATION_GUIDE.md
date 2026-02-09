# 评论系统迁移指南

## 问题说明

当前评论系统依赖外键约束，但由于模型修改（移除所有外键），以下代码无法正常工作：

1. `comment_repo.go` 第 103 行的 `Preload("User")` 和 `Preload("Article")`
2. 模型中的 `foreignKey` 标签失效
3. 新增的 `user_comment_article` 中间表未被使用

## 解决方案

### 1. 数据库迁移

执行 SQL 文件创建表结构：

```bash
mysql -u root -p your_database < migrations/create_user_comment_tables.sql
```

或者手动执行：

```sql
-- 创建用户表（如果不存在）
CREATE TABLE IF NOT EXISTS `user` (
  -- ... 参考 migrations/create_user_comment_tables.sql
);

-- 创建评论表（如果不存在）
CREATE TABLE IF NOT EXISTS `comment` (
  -- ... 参考 migrations/create_user_comment_tables.sql
);

-- 创建用户评论中间表
CREATE TABLE IF NOT EXISTS `user_comment_article` (
  -- ... 参考 migrations/create_user_comment_tables.sql
);
```

### 2. 代码迁移步骤

#### 步骤 1：替换模型文件

```bash
# 备份原有文件
mv internal/models/user.go internal/models/user.go.backup
mv internal/models/comment.go internal/models/comment.go.backup
mv internal/models/article.go internal/models/article.go.backup

# 使用新的模型文件
mv internal/models/user_updated.go internal/models/user.go
mv internal/models/comment_updated.go internal/models/comment.go
mv internal/models/article_updated.go internal/models/article.go
```

#### 步骤 2：替换 Repository 实现

```bash
# 备份原有文件
mv internal/repository/comment_repository.go internal/repository/comment_repository.go.backup
mv internal/repository/impl/mysql/comment_repo.go internal/repository/impl/mysql/comment_repo.go.backup

# 使用新的实现
mv internal/repository/comment_repository_updated.go internal/repository/comment_repository.go
mv internal/repository/impl/mysql/comment_repo_updated.go internal/repository/impl/mysql/comment_repo.go
```

#### 步骤 3：更新 Service（如需使用中间表）

如果需要在 Service 层使用中间表，添加新的服务方法：

```go
// 在 comment_service.go 中添加
ListUserCommentedArticles(userID uint, page, pageSize int) ([]*models.Article, int64, error)
```

```go
// 在 comment_service.go 中实现
func (s *commentService) ListUserCommentedArticles(userID uint, page, pageSize int) ([]*models.Article, int64, error) {
    return s.commentRepo.ListUserCommentedArticles(userID, page, pageSize)
}
```

#### 步骤 4：添加新的 API 接口（可选）

在 `comment_handler.go` 中添加：

```go
// GetUserCommentedArticles 获取用户评论过的文章列表
func (h *CommentHandler) GetUserCommentedArticles(c *gin.Context) {
    userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "无效的用户ID", err)
        return
    }

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 10
    }

    articles, total, err := h.commentService.ListUserCommentedArticles(uint(userID), page, pageSize)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "获取文章列表失败", err)
        return
    }

    response.Success(c, gin.H{
        "items":     articles,
        "total":     total,
        "page":      page,
        "page_size": pageSize,
    }, "")
}
```

#### 步骤 5：更新路由配置

在 `router.go` 中添加新路由：

```go
comments.GET("/user/:user_id/articles", handlers.Comment.GetUserCommentedArticles)
```

## 主要变化

### 模型层变化

**原 user.go**：
```go
Comments []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"`
```

**新 user.go**：
```go
Comments []Comment `gorm:"-" json:"comments,omitempty"`  // 手动查询
```

### Repository 层变化

**原 comment_repo.go**：
```go
err := query.Preload("User").Preload("Article").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&comments).Error
```

**新 comment_repo.go**：
```go
// 手动加载关联数据
var users []models.User
r.db.Where("id IN ?", userIDs).Find(&users)
// ... 组装数据
```

### 中间表使用

**创建评论时自动维护中间表**：
```go
// 插入评论者记录
INSERT IGNORE INTO user_comment_article (user_id, article_id, comment_id, is_author)
VALUES (userID, articleID, commentID, 0)

// 插入作者记录
INSERT IGNORE INTO user_comment_article (user_id, article_id, comment_id, is_author)
VALUES (authorID, articleID, commentID, 1)
```

**使用中间表查询**：
```go
// 查询用户评论过的文章
SELECT a.* FROM articles a
INNER JOIN user_comment_article uca ON a.id = uca.article_id
WHERE uca.user_id = ? AND a.status = 1
```

## 性能优化

1. **批量加载**：使用 `WHERE id IN ?` 批量加载关联数据，减少数据库查询次数
2. **索引优化**：中间表已有 `(user_id, article_id, comment_id)` 唯一索引和三个单独索引
3. **事务处理**：创建和删除评论时使用事务，保证数据一致性

## 测试清单

- [ ] 创建评论时中间表正确插入
- [ ] 删除评论时中间表正确删除
- [ ] 查询评论列表时用户信息正确加载
- [ ] 通过中间表查询用户评论过的文章
- [ ] 触发器正确工作
- [ ] 性能测试（对比迁移前后）

## 回滚方案

如需回滚，执行以下步骤：

```bash
# 恢复备份的文件
mv internal/models/user.go.backup internal/models/user.go
mv internal/models/comment.go.backup internal/models/comment.go
mv internal/models/article.go.backup internal/models/article.go
mv internal/repository/comment_repository.go.backup internal/repository/comment_repository.go
mv internal/repository/impl/mysql/comment_repo.go.backup internal/repository/impl/mysql/comment_repo.go

# 可选：删除中间表
DROP TABLE IF EXISTS user_comment_article;
DROP TRIGGER IF EXISTS trigger_comment_insert;
DROP TRIGGER IF EXISTS trigger_comment_delete;
```

## 注意事项

1. **数据迁移**：如果已有评论数据，需要手动插入到中间表：

```sql
INSERT INTO user_comment_article (user_id, article_id, comment_id, is_author)
SELECT
  c.user_id,
  c.article_id,
  c.id,
  CASE WHEN c.user_id = a.author_id THEN 1 ELSE 0 END
FROM comment c
INNER JOIN articles a ON c.article_id = a.id;
```

2. **外键检查**：确保原表没有外键约束，如有需要先删除：

```sql
-- 查看外键约束
SHOW CREATE TABLE comment;

-- 删除外键约束（如果存在）
ALTER TABLE comment DROP FOREIGN KEY fk_comment_article;
ALTER TABLE comment DROP FOREIGN KEY fk_comment_user;
```

3. **触发器备份**：迁移时注意保留或重建其他触发器
