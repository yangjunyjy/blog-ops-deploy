import { createRouter, createWebHistory } from 'vue-router'
import { useRbacStore, MenuType } from '../store/rbac'
import { storeToRefs } from 'pinia'

// 组件映射表，用于动态导入


// 从菜单数据生成路由配置
function generateRoutesFromMenus(menus) {
  const routes = []

  const processMenu = (menu, parentPath = '') => {
    console.log('处理菜单:', {
      title: menu.title,
      type: menu.type,
      path: menu.path,
      component: menu.component,
      hasChildren: menu.children?.length || 0
    })

    // 处理类型为 MENU 和 DIRECTORY 的菜单项
    if (menu.type !== MenuType.MENU && menu.type !== MenuType.DIRECTORY) {
      console.log('跳过按钮类型菜单:', menu.title)
      return null
    }

    // 对于目录类型，如果 path 为空，跳过路由生成，只处理子菜单
    if (menu.type === MenuType.DIRECTORY) {
      if (!menu.path || menu.path.trim() === '') {
        console.log('目录无路径，处理子菜单:', menu.title)
        // 处理子菜单
        if (menu.children && menu.children.length > 0) {
          menu.children.forEach(child => {
            const childRoute = processMenu(child, parentPath)
            if (childRoute) {
              routes.push(childRoute)
            }
          })
        }
        return null
      }
    } else {
      // 菜单类型必须有路径
      if (!menu.path || typeof menu.path !== 'string' || menu.path.trim() === '') {
        console.warn('菜单项缺少路径:', menu)
        return null
      }
    }

    // 构建完整路径 - 修复路径重复问题
    let fullPath
    if (parentPath && menu.path) {
      // 检查 menu.path 是否已经是绝对路径
      if (menu.path.startsWith('/')) {
        // 已经是绝对路径，直接使用
        fullPath = menu.path
      } else {
        // 相对路径，拼接父路径
        fullPath = `${parentPath}/${menu.path}`
      }
    } else {
      fullPath = menu.path || parentPath
    }

    // 对于目录类型，如果没有组件，使用 DirectoryWrapper
    let componentKey = menu.component
    if (menu.type === MenuType.DIRECTORY && (!componentKey || componentKey.trim() === '')) {
      componentKey = 'DirectoryWrapper'
    }

    // 菜单类型必须有组件
    if (!componentKey || componentKey.trim() === '') {
      console.warn('菜单项缺少组件:', menu.title, menu)
      return null
    }

    // 标准化组件路径
    let normalizedComponentKey = componentKey.trim()
    if (!normalizedComponentKey.endsWith('.vue')) {
      normalizedComponentKey += '.vue'
    }

    // 确保路径以 views/ 开头
    if (!normalizedComponentKey.startsWith('views/') && normalizedComponentKey !== 'DirectoryWrapper') {
      normalizedComponentKey = `views/${normalizedComponentKey}`
    }

    console.log('组件路径:', componentKey, '->', normalizedComponentKey)
    // 创建路由对象
    const route = {
      path: fullPath,
      name: menu.menu_name || menu.menuCode || menu.path.replace(/\//g, '-').replace(/^-/, ''),
      component: () => import(/* @vite-ignore */`../${normalizedComponentKey}`),
      meta: {
        title: menu.title || '未命名',
        icon: menu.icon,
        permission: menu.permission,
        menuCode: menu.menuCode
      },
      _dynamic: true // 标记为动态路由，便于重置时移除
    }

    console.log('生成路由:', route.path, route.name)

    // 如果有子菜单，递归处理
    if (menu.children && menu.children.length > 0) {
      route.children = []
      menu.children.forEach(child => {
        const childRoute = processMenu(child, fullPath)
        if (childRoute) {
          route.children.push(childRoute)
        }
      })
    }

    return route
  }
  
  // 处理所有菜单项
  menus.forEach(menu => {
    const route = processMenu(menu)
    if (route) {
      routes.push(route)
    }
  })
  
  console.log('从菜单生成路由，数量:', routes.length)
  return routes
}


const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    redirect: '/dashboard/index', // 重定向到子路由
    component: () => import('../layout/Layout.vue'),
    meta: { hidden: true }, // 隐藏，不显示在菜单/标签中
    children: [
      {
        path: 'index',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { title: '首页' }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('../views/profile/Profile.vue'),
        meta: { title: '个人中心' }
      }
    ]
  },
  {
    path: '/terminal/:hostId',
    name: 'Terminal',
    component: () => import('../views/host/Terminal.vue'),
    meta: { title: 'SSH终端' }
  },
  {
    path: '/:pathMatch(.*)*',
    component: () => import('../views/404.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})
// 更精确的单层路径检查函数
function isValidSingleLevelPath(path) {
  if (typeof path !== 'string') return false
  
  // 必须是 / 开头且长度大于1
  if (!path.startsWith('/') || path.length <= 1) return false
  
  // 不能有第二个斜杠
  if (path.indexOf('/', 1) !== -1) return false
  
  // 也不能以斜杠结尾（虽然是单层，但以防万一）
  if (path.endsWith('/')) return false
  
  return true
}
// 添加动态路由（根据菜单生成）
export function addRoutes(menus) {
  console.log('开始添加动态路由，菜单数量:', menus?.length || 0)

  // 从菜单生成路由
  const menuRoutes = generateRoutesFromMenus(menus || [])

  // 添加菜单路由到 dashboard 布局
  menuRoutes.forEach(route => {
    // 判断是否需要隐藏布局：路径是 /create 或 /login
    console.log(route.path);
    console.log(isValidSingleLevelPath(route.path));
    
    const shouldHideLayout = route.path === '/create' || route.path === '/editor/:id';
    if (shouldHideLayout) {
      console.log('添加路由:', route.path, '到根路由 (hideLayout)')
      router.addRoute(route)
    } else {
      console.log('添加路由:', route.path, '到 dashboard')
      router.addRoute('dashboard', route)
    }
  })

  console.log('动态路由添加完成，共:', menuRoutes.length, '个路由')
}

// 清除动态路由
export function resetRoutes() {
  // 获取所有动态路由名称
  const dynamicRouteNames = []
  
  const removeDynamicRoutes = (routes) => {
    routes.forEach(route => {
      if (route._dynamic) {
        if (router.hasRoute(route.name)) {
          try {
            router.removeRoute(route.name)
            dynamicRouteNames.push(route.name)
          } catch (e) {
            console.warn('移除路由失败:', route.name, e)
          }
        }
      }
      // 递归处理子路由
      if (route.children) {
        removeDynamicRoutes(route.children)
      }
    })
  }
  
  // 遍历所有路由，移除动态路由
  router.getRoutes().forEach(route => {
    if (route._dynamic) {
      if (router.hasRoute(route.name)) {
        try {
          router.removeRoute(route.name)
          dynamicRouteNames.push(route.name)
        } catch (e) {
          console.warn('移除路由失败:', route.name, e)
        }
      }
    }
  })
  
  console.log('清除动态路由完成，共移除:', dynamicRouteNames.length, '个路由')
}

// 存储上次的菜单数据，用于检测菜单变化
let lastMenus = []
let routesInitialized = false


// 检查是否需要重新初始化路由
const shouldReinitRoutes = () => {
  const flag = sessionStorage.getItem('routeInitialized')
  return flag === 'false' || flag === null
}

// 路由守卫
router.beforeEach(async (to, _from, next) => {
  document.title = to.meta.title? `${to.meta.title}`:"博客后台"
  const token = sessionStorage.getItem('token')
    const whiteList = ['/login', '/terminal']

    if (whiteList.includes(to.path) || to.path.startsWith('/terminal/')) {
      if (token) {
        next()
      } else {
      next()
    }
    return
  }

  if (!token) {
    next('/login')
    return
  }

  try {
    const userInfo = JSON.parse(sessionStorage.getItem('userInfo') || '{}')
    // 验证用户信息是否有email
    if (!userInfo.email) {
      console.error('用户信息不完整，跳转到登录页')
      sessionStorage.clear()
      next('/login')
      return
    }

    const rbacStore = useRbacStore()
    const { menus, menusLoaded } = storeToRefs(rbacStore)

    // 检查是否需要重新初始化路由
    if (shouldReinitRoutes()) {
      console.log('检测到路由重新初始化标志,强制重置')
      routesInitialized = false
      sessionStorage.setItem('routeInitialized', 'true')
    }

    // 首次加载或菜单未加载时，从后端获取菜单并初始化路由
    if (!routesInitialized || !menusLoaded.value) {
      console.log('初始化路由，菜单已加载:', menusLoaded.value, '路由已初始化:', routesInitialized)

      try {
        // 如果菜单未加载，先加载菜单
        if (!menusLoaded.value) {
          console.log('菜单未加载，从后端获取菜单数据')
          await rbacStore.loadAllDataFromBackend()
        }

        // 清除旧的动态路由
        resetRoutes()

        // 根据菜单生成并添加新路由
        addRoutes(menus.value)

        // 记录当前菜单数据
        lastMenus = JSON.parse(JSON.stringify(menus.value))
        routesInitialized = true

        // 首次加载时，总是重新导航以确保动态路由生效
        console.log('路由初始化完成，重新导航到:', to.path)
        next({ ...to, replace: true })
        return
      } catch (error) {
        console.error('路由初始化失败:', error)
        sessionStorage.clear()
        next('/login')
        return
      }
    }

    // 检查菜单是否发生变化
    const menusChanged = JSON.stringify(menus.value) !== JSON.stringify(lastMenus)

    // 菜单变化时，重新加载路由
    if (menusChanged) {
      console.log('菜单数据变化，重新加载路由')

      // 清除旧的动态路由
      resetRoutes()

      // 根据菜单生成并添加新路由
      addRoutes(menus.value)

      // 记录当前菜单数据
      lastMenus = JSON.parse(JSON.stringify(menus.value))
      // 菜单变化时不重新导航，避免闪烁
    }
  } catch (error) {
    console.error('路由初始化失败:', error)
  }

  next()
})

export default router
