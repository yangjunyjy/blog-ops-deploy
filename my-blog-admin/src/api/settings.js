import request from '@/utils/request'

/**
 * 获取系统设置
 */
export function getSystemSettings() {
  return request({
    url: '/admin/settings',
    method: 'get'
  })
}

/**
 * 更新系统设置
 * @param {Object} data - 设置数据
 */
export function updateSystemSettings(data) {
  return request({
    url: '/admin/settings',
    method: 'put',
    data
  })
}

/**
 * 重置系统设置
 * @param {Object} data - 重置选项
 * @param {Array<string>} data.groups - 要重置的设置组(site, seo, system, email, storage)
 */
export function resetSystemSettings(data) {
  return request({
    url: '/admin/settings/reset',
    method: 'post',
    data
  })
}

/**
 * 获取网站信息设置
 */
export function getSiteSettings() {
  return request({
    url: '/admin/settings/site',
    method: 'get'
  })
}

/**
 * 更新网站信息设置
 * @param {Object} data - 网站信息
 */
export function updateSiteSettings(data) {
  return request({
    url: '/admin/settings/site',
    method: 'put',
    data
  })
}

/**
 * 获取SEO设置
 */
export function getSEOSettings() {
  return request({
    url: '/admin/settings/seo',
    method: 'get'
  })
}

/**
 * 更新SEO设置
 * @param {Object} data - SEO设置
 */
export function updateSEOSettings(data) {
  return request({
    url: '/admin/settings/seo',
    method: 'put',
    data
  })
}

/**
 * 获取系统配置
 */
export function getSystemConfig() {
  return request({
    url: '/admin/settings/system',
    method: 'get'
  })
}

/**
 * 更新系统配置
 * @param {Object} data - 系统配置
 */
export function updateSystemConfig(data) {
  return request({
    url: '/admin/settings/system',
    method: 'put',
    data
  })
}

/**
 * 获取邮件设置
 */
export function getEmailSettings() {
  return request({
    url: '/admin/settings/email',
    method: 'get'
  })
}

/**
 * 更新邮件设置
 * @param {Object} data - 邮件设置
 */
export function updateEmailSettings(data) {
  return request({
    url: '/admin/settings/email',
    method: 'put',
    data
  })
}

/**
 * 发送测试邮件
 * @param {Object} data - 邮件数据
 * @param {string} data.to - 收件人邮箱
 * @param {string} data.subject - 邮件主题
 * @param {string} data.content - 邮件内容
 */
export function sendTestEmail(data) {
  return request({
    url: '/admin/settings/test-email',
    method: 'post',
    data
  })
}

/**
 * 获取存储设置
 */
export function getStorageSettings() {
  return request({
    url: '/admin/settings/storage',
    method: 'get'
  })
}

/**
 * 更新存储设置
 * @param {Object} data - 存储设置
 */
export function updateStorageSettings(data) {
  return request({
    url: '/admin/settings/storage',
    method: 'put',
    data
  })
}

/**
 * 创建备份
 * @param {Object} data - 备份数据
 * @param {string} data.type - 备份类型(full-完整备份, incremental-增量备份)
 * @param {string} data.name - 备份名称
 * @param {string} data.description - 备份描述
 */
export function createBackup(data) {
  return request({
    url: '/admin/backup/create',
    method: 'post',
    data
  })
}

/**
 * 获取备份列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.page_size - 每页数量
 * @param {string} params.type - 备份类型
 */
export function getBackupList(params) {
  return request({
    url: '/admin/backup/list',
    method: 'get',
    params
  })
}

/**
 * 恢复备份
 * @param {number} id - 备份ID
 */
export function restoreBackup(id) {
  return request({
    url: `/admin/backup/${id}/restore`,
    method: 'post'
  })
}

/**
 * 删除备份
 * @param {number} id - 备份ID
 */
export function deleteBackup(id) {
  return request({
    url: `/admin/backup/${id}`,
    method: 'delete'
  })
}

/**
 * 下载备份
 * @param {number} id - 备份ID
 */
export function downloadBackup(id) {
  return request({
    url: `/admin/backup/${id}/download`,
    method: 'get',
    responseType: 'blob'
  })
}
