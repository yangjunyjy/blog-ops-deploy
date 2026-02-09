package api

import (
	"fmt"
	"my-blog-backend/internal/api/v1/dto/response"
	_ "my-blog-backend/internal/models"
	"my-blog-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserActivityHandler struct {
	activityService *services.UserActivityService
}

func NewUserActivityHandler(activityService *services.UserActivityService) *UserActivityHandler {
	return &UserActivityHandler{
		activityService: activityService,
	}
}

// GetUserActivitiesRequest 获取用户活动记录请求
type GetUserActivitiesRequest struct {
	Page    int    `form:"page,default=1"`
	PerPage int    `form:"perPage,default=10"`
	Type    string `form:"type"` // all, like, comment, share, favorite
}

// GetUserActivities 获取用户活动记录
// @Summary 获取用户活动记录
// @Tags 用户活动
// @Accept json
// @Produce json
// @Param userId path int true "用户ID"
// @Param page query int false "页码"
// @Param perPage query int false "每页数量"
// @Param type query string false "活动类型"
// @Success 200 {object} response.Response{data=[]models.UserActivity}
// @Router /api/v1/front/users/{userId}/activities [get]
func (h *UserActivityHandler) GetUserActivities(c *gin.Context) {
	var req GetUserActivitiesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 从路径参数获取用户ID
	userIDStr := c.Param("id")
	if userIDStr == "" {
		response.Error(c, http.StatusBadRequest, "用户ID不能为空", nil)
		return
	}

	var userID uint
	_, err := fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil || userID == 0 {
		response.Error(c, http.StatusBadRequest, "用户ID格式错误", err)
		return
	}

	pageSize := req.PerPage
	if pageSize > 50 {
		pageSize = 50
	}

	activities, total, err := h.activityService.GetUserActivities(userID, req.Type, req.Page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取活动记录失败", err)
		return
	}

	// 获取统计数据
	stats, _ := h.activityService.GetUserStats(userID)

	response.Success(c, gin.H{
		"list":  activities,
		"total": total,
		"stats": stats,
	}, "获取数据成功")
}

// GetUserStats 获取用户活动统计
// @Summary 获取用户活动统计
// @Tags 用户活动
// @Accept json
// @Produce json
// @Param userId path int true "用户ID"
// @Success 200 {object} response.Response{data=map[string]int64}
// @Router /api/v1/front/users/{userId}/activities/stats [get]
func (h *UserActivityHandler) GetUserStats(c *gin.Context) {
	// 从路径参数获取用户ID
	userIDStr := c.Param("id")
	if userIDStr == "" {
		response.Error(c, http.StatusBadRequest, "用户ID不能为空", nil)
		return
	}

	var userID uint
	_, err := fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil || userID == 0 {
		response.Error(c, http.StatusBadRequest, "用户ID格式错误", err)
		return
	}

	stats, err := h.activityService.GetUserStats(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取统计信息失败", err)
		return
	}

	response.Success(c, stats, "获取统计信息成功")
}
