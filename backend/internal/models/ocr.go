package models

import (
	"time"

	"github.com/google/uuid"
)

// OCRRequest はOCR処理リクエスト
type OCRRequest struct {
	ImageURL string     `json:"image_url" binding:"required"`
	Language string     `json:"language" binding:"required"` // 学習先言語
	Options  OCROptions `json:"options"`
}

// OCROptions はOCR処理オプション
type OCROptions struct {
	DetectOrientation bool     `json:"detect_orientation"` // 向き検出
	DetectLanguage    bool     `json:"detect_language"`    // 言語自動検出
	Languages         []string `json:"languages"`          // 検出対象言語リスト
}

// OCRJobDetail はOCR処理ジョブの詳細情報（APIレスポンス用）
type OCRJobDetail struct {
	ID          string      `json:"id"`
	BookID      string      `json:"book_id"`
	PageNumber  int         `json:"page_number"`
	ImageURL    string      `json:"image_url"`
	Status      OCRStatus   `json:"status"`
	Progress    int         `json:"progress"` // 0-100
	Result      *OCRResult  `json:"result,omitempty"`
	Error       string      `json:"error,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	CompletedAt *time.Time  `json:"completed_at,omitempty"`
}

// OCRResult はOCR処理結果
type OCRResult struct {
	Text             string     `json:"text"`              // 抽出されたテキスト
	DetectedLanguage string     `json:"detected_language"` // 検出された言語コード
	Confidence       float64    `json:"confidence"`        // 全体の信頼度 (0-1)
	Words            []OCRWord  `json:"words"`             // 単語レベルの結果
	Lines            []OCRLine  `json:"lines"`             // 行レベルの結果
	Blocks           []OCRBlock `json:"blocks"`            // ブロックレベルの結果
	HasRuby          bool       `json:"has_ruby"`          // ルビ（ふりがな）の有無
	Orientation      int        `json:"orientation"`       // 画像の向き（度数）
	ProcessingTime   int        `json:"processing_time"`   // 処理時間（ミリ秒）
}

// OCRWord は単語レベルのOCR結果
type OCRWord struct {
	Text        string      `json:"text"`
	Confidence  float64     `json:"confidence"`
	BoundingBox BoundingBox `json:"bounding_box"`
	Language    string      `json:"language,omitempty"`
}

// OCRLine は行レベルのOCR結果
type OCRLine struct {
	Text        string      `json:"text"`
	Confidence  float64     `json:"confidence"`
	BoundingBox BoundingBox `json:"bounding_box"`
	Words       []OCRWord   `json:"words"`
}

// OCRBlock はブロックレベルのOCR結果
type OCRBlock struct {
	Type        string      `json:"type"` // "text", "table", "image", etc.
	Text        string      `json:"text"`
	Confidence  float64     `json:"confidence"`
	BoundingBox BoundingBox `json:"bounding_box"`
	Lines       []OCRLine   `json:"lines"`
}

// BoundingBox は矩形領域
type BoundingBox struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// ProcessPageOCRRequest は特定ページのOCR処理リクエスト
type ProcessPageOCRRequest struct {
	PageNumber int        `json:"page_number" binding:"required,min=1"`
	Language   string     `json:"language" binding:"required"`
	Options    OCROptions `json:"options"`
}

// OCRJobResponse はOCR処理ジョブレスポンス
type OCRJobResponse struct {
	JobID      string    `json:"job_id"`
	BookID     string    `json:"book_id,omitempty"`
	PageNumber int       `json:"page_number,omitempty"`
	Status     OCRStatus `json:"status"`
	Progress   int       `json:"progress"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// OCRResultResponse はOCR処理結果レスポンス
type OCRResultResponse struct {
	JobID       string     `json:"job_id"`
	BookID      string     `json:"book_id,omitempty"`
	PageNumber  int        `json:"page_number,omitempty"`
	Status      OCRStatus  `json:"status"`
	Result      *OCRResult `json:"result,omitempty"`
	Error       string     `json:"error,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// BatchOCRRequest はバッチOCR処理リクエスト
type BatchOCRRequest struct {
	BookID   string     `json:"book_id" binding:"required"`
	Language string     `json:"language" binding:"required"`
	Options  OCROptions `json:"options"`
}

// BatchOCRResponse はバッチOCR処理レスポンス
type BatchOCRResponse struct {
	BookID     string    `json:"book_id"`
	TotalPages int       `json:"total_pages"`
	JobIDs     []string  `json:"job_ids"`
	CreatedAt  time.Time `json:"created_at"`
}

// OCRStatistics はOCR統計情報
type OCRStatistics struct {
	TotalJobs           int     `json:"total_jobs"`
	CompletedJobs       int     `json:"completed_jobs"`
	FailedJobs          int     `json:"failed_jobs"`
	PendingJobs         int     `json:"pending_jobs"`
	ProcessingJobs      int     `json:"processing_jobs"`
	AverageConfidence   float64 `json:"average_confidence"`
	TotalProcessingTime int     `json:"total_processing_time"` // ミリ秒
}

// OCRJobRecord はOCR処理ジョブのデータベースレコード
type OCRJobRecord struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	BookID      uuid.UUID
	PageNumber  int
	ImageURL    string
	Status      OCRStatus
	Progress    int
	ResultJSON  string // JSON化されたOCRResult
	Error       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
}
