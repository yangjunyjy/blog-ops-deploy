# RBAC权限管理系统实现指南

## 1. 数据库设计

已完成迁移文件：`migrations/20250123000001_create_rbac_tables.sql`

包含以下表：
- `sys_user` - 用户表
- `sys_role` - 角色表
- `sys_menu` - 菜单表
- `sys_dept` - 部门表
- `sys_user_role` - 用户角色关联表
- `sys_role_menu` - 角色菜单关联表
- `sys_dict_type` - 字典类型表
- `sys_dict_data` - 字典数据表
- `sys_operation_log` - 操作日志表
- `sys_login_log` - 登录日志表
- `sys_post` - 岗位表
- `sys_user_post` - 用户岗位关联表

## 2. 模型定义 (Models)

已完成：
- `internal/models/common.go` - 基础模型
- `internal/models/sys_user.go` - 用户模型
- `internal/models/sys_role.go` - 角色模型
- `internal/models/constants.go` - 常量定义

待完成：
- `sys_menu.go` - 菜单模型
- `sys_dept.go` - 部门模型
- `sys_dict_type.go` - 字典类型模型
- `sys_dict_data.go` - 字典数据模型
- `sys_post.go` - 岗位模型
- `sys_operation_log.go` - 操作日志模型
- `sys_login_log.go` - 登录日志模型
- `sys_user_role.go` - 用户角色关联模型
- `sys_role_menu.go` - 角色菜单关联模型
- `sys_user_post.go` - 用户岗位关联模型

## 3. 会话和Token管理

### Session管理 (Cookie存储)

实现位置：`internal/pkg/session/`

**功能需求**：
1. Session ID生成和验证
2. Cookie管理 (设置和读取)
3. Session过期处理
4. Session数据存储 (可选Redis或内存)

### One-time Token机制 (Redis存储)

实现位置：`internal/pkg/token/`

**功能需求**：
1. Token生成 (UUID)
2. Redis存储Token (key: `once_token:{token_id}`, value: 用户信息)
3. Token验证 (使用后立即删除)
4. Token过期处理 (默认5分钟)

**请求头要求**：
- Header: `Once-Token: {token_id}`

**流程**：
```
1. 用户登录 -> 生成Session ID -> 设置Cookie
2. 每次请求携带Cookie中的Session ID
3. 敏感操作需要携带Once-Token请求头
4. 验证Token后立即删除，确保一次性使用
```

## 4. Repository层实现

实现位置：`internal/repository/impl/mysql/`

### 需要实现的Repository：
1. `sys_user_repo.go` - 用户仓储
2. `sys_role_repo.go` - 角色仓储
3. `sys_menu_repo.go` - 菜单仓储
4. `sys_dept_repo.go` - 部门仓储
5. `sys_post_repo.go` - 岗位仓储
6. `sys_dict_repo.go` - 字典仓储
7. `sys_log_repo.go` - 日志仓储
8. `sys_user_role_repo.go` - 用户角色关联仓储
9. `sys_role_menu_repo.go` - 角色菜单关联仓储
10. `sys_user_post_repo.go` - 用户岗位关联仓储

### 主要方法：
```go
// 用户Repository
FindByUsername(username string) (*SysUser, error)
List(page, pageSize int) ([]*SysUser, int64, error)
Create(user *SysUser) error
Update(user *SysUser) error
Delete(id uint64) error
GetUserRoles(userID uint64) ([]*SysRole, error)

// 角色Repository
List(page, pageSize int) ([]*SysRole, int64, error)
Create(role *SysRole) error
Update(role *SysRole) error
Delete(id uint64) error
AssignMenus(roleID uint64, menuIDs []uint64) error

// 菜单Repository
GetUserMenus(userID uint64) ([]*SysMenu, error)
GetMenuTree() ([]*SysMenu, error)
Create(menu *SysMenu) error
Update(menu *SysMenu) error
Delete(id uint64) error
```

## 5. Service层实现

实现位置：`internal/services/`

### 需要实现的Service：
1. `user_service.go` - 用户服务
2. `role_service.go` - 角色服务
3. `menu_service.go` - 菜单服务
4. `dept_service.go` - 部门服务
5. `dict_service.go` - 字典服务
6. `log_service.go` - 日志服务
7. `permission_service.go` - 权限服务
8. `session_service.go` - 会话服务
9. `once_token_service.go` - 一次性Token服务

### 用户服务主要方法：
```go
Login(username, password, ip string) (*SysUser, string, error)
Logout(userID uint64) error
GetUserInfo(userID uint64) (*SysUser, error)
CreateUser(user *SysUser, creatorID uint64) error
UpdateUser(user *SysUser, updaterID uint64) error
DeleteUser(id, operatorID uint64) error
ResetPassword(id uint64, newPassword string) error
```

### 权限服务主要方法：
```go
CheckPermission(userID uint64, menuCode string) bool
GetUserPermissions(userID uint64) ([]string, error)
GetUserMenus(userID uint64) ([]*SysMenu, error)
```

## 6. API Handler实现

实现位置：`internal/api/v1/`

