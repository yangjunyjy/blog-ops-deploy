import request from '@/utils/request'

/**
 * 获取菜单列表
 * @param {Object} params - 查询参数
 * @param {string} params.keyword - 搜索关键词
 * @param {number} params.status - 状态(0-禁用, 1-启用)
 */
export function getMenuList(params) {
  return request({
    url: '/admin/menus',
    method: 'get',
    params
  })
}

/**
 * 获取菜单树
 */
export function getMenuTree() {
  return request({
    url: '/rbac/menus/tree',
    method: 'get'
  })
}

/**
 * 获取菜单详情
 * @param {number} id - 菜单ID
 */
export function getMenuDetail(id) {
  return request({
    url: `/admin/menus/${id}`,
    method: 'get'
  })
}

/**
 * 创建菜单
 * @param {Object} data - 菜单数据
 * @param {number} data.parent_id - 父菜单ID
 * @param {string} data.title - 菜单标题
 * @param {string} data.icon - 菜单图标
 * @param {string} data.path - 路由路径
 * @param {string} data.component - 组件路径
 * @param {number} data.type - 类型(0-目录, 1-菜单, 2-按钮)
 * @param {string} data.permission - 权限标识
 * @param {number} data.sort - 排序
 * @param {number} data.status - 状态(0-禁用, 1-启用)
 */
export function createMenu(data) {
  return request({
    url: '/rbac/menus/create',
    method: 'post',
    data
  })
}

/**
 * 更新菜单
 * @param {number} id - 菜单ID
 * @param {Object} data - 菜单数据
 */
export function updateMenu(data) {
  return request({
    url: '/rbac/menus/update',
    method: 'post',
    data
  })
}

/**
 * 删除菜单
 * @param {number} id - 菜单ID
 */
export function deleteMenu(id) {
  return request({
    url: `/rbac/menus/${id}`,
    method: 'delete'
  })
}

/**
 * 批量删除菜单
 * @param {Array<number>} ids - 菜单ID列表
 */
export function batchDeleteMenus(ids) {
  return request({
    url: '/admin/menus/batch',
    method: 'delete',
    data: { ids }
  })
}

export function getAllMenus(){
  return request({
    url:'/rbac/menus',
    method:'get'
  })
}