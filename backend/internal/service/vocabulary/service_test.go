package vocabulary

import (
	"context"
	"testing"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestVocabularyService_AutoCollectWords は自動単語収集のテスト
func TestVocabularyService_AutoCollectWords(t *testing.T) {
	service := NewMockVocabularyService()
	ctx := context.Background()

	userID := "user-1"
	bookID := "book-1"
	pageNumber := 1
	text := "Здравствуйте! Как дела? Меня зовут Иван."
	language := "ru"

	err := service.AutoCollectWords(ctx, userID, bookID, pageNumber, text, language)
	require.NoError(t, err)

	// 収集された単語を確認
	filter := &models.WordFilter{
		UserID: userID,
		BookID: bookID,
	}
	words, err := service.GetWords(ctx, filter)
	require.NoError(t, err)
	assert.NotEmpty(t, words, "単語が自動収集されること")
}

// TestVocabularyService_AddWord は単語追加のテスト
func TestVocabularyService_AddWord(t *testing.T) {
	service := NewMockVocabularyService()
	ctx := context.Background()

	word := &models.Word{
		UserID:     "user-1",
		BookID:     "book-1",
		PageNumber: 1,
		Text:       "Здравствуйте",
		Meaning:    "こんにちは",
		Language:   "ru",
	}

	err := service.AddWord(ctx, word)
	require.NoError(t, err)
	assert.NotEmpty(t, word.ID, "単語IDが生成されること")
}

// TestVocabularyService_AddWordDuplicate は重複単語追加のエラーテスト
func TestVocabularyService_AddWordDuplicate(t *testing.T) {
	service := NewMockVocabularyService()
	ctx := context.Background()

	word := &models.Word{
		UserID:     "user-1",
		BookID:     "book-1",
		PageNumber: 1,
		Text:       "Здравствуйте",
		Meaning:    "こんにちは",
		Language:   "ru",
	}

	// 最初の追加は成功
	err := service.AddWord(ctx, word)
	require.NoError(t, err)

	// 同じ単語を再度追加するとエラー
	duplicateWord := &models.Word{
		UserID:     "user-1",
		BookID:     "book-1",
		PageNumber: 1,
		Text:       "Здравствуйте",
		Meaning:    "こんにちは",
		Language:   "ru",
	}
	err = service.AddWord(ctx, duplicateWord)
	assert.Error(t, err, "重複単語の追加はエラーになること")
}

// TestVocabularyService_GetWords は単語取得のテスト
func TestVocabularyService_GetWords(t *testing.T) {
	service := NewMockVocabularyService()
	ctx := context.Background()

	// テストデータの作成
	words := []*models.Word{
		{
			UserID:     "user-1",
			BookID:     "book-1",
			PageNumber: 1,
			Text:       "Здравствуйте",
			Meaning:    "こんにちは",
			Language:   "ru",
		},
		{
			UserID:     "user-1",
			BookID:     "book-1",
			PageNumber: 1,
			Text:       "Привет",
			Meaning:    "やあ",
			Language:   "ru",
		},
	}

	for _, word := range words {
		err := service.AddWord(ctx, word)
		require.NoError(t, err)
	}

	// 取得
	filter := &models.WordFilter{
		UserID: "user-1",
		BookID: "book-1",
	}
	result, err := service.GetWords(ctx, filter)
	require.NoError(t, err)
	assert.Len(t, result, 2, "追加した単語が取得できること")
}

// TestVocabularyService_SearchWords は単語検索のテスト
func TestVocabularyService_SearchWords(t *testing.T) {
	service := NewMockVocabularyService()
	ctx := context.Background()

	// テストデータの作成
	words := []*models.Word{
		{
			UserID:     "user-1",
			BookID:     "book-1",
			PageNumber: 1,
			Text:       "Здравствуйте",
			Meaning:    "こんにちは",
			Language:   "ru",
		},
		{
			UserID:     "user-1",
			BookID:     "book-1",
			PageNumber: 1,
			Text:       "Привет",
			Meaning:    "やあ",
			Language:   "ru",
		},
	}

	for _, word := range words {
		err := service.AddWord(ctx, word)
		require.NoError(t, err)
	}

	// 検索
	filter := &models.WordFilter{
		UserID: "user-1",
		Query:  "Здрав",
	}
	result, err := service.GetWords(ctx, filter)
	require.NoError(t, err)
	assert.Len(t, result, 1, "検索クエリにマッチする単語が取得できること")
	assert.Equal(t, "Здравствуйте", result[0].Text)
}

// TestVocabularyService_UpdateWord は単語更新のテスト
func TestVocabularyService_UpdateWord(t *testing.T) {
	service := NewMockVocabularyService()
	ctx := context.Background()

	word := &models.Word{
		UserID:     "user-1",
		BookID:     "book-1",
		PageNumber: 1,
		Text:       "Здравствуйте",
		Meaning:    "こんにちは",
		Language:   "ru",
		Mastery:    50.0,
	}

	err := service.AddWord(ctx, word)
	require.NoError(t, err)

	// 更新
	word.Mastery = 80.0
	word.Meaning = "こんにちは（丁寧）"
	err = service.UpdateWord(ctx, word)
	require.NoError(t, err)

	// 確認
	updated, err := service.GetWordByID(ctx, word.ID)
	require.NoError(t, err)
	assert.Equal(t, 80.0, updated.Mastery)
	assert.Equal(t, "こんにちは（丁寧）", updated.Meaning)
}

// TestVocabularyService_DeleteWord は単語削除のテスト
func TestVocabularyService_DeleteWord(t *testing.T) {
	service := NewMockVocabularyService()
	ctx := context.Background()

	word := &models.Word{
		UserID:     "user-1",
		BookID:     "book-1",
		PageNumber: 1,
		Text:       "Здравствуйте",
		Meaning:    "こんにちは",
		Language:   "ru",
	}

	err := service.AddWord(ctx, word)
	require.NoError(t, err)

	// 削除
	err = service.DeleteWord(ctx, word.ID)
	require.NoError(t, err)

	// 削除確認
	_, err = service.GetWordByID(ctx, word.ID)
	assert.Error(t, err, "削除後は取得できないこと")
}

// TestVocabularyService_RecordReview は学習記録のテスト
func TestVocabularyService_RecordReview(t *testing.T) {
	service := NewMockVocabularyService()
	ctx := context.Background()

	word := &models.Word{
		UserID:      "user-1",
		BookID:      "book-1",
		PageNumber:  1,
		Text:        "Здравствуйте",
		Meaning:     "こんにちは",
		Language:    "ru",
		ReviewCount: 0,
		Mastery:     0.0,
	}

	err := service.AddWord(ctx, word)
	require.NoError(t, err)

	// 学習記録
	score := 85.0
	err = service.RecordReview(ctx, word.ID, score)
	require.NoError(t, err)

	// 確認
	updated, err := service.GetWordByID(ctx, word.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, updated.ReviewCount, "学習回数が増加すること")
	assert.Greater(t, updated.Mastery, 0.0, "習得度が更新されること")
}

// TestVocabularyService_CalculateMastery は習得度計算のテスト
func TestVocabularyService_CalculateMastery(t *testing.T) {
	tests := []struct {
		name         string
		reviewCount  int
		averageScore float64
		expected     float64
	}{
		{
			name:         "初回学習（高スコア）",
			reviewCount:  1,
			averageScore: 90.0,
			expected:     9.0, // 初回なので1/10程度（90.0 * 0.1）
		},
		{
			name:         "5回学習（高スコア）",
			reviewCount:  5,
			averageScore: 85.0,
			expected:     42.5, // 5回なので半分程度（85.0 * 0.5）
		},
		{
			name:         "10回学習（完璧）",
			reviewCount:  10,
			averageScore: 95.0,
			expected:     95.0, // 10回で完全習得（95.0 * 1.0）
		},
		{
			name:         "3回学習（低スコア）",
			reviewCount:  3,
			averageScore: 50.0,
			expected:     15.0, // 3回なので30%程度（50.0 * 0.3）
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mastery := CalculateMastery(tt.reviewCount, tt.averageScore)
			assert.InDelta(t, tt.expected, mastery, 10.0, "習得度が期待値に近いこと")
		})
	}
}

