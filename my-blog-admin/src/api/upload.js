import request from '@/utils/request'

// src/api/upload.js

/**
 * 上传并解析 Markdown 文件
 * @param {FormData} formData - 表单数据，包含文件
 * @returns {Promise} 解析后的文章数据
 */
export function uploadMarkdown(formData) {
  return request({
    url: '/rbac/upload/markdown',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}


/**
 * 上传图片
 * @param {FormData} formData - 表单数据
 * @param {File} formData.file - 文件
 * @param {string} formData.type - 类型(cover-封面, avatar-头像, logo-Logo)
 */
export function uploadImage(formData) {
  return request({
    url: '/rbac/upload/image',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

/**
 * 上传文件
 * @param {FormData} formData - 表单数据
 * @param {File} formData.file - 文件
 */
export function uploadFile(formData) {
  return request({
    url: '/rbac/upload/file',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

/**
 * 获取文件上传列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.page_size - 每页数量
 * @param {string} params.type - 文件类型
 */
export function getUploadList(params) {
  return request({
    url: '/rbac/upload/list',
    method: 'get',
    params
  })
}

/**
 * 删除文件
 * @param {number} id - 文件ID
 */
export function deleteUpload(id) {
  return request({
    url: `/rbac/upload/${id}`,
    method: 'delete'
  })
}

/**
 * 批量删除文件
 * @param {Array<number>} ids - 文件ID列表
 */
export function batchDeleteUploads(ids) {
  return request({
    url: '/rbac/upload/batch',
    method: 'delete',
    data: { ids }
  })
}
