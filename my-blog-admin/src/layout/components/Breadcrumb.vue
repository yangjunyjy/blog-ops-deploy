<template>
  <el-breadcrumb separator="/">
    <transition-group name="breadcrumb">
      <el-breadcrumb-item v-for="(item, index) in levelList" :key="item.path">
        <span v-if="index === levelList.length - 1" class="no-redirect">
          {{ item.meta.title }}
        </span>
        <span v-else>
          {{ item.meta.title }}
        </span>
      </el-breadcrumb-item>
    </transition-group>
  </el-breadcrumb>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const levelList = ref([])

const getBreadcrumb = () => {
  let matched = route.matched.filter(item => item.meta && item.meta.title)

  // 过滤掉 dashboard 路由本身（父路由）
  matched = matched.filter(item => item.path !== '/dashboard')

  // 如果没有匹配的路由，说明当前路由就是首页，显示首页
  if (matched.length === 0) {
    matched = [{ path: '/dashboard', meta: { title: '首页' } }]
  } else {
    // 如果第一个路由不是首页，添加首页作为第一项
    const first = matched[0]
    // 判断是否为首页路由（包括 /dashboard 和 /dashboard/index）
    const isDashboardRoute = first.path === '/dashboard' ||
                           first.path === '/dashboard/index' ||
                           first.name === 'Dashboard'

    if (!isDashboardRoute) {
      matched = [{ path: '/dashboard', meta: { title: '首页' } }].concat(matched)
    }
  }

  levelList.value = matched.filter(item => item.meta && item.meta.title && item.meta.breadcrumb !== false)
}

const handleLink = (item) => {
  const { redirect, path } = item
  if (redirect) {
    router.push(redirect)
    return
  }
  router.push(path)
}

watch(route, getBreadcrumb, { immediate: true })
</script>

<style scoped lang="scss">
.el-breadcrumb {
  display: inline-block;
  font-size: 14px;
  margin-left: 12px;

  :deep(.el-breadcrumb__separator) {
    margin: 0 8px;
    color: #d9d9d9;
  }

  :deep(.el-breadcrumb__inner) {
    font-weight: 400;
    color: #595959;
    transition: all 0.3s;

    &.is-link {
      &:hover {
        color: #1890ff;
      }
    }
  }
}

.no-redirect {
  color: #262626;
  font-weight: 500;
}

.breadcrumb-enter-active,
.breadcrumb-leave-active {
  transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
}

.breadcrumb-enter-from,
.breadcrumb-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
</style>
