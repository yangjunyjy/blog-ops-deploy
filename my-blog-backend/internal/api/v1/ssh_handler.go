package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	dtoResponse "my-blog-backend/internal/api/v1/dto/response"
	"my-blog-backend/internal/services"
	"my-blog-backend/internal/ssh"

	"github.com/gin-gonic/gin"

	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	wsReadLimit  = 64 * 1024
	wsWriteWait  = 10 * time.Second
	wsPongWait   = 70 * time.Second
	wsPingPeriod = 30 * time.Second
	wsFlushDelay = 100 * time.Millisecond
)

type SshHandler struct {
	hostService *services.HostService
	pool        *ssh.Pool
	sessions    map[string]*ssh.Session
	mu          sync.RWMutex
}

func NewSshHandler(hostService *services.HostService, pool *ssh.Pool) *SshHandler {
	return &SshHandler{
		hostService: hostService,
		pool:        pool,
		sessions:    make(map[string]*ssh.Session),
	}
}

// GetSessions 获取 sessions map 供其他 handler 共享
func (h *SshHandler) GetSessions() map[string]*ssh.Session {
	return h.sessions
}

// WebSocketConnect WebSocket 连接
// @Summary WebSocket 连接
// @Tags SSH终端
// @Param host_id path string true "主机ID"
// @Param session_id query string true "会话ID"
// @Success 101
// @Router /api/v1/ssh/connect/{host_id} [get]
func (h *SshHandler) WebSocketConnect(c *gin.Context) {
	hostIDStr := c.Param("host_id")
	hostID, err := parseUint(hostIDStr)
	if err != nil {
		log.Printf("WebSocket connect error: invalid host_id: %s", hostIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的主机ID"})
		return
	}

	sessionID := c.Query("session_id")
	if sessionID == "" {
		log.Printf("WebSocket connect error: missing session_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少会话ID"})
		return
	}

	log.Printf("WebSocket connection request: hostID=%d, sessionID=%s", hostID, sessionID)

	// 升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	log.Printf("WebSocket connection established for hostID=%d", hostID)

	// 获取 SSH 配置
	sshConfig, err := h.hostService.GetSSHConfig(hostID)
	if err != nil {
		log.Printf("Get SSH config error: %v", err)
		conn.WriteMessage(ws.TextMessage, []byte("获取主机配置失败: "+err.Error()))
		conn.Close()
		return
	}

	log.Printf("SSH config: host=%s, port=%d, user=%s, authType=%v",
		sshConfig.Host, sshConfig.Port, sshConfig.Username, sshConfig.AuthType)

	// 从连接池获取 SSH 客户端
	sshClient, err := h.pool.Get(c.Request.Context(), sshConfig, hostID)
	if err != nil {
		log.Printf("Get SSH client from pool error: %v", err)
		conn.WriteMessage(ws.TextMessage, []byte("SSH 连接失败: "+err.Error()))
		conn.Close()
		return
	}

	log.Printf("SSH client obtained successfully for hostID=%d", hostID)

	// 创建会话
	session := ssh.NewSession(sshClient, conn, sessionID)

	h.mu.Lock()
	h.sessions[sessionID] = session
	h.mu.Unlock()

	log.Printf("SSH session created: sessionID=%s", sessionID)

	// 启动会话
	if err := session.Start(ssh.PtyConfig{
		Term: "xterm",
		Rows: 50,
		Cols: 150,
	}); err != nil {
		log.Printf("Session start error: %v", err)
		conn.WriteMessage(ws.TextMessage, []byte("启动会话失败: "+err.Error()))
		session.Close()
		return
	}

	log.Printf("SSH session started successfully for hostID=%d", hostID)

	// 发送连接成功消息
	conn.WriteMessage(ws.TextMessage, []byte("SSH 会话已建立，连接到远程主机...\r\n"))
	log.Printf("Connection message sent to client")

	// 使用 WaitGroup 管理协程
	var wg sync.WaitGroup

	// 启动 WebSocket 读取协程
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("readWebSocket goroutine started")
		h.readWebSocket(conn, session)
	}()

	// 启动 SSH 输出协程
	wg.Add(1)
	go func() {
		defer wg.Done()
		h.writeWebSocket(conn, session)
	}()

	// 注意：窗口大小调整功能已移除，因为 xterm.js 会自动处理
	// handleResize 与 readWebSocket 冲突，都尝试从同一个 conn 读取消息
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	h.handleResize(conn, session)
	// }()

	// 等待协程完成
	wg.Wait()

	// 清理会话
	h.mu.Lock()
	delete(h.sessions, sessionID)
	h.mu.Unlock()
	session.Close()
}

