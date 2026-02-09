package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	dtoRequest "my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/api/v1/dto/response"
	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/services"
)

// SeriesHandler 系列处理器
type SeriesHandler struct {
	seriesService *services.SeriesService
}

// NewSeriesHandler 创建系列处理器
func NewSeriesHandler(seriesService *services.SeriesService) *SeriesHandler {
	return &SeriesHandler{seriesService: seriesService}
}

// ==================== 系列相关 ====================

// CreateSeries 创建系列
// @Summary 创建系列
// @Tags 系列管理
// @Accept json
// @Produce json
// @Param request body dtoRequest.CreateSeriesRequest true "系列信息"
// @Success 200 {object} response.Response
// @Router /api/v1/series [post]
func (h *SeriesHandler) CreateSeries(c *gin.Context) {
	var req dtoRequest.CreateSeriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error(), err)
		return
	}

	series := &models.Series{
		Name:        req.Name,
		Slug:        req.Slug,
		Icon:        response.SafeDerefString(req.Icon),
		Description: response.SafeDerefString(req.Description),
		Cover:       response.SafeDerefString(req.Cover),
		SortOrder:   req.SortOrder,
		Status:      1,
	}

	if err := h.seriesService.CreateSeries(series); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建系列失败", err)
		return
	}

	response.Success(c, series, "创建成功")
}

// UpdateSeries 更新系列
// @Summary 更新系列
// @Tags 系列管理
// @Accept json
// @Produce json
// @Param id path int true "系列ID"
// @Param request body dtoRequest.UpdateSeriesRequest true "系列信息"
// @Success 200 {object} response.Response
// @Router /api/v1/series/{id} [put]
func (h *SeriesHandler) UpdateSeries(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID", err)
		return
	}

	var req dtoRequest.UpdateSeriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	series, err := h.seriesService.GetSeriesByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "系列不存在", err)
		return
	}

	if req.Name != nil {
		series.Name = *req.Name
	}
	if req.Slug != nil {
		series.Slug = *req.Slug
	}
	if req.Icon != nil {
		series.Icon = *req.Icon
	}
	if req.Description != nil {
		series.Description = *req.Description
	}
	if req.Cover != nil {
		series.Cover = *req.Cover
	}
	if req.SortOrder != nil {
		series.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		series.Status = *req.Status
	}

	if err := h.seriesService.UpdateSeries(series); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败", err)
		return
	}

	response.Success(c, series, "更新成功")
}

// DeleteSeries 删除系列
// @Summary 删除系列
// @Tags 系列管理
// @Produce json
// @Param id path int true "系列ID"
// @Success 200 {object} response.Response
// @Router /api/v1/series/{id} [delete]
func (h *SeriesHandler) DeleteSeries(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID", err)
		return
	}

	if err := h.seriesService.DeleteSeries(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败", err)
		return
	}

	response.Success(c, nil, "删除成功")
}

