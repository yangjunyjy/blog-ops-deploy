import request from '@/utils/request'

/**
 * 获取待审核文章列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.page_size - 每页数量
 * @param {string} params.keyword - 搜索关键词
 */
export function getPendingArticles(params) {
  return request({
    url: '/admin/articles/pending',
    method: 'get',
    params
  })
}

/**
 * 获取已审核通过文章列表
 * @param {Object} params - 查询参数
 */
export function getApprovedArticles(params) {
  return request({
    url: '/admin/articles/approved',
    method: 'get',
    params
  })
}

/**
 * 获取已审核拒绝文章列表
 * @param {Object} params - 查询参数
 */
export function getRejectedArticles(params) {
  return request({
    url: '/admin/articles/rejected',
    method: 'get',
    params
  })
}

/**
 * 审核通过文章
 * @param {number} id - 文章ID
 * @param {Object} data - 审核数据
 * @param {string} data.remark - 审核备注
 */
export function approveArticle(id, data) {
  return request({
    url: `/admin/articles/${id}/approve`,
    method: 'put',
    data
  })
}

/**
 * 审核拒绝文章
 * @param {number} id - 文章ID
 * @param {Object} data - 拒绝数据
 * @param {string} data.reason - 拒绝原因
 * @param {string} data.remark - 备注
 */
export function rejectArticle(id, data) {
  return request({
    url: `/admin/articles/${id}/reject`,
    method: 'put',
    data
  })
}

/**
 * 获取文章审核日志
 * @param {number} id - 文章ID
 * @param {Object} params - 查询参数
 */
export function getArticleAuditLogs(id, params) {
  return request({
    url: `/admin/articles/${id}/audit-log`,
    method: 'get',
    params
  })
}

/**
 * 获取待审核评论列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.page_size - 每页数量
 * @param {number} params.level - 评论层级(0-顶级, 1-二级, 2-三级)
 */
export function getPendingComments(params) {
  return request({
    url: '/admin/comments/pending',
    method: 'get',
    params
  })
}

/**
 * 获取已审核通过评论列表
 * @param {Object} params - 查询参数
 */
export function getApprovedComments(params) {
  return request({
    url: '/admin/comments/approved',
    method: 'get',
    params
  })
}

/**
 * 获取已审核拒绝评论列表
 * @param {Object} params - 查询参数
 */
export function getRejectedComments(params) {
  return request({
    url: '/admin/comments/rejected',
    method: 'get',
    params
  })
}

/**
 * 审核通过评论
 * @param {number} id - 评论ID
 */
export function approveComment(id) {
  return request({
    url: `/admin/comments/${id}/approve`,
    method: 'put'
  })
}

/**
 * 审核拒绝评论
 * @param {number} id - 评论ID
 * @param {Object} data - 拒绝数据
 * @param {string} data.reason - 拒绝原因
 */
export function rejectComment(id, data) {
  return request({
    url: `/admin/comments/${id}/reject`,
    method: 'put',
    data
  })
}

/**
 * 批量审核通过评论
 * @param {Array<number>} ids - 评论ID列表
 */
export function batchApproveComments(ids) {
  return request({
    url: '/admin/comments/batch/approve',
    method: 'put',
    data: { ids }
  })
}

/**
 * 批量审核拒绝评论
 * @param {Array<number>} ids - 评论ID列表
 * @param {string} reason - 拒绝原因
 */
export function batchRejectComments(ids, reason) {
  return request({
    url: '/admin/comments/batch/reject',
    method: 'put',
    data: { ids, reason }
  })
}

/**
 * 获取待审核用户列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.page_size - 每页数量
 * @param {string} params.keyword - 搜索关键词
 */
export function getPendingUsers(params) {
  return request({
    url: '/admin/users/pending',
    method: 'get',
    params
  })
}

/**
 * 获取已审核通过用户列表
 * @param {Object} params - 查询参数
 */
export function getApprovedUsers(params) {
  return request({
    url: '/admin/users/approved',
    method: 'get',
    params
  })
}

/**
 * 获取已审核拒绝用户列表
 * @param {Object} params - 查询参数
 */
export function getRejectedUsers(params) {
  return request({
    url: '/admin/users/rejected',
    method: 'get',
    params
  })
}

/**
 * 审核通过用户
 * @param {number} id - 用户ID
 */
export function approveUser(id) {
  return request({
    url: `/admin/users/${id}/approve`,
    method: 'put'
  })
}

/**
 * 审核拒绝用户
 * @param {number} id - 用户ID
 * @param {Object} data - 拒绝数据
 * @param {string} data.reason - 拒绝原因
 */
export function rejectUser(id, data) {
  return request({
    url: `/admin/users/${id}/reject`,
    method: 'put',
    data
  })
}
