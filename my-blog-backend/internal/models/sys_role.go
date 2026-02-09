package models

var SuperAdminRoleID uint64 = 1

// SysRole 角色模型
type SysRole struct {
	BaseModel
	RoleCode  string `gorm:"column:role_code;type:varchar(50);not null;uniqueIndex:uk_role_code;comment:角色编码" json:"roleCode"`
	RoleName  string `gorm:"column:role_name;type:varchar(50);not null;comment:角色名称" json:"roleName"`
	RoleDesc  string `gorm:"column:role_desc;type:varchar(200);comment:角色描述" json:"roleDesc"`
	DataScope int8   `gorm:"column:data_scope;default:1;comment:数据权限范围 1:全部数据 2:本部门及以下 3:本部门 4:仅本人" json:"dataScope"`
	Status    int8   `gorm:"column:status;default:1;index:idx_status;comment:状态 0:禁用 1:启用" json:"status"`
	Sort      int    `gorm:"column:sort;default:0;comment:排序" json:"sort"`

	// 关联关系
	Users []*SysUser `gorm:"many2many:sys_user_role;joinForeignKey:role_id;joinReferences:user_id" json:"users,omitempty"`
	Menus []*SysMenu `gorm:"many2many:sys_role_menu;joinForeignKey:role_id;joinReferences:menu_id" json:"menus,omitempty"`
}

// TableName 指定表名
func (SysRole) TableName() string {
	return "sys_role"
}
