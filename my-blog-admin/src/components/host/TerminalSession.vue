<template>
  <div class="terminal-session">
    <!-- 会话标签页 -->
    <div class="session-tabs">
      <div v-for="tab in tabs" :key="tab.id"
        :class="['tab-item', { active: activeTabId === tab.id, 'connecting': tab.connecting }]"
        @click="handleTabClick(tab)">
        <el-icon class="tab-icon">
          <Monitor />
        </el-icon>
        <span class="tab-title">{{ tab.title }}</span>
        <span class="tab-status" v-if="tab.connecting">连接中...</span>
        <el-icon class="tab-close" @click.stop="handleTabClose(tab.id)" v-if="tabs.length > 1">
          <Close />
        </el-icon>
      </div>
      <div class="tab-add" @click="handleAddTab" title="新建会话">
        <el-icon>
          <Plus />
        </el-icon>
      </div>
    </div>

    <!-- 会话内容区 -->
    <div class="session-content">
      <div v-for="tab in tabs" :key="tab.id" v-show="activeTabId === tab.id" class="tab-pane">
        <div v-if="tab.connecting" class="connecting-state">
          <el-icon class="loading-icon" :size="40">
            <Loading />
          </el-icon>
          <p>正在连接到 {{ tab.host.name }} ({{ tab.host.ip }}:{{ tab.host.port }})</p>
          <p class="auth-type">认证方式: {{ tab.host.auth_type === 'password' ? '密码' : '密钥' }}</p>
        </div>
        <div v-else class="terminal-container">
          <!-- 连接状态指示器 -->
          <div class="connection-status" :class="getStatusClass(tab)">
            <el-icon>
              <component :is="getStatusIcon(tab)" />
            </el-icon>
            <span>{{ getStatusText(tab) }}</span>
          </div>

          <!-- xterm.js 终端容器 -->
          <div class="terminal-wrapper" :id="`terminal-${tab.id}`" ref="terminalWrapper"></div>
          <div class="terminal-footer">
            <el-button size="small" :icon="Connection" circle @click="handleReconnect(tab)" title="重连" />
            <el-button size="small" :icon="Folder" @click="handleOpenFileTransfer(tab)" title="文件传输" />
          </div>
          <!-- 添加 FileTransferPanel 组件 -->
          <FileTransferPanel v-if="currentTab" v-model:visible="fileTransferVisible" :session-id="currentTab.sshSessionId" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineOptions({
  name: 'TerminalSession'
})

import { ref, nextTick, watch, onUnmounted, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Monitor, Close, Plus, Loading, Connection, SuccessFilled, CircleCloseFilled, Folder
} from '@element-plus/icons-vue'
import SSHWebSocket from '@/utils/sshWebSocket'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import FileTransferPanel from '@/views/host/components/FileTransferPanel.vue'
import '@xterm/xterm/css/xterm.css'

const fileTransferVisible = ref(false)
const currentTab = ref(null)

const handleOpenFileTransfer = (tab) => {
  currentTab.value = tab
  fileTransferVisible.value = true
  console.log('打开文件传输面板，SSH Session ID:', tab.sshSessionId)
}

// 简单的 UUID 生成函数
function generateUUID() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    const r = Math.random() * 16 | 0
    const v = c === 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}

const props = defineProps({
  initialHost: {
    type: Object,
    default: () => null
  }
})

const emit = defineEmits(['close'])

// 状态管理
const tabs = ref([])
const activeTabId = ref('')
const tabCounter = ref(0)
const terminalWrapper = ref(null)
const wsConnections = ref(new Map()) // 存储每个标签的 WebSocket 连接
const terminals = ref(new Map()) // 存储每个标签的 xterm 实例
const fitAddons = ref(new Map()) // 存储每个标签的 FitAddon 实例