// GetSeries 获取系列详情
// @Summary 获取系列详情
// @Tags 系列管理
// @Produce json
// @Param id path int true "系列ID"
// @Success 200 {object} response.Response
// @Router /api/v1/series/{id} [get]
func (h *SeriesHandler) GetSeries(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID", err)
		return
	}

	series, err := h.seriesService.GetSeriesDetail(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "系列不存在", err)
		return
	}

	// 转换章节信息
	sections := make([]response.SectionInfo, 0, len(series.Sections))
	chapters := make([]response.ChapterInfo, 0, len(series.Sections))

	for _, section := range series.Sections {
		subchapters := make([]response.SubchapterInfo, 0, len(section.Subchapters))

		// 为前端兼容，创建章节对象（包含所有子章节的文章）
		chapterArticles := make([]response.ArticleListItem, 0)

		for _, subchapter := range section.Subchapters {
			// 收集文章ID
			articleIDs := make([]uint, 0, len(subchapter.Articles))
			// 收集文章列表
			articles := make([]response.ArticleListItem, 0, len(subchapter.Articles))

			for _, article := range subchapter.Articles {
				articleIDs = append(articleIDs, uint(article.ID))

				// 转换文章列表项
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

				tags := make([]response.TagInfo, 0, len(article.Tags))
				for _, tag := range article.Tags {
					tags = append(tags, response.TagInfo{
						Name: tag.Name,
					})
				}

				articleItem := response.ArticleListItem{
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
				}

				chapterArticles = append(chapterArticles, articleItem)
				articles = append(articles, articleItem)
			}

			subchapters = append(subchapters, response.SubchapterInfo{
				ID:          uint(subchapter.ID),
				Name:        subchapter.Name,
				Description: subchapter.Description,
				SortOrder:   subchapter.SortOrder,
				ArticleIDs:  articleIDs,
				Articles:    articles,
			})
		}

		sections = append(sections, response.SectionInfo{
			ID:          uint(section.ID),
			Name:        section.Name,
			Description: section.Description,
			SortOrder:   section.SortOrder,
			Subchapters: subchapters,
		})

		// 添加到 chapters（前端兼容）
		chapters = append(chapters, response.ChapterInfo{
			ID:          uint(section.ID),
			Name:        section.Name,
			Description: section.Description,
			SortOrder:   section.SortOrder,
			Articles:    chapterArticles,
		})
	}

	seriesDetail := response.SeriesDetail{
		ID:          uint(series.ID),
		Name:        series.Name,
		Icon:        series.Icon,
		Description: series.Description,
		Cover:       response.StringPtr(series.Cover),
		SortOrder:   series.SortOrder,
		Status:      series.Status,
		Sections:    sections,
		Chapters:    chapters,
	}

	response.Success(c, seriesDetail, "")
}

// ListSeries 分页获取系列列表
// @Summary 获取系列列表
// @Tags 系列管理
// @Produce json
// @Param request query dtoRequest.SeriesListRequest false "查询参数"
// @Success 200 {object} response.Response
// @Router /api/v1/series [get]
func (h *SeriesHandler) ListSeries(c *gin.Context) {
	var req dtoRequest.SeriesListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 10
	}

	series, total, err := h.seriesService.ListSeries(req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取列表失败", err)
		return
	}

	// 转换为响应格式
	seriesListItems := make([]response.SeriesListItem, 0, len(series))
	for _, s := range series {
		seriesListItems = append(seriesListItems, response.SeriesListItem{
			ID:          uint(s.ID),
			Name:        s.Name,
			Icon:        s.Icon,
			Description: s.Description,
			Cover:       response.StringPtr(s.Cover),
			SortOrder:   s.SortOrder,
			Status:      s.Status,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		})
	}

	seriesListResp := response.SeriesListResponse{
		List:     seriesListItems,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	response.Success(c, seriesListResp, "")
}

// ==================== 章节相关 ====================

// CreateSection 创建章节
// @Summary 创建章节
// @Tags 章节管理
// @Accept json
// @Produce json
// @Param id path int true "系列ID"
// @Param request body dtoRequest.CreateSectionRequest true "章节信息"
// @Success 200 {object} response.Response
// @Router /api/v1/series/{id}/sections [post]
func (h *SeriesHandler) CreateSection(c *gin.Context) {
	seriesID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的系列ID", err)
		return
	}

	var req dtoRequest.CreateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	section := &models.SeriesSection{
		SeriesID:    seriesID,
		Name:        req.Name,
		Description: response.SafeDerefString(req.Description),
		SortOrder:   req.SortOrder,
	}

	if err := h.seriesService.CreateSection(section); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建章节失败", err)
		return
	}

	response.Success(c, section, "创建成功")
}

