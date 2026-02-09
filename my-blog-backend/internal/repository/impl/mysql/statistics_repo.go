package mysql

import (
	"time"

	models "my-blog-backend/internal/models/frontendModel"
	"my-blog-backend/internal/repository"

	"gorm.io/gorm"
)

type StatisticsRepositoryImpl struct {
	db *gorm.DB
}

func NewStatisticsRepository(db *gorm.DB) repository.StatisticsRepository {
	return &StatisticsRepositoryImpl{db: db}
}

// GetArticleCount 获取文章数量
func (r *StatisticsRepositoryImpl) GetArticleCount(status *int) (int64, error) {
	query := r.db.Model(&models.Article{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

// GetTotalViews 获取总浏览量
func (r *StatisticsRepositoryImpl) GetTotalViews() (int64, error) {
	var count int64
	err := r.db.Model(&models.ArticleView{}).Count(&count).Error
	return count, err
}

// GetTotalLikes 获取总点赞数
func (r *StatisticsRepositoryImpl) GetTotalLikes() (int64, error) {
	var count int64
	err := r.db.Model(&models.ArticleLike{}).Count(&count).Error
	return count, err
}

// GetCommentCount 获取评论数量
func (r *StatisticsRepositoryImpl) GetCommentCount(status *int) (int64, error) {
	query := r.db.Model(&models.Comment{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

// GetCategoryCount 获取分类数量
func (r *StatisticsRepositoryImpl) GetCategoryCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Category{}).Count(&count).Error
	return count, err
}

// GetTagCount 获取标签数量
func (r *StatisticsRepositoryImpl) GetTagCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Tag{}).Count(&count).Error
	return count, err
}

// GetSeriesCount 获取系列数量
func (r *StatisticsRepositoryImpl) GetSeriesCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Series{}).Count(&count).Error
	return count, err
}

// GetArticleTrend 获取文章发布趋势
func (r *StatisticsRepositoryImpl) GetArticleTrend(days int) ([]int, error) {
	var trend []int
	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		err := r.db.Model(&models.Article{}).Where("DATE(created_at) = ?", date).Count(&count).Error
		if err != nil {
			return nil, err
		}
		trend = append(trend, int(count))
	}
	return trend, nil
}

// GetViewsTrend 获取浏览量趋势（总数）
func (r *StatisticsRepositoryImpl) GetViewsTrend(days int) (int64, error) {
	startDate := time.Now().AddDate(0, 0, -days)
	var count int64
	err := r.db.Model(&models.ArticleView{}).Where("created_at >= ?", startDate).Count(&count).Error
	return count, err
}

// GetViewsTrendByDay 获取按天统计的浏览量趋势
func (r *StatisticsRepositoryImpl) GetViewsTrendByDay(days int) ([]int, error) {
	var trend []int
	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		err := r.db.Model(&models.ArticleView{}).Where("DATE(created_at) = ?", date).Count(&count).Error
		if err != nil {
			return nil, err
		}
		trend = append(trend, int(count))
	}
	return trend, nil
}

// GetUserGrowthTrend 获取用户增长趋势
func (r *StatisticsRepositoryImpl) GetUserGrowthTrend(days int) ([]string, []int, error) {
	var dates []string
	var counts []int
	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		dates = append(dates, date)
		var count int64
		err := r.db.Model(&models.User{}).Where("DATE(created_at) = ?", date).Count(&count).Error
		if err != nil {
			return nil, nil, err
		}
		counts = append(counts, int(count))
	}
	return dates, counts, nil
}

// GetActiveUserTrend 获取活跃用户趋势
func (r *StatisticsRepositoryImpl) GetActiveUserTrend(days int) ([]string, []int, error) {
	var dates []string
	var counts []int
	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		dates = append(dates, date)
		startTime := date + " 00:00:00"
		endTime := date + " 23:59:59"
		var count int64
		err := r.db.Model(&models.UserActivity{}).Where("created_at >= ? AND created_at <= ?", startTime, endTime).
			Distinct("user_id").
			Count(&count).Error
		if err != nil {
			return nil, nil, err
		}
		counts = append(counts, int(count))
	}
	return dates, counts, nil
}

