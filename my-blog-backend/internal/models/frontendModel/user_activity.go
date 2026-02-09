package models

import "time"

// UserActivity 用户活动记录
type UserActivity struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`                          // 用户ID
	ArticleID    uint      `gorm:"not null;index" json:"article_id"`                         // 文章ID
	Type         string    `gorm:"size:20;not null;index" json:"type"`                      // 活动类型：like(点赞)、comment(评论)、share(分享)、favorite(收藏)
	Content      string    `gorm:"type:text" json:"content"`                                 // 评论内容（仅评论类型）
	Platform     string    `gorm:"size:50" json:"platform"`                                  // 分享平台（仅分享类型）
	FolderName   string    `gorm:"size:100" json:"folder_name"`                              // 收藏文件夹名称（仅收藏类型)
	ArticleTitle string    `gorm:"size:200;not null" json:"article_title"`                   // 文章标题（冗余字段，便于查询）
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`                         // 创建时间
}

func (UserActivity) TableName() string {
	return "user_activities"
}
