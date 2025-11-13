package repository

import (
	"context"
	"testing"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWordRepository_Create は単語の作成テスト
func TestWordRepository_Create(t *testing.T) {
	repo := NewMockWordRepository()
	ctx := context.Background()

	word := &models.Word{
		UserID:     "user-1",
		BookID:     "book-1",
		PageNumber: 1,
		Text:       "Здравствуйте",
		Meaning:    "こんにちは",
		Language:   "ru",
	}

	err := repo.Create(ctx, word)
	require.NoError(t, err)
	assert.NotEmpty(t, word.ID, "単語IDが生成されること")
}

// TestWordRepository_CreateDuplicate は重複単語の作成エラーテスト
func TestWordRepository_CreateDuplicate(t *testing.T) {
	repo := NewMockWordRepository()
	ctx := context.Background()

	word := &models.Word{
		UserID:     "user-1",
		BookID:     "book-1",
		PageNumber: 1,
		Text:       "Здравствуйте",
		Meaning:    "こんにちは",
		Language:   "ru",
	}

	// 最初の作成は成功
	err := repo.Create(ctx, word)
	require.NoError(t, err)

	// 同じ単語を再度作成するとエラー
	err = repo.Create(ctx, word)
	assert.Error(t, err, "重複単語の作成はエラーになること")
}

// TestWordRepository_GetByID は単語IDによる取得テスト
func TestWordRepository_GetByID(t *testing.T) {
	repo := NewMockWordRepository()
	ctx := context.Background()

	word := &models.Word{
		UserID:     "user-1",
		BookID:     "book-1",
		PageNumber: 1,
		Text:       "Здравствуйте",
		Meaning:    "こんにちは",
		Language:   "ru",
	}

	err := repo.Create(ctx, word)
	require.NoError(t, err)

	// 取得
	retrieved, err := repo.GetByID(ctx, word.ID)
	require.NoError(t, err)
	assert.Equal(t, word.Text, retrieved.Text)
	assert.Equal(t, word.Meaning, retrieved.Meaning)
}

// TestWordRepository_GetByIDNotFound は存在しない単語IDの取得テスト
func TestWordRepository_GetByIDNotFound(t *testing.T) {
	repo := NewMockWordRepository()
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "non-existent-id")
	assert.Error(t, err, "存在しない単語IDはエラーになること")
}

// TestWordRepository_List は単語一覧取得テスト
func TestWordRepository_List(t *testing.T) {
	repo := NewMockWordRepository()
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
		{
			UserID:     "user-1",
			BookID:     "book-2",
			PageNumber: 1,
			Text:       "Hello",
			Meaning:    "こんにちは",
			Language:   "en",
		},
	}

	for _, word := range words {
		err := repo.Create(ctx, word)
		require.NoError(t, err)
	}

	// フィルタなしで全取得
	filter := &models.WordFilter{
		UserID: "user-1",
	}
	result, total, err := repo.List(ctx, filter)
	require.NoError(t, err)
	assert.Equal(t, 3, total, "全単語が取得できること")
	assert.Len(t, result, 3)
}

// TestWordRepository_ListWithFilter はフィルタ付き単語一覧取得テスト
func TestWordRepository_ListWithFilter(t *testing.T) {
	repo := NewMockWordRepository()
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
			Mastery:    85.0,
		},
		{
			UserID:     "user-1",
			BookID:     "book-1",
			PageNumber: 1,
			Text:       "Привет",
			Meaning:    "やあ",
			Language:   "ru",
			Mastery:    60.0,
		},
		{
			UserID:     "user-1",
			BookID:     "book-2",
			PageNumber: 1,
			Text:       "Hello",
			Meaning:    "こんにちは",
			Language:   "en",
			Mastery:    40.0,
		},
	}

	for _, word := range words {
		err := repo.Create(ctx, word)
		require.NoError(t, err)
	}

	// BookIDでフィルタ
	filter := &models.WordFilter{
		UserID: "user-1",
		BookID: "book-1",
	}
	result, total, err := repo.List(ctx, filter)
	require.NoError(t, err)
	assert.Equal(t, 2, total, "book-1の単語のみ取得できること")

	// 言語でフィルタ
	filter = &models.WordFilter{
		UserID:   "user-1",
		Language: "ru",
	}
	result, total, err = repo.List(ctx, filter)
	require.NoError(t, err)
	assert.Equal(t, 2, total, "ロシア語の単語のみ取得できること")

	// 習得度でフィルタ
	filter = &models.WordFilter{
		UserID:     "user-1",
		MinMastery: 70.0,
	}
	result, total, err = repo.List(ctx, filter)
	require.NoError(t, err)
	assert.Equal(t, 1, total, "習得度70%以上の単語のみ取得できること")
	assert.Equal(t, 85.0, result[0].Mastery)
}

