<template>
  <!-- 回复评论对话框 -->
  <el-dialog
    :model-value="reployDialogVisible"
    @update:model-value="(val) => emit('update:reployDialogVisible', val)"
    title="回复评论"
    width="500px"
    :modal="true"
    :close-on-click-modal="true"
    :close-on-press-escape="true"
    :show-close="true"
    class="reply-comment-dialog"
    transition="el-fade-in-linear"
  >
    <!-- 评论引用区 -->
    <div class="comment-quote">
      <p class="quote-label">回复：</p>
      <p class="quote-content">{{ replyTargetContent || "暂无评论内容" }}</p>
    </div>

    <!-- 回复输入区 -->
    <el-form :model="replyForm" label-width="0" class="reply-form">
      <el-form-item>
        <el-input
          v-model="replyForm.content"
          type="textarea"
          :rows="5"
          placeholder="请输入回复内容..."
          resize="none"
          class="reply-input"
        />
      </el-form-item>
    </el-form>

    <!-- 底部操作区 -->
    <template #footer>
      <div class="dialog-footer">
        <el-button
          @click="emit('update:reployDialogVisible', false)"
          class="cancel-btn"
        >
          取消
        </el-button>
        <el-button
          type="primary"
          @click="handleSubmit"
          class="submit-btn"
          :loading="submitLoading"
        >
          提交回复
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, defineProps, defineEmits } from 'vue'
import { ElMessage } from 'element-plus'
import { addComment } from '@/api'

// 接收父组件参数
const props = defineProps({
  // 控制对话框显隐
  reployDialogVisible: {
    type: Boolean,
    default: false
  },
  // 回复的目标评论内容
  replyTargetContent: {
    type: String,
    default: ''
  },
  // 回复的目标评论ID
  replyTargetId: {
    type: [String, Number],
    default: ''
  },
  // 文章ID
  articleId: {
    type: Number,
    required: true
  }
})

// 定义事件
const emit = defineEmits(['close', 'submit', 'update:reployDialogVisible'])

// 响应式数据
const replyForm = ref({
  content: '' // 回复内容
})
const submitLoading = ref(false) // 提交加载状态

// 提交回复
const handleSubmit = async () => {
  if (!replyForm.value.content.trim()) {
    ElMessage.warning('请输入回复内容！')
    return
  }
  submitLoading.value = true

  try {
    const data = await addComment({
      articleId: props.articleId,
      content: replyForm.value.content,
      parentId: props.replyTargetId
    })

    // 重置表单并关闭对话框
    replyForm.value.content = ''
    submitLoading.value = false
    emit('update:reployDialogVisible', false)

    // 根据评论状态显示不同提示
    if (data.data?.status === 0) {
      ElMessage.info('评论已提交，等待审核')
    } else {
      ElMessage.success('回复成功！')
    }

    // 向父组件传递回复数据（在提示之后）
    emit('submit', data.data)
  } catch (error) {
    console.error('提交回复失败:', error)
    ElMessage.error(error.response?.data?.message || '提交回复失败，请重试')
    submitLoading.value = false
  }
}
</script>

<style scoped>
/* 整体对话框样式 - 淡色系 */
.reply-comment-dialog {
  --el-dialog-bg-color: #fafbfc;
  --el-dialog-title-color: #333;
}

/* 对话框标题样式 */
:deep(.el-dialog__header) {
  border-bottom: 1px solid #eef1f5;
  padding-bottom: 10px;
}

:deep(.el-dialog__title) {
  font-size: 16px;
  font-weight: 500;
  color: #404958;
}

/* 评论引用区 */
.comment-quote {
  padding: 12px 15px;
  background-color: #f5f7fa;
  border-radius: 6px;
  margin-bottom: 15px;
}

.quote-label {
  margin: 0 0 5px 0;
  font-size: 14px;
  color: #6e7681;
  font-weight: 500;
}

.quote-content {
  margin: 0;
  font-size: 14px;
  color: #404958;
  line-height: 1.5;
  word-break: break-all;
}

/* 回复表单样式 */
.reply-form {
  margin-top: 10px;
}

.reply-input {
  --el-input-bg-color: #ffffff;
  --el-input-border-color: #d8e0e8;
  --el-input-hover-border-color: #b8c6d4;
  border-radius: 6px;
}

:deep(.el-textarea__inner) {
  padding: 12px;
  font-size: 14px;
  color: #333;
}

/* 底部按钮区 */
.dialog-footer {
  text-align: right;
  padding-top: 10px;
  border-top: 1px solid #eef1f5;
}

.cancel-btn {
  background-color: #f0f2f5;
  color: #6e7681;
  border: none;
  margin-right: 10px;
}

.cancel-btn:hover {
  background-color: #e5e8ed;
  color: #404958;
}

.submit-btn {
  background-color: #4e89e8;
  border: none;
  --el-button-hover-bg-color: #3a79d8;
}

/* 动画优化 - 让弹出更优雅 */
:deep(.el-fade-in-linear-enter-from) {
  opacity: 0;
  transform: translateY(-20px) scale(0.98);
}

:deep(.el-fade-in-linear-enter-active) {
  transition: all 0.3s cubic-bezier(0.38, 0.99, 0.42, 1);
}

:deep(.el-fade-in-linear-leave-to) {
  opacity: 0;
  transform: translateY(10px);
}

:deep(.el-fade-in-linear-leave-active) {
  transition: all 0.2s cubic-bezier(0.38, 0.99, 0.42, 1);
}
</style>