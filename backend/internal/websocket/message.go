package websocket

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// MessageType は通知メッセージのタイプ
type MessageType string

const (
	// MessageTypeOCRProgress はOCR処理の進捗通知
	MessageTypeOCRProgress MessageType = "ocr_progress"

	// MessageTypeBookReady は書籍の準備完了通知
	MessageTypeBookReady MessageType = "book_ready"

	// MessageTypeReviewReminder は復習リマインダー通知
	MessageTypeReviewReminder MessageType = "review_reminder"

	// MessageTypeLearningUpdate は学習状態の更新通知
	MessageTypeLearningUpdate MessageType = "learning_update"

	// MessageTypeNotification は一般的な通知
	MessageTypeNotification MessageType = "notification"

	// MessageTypeError はエラー通知
	MessageTypeError MessageType = "error"

	// MessageTypeConnectionEstablished は接続確立通知
	MessageTypeConnectionEstablished MessageType = "connection_established"
)

// Message はWebSocketメッセージの基本構造
type Message struct {
	Type      MessageType     `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp time.Time       `json:"timestamp"`
}

// NewMessage は新しいメッセージを作成する
func NewMessage(msgType MessageType, payload interface{}) (Message, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return Message{}, err
	}

	return Message{
		Type:      msgType,
		Payload:   data,
		Timestamp: time.Now(),
	}, nil
}

// OCRProgressPayload はOCR処理の進捗ペイロード
type OCRProgressPayload struct {
	BookID         uuid.UUID `json:"bookId"`
	TotalPages     int       `json:"totalPages"`
	ProcessedPages int       `json:"processedPages"`
	Progress       float64   `json:"progress"`
	CurrentPage    int       `json:"currentPage,omitempty"`
	Status         string    `json:"status"`
	Message        string    `json:"message,omitempty"`
}

// BookReadyPayload は書籍準備完了のペイロード
type BookReadyPayload struct {
	BookID     uuid.UUID `json:"bookId"`
	Title      string    `json:"title"`
	TotalPages int       `json:"totalPages"`
	Message    string    `json:"message"`
}

// ReviewItem は復習アイテム
type ReviewItem struct {
	ID          uuid.UUID `json:"id"`
	Content     string    `json:"content"`
	Translation string    `json:"translation,omitempty"`
	DueDate     time.Time `json:"dueDate"`
	Priority    string    `json:"priority"` // urgent, recommended, optional
}

// ReviewReminderPayload は復習リマインダーのペイロード
type ReviewReminderPayload struct {
	Count   int          `json:"count"`
	Items   []ReviewItem `json:"items"`
	Message string       `json:"message"`
}

// LearningStats は学習統計
type LearningStats struct {
	TotalTime         int       `json:"totalTime"`         // 秒単位
	CompletedPages    int       `json:"completedPages"`
	MasteredWords     int       `json:"masteredWords"`
	PronunciationScore float64   `json:"pronunciationScore,omitempty"`
	StreakDays        int       `json:"streakDays"`
	LastStudiedAt     time.Time `json:"lastStudiedAt,omitempty"`
}

// LearningUpdatePayload は学習更新のペイロード
type LearningUpdatePayload struct {
	SessionID uuid.UUID     `json:"sessionId"`
	BookID    uuid.UUID     `json:"bookId,omitempty"`
	Stats     LearningStats `json:"stats"`
	Message   string        `json:"message,omitempty"`
}

// NotificationLevel は通知レベル
type NotificationLevel string

const (
	// NotificationLevelInfo は情報レベル
	NotificationLevelInfo NotificationLevel = "info"

	// NotificationLevelSuccess は成功レベル
	NotificationLevelSuccess NotificationLevel = "success"

	// NotificationLevelWarning は警告レベル
	NotificationLevelWarning NotificationLevel = "warning"

	// NotificationLevelError はエラーレベル
	NotificationLevelError NotificationLevel = "error"
)

// NotificationPayload は一般的な通知のペイロード
type NotificationPayload struct {
	Title   string            `json:"title"`
	Message string            `json:"message"`
	Level   NotificationLevel `json:"level"`
	Action  *NotificationAction `json:"action,omitempty"`
}

// NotificationAction は通知に関連するアクション
type NotificationAction struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// ErrorPayload はエラー通知のペイロード
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ConnectionEstablishedPayload は接続確立のペイロード
type ConnectionEstablishedPayload struct {
	UserID    uuid.UUID `json:"userId"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// Helper functions for creating typed messages

// NewOCRProgressMessage はOCR進捗メッセージを作成する
func NewOCRProgressMessage(bookID uuid.UUID, totalPages, processedPages int, status, message string) (Message, error) {
	progress := 0.0
	if totalPages > 0 {
		progress = float64(processedPages) / float64(totalPages) * 100
	}

	payload := OCRProgressPayload{
		BookID:         bookID,
		TotalPages:     totalPages,
		ProcessedPages: processedPages,
		Progress:       progress,
		Status:         status,
		Message:        message,
	}

	return NewMessage(MessageTypeOCRProgress, payload)
}

// NewBookReadyMessage は書籍準備完了メッセージを作成する
func NewBookReadyMessage(bookID uuid.UUID, title string, totalPages int) (Message, error) {
	payload := BookReadyPayload{
		BookID:     bookID,
		Title:      title,
		TotalPages: totalPages,
		Message:    "Book is ready for learning!",
	}

	return NewMessage(MessageTypeBookReady, payload)
}

// NewReviewReminderMessage は復習リマインダーメッセージを作成する
func NewReviewReminderMessage(count int, items []ReviewItem) (Message, error) {
	payload := ReviewReminderPayload{
		Count:   count,
		Items:   items,
		Message: "You have items to review",
	}

	return NewMessage(MessageTypeReviewReminder, payload)
}

// NewLearningUpdateMessage は学習更新メッセージを作成する
func NewLearningUpdateMessage(sessionID uuid.UUID, stats LearningStats) (Message, error) {
	payload := LearningUpdatePayload{
		SessionID: sessionID,
		Stats:     stats,
	}

	return NewMessage(MessageTypeLearningUpdate, payload)
}

// NewNotificationMessage は一般通知メッセージを作成する
func NewNotificationMessage(title, message string, level NotificationLevel) (Message, error) {
	payload := NotificationPayload{
		Title:   title,
		Message: message,
		Level:   level,
	}

	return NewMessage(MessageTypeNotification, payload)
}

// NewErrorMessage はエラーメッセージを作成する
func NewErrorMessage(code, message, details string) (Message, error) {
	payload := ErrorPayload{
		Code:    code,
		Message: message,
		Details: details,
	}

	return NewMessage(MessageTypeError, payload)
}

// NewConnectionEstablishedMessage は接続確立メッセージを作成する
func NewConnectionEstablishedMessage(userID uuid.UUID) (Message, error) {
	payload := ConnectionEstablishedPayload{
		UserID:    userID,
		Message:   "WebSocket connection established",
		Timestamp: time.Now(),
	}

	return NewMessage(MessageTypeConnectionEstablished, payload)
}