// 获取 WebSocket URL
const getWebSocketUrl = (hostId, sshSessionId) => {
  const isDev = import.meta.env.DEV

  if (isDev) {
    return `ws://localhost:3000/rbac/ssh/connect/${hostId}?session_id=${sshSessionId}`
  } else {
    const apiBaseUrl = 'http://localhost:8081/api/v1'

    let protocol = 'ws:'
    let host = 'localhost:8081'

    try {
      const url = new URL(apiBaseUrl)
      protocol = url.protocol === 'https:' ? 'wss:' : 'ws:'
      host = url.host
    } catch (e) {
      console.error('解析 API URL 失败:', e)
    }

    return `${protocol}//${host}/api/v1/rbac/ssh/connect/${hostId}?session_id=${sshSessionId}`
  }
}

// 创建 xterm 实例
const createTerminal = (tab) => {
  const container = document.getElementById(`terminal-${tab.id}`)
  if (!container) {
    console.error('Terminal container not found:', tab.id)
    return
  }

  // 获取容器的实际宽度
  const cont = document.getElementById(`terminal-${tab.id}`)
  const containerWidth = cont.offsetWidth
  // 计算终端列数（根据屏幕宽度）
  const calculateCols = (width) => {
    const charWidth = 8.5 // 每个字符的大致宽度（像素）
    const data = Math.floor(width / charWidth)
    console.log("屏幕宽度为", data);
    return data
  }

  // 清空容器
  container.innerHTML = ''

  // 创建 Terminal 实例
  const term = new Terminal({
    cursorBlink: true,
    cursorStyle: 'block',
    fontSize: 13,
    fontFamily: 'Consolas, "Courier New", monospace',
    scrollback: 1000,
    tabStopWidth: 4,
    allowProposedApi: true,
    fastScrollSensitivity: 5,
    fastScrollModifier: 'alt',
    cols: calculateCols(containerWidth),
    disableStdin: false, // 启用输入
    rendererType: 'canvas',
    theme: {
      background: '#0c0c0c',
      foreground: '#d4d4d4',
      cursor: '#4ec9b0',
      cursorAccent: '#0c0c0c',
      selection: '#264f78',
      black: '#000000',
      red: '#cd3131',
      green: '#0dbc79',
      yellow: '#e5e510',
      blue: '#2472c8',
      magenta: '#bc3fbc',
      cyan: '#11a8cd',
      white: '#e5e5e5',
      brightBlack: '#666666',
      brightRed: '#f14c4c',
      brightGreen: '#23d18b',
      brightYellow: '#f5f543',
      brightBlue: '#3b8eea',
      brightMagenta: '#d670d6',
      brightCyan: '#29b8db',
      brightWhite: '#ffffff'
    }
  })

  // 创建 FitAddon
  const fitAddon = new FitAddon()
  term.loadAddon(fitAddon)

  // 打开终端
  term.open(container)

  // 适应容器大小
  fitAddon.fit()

  // 获取实际大小并发送到 SSH
  const actualCols = term.cols
  const actualRows = term.rows
  console.log('Terminal fitted size:', { cols: actualCols, rows: actualRows })

  // 聚焦终端，确保可以接收键盘输入
  term.focus()

  // 点击容器时聚焦终端
  container.addEventListener('click', () => {
    term.focus()
  })

  // 终端输入处理 - 实时发送所有字符
  // 刷新机制已内置在 sshWebSocket.send() 中，检测回车键会自动触发
  term.onData((data) => {
    const connection = wsConnections.value.get(tab.id)

    if (!connection || !connection.sshWs.isConnected()) {
      console.warn('SSH not connected, cannot send data')
      return
    }

    // 立即发送所有字符到 SSH
    // SSH 会处理回显和编辑
    connection.sshWs.send(data)
  })

  // 保存实例
  terminals.value.set(tab.id, term)
  fitAddons.value.set(tab.id, fitAddon)

  console.log('Terminal created for tab:', tab.id)

  const connection = wsConnections.value.get(tab.id)
  if (connection && connection.sshWs.isConnected()) {
    // 使用 fit 后的实际大小
    const resizeMsg = JSON.stringify({
      type: 'resize',
      cols: term.cols,
      rows: term.rows
    })
    connection.sshWs.send(resizeMsg)
  }
}

