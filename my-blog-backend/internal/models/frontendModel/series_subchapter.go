package models

import "time"

// SeriesSubchapter 系列子章节模型，一个章节包含多个子章节
type SeriesSubchapter struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`                                 // 子章节唯一标识ID
	SectionID   uint64    `gorm:"not null;index" json:"section_id"`                                   // 所属章节的ID
	Name        string    `gorm:"size:100;not null" json:"name" validate:"required"`                 // 子章节名称，必填
	Description string    `gorm:"type:text" json:"description"`                                      // 子章节描述说明
	SortOrder   int       `gorm:"not null;default:0" json:"sort_order"`                              // 子章节排序字段，数值越小越靠前
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`                                  // 子章节创建时间
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`                                  // 子章节最后更新时间

	// 关联关系
	Section  *SeriesSection `gorm:"foreignKey:SectionID;constraint:OnDelete:CASCADE" json:"section,omitempty"` // 所属的章节，删除章节时级联删除子章节
	Articles []Article      `gorm:"many2many:subchapter_articles;foreignKey:ID;joinForeignKey:SubchapterID;joinReferences:ArticleID" json:"articles,omitempty"` // 该子章节包含的所有文章列表，多对多关系
}

func (SeriesSubchapter) TableName() string {
	return "series_subchapters"
}
