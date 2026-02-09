package models

import "time"

// SeriesSection 系列章节模型，一个系列包含多个章节
type SeriesSection struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`                                 // 章节唯一标识ID
	SeriesID    uint64    `gorm:"not null;index" json:"series_id"`                                   // 所属系列的ID
	Name        string    `gorm:"size:100;not null" json:"name" validate:"required"`                 // 章节名称，必填
	Description string    `gorm:"type:text" json:"description"`                                      // 章节描述说明
	SortOrder   int       `gorm:"not null;default:0" json:"sort_order"`                              // 章节排序字段，数值越小越靠前
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`                                  // 章节创建时间
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`                                  // 章节最后更新时间

	// 关联关系
	Series      *Series            `gorm:"foreignKey:SeriesID;constraint:OnDelete:CASCADE" json:"series,omitempty"`       // 所属的系列，删除系列时级联删除章节
	Subchapters []SeriesSubchapter `gorm:"foreignKey:SectionID" json:"subchapters,omitempty"` // 该章节包含的所有子章节列表，一对多关系
}

func (SeriesSection) TableName() string {
	return "series_sections"
}
