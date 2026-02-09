package models

import "time"

// CommentLike 评论点赞模型，记录用户对评论的点赞
type CommentLike struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`                                    // 点赞记录唯一标识ID
	CommentID uint      `gorm:"not null;index:idx_comment_id,index:idx_comment_user" json:"comment_id"` // 被点赞评论的ID
	UserID    uint      `gorm:"not null;index:idx_user_id,index:idx_comment_user" json:"user_id"`       // 点赞用户的ID
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`                                      // 点赞时间
}

func (CommentLike) TableName() string {
	return "comment_like"
}
