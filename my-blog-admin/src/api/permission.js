import request from '@/utils/request'

/**
 * 获取当前用户菜单列表
 */
export function getUserMenus() {
  return request({
    url: '/rbac/auth/menu',
    method: 'get'
  })
}

/**
 * 检查权限（前端本地验证，无需请求后端）
 * @param {string|Array<string>} permission - 权限标识
 */
export function hasPermission(permission) {
  // 从 store 获取用户权限
  const permissions = JSON.parse(localStorage.getItem('user_permissions') || '[]')

  if (Array.isArray(permission)) {
    // 检查是否拥有任一权限
    return permission.some(p => permissions.includes(p))
  } else {
    // 检查是否拥有该权限
    return permissions.includes(permission)
  }
}

/**
 * 检查是否拥有所有权限
 * @param {Array<string>} permissions - 权限标识列表
 */
export function hasAllPermissions(permissions) {
  const userPermissions = JSON.parse(localStorage.getItem('user_permissions') || '[]')
  return permissions.every(p => userPermissions.includes(p))
}

/**
 * 检查是否拥有任一权限
 * @param {Array<string>} permissions - 权限标识列表
 */
export function hasAnyPermission(permissions) {
  const userPermissions = JSON.parse(localStorage.getItem('user_permissions') || '[]')
  return permissions.some(p => userPermissions.includes(p))
}
