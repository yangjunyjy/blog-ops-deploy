package mysql

import (
	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"

	"gorm.io/gorm"
)

type UserActivityRepositoryImpl struct {
	db *gorm.DB
}

func NewUserActivityRepository(db *gorm.DB) repository.UserActivityRepository {
	return &UserActivityRepositoryImpl{db: db}
}

// Create 创建用户活动记录
func (r *UserActivityRepositoryImpl) Create(activity *models.UserActivity) error {
	return r.db.Create(activity).Error
}

// GetByUserID 获取用户的互动记录列表
func (r *UserActivityRepositoryImpl) GetByUserID(userID uint, activityType string, page, pageSize int) ([]*models.UserActivity, int64, error) {
	var activities []*models.UserActivity
	var total int64

	query := r.db.Model(&models.UserActivity{}).Where("user_id = ?", userID)

	// 按类型过滤
	if activityType != "" && activityType != "all" {
		query = query.Where("type = ?", activityType)
	}

	// 统计总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&activities).Error

	return activities, total, err
}

// GetStats 获取用户互动统计
func (r *UserActivityRepositoryImpl) GetStats(userID uint) (map[string]int64, error) {
	var stats []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}

	err := r.db.Model(&models.UserActivity{}).
		Select("type, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("type").
		Find(&stats).Error

	if err != nil {
		return nil, err
	}

	// 构建统计结果
	result := map[string]int64{
		"total":    0,
		"like":     0,
		"comment":  0,
		"share":    0,
		"favorite": 0,
	}

	for _, stat := range stats {
		result[stat.Type] = stat.Count
		result["total"] += stat.Count
	}

	return result, nil
}

// Delete 删除活动记录
func (r *UserActivityRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.UserActivity{}, id).Error
}
