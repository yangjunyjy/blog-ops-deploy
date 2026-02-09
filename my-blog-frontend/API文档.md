# 博客系统 API 文档

## 基础信息

- Base URL: `/api`
- 响应格式:
  ```json
  {
    "code": 200,
    "message": "success",
    "data": {}
  }
  ```
  错误响应:
  ```json
  {
    "code": 500,
    "message": "错误信息",
    "data": null
  }
  ```

---

## 1. 文章相关接口

### 1.1 获取文章列表

**接口路径:** `GET /articles`

**请求参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | number | 否 | 页码，默认 1 |
| pageSize | number | 否 | 每页数量，默认 10 |
| categoryId | number | 否 | 分类 ID，按分类筛选 |
| tag | string | 否 | 标签名称，按标签筛选 |
| keyword | string | 否 | 搜索关键词 |
| sort | string | 否 | 排序方式：latest(最新)、hot(热门) |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "Vue 3 组合式 API 完全指南",
        "summary": "深入了解 Vue 3 的组合式 API",
        "cover": "https://example.com/cover.jpg",
        "category": {
          "id": 1,
          "name": "前端开发"
        },
        "tags": [
          {
            "name": "Vue.js",
            "articleCount": 15
          }
        ],
        "author": {
          "id": 1,
          "name": "技术博主",
          "avatar": "https://example.com/avatar.jpg",
          "bio": "专注前端技术分享"
        },
        "views": 1250,
        "likes": 86,
        "favorites": 42,
        "comments": 15,
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "pageSize": 10
  }
}
```

---

### 1.2 获取文章详情

**接口路径:** `GET /articles/:id`

**路径参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | number | 是 | 文章 ID |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "title": "Vue 3 组合式 API 完全指南",
    "summary": "深入了解 Vue 3 的组合式 API",
    "content": "# Vue 3 组合式 API\n\n详细内容...",
    "cover": "https://example.com/cover.jpg",
    "category": {
      "id": 1,
      "name": "前端开发"
    },
    "tags": [
      {
        "name": "Vue.js",
        "articleCount": 15
      }
    ],
    "author": {
      "id": 1,
      "name": "技术博主",
      "avatar": "https://example.com/avatar.jpg",
      "bio": "专注前端技术分享"
    },
    "views": 1250,
    "likes": 86,
    "favorites": 42,
    "comments": 15,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z",
    "prevArticle": {
      "id": 2,
      "title": "上一篇标题"
    },
    "nextArticle": {
      "id": 3,
      "title": "下一篇标题"
    }
  }
}
```

---

### 1.3 获取热门文章

**接口路径:** `GET /articles/hot`

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "title": "Vue 3 组合式 API 完全指南",
      "views": 1250,
      "cover": "https://example.com/cover.jpg"
    }
  ]
}
```

---

### 1.4 获取最新文章

**接口路径:** `GET /articles/latest`

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "title": "Vue 3 组合式 API 完全指南",
      "createdAt": "2024-01-01T00:00:00Z",
      "cover": "https://example.com/cover.jpg"
    }
  ]
}
```

---

### 1.5 点赞文章

**接口路径:** `POST /articles/:id/like`

**路径参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | number | 是 | 文章 ID |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "likes": 87
  }
}
```

---

### 1.6 取消点赞文章

**接口路径:** `DELETE /articles/:id/like`

**路径参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | number | 是 | 文章 ID |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "likes": 86
  }
}
```

---

### 1.7 收藏文章

**接口路径:** `POST /articles/:id/favorite`

**路径参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | number | 是 | 文章 ID |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "favorites": 43
  }
}
```

---

### 1.8 取消收藏文章

**接口路径:** `DELETE /articles/:id/favorite`

**路径参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | number | 是 | 文章 ID |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "favorites": 42
  }
}
```

---

## 2. 分类相关接口

### 2.1 获取分类列表

