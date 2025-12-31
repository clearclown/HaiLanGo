package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	GetByID(ctx context.Context, id uuid.UUID) (*models.Word, error)
	List(ctx context.Context, filter *models.WordFilter) ([]*models.Word, int, error)
	Update(ctx context.Context, word *models.Word) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetStats(ctx context.Context, userID, bookID uuid.UUID) (*models.WordStats, error)
	BulkCreate(ctx context.Context, words []*models.Word) error
}

// NewWordRepository は環境変数に基づいて適切なWordRepositoryを返す
func NewWordRepository(db *sql.DB) WordRepository {
	useMock := os.Getenv("USE_MOCK_APIS") == "true" ||
		os.Getenv("TEST_USE_MOCKS") == "true"

	if useMock || db == nil {
		return NewMockWordRepository()
	}

	return NewWordRepositoryPostgres(db)
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
	word.ID = uuid.New()
	word.CreatedAt = time.Now()
	word.UpdatedAt = time.Now()

	// コピーして保存
	wordCopy := *word
	r.words[word.ID.String()] = &wordCopy

	return nil
}

// GetByID はIDで単語を取得する
func (r *MockWordRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Word, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	word, ok := r.words[id.String()]
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

	if _, ok := r.words[word.ID.String()]; !ok {
		return ErrWordNotFound
	}

	word.UpdatedAt = time.Now()
	wordCopy := *word
	r.words[word.ID.String()] = &wordCopy

	return nil
}

// Delete は単語を削除する
func (r *MockWordRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.words[id.String()]; !ok {
		return ErrWordNotFound
	}

	delete(r.words, id.String())
	return nil
}

