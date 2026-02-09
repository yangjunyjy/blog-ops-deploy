package models

import (
	"time"

	"gorm.io/gorm"
)

// FavoriteFolder 用户收藏文件夹，用于分类管理收藏的文章
type FavoriteFolder struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint           `gorm:"not null;index:idx_folder_user" json:"user_id"`
	Name        string         `gorm:"size:50;not null" json:"name"`
	Description string         `gorm:"size:200" json:"description"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系（无外键约束，仅用于 ORM 查询）
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (FavoriteFolder) TableName() string {
	return "favorite_folder"
}
