# 博客系统缺失的 API 接口和模型文档

## 概述

本文档详细列出了博客管理平台前端需要但后端尚未实现的所有 API 接口、数据库模型和相关功能。

---

## 一、缺失的核心功能模块

### 1. RBAC 权限管理系统

#### 1.1 角色管理

**前端组件**: `src/views/roles/RoleList.vue`
**前端路由**: `/system/roles`
**权限标识**: `system:role:read`, `system:role:create`, `system:role:update`, `system:role:delete`

**需要的 API 接口**:

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| GET | `/api/v1/admin/roles` | 获取角色列表 | ❌ |
| GET | `/api/v1/admin/roles/:id` | 获取角色详情 | ❌ |
| POST | `/api/v1/admin/roles` | 创建角色 | ❌ |
| PUT | `/api/v1/admin/roles/:id` | 更新角色 | ❌ |
| DELETE | `/api/v1/admin/roles/:id` | 删除角色 | ❌ |
| POST | `/api/v1/admin/roles/:id/permissions` | 分配权限给角色 | ❌ |
| GET | `/api/v1/admin/roles/:id/permissions` | 获取角色权限 | ❌ |

**请求示例**:

```json
// 创建角色
POST /api/v1/admin/roles
{
  "name": "编辑",
  "code": "editor",
  "description": "内容编辑角色",
  "status": 1
}

// 分配权限
POST /api/v1/admin/roles/:id/permissions
{
  "menu_ids": [1, 2, 3, 4, 5, 6]
}
```

---

#### 1.2 菜单管理

**前端组件**: `src/views/menus/MenuList.vue`
**前端路由**: `/system/menus`
**权限标识**: `system:menu:read`, `system:menu:create`, `system:menu:update`, `system:menu:delete`

**需要的 API 接口**:

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| GET | `/api/v1/admin/menus` | 获取菜单列表 | ❌ |
| GET | `/api/v1/admin/menus/tree` | 获取菜单树结构 | ❌ |
| GET | `/api/v1/admin/menus/:id` | 获取菜单详情 | ❌ |
| POST | `/api/v1/admin/menus` | 创建菜单 | ❌ |
| PUT | `/api/v1/admin/menus/:id` | 更新菜单 | ❌ |
| DELETE | `/api/v1/admin/menus/:id` | 删除菜单 | ❌ |

**请求示例**:

```json
// 创建菜单
POST /api/v1/admin/menus
{
  "parent_id": null,
  "title": "内容管理",
  "icon": "el-icon-document",
  "path": "/content",
  "component": null,
  "type": 0,  // 0-目录, 1-菜单, 2-按钮
  "permission": "content:read",
  "sort": 1,
  "status": 1
}
```

---

#### 1.3 用户角色管理

**前端组件**: `src/views/users/UserList.vue`
**前端路由**: `/system/users`
**权限标识**: `system:user:assign-role`

**需要的 API 接口**:

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| GET | `/api/v1/admin/users` | 获取用户列表 | ✅ |
| GET | `/api/v1/admin/users/:id` | 获取用户详情 | ✅ |
| POST | `/api/v1/admin/users` | 创建用户 | ✅ |
| PUT | `/api/v1/admin/users/:id` | 更新用户 | ✅ |
| DELETE | `/api/v1/admin/users/:id` | 删除用户 | ✅ |
| PUT | `/api/v1/admin/users/:id/status` | 更新用户状态 | ✅ |
| POST | `/api/v1/admin/users/:id/reset-password` | 重置用户密码 | ❌ |
| POST | `/api/v1/admin/users/:id/assign-roles` | 分配角色给用户 | ❌ |
| GET | `/api/v1/admin/users/:id/roles` | 获取用户角色 | ❌ |

**请求示例**:

```json
// 重置密码
POST /api/v1/admin/users/:id/reset-password
{
  "new_password": "123456"
}

// 分配角色
POST /api/v1/admin/users/:id/assign-roles
{
  "role_ids": [1, 2]
}
```

---

#### 1.4 权限验证接口

**需要的 API 接口**:

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| GET | `/api/v1/admin/user/permissions` | 获取当前用户权限列表 | ❌ |
| GET | `/api/v1/admin/user/menus` | 获取当前用户菜单列表 | ❌ |

**响应示例**:

