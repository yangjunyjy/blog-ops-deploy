<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Search, Menu, User, SwitchButton, ArrowDown, Sunny, Moon } from '@element-plus/icons-vue'
import { getCategories, getSeries } from '@/api'
import { useUserStore } from '@/stores/user'
import { useThemeStore } from '@/stores/theme'
import { ElMessageBox, ElMessage } from 'element-plus'
import HaloButton from './HaloButton.vue'

const router = useRouter()
const route = useRoute()

const userStore = useUserStore()
const themeStore = useThemeStore()

const searchKeyword = ref('')
const isMobileMenuOpen = ref(false)
const categories = ref([])
const seriesList = ref([])
const userDropdownVisible = ref(false)

const avatarUrl = computed(() => {
  const avatar = userStore.user?.avatar
  // 如果用户有上传头像且不为空，使用上传的头像
  // 否则使用默认头像
  return avatar && avatar.trim() ? avatar : '/images/default.jpg'
})
const userName = computed(() => userStore.user?.nickname || userStore.user?.username || '用户')
const userEmail = computed(() => userStore.user?.email || '')

const isSeriesActive = computed(() => route.path.startsWith('/series'))

onMounted(async () => {
  try {
    const [catRes, seriesRes] = await Promise.all([
      getCategories(),
      getSeries()
    ])
    categories.value = catRes.data || []
    seriesList.value = seriesRes.data || []
  } catch (error) {
    console.error('加载数据失败:', error)
  }
})

const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push({ path: '/search', query: { keyword: searchKeyword.value } })
    searchKeyword.value = ''
    isMobileMenuOpen.value = false
  }
}

const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}

const navigateTo = (path) => {
  router.push(path)
  isMobileMenuOpen.value = false
}

const isActive = (path) => {
  return route.path === path
}

const isTagActive = () => {
  return route.path.startsWith('/tag')
}

const handleLogout = () => {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    userStore.logout()
    ElMessage.success('已退出登录~ 👋')
    router.push('/')
  }).catch(() => {
    // 用户取消
  })
}

const goToProfile = () => {
  userDropdownVisible.value = false
  router.push('/profile')
}

const goToLogin = () => {
  router.push('/login')
}

const goToRegister = () => {
  router.push('/register')
}

const toggleTheme = () => {
  themeStore.toggleTheme()
}


</script>

