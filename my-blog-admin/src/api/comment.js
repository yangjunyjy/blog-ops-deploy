import request from '@/utils/request'
import { mockApi } from '@/mock'

// 是否使用Mock数据
const USE_MOCK = false

/**
 * 获取评论列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.page_size - 每页数量
 * @param {number} params.article_id - 文章ID
 * @param {number} params.status - 状态(0-待审核, 1-已通过, 2-已拒绝)
 */
export const getCommentList = (params) => {
  if (USE_MOCK) {
    return mockApi.getCommentList(params)
  }
  return request.get('/rbac/comments', { params })
}

/**
 * 获取评论详情
 * @param {number} id - 评论ID
 */
export const getCommentDetail = (id) => {
  if (USE_MOCK) {
    return mockApi.getCommentDetail(id)
  }
  return request.get(`/rbac/comments/${id}`)
}

/**
 * 审核通过评论
 * @param {number} id - 评论ID
 */
export const approveComment = (id) => {
  if (USE_MOCK) {
    return mockApi.approveComment(id)
  }
  return request.put(`/rbac/comments/${id}/approve`)
}

/**
 * 审核拒绝评论
 * @param {number} id - 评论ID
 * @param {Object} data - 拒绝数据
 * @param {string} data.reason - 拒绝原因
 */
export const rejectComment = (id, data) => {
  if (USE_MOCK) {
    return mockApi.rejectComment(id, data)
  }
  return request.put(`/rbac/comments/${id}/reject`, data)
}

/**
 * 删除评论
 * @param {number} id - 评论ID
 */
export const deleteComment = (id) => {
  if (USE_MOCK) {
    return mockApi.deleteComment(id)
  }
  return request.delete(`/rbac/comments/${id}`)
}

/**
 * 批量删除评论
 * @param {Array<number>} ids - 评论ID列表
 */
export const batchDeleteComments = (ids) => {
  if (USE_MOCK) {
    return mockApi.batchDeleteComments(ids)
  }
  return request.delete('/rbac/comments/batch', { data: { ids } })
}

/**
 * 批量审核通过评论
 * @param {Array<number>} ids - 评论ID列表
 */
export const batchApproveComments = (ids) => {
  if (USE_MOCK) {
    return mockApi.batchApproveComments(ids)
  }
  return request.put('/rbac/comments/batch/approve', { ids })
}

/**
 * 批量审核拒绝评论
 * @param {Array<number>} ids - 评论ID列表
 * @param {string} reason - 拒绝原因
 */
export const batchRejectComments = (ids, reason) => {
  if (USE_MOCK) {
    return mockApi.batchRejectComments(ids, reason)
  }
  return request.put('/rbac/comments/batch/reject', { ids, reason })
}
