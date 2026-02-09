package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[*Client]bool   // 活跃的 WebSocket 客户端
	sessions   map[string]*Client // sessionID -> Client
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type Client struct {
	SessionID string
	HostID    uint
	Conn      *websocket.Conn
	Send      chan []byte
	Hub       *Hub
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		sessions:   make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// 使用通道处理客户端注册/注销，避免竞态条件
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.sessions[client.SessionID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.sessions, client.SessionID)
				close(client.Send)
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) GetClient(sessionID string) (*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	client, ok := h.sessions[sessionID]
	return client, ok
}

func (h *Hub) Broadcast(hostID uint, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.HostID == hostID {
			select {
			case client.Send <- message:
			default:
				// 缓冲区满，断开连接
				h.unregister <- client
			}
		}
	}
}
