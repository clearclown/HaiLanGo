package vocabulary

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/clearclown/HaiLanGo/backend/pkg/vocabulary"
)

// VocabularyService は単語帳サービスのインターフェース
type VocabularyService interface {
	AutoCollectWords(ctx context.Context, userID, bookID string, pageNumber int, text, language string) error
	AddWord(ctx context.Context, word *models.Word) error
	GetWords(ctx context.Context, filter *models.WordFilter) ([]*models.Word, error)
	GetWordByID(ctx context.Context, id string) (*models.Word, error)
	UpdateWord(ctx context.Context, word *models.Word) error
	DeleteWord(ctx context.Context, id string) error
	RecordReview(ctx context.Context, wordID string, score float64) error
	GetStats(ctx context.Context, userID, bookID string) (*models.WordStats, error)
	ExportWordsToCSV(ctx context.Context, filter *models.WordFilter) ([]byte, error)
	AddTags(ctx context.Context, wordID string, tags []string) error
}

// vocabularyService は単語帳サービスの実装
type vocabularyService struct {
	repo repository.WordRepository
}

// NewVocabularyService は新しい単語帳サービスを作成する
func NewVocabularyService(repo repository.WordRepository) VocabularyService {
	return &vocabularyService{
		repo: repo,
	}
}

// NewMockVocabularyService はモック単語帳サービスを作成する
func NewMockVocabularyService() VocabularyService {
	return &vocabularyService{
		repo: repository.NewMockWordRepository(),
	}
}

// AutoCollectWords はテキストから単語を自動収集する
func (s *vocabularyService) AutoCollectWords(ctx context.Context, userID, bookID string, pageNumber int, text, language string) error {
	// 単語を抽出
	words := vocabulary.ExtractWords(text, language)

	// 各単語を保存
	for _, wordText := range words {
		word := &models.Word{
			UserID:     userID,
			BookID:     bookID,
			PageNumber: pageNumber,
			Text:       wordText,
			Language:   language,
			// 意味は後で辞書APIで取得する（今は空）
			Meaning:  "",
			Mastery:  0.0,
			Tags:     []string{},
		}

		// 重複チェックのため、既存単語を確認
		filter := &models.WordFilter{
			UserID: userID,
			BookID: bookID,
			Query:  wordText,
		}
		existing, _, err := s.repo.List(ctx, filter)
		if err != nil {
			return fmt.Errorf("failed to check existing word: %w", err)
		}

		// 既に存在する場合はスキップ
		if len(existing) > 0 {
			continue
		}

		// 単語を作成
		if err := s.repo.Create(ctx, word); err != nil {
			// 重複エラーは無視
			if err != repository.ErrWordAlreadyExists {
				return fmt.Errorf("failed to create word: %w", err)
			}
		}
	}

	return nil
}

// AddWord は単語を追加する
func (s *vocabularyService) AddWord(ctx context.Context, word *models.Word) error {
	if err := s.repo.Create(ctx, word); err != nil {
		return fmt.Errorf("failed to add word: %w", err)
	}
	return nil
}

// GetWords はフィルタ条件に基づいて単語一覧を取得する
func (s *vocabularyService) GetWords(ctx context.Context, filter *models.WordFilter) ([]*models.Word, error) {
	words, _, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get words: %w", err)
	}
	return words, nil
}

// GetWordByID はIDで単語を取得する
func (s *vocabularyService) GetWordByID(ctx context.Context, id string) (*models.Word, error) {
	word, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get word: %w", err)
	}
	return word, nil
}

// UpdateWord は単語を更新する
func (s *vocabularyService) UpdateWord(ctx context.Context, word *models.Word) error {
	if err := s.repo.Update(ctx, word); err != nil {
		return fmt.Errorf("failed to update word: %w", err)
	}
	return nil
}

// DeleteWord は単語を削除する
func (s *vocabularyService) DeleteWord(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete word: %w", err)
	}
	return nil
}

// RecordReview は学習記録を保存し、習得度を更新する
func (s *vocabularyService) RecordReview(ctx context.Context, wordID string, score float64) error {
	// 単語を取得
	word, err := s.repo.GetByID(ctx, wordID)
	if err != nil {
		return fmt.Errorf("failed to get word: %w", err)
	}

	// 平均スコアを更新
	totalScore := word.AverageScore*float64(word.ReviewCount) + score
	word.ReviewCount++
	word.AverageScore = totalScore / float64(word.ReviewCount)

	// 習得度を計算
	word.Mastery = CalculateMastery(word.ReviewCount, word.AverageScore)

	// 最終学習日時を更新
	word.LastReviewedAt = time.Now()

	// 更新
	if err := s.repo.Update(ctx, word); err != nil {
		return fmt.Errorf("failed to update word: %w", err)
	}

	return nil
}

// CalculateMastery は習得度を計算する
func CalculateMastery(reviewCount int, averageScore float64) float64 {
	// 学習回数と平均スコアから習得度を計算
	// 学習回数が増えるほど習得度が上がる
	// 平均スコアが高いほど習得度が上がる

	// 学習回数係数（最大10回まで）
	reviewFactor := math.Min(float64(reviewCount)/10.0, 1.0)

	// 習得度 = 平均スコア × 学習回数係数
	mastery := averageScore * reviewFactor

	// 0-100の範囲に制限
	if mastery > 100.0 {
		mastery = 100.0
	}
	if mastery < 0.0 {
		mastery = 0.0
	}

	return mastery
}

// GetStats は単語統計を取得する
func (s *vocabularyService) GetStats(ctx context.Context, userID, bookID string) (*models.WordStats, error) {
	stats, err := s.repo.GetStats(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}
	return stats, nil
}

// ExportWordsToCSV は単語をCSV形式でエクスポートする
func (s *vocabularyService) ExportWordsToCSV(ctx context.Context, filter *models.WordFilter) ([]byte, error) {
	// 単語を取得
	words, _, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get words: %w", err)
	}

	// CSVを生成
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// ヘッダー
	header := []string{"単語", "意味", "言語", "品詞", "発音", "例文", "習得度", "学習回数", "最終学習日"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write csv header: %w", err)
	}

	// データ
	for _, word := range words {
		lastReviewedAt := ""
		if !word.LastReviewedAt.IsZero() {
			lastReviewedAt = word.LastReviewedAt.Format("2006-01-02 15:04:05")
		}

		row := []string{
			word.Text,
			word.Meaning,
			word.Language,
			word.PartOfSpeech,
			word.Pronunciation,
			word.Example,
			fmt.Sprintf("%.1f%%", word.Mastery),
			fmt.Sprintf("%d", word.ReviewCount),
			lastReviewedAt,
		}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("failed to write csv row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("failed to flush csv writer: %w", err)
	}

	return buf.Bytes(), nil
}

// AddTags は単語にタグを追加する
func (s *vocabularyService) AddTags(ctx context.Context, wordID string, tags []string) error {
	// 単語を取得
	word, err := s.repo.GetByID(ctx, wordID)
	if err != nil {
		return fmt.Errorf("failed to get word: %w", err)
	}

	// タグを追加（重複を除外）
	existingTags := make(map[string]bool)
	for _, tag := range word.Tags {
		existingTags[tag] = true
	}

	for _, tag := range tags {
		if !existingTags[tag] {
			word.Tags = append(word.Tags, tag)
			existingTags[tag] = true
		}
	}

	// 更新
	if err := s.repo.Update(ctx, word); err != nil {
		return fmt.Errorf("failed to update word: %w", err)
	}

	return nil
}
