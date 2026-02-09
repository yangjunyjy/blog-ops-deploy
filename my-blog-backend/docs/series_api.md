# Series API 文档

## 概述

Series（系列）模块提供了完整的三级结构管理功能：系列 -> 章节 -> 子章节。

## 基础路径

所有API的基础路径为：`/api/v1/series`

## 认证

大部分接口需要认证，请在请求头中携带JWT Token：
```
Authorization: Bearer <token>
```

---

## 系列管理 (Series)

### 获取系列列表

**请求**
```http
GET /api/v1/series?page=1&page_size=10
```

**参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认10 |

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "Go语言入门",
        "slug": "go-intro",
        "icon": "go",
        "description": "Go语言从入门到精通",
        "cover": "https://...",
        "sort_order": 1,
        "status": 1,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10
  }
}
```

### 获取系列详情

**请求**
```http
GET /api/v1/series/:id
```

**响应**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 1,
    "name": "Go语言入门",
    "slug": "go-intro",
    "sections": [...]
  }
}
```

### 创建系列 [需要认证]

**请求**
```http
POST /api/v1/series
Content-Type: application/json

{
  "name": "Go语言入门",
  "slug": "go-intro",
  "icon": "go",
  "description": "Go语言从入门到精通",
  "cover": "https://...",
  "sort_order": 1
}
```

**响应**
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": 1,
    ...
  }
}
```

### 更新系列 [需要认证]

**请求**
```http
PUT /api/v1/series/:id
Content-Type: application/json

{
  "name": "Go语言入门（修订版）",
  "description": "新描述"
}
```

**响应**
```json
{
  "code": 0,
  "message": "更新成功",
  "data": {...}
}
```

### 删除系列 [需要认证]

**请求**
```http
DELETE /api/v1/series/:id
```

**响应**
```json
{
  "code": 0,
  "message": "删除成功",
  "data": null
}
```

---

## 章节管理 (Sections)

### 获取系列章节列表

**请求**
```http
GET /api/v1/series/:id/sections?page=1&page_size=10
```

**响应**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "items": [
      {
        "id": 1,
        "series_id": 1,
        "name": "基础语法",
        "description": "Go语言基础语法",
        "sort_order": 1,
        "created_at": "...",
        "updated_at": "..."
      }
    ],
    "total": 3,
    "page": 1,
    "page_size": 10
  }
}
```

### 获取章节详情

**请求**
```http
GET /api/v1/series/sections/:id
```

**响应**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 1,
    "name": "基础语法",
    "subchapters": [...]
  }
}
```

### 创建章节 [需要认证]

**请求**
```http
POST /api/v1/series/:id/sections
Content-Type: application/json

{
  "series_id": 1,
  "name": "基础语法",
  "description": "Go语言基础语法",
  "sort_order": 1
}
```

### 更新章节 [需要认证]

**请求**
```http
PUT /api/v1/series/sections/:id
Content-Type: application/json

{
  "name": "基础语法（修订版）",
  "description": "新描述"
}
```

### 删除章节 [需要认证]

**请求**
```http
DELETE /api/v1/series/sections/:id
```

---

## 子章节管理 (Subchapters)

### 获取章节子章节列表

**请求**
```http
GET /api/v1/series/sections/:id/subchapters?page=1&page_size=10
```

### 获取子章节详情

**请求**
```http
GET /api/v1/series/subchapters/:id
```

**响应**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 1,
    "section_id": 1,
    "name": "变量与常量",
    "description": "Go语言变量和常量的定义",
    "articles": [...]
  }
}
```

### 创建子章节 [需要认证]

**请求**
```http
POST /api/v1/series/subchapters
Content-Type: application/json

{
  "section_id": 1,
  "name": "变量与常量",
  "description": "Go语言变量和常量的定义",
  "sort_order": 1
}
```

### 更新子章节 [需要认证]

**请求**
```http
PUT /api/v1/series/subchapters/:id
Content-Type: application/json

{
  "name": "变量与常量（修订版）",
  "description": "新描述"
}
```

### 删除子章节 [需要认证]

**请求**
```http
DELETE /api/v1/series/subchapters/:id
```

---

## 子章节文章关联

### 添加文章到子章节 [需要认证]

**请求**
```http
POST /api/v1/series/subchapters/:id/articles
Content-Type: application/json

{
  "article_id": 10,
  "sort_order": 1
}
```

### 从子章节移除文章 [需要认证]

**请求**
```http
DELETE /api/v1/series/subchapters/:subchapter_id/articles/:article_id
```

### 获取子章节的文章列表

**请求**
```http
GET /api/v1/series/subchapters/:id/articles
```

**响应**
```json
{
  "code": 0,
  "message": "",
  "data": [
    {
      "id": 10,
      "title": "文章标题",
      "slug": "article-slug",
      ...
    }
  ]
}
```

---

## 错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 级联删除

- 删除系列时会自动删除所有章节
- 删除章节时会自动删除所有子章节
- 删除子章节时会自动移除所有文章关联