// UpdateSection 更新章节
// @Summary 更新章节
// @Tags 章节管理
// @Accept json
// @Produce json
// @Param id path int true "章节ID"
// @Param request body dtoRequest.UpdateSectionRequest true "章节信息"
// @Success 200 {object} response.Response
// @Router /api/v1/series/sections/{id} [put]
func (h *SeriesHandler) UpdateSection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID", err)
		return
	}

	var req dtoRequest.UpdateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	section, err := h.seriesService.GetSectionByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "章节不存在", err)
		return
	}

	if req.Name != nil {
		section.Name = *req.Name
	}
	if req.Description != nil {
		section.Description = *req.Description
	}
	if req.SortOrder != nil {
		section.SortOrder = *req.SortOrder
	}

	if err := h.seriesService.UpdateSection(section); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败", err)
		return
	}

	response.Success(c, section, "更新成功")
}

// DeleteSection 删除章节
// @Summary 删除章节
// @Tags 章节管理
// @Produce json
// @Param id path int true "章节ID"
// @Success 200 {object} response.Response
// @Router /api/v1/series/sections/{id} [delete]
func (h *SeriesHandler) DeleteSection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID", err)
		return
	}

	if err := h.seriesService.DeleteSection(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败", err)
		return
	}

	response.Success(c, nil, "删除成功")
}

// GetSection 获取章节详情
// @Summary 获取章节详情
// @Tags 章节管理
// @Produce json
// @Param id path int true "章节ID"
// @Success 200 {object} response.Response
// @Router /api/v1/series/sections/{id} [get]
func (h *SeriesHandler) GetSection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID", err)
		return
	}

	section, err := h.seriesService.GetSectionDetail(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "章节不存在", err)
		return
	}

	response.Success(c, section, "")
}

// ListSections 分页获取章节列表
// @Summary 获取章节列表
// @Tags 章节管理
// @Produce json
// @Param series_id path int true "系列ID"
// @Param request query dtoRequest.SectionListRequest false "查询参数"
// @Success 200 {object} response.Response
// @Router /api/v1/series/{id}/sections [get]
func (h *SeriesHandler) ListSections(c *gin.Context) {
	seriesID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的系列ID", err)
		return
	}

	var req dtoRequest.SectionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 10
	}

	sections, total, err := h.seriesService.ListSectionsBySeriesID(uint(seriesID), req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取列表失败", err)
		return
	}

	response.Success(c, gin.H{
		"items":     sections,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	}, "")
}

// ==================== 子章节相关 ====================

// CreateSubchapter 创建子章节
// @Summary 创建子章节
// @Tags 子章节管理
// @Accept json
// @Produce json
// @Param request body dtoRequest.CreateSubchapterRequest true "子章节信息"
// @Success 200 {object} response.Response
// @Router /api/v1/series/subchapters [post]
func (h *SeriesHandler) CreateSubchapter(c *gin.Context) {
	var req dtoRequest.CreateSubchapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	subchapter := &models.SeriesSubchapter{
		SectionID:   req.SectionID,
		Name:        req.Name,
		Description: response.SafeDerefString(req.Description),
		SortOrder:   req.SortOrder,
	}

	if err := h.seriesService.CreateSubchapter(subchapter); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建子章节失败", err)
		return
	}

	response.Success(c, subchapter, "创建成功")
}

// GetSubchapter 获取子章节详情
// @Summary 获取子章节详情
// @Tags 子章节管理
// @Produce json
// @Param id path int true "子章节ID"
// @Success 200 {object} response.Response
// @Router /api/v1/series/subchapters/{id} [get]
func (h *SeriesHandler) GetSubchapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID", err)
		return
	}

	subchapter, err := h.seriesService.GetSubchapterByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "子章节不存在", err)
		return
	}

	response.Success(c, subchapter, "")
}

// UpdateSubchapter 更新子章节
// @Summary 更新子章节
// @Tags 子章节管理
// @Accept json
// @Produce json
// @Param id path int true "子章节ID"
// @Param request body dtoRequest.UpdateSubchapterRequest true "子章节信息"
// @Success 200 {object} response.Response
// @Router /api/v1/series/subchapters/{id} [put]
func (h *SeriesHandler) UpdateSubchapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID", err)
		return
	}

	var req dtoRequest.UpdateSubchapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	subchapter, err := h.seriesService.GetSubchapterByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "子章节不存在", err)
		return
	}

	if req.Name != nil {
		subchapter.Name = *req.Name
	}
	if req.Description != nil {
		subchapter.Description = *req.Description
	}
	if req.SortOrder != nil {
		subchapter.SortOrder = *req.SortOrder
	}

	if err := h.seriesService.UpdateSubchapter(subchapter); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败", err)
		return
	}

	response.Success(c, subchapter, "更新成功")
}

