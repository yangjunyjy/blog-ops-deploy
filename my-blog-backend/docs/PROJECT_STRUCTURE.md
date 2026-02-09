# 项目结构说明

## 目录树

```
my-blog-backend/
├── cmd/                              # 应用入口
│   └── main.go                       # 主入口文件
│
├── internal/                          # 内部代码(不对外暴露)
│   ├── app/                         # 应用层
│   │   ├── app.go                   # 应用核心逻辑
│   │   ├── router.go                # 路由配置
│   │   └── wire.go                  # Wire依赖注入配置
│   │
│   ├── api/                         # API处理器层
│   │   ├── auth/                    # 认证相关API
│   │   └── v1/                      # API版本1
│   │
│   ├── config/                      # 配置管理
│   │   ├── config.go                # 配置结构体定义
│   │   └── loader.go                # 配置加载器
│   │
│   ├── pkg/                         # 通用工具包
│   │   ├── database/                # 数据库管理
│   │   │   ├── driver/              # 数据库驱动
│   │   │   │   ├── driver.go       # 驱动接口
│   │   │   │   ├── mysql.go        # MySQL驱动
│   │   │   │   ├── postgres.go      # PostgreSQL驱动
│   │   │   │   └── sqlite.go        # SQLite驱动
│   │   │   ├── manager.go            # 数据库管理器
│   │   │   ├── factory.go            # 工厂模式
│   │   │   ├── factory_registry.go    # 工厂注册表
│   │   │   ├── generic_factory.go    # 通用工厂
│   │   │   ├── repository.go         # 仓储接口
│   │   │   └── transaction.go       # 事务管理
│   │   │
│   │   ├── logger/                 # 日志管理
│   │   │   ├── logger.go           # 日志接口
│   │   │   └── zap_logger.go       # Zap实现
│   │   │
│   │   ├── middleware/             # HTTP中间件
│   │   │   ├── auth.go            # 认证中间件
│   │   │   ├── cors.go            # 跨域中间件
│   │   │   ├── recovery.go        # 恢复中间件
│   │   │   └── request_id.go      # 请求ID中间件
│   │   │
│   │   └── cache/                  # 缓存管理
│   │       ├── redis.go           # Redis缓存
│   │       └── memory.go          # 内存缓存
│   │
│   ├── repository/                  # 数据访问层(接口定义)
│   │   ├── user_repository.go       # 用户仓储接口
│   │   ├── role_repository.go       # 角色仓储接口
│   │   └── impl/                   # 仓储实现
│   │       ├── generic/             # 通用仓储
│   │       │   └── generic_repo.go # 通用仓储实现
│   │       ├── mysql/               # MySQL实现
│   │       ├── postgresql/          # PostgreSQL实现
│   │       └── sqlite/              # SQLite实现
│   │
│   ├── services/                    # 业务逻辑层
│   │   ├── user_service.go        # 用户服务
│   │   └── auth_service.go        # 认证服务
│   │
│   └── models/                      # 数据模型
│       ├── user.go                # 用户模型
│       ├── role.go                # 角色模型
│       └── ...
│
├── api/                            # API文档和定义
│   └── swagger.yaml                 # Swagger文档
│
├── migrations/                      # 数据库迁移脚本
│   ├── mysql/                       # MySQL迁移
│   ├── postgresql/                  # PostgreSQL迁移
│   └── sqlite/                      # SQLite迁移
│
├── deployments/                    # 部署配置
│   ├── docker-compose.yaml          # Docker Compose
│   ├── kubernetes/                 # K8s配置
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── ingress.yaml
│   └── scripts/                    # 部署脚本
│
├── tests/                          # 测试代码
│   ├── integration/                 # 集成测试
│   └── e2e/                        # 端到端测试
│
├── docs/                           # 项目文档
│   └── PROJECT_STRUCTURE.md         # 本文件
│
├── .env.example                    # 环境变量示例
├── .gitignore                      # Git忽略文件
├── go.mod                         # Go模块定义
├── go.sum                         # 依赖锁定文件
├── Makefile                       # 构建脚本
└── README.md                      # 项目说明
```

## 分层架构

### 1. 应用层 (internal/app)
- 负责应用启动和路由配置
- 使用 Wire 进行依赖注入
- 组装各个组件

### 2. API层 (internal/api)
- 处理HTTP请求和响应
- 参数验证
- 调用服务层

### 3. 服务层 (internal/services)
- 业务逻辑处理
- 事务管理
- 调用仓储层

### 4. 仓储层 (internal/repository)
- 数据访问接口定义
- 具体实现按数据库类型分离
- 支持泛型仓储

### 5. 模型层 (internal/models)
- 数据模型定义
- ORM映射
- 数据验证

### 6. 基础设施层 (internal/pkg)
- 数据库驱动
- 日志管理
- 缓存管理
- 中间件

## 核心设计模式

### 1. 依赖注入
使用 Wire 自动生成依赖注入代码:
```go
//go:build wireinject
func InitializeApplication(cfg *config.Config) (*Application, error) {
    panic("wire: compile")
}
```

### 2. 仓储模式
接口定义与实现分离:
```go
// 接口
type UserRepository interface {
    Create(user *User) error
}

// 实现
type mysqlUserRepository struct {
    db *gorm.DB
}
```

### 3. 工厂模式
数据库驱动工厂:
```go
func getDriver(name string) driver.Driver {
    switch name {
    case "mysql":
        return driver.NewMySQLDriver()
    case "postgresql":
        return driver.NewPostgreSQLDriver()
    }
}
```

### 4. 泛型仓储
减少重复代码:
```go
type GenericRepository[T any] interface {
    Create(entity *T) error
    Update(entity *T) error
}
```

## 多数据库支持

项目支持三种数据库:
- **MySQL**: 生产环境推荐
- **PostgreSQL**: 复杂查询场景
- **SQLite**: 开发和测试

切换数据库只需修改配置:
```yaml
database:
  driver: mysql  # 或 postgresql, sqlite
  dsn: ...
```

## 扩展指南

### 添加新模块

1. 创建模型 (internal/models)
2. 定义仓储接口 (internal/repository)
3. 实现仓储 (internal/repository/impl)
4. 创建服务 (internal/services)
5. 实现API (internal/api/v1)
6. 注册路由 (internal/app/router.go)

### 添加新的数据库支持

1. 创建驱动 (internal/pkg/database/driver)
2. 实现Driver接口
3. 在manager.go中注册
4. 创建迁移脚本 (migrations/new_db)
5. 实现仓储 (internal/repository/impl/new_db)

## 注意事项

1. **internal包**: 只能被项目内部引用,外部无法导入
2. **泛型使用**: Go 1.18+ 支持,利用泛型减少代码重复
3. **依赖注入**: 使用 Wire 管理复杂依赖关系
4. **配置管理**: 使用 Viper 支持多环境配置
5. **错误处理**: 统一的错误处理机制(待实现)
