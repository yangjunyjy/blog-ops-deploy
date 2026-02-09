# 博客系统后端资产模型文档

## 模型概述

本系统采用 Go 语言 + GORM 框架实现，所有模型位于 `internal/models/` 目录。

模型设计遵循以下原则：
- 使用 GORM 标签定义字段属性
- 支持软删除
- 包含 JSON 序列化标签
- 包含数据验证标签
- 支持关联关系（一对一、一对多、多对多）

---

## 用户模型

### 文件位置
`internal/models/user.go`

### 结构体定义

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

// User 用户模型
type User struct {
    ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
    Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
    Email     string         `gorm:"type:varchar(100);uniqueIndex" json:"email"`
    Password  string         `gorm:"type:varchar(255);not null" json:"-"` // 不返回密码
    Nickname  string         `gorm:"type:varchar(50)" json:"nickname"`
    Avatar    string         `gorm:"type:varchar(255)" json:"avatar"`
    Bio       string         `gorm:"type:text" json:"bio"`
    Website   string         `gorm:"type:varchar(255)" json:"website"`
    GitHub    string         `gorm:"type:varchar(255)" json:"github"`
    Role      int            `gorm:"type:int;default:0;index" json:"role"` // 0-普通用户, 1-管理员
    Status    int            `gorm:"type:int;default:1;index" json:"status"` // 0-禁用, 1-启用
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

    // 关联关系
    Articles      []Article        `gorm:"foreignKey:AuthorID" json:"articles,omitempty"`
    Comments      []Comment        `gorm:"foreignKey:UserID" json:"comments,omitempty"`
    ArticleLikes  []ArticleLike    `gorm:"foreignKey:UserID" json:"article_likes,omitempty"`
    Favorites     []ArticleFavorite `gorm:"foreignKey:UserID" json:"favorites,omitempty"`
    ArticleViews  []ArticleView    `gorm:"foreignKey:UserID" json:"article_views,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
    return "users"
}
```

### 关联关系
- 一对多：用户 -> 文章
- 一对多：用户 -> 评论
- 一对多：用户 -> 文章点赞
- 一对多：用户 -> 文章收藏
- 一对多：用户 -> 文章浏览

### 常量定义

```go
const (
    UserRoleNormal  = 0 // 普通用户
    UserRoleAdmin   = 1 // 管理员

    UserStatusDisabled = 0 // 禁用
    UserStatusEnabled  = 1 // 启用
)
```

---

## 角色模型

### 文件位置
`internal/models/role.go`

### 结构体定义

```go
package models

import "time"

// Role 角色模型
type Role struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name      string    `gorm:"type:varchar(50);not null" json:"name"`
    Code      string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
    Status    int       `gorm:"type:int;default:1" json:"status"` // 0-禁用, 1-启用
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Role) TableName() string {
    return "roles"
}
```

### 常量定义

```go
const (
    RoleStatusDisabled = 0 // 禁用
    RoleStatusEnabled  = 1 // 启用
)
```

---

## 文章模型

### 文件位置
`internal/models/article.go`

### 结构体定义

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

// Article 文章模型
type Article struct {
    ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
    Title        string         `gorm:"type:varchar(200);not null" json:"title"`
    Slug         string         `gorm:"type:varchar(200);uniqueIndex;not null" json:"slug"`
    Summary      string         `gorm:"type:text" json:"summary"`
    Content      string         `gorm:"type:text;not null" json:"content"`
    Cover        string         `gorm:"type:varchar(255)" json:"cover"`
    CategoryID   *uint          `gorm:"index" json:"category_id"`
    AuthorID     uint           `gorm:"not null;index" json:"author_id"`
    Views        int            `gorm:"default:0" json:"views"`
    Likes        int            `gorm:"default:0" json:"likes"`
    Favorites    int            `gorm:"default:0" json:"favorites"`
    CommentCount int            `gorm:"default:0" json:"comment_count"`
    Status       int            `gorm:"default:0;index" json:"status"` // 0-草稿, 1-已发布
    IsTop        bool           `gorm:"default:false;index" json:"is_top"`
    SortOrder    int            `gorm:"default:0" json:"sort_order"`
    PublishedAt  *time.Time     `gorm:"index" json:"published_at"`
    CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

    // 关联关系
    Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
    Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
    Tags     []Tag     `gorm:"many2many:article_tags" json:"tags,omitempty"`
    Comments []Comment `gorm:"foreignKey:ArticleID" json:"comments,omitempty"`
    Likes    []ArticleLike `gorm:"foreignKey:ArticleID" json:"-"`
    FavoritesList []ArticleFavorite `gorm:"foreignKey:ArticleID" json:"-"`
    Views    []ArticleView `gorm:"foreignKey:ArticleID" json:"-"`
}

// TableName 指定表名
func (Article) TableName() string {
    return "articles"
}
```

### 关联关系
- 多对一：文章 -> 作者
- 多对一：文章 -> 分类
- 多对多：文章 <-> 标签
- 一对多：文章 -> 评论
- 一对多：文章 -> 点赞
- 一对多：文章 -> 收藏
- 一对多：文章 -> 浏览记录

### 常量定义

```go
const (
    ArticleStatusDraft     = 0 // 草稿
    ArticleStatusPublished = 1 // 已发布
)
```

---

## 文章标签关联模型

### 文件位置
`internal/models/article_tag.go`

### 结构体定义

```go
package models

import "time"

// ArticleTag 文章标签关联模型
type ArticleTag struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    ArticleID uint      `gorm:"not null;index:idx_article_tag;index" json:"article_id"`
    TagID     uint      `gorm:"not null;index:idx_article_tag;index" json:"tag_id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

    // 关联关系
    Article Article `gorm:"foreignKey:ArticleID" json:"-"`
    Tag     Tag     `gorm:"foreignKey:TagID" json:"-"`
}

// TableName 指定表名
func (ArticleTag) TableName() string {
    return "article_tags"
}
```

---

## 文章点赞模型

### 文件位置
`internal/models/article_like.go`

### 结构体定义

```go
package models

import "time"

// ArticleLike 文章点赞模型
type ArticleLike struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    ArticleID uint      `gorm:"not null;uniqueIndex:idx_article_user;index" json:"article_id"`
    UserID    uint      `gorm:"not null;uniqueIndex:idx_article_user;index" json:"user_id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

    // 关联关系
    Article Article `gorm:"foreignKey:ArticleID" json:"-"`
    User    User    `gorm:"foreignKey:UserID" json:"-"`
}

