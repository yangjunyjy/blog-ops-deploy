/**
 * 主机数据适配器
 * 用于处理前后端字段名不一致的问题
 */

/**
 * 将后端响应的数据转换为前端格式
 * @param {Object} host - 后端返回的主机数据
 * @returns {Object} 前端使用的主机数据
 */
export function adaptHostFromBackend(host) {
  return {
    id: host.id,
    name: host.name,
    ip: host.address,           // 后端 address -> 前端 ip
    port: host.port,
    username: host.username,
    password: '',               // 不返回密码
    private_key: '',           // 不返回私钥
    auth_type: host.type,       // 后端 type -> 前端 auth_type
    status: host.status === 'active' ? 1 : 0,  // 'active'/'inactive' -> 1/0
    created_at: host.created_at,
    updated_at: host.updated_at
  }
}

/**
 * 将前端表单数据转换为后端格式
 * @param {Object} formData - 前端表单数据
 * @returns {Object} 后端需要的数据格式
 */
export function adaptHostToBackend(formData) {
  return {
    name: formData.name,
    address: formData.ip,              // 前端 ip -> 后端 address
    port: formData.port,
    username: formData.username,
    password: formData.password,
    secret_key: formData.private_key,    // 前端 private_key -> 后端 secret_key
    type: formData.auth_type,            // 前端 auth_type -> 后端 type
    status: formData.status === 1 ? 'active' : 'inactive'  // 1/0 -> 'active'/'inactive'
  }
}

/**
 * 适配主机列表
 * @param {Object} response - 后端返回的列表响应
 * @returns {Object} 前端使用的列表数据
 */
export function adaptHostListFromBackend(response) {
  return {
    list: (response.items || []).map(host => adaptHostFromBackend(host)),
    total: response.total || 0
  }
}

/**
 * 适配测试连接响应
 * @param {Object} response - 后端返回的测试结果
 * @returns {Object} 前端使用的测试结果
 */
export function adaptTestConnectionFromBackend(response) {
  return {
    success: response.success || false,
    message: response.message || ''
  }
}
