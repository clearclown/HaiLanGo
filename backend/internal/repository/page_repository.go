package repository

import (
	"context"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// PageRepository はページのリポジトリインターフェース
type PageRepository interface {
	// Create はページを作成する
	Create(ctx context.Context, page *models.Page) error

	// Update はページを更新する
	Update(ctx context.Context, page *models.Page) error

	// FindByID はIDでページを取得する
	FindByID(ctx context.Context, id uuid.UUID) (*models.Page, error)

	// FindByBookID は書籍IDで全ページを取得する
	FindByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.Page, error)

	// UpdateOCRResult はOCR結果を更新する
	UpdateOCRResult(ctx context.Context, pageID uuid.UUID, ocrText string, confidence float64, detectedLang string) error
}

// MockPageRepository はモックページリポジトリ
type MockPageRepository struct {
	pages map[uuid.UUID]*models.Page
}

// NewMockPageRepository は新しいモックページリポジトリを作成する
func NewMockPageRepository() *MockPageRepository {
	return &MockPageRepository{
		pages: make(map[uuid.UUID]*models.Page),
	}
}

// Create はページを作成する
func (r *MockPageRepository) Create(ctx context.Context, page *models.Page) error {
	r.pages[page.ID] = page
	return nil
}

// Update はページを更新する
func (r *MockPageRepository) Update(ctx context.Context, page *models.Page) error {
	r.pages[page.ID] = page
	return nil
}

// FindByID はIDでページを取得する
func (r *MockPageRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Page, error) {
	page, ok := r.pages[id]
	if !ok {
		return nil, nil
	}
	return page, nil
}

// FindByBookID は書籍IDで全ページを取得する
func (r *MockPageRepository) FindByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.Page, error) {
	var pages []*models.Page
	for _, page := range r.pages {
		if page.BookID == bookID {
			pages = append(pages, page)
		}
	}
	return pages, nil
}

// UpdateOCRResult はOCR結果を更新する
func (r *MockPageRepository) UpdateOCRResult(ctx context.Context, pageID uuid.UUID, ocrText string, confidence float64, detectedLang string) error {
	page, ok := r.pages[pageID]
	if !ok {
		return nil
	}

	page.OCRText = ocrText
	page.OCRConfidence = confidence
	page.DetectedLang = detectedLang
	page.OCRStatus = models.OCRStatusCompleted

	return nil
}
