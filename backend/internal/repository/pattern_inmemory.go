package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/service/pattern"
	"github.com/google/uuid"
)

// PatternRepositoryInterface はパターンリポジトリのインターフェース
type PatternRepositoryInterface interface {
	// ExtractPatterns は書籍からパターンを抽出
	ExtractPatterns(ctx context.Context, bookID uuid.UUID, pages []pattern.PageText, minFrequency int) ([]models.Pattern, error)

	// GetPatternsByBookID は書籍IDでパターンを取得
	GetPatternsByBookID(ctx context.Context, bookID uuid.UUID) ([]models.Pattern, error)

	// GetPatternByID はIDでパターンを取得
	GetPatternByID(ctx context.Context, patternID uuid.UUID) (*models.Pattern, error)

	// GetPatternExamples はパターンの使用例を取得
	GetPatternExamples(ctx context.Context, patternID uuid.UUID, limit int) ([]models.PatternExample, error)

	// GetPatternPractice はパターンの練習問題を取得
	GetPatternPractice(ctx context.Context, patternID uuid.UUID, count int) ([]models.PatternPractice, error)

	// UpdatePatternProgress はユーザーのパターン学習進捗を更新
	UpdatePatternProgress(ctx context.Context, userID uuid.UUID, patternID uuid.UUID, correct bool) (*models.PatternProgress, error)

	// GetUserPatternProgress はユーザーのパターン学習進捗を取得
	GetUserPatternProgress(ctx context.Context, userID uuid.UUID, patternID uuid.UUID) (*models.PatternProgress, error)
}

// InMemoryPatternRepository はインメモリパターンリポジトリ
type InMemoryPatternRepository struct {
	mu              sync.RWMutex
	patterns        map[uuid.UUID]*models.Pattern           // PatternID -> Pattern
	patternsByBook  map[uuid.UUID][]uuid.UUID               // BookID -> PatternIDs
	examples        map[uuid.UUID][]models.PatternExample   // PatternID -> Examples
	practices       map[uuid.UUID][]models.PatternPractice  // PatternID -> Practices
	userProgress    map[string]*models.PatternProgress      // UserID:PatternID -> Progress
	extractor       *pattern.Extractor
}

// NewInMemoryPatternRepository はインメモリパターンリポジトリを作成
func NewInMemoryPatternRepository() *InMemoryPatternRepository {
	repo := &InMemoryPatternRepository{
		patterns:        make(map[uuid.UUID]*models.Pattern),
		patternsByBook:  make(map[uuid.UUID][]uuid.UUID),
		examples:        make(map[uuid.UUID][]models.PatternExample),
		practices:       make(map[uuid.UUID][]models.PatternPractice),
		userProgress:    make(map[string]*models.PatternProgress),
		extractor:       pattern.NewExtractor(),
	}

	// サンプルデータを初期化
	repo.initSampleData()

	return repo
}

