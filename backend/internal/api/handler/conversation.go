package handler

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/backend/pkg/realtime"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// ConversationHandler handles real-time voice conversation endpoints
type ConversationHandler struct {
	upgrader    websocket.Upgrader
	sessions    map[string]*ConversationSession
	sessionsMu  sync.RWMutex
}

// ConversationSession represents an active conversation session
type ConversationSession struct {
	ID              string
	UserID          string
	TargetLanguage  string
	NativeLanguage  string
	RealtimeSession realtime.Session
	WSConn          *websocket.Conn
	StartedAt       time.Time
	done            chan struct{}
}

// NewConversationHandler creates a new conversation handler
func NewConversationHandler() *ConversationHandler {
	return &ConversationHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins in development
				// TODO: Restrict in production
				return true
			},
		},
		sessions: make(map[string]*ConversationSession),
	}
}

// RegisterRoutes registers conversation API routes
func (h *ConversationHandler) RegisterRoutes(r *gin.RouterGroup) {
	conversation := r.Group("/conversation")
	{
		conversation.GET("/ws", h.HandleWebSocket)
		conversation.POST("/start", h.StartSession)
		conversation.POST("/stop", h.StopSession)
		conversation.GET("/status", h.GetSessionStatus)
	}
}

// StartSessionRequest represents a request to start a conversation session
type StartSessionRequest struct {
	TargetLanguage string `json:"target_language" binding:"required"`
	NativeLanguage string `json:"native_language" binding:"required"`
}

// StartSessionResponse represents the response for starting a session
type StartSessionResponse struct {
	SessionID string `json:"session_id"`
	WebSocketURL string `json:"websocket_url"`
}

// StartSession starts a new conversation session
func (h *ConversationHandler) StartSession(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req StartSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Create realtime client
	client, err := realtime.NewClientWithLanguage(req.TargetLanguage, req.NativeLanguage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create realtime client"})
		return
	}

	// Configure session for language learning
	config := realtime.LanguageLearningConfig(req.TargetLanguage, req.NativeLanguage)

	// Connect to realtime API
	session, err := client.Connect(c.Request.Context(), config)
	if err != nil {
		client.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to realtime API"})
		return
	}

	// Create conversation session
	convSession := &ConversationSession{
		ID:              session.ID(),
		UserID:          userID.(string),
		TargetLanguage:  req.TargetLanguage,
		NativeLanguage:  req.NativeLanguage,
		RealtimeSession: session,
		StartedAt:       time.Now(),
		done:            make(chan struct{}),
	}

	h.sessionsMu.Lock()
	h.sessions[convSession.ID] = convSession
	h.sessionsMu.Unlock()

	// WebSocket URL for this session
	wsURL := "/api/v1/conversation/ws?session_id=" + convSession.ID

	c.JSON(http.StatusOK, StartSessionResponse{
		SessionID:    convSession.ID,
		WebSocketURL: wsURL,
	})
}

