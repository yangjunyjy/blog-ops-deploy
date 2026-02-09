<template>
  <div class="comment-item" :class="{ 'is-reply': isReply }">
    <!-- 评论头像 -->
    <div class="comment-avatar">
      <img :src="comment.author?.avatar || '/images/default.jpg'" :alt="comment.author?.name" />
    </div>

    <!-- 评论内容 -->
    <div class="comment-content">
      <!-- 评论头部 -->
      <div class="comment-header">
        <div class="header-left">
          <span class="comment-author">{{ comment.author?.name || '匿名用户' }}</span>
          <span class="comment-date">{{ formatDate(comment.createdAt) }}</span>
        </div>
        <div class="header-right">
          <el-button
            text
            size="small"
            :type="comment.isLiked ? 'primary' : 'default'"
            @click="handleLike"
          >
            <el-icon v-if="comment.isLiked"><StarFilled /></el-icon>
            <el-icon v-else><Star /></el-icon>
            {{ comment.likes || 0 }}
          </el-button>
          <el-button text size="small" @click="handleReply">回复</el-button>
          <el-button
            v-if="canDelete"
            text
            size="small"
            type="danger"
            @click="handleDelete"
          >
            删除
          </el-button>
        </div>
      </div>

      <!-- 评论内容 -->
      <div class="comment-text">{{ comment.content }}</div>

      <!-- 子评论区域 -->
      <div v-if="hasReplies" class="replies-section">
        <!-- 折叠/展开按钮 -->
        <div v-if="showCollapse" class="collapse-trigger" @click="toggleReplies">
          <span class="collapse-text">
            {{ isExpanded ? '收起' : `展开全部 ${replyCount} 条回复` }}
          </span>
        </div>

        <!-- 子评论列表 -->
        <div v-if="hasReplies" class="replies-list">
          <div v-loading="repliesLoading">
            <div v-for="reply in displayReplies" :key="reply.id" class="reply-item">
              <span class="reply-user-name">{{ reply.author?.name || '匿名用户' }}：</span>
              <span v-if="reply.replyTarget" class="reply-target">回复 @{{ reply.replyTarget.author?.name }}</span>
              <span class="reply-content">{{ reply.content }}</span>
              <span class="reply-meta">
                <span class="reply-date">{{ formatDate(reply.createdAt) }}</span>
                <span class="reply-actions">
                  <el-button
                    text
                    size="small"
                    :type="reply.isLiked ? 'primary' : 'default'"
                    @click="handleReplyLike(reply)"
                  >
                    <el-icon v-if="reply.isLiked"><StarFilled /></el-icon>
                    <el-icon v-else><Star /></el-icon>
                    {{ reply.likes || 0 }}
                  </el-button>
                  <el-button text size="small" @click="handleReplyClick(reply)">回复</el-button>
                  <el-button
                    v-if="canDeleteReply(reply)"
                    text
                    size="small"
                    type="danger"
                    @click="handleReplyDelete(reply)"
                  >
                    删除
                  </el-button>
                </span>
              </span>
            </div>

            <!-- 加载更多按钮 -->
            <div v-if="hasMoreReplies && isExpanded" class="load-more">
              <el-button text size="small" @click="loadMoreReplies" :loading="loadingMore">
                加载更多
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Star, StarFilled, ArrowUp, ArrowDown } from '@element-plus/icons-vue'
import { likeComment, unlikeComment, getCommentReplies, getReplyCount, deleteComment } from '@/api'
import { formatDate } from '@/utils/format'

