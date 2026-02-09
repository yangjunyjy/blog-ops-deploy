package models

// SysPost 岗位表
type SysPost struct {
	ID        uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PostName  string `gorm:"column:post_name;type:varchar(50);not null;comment:岗位名称" json:"post_name"`
	PostCode  string `gorm:"column:post_code;type:varchar(50);not null;comment:岗位编码" json:"post_code"`
	PostSort  int    `gorm:"column:post_sort;type:int;default:0;comment:显示顺序" json:"post_sort"`
	Status    int8   `gorm:"column:status;type:tinyint;default:1;comment:状态 0:停用 1:正常" json:"status"`
	Remark     string `gorm:"column:remark;type:varchar(500);comment:备注" json:"remark"`

	// 关联关系
	Users []*SysUser `gorm:"many2many:sys_user_post" json:"users,omitempty"`
}

// TableName 指定表名
func (SysPost) TableName() string {
	return "sys_post"
}
