# 博客管理系统前端 API 模块文档

## 概述

本文档描述了博客管理系统前端 API 包的所有接口和模块。

## API 模块列表

### 1. 认证模块 (auth.js)

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| login | POST | /auth/login | 用户登录 |
| logout | POST | /auth/logout | 用户登出 |
| getUserInfo | GET | /user/info | 获取当前用户信息 |
| refreshToken | POST | /auth/refresh | 刷新访问令牌 |

---

### 2. 文章管理 (article.js)

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getArticleList | GET | /admin/articles | 获取文章列表 |
| getArticleDetail | GET | /admin/articles/:id | 获取文章详情 |
| createArticle | POST | /admin/articles | 创建文章 |
| updateArticle | PUT | /admin/articles/:id | 更新文章 |
| deleteArticle | DELETE | /admin/articles/:id | 删除文章 |
| batchDeleteArticles | DELETE | /admin/articles/batch | 批量删除文章 |
| publishArticle | PUT | /admin/articles/:id/publish | 发布文章 |
| withdrawArticle | PUT | /admin/articles/:id/withdraw | 撤回文章 |
| batchPublishArticles | PUT | /admin/articles/batch/publish | 批量发布文章 |
| batchWithdrawArticles | PUT | /admin/articles/batch/withdraw | 批量撤回文章 |
| updateArticleStatus | PUT | /admin/articles/:id/status | 更新文章状态 |
| batchUpdateArticleStatus | PUT | /admin/articles/batch/status | 批量更新文章状态 |

---

### 3. 分类管理 (category.js)

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getCategories | GET | /admin/categories | 获取所有分类（不分页） |
| getCategoryList | GET | /admin/categories | 获取分类列表（分页） |
| getCategoryTree | GET | /admin/categories/tree | 获取分类树 |
| getCategoryDetail | GET | /admin/categories/:id | 获取分类详情 |
| createCategory | POST | /admin/categories | 创建分类 |
| updateCategory | PUT | /admin/categories/:id | 更新分类 |
| deleteCategory | DELETE | /admin/categories/:id | 删除分类 |
| batchDeleteCategories | DELETE | /admin/categories/batch | 批量删除分类 |

---

### 4. 标签管理 (tag.js)

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getTags | GET | /admin/tags | 获取所有标签（不分页） |
| getTagList | GET | /admin/tags | 获取标签列表（分页） |
| getTagDetail | GET | /admin/tags/:id | 获取标签详情 |
| createTag | POST | /admin/tags | 创建标签 |
| updateTag | PUT | /admin/tags/:id | 更新标签 |
| deleteTag | DELETE | /admin/tags/:id | 删除标签 |
| batchDeleteTags | DELETE | /admin/tags/batch | 批量删除标签 |

---

### 5. 评论管理 (comment.js)

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getCommentList | GET | /admin/comments | 获取评论列表 |
| getCommentDetail | GET | /admin/comments/:id | 获取评论详情 |
| approveComment | PUT | /admin/comments/:id/approve | 审核通过评论 |
| rejectComment | PUT | /admin/comments/:id/reject | 审核拒绝评论 |
| deleteComment | DELETE | /admin/comments/:id | 删除评论 |
| batchDeleteComments | DELETE | /admin/comments/batch | 批量删除评论 |
| batchApproveComments | PUT | /admin/comments/batch/approve | 批量审核通过评论 |
| batchRejectComments | PUT | /admin/comments/batch/reject | 批量审核拒绝评论 |

---

### 6. 用户管理 (user.js)

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getUserList | GET | /admin/users | 获取用户列表 |
| getUserDetail | GET | /admin/users/:id | 获取用户详情 |
| createUser | POST | /admin/users | 创建用户 |
| updateUser | PUT | /admin/users/:id | 更新用户 |
| deleteUser | DELETE | /admin/users/:id | 删除用户 |
| batchDeleteUsers | DELETE | /admin/users/batch | 批量删除用户 |
| resetPassword | POST | /admin/users/:id/reset-password | 重置用户密码 |
| updateUserStatus | PUT | /admin/users/:id/status | 更新用户状态 |
| assignUserRoles | POST | /admin/users/:id/assign-roles | 分配角色给用户 |
| getUserRoles | GET | /admin/users/:id/roles | 获取用户角色 |
| batchUpdateUserStatus | PUT | /admin/users/batch/status | 批量更新用户状态 |

---

### 7. 角色管理 (role.js) 🆕

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getRoleList | GET | /admin/roles | 获取角色列表 |
| getRoleDetail | GET | /admin/roles/:id | 获取角色详情 |
| createRole | POST | /admin/roles | 创建角色 |
| updateRole | PUT | /admin/roles/:id | 更新角色 |
| deleteRole | DELETE | /admin/roles/:id | 删除角色 |
| assignRolePermissions | POST | /admin/roles/:id/permissions | 分配权限给角色 |
| getRolePermissions | GET | /admin/roles/:id/permissions | 获取角色权限 |
| batchDeleteRoles | DELETE | /admin/roles/batch | 批量删除角色 |

