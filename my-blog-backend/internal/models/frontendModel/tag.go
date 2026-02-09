package models

import "time"

// Tag 文章标签模型，用于对文章进行标签标记和分类
type Tag struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`                                 // 标签唯一标识ID
	Name         string    `gorm:"size:50;not null;uniqueIndex" json:"name" validate:"required"`       // 标签名称，必填，全局唯一
	Slug         string    `gorm:"size:100;not null;uniqueIndex" json:"slug" validate:"required,slug"` // 标签URL别名，用于生成友好的URL，全局唯一
	Description  string    `gorm:"size:200" json:"description"`                                        // 标签描述说明
	ArticleCount int64     `gorm:"-" json:"article_count"`                                           // 该标签下的文章数量统计，不从数据库加载，运行时计算
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`                                  // 标签创建时间

	// 关联关系
	Articles []Article `gorm:"many2many:article_tags;" json:"articles,omitempty"` // 使用该标签的所有文章列表，多对多关系
}

func (Tag) TableName() string {
	return "tag"
}
