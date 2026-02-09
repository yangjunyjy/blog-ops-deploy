import request from '@/utils/request'

/**
 * 上传并解析 Word 文档
 * @param {FormData} formData - 表单数据，包含文件
 * @returns {Promise} 解析后的文章数据
 */
export function uploadWord(formData) {
  return request({
    url: '/admin/upload/word',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
