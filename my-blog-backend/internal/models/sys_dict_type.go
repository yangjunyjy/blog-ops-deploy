package models

// SysDictType 字典类型表
type SysDictType struct {
	BaseModel
	Name    string `gorm:"column:name;type:varchar(100);not null;uniqueIndex;comment:字典名称" json:"name"`
	Type    string `gorm:"column:type;type:varchar(100);not null;uniqueIndex;comment:字典类型" json:"type"`
	Status  Status `gorm:"column:status;type:tinyint;not null;default:1;comment:状态:1正常 2停用" json:"status"`
	Remark  string `gorm:"column:remark;type:varchar(500);comment:备注" json:"remark"`

	// 关联关系
	DictData []*SysDictData `gorm:"foreignKey:DictTypeID" json:"dict_data,omitempty"`
}

// TableName 指定表名
func (SysDictType) TableName() string {
	return "sys_dict_type"
}