```json
// 获取当前用户权限
GET /api/v1/admin/user/permissions
{
  "code": 0,
  "message": "success",
  "data": {
    "permissions": [
      "content:article:read",
      "content:article:create",
      "content:category:read",
      "system:user:read",
      "system:role:read"
    ]
  }
}

// 获取当前用户菜单
GET /api/v1/admin/user/menus
{
  "code": 0,
  "message": "success",
  "data": {
    "menus": [
      {
        "id": 1,
        "parent_id": null,
        "title": "内容管理",
        "icon": "el-icon-document",
        "path": "/content",
        "type": 0,
        "sort": 1,
        "children": [
          {
            "id": 2,
            "parent_id": 1,
            "title": "文章管理",
            "icon": "el-icon-document",
            "path": "/articles",
            "component": "articles/ArticleList",
            "type": 1,
            "permission": "content:article:read",
            "sort": 1
          }
        ]
      }
    ]
  }
}
```

---

### 2. 审核管理系统

#### 2.1 文章审核

**前端组件**: `src/views/audit/ArticleAudit.vue`
**前端路由**: `/audit/articles`
**权限标识**: `audit:article:read`, `audit:article:approve`, `audit:article:reject`

**需要的 API 接口**:

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| GET | `/api/v1/admin/articles/pending` | 获取待审核文章列表 | ❌ |
| GET | `/api/v1/admin/articles/approved` | 获取已审核文章列表 | ❌ |
| GET | `/api/v1/admin/articles/rejected` | 获取已拒绝文章列表 | ❌ |
| PUT | `/api/v1/admin/articles/:id/approve` | 审核通过文章 | ❌ |
| PUT | `/api/v1/admin/articles/:id/reject` | 审核拒绝文章 | ❌ |
| GET | `/api/v1/admin/articles/:id/audit-log` | 获取文章审核日志 | ❌ |

**请求示例**:

```json
// 审核通过
PUT /api/v1/admin/articles/:id/approve
{
  "remark": "内容符合规范"
}

// 审核拒绝
PUT /api/v1/admin/articles/:id/reject
{
  "reason": "内容包含违规信息",
  "remark": "请修改后重新提交"
}
```

---

#### 2.2 评论审核

**前端组件**: `src/views/audit/CommentAudit.vue`
**前端路由**: `/audit/comments`
**权限标识**: `audit:comment:read`, `audit:comment:approve`, `audit:comment:reject`

**需要的 API 接口**:

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| GET | `/api/v1/admin/comments/pending` | 获取待审核评论列表 | ✅ 部分实现 |
| GET | `/api/v1/admin/comments/approved` | 获取已审核评论列表 | ❌ |
| GET | `/api/v1/admin/comments/rejected` | 获取已拒绝评论列表 | ❌ |
| PUT | `/api/v1/admin/comments/:id/approve` | 审核通过评论 | ✅ |
| PUT | `/api/v1/admin/comments/:id/reject` | 审核拒绝评论 | ✅ |
| DELETE | `/api/v1/admin/comments/:id` | 删除评论 | ✅ |
| DELETE | `/api/v1/admin/comments/batch` | 批量删除评论 | ❌ |

---

#### 2.3 用户审核

**前端组件**: `src/views/audit/UserAudit.vue`
**前端路由**: `/audit/users`
**权限标识**: `audit:user:read`, `audit:user:approve`, `audit:user:reject`

**需要的 API 接口**:

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| GET | `/api/v1/admin/users/pending` | 获取待审核用户列表 | ❌ |
| GET | `/api/v1/admin/users/approved` | 获取已审核用户列表 | ❌ |
| GET | `/api/v1/admin/users/rejected` | 获取已拒绝用户列表 | ❌ |
| PUT | `/api/v1/admin/users/:id/approve` | 审核通过用户 | ❌ |
| PUT | `/api/v1/admin/users/:id/reject` | 审核拒绝用户 | ❌ |
| PUT | `/api/v1/admin/users/:id/status` | 启用/禁用用户 | ✅ |

---

### 3. 系统设置管理

**前端组件**: `src/views/system/Settings.vue`
**前端路由**: `/system/settings`
**权限标识**: `system:settings:read`, `system:settings:update`

#### 3.1 系统设置

**需要的 API 接口**:

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| GET | `/api/v1/admin/settings` | 获取系统设置 | ❌ |
| PUT | `/api/v1/admin/settings` | 更新系统设置 | ❌ |
| POST | `/api/v1/admin/settings/reset` | 重置系统设置 | ❌ |

**设置项说明**:

