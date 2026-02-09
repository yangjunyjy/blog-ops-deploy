package models

import "gorm.io/gorm"

// SysMenu 系统菜单表
type SysMenu struct {
	BaseModel
	ParentID   uint64   `gorm:"column:parent_id;type:bigint;default:0;comment:父菜单ID，0为一级菜单" json:"parent_id"`
	MenuName   string   `gorm:"column:menu_name;type:varchar(50);not null;comment:菜单名称" json:"menu_name"`
	MenuCode   string   `gorm:"column:menu_code;type:varchar(100);not null;uniqueIndex;comment:菜单编码，权限标识" json:"menu_code"`
	MenuType   MenuType `gorm:"column:menu_type;type:tinyint;not null;comment:菜单类型:1目录 2菜单 3按钮" json:"menu_type"`
	Icon       string   `gorm:"column:icon;type:varchar(100);comment:菜单图标" json:"icon"`
	Sort       int      `gorm:"column:sort;type:int;default:0;comment:显示顺序" json:"sort"`
	Path       string   `gorm:"column:path;type:varchar(200);comment:路由地址" json:"path"`
	Component  string   `gorm:"column:component;type:varchar(200);comment:组件路径" json:"component"`
	Permission string   `gorm:"column:perms;type:varchar(500);comment:权限标识" json:"permission"`
	IsVisible  bool     `gorm:"column:is_visible;type:tinyint(1);default:1;comment:是否显示:1显示 0隐藏" json:"is_visible"`
	Status     Status   `gorm:"column:status;type:tinyint;not null;default:1;comment:状态:1正常 2停用" json:"status"`
	Remark     string   `gorm:"column:remark;type:varchar(500);comment:备注" json:"remark"`

	// 关联关系
	Children []*SysMenu `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Roles    []*SysRole `gorm:"many2many:sys_role_menu" json:"roles,omitempty"`
}

// TableName 指定表名
func (SysMenu) TableName() string {
	return "sys_menu"
}

// AfterCreate 创建后钩子
func (m *SysMenu) AfterCreate(tx *gorm.DB) error {
	return nil
}
