package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// WriteWait はメッセージ送信のタイムアウト
	WriteWait = 10 * time.Second

	// PongWait はpongメッセージを待つタイムアウト
	PongWait = 60 * time.Second

	// PingPeriod はpingメッセージを送信する間隔（PongWaitより短くする必要がある）
	PingPeriod = (PongWait * 9) / 10

	// MaxMessageSize は受信するメッセージの最大サイズ
	MaxMessageSize = 512 * 1024 // 512KB
)

// Client はWebSocket接続のクライアント
type Client struct {
	// Hub はこのクライアントが所属するHub
	Hub *Hub

	// Conn はWebSocket接続
	Conn *websocket.Conn

	// UserID はクライアントのユーザーID
	UserID uuid.UUID

	// Send はクライアントへの送信メッセージチャネル
	Send chan []byte

	// mu は並行アクセスの保護用
	mu sync.Mutex
}

// Hub はアクティブなクライアントを管理し、メッセージをブロードキャストする
type Hub struct {
	// clients は登録されたクライアント
	clients map[*Client]bool

	// userClients はユーザーIDごとのクライアントマッピング
	userClients map[uuid.UUID]map[*Client]bool

	// broadcast はすべてのクライアントへのブロードキャストメッセージチャネル
	broadcast chan []byte

	// register は新しいクライアントの登録リクエストチャネル
	register chan *Client

	// unregister はクライアントの登録解除リクエストチャネル
	unregister chan *Client

	// mu は並行アクセスの保護用
	mu sync.RWMutex
}

// NewHub は新しいHubを作成する
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		userClients: make(map[uuid.UUID]map[*Client]bool),
		broadcast:   make(chan []byte, 256),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
	}
}

// Run はHubのメインループを実行する
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient はクライアントを登録する
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true

	// ユーザーIDごとのマッピングに追加
	if _, ok := h.userClients[client.UserID]; !ok {
		h.userClients[client.UserID] = make(map[*Client]bool)
	}
	h.userClients[client.UserID][client] = true

	log.Printf("Client registered: userID=%s, total clients=%d", client.UserID, len(h.clients))
}

// Register はクライアントを登録する（エクスポート版）
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// unregisterClient はクライアントの登録を解除する
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)

		// ユーザーIDごとのマッピングから削除
		if userClients, ok := h.userClients[client.UserID]; ok {
			delete(userClients, client)
			if len(userClients) == 0 {
				delete(h.userClients, client.UserID)
			}
		}

		close(client.Send)
		log.Printf("Client unregistered: userID=%s, total clients=%d", client.UserID, len(h.clients))
	}
}

// broadcastMessage はすべてのクライアントにメッセージをブロードキャストする
func (h *Hub) broadcastMessage(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.Send <- message:
		default:
			// チャネルがフルの場合はクライアントを切断
			close(client.Send)
			delete(h.clients, client)
		}
	}
}

// SendToUser は特定のユーザーにメッセージを送信する
func (h *Hub) SendToUser(userID uuid.UUID, message Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	userClients, ok := h.userClients[userID]
	if !ok {
		// ユーザーが接続していない場合はエラーではなくログのみ
		log.Printf("No active connections for user: %s", userID)
		return nil
	}

	for client := range userClients {
		select {
		case client.Send <- data:
		default:
			// チャネルがフルの場合はスキップ
			log.Printf("Failed to send message to client: userID=%s", userID)
		}
	}

	return nil
}

// BroadcastToAll はすべてのクライアントにメッセージをブロードキャストする
func (h *Hub) BroadcastToAll(message Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	h.broadcast <- data
	return nil
}

// GetConnectedUserCount は接続中のユーザー数を返す
func (h *Hub) GetConnectedUserCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.userClients)
}

// GetTotalConnectionCount は総接続数を返す
func (h *Hub) GetTotalConnectionCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// IsUserConnected は指定されたユーザーが接続しているかを確認する
func (h *Hub) IsUserConnected(userID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	userClients, ok := h.userClients[userID]
	return ok && len(userClients) > 0
}

// ReadPump はクライアントからのメッセージを読み取る
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(PongWait))
	c.Conn.SetReadLimit(MaxMessageSize)
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// クライアントからのメッセージを処理（必要に応じて）
		log.Printf("Received message from user %s: %s", c.UserID, string(message))

		// 現時点ではクライアントからのメッセージは処理しないが、
		// 将来的にはここでメッセージを処理することができる
	}
}

// WritePump はクライアントへメッセージを書き込む
func (c *Client) WritePump() {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				// Hubがチャネルを閉じた
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// キューに溜まっているメッセージをバッチで送信
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