**接口路径:** `GET /categories`

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "前端开发",
      "articleCount": 15
    },
    {
      "id": 2,
      "name": "后端开发",
      "articleCount": 10
    }
  ]
}
```

---

### 2.2 获取分类详情

**接口路径:** `GET /categories/:id`

**路径参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | number | 是 | 分类 ID |

**请求参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | number | 否 | 页码，默认 1 |
| pageSize | number | 否 | 每页数量，默认 12 |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "category": {
      "id": 1,
      "name": "前端开发"
    },
    "articles": {
      "list": [
        {
          "id": 1,
          "title": "Vue 3 组合式 API 完全指南",
          "cover": "https://example.com/cover.jpg"
        }
      ],
      "total": 15,
      "page": 1,
      "pageSize": 12
    }
  }
}
```

---

## 3. 标签相关接口

### 3.1 获取标签列表

**接口路径:** `GET /tags`

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "name": "Vue.js",
      "articleCount": 15
    },
    {
      "name": "React",
      "articleCount": 12
    }
  ]
}
```

---

### 3.2 获取标签详情

**接口路径:** `GET /tags/:name`

**路径参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 标签名称 |

**请求参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | number | 否 | 页码，默认 1 |
| pageSize | number | 否 | 每页数量，默认 12 |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "tag": {
      "name": "Vue.js"
    },
    "articles": {
      "list": [
        {
          "id": 1,
          "title": "Vue 3 组合式 API 完全指南",
          "cover": "https://example.com/cover.jpg"
        }
      ],
      "total": 15,
      "page": 1,
      "pageSize": 12
    }
  }
}
```

---

## 4. 搜索相关接口

### 4.1 搜索文章

**接口路径:** `GET /search`

**请求参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| keyword | string | 是 | 搜索关键词 |
| page | number | 否 | 页码，默认 1 |
| pageSize | number | 否 | 每页数量，默认 12 |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "Vue 3 组合式 API 完全指南",
        "summary": "深入了解 Vue 3 的组合式 API",
        "cover": "https://example.com/cover.jpg"
      }
    ],
    "total": 5,
    "page": 1,
    "pageSize": 12
  }
}
```

---

## 5. 评论相关接口

### 5.1 获取文章评论

**接口路径:** `GET /articles/:articleId/comments`

**路径参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| articleId | number | 是 | 文章 ID |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "articleId": 1,
      "content": "这篇文章写得非常好！",
      "author": {
        "id": 1,
        "name": "张三",
        "avatar": "https://example.com/avatar.jpg"
      },
      "createdAt": "2024-01-01T00:00:00Z",
      "parentId": null,
      "replies": [
        {
          "id": 2,
          "articleId": 1,
          "content": "感谢支持！",
          "author": {
            "id": 2,
            "name": "博主",
            "avatar": "https://example.com/avatar.jpg"
          },
          "createdAt": "2024-01-01T01:00:00Z",
          "parentId": 1,
          "replies": []
        }
      ]
    }
  ]
}
```

---

### 5.2 添加评论

**接口路径:** `POST /comments`

**请求体:**
```json
{
  "articleId": 1,
  "content": "这是一条评论",
  "parentId": null
}
```

**字段说明:**
| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| articleId | number | 是 | 文章 ID |
| content | string | 是 | 评论内容 |
| parentId | number | 否 | 父评论 ID（回复评论时填写） |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 3,
    "articleId": 1,
    "content": "这是一条评论",
    "author": {
      "id": 99,
      "name": "当前用户",
      "avatar": "https://example.com/avatar.jpg"
    },
    "createdAt": "2024-01-01T00:00:00Z",
    "parentId": null,
    "replies": []
  }
}
```

---

### 5.3 删除评论

**接口路径:** `DELETE /comments/:id`

**路径参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | number | 是 | 评论 ID |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "message": "删除成功"
  }
}
```

---

## 6. 专栏系列接口

### 6.1 获取专栏列表

**接口路径:** `GET /series`

