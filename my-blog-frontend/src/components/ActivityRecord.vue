<template>
  <div class="activity-record">
    <!-- 筛选标签 -->
    <div class="activity-filters">
      <button
        class="filter-btn"
        :class="{ active: activeFilter === 'all' }"
        @click="handleFilterChange('all')"
      >
        全部 ({{ stats.total }})
      </button>
      <button
        class="filter-btn"
        :class="{ active: activeFilter === 'like' }"
        @click="handleFilterChange('like')"
      >
        <el-icon><Star /></el-icon>
        点赞 ({{ stats.likes }})
      </button>
      <button
        class="filter-btn"
        :class="{ active: activeFilter === 'comment' }"
        @click="handleFilterChange('comment')"
      >
        <el-icon><ChatDotRound /></el-icon>
        评论 ({{ stats.comments }})
      </button>
      <button
        class="filter-btn"
        :class="{ active: activeFilter === 'share' }"
        @click="handleFilterChange('share')"
      >
        <el-icon><Share /></el-icon>
        分享 ({{ stats.shares }})
      </button>
    </div>

    <!-- 活动列表 -->
    <div class="activity-list">
      <div
        v-for="item in filteredActivities"
        :key="item.id"
        class="activity-item"
      >
        <div class="activity-icon" :class="`activity-${item.type}`">
          <el-icon v-if="item.type === 'like'"><Star /></el-icon>
          <el-icon v-else-if="item.type === 'comment'"><ChatDotRound /></el-icon>
          <el-icon v-else-if="item.type === 'share'"><Share /></el-icon>
          <el-icon v-else-if="item.type === 'favorite'"><Collection /></el-icon>
        </div>
        <div class="activity-content" @click="handleGoToArticle(item.article_id)">
          <div class="activity-header">
            <span class="activity-type">
              <el-icon v-if="item.type === 'like'"><Star /></el-icon>
              <el-icon v-else-if="item.type === 'comment'"><ChatDotRound /></el-icon>
              <el-icon v-else-if="item.type === 'share'"><Share /></el-icon>
              <el-icon v-else-if="item.type === 'favorite'"><Collection /></el-icon>
              {{ getTypeText(item.type) }}
            </span>
            <span class="activity-time">{{ formatTime(item.created_at) }}</span>
          </div>
          <h4 class="activity-article-title">{{ item.article_title }}</h4>
          <p v-if="item.content" class="activity-comment">{{ item.content }}</p>
          <p v-if="item.platform" class="activity-platform">分享到：{{ item.platform }}</p>
          <p v-if="item.folder_name" class="activity-folder">收藏到：{{ item.folder_name }}</p>
        </div>
        <el-icon class="activity-arrow"><Right /></el-icon>
      </div>
      <el-empty v-if="!filteredActivities.length" description="暂无互动记录" />
    </div>

    <!-- 分页 -->
    <Pagination
      v-if="pagination.total > 0"
      :current-page="pagination.page"
      :total="pagination.total"
      :page-size="pagination.pageSize"
      :total-pages="totalPages"
      @change="handlePageChange"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { Star, ChatDotRound, Share, Collection, Right } from '@element-plus/icons-vue'
import Pagination from './Pagination.vue'
import { getUserActivities, getUserActivityStats } from '@/api'

const props = defineProps({
  userId: {
    type: Number,
    required: true
  },
  pageSize: {
    type: Number,
    default: 10
  }
})

const router = useRouter()

// 活动记录列表
const activities = ref([])

// 筛选条件
const activeFilter = ref('all')

// 分页
const pagination = ref({
  page: 1,
  total: 0,
  pageSize: props.pageSize
})

// 统计数据
const stats = ref({
  total: 0,
  likes: 0,
  comments: 0,
  shares: 0,
  favorites: 0
})

// 计算总页数
const totalPages = computed(() => {
  return Math.ceil(pagination.value.total / pagination.value.pageSize) || 1
})

// 过滤后的活动记录
const filteredActivities = computed(() => {
  let filtered = activities.value

  // 按类型过滤
  if (activeFilter.value !== 'all') {
    filtered = filtered.filter(item => item.type === activeFilter.value)
  }

  // 分页
  const startIndex = (pagination.value.page - 1) * pagination.value.pageSize
  const endIndex = startIndex + pagination.value.pageSize
  return filtered.slice(startIndex, endIndex)
})

// 获取类型文本
const getTypeText = (type) => {
  const typeMap = {
    like: '点赞',
    comment: '评论',
    share: '分享',
    favorite: '收藏'
  }
  return typeMap[type] || '互动'
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return ''
  // 转换为 Date 对象（支持字符串和 Date 对象）
  const date = typeof time === 'string' ? new Date(time) : time
  if (isNaN(date.getTime())) return ''

  const now = new Date()
  const diff = now - date

  // 如果是今天，只显示时间
  if (diff < 86400000 && date.getDate() === now.getDate()) {
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  // 否则显示日期
  return date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })
}

