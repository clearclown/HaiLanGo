package models

import "time"

// PageLearning はページ学習データ
type PageLearning struct {
	Page       PageWithOCR          `json:"page"`
	Progress   PageProgressDetail   `json:"progress"`
	Phrases    []Phrase             `json:"phrases"`
	Vocabulary []VocabularyItem     `json:"vocabulary"`
	Navigation NavigationInfo       `json:"navigation"`
}

// PageWithOCR はOCRデータを含むページ
type PageWithOCR struct {
	ID          string `json:"id"`
	BookID      string `json:"book_id"`
	PageNumber  int    `json:"page_number"`
	ImageURL    string `json:"image_url"`
	OCRText     string `json:"ocr_text"`
	Translation string `json:"translation"`
	Language    string `json:"language"`
	HasAudio    bool   `json:"has_audio"`
	AudioURL    string `json:"audio_url,omitempty"`
}

// PageProgressDetail はページ進捗の詳細
type PageProgressDetail struct {
	IsCompleted   bool       `json:"is_completed"`
	CompletedAt   *time.Time `json:"completed_at"`
	StudyTime     int        `json:"study_time"` // 秒
	ReviewCount   int        `json:"review_count"`
	LastStudiedAt *time.Time `json:"last_studied_at"`
}

// Phrase はフレーズ
type Phrase struct {
	ID            string `json:"id"`
	Text          string `json:"text"`
	Translation   string `json:"translation"`
	Pronunciation string `json:"pronunciation,omitempty"`
	AudioURL      string `json:"audio_url,omitempty"`
}

// VocabularyItem は単語
type VocabularyItem struct {
	Word         string `json:"word"`
	Translation  string `json:"translation"`
	PartOfSpeech string `json:"part_of_speech,omitempty"`
	Frequency    string `json:"frequency,omitempty"`
}

// NavigationInfo はナビゲーション情報
type NavigationInfo struct {
	HasPrevious bool `json:"has_previous"`
	HasNext     bool `json:"has_next"`
	TotalPages  int  `json:"total_pages"`
	CurrentPage int  `json:"current_page"`
}

// CompletePageRequest はページ完了リクエスト
type CompletePageRequest struct {
	StudyTime int    `json:"study_time" binding:"required,min=1"`
	Notes     string `json:"notes"`
}

// CompletePageResponse はページ完了レスポンス
type CompletePageResponse struct {
	Message  string             `json:"message"`
	Progress PageProgressDetail `json:"progress"`
	NextPage int                `json:"next_page"`
}

// SessionRequest は学習セッションリクエスト
type SessionRequest struct {
	Action    string    `json:"action" binding:"required,oneof=start end"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
}

// SessionResponse は学習セッションレスポンス
type SessionResponse struct {
	SessionID string     `json:"session_id"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
}

// BookProgressSummary は書籍進捗サマリー
type BookProgressSummary struct {
	BookID               string                  `json:"book_id"`
	TotalPages           int                     `json:"total_pages"`
	CompletedPages       int                     `json:"completed_pages"`
	CompletionPercentage float64                 `json:"completion_percentage"`
	TotalStudyTime       int                     `json:"total_study_time"`
	AverageTimePerPage   float64                 `json:"average_time_per_page"`
	CurrentPage          int                     `json:"current_page"`
	LastStudiedAt        *time.Time              `json:"last_studied_at"`
	Pages                []PageProgressSummaryItem `json:"pages"`
}

// PageProgressSummaryItem はページ進捗サマリー項目
type PageProgressSummaryItem struct {
	PageNumber  int  `json:"page_number"`
	IsCompleted bool `json:"is_completed"`
	StudyTime   int  `json:"study_time"`
	ReviewCount int  `json:"review_count"`
}

// PageProgressRecord はデータベース用のページ進捗レコード
type PageProgressRecord struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	BookID        string     `json:"book_id"`
	PageNumber    int        `json:"page_number"`
	IsCompleted   bool       `json:"is_completed"`
	CompletedAt   *time.Time `json:"completed_at"`
	StudyTime     int        `json:"study_time"` // 秒
	ReviewCount   int        `json:"review_count"`
	LastStudiedAt *time.Time `json:"last_studied_at"`
	Notes         string     `json:"notes,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