// TableName 指定表名
func (ArticleLike) TableName() string {
    return "article_likes"
}
```

---

## 文章收藏模型

### 文件位置
`internal/models/article_favorite.go`

### 结构体定义

```go
package models

import "time"

// ArticleFavorite 文章收藏模型
type ArticleFavorite struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    ArticleID uint      `gorm:"not null;uniqueIndex:idx_article_user;index" json:"article_id"`
    UserID    uint      `gorm:"not null;uniqueIndex:idx_article_user;index" json:"user_id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

    // 关联关系
    Article Article `gorm:"foreignKey:ArticleID" json:"-"`
    User    User    `gorm:"foreignKey:UserID" json:"-"`
}

// TableName 指定表名
func (ArticleFavorite) TableName() string {
    return "article_favorites"
}
```

---

## 文章浏览模型

### 文件位置
`internal/models/article_view.go`

### 结构体定义

```go
package models

import "time"

// ArticleView 文章浏览记录模型
type ArticleView struct {
    ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    ArticleID  uint      `gorm:"not null;index" json:"article_id"`
    UserID     *uint     `gorm:"index" json:"user_id"` // 可为空, 记录游客浏览
    IPAddress  string    `gorm:"type:varchar(45)" json:"ip_address"`
    UserAgent  string    `gorm:"type:varchar(500)" json:"user_agent"`
    ViewedAt   time.Time `gorm:"autoCreateTime" json:"viewed_at"`

    // 关联关系
    Article Article `gorm:"foreignKey:ArticleID" json:"-"`
    User    *User   `gorm:"foreignKey:UserID" json:"-"`
}

// TableName 指定表名
func (ArticleView) TableName() string {
    return "article_views"
}
```

---

## 分类模型

### 文件位置
`internal/models/category.go`

### 结构体定义

