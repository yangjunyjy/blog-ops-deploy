package models

import (
	"time"

	"gorm.io/gorm"
)

// SysUser 用户模型
type SysUser struct {
	BaseModel
	Username      string    `gorm:"column:username;type:varchar(50);not null;uniqueIndex:uk_username;comment:用户名" json:"username"`
	Password      string    `gorm:"column:password;type:varchar(255);not null;comment:密码" json:"-"`
	RealName      string    `gorm:"column:real_name;type:varchar(50);comment:真实姓名" json:"realName"`
	Nickname      string    `gorm:"column:nickname;type:varchar(50);comment:昵称" json:"nickname"`
	Avatar        string    `gorm:"column:avatar;type:varchar(500);comment:头像URL" json:"avatar"`
	Email         string    `gorm:"column:email;type:varchar(100);comment:邮箱" json:"email"`
	Phone         string    `gorm:"column:phone;type:varchar(20);comment:手机号" json:"phone"`
	Gender        int8      `gorm:"column:gender;default:0;comment:性别 0:未知 1:男 2:女" json:"gender"`
	DeptID        uint64    `gorm:"column:dept_id;index:idx_dept_id;comment:部门ID" json:"deptId"`
	Status        int8      `gorm:"column:status;default:1;index:idx_status;comment:状态 0:禁用 1:启用" json:"status"`
	IsAdmin       int8      `gorm:"column:is_admin;default:0;comment:是否超级管理员 0:否 1:是" json:"isAdmin"`
	LastLoginTime *time.Time `gorm:"column:last_login_time;comment:最后登录时间" json:"lastLoginTime"`
	LastLoginIP   string    `gorm:"column:last_login_ip;type:varchar(50);comment:最后登录IP" json:"lastLoginIp"`

	// 关联关系（虚拟外键）
	Dept  *SysDept   `gorm:"foreignKey:DeptID;references:ID" json:"dept,omitempty"`
	Roles []*SysRole `gorm:"many2many:sys_user_role;joinForeignKey:user_id;joinReferences:role_id" json:"roles,omitempty"`
	Posts []*SysPost `gorm:"many2many:sys_user_post;joinForeignKey:user_id;joinReferences:post_id" json:"posts,omitempty"`
}

// TableName 指定表名
func (SysUser) TableName() string {
	return "sys_user"
}

// BeforeCreate 创建前钩子
func (u *SysUser) BeforeCreate(tx *gorm.DB) error {
	if u.Gender == 0 {
		u.Gender = 0
	}
	if u.Status == 0 {
		u.Status = 1
	}
	if u.IsAdmin == 0 {
		u.IsAdmin = 0
	}
	return nil
}
