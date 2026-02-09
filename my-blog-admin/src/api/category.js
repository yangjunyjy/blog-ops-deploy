import request from '@/utils/request'

/**
 * 获取分类列表（后台管理）
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.page_size - 每页数量
 */
export const getCategoryList = (params) => {
  return request.get('/public/categories', { params })
}

/**
 * 获取分类树
 */
export const getCategoryTree = () => {
  return request.get('/public/categories/tree')
}

/**
 * 获取分类详情
 * @param {number} id - 分类ID
 */
export const getCategoryDetail = (id) => {
  return request.get(`/public/categories/${id}`)
}

/**
 * 创建分类
 * @param {Object} data - 分类数据
 * @param {string} data.name - 分类名称
 * @param {string} data.slug - URL别名
 * @param {string} data.description - 描述
 * @param {string} data.icon - 图标
 * @param {number} data.sort_order - 排序
 */
export const createCategory = (data) => {
  return request.post('/rbac/categories', data)
}

/**
 * 更新分类
 * @param {number} id - 分类ID
 * @param {Object} data - 分类数据
 */
export const updateCategory = (id, data) => {
  return request.put(`/rbac/categories/${id}`, data)
}

/**
 * 删除分类
 * @param {number} id - 分类ID
 */
export const deleteCategory = (id) => {
  return request.delete(`/rbac/categories/${id}`)
}
