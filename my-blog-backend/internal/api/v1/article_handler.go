package api

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	dtoRequest "my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/api/v1/dto/response"
	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/pkg/logger"
	"my-blog-backend/internal/pkg/middleware"
	"my-blog-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleService services.ArticleService
}

func NewArticleHandler(articleService services.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

// CreateArticle 创建文章
// @Summary 创建文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param request body dtoRequest.CreateArticleRequest true "文章信息"
// @Success 200 {object} response.Response
// @Router /api/v1/articles [post]
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req dtoRequest.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 从JWT中获取当前用户ID
	var userID uint = 1

	article := &models.Article{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Cover:      response.SafeDerefString(req.Cover),
		CategoryID: response.SafeDerefUint64(req.CategoryID),
		AuthorID:   uint64(userID),
		Status:     req.Status,
		IsTop:      req.IsTop,
		Slug:       req.Slug,
	}

	if err := h.articleService.CreateArticle(article, req.TagIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建文章失败", err)
		return
	}

	response.Success(c, article, "创建成功")
}

// UpdateArticle 更新文章
// @Summary 更新文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param request body dtoRequest.UpdateArticleRequest true "文章信息"
// @Success 200 {object} response.Response
// @Router /api/v1/articles/{id} [put]
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文章ID", err)
		return
	}

	var req dtoRequest.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	article := &models.Article{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Cover:      response.SafeDerefString(req.Cover),
		CategoryID: response.SafeDerefUint64(req.CategoryID),
		Status:     req.Status,
		IsTop:      req.IsTop,
	}

	if err := h.articleService.UpdateArticle(uint(id), article, req.TagIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新文章失败", err)
		return
	}

	response.Success(c, article, "更新成功")
}

// DeleteArticle 删除文章
// @Summary 删除文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Router /api/v1/articles/{id} [delete]
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文章ID", err)
		return
	}

	if err := h.articleService.DeleteArticle(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除文章失败", err)
		return
	}

	response.Success(c, nil, "删除成功")
}

// GetArticle 获取文章详情
// @Summary 获取文章详情
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Router /api/v1/articles/{id} [get]
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文章ID", err)
		return
	}

	article, err := h.articleService.GetArticle(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "文章不存在", err)
		return
	}

	// 构建完整的封面URL
	if article.Cover != "" {
		article.Cover = response.BuildFullURL(article.Cover)
	}

	response.Success(c, article, "")
}

// ListArticles 获取文章列表
// @Summary 获取文章列表
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param status query int false "状态"
// @Param categoryID query int false "分类ID"
// @Success 200 {object} response.Response
// @Router /api/v1/articles [get]
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	var req dtoRequest.ArticleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 50 {
		req.PageSize = 10
	}

	var status *int = new(int)
	*status = int(req.Status)

	var categoryID *uint
	if req.CategoryID != 0 {
		categoryID = &req.CategoryID
	}

	articles, total, err := h.articleService.ListArticles(req.Page, req.PageSize, status, categoryID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取文章列表失败", err)
		return
	}

	// 构建完整的封面URL
	for i := range articles {
		if articles[i].Cover != "" {
			articles[i].Cover = response.BuildFullURL(articles[i].Cover)
		}
	}

	//logger.Info("articles", logger.Any("aticles", articles))
	response.Success(c, gin.H{
		"items":     articles,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	}, "")
}

// SearchArticles 搜索文章
// @Summary 搜索文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param keyword query string false "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/articles/search [get]
func (h *ArticleHandler) SearchArticles(c *gin.Context) {
	var req dtoRequest.ArticleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 50 {
		req.PageSize = 10
	}

	articles, total, err := h.articleService.SearchArticles(req.Keyword, req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "搜索失败", err)
		return
	}

	response.Success(c, gin.H{
		"items":     articles,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
		"keyword":   req.Keyword,
	}, "")
}

// ViewArticle 查看文章（增加阅读量）
// @Summary 查看文章（增加阅读量）
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Router /api/v1/articles/{id}/view [post]
func (h *ArticleHandler) ViewArticle(c *gin.Context) {
	var req dtoRequest.LikeArticleRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	if err := h.articleService.ViewArticle(req.ArticleID); err != nil {
		response.Error(c, http.StatusInternalServerError, "操作失败", err)
		return
	}

	response.Success(c, nil, "成功")
}

// GetHotArticles 获取热门文章
// @Summary 获取热门文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param limit query int false "限制数量" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/articles/hot [get]
func (h *ArticleHandler) GetHotArticles(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	articles, err := h.articleService.GetHotArticles(limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取热门文章失败", err)
		return
	}

	response.Success(c, articles, "")
}

// GetRecentArticles 获取最新文章
// @Summary 获取最新文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param limit query int false "限制数量" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/articles/recent [get]
func (h *ArticleHandler) GetRecentArticles(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	articles, err := h.articleService.GetRecentArticles(limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取最新文章失败", err)
		return
	}

	response.Success(c, articles, "")
}