```json
// 网站信息
{
  "site_title": "我的博客",
  "site_description": "一个简洁优雅的博客系统",
  "site_keywords": "博客,技术,分享",
  "site_logo": "/uploads/logo.png"
}

// SEO 设置
{
  "seo_enabled": true,
  "structured_data": true,
  "sitemap_enabled": true,
  "home_seo_title": "首页标题",
  "home_seo_description": "首页描述",
  "home_seo_keywords": "首页关键词"
}

// 系统配置
{
  "page_size": 10,
  "comment_audit_enabled": true,
  "article_audit_enabled": false,
  "registration_enabled": true,
  "timezone": "Asia/Shanghai"
}

// 邮件设置
{
  "smtp_host": "smtp.example.com",
  "smtp_port": 587,
  "smtp_username": "noreply@example.com",
  "smtp_password": "password",
  "smtp_from": "noreply@example.com"
}
```

---

#### 3.2 邮件测试

**需要的 API 接口**:

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| POST | `/api/v1/admin/settings/test-email` | 发送测试邮件 | ❌ |

**请求示例**:

```json
POST /api/v1/admin/settings/test-email
{
  "to": "test@example.com",
  "subject": "测试邮件",
  "content": "这是一封测试邮件"
}
```

---

#### 3.3 备份恢复

**需要的 API 接口**:

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| POST | `/api/v1/admin/backup/create` | 创建备份 | ❌ |
| GET | `/api/v1/admin/backup/list` | 获取备份列表 | ❌ |
| POST | `/api/v1/admin/backup/:id/restore` | 恢复备份 | ❌ |
| DELETE | `/api/v1/admin/backup/:id` | 删除备份 | ❌ |
| GET | `/api/v1/admin/backup/:id/download` | 下载备份 | ❌ |

**请求示例**:

```json
// 创建备份
POST /api/v1/admin/backup/create
{
  "type": "full",  // full-完整备份, incremental-增量备份
  "name": "自动备份-2024-01-01",
  "description": "自动创建的完整备份"
}
```

---

### 4. 数据统计分析

#### 4.1 仪表盘统计

**需要的 API 接口**:

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| GET | `/api/v1/admin/statistics/dashboard` | 获取仪表盘统计数据 | ❌ |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_articles": 120,
    "published_articles": 100,
    "draft_articles": 20,
    "total_views": 15000,
    "total_likes": 500,
    "total_comments": 300,
    "total_users": 50,
    "active_users": 30,
    "today_new_users": 5,
    "online_users": 10
  }
}
```

---

#### 4.2 文章统计

**需要的 API 接口**:

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| GET | `/api/v1/admin/statistics/articles` | 获取文章统计数据 | ❌ |
| GET | `/api/v1/admin/statistics/articles/trend` | 获取文章发布趋势 | ❌ |
| GET | `/api/v1/admin/statistics/views/trend` | 获取浏览量趋势 | ❌ |
| GET | `/api/v1/admin/statistics/articles/hot` | 获取热门文章 | ❌ |
| GET | `/api/v1/admin/statistics/categories` | 获取分类统计 | ❌ |
| GET | `/api/v1/admin/statistics/tags` | 获取标签统计 | ❌ |

**响应示例**:

```json
// 文章发布趋势
GET /api/v1/admin/statistics/articles/trend?days=30
{
  "code": 0,
  "message": "success",
  "data": {
    "dates": ["2024-01-01", "2024-01-02", "..."],
    "counts": [5, 8, 3, 10, ...]
  }
}

// 浏览量趋势
GET /api/v1/admin/statistics/views/trend?days=30
{
  "code": 0,
  "message": "success",
  "data": {
    "dates": ["2024-01-01", "2024-01-02", "..."],
    "views": [100, 150, 200, 180, ...]
  }
}

