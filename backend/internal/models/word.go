package models

import (
	"time"
)

// Word は単語データを表すモデル
type Word struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"user_id" db:"user_id"`
	BookID         string    `json:"book_id" db:"book_id"`
	PageNumber     int       `json:"page_number" db:"page_number"`
	Text           string    `json:"text" db:"text"`                     // 学習先言語の単語
	Meaning        string    `json:"meaning" db:"meaning"`               // 母国語での意味
	Pronunciation  string    `json:"pronunciation" db:"pronunciation"`   // 発音記号（オプション）
	PartOfSpeech   string    `json:"part_of_speech" db:"part_of_speech"` // 品詞
	Example        string    `json:"example" db:"example"`               // 例文（オプション）
	Language       string    `json:"language" db:"language"`             // 言語コード（例: "ru", "en", "ja"）
	ReviewCount    int       `json:"review_count" db:"review_count"`     // 学習回数
	AverageScore   float64   `json:"average_score" db:"average_score"`   // 平均スコア
	Mastery        float64   `json:"mastery" db:"mastery"`               // 習得度（0-100%）
	Tags           []string  `json:"tags" db:"tags"`                     // タグ（グループ化用）
	LastReviewedAt time.Time `json:"last_reviewed_at" db:"last_reviewed_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// WordFilter は単語検索のフィルタ条件
type WordFilter struct {
	UserID     string   `json:"user_id"`
	BookID     string   `json:"book_id"`
	Language   string   `json:"language"`
	Query      string   `json:"query"`       // 単語のテキスト検索
	Tags       []string `json:"tags"`        // タグでフィルタ
	MinMastery float64  `json:"min_mastery"` // 最小習得度
	MaxMastery float64  `json:"max_mastery"` // 最大習得度
	Limit      int      `json:"limit"`
	Offset     int      `json:"offset"`
	SortBy     string   `json:"sort_by"` // "created_at", "mastery", "review_count"
	SortOrder  string   `json:"sort_order"` // "asc", "desc"
}

// WordStats は単語統計情報
type WordStats struct {
	TotalWords     int     `json:"total_words"`
	MasteredWords  int     `json:"mastered_words"` // 習得度80%以上
	AverageMastery float64 `json:"average_mastery"`
	TotalReviews   int     `json:"total_reviews"`
}
