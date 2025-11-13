package models

import "time"

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
