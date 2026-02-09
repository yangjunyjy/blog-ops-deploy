package router

import (
	apiv1 "my-blog-backend/internal/api/v1"
	"my-blog-backend/internal/pkg/middleware"

	_ "my-blog-backend/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handlers 路由处理器集合
type Handlers struct {
	Auth         *apiv1.AuthHandler
	User         *apiv1.UserHandler
	Article      *apiv1.ArticleHandler
	Category     *apiv1.CategoryHandler
	Tag          *apiv1.TagHandler
	Comment      *apiv1.CommentHandler
	Series       *apiv1.SeriesHandler
	Favorite     *apiv1.FavoriteHandler
	Upload       *apiv1.UploadHandler
	UserActivity *apiv1.UserActivityHandler
	SysAuth      *apiv1.SysAuthHandler
	SysRole      *apiv1.SysRoleHandler
	SysMenu      *apiv1.SysMenuHandler
	SysUser      *apiv1.SysUserHandler
	Statistics   *apiv1.StatisticsHandler
	Host         *apiv1.HostHandler
	Ssh          *apiv1.SshHandler
	Sftp         *apiv1.SshFileHandler
}

// SetupRouter 设置路由
func SetupRouter(engine *gin.Engine, handlers *Handlers) {
	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 静态文件访问（上传的图片）
	engine.Static("/uploads", "./uploads")

	// swagger文档路由
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 路由
	v1 := engine.Group("/api/v1")
	{
		// 公开路由 - 无需认证
		setupPublicRoutes(v1, handlers)

		// 前台路由 - 需要登录
		setupFrontRoutes(v1, handlers)

		// 后台路由 - 需要管理员权限
		setupAdminRoutes(v1, handlers)

	}
}

// 公开路由 - 无需认证，前台和后台都可以访问
func setupPublicRoutes(router *gin.RouterGroup, handlers *Handlers) {
	public := router.Group("/public")
	{
		// 认证相关
		auth := public.Group("/auth")
		{
			auth.GET("/captcha", handlers.Auth.GetCaptcha)
			auth.POST("/login", handlers.Auth.LoginWithCode)
			auth.POST("/refresh", handlers.Auth.RefreshToken)
			auth.POST("/send-email-code", handlers.Auth.SendEmailCaptcha)
			auth.POST("/register", handlers.Auth.Register)
			auth.POST("/verify", handlers.Auth.VerifyEmailCaptcha)
			auth.POST("/loginByEmail", handlers.Auth.LoginWithEmail)
		}

		// 文章公开查询
		articles := public.Group("/articles")
		{
			articles.GET("", handlers.Article.ListArticles)             // 获取文章列表
			articles.GET("/search", handlers.Article.SearchArticles)    // 搜索文章
			articles.GET("/hot", handlers.Article.GetHotArticles)       // 获取热门文章
			articles.GET("/recent", handlers.Article.GetRecentArticles) // 获取最新文章
			articles.GET("/:id", handlers.Article.GetArticle)           // 获取文章详情
		}

		// 分类公开查询
		categories := public.Group("/categories")
		{
			categories.GET("", handlers.Category.ListCategories)       // 获取分类列表
			categories.GET("/tree", handlers.Category.GetCategoryTree) // 获取分类树
			categories.GET("/:id", handlers.Category.GetCategory)      // 获取分类详情
		}

		// 标签公开查询
		tags := public.Group("/tags")
		{
			tags.GET("", handlers.Tag.ListTags)                           // 获取标签列表
			tags.GET("/:id", handlers.Tag.GetTag)                         // 获取标签详情
			tags.GET("/article/:article_id", handlers.Tag.GetArticleTags) // 获取文章标签
		}

		// 评论公开查询
		comments := public.Group("/comments")
		{
			comments.GET("/article/:articleId", handlers.Comment.ListComments) // 获取文章评论
			comments.GET("/:id/replies", handlers.Comment.GetReplies)          // 获取子评论列表
			comments.GET("/:id/count", handlers.Comment.GetReplyCount)         // 获取子评论数量
		}

		// 系列公开查询
		series := public.Group("/series")
		{
			series.GET("", handlers.Series.ListSeries)                                         // 获取系列列表
			series.GET("/:id", handlers.Series.GetSeries)                                      // 获取系列详情
			series.GET("/:id/sections", handlers.Series.ListSections)                          // 获取系列章节列表
			series.GET("/sections/:id", handlers.Series.GetSection)                            // 获取章节详情
			series.GET("/sections/:id/subchapters", handlers.Series.ListSubchapters)           // 获取章节子章节列表
			series.GET("/subchapters/:id", handlers.Series.GetSubchapter)                      // 获取子章节详情
			series.GET("/subchapters/:id/articles", handlers.Series.GetArticlesBySubchapterID) // 获取子章节文章列表
		}
	}
}

// 前台路由 - 需要登录，普通用户和管理员都可以访问
func setupFrontRoutes(router *gin.RouterGroup, handlers *Handlers) {
	front := router.Group("/front")
	front.Use(middleware.AuthRequired())
	{
		// 文章互动
		articles := front.Group("/articles")
		{
			articles.POST("/:id/view", handlers.Article.ViewArticle)                  // 增加阅读量
			articles.GET("/:id/like/status", handlers.Article.CheckArticleLikeStatus) // 检查点赞状态
			articles.POST("/:id/like", handlers.Article.LikeArticle)                  // 点赞文章
			articles.DELETE("/:id/like", handlers.Article.UnlikeArticle)              // 取消点赞
			articles.POST("/:id/favorite", handlers.Article.FavoriteArticle)          // 收藏文章
			articles.DELETE("/:id/favorite", handlers.Article.UnfavoriteArticle)      // 取消收藏
		}

		// 评论操作
		comments := front.Group("/comments")
		{
			comments.GET("/:id", handlers.Comment.GetComment)                 // 获取评论详情
			comments.POST("", handlers.Comment.CreateComment)                 // 创建评论
			comments.PUT("/:id", handlers.Comment.UpdateComment)              // 更新自己的评论
			comments.DELETE("/:id", handlers.Comment.DeleteComment)           // 删除自己的评论
			comments.GET("/user/:user_id", handlers.Comment.ListUserComments) // 获取用户评论
			comments.POST("/:id/like", handlers.Comment.LikeComment)          // 点赞评论
			comments.DELETE("/:id/like", handlers.Comment.UnlikeComment)      // 取消点赞评论
		}

		// 用户信息
		users := front.Group("/users")
		{
			users.GET("/me", handlers.User.GetCurrentUser)                         // 获取当前用户信息
			users.PUT("/me", handlers.User.UpdateCurrentUser)                      // 更新当前用户信息
			users.GET("/:id", handlers.User.GetUser)                               // 获取用户详情
			users.GET("/:id/activities", handlers.UserActivity.GetUserActivities)  // 获取用户活动记录
			users.GET("/:id/activities/stats", handlers.UserActivity.GetUserStats) // 获取用户活动统计
		}

		// 收藏操作
		favorite := front.Group("/favorite")
		{
			favorite.GET("/userFolder", handlers.Favorite.GetUserFoldersWithCount)   // 获取用户收藏文件夹
			favorite.POST("/create", handlers.Favorite.CreateFolder)                 // 创建用户收藏文件夹
			favorite.PUT("/userFolder/:id", handlers.Favorite.UpdateFolder)          // 更新用户收藏文件夹
			favorite.DELETE("/userFolder/:id", handlers.Favorite.DeleteFolder)       // 删除用户收藏文件夹
			favorite.GET("/list", handlers.Favorite.GetUserFavorites)                // 获取收藏物品列表
			favorite.POST("/add", handlers.Favorite.AddFavorite)                     // 添加收藏物品
			favorite.POST("/remove", handlers.Favorite.RemoveFavorite)               // 移除收藏物品
			favorite.POST("/move", handlers.Favorite.MoveFavorite)                   // 移动收藏物品
			favorite.GET("/check/:articleId", handlers.Favorite.CheckFavoriteStatus) // 检查文章收藏状态
		}
	}
}

// 管理后台路由 - 统一使用 RBAC 认证
func setupAdminRoutes(router *gin.RouterGroup, handlers *Handlers) {
	// RBAC 公开路由 - 无需认证 (必须在认证路由之前注册)
	router.POST("/rbac/auth/login", handlers.SysAuth.Login)
	router.POST("/rbac/auth/send-email-code", handlers.SysAuth.SendEmailCode)
	// TOTP配置信息（公开可访问，客户端需要知道是否启用TOTP）
	router.GET("/rbac/auth/totp-config", middleware.GetTOTPConfig)
	router.GET("/rbac/auth/totp-status", middleware.CheckTOTPStatus)

	// RBAC 基础认证路由（只需要Session认证）
	rbacAuth := router.Group("/rbac")
	rbacAuth.Use(middleware.RBACAuth())
	{
		// 生成一次性Token的接口（只需要Session认证，不需要Once-Token）
		rbacAuth.POST("/auth/token", handlers.SysAuth.GenerateToken)
	}

	// RBAC 安全路由（需要双重认证：Session + Once-Token）
	rbacSecure := rbacAuth.Group("")
	rbacSecure.Use(middleware.OnceTokenMiddleware())
	{
		// 用户信息
		rbacSecure.POST("/auth/logout", handlers.SysAuth.Logout)
		rbacSecure.GET("/auth/info", handlers.SysAuth.GetInfo)
		rbacSecure.GET("/auth/menu", handlers.SysAuth.GetMenu)
		rbacSecure.GET("/auth/permission", handlers.SysAuth.CheckPermission)

		// 用户管理
		rbacSecure.DELETE("/users/:id", handlers.SysUser.DeleteUser)
		rbacSecure.GET("/users", handlers.SysUser.List)
		rbacSecure.POST("/users/update", handlers.SysUser.Update)
		rbacSecure.POST("/users/create", handlers.SysUser.Create)
		rbacSecure.POST("/users/assign", handlers.SysUser.Assign)
		rbacSecure.GET("/users/:id", handlers.SysUser.GetUserByID)
		rbacSecure.GET("/users/all", handlers.SysUser.GetAllUser)

		// 角色路由
		rbacSecure.GET("/roles", handlers.SysRole.List)
		rbacSecure.POST("/roles/assign", handlers.SysRole.AssignMenus)
		rbacSecure.GET("/roles/:id", handlers.SysRole.GetMenus)
		rbacSecure.DELETE("/roles/:id", handlers.SysRole.Delete)
		rbacSecure.POST("/roles/update", handlers.SysRole.Update)
		rbacSecure.POST("/roles/create", handlers.SysRole.Create)

		// 菜单管理
		rbacSecure.GET("/menus", handlers.SysMenu.GetAllMenus)
		rbacSecure.GET("/menus/byPaging", handlers.SysMenu.List)
		rbacSecure.GET("/menus/:id", handlers.SysMenu.GetByID)
		rbacSecure.GET("/menus/tree", handlers.SysMenu.GetTree)
		rbacSecure.DELETE("/menus/:id", handlers.SysMenu.Delete)
		rbacSecure.POST("/menus/create", handlers.SysMenu.Create)
		rbacSecure.POST("/menus/update", handlers.SysMenu.Update)

		// 文章管理
		rbacSecure.POST("/articles/import", handlers.Article.ImportMarkdownArticle)
		rbacSecure.PUT("/articles/:id/status", handlers.Article.UpdateArticleStatus)
		rbacSecure.POST("/articles", handlers.Article.CreateArticle)
		rbacSecure.PUT("/articles/:id", handlers.Article.UpdateArticle)
		rbacSecure.DELETE("/articles/:id", handlers.Article.DeleteArticle)

		// 评论管理
		rbacSecure.GET("/comments", handlers.Comment.ListAllComments)
		rbacSecure.GET("/comments/pending", handlers.Comment.ListPendingComments)
		rbacSecure.PUT("/comments/:id/approve", handlers.Comment.ApproveComment)
		rbacSecure.PUT("/comments/:id/reject", handlers.Comment.RejectComment)
		rbacSecure.DELETE("/comments/:id", handlers.Comment.DeleteComment)

		// 分类管理
		rbacSecure.PUT("/categories/:id", handlers.Category.UpdateCategory)
		rbacSecure.POST("/categories", handlers.Category.CreateCategory)
		rbacSecure.DELETE("/categories/:id", handlers.Category.DeleteCategory)

		// 标签管理
		rbacSecure.PUT("/tags/:id", handlers.Tag.UpdateTag)
		rbacSecure.DELETE("/tags/:id", handlers.Tag.DeleteTag)
		rbacSecure.POST("/tags", handlers.Tag.CreateTag)

		// 系列管理
		rbacSecure.GET("/series", handlers.Series.ListSeries)
		rbacSecure.PUT("/series/sections/:id", handlers.Series.UpdateSection)
		rbacSecure.DELETE("/series/sections/:id", handlers.Series.DeleteSection)
		rbacSecure.POST("/series/subchapters", handlers.Series.CreateSubchapter)
		rbacSecure.PUT("/series/subchapters/:id", handlers.Series.UpdateSubchapter)
		rbacSecure.DELETE("/series/subchapters/:id", handlers.Series.DeleteSubchapter)
		rbacSecure.POST("/series", handlers.Series.CreateSeries)
		rbacSecure.PUT("/series/:id", handlers.Series.UpdateSeries)
		rbacSecure.DELETE("/series/:id", handlers.Series.DeleteSeries)
		rbacSecure.POST("/series/:id/sections", handlers.Series.CreateSection)
		rbacSecure.POST("/series/subchapters/:id/articles", handlers.Series.AddArticleToSubchapter)

		// 文件上传
		rbacSecure.POST("/upload/markdown", handlers.Upload.UploadMarkdown)
		rbacSecure.POST("/upload/word", handlers.Upload.UploadWord)
		rbacSecure.POST("/upload/image", handlers.Upload.UploadImage)

		// 统计数据
		rbacSecure.GET("/statistics/dashboard", handlers.Statistics.GetDashboardStats)
		rbacSecure.GET("/statistics/articles", handlers.Statistics.GetArticleStats)
		rbacSecure.GET("/statistics/articles/trend", handlers.Statistics.GetArticleTrend)
		rbacSecure.GET("/statistics/views/trend", handlers.Statistics.GetViewTrend)
		rbacSecure.GET("/statistics/articles/hot", handlers.Statistics.GetHotArticles)
		rbacSecure.GET("/statistics/categories", handlers.Statistics.GetCategoryStats)
		rbacSecure.GET("/statistics/tags", handlers.Statistics.GetTagStats)
		rbacSecure.GET("/statistics/tags/hot", handlers.Statistics.GetHotTags)
		rbacSecure.GET("/statistics/users", handlers.Statistics.GetUserStats)
		rbacSecure.GET("/statistics/users/growth", handlers.Statistics.GetUserGrowth)
		rbacSecure.GET("/statistics/users/active", handlers.Statistics.GetActiveUsers)
		rbacSecure.GET("/statistics/users/distribution", handlers.Statistics.GetUserDistribution)
		rbacSecure.GET("/statistics/users/online", handlers.Statistics.GetOnlineUsers)

		// 主机管理
		rbacSecure.GET("/hosts", handlers.Host.ListHosts)
		rbacSecure.GET("/hosts/all", handlers.Host.GetAllHosts)
		rbacSecure.GET("/hosts/:id", handlers.Host.GetHost)
		rbacSecure.POST("/hosts", handlers.Host.CreateHost)
		rbacSecure.PUT("/hosts", handlers.Host.UpdateHost)
		rbacSecure.DELETE("/hosts/:id", handlers.Host.DeleteHost)
		rbacSecure.POST("/hosts/:id/test", handlers.Host.TestConnection)

		// SSH 终端（只需要 RBAC 认证，WebSocket 无法携带 Once-Token）
		rbacAuth.GET("/ssh/connect/:host_id", handlers.Ssh.WebSocketConnect)
		rbacSecure.GET("/ssh/sessions", handlers.Ssh.ListSessions)
		rbacSecure.DELETE("/ssh/sessions/:session_id", handlers.Ssh.CloseSession)

		//sftp终端
		rbacSecure.POST("/sftp/uploadFile", handlers.Sftp.UploadFile)
		rbacSecure.GET("/sftp/list", handlers.Sftp.List)
	}
}
