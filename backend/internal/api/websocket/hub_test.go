package websocket

import (
	"testing"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestNewHub(t *testing.T) {
	hub := NewHub()
	assert.NotNil(t, hub)
	assert.NotNil(t, hub.clients)
	assert.NotNil(t, hub.broadcast)
	assert.NotNil(t, hub.register)
	assert.NotNil(t, hub.unregister)
}

func TestHubRegisterClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create a mock client
	mockConn := &websocket.Conn{} // This will be nil, but sufficient for testing
	client := &Client{
		hub:    hub,
		conn:   mockConn,
		send:   make(chan []byte, 256),
		userID: "test-user-1",
	}

	// Register the client
	hub.register <- client

	// Give it time to process
	time.Sleep(100 * time.Millisecond)

	// Check if user is connected
	assert.True(t, hub.IsUserConnected("test-user-1"))
	assert.Equal(t, 1, hub.GetConnectedUsers())
}

func TestHubUnregisterClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create a mock client
	client := &Client{
		hub:    hub,
		conn:   nil,
		send:   make(chan []byte, 256),
		userID: "test-user-2",
	}

	// Register the client
	hub.register <- client
	time.Sleep(100 * time.Millisecond)

	assert.True(t, hub.IsUserConnected("test-user-2"))

	// Unregister the client
	hub.unregister <- client
	time.Sleep(100 * time.Millisecond)

	assert.False(t, hub.IsUserConnected("test-user-2"))
	assert.Equal(t, 0, hub.GetConnectedUsers())
}

func TestHubSendToUser(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create a mock client with a channel we can read from
	client := &Client{
		hub:    hub,
		conn:   nil,
		send:   make(chan []byte, 256),
		userID: "test-user-3",
	}

	// Register the client
	hub.register <- client
	time.Sleep(100 * time.Millisecond)

	// Send a notification
	notification := &models.Notification{
		Type:      models.NotificationTypePing,
		Data:      nil,
		Timestamp: time.Now(),
	}

	hub.SendToUser("test-user-3", notification)

	// Wait for message to be sent
	time.Sleep(100 * time.Millisecond)

	// Check if message was received
	select {
	case msg := <-client.send:
		assert.NotNil(t, msg)
		assert.Contains(t, string(msg), "ping")
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for message")
	}
}

func TestHubMultipleClients(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create multiple clients for the same user
	client1 := &Client{
		hub:    hub,
		conn:   nil,
		send:   make(chan []byte, 256),
		userID: "test-user-4",
	}

	client2 := &Client{
		hub:    hub,
		conn:   nil,
		send:   make(chan []byte, 256),
		userID: "test-user-4",
	}

	// Register both clients
	hub.register <- client1
	hub.register <- client2
	time.Sleep(100 * time.Millisecond)

	assert.True(t, hub.IsUserConnected("test-user-4"))
	assert.Equal(t, 1, hub.GetConnectedUsers()) // Same user, so count is 1

	// Send notification to user
	notification := &models.Notification{
		Type:      models.NotificationTypePing,
		Data:      nil,
		Timestamp: time.Now(),
	}

	hub.SendToUser("test-user-4", notification)
	time.Sleep(100 * time.Millisecond)

	// Both clients should receive the message
	select {
	case <-client1.send:
		// OK
	case <-time.After(1 * time.Second):
		t.Fatal("Client 1 did not receive message")
	}

	select {
	case <-client2.send:
		// OK
	case <-time.After(1 * time.Second):
		t.Fatal("Client 2 did not receive message")
	}
}