func (r *InMemoryPatternRepository) initSampleData() {
	// サンプル書籍ID
	sampleBookID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	// サンプルパターン1: 挨拶
	pattern1 := &models.Pattern{
		ID:          uuid.New(),
		BookID:      sampleBookID,
		Type:        models.PatternTypeGreeting,
		Pattern:     "Здравствуйте!",
		Translation: "こんにちは！",
		Frequency:   5,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	r.patterns[pattern1.ID] = pattern1
	r.patternsByBook[sampleBookID] = append(r.patternsByBook[sampleBookID], pattern1.ID)

	// サンプルパターン2: 質問
	pattern2 := &models.Pattern{
		ID:          uuid.New(),
		BookID:      sampleBookID,
		Type:        models.PatternTypeQuestion,
		Pattern:     "Как дела?",
		Translation: "調子はどう？",
		Frequency:   3,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	r.patterns[pattern2.ID] = pattern2
	r.patternsByBook[sampleBookID] = append(r.patternsByBook[sampleBookID], pattern2.ID)

	// サンプル例文
	r.examples[pattern1.ID] = []models.PatternExample{
		{
			ID:             uuid.New(),
			PatternID:      pattern1.ID,
			PageNumber:     1,
			OriginalText:   "Здравствуйте! Меня зовут Иван.",
			TranslatedText: "こんにちは！私の名前はイワンです。",
			Context:        "会話の始まり",
			CreatedAt:      time.Now(),
		},
	}

	// サンプル練習問題
	r.practices[pattern1.ID] = []models.PatternPractice{
		{
			ID:                 uuid.New(),
			PatternID:          pattern1.ID,
			Question:           "「こんにちは」をロシア語で言ってください。",
			CorrectAnswer:      "Здравствуйте",
			AlternativeAnswers: []string{"Привет", "Добрый день"},
			Difficulty:         1,
			CreatedAt:          time.Now(),
		},
	}
}

func (r *InMemoryPatternRepository) ExtractPatterns(ctx context.Context, bookID uuid.UUID, pages []pattern.PageText, minFrequency int) ([]models.Pattern, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// パターン抽出サービスを使用
	extractedPatterns, err := r.extractor.ExtractPatterns(ctx, bookID, pages, minFrequency)
	if err != nil {
		return nil, err
	}

	// 抽出されたパターンを保存
	for i := range extractedPatterns {
		p := &extractedPatterns[i]
		r.patterns[p.ID] = p
		r.patternsByBook[bookID] = append(r.patternsByBook[bookID], p.ID)
	}

	return extractedPatterns, nil
}

func (r *InMemoryPatternRepository) GetPatternsByBookID(ctx context.Context, bookID uuid.UUID) ([]models.Pattern, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	patternIDs, exists := r.patternsByBook[bookID]
	if !exists {
		return []models.Pattern{}, nil
	}

	patterns := make([]models.Pattern, 0, len(patternIDs))
	for _, patternID := range patternIDs {
		if pattern, exists := r.patterns[patternID]; exists {
			patterns = append(patterns, *pattern)
		}
	}

	return patterns, nil
}

func (r *InMemoryPatternRepository) GetPatternByID(ctx context.Context, patternID uuid.UUID) (*models.Pattern, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pattern, exists := r.patterns[patternID]
	if !exists {
		return nil, nil
	}

	return pattern, nil
}

func (r *InMemoryPatternRepository) GetPatternExamples(ctx context.Context, patternID uuid.UUID, limit int) ([]models.PatternExample, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	examples, exists := r.examples[patternID]
	if !exists {
		return []models.PatternExample{}, nil
	}

	if limit > 0 && limit < len(examples) {
		return examples[:limit], nil
	}

	return examples, nil
}

func (r *InMemoryPatternRepository) GetPatternPractice(ctx context.Context, patternID uuid.UUID, count int) ([]models.PatternPractice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	practices, exists := r.practices[patternID]
	if !exists {
		return []models.PatternPractice{}, nil
	}

	if count > 0 && count < len(practices) {
		return practices[:count], nil
	}

	return practices, nil
}

func (r *InMemoryPatternRepository) UpdatePatternProgress(ctx context.Context, userID uuid.UUID, patternID uuid.UUID, correct bool) (*models.PatternProgress, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := userID.String() + ":" + patternID.String()
	progress, exists := r.userProgress[key]

	if !exists {
		now := time.Now()
		progress = &models.PatternProgress{
			ID:              uuid.New(),
			UserID:          userID,
			PatternID:       patternID,
			MasteryLevel:    0,
			PracticeCount:   0,
			CorrectCount:    0,
			LastPracticedAt: &now,
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		r.userProgress[key] = progress
	}

	// 進捗を更新
	progress.PracticeCount++
	if correct {
		progress.CorrectCount++
	}

	// 習熟度を計算 (正答率ベース)
	if progress.PracticeCount > 0 {
		progress.MasteryLevel = (progress.CorrectCount * 100) / progress.PracticeCount
	}

	now := time.Now()
	progress.LastPracticedAt = &now
	progress.UpdatedAt = now

	return progress, nil
}

func (r *InMemoryPatternRepository) GetUserPatternProgress(ctx context.Context, userID uuid.UUID, patternID uuid.UUID) (*models.PatternProgress, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := userID.String() + ":" + patternID.String()
	progress, exists := r.userProgress[key]
	if !exists {
		return nil, nil
	}

	return progress, nil
}

// PostgreSQL Implementation

type PatternRepositoryPostgres struct {
	db        *sql.DB
	extractor *pattern.Extractor
}

func NewPatternRepositoryPostgres(db *sql.DB) PatternRepositoryInterface {
	return &PatternRepositoryPostgres{
		db:        db,
		extractor: pattern.NewExtractor(),
	}
}

func (r *PatternRepositoryPostgres) ExtractPatterns(ctx context.Context, bookID uuid.UUID, pages []pattern.PageText, minFrequency int) ([]models.Pattern, error) {
	// パターン抽出サービスを使用
	extractedPatterns, err := r.extractor.ExtractPatterns(ctx, bookID, pages, minFrequency)
	if err != nil {
		return nil, err
	}

	// 抽出されたパターンをデータベースに保存
	for i := range extractedPatterns {
		p := &extractedPatterns[i]
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO patterns (id, book_id, type, pattern, translation, frequency, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (id) DO UPDATE SET
				frequency = patterns.frequency + EXCLUDED.frequency,
				updated_at = EXCLUDED.updated_at
		`, p.ID, p.BookID, p.Type, p.Pattern, p.Translation, p.Frequency, p.CreatedAt, p.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}

	return extractedPatterns, nil
}

func (r *PatternRepositoryPostgres) GetPatternsByBookID(ctx context.Context, bookID uuid.UUID) ([]models.Pattern, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, book_id, type, pattern, translation, frequency, created_at, updated_at
		FROM patterns
		WHERE book_id = $1
		ORDER BY frequency DESC, created_at DESC
	`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	patterns := []models.Pattern{}
	for rows.Next() {
		var p models.Pattern
		err := rows.Scan(&p.ID, &p.BookID, &p.Type, &p.Pattern, &p.Translation, &p.Frequency, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			continue
		}
		patterns = append(patterns, p)
	}

	return patterns, nil
}

func (r *PatternRepositoryPostgres) GetPatternByID(ctx context.Context, patternID uuid.UUID) (*models.Pattern, error) {
	var p models.Pattern
	err := r.db.QueryRowContext(ctx, `
		SELECT id, book_id, type, pattern, translation, frequency, created_at, updated_at
		FROM patterns
		WHERE id = $1
	`, patternID).Scan(&p.ID, &p.BookID, &p.Type, &p.Pattern, &p.Translation, &p.Frequency, &p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PatternRepositoryPostgres) GetPatternExamples(ctx context.Context, patternID uuid.UUID, limit int) ([]models.PatternExample, error) {
	query := `
		SELECT id, pattern_id, page_number, original_text, translated_text, context, created_at
		FROM pattern_examples
		WHERE pattern_id = $1
		ORDER BY page_number ASC
	`
	args := []interface{}{patternID}

	if limit > 0 {
		query += " LIMIT $2"
		args = append(args, limit)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	examples := []models.PatternExample{}
	for rows.Next() {
		var ex models.PatternExample
		var context sql.NullString
		err := rows.Scan(&ex.ID, &ex.PatternID, &ex.PageNumber, &ex.OriginalText, &ex.TranslatedText, &context, &ex.CreatedAt)
		if err != nil {
			continue
		}
		if context.Valid {
			ex.Context = context.String
		}
		examples = append(examples, ex)
	}

	return examples, nil
}

func (r *PatternRepositoryPostgres) GetPatternPractice(ctx context.Context, patternID uuid.UUID, count int) ([]models.PatternPractice, error) {
	query := `
		SELECT id, pattern_id, question, correct_answer, alternative_answers, difficulty, created_at
		FROM pattern_practices
		WHERE pattern_id = $1
		ORDER BY difficulty ASC, created_at ASC
	`
	args := []interface{}{patternID}

	if count > 0 {
		query += " LIMIT $2"
		args = append(args, count)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	practices := []models.PatternPractice{}
	for rows.Next() {
		var pr models.PatternPractice
		var alternativesJSON []byte
		err := rows.Scan(&pr.ID, &pr.PatternID, &pr.Question, &pr.CorrectAnswer, &alternativesJSON, &pr.Difficulty, &pr.CreatedAt)
		if err != nil {
			continue
		}

		// JSONB配列をデコード
		if alternativesJSON != nil {
			json.Unmarshal(alternativesJSON, &pr.AlternativeAnswers)
		}

		practices = append(practices, pr)
	}

	return practices, nil
}

func (r *PatternRepositoryPostgres) UpdatePatternProgress(ctx context.Context, userID uuid.UUID, patternID uuid.UUID, correct bool) (*models.PatternProgress, error) {
	// 既存の進捗を取得
	var progress models.PatternProgress
	var lastPracticed sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, pattern_id, mastery_level, practice_count, correct_count, last_practiced_at, created_at, updated_at
		FROM pattern_progress
		WHERE user_id = $1 AND pattern_id = $2
	`, userID, patternID).Scan(&progress.ID, &progress.UserID, &progress.PatternID, &progress.MasteryLevel,
		&progress.PracticeCount, &progress.CorrectCount, &lastPracticed, &progress.CreatedAt, &progress.UpdatedAt)

	now := time.Now()

	if err == sql.ErrNoRows {
		// 新規作成
		progress = models.PatternProgress{
			ID:              uuid.New(),
			UserID:          userID,
			PatternID:       patternID,
			MasteryLevel:    0,
			PracticeCount:   1,
			CorrectCount:    0,
			LastPracticedAt: &now,
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		if correct {
			progress.CorrectCount = 1
		}

		// 習熟度を計算
		if progress.PracticeCount > 0 {
			progress.MasteryLevel = (progress.CorrectCount * 100) / progress.PracticeCount
		}

		_, err = r.db.ExecContext(ctx, `
			INSERT INTO pattern_progress (id, user_id, pattern_id, mastery_level, practice_count, correct_count, last_practiced_at, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`, progress.ID, progress.UserID, progress.PatternID, progress.MasteryLevel, progress.PracticeCount,
			progress.CorrectCount, progress.LastPracticedAt, progress.CreatedAt, progress.UpdatedAt)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// 既存の進捗を更新
		if lastPracticed.Valid {
			progress.LastPracticedAt = &lastPracticed.Time
		}

		progress.PracticeCount++
		if correct {
			progress.CorrectCount++
		}

		// 習熟度を計算
		if progress.PracticeCount > 0 {
			progress.MasteryLevel = (progress.CorrectCount * 100) / progress.PracticeCount
		}

		progress.LastPracticedAt = &now
		progress.UpdatedAt = now

		_, err = r.db.ExecContext(ctx, `
			UPDATE pattern_progress
			SET mastery_level = $1, practice_count = $2, correct_count = $3, last_practiced_at = $4, updated_at = $5
			WHERE id = $6
		`, progress.MasteryLevel, progress.PracticeCount, progress.CorrectCount, progress.LastPracticedAt, progress.UpdatedAt, progress.ID)
		if err != nil {
			return nil, err
		}
	}

	return &progress, nil
}

func (r *PatternRepositoryPostgres) GetUserPatternProgress(ctx context.Context, userID uuid.UUID, patternID uuid.UUID) (*models.PatternProgress, error) {
	var progress models.PatternProgress
	var lastPracticed sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, pattern_id, mastery_level, practice_count, correct_count, last_practiced_at, created_at, updated_at
		FROM pattern_progress
		WHERE user_id = $1 AND pattern_id = $2
	`, userID, patternID).Scan(&progress.ID, &progress.UserID, &progress.PatternID, &progress.MasteryLevel,
		&progress.PracticeCount, &progress.CorrectCount, &lastPracticed, &progress.CreatedAt, &progress.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if lastPracticed.Valid {
		progress.LastPracticedAt = &lastPracticed.Time
	}

	return &progress, nil
}
