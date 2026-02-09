# 角色管理 API 接口设计文档

## 概述

角色管理模块提供完整的角色 CRUD 操作，包括角色查询、创建、更新、删除以及菜单权限分配功能。

## 数据模型

### SysRole 角色模型

| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | uint64 | 是 | 角色ID（主键，自增） |
| roleCode | string | 是 | 角色编码，唯一 |
| roleName | string | 是 | 角色名称 |
| roleDesc | string | 否 | 角色描述 |
| dataScope | int8 | 否 | 数据权限范围：1-全部数据 2-本部门及以下 3-本部门 4-仅本人 |
| status | int8 | 否 | 状态：0-禁用 1-启用 |
| sort | int | 否 | 排序值 |
| createdAt | time.Time | 是 | 创建时间 |
| updatedAt | time.Time | 是 | 更新时间 |
| deletedAt | gorm.DeletedAt | 否 | 软删除时间 |

## API 接口列表

### 1. 获取角色列表（分页）

**接口描述：** 分页查询角色列表，支持关键字搜索、状态筛选、排序

**请求方式：** `GET /api/v1/roles`

**请求参数（Query Params）：**

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 10 | 每页数量 |
| keyword | string | 否 | - | 搜索关键字（角色名称或编码） |
| status | int | 否 | - | 状态筛选：0-禁用 1-启用 |
| sort | string | 否 | id | 排序字段 |
| order | string | 否 | desc | 排序方向：asc-升序 desc-降序 |

**请求示例：**
```http
GET /api/v1/roles?page=1&page_size=10&keyword=admin&status=1&sort=id&order=desc
```

**响应示例：**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "list": [
      {
        "id": 1,
        "roleCode": "super_admin",
        "roleName": "超级管理员",
        "roleDesc": "系统超级管理员，拥有所有权限",
        "dataScope": 1,
        "status": 1,
        "sort": 0,
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 10
  }
}
```

---

### 2. 获取所有角色列表（不分页）

**接口描述：** 获取所有启用的角色列表，用于下拉选择等场景

**请求方式：** `GET /api/v1/roles/all`

**请求参数：** 无

**响应示例：**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": 1,
      "roleCode": "super_admin",
      "roleName": "超级管理员",
      "roleDesc": "系统超级管理员，拥有所有权限",
      "dataScope": 1,
      "status": 1,
      "sort": 0
    },
    {
      "id": 2,
      "roleCode": "admin",
      "roleName": "管理员",
      "roleDesc": "普通管理员",
      "dataScope": 2,
      "status": 1,
      "sort": 1
    }
  ]
}
```

---

### 3. 获取角色详情

**接口描述：** 根据角色ID获取角色详细信息

**请求方式：** `GET /api/v1/roles/{id}`

**请求参数（Path Params）：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | uint64 | 是 | 角色ID |

**请求示例：**
```http
GET /api/v1/roles/1
```

