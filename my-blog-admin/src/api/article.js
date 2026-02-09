import request from '@/utils/request'
import { mockApi } from '@/mock'

// 是否使用Mock数据
const USE_MOCK = false

/**
 * 获取文章列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.pageSize - 每页数量
 * @param {string} params.keyword - 搜索关键词
 * @param {number} params.category_id - 分类ID
 * @param {number} params.tag_id - 标签ID
 * @param {number} params.status - 状态(0-草稿, 1-已发布)
 */
export const getArticleList = (params) => {
  if (USE_MOCK) {
    return mockApi.getArticleList(params)
  }
  return request.get('/public/articles', { params })
}

/**
 * 获取文章详情
 * @param {number} id - 文章ID
 */
export const getArticleDetail = (id) => {
  if (USE_MOCK) {
    return mockApi.getArticleDetail(id)
  }
  return request.get(`/public/articles/${id}`)
}

/**
 * 创建文章
 * @param {Object} data - 文章数据
 * @param {string} data.title - 文章标题
 * @param {string} data.slug - URL别名
 * @param {string} data.summary - 文章摘要
 * @param {string} data.content - 文章内容
 * @param {string} data.cover - 封面URL
 * @param {number} data.category_id - 分类ID
 * @param {Array<number>} data.tag_ids - 标签ID列表
 * @param {number} data.status - 状态(0-草稿, 1-已发布)
 * @param {boolean} data.is_top - 是否置顶
 */
export const createArticle = (data) => {
  if (USE_MOCK) {
    return mockApi.createArticle(data)
  }
  return request.post('/rbac/articles', data)
}

/**
 * 更新文章
 * @param {number} id - 文章ID
 * @param {Object} data - 文章数据
 */
export const updateArticle = (id, data) => {
  if (USE_MOCK) {
    return mockApi.updateArticle(id, data)
  }
  return request.put(`/rbac/articles/${id}`, data)
}

/**
 * 删除文章
 * @param {number} id - 文章ID
 */
export const deleteArticle = (id) => {
  if (USE_MOCK) {
    return mockApi.deleteArticle(id)
  }
  return request.delete(`/rbac/articles/${id}`)
}

/**
 * 批量删除文章
 * @param {Array<number>} ids - 文章ID列表
 */
export const batchDeleteArticles = (ids) => {
  if (USE_MOCK) {
    return mockApi.batchDeleteArticles(ids)
  }
  return request.delete('/rbac/articles/batch', { data: { ids } })
}

/**
 * 发布文章
 * @param {number} id - 文章ID
 */
export const publishArticle = (id) => {
  if (USE_MOCK) {
    return mockApi.publishArticle(id)
  }
  return request.put(`/rbac/articles/${id}/publish`)
}

/**
 * 撤回文章
 * @param {number} id - 文章ID
 */
export const withdrawArticle = (id) => {
  if (USE_MOCK) {
    return mockApi.withdrawArticle(id)
  }
  return request.put(`/rbac/articles/${id}/withdraw`)
}

/**
 * 批量发布文章
 * @param {Array<number>} ids - 文章ID列表
 */
export const batchPublishArticles = (ids) => {
  if (USE_MOCK) {
    return mockApi.batchPublishArticles(ids)
  }
  return request.put('/rbac/articles/batch/publish', { ids })
}

/**
 * 批量撤回文章
 * @param {Array<number>} ids - 文章ID列表
 */
export const batchWithdrawArticles = (ids) => {
  if (USE_MOCK) {
    return mockApi.batchWithdrawArticles(ids)
  }
  return request.put('/rbac/articles/batch/withdraw', { ids })
}

/**
 * 更新文章状态
 * @param {number} id - 文章ID
 * @param {Object} data - 状态数据
 * @param {number} data.status - 状态(0-草稿, 1-已发布)
 */
export const updateArticleStatus = (id, data) => {
  if (USE_MOCK) {
    return mockApi.updateArticleStatus(id, data)
  }
  return request.put(`/rbac/articles/${id}/status`, data)
}

/**
 * 批量更新文章状态
 * @param {Array<number>} ids - 文章ID列表
 * @param {number} status - 状态(0-草稿, 1-已发布)
 */
export const batchUpdateArticleStatus = (ids, status) => {
  if (USE_MOCK) {
    return mockApi.batchUpdateArticleStatus(ids, status)
  }
  return request.put('/rbac/articles/batch/status', { ids, status })
}
