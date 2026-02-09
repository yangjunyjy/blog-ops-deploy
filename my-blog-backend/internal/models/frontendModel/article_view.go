package models

import "time"

// ArticleView 文章浏览记录模型，记录文章的浏览情况
type ArticleView struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`                             // 浏览记录唯一标识ID
	ArticleID uint64    `gorm:"not null;index" json:"article_id"`                               // 被浏览文章的ID
	IPAddress string    `gorm:"size:50" json:"ip_address"`                                     // 访问者的IP地址
	UserAgent string    `gorm:"size:500" json:"user_agent"`                                    // 访问者的User-Agent信息
	UserID    *uint64   `gorm:"index" json:"user_id"`                                          // 访问用户的ID，NULL表示匿名用户
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`                              // 浏览时间

	// 关联关系
	Article *Article `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"article,omitempty"` // 被浏览的文章，删除文章时级联删除浏览记录
	User    *User    `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`   // 访问的用户，删除用户时保留浏览记录但置空外键
}

func (ArticleView) TableName() string {
	return "article_view"
}
