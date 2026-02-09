package models

// SysUserPost 用户岗位关联表
type SysUserPost struct {
	ID     uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID uint64 `gorm:"column:user_id;type:bigint;not null;index:idx_user_id,priority:1;comment:用户ID" json:"user_id"`
	PostID uint64 `gorm:"column:post_id;type:bigint;not null;index:idx_post_id,priority:2;comment:岗位ID" json:"post_id"`
}

// TableName 指定表名
func (SysUserPost) TableName() string {
	return "sys_user_post"
}
