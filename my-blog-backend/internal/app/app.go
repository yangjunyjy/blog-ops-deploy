package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	apiV1 "my-blog-backend/internal/api/v1"
	"my-blog-backend/internal/config"
	"my-blog-backend/internal/infrastructure/database"
	"my-blog-backend/internal/infrastructure/smtp"
	"my-blog-backend/internal/pkg/logger"
	"my-blog-backend/internal/pkg/middleware"
	"my-blog-backend/internal/pkg/session"
	"my-blog-backend/internal/pkg/token"
	impl "my-blog-backend/internal/repository/impl"
	implMysql "my-blog-backend/internal/repository/impl/mysql"
	"my-blog-backend/internal/router"
	"my-blog-backend/internal/services"
	"my-blog-backend/internal/ssh"

	"my-blog-backend/internal/api/v1/dto/response"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Application 应用主结构
type Application struct {
	config     *config.Config
	logger     logger.Logger
	dbManager  database.Manager
	httpServer *http.Server
	router     *gin.Engine
	handlers   *router.Handlers
}

// NewApplication 创建应用实例
func NewApplication() (*Application, error) {
	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 2. 初始化日志
	errlogger := logger.InitGlobalLogger(cfg.Log)
	if errlogger != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	app := &Application{
		config: cfg,
		logger: logger.L(),
	}

	// 设置JWT Secret到中间件
	middleware.JWTSecret = []byte(cfg.Auth.JWTSecret)

	// 初始化URL前缀（用于构建完整URL）
	response.InitURLPrefix(cfg)

	// 打印JWT配置用于调试
	app.logger.Info("JWT配置信息",
		logger.String("jwt_secret", cfg.Auth.JWTSecret[:min(10, len(cfg.Auth.JWTSecret))]+"..."),
		logger.String("jwt_expire", cfg.Auth.JWTExpire.String()),
		logger.Duration("jwt_expire_duration", cfg.Auth.JWTExpire),
	)

	// 4. 初始化smtp服务器
	if err := app.initSmtpServer(); err != nil {
		return nil, fmt.Errorf("failed to init smtp server: %w", err)
	}

	// 3. 初始化基础设施
	if err := app.initInfrastructure(); err != nil {
		return nil, fmt.Errorf("failed to init infrastructure: %w", err)
	}

	// 5. 初始化HTTP服务器
	if err := app.initHTTPServer(); err != nil {
		return nil, fmt.Errorf("failed to init HTTP server: %w", err)
	}
	return app, nil
}

