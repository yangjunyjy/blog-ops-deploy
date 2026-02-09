package ssh

import (
	"fmt"
	"io"
	"log"
	"my-blog-backend/internal/pkg/logger"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type Session struct {
	ID            string
	Client        *SSHClient
	SSHClient     *ssh.Session
	stdin         io.WriteCloser
	outputReader  io.Reader
	WsConn        *websocket.Conn
	InputChan     chan []byte
	OutputChan    chan []byte
	FlushChan     chan struct{} // 新增：用于立即刷新批量缓冲区
	Done          chan struct{}
	mu            sync.Mutex
	wg            sync.WaitGroup
	active        bool
	lastInputTime time.Time // 新增：记录最后输入时间
}

type PtyConfig struct {
	Term string
	Rows int
	Cols int
}

func NewSession(client *SSHClient, conn *websocket.Conn, id string) *Session {
	return &Session{
		ID:            id,
		Client:        client,
		WsConn:        conn,
		InputChan:     make(chan []byte, 2048),
		OutputChan:    make(chan []byte, 65536), // 增加到 64KB 缓冲区
		FlushChan:     make(chan struct{}, 10),  // 用于立即刷新批量缓冲区，增加容量避免阻塞
		Done:          make(chan struct{}),
		active:        true,
		lastInputTime: time.Now(), // 初始化为当前时间
	}
}

func (s *Session) Start(cfg PtyConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.active {
		return fmt.Errorf("连接已关闭")
	}

	// 创建会话
	session, err := s.Client.client.NewSession()
	if err != nil {
		return fmt.Errorf("创建会话失败: %v", err)
	}
	s.SSHClient = session

	// 设置伪终端
	logger.Info("设置伪终端格式", logger.String("format", fmt.Sprintf("RequestPty: term=%s, rows=%d, cols=%d", cfg.Term, cfg.Rows, cfg.Cols)))
	if err := session.RequestPty(cfg.Term, cfg.Rows, cfg.Cols, ssh.TerminalModes{
		// ssh.ECHO:          0, // 不设置 ECHO，让服务器自动处理
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}); err != nil {
		session.Close()
		return fmt.Errorf("设置伪终端失败: %v", err)
	}

	// 获取标准输入用于写入命令
	s.stdin, err = session.StdinPipe()
	if err != nil {
		session.Close()
		return fmt.Errorf("获取标准输入出错：%v", err)
	}

	// 获取标准输出和标准错误用于读取输出
	stdout, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		return fmt.Errorf("获取标准输出出错：%v", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		session.Close()
		return fmt.Errorf("获取标准错误出错：%v", err)
	}

	// 将输出和错误合并到一个 io.Reader
	s.outputReader = io.MultiReader(stdout, stderr)

	// 启动会话
	if err3 := session.Shell(); err3 != nil {
		session.Close()
		return fmt.Errorf("启动会话失败: %v", err3)
	}

	// 启动携程
	s.wg.Add(2)
	go s.HandleInput()
	go s.handleOutput()

	return nil
}

// 处理websocket 输入并写入ssh
func (s *Session) HandleInput() {
	defer s.wg.Done()
	inputCount := 0

	// 启动超时检测协程
	go s.detectInputTimeout()

	for {
		select {
		case data, ok := <-s.InputChan:
			if !ok {
				return
			}
			inputCount++

			// 更新最后输入时间
			s.mu.Lock()
			s.lastInputTime = time.Now()
			s.mu.Unlock()

			// 只在数据较大时记录日志
			if len(data) > 10 {
				log.Printf("SSH input [%d]: %d bytes, data: %q", inputCount, len(data), string(data))
			}
			if _, err := s.stdin.Write(data); err != nil {
				log.Printf("SSH input write error: %v", err)
				// SSH 写入失败，立即关闭整个会话
				s.Close()
				return
			}
		case <-s.Done:
			log.Printf("SSH input handler stopped for session: %s, total inputs: %d", s.ID, inputCount)
			return
		}
	}
}

// 检测输入超时，超过 5 分钟无输入则断开连接
func (s *Session) detectInputTimeout() {
	ticker := time.NewTicker(60 * time.Second) // 每分钟检查一次
	defer ticker.Stop()

	timeoutDuration := 5 * time.Minute // 5 分钟超时

	for {
		select {
		case <-s.Done:
			return
		case <-ticker.C:
			s.mu.Lock()
			lastInput := s.lastInputTime
			timeSinceLastInput := time.Since(lastInput)
			s.mu.Unlock()

			log.Printf("Session %s timeout check: %v since last input", s.ID, timeSinceLastInput)

			// 超过 5 分钟无输入，断开连接
			if timeSinceLastInput > timeoutDuration {
				log.Printf("Session %s timeout after %v of inactivity, closing connection", s.ID, timeoutDuration)

				// 发送超时消息到前端
				if s.WsConn != nil {
					timeoutMsg := fmt.Sprintf("\r\n\033[31m[会话超时] 检测到 %d 分钟无操作，连接已断开\033[0m\r\n", int(timeoutDuration.Minutes()))
					s.WsConn.SetWriteDeadline(time.Now().Add(10 * time.Second))
					s.WsConn.WriteMessage(websocket.TextMessage, []byte(timeoutMsg))
				}

				// 关闭会话
				s.Close()
				return
			}
		}
	}
}

// 处理ssh输出并写入到OutputChan中
func (s *Session) handleOutput() {
	defer s.wg.Done()

	log.Printf("SSH session output handler started for session: %s", s.ID)

	buf := make([]byte, 8192)
	outputCount := 0

	for {
		// 阻塞读取数据，直到有数据或连接关闭
		n, err := s.outputReader.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("SSH output read error for session %s: %v", s.ID, err)
			}
			return
		}

		if n > 0 {
			output := make([]byte, n)
			copy(output, buf[:n])
			outputCount++

			// 只在数据较小时记录日志，避免长文本日志淹没
			if n <= 100 {
				log.Printf("SSH output [%d]: %d bytes, content: %q", outputCount, n, string(output))
			} else {
				log.Printf("SSH output [%d]: %d bytes (truncated)", outputCount, n)
			}

			// 发送到通道，如果通道满了就丢弃旧数据
			select {
			case s.OutputChan <- output:
			case <-s.Done:
				return
			default:
				// 通道满时，丢弃这个数据包
				log.Printf("Output channel full, dropping packet [%d]", outputCount)
			}
		}
	}

}

func (s *Session) ReSize(rows, cols int) error {
	logger.Info("调整窗口大小", logger.Int("rows", rows), logger.Int("cols", cols))
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.SSHClient == nil {
		return fmt.Errorf("ssh会话未初始化")
	}

	return s.SSHClient.WindowChange(rows, cols)
}

func (s *Session) Close() error {
	s.mu.Lock()

	if !s.active {
		s.mu.Unlock()
		return nil
	}

	s.active = false

	// 关闭 Done channel
	select {
	case <-s.Done:
	default:
		close(s.Done)
	}

	s.mu.Unlock()

	// 关闭 stdin
	if s.stdin != nil {
		s.stdin.Close()
	}

	// 关闭 InputChan
	close(s.InputChan)
	s.WsConn.Close()

	// 等待所有协程退出
	s.wg.Wait()

	// 关闭 SSH 会话
	if s.SSHClient != nil {
		s.SSHClient.Close()
	}

	return nil
}

func (s *Session) IsActive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.active
}
