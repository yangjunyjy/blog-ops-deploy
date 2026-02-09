package api

import (
	"net/http"
	"strconv"

	dtoRequest "my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/api/v1/dto/response"
	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory 创建分类
// @Summary 创建分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param request body dtoRequest.CreateCategoryRequest true "分类信息"
// @Success 200 {object} response.Response
// @Router /api/v1/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dtoRequest.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: response.SafeDerefString(req.Description),
		Icon:        response.SafeDerefString(req.Icon),
	}

	if err := h.categoryService.CreateCategory(category); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建分类失败", err)
		return
	}

	response.Success(c, category, "创建成功")
}

// UpdateCategory 更新分类
// @Summary 更新分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Param request body dtoRequest.CreateCategoryRequest true "分类信息"
// @Success 200 {object} response.Response
// @Router /api/v1/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分类ID", err)
		return
	}

	var req dtoRequest.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: response.SafeDerefString(req.Description),
		Icon:        response.SafeDerefString(req.Icon),
	}

	if err := h.categoryService.UpdateCategory(uint(id), category); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新分类失败", err)
		return
	}

	response.Success(c, category, "更新成功")
}

// DeleteCategory 删除分类
// @Summary 删除分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response
// @Router /api/v1/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分类ID", err)
		return
	}

	if err := h.categoryService.DeleteCategory(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除分类失败", err)
		return
	}

	response.Success(c, nil, "删除成功")
}

// GetCategory 获取分类详情
// @Summary 获取分类详情
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(12)
// @Success 200 {object} response.Response
// @Router /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分类ID", err)
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "12"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 12
	}

	category, err := h.categoryService.GetCategory(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "分类不存在", err)
		return
	}

	// 获取该分类下的文章列表
	articles, total, err := h.categoryService.GetCategoryArticles(uint(id), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取文章列表失败", err)
		return
	}

	// 转换为响应格式
	articleItems := make([]response.ArticleListItem, 0, len(articles))
	for _, article := range articles {
		var categoryInfo response.CategoryInfo
		if article.Category != nil {
			categoryInfo = response.CategoryInfo{
				ID:   article.Category.ID,
				Name: article.Category.Name,
			}
		}

		var authorInfo response.AuthorInfo
		if article.Author != nil {
			authorInfo = response.AuthorInfo{
				ID:     article.Author.ID,
				Name:   article.Author.Username,
				Avatar: response.StringPtr(article.Author.Avatar),
				Bio:    response.StringPtr(article.Author.Bio),
			}
		}

		articleItems = append(articleItems, response.ArticleListItem{
			ID:        uint(article.ID),
			Title:     article.Title,
			Summary:   article.Summary,
			Cover:     response.StringPtr(article.Cover),
			Category:  categoryInfo,
			Author:    authorInfo,
			Views:     int64(article.Views),
			Likes:     int64(article.Likes),
			Favorites: int64(article.Favorites),
			Comments:  int64(article.CommentCount),
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		})
	}

	categoryResp := response.CategoryInfo{
		ID:   category.ID,
		Name: category.Name,
	}

	articleListResp := &response.ArticleListResponse{
		List:     articleItems,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	categoryDetailResp := &response.CategoryDetailResponse{
		Category: categoryResp,
		Articles: articleListResp,
	}

	response.Success(c, categoryDetailResp, "")
}

// ListCategories 获取分类列表
// @Summary 获取分类列表
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	categories, total, err := h.categoryService.ListCategories(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取分类列表失败", err)
		return
	}

	// 转换为响应格式
	categoryItems := make([]response.CategoryItem, 0, len(categories))
	for _, category := range categories {
		categoryItems = append(categoryItems, response.CategoryItem{
			ID:           category.ID,
			Name:         category.Name,
			ArticleCount: category.ArticleCount,
		})
	}

	response.Success(c, gin.H{
		"items":     categoryItems,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, "")
}

// GetCategoryTree 获取分类树
// @Summary 获取分类树
// @Tags 分类管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/categories/tree [get]
func (h *CategoryHandler) GetCategoryTree(c *gin.Context) {
	categories, err := h.categoryService.GetCategoryTree()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取分类树失败", err)
		return
	}

	response.Success(c, categories, "")
}
