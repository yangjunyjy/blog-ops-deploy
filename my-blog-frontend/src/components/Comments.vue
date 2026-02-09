<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { ChatDotRound, CirclePlus } from '@element-plus/icons-vue'
import { getComments, addComment, deleteComment } from '@/api'

const props = defineProps({
  articleId: {
    type: Number,
    required: true
  }
})

const emit = defineEmits(['update-count'])

const comments = ref([])
const loading = ref(false)
const submitting = ref(false)
const commentText = ref('')
const replyText = ref('')
const replyingTo = ref(null)

const loadComments = async () => {
  loading.value = true
  try {
    const res = await getComments(props.articleId)
    comments.value = res.data || []
    emit('update-count', comments.value.length)
  } catch (error) {
    ElMessage.error('加载评论失败')
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  if (!commentText.value.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }

  submitting.value = true
  try {
    await addComment({
      articleId: props.articleId,
      content: commentText.value
    })
    ElMessage.success('评论成功')
    commentText.value = ''
    await loadComments()
  } catch (error) {
    ElMessage.error('评论失败')
  } finally {
    submitting.value = false
  }
}

const handleReply = (comment) => {
  replyingTo.value = comment
  replyText.value = `@${comment.author.name} `
}

const handleCancelReply = () => {
  replyingTo.value = null
  replyText.value = ''
}

const handleSubmitReply = async (comment) => {
  if (!replyText.value.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }

  submitting.value = true
  try {
    await addComment({
      articleId: props.articleId,
      content: replyText.value,
      parentId: comment.id
    })
    ElMessage.success('回复成功')
    handleCancelReply()
    await loadComments()
  } catch (error) {
    ElMessage.error('回复失败')
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (commentId) => {
  try {
    await deleteComment(commentId)
    ElMessage.success('删除成功')
    await loadComments()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

onMounted(() => {
  loadComments()
})
</script>

<template>
  <div class="comments-section">
    <div class="comments-header">
      <el-icon><ChatDotRound /></el-icon>
      <h3>评论 ({{ comments.length }})</h3>
    </div>

    <!-- 发表评论 -->
    <div class="comment-form">
      <el-input
        v-model="commentText"
        type="textarea"
        :rows="3"
        placeholder="发表你的评论..."
        maxlength="500"
        show-word-limit
      />
      <div class="form-actions">
        <el-button
          type="primary"
          :loading="submitting"
          @click="handleSubmit"
        >
          发表评论
        </el-button>
      </div>
    </div>

    <!-- 评论列表 -->
    <div class="comments-list" v-loading="loading">
      <div v-if="comments.length === 0" class="no-comments">
        <el-empty description="暂无评论，快来抢沙发吧！" />
      </div>

      <div
        v-for="comment in comments"
        :key="comment.id"
        class="comment-item"
      >
        <div class="comment-avatar">
          <img :src="comment.author.avatar" :alt="comment.author.name" />
        </div>

        <div class="comment-content-wrapper">
          <div class="comment-info">
            <span class="author-name">{{ comment.author.name }}</span>
            <span class="comment-time">{{ $filters.formatRelativeTime(comment.createdAt) }}</span>
          </div>

          <div class="comment-text">{{ comment.content }}</div>

          <div class="comment-actions">
            <el-button
              type="primary"
              text
              size="small"
              :icon="CirclePlus"
              @click="handleReply(comment)"
            >
              回复
            </el-button>
            <el-button
              type="danger"
              text
              size="small"
              @click="handleDelete(comment.id)"
            >
              删除
            </el-button>
          </div>

          <!-- 回复�?-->
          <div v-if="replyingTo && replyingTo.id === comment.id" class="reply-form">
            <el-input
              v-model="replyText"
              type="textarea"
              :rows="2"
              :placeholder="`回复 ${comment.author.name}`"
              maxlength="500"
              show-word-limit
            />
            <div class="reply-actions">
              <el-button size="small" @click="handleCancelReply">取消</el-button>
              <el-button
                type="primary"
                size="small"
                :loading="submitting"
                @click="handleSubmitReply(comment)"
              >
                回复
              </el-button>
            </div>
          </div>

          <!-- 子评�?-->
          <div v-if="comment.replies && comment.replies.length > 0" class="replies-list">
            <div
              v-for="reply in comment.replies"
              :key="reply.id"
              class="reply-item"
            >
              <div class="reply-avatar">
                <img :src="reply.author.avatar" :alt="reply.author.name" />
              </div>

              <div class="reply-content-wrapper">
                <div class="reply-info">
                  <span class="author-name">{{ reply.author.name }}</span>
                  <span class="comment-time">{{ $filters.formatRelativeTime(reply.createdAt) }}</span>
                </div>

                <div class="reply-text">{{ reply.content }}</div>

                <div class="reply-actions">
                  <el-button
                    type="primary"
                    text
                    size="small"
                    @click="handleDelete(reply.id)"
                  >
                    删除
                  </el-button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.comments-section {
  margin-top: 50px;
  padding-top: 50px;
  border-top: 2px solid #e2e8f0;
  position: relative;
}

.comments-section::before {
  content: '';
  position: absolute;
  top: -2px;
  left: 0;
  width: 50px;
  height: 2px;
  background: #3b82f6;
}

.comments-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 35px;
}

.comments-header h3 {
  font-size: 26px;
  color: #1a202c;
  margin: 0;
  font-weight: 700;
  letter-spacing: -0.5px;
}

.comments-header .el-icon {
  font-size: 28px;
  color: #667eea;
}

.comment-form {
  margin-bottom: 35px;
  padding: 24px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.5);
}

.comment-form :deep(.el-textarea__inner) {
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  transition: all 0.3s;
}

.comment-form :deep(.el-textarea__inner:focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-actions {
  margin-top: 18px;
  text-align: right;
}

.form-actions :deep(.el-button--primary) {
  background: #3b82f6;
  border: none;
  border-radius: 10px;
  padding: 12px 28px;
  font-weight: 600;
  transition: all 0.3s;
}

.form-actions :deep(.el-button--primary:hover) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.no-comments {
  text-align: center;
  padding: 60px 20px;
}

.comment-item {
  display: flex;
  gap: 18px;
  padding: 24px;
  background: #f8fafc;
  border-radius: 16px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.5);
  transition: all 0.3s;
}

.comment-item:hover {
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.comment-avatar {
  flex-shrink: 0;
  width: 52px;
  height: 52px;
}

.comment-avatar img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #fff;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.comment-content-wrapper {
  flex: 1;
  min-width: 0;
}

.comment-info {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 12px;
}

.author-name {
  font-weight: 700;
  color: #1a202c;
  font-size: 15px;
}

.comment-time {
  font-size: 13px;
  color: #a0aec0;
  font-weight: 500;
}

.comment-text {
  color: #4a5568;
  line-height: 1.7;
  margin-bottom: 12px;
  font-size: 15px;
}

.comment-actions {
  display: flex;
  gap: 12px;
}

.comment-actions :deep(.el-button) {
  padding: 6px 14px;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s;
}

.reply-form {
  margin-top: 18px;
  padding: 20px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}

.reply-actions {
  margin-top: 12px;
  display: flex;
  gap: 12px;
}


.replies-list {
  margin-top: 20px;
  padding-left: 20px;
  border-left: 3px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.reply-item {
  display: flex;
  gap: 14px;
  padding: 16px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 12px;
  transition: all 0.3s;
}

.reply-item:hover {
  background: linear-gradient(135deg, #f0f4ff 0%, #ede9fe 100%);
}

.reply-avatar {
  flex-shrink: 0;
  width: 44px;
  height: 44px;
}

.reply-avatar img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.reply-content-wrapper {
  flex: 1;
  min-width: 0;
}

.reply-info {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
}

.reply-text {
  color: #4a5568;
  line-height: 1.6;
  margin-bottom: 10px;
  font-size: 14px;
}

@media (max-width: 768px) {
  .comments-section {
    padding-top: 35px;
  }

  .comment-item {
    padding: 18px;
    gap: 14px;
  }

  .comment-avatar {
    width: 44px;
    height: 44px;
  }

  .replies-list {
    padding-left: 12px;
    gap: 14px;
  }

  .reply-item {
    padding: 14px;
  }

  .reply-avatar {
    width: 36px;
    height: 36px;
  }
}
</style>
