import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: () => import('@/views/Home.vue'),
      meta: { title: '首页' }
    },
    {
      path: '/articles',
      name: 'Articles',
      component: () => import('@/views/Articles.vue'),
      meta: { title: '文章列表' }
    },
    {
      path: '/article/:id',
      name: 'ArticleDetail',
      component: () => import('@/views/ArticleDetail.vue'),
      meta: { title: '文章详情' }
    },
    {
      path: '/categories',
      name: 'Categories',
      component: () => import('@/views/Categories.vue'),
      meta: { title: '所有分类' }
    },
    {
      path: '/category/:id',
      name: 'Category',
      component: () => import('@/views/Category.vue'),
      meta: { title: '分类详情' }
    },
    {
      path: '/tag',
      name: 'Tags',
      component: () => import('@/views/Tags.vue'),
      meta: { title: '所有标签' }
    },
    {
      path: '/tag/:name',
      name: 'Tag',
      component: () => import('@/views/Tag.vue'),
      meta: { title: '标签详情' }
    },
    {
      path: '/search',
      name: 'Search',
      component: () => import('@/views/Search.vue'),
      meta: { title: '搜索' }
    },
    {
      path: '/about',
      name: 'About',
      component: () => import('@/views/About.vue'),
      meta: { title: '关于我' }
    },
    {
      path: '/series',
      name: 'Series',
      component: () => import('@/views/Series.vue'),
      meta: { title: '专栏' }
    },
    {
      path: '/series/:id',
      name: 'SeriesDetail',
      component: () => import('@/views/Series.vue'),
      meta: { title: '专栏详情' }
    },
    {
      path: '/series/:seriesId/chapter/:chapterId',
      name: 'ChapterDetail',
      component: () => import('@/views/ChapterDetail.vue'),
      meta: { title: '章节详情', hideLayout: true }
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue'),
      meta: { title: '登录' }
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/Register.vue'),
      meta: { title: '注册' }
    },
    {
      path: '/forgot',
      name: 'Forgot',
      component: () => import('@/views/Forgot.vue'),
      meta: { title: '忘记密码' }
    },
    {
      path: '/profile',
      name: 'Profile',
      component: () => import('@/views/Profile.vue'),
      meta: { title: '个人中心', requiresAuth: true }
    },
    // {
    //   path: '/create',
    //   name: 'Create',
    //   component: () => import('@/views/Create.vue'),
    //   meta: { title: '创作文章', requiresAuth: true, hideLayout: true }
    // }
  ]
})

// 全局路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title}` : '杨骏的博客'

  // 检查是否需要登录
  const userStore = useUserStore()
  userStore.loadUserFromStorage()

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    // 未登录跳转到登录页
    next('/login')
  } else if ((to.path === '/login' || to.path === '/register') && userStore.isLoggedIn) {
    // 已登录访问登录/注册页，跳转到首页
    next('/')
  } else {
    next()
  }
})

export default router
