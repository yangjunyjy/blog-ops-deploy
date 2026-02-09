# 博客系统 API 文档

## 基础信息

- **基础URL**: `http://localhost:8080`
- **API版本**: v1
- **认证方式**: JWT Bearer Token
- **数据格式**: JSON

## 通用响应格式

### 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 错误响应
```json
{
  "code": 1,
  "message": "error message",
  "data": null
}
```

### 分页响应
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

## 认证模块 (`/api/v1/public/auth`)

### 1. 获取验证码
```
GET /api/v1/public/auth/captcha
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "captcha_id": "uuid",
    "captcha_image": "base64_image"
  }
}
```

### 2. 用户登录
```
POST /api/v1/public/auth/login
```

**请求体**:
```json
{
  "username": "admin",
  "password": "password",
  "captcha_id": "uuid",
  "captcha_code": "1234"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "jwt_token",
    "refresh_token": "refresh_token",
    "expires_in": 3600,
    "user": {
      "id": 1,
      "username": "admin",
      "nickname": "管理员",
      "email": "admin@example.com",
      "avatar": "",
      "role": 1
    }
  }
}
```

### 3. 刷新令牌
```
POST /api/v1/public/auth/refresh
```

**请求体**:
```json
{
  "refresh_token": "refresh_token"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "new_jwt_token",
    "expires_in": 3600
  }
}
```

### 4. 用户注册
```
POST /api/v1/public/auth/register
```

**请求体**:
```json
{
  "username": "newuser",
  "password": "password",
  "email": "user@example.com",
  "captcha_id": "uuid",
  "captcha_code": "1234"
}
```

### 5. 发送邮箱验证码
```
POST /api/v1/public/auth/send-email-captcha
```

**请求体**:
```json
{
  "email": "user@example.com"
}
```

### 6. 验证邮箱验证码
```
POST /api/v1/public/auth/verify-email-captcha
```

**请求体**:
```json
{
  "email": "user@example.com",
  "captcha_code": "123456"
}
```

---

## 文章模块 (`/api/v1`)

### 公开接口 (无需认证)

#### 1. 获取文章列表
```
GET /api/v1/public/articles
```

**查询参数**:
- `page` (int, 可选): 页码, 默认1
- `page_size` (int, 可选): 每页数量, 默认10
- `category_id` (int, 可选): 分类ID
- `tag_id` (int, 可选): 标签ID
- `status` (int, 可选): 状态(0-草稿,1-已发布)
- `keyword` (string, 可选): 搜索关键词

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "title": "文章标题",
        "slug": "article-slug",
        "summary": "文章摘要",
        "cover": "封面URL",
        "author": {
          "id": 1,
          "username": "admin",
          "nickname": "管理员"
        },
        "category": {
          "id": 1,
          "name": "分类名"
        },
        "tags": [
          {
            "id": 1,
            "name": "标签名"
          }
        ],
        "views": 100,
        "likes": 10,
        "favorites": 5,
        "comment_count": 3,
        "status": 1,
        "is_top": false,
        "published_at": "2024-01-01T00:00:00Z",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

#### 2. 获取文章详情
```
GET /api/v1/public/articles/:id
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "title": "文章标题",
    "slug": "article-slug",
    "summary": "文章摘要",
    "content": "文章内容",
    "cover": "封面URL",
    "author": {
      "id": 1,
      "username": "admin",
      "nickname": "管理员",
      "avatar": "头像URL"
    },
    "category": {
      "id": 1,
      "name": "分类名",
      "slug": "category-slug"
    },
    "tags": [
      {
        "id": 1,
        "name": "标签名",
        "slug": "tag-slug"
      }
    ],
    "views": 100,
    "likes": 10,
    "favorites": 5,
    "comment_count": 3,
    "status": 1,
    "is_top": false,
    "sort_order": 0,
    "published_at": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 3. 查看文章（增加浏览量）
```
POST /api/v1/public/articles/:id/view
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "article_id": 1,
    "views": 101
  }
}
```

#### 4. 获取热门文章
```
GET /api/v1/public/articles/hot
```

**查询参数**:
- `limit` (int, 可选): 数量, 默认10

#### 5. 获取最新文章
```
GET /api/v1/public/articles/recent
```

**查询参数**:
- `limit` (int, 可选): 数量, 默认10

#### 6. 搜索文章
```
GET /api/v1/public/articles/search
```

**查询参数**:
- `keyword` (string, 必填): 搜索关键词
- `page` (int, 可选): 页码
- `page_size` (int, 可选): 每页数量

### 前台接口 (需要登录)

#### 7. 点赞文章
```
POST /api/v1/front/articles/:id/like
```

**请求头**:
```
Authorization: Bearer {access_token}
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "article_id": 1,
    "user_id": 1,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 8. 取消点赞