const props = defineProps({
  comment: {
    type: Object,
    required: true
  },
  replyTarget: {
    type: Object,
    default: null
  },
  articleId: {
    type: Number,
    required: true
  },
  isReply: {
    type: Boolean,
    default: false
  },
  currentUserId: {
    type: Number,
    default: null
  },
  maxDepth: {
    type: Number,
    default: 3
  },
  currentDepth: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['like', 'reply', 'delete', 'reply-added'])

// 响应式数据
const isExpanded = ref(false)
const replies = ref([])
const repliesLoading = ref(false)
const loadingMore = ref(false)
const replyPage = ref(1)
const replyCount = ref(0)

// 抖音模式：默认显示2条回复预览
const previewCount = 2

// 计算属性
const showCollapse = computed(() => {
  return hasReplies.value && replyCount.value > previewCount
})

const hasReplies = computed(() => replyCount.value > 0)

const hasMoreReplies = computed(() => {
  return replies.value.length < replyCount.value
})

// 显示的回复列表：展开时显示全部，折叠时只显示预览数
const displayReplies = computed(() => {
  if (isExpanded.value || !showCollapse.value) {
    return replies.value
  }
  return replies.value.slice(0, previewCount)
})

const canDelete = computed(() => {
  return props.currentUserId && props.comment.userId === props.currentUserId
})

const canDeleteReply = (reply) => {
  return props.currentUserId && reply.userId === props.currentUserId
}

const canLoadMoreReplies = computed(() => {
  // 抖音模式：顶级评论才能加载子回复，子评论本身不递归显示
  return props.currentDepth === 0
})

// 加载子评论数量
const loadReplyCount = async () => {
  try {
    const res = await getReplyCount(props.comment.id)
    replyCount.value = res.data?.count || 0
  } catch (error) {
    console.error('获取子评论数量失败:', error)
  }
}

// 加载子评论列表
const loadReplies = async (isPreview = false) => {
  if (!canLoadMoreReplies.value || repliesLoading.value) {
    console.log(`评论 ${props.comment.id} 加载被阻止: canLoadMoreReplies=${canLoadMoreReplies.value}, repliesLoading=${repliesLoading.value}`)
    return
  }

  console.log(`评论 ${props.comment.id} 开始加载回复, isPreview=${isPreview}`)
  repliesLoading.value = true
  try {
    // 如果是预览模式，只加载previewCount条
    const pageSize = isPreview ? Math.min(previewCount, 10) : 10
    const res = await getCommentReplies(props.comment.id, replyPage.value, pageSize)
    console.log(`评论 ${props.comment.id} 加载结果:`, res.data)
    if (res.code === 200) {
      // 后端已经返回了 reply_to 和 reply_to_user 信息
      replies.value = (res.data?.items || []).map(item => ({
        ...item,
        userId: item.user_id,
        replyTarget: item.reply_to_user ? {
          author: {
            name: item.reply_to_user.name,
            id: item.reply_to_user.id,
            avatar: item.reply_to_user.avatar
          }
        } : null
      }))
      replyCount.value = res.data?.total || 0
      console.log(`评论 ${props.comment.id} 设置完成, replies=${replies.value.length}, replyCount=${replyCount.value}, showCollapse=${showCollapse.value}`)
    }
  } catch (error) {
    console.error('加载子评论失败:', error)
  } finally {
    repliesLoading.value = false
  }
}

// 加载更多回复
const loadMoreReplies = async () => {
  if (loadingMore.value || !hasMoreReplies.value) {
    return
  }

  loadingMore.value = true
  try {
    const res = await getCommentReplies(props.comment.id, replyPage.value + 1, 10)
    if (res.code === 200) {
      const newItems = (res.data?.items || []).map(item => ({
        ...item,
        userId: item.user_id,
        replyTarget: item.reply_to_user ? {
          author: {
            name: item.reply_to_user.name,
            id: item.reply_to_user.id,
            avatar: item.reply_to_user.avatar
          }
        } : null
      }))
      replies.value = [...replies.value, ...newItems]
      replyPage.value++
    }
  } catch (error) {
    console.error('加载更多回复失败:', error)
  } finally {
    loadingMore.value = false
  }
}

// 切换子评论显示状态
const toggleReplies = () => {
  isExpanded.value = !isExpanded.value
  if (isExpanded.value && replies.value.length === 0) {
    loadReplies(true)
  } else if (isExpanded.value && replies.value.length < replyCount.value) {
    // 如果是展开状态但数据不足，加载更多
    loadReplies(false)
  }
}

// 处理回复点击
const handleReplyClick = (reply) => {
  emit('reply', reply)
}

// 点赞评论
const handleLike = async () => {
  try {
    const comment = props.comment
    if (comment.isLiked) {
      await unlikeComment(comment.id)
      comment.likes--
      comment.isLiked = false
      ElMessage.success('取消点赞成功')
    } else {
      await likeComment(comment.id)
      comment.likes++
      comment.isLiked = true
      ElMessage.success('点赞成功')
    }
  } catch (error) {
    console.error('点赞操作失败:', error)
    ElMessage.error(error.response?.data?.message || '操作失败，请重试')
  }
}

// 处理子评论点赞
const handleReplyLike = async (reply) => {
  try {
    if (reply.isLiked) {
      await unlikeComment(reply.id)
      reply.likes--
      reply.isLiked = false
      ElMessage.success('取消点赞成功')
    } else {
      await likeComment(reply.id)
      reply.likes++
      reply.isLiked = true
      ElMessage.success('点赞成功')
    }
  } catch (error) {
    console.error('点赞操作失败:', error)
    ElMessage.error(error.response?.data?.message || '操作失败，请重试')
  }
}

// 回复评论
const handleReply = () => {
  emit('reply', props.comment)
}

// 处理回复子评论
const handleReplyReply = (reply) => {
  emit('reply', reply)
}

// 删除评论
const handleDelete = async () => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？删除后所有回复也会被删除。', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    try {
      await deleteComment(props.comment.id)
      ElMessage.success('删除成功')
      emit('delete', props.comment)
    } catch (error) {
      console.error('删除评论失败:', error)
      ElMessage.error(error.response?.data?.message || '删除失败，请重试')
    }
  } catch (error) {
    // 用户取消删除
  }
}