```go
package models

import "time"

// Category 分类模型
type Category struct {
    ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string    `gorm:"type:varchar(50);not null" json:"name"`
    Slug        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
    Description string    `gorm:"type:text" json:"description"`
    Icon        string    `gorm:"type:varchar(100)" json:"icon"`
    SortOrder   int       `gorm:"default:0" json:"sort_order"`
    Status      int       `gorm:"default:1;index" json:"status"` // 0-禁用, 1-启用
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

    // 关联关系
    Articles []Article `gorm:"foreignKey:CategoryID" json:"articles,omitempty"`
}

// TableName 指定表名
func (Category) TableName() string {
    return "categories"
}
```

### 关联关系
- 一对多：分类 -> 文章

### 常量定义

```go
const (
    CategoryStatusDisabled = 0 // 禁用
    CategoryStatusEnabled  = 1 // 启用
)
```

---

## 标签模型

### 文件位置
`internal/models/tag.go`

### 结构体定义

```go
package models

import "time"

// Tag 标签模型
type Tag struct {
    ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
    Slug        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
    Description string    `gorm:"type:text" json:"description"`
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

    // 关联关系
    Articles []Article `gorm:"many2many:article_tags" json:"articles,omitempty"`
}

// TableName 指定表名
func (Tag) TableName() string {
    return "tags"
}
```

### 关联关系
- 多对多：标签 <-> 文章

---

## 评论模型

### 文件位置
`internal/models/comment.go`

### 结构体定义

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
    ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
    ArticleID uint           `gorm:"not null;index" json:"article_id"`
    UserID    uint           `gorm:"not null;index" json:"user_id"`
    ParentID  *uint          `gorm:"index" json:"parent_id"` // 支持多级评论
    Content   string         `gorm:"type:text;not null" json:"content"`
    Status    int            `gorm:"default:0;index" json:"status"` // 0-待审核, 1-已通过, 2-已拒绝
    IPAddress string         `gorm:"type:varchar(45)" json:"ip_address"`
    UserAgent string         `gorm:"type:varchar(500)" json:"user_agent"`
    CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`

    // 关联关系
    Article    Article     `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
    User       User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Parent     *Comment    `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
    ParentUser *User       `gorm:"foreignKey:ParentID;references:UserID" json:"parent_user,omitempty"`
    Replies    []Comment   `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// TableName 指定表名
func (Comment) TableName() string {
    return "comments"
}
```

### 关联关系
- 多对一：评论 -> 文章
- 多对一：评论 -> 用户
- 自关联：评论 -> 父评论（支持多级评论）
- 一对多：评论 -> 子评论

### 常量定义

```go
const (
    CommentStatusPending  = 0 // 待审核
    CommentStatusApproved = 1 // 已通过
    CommentStatusRejected = 2 // 已拒绝
)
```

---

## 专题系列模型

### 文件位置
`internal/models/series.go`

### 结构体定义

```go
package models

import "time"

// Series 专题系列模型
type Series struct {
    ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string    `gorm:"type:varchar(100);not null" json:"name"`
    Slug        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
    Icon        string    `gorm:"type:varchar(100)" json:"icon"`
    Description string    `gorm:"type:text" json:"description"`
    Cover       string    `gorm:"type:varchar(255)" json:"cover"`
    SortOrder   int       `gorm:"default:0" json:"sort_order"`
    Status      int       `gorm:"default:1;index" json:"status"` // 0-禁用, 1-启用
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

    // 关联关系
    Sections []SeriesSection `gorm:"foreignKey:SeriesID" json:"sections,omitempty"`
}

// TableName 指定表名
func (Series) TableName() string {
    return "series"
}
```

### 关联关系
- 一对多：系列 -> 章节

### 常量定义

```go
const (
    SeriesStatusDisabled = 0 // 禁用
    SeriesStatusEnabled  = 1 // 启用
)
```

---

## 系列章节模型

### 文件位置
`internal/models/serie_section.go`

### 结构体定义

```go
package models

import "time"

// SeriesSection 系列章节模型
type SeriesSection struct {
    ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    SeriesID    uint      `gorm:"not null;index" json:"series_id"`
    Name        string    `gorm:"type:varchar(100);not null" json:"name"`
    Description string    `gorm:"type:text" json:"description"`
    SortOrder   int       `gorm:"default:0;index" json:"sort_order"`
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

    // 关联关系
    Series      Series             `gorm:"foreignKey:SeriesID" json:"series,omitempty"`
    Subchapters []SeriesSubchapter `gorm:"foreignKey:SectionID" json:"subchapters,omitempty"`
}