// 跳转到文章详情
const handleGoToArticle = (articleId) => {
  router.push(`/article/${articleId}`)
}

// 筛选变化
const handleFilterChange = (type) => {
  activeFilter.value = type
  pagination.value.page = 1
}

// 分页变化
const handlePageChange = (page) => {
  pagination.value.page = page
}

// 加载活动记录
const loadActivities = async () => {
  try {
    const res = await getUserActivities(props.userId, {
      page: pagination.value.page,
      perPage: pagination.value.pageSize,
      type: activeFilter.value
    })

    if (res.code === 200 && res.data) {
      activities.value = res.data.list || []
      pagination.value.total = res.data.total || 0

      // 如果有统计数据，更新统计
      if (res.data.stats) {
        stats.value = {
          total: res.data.stats.total || 0,
          likes: res.data.stats.like || 0,
          comments: res.data.stats.comment || 0,
          shares: res.data.stats.share || 0,
          favorites: res.data.stats.favorite || 0
        }
      }
    }
  } catch (error) {
    console.error('加载活动记录失败:', error)
  }
}

// 加载统计数据
const loadStats = async () => {
  try {
    const res = await getUserActivityStats(props.userId)
    if (res.code === 200 && res.data) {
      stats.value = {
        total: res.data.total || 0,
        likes: res.data.like || 0,
        comments: res.data.comment || 0,
        shares: res.data.share || 0,
        favorites: res.data.favorite || 0
      }
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

// 监听筛选条件变化
watch(activeFilter, () => {
  loadActivities()
})

// 监听用户 ID 变化
watch(() => props.userId, () => {
  loadActivities()
  loadStats()
}, { immediate: true })

// 组件挂载时加载
onMounted(() => {
  loadActivities()
  loadStats()
})

// 暴露方法供父组件调用
defineExpose({
  loadActivities,
  loadStats,
  refresh: () => {
    loadActivities()
    loadStats()
  }
})
</script>

<style scoped>
.activity-record {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 筛选标签样式 */
.activity-filters {
  display: flex;
  gap: 10px;
  padding-bottom: 15px;
  border-bottom: 1px solid #e4e7ed;
}

.filter-btn {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 8px 16px;
  background: #f5f7fa;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  color: #606266;
  cursor: pointer;
  transition: all 0.3s;
}

.filter-btn:hover {
  background: #ecf5ff;
  color: #409eff;
}

.filter-btn.active {
  background: #409eff;
  color: #fff;
}

.filter-btn .el-icon {
  font-size: 16px;
}

/* 活动列表样式 */
.activity-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: 15px;
  padding: 16px;
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 12px;
  transition: all 0.3s;
  cursor: pointer;
}

.activity-item:hover {
  border-color: #409eff;
  box-shadow: 0 2px 12px rgba(64, 158, 255, 0.1);
  transform: translateX(4px);
}

/* 活动图标 */
.activity-icon {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  font-size: 20px;
}

.activity-icon.activity-like {
  background: linear-gradient(135deg, #ffeaa7, #fdcb6e);
  color: #fff;
}

.activity-icon.activity-comment {
  background: linear-gradient(135deg, #74b9ff, #0984e3);
  color: #fff;
}

.activity-icon.activity-share {
  background: linear-gradient(135deg, #a29bfe, #6c5ce7);
  color: #fff;
}

.activity-icon.activity-favorite {
  background: linear-gradient(135deg, #fd79a8, #e84393);
  color: #fff;
}

/* 活动内容 */
.activity-content {
  flex: 1;
  min-width: 0;
}

.activity-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.activity-type {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 14px;
  font-weight: 500;
  color: #409eff;
}

.activity-type .el-icon {
  font-size: 14px;
}

.activity-time {
  font-size: 12px;
  color: #909399;
}

.activity-article-title {
  margin: 8px 0;
  font-size: 15px;
  font-weight: 600;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.activity-comment,
.activity-platform,
.activity-folder {
  margin: 6px 0 0;
  font-size: 13px;
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.activity-comment {
  padding: 8px 12px;
  background: #f5f7fa;
  border-left: 3px solid #409eff;
  border-radius: 4px;
}

.activity-platform {
  color: #e67e22;
}

.activity-folder {
  color: #e84393;
}

/* 箭头 */
.activity-arrow {
  flex-shrink: 0;
  margin-top: 10px;
  font-size: 16px;
  color: #c0c4cc;
  transition: transform 0.3s;
}

.activity-item:hover .activity-arrow {
  transform: translateX(4px);
  color: #409eff;
}

/* 响应式 */
@media (max-width: 768px) {
  .activity-filters {
    flex-wrap: wrap;
  }

  .filter-btn {
    flex: 1;
    min-width: calc(50% - 5px);
  }

  .activity-item {
    padding: 12px;
  }

  .activity-icon {
    width: 36px;
    height: 36px;
    font-size: 18px;
  }
}
</style>
