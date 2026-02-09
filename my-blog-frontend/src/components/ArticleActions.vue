<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Star, StarFilled, Share, Top, Sunny, Moon, Folder, Document } from '@element-plus/icons-vue'
import { likeArticle, unlikeArticle, favoriteArticle, unfavoriteArticle, getFavoriteFolders, checkFavoriteStatus, checkArticleLikeStatus } from '@/api'
import { useUserStore } from '@/stores/user'
import ShareDialog from './ShareDialog.vue'

const userStore = useUserStore()

const props = defineProps({
  articleId: {
    type: Number,
    required: true
  },
  likes: {
    type: Number,
    default: 0
  },
  favorites: {
    type: Number,
    default: 0
  },
  title: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:likes', 'update:favorites'])

const isLiked = ref(false)
const isFavorited = ref(false)
const shareDialogVisible = ref(false)
const folderDialogVisible = ref(false)
const selectedFolderId = ref(null)
const folders = ref([])
const folderPagination = ref({
  page: 1,
  pageSize: 5,
  total: 0
})

const articleTitle = computed(() => props.title || document.title)

const handleLike = async () => {
  try {
    if (isLiked.value) {
      await unlikeArticle(props.articleId)
      emit('update:likes', props.likes - 1)
    } else {
      await likeArticle(props.articleId)
      emit('update:likes', props.likes + 1)
      ElMessage.success('点赞成功')
    }
    isLiked.value = !isLiked.value
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleFavorite = async () => {
  try {
    if (isFavorited.value) {
      await unfavoriteArticle({
        articleId: props.articleId,
        folderId: selectedFolderId.value
      })
      emit('update:favorites', props.favorites - 1)
      ElMessage.success('已取消藏')
      isFavorited.value = !isFavorited.value
    } else {
      // 加载文件夹列表
      const res = await getFavoriteFolders({
        page:folderPagination.value.page,
        pageSize:folderPagination.value.pageSize
      })
      if (res.code === 200 && res.data) {
        folders.value = res.data.list

        // 如果只有一个默认文件夹，直接收藏
        if (folders.value.length === 1 && folders.value[0].isDefault) {
          await performFavorite(folders.value[0].id)
        } else if (folders.value.length === 1) {
          await performFavorite(folders.value[0].id)
        } else {
          // 多个文件夹，弹出选择对话框
          folderDialogVisible.value = true
          selectedFolderId.value = folders.value[0]?.id
        }
      }
    }
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const performFavorite = async (folderId) => {
  try {
    await favoriteArticle({
      articleId: props.articleId,
      folderId: folderId
    })
    emit('update:favorites', props.favorites + 1)
    ElMessage.success('收藏成功')
    isFavorited.value = true
    selectedFolderId.value = folderId
    folderDialogVisible.value = false
  } catch (error) {
    ElMessage.error('收藏失败')
  }
}

const handleShare = () => {
  shareDialogVisible.value = true
}

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const currentUrl = computed(() => window.location.href)

// 检查收藏状态
const checkArticleFavoriteStatus = async () => {
  // 未登录时不检查收藏状态
  if (!userStore.isLoggedIn) {
    return
  }

  try {
    const res = await checkFavoriteStatus(props.articleId)
    if (res.code === 200 && res.data) {
      isFavorited.value = res.data.isFavorited
      selectedFolderId.value = res.data.folderId
    }
  } catch (error) {
    // 只在控制台记录，不显示错误提示
    console.error('检查收藏状态失败:', error)
  }
}

// 检查点赞状态
const checkArticleLike = async () => {
  // 未登录时不检查点赞状态
  if (!userStore.isLoggedIn) {
    return
  }

  try {
    const res = await checkArticleLikeStatus(props.articleId)
    if (res.code === 200 && res.data) {
      isLiked.value = res.data.isLiked
    }
  } catch (error) {
    // 只在控制台记录，不显示错误提示
    console.error('检查点赞状态失败:', error)
  }
}

// 监听登录状态变化，登录成功后检查点赞和收藏状态
watch(() => userStore.isLoggedIn, (isLoggedIn) => {
  if (isLoggedIn) {
    checkArticleLike()
    checkArticleFavoriteStatus()
  } else {
    // 退出登录时重置状态
    isLiked.value = false
    isFavorited.value = false
    selectedFolderId.value = null
  }
})

// 组件挂载时检查点赞和收藏状态
onMounted(() => {
  checkArticleLike()
  checkArticleFavoriteStatus()
})

// 监听文章ID变化，重新检查点赞和收藏状态
watch(() => props.articleId, () => {
  checkArticleLike()
  checkArticleFavoriteStatus()
})

</script>

<template>
  <div class="article-actions">
    <div class="action-item" :class="{ active: isLiked }" @click="handleLike">
      <el-icon>
        <Sunny v-if="isLiked" />
        <Moon v-else />
      </el-icon>
      <span>{{ likes }}</span>
    </div>

    <div class="action-item" :class="{ active: isFavorited }" @click="handleFavorite">
      <el-icon>
        <StarFilled v-if="isFavorited" />
        <Star v-else />
      </el-icon>
      <span>{{ favorites }}</span>
    </div>

    <div class="action-item" @click="handleShare">
      <el-icon><Share /></el-icon>
      <span>分享</span>
    </div>

    <div class="action-item" @click="scrollToTop">
      <el-icon><Top /></el-icon>
      <span>顶部</span>
    </div>
  </div>

  <ShareDialog
    v-model:visible="shareDialogVisible"
    :title="articleTitle"
    :url="currentUrl"
  />

  <!-- 文件夹选择对话框 -->
  <el-dialog 
    v-model="folderDialogVisible" 
    title="选择收藏文件夹" 
    width="500px"
    :close-on-click-modal="false"
  >
    <div class="folder-select-container">
      <el-radio-group v-model="selectedFolderId" class="folder-list">
        <div 
          v-for="folder in folders" 
          :key="folder.id" 
          class="folder-option"
          :class="{ 'is-selected': selectedFolderId === folder.id }"
        >
          <el-radio :label="folder.id" class="folder-radio">
            <div class="folder-option-content">
              <div class="folder-option-main">
                <el-icon class="folder-icon">
                  <Folder />
                </el-icon>
                <div class="folder-option-info">
                  <div class="folder-option-name">
                    {{ folder.name }}
                    <el-tag 
                      v-if="folder.isDefault" 
                      size="small" 
                      type="info"
                      style="margin-left: 8px;"
                    >
                      默认
                    </el-tag>
                  </div>
                  <div class="folder-option-desc">{{ folder.description }}</div>
                </div>
              </div>
              <div class="folder-option-count">
                <el-icon><Document /></el-icon>
                <span>{{ folder.articleCount || 0 }} 篇文章</span>
              </div>
            </div>
          </el-radio>
        </div>
      </el-radio-group>
      
      <!-- 空状态提示 -->
      <el-empty 
        v-if="!folders.length" 
        description="暂无收藏夹"
        :image-size="80"
      />
    </div>
    
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="folderDialogVisible = false">取消</el-button>
        <el-button 
          type="primary" 
          :disabled="!selectedFolderId"
          @click="performFavorite(selectedFolderId)"
        >
          <el-icon style="margin-right: 4px;"><Star /></el-icon>
          确定收藏
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<style scoped>
.article-actions {
  position: fixed;
  right: 40px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  gap: 18px;
  z-index: 100;
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 50%;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
  cursor: pointer;
  transition: all 0.4s ease;
  color: #a0aec0;
  border: 1px solid rgba(255, 255, 255, 0.5);
}

html.dark .action-item {
  background: rgba(45, 55, 72, 0.95);
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.action-item:hover {
  transform: translateY(-5px) scale(1.1);
  box-shadow: 0 8px 30px rgba(102, 126, 234, 0.35);
  color: #667eea;
  background: linear-gradient(135deg, #fff 0%, #f0f4ff 100%);
}

html.dark .action-item:hover {
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
}

.action-item.active {
  color: #f56565;
  background: linear-gradient(135deg, #fff5f5 0%, #fed7d7 100%);
  box-shadow: 0 4px 20px rgba(245, 101, 101, 0.3);
}

html.dark .action-item.active {
  background: linear-gradient(135deg, #3a1d1d 0%, #2d1515 100%);
}

.action-item.active:hover {
  transform: translateY(-5px) scale(1.1);
  box-shadow: 0 8px 30px rgba(245, 101, 101, 0.4);
}

.action-item .el-icon {
  font-size: 22px;
  margin-bottom: 3px;
}

.action-item span {
  font-size: 11px;
  font-weight: 500;
}

@media (max-width: 768px) {
  .article-actions {
    right: 20px;
    bottom: 90px;
    top: auto;
    transform: none;
    flex-direction: row;
    gap: 12px;
  }

  .action-item {
    width: 50px;
    height: 50px;
  }
}

/* 文件夹选择对话框样式 */
.folder-select-container {
  padding: 8px 0;
  max-height: 400px;
  overflow-y: auto;
}

.folder-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.folder-option {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  transition: all 0.3s;
  background-color: #fff;
  min-height: 68px;
  width: 100%;
  box-sizing: border-box;
}

.folder-option:hover {
  border-color: #409eff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.folder-option.is-selected {
  border-color: #409eff;
  background-color: #ecf5ff;
}

.folder-option-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.folder-option-main {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.folder-radio {
  width: 100%;
}

.folder-radio :deep(.el-radio__label) {
  width: 100%;
  padding-left: 8px;
}

.folder-icon {
  font-size: 20px;
  color: #409eff;
  flex-shrink: 0;
}

.folder-option-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
}

.folder-option-name {
  font-size: 15px;
  font-weight: 500;
  color: #303133;
  display: flex;
  align-items: center;
}

.folder-option-desc {
  font-size: 12px;
  color: #909399;
  max-width: 280px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-height: 18px;
  height: 18px;
}

.folder-option-count {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #606266;
  background-color: #f5f7fa;
  padding: 4px 10px;
  border-radius: 12px;
  flex-shrink: 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid #e4e7ed;
}

/* 暗黑模式适配 */
html.dark .folder-option {
  background-color: #2d3748;
  border-color: #4a5568;
}

html.dark .folder-option:hover {
  border-color: #409eff;
  background-color: #2c5282;
}

html.dark .folder-option.is-selected {
  border-color: #409eff;
  background-color: #2c5282;
}

html.dark .folder-option-name {
  color: #e2e8f0;
}

html.dark .folder-option-desc {
  color: #a0aec0;
}

html.dark .folder-option-count {
  background-color: #1a202c;
  color: #cbd5e0;
}

html.dark .dialog-footer {
  border-top-color: #4a5568;
}

/* 滚动条样式优化 */
.folder-select-container::-webkit-scrollbar {
  width: 6px;
}

.folder-select-container::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.folder-select-container::-webkit-scrollbar-thumb {
  background: #c0c4cc;
  border-radius: 3px;
}

.folder-select-container::-webkit-scrollbar-thumb:hover {
  background: #909399;
}

html.dark .folder-select-container::-webkit-scrollbar-track {
  background: #1a202c;
}

html.dark .folder-select-container::-webkit-scrollbar-thumb {
  background: #4a5568;
}

html.dark .folder-select-container::-webkit-scrollbar-thumb:hover {
  background: #718096;
}
</style>
