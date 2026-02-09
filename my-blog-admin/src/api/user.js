import request from '@/utils/request'
import { mockApi } from '@/mock'

// 是否使用Mock数据
const USE_MOCK = false

/**
 * 获取用户列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.page_size - 每页数量
 * @param {string} params.keyword - 搜索关键词
 * @param {number} params.status - 状态(0-待审核, 1-正常, 2-禁用)
 */
export const getUserList = (params) => {
  if (USE_MOCK) {
    return mockApi.getUserList(params)
  }
  return request.get('/rbac/users',{
    params
  })
}

/**
 * 获取用户详情
 * @param {number} id - 用户ID
 */
export const getUserDetail = (id) => {
  if (USE_MOCK) {
    return mockApi.getUserDetail(id)
  }
  return request.get(`/rbac/users/${id}`)
}

/**
 * 创建用户
 * @param {Object} data - 用户数据
 * @param {string} data.username - 用户名
 * @param {string} data.password - 密码
 * @param {string} data.email - 邮箱
 * @param {string} data.nickname - 昵称
 * @param {number} data.role - 角色(0-普通用户, 1-管理员)
 * @param {number} data.status - 状态(0-待审核, 1-正常, 2-禁用)
 */
export const createUser = (data) => {
  if (USE_MOCK) {
    return mockApi.createUser(data)
  }
  return request.post('/rbac/users/create', data)
}

/**
 * 更新用户
 * @param {number} id - 用户ID
 * @param {Object} data - 用户数据
 */
export const updateUser = (data) => {
  return request.post('/rbac/users/update', data)
}

/**
 * 删除用户
 * @param {number} id - 用户ID
 */
export const deleteUser = (id) => {
  if (USE_MOCK) {
    return mockApi.deleteUser(id)
  }
  return request.delete(`/rbac/users/${id}`)
}

/**
 * 批量删除用户
 * @param {Array<number>} ids - 用户ID列表
 */
export const batchDeleteUsers = (ids) => {
  if (USE_MOCK) {
    return mockApi.batchDeleteUsers(ids)
  }
  return request.delete('/admin/users/batch', { data: { ids } })
}

/**
 * 重置用户密码
 * @param {number} id - 用户ID
 * @param {Object} data - 密码数据
 * @param {string} data.new_password - 新密码
 */
export const resetPassword = (id, data) => {
  if (USE_MOCK) {
    return mockApi.resetPassword(id)
  }
  return request.post(`/admin/users/${id}/reset-password`, data)
}

/**
 * 修改用户状态
 * @param {number} id - 用户ID
 * @param {Object} data - 状态数据
 * @param {number} data.status - 状态(0-待取, 1-正常, 2-禁用)
 */
export const updateUserStatus = (id, data) => {
  if (USE_MOCK) {
    return mockApi.updateUserStatus(id, data)
  }
  return request.put(`/admin/users/${id}/status`, data)
}

/**
 * 分配角色给用户
 * @param {number} userId - 用户ID
 * @param {Object} data - 角色数据
 * @param {Array<number>} roleIds - 角色ID列表
 */
export const assignUserRoles = (data) => {
  return request.post('/rbac/users/assign', data)
}

/**
 * 获取用户角色
 * @param {number} id - 用户ID
 */
export const getUserRoles = (id) => {
  if (USE_MOCK) {
    return mockApi.getUserRoles(id)
  }
  return request.get(`/admin/users/${id}/roles`)
}

/**
 * 批量更新用户状态
 * @param {Array<number>} ids - 用户ID列表
 * @param {number} status - 状态(0-待审核, 1-正常, 2-禁用)
 */
export const batchUpdateUserStatus = (ids, status) => {
  if (USE_MOCK) {
    return mockApi.batchUpdateUserStatus(ids, status)
  }
  return request.put('/admin/users/batch/status', { ids, status })
}
