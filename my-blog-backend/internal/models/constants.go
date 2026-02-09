package models

// 性别枚举
const (
	GenderUnknown = iota // 未知
	GenderMale           // 男
	GenderFemale         // 女
)

// 菜单类型枚举
const (
	MenuTypeDir  = 1 // 目录
	MenuTypeMenu = 2 // 菜单
	MenuTypeBtn  = 3 // 按钮
)

// 状态枚举
const (
	StatusDisabled = 0 // 禁用/停用
	StatusEnabled  = 1 // 启用/正常
)

// 删除标记
const (
	DeletedNo  = 0 // 未删除
	DeletedYes = 1 // 已删除
)

// 数据权限范围
const (
	DataScopeAll        = 1 // 全部数据
	DataScopeDeptAndSub = 2 // 本部门及以下
	DataScopeDept       = 3 // 本部门
	DataScopeSelf       = 4 // 仅本人
)

// 是否超级管理员
const (
	AdminNo  = 0 // 否
	AdminYes = 1 // 是
)

// 是否系统字典
const (
	SystemDictNo  = 0 // 否
	SystemDictYes = 1 // 是
)

// 字典数据来源
const (
	DictSourceSystem = 1 // 系统内置
	DictSourceCustom = 2 // 用户自定义
)

// 业务操作类型
const (
	BusinessTypeAdd    = 1 // 新增
	BusinessTypeEdit   = 2 // 修改
	BusinessTypeDelete = 3 // 删除
	BusinessTypeAuth   = 4 // 授权
	BusinessTypeExport = 5 // 导出
	BusinessTypeImport = 6 // 导入
)

// 登录状态
const (
	LoginSuccess = 0 // 成功
	LoginFailed  = 1 // 失败
)

// 操作状态
const (
	OperationSuccess = 0 // 成功
	OperationFailed  = 1 // 失败
)
