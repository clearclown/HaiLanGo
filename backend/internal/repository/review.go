package repository

import (
	"context"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

var (
	ErrReviewItemNotFound = &RepositoryError{
		Code:    "REVIEW_ITEM_NOT_FOUND",
		Message: "review item not found",
	}
)

type ReviewRepository interface {
	Create(ctx context.Context, item *models.ReviewItem) error
	FindByID(ctx context.Context, id string) (*models.ReviewItem, error)
	FindByUserID(ctx context.Context, userID string) ([]*models.ReviewItem, error)
	Update(ctx context.Context, item *models.ReviewItem) error
	Delete(ctx context.Context, id string) error

	// 統計用
	CountCompletedToday(ctx context.Context, userID string, since time.Time) (int, error)
	CountCompletedSince(ctx context.Context, userID string, since time.Time) (int, error)

	// 履歴
	SaveHistory(ctx context.Context, history *models.ReviewHistory) error
}
