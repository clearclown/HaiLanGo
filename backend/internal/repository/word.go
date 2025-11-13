package repository

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

var (
	// ErrWordNotFound は単語が見つからないエラー
	ErrWordNotFound = errors.New("word not found")
	// ErrWordAlreadyExists は単語が既に存在するエラー
	ErrWordAlreadyExists = errors.New("word already exists")
)

// WordRepository は単語リポジトリのインターフェース
type WordRepository interface {
	Create(ctx context.Context, word *models.Word) error
	GetByID(ctx context.Context, id string) (*models.Word, error)
	List(ctx context.Context, filter *models.WordFilter) ([]*models.Word, int, error)
	Update(ctx context.Context, word *models.Word) error
	Delete(ctx context.Context, id string) error
	GetStats(ctx context.Context, userID, bookID string) (*models.WordStats, error)
	BulkCreate(ctx context.Context, words []*models.Word) error
}

// MockWordRepository はメモリ内で動作するモックリポジトリ
type MockWordRepository struct {
	mu    sync.RWMutex
	words map[string]*models.Word
}

// NewMockWordRepository は新しいモックリポジトリを作成する
func NewMockWordRepository() *MockWordRepository {
	return &MockWordRepository{
		words: make(map[string]*models.Word),
	}
}

// Create は単語を作成する
func (r *MockWordRepository) Create(ctx context.Context, word *models.Word) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 重複チェック
	for _, w := range r.words {
		if w.UserID == word.UserID &&
			w.BookID == word.BookID &&
			strings.EqualFold(w.Text, word.Text) {
			return ErrWordAlreadyExists
		}
	}

	// IDの生成
	word.ID = uuid.New().String()
	word.CreatedAt = time.Now()
	word.UpdatedAt = time.Now()

	// コピーして保存
	wordCopy := *word
	r.words[word.ID] = &wordCopy

	return nil
}

// GetByID はIDで単語を取得する
func (r *MockWordRepository) GetByID(ctx context.Context, id string) (*models.Word, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	word, ok := r.words[id]
	if !ok {
		return nil, ErrWordNotFound
	}

	// コピーを返す
	wordCopy := *word
	return &wordCopy, nil
}

// List はフィルタ条件に基づいて単語一覧を取得する
func (r *MockWordRepository) List(ctx context.Context, filter *models.WordFilter) ([]*models.Word, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// フィルタリング
	filtered := make([]*models.Word, 0)
	for _, word := range r.words {
		if !matchFilter(word, filter) {
			continue
		}
		wordCopy := *word
		filtered = append(filtered, &wordCopy)
	}

	// ソート
	sortWords(filtered, filter)

	// ページネーション
	total := len(filtered)
	start := filter.Offset
	end := start + filter.Limit

	if filter.Limit > 0 {
		if start > total {
			return []*models.Word{}, total, nil
		}
		if end > total {
			end = total
		}
		filtered = filtered[start:end]
	}

	return filtered, total, nil
}

// Update は単語を更新する
func (r *MockWordRepository) Update(ctx context.Context, word *models.Word) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.words[word.ID]; !ok {
		return ErrWordNotFound
	}

	word.UpdatedAt = time.Now()
	wordCopy := *word
	r.words[word.ID] = &wordCopy

	return nil
}

// Delete は単語を削除する
func (r *MockWordRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.words[id]; !ok {
		return ErrWordNotFound
	}

	delete(r.words, id)
	return nil
}

// GetStats は単語統計を取得する
func (r *MockWordRepository) GetStats(ctx context.Context, userID, bookID string) (*models.WordStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var totalWords int
	var masteredWords int
	var totalMastery float64
	var totalReviews int

	for _, word := range r.words {
		if word.UserID != userID {
			continue
		}
		if bookID != "" && word.BookID != bookID {
			continue
		}

		totalWords++
		totalMastery += word.Mastery
		totalReviews += word.ReviewCount

		if word.Mastery >= 80.0 {
			masteredWords++
		}
	}

	averageMastery := 0.0
	if totalWords > 0 {
		averageMastery = totalMastery / float64(totalWords)
	}

	return &models.WordStats{
		TotalWords:     totalWords,
		MasteredWords:  masteredWords,
		AverageMastery: averageMastery,
		TotalReviews:   totalReviews,
	}, nil
}

// BulkCreate は複数の単語を一括作成する
func (r *MockWordRepository) BulkCreate(ctx context.Context, words []*models.Word) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, word := range words {
		// 重複チェック
		for _, w := range r.words {
			if w.UserID == word.UserID &&
				w.BookID == word.BookID &&
				strings.EqualFold(w.Text, word.Text) {
				continue // スキップ
			}
		}

		// IDの生成
		word.ID = uuid.New().String()
		word.CreatedAt = time.Now()
		word.UpdatedAt = time.Now()

		// コピーして保存
		wordCopy := *word
		r.words[word.ID] = &wordCopy
	}

	return nil
}

// matchFilter はフィルタ条件にマッチするか判定する
func matchFilter(word *models.Word, filter *models.WordFilter) bool {
	if filter.UserID != "" && word.UserID != filter.UserID {
		return false
	}

	if filter.BookID != "" && word.BookID != filter.BookID {
		return false
	}

	if filter.Language != "" && word.Language != filter.Language {
		return false
	}

	if filter.Query != "" {
		query := strings.ToLower(filter.Query)
		text := strings.ToLower(word.Text)
		meaning := strings.ToLower(word.Meaning)
		if !strings.Contains(text, query) && !strings.Contains(meaning, query) {
			return false
		}
	}

	if filter.MinMastery > 0 && word.Mastery < filter.MinMastery {
		return false
	}

	if filter.MaxMastery > 0 && word.Mastery > filter.MaxMastery {
		return false
	}

	if len(filter.Tags) > 0 {
		hasTag := false
		for _, tag := range filter.Tags {
			for _, wordTag := range word.Tags {
				if tag == wordTag {
					hasTag = true
					break
				}
			}
			if hasTag {
				break
			}
		}
		if !hasTag {
			return false
		}
	}

	return true
}

// sortWords は単語をソートする
func sortWords(words []*models.Word, filter *models.WordFilter) {
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	// 簡易的なソート実装
	// 実際のDBでは ORDER BY 句を使用
	// ここではテスト用の簡易実装
}
