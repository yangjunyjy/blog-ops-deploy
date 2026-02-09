<script setup>
import { formatDate, truncateText } from '@/utils/format'
import { View, Calendar, Star } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, inject } from 'vue'
import { useRouter } from 'vue-router'
import { favoriteArticle, unfavoriteArticle, getFavoriteFolders } from '@/api'

const props = defineProps({
  article: {
    type: Object,
    required: true
  }
})

const router = useRouter()
const userStore = useUserStore()
const isFavorited = ref(false)
const isLiked = ref(false)
const likesCount = ref(props.article.likes || 0)

// 收藏文件夹相关
const folderDialogVisible = ref(false)
const selectedFolderId = ref(null)
const folders = ref([])

const goToDetail = () => {
  const url = router.resolve({
    name: 'ArticleDetail',
    params: { id: props.article.id }
  }).href
  window.open(url, '_blank')
}

const handleLike = (event) => {
  event.stopPropagation()
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }
  isLiked.value = !isLiked.value
  likesCount.value += isLiked.value ? 1 : -1
  ElMessage.success(isLiked.value ? '已点赞' : '已取消点赞')
}

const handleFavorite = async (event) => {
  event.stopPropagation()
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  if (!isFavorited.value) {
    // 收藏：检查是否需要选择文件夹
    try {
      const res = await getFavoriteFolders()
      if (res.code === 200 && res.data) {
        folders.value = res.data.list

        // 如果只有一个默认文件夹，直接收藏
        if (folders.value.length === 1 && folders.value[0].isDefault) {
          await performFavorite(folders.value[0].id)
        } else if (folders.value.length === 1) {
          // 只有一个非默认文件夹
          await performFavorite(folders.value[0].id)
        } else {
          // 多个文件夹，弹出选择对话框
          folderDialogVisible.value = true
          selectedFolderId.value = folders.value[0]?.id
        }
      }
    } catch (error) {
      console.error('加载文件夹失败', error)
      ElMessage.error('操作失败')
    }
  } else {
    // 取消收藏
    try {
      await unfavoriteArticle({
        articleId: props.article.id,
        folderId: selectedFolderId.value
      })
      isFavorited.value = false
      ElMessage.success('已取消收藏')
    } catch (error) {
      ElMessage.error('操作失败')
    }
  }
}

const performFavorite = async (folderId) => {
  try {
    await favoriteArticle({
      articleId: props.article.id,
      folderId: folderId
    })
    isFavorited.value = true
    selectedFolderId.value = folderId
    ElMessage.success('已收藏')
    folderDialogVisible.value = false
  } catch (error) {
    ElMessage.error('收藏失败')
  }
}
</script>

<template>
  <div class="article-card" @click="goToDetail">
    <div class="article-cover" v-if="article.cover">
      <img v-lazy="article.cover" :alt="article.title" />
    </div>

    <div class="article-content">
      <div class="article-meta">
        <span class="category">{{ article.category?.name || '未分类' }}</span>
        <span class="date">
          <el-icon><Calendar /></el-icon>
          {{ formatDate(article.createdAt) }}
        </span>
      </div>

      <h3 class="article-title">{{ article.title }}</h3>

      <p class="article-summary">{{ truncateText(article.summary, 120) }}</p>

      <div class="article-tags" v-if="article.tags && article.tags.length">
        <el-tag
          v-for="tag in article.tags.slice(0, 3)"
          :key="tag.name"
          size="small"
          type="info"
          effect="plain"
        >
          {{ tag.name }}
        </el-tag>
      </div>

      <div class="article-footer">
        <div class="article-stats">
          <span class="views">
            <el-icon><View /></el-icon>
            {{ article.views || 0 }}
          </span>
          <span class="likes" @click="handleLike" :class="{ liked: isLiked }">
            <el-icon><Star /></el-icon>
            {{ likesCount }}
          </span>
        </div>
        <div class="article-actions">
          <el-button
            text
            :icon="isFavorited ? Star : Star"
            :class="{ favorited: isFavorited }"
            @click="handleFavorite"
          >
            {{ isFavorited ? '已收藏' : '收藏' }}
          </el-button>
          <el-button type="primary" text @click="goToDetail">阅读更多</el-button>
        </div>
      </div>
    </div>
  </div>
  <!-- 文件夹选择对话框-->

  <el-dialog v-model="folderDialogVisible" title="选择收藏文件夹" width="400px">
    <div class="folder-select-container">
      <el-radio-group v-model="selectedFolderId" class="folder-list">
        <div v-for="folder in folders" :key="folder.id" class="folder-option">
          <el-radio :label="folder.id" :disabled="folder.isDefault && folders.length > 1">


            <div class="folder-option-content">

              <div class="folder-option-name">{{ folder.name }}</div>

              <div class="folder-option-count">{{ folder.articleCount || 0 }} 篇文章</div>
            </div>
          </el-radio>
        </div>
      </el-radio-group>
    </div>
    <template #footer>
      <el-button @click="folderDialogVisible = false">取消</el-button>
      <el-button type="primary" @click="performFavorite(selectedFolderId)">确定收藏</el-button>
    </template>
  </el-dialog>