// 销毁终端
const destroyTerminal = (tabId) => {
  const term = terminals.value.get(tabId)
  const fitAddon = fitAddons.value.get(tabId)

  if (term) {
    try {
      term.dispose()
    } catch (e) {
      console.error('Error disposing terminal:', e)
    }
    terminals.value.delete(tabId)
  }

  if (fitAddon) {
    fitAddons.value.delete(tabId)
  }
}

// 创建新标签
const createTab = (host) => {
  tabCounter.value++
  const tabId = `tab-${tabCounter.value}`

  const sameHostTabs = tabs.value.filter(t => t.host.id === host.id)
  const sessionIndex = sameHostTabs.length + 1
  const title = sessionIndex === 1 ? host.name : `${host.name} (${sessionIndex})`

  // 生成真正的 SSH session ID（UUID）
  const sshSessionId = generateUUID()

  const newTab = {
    id: tabId,
    title: title,
    host: host,
    connecting: true,
    connected: false,
    sessionIndex: sessionIndex,
    sshSessionId: sshSessionId // 保存真正的 SSH session ID
  }
  tabs.value.push(newTab)
  activeTabId.value = tabId

  // 建立 WebSocket 连接
  connectSSH(newTab)
}

// 建立 SSH 连接
const connectSSH = (tab) => {
  const wsUrl = getWebSocketUrl(tab.host.id, tab.sshSessionId)
  console.log('Connecting to SSH WebSocket:', wsUrl, 'sshSessionId:', tab.sshSessionId)

  const sshWs = new SSHWebSocket({
    url: wsUrl,
    onOpen: () => {
      console.log('WebSocket connection established, waiting for SSH session...')
      // WebSocket 打开，但 SSH session 还未建立
      tab.connecting = true
      tab.connected = false
    },
    onMessage: (data) => {
      handleSSHMessage(tab.id, data)
    },
    onClose: (event) => {
      console.log('SSH WebSocket closed:', event)
      handleConnectionClose(tab.id)
    },
    onError: (error) => {
      console.error('SSH WebSocket error:', error)
      tab.connecting = false
      tab.connected = false
      ElMessage.error('SSH 连接错误: ' + error.message)
    }
  })

  sshWs.connect()
  wsConnections.value.set(tab.id, { sshWs, tab })
}

// 处理 SSH 消息
const handleSSHMessage = (tabId, data) => {
  const tab = tabs.value.find(t => t.id === tabId)
  if (!tab) {
    console.warn('Tab not found:', tabId)
    return
  }

  console.log('SSH message received:', tabId, 'length:', data.length, 'content:', JSON.stringify(data))

  // 任何消息都说明 SSH 会话已经启动，改变状态
  if (tab.connecting) {
    tab.connecting = false
    tab.connected = true
    ElMessage.success(`已连接到 ${tab.host.name}`)

    // 创建 xterm 终端实例
    nextTick(() => {
      createTerminal(tab)
    })
  }

  // 直接写入终端
  const term = terminals.value.get(tabId)
  if (term) {
    console.log('Writing to xterm:', data.length, 'bytes')
    term.write(data)
  } else {
    console.warn('Terminal not found for tab:', tabId)
  }
}

// 处理连接关闭
const handleConnectionClose = (tabId) => {
  const tab = tabs.value.find(t => t.id === tabId)
  if (tab) {
    tab.connecting = false
    tab.connected = false

    // 在终端显示断开连接消息
    const term = terminals.value.get(tabId)
    if (term) {
      term.writeln('\r\n\x1b[31m连接已断开\x1b[0m\r\n')
    }
  }
}

// 获取连接状态文字
const getStatusText = (tab) => {
  if (tab.connecting) return '连接中'
  if (!tab.connected) return '已断开'
  return '已连接'
}

// 获取连接状态类名
const getStatusClass = (tab) => {
  if (tab.connecting) return 'status-connecting'
  if (!tab.connected) return 'status-disconnected'
  return 'status-connected'
}

// 获取状态图标
const getStatusIcon = (tab) => {
  if (tab.connecting) return Loading
  if (!tab.connected) return CircleCloseFilled
  return SuccessFilled
}

