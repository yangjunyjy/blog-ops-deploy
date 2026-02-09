<template>
  <div class="tags-view">
    <el-scrollbar class="tags-scrollbar" ref="scrollbarRef" @mousemove="handleMouseMove" @mouseleave="handleMouseLeave">
      <div class="tags-list">
        <div
          v-for="tag in visitedViews"
          :key="tag.path"
          :ref="el => setTagRef(el, tag.path)"
          :class="['tag-item', { active: isActive(tag.path) }]"
          @click="handleTagClick(tag.path)"
        >
          <span class="tag-title">{{ tag.title }}</span>
          <el-icon
            v-if="!isAffix(tag)"
            class="tag-close"
            @click.stop="closeTag(tag.path)"
          >
            <Close />
          </el-icon>
        </div>
      </div>
    </el-scrollbar>
    <div class="tags-actions">
      <el-dropdown @command="handleTagCommand" trigger="click">
        <el-icon class="tags-action-btn">
          <Operation />
        </el-icon>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="closeOthers">关闭其他</el-dropdown-item>
            <el-dropdown-item command="closeAll">关闭所有</el-dropdown-item>
            <el-dropdown-item command="refresh">刷新当前</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Close, Operation } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

const visitedViews = ref([])
const affixTags = ref(['/dashboard'])
const scrollbarRef = ref(null)
const tagRefs = ref({})
const scrollTimer = ref(null)

// 判断是否为固定标签
const isAffix = (tag) => {
  return affixTags.value.includes(tag.path)
}

// 判断是否激活
const isActive = (path) => {
  return path === route.path
}

// 设置标签引用
const setTagRef = (el, path) => {
  if (el) {
    tagRefs.value[path] = el
  }
}

// 滚动到当前激活的标签
const scrollToActiveTag = () => {
  nextTick(() => {
    const activeTag = tagRefs.value[route.path]

    if (activeTag) {
      // 获取滚动容器
      const scrollbarWrap = document.querySelector('.tags-scrollbar .el-scrollbar__wrap')

      if (scrollbarWrap) {
        const scrollLeft = scrollbarWrap.scrollLeft
        const clientWidth = scrollbarWrap.clientWidth
        const tagLeft = activeTag.offsetLeft
        const tagWidth = activeTag.offsetWidth

        // 将标签滚动到视口中间
        const targetScrollLeft = tagLeft - (clientWidth - tagWidth) / 2

        scrollbarWrap.scrollTo({
          left: targetScrollLeft,
          behavior: 'smooth'
        })
      }
    }
  })
}

// 鼠标移动事件处理
const handleMouseMove = (event) => {
  const scrollbarWrap = document.querySelector('.tags-scrollbar .el-scrollbar__wrap')

  if (!scrollbarWrap) return

  const rect = scrollbarWrap.getBoundingClientRect()
  const mouseX = event.clientX - rect.left
  const edgeSize = 50 // 左右边缘区域大小

  // 清除之前的动画
  if (scrollTimer.value) {
    cancelAnimationFrame(scrollTimer.value)
    scrollTimer.value = null
  }

  // 根据鼠标位置计算滚动速度
  const scrollSpeed = (offset, edgeSize) => {
    const progress = 1 - (offset / edgeSize)
    return progress * progress * 8 // 速度随接近边缘而增加
  }

  const animate = () => {
    const rect = scrollbarWrap.getBoundingClientRect()
    const mouseX = event.clientX - rect.left

    // 鼠标在左侧边缘，向左滚动（显示左边隐藏的）
    if (mouseX < edgeSize && scrollbarWrap.scrollLeft > 0) {
      scrollbarWrap.scrollLeft -= scrollSpeed(mouseX, edgeSize)
      scrollTimer.value = requestAnimationFrame(animate)
    }
    // 鼠标在右侧边缘，向右滚动（显示右边隐藏的）
    else if (mouseX > rect.width - edgeSize && scrollbarWrap.scrollLeft < scrollbarWrap.scrollWidth - rect.width) {
      scrollbarWrap.scrollLeft += scrollSpeed(rect.width - mouseX, edgeSize)
      scrollTimer.value = requestAnimationFrame(animate)
    }
  }

  // 开始动画
  scrollTimer.value = requestAnimationFrame(animate)
}

// 鼠标离开事件处理
const handleMouseLeave = () => {
  if (scrollTimer.value) {
    cancelAnimationFrame(scrollTimer.value)
    scrollTimer.value = null
  }
}

// 点击标签
const handleTagClick = (path) => {
  if (path !== route.path) {
    router.push(path)
  }
}