</template>

<style scoped>
.article-card {
  background: #f8fafc;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
  cursor: pointer;
  height: 100%;
  display: flex;
  flex-direction: column;
  border: 1px solid rgba(255, 255, 255, 0.5);
  position: relative;
}

html.dark .article-card {
  background: #2d3748;
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.25);
}

.article-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.08) 0%, rgba(118, 75, 162, 0.08) 100%);
  opacity: 0;
  transition: opacity 0.3s;
  z-index: 1;
}

.article-card:hover {
  transform: translateY(-4px) scale(1.01);
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.15);
  border-color: rgba(102, 126, 234, 0.25);
}

.article-card:hover::before {
  opacity: 1;
}

.article-cover {
  width: 100%;
  height: 160px;
  overflow: hidden;
  position: relative;
}

.article-cover::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 40px;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.08));
  pointer-events: none;
}

.article-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.article-card:hover .article-cover img {
  transform: scale(1.08);
}

.article-content {
  padding: 16px;
  flex: 1;
  display: flex;
  flex-direction: column;
  position: relative;
  z-index: 2;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
  font-size: 12px;
}

.category {
  color: #667eea;
  font-weight: 600;
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.3px;
  padding: 3px 10px;
  background: linear-gradient(135deg, #f0f4ff 0%, #ede9fe 100%);
  border-radius: 8px;
}

.date {
  color: #a0aec0;
  display: flex;
  align-items: center;
  gap: 3px;
  font-size: 11px;
}

.article-title {
  font-size: 16px;
  font-weight: 700;
  margin: 0 0 10px;
  color: #1a202c;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  letter-spacing: -0.2px;
  transition: color 0.3s;
}

html.dark .article-title {
  color: #f7fafc;
}

.article-card:hover .article-title {
  color: #667eea;
}

.article-summary {
  color: #4a5568;
  font-size: 13px;
  line-height: 1.6;
  margin: 0 0 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

html.dark .article-summary {
  color: #a0aec0;
}

.article-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 14px;
}

.article-footer {
  margin-top: auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

html.dark .article-footer {
  border-top-color: #4a5568;
}

.article-stats {
  display: flex;
  align-items: center;
  gap: 12px;
}

.views {
  color: #a0aec0;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.likes {
  color: #a0aec0;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.likes:hover {
  color: #667eea;
}

.likes.liked {
  color: #f56565;
}

.likes.liked .el-icon {
  animation: likePulse 0.3s ease-out;
}

@keyframes likePulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.3); }
  100% { transform: scale(1); }
}

.article-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.article-actions .el-button.favorited {
  color: #f56565;
}

.article-actions .el-button.favorited:hover {
  color: #e53e3e;
}

/* 文件夹选择对话框样�?*/
.folder-select-container {
  padding: 8px 0;
}

.folder-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.folder-option {
  padding: 8px;
  border-radius: 8px;
  transition: background-color 0.3s;
}

.folder-option:hover {
  background-color: #f7fafc;
}

html.dark .folder-option:hover {
  background-color: #2d3748;
}

.folder-option-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex: 1;
}

.folder-option-name {
  font-size: 14px;
  color: #1a202c;
}

html.dark .folder-option-name {
  color: #e2e8f0;
}

.folder-option-count {
  font-size: 12px;
  color: #a0aec0;
}
</style>