// LikeArticle 点赞文章
// @Summary 点赞文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Router /api/v1/articles/{id}/like [post]
func (h *ArticleHandler) LikeArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文章ID", err)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	if err := h.articleService.LikeArticle(uint(id), userID.(uint)); err != nil {
		response.Error(c, http.StatusInternalServerError, "点赞失败", err)
		return
	}

	response.Success(c, nil, "点赞成功")
}

// UnlikeArticle 取消点赞文章
// @Summary 取消点赞文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Router /api/v1/articles/{id}/unlike [post]
func (h *ArticleHandler) UnlikeArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文章ID", err)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	if err := h.articleService.UnlikeArticle(uint(id), userID.(uint)); err != nil {
		response.Error(c, http.StatusInternalServerError, "取消点赞失败", err)
		return
	}

	response.Success(c, nil, "取消点赞成功")
}

// FavoriteArticle 收藏文章
// @Summary 收藏文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Router /api/v1/articles/{id}/favorite [post]
func (h *ArticleHandler) FavoriteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文章ID", err)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	if err := h.articleService.FavoriteArticle(uint(id), userID.(uint)); err != nil {
		response.Error(c, http.StatusInternalServerError, "收藏失败", err)
		return
	}

	response.Success(c, nil, "收藏成功")
}

// UnfavoriteArticle 取消收藏文章
// @Summary 取消收藏文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Router /api/v1/articles/{id}/unfavorite [post]
func (h *ArticleHandler) UnfavoriteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文章ID", err)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	if err := h.articleService.UnfavoriteArticle(uint(id), userID.(uint)); err != nil {
		response.Error(c, http.StatusInternalServerError, "取消收藏失败", err)
		return
	}

	response.Success(c, nil, "取消收藏成功")
}

// CheckArticleLikeStatus 检查文章点赞状态
// @Summary 检查文章点赞状态
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Router /api/v1/articles/{id}/like-status [get]
func (h *ArticleHandler) CheckArticleLikeStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文章ID", err)
		return
	}

	// 尝试从token获取用户ID（可选）
	var userID uint
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			token := parts[1]
			// 验证token但不强制要求
			if userClaims, err := middleware.ValidateToken(token); err == nil {
				userID = userClaims.UserID
			}
		}
	}

	// 未登录用户返回未点赞
	if userID == 0 {
		response.Success(c, gin.H{
			"isLiked": false,
		}, "")
		return
	}

	isLiked, err := h.articleService.CheckArticleLiked(uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "检查点赞状态失败", err)
		return
	}

	response.Success(c, gin.H{
		"isLiked": isLiked,
	}, "")
}

// UpdateArticleStatus 更新文章状态
// @Summary 更新文章状态
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param request body object true "状态信息"
// @Success 200 {object} response.Response
// @Router /api/v1/articles/{id}/status [put]
func (h *ArticleHandler) UpdateArticleStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文章ID", err)
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	article := &models.Article{
		Status: uint8(req.Status),
	}

	if err := h.articleService.UpdateArticle(uint(id), article, nil); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新文章状态失败", err)
		return
	}

	response.Success(c, nil, "更新成功")
}

// ImportMarkdownArticle 从Markdown文件导入文章内容
// @Summary 从Markdown文件导入文章内容
// @Tags 文章管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Markdown文件"
// @Success 200 {object} response.Response
// @Router /api/v1/articles/import/markdown [post]
func (h *ArticleHandler) ImportMarkdownArticle(c *gin.Context) {
	// 从表单获取文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "获取文件失败", err)
		return
	}

	// 验证文件类型
	ext := filepath.Ext(file.Filename)
	if ext != ".md" && ext != ".markdown" {
		response.Error(c, http.StatusBadRequest, "只支持.md或.markdown文件", nil)
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		logger.Error("打开文件失败", logger.Err("error", err))
		response.Error(c, http.StatusInternalServerError, "打开文件失败", err)
		return
	}
	defer src.Close()

	// 读取文件内容
	fileContent, err := io.ReadAll(src)
	if err != nil {
		logger.Error("读取Markdown文件失败", logger.Err("error", err))
		response.Error(c, http.StatusInternalServerError, "读取文件失败", err)
		return
	}

	content := string(fileContent)

	// 存储Markdown文件到本地目录
	uploadDir := "./markdown_files"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Error("创建目录失败", logger.Err("error", err))
		response.Error(c, http.StatusInternalServerError, "创建目录失败", err)
		return
	}

	uploadPath := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		logger.Error("保存文件失败", logger.Err("error", err))
		response.Error(c, http.StatusInternalServerError, "保存文件失败", err)
		return
	}

	// TODO: 根据配置选择存储方式（本地或远程）
	// 1. 本地存储: 直接保存到服务器本地目录
	// 2. 远程存储: 上传到OSS/COS等对象存储

	// 返回文件路径和内容供前端使用
	result := map[string]interface{}{
		"filename": file.Filename,
		"path":     uploadPath,
		"content":  content,
		"size":     file.Size,
	}

	logger.Info("Markdown文件上传成功",
		logger.String("filename", file.Filename),
		logger.Int64("size", file.Size),
		logger.String("path", uploadPath),
	)

	response.Success(c, result, "文件上传成功")
}
