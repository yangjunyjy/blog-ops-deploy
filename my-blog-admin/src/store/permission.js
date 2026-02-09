import { ref } from 'vue'

// RBAC权限数据
const mockRoles = {
  admin: {
    id: 1,
    name: '超级管理员',
    permissions: ['*'], // 所有权限
    menus: [
      {
        id: 'dashboard',
        path: '/dashboard',
        title: '首页',
        icon: 'HomeFilled',
        type: 'menu'
      },
      {
        id: 'content',
        title: '内容管理',
        icon: 'Document',
        type: 'submenu',
        children: [
          { id: 'articles', path: '/articles', title: '文章管理', icon: 'Document', type: 'menu' },
          { id: 'categories', path: '/categories', title: '分类管理', icon: 'Folder', type: 'menu' },
          { id: 'tags', path: '/tags', title: '标签管理', icon: 'PriceTag', type: 'menu' },
          { id: 'series', path: '/series', title: '系列管理', icon: 'Collection', type: 'menu' }
        ]
      },
      {
        id: 'interaction',
        title: '互动管理',
        icon: 'ChatDotRound',
        type: 'submenu',
        children: [
          { id: 'comments', path: '/comments', title: '评论管理', icon: 'ChatDotRound', type: 'menu' }
        ]
      },
      {
        id: 'user',
        title: '用户管理',
        icon: 'User',
        type: 'submenu',
        children: [
          { id: 'users', path: '/users', title: '用户列表', icon: 'User', type: 'menu' }
        ]
      },
      {
        id: 'statistics',
        title: '数据统计',
        icon: 'DataLine',
        type: 'submenu',
        children: [
          { id: 'statistics-overview', path: '/statistics/overview', title: '数据概览', icon: 'DataLine', type: 'menu' },
          { id: 'statistics-content', path: '/statistics/content', title: '内容统计', icon: 'TrendCharts', type: 'menu' },
          { id: 'statistics-user', path: '/statistics/user', title: '用户统计', icon: 'UserFilled', type: 'menu' }
        ]
      }
    ]
  },
  editor: {
    id: 2,
    name: '编辑',
    permissions: ['article:read', 'article:create', 'article:update', 'category:read', 'tag:read', 'series:read', 'comment:read'],
    menus: [
      {
        id: 'dashboard',
        path: '/dashboard',
        title: '首页',
        icon: 'HomeFilled',
        type: 'menu'
      },
      {
        id: 'content',
        title: '内容管理',
        icon: 'Document',
        type: 'submenu',
        children: [
          { id: 'articles', path: '/articles', title: '文章管理', icon: 'Document', type: 'menu' },
          { id: 'categories', path: '/categories', title: '分类管理', icon: 'Folder', type: 'menu' },
          { id: 'tags', path: '/tags', title: '标签管理', icon: 'PriceTag', type: 'menu' },
          { id: 'series', path: '/series', title: '系列管理', icon: 'Collection', type: 'menu' }
        ]
      },
      {
        id: 'interaction',
        title: '互动管理',
        icon: 'ChatDotRound',
        type: 'submenu',
        children: [
          { id: 'comments', path: '/comments', title: '评论管理', icon: 'ChatDotRound', type: 'menu' }
        ]
      }
    ]
  },
  visitor: {
    id: 3,
    name: '访客',
    permissions: ['article:read'],
    menus: [
      {
        id: 'dashboard',
        path: '/dashboard',
        title: '首页',
        icon: 'HomeFilled',
        type: 'menu'
      },
      {
        id: 'content',
        title: '内容管理',
        icon: 'Document',
        type: 'submenu',
        children: [
          { id: 'articles', path: '/articles', title: '文章管理', icon: 'Document', type: 'menu' }
        ]
      }
    ]
  }
}

// 用户角色映射
const mockUsers = {
  'admin@blog.com': { role: 'admin' },
  'editor@blog.com': { role: 'editor' },
  'visitor@blog.com': { role: 'visitor' }
}

const state = ref({
  routes: [],
  menus: [],
  permissions: []
})

// 根据用户权限获取菜单
const getMenusByRole = (role) => {
  const roleData = mockRoles[role] || mockRoles.visitor
  return roleData.menus || []
}

// 根据用户权限获取路由
const getRoutesByPermission = (role) => {
  const roleData = mockRoles[role] || mockRoles.visitor
  const permissions = roleData.permissions

  // 如果有所有权限，返回所有路由
  if (permissions.includes('*')) {
    return []
  }

  // 根据权限过滤路由
  const permissionMap = {
    'article:read': ['/articles'],
    'category:read': ['/categories'],
    'tag:read': ['/tags'],
    'series:read': ['/series'],
    'comment:read': ['/comments'],
    'user:read': ['/users'],
    'statistics:read': ['/statistics/overview', '/statistics/content', '/statistics/user']
  }

  const allowedRoutes = []
  Object.keys(permissionMap).forEach(permission => {
    if (permissions.includes(permission)) {
      allowedRoutes.push(...permissionMap[permission])
    }
  })

  return allowedRoutes
}

// 根据用户邮箱获取角色
const getRoleByEmail = (email) => {
  const userData = mockUsers[email]
  return userData?.role || 'visitor'
}

// 生成权限状态
const generatePermission = (role) => {
  const roleData = mockRoles[role] || mockRoles.visitor
  state.value.menus = roleData.menus
  state.value.permissions = roleData.permissions
  return state.value
}

// 权限检查函数
const hasPermission = (permission) => {
  const permissions = state.value.permissions
  return permissions.includes('*') || permissions.includes(permission)
}

// 重置权限
const resetPermission = () => {
  state.value = {
    routes: [],
    menus: [],
    permissions: []
  }
}

export {
  state,
  getMenusByRole,
  getRoutesByPermission,
  getRoleByEmail,
  generatePermission,
  hasPermission,
  resetPermission
}
