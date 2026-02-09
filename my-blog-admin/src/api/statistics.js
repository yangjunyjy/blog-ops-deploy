import request from '@/utils/request'
import { mockApi } from '@/mock'

// 是否使用Mock数据
const USE_MOCK = true

/**
 * 获取仪表盘统计数据
 */
export const getDashboardStats = () => {
  if (USE_MOCK) {
    return mockApi.getDashboardStats()
  }
  return request.get('/rbac/statistics/dashboard')
}

/**
 * 获取文章统计数据
 * @param {Object} params - 查询参数
 * @param {number} params.days - 天数
 */
export const getArticleStats = (params) => {
  if (USE_MOCK) {
    return mockApi.getArticleStats(params)
  }
  return request.get('/rbac/statistics/articles', { params })
}

/**
 * 获取文章发布趋势
 * @param {Object} params - 查询参数
 * @param {number} params.days - 天数
 */
export const getArticleTrend = (params) => {
  if (USE_MOCK) {
    return mockApi.getArticleTrend(params)
  }
  return request.get('/rbac/statistics/articles/trend', { params })
}

/**
 * 获取浏览量趋势
 * @param {Object} params - 查询参数
 * @param {number} params.days - 天数
 */
export const getViewTrend = (params) => {
  if (USE_MOCK) {
    return mockApi.getViewTrend(params)
  }
  return request.get('/rbac/statistics/views/trend', { params })
}

/**
 * 获取热门文章
 * @param {Object} params - 查询参数
 * @param {number} params.limit - 数量
 * @param {number} params.days - 天数
 */
export const getHotArticles = (params) => {
  if (USE_MOCK) {
    return mockApi.getHotArticles(params)
  }
  return request.get('/rbac/statistics/articles/hot', { params })
}

/**
 * 获取分类统计
 */
export const getCategoryStats = () => {
  if (USE_MOCK) {
    return mockApi.getCategoryStats()
  }
  return request.get('/rbac/statistics/categories')
}

/**
 * 获取标签统计
 */
export const getTagStats = () => {
  if (USE_MOCK) {
    return mockApi.getTagStats()
  }
  return request.get('/rbac/statistics/tags')
}

/**
 * 获取热门标签
 * @param {Object} params - 查询参数
 * @param {number} params.limit - 数量
 */
export const getHotTags = (params) => {
  if (USE_MOCK) {
    return mockApi.getHotTags(params)
  }
  return request.get('/rbac/statistics/tags/hot', { params })
}

/**
 * 获取用户增长统计
 * @param {Object} params - 查询参数
 * @param {number} params.days - 天数
 */
export const getUserGrowth = (params) => {
  if (USE_MOCK) {
    return mockApi.getUserGrowth(params)
  }
  return request.get('/rbac/statistics/users/growth', { params })
}

/**
 * 获取活跃用户统计
 * @param {Object} params - 查询参数
 * @param {number} params.days - 天数
 */
export const getActiveUsers = (params) => {
  if (USE_MOCK) {
    return mockApi.getActiveUsers(params)
  }
  return request.get('/rbac/statistics/users/active', { params })
}

/**
 * 获取用户分布统计
 */
export const getUserDistribution = () => {
  if (USE_MOCK) {
    return mockApi.getUserDistribution()
  }
  return request.get('/rbac/statistics/users/distribution')
}

/**
 * 获取在线用户统计
 */
export const getOnlineUsers = () => {
  if (USE_MOCK) {
    return mockApi.getOnlineUsers()
  }
  return request.get('/rbac/statistics/users/online')
}
