<script setup>
import { useRouter } from 'vue-router'
import { Plus, House } from '@element-plus/icons-vue'

const props = defineProps({
  message: {
    type: String,
    default: '暂无数据'
  },
  showAction: {
    type: Boolean,
    default: true
  }
})

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<template>
  <div class="empty-state">
    <div class="empty-content">
      <div class="empty-icon">
        <el-icon :size="80"><Box /></el-icon>
      </div>
      <h2 class="empty-title">{{ message }}</h2>
      <div class="empty-actions" v-if="showAction">
        <el-button type="primary" :icon="House" @click="goHome">
          返回首页
        </el-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.empty-state {
  min-height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.empty-content {
  text-align: center;
}

.empty-icon {
  color: #c0c4cc;
  margin-bottom: 20px;
}

.empty-title {
  font-size: 18px;
  color: #909399;
  margin-bottom: 30px;
}

.empty-actions {
  display: flex;
  gap: 15px;
  justify-content: center;
}
</style>