// TableName 指定表名
func (SeriesSection) TableName() string {
    return "series_sections"
}
```

### 关联关系
- 多对一：章节 -> 系列
- 一对多：章节 -> 子章节

---

## 系列子章节模型

### 文件位置
`internal/models/series_subchapter.go`

### 结构体定义

```go
package models

import "time"

// SeriesSubchapter 系列子章节模型
type SeriesSubchapter struct {
    ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    SectionID   uint      `gorm:"not null;index" json:"section_id"`
    Name        string    `gorm:"type:varchar(100);not null" json:"name"`
    Description string    `gorm:"type:text" json:"description"`
    SortOrder   int       `gorm:"default:0;index" json:"sort_order"`
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

    // 关联关系
    Section  *SeriesSection `gorm:"foreignKey:SectionID" json:"section,omitempty"`
    Articles []Article      `gorm:"many2many:subchapter_articles" json:"articles,omitempty"`
}

// TableName 指定表名
func (SeriesSubchapter) TableName() string {
    return "series_subchapters"
}
```

### 关联关系
- 多对一：子章节 -> 章节
- 多对多：子章节 <-> 文章

---

## 子章节文章关联模型

### 文件位置
`internal/models/subchapter_article.go`

### 结构体定义

```go
package models

import "time"

// SubchapterArticle 子章节文章关联模型
type SubchapterArticle struct {
    ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    SubchapterID uint      `gorm:"not null;uniqueIndex:idx_subchapter_article;index" json:"subchapter_id"`
    ArticleID    uint      `gorm:"not null;uniqueIndex:idx_subchapter_article;index" json:"article_id"`
    SortOrder    int       `gorm:"default:0" json:"sort_order"`
    CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

    // 关联关系
    Subchapter SeriesSubchapter `gorm:"foreignKey:SubchapterID" json:"-"`
    Article    Article           `gorm:"foreignKey:ArticleID" json:"-"`
}

// TableName 指定表名
func (SubchapterArticle) TableName() string {
    return "subchapter_articles"
}
```

---

## 模型关系总结

### 用户中心关系
- **User (用户)**
  - 一对多 → Article (文章)
  - 一对多 → Comment (评论)
  - 一对多 → ArticleLike (文章点赞)
  - 一对多 → ArticleFavorite (文章收藏)
  - 一对多 → ArticleView (文章浏览)

### 内容组织关系
- **Article (文章)**
  - 多对一 → User (作者)
  - 多对一 → Category (分类)
  - 多对多 → Tag (标签)
  - 一对多 → Comment (评论)
  - 一对多 → ArticleLike (点赞)
  - 一对多 → ArticleFavorite (收藏)
  - 一对多 → ArticleView (浏览记录)

- **Category (分类)**
  - 一对多 → Article (文章)

- **Tag (标签)**
  - 多对多 → Article (文章)

### 评论关系
- **Comment (评论)**
  - 多对一 → Article (文章)
  - 多对一 → User (用户)
  - 自关联 → Comment (父评论/子评论)

### 专题系列关系
- **Series (系列)**
  - 一对多 → SeriesSection (章节)

- **SeriesSection (章节)**
  - 多对一 → Series (系列)
  - 一对多 → SeriesSubchapter (子章节)

- **SeriesSubchapter (子章节)**
  - 多对一 → SeriesSection (章节)
  - 多对多 → Article (文章)

---

## 模型使用示例

### 创建文章
```go
article := &models.Article{
    Title:      "我的第一篇文章",
    Slug:       "my-first-article",
    Summary:    "这是文章摘要",
    Content:    "这是文章内容",
    AuthorID:   1,
    CategoryID: &[]uint{1}[0],
    Status:     models.ArticleStatusPublished,
}

// 添加标签
tags := []models.Tag{
    {ID: 1},
    {ID: 2},
}
article.Tags = tags

db.Create(article)
```

### 获取文章及其关联
```go
var article models.Article
db.Preload("Author").
   Preload("Category").
   Preload("Tags").
   Preload("Comments", func(db *gorm.DB) *gorm.DB {
       return db.Where("status = ?", models.CommentStatusApproved)
   }).
   Where("id = ?", 1).
   First(&article)
```

### 点赞文章
```go
like := &models.ArticleLike{
    ArticleID: 1,
    UserID:    1,
}
db.Create(like)