// TestVocabularyService_GetStats は統計取得のテスト
func TestVocabularyService_GetStats(t *testing.T) {
	service := NewMockVocabularyService()
	ctx := context.Background()

	// テストデータの作成
	words := []*models.Word{
		{
			UserID:      "user-1",
			BookID:      "book-1",
			PageNumber:  1,
			Text:        "Здравствуйте",
			Meaning:     "こんにちは",
			Language:    "ru",
			Mastery:     85.0,
			ReviewCount: 5,
		},
		{
			UserID:      "user-1",
			BookID:      "book-1",
			PageNumber:  1,
			Text:        "Привет",
			Meaning:     "やあ",
			Language:    "ru",
			Mastery:     60.0,
			ReviewCount: 3,
		},
		{
			UserID:      "user-1",
			BookID:      "book-1",
			PageNumber:  1,
			Text:        "Спасибо",
			Meaning:     "ありがとう",
			Language:    "ru",
			Mastery:     90.0,
			ReviewCount: 7,
		},
	}

	for _, word := range words {
		err := service.AddWord(ctx, word)
		require.NoError(t, err)
	}

	// 統計取得
	stats, err := service.GetStats(ctx, "user-1", "book-1")
	require.NoError(t, err)
	assert.Equal(t, 3, stats.TotalWords, "全単語数が正しいこと")
	assert.Equal(t, 2, stats.MasteredWords, "習得単語数が正しいこと（習得度80%以上）")
	assert.InDelta(t, 78.33, stats.AverageMastery, 1.0, "平均習得度が正しいこと")
	assert.Equal(t, 15, stats.TotalReviews, "合計復習回数が正しいこと")
}

