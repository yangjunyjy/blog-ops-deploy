# 专栏数据结构文档

## 📚 概述

本文档描述博客系统专栏模块的完整数据结构,采用四层层级结构:Series → Section → Subchapter → Articles。

## 🎯 数据层级

```
Series (系列)
  └── Section (章节标题/大章节)
       └── Subchapter (子章节/目录卡片)
            └── Articles (文章)
```

## 📊 数据结构定义

### 1. Series (系列)

系列的顶级容器,包含整个专题的信息。

```typescript
interface Series {
  id: number              // 系列ID,主键
  name: string           // 系列名称,如"Python基础知识"
  icon: string           // 系列图标(emoji),如"🐍"
  description: string    // 系列描述
  sections: Section[]     // 章节标题列表
}
```

**示例:**
```json
{
  "id": 1,
  "name": "Python基础知识",
  "icon": "🐍",
  "description": "从入门到精通的Python学习路线",
  "sections": [
    {
      "id": 1,
      "name": "基础篇",
      "description": "Python入门必备知识",
      "subchapters": [...]
    },
    {
      "id": 2,
      "name": "进阶篇",
      "description": "进阶编程技巧",
      "subchapters": [...]
    }
  ]
}
```

---

### 2. Section (章节标题)

系列下的大章节,用于组织和归类子章节。

```typescript
interface Section {
  id: number                  // 章节标题ID
  name: string               // 章节标题名称,如"基础篇"、"进阶篇"
  description: string        // 章节标题描述
  subchapters: Subchapter[]  // 子章节列表
}
```

**示例:**
```json
{
  "id": 1,
  "name": "基础篇",
  "description": "Python入门必备知识",
  "subchapters": [
    {
      "id": 1,
      "name": "Python3面向对象",
      "description": "面向对象编程详解",
      "articleIds": [1, 2]
    },
    {
      "id": 2,
      "name": "数据结构",
      "description": "Python内置数据结构",
      "articleIds": [3, 4]
    }
  ]
}
```

---

### 3. Subchapter (子章节)

具体的目录卡片,点击后显示该目录下的所有文章。

```typescript
interface Subchapter {
  id: number          // 子章节ID
  name: string       // 子章节名称,如"Python3面向对象"、"数据结构"
  description: string // 子章节描述
  articleIds: number[] // 文章ID数组
  articles?: Article[] // 文章详情列表(可选,按需加载)
}
```

**示例:**
```json
{
  "id": 1,
  "name": "Python3面向对象",
  "description": "面向对象编程详解",
  "articleIds": [1, 2],
  "articles": [
    {
      "id": 1,
      "title": "Python面向对象编程详解",
      "summary": "详细讲解Python的面向对象特性",
      // ... 文章其他字段
    }
  ]
}
```

---

### 4. Article (文章)

文章的完整信息。

```typescript
interface Article {
  id: number              // 文章ID
  title: string          // 文章标题
  summary: string         // 文章摘要
  content: string         // 文章内容(Markdown)
  cover: string          // 封面图片URL
  category: {            // 分类信息
    id: number
    name: string
  }
  tags: {                // 标签数组
    name: string
    articleCount: number
  }[]
  author: {              // 作者信息
    id: number
    name: string
    avatar: string
    bio: string
  }
  views: number          // 浏览次数
  likes: number          // 点赞数
  favorites: number      // 收藏数
  comments: number      // 评论数
  createdAt: string      // 创建时间(ISO 8601)
  updatedAt: string      // 更新时间(ISO 8601)
}
```

**示例:**
```json
{
  "id": 1,
  "title": "Python面向对象编程详解",
  "summary": "详细讲解Python的面向对象特性",
  "content": "# Python面向对象编程\n\n详细内容...",
  "cover": "https://example.com/cover.jpg",
  "category": {
    "id": 2,
    "name": "后端开发"
  },
  "tags": [
    {
      "name": "Python",
      "articleCount": 8
    }
  ],
  "author": {
    "id": 1,
    "name": "技术博主",
    "avatar": "https://example.com/avatar.jpg",
    "bio": "专注前端技术分享"
  },
  "views": 100,
  "likes": 20,
  "favorites": 10,
  "comments": 5,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

---

## 🔌 API 接口

### 1. 获取系列列表

```http
GET /api/series
```

**响应:**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "Python基础知识",
      "icon": "🐍",
      "description": "从入门到精通的Python学习路线",
      "sections": [...]
    }
  ]
}
```

---

### 2. 获取系列详情

```http
GET /api/series/:id
```

**路径参数:**
- `id`: 系列ID

**响应:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "Python基础知识",
    "icon": "🐍",
    "description": "从入门到精通的Python学习路线",
    "sections": [
      {
        "id": 1,
        "name": "基础篇",
        "description": "Python入门必备知识",
        "subchapters": [
          {
            "id": 1,
            "name": "Python3面向对象",
            "description": "面向对象编程详解",
            "articleIds": [1, 2]
          }
        ]
      }
    ]
  }
}
```

---

### 3. 获取子章节文章列表

```http
GET /api/series/:seriesId/chapters/:chapterId
```

**路径参数:**
- `seriesId`: 系列ID
- `chapterId`: 子章节ID

**响应:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "Python3面向对象",
    "description": "面向对象编程详解",
    "articleIds": [1, 2],
    "articles": [
      {
        "id": 1,
        "title": "Python面向对象编程详解",
        "summary": "详细讲解Python的面向对象特性",
        "content": "# Python面向对象...",
        // ... 其他字段
      }
    ],
    "series": {
      "id": 1,
      "name": "Python基础知识",
      "icon": "🐍"
    }
  }
}
```

