<script setup>
import { useRouter } from 'vue-router'
import { RefreshLeft, House } from '@element-plus/icons-vue'

const props = defineProps({
  message: {
    type: String,
    default: '出错�?
  }
})

const emit = defineEmits(['retry'])

const router = useRouter()

const handleRetry = () => {
  emit('retry')
}

const goHome = () => {
  router.push('/')
}
</script>

<template>
  <div class="error-state">
    <div class="error-content">
      <div class="error-icon">
        <el-icon :size="80"><Warning /></el-icon>
      </div>
      <h2 class="error-title">{{ message }}</h2>
      <p class="error-desc">抱歉，加载失败，请稍后再�?/p>
      <div class="error-actions">
        <el-button type="primary" :icon="RefreshLeft" @click="handleRetry">
          重试
        </el-button>
        <el-button :icon="House" @click="goHome">
          返回首页
        </el-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.error-state {
  min-height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.error-content {
  text-align: center;
  padding: 40px;
}

.error-icon {
  color: #f56c6c;
  margin-bottom: 20px;
}

.error-title {
  font-size: 24px;
  color: #303133;
  margin-bottom: 12px;
}

.error-desc {
  color: #909399;
  margin-bottom: 30px;
}

.error-actions {
  display: flex;
  gap: 15px;
  justify-content: center;
}
</style>