// DeleteSubchapter 删除子章节
// @Summary 删除子章节
// @Tags 子章节管理
// @Produce json
// @Param id path int true "子章节ID"
// @Success 200 {object} response.Response
// @Router /api/v1/series/subchapters/{id} [delete]
func (h *SeriesHandler) DeleteSubchapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID", err)
		return
	}

	if err := h.seriesService.DeleteSubchapter(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败", err)
		return
	}

	response.Success(c, nil, "删除成功")
}

// ListSubchapters 获取子章节列表
// @Summary 获取子章节列表
// @Tags 子章节管理
// @Produce json
// @Param section_id path int true "章节ID"
// @Param request query dtoRequest.SubchapterListRequest false "查询参数"
// @Success 200 {object} response.Response
// @Router /api/v1/series/sections/{id}/subchapters [get]
func (h *SeriesHandler) ListSubchapters(c *gin.Context) {
	sectionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的章节ID", err)
		return
	}

	var req dtoRequest.SubchapterListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 10
	}

	subchapters, total, err := h.seriesService.ListSubchaptersBySectionID(uint(sectionID), req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取列表失败", err)
		return
	}

	response.Success(c, gin.H{
		"items":     subchapters,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	}, "")
}

// AddArticleToSubchapter 添加文章到子章节
// @Summary 添加文章到子章节
// @Tags 子章节管理
// @Accept json
// @Produce json
// @Param id path int true "子章节ID"
// @Param request body dtoRequest.AddArticleToSubchapterRequest true "文章信息"
// @Success 200 {object} response.Response
// @Router /api/v1/series/subchapters/{id}/articles [post]
func (h *SeriesHandler) AddArticleToSubchapter(c *gin.Context) {
	subchapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的子章节ID", err)
		return
	}

	var req dtoRequest.AddArticleToSubchapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	if err := h.seriesService.AddArticleToSubchapter(uint(subchapterID), req.ArticleID, req.SortOrder); err != nil {
		response.Error(c, http.StatusInternalServerError, "添加失败", err)
		return
	}

	response.Success(c, nil, "添加成功")
}

// RemoveArticleFromSubchapter 从子章节移除文章
// @Summary 从子章节移除文章
// @Tags 子章节管理
// @Produce json
// @Param subchapter_id path int true "子章节ID"
// @Param article_id path int true "文章ID"
// @Success 200 {object} response.Response
// @Router /api/v1/series/subchapters/{subchapter_id}/articles/{article_id} [delete]
func (h *SeriesHandler) RemoveArticleFromSubchapter(c *gin.Context) {
	subchapterID, err := strconv.ParseUint(c.Param("subchapter_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的子章节ID", err)
		return
	}

	articleID, err := strconv.ParseUint(c.Param("article_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文章ID", err)
		return
	}

	if err := h.seriesService.RemoveArticleFromSubchapter(uint(subchapterID), uint(articleID)); err != nil {
		response.Error(c, http.StatusInternalServerError, "移除失败", err)
		return
	}

	response.Success(c, nil, "移除成功")
}

// GetArticlesBySubchapterID 获取子章节的文章列表
// @Summary 获取子章节的文章列表
// @Tags 子章节管理
// @Produce json
// @Param id path int true "子章节ID"
// @Success 200 {object} response.Response
// @Router /api/v1/series/subchapters/{id}/articles [get]
func (h *SeriesHandler) GetArticlesBySubchapterID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的子章节ID", err)
		return
	}

	articles, err := h.seriesService.GetArticlesBySubchapterID(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取失败", err)
		return
	}

	response.Success(c, articles, "")
}
