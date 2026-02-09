package repository

import (
	models "my-blog-backend/internal/models/frontendModel"
)

type UserActivityRepository interface {
	// Create 创建用户活动记录
	Create(activity *models.UserActivity) error

	// GetByUserID 获取用户的互动记录列表
	GetByUserID(userID uint, activityType string, page, pageSize int) ([]*models.UserActivity, int64, error)

	// GetStats 获取用户互动统计
	GetStats(userID uint) (map[string]int64, error)

	// Delete 删除活动记录
	Delete(id uint) error
}