// readWebSocket 从 WebSocket 读取数据并发送到 SSH
func (h *SshHandler) readWebSocket(conn *ws.Conn, session *ssh.Session) {
	log.Printf("readWebSocket: starting to read...")
	for {
		log.Printf("readWebSocket: waiting for message...")
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			return
		}

		// 记录输入，方便调试
		log.Printf("Received from WebSocket: %d bytes, content: %q", len(message), string(message))

		// 检测特殊刷新字符 U+2400（空白），用于触发输出刷新
		// UTF-8 编码: 0xEF 0xBF 0x80
		isFlushSignal := len(message) >= 3 && message[0] == 0xEF && message[1] == 0xBF && message[2] == 0x80
		if isFlushSignal {
			log.Printf("Flush signal detected via special character")
			select {
			case session.FlushChan <- struct{}{}:
				log.Printf("Flush signal sent to channel")
			default:
				log.Printf("Flush channel full, skipping")
			}
			continue // 不发送到 SSH，只用于刷新
		} else {
			// 处理 JSON 消息（resize）
			var msg struct {
				Type string `json:"type"`
				Rows int    `json:"rows"`
				Cols int    `json:"cols"`
			}
			if err := json.Unmarshal(message, &msg); err == nil && msg.Type == "resize" {
				// 调整窗口大小
				if err := session.ReSize(msg.Rows, msg.Cols); err != nil {
					log.Printf("Resize window failed: %v", err)
				} else {
					log.Printf("Resized terminal to %d cols x %d rows", msg.Cols, msg.Rows)
				}
				continue // 不发送到 InputChan，直接 continue
			}

			// 发送到 SSH
			log.Printf("readWebSocket: sending to InputChan...")
			select {
			case session.InputChan <- message:
				log.Printf("readWebSocket: data sent to InputChan")
			case <-session.Done:
				log.Printf("readWebSocket: session done, returning")
				return
			}
		}
	}
}

// writeWebSocket 从 SSH 读取数据并发送到 WebSocket
func (h *SshHandler) writeWebSocket(conn *ws.Conn, session *ssh.Session) {
	// 增加心跳间隔到 30 秒
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// 设置写超时，避免长文本导致阻塞
	conn.SetWriteDeadline(time.Now().Add(60 * time.Second))

	dataSent := 0
	// 批量发送缓冲区（限制为 16KB，避免单次发送过大）
	batchBuffer := make([]byte, 0, 16384)
	lastFlush := time.Now()

	// 刷新延迟定时器
	flushTimer := time.NewTimer(0)
	defer flushTimer.Stop()

	for {
		select {
		case <-session.Done:
			log.Printf("writeWebSocket stopped for session: %s, total data sent: %d", session.ID, dataSent)
			return

		case <-session.FlushChan:
			// 收到刷新信号，启动一个定时器，100ms 后刷新
			// 这样可以让 SSH 输出的数据有时间进入缓冲区
			log.Printf("Flush signal received, scheduling flush in 100ms")
			flushTimer.Stop()
			flushTimer.Reset(100 * time.Millisecond)

		case <-flushTimer.C:
			// 定时器到期，检查并刷新缓冲区
			if len(batchBuffer) > 0 {
				conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
				if err := conn.WriteMessage(ws.TextMessage, batchBuffer); err != nil {
					log.Printf("WebSocket write error (timer flush): %v", err)
					return
				}
				log.Printf("Sent buffered data on timer [%d]: %d bytes", dataSent, len(batchBuffer))
				batchBuffer = batchBuffer[:0]
				lastFlush = time.Now()
			}

		case data := <-session.OutputChan:
			dataSent++

			// 如果单条数据超过 8KB，直接发送不批量
			if len(data) > 8192 {
				conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
				if err := conn.WriteMessage(ws.TextMessage, data); err != nil {
					log.Printf("WebSocket write error (large data): %v", err)
					return
				}
				if len(data) > 100 {
					log.Printf("Sent large data directly [%d]: %d bytes", dataSent, len(data))
				}
				continue
			}

			// 添加到批量缓冲区
			if len(batchBuffer)+len(data) <= cap(batchBuffer) {
				batchBuffer = append(batchBuffer, data...)
			} else {
				// 缓冲区满，先发送现有数据
				conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
				if err := conn.WriteMessage(ws.TextMessage, batchBuffer); err != nil {
					log.Printf("WebSocket write error (buffer flush): %v", err)
					return
				}
				if len(batchBuffer) > 100 {
					log.Printf("Sent buffer flush [%d]: %d bytes", dataSent, len(batchBuffer))
				}
				batchBuffer = batchBuffer[:0]
				batchBuffer = append(batchBuffer, data...)
			}

			// 计算是否需要刷新
			shouldFlush := len(batchBuffer) >= 2048 || // 缓冲区 >= 2KB
				time.Since(lastFlush) >= 50*time.Millisecond || // 或超过 50ms
				len(data) < 256 // 小数据（如按键）立即发送

			if shouldFlush && len(batchBuffer) > 0 {
				// 每次写入前更新超时时间
				conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
				if err := conn.WriteMessage(ws.TextMessage, batchBuffer); err != nil {
					log.Printf("WebSocket write error: %v", err)
					return
				}
				// 只在数据较大时记录日志
				if len(batchBuffer) > 100 {
					log.Printf("Sent data to WebSocket [%d]: %d bytes", dataSent, len(batchBuffer))
				}
				batchBuffer = batchBuffer[:0] // 清空缓冲区
				lastFlush = time.Now()
			}

		case <-ticker.C:
			// ticker 只用于定期刷新缓冲区，不需要发送 Ping
			// 保活由前端通过发送数据实现
		}
	}
}

