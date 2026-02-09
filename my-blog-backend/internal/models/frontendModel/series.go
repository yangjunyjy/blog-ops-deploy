package models

import "time"

// Series 专题系列模型，用于组织相关的文章成为一个系列
type Series struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`                                 // 系列唯一标识ID
	Name        string    `gorm:"size:100;not null" json:"name" validate:"required"`                  // 系列名称，必填
	Slug        string    `gorm:"size:100;not null;uniqueIndex" json:"slug" validate:"required,slug"` // 系列URL别名，用于生成友好的URL，全局唯一
	Icon        string    `gorm:"size:50" json:"icon"`                                                // 系列图标标识
	Description string    `gorm:"type:text" json:"description"`                                       // 系列描述说明
	Cover       string    `gorm:"size:500" json:"cover"`                                              // 系列封面图片URL
	SortOrder   int       `gorm:"not null;default:0" json:"sort_order"`                               // 系列排序字段，数值越小越靠前
	Status      uint8     `gorm:"not null;default:1" json:"status"`                                   // 系列状态：0-禁用，1-启用
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`                                   // 系列创建时间
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`                                   // 系列最后更新时间

	// 关联关系
	Sections []SeriesSection `gorm:"foreignKey:SeriesID" json:"sections,omitempty"` // 该系列包含的所有章节列表，一对多关系
}

func (Series) TableName() string {
	return "series"
}
