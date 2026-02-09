import request from '@/utils/request'

/**
 * 获取角色列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.pageSize - 每页数量
 * @param {string} params.keyword - 搜索关键词
 * @param {number} params.status - 状态(0-禁用, 1-启用)
 */
export function getRoleList(params) {
  return request({
    url: '/rbac/roles',
    method: 'get',
    params
  })
}

/**
 * 获取角色详情
 * @param {number} id - 角色ID
 */
export function getRoleDetail(id) {
  return request({
    url: `/admin/roles/${id}`,
    method: 'get'
  })
}

/**
 * 创建角色
 * @param {Object} data - 角色数据
 * @param {string} data.name - 角色名称
 * @param {string} data.code - 角色代码
 * @param {string} data.description - 角色描述
 * @param {number} data.status - 状态(0-禁用, 1-启用)
 */
export function createRole(data) {
  return request({
    url: '/rbac/roles/create',
    method: 'post',
    data
  })
}

/**
 * 更新角色
 * @param {number} id - 角色ID
 * @param {Object} data - 角色数据
 */
export function updateRole(data) {
  return request({
    url: '/rbac/roles/update',
    method: 'post',
    data
  })
}

/**
 * 删除角色
 * @param {number} id - 角色ID
 */
export function deleteRole(id) {
  return request({
    url: `/rbac/roles/${id}`,
    method: 'delete'
  })
}

/**
 * 分配菜单给角色
 * @param {number} id - 角色ID
 * @param {Object} data - 权限数据
 * @param {Array<number>} data.menu_ids - 菜单ID列表
 */
export function assignRoleMenus(data) {
  return request({
    url: `/rbac/roles/assign`,
    method: 'post',
    data
  })
}

/**
 * 批量删除角色
 * @param {Array<number>} ids - 角色ID列表
 */
export function batchDeleteRoles(ids) {
  return request({
    url: '/admin/roles/batch',
    method: 'delete',
    data: { ids }
  })
}

export function GetMenusByRoleID(id){
  return request({
    url: `/rbac/roles/${id}`,
    method:'get'
  })
}
