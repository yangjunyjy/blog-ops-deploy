import request from '@/utils/request'

/**
 * 获取系列列表（后台管理）
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.page_size - 每页数量
 * @param {string} params.keyword - 搜索关键词
 */
export const getSeriesList = (params) => {
  return request.get('/rbac/series', { params })
}

/**
 * 获取系列详情
 * @param {number} id - 系列ID
 */
export const getSeriesDetail = (id) => {
  return request.get(`/public/series/${id}`)
}

/**
 * 创建系列
 * @param {Object} data - 系列数据
 * @param {string} data.name - 系列名称
 * @param {string} data.slug - URL别名
 * @param {string} data.icon - 图标
 * @param {string} data.description - 描述
 * @param {string} data.cover - 封面URL
 * @param {number} data.sort_order - 排序
 * @param {number} data.status - 状态(0-禁用, 1-启用)
 */
export const createSeries = (data) => {
  return request.post('/rbac/series', data)
}

/**
 * 更新系列
 * @param {number} id - 系列ID
 * @param {Object} data - 系列数据
 */
export const updateSeries = (id, data) => {
  return request.put(`/rbac/series/${id}`, data)
}

/**
 * 删除系列
 * @param {number} id - 系列ID
 */
export const deleteSeries = (id) => {
  return request.delete(`/rbac/series/${id}`)
}

/**
 * 获取系列章节列表
 * @param {number} id - 系列ID
 * @param {Object} params - 查询参数
 */
export const getSeriesSections = (id, params) => {
  return request.get(`/public/series/${id}/sections`, { params })
}

/**
 * 获取章节详情
 * @param {number} id - 章节ID
 */
export const getSectionDetail = (id) => {
  return request.get(`/public/series/sections/${id}`)
}

/**
 * 创建章节
 * @param {number} seriesId - 系列ID
 * @param {Object} data - 章节数据
 * @param {string} data.name - 章节名称
 * @param {string} data.description - 章节描述
 * @param {number} data.sort_order - 排序
 */
export const createSection = (seriesId, data) => {
  return request.post(`/rbac/series/${seriesId}/sections`, data)
}

/**
 * 更新章节
 * @param {number} id - 章节ID
 * @param {Object} data - 章节数据
 */
export const updateSection = (id, data) => {
  return request.put(`/rbac/series/sections/${id}`, data)
}

/**
 * 删除章节
 * @param {number} id - 章节ID
 */
export const deleteSection = (id) => {
  return request.delete(`/rbac/series/sections/${id}`)
}

/**
 * 获取子章节列表
 * @param {number} id - 章节ID
 * @param {Object} params - 查询参数
 */
export const getSectionSubchapters = (id, params) => {
  return request.get(`/public/series/sections/${id}/subchapters`, { params })
}

/**
 * 获取子章节详情
 * @param {number} id - 子章节ID
 */
export const getSubchapterDetail = (id) => {
  return request.get(`/public/series/subchapters/${id}`)
}

/**
 * 创建子章节
 * @param {Object} data - 子章节数据
 * @param {string} data.name - 子章节名称
 * @param {string} data.description - 子章节描述
 * @param {number} data.sort_order - 排序
 */
export const createSubchapter = (data) => {
  return request.post('/rbac/series/subchapters', data)
}

/**
 * 更新子章节
 * @param {number} id - 子章节ID
 * @param {Object} data - 子章节数据
 */
export const updateSubchapter = (id, data) => {
  return request.put(`/rbac/series/subchapters/${id}`, data)
}

/**
 * 删除子章节
 * @param {number} id - 子章节ID
 */
export const deleteSubchapter = (id) => {
  return request.delete(`/rbac/series/subchapters/${id}`)
}

/**
 * 获取子章节文章列表
 * @param {number} id - 子章节ID
 */
export const getSubchapterArticles = (id) => {
  return request.get(`/public/series/subchapters/${id}/articles`)
}

/**
 * 添加文章到子章节
 * @param {number} id - 子章节ID
 * @param {Object} data - 文章数据
 * @param {number} data.article_id - 文章ID
 * @param {number} data.sort_order - 排序
 */
export const addArticleToSubchapter = (id, data) => {
  return request.post(`/rbac/series/subchapters/${id}/articles`, data)
}

/**
 * 从子章节移除文章
 * @param {number} id - 子章节ID
 * @param {number} articleId - 文章ID
 */
export const removeArticleFromSubchapter = (id, articleId) => {
  return request.delete(`/rbac/series/subchapters/${id}/articles/${articleId}`)
}
