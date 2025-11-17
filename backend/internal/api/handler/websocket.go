package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/clearclown/HaiLanGo/backend/internal/websocket"
	"github.com/clearclown/HaiLanGo/backend/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	gorillaws "github.com/gorilla/websocket"
)

var upgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 本番環境では適切なオリジンチェックを実装
		// 開発環境では全てのオリジンを許可
		return true
	},
}

// WebSocketHandler はWebSocket接続を処理する
type WebSocketHandler struct {
	hub *websocket.Hub
}

// NewWebSocketHandler は新しいWebSocketHandlerを作成する
func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// HandleWebSocket はWebSocket接続をアップグレードして処理する
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// トークンの取得（クエリパラメータまたはヘッダー）
	token := c.Query("token")
	if token == "" {
		// ヘッダーからトークンを取得
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

	// トークンの検証
	claims, err := jwt.VerifyToken(token)
	if err != nil {
		log.Printf("Token verification failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// ユーザーIDをUUIDに変換
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		log.Printf("Invalid user ID in token: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// WebSocket接続へアップグレード
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}

	// クライアントを作成
	client := &websocket.Client{
		// hub, conn, userID は外部からアクセス可能なフィールドとして定義されている必要がある
		// hub.goで定義されているClient構造体を使用
	}

	// クライアント構造体の初期化
	// 注: hub.goのClient構造体のフィールドがエクスポートされていることを前提
	client = h.createClient(conn, userID)

	// Hubにクライアントを登録
	h.hub.Register(client)

	// 接続確立メッセージを送信
	connectionMsg, err := websocket.NewConnectionEstablishedMessage(userID)
	if err != nil {
		log.Printf("Failed to create connection message: %v", err)
	} else {
		if err := h.hub.SendToUser(userID, connectionMsg); err != nil {
			log.Printf("Failed to send connection message: %v", err)
		}
	}

	// ゴルーチンで読み書きを開始
	go client.WritePump()
	go client.ReadPump()

	log.Printf("WebSocket connection established for user: %s", userID)
}

// createClient は新しいクライアントを作成する（ヘルパー関数）
func (h *WebSocketHandler) createClient(conn *gorillaws.Conn, userID uuid.UUID) *websocket.Client {
	return &websocket.Client{
		Hub:    h.hub,
		Conn:   conn,
		UserID: userID,
		Send:   make(chan []byte, 256),
	}
}

// RegisterRoutes はWebSocketルートを登録する
func (h *WebSocketHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/ws", h.HandleWebSocket)
}

// GetStats はWebSocket接続の統計情報を返す（デバッグ用）
func (h *WebSocketHandler) GetStats(c *gin.Context) {
	stats := gin.H{
		"connected_users":   h.hub.GetConnectedUserCount(),
		"total_connections": h.hub.GetTotalConnectionCount(),
	}

	c.JSON(http.StatusOK, stats)
}