**响应示例:**
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
            },
            {
              "id": 2,
              "name": "数据结构",
              "description": "Python内置数据结构",
              "articleIds": [3, 4]
            }
          ]
        },
        {
          "id": 2,
          "name": "进阶篇",
          "description": "进阶编程技巧",
          "subchapters": [
            {
              "id": 3,
              "name": "装饰器",
              "description": "Python装饰器详解",
              "articleIds": [5, 6]
            }
          ]
        }
      ]
    }
  ]
}
```

---

### 6.2 获取专栏详情

**接口路径:** `GET /series/:id`

**路径参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | number | 是 | 系列 ID |

**响应示例:**
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
            "articleIds": [1, 2],
            "articles": [
              {
                "id": 1,
                "title": "文章标题1",
                "summary": "文章摘要",
                "cover": "https://example.com/cover.jpg",
                "views": 100,
                "createdAt": "2024-01-01T00:00:00Z"
              },
              {
                "id": 2,
                "title": "文章标题2",
                "summary": "文章摘要",
                "cover": "https://example.com/cover.jpg",
                "views": 150,
                "createdAt": "2024-01-02T00:00:00Z"
              }
            ]
          }
        ]
      }
    ]
  }
}
```

---

### 6.3 获取子章节文章

**接口路径:** `GET /series/:seriesId/chapters/:chapterId`

**路径参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| seriesId | number | 是 | 系列 ID |
| chapterId | number | 是 | 子章节 ID |

**响应示例:**
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
        "content": "完整内容...",
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

## 7. 关于页面

### 7.1 获取关于信息

**接口路径:** `GET /about`

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "name": "博客作者",
    "avatar": "https://example.com/avatar.jpg",
    "bio": "热爱技术，分享知识，持续学习",
    "email": "example@email.com",
    "github": "https://github.com/username",
    "website": "https://example.com",
    "skills": ["JavaScript", "Vue.js", "React", "Node.js", "Python", "TypeScript"],
    "description": "你好！欢迎来到我的博客。\n\n我是一名热爱技术的开发者..."
  }
}
```

---

## 8. 用户认证接口

### 8.1 用户注册

**接口路径:** `POST /auth/register`

**请求体:**
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "nickname": "测试用户"
}
```

**字段说明:**
| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名，3-50个字符，唯一 |
| email | string | 是 | 邮箱地址，唯一 |
| password | string | 是 | 密码，6-32个字符 |
| nickname | string | 否 | 昵称，可选 |

**响应示例:**
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "nickname": "测试用户",
      "avatar": null,
      "bio": null
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

### 8.2 用户登录

**接口路径:** `POST /auth/login`

**请求体:**
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**字段说明:**
| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名或邮箱 |
| password | string | 是 | 密码 |

**响应示例:**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "nickname": "测试用户",
      "avatar": "https://example.com/avatar.jpg",
      "bio": "这是我的个人简介"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

### 8.3 退出登录

**接口路径:** `POST /auth/logout`

**请求头:**
```
Authorization: Bearer <token>
```

**响应示例:**
```json
{
  "code": 200,
  "message": "退出成功",
  "data": null
}
```

---

### 8.4 刷新 Token

**接口路径:** `POST /auth/refresh`

**请求头:**
```
Authorization: Bearer <token>
```

