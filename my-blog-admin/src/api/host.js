import request from '@/utils/request'
import { mockApi } from '@/mock'

// 是否使用Mock数据
const USE_MOCK = false

/**
 * 获取主机列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.page_size - 每页数量
 * @param {string} params.keyword - 搜索关键词（主机名称或IP）
 * @param {number} params.status - 状态(0-禁用, 1-启用)
 */
export function getHostList(params) {
  if (USE_MOCK) {
    return mockApi.getHostList(params)
  }
  return request({
    url: '/rbac/hosts',
    method: 'get',
    params: {
      page: params.page || 1,
      page_size: params.page_size || 10,
      name: params.keyword || '',
      status: params.status || ''
    }
  })
}

/**
 * 获取主机详情
 * @param {number} id - 主机ID
 */
export function getHostDetail(id) {
  if (USE_MOCK) {
    return mockApi.getHostDetail(id)
  }
  return request({
    url: `/rbac/hosts/${id}`,
    method: 'get'
  })
}

/**
 * 创建主机
 * @param {Object} data - 主机数据
 * @param {string} data.name - 主机名称
 * @param {string} data.ip - 主机IP
 * @param {number} data.port - 端口
 * @param {string} data.username - 用户名
 * @param {string} data.password - 密码（可选，如果使用密钥）
 * @param {string} data.private_key - 私钥（可选，如果使用密码）
 * @param {string} data.auth_type - 认证类型（password/key）
 * @param {string} data.description - 描述
 * @param {number} data.status - 状态(0-禁用, 1-启用)
 */
export function createHost(data) {
  if (USE_MOCK) {
    return mockApi.createHost(data)
  }
  return request({
    url: '/rbac/hosts',
    method: 'post',
    data: {
      name: data.name,
      address: data.address,
      port: data.port,
      username: data.username,
      password: data.password,
      secret_key: data.secret_key,
      type: data.type,
      status: data.status
    }
  })
}

/**
 * 更新主机
 * @param {Object} data - 主机数据（包含id）
 */
export function updateHost(data) {
  if (USE_MOCK) {
    return mockApi.updateHost(data)
  }
  return request({
    url: '/rbac/hosts',
    method: 'put',
    data: {
      id: data.id,
      name: data.name,
      address: data.ip,
      port: data.port,
      username: data.username,
      password: data.password,
      secret_key: data.private_key,
      type: data.auth_type,
      status: data.status === 1 ? 'active' : 'inactive'
    }
  })
}

/**
 * 删除主机
 * @param {number} id - 主机ID
 */
export function deleteHost(id) {
  if (USE_MOCK) {
    return mockApi.deleteHost(id)
  }
  return request({
    url: `/rbac/hosts/${id}`,
    method: 'delete'
  })
}

/**
 * 获取所有主机（用于下拉选择）
 */
export function getAllHosts() {
  if (USE_MOCK) {
    return mockApi.getAllHosts()
  }
  return request({
    url: '/rbac/hosts/all',
    method: 'get'
  })
}

/**
 * 获取活跃会话列表
 */
export function getSessions() {
  if (USE_MOCK) {
    return mockApi.getSessions()
  }
  return request({
    url: '/rbac/ssh/sessions',
    method: 'get'
  })
}

/**
 * 关闭SSH会话
 * @param {string} sessionId - 会话ID
 */
export function closeSession(sessionId) {
  if (USE_MOCK) {
    return mockApi.closeSession(sessionId)
  }
  return request({
    url: `/rbac/ssh/sessions/${sessionId}`,
    method: 'delete'
  })
}

/**
 * 测试主机连接
 * @param {number} id - 主机ID
 */
export function testHostConnection(id) {
  if (USE_MOCK) {
    return mockApi.testHostConnection(id)
  }
  return request({
    url: `/rbac/hosts/${id}/test`,
    method: 'post'
  })
}

/**
 * 远程连接主机（返回连接URL或信息）
 * @param {number} id - 主机ID
 * @param {string} auth_type - 认证类型（password/key）
 */
/**
 * 远程连接主机（返回WebSocket连接信息）
 * @param {number} id - 主机ID
 * @param {string} sessionId - 会话ID
 */
export function remoteConnectHost(id, sessionId) {
  if (USE_MOCK) {
    return mockApi.remoteConnectHost(id, sessionId)
  }
  // 返回WebSocket连接URL，不是HTTP请求
  return {
    url: `/rbac/ssh/connect/${id}?session_id=${sessionId}`,
    method: 'websocket'
  }
}