```
DELETE /api/v1/front/articles/:id/like
```

#### 9. 收藏文章
```
POST /api/v1/front/articles/:id/favorite
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "article_id": 1,
    "user_id": 1,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 10. 取消收藏
```
DELETE /api/v1/front/articles/:id/favorite
```

#### 11. 创建文章
```
POST /api/v1/front/articles
```

**请求体**:
```json
{
  "title": "文章标题",
  "slug": "article-slug",
  "summary": "文章摘要",
  "content": "文章内容",
  "cover": "封面URL",
  "category_id": 1,
  "tag_ids": [1, 2, 3],
  "status": 1,
  "is_top": false,
  "sort_order": 0
}
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "title": "文章标题"
  }
}
```

#### 12. 更新文章
```
PUT /api/v1/front/articles/:id
```

**请求体**: 同创建文章

#### 13. 删除文章
```
DELETE /api/v1/front/articles/:id
```

### 后台管理接口 (需要管理员权限)

#### 14. 更新文章状态
```
PUT /api/v1/admin/articles/:id/status
```

**请求体**:
```json
{
  "status": 1
}
```

---

## 分类模块 (`/api/v1`)

### 公开接口

#### 1. 获取分类列表
```
GET /api/v1/public/categories
```

**查询参数**:
- `page` (int, 可选): 页码
- `page_size` (int, 可选): 每页数量

#### 2. 获取分类树
```
GET /api/v1/public/categories/tree
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "分类名",
      "slug": "category-slug",
      "description": "描述",
      "icon": "图标",
      "sort_order": 0,
      "status": 1,
      "article_count": 10,
      "children": []
    }
  ]
}
```

#### 3. 获取分类详情
```
GET /api/v1/public/categories/:id
```

### 前台接口

#### 4. 创建分类
```
POST /api/v1/front/categories
```

**请求体**:
```json
{
  "name": "分类名",
  "slug": "category-slug",
  "description": "描述",
  "icon": "图标",
  "sort_order": 0,
  "status": 1
}
```

#### 5. 更新分类
```
PUT /api/v1/front/categories/:id
```

**请求体**: 同创建分类

#### 6. 删除分类
```
DELETE /api/v1/front/categories/:id
```

---

## 标签模块 (`/api/v1`)

### 公开接口

#### 1. 获取标签列表
```
GET /api/v1/public/tags
```

**查询参数**:
- `page` (int, 可选): 页码
- `page_size` (int, 可选): 每页数量

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "标签名",
        "slug": "tag-slug",
        "description": "描述",
        "article_count": 5,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

#### 2. 获取标签详情
```
GET /api/v1/public/tags/:id
```

### 前台接口

#### 3. 创建标签
```
POST /api/v1/front/tags
```

**请求体**:
```json
{
  "name": "标签名",
  "slug": "tag-slug",
  "description": "描述"
}
```

#### 4. 更新标签
```
PUT /api/v1/front/tags/:id
```

**请求体**: 同创建标签

#### 5. 删除标签
```
DELETE /api/v1/front/tags/:id
```

---

## 评论模块 (`/api/v1`)

### 公开接口

#### 1. 获取文章评论
```
GET /api/v1/public/articles/:articleId/comments
```

**查询参数**:
- `page` (int, 可选): 页码
- `page_size` (int, 可选): 每页数量
- `status` (int, 可选): 状态(0-待审核,1-已通过,2-已拒绝)

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "content": "评论内容",
        "article_id": 1,
        "user": {
          "id": 1,
          "username": "user",
          "nickname": "用户名",
          "avatar": "头像URL"
        },
        "parent_id": null,
        "parent_user": null,
        "status": 1,
        "ip_address": "192.168.1.1",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10
  }
}
```