// Run 启动应用
func (app *Application) Run() error {
	app.logger.Info("Starting application",
		logger.String("name", app.config.App.Name),
		logger.String("version", app.config.App.Version),
		logger.String("env", app.config.App.Env),
	)

	// 启动HTTP服务器
	go func() {
		app.logger.Info("HTTP server starting",
			logger.String("address", app.httpServer.Addr),
		)

		var err error
		if app.config.Server.TLS.Enabled {
			err = app.httpServer.ListenAndServeTLS(
				app.config.Server.TLS.CertFile,
				app.config.Server.TLS.KeyFile,
			)
		} else {
			err = app.httpServer.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			app.logger.Error("HTTP server failed", logger.Err("error", err))
			return
		}
	}()

	// 注册健康检查路由
	if app.config.Server.HealthCheck.Enabled {
		app.router.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	// 等待关闭信号
	return app.waitForShutdown()
}

// initInfrastructure 初始化基础设施
func (app *Application) initInfrastructure() error {
	// 初始化数据库管理器
	dbManager, err := database.NewManager(app.config)
	if err != nil {
		return fmt.Errorf("failed to create database manager: %w", err)
	}

	// 连接数据库
	if err := dbManager.Connect(); err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	app.dbManager = dbManager

	// 初始化服务
	app.initServices()

	return nil
}

// initServices 初始化业务服务
func (app *Application) initServices() {
	// 获取DB实例
	db := app.dbManager.GetPrimaryDB()

	// 创建Repository - 使用泛型仓储实现
	userRepo := implMysql.NewUserRepositoryImpl(db)
	articleRepo := implMysql.NewArticleRepositoryImpl(db)
	categoryRepo := implMysql.NewCategoryRepositoryImpl(db)
	tagRepo := implMysql.NewTagRepositoryImpl(db)
	commentRepo := implMysql.NewCommentRepositoryImpl(db)
	commentLikeRepo := implMysql.NewCommentLikeRepositoryImpl(db)
	articleLikeRepo := implMysql.NewArticleLikeRepositoryImpl(db)
	seriesRepo := implMysql.NewSeriesRepositoryImpl(db)
	sectionRepo := implMysql.NewSeriesSectionRepositoryImpl(db)
	subchapterRepo := implMysql.NewSeriesSubchapterRepositoryImpl(db)
	favoriteRepo := implMysql.NewFavoriteRepositoryImpl(db)

	// 创建CacheRepo
	redisCacheRepo := impl.NewRedisCacheRepo(app.dbManager.GetRedisClient())

	// 创建Service
	userService := services.NewUserService(userRepo)
	favoriteService := services.NewFavoriteService(favoriteRepo)
	userActivityService := services.NewUserActivityService(db)
	authService := services.NewAuthServiceWithCache(userService, favoriteService, redisCacheRepo, app.config.Auth.JWTSecret, app.config.Auth.JWTExpire)
	articleService := services.NewArticleService(articleRepo, tagRepo, commentRepo, commentLikeRepo, favoriteRepo, articleLikeRepo, userActivityService)
	categoryService := services.NewCategoryService(categoryRepo)
	tagService := services.NewTagService(tagRepo)

	// 创建敏感词服务
	sensitiveWordService := services.NewSensitiveWordService(app.config)

	// 创建评论服务（注入敏感词服务、点赞仓储和用户活动服务）
	commentService := services.NewCommentService(commentRepo, commentLikeRepo, app.config, sensitiveWordService, userActivityService)
	seriesService := services.NewSeriesService(seriesRepo, sectionRepo, subchapterRepo)

	// 创建上传服务
	uploadService := services.NewUploadService(&app.config.Upload)

	// 创建Handler
	app.handlers = &router.Handlers{
		User:         apiV1.NewUserHandler(userService),
		Auth:         apiV1.NewAuthHandler(authService),
		Article:      apiV1.NewArticleHandler(articleService),
		Category:     apiV1.NewCategoryHandler(categoryService),
		Tag:          apiV1.NewTagHandler(tagService),
		Comment:      apiV1.NewCommentHandler(commentService),
		Series:       apiV1.NewSeriesHandler(seriesService),
		Favorite:     apiV1.NewFavoriteHandler(favoriteService),
		Upload:       apiV1.NewUploadHandler(uploadService),
		UserActivity: apiV1.NewUserActivityHandler(userActivityService),
	}

	// 初始化 RBAC 服务（需要在 handlers 初始化之后）
	app.initRBACServices(db)
}

// initRBACServices 初始化 RBAC 服务
func (app *Application) initRBACServices(db *gorm.DB) {
	// 创建 RBAC Repository
	sysUserRepo := implMysql.NewSysUserRepositoryImpl(db)
	sysRoleRepo := implMysql.NewSysRoleRepositoryImpl(db)
	sysMenuRepo := implMysql.NewSysMenuRepositoryImpl(db)
	sysDeptRepo := implMysql.NewSysDeptRepositoryImpl(db)
	sysLogRepo := implMysql.NewSysLogRepositoryImpl(db)
	_ = implMysql.NewSysPostRepositoryImpl(db)

	// 创建 Redis Cache Repo
	redisCacheRepo := impl.NewRedisCacheRepo(app.dbManager.GetRedisClient())

	// 创建 RBAC Service
	permissionService := services.NewSysPermissionService(sysUserRepo, sysRoleRepo, sysMenuRepo)
	sysUserService := services.NewSysUserService(sysUserRepo, sysLogRepo, redisCacheRepo)
	sysRoleService := services.NewSysRoleService(sysRoleRepo)
	sysMenuService := services.NewSysMenuService(sysMenuRepo)
	_ = services.NewSysDeptService(sysDeptRepo)
	_ = services.NewSysLogService(sysLogRepo)

	// 创建 Session 和 Token 管理器
	redisClient := app.dbManager.GetRedisClient()

	sessionConfig := &session.Config{
		CookieName: "rbac_session_id",
		MaxAge:     86400,
		Path:       "/",
		Domain:     "",
		Secure:     false,
		HttpOnly:   true,
	}
	sessionManager := session.NewRedisSession(redisClient, sessionConfig)

	tokenConfig := token.DefaultConfig()
	tokenConfig.ExpireTime = 300
	tokenManager := token.NewRedisTokenManager(redisClient, tokenConfig)

	// 初始化中间件
	middleware.InitConfig(app.config)
	middleware.InitPermissionService(permissionService)
	middleware.InitSessionManager(sessionManager)
	middleware.InitTokenManager(tokenManager)
	middleware.InitOperationLogService(nil)

	// 创建 RBAC Handlers
	sysAuthHandler := apiV1.NewSysAuthHandler(
		sysUserService,
		permissionService,
		sessionManager,
		tokenManager,
	)
	sysRoleHandler := apiV1.NewSysRoleHandler(sysRoleService)
	sysMenuHandler := apiV1.NewSysMenuHandler(sysMenuService)
	sysUerHandler := apiV1.NewSysUserHandler(sysUserService)

	// 创建统计服务
	statisticsRepo := implMysql.NewStatisticsRepository(db)
	statisticsService := services.NewStatisticsService(statisticsRepo)
	statisticsHandler := apiV1.NewStatisticsHandler(statisticsService)

	// 创建主机管理服务和Handler
	hostRepo := implMysql.NewHostRepository(db)
	sshPool := ssh.NewPool(5 * time.Minute)
	hostService := services.NewHostService(hostRepo, sshPool)
	hostHandler := apiV1.NewHostHandler(hostService)
	sshHandler := apiV1.NewSshHandler(hostService, sshPool)
	sftpHandler := apiV1.NewSshFileHandler(hostService, sshPool, sshHandler.GetSessions())

	// 将 RBAC Handlers 添加到 handlers 结构体
	app.handlers.SysAuth = sysAuthHandler
	app.handlers.SysRole = sysRoleHandler
	app.handlers.SysMenu = sysMenuHandler
	app.handlers.Statistics = statisticsHandler
	app.handlers.SysUser = sysUerHandler

	// 添加主机管理和SSH Handlers
	app.handlers.Host = hostHandler
	app.handlers.Ssh = sshHandler

	app.handlers.Sftp = sftpHandler

	app.logger.Info("RBAC services initialized successfully")
}

// initHTTPServer 初始化HTTP服务器
func (app *Application) initHTTPServer() error {
	// 设置Gin模式
	app.setGinMode()

	// 创建Gin引擎
	app.router = gin.New()

	// 注册中间件
	app.registerMiddlewares()

	// 注册路由
	router.SetupRouter(app.router, app.handlers)

	// 创建HTTP服务器
	app.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", app.config.Server.Host, app.config.Server.Port),
		Handler:      app.router,
		ReadTimeout:  app.config.Server.ReadTimeout,
		WriteTimeout: app.config.Server.WriteTimeout,
		IdleTimeout:  app.config.Server.IdleTimeout,
	}

	return nil
}