<template>
  <header class="header">
    <div class="container">
      <div class="header-inner">
        <div class="logo" @click="navigateTo('/')">
          <h1 data-text="博客平台">博客平台</h1>
        </div>

        <nav class="nav-desktop">
          <router-link to="/" :class="{ active: isActive('/') }">首页</router-link>
          <router-link to="/articles" :class="{ active: isActive('/articles') }">文章</router-link>
          <router-link to="/series" :class="{ active: isSeriesActive }">专栏</router-link>
          <router-link to="/categories" :class="{ active: isActive('/categories') }">分类</router-link>
          <router-link to="/tag" :class="{ active: isTagActive() }">标签</router-link>
          <router-link to="/about" :class="{ active: isActive('/about') }">关于</router-link>
        </nav>

        <div class="header-right">
          <!-- 主题切换按钮 -->
          <el-button
            class="theme-toggle"
            :icon="themeStore.isDark ? Sunny : Moon"
            circle
            @click="toggleTheme"
          />

          <div class="search-box">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索文章..."
              @keyup.enter="handleSearch"
              clearable
            >
              <template #suffix>
                <el-icon @click="handleSearch"><Search /></el-icon>
              </template>
            </el-input>
          </div>

          <!-- 未登录状态-->
          <div v-if="!userStore.isLoggedIn" class="auth-buttons">
            <!-- <el-button type="primary" > -->
            <HaloButton content="登录" size="large" @click="goToLogin">
              <el-icon><User /></el-icon>
              <span>登录</span>
            <!-- </el-button> -->
            </HaloButton>
          </div>

          <!-- 已登录状态-->
          <div v-else class="logged-in-section">
            <el-dropdown trigger="hover" @visible-change="userDropdownVisible = $event">
              <div class="user-info">
                <img :src="avatarUrl" :alt="userName" class="user-avatar" />
                <span class="user-name">{{ userName }}</span>
                <el-icon class="dropdown-arrow"><ArrowDown /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu class="user-dropdown-menu">
                  <el-dropdown-item @click="goToProfile">
                    <el-icon><User /></el-icon>
                    <span>个人中心</span>
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleLogout">
                    <el-icon><SwitchButton /></el-icon>
                    <span>退出登录</span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>

          <el-button
            class="mobile-menu-btn"
            :icon="Menu"
            circle
            @click="toggleMobileMenu"
          />
        </div>
      </div>
    </div>

    <!-- 移动端菜单-->
    <div class="mobile-menu" :class="{ open: isMobileMenuOpen }">
      <router-link to="/" @click="navigateTo('/')">首页</router-link>
      <router-link to="/articles" @click="navigateTo('/articles')">文章</router-link>
      <router-link to="/categories" @click="navigateTo('/categories')" :class="{ active: isActive('/categories') }">分类</router-link>
      <router-link to="/tag" @click="navigateTo('/tag')" :class="{ active: isTagActive() }">标签</router-link>
      <router-link to="/about" @click="navigateTo('/about')">关于</router-link>

      <div class="mobile-auth-section" v-if="!userStore.isLoggedIn">
        <router-link to="/login" @click="navigateTo('/login')" class="mobile-login-btn">
          <el-icon><User /></el-icon>
          登录
        </router-link>
        <router-link to="/register" @click="navigateTo('/register')" class="mobile-register-btn">注册</router-link>
      </div>

      <div class="mobile-user-section" v-else>
        <div class="mobile-user-info">
          <img :src="avatarUrl" :alt="userName" class="mobile-user-avatar" />
          <div class="mobile-user-details">
            <div class="mobile-user-name">{{ userName }}</div>
            <div class="mobile-user-email">{{ userEmail }}</div>
          </div>
        </div>
        <router-link to="/profile" @click="navigateTo('/profile')" class="mobile-profile-link">个人中心</router-link>
        <div class="mobile-logout-link" @click="handleLogout">
          <el-icon><SwitchButton /></el-icon>
          退出登录
        </div>
      </div>
    </div>
  </header>
</template>

<style scoped>
.header {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  position: sticky;
  top: 0;
  z-index: 1000;
  border-bottom: 1px solid rgba(255, 255, 255, 0.3);
  transition: background 0.3s ease, border-color 0.3s ease;
}

html.dark .header {
  background: rgba(26, 32, 44, 0.85);
  border-bottom-color: rgba(255, 255, 255, 0.1);
}

.header-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 80px;
}

.logo {
  cursor: pointer;
  flex-shrink: 0;
  transition: transform 0.3s ease;
}

.logo h1 {
  position: relative;
  font-size: 26px;
  font-weight: 700;
  margin: 0;
  letter-spacing: -0.5px;
}

.logo h1::before {
  content: attr(data-text);
  position: absolute;
  top: 0;
  left: 0;
  color: skyblue;
}

html.dark .logo h1::before {
  color: skyblue;
}

.logo h1::after {
  content: attr(data-text);
  color: transparent;
  position: absolute;
  top: 0;
  left: 0;
  clip-path: circle(20px at 0% 50%);
  animation: movelight 5s infinite;
  background: url(/images/1.jpg) center;
  background-size: 200%;
  background-clip: text;
  -webkit-background-clip: text;
}

@keyframes movelight {
  0% {
    clip-path: circle(30px at 0% 50%);
  }
  50% {
    clip-path: circle(30px at 100% 50%);
  }
  100% {
    clip-path: circle(30px at 0% 50%);
  }
}

.logo:hover {
  transform: scale(1.05);
}

/* .logo h1 {
  font-size: 26px;
  font-weight: 700;
  background: skyblue;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
  letter-spacing: -0.5px;
} */