### 前台接口

#### 2. 创建评论
```
POST /api/v1/front/comments
```

**请求体**:
```json
{
  "article_id": 1,
  "content": "评论内容",
  "parent_id": null
}
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "content": "评论内容",
    "status": 1
  }
}
```

#### 3. 更新评论
```
PUT /api/v1/front/comments/:id
```

**请求体**:
```json
{
  "content": "更新后的内容"
}
```

#### 4. 删除评论
```
DELETE /api/v1/front/comments/:id
```

### 后台管理接口

#### 5. 获取待审核评论
```
GET /api/v1/admin/comments/pending
```

**查询参数**:
- `page` (int, 可选): 页码
- `page_size` (int, 可选): 每页数量
- `level` (int, 可选): 评论层级(0-顶级,1-二级,2-三级)

#### 6. 审核通过评论
```
PUT /api/v1/admin/comments/:id/approve
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "status": 1
  }
}
```

#### 7. 审核拒绝评论
```
PUT /api/v1/admin/comments/:id/reject
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "status": 2
  }
}
```

---

## 用户模块 (`/api/v1`)

### 公开接口

#### 1. 获取用户列表
```
GET /api/v1/public/users
```

**查询参数**:
- `page` (int, 可选): 页码
- `page_size` (int, 可选): 每页数量
- `keyword` (string, 可选): 搜索关键词
- `status` (int, 可选): 状态(0-禁用,1-启用)

### 前台接口

#### 2. 获取当前用户信息
```
GET /api/v1/front/user/info
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "nickname": "管理员",
    "avatar": "头像URL",
    "bio": "个人简介",
    "website": "个人网站",
    "github": "GitHub",
    "role": 1,
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 3. 更新当前用户信息
```
PUT /api/v1/front/user/info
```

**请求体**:
```json
{
  "nickname": "新昵称",
  "email": "newemail@example.com",
  "avatar": "新头像URL",
  "bio": "新简介",
  "website": "新网站",
  "github": "新GitHub"
}
```

#### 4. 获取用户评论
```
GET /api/v1/front/user/comments
```

**查询参数**:
- `page` (int, 可选): 页码
- `page_size` (int, 可选): 每页数量

### 后台管理接口

#### 5. 创建用户
```
POST /api/v1/admin/users
```

**请求体**:
```json
{
  "username": "newuser",
  "password": "password",
  "email": "user@example.com",
  "nickname": "用户名",
  "role": 0,
  "status": 1
}
```

#### 6. 更新用户
```
PUT /api/v1/admin/users/:id
```

**请求体**: 同创建用户(密码可选)

#### 7. 删除用户
```
DELETE /api/v1/admin/users/:id
```

#### 8. 更新用户状态
```
PUT /api/v1/admin/users/:id/status
```

**请求体**:
```json
{
  "status": 1
}
```

---

## 专题系列模块 (`/api/v1`)

### 公开接口

#### 1. 获取系列列表
```
GET /api/v1/public/series
```

**查询参数**:
- `page` (int, 可选): 页码
- `page_size` (int, 可选): 每页数量

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "系列名称",
        "slug": "series-slug",
        "icon": "图标",
        "description": "描述",
        "cover": "封面URL",
        "sort_order": 0,
        "status": 1,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10
  }
}
```

