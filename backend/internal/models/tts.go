package models

import (
	"time"

	"github.com/google/uuid"
)

// TTSRequest はTTS音声合成リクエスト
type TTSRequest struct {
	Text     string              `json:"text" binding:"required"`
	Language string              `json:"language" binding:"required"` // 言語コード（例: "ja", "en", "ru"）
	Options  TTSSynthesizeOptions `json:"options"`
}

// TTSSynthesizeOptions はTTS音声合成オプション
type TTSSynthesizeOptions struct {
	Speed       float64       `json:"speed"`        // 再生速度（0.5-2.0、デフォルト: 1.0）
	Quality     TTSQuality    `json:"quality"`      // 音質（standard/premium）
	Voice       string        `json:"voice"`        // 声のID（オプション）
	Format      TTSAudioFormat `json:"format"`      // 音声フォーマット（mp3/wav/ogg）
	CacheEnable bool          `json:"cache_enable"` // キャッシュを使用するか（デフォルト: true）
}

// TTSQuality はTTS音質レベル
type TTSQuality string

const (
	TTSQualityStandard TTSQuality = "standard" // 標準品質（無料プラン）
	TTSQualityPremium  TTSQuality = "premium"  // 高品質（プレミアムプラン）
)

// TTSAudioFormat は音声フォーマット
type TTSAudioFormat string

const (
	TTSAudioFormatMP3 TTSAudioFormat = "mp3"
	TTSAudioFormatWAV TTSAudioFormat = "wav"
	TTSAudioFormatOGG TTSAudioFormat = "ogg"
)

// TTSAudioResponse はTTS音声レスポンス
type TTSAudioResponse struct {
	AudioID     string         `json:"audio_id"`
	AudioURL    string         `json:"audio_url"`
	Duration    int            `json:"duration"`     // 音声の長さ（秒）
	Format      TTSAudioFormat `json:"format"`
	Language    string         `json:"language"`
	Quality     TTSQuality     `json:"quality"`
	Speed       float64        `json:"speed"`
	CachedAt    *time.Time     `json:"cached_at,omitempty"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
	ProcessingTime int         `json:"processing_time"` // 処理時間（ミリ秒）
}

// TTSLanguage はサポート言語情報
type TTSLanguage struct {
	Code        string   `json:"code"`         // 言語コード（例: "ja"）
	Name        string   `json:"name"`         // 言語名（例: "Japanese"）
	NativeName  string   `json:"native_name"`  // ネイティブ名（例: "日本語"）
	Voices      []string `json:"voices"`       // 利用可能な声のID
	IsSupported bool     `json:"is_supported"` // サポートされているか
}

// TTSBatchRequest はバッチTTS合成リクエスト
type TTSBatchRequest struct {
	BookID   string              `json:"book_id" binding:"required"`
	Language string              `json:"language" binding:"required"`
	Options  TTSSynthesizeOptions `json:"options"`
}

// TTSBatchResponse はバッチTTS合成レスポンス
type TTSBatchResponse struct {
	BatchID    string    `json:"batch_id"`
	BookID     string    `json:"book_id"`
	TotalPages int       `json:"total_pages"`
	JobIDs     []string  `json:"job_ids"`
	CreatedAt  time.Time `json:"created_at"`
}

// TTSJobDetail はTTS合成ジョブの詳細情報
type TTSJobDetail struct {
	ID          string       `json:"id"`
	BookID      string       `json:"book_id"`
	PageNumber  int          `json:"page_number"`
	Text        string       `json:"text"`
	Language    string       `json:"language"`
	Status      TTSStatus    `json:"status"`
	Progress    int          `json:"progress"` // 0-100
	AudioID     string       `json:"audio_id,omitempty"`
	AudioURL    string       `json:"audio_url,omitempty"`
	Error       string       `json:"error,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
}

// TTSStatus はTTS処理ステータス
type TTSStatus string

const (
	TTSStatusPending    TTSStatus = "pending"
	TTSStatusProcessing TTSStatus = "processing"
	TTSStatusCompleted  TTSStatus = "completed"
	TTSStatusFailed     TTSStatus = "failed"
)

// TTSCacheStats はTTSキャッシュ統計情報
type TTSCacheStats struct {
	TotalCached     int     `json:"total_cached"`      // キャッシュされた音声の総数
	CacheHitRate    float64 `json:"cache_hit_rate"`    // キャッシュヒット率（0-1）
	TotalSize       int64   `json:"total_size"`        // 総サイズ（バイト）
	AvgDuration     int     `json:"avg_duration"`      // 平均音声長（秒）
	Languages       map[string]int `json:"languages"` // 言語ごとのキャッシュ数
	OldestCachedAt  *time.Time `json:"oldest_cached_at,omitempty"`
	LatestCachedAt  *time.Time `json:"latest_cached_at,omitempty"`
}

// TTSJobResponse はTTS処理ジョブレスポンス
type TTSJobResponse struct {
	JobID      string    `json:"job_id"`
	BookID     string    `json:"book_id,omitempty"`
	PageNumber int       `json:"page_number,omitempty"`
	Status     TTSStatus `json:"status"`
	Progress   int       `json:"progress"`
	AudioID    string    `json:"audio_id,omitempty"`
	AudioURL   string    `json:"audio_url,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TTSJobRecord はTTS処理ジョブのデータベースレコード
type TTSJobRecord struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	BookID      uuid.UUID
	PageNumber  int
	Text        string
	Language    string
	Status      TTSStatus
	Progress    int
	AudioID     string
	AudioURL    string
	Error       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
}
