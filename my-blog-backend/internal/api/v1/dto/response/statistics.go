package response

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	Overview struct {
		Articles int `json:"articles"`
		Views    int `json:"views"`
		Likes    int `json:"likes"`
		Comments int `json:"comments"`
	} `json:"overview"`
	Trend struct {
		Views    []int `json:"views"`
		Articles []int `json:"articles"`
	} `json:"trend"`
	Content struct {
		Categories int `json:"categories"`
		Tags       int `json:"tags"`
		Series     int `json:"series"`
	} `json:"content"`
}

// ArticleStats 文章统计
type ArticleStats struct {
	Total      int    `json:"total"`
	Published  int    `json:"published"`
	Draft      int    `json:"draft"`
	TotalViews int    `json:"totalViews"`
	Trend      []int  `json:"trend"`
}

// UserStats 用户统计
type UserStats struct {
	TotalUsers  int `json:"totalUsers"`
	ActiveUsers int `json:"activeUsers"`
	NewUsers    int `json:"newUsers"`
	OnlineUsers int `json:"onlineUsers"`
}

// CategoryStat 分类统计
type CategoryStat struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ArticleCount int    `json:"articleCount"`
}

// TagStat 标签统计
type TagStat struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ArticleCount int    `json:"articleCount"`
}

// HotArticle 热门文章
type HotArticle struct {
	ID           uint   `json:"id"`
	Title        string `json:"title"`
	ViewCount    int    `json:"views"`
	LikeCount    int    `json:"likes"`
	CommentCount int    `json:"commentCount"`
}

// GrowthData 增长数据
type GrowthData struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// ActiveUserData 活跃用户数据
type ActiveUserData struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// Distribution 分布数据
type Distribution struct {
	Role  string `json:"role"`
	Count int    `json:"count"`
}