// startHealthServer 启动健康检查服务器
func (app *Application) startHealthServer() {
	// 健康检查路由已经在Run方法中注册
}

func (app *Application) initSmtpServer() error {
	// 初始化smtp服务器
	err := smtp.InitEmailServer(app.config.EmailServer.SmtpServerConfig, app.config.EmailServer.CodeConfig)
	if err != nil {
		app.logger.Error("Failed to initialize SMTP server", zap.Error(err))
		return err
	}
	return nil
}

// waitForShutdown 等待关闭信号
func (app *Application) waitForShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	app.logger.Info("Received shutdown signal", zap.String("signal", sig.String()))

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.shutdown(ctx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	app.logger.Info("Application shutdown completed")
	return nil
}

// shutdown 优雅关闭
func (app *Application) shutdown(ctx context.Context) error {
	var errs []error

	// 关闭HTTP服务器
	if app.httpServer != nil {
		if err := app.httpServer.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("HTTP server shutdown error: %w", err))
		}
	}

	// 关闭数据库连接
	if app.dbManager != nil {
		if err := app.dbManager.Close(); err != nil {
			errs = append(errs, fmt.Errorf("database disconnect error: %w", err))
		}
	}

	// 同步日志
	if app.logger != nil {
		_ = app.logger.Sync()
	}

	// 返回组合错误
	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	return nil
}

// 设置Gin模式
func (app *Application) setGinMode() {
	switch app.config.App.Env {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

// 注册中间件
func (app *Application) registerMiddlewares() {
	// 恢复中间件
	app.router.Use(gin.Recovery())

	// 日志中间件
	app.router.Use(gin.Logger())

	// CORS中间件
	app.router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 允许的源列表
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:8080",
			"ws://localhost:3000",
			"ws://127.0.0.1:3000",
		}

		// 检查是否在允许的源列表中
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
}