// 热门文章
GET /api/v1/admin/statistics/articles/hot?limit=10&days=30
{
  "code": 0,
  "message": "success",
  "data": {
    "articles": [
      {
        "id": 1,
        "title": "热门文章标题",
        "views": 1000,
        "likes": 100,
        "comments": 50,
        "published_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}

// 分类统计
GET /api/v1/admin/statistics/categories
{
  "code": 0,
  "message": "success",
  "data": {
    "categories": [
      {
        "id": 1,
        "name": "技术",
        "count": 50
      },
      {
        "id": 2,
        "name": "生活",
        "count": 30
      }
    ]
  }
}

// 标签统计
GET /api/v1/admin/statistics/tags
{
  "code": 0,
  "message": "success",
  "data": {
    "tags": [
      {
        "id": 1,
        "name": "Go",
        "count": 40
      },
      {
        "id": 2,
        "name": "Vue",
        "count": 35
      }
    ]
  }
}
```

---

#### 4.3 用户统计

**需要的 API 接口**:

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| GET | `/api/v1/admin/statistics/users/growth` | 获取用户增长统计 | ❌ |
| GET | `/api/v1/admin/statistics/users/active` | 获取活跃用户统计 | ❌ |
| GET | `/api/v1/admin/statistics/users/distribution` | 获取用户分布统计 | ❌ |
| GET | `/api/v1/admin/statistics/users/online` | 获取在线用户统计 | ❌ |

**响应示例**:

```json
// 用户增长统计
GET /api/v1/admin/statistics/users/growth?days=30
{
  "code": 0,
  "message": "success",
  "data": {
    "dates": ["2024-01-01", "2024-01-02", "..."],
    "counts": [2, 5, 3, 8, ...],
    "total": 100
  }
}

// 活跃用户统计
GET /api/v1/admin/statistics/users/active?days=7
{
  "code": 0,
  "message": "success",
  "data": {
    "active_users": [
      {
        "id": 1,
        "username": "user1",
        "nickname": "用户1",
        "login_count": 10,
        "last_login": "2024-01-01T12:00:00Z"
      }
    ],
    "total": 20
  }
}

// 用户分布统计
GET /api/v1/admin/statistics/users/distribution
{
  "code": 0,
  "message": "success",
  "data": {
    "by_role": [
      {
        "role": "管理员",
        "count": 5
      },
      {
        "role": "编辑",
        "count": 10
      },
      {
        "role": "普通用户",
        "count": 85
      }
    ],
    "by_status": [
      {
        "status": "正常",
        "count": 90
      },
      {
        "status": "禁用",
        "count": 10
      }
    ]
  }
}
```

---

### 5. 其他缺失的 API 接口

#### 5.1 文章批量操作

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| DELETE | `/api/v1/admin/articles/batch` | 批量删除文章 | ❌ |
| PUT | `/api/v1/admin/articles/batch/status` | 批量更新文章状态 | ❌ |
| PUT | `/api/v1/admin/articles/batch/publish` | 批量发布文章 | ❌ |
| PUT | `/api/v1/admin/articles/batch/withdraw` | 批量撤回文章 | ❌ |

**请求示例**:

```json
// 批量删除
DELETE /api/v1/admin/articles/batch
{
  "ids": [1, 2, 3]
}

// 批量发布
PUT /api/v1/admin/articles/batch/publish
{
  "ids": [1, 2, 3]
}
```

---

#### 5.2 评论批量操作

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| DELETE | `/api/v1/admin/comments/batch` | 批量删除评论 | ❌ |
| PUT | `/api/v1/admin/comments/batch/approve` | 批量审核通过评论 | ✅ |
| PUT | `/api/v1/admin/comments/batch/reject` | 批量审核拒绝评论 | ❌ |

---

#### 5.3 用户批量操作

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| DELETE | `/api/v1/admin/users/batch` | 批量删除用户 | ✅ |
| PUT | `/api/v1/admin/users/batch/status` | 批量更新用户状态 | ❌ |
| POST | `/api/v1/admin/users/batch/assign-roles` | 批量分配角色 | ❌ |

---

#### 5.4 标签批量操作

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| DELETE | `/api/v1/admin/tags/batch` | 批量删除标签 | ✅ |

---

#### 5.5 分类批量操作

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| DELETE | `/api/v1/admin/categories/batch` | 批量删除分类 | ❌ |

---

#### 5.6 文件上传

| 方法 | 路径 | 说明 | 是否已实现 |
|------|------|------|-----------|
| POST | `/api/v1/admin/upload/image` | 上传图片 | ❌ |
| POST | `/api/v1/admin/upload/file` | 上传文件 | ❌ |
| DELETE | `/api/v1/admin/upload/:id` | 删除文件 | ❌ |

**请求示例**:

```json
// 上传图片（FormData）
POST /api/v1/admin/upload/image
Content-Type: multipart/form-data

file: <binary>
type: "cover" | "avatar" | "logo"
```

---

## 二、缺失的数据库模型

### 1. RBAC 相关模型

#### 1.1 角色表 (roles) - 已有但不完整

**当前状态**: 已创建但缺少字段
**文件位置**: `internal/models/role.go`

**需要添加的字段**:

```go
type Role struct {
    ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string    `gorm:"type:varchar(50);not null" json:"name"`
    Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
    Description string    `gorm:"type:text" json:"description"`  // 新增
    Status      int       `gorm:"type:int;default:1" json:"status"`
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

    // 关联关系
    Menus []Menu `gorm:"many2many:role_menus" json:"menus,omitempty"`
    Users []User `gorm:"many2many:user_roles" json:"users,omitempty"`
}
```

---

#### 1.2 菜单表 (menus) - 缺失

**文件位置**: 需要新建 `internal/models/menu.go`

**模型定义**:

```go
package models

import (
    "time"
)

// Menu 菜单模型
type Menu struct {
    ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
    ParentID  *uint      `gorm:"index" json:"parent_id"`           // 父菜单ID
    Title     string     `gorm:"type:varchar(100);not null" json:"title"`      // 菜单标题
    Icon      string     `gorm:"type:varchar(100)" json:"icon"`                // 菜单图标
    Path      string     `gorm:"type:varchar(255)" json:"path"`                // 路由路径
    Component string     `gorm:"type:varchar(255)" json:"component"`           // 组件路径
    Type      int        `gorm:"type:int;not null;default:1" json:"type"`    // 类型: 0-目录, 1-菜单, 2-按钮
    Permission string   `gorm:"type:varchar(100);uniqueIndex" json:"permission"` // 权限标识
    Sort      int        `gorm:"default:0" json:"sort"`                       // 排序
    Status    int        `gorm:"type:int;default:1;index" json:"status"`       // 状态: 0-禁用, 1-启用
    CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

    // 关联关系
    Parent   *Menu   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
    Children []Menu  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
    Roles    []Role  `gorm:"many2many:role_menus" json:"roles,omitempty"`
}

// TableName 指定表名
func (Menu) TableName() string {
    return "menus"
}

// 常量定义
const (
    MenuTypeDirectory = 0 // 目录
    MenuTypeMenu      = 1 // 菜单
    MenuTypeButton    = 2 // 按钮

    MenuStatusDisabled = 0 // 禁用
    MenuStatusEnabled  = 1 // 启用
)
```

---

#### 1.3 角色菜单关联表 (role_menus) - 缺失

**文件位置**: 需要新建 `internal/models/role_menu.go`

**模型定义**:

```go
package models

import "time"

// RoleMenu 角色菜单关联模型
type RoleMenu struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    RoleID    uint      `gorm:"not null;uniqueIndex:idx_role_menu;index" json:"role_id"`
    MenuID    uint      `gorm:"not null;uniqueIndex:idx_role_menu;index" json:"menu_id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

    // 关联关系
    Role Role `gorm:"foreignKey:RoleID" json:"-"`
    Menu Menu `gorm:"foreignKey:MenuID" json:"-"`
}

// TableName 指定表名
func (RoleMenu) TableName() string {
    return "role_menus"
}
```

---

#### 1.4 用户角色关联表 (user_roles) - 缺失

**文件位置**: 需要新建 `internal/models/user_role.go`

**模型定义**:

```go
package models

import "time"

// UserRole 用户角色关联模型
type UserRole struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID    uint      `gorm:"not null;uniqueIndex:idx_user_role;index" json:"user_id"`
    RoleID    uint      `gorm:"not null;uniqueIndex:idx_user_role;index" json:"role_id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

    // 关联关系
    User User `gorm:"foreignKey:UserID" json:"-"`
    Role Role `gorm:"foreignKey:RoleID" json:"-"`
}

// TableName 指定表名
func (UserRole) TableName() string {
    return "user_roles"
}
```

---

### 2. 系统设置相关模型

#### 2.1 系统设置表 (system_settings) - 缺失

**文件位置**: 需要新建 `internal/models/system_setting.go`

**模型定义**:

```go
package models

import "time"

// SystemSetting 系统设置模型
type SystemSetting struct {
    ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Key         string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
    Value       string    `gorm:"type:text" json:"value"`
    Type        string    `gorm:"type:varchar(20);default:'string'" json:"type"` // string, number, boolean, json
    Description string    `gorm:"type:text" json:"description"`
    Group       string    `gorm:"type:varchar(50)" json:"group"` // site, seo, system, email, storage
    UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (SystemSetting) TableName() string {
    return "system_settings"
}

// 常量定义
const (
    SettingTypeString   = "string"
    SettingTypeNumber  = "number"
    SettingTypeBoolean = "boolean"
    SettingTypeJSON    = "json"

    SettingGroupSite     = "site"     // 网站信息
    SettingGroupSEO      = "seo"      // SEO 设置
    SettingGroupSystem   = "system"   // 系统配置
    SettingGroupEmail    = "email"    // 邮件设置
    SettingGroupStorage  = "storage"  // 存储设置
)
```

---

#### 2.2 备份记录表 (backups) - 缺失

**文件位置**: 需要新建 `internal/models/backup.go`

**模型定义**:

```go
package models

import "time"

// Backup 备份记录模型
type Backup struct {
    ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string    `gorm:"type:varchar(255);not null" json:"name"`
    Type        string    `gorm:"type:varchar(20);not null" json:"type"` // full, incremental
    Size        int64     `gorm:"default:0" json:"size"`                 // 备份大小（字节）
    Path        string    `gorm:"type:varchar(500);not null" json:"path"`
    Description string    `gorm:"type:text" json:"description"`
    CreatedBy   uint      `gorm:"index" json:"created_by"`
    CreatedAt   time.Time `gorm:"autoCreateTime;index" json:"created_at"`

    // 关联关系
    Creator *User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TableName 指定表名
func (Backup) TableName() string {
    return "backups"
}

// 常量定义
const (
    BackupTypeFull        = "full"         // 完整备份
    BackupTypeIncremental = "incremental" // 增量备份
)
```

---

### 3. 审核相关模型

#### 3.1 审核日志表 (audit_logs) - 缺失

**文件位置**: 需要新建 `internal/models/audit_log.go`

**模型定义**:

```go
package models

import "time"

// AuditLog 审核日志模型
type AuditLog struct {
    ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    ResourceType string  `gorm:"type:varchar(50);not null;index" json:"resource_type"` // article, comment, user
    ResourceID  uint     `gorm:"not null;index" json:"resource_id"`
    Action      string   `gorm:"type:varchar(20);not null" json:"action"` // approve, reject
    Status      int      `gorm:"type:int;not null" json:"status"` // 1-通过, 2-拒绝
    Reason      string   `gorm:"type:text" json:"reason"`
    Remark      string   `gorm:"type:text" json:"remark"`
    AuditorID   uint     `gorm:"index" json:"auditor_id"`
    CreatedAt   time.Time `gorm:"autoCreateTime;index" json:"created_at"`

    // 关联关系
    Auditor *User `gorm:"foreignKey:AuditorID" json:"auditor,omitempty"`
}

// TableName 指定表名
func (AuditLog) TableName() string {
    return "audit_logs"
}

// 常量定义
const (
    AuditResourceTypeArticle = "article"
    AuditResourceTypeComment = "comment"
    AuditResourceTypeUser    = "user"

    AuditActionApprove = "approve"
    AuditActionReject  = "reject"

    AuditStatusApproved = 1
    AuditStatusRejected = 2
)
```

---

### 4. 统计相关模型

#### 4.1 访问统计表 (visit_stats) - 缺失

**文件位置**: 需要新建 `internal/models/visit_stat.go`

**模型定义**:

```go
package models

import "time"

// VisitStat 访问统计模型
type VisitStat struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Date      time.Time `gorm:"type:date;uniqueIndex;not null" json:"date"`
    Views     int       `gorm:"default:0" json:"views"`       // 浏览量
    Visitors  int       `gorm:"default:0" json:"visitors"`    // 访客数
    PageViews int       `gorm:"default:0" json:"page_views"`  // 页面浏览数
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (VisitStat) TableName() string {
    return "visit_stats"
}
```

---

#### 4.2 文件上传记录表 (uploads) - 缺失

**文件位置**: 需要新建 `internal/models/upload.go`

**模型定义**:

```go
package models

import "time"

// Upload 文件上传记录模型
type Upload struct {
    ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Filename    string    `gorm:"type:varchar(255);not null" json:"filename"`
    OriginalName string   `gorm:"type:varchar(255);not null" json:"original_name"`
    Path        string    `gorm:"type:varchar(500);not null" json:"path"`
    URL         string    `gorm:"type:varchar(500)" json:"url"`
    Size        int64     `gorm:"not null" json:"size"`        // 文件大小（字节）
    MimeType    string    `gorm:"type:varchar(100)" json:"mime_type"`
    Type        string    `gorm:"type:varchar(20);index" json:"type"` // image, file, document, etc.
    UploadedBy  uint      `gorm:"index" json:"uploaded_by"`
    CreatedAt   time.Time `gorm:"autoCreateTime;index" json:"created_at"`

    // 关联关系
    Uploader *User `gorm:"foreignKey:UploadedBy" json:"uploader,omitempty"`
}

// TableName 指定表名
func (Upload) TableName() string {
    return "uploads"
}

// 常量定义
const (
    UploadTypeImage    = "image"
    UploadTypeFile     = "file"
    UploadTypeDocument = "document"
    UploadTypeVideo    = "video"
    UploadTypeAudio    = "audio"
)
```

---

### 5. 需要修改的现有模型

#### 5.1 用户表 (users) - 需要添加字段

**当前状态**: 已有
**文件位置**: `internal/models/user.go`

**需要修改**:

```go
type User struct {
    // ... 现有字段 ...

    // 新增字段
    LastLoginIP   string    `gorm:"type:varchar(45)" json:"last_login_ip"`
    LastLoginTime time.Time `json:"last_login_time"`
    LoginCount    int       `gorm:"default:0" json:"login_count"`

    // 关联关系新增
    Roles []Role `gorm:"many2many:user_roles" json:"roles,omitempty"`
}
```

---

#### 5.2 文章表 (articles) - 需要添加字段

**当前状态**: 已有
**文件位置**: `internal/models/article.go`

**需要修改**:

```go
type Article struct {
    // ... 现有字段 ...

    // 新增字段
    AuditStatus  int        `gorm:"default:1;index" json:"audit_status"` // 0-待审核, 1-已通过, 2-已拒绝
    AuditReason  string     `gorm:"type:text" json:"audit_reason"`
    AuditRemark  string     `gorm:"type:text" json:"audit_remark"`
    AuditedAt    *time.Time `json:"audited_at"`
    AuditorID    *uint      `gorm:"index" json:"auditor_id"`

    // 关联关系新增
    Auditor *User `gorm:"foreignKey:AuditorID" json:"auditor,omitempty"`
}
```

---

## 三、实施优先级建议

### 🔴 第一阶段：核心 RBAC 系统（最高优先级）

**目标**: 实现完整的权限管理系统

**任务清单**:
1. ✅ 创建菜单表 (menus)
2. ✅ 创建角色菜单关联表 (role_menus)
3. ✅ 创建用户角色关联表 (user_roles)
4. ✅ 修改角色模型，添加 description 字段和关联关系
5. ✅ 修改用户模型，添加角色关联和登录统计字段
6. ✅ 实现菜单管理 API（CRUD）
7. ✅ 实现角色管理 API（CRUD + 权限分配）
8. ✅ 实现用户角色分配 API
9. ✅ 实现权限验证中间件
10. ✅ 实现获取用户权限和菜单的 API

**预计工时**: 3-5 天

---

### 🟡 第二阶段：审核管理系统

**目标**: 实现文章、评论、用户审核流程

**任务清单**:
1. ✅ 创建审核日志表 (audit_logs)
2. ✅ 修改文章模型，添加审核相关字段
3. ✅ 实现文章审核 API
4. ✅ 完善评论审核 API
5. ✅ 实现用户审核 API
6. ✅ 实现审核日志查询 API

**预计工时**: 2-3 天

---

### 🟢 第三阶段：数据统计系统

**目标**: 实现数据统计和可视化

**任务清单**:
1. ✅ 创建访问统计表 (visit_stats)
2. ✅ 实现仪表盘统计 API
3. ✅ 实现文章统计 API（趋势、热门、分类、标签）
4. ✅ 实现用户统计 API（增长、活跃、分布）
5. ✅ 实现浏览量趋势 API
6. ✅ 定时任务：更新统计数据

**预计工时**: 2-3 天

---

### 🔵 第四阶段：系统配置和备份

**目标**: 实现系统设置和备份恢复功能

**任务清单**:
1. ✅ 创建系统设置表 (system_settings)
2. ✅ 创建备份记录表 (backups)
3. ✅ 创建文件上传记录表 (uploads)
4. ✅ 实现系统设置 API
5. ✅ 实现文件上传 API
6. ✅ 实现备份恢复 API
7. ✅ 实现邮件测试 API

**预计工时**: 3-4 天

---

### ⚪ 第五阶段：完善现有功能

**目标**: 完善现有功能的批量操作和细节

**任务清单**:
1. ✅ 实现文章批量操作 API
2. ✅ 实现评论批量操作 API
3. ✅ 实现用户批量操作 API
4. ✅ 实现分类批量操作 API
5. ✅ 完善错误处理和日志记录
6. ✅ 添加 API 文档注释

**预计工时**: 1-2 天

---

## 四、数据库迁移 SQL 脚本

### 创建缺失表的 SQL

```sql
-- ============================================
-- RBAC 相关表
-- ============================================

-- 菜单表
CREATE TABLE IF NOT EXISTS menus (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_id INTEGER,
    title VARCHAR(100) NOT NULL,
    icon VARCHAR(100),
    path VARCHAR(255),
    component VARCHAR(255),
    type INTEGER NOT NULL DEFAULT 1,
    permission VARCHAR(100) UNIQUE,
    sort INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES menus(id) ON DELETE CASCADE
);
CREATE INDEX idx_menus_parent_id ON menus(parent_id);
CREATE INDEX idx_menus_status ON menus(status);

-- 角色菜单关联表
CREATE TABLE IF NOT EXISTS role_menus (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER NOT NULL,
    menu_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role_id, menu_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE
);
CREATE INDEX idx_role_menus_role_id ON role_menus(role_id);
CREATE INDEX idx_role_menus_menu_id ON role_menus(menu_id);

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);

-- ============================================
-- 系统设置相关表
-- ============================================

-- 系统设置表
CREATE TABLE IF NOT EXISTS system_settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT,
    type VARCHAR(20) DEFAULT 'string',
    description TEXT,
    group_name VARCHAR(50),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_system_settings_group ON system_settings(group_name);

-- 备份记录表
CREATE TABLE IF NOT EXISTS backups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL,
    size INTEGER DEFAULT 0,
    path VARCHAR(500) NOT NULL,
    description TEXT,
    created_by INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX idx_backups_created_at ON backups(created_at);
CREATE INDEX idx_backups_type ON backups(type);

-- 文件上传记录表
CREATE TABLE IF NOT EXISTS uploads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filename VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    path VARCHAR(500) NOT NULL,
    url VARCHAR(500),
    size INTEGER NOT NULL,
    mime_type VARCHAR(100),
    type VARCHAR(20),
    uploaded_by INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (uploaded_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX idx_uploads_type ON uploads(type);
CREATE INDEX idx_uploads_created_at ON uploads(created_at);

-- ============================================
-- 审核相关表
-- ============================================

-- 审核日志表
CREATE TABLE IF NOT EXISTS audit_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    resource_type VARCHAR(50) NOT NULL,
    resource_id INTEGER NOT NULL,
    action VARCHAR(20) NOT NULL,
    status INTEGER NOT NULL,
    reason TEXT,
    remark TEXT,
    auditor_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (auditor_id) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_auditor_id ON audit_logs(auditor_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- ============================================
-- 统计相关表
-- ============================================

-- 访问统计表
CREATE TABLE IF NOT EXISTS visit_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date DATE UNIQUE NOT NULL,
    views INTEGER DEFAULT 0,
    visitors INTEGER DEFAULT 0,
    page_views INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_visit_stats_date ON visit_stats(date);

-- ============================================
-- 修改现有表
-- ============================================

-- 修改 users 表，添加新字段
ALTER TABLE users ADD COLUMN last_login_ip VARCHAR(45);
ALTER TABLE users ADD COLUMN last_login_time DATETIME;
ALTER TABLE users ADD COLUMN login_count INTEGER DEFAULT 0;

-- 修改 roles 表，添加 description 字段
ALTER TABLE roles ADD COLUMN description TEXT;

-- 修改 articles 表，添加审核相关字段
ALTER TABLE articles ADD COLUMN audit_status INTEGER DEFAULT 1;
ALTER TABLE articles ADD COLUMN audit_reason TEXT;
ALTER TABLE articles ADD COLUMN audit_remark TEXT;
ALTER TABLE articles ADD COLUMN audited_at DATETIME;
ALTER TABLE articles ADD COLUMN auditor_id INTEGER;
CREATE INDEX idx_articles_audit_status ON articles(audit_status);
CREATE INDEX idx_articles_auditor_id ON articles(auditor_id);
```

---

## 五、总结

### 前端完成度: 85%
- ✅ 完整的 UI 界面
- ✅ 完整的交互逻辑
- ✅ 完善的权限控制机制
- ⚠️ API 接口使用 Mock 数据

### 后端完成度: 40%
- ✅ 基础用户认证
- ✅ 基础内容管理（文章、分类、标签、评论）
- ✅ 基础用户管理
- ❌ RBAC 权限系统（最大缺口）
- ❌ 审核系统
- ❌ 数据统计
- ❌ 系统设置

### 最紧急需要实现的功能

1. **RBAC 权限管理系统** - 前端已完整实现，后端完全缺失
2. **审核管理系统** - 评论审核部分实现，文章和用户审核未实现
3. **数据统计接口** - 前端统计页面需要后端数据支持

### 建议实施顺序

1. 先实现 RBAC 系统，因为这是整个系统权限控制的基础
2. 再实现审核系统，完善内容管理流程
3. 然后实现数据统计，提供数据分析能力
4. 最后完善系统设置和备份功能

---

**文档版本**: v1.0
**最后更新**: 2024-01-16
