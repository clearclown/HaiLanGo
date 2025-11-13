package websocket

import (
	"sync"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients mapped by user ID
	clients map[string]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *BroadcastMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Mutex to protect clients map
	mu sync.RWMutex
}

// BroadcastMessage represents a message to be broadcast
type BroadcastMessage struct {
	UserID       string
	Notification *models.Notification
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[*Client]bool),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; !ok {
				h.clients[client.userID] = make(map[*Client]bool)
			}
			h.clients[client.userID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.userID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.userID)
					}
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			if clients, ok := h.clients[message.UserID]; ok {
				for client := range clients {
					client.SendNotification(message.Notification)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// SendToUser sends a notification to a specific user
func (h *Hub) SendToUser(userID string, notification *models.Notification) {
	h.broadcast <- &BroadcastMessage{
		UserID:       userID,
		Notification: notification,
	}
}

// GetConnectedUsers returns the number of connected users
func (h *Hub) GetConnectedUsers() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// IsUserConnected checks if a user is connected
func (h *Hub) IsUserConnected(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients, ok := h.clients[userID]
	return ok && len(clients) > 0
}

// RegisterClient registers a client to the hub (for testing)
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient unregisters a client from the hub (for testing)
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}