**响应示例：**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": 1,
    "roleCode": "super_admin",
    "roleName": "超级管理员",
    "roleDesc": "系统超级管理员，拥有所有权限",
    "dataScope": 1,
    "status": 1,
    "sort": 0,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z",
    "menus": [
      {
        "id": 1,
        "title": "系统管理",
        "path": "/system",
        "icon": "Setting",
        "type": "directory"
      }
    ],
    "users": [
      {
        "id": 1,
        "username": "admin",
        "nickname": "管理员"
      }
    ]
  }
}
```

---

### 4. 创建角色

**接口描述：** 创建新的角色

**请求方式：** `POST /api/v1/roles`

**请求参数（Request Body）：**

| 参数名 | 类型 | 必填 | 校验规则 | 说明 |
|--------|------|------|----------|------|
| roleCode | string | 是 | required,min=2,max=50 | 角色编码，唯一 |
| roleName | string | 是 | required,min=2,max=50 | 角色名称 |
| roleDesc | string | 否 | max=200 | 角色描述 |
| dataScope | int8 | 否 | oneof=1 2 3 4, default=1 | 数据权限范围 |
| status | int8 | 否 | oneof=0 1, default=1 | 状态 |
| sort | int | 否 | default=0 | 排序值 |
| menuIds | []uint64 | 否 | - | 关联的菜单ID列表 |

**请求示例：**
```json
{
  "roleCode": "editor",
  "roleName": "编辑员",
  "roleDesc": "内容编辑员，只能管理文章和评论",
  "dataScope": 3,
  "status": 1,
  "sort": 2,
  "menuIds": [10, 11, 12, 20, 21]
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": 3,
    "roleCode": "editor",
    "roleName": "编辑员",
    "roleDesc": "内容编辑员，只能管理文章和评论",
    "dataScope": 3,
    "status": 1,
    "sort": 2,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

---

### 5. 更新角色

**接口描述：** 更新角色信息

**请求方式：** `PUT /api/v1/roles/{id}`

**请求参数（Path Params）：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | uint64 | 是 | 角色ID |

**请求参数（Request Body）：**

| 参数名 | 类型 | 必填 | 校验规则 | 说明 |
|--------|------|------|----------|------|
| roleCode | string | 否 | min=2,max=50 | 角色编码 |
| roleName | string | 否 | min=2,max=50 | 角色名称 |
| roleDesc | string | 否 | max=200 | 角色描述 |
| dataScope | int8 | 否 | oneof=1 2 3 4 | 数据权限范围 |
| status | int8 | 否 | oneof=0 1 | 状态 |
| sort | int | 否 | - | 排序值 |
| menuIds | []uint64 | 否 | - | 关联的菜单ID列表（会覆盖原有菜单） |

**请求示例：**
```json
{
  "roleName": "高级编辑员",
  "roleDesc": "高级内容编辑员",
  "dataScope": 2,
  "status": 1,
  "sort": 2,
  "menuIds": [10, 11, 12, 20, 21, 30, 31]
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "id": 3,
    "roleCode": "editor",
    "roleName": "高级编辑员",
    "roleDesc": "高级内容编辑员",
    "dataScope": 2,
    "status": 1,
    "sort": 2,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T10:00:00Z"
  }
}
```

---

### 6. 删除角色

**接口描述：** 删除指定角色（软删除）

**请求方式：** `DELETE /api/v1/roles/{id}`

**请求参数（Path Params）：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | uint64 | 是 | 角色ID |

**请求示例：**
```http
DELETE /api/v1/roles/3
```

**响应示例：**
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

**业务规则：**
- 超级管理员角色（ID=1）不可删除
- 如果角色下有关联用户，需提示先解除用户关联
- 如果角色下有关联菜单，会自动解除关联

---

### 7. 批量删除角色

**接口描述：** 批量删除多个角色

**请求方式：** `DELETE /api/v1/roles/batch`

**请求参数（Request Body）：**

| 参数名 | 类型 | 必填 | 校验规则 | 说明 |
|--------|------|------|----------|------|
| ids | []uint64 | 是 | required,min=1 | 角色ID列表 |

**请求示例：**
```json
{
  "ids": [3, 4, 5]
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "删除成功",
  "data": {
    "successCount": 3,
    "failCount": 0,
    "failIds": []
  }
}
```

---

### 8. 更新角色状态

**接口描述：** 启用或禁用角色

**请求方式：** `PUT /api/v1/roles/{id}/status`

**请求参数（Path Params）：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | uint64 | 是 | 角色ID |

**请求参数（Request Body）：**

| 参数名 | 类型 | 必填 | 校验规则 | 说明 |
|--------|------|------|----------|------|
| status | int8 | 是 | required,oneof=0 1 | 状态：0-禁用 1-启用 |

**请求示例：**
```json
{
  "status": 0
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "状态更新成功",
  "data": null
}
```

---

### 9. 分配角色菜单权限

**接口描述：** 为角色分配菜单权限

**请求方式：** `POST /api/v1/roles/{id}/menus`

**请求参数（Path Params）：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | uint64 | 是 | 角色ID |

**请求参数（Request Body）：**

| 参数名 | 类型 | 必填 | 校验规则 | 说明 |
|--------|------|------|----------|------|
| menuIds | []uint64 | 是 | required,min=1 | 菜单ID列表 |

**请求示例：**
```json
{
  "menuIds": [1, 2, 10, 11, 12, 20, 21, 30, 31, 40, 41]
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "菜单权限分配成功",
  "data": null
}
```

---

### 10. 获取角色菜单权限

**接口描述：** 获取角色已分配的菜单权限列表

**请求方式：** `GET /api/v1/roles/{id}/menus`

**请求参数（Path Params）：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | uint64 | 是 | 角色ID |

**请求参数（Query Params）：**

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| tree | bool | 否 | false | 是否返回树形结构 |
| type | string | 否 | - | 菜单类型筛选：directory-目录 menu-菜单 button-按钮 |

**请求示例：**
```http
GET /api/v1/roles/1/menus?tree=true&type=menu
```

**响应示例（tree=false）：**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": 1,
      "parentId": 0,
      "title": "系统管理",
      "path": "/system",
      "icon": "Setting",
      "type": "directory",
      "sort": 1,
      "status": 1
    },
    {
      "id": 10,
      "parentId": 1,
      "title": "用户管理",
      "path": "/system/users",
      "icon": "User",
      "type": "menu",
      "sort": 1,
      "status": 1
    }
  ]
}
```

**响应示例（tree=true）：**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": 1,
      "parentId": 0,
      "title": "系统管理",
      "path": "/system",
      "icon": "Setting",
      "type": "directory",
      "sort": 1,
      "status": 1,
      "children": [
        {
          "id": 10,
          "parentId": 1,
          "title": "用户管理",
          "path": "/system/users",
          "icon": "User",
          "type": "menu",
          "sort": 1,
          "status": 1,
          "children": []
        }
      ]
    }
  ]
}
```

