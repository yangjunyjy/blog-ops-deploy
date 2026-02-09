package models

// SysDept 系统部门表
type SysDept struct {
	BaseModel
	ParentID  uint64 `gorm:"column:parent_id;type:bigint;default:0;comment:父部门ID，0为顶级部门" json:"parent_id"`
	Ancestors string `gorm:"column:ancestors;type:varchar(500);default:0;comment:祖级列表，如:0,1,2" json:"ancestors"`
	Name      string `gorm:"column:name;type:varchar(50);not null;uniqueIndex;comment:部门名称" json:"name"`
	Sort      int    `gorm:"column:sort;type:int;default:0;comment:显示顺序" json:"sort"`
	Leader    string `gorm:"column:leader;type:varchar(50);comment:负责人" json:"leader"`
	Phone     string `gorm:"column:phone;type:varchar(20);comment:联系电话" json:"phone"`
	Email     string `gorm:"column:email;type:varchar(50);comment:邮箱" json:"email"`
	Status    Status `gorm:"column:status;type:tinyint;not null;default:1;comment:状态:1正常 2停用" json:"status"`
	Remark    string `gorm:"column:remark;type:varchar(500);comment:备注" json:"remark"`

	// 关联关系
	Children []*SysDept `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Users    []*SysUser `gorm:"foreignKey:DeptID" json:"users,omitempty"`
}

// TableName 指定表名
func (SysDept) TableName() string {
	return "sys_dept"
}
