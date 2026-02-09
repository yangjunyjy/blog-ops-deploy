<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Share, Link, DocumentCopy, Close, ChatDotRound, Share as ShareIcon, Promotion, User, Files, Connection } from '@element-plus/icons-vue'

const iconMap = {
  ChatDotRound,
  ShareIcon,
  Promotion,
  User,
  Files,
  DocumentCopy,
  Connection
}

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: ''
  },
  url: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:visible', 'close'])

const sharePlatforms = [
  {
    name: '微信',
    color: '#07c160',
    action: 'weixin',
    icon: 'ChatDotRound'
  },
  {
    name: '朋友圈',
    color: '#07c160',
    action: 'friend',
    icon: 'Connection'
  },
  {
    name: '微博',
    color: '#e6162d',
    action: 'weibo',
    icon: 'Promotion'
  },
  {
    name: 'QQ',
    color: '#12b7f5',
    action: 'qq',
    icon: 'User'
  },
  {
    name: 'QQ空间',
    color: '#ffce00',
    action: 'qzone',
    icon: 'Files'
  },
  {
    name: '复制链接',
    color: '#1890ff',
    action: 'copy',
    icon: 'DocumentCopy'
  }
]

const currentUrl = ref('')

const handleShare = async (platform) => {
  currentUrl.value = props.url || window.location.href
  const encodedUrl = encodeURIComponent(currentUrl.value)
  const encodedTitle = encodeURIComponent(props.title || document.title)

  switch (platform.action) {
    case 'weixin':
      showQRCode(currentUrl.value)
      break
    case 'friend':
      showQRCode(currentUrl.value)
      break
    case 'weibo':
      window.open(`https://service.weibo.com/share/share.php?url=${encodedUrl}&title=${encodedTitle}`, '_blank')
      break
    case 'qq':
      window.open(`https://connect.qq.com/widget/shareqq/index.html?url=${encodedUrl}&title=${encodedTitle}`, '_blank')
      break
    case 'qzone':
      window.open(`https://sns.qzone.qq.com/cgi-bin/qzshare/cgi_qzshare_onekey?url=${encodedUrl}&title=${encodedTitle}`, '_blank')
      break
    case 'copy':
      await copyToClipboard(currentUrl.value)
      break
  }
}

const showQRCode = (url) => {
  ElMessage.info('请使用微信扫一扫分享')
  window.open(`https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(url)}`, '_blank')
}

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('链接已复制到剪切板')
  } catch (error) {
    const input = document.createElement('input')
    input.value = text
    document.body.appendChild(input)
    input.select()
    document.execCommand('copy')
    document.body.removeChild(input)
    ElMessage.success('链接已复制到剪切板')
  }
}

const handleClose = () => {
  emit('update:visible', false)
  emit('close')
}

const handleCopyUrl = async () => {
  currentUrl.value = props.url || window.location.href
  await copyToClipboard(currentUrl.value)
}
</script>

<template>
  <el-dialog
    :model-value="visible"
    title="分享文章"
    :width="480"
    @close="handleClose"
    class="share-dialog"
  >
    <div class="share-content">
      <div class="share-header">
        <h4>分享</h4>
      </div>

      <div class="share-platforms">
        <div
          v-for="platform in sharePlatforms"
          :key="platform.name"
          class="platform-item"
          :style="{ '--platform-color': platform.color }"
          @click="handleShare(platform)"
        >
          <div class="platform-icon">
            <el-icon :size="28">
              <component :is="iconMap[platform.icon]" />
            </el-icon>
          </div>
          <span class="platform-name">{{ platform.name }}</span>
        </div>
      </div>

      <div class="share-url-section">
        <div class="url-label">
          <el-icon><Link /></el-icon>
          <span>文章链接</span>
        </div>
        <div class="url-input-wrapper">
          <el-input
            :model-value="url || window.location.href"
            readonly
            size="large"
            class="url-input"
          />
          <el-button
            type="primary"
            :icon="DocumentCopy"
            @click="handleCopyUrl"
            class="copy-btn"
          >
            复制
          </el-button>
        </div>
      </div>
    </div>

    <template #footer>
      <el-button :icon="Close" @click="handleClose">关闭</el-button>
    </template>
  </el-dialog>
</template>

<style scoped>
.share-content {
  padding: 10px 0;
}

.share-header {
  text-align: center;
  margin-bottom: 30px;
}

.share-header h4 {
  font-size: 18px;
  font-weight: 600;
  color: #1a202c;
  margin: 0;
}

html.dark .share-header h4 {
  color: #f7fafc;
}

.share-platforms {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 30px;
}

.platform-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px 16px;
  background: #f8fafc;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;
}

html.dark .platform-item {
  background: #2d3748;
  border-color: #4a5568;
}

.platform-item:hover {
  transform: translateY(-4px);
  border-color: var(--platform-color);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.12);
}

.platform-icon {
  margin-bottom: 10px;
  color: var(--platform-color);
  transition: transform 0.3s;
}

.platform-item:hover .platform-icon {
  transform: scale(1.1);
}

.platform-name {
  font-size: 14px;
  font-weight: 500;
  color: #4a5568;
}

html.dark .platform-name {
  color: #cbd5e0;
}

.share-url-section {
  background: #f8fafc;
  border-radius: 12px;
  padding: 20px;
}

html.dark .share-url-section {
  background: #1a202c;
}

.url-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #718096;
  margin-bottom: 12px;
}

html.dark .url-label {
  color: #a0aec0;
}

.url-input-wrapper {
  display: flex;
  gap: 12px;
}

.url-input {
  flex: 1;
}

.url-input :deep(.el-input__inner) {
  color: #4a5568;
}

html.dark .url-input :deep(.el-input__inner) {
  color: #cbd5e0;
  background: #2d3748;
}

.copy-btn {
  white-space: nowrap;
}

@media (max-width: 768px) {
  .share-platforms {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
