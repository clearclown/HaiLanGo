package notification

import (
	"testing"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/api/websocket"
	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	hub := websocket.NewHub()
	service := NewService(hub)
	assert.NotNil(t, service)
	assert.NotNil(t, service.hub)
}

func TestNotifyOCRProgress(t *testing.T) {
	hub := websocket.NewHub()
	go hub.Run()

	service := NewService(hub)

	// Create a mock client
	client := websocket.NewClient(hub, nil, "test-user")
	hub.RegisterClient(client)
	time.Sleep(100 * time.Millisecond)

	// Send OCR progress notification
	data := &models.OCRProgressData{
		BookID:         "book-123",
		TotalPages:     100,
		ProcessedPages: 50,
		CurrentPage:    50,
		Progress:       50.0,
		Status:         "processing",
	}

	err := service.NotifyOCRProgress("test-user", data)
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
}

func TestNotifyTTSProgress(t *testing.T) {
	hub := websocket.NewHub()
	go hub.Run()

	service := NewService(hub)

	// Create a mock client
	client := websocket.NewClient(hub, nil, "test-user")
	hub.RegisterClient(client)
	time.Sleep(100 * time.Millisecond)

	// Send TTS progress notification
	data := &models.TTSProgressData{
		BookID:            "book-123",
		PageNumber:        10,
		TotalSegments:     5,
		ProcessedSegments: 3,
		Progress:          60.0,
		Status:            "processing",
	}

	err := service.NotifyTTSProgress("test-user", data)
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
}

func TestNotifyLearningUpdate(t *testing.T) {
	hub := websocket.NewHub()
	go hub.Run()

	service := NewService(hub)

	// Create a mock client
	client := websocket.NewClient(hub, nil, "test-user")
	hub.RegisterClient(client)
	time.Sleep(100 * time.Millisecond)

	// Send learning update notification
	data := &models.LearningUpdateData{
		UserID:         "test-user",
		BookID:         "book-123",
		PageNumber:     20,
		CompletedPages: 20,
		TotalPages:     100,
		LearnedWords:   150,
		StudyTimeMS:    3600000,
	}

	err := service.NotifyLearningUpdate("test-user", data)
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
}

func TestNotifyError(t *testing.T) {
	hub := websocket.NewHub()
	go hub.Run()

	service := NewService(hub)

	// Create a mock client
	client := websocket.NewClient(hub, nil, "test-user")
	hub.RegisterClient(client)
	time.Sleep(100 * time.Millisecond)

	// Send error notification
	data := &models.ErrorData{
		Code:    "OCR_FAILED",
		Message: "Failed to process image",
		Details: "OCR service is temporarily unavailable",
	}

	err := service.NotifyError("test-user", data)
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
}

func TestIsUserConnected(t *testing.T) {
	hub := websocket.NewHub()
	go hub.Run()

	service := NewService(hub)

	// Initially, user should not be connected
	assert.False(t, service.IsUserConnected("test-user"))

	// Register a client
	client := websocket.NewClient(hub, nil, "test-user")
	hub.RegisterClient(client)
	time.Sleep(100 * time.Millisecond)

	// Now user should be connected
	assert.True(t, service.IsUserConnected("test-user"))
}

func TestGetConnectedUsersCount(t *testing.T) {
	hub := websocket.NewHub()
	go hub.Run()

	service := NewService(hub)

	// Initially, no users should be connected
	assert.Equal(t, 0, service.GetConnectedUsersCount())

	// Register multiple clients
	client1 := websocket.NewClient(hub, nil, "user-1")
	client2 := websocket.NewClient(hub, nil, "user-2")
	hub.RegisterClient(client1)
	hub.RegisterClient(client2)
	time.Sleep(100 * time.Millisecond)

	// Should have 2 connected users
	assert.Equal(t, 2, service.GetConnectedUsersCount())
}