---

## 🗄️ 数据库表结构

### 1. series (系列表)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGINT | 系列ID,主键 |
| name | VARCHAR(100) | 系列名称 |
| slug | VARCHAR(100) | 系列别名,唯一 |
| icon | VARCHAR(50) | 系列图标(emoji) |
| description | TEXT | 系列描述 |
| cover | VARCHAR(500) | 系列封面图 |
| sort_order | INT | 排序顺序 |
| status | TINYINT | 状态:0-禁用,1-启用 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

---

### 2. series_sections (章节标题表)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGINT | 章节标题ID,主键 |
| series_id | BIGINT | 所属系列ID,外键 |
| name | VARCHAR(100) | 章节标题名称 |
| description | TEXT | 章节标题描述 |
| sort_order | INT | 排序顺序 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

**外键:** `series_id` → `series(id)`

---

### 3. series_subchapters (子章节表)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGINT | 子章节ID,主键 |
| section_id | BIGINT | 所属章节标题ID,外键 |
| name | VARCHAR(100) | 子章节名称 |
| description | TEXT | 子章节描述 |
| sort_order | INT | 排序顺序 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

**外键:** `section_id` → `series_sections(id)`

---

### 4. subchapter_articles (子章节文章关联表)

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGINT | 主键ID |
| subchapter_id | BIGINT | 子章节ID,外键 |
| article_id | BIGINT | 文章ID,外键 |
| sort_order | INT | 在子章节中的排序 |
| created_at | TIMESTAMP | 创建时间 |

**外键:**
- `subchapter_id` → `series_subchapters(id)`
- `article_id` → `articles(id)`

---

## 🔄 数据查询示例

### 1. 获取系列完整结构

```sql
SELECT
  s.id as series_id,
  s.name as series_name,
  s.icon as series_icon,
  sec.id as section_id,
  sec.name as section_name,
  sub.id as subchapter_id,
  sub.name as subchapter_name,
  a.id as article_id,
  a.title as article_title
FROM series s
LEFT JOIN series_sections sec ON s.id = sec.series_id
LEFT JOIN series_subchapters sub ON sec.id = sub.section_id
LEFT JOIN subchapter_articles sa ON sub.id = sa.subchapter_id
LEFT JOIN articles a ON sa.article_id = a.id
WHERE s.id = ?
ORDER BY sec.sort_order, sub.sort_order, sa.sort_order;
```

---

### 2. 获取子章节的文章列表

```sql
SELECT
  sub.*,
  a.id as article_id,
  a.title,
  a.summary,
  a.cover,
  a.views,
  a.created_at
FROM series_subchapters sub
JOIN subchapter_articles sa ON sub.id = sa.subchapter_id
JOIN articles a ON sa.article_id = a.id
WHERE sub.id = ?
ORDER BY sa.sort_order;
```

---

## 💡 使用场景

### 场景1: 用户浏览系列

1. 用户点击"Python基础知识"系列
2. 前端调用 `GET /api/series/1` 获取系列详情
3. 页面显示章节标题(Section)列表
4. 每个章节标题下显示子章节(Subchapter)卡片

### 场景2: 用户点击子章节

1. 用户点击"Python3面向对象"子章节卡片
2. 前端调用 `GET /api/series/1/chapters/1` 获取子章节文章
3. 在新窗口打开文章列表页面
4. 左侧显示该子章节的所有文章目录
5. 右侧显示选中文章的详细内容

---

## 🎨 前端数据流

```
Series.vue (系列列表页)
  └─> 点击系列
      └─> SeriesDetail.vue (系列详情页)
          ├─> 显示 sections (章节标题)
          │   └─> 显示 subchapters (子章节卡片)
          │       └─> 点击子章节
          │           └─> ChapterDetail.vue (子章节文章页)
          │               ├─> 左侧: 文章目录列表
          │               └─> 右侧: 文章详情
```

---

## 🔧 开发注意事项

1. **ID关联**: 确保 `articleIds` 与数据库中实际的文章ID一致
2. **层级遍历**: 在系列中查找子章节时,需要遍历 `sections` 数组
3. **懒加载**: 文章详情可以按需加载,不在系列列表中返回完整文章内容
4. **排序**: 每层结构都应该支持自定义排序(`sort_order`)
5. **性能**: 考虑使用缓存机制,避免重复查询数据库

---

## 📝 Mock数据示例

完整Mock数据请参考: `src/mock/data.js`

```javascript
const series = [
  {
    id: 1,
    name: 'Python基础知识',
    icon: '🐍',
    description: '从入门到精通的Python学习路线',
    sections: [
      {
        id: 1,
        name: '基础篇',
        description: 'Python入门必备知识',
        subchapters: [
          {
            id: 1,
            name: 'Python3面向对象',
            description: '面向对象编程详解',
            articleIds: [1, 2]
          },
          {
            id: 2,
            name: '数据结构',
            description: 'Python内置数据结构',
            articleIds: [3, 4]
          }
        ]
      }
    ]
  }
]
```

---

## 🚀 未来扩展

- [ ] 支持用户自定义学习路径
- [ ] 添加学习进度追踪
- [ ] 支持子章节间的依赖关系
- [ ] 添加学习笔记功能
- [ ] 支持多语言系列