// handleResize 处理窗口大小调整
func (h *SshHandler) handleResize(conn *ws.Conn, session *ssh.Session) {
	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// 只处理文本消息
		if msgType == ws.TextMessage {
			// 解析调整窗口大小的指令
			// 格式: {"type": "resize", "rows": 24, "cols": 80}
			type ResizeMsg struct {
				Type string `json:"type"`
				Rows int    `json:"rows"`
				Cols int    `json:"cols"`
			}

			var msg ResizeMsg
			if err := json.Unmarshal(data, &msg); err != nil {
				continue // 忽略无效消息
			}

			// 处理调整窗口大小的请求
			if msg.Type == "resize" && msg.Rows > 0 && msg.Cols > 0 {
				if err := session.ReSize(msg.Rows, msg.Cols); err != nil {
					log.Printf("Resize window failed: %v", err)
				}
			}
		}
	}
}

// CloseSession 关闭会话
// @Summary 关闭会话
// @Tags SSH终端
// @Param session_id path string true "会话ID"
// @Success 200 {object} dtoResponse.Response
// @Router /api/v1/ssh/sessions/{session_id} [delete]
func (h *SshHandler) CloseSession(c *gin.Context) {
	sessionID := c.Param("session_id")

	h.mu.RLock()
	session, exists := h.sessions[sessionID]
	h.mu.RUnlock()

	if !exists {
		dtoResponse.Error(c, 404, "会话不存在", nil)
		return
	}

	if err := session.Close(); err != nil {
		dtoResponse.Error(c, 500, "关闭会话失败", err)
		return
	}

	h.mu.Lock()
	delete(h.sessions, sessionID)
	h.mu.Unlock()

	dtoResponse.Success(c, nil, "会话已关闭")
}

// ListSessions 列出活跃会话
// @Summary 列出活跃会话
// @Tags SSH终端
// @Success 200 {object} dtoResponse.Response
// @Router /api/v1/ssh/sessions [get]
func (h *SshHandler) ListSessions(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	sessions := make([]map[string]interface{}, 0, len(h.sessions))
	for sessionID, session := range h.sessions {
		sessions = append(sessions, map[string]interface{}{
			"session_id": sessionID,
			"host_id":    session.Client.GetHostID(),
			"active":     session.IsActive(),
		})
	}

	dtoResponse.Success(c, sessions, "获取成功")
}

// 辅助函数
func parseUint(s string) (uint, error) {
	var id uint64
	_, err := fmt.Sscanf(s, "%d", &id)
	return uint(id), err
}
