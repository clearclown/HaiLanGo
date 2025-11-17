package models

import (
	"time"

	"github.com/google/uuid"
)

// STTResult は音声認識の結果を表す
type STTResult struct {
	Text       string    `json:"text"`        // 認識されたテキスト
	Language   string    `json:"language"`    // 認識された言語コード
	Confidence float64   `json:"confidence"`  // 認識の信頼度（0.0-1.0）
	Duration   float64   `json:"duration"`    // 音声の長さ（秒）
	Words      []WordInfo `json:"words"`      // 単語レベルの情報
	CreatedAt  time.Time `json:"created_at"`
}

// WordInfo は単語レベルの情報を表す
type WordInfo struct {
	Word       string  `json:"word"`       // 単語
	StartTime  float64 `json:"start_time"` // 開始時間（秒）
	EndTime    float64 `json:"end_time"`   // 終了時間（秒）
	Confidence float64 `json:"confidence"` // 信頼度（0.0-1.0）
}

// PronunciationScore は発音評価のスコアを表す
type PronunciationScore struct {
	TotalScore      int                `json:"total_score"`       // 総合スコア（0-100）
	AccuracyScore   int                `json:"accuracy_score"`    // 正確性スコア（0-100）
	FluencyScore    int                `json:"fluency_score"`     // 流暢性スコア（0-100）
	PronuncScore    int                `json:"pronunc_score"`     // 発音スコア（0-100）
	WordScores      []WordScore        `json:"word_scores"`       // 単語ごとのスコア
	Feedback        *Feedback          `json:"feedback"`          // フィードバック
	ExpectedText    string             `json:"expected_text"`     // 期待されるテキスト
	RecognizedText  string             `json:"recognized_text"`   // 認識されたテキスト
	EvaluationID    string             `json:"evaluation_id"`     // 評価ID
	UserID          string             `json:"user_id"`           // ユーザーID
	CreatedAt       time.Time          `json:"created_at"`
}

// WordScore は単語ごとの評価スコアを表す
type WordScore struct {
	Word           string  `json:"word"`            // 単語
	Score          int     `json:"score"`           // スコア（0-100）
	IsCorrect      bool    `json:"is_correct"`      // 正しく発音されたか
	ExpectedWord   string  `json:"expected_word"`   // 期待される単語
	RecognizedWord string  `json:"recognized_word"` // 認識された単語
}

// Feedback は発音評価のフィードバックを表す
type Feedback struct {
	Level           string   `json:"level"`            // レベル（excellent, good, fair, poor）
	Message         string   `json:"message"`          // メッセージ
	PositivePoints  []string `json:"positive_points"`  // 良かった点
	Improvements    []string `json:"improvements"`     // 改善ポイント
	SpecificAdvice  []string `json:"specific_advice"`  // 具体的なアドバイス
}

// AudioProcessingResult は音声処理の結果を表す
type AudioProcessingResult struct {
	ProcessedAudio []byte  `json:"processed_audio"`  // 処理後の音声データ
	SampleRate     int     `json:"sample_rate"`      // サンプリングレート
	Channels       int     `json:"channels"`         // チャンネル数
	Duration       float64 `json:"duration"`         // 長さ（秒）
	NoiseLevel     float64 `json:"noise_level"`      // ノイズレベル
	IsLowQuality   bool    `json:"is_low_quality"`   // 低品質かどうか
}

// STTRequest はSTT音声認識リクエスト
type STTRequest struct {
	AudioData     string          `json:"audio_data"`             // 音声データ（Base64エンコード）
	AudioURL      string          `json:"audio_url"`              // 音声URL（代替）
	Language      string          `json:"language" binding:"required"` // 言語コード（例: "ja", "en", "ru"）
	ReferenceText string          `json:"reference_text"`         // 参照テキスト（発音評価用）
	Options       STTRecognizeOptions `json:"options"`
}

// STTRecognizeOptions はSTT音声認識オプション
type STTRecognizeOptions struct {
	Format            string `json:"format"`             // 音声フォーマット（mp3/wav/ogg）
	EnablePunctuation bool   `json:"enable_punctuation"` // 句読点を有効にするか
	EnableWordTiming  bool   `json:"enable_word_timing"` // 単語タイミングを取得するか
	Evaluate          bool   `json:"evaluate"`           // 発音評価を行うか
}

// STTJobDetail はSTT処理ジョブの詳細情報
type STTJobDetail struct {
	ID             string       `json:"id"`
	BookID         string       `json:"book_id,omitempty"`
	PageNumber     int          `json:"page_number,omitempty"`
	AudioURL       string       `json:"audio_url"`
	ReferenceText  string       `json:"reference_text,omitempty"`
	Language       string       `json:"language"`
	Status         STTStatus    `json:"status"`
	Progress       int          `json:"progress"` // 0-100
	Result         *STTResult   `json:"result,omitempty"`
	Score          *PronunciationScore `json:"score,omitempty"`
	Error          string       `json:"error,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	CompletedAt    *time.Time   `json:"completed_at,omitempty"`
}

// STTStatus はSTT処理ステータス
type STTStatus string

const (
	STTStatusPending    STTStatus = "pending"
	STTStatusProcessing STTStatus = "processing"
	STTStatusCompleted  STTStatus = "completed"
	STTStatusFailed     STTStatus = "failed"
)

// STTLanguage はサポート言語情報
type STTLanguage struct {
	Code                  string `json:"code"`                    // 言語コード（例: "ja"）
	Name                  string `json:"name"`                    // 言語名（例: "Japanese"）
	NativeName            string `json:"native_name"`             // ネイティブ名（例: "日本語"）
	IsSupported           bool   `json:"is_supported"`            // サポートされているか
	SupportsPronunciation bool   `json:"supports_pronunciation"`  // 発音評価サポート
}

// STTJobResponse はSTT処理ジョブレスポンス
type STTJobResponse struct {
	JobID      string              `json:"job_id"`
	BookID     string              `json:"book_id,omitempty"`
	PageNumber int                 `json:"page_number,omitempty"`
	Status     STTStatus           `json:"status"`
	Progress   int                 `json:"progress"`
	Result     *STTResult          `json:"result,omitempty"`
	Score      *PronunciationScore `json:"score,omitempty"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

// STTJobRecord はSTT処理ジョブのデータベースレコード
type STTJobRecord struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	BookID        uuid.UUID
	PageNumber    int
	AudioURL      string
	ReferenceText string
	Language      string
	Status        STTStatus
	Progress      int
	ResultText    string
	ResultScore   int
	Error         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CompletedAt   *time.Time
}

// STTStatistics はSTT統計情報
type STTStatistics struct {
	TotalRecognitions  int                `json:"total_recognitions"`   // 総認識数
	TotalEvaluations   int                `json:"total_evaluations"`    // 総発音評価数
	AverageScore       float64            `json:"average_score"`        // 平均スコア
	BestScore          int                `json:"best_score"`           // 最高スコア
	LanguageStats      map[string]int     `json:"language_stats"`       // 言語ごとの統計
	RecentJobs         []*STTJobDetail    `json:"recent_jobs"`          // 最近のジョブ
	TotalDuration      int                `json:"total_duration"`       // 総音声時間（秒）
}