// TestWordRepository_Search は単語検索テスト
func TestWordRepository_Search(t *testing.T) {
	repo := NewMockWordRepository()
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
		err := repo.Create(ctx, word)
		require.NoError(t, err)
	}

	// 検索
	filter := &models.WordFilter{
		UserID: "user-1",
		Query:  "Здрав",
	}
	result, total, err := repo.List(ctx, filter)
	require.NoError(t, err)
	assert.Equal(t, 1, total, "検索クエリにマッチする単語が取得できること")
	assert.Equal(t, "Здравствуйте", result[0].Text)
}

// TestWordRepository_Update は単語更新テスト
func TestWordRepository_Update(t *testing.T) {
	repo := NewMockWordRepository()
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

	err := repo.Create(ctx, word)
	require.NoError(t, err)

	// 更新
	word.Mastery = 80.0
	word.ReviewCount = 5
	err = repo.Update(ctx, word)
	require.NoError(t, err)

	// 確認
	updated, err := repo.GetByID(ctx, word.ID)
	require.NoError(t, err)
	assert.Equal(t, 80.0, updated.Mastery)
	assert.Equal(t, 5, updated.ReviewCount)
}

// TestWordRepository_Delete は単語削除テスト
func TestWordRepository_Delete(t *testing.T) {
	repo := NewMockWordRepository()
	ctx := context.Background()

	word := &models.Word{
		UserID:     "user-1",
		BookID:     "book-1",
		PageNumber: 1,
		Text:       "Здравствуйте",
		Meaning:    "こんにちは",
		Language:   "ru",
	}

	err := repo.Create(ctx, word)
	require.NoError(t, err)

	// 削除
	err = repo.Delete(ctx, word.ID)
	require.NoError(t, err)

	// 削除確認
	_, err = repo.GetByID(ctx, word.ID)
	assert.Error(t, err, "削除後は取得できないこと")
}

// TestWordRepository_GetStats は単語統計取得テスト
func TestWordRepository_GetStats(t *testing.T) {
	repo := NewMockWordRepository()
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
			BookID:      "book-2",
			PageNumber:  1,
			Text:        "Hello",
			Meaning:     "こんにちは",
			Language:    "en",
			Mastery:     90.0,
			ReviewCount: 7,
		},
	}

	for _, word := range words {
		err := repo.Create(ctx, word)
		require.NoError(t, err)
	}

	// 統計取得
	stats, err := repo.GetStats(ctx, "user-1", "book-1")
	require.NoError(t, err)
	assert.Equal(t, 2, stats.TotalWords, "全単語数が正しいこと")
	assert.Equal(t, 1, stats.MasteredWords, "習得単語数が正しいこと（習得度80%以上）")
	assert.Equal(t, 72.5, stats.AverageMastery, "平均習得度が正しいこと")
	assert.Equal(t, 8, stats.TotalReviews, "合計復習回数が正しいこと")
}

// TestWordRepository_BulkCreate は一括作成テスト
func TestWordRepository_BulkCreate(t *testing.T) {
	repo := NewMockWordRepository()
	ctx := context.Background()

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

	err := repo.BulkCreate(ctx, words)
	require.NoError(t, err)

	// 確認
	filter := &models.WordFilter{
		UserID: "user-1",
		BookID: "book-1",
	}
	result, total, err := repo.List(ctx, filter)
	require.NoError(t, err)
	assert.Equal(t, 2, total, "一括作成された単語が取得できること")
	assert.Len(t, result, 2)
}

// Benchmark tests
func BenchmarkWordRepository_Create(b *testing.B) {
	repo := NewMockWordRepository()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		word := &models.Word{
			UserID:     "user-1",
			BookID:     "book-1",
			PageNumber: 1,
			Text:       "test-word",
			Meaning:    "テスト",
			Language:   "en",
		}
		repo.Create(ctx, word)
	}
}

func BenchmarkWordRepository_List(b *testing.B) {
	repo := NewMockWordRepository()
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
		repo.Create(ctx, word)
	}

	filter := &models.WordFilter{
		UserID: "user-1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.List(ctx, filter)
	}
}