// 标签点击
const handleTabClick = (tab) => {
  activeTabId.value = tab.id
  nextTick(() => {
    // 调整终端大小以适应容器
    const fitAddon = fitAddons.value.get(tab.id)
    const term = terminals.value.get(tab.id)
    if (fitAddon && term) {
      fitAddon.fit()

      // 发送新的尺寸到 SSH
      const connection = wsConnections.value.get(tab.id)
      if (connection && connection.sshWs.isConnected()) {
        const resizeMsg = JSON.stringify({
          type: 'resize',
          cols: term.cols,
          rows: term.rows
        })
        connection.sshWs.send(resizeMsg)
      }
    }

    // 聚焦终端
    if (term) {
      term.focus()
    }
  })
}

// 关闭标签
const handleTabClose = async (tabId) => {
  const tab = tabs.value.find(t => t.id === tabId)
  if (tab && tab.connected) {
    try {
      await ElMessageBox.confirm(
        `确定要断开与 ${tab.host.name} 的连接吗？`,
        '确认断开',
        { type: 'warning' }
      )
    } catch {
      return
    }
  }

  // 销毁终端
  destroyTerminal(tabId)

  // 断开 WebSocket 连接
  const connection = wsConnections.value.get(tabId)
  if (connection) {
    connection.sshWs.disconnect()
    wsConnections.value.delete(tabId)
  }

  const index = tabs.value.findIndex(t => t.id === tabId)
  tabs.value.splice(index, 1)

  // 如果关闭的是当前激活的标签，切换到其他标签
  if (activeTabId.value === tabId) {
    if (tabs.value.length > 0) {
      activeTabId.value = tabs.value[Math.max(0, index - 1)].id
      nextTick(() => {
        // 调整新激活标签的终端大小
        const newTab = tabs.value[Math.max(0, index - 1)]
        const fitAddon = fitAddons.value.get(newTab.id)
        if (fitAddon) {
          fitAddon.fit()
        }
        // 聚焦新激活标签的终端
        const term = terminals.value.get(newTab.id)
        if (term) {
          term.focus()
        }
      })
    } else {
      emit('close')
    }
  }
}

// 添加新标签
const handleAddTab = () => {
  if (!props.initialHost) {
    ElMessage.warning('没有可连接的主机')
    return
  }
  createTab(props.initialHost)
}

// 重连
const handleReconnect = (tab) => {
  ElMessageBox.confirm('确定要重新连接吗？', '确认重连', { type: 'info' })
    .then(() => {
      // 销毁旧终端
      destroyTerminal(tab.id)

      // 断开旧连接
      const oldConnection = wsConnections.value.get(tab.id)
      if (oldConnection) {
        oldConnection.sshWs.disconnect()
        wsConnections.value.delete(tab.id)
      }

      // 重新设置连接状态
      tab.connecting = true
      tab.connected = false

      // 重新连接
      connectSSH(tab)
    })
    .catch(() => { })
}

// 添加会话（外部调用）
const addSession = (host) => {
  createTab(host)
}

// 监听初始主机
watch(() => props.initialHost, (host) => {
  if (host && host.id) {
    createTab(host)
  }
}, { immediate: true })

// 窗口大小改变时调整终端
onMounted(() => {
  window.addEventListener('resize', handleWindowResize)
})

const handleWindowResize = () => {
  tabs.value.forEach(tab => {
    if (activeTabId.value === tab.id) {
      const fitAddon = fitAddons.value.get(tab.id)
      const term = terminals.value.get(tab.id)
      if (fitAddon && term) {
        fitAddon.fit()

        // 发送新的尺寸到 SSH
        const connection = wsConnections.value.get(tab.id)
        if (connection && connection.sshWs.isConnected()) {
          const resizeMsg = JSON.stringify({
            type: 'resize',
            cols: term.cols,
            rows: term.rows
          })
          connection.sshWs.send(resizeMsg)
        }
      }
    }
  })
}

// 组件卸载时清理所有连接
onUnmounted(() => {
  window.removeEventListener('resize', handleWindowResize)
  wsConnections.value.forEach((connection) => {
    connection.sshWs.disconnect()
  })
  wsConnections.value.clear()

  terminals.value.forEach((term, tabId) => {
    try {
      term.dispose()
    } catch (e) {
      console.error('Error disposing terminal:', e)
    }
  })
  terminals.value.clear()
  fitAddons.value.clear()
})

