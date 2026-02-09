import request from '@/utils/request'

/**
 * 列出指定路径下的文件和目录
 * @param {string} sessionId - SSH会话ID
 * @param {string} path - 要列出的目录路径
 * @returns {Promise} 文件列表数据
 */
export function listFiles(sessionId, path) {
  return request({
    url: '/rbac/sftp/list',
    method: 'get',
    params: { session_id: sessionId, path }
  })
}

/**
 * 上传文件（分片并发上传）
 * 将大文件分割为多个小分片，并发上传以提高速度
 * @param {string} sessionId - SSH会话ID
 * @param {string} path - 目标上传路径（如："/home/user"）
 * @param {File} file - 要上传的文件对象
 * @param {Function} onProgress - 进度回调函数，参数为进度百分比(0-100)
 * @returns {Promise<void>} 上传完成Promise
 */
export async function uploadFile(sessionId, path, file, onProgress) {
  // ==================== 步骤1: 计算分片参数 ====================
  const CHUNK_SIZE = 5 * 1024 * 1024 // 每个分片大小：5MB
  const totalSize = file.size           // 文件总大小（字节）
  const totalChunks = Math.ceil(totalSize / CHUNK_SIZE) // 总分片数
  const CONCURRENT_LIMIT = 3           // 并发上传的分片数限制

  // 使用 Map 来追踪每个分片的上传字节数
  const chunkUploadProgress = new Map() // chunkIndex -> 已上传字节数

  // 获取已上传的总字节数
  const getUploadedBytes = () => {
    let total = 0
    chunkUploadProgress.forEach((uploaded) => {
      total += uploaded
    })
    return total
  }

  // ==================== 步骤2: 定义单个分片上传函数 ====================
  const uploadSingleChunk = async (chunkIndex) => {
    // 步骤2.1: 计算当前分片的字节范围
    const start = chunkIndex * CHUNK_SIZE
    const end = Math.min(start + CHUNK_SIZE, totalSize)

    // 步骤2.2: 从原文件中切片数据（使用 slice 方法，不会创建副本）
    const chunk = file.slice(start, end)

    // 步骤2.3: 构建 FormData，包含分片和元数据
    const formData = new FormData()
    formData.append('file', chunk, file.name) // 分片文件数据（第三个参数为文件名）
    formData.append('path', path)             // 目标路径
    formData.append('file_name', file.name)   // 原始文件名（用于最终合并）
    formData.append('chunk_index', chunkIndex) // 当前分片索引（从0开始）
    formData.append('total_chunks', totalChunks) // 总分片数
    formData.append('total_size', totalSize)   // 文件总大小

    // 步骤2.4: 发起上传请求
    try {
      const response = await request({
        url: '/rbac/sftp/uploadFile',
        method: 'post',
        data: formData,
        params: { session_id: sessionId },
        onUploadProgress: (progressEvent) => {
          // 计算当前分片的已上传字节数
          const chunkUploaded = progressEvent.loaded

          // 更新当前分片的上传进度
          chunkUploadProgress.set(chunkIndex, chunkUploaded)

          // 计算总进度：所有分片已上传字节数
          const totalUploaded = getUploadedBytes()
          const totalProgress = (totalUploaded / totalSize) * 100

          console.log(`分片 ${chunkIndex} 上传进度: ${chunkUploaded}/${CHUNK_SIZE}字节, 总已上传: ${totalUploaded}/${totalSize}字节, 总进度: ${totalProgress.toFixed(2)}%`)

          // 回调通知进度更新
          if (onProgress) {
            onProgress(totalProgress)
          }
        }
      })

      return response.data
    } catch (error) {
      console.error(`切片 ${chunkIndex} 上传失败:`, error)
      throw error
    }
  }

  // ==================== 步骤3: 批量并发上传所有分片 ====================
  // 每次取 CONCURRENT_LIMIT 个分片作为一个批次，并发上传
  for (let i = 0; i < totalChunks; i += CONCURRENT_LIMIT) {
    const batch = []

    // 构建当前批次的所有分片上传任务
    for (let j = 0; j < CONCURRENT_LIMIT && i + j < totalChunks; j++) {
      batch.push(uploadSingleChunk(i + j))
    }

    // 步骤3.1: 并发执行当前批次的所有分片上传
    await Promise.all(batch)
  }

  console.log('所有切片上传完成')

  // 所有分片上传完成，确保进度显示为 100%
  if (onProgress) {
    onProgress(100)
  }
}

// 下载文件
export function downloadFile(sessionId, path) {
  return request({
    url: `/api/v1/ssh/files/download`,
    method: 'get',
    params: { session_id: sessionId, path },
    responseType: 'blob'
  }).then(response => {
    // 从响应头获取文件名
    const contentDisposition = response.headers['content-disposition']
    let filename = 'download'
    if (contentDisposition) {
      const match = contentDisposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/)
      if (match && match[1]) {
        filename = match[1].replace(/['"]/g, '')
      }
    }

    // 创建下载链接
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', filename)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  })
}

// 删除文件
export function deleteFile(sessionId, path) {
  return request({
    url: `/api/v1/ssh/files/delete`,
    method: 'delete',
    params: { session_id: sessionId, path }
  })
}

// 创建目录
export function createDir(sessionId, path) {
  return request({
    url: `/api/v1/ssh/files/mkdir`,
    method: 'post',
    params: { session_id: sessionId, path }
  })
}

// 重命名文件
export function renameFile(sessionId, oldPath, newPath) {
  return request({
    url: `/api/v1/ssh/files/rename`,
    method: 'post',
    data: { session_id: sessionId, old_path: oldPath, new_path: newPath }
  })
}
