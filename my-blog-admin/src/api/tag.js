import request from '@/utils/request'

/**
 * 获取标签列表（后台管理）
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.page_size - 每页数量
 */
export const getTagList = (params) => {
  return request.get('/public/tags', { params })
}

/**
 * 获取标签详情
 * @param {number} id - 标签ID
 */
export const getTagDetail = (id) => {
  return request.get(`/public/tags/${id}`)
}

/**
 * 创建标签
 * @param {Object} data - 标签数据
 * @param {string} data.name - 标签名称
 * @param {string} data.slug - URL别名
 * @param {string} data.description - 描述
 */
export const createTag = (data) => {
  return request.post('/rbac/tags', data)
}

/**
 * 更新标签
 * @param {number} id - 标签ID
 * @param {Object} data - 标签数据
 */
export const updateTag = (id, data) => {
  return request.put(`/rbac/tags/${id}`, data)
}

/**
 * 删除标签
 * @param {number} id - 标签ID
 */
export const deleteTag = (id) => {
  return request.delete(`/rbac/tags/${id}`)
}
