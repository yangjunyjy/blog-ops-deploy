package models

import "time"

// Category 文章分类模型，用于对文章进行分类组织
type Category struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`                                 // 分类唯一标识ID
	Name         string    `gorm:"size:50;not null" json:"name" validate:"required"`                   // 分类名称，必填
	Slug         string    `gorm:"size:100;not null;uniqueIndex" json:"slug" validate:"required,slug"` // 分类URL别名，用于生成友好的URL，全局唯一
	Description  string    `gorm:"size:200" json:"description"`                                        // 分类描述说明
	Icon         string    `gorm:"size:50" json:"icon"`                                               // 分类图标标识
	SortOrder    int       `gorm:"not null;default:0" json:"sort_order"`                              // 分类排序字段，数值越小越靠前
	Status       uint8     `gorm:"not null;default:1" json:"status"`                                  // 分类状态：0-禁用，1-启用
	ArticleCount int64     `gorm:"-" json:"article_count"`                                           // 该分类下的文章数量统计，不从数据库加载，运行时计算
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`                                  // 分类创建时间
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`                                  // 分类最后更新时间

	// 关联关系
	Articles []Article `gorm:"foreignKey:CategoryID" json:"articles,omitempty"` // 该分类下的所有文章列表，一对多关系
}

func (Category) TableName() string {
	return "categorie"
}
