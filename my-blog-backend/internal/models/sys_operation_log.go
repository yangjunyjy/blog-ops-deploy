package models

// SysOperationLog 操作日志表
type SysOperationLog struct {
	BaseModel
	Module      string `gorm:"column:module;type:varchar(50);comment:模块标题" json:"module"`
	Business    string `gorm:"column:business;type:varchar(50);comment:业务类型" json:"business"`
	Method      string `gorm:"column:method;type:varchar(100);comment:方法名称" json:"method"`
	Request     string `gorm:"column:request;type:text;comment:请求方式" json:"request"`
	Params      string `gorm:"column:params;type:text;comment:请求参数" json:"params"`
	IP          string `gorm:"column:ip;type:varchar(128);comment:操作IP" json:"ip"`
	Location    string `gorm:"column:location;type:varchar(255);comment:操作地点" json:"location"`
	Agent       string `gorm:"column:agent;type:varchar(500);comment:用户代理" json:"agent"`
	Operator    string `gorm:"column:operator;type:varchar(50);comment:操作人员" json:"operator"`
	Status      int    `gorm:"column:status;type:int;default:0;comment:操作状态:0成功 1失败" json:"status"`
	ErrorMsg    string `gorm:"column:error_msg;type:text;comment:错误消息" json:"error_msg"`
	CostTime    int64  `gorm:"column:cost_time;type:bigint;default:0;comment:消耗时间(毫秒)" json:"cost_time"`
}

// TableName 指定表名
func (SysOperationLog) TableName() string {
	return "sys_operation_log"
}