// 更新文章点赞数
db.Model(&models.Article{ID: 1}).
   UpdateColumn("likes", gorm.Expr("likes + ?", 1))
```

### 创建嵌套评论
```go
comment := &models.Comment{
    ArticleID: 1,
    UserID:    1,
    ParentID:  &[]uint{10}[0], // 回复ID为10的评论
    Content:   "这是回复",
    Status:    models.CommentStatusPending,
}
db.Create(comment)
```

---

## GORM 标签说明

### 常用标签
- `primaryKey`: 主键
- `autoIncrement`: 自增
- `not null`: 非空
- `uniqueIndex`: 唯一索引
- `index`: 普通索引
- `default`: 默认值
- `size`: 字段长度
- `type`: 字段类型
- `foreignKey`: 外键
- `many2many`: 多对多关系
- `autoCreateTime`: 自动创建时间
- `autoUpdateTime`: 自动更新时间

### JSON 标签
- `json:"field_name"`: JSON 字段名
- `json:"field_name,omitempty"`: 如果字段为零值则省略
- `json:"-"`: 不在 JSON 中序列化（如密码）

---

## 软删除

所有支持软删除的模型都包含 `DeletedAt gorm.DeletedAt` 字段：

```go
DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
```

使用 GORM 的 `Delete` 方法时，会自动设置 `deleted_at` 字段，而不是物理删除。

### 查询时排除已删除记录
```go
db.Where("id = ?", 1).First(&user) // 自动排除 deleted_at 不为 NULL 的记录
```

### 包含已删除记录
```go
db.Unscoped().Where("id = ?", 1).First(&user) // 包含所有记录
```

### 永久删除
```go
db.Unscoped().Delete(&user) // 物理删除
```

---

## 数据验证

GORM 支持通过结构体标签进行数据验证：

```go
type User struct {
    Username string `gorm:"type:varchar(50);not null" validate:"required,min=3,max=50"`
    Email    string `gorm:"type:varchar(100);uniqueIndex" validate:"required,email"`
    Password string `gorm:"type:varchar(255);not null" validate:"required,min=6"`
}
```

使用验证器：
```go
import "github.com/go-playground/validator/v10"

validate := validator.New()
err := validate.Struct(user)
```

---

## 迁移说明

### 自动迁移
```go
import (
    "gorm.io/gorm"
    "your-project/internal/models"
)

func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
        &models.Role{},
        &models.Article{},
        &models.Category{},
        &models.Tag{},
        &models.ArticleTag{},
        &models.ArticleLike{},
        &models.ArticleFavorite{},
        &models.ArticleView{},
        &models.Comment{},
        &models.Series{},
        &models.SeriesSection{},
        &models.SeriesSubchapter{},
        &models.SubchapterArticle{},
    )
}
```

### 手动迁移
使用 SQL 迁移脚本：
```bash
sqlite3 myblog.db < migrations/001_init_schema.sql
```

---

## 性能优化建议

1. **合理使用索引**：为频繁查询的字段添加索引
2. **预加载关联**：使用 `Preload` 避免 N+1 查询问题
3. **分页查询**：使用 `Limit` 和 `Offset` 进行分页
4. **选择字段**：使用 `Select` 只查询需要的字段
5. **缓存热点数据**：使用 Redis 缓存热点数据
6. **批量操作**：使用 `CreateInBatches` 和 `Update` 进行批量操作

### 分页查询示例
```go
var articles []models.Article
db.Model(&models.Article{}).
   Count(&total).
   Offset((page - 1) * pageSize).
   Limit(pageSize).
   Find(&articles)
```

### 选择字段示例
```go
var articles []models.Article
db.Model(&models.Article{}).
   Select("id", "title", "slug", "summary", "created_at").
   Find(&articles)
```

---

## 总结

本系统的模型设计具有以下特点：

1. **清晰的分层架构**：模型、仓储、服务、处理各层分离
2. **完整的关联关系**：支持一对一、一对多、多对多关系
3. **软删除支持**：重要数据支持软删除，可恢复
4. **完善的索引**：为查询优化建立了合理的索引
5. **数据验证**：通过标签实现数据验证
6. **JSON 序列化**：支持 API 响应的 JSON 序列化
7. **类型安全**：使用 Go 的类型系统保证代码安全
8. **易于扩展**：模型结构清晰，易于添加新字段和关系