---

### 8. 菜单管理 (menu.js) 🆕

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getMenuList | GET | /admin/menus | 获取菜单列表 |
| getMenuTree | GET | /admin/menus/tree | 获取菜单树 |
| getMenuDetail | GET | /admin/menus/:id | 获取菜单详情 |
| createMenu | POST | /admin/menus | 创建菜单 |
| updateMenu | PUT | /admin/menus/:id | 更新菜单 |
| deleteMenu | DELETE | /admin/menus/:id | 删除菜单 |
| batchDeleteMenus | DELETE | /admin/menus/batch | 批量删除菜单 |

---

### 9. 审核管理 (audit.js) 🆕

#### 文章审核
| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getPendingArticles | GET | /admin/articles/pending | 获取待审核文章 |
| getApprovedArticles | GET | /admin/articles/approved | 获取已审核通过文章 |
| getRejectedArticles | GET | /admin/articles/rejected | 获取已审核拒绝文章 |
| approveArticle | PUT | /admin/articles/:id/approve | 审核通过文章 |
| rejectArticle | PUT | /admin/articles/:id/reject | 审核拒绝文章 |
| getArticleAuditLogs | GET | /admin/articles/:id/audit-log | 获取文章审核日志 |

#### 评论审核
| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getPendingComments | GET | /admin/comments/pending | 获取待审核评论 |
| getApprovedComments | GET | /admin/comments/approved | 获取已审核通过评论 |
| getRejectedComments | GET | /admin/comments/rejected | 获取已审核拒绝评论 |

#### 用户审核
| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getPendingUsers | GET | /admin/users/pending | 获取待审核用户 |
| getApprovedUsers | GET | /admin/users/approved | 获取已审核通过用户 |
| getRejectedUsers | GET | /admin/users/rejected | 获取已审核拒绝用户 |
| approveUser | PUT | /admin/users/:id/approve | 审核通过用户 |
| rejectUser | PUT | /admin/users/:id/reject | 审核拒绝用户 |

---

### 10. 专题系列管理 (series.js)

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getSeriesList | GET | /admin/series | 获取系列列表 |
| getSeriesDetail | GET | /admin/series/:id | 获取系列详情 |
| createSeries | POST | /admin/series | 创建系列 |
| updateSeries | PUT | /admin/series/:id | 更新系列 |
| deleteSeries | DELETE | /admin/series/:id | 删除系列 |
| createSection | POST | /admin/series/:id/sections | 创建章节 |
| updateSection | PUT | /admin/series/:id/sections/:sectionId | 更新章节 |
| deleteSection | DELETE | /admin/series/:id/sections/:sectionId | 删除章节 |
| createSubchapter | POST | /admin/series/:id/sections/:sectionId/subchapters | 创建子章节 |
| updateSubchapter | PUT | /admin/series/:id/sections/:sectionId/subchapters/:subchapterId | 更新子章节 |
| deleteSubchapter | DELETE | /admin/series/:id/sections/:sectionId/subchapters/:subchapterId | 删除子章节 |
| getSubchapterArticles | GET | /admin/series/:id/sections/:sectionId/subchapters/:subchapterId/articles | 获取子章节文章 |

---

### 11. 数据统计 (statistics.js)

#### 仪表盘统计
| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getDashboardStats | GET | /admin/statistics/dashboard | 获取仪表盘统计数据 |

#### 文章统计
| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getArticleStats | GET | /admin/statistics/articles | 获取文章统计数据 |
| getArticleTrend | GET | /admin/statistics/articles/trend | 获取文章发布趋势 |
| getHotArticles | GET | /admin/statistics/articles/hot | 获取热门文章 |
| getCategoryStats | GET | /admin/statistics/categories | 获取分类统计 |
| getTagStats | GET | /admin/statistics/tags | 获取标签统计 |
| getHotTags | GET | /admin/statistics/tags/hot | 获取热门标签 |

#### 浏览统计
| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getViewTrend | GET | /admin/statistics/views/trend | 获取浏览量趋势 |

#### 用户统计
| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getUserGrowth | GET | /admin/statistics/users/growth | 获取用户增长统计 |
| getActiveUsers | GET | /admin/statistics/users/active | 获取活跃用户统计 |
| getUserDistribution | GET | /admin/statistics/users/distribution | 获取用户分布统计 |
| getOnlineUsers | GET | /admin/statistics/users/online | 获取在线用户统计 |

---

### 12. 系统设置 (settings.js) 🆕