// 处理子评论删除
const handleReplyDelete = async (reply) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？删除后所有回复也会被删除。', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    try {
      await deleteComment(reply.id)
      ElMessage.success('删除成功')
      // 从本地列表中移除
      const index = replies.value.findIndex(r => r.id === reply.id)
      if (index !== -1) {
        replies.value.splice(index, 1)
        replyCount.value--
      }
    } catch (error) {
      console.error('删除评论失败:', error)
      ElMessage.error(error.response?.data?.message || '删除失败，请重试')
    }
  } catch (error) {
    // 用户取消删除
  }
}

// 添加新回复到当前评论
const addNewReply = (replyData) => {
  // replyData 是后端返回的完整数据：{ comment, status, message, user, reply_to, reply_to_user }
  const comment = replyData.comment

  const newReply = {
    id: comment.id,
    userId: comment.user_id,
    content: comment.content,
    author: replyData.user ? {
      id: replyData.user.id,
      name: replyData.user.nickname || replyData.user.username || '用户',
      avatar: replyData.user.avatar || '/images/default.jpg'
    } : {
      id: comment.user_id,
      name: comment.user?.nickname || comment.user?.username || '用户',
      avatar: comment.user?.avatar || '/images/default.jpg'
    },
    replyTarget: replyData.reply_to_user ? {
      author: {
        name: replyData.reply_to_user.nickname || replyData.reply_to_user.username,
        id: replyData.reply_to_user.id,
        avatar: replyData.reply_to_user.avatar
      }
    } : null,
    createdAt: comment.created_at || new Date().toISOString(),
    likes: comment.likes || 0,
    isLiked: false
  }

  replies.value.unshift(newReply)
  replyCount.value++
  isExpanded.value = true
}

// 初始化
onMounted(async () => {
  console.log(`评论 ${props.comment.id} 初始化, currentDepth=${props.currentDepth}`)
  await loadReplyCount()
  console.log(`评论 ${props.comment.id} 子评论数量: ${replyCount.value}`)
  // 如果有子评论，自动加载第一页预览（只加载2条）
  if (replyCount.value > 0) {
    console.log(`评论 ${props.comment.id} 开始加载预览`)
    await loadReplies(true)
  }
})

// 暴露方法给父组件
defineExpose({
  addNewReply
})
</script>

<style scoped>
.comment-item {
  display: flex;
  gap: 16px;
  padding: 20px;
  background: #fff;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  transition: all 0.3s;
  margin-bottom: 16px;
}

.comment-item:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.comment-item.is-reply {
  background: #f9fafb;
  margin-top: 12px;
  padding: 16px;
}

