package api

import (
	"strconv"

	"my-blog-backend/internal/api/v1/dto/response"
	"my-blog-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type StatisticsHandler struct {
	statsService *services.StatisticsService
}

func NewStatisticsHandler(statsService *services.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		statsService: statsService,
	}
}

// GetDashboardStats 获取仪表盘统计数据
// @Summary 获取仪表盘统计数据
// @Tags 统计管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/stats/dashboard [get]
func (h *StatisticsHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.statsService.GetDashboardStats()
	if err != nil {
		response.Error(c, 500, "获取统计数据失败", err)
		return
	}
	response.Success(c, stats, "获取统计数据成功")
}

// GetArticleStats 获取文章统计数据
// @Summary 获取文章统计数据
// @Tags 统计管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/stats/articles [get]
func (h *StatisticsHandler) GetArticleStats(c *gin.Context) {
	stats, err := h.statsService.GetArticleStats()
	if err != nil {
		response.Error(c, 500, "获取文章统计失败", err)
		return
	}
	response.Success(c, stats, "获取文章统计成功")
}

// GetArticleTrend 获取文章发布趋势
// @Summary 获取文章发布趋势
// @Tags 统计管理
// @Accept json
// @Produce json
// @Param days query int false "天数" default(7)
// @Success 200 {object} response.Response
// @Router /api/v1/stats/articles/trend [get]
func (h *StatisticsHandler) GetArticleTrend(c *gin.Context) {
	days := 7
	if d := c.Query("days"); d != "" {
		if dInt, err := strconv.Atoi(d); err == nil && dInt > 0 {
			days = dInt
		}
	}

	trend, err := h.statsService.GetArticleTrend(days)
	if err != nil {
		response.Error(c, 500, "获取文章趋势失败", err)
		return
	}
	response.Success(c, trend, "获取文章趋势成功")
}

// GetViewTrend 获取浏览量趋势
// @Summary 获取浏览量趋势
// @Tags 统计管理
// @Accept json
// @Produce json
// @Param days query int false "天数" default(7)
// @Success 200 {object} response.Response
// @Router /api/v1/stats/views/trend [get]
func (h *StatisticsHandler) GetViewTrend(c *gin.Context) {
	days := 7
	if d := c.Query("days"); d != "" {
		if dInt, err := strconv.Atoi(d); err == nil && dInt > 0 {
			days = dInt
		}
	}

	trend, err := h.statsService.GetViewTrend(days)
	if err != nil {
		response.Error(c, 500, "获取浏览量趋势失败", err)
		return
	}
	response.Success(c, trend, "获取浏览量趋势成功")
}

// GetHotArticles 获取热门文章
// @Summary 获取热门文章
// @Tags 统计管理
// @Accept json
// @Produce json
// @Param limit query int false "限制数量" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/stats/articles/hot [get]
func (h *StatisticsHandler) GetHotArticles(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if lInt, err := strconv.Atoi(l); err == nil && lInt > 0 {
			limit = lInt
		}
	}

	articles, err := h.statsService.GetHotArticles(limit)
	if err != nil {
		response.Error(c, 500, "获取热门文章失败", err)
		return
	}
	response.Success(c, articles, "获取热门文章成功")
}

// GetCategoryStats 获取分类统计
// @Summary 获取分类统计
// @Tags 统计管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/stats/categories [get]
func (h *StatisticsHandler) GetCategoryStats(c *gin.Context) {
	stats, err := h.statsService.GetCategoryStats()
	if err != nil {
		response.Error(c, 500, "获取分类统计失败", err)
		return
	}
	response.Success(c, stats, "获取分类统计成功")
}

// GetTagStats 获取标签统计
func (h *StatisticsHandler) GetTagStats(c *gin.Context) {
	stats, err := h.statsService.GetTagStats()
	if err != nil {
		response.Error(c, 500, "获取标签统计失败", err)
		return
	}
	response.Success(c, stats, "获取标签统计成功")
}

// GetHotTags 获取热门标签
// @Summary 获取热门标签
// @Tags 统计管理
// @Accept json
// @Produce json
// @Param limit query int false "限制数量" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/stats/tags/hot [get]
func (h *StatisticsHandler) GetHotTags(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if lInt, err := strconv.Atoi(l); err == nil && lInt > 0 {
			limit = lInt
		}
	}

	tags, err := h.statsService.GetHotTags(limit)
	if err != nil {
		response.Error(c, 500, "获取热门标签失败", err)
		return
	}
	response.Success(c, tags, "获取热门标签成功")
}

// GetUserStats 获取用户统计
// @Summary 获取用户统计
// @Tags 统计管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/stats/users [get]
func (h *StatisticsHandler) GetUserStats(c *gin.Context) {
	stats, err := h.statsService.GetUserStats()
	if err != nil {
		response.Error(c, 500, "获取用户统计失败", err)
		return
	}
	response.Success(c, stats, "获取用户统计成功")
}

// GetUserGrowth 获取用户增长趋势
// @Summary 获取用户增长趋势
// @Tags 统计管理
// @Accept json
// @Produce json
// @Param days query int false "天数" default(7)
// @Success 200 {object} response.Response
// @Router /api/v1/stats/users/growth [get]
func (h *StatisticsHandler) GetUserGrowth(c *gin.Context) {
	days := 7
	if d := c.Query("days"); d != "" {
		if dInt, err := strconv.Atoi(d); err == nil && dInt > 0 {
			days = dInt
		}
	}

	growth, err := h.statsService.GetUserGrowth(days)
	if err != nil {
		response.Error(c, 500, "获取用户增长趋势失败", err)
		return
	}
	response.Success(c, growth, "获取用户增长趋势成功")
}

// GetActiveUsers 获取活跃用户趋势
// @Summary 获取活跃用户趋势
// @Tags 统计管理
// @Accept json
// @Produce json
// @Param days query int false "天数" default(7)
// @Success 200 {object} response.Response
// @Router /api/v1/stats/users/active [get]
func (h *StatisticsHandler) GetActiveUsers(c *gin.Context) {
	days := 7
	if d := c.Query("days"); d != "" {
		if dInt, err := strconv.Atoi(d); err == nil && dInt > 0 {
			days = dInt
		}
	}

	users, err := h.statsService.GetActiveUsers(days)
	if err != nil {
		response.Error(c, 500, "获取活跃用户趋势失败", err)
		return
	}
	response.Success(c, users, "获取活跃用户趋势成功")
}

// GetUserDistribution 获取用户分布
// @Summary 获取用户分布
// @Tags 统计管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/stats/users/distribution [get]
func (h *StatisticsHandler) GetUserDistribution(c *gin.Context) {
	distribution, err := h.statsService.GetUserDistribution()
	if err != nil {
		response.Error(c, 500, "获取用户分布失败", err)
		return
	}
	response.Success(c, distribution, "获取用户分布成功")
}

// GetOnlineUsers 获取在线用户数
// @Summary 获取在线用户数
// @Tags 统计管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/stats/users/online [get]
func (h *StatisticsHandler) GetOnlineUsers(c *gin.Context) {
	count, err := h.statsService.GetOnlineUsers()
	if err != nil {
		response.Error(c, 500, "获取在线用户数失败", err)
		return
	}
	response.Success(c, count, "获取在线用户数成功")
}

