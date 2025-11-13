// Package teachermode provides teacher mode functionality
package teachermode

import (
	"context"
	"errors"
	"fmt"
)

// TeacherModeSettings 教師モード設定
type TeacherModeSettings struct {
	Speed        float64            `json:"speed"`
	PageInterval int                `json:"pageInterval"`
	RepeatCount  int                `json:"repeatCount"`
	AudioQuality string             `json:"audioQuality"`
	Content      TeacherModeContent `json:"content"`
}

// TeacherModeContent 学習内容設定
type TeacherModeContent struct {
	IncludeTranslation           bool `json:"includeTranslation"`
	IncludeWordExplanation       bool `json:"includeWordExplanation"`
	IncludeGrammarExplanation    bool `json:"includeGrammarExplanation"`
	IncludePronunciationPractice bool `json:"includePronunciationPractice"`
	IncludeExampleSentences      bool `json:"includeExampleSentences"`
}

// AudioSegment 音声セグメント
type AudioSegment struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	AudioURL string `json:"audioUrl"`
	Duration int64  `json:"duration"`
	Text     string `json:"text"`
	Language string `json:"language"`
}

// PageAudio ページ音声
type PageAudio struct {
	PageNumber    int            `json:"pageNumber"`
	Segments      []AudioSegment `json:"segments"`
	TotalDuration int64          `json:"totalDuration"`
}

// TeacherModePlaylist 教師モードプレイリスト
type TeacherModePlaylist struct {
	ID            string                `json:"id"`
	BookID        string                `json:"bookId"`
	Pages         []PageAudio           `json:"pages"`
	Settings      *TeacherModeSettings  `json:"settings"`
	TotalDuration int64                 `json:"totalDuration"`
}

// PageData ページデータ（OCRから取得）
type PageData struct {
	PageNumber      int
	PhraseText      string
	TranslationText string
	Language        string
}

// Service 教師モードサービス
type Service struct {
	// 将来的にTTS、STT、辞書APIクライアントを追加
}

// NewService 新しいサービスインスタンスを作成
func NewService() *Service {
	return &Service{}
}

// ValidateSettings 設定の検証
func (s *Service) ValidateSettings(settings *TeacherModeSettings) error {
	if settings == nil {
		return errors.New("settings is required")
	}

	// 速度の検証
	validSpeeds := []float64{0.5, 0.75, 1.0, 1.25, 1.5, 2.0}
	speedValid := false
	for _, speed := range validSpeeds {
		if settings.Speed == speed {
			speedValid = true
			break
		}
	}
	if !speedValid {
		return errors.New("invalid speed: must be one of 0.5, 0.75, 1.0, 1.25, 1.5, 2.0")
	}

	// ページ間隔の検証
	if settings.PageInterval < 0 || settings.PageInterval > 30 {
		return errors.New("invalid page interval: must be between 0 and 30 seconds")
	}

	// リピート回数の検証
	if settings.RepeatCount < 1 || settings.RepeatCount > 3 {
		return errors.New("invalid repeat count: must be between 1 and 3")
	}

	// 音質の検証
	if settings.AudioQuality != "standard" && settings.AudioQuality != "premium" {
		return errors.New("invalid audio quality: must be 'standard' or 'premium'")
	}

	return nil
}

// CalculateTotalDuration セグメントの総再生時間を計算
func (s *Service) CalculateTotalDuration(segments []AudioSegment) int64 {
	var total int64
	for _, seg := range segments {
		total += seg.Duration
	}
	return total
}

// GenerateAudioSegments ページの音声セグメントを生成
func (s *Service) GenerateAudioSegments(ctx context.Context, page *PageData, settings *TeacherModeSettings) ([]AudioSegment, error) {
	var segments []AudioSegment

	// 1. 学習先言語のフレーズ（必須）
	phraseSegment := AudioSegment{
		ID:       fmt.Sprintf("phrase-%d", page.PageNumber),
		Type:     "phrase",
		AudioURL: fmt.Sprintf("http://example.com/audio/phrase-%d.mp3", page.PageNumber),
		Duration: 2000, // 2秒（実際はTTS APIから取得）
		Text:     page.PhraseText,
		Language: page.Language,
	}
	segments = append(segments, phraseSegment)

	// 2. 母国語訳（オプション）
	if settings.Content.IncludeTranslation && page.TranslationText != "" {
		translationSegment := AudioSegment{
			ID:       fmt.Sprintf("translation-%d", page.PageNumber),
			Type:     "translation",
			AudioURL: fmt.Sprintf("http://example.com/audio/translation-%d.mp3", page.PageNumber),
			Duration: 1500, // 1.5秒
			Text:     page.TranslationText,
			Language: "ja", // 母国語（日本語）
		}
		segments = append(segments, translationSegment)
	}

	// 3. 単語解説（オプション）
	if settings.Content.IncludeWordExplanation {
		explanationSegment := AudioSegment{
			ID:       fmt.Sprintf("explanation-%d", page.PageNumber),
			Type:     "explanation",
			AudioURL: fmt.Sprintf("http://example.com/audio/explanation-%d.mp3", page.PageNumber),
			Duration: 3000, // 3秒
			Text:     fmt.Sprintf("%sの単語解説", page.PhraseText),
			Language: "ja",
		}
		segments = append(segments, explanationSegment)
	}

	return segments, nil
}