---

### 11. 获取角色用户列表

**接口描述：** 获取角色下的用户列表

**请求方式：** `GET /api/v1/roles/{id}/users`

**请求参数（Path Params）：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | uint64 | 是 | 角色ID |

**请求参数（Query Params）：**

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 10 | 每页数量 |
| keyword | string | 否 | - | 搜索关键字（用户名或昵称） |
| status | int | 否 | - | 状态筛选：0-禁用 1-启用 |

**请求示例：**
```http
GET /api/v1/roles/1/users?page=1&page_size=10
```

**响应示例：**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "admin",
        "nickname": "超级管理员",
        "email": "admin@example.com",
        "avatar": "/uploads/avatar/admin.jpg",
        "status": 1,
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 10
  }
}
```

---

### 12. 检查角色编码是否存在

**接口描述：** 检查角色编码是否已存在（用于前端校验）

**请求方式：** `GET /api/v1/roles/check-code`

**请求参数（Query Params）：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| code | string | 是 | 角色编码 |
| exclude_id | uint64 | 否 | 排除的角色ID（用于编辑时校验） |

**请求示例：**
```http
GET /api/v1/roles/check-code?code=admin&exclude_id=1
```

**响应示例：**
```json
{
  "code": 200,
  "message": "检查成功",
  "data": {
    "exists": false,
    "message": "角色编码可用"
  }
}
```

---

## DTO 设计

### Request DTO（请求DTO）

#### 1. CreateRoleRequest 创建角色请求
```go
package request

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
    RoleCode  string   `json:"roleCode" binding:"required,min=2,max=50"`
    RoleName  string   `json:"roleName" binding:"required,min=2,max=50"`
    RoleDesc  string   `json:"roleDesc" binding:"omitempty,max=200"`
    DataScope int8     `json:"dataScope" binding:"omitempty,oneof=1 2 3 4"`
    Status    int8     `json:"status" binding:"omitempty,oneof=0 1"`
    Sort      int      `json:"sort" binding:"omitempty,gte=0"`
    MenuIDs   []uint64 `json:"menuIds" binding:"omitempty"`
}
```

#### 2. UpdateRoleRequest 更新角色请求
```go
// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
    RoleCode  *string  `json:"roleCode" binding:"omitempty,min=2,max=50"`
    RoleName  *string  `json:"roleName" binding:"omitempty,min=2,max=50"`
    RoleDesc  *string  `json:"roleDesc" binding:"omitempty,max=200"`
    DataScope *int8    `json:"dataScope" binding:"omitempty,oneof=1 2 3 4"`
    Status    *int8    `json:"status" binding:"omitempty,oneof=0 1"`
    Sort      *int     `json:"sort" binding:"omitempty,gte=0"`
    MenuIDs   []uint64 `json:"menuIds" binding:"omitempty"`
}
```

#### 3. RoleListQueryRequest 角色列表查询请求
```go
// RoleListQueryRequest 角色列表查询请求
type RoleListQueryRequest struct {
    Page     int    `form:"page" binding:"omitempty,gte=1"`
    PageSize int    `form:"page_size" binding:"omitempty,gte=1,lte=100"`
    Keyword  string `form:"keyword" binding:"omitempty,max=50"`
    Status   *int8  `form:"status" binding:"omitempty,oneof=0 1"`
    Sort     string `form:"sort" binding:"omitempty,oneof=id role_code role_name sort"`
    Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}
```

#### 4. UpdateRoleStatusRequest 更新角色状态请求
```go
// UpdateRoleStatusRequest 更新角色状态请求
type UpdateRoleStatusRequest struct {
    Status int8 `json:"status" binding:"required,oneof=0 1"`
}
```

#### 5. AssignRoleMenusRequest 分配角色菜单请求
```go
// AssignRoleMenusRequest 分配角色菜单请求
type AssignRoleMenusRequest struct {
    MenuIDs []uint64 `json:"menuIds" binding:"required,min=1"`
}
```

#### 6. BatchDeleteRequest 批量删除请求（通用）
```go
// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
    IDs []uint64 `json:"ids" binding:"required,min=1"`
}
```

#### 7. CheckRoleCodeRequest 检查角色编码请求
```go
// CheckRoleCodeRequest 检查角色编码请求
type CheckRoleCodeRequest struct {
    Code      string  `form:"code" binding:"required,min=2,max=50"`
    ExcludeID *uint64 `form:"exclude_id" binding:"omitempty"`
}
```

---

### Response DTO（响应DTO）

#### 1. RoleResponse 角色响应
```go
package response

