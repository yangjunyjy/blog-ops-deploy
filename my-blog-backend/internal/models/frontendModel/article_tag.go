package models

import "time"

// ArticleTag 文章标签关联模型，实现文章和标签的多对多关系
type ArticleTag struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`                                 // 关联记录唯一标识ID
	ArticleID uint64    `gorm:"not null;uniqueIndex:uk_article_tag;index" json:"article_id"`       // 文章ID
	TagID     uint64    `gorm:"not null;uniqueIndex:uk_article_tag;index" json:"tag_id"`           // 标签ID
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`                                  // 关联创建时间

	// 关联关系
	Article Article `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE"` // 关联的文章，删除文章时级联删除关联记录
	Tag     Tag     `gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE"`     // 关联的标签，删除标签时级联删除关联记录
}

func (ArticleTag) TableName() string {
	return "article_tag"
}
