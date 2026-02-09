package models

import "time"

// SubchapterArticle 子章节文章关联模型，实现系列子章节和文章的多对多关系
type SubchapterArticle struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`                                             // 关联记录唯一标识ID
	SubchapterID uint64    `gorm:"not null;uniqueIndex:uk_subchapter_article;index" json:"subchapter_id"`           // 子章节ID
	ArticleID    uint64    `gorm:"not null;uniqueIndex:uk_subchapter_article;index" json:"article_id"`              // 文章ID
	SortOrder    int       `gorm:"not null;default:0" json:"sort_order"`                                            // 文章在子章节中的排序
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`                                                // 关联创建时间

	// 关联关系
	Subchapter *SeriesSubchapter `gorm:"foreignKey:SubchapterID;constraint:OnDelete:CASCADE"` // 关联的子章节，删除子章节时级联删除关联记录
	Article    *Article          `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE"`     // 关联的文章，删除文章时级联删除关联记录
}

func (SubchapterArticle) TableName() string {
	return "subchapter_articles"
}
