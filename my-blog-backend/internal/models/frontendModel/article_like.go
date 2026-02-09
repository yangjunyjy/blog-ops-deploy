package models

import "time"

// ArticleLike 文章点赞记录模型，记录用户对文章的点赞行为
type ArticleLike struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`                           // 点赞记录唯一标识ID
	ArticleID uint64    `gorm:"not null;uniqueIndex:uk_article_user;index" json:"article_id"` // 被点赞文章的ID
	UserID    uint64    `gorm:"not null;uniqueIndex:uk_article_user;index" json:"user_id"`    // 点赞用户的ID
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`                             // 点赞时间

	// 关联关系
	Article *Article `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"article,omitempty"` // 被点赞的文章，删除文章时级联删除点赞记录
	User    *User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`       // 点赞的用户，删除用户时级联删除点赞记录
}

func (ArticleLike) TableName() string {
	return "article_likes"
}
