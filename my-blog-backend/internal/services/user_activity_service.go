package services

import (
	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"
	implMysql "my-blog-backend/internal/repository/impl/mysql"

	"gorm.io/gorm"
)

type UserActivityService struct {
	activityRepo repository.UserActivityRepository
	articleRepo  repository.ArticleRepository
}

func NewUserActivityService(db *gorm.DB) *UserActivityService {
	return &UserActivityService{
		activityRepo: implMysql.NewUserActivityRepository(db),
		articleRepo:  implMysql.NewArticleRepositoryImpl(db),
	}
}

// CreateActivity 创建活动记录
func (s *UserActivityService) CreateActivity(activity *models.UserActivity) error {
	// 如果有文章ID，查询文章标题
	if activity.ArticleID > 0 {
		article, err := s.articleRepo.GetByID(activity.ArticleID)
		if err == nil && article != nil {
			activity.ArticleTitle = article.Title
		}
	}
	return s.activityRepo.Create(activity)
}

// GetUserActivities 获取用户活动列表
func (s *UserActivityService) GetUserActivities(userID uint, activityType string, page, pageSize int) ([]*models.UserActivity, int64, error) {
	return s.activityRepo.GetByUserID(userID, activityType, page, pageSize)
}

// GetUserStats 获取用户活动统计
func (s *UserActivityService) GetUserStats(userID uint) (map[string]int64, error) {
	return s.activityRepo.GetStats(userID)
}

// DeleteActivity 删除活动记录
func (s *UserActivityService) DeleteActivity(id uint) error {
	return s.activityRepo.Delete(id)
}

// RecordLike 记录点赞活动
func (s *UserActivityService) RecordLike(userID uint, articleID uint) error {
	activity := &models.UserActivity{
		UserID:    userID,
		ArticleID: articleID,
		Type:      "like",
	}
	return s.CreateActivity(activity)
}

// RecordComment 记录评论活动
func (s *UserActivityService) RecordComment(userID uint, articleID uint, content string) error {
	activity := &models.UserActivity{
		UserID:    userID,
		ArticleID: articleID,
		Type:      "comment",
		Content:   content,
	}
	return s.CreateActivity(activity)
}

// RecordShare 记录分享活动
func (s *UserActivityService) RecordShare(userID uint, articleID uint, platform string) error {
	activity := &models.UserActivity{
		UserID:    userID,
		ArticleID: articleID,
		Type:      "share",
		Platform:  platform,
	}
	return s.CreateActivity(activity)
}

// RecordFavorite 记录收藏活动
func (s *UserActivityService) RecordFavorite(userID uint, articleID uint, folderName string) error {
	activity := &models.UserActivity{
		UserID:     userID,
		ArticleID:  articleID,
		Type:       "favorite",
		FolderName: folderName,
	}
	return s.CreateActivity(activity)
}