import "time"

// RoleResponse 角色响应
type RoleResponse struct {
    ID        uint64    `json:"id"`
    RoleCode  string    `json:"roleCode"`
    RoleName  string    `json:"roleName"`
    RoleDesc  string    `json:"roleDesc"`
    DataScope int8      `json:"dataScope"`
    Status    int8      `json:"status"`
    Sort      int       `json:"sort"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

#### 2. RoleDetailResponse 角色详情响应
```go
// RoleDetailResponse 角色详情响应
type RoleDetailResponse struct {
    ID        uint64              `json:"id"`
    RoleCode  string              `json:"roleCode"`
    RoleName  string              `json:"roleName"`
    RoleDesc  string              `json:"roleDesc"`
    DataScope int8                `json:"dataScope"`
    Status    int8                `json:"status"`
    Sort      int                 `json:"sort"`
    CreatedAt time.Time           `json:"createdAt"`
    UpdatedAt time.Time           `json:"updatedAt"`
    Menus     []MenuTreeResponse  `json:"menus,omitempty"`
    Users     []UserRoleResponse  `json:"users,omitempty"`
}

// UserRoleResponse 角色用户简略响应
type UserRoleResponse struct {
    ID       uint64 `json:"id"`
    Username string `json:"username"`
    Nickname string `json:"nickname"`
    Email    string `json:"email"`
    Avatar   string `json:"avatar"`
}

// MenuTreeResponse 菜单树响应
type MenuTreeResponse struct {
    ID       uint64           `json:"id"`
    ParentID uint64           `json:"parentId"`
    Title    string           `json:"title"`
    Path     string           `json:"path"`
    Icon     string           `json:"icon"`
    Type     string           `json:"type"`
    Sort     int              `json:"sort"`
    Status   int8             `json:"status"`
    Children []MenuTreeResponse `json:"children,omitempty"`
}
```

#### 3. RoleListResponse 角色列表响应
```go
// RoleListResponse 角色列表响应
type RoleListResponse struct {
    List     []RoleResponse `json:"list"`
    Total    int64          `json:"total"`
    Page     int            `json:"page"`
    PageSize int            `json:"pageSize"`
}
```

#### 4. CheckCodeResponse 检查编码响应
```go
// CheckCodeResponse 检查编码响应
type CheckCodeResponse struct {
    Exists  bool   `json:"exists"`
    Message string `json:"message"`
}
```

#### 5. BatchDeleteResponse 批量删除响应
```go
// BatchDeleteResponse 批量删除响应
type BatchDeleteResponse struct {
    SuccessCount int      `json:"successCount"`
    FailCount   int      `json:"failCount"`
    FailIDs     []uint64 `json:"failIds"`
}
```

---

## 业务规则

1. **超级管理员角色保护**
   - ID为1的角色是超级管理员角色，不可删除
   - 超级管理员角色的状态不可禁用

2. **角色编码唯一性**
   - 角色编码必须全局唯一
   - 创建和更新时需要进行唯一性校验

3. **角色删除约束**
   - 如果角色下有用户关联，不能删除
   - 如果角色下有菜单权限，删除时会自动解除关联
   - 软删除不会物理删除数据，仅标记deleted_at字段

4. **数据权限说明**
   - 1-全部数据：可以查看所有数据
   - 2-本部门及以下：可以查看本部门及子部门数据
   - 3-本部门：只能查看本部门数据
   - 4-仅本人：只能查看自己的数据

5. **菜单权限分配**
   - 分配菜单时会覆盖原有菜单权限
   - 如果分配的是父级菜单，子级菜单不会自动分配（需要显式指定）
   - 支持分配目录、菜单、按钮三种类型的菜单

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 无权限操作 |
| 404 | 角色不存在 |
| 409 | 角色编码已存在 |
| 500 | 服务器内部错误 |

---

## 注意事项

1. 所有接口都需要进行身份认证（JWT Token）
2. 超级管理员拥有所有权限，其他角色根据分配的菜单权限进行访问控制
3. 删除角色时建议前端二次确认
4. 分页查询建议限制最大每页数量（如100条）
5. 角色名称建议控制在50个字符以内，避免前端显示过长
6. 角色描述建议使用简洁明了的语言
