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

type TagHandler struct {
	tagService services.TagService
}

func NewTagHandler(tagService services.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// CreateTag 创建标签
// @Summary 创建标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param request body dtoRequest.CreateTagRequest true "标签信息"
// @Success 200 {object} response.Response
// @Router /api/v1/tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req dtoRequest.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	tag := &models.Tag{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: response.SafeDerefString(req.Description),
	}

	if err := h.tagService.CreateTag(tag); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建标签失败", err)
		return
	}

	response.Success(c, tag, "创建成功")
}

// UpdateTag 更新标签
// @Summary 更新标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Param request body dtoRequest.CreateTagRequest true "标签信息"
// @Success 200 {object} response.Response
// @Router /api/v1/tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID", err)
		return
	}

	var req dtoRequest.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	tag := &models.Tag{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: response.SafeDerefString(req.Description),
	}

	if err := h.tagService.UpdateTag(uint(id), tag); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新标签失败", err)
		return
	}

	response.Success(c, tag, "更新成功")
}

// DeleteTag 删除标签
// @Summary 删除标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} response.Response
// @Router /api/v1/tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID", err)
		return
	}

	if err := h.tagService.DeleteTag(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除标签失败", err)
		return
	}

	response.Success(c, nil, "删除成功")
}

// GetTag 获取标签详情
// @Summary 获取标签详情
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(12)
// @Success 200 {object} response.Response
// @Router /api/v1/tags/{id} [get]
func (h *TagHandler) GetTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID", err)
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

	tag, err := h.tagService.GetTag(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "标签不存在", err)
		return
	}

	// 获取该标签下的文章列表
	articles, total, err := h.tagService.GetTagArticles(uint(id), page, pageSize)
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

		// 获取文章的标签
		tags := make([]response.TagInfo, 0)
		for _, tag := range article.Tags {
			tags = append(tags, response.TagInfo{
				Name: tag.Name,
			})
		}

		articleItems = append(articleItems, response.ArticleListItem{
			ID:        uint(article.ID),
			Title:     article.Title,
			Summary:   article.Summary,
			Cover:     response.StringPtr(article.Cover),
			Category:  categoryInfo,
			Tags:      tags,
			Author:    authorInfo,
			Views:     int64(article.Views),
			Likes:     int64(article.Likes),
			Favorites: int64(article.Favorites),
			Comments:  int64(article.CommentCount),
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		})
	}

	tagInfo := response.TagInfo{
		Name: tag.Name,
	}

	articleListResp := &response.ArticleListResponse{
		List:     articleItems,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	tagDetailResp := &response.TagDetailResponse{
		Tag:      tagInfo,
		Articles: articleListResp,
	}

	response.Success(c, tagDetailResp, "")
}

// ListTags 获取标签列表
// @Summary 获取标签列表
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/tags [get]
func (h *TagHandler) ListTags(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	tags, total, err := h.tagService.ListTags(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取标签列表失败", err)
		return
	}

	// 转换为响应格式
	tagItems := make([]response.TagItem, 0, len(tags))
	for _, tag := range tags {
		tagItems = append(tagItems, response.TagItem{
			ID:           uint(tag.ID),
			Name:         tag.Name,
			ArticleCount: tag.ArticleCount,
		})
	}

	response.Success(c, gin.H{
		"items":     tagItems,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, "")
}

// GetArticleTags 获取文章的标签
// @Summary 获取文章的标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Router /api/v1/articles/{id}/tags [get]
func (h *TagHandler) GetArticleTags(c *gin.Context) {
	var req dtoRequest.ArticleIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	tags, err := h.tagService.GetArticleTags(req.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取标签失败", err)
		return
	}

	response.Success(c, tags, "")
}
