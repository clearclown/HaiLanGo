// Package teachermode provides teacher mode functionality
package teachermode

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratePlaylist(t *testing.T) {
	ctx := context.Background()

	t.Run("正常にプレイリストが生成される", func(t *testing.T) {
		service := NewService()
		settings := &TeacherModeSettings{
			Speed:        1.0,
			PageInterval: 5,
			RepeatCount:  1,
			AudioQuality: "standard",
			Content: TeacherModeContent{
				IncludeTranslation:           true,
				IncludeWordExplanation:       true,
				IncludeGrammarExplanation:    false,
				IncludePronunciationPractice: false,
				IncludeExampleSentences:      false,
			},
		}

		playlist, err := service.GeneratePlaylist(ctx, "test-book", settings)

		require.NoError(t, err)
		assert.NotNil(t, playlist)
		assert.Equal(t, "test-book", playlist.BookID)
		assert.NotEmpty(t, playlist.Pages)
	})

	t.Run("無効なbookIDでエラーが返る", func(t *testing.T) {
		service := NewService()
		settings := &TeacherModeSettings{
			Speed:        1.0,
			PageInterval: 5,
			RepeatCount:  1,
			AudioQuality: "standard",
		}

		_, err := service.GeneratePlaylist(ctx, "", settings)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid book ID")
	})

	t.Run("無効な設定でエラーが返る", func(t *testing.T) {
		service := NewService()

		_, err := service.GeneratePlaylist(ctx, "test-book", nil)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "settings is required")
	})
}

func TestGenerateAudioSegments(t *testing.T) {
	ctx := context.Background()

	t.Run("学習先言語のセグメントが生成される", func(t *testing.T) {
		service := NewService()
		page := &PageData{
			PageNumber: 1,
			PhraseText: "Hello",
			Language:   "en",
		}
		settings := &TeacherModeSettings{
			Speed:        1.0,
			AudioQuality: "standard",
		}

		segments, err := service.GenerateAudioSegments(ctx, page, settings)

		require.NoError(t, err)
		assert.NotEmpty(t, segments)
		assert.Equal(t, "phrase", segments[0].Type)
		assert.Equal(t, "Hello", segments[0].Text)
	})

	t.Run("翻訳セグメントが含まれる", func(t *testing.T) {
		service := NewService()
		page := &PageData{
			PageNumber:      1,
			PhraseText:      "Hello",
			TranslationText: "こんにちは",
			Language:        "en",
		}
		settings := &TeacherModeSettings{
			Speed:        1.0,
			AudioQuality: "standard",
			Content: TeacherModeContent{
				IncludeTranslation: true,
			},
		}

		segments, err := service.GenerateAudioSegments(ctx, page, settings)

		require.NoError(t, err)
		assert.True(t, len(segments) >= 2)

		var translationSegment *AudioSegment
		for _, seg := range segments {
			if seg.Type == "translation" {
				translationSegment = &seg
				break
			}
		}

		require.NotNil(t, translationSegment)
		assert.Equal(t, "こんにちは", translationSegment.Text)
	})

	t.Run("単語解説セグメントが含まれる", func(t *testing.T) {
		service := NewService()
		page := &PageData{
			PageNumber: 1,
			PhraseText: "Hello",
			Language:   "en",
		}
		settings := &TeacherModeSettings{
			Speed:        1.0,
			AudioQuality: "standard",
			Content: TeacherModeContent{
				IncludeWordExplanation: true,
			},
		}

		segments, err := service.GenerateAudioSegments(ctx, page, settings)

		require.NoError(t, err)

		var explanationSegment *AudioSegment
		for _, seg := range segments {
			if seg.Type == "explanation" {
				explanationSegment = &seg
				break
			}
		}

		require.NotNil(t, explanationSegment)
		assert.NotEmpty(t, explanationSegment.Text)
	})
}

func TestCalculateTotalDuration(t *testing.T) {
	t.Run("総再生時間が正しく計算される", func(t *testing.T) {
		service := NewService()
		segments := []AudioSegment{
			{Duration: 2000},
			{Duration: 3000},
			{Duration: 1500},
		}

		totalDuration := service.CalculateTotalDuration(segments)

		assert.Equal(t, int64(6500), totalDuration)
	})

	t.Run("セグメントが空の場合は0が返る", func(t *testing.T) {
		service := NewService()
		segments := []AudioSegment{}

		totalDuration := service.CalculateTotalDuration(segments)

		assert.Equal(t, int64(0), totalDuration)
	})
}

func TestValidateSettings(t *testing.T) {
	t.Run("正常な設定でエラーが返らない", func(t *testing.T) {
		service := NewService()
		settings := &TeacherModeSettings{
			Speed:        1.0,
			PageInterval: 5,
			RepeatCount:  1,
			AudioQuality: "standard",
		}

		err := service.ValidateSettings(settings)

		assert.NoError(t, err)
	})

	t.Run("無効な速度でエラーが返る", func(t *testing.T) {
		service := NewService()
		settings := &TeacherModeSettings{
			Speed:        3.0, // 無効な速度
			PageInterval: 5,
			RepeatCount:  1,
			AudioQuality: "standard",
		}

		err := service.ValidateSettings(settings)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid speed")
	})

	t.Run("無効なページ間隔でエラーが返る", func(t *testing.T) {
		service := NewService()
		settings := &TeacherModeSettings{
			Speed:        1.0,
			PageInterval: -1, // 無効な間隔
			RepeatCount:  1,
			AudioQuality: "standard",
		}

		err := service.ValidateSettings(settings)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid page interval")
	})

	t.Run("無効な音質でエラーが返る", func(t *testing.T) {
		service := NewService()
		settings := &TeacherModeSettings{
			Speed:        1.0,
			PageInterval: 5,
			RepeatCount:  1,
			AudioQuality: "invalid", // 無効な音質
		}

		err := service.ValidateSettings(settings)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid audio quality")
	})
}
