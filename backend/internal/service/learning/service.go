package learning

import (
	"context"
	"errors"

	"github.com/clearclown/HaiLanGo/internal/models"
	"github.com/google/uuid"
)

var (
	// ErrPageNotFound はページが見つからない場合のエラー
	ErrPageNotFound = errors.New("page not found")
	// ErrInvalidPageNumber はページ番号が不正な場合のエラー
	ErrInvalidPageNumber = errors.New("invalid page number")
)

// Repository は学習データのリポジトリインターフェース
type Repository interface {
	GetPage(ctx context.Context, bookID uuid.UUID, pageNumber int) (*models.PageWithProgress, error)
	MarkPageCompleted(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, studyTime int) error
	GetProgress(ctx context.Context, userID, bookID uuid.UUID) (*models.LearningProgress, error)
}

// Service は学習機能を提供するサービス
type Service struct {
	repo Repository
}

// NewService は新しいServiceを作成する
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetPage は指定されたページを取得する
func (s *Service) GetPage(ctx context.Context, bookID uuid.UUID, pageNumber int) (*models.PageWithProgress, error) {
	if pageNumber < 1 {
		return nil, ErrInvalidPageNumber
	}

	page, err := s.repo.GetPage(ctx, bookID, pageNumber)
	if err != nil {
		return nil, err
	}

	return page, nil
}

// MarkPageCompleted はページを完了としてマークする
func (s *Service) MarkPageCompleted(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, studyTime int) error {
	if pageNumber < 1 {
		return ErrInvalidPageNumber
	}

	if studyTime < 0 {
		studyTime = 0
	}

	return s.repo.MarkPageCompleted(ctx, userID, bookID, pageNumber, studyTime)
}

// GetProgress は書籍の学習進捗を取得する
func (s *Service) GetProgress(ctx context.Context, userID, bookID uuid.UUID) (*models.LearningProgress, error) {
	progress, err := s.repo.GetProgress(ctx, userID, bookID)
	if err != nil {
		return nil, err
	}

	// パーセンテージを計算
	if progress.TotalPages > 0 {
		progress.Progress = float64(progress.CompletedPages) / float64(progress.TotalPages) * 100
	}

	return progress, nil
}