**响应示例:**
```json
{
  "code": 200,
  "message": "刷新成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

## 9. 用户信息接口

### 9.1 获取当前用户信息

**接口路径:** `GET /user/me`

**请求头:**
```
Authorization: Bearer <token>
```

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "测试用户",
    "avatar": "https://example.com/avatar.jpg",
    "bio": "这是我的个人简介",
    "website": "https://example.com",
    "github": "https://github.com/username",
    "role": 0,
    "status": 1,
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

---

### 9.2 更新用户信息

**接口路径:** `PUT /user/me`

**请求头:**
```
Authorization: Bearer <token>
```

**请求体:**
```json
{
  "nickname": "新昵称",
  "bio": "新的个人简介",
  "website": "https://newwebsite.com",
  "github": "https://github.com/newusername"
}
```

**字段说明:**
| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| nickname | string | 否 | 昵称 |
| bio | string | 否 | 个人简介 |
| website | string | 否 | 个人网站 |
| github | string | 否 | GitHub 地址 |

**响应示例:**
```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "新昵称",
    "avatar": "https://example.com/avatar.jpg",
    "bio": "新的个人简介",
    "website": "https://newwebsite.com",
    "github": "https://github.com/newusername"
  }
}
```

---

### 9.3 更改密码

**接口路径:** `PUT /user/password`

**请求头:**
```
Authorization: Bearer <token>
```

**请求体:**
```json
{
  "oldPassword": "oldpass123",
  "newPassword": "newpass123"
}
```

**字段说明:**
| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| oldPassword | string | 是 | 旧密码 |
| newPassword | string | 是 | 新密码（6-32个字符） |

**响应示例:**
```json
{
  "code": 200,
  "message": "密码修改成功",
  "data": null
}
```

---

### 9.4 上传头像

**接口路径:** `POST /user/avatar`

**请求头:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**请求体:**
```
avatar: <文件>
```

**响应示例:**
```json
{
  "code": 200,
  "message": "头像上传成功",
  "data": {
    "url": "https://example.com/avatars/user1_new.jpg"
  }
}
```

---

### 9.5 获取用户收藏列表

**接口路径:** `GET /user/favorites`

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | number | 否 | 页码，默认 1 |
| pageSize | number | 否 | 每页数量，默认 12 |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "Vue 3 组合式 API 完全指南",
        "summary": "深入了解 Vue 3 的组合式 API",
        "cover": "https://example.com/cover.jpg",
        "author": {
          "id": 1,
          "name": "技术博主"
        },
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 20,
    "page": 1,
    "pageSize": 12
  }
}
```

---

### 9.6 获取用户点赞列表

**接口路径:** `GET /user/likes`

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | number | 否 | 页码，默认 1 |
| pageSize | number | 否 | 每页数量，默认 12 |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "Vue 3 组合式 API 完全指南",
        "summary": "深入了解 Vue 3 的组合式 API",
        "cover": "https://example.com/cover.jpg",
        "author": {
          "id": 1,
          "name": "技术博主"
        },
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 15,
    "page": 1,
    "pageSize": 12
  }
}
```

---

### 9.7 获取用户评论列表

**接口路径:** `GET /user/comments`

**请求头:**
```
Authorization: Bearer <token>
```

**请求参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | number | 否 | 页码，默认 1 |
| pageSize | number | 否 | 每页数量，默认 12 |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "content": "这篇文章写得非常好！",
        "article": {
          "id": 1,
          "title": "Vue 3 组合式 API 完全指南"
        },
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "pageSize": 12
  }
}
```

---

### 9.8 检查文章是否已收藏

**接口路径:** `GET /user/favorite/check/:articleId`

**请求头:**
```
Authorization: Bearer <token>
```

**路径参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| articleId | number | 是 | 文章 ID |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "isFavorite": true
  }
}
```

---

### 9.9 检查文章是否已点赞

**接口路径:** `GET /user/like/check/:articleId`

**请求头:**
```
Authorization: Bearer <token>
```

**路径参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| articleId | number | 是 | 文章 ID |

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "isLiked": true
  }
}
```

---

## 10. 管理员接口

### 10.1 发布文章

**接口路径:** `POST /admin/articles`

**请求头:**
```
Authorization: Bearer <token>
```

**请求体:**
```json
{
  "title": "Vue 3 组合式 API 完全指南",
  "summary": "深入了解 Vue 3 的组合式 API",
  "content": "# Vue 3 组合式 API\n\n详细内容...",
  "cover": "https://example.com/cover.jpg",
  "categoryId": 1,
  "tags": ["Vue.js", "JavaScript"],
  "status": 1,
  "publishedAt": "2024-01-01T00:00:00Z"
}
```