// 关闭标签
const closeTag = (path) => {
  const index = visitedViews.value.findIndex(v => v.path === path)
  if (index > -1) {
    let nextPath = null
    if (path === route.path && visitedViews.value.length > 1) {
      nextPath = visitedViews.value[index - 1]?.path || visitedViews.value[index + 1]?.path
    }

    visitedViews.value.splice(index, 1)

    if (nextPath) {
      router.push(nextPath)
    } else if (path === route.path && visitedViews.value.length === 0) {
      router.push('/dashboard')
    }
  }
}

// 标签操作
const handleTagCommand = (command) => {
  switch (command) {
    case 'closeOthers':
      closeOthersTags()
      break
    case 'closeAll':
      closeAllTags()
      break
    case 'refresh':
      refreshCurrentPage()
      break
  }
}

// 关闭其他标签
const closeOthersTags = () => {
  const currentTag = visitedViews.value.find(v => v.path === route.path)
  if (currentTag) {
    visitedViews.value = affixTags.value.map(path => ({ path, title: getTitle(path) }))
    if (!affixTags.value.includes(route.path)) {
      visitedViews.value.push(currentTag)
    }
  }
}

// 关闭所有标签
const closeAllTags = () => {
  visitedViews.value = affixTags.value.map(path => ({ path, title: getTitle(path) }))
  router.push('/dashboard')
}

// 刷新当前页面
const refreshCurrentPage = () => {
  router.go(0)
}

// 获取路由标题
const getTitle = (path) => {
  const titles = {
    '/dashboard': '首页',
    '/dashboard/index': '首页',
    '/articles': '文章管理',
    '/categories': '分类管理',
    '/tags': '标签管理',
    '/series': '系列管理',
    '/comments': '评论管理',
    '/users': '用户管理',
    '/statistics/overview': '数据概览',
    '/statistics/content': '内容统计',
    '/statistics/user': '用户统计'
  }
  return titles[path] || '未知'
}

// 添加访问的路由
const addVisitedView = (route) => {
  const title = route.meta?.title || getTitle(route.path)

  // 跳过空标题（如 /dashboard 父路由）
  if (!title || title === '') {
    return
  }

  if (!visitedViews.value.some(v => v.path === route.path)) {
    visitedViews.value.push({
      path: route.path,
      title: title
    })
  }
}

// 监听路由变化
watch(
  () => route.path,
  (newPath, oldPath) => {
    addVisitedView(route)
    // 滚动到当前激活的标签
    if (newPath !== oldPath) {
      scrollToActiveTag()
    }
  },
  { immediate: true }
)
</script>

<style scoped lang="scss">
.tags-view {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
  height: 50px;

  .tags-scrollbar {
    flex: 1;
    overflow: hidden;
  }

  :deep(.tags-scrollbar .el-scrollbar__wrap) {
    overflow-x: auto;
    overflow-y: hidden;
  }

  :deep(.tags-scrollbar .el-scrollbar__bar.is-horizontal) {
    height: 4px;

    .el-scrollbar__thumb {
      background: rgba(0, 0, 0, 0.2);
      border-radius: 2px;
    }
  }

  .tags-list {
    display: flex;
    align-items: center;
    gap: 6px;
    white-space: nowrap;
    height: 100%;
    padding: 0 0;

    .tag-item {
      margin-top: 3px;
      display: inline-flex;
      align-items: center;
      gap: 3px;
      padding: 4px 8px;
      font-size: 13px;
      color: #303030;
      background: #f8efef;
      border: 1px solid #f0f0f0;
      border-radius: 6px;
      cursor: pointer;
      transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
      white-space: nowrap;
      flex-shrink: 0;

      &:hover {
        background: linear-gradient(135deg, #9fb041 0%, #5cc8cc 100%);
        border-color: white;
        color: #7e5580;
        transform: translateY(-1px);
        box-shadow: 0 1px 3px rgba(24, 144, 255, 0.15);

        .tag-close {
          opacity: 1;
        }
      }

      &.active {
        background: skyblue;
        border-color: skyblue;
        color: #fff;
        font-weight: 500;

        .tag-close {
          color: #fff;
          opacity: 1;
        }
      }

      .tag-title {
        max-width: 100px;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .tag-close {
        margin-left: 4px;
        font-size: 12px;
        opacity: 0;
        transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);

        &:hover {
          transform: scale(1.1);
        }
      }
    }
  }

  .tags-actions {
    margin-left: 12px;
    flex-shrink: 0;

    .tags-action-btn {
      font-size: 16px;
      cursor: pointer;
      color: #8c8c8c;
      transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
      padding: 6px;
      border-radius: 2px;

      &:hover {
        color: #1890ff;
        background: rgba(24, 144, 255, 0.08);
        transform: rotate(90deg);
      }
    }
  }
}
</style>
