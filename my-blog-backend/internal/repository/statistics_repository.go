package repository

// StatisticsRepository 统计数据仓储接口
type StatisticsRepository interface {
	// 文章统计
	GetArticleCount(status *int) (int64, error)
	GetTotalViews() (int64, error)
	GetTotalLikes() (int64, error)
	GetCommentCount(status *int) (int64, error)
	GetCategoryCount() (int64, error)
	GetTagCount() (int64, error)
	GetSeriesCount() (int64, error)

	// 趋势数据
	GetArticleTrend(days int) ([]int, error)
	GetViewsTrend(days int) (int64, error)
	GetViewsTrendByDay(days int) ([]int, error)
	GetUserGrowthTrend(days int) ([]string, []int, error)
	GetActiveUserTrend(days int) ([]string, []int, error)

	// 热门数据
	GetHotArticles(limit int) ([]interface{}, error)
	GetHotTags(limit int) ([]interface{}, error)

	// 分类和标签统计
	GetCategoryStats() ([]interface{}, error)
	GetTagStats() ([]interface{}, error)

	// 用户统计
	GetUserCount() (int64, error)
	GetNewUserCount(date string) (int64, error)
	GetActiveUserCount(days int) (int64, error)
	GetUserDistribution() ([]interface{}, error)
}
