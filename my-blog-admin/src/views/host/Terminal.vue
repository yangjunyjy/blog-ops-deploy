<template>
  <div class="terminal-page">
    <!-- 顶部工具栏 -->
    <div class="page-header">
      <div class="header-left">
        <el-icon class="header-icon" color="#007acc">
          <Monitor />
        </el-icon>
        <h2 class="header-title">远程终端</h2>
        <span class="host-info" v-if="currentHost">
          {{ currentHost.name }}@{{ currentHost.ip }}:{{ currentHost.port }}
        </span>
      </div>
      <div class="header-right">
        <el-button size="small" :icon="Close" circle @click="handleClose" />
      </div>
    </div>

    <!-- 多标签会话管理 -->
    <terminal-session :initial-host="currentHost" @close="handleSessionClose" />
  </div>
</template>

<script setup>
defineOptions({
  name: 'Terminal'
})

import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Monitor, Close } from '@element-plus/icons-vue'
import TerminalSession from '@/components/host/TerminalSession.vue'
import { getHostDetail } from '@/api/host'

const router = useRouter()
const route = useRoute()
const currentHost = ref(null)

// 加载主机信息
const loadHostInfo = async () => {
  const hostId = route.params.hostId
  console.log('Terminal.vue 获取到的主机 ID:', hostId, '所有参数:', route.params)

  if (!hostId) {
    ElMessage.error('缺少主机ID')
    router.back()
    return
  }

  try {
    // 从后端获取主机信息
    console.log('开始加载主机信息...')
    const data = await getHostDetail(hostId)
    console.log('主机信息响应:', data)
    const hostdata = data.data

    if (!hostdata) {
      ElMessage.error('未找到主机信息')
      router.back()
      return
    }

    currentHost.value = {
      id: hostdata.id,
      name: hostdata.name,
      ip: hostdata.address,
      port: hostdata.port,
      auth_type: hostdata.type,
      username: hostdata.username,
      password: hostdata.password,
      connect_type: 'web'
    }

    console.log('设置当前主机:', currentHost.value)
  } catch (error) {
    console.error('加载主机信息失败:', error)
    ElMessage.error('加载主机信息失败: ' + (error.message || '未知错误'))
  }
}

// 会话关闭处理
const handleSessionClose = () => {
  ElMessageBox.confirm('确定要关闭所有终端会话吗？', '确认关闭', { type: 'warning' })
    .then(() => {
      handleClose()
    })
    .catch(() => { })
}

// 关闭页面
const handleClose = () => {
  window.close()
  // 如果 window.close() 无法关闭（在新标签页中），则返回上一页
  setTimeout(() => {
    router.push('/host')
  }, 300)
}

// 监听页面关闭
onMounted(() => {
  loadHostInfo()

  // 监听键盘快捷键
  window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})

// 键盘快捷键
const handleKeyDown = (e) => {
  // Ctrl+T: 新建会话
  if (e.ctrlKey && e.key === 't') {
    e.preventDefault()
    ElMessage.info('请从主机列表选择要连接的主机')
  }
  // Ctrl+W: 关闭当前标签
  if (e.ctrlKey && e.key === 'w') {
    e.preventDefault()
    ElMessage.info('请使用标签上的关闭按钮')
  }
  // Ctrl+Tab: 切换标签
  if (e.ctrlKey && e.key === 'Tab') {
    e.preventDefault()
    ElMessage.info('请使用鼠标点击标签切换')
  }
}
</script>

<style scoped lang="scss">
.terminal-page {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #1e1e1e;
  overflow: hidden;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background-color: #1e1e1e;
  border-bottom: 1px solid #3e3e3e;
  flex-shrink: 0;

  .header-left {
    display: flex;
    align-items: center;
    gap: 12px;

    .header-icon {
      font-size: 24px;
      color: #cccccc;
    }

    .header-title {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      color: #cccccc;
    }

    .host-info {
      font-size: 13px;
      color: #4ec9b0;
      padding: 4px 12px;
      background-color: rgba(78, 201, 176, 0.15);
      border-radius: 4px;
      border: 1px solid rgba(78, 201, 176, 0.4);
    }
  }

  .header-right {
    display: flex;
    gap: 8px;
  }
}

:deep(.el-button) {
  background-color: #007acc;
  border-color: #007acc;
  color: #ffffff;
  height: 32px;
  padding: 0 16px;
  font-size: 13px;

  &:hover {
    background-color: #0062a3;
    border-color: #0062a3;
  }

  &.el-button--default {
    background-color: #3c3c3c;
    border-color: #3c3c3c;

    &:hover {
      background-color: #4a4a4a;
      border-color: #4a4a4a;
    }
  }
}
</style>
