package stt

// エラーメッセージ定数
const (
	ErrEmptyAudioData     = "音声データが空です"
	ErrEmptyExpectedText  = "期待されるテキストが空です"
	ErrInvalidLanguage    = "無効な言語コード"
	ErrSTTRecognition     = "音声認識に失敗しました"
	ErrAudioProcessing    = "音声処理に失敗しました"
)

// 評価スコアの境界値
const (
	ScoreExcellentThreshold = 90 // 優秀
	ScoreGoodThreshold      = 75 // 良好
	ScoreFairThreshold      = 45 // 改善が必要
	// それ以下は poor (要練習)
)

// スコアの重み付け（総合スコア計算用）
const (
	AccuracyWeight      = 40 // 正確性の重み（%）
	FluencyWeight       = 30 // 流暢性の重み（%）
	PronunciationWeight = 30 // 発音の重み（%）
)

// 音声処理の定数
const (
	DefaultSampleRate     = 16000 // 推奨サンプリングレート（Hz）
	MinAudioDataSize      = 10    // 最小音声データサイズ（バイト）
	NoiseThresholdHigh    = 0.3   // 高ノイズレベルの閾値
	IdealTimePerWord      = 0.5   // 理想的な単語あたりの時間（秒）
	WordScoreCorrectMin   = 80    // 単語が正しいと判断される最小スコア
)

// フィードバックレベル
const (
	FeedbackLevelExcellent = "excellent"
	FeedbackLevelGood      = "good"
	FeedbackLevelFair      = "fair"
	FeedbackLevelPoor      = "poor"
)