#### 2. 获取系列详情
```
GET /api/v1/public/series/:id
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "系列名称",
    "slug": "series-slug",
    "icon": "图标",
    "description": "描述",
    "cover": "封面URL",
    "sort_order": 0,
    "status": 1,
    "sections": [
      {
        "id": 1,
        "series_id": 1,
        "name": "章节名称",
        "description": "章节描述",
        "sort_order": 0,
        "subchapters": [
          {
            "id": 1,
            "section_id": 1,
            "name": "子章节名称",
            "description": "子章节描述",
            "sort_order": 0,
            "articles": [
              {
                "id": 1,
                "title": "文章标题",
                "slug": "article-slug",
                "summary": "摘要",
                "cover": "封面"
              }
            ]
          }
        ]
      }
    ],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 前台接口

#### 3. 创建系列
```
POST /api/v1/front/series
```

**请求体**:
```json
{
  "name": "系列名称",
  "slug": "series-slug",
  "icon": "图标",
  "description": "描述",
  "cover": "封面URL",
  "sort_order": 0,
  "status": 1
}
```

#### 4. 更新系列
```
PUT /api/v1/front/series/:id
```

**请求体**: 同创建系列

#### 5. 删除系列
```
DELETE /api/v1/front/series/:id
```

#### 6. 创建章节
```
POST /api/v1/front/series/:seriesId/sections
```

**请求体**:
```json
{
  "name": "章节名称",
  "description": "章节描述",
  "sort_order": 0
}
```

#### 7. 更新章节
```
PUT /api/v1/front/series/:seriesId/sections/:sectionId
```

#### 8. 删除章节
```
DELETE /api/v1/front/series/:seriesId/sections/:sectionId
```

#### 9. 创建子章节
```
POST /api/v1/front/series/:seriesId/sections/:sectionId/subchapters
```

**请求体**:
```json
{
  "name": "子章节名称",
  "description": "子章节描述",
  "sort_order": 0
}
```

#### 10. 更新子章节
```
PUT /api/v1/front/series/:seriesId/sections/:sectionId/subchapters/:subchapterId
```

#### 11. 删除子章节
```
DELETE /api/v1/front/series/:seriesId/sections/:sectionId/subchapters/:subchapterId
```

#### 12. 获取子章节文章
```
GET /api/v1/front/series/:seriesId/sections/:sectionId/subchapters/:subchapterId/articles
```

---

## 统计数据模块 (`/api/v1`)

### 前台接口

#### 1. 获取仪表盘统计
```
GET /api/v1/front/statistics/dashboard
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_articles": 100,
    "total_views": 10000,
    "total_likes": 500,
    "total_comments": 200,
    "recent_articles": 10,
    "recent_comments": 20
  }
}
```

#### 2. 获取文章统计
```
GET /api/v1/front/statistics/articles
```

**查询参数**:
- `days` (int, 可选): 天数, 默认30

#### 3. 获取浏览量趋势
```
GET /api/v1/front/statistics/view-trend
```

**查询参数**:
- `days` (int, 可选): 天数, 默认30

#### 4. 获取用户增长
```
GET /api/v1/front/statistics/user-growth
```

**查询参数**:
- `days` (int, 可选): 天数, 默认30

#### 5. 获取热门文章
```
GET /api/v1/front/statistics/hot-articles
```

**查询参数**:
- `limit` (int, 可选): 数量, 默认10
- `days` (int, 可选): 天数, 默认30

#### 6. 获取热门标签
```
GET /api/v1/front/statistics/hot-tags
```

**查询参数**:
- `limit` (int, 可选): 数量, 默认10

---

## 状态码说明

| code | message | 说明 |
|------|---------|------|
| 0 | success | 成功 |
| 1 | error | 通用错误 |
| 1001 | 参数错误 | 请求参数验证失败 |
| 1002 | 未授权 | 需要登录 |
| 1003 | 权限不足 | 需要管理员权限 |
| 1004 | 资源不存在 | 请求的资源未找到 |
| 1005 | 资源已存在 | 资源冲突 |
| 2001 | 用户名或密码错误 | 登录失败 |
| 2002 | 验证码错误 | 验证码不正确 |
| 2003 | 验证码已过期 | 验证码超时 |
| 2004 | 用户已存在 | 注册时用户名已存在 |
| 2005 | 邮箱验证码错误 | 邮箱验证码不正确 |
| 3001 | 文章不存在 | 文章未找到 |
| 3002 | 评论不存在 | 评论未找到 |
| 4001 | 评论审核已禁用 | 系统未开启评论审核 |