// GetHotArticles 获取热门文章
func (r *StatisticsRepositoryImpl) GetHotArticles(limit int) ([]interface{}, error) {
	type HotArticle struct {
		ID           uint   `json:"id"`
		Title        string `json:"title"`
		ViewCount    int    `json:"views"`
		LikeCount    int    `json:"likes"`
		CommentCount int    `json:"commentCount"`
	}

	var articles []models.Article
	err := r.db.Order("view_count DESC").Limit(limit).Find(&articles).Error
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(articles))
	for i, article := range articles {
		var likeCount int64
		r.db.Model(&models.ArticleLike{}).Where("article_id = ?", article.ID).Count(&likeCount)

		var commentCount int64
		r.db.Model(&models.Comment{}).Where("article_id = ? AND status = ?", article.ID, 1).Count(&commentCount)

		result[i] = HotArticle{
			ID:           uint(article.ID),
			Title:        article.Title,
			ViewCount:    int(article.Views),
			LikeCount:    int(likeCount),
			CommentCount: int(commentCount),
		}
	}
	return result, nil
}

// GetHotTags 获取热门标签
func (r *StatisticsRepositoryImpl) GetHotTags(limit int) ([]interface{}, error) {
	type HotTag struct {
		ID           uint   `json:"id"`
		Name         string `json:"name"`
		ArticleCount int    `json:"articleCount"`
	}

	var tags []HotTag
	err := r.db.Model(&models.Tag{}).
		Select("tags.id, tags.name, COUNT(article_tags.article_id) as article_count").
		Joins("LEFT JOIN article_tags ON tags.id = article_tags.tag_id").
		Group("tags.id, tags.name").
		Order("article_count DESC").
		Limit(limit).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(tags))
	for i, tag := range tags {
		result[i] = tag
	}
	return result, nil
}

// GetCategoryStats 获取分类统计
func (r *StatisticsRepositoryImpl) GetCategoryStats() ([]interface{}, error) {
	type CategoryStat struct {
		ID           uint   `json:"id"`
		Name         string `json:"name"`
		ArticleCount int    `json:"articleCount"`
	}

	var stats []CategoryStat
	err := r.db.Model(&models.Category{}).
		Select("id, name, (SELECT COUNT(*) FROM articles WHERE category_id = categories.id AND status = 1) as article_count").
		Find(&stats).Error
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(stats))
	for i, stat := range stats {
		result[i] = stat
	}
	return result, nil
}

// GetTagStats 获取标签统计
func (r *StatisticsRepositoryImpl) GetTagStats() ([]interface{}, error) {
	type TagStat struct {
		ID           uint   `json:"id"`
		Name         string `json:"name"`
		ArticleCount int    `json:"articleCount"`
	}

	var stats []TagStat
	err := r.db.Model(&models.Tag{}).
		Select("tags.id, tags.name, COUNT(article_tags.article_id) as article_count").
		Joins("LEFT JOIN article_tags ON tags.id = article_tags.tag_id").
		Group("tags.id, tags.name").
		Find(&stats).Error
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(stats))
	for i, stat := range stats {
		result[i] = stat
	}
	return result, nil
}

// GetUserCount 获取用户总数
func (r *StatisticsRepositoryImpl) GetUserCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}

// GetNewUserCount 获取新用户数
func (r *StatisticsRepositoryImpl) GetNewUserCount(date string) (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("DATE(created_at) = ?", date).Count(&count).Error
	return count, err
}

// GetActiveUserCount 获取活跃用户数
func (r *StatisticsRepositoryImpl) GetActiveUserCount(days int) (int64, error) {
	sevenDaysAgo := time.Now().AddDate(0, 0, -days)
	var count int64
	err := r.db.Model(&models.UserActivity{}).Where("created_at >= ?", sevenDaysAgo).
		Distinct("user_id").
		Count(&count).Error
	return count, err
}

// GetUserDistribution 获取用户分布
func (r *StatisticsRepositoryImpl) GetUserDistribution() ([]interface{}, error) {
	type Distribution struct {
		Role  string `json:"role"`
		Count int    `json:"count"`
	}

	var result []Distribution

	// 超级管理员
	var adminCount int64
	r.db.Model(&models.User{}).Where("role = ?", 1).Count(&adminCount)
	result = append(result, Distribution{Role: "管理员", Count: int(adminCount)})

	// 普通用户
	var userCount int64
	r.db.Model(&models.User{}).Where("role = ?", 2).Count(&userCount)
	result = append(result, Distribution{Role: "普通用户", Count: int(userCount)})

	result2 := make([]interface{}, len(result))
	for i, r := range result {
		result2[i] = r
	}
	return result2, nil
}