// 导出方法供父组件调用
defineExpose({
  addSession
})
</script>

<style scoped lang="scss">
.terminal-session {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #1e1e1e;
  width: 100%;
}

.session-tabs {
  display: flex;
  align-items: center;
  background-color: #2d2d2d;
  border-bottom: 1px solid #3e3e3e;
  padding: 4px 8px;
  gap: 2px;
  overflow-x: auto;
  width: 100%;

  &::-webkit-scrollbar {
    height: 4px;
  }

  &::-webkit-scrollbar-thumb {
    background: #555;
    border-radius: 2px;
  }

  .tab-item {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 12px;
    background-color: #3c3c3c;
    border-radius: 4px 4px 0 0;
    cursor: pointer;
    font-size: 13px;
    color: #cccccc;
    transition: all 0.3s;
    min-width: 120px;
    max-width: 200px;
    position: relative;

    &:hover {
      background-color: #4a4a4a;
    }

    &.active {
      background-color: #1e1e1e;
      color: #ffffff;
      border-top: 2px solid #007acc;
    }

    &.connecting {
      .tab-status {
        display: inline;
        color: #ffcc00;
        font-size: 11px;
        animation: pulse 1.5s infinite;
      }
    }

    .tab-icon {
      font-size: 14px;
    }

    .tab-title {
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .tab-status {
      display: none;
    }

    .tab-close {
      font-size: 14px;
      opacity: 0.6;
      transition: opacity 0.3s;

      &:hover {
        opacity: 1;
        color: #ff6b6b;
      }
    }
  }

  .tab-add {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border-radius: 4px;
    cursor: pointer;
    color: #cccccc;
    transition: all 0.3s;
    flex-shrink: 0;

    &:hover {
      background-color: #4a4a4a;
      color: #ffffff;
    }
  }
}

@keyframes pulse {

  0%,
  100% {
    opacity: 1;
  }

  50% {
    opacity: 0.5;
  }
}

.session-content {
  flex: 1;
  overflow: hidden;
  width: 100%;
}

.tab-pane {
  height: 100%;
  display: flex;
  flex-direction: column;
  width: 100%;
}

.connecting-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #cccccc;

  .loading-icon {
    animation: rotate 1s linear infinite;
    margin-bottom: 20px;
    color: #007acc;
  }
}

.terminal-container {
  flex: 1;
  position: relative;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.connection-status {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 4px;
  font-size: 12px;
  color: #999;
  margin-bottom: 8px;

  &.status-connected {
    color: #67c23a;
    background: rgba(103, 194, 58, 0.1);
  }

  &.status-connecting {
    color: #409eff;
    background: rgba(64, 158, 255, 0.1);
  }

  &.status-disconnected {
    color: #f56c6c;
    background: rgba(245, 108, 108, 0.1);
  }
}

p {
  margin: 8px 0;
  font-size: 14px;
}

.auth-type {
  font-size: 12px;
  color: #999999;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
}

.terminal-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: #0c0c0c;
}

.terminal-wrapper {
  flex: 1;
  overflow: hidden;
  padding: 0;
  position: relative;

  :deep(.xterm) {
    height: 100%;
    padding: 8px;
  }

  :deep(.xterm-viewport) {
    overflow-y: auto !important;
  }

  :deep(.xterm-screen) {
    position: relative;
  }
}

.terminal-footer {
  display: flex;
  justify-content: flex-end;
  padding: 8px 16px;
  background-color: #1e1e1e;
  border-top: 1px solid #3e3e3e;
  flex-shrink: 0;

  :deep(.el-button) {
    background-color: #2d2d2d;
    border-color: #3e3e3e;
    color: #cccccc;

    &:hover {
      background-color: #3e3e3e;
      border-color: #4a4a4a;
      color: #ffffff;
    }

    .el-icon {
      color: inherit;
    }
  }
}

:deep(.el-button) {
  height: 28px;
  padding: 0 12px;
  font-size: 12px;
}
</style>
