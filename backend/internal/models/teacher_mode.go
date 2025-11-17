package models

import (
	"time"

	"github.com/google/uuid"
)

// TeacherModeSettings は教師モードの設定を保持する
type TeacherModeSettings struct {
	Speed                float64               `json:"speed"`                  // 再生速度 (0.5, 0.75, 1.0, 1.25, 1.5, 2.0)
	PageInterval         int                   `json:"page_interval"`          // ページ間隔（秒）
	RepeatCount          int                   `json:"repeat_count"`           // リピート回数 (1, 2, 3)
	AudioQuality         string                `json:"audio_quality"`          // 音質 ("standard", "premium")
	Content              TeacherModeContent    `json:"content"`                // 学習内容設定
}

// TeacherModeContent は教師モードの学習内容設定を保持する
type TeacherModeContent struct {
	IncludeTranslation           bool `json:"include_translation"`            // 母国語訳を含む
	IncludeWordExplanation       bool `json:"include_word_explanation"`       // 単語解説を含む
	IncludeGrammarExplanation    bool `json:"include_grammar_explanation"`    // 文法解説を含む
	IncludePronunciationPractice bool `json:"include_pronunciation_practice"` // 発音練習を含む
	IncludeExampleSentences      bool `json:"include_example_sentences"`      // 例文を含む
}

// AudioSegmentType は音声セグメントのタイプ
type AudioSegmentType string

const (
	AudioSegmentTypePhrase      AudioSegmentType = "phrase"      // フレーズ
	AudioSegmentTypeTranslation AudioSegmentType = "translation" // 翻訳
	AudioSegmentTypeExplanation AudioSegmentType = "explanation" // 解説
	AudioSegmentTypePause       AudioSegmentType = "pause"       // 一時停止
)

// AudioSegment は音声セグメントを表す
type AudioSegment struct {
	ID       string           `json:"id"`        // セグメントID
	Type     AudioSegmentType `json:"type"`      // タイプ
	AudioURL string           `json:"audio_url"` // 音声URL
	Duration int              `json:"duration"`  // 長さ（ミリ秒）
	Text     string           `json:"text"`      // テキスト
	Language string           `json:"language"`  // 言語
}

// PageAudio はページの音声情報を表す
type PageAudio struct {
	PageNumber    int            `json:"page_number"`    // ページ番号
	Segments      []AudioSegment `json:"segments"`       // セグメント
	TotalDuration int            `json:"total_duration"` // 総長さ（ミリ秒）
}

// TeacherModePlaylist は教師モードのプレイリストを表す
type TeacherModePlaylist struct {
	ID            string              `json:"id"`             // プレイリストID
	BookID        uuid.UUID           `json:"book_id"`        // 書籍ID
	Pages         []PageAudio         `json:"pages"`          // ページ
	Settings      TeacherModeSettings `json:"settings"`       // 設定
	TotalDuration int                 `json:"total_duration"` // 総長さ（ミリ秒）
}

// PlaybackStatus は再生状態
type PlaybackStatus string

const (
	PlaybackStatusStopped PlaybackStatus = "stopped" // 停止
	PlaybackStatusPlaying PlaybackStatus = "playing" // 再生中
	PlaybackStatusPaused  PlaybackStatus = "paused"  // 一時停止
)

// PlaybackState は再生状態を表す
type PlaybackState struct {
	Status               PlaybackStatus `json:"status"`                  // 状態
	CurrentPage          int            `json:"current_page"`            // 現在のページ
	CurrentSegmentIndex  int            `json:"current_segment_index"`   // 現在のセグメントインデックス
	ElapsedTime          int            `json:"elapsed_time"`            // 経過時間（ミリ秒）
	TotalDuration        int            `json:"total_duration"`          // 総長さ（ミリ秒）
}

// TeacherModeDownload は教師モードのダウンロード履歴を表す
type TeacherModeDownload struct {
	ID             uuid.UUID           `json:"id" db:"id"`
	UserID         uuid.UUID           `json:"user_id" db:"user_id"`
	BookID         uuid.UUID           `json:"book_id" db:"book_id"`
	Settings       TeacherModeSettings `json:"settings" db:"settings"`           // JSONBとして保存
	TotalSizeBytes int64               `json:"total_size_bytes" db:"total_size_bytes"`
	DownloadedAt   time.Time           `json:"downloaded_at" db:"downloaded_at"`
	ExpiresAt      *time.Time          `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt      time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at" db:"updated_at"`
}

// TeacherModePlaybackHistory は教師モードの再生履歴を表す
type TeacherModePlaybackHistory struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	UserID               uuid.UUID `json:"user_id" db:"user_id"`
	BookID               uuid.UUID `json:"book_id" db:"book_id"`
	CurrentPage          int       `json:"current_page" db:"current_page"`
	CurrentSegmentIndex  int       `json:"current_segment_index" db:"current_segment_index"`
	ElapsedTime          int       `json:"elapsed_time" db:"elapsed_time"` // 秒単位
	TotalPlayTimeSeconds int       `json:"total_play_time_seconds" db:"total_play_time_seconds"`
	LastPlayedAt         time.Time `json:"last_played_at" db:"last_played_at"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// PageRange はページ範囲を表す
type PageRange struct {
	Start int `json:"start"` // 開始ページ
	End   int `json:"end"`   // 終了ページ
}