// GetStats は単語統計を取得する
func (r *MockWordRepository) GetStats(ctx context.Context, userID, bookID uuid.UUID) (*models.WordStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var totalWords int
	var masteredWords int
	var totalMastery float64
	var totalReviews int

	emptyUUID := uuid.UUID{}
	for _, word := range r.words {
		if word.UserID != userID {
			continue
		}
		if bookID != emptyUUID && word.BookID != bookID {
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
		word.ID = uuid.New()
		word.CreatedAt = time.Now()
		word.UpdatedAt = time.Now()

		// コピーして保存
		wordCopy := *word
		r.words[word.ID.String()] = &wordCopy
	}

	return nil
}

// matchFilter はフィルタ条件にマッチするか判定する
func matchFilter(word *models.Word, filter *models.WordFilter) bool {
	emptyUUID := uuid.UUID{}
	if filter.UserID != emptyUUID && word.UserID != filter.UserID {
		return false
	}

	if filter.BookID != emptyUUID && word.BookID != filter.BookID {
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

// PostgreSQL実装

// wordRepositoryPostgres はPostgreSQLベースの単語リポジトリ実装
type wordRepositoryPostgres struct {
	db *sql.DB
}

// NewWordRepositoryPostgres は新しいPostgreSQL実装のWordRepositoryを作成する
func NewWordRepositoryPostgres(db *sql.DB) WordRepository {
	return &wordRepositoryPostgres{db: db}
}

// Create は単語を作成する
func (r *wordRepositoryPostgres) Create(ctx context.Context, word *models.Word) error {
	// 重複チェック
	var exists bool
	checkQuery := `
		SELECT EXISTS(
			SELECT 1 FROM words
			WHERE user_id = $1 AND book_id = $2 AND LOWER(text) = LOWER($3)
		)
	`
	err := r.db.QueryRowContext(ctx, checkQuery, word.UserID, word.BookID, word.Text).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return ErrWordAlreadyExists
	}

	// IDの生成
	if word.ID == uuid.Nil {
		word.ID = uuid.New()
	}
	now := time.Now()
	word.CreatedAt = now
	word.UpdatedAt = now

	query := `
		INSERT INTO words (id, user_id, book_id, page_number, text, meaning, pronunciation,
		                   part_of_speech, example, language, review_count, average_score,
		                   mastery, tags, last_reviewed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		word.ID,
		word.UserID,
		word.BookID,
		word.PageNumber,
		word.Text,
		word.Meaning,
		word.Pronunciation,
		word.PartOfSpeech,
		word.Example,
		word.Language,
		word.ReviewCount,
		word.AverageScore,
		word.Mastery,
		pq.Array(word.Tags),
		word.LastReviewedAt,
		word.CreatedAt,
		word.UpdatedAt,
	)

	return err
}

// GetByID はIDで単語を取得する
func (r *wordRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (*models.Word, error) {
	query := `
		SELECT id, user_id, book_id, page_number, text, meaning, pronunciation,
		       part_of_speech, example, language, review_count, average_score,
		       mastery, tags, last_reviewed_at, created_at, updated_at
		FROM words
		WHERE id = $1
	`

	word := &models.Word{}
	var tags pq.StringArray
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&word.ID,
		&word.UserID,
		&word.BookID,
		&word.PageNumber,
		&word.Text,
		&word.Meaning,
		&word.Pronunciation,
		&word.PartOfSpeech,
		&word.Example,
		&word.Language,
		&word.ReviewCount,
		&word.AverageScore,
		&word.Mastery,
		&tags,
		&word.LastReviewedAt,
		&word.CreatedAt,
		&word.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrWordNotFound
		}
		return nil, err
	}

	word.Tags = []string(tags)
	return word, nil
}

// List はフィルタ条件に基づいて単語一覧を取得する
func (r *wordRepositoryPostgres) List(ctx context.Context, filter *models.WordFilter) ([]*models.Word, int, error) {
	// ベースクエリ
	baseQuery := `FROM words WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	// フィルタ条件の構築
	emptyUUID := uuid.UUID{}
	if filter.UserID != emptyUUID {
		baseQuery += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, filter.UserID)
		argIndex++
	}

	if filter.BookID != emptyUUID {
		baseQuery += fmt.Sprintf(" AND book_id = $%d", argIndex)
		args = append(args, filter.BookID)
		argIndex++
	}

	if filter.Language != "" {
		baseQuery += fmt.Sprintf(" AND language = $%d", argIndex)
		args = append(args, filter.Language)
		argIndex++
	}

	if filter.Query != "" {
		baseQuery += fmt.Sprintf(" AND (LOWER(text) LIKE $%d OR LOWER(meaning) LIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+strings.ToLower(filter.Query)+"%")
		argIndex++
	}

	if filter.MinMastery > 0 {
		baseQuery += fmt.Sprintf(" AND mastery >= $%d", argIndex)
		args = append(args, filter.MinMastery)
		argIndex++
	}

	if filter.MaxMastery > 0 {
		baseQuery += fmt.Sprintf(" AND mastery <= $%d", argIndex)
		args = append(args, filter.MaxMastery)
		argIndex++
	}

	if len(filter.Tags) > 0 {
		baseQuery += fmt.Sprintf(" AND tags && $%d", argIndex)
		args = append(args, pq.Array(filter.Tags))
		argIndex++
	}

	// カウントクエリ
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// ソート順
	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	// SQL injection対策
	validSortColumns := map[string]bool{
		"created_at":   true,
		"mastery":      true,
		"review_count": true,
		"text":         true,
		"updated_at":   true,
	}
	if !validSortColumns[sortBy] {
		sortBy = "created_at"
	}

	sortOrder := strings.ToUpper(filter.SortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	// データ取得クエリ
	selectQuery := `
		SELECT id, user_id, book_id, page_number, text, meaning, pronunciation,
		       part_of_speech, example, language, review_count, average_score,
		       mastery, tags, last_reviewed_at, created_at, updated_at
	` + baseQuery + fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// ページネーション
	if filter.Limit > 0 {
		selectQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, filter.Limit, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var words []*models.Word
	for rows.Next() {
		word := &models.Word{}
		var tags pq.StringArray
		err := rows.Scan(
			&word.ID,
			&word.UserID,
			&word.BookID,
			&word.PageNumber,
			&word.Text,
			&word.Meaning,
			&word.Pronunciation,
			&word.PartOfSpeech,
			&word.Example,
			&word.Language,
			&word.ReviewCount,
			&word.AverageScore,
			&word.Mastery,
			&tags,
			&word.LastReviewedAt,
			&word.CreatedAt,
			&word.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		word.Tags = []string(tags)
		words = append(words, word)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return words, total, nil
}

// Update は単語を更新する
func (r *wordRepositoryPostgres) Update(ctx context.Context, word *models.Word) error {
	query := `
		UPDATE words
		SET text = $1, meaning = $2, pronunciation = $3, part_of_speech = $4,
		    example = $5, language = $6, review_count = $7, average_score = $8,
		    mastery = $9, tags = $10, last_reviewed_at = $11, updated_at = NOW()
		WHERE id = $12
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		word.Text,
		word.Meaning,
		word.Pronunciation,
		word.PartOfSpeech,
		word.Example,
		word.Language,
		word.ReviewCount,
		word.AverageScore,
		word.Mastery,
		pq.Array(word.Tags),
		word.LastReviewedAt,
		word.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrWordNotFound
	}

	return nil
}

// Delete は単語を削除する
func (r *wordRepositoryPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM words WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrWordNotFound
	}

	return nil
}

// GetStats は単語統計を取得する
func (r *wordRepositoryPostgres) GetStats(ctx context.Context, userID, bookID uuid.UUID) (*models.WordStats, error) {
	var query string
	var args []interface{}

	emptyUUID := uuid.UUID{}
	if bookID != emptyUUID {
		query = `
			SELECT
				COUNT(*) as total_words,
				COALESCE(SUM(CASE WHEN mastery >= 80 THEN 1 ELSE 0 END), 0) as mastered_words,
				COALESCE(AVG(mastery), 0) as average_mastery,
				COALESCE(SUM(review_count), 0) as total_reviews
			FROM words
			WHERE user_id = $1 AND book_id = $2
		`
		args = []interface{}{userID, bookID}
	} else {
		query = `
			SELECT
				COUNT(*) as total_words,
				COALESCE(SUM(CASE WHEN mastery >= 80 THEN 1 ELSE 0 END), 0) as mastered_words,
				COALESCE(AVG(mastery), 0) as average_mastery,
				COALESCE(SUM(review_count), 0) as total_reviews
			FROM words
			WHERE user_id = $1
		`
		args = []interface{}{userID}
	}

	stats := &models.WordStats{}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&stats.TotalWords,
		&stats.MasteredWords,
		&stats.AverageMastery,
		&stats.TotalReviews,
	)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// BulkCreate は複数の単語を一括作成する
func (r *wordRepositoryPostgres) BulkCreate(ctx context.Context, words []*models.Word) error {
	if len(words) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO words (id, user_id, book_id, page_number, text, meaning, pronunciation,
		                   part_of_speech, example, language, review_count, average_score,
		                   mastery, tags, last_reviewed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		ON CONFLICT (user_id, book_id, text) DO NOTHING
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()
	for _, word := range words {
		if word.ID == uuid.Nil {
			word.ID = uuid.New()
		}
		word.CreatedAt = now
		word.UpdatedAt = now

		_, err = stmt.ExecContext(
			ctx,
			word.ID,
			word.UserID,
			word.BookID,
			word.PageNumber,
			word.Text,
			word.Meaning,
			word.Pronunciation,
			word.PartOfSpeech,
			word.Example,
			word.Language,
			word.ReviewCount,
			word.AverageScore,
			word.Mastery,
			pq.Array(word.Tags),
			word.LastReviewedAt,
			word.CreatedAt,
			word.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Ensure wordRepositoryPostgres implements WordRepository interface
var _ WordRepository = (*wordRepositoryPostgres)(nil)