// StopSession stops an active conversation session
func (h *ConversationHandler) StopSession(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sessionID := c.Query("session_id")
	if sessionID == "" {
		var req struct {
			SessionID string `json:"session_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
			return
		}
		sessionID = req.SessionID
	}

	h.sessionsMu.Lock()
	session, exists := h.sessions[sessionID]
	if !exists {
		h.sessionsMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	if session.UserID != userID.(string) {
		h.sessionsMu.Unlock()
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	delete(h.sessions, sessionID)
	h.sessionsMu.Unlock()

	// Close session
	close(session.done)
	if session.RealtimeSession != nil {
		session.RealtimeSession.Close()
	}
	if session.WSConn != nil {
		session.WSConn.Close()
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session stopped"})
}

// GetSessionStatus returns the status of an active session
func (h *ConversationHandler) GetSessionStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	h.sessionsMu.RLock()
	session, exists := h.sessions[sessionID]
	h.sessionsMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	if session.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":      session.ID,
		"target_language": session.TargetLanguage,
		"native_language": session.NativeLanguage,
		"started_at":      session.StartedAt,
		"duration_seconds": time.Since(session.StartedAt).Seconds(),
		"active":          session.WSConn != nil,
	})
}

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type    string          `json:"type"`
	Audio   string          `json:"audio,omitempty"`   // Base64 encoded audio
	Text    string          `json:"text,omitempty"`
	Error   string          `json:"error,omitempty"`
	Event   json.RawMessage `json:"event,omitempty"`
}

// HandleWebSocket handles WebSocket connections for real-time conversation
func (h *ConversationHandler) HandleWebSocket(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	h.sessionsMu.RLock()
	session, exists := h.sessions[sessionID]
	h.sessionsMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	// Upgrade to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	session.WSConn = conn

	// Create channels for coordination
	done := make(chan struct{})
	defer close(done)

	// Forward events from realtime API to WebSocket client
	go h.forwardEvents(session, done)

	// Handle incoming messages from WebSocket client
	for {
		select {
		case <-session.done:
			return
		case <-done:
			return
		default:
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			return
		}

		var msg WebSocketMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			h.sendError(conn, "Invalid message format")
			continue
		}

		switch msg.Type {
		case "audio":
			// Decode and forward audio to realtime API
			audio, err := base64.StdEncoding.DecodeString(msg.Audio)
			if err != nil {
				h.sendError(conn, "Invalid audio encoding")
				continue
			}
			if err := session.RealtimeSession.SendAudio(c.Request.Context(), audio); err != nil {
				h.sendError(conn, "Failed to send audio")
			}

		case "text":
			// Forward text to realtime API
			if err := session.RealtimeSession.SendText(c.Request.Context(), msg.Text); err != nil {
				h.sendError(conn, "Failed to send text")
			}

		case "commit":
			// Commit audio buffer (trigger processing)
			// The realtime session handles buffer commits internally via SendAudio
			// No additional action needed

		default:
			h.sendError(conn, "Unknown message type: "+msg.Type)
		}
	}
}

// forwardEvents forwards events from the realtime API to the WebSocket client
func (h *ConversationHandler) forwardEvents(session *ConversationSession, done <-chan struct{}) {
	events := session.RealtimeSession.Receive()

	for {
		select {
		case <-done:
			return
		case <-session.done:
			return
		case event, ok := <-events:
			if !ok {
				return
			}

			msg := WebSocketMessage{
				Type: event.Type,
			}

			// Handle different event types
			switch event.Type {
			case realtime.EventResponseAudioDelta:
				if event.Delta != nil && len(event.Delta.Audio) > 0 {
					msg.Audio = base64.StdEncoding.EncodeToString(event.Delta.Audio)
				}
			case realtime.EventResponseTextDelta, realtime.EventResponseAudioTranscriptDelta:
				if event.Delta != nil {
					msg.Text = event.Delta.Text
				}
			case realtime.EventError:
				if event.Error != nil {
					msg.Error = event.Error.Message
				}
			}

			if session.WSConn != nil {
				if err := session.WSConn.WriteJSON(msg); err != nil {
					log.Printf("WebSocket write error: %v", err)
					return
				}
			}
		}
	}
}

// sendError sends an error message over WebSocket
func (h *ConversationHandler) sendError(conn *websocket.Conn, message string) {
	msg := WebSocketMessage{
		Type:  "error",
		Error: message,
	}
	conn.WriteJSON(msg)
}

// CleanupSessions removes stale sessions
func (h *ConversationHandler) CleanupSessions(maxAge time.Duration) {
	h.sessionsMu.Lock()
	defer h.sessionsMu.Unlock()

	now := time.Now()
	for id, session := range h.sessions {
		if now.Sub(session.StartedAt) > maxAge {
			close(session.done)
			if session.RealtimeSession != nil {
				session.RealtimeSession.Close()
			}
			if session.WSConn != nil {
				session.WSConn.Close()
			}
			delete(h.sessions, id)
		}
	}
}