### 需要实现的Handler：
1. `user_handler.go` - 用户管理
2. `role_handler.go` - 角色管理
3. `menu_handler.go` - 菜单管理
4. `dept_handler.go` - 部门管理
5. `dict_handler.go` - 字典管理
6. `log_handler.go` - 日志管理
7. `auth_handler.go` - 认证授权

### 认证Handler方法：
```go
POST   /api/v1/auth/login       - 用户登录
POST   /api/v1/auth/logout      - 用户登出
GET    /api/v1/auth/info        - 获取当前用户信息
GET    /api/v1/auth/menu        - 获取当前用户菜单
POST   /api/v1/auth/token       - 生成一次性Token
```

### 用户管理Handler方法：
```go
GET    /api/v1/users             - 用户列表
POST   /api/v1/users             - 创建用户
GET    /api/v1/users/:id        - 用户详情
PUT    /api/v1/users/:id        - 更新用户
DELETE /api/v1/users/:id        - 删除用户
PUT    /api/v1/users/:id/reset   - 重置密码
```

## 7. 中间件实现

实现位置：`internal/pkg/middleware/`

### 需要实现的中间件：
1. `session.go` - Session中间件
2. `auth.go` - 认证中间件
3. `permission.go` - 权限中间件
4. `once_token.go` - 一次性Token中间件
5. `operation_log.go` - 操作日志中间件

### Session中间件
```go
func SessionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从Cookie读取Session ID
        // 2. 从存储获取Session数据
        // 3. 验证Session有效性
        // 4. 将用户信息存入context
    }
}
```

### 权限中间件
```go
func PermissionMiddleware(menuCode string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从context获取用户信息
        // 2. 检查是否为超级管理员
        // 3. 检查用户权限
        // 4. 无权限则返回403
    }
}
```

### 一次性Token中间件
```go
func OnceTokenMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从Header读取Once-Token
        // 2. 从Redis验证Token
        // 3. 验证后删除Token
        // 4. 无效则返回401
    }
}
```

## 8. 配置文件修改

需要在 `config/config.yaml` 添加：

```yaml
session:
  cookieName: "session_id"
  maxAge: 86400  # 24小时
  path: "/"
  domain: ""
  secure: false
  httpOnly: true

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  poolSize: 10

onceToken:
  expireTime: 300  # 5分钟
```

## 9. 路由注册

修改 `internal/router/router.go`：

```go
// 公开路由
publicGroup := r.Group("/api/v1/auth")
{
    publicGroup.POST("/login", authHandler.Login)
    publicGroup.POST("/logout", authHandler.Logout)
}

// 需要认证的路由
authGroup := r.Group("/api/v1")
authGroup.Use(middleware.SessionMiddleware())
{
    authGroup.GET("/auth/info", authHandler.GetInfo)
    authGroup.GET("/auth/menu", authHandler.GetMenu)
    authGroup.POST("/auth/token", authHandler.GenerateOnceToken)

    // 用户管理
    userGroup := authGroup.Group("/users")
    {
        userGroup.GET("", userHandler.List)
        userGroup.POST("", userHandler.Create)
        userGroup.GET("/:id", userHandler.Get)
        userGroup.PUT("/:id", userHandler.Update)
        userGroup.DELETE("/:id", userHandler.Delete)
    }

    // 需要一次性Token的路由
    sensitiveGroup := authGroup.Group("")
    sensitiveGroup.Use(middleware.OnceTokenMiddleware())
    {
        userGroup.PUT("/:id/reset", userHandler.ResetPassword)
        userGroup.DELETE("/:id", userHandler.Delete)
    }
}
```

## 10. 前端集成

### Cookie处理
前端需要自动携带Cookie，使用Axios的 `withCredentials: true`

```javascript
axios.defaults.withCredentials = true
```

### 一次性Token处理
敏感操作需要先获取Token，然后在请求头携带：

```javascript
// 1. 获取Token
const token = await generateOnceToken()

// 2. 请求携带Token
axios.delete('/api/v1/users/1', {
  headers: {
    'Once-Token': token
  }
})
```

## 11. 实施步骤

1. ✅ 创建数据库迁移文件和基础模型
2. ⬜ 完成所有模型定义
3. ⬜ 实现Session管理
4. ⬜ 实现One-time Token机制
5. ⬜ 实现Repository层
6. ⬜ 实现Service层
7. ⬜ 实现API Handler
8. ⬜ 实现中间件
9. ⬜ 配置Redis连接
10. ⬜ 注册路由
11. ⬜ 测试功能
12. ⬜ 前端集成

## 12. 注意事项

1. **密码加密**：使用bcrypt加密，默认密码需加密
2. **Session安全**：使用HttpOnly Cookie防止XSS
3. **一次性Token**：使用后立即删除，防止重放攻击
4. **权限缓存**：可以考虑缓存用户权限信息
5. **日志记录**：重要操作需要记录操作日志
6. **数据权限**：根据角色的data_scope过滤数据
7. **超级管理员**：跳过权限检查，拥有所有权限

## 13. Redis依赖

需要安装Redis客户端库：

```bash
go get github.com/go-redis/redis/v8
```

配置示例：

```go
redisClient := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password set
    DB:       0,  // use default DB
})
```