.nav-desktop {
  display: flex;
  gap: 40px;
  align-items: center;
}

.nav-desktop a {
  color: #4a5568;
  font-size: 15px;
  font-weight: 500;
  transition: all 0.3s ease;
  position: relative;
  padding: 8px 0;
}

html.dark .nav-desktop a {
  color: #a0aec0;
}

.nav-desktop a::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  width: 0;
  height: 2px;
  background: skyblue;
  transition: all 0.3s ease;
  transform: translateX(-50%);
}

.nav-desktop a:hover,
.nav-desktop a.active {
  color: #70d2b0;
}

.nav-desktop a:hover::after,
.nav-desktop a.active::after {
  width: 100%;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.search-box {
  width: 240px;
  transition: width 0.3s;
}

.search-box :deep(.el-input__wrapper) {
  border-radius: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.9);
  transition: all 0.3s;
}

html.dark .search-box :deep(.el-input__wrapper) {
  background: rgba(45, 55, 72, 0.9);
  border-color: rgba(255, 255, 255, 0.1);
}

html.dark .search-box :deep(.el-input__inner) {
  color: #e2e8f0;
}

html.dark .search-box :deep(.el-input__inner::placeholder) {
  color: #a0aec0;
}

.search-box :deep(.el-input__wrapper:hover),
.search-box :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
  border-color: #667eea;
}

.auth-buttons {
  display: flex;
  gap: 10px;
  align-items: center;
}

.auth-buttons :deep(.el-button) {
  padding: 8px 20px;
  border-radius: 10px;
  font-weight: 500;
  background-color: skyblue;
  border: none;
}

.logged-in-section {
  display: flex;
  align-items: center;
  gap: 16px;
}

.create-btn {
  background: skyblue;
  border: none;
  font-weight: 600;
  padding: 10px 24px;
  border-radius: 20px;
}

.create-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(102, 126, 234, 0.35);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 12px;
  transition: all 0.3s;
  background: rgba(255, 255, 255, 0.5);
}

html.dark .user-info {
  background: rgba(255, 255, 255, 0.1);
}

.user-info:hover {
  background: rgba(102, 126, 234, 0.1);
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #667eea;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.2);
}

.user-name {
  font-size: 15px;
  font-weight: 600;
  color: #4a5568;
  max-width: 100px;
  /* overflow: hidden; */
  text-overflow: ellipsis;
  white-space: nowrap;
}

html.dark .user-name {
  color: #e2e8f0;
}

.theme-toggle {
  background: skyblue;
  border: none;
  color: white;
  transition: all 0.3s;
}

