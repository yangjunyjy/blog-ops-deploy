package models

import "time"

// ArticleFavorite 文章收藏记录模型，记录用户对文章的收藏行为
type ArticleFavorite struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`                                                // 收藏记录唯一标识ID
	ArticleID uint      `gorm:"not null;uniqueIndex:uk_article_user;index:idx_favorite_article" json:"article_id"` // 被收藏文章的ID
	UserID    uint      `gorm:"not null;uniqueIndex:uk_article_user;index:idx_favorite_user" json:"user_id"`       // 收藏用户的ID
	FolderID  uint      `gorm:"not null;index:idx_favorite_folder" json:"folder_id"`                                  // 所属文件夹ID
	Note      string    `gorm:"type:text" json:"note"`                                                             // 收藏备注
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`                                                  // 收藏时间

	// 关联关系（无外键约束，仅用于 ORM 查询）
	Article *Article        `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
	User    *User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Folder  *FavoriteFolder `gorm:"foreignKey:FolderID" json:"folder,omitempty"`
}

func (ArticleFavorite) TableName() string {
	return "article_favorites"
}
