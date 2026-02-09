<template>
  <div class="layout-container">
    <!-- 侧边栏 -->
    <sidebar :is-collapse="isCollapse" />

    <!-- 主体内容 -->
    <div class="main-container">
      <!-- 顶部导航 -->
      <navbar :is-collapse="isCollapse" @toggle-collapse="toggleCollapse" />

      <!-- 标签栏 -->
      <tags-view />

      <!-- 页面内容 -->
      <div class="app-main">
        <router-view v-slot="{ Component }">
          <keep-alive :include="cachedViews">
            <transition name="slide" mode="out-in">
              <component :is="Component" />
            </transition>
          </keep-alive>
        </router-view>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import Sidebar from './components/Sidebar.vue'
import Navbar from './components/Navbar.vue'
import TagsView from './components/TagsView.vue'

const isCollapse = ref(false)

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

// 需要缓存的组件名称列表
const cachedViews = computed(() => {
  return [
    'Dashboard',
    'ArticleList',
    'CategoryList',
    'TagList',
    'SeriesList',
    'CommentList',
    'UserList',
    'RoleList',
    'MenuList',
    'ArticleAudit',
    'CommentAudit',
    'UserAudit',
    'Settings',
    'Overview',
    'Content',
    'User'
  ]
})
</script>

<style lang="scss" scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background: #f5f5f5;
}

.layout-container {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #f5f5f5;
}

.app-main {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
  position: relative;
  min-height: 0;
  height: 0; // 关键：让 flex: 1 生效

  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-thumb {
    background: #d9d9d9;
    border-radius: 3px;
  }

  &::-webkit-scrollbar-thumb:hover {
    background: #bfbfbf;
  }

  &::-webkit-scrollbar-track {
    background: #f0f0f0;
    border-radius: 3px;
  }
}

// 组件切换时的平滑效果
.router-view-content {
  animation: fadeIn 0.2s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0.95;
  }

  to {
    opacity: 1;
  }
}

.slide-enter-active {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-leave-active {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  z-index: 1;
}

.slide-enter-from {
  opacity: 0;
  transform: translateX(60px);
}

.slide-leave-to {
  opacity: 0;
  transform: translateX(-60px);
}
</style>

<style lang="scss">
// 全局样式 - 优化卡片容器避免底部留白
.app-main {

  // 让所有 el-card 的内容区域更紧凑
  :deep(.el-card) {
    .el-card__body {
      // 确保卡片内容不会有过大的空白
      display: flex;
      flex-direction: column;
      min-height: 0;
    }
  }

  // 减少分页器和表格之间的间距
  :deep(.pagination-wrapper) {
    padding-top: 12px;
  }

  // 优化表格包装器，让内容紧凑
  :deep(.table-wrapper) {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
  }
}
</style>
