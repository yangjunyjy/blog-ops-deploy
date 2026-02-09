/**
 * SSH WebSocket 连接工具类
 */
export class SSHWebSocket {
  constructor(options) {
    this.ws = null
    this.url = options.url
    this.onMessage = options.onMessage
    this.onOpen = options.onOpen
    this.onClose = options.onClose
    this.onError = options.onError
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 3
    // this.keepAliveInterval = null
    this.lastInputTime = Date.now()
  }

  /**
   * 连接 WebSocket
   */
  connect() {
    try {
      this.ws = new WebSocket(this.url)

      this.ws.onopen = () => {
        console.log('SSH WebSocket connected')
        this.reconnectAttempts = 0
        if (this.onOpen) {
          this.onOpen()
        }
        // 启动保活机制：每 30 秒发送一次
        this.startKeepAlive()
      }

      this.ws.onmessage = (event) => {
        if (this.onMessage) {
          this.onMessage(event.data)
        }
      }

      this.ws.onclose = (event) => {
        console.log('SSH WebSocket closed:', event.code, event.reason)
        this.stopKeepAlive()
        if (this.onClose) {
          this.onClose(event)
        }
        // 尝试重连
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
          this.reconnectAttempts++
          setTimeout(() => {
            console.log(`Reconnecting... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`)
            this.connect()
          }, 3000)
        }
      }

      this.ws.onerror = (error) => {
        console.error('SSH WebSocket error:', error)
        this.stopKeepAlive()
        if (this.onError) {
          this.onError(error)
        }
      }
    } catch (error) {
      console.error('Failed to create WebSocket:', error)
      if (this.onError) {
        this.onError(error)
      }
    }
  }

  /**
   * 发送数据到 SSH
   * @param {string|Uint8Array} data - 要发送的数据
   */
  send(data) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      console.log('Sending data to SSH:', data.length, 'bytes, first char code:', data.charCodeAt(0))
      this.lastInputTime = Date.now()

      // 检测是否是回车键，如果是则发送特殊字节序列触发刷新
      if (typeof data === 'string' && (data === '\r' || data === '\n')) {
        // 发送特殊字节序列 0xEF 0xBF 0x80（U+2400 在 UTF-8 中的编码）
        const flushSignal = new Uint8Array([0xEF, 0xBF, 0x80])
        this.ws.send(flushSignal)
        this.ws.send(data)
      } else {
        this.ws.send(data)
      }
    } else {
      console.warn('WebSocket is not connected, state:', this.getReadyState())
    }
  }

  /**
   * 启动保活机制
   */
  startKeepAlive() {
    this.stopKeepAlive()
    this.keepAliveInterval = setInterval(() => {
      const timeSinceLastInput = Date.now() - this.lastInputTime
      console.log(`Keep alive check: ${timeSinceLastInput}ms since last input`)

      // 发送空字符串保持连接活跃（WebSocket 层面的保活）
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        try {
          this.ws.send('')
        } catch (e) {
          console.warn('Keep alive send failed:', e)
        }
      }
    }, 30000) // 每 30 秒发送一次
  }

  
  //停止保活机制
  
  stopKeepAlive() {
    if (this.keepAliveInterval) {
      clearInterval(this.keepAliveInterval)
      this.keepAliveInterval = null
    }
  }

  /**
   * 调整终端窗口大小
   * @param {number} rows - 行数
   * @param {number} cols - 列数
   */
  resize(rows, cols) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      const message = JSON.stringify({
        type: 'resize',
        rows: rows,
        cols: cols
      })
      this.ws.send(message)
    }
  }

  /**
   * 断开连接
   */
  disconnect() {
    this.stopKeepAlive()
    if (this.ws) {
      this.reconnectAttempts = this.maxReconnectAttempts // 阻止自动重连
      this.ws.close()
      this.ws = null
    }
  }

  /**
   * 获取连接状态
   * @returns {boolean} 是否已连接
   */
  isConnected() {
    return this.ws && this.ws.readyState === WebSocket.OPEN
  }

  /**
   * 获取连接状态字符串
   * @returns {string} 连接状态
   */
  getReadyState() {
    if (!this.ws) return 'DISCONNECTED'
    switch (this.ws.readyState) {
      case WebSocket.CONNECTING:
        return 'CONNECTING'
      case WebSocket.OPEN:
        return 'CONNECTED'
      case WebSocket.CLOSING:
        return 'CLOSING'
      case WebSocket.CLOSED:
        return 'CLOSED'
      default:
        return 'UNKNOWN'
    }
  }
}

export default SSHWebSocket
