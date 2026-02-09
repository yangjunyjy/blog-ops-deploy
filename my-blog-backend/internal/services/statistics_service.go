package services

import (
	"my-blog-backend/internal/api/v1/dto/response"
	"my-blog-backend/internal/repository"

	"time"
)

// StatisticsService 统计服务
type StatisticsService struct {
	statsRepo repository.StatisticsRepository
}

// NewStatisticsService 创建统计服务
func NewStatisticsService(statsRepo repository.StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		statsRepo: statsRepo,
	}
}

// GetDashboardStats 获取仪表盘统计数据
func (s *StatisticsService) GetDashboardStats() (*response.DashboardStats, error) {
	stats := &response.DashboardStats{}

	// 获取文章总数
	articleCount, _ := s.statsRepo.GetArticleCount(nil)
	stats.Overview.Articles = int(articleCount)

	// 获取总浏览量
	totalViews, _ := s.statsRepo.GetTotalViews()
	stats.Overview.Views = int(totalViews)

	// 获取总点赞数
	totalLikes, _ := s.statsRepo.GetTotalLikes()
	stats.Overview.Likes = int(totalLikes)

	// 获取总评论数
	commentCount, _ := s.statsRepo.GetCommentCount(nil)
	stats.Overview.Comments = int(commentCount)

	// 获取分类、标签、系列数量
	categoryCount, _ := s.statsRepo.GetCategoryCount()
	tagCount, _ := s.statsRepo.GetTagCount()
	seriesCount, _ := s.statsRepo.GetSeriesCount()
	stats.Content.Categories = int(categoryCount)
	stats.Content.Tags = int(tagCount)
	stats.Content.Series = int(seriesCount)

	// 获取趋势数据
	viewsTrend := s.getViewsTrend(7)
	articlesTrend, _ := s.statsRepo.GetArticleTrend(7)
	stats.Trend.Views = viewsTrend
	stats.Trend.Articles = articlesTrend

	return stats, nil
}

// GetArticleStats 获取文章统计数据
func (s *StatisticsService) GetArticleStats() (*response.ArticleStats, error) {
	status1 := 1
	status0 := 0

	total, _ := s.statsRepo.GetArticleCount(nil)
	published, _ := s.statsRepo.GetArticleCount(&status1)
	draft, _ := s.statsRepo.GetArticleCount(&status0)
	totalViews, _ := s.statsRepo.GetTotalViews()
	trend, _ := s.statsRepo.GetArticleTrend(7)

	stats := &response.ArticleStats{
		Total:      int(total),
		Published:  int(published),
		Draft:      int(draft),
		TotalViews: int(totalViews),
		Trend:      trend,
	}
	return stats, nil
}

// GetArticleTrend 获取文章发布趋势
func (s *StatisticsService) GetArticleTrend(days int) ([]int, error) {
	return s.statsRepo.GetArticleTrend(days)
}

// GetViewTrend 获取浏览量趋势
func (s *StatisticsService) GetViewTrend(days int) ([]int, error) {
	viewsTrend := s.getViewsTrend(days)
	return viewsTrend, nil
}

// GetHotArticles 获取热门文章
func (s *StatisticsService) GetHotArticles(limit int) ([]interface{}, error) {
	return s.statsRepo.GetHotArticles(limit)
}

// GetCategoryStats 获取分类统计
func (s *StatisticsService) GetCategoryStats() ([]interface{}, error) {
	return s.statsRepo.GetCategoryStats()
}

// GetTagStats 获取标签统计
func (s *StatisticsService) GetTagStats() ([]interface{}, error) {
	return s.statsRepo.GetTagStats()
}

// GetHotTags 获取热门标签
func (s *StatisticsService) GetHotTags(limit int) ([]interface{}, error) {
	return s.statsRepo.GetHotTags(limit)
}

// GetUserStats 获取用户统计
func (s *StatisticsService) GetUserStats() (*response.UserStats, error) {
	totalUsers, _ := s.statsRepo.GetUserCount()

	today := time.Now().Format("2006-01-02")
	newUsers, _ := s.statsRepo.GetNewUserCount(today)

	activeUsers, _ := s.statsRepo.GetActiveUserCount(7)

	stats := &response.UserStats{
		TotalUsers:  int(totalUsers),
		ActiveUsers: int(activeUsers),
		NewUsers:    int(newUsers),
		OnlineUsers: 0, // 需要实现在线用户功能
	}
	return stats, nil
}

// GetUserGrowth 获取用户增长趋势
func (s *StatisticsService) GetUserGrowth(days int) ([]response.GrowthData, error) {
	dates, counts, err := s.statsRepo.GetUserGrowthTrend(days)
	if err != nil {
		return nil, err
	}

	result := make([]response.GrowthData, len(dates))
	for i, date := range dates {
		result[i] = response.GrowthData{
			Date:  date,
			Count: counts[i],
		}
	}
	return result, nil
}

// GetActiveUsers 获取活跃用户趋势
func (s *StatisticsService) GetActiveUsers(days int) ([]response.ActiveUserData, error) {
	dates, counts, err := s.statsRepo.GetActiveUserTrend(days)
	if err != nil {
		return nil, err
	}

	result := make([]response.ActiveUserData, len(dates))
	for i, date := range dates {
		result[i] = response.ActiveUserData{
			Date:  date,
			Count: counts[i],
		}
	}
	return result, nil
}

// GetUserDistribution 获取用户分布
func (s *StatisticsService) GetUserDistribution() ([]interface{}, error) {
	return s.statsRepo.GetUserDistribution()
}

// GetOnlineUsers 获取在线用户数
func (s *StatisticsService) GetOnlineUsers() (int, error) {
	// 需要实现在线用户功能，暂时返回0
	return 0, nil
}

// getViewsTrend 获取浏览量趋势数据
func (s *StatisticsService) getViewsTrend(days int) []int {
	trend, err := s.statsRepo.GetViewsTrendByDay(days)
	if err != nil {
		// 如果查询失败，返回空数组
		return []int{}
	}
	return trend
}