#### 系统设置
| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getSystemSettings | GET | /admin/settings | 获取系统设置 |
| updateSystemSettings | PUT | /admin/settings | 更新系统设置 |
| resetSystemSettings | POST | /admin/settings/reset | 重置系统设置 |
| getSiteSettings | GET | /admin/settings/site | 获取网站信息设置 |
| updateSiteSettings | PUT | /admin/settings/site | 更新网站信息设置 |
| getSEOSettings | GET | /admin/settings/seo | 获取SEO设置 |
| updateSEOSettings | PUT | /admin/settings/seo | 更新SEO设置 |
| getSystemConfig | GET | /admin/settings/system | 获取系统配置 |
| updateSystemConfig | PUT | /admin/settings/system | 更新系统配置 |
| getEmailSettings | GET | /admin/settings/email | 获取邮件设置 |
| updateEmailSettings | PUT | /admin/settings/email | 更新邮件设置 |
| sendTestEmail | POST | /admin/settings/test-email | 发送测试邮件 |
| getStorageSettings | GET | /admin/settings/storage | 获取存储设置 |
| updateStorageSettings | PUT | /admin/settings/storage | 更新存储设置 |

#### 备份恢复
| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| createBackup | POST | /admin/backup/create | 创建备份 |
| getBackupList | GET | /admin/backup/list | 获取备份列表 |
| restoreBackup | POST | /admin/backup/:id/restore | 恢复备份 |
| deleteBackup | DELETE | /admin/backup/:id | 删除备份 |
| downloadBackup | GET | /admin/backup/:id/download | 下载备份 |

---

### 13. 文件上传 (upload.js) 🆕

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| uploadImage | POST | /admin/upload/image | 上传图片 |
| uploadFile | POST | /admin/upload/file | 上传文件 |
| getUploadList | GET | /admin/upload/list | 获取文件上传列表 |
| deleteUpload | DELETE | /admin/upload/:id | 删除文件 |
| batchDeleteUploads | DELETE | /admin/upload/batch | 批量删除文件 |

---

### 14. 权限管理 (permission.js) 🆕

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| getUserPermissions | GET | /admin/user/permissions | 获取当前用户权限列表 |
| getUserMenus | GET | /admin/user/menus | 获取当前用户菜单列表 |
| hasPermission | - | - | 检查权限（前端本地验证） |
| hasAllPermissions | - | - | 检查是否拥有所有权限 |
| hasAnyPermission | - | - | 检查是否拥有任一权限 |

---

## API 响应格式

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

---

## Mock 数据切换

所有 API 模块都支持 Mock 数据模式，通过设置 `USE_MOCK` 常量来切换：

```javascript
const USE_MOCK = true  // true 使用 Mock 数据，false 使用真实 API
```

---

## 使用示例

### 基本使用
```javascript
import { getUserList, createUser } from '@/api/user'

// 获取用户列表
const res = await getUserList({ page: 1, page_size: 10 })

// 创建用户
const user = await createUser({
  username: 'test',
  email: 'test@example.com',
  password: '123456'
})
```

### 批量操作
```javascript
import { batchDeleteUsers, batchUpdateUserStatus } from '@/api/user'

// 批量删除用户
await batchDeleteUsers([1, 2, 3])

// 批量更新用户状态
await batchUpdateUserStatus([1, 2, 3], 1)
```

### 文件上传
```javascript
import { uploadImage } from '@/api/upload'

const formData = new FormData()
formData.append('file', file)
formData.append('type', 'cover')

const res = await uploadImage(formData)
```

### 权限检查
```javascript
import { hasPermission, hasAllPermissions, hasAnyPermission } from '@/api/permission'

// 检查单个权限
if (hasPermission('content:article:create')) {
  // 有权限
}

// 检查是否拥有所有权限
if (hasAllPermissions(['content:article:create', 'content:article:update'])) {
  // 有所有权限
}

// 检查是否拥有任一权限
if (hasAnyPermission(['content:article:create', 'content:article:delete'])) {
  // 至少有一个权限
}
```

---

## 状态码说明

### 用户状态
- `0` - 待审核
- `1` - 正常/启用
- `2` - 禁用

### 文章状态
- `0` - 草稿
- `1` - 已发布

### 评论状态
- `0` - 待审核
- `1` - 已通过
- `2` - 已拒绝

### 角色状态
- `0` - 禁用
- `1` - 启用

### 菜单类型
- `0` - 目录
- `1` - 菜单
- `2` - 按钮

---

## 总结

本 API 包包含 14 个模块，共 100+ 个接口，涵盖：

- ✅ 认证授权
- ✅ 内容管理（文章、分类、标签、系列）
- ✅ 评论管理
- ✅ 用户管理
- ✅ 角色管理 🆕
- ✅ 菜单管理 🆕
- ✅ 审核管理 🆕
- ✅ 数据统计 🆕
- ✅ 系统设置 🆕
- ✅ 文件上传 🆕
- ✅ 权限管理 🆕

所有接口都支持 Mock 数据模式，方便前端开发和测试。