// TestVocabularyService_ExportWords は単語エクスポートのテスト
func TestVocabularyService_ExportWords(t *testing.T) {
	service := NewMockVocabularyService()
	ctx := context.Background()

	// テストデータの作成
	words := []*models.Word{
		{
			UserID:     "user-1",
			BookID:     "book-1",
			PageNumber: 1,
			Text:       "Здравствуйте",
			Meaning:    "こんにちは",
			Language:   "ru",
		},
		{
			UserID:     "user-1",
			BookID:     "book-1",
			PageNumber: 1,
			Text:       "Привет",
			Meaning:    "やあ",
			Language:   "ru",
		},
	}

	for _, word := range words {
		err := service.AddWord(ctx, word)
		require.NoError(t, err)
	}

	// エクスポート
	filter := &models.WordFilter{
		UserID: "user-1",
		BookID: "book-1",
	}
	csvData, err := service.ExportWordsToCSV(ctx, filter)
	require.NoError(t, err)
	assert.NotEmpty(t, csvData, "CSVデータが生成されること")
	assert.Contains(t, string(csvData), "Здравствуйте", "単語がCSVに含まれること")
	assert.Contains(t, string(csvData), "こんにちは", "意味がCSVに含まれること")
}

// TestVocabularyService_AddTags はタグ追加のテスト
func TestVocabularyService_AddTags(t *testing.T) {
	service := NewMockVocabularyService()
	ctx := context.Background()

	word := &models.Word{
		UserID:     "user-1",
		BookID:     "book-1",
		PageNumber: 1,
		Text:       "Здравствуйте",
		Meaning:    "こんにちは",
		Language:   "ru",
		Tags:       []string{},
	}

	err := service.AddWord(ctx, word)
	require.NoError(t, err)

	// タグ追加
	tags := []string{"挨拶", "基本"}
	err = service.AddTags(ctx, word.ID, tags)
	require.NoError(t, err)

	// 確認
	updated, err := service.GetWordByID(ctx, word.ID)
	require.NoError(t, err)
	assert.ElementsMatch(t, tags, updated.Tags, "タグが追加されること")
}

// TestVocabularyService_FilterByTags はタグでのフィルタテスト
func TestVocabularyService_FilterByTags(t *testing.T) {
	service := NewMockVocabularyService()
	ctx := context.Background()

	// テストデータの作成
	words := []*models.Word{
		{
			UserID:     "user-1",
			BookID:     "book-1",
			PageNumber: 1,
			Text:       "Здравствуйте",
			Meaning:    "こんにちは",
			Language:   "ru",
			Tags:       []string{"挨拶", "基本"},
		},
		{
			UserID:     "user-1",
			BookID:     "book-1",
			PageNumber: 1,
			Text:       "Спасибо",
			Meaning:    "ありがとう",
			Language:   "ru",
			Tags:       []string{"感謝", "基本"},
		},
		{
			UserID:     "user-1",
			BookID:     "book-1",
			PageNumber: 1,
			Text:       "Пожалуйста",
			Meaning:    "どういたしまして",
			Language:   "ru",
			Tags:       []string{"感謝"},
		},
	}

	for _, word := range words {
		err := service.AddWord(ctx, word)
		require.NoError(t, err)
	}

	// タグでフィルタ
	filter := &models.WordFilter{
		UserID: "user-1",
		Tags:   []string{"基本"},
	}
	result, err := service.GetWords(ctx, filter)
	require.NoError(t, err)
	assert.Len(t, result, 2, "「基本」タグを持つ単語が取得できること")
}

// Benchmark tests
func BenchmarkVocabularyService_AutoCollectWords(b *testing.B) {
	service := NewMockVocabularyService()
	ctx := context.Background()

	text := "Здравствуйте! Как дела? Меня зовут Иван. Я изучаю русский язык."
	language := "ru"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.AutoCollectWords(ctx, "user-1", "book-1", 1, text, language)
	}
}

func BenchmarkVocabularyService_GetWords(b *testing.B) {
	service := NewMockVocabularyService()
	ctx := context.Background()

	// テストデータの準備
	for i := 0; i < 100; i++ {
		word := &models.Word{
			UserID:     "user-1",
			BookID:     "book-1",
			PageNumber: i,
			Text:       "test-word",
			Meaning:    "テスト",
			Language:   "en",
		}
		service.AddWord(ctx, word)
	}

	filter := &models.WordFilter{
		UserID: "user-1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetWords(ctx, filter)
	}
}