/* 暗黑模式 */
html.dark .comment-item {
  background: #2d3748;
  border-color: #4a5568;
}

html.dark .comment-item:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

html.dark .comment-item.is-reply {
  background: #1a202c;
}

.comment-avatar {
  flex-shrink: 0;
}

.comment-avatar img {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #667eea;
}

html.dark .comment-avatar img {
  border-color: #764ba2;
}

.comment-content {
  flex: 1;
  min-width: 0;
}

.comment-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  gap: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.comment-author {
  font-weight: 600;
  font-size: 15px;
  color: #1a202c;
}

html.dark .comment-author {
  color: #f7fafc;
}

.comment-date {
  font-size: 13px;
  color: #a0aec0;
}

.reply-to {
  margin-bottom: 8px;
  padding: 8px 12px;
  background: linear-gradient(135deg, #f0f4ff 0%, #ede9fe 100%);
  border-radius: 8px;
  display: inline-block;
}

html.dark .reply-to {
  background: linear-gradient(135deg, #2c5282 0%, #1a365d 100%);
}

.reply-label {
  font-size: 12px;
  color: #667eea;
  font-weight: 500;
}

html.dark .reply-label {
  color: #764ba2;
}

.reply-user {
  font-size: 13px;
  color: #4a5568;
  margin-left: 6px;
}

html.dark .reply-user {
  color: #cbd5e0;
}

.reply-item {
  padding: 12px 0;
  border-bottom: 1px solid #e2e8f0;
  font-size: 14px;
  line-height: 1.8;
}

html.dark .reply-item {
  border-bottom-color: #4a5568;
}

.reply-item:last-child {
  border-bottom: none;
}

.reply-user-name {
  font-weight: 600;
  color: #1a202c;
}

html.dark .reply-user-name {
  color: #f7fafc;
}

.reply-target {
  color: #667eea;
  margin: 0 4px;
}

.reply-content {
  color: #4a5568;
}

html.dark .reply-content {
  color: #cbd5e0;
}

.reply-meta {
  display: inline-flex;
  align-items: center;
  margin-top: 6px;
  font-size: 12px;
  width: 100%;
}

.reply-date {
  color: #a0aec0;
  margin-right: 8px;
}

.reply-actions {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.reply-actions .el-button {
  color: #a0aec0;
  padding: 2px 8px;
  font-size: 12px;
}

.reply-actions .el-button:hover {
  color: #667eea;
}

.comment-text {
  color: #4a5568;
  line-height: 1.6;
  margin-bottom: 12px;
  word-break: break-word;
}

html.dark .comment-text {
  color: #cbd5e0;
}

.header-right .el-button {
  color: #718096;
  padding: 4px 8px;
}

html.dark .header-right .el-button {
  color: #a0aec0;
}

.header-right .el-button:hover {
  color: #667eea;
}

.replies-section {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px dashed #e2e8f0;
}

html.dark .replies-section {
  border-top-color: #4a5568;
}

.collapse-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 8px 12px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 6px;
  transition: all 0.3s;
}

.collapse-trigger:hover {
  background: linear-gradient(135deg, #e2e8f0 0%, #e0e7ff 100%);
  transform: translateX(4px);
}

html.dark .collapse-trigger {
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
}

html.dark .collapse-trigger:hover {
  background: linear-gradient(135deg, #2d3748 0%, #374151 100%);
}

.collapse-text {
  font-size: 13px;
  color: #4a5568;
}

html.dark .collapse-text {
  color: #cbd5e0;
}

.collapse-icon {
  color: #667eea;
  transition: transform 0.3s;
}

.collapse-trigger:hover .collapse-icon {
  transform: scale(1.1);
}

.replies-list {
  margin-top: 12px;
}

.load-more {
  text-align: center;
  margin-top: 12px;
}

.load-more .el-button {
  color: #667eea;
}

@media (max-width: 768px) {
  .comment-item {
    padding: 16px;
    gap: 12px;
  }

  .comment-avatar img {
    width: 40px;
    height: 40px;
  }

  .comment-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .header-left {
    width: 100%;
  }

  .header-right {
    width: 100%;
    justify-content: flex-start;
  }
}
</style>
