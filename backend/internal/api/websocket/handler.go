package websocket

import (
	"log"
	"net/http"
	"strings"

	"github.com/clearclown/HaiLanGo/backend/internal/api/middleware"
	"github.com/clearclown/HaiLanGo/backend/pkg/jwt"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if middleware.IsOriginAllowed(origin) {
			return true
		}
		log.Printf("WebSocket origin rejected (legacy handler): origin=%q host=%q path=%q", origin, r.Host, r.URL.Path)
		return false
	},
}

// Handler handles WebSocket connections
type Handler struct {
	hub *Hub
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

// ServeWS handles websocket requests from the peer.
func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	// Extract JWT token from query parameter or Authorization header
	token := r.URL.Query().Get("token")
	if token == "" {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}
	}

	if token == "" {
		http.Error(w, "token is required", http.StatusUnauthorized)
		return
	}

	// Verify the JWT token
	claims, err := jwt.VerifyToken(token)
	if err != nil {
		log.Printf("WebSocket token verification failed: %v", err)
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID
	if userID == "" {
		http.Error(w, "invalid token: missing user ID", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(h.hub, conn, userID)
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// GetHub returns the hub
func (h *Handler) GetHub() *Hub {
	return h.hub
}