**字段说明:**
| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| title | string | 是 | 文章标题 |
| summary | string | 是 | 文章摘要 |
| content | string | 是 | 文章内容（Markdown） |
| cover | string | 否 | 封面图片 URL |
| categoryId | number | 否 | 分类 ID |
| tags | string[] | 否 | 标签数组 |
| status | number | 否 | 状态：0-草稿，1-已发布，默认 1 |
| publishedAt | string | 否 | 发布时间 |

**响应示例:**
```json
{
  "code": 200,
  "message": "文章发布成功",
  "data": {
    "id": 1,
    "title": "Vue 3 组合式 API 完全指南",
    "slug": "vue-3-composition-api-guide",
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

---

### 10.2 更新文章

**接口路径:** `PUT /admin/articles/:id`

**请求头:**
```
Authorization: Bearer <token>
```

**路径参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | number | 是 | 文章 ID |

**请求体:** 同发布文章

**响应示例:**
```json
{
  "code": 200,
  "message": "文章更新成功",
  "data": {
    "id": 1,
    "title": "Vue 3 组合式 API 完全指南",
    "updatedAt": "2024-01-02T00:00:00Z"
  }
}
```

---

### 10.3 删除文章

**接口路径:** `DELETE /admin/articles/:id`

**请求头:**
```
Authorization: Bearer <token>
```

**路径参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | number | 是 | 文章 ID |

**响应示例:**
```json
{
  "code": 200,
  "message": "文章删除成功",
  "data": null
}
```

---

### 10.4 创建分类

**接口路径:** `POST /admin/categories`

**请求头:**
```
Authorization: Bearer <token>
```

**请求体:**
```json
{
  "name": "前端开发",
  "slug": "frontend",
  "description": "前端开发相关文章",
  "icon": "💻"
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "分类创建成功",
  "data": {
    "id": 1,
    "name": "前端开发"
  }
}
```

---

### 10.5 创建标签

**接口路径:** `POST /admin/tags`

**请求头:**
```
Authorization: Bearer <token>
```

**请求体:**
```json
{
  "name": "Vue.js",
  "slug": "vuejs",
  "description": "Vue框架相关"
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "标签创建成功",
  "data": {
    "id": 1,
    "name": "Vue.js"
  }
}
```

---

### 10.6 审核评论

**接口路径:** `PUT /admin/comments/:id/status`

**请求头:**
```
Authorization: Bearer <token>
```

**路径参数:**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | number | 是 | 评论 ID |

**请求体:**
```json
{
  "status": 1
}
```

**字段说明:**
| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| status | number | 是 | 状态：1-通过，2-拒绝 |

**响应示例:**
```json
{
  "code": 200,
  "message": "评论审核完成",
  "data": null
}
```

---

## 数据结构说明

### 文章 (Article)
```typescript
interface Article {
  id: number
  title: string
  summary: string
  content: string
  cover: string
  category: Category
  tags: Tag[]
  author: User
  views: number
  likes: number
  favorites: number
  comments: number
  createdAt: string
  updatedAt: string
  prevArticle?: Article
  nextArticle?: Article
}
```

### 分类 (Category)
```typescript
interface Category {
  id: number
  name: string
  articleCount: number
}
```

### 标签 (Tag)
```typescript
interface Tag {
  name: string
  articleCount: number
}
```

### 用户 (User)
```typescript
interface User {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  bio: string
  website: string
  github: string
  role: number
  status: number
  createdAt: string
}
```

### 评论 (Comment)
```typescript
interface Comment {
  id: number
  articleId: number
  content: string
  author: User
  createdAt: string
  parentId: number | null
  replies: Comment[]
}
```

### 系列 (Series)
```typescript
interface Series {
  id: number
  name: string
  icon: string
  description: string
  sections: Section[]
}
```

### 章节标题 (Section)
```typescript
interface Section {
  id: number
  name: string
  description: string
  subchapters: Subchapter[]
}
```

### 子章节 (Subchapter)
```typescript
interface Subchapter {
  id: number
  name: string
  description: string
  articleIds: number[]
  articles?: Article[]
}
```

### 用户认证响应 (AuthResponse)
```typescript
interface AuthResponse {
  user: User
  token: string
}
```
