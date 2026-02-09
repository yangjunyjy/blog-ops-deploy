<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { User, Plus } from '@element-plus/icons-vue'

const props = defineProps({
  author: {
    type: Object,
    required: true
  }
})

const userStore = useUserStore()
const isFollowing = ref(false)
const followersCount = ref(props.author.followers || 0)

const canFollow = computed(() => {
  return userStore.isLoggedIn && userStore.user?.id !== props.author.id
})

const handleFollow = () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  isFollowing.value = !isFollowing.value
  followersCount.value += isFollowing.value ? 1 : -1

  if (isFollowing.value) {
    ElMessage.success(`已关�?${props.author.name}`)
  } else {
    ElMessage.success(`已取消关�?${props.author.name}`)
  }
}
</script>

<template>
  <div class="author-card">
    <div class="author-info">
      <div class="author-avatar-wrapper">
        <img :src="author.avatar" :alt="author.name" class="author-avatar" />
      </div>
      <div class="author-details">
        <h4 class="author-name">{{ author.name }}</h4>
        <p class="author-bio">{{ author.bio || '暂无简�? }}</p>
        <div class="author-stats">
          <span class="stat-item">
            <strong>{{ author.articles || 0 }}</strong> 文章
          </span>
          <span class="stat-item">
            <strong>{{ followersCount }}</strong> 粉丝
          </span>
        </div>
      </div>
    </div>
    <el-button
      v-if="canFollow"
      :type="isFollowing ? 'default' : 'primary'"
      :icon="isFollowing ? null : Plus"
      class="follow-btn"
      :class="{ 'is-following': isFollowing }"
      @click="handleFollow"
    >
      {{ isFollowing ? '已关�? : '关注' }}
    </el-button>
    <div v-else-if="userStore.isLoggedIn && userStore.user?.id === author.id" class="is-you">
      这就是我自己
    </div>
  </div>
</template>

<style scoped>
.author-card {
  background: #f8fafc;
  border-radius: 16px;
  padding: 20px;
  margin-top: 20px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.5);
  transition: all 0.3s;
}

html.dark .author-card {
  background: #2d3748;
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3);
}

.author-info {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.author-avatar-wrapper {
  flex-shrink: 0;
}

.author-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #667eea;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}

.author-details {
  flex: 1;
  min-width: 0;
}

.author-name {
  font-size: 16px;
  font-weight: 600;
  color: #1a202c;
  margin: 0 0 8px;
}

html.dark .author-name {
  color: #f7fafc;
}

.author-bio {
  font-size: 13px;
  color: #718096;
  margin: 0 0 12px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

html.dark .author-bio {
  color: #a0aec0;
}

.author-stats {
  display: flex;
  gap: 16px;
  font-size: 13px;
  color: #4a5568;
}

html.dark .author-stats {
  color: #718096;
}

.stat-item strong {
  color: #667eea;
  font-weight: 700;
  margin-right: 4px;
}

.follow-btn {
  width: 100%;
  transition: all 0.3s;
}

.follow-btn.is-following {
  background: #f7fafc;
  border-color: #e2e8f0;
  color: #718096;
}

html.dark .follow-btn.is-following {
  background: #1a202c;
  border-color: #4a5568;
  color: #a0aec0;
}

.is-you {
  text-align: center;
  font-size: 14px;
  color: #a0aec0;
  padding: 8px;
  background: #f7fafc;
  border-radius: 8px;
}

html.dark .is-you {
  background: #1a202c;
  color: #718096;
}
</style>
