package notification

import (
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/api/websocket"
	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

// Service handles notification operations
type Service struct {
	hub *websocket.Hub
}

// NewService creates a new notification service
func NewService(hub *websocket.Hub) *Service {
	return &Service{
		hub: hub,
	}
}

// NotifyOCRProgress sends OCR progress notification to a user
func (s *Service) NotifyOCRProgress(userID string, data *models.OCRProgressData) error {
	notification := &models.Notification{
		Type:      models.NotificationTypeOCRProgress,
		Data:      data,
		Timestamp: time.Now(),
	}

	s.hub.SendToUser(userID, notification)
	return nil
}

// NotifyTTSProgress sends TTS progress notification to a user
func (s *Service) NotifyTTSProgress(userID string, data *models.TTSProgressData) error {
	notification := &models.Notification{
		Type:      models.NotificationTypeTTSProgress,
		Data:      data,
		Timestamp: time.Now(),
	}

	s.hub.SendToUser(userID, notification)
	return nil
}

// NotifyLearningUpdate sends learning progress update to a user
func (s *Service) NotifyLearningUpdate(userID string, data *models.LearningUpdateData) error {
	notification := &models.Notification{
		Type:      models.NotificationTypeLearningUpdate,
		Data:      data,
		Timestamp: time.Now(),
	}

	s.hub.SendToUser(userID, notification)
	return nil
}

// NotifyError sends error notification to a user
func (s *Service) NotifyError(userID string, data *models.ErrorData) error {
	notification := &models.Notification{
		Type:      models.NotificationTypeError,
		Data:      data,
		Timestamp: time.Now(),
	}

	s.hub.SendToUser(userID, notification)
	return nil
}

// IsUserConnected checks if a user is connected to WebSocket
func (s *Service) IsUserConnected(userID string) bool {
	return s.hub.IsUserConnected(userID)
}

// GetConnectedUsersCount returns the number of connected users
func (s *Service) GetConnectedUsersCount() int {
	return s.hub.GetConnectedUsers()
}