.theme-toggle:hover {
  transform: rotate(180deg);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.dropdown-arrow {
  font-size: 14px;
  color: #a0aec0;
  transition: transform 0.3s;
}

.user-info:hover .dropdown-arrow {
  transform: rotate(180deg);
}

/* 下拉菜单样式 */
.user-dropdown-menu {
  min-width: 200px;
  padding: 0;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  border: none;
}

.user-dropdown-menu :deep(.el-dropdown-menu__item) {
  padding: 12px 20px;
  font-size: 14px;
  color: #4a5568;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-dropdown-menu :deep(.el-dropdown-menu__item .el-icon) {
  font-size: 16px;
  flex-shrink: 0;
}

.user-dropdown-menu :deep(.el-dropdown-menu__item span) {
  flex: 1;
}

.user-dropdown-menu :deep(.el-dropdown-menu__item:hover) {
  background: linear-gradient(135deg, #f0f4ff 0%, #f5f3ff 100%);
  color: #667eea;
}

.user-dropdown-menu :deep(.el-dropdown-menu__item.is-divided) {
  border-top: 1px solid #e2e8f0;
  margin-top: 4px;
}

html.dark .user-dropdown-menu :deep(.el-dropdown-menu__item) {
  color: #a0aec0;
}

html.dark .user-dropdown-menu :deep(.el-dropdown-menu__item:hover) {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: #fff;
}

html.dark .user-dropdown-menu :deep(.el-dropdown-menu__item.is-divided) {
  border-top-color: rgba(255, 255, 255, 0.1);
}

.user-dropdown-header {
  padding: 16px;
  background: linear-gradient(135deg, #f0f4ff 0%, #ede9fe 100%);
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.dropdown-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: 3px solid #fff;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  flex-shrink: 0;
}

.dropdown-user-info {
  flex: 1;
  overflow: hidden;
}

.dropdown-username {
  font-size: 16px;
  font-weight: 700;
  color: #1a202c;
  margin-bottom: 4px;
}

html.dark .dropdown-username {
  color: #f7fafc;
}

.dropdown-email {
  font-size: 13px;
  color: #718096;
}

html.dark .dropdown-email {
  color: #a0aec0;
}

.mobile-menu-btn {
  display: none;
  background: skyblue;
  border: none;
  color: white;
}

.mobile-menu-btn:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.mobile-auth-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 12px 0;
  border-top: 1px solid #e2e8f0;
  border-bottom: 1px solid #e2e8f0;
  margin: 12px 0;
}

.mobile-login-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-weight: 600;
}

.mobile-register-btn {
  background: linear-gradient(135deg, #667eea 0%, #f093fb 100%);
  color: #fff;
  font-weight: 600;
  padding: 14px 20px;
  text-align: center;
}

.mobile-user-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 12px 0;
  border-top: 1px solid #e2e8f0;
  border-bottom: 1px solid #e2e8f0;
  margin: 12px 0;
}

.mobile-user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: linear-gradient(135deg, #f0f4ff 0%, #ede9fe 100%);
  border-radius: 12px;
}

.mobile-user-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: 2px solid #667eea;
}

.mobile-user-details {
  flex: 1;
}

.mobile-user-name {
  font-size: 16px;
  font-weight: 700;
  color: #1a202c;
  margin-bottom: 4px;
}

.mobile-user-email {
  font-size: 13px;
  color: #718096;
}

.mobile-create-link {
  background: skyblue;
  color: #fff;
  font-weight: 600;
  text-align: center;
  padding: 12px 20px;
}

.mobile-profile-link {
  background: linear-gradient(135deg, #667eea 0%, #f093fb 100%);
  color: #fff;
  font-weight: 600;
  text-align: center;
  padding: 12px 20px;
}

.mobile-logout-link {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #f56565;
  font-weight: 500;
  padding: 12px 20px;
}

.mobile-menu {
  position: fixed;
  top: 80px;
  left: 0;
  right: 0;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  transform: translateY(-100%);
  opacity: 0;
  visibility: hidden;
  transition: all 0.4s ease;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  border-top: 1px solid rgba(255, 255, 255, 0.3);
}

html.dark .mobile-menu {
  background: rgba(26, 32, 44, 0.98);
  border-top-color: rgba(255, 255, 255, 0.1);
}

.mobile-menu.open {
  transform: translateY(0);
  opacity: 1;
  visibility: visible;
}

.mobile-menu a {
  color: #4a5568;
  font-size: 16px;
  font-weight: 500;
  padding: 14px 20px;
  border-radius: 12px;
  transition: all 0.3s;
}

html.dark .mobile-menu a {
  color: #a0aec0;
}

.mobile-menu a:hover,
.mobile-menu a.active {
  background: linear-gradient(135deg, #f0f4ff 0%, #f5f3ff 100%);
  color: #667eea;
  transform: translateX(8px);
}

@media (max-width: 768px) {
  .nav-desktop {
    display: none;
  }

  .search-box {
    display: none;
  }

  .auth-buttons {
    display: none;
  }

  .user-info {
    display: none;
  }

  .logged-in-section {
    gap: 12px;
  }

  .create-btn {
    padding: 8px 16px;
    font-size: 13px;
  }

  .mobile-menu-btn {
    display: flex;
  }
}
</style>

