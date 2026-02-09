package models

// SysDictData 字典数据表
type SysDictData struct {
	BaseModel
	DictTypeID uint64 `gorm:"column:dict_type_id;type:bigint;not null;index;comment:字典类型ID" json:"dict_type_id"`
	Label      string `gorm:"column:label;type:varchar(100);not null;comment:字典标签" json:"label"`
	Value      string `gorm:"column:value;type:varchar(100);not null;comment:字典键值" json:"value"`
	Sort       int    `gorm:"column:sort;type:int;default:0;comment:字典排序" json:"sort"`
	Status     Status `gorm:"column:status;type:tinyint;not null;default:1;comment:状态:1正常 2停用" json:"status"`
	Remark     string `gorm:"column:remark;type:varchar(500);comment:备注" json:"remark"`

	// 关联关系
	DictType *SysDictType `gorm:"foreignKey:DictTypeID" json:"dict_type,omitempty"`
}

// TableName 指定表名
func (SysDictData) TableName() string {
	return "sys_dict_data"
}
