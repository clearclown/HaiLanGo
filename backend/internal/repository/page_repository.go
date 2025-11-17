package repository

import (
	"context"
	"database/sql"

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

// PostgreSQL Implementation

// pageRepositoryPostgres はPostgreSQLベースのページリポジトリ実装
type pageRepositoryPostgres struct {
	db *sql.DB
}

// NewPageRepositoryPostgres は新しいPostgreSQL実装のPageRepositoryを作成する
func NewPageRepositoryPostgres(db *sql.DB) PageRepository {
	return &pageRepositoryPostgres{db: db}
}

// Create はページを作成する
func (r *pageRepositoryPostgres) Create(ctx context.Context, page *models.Page) error {
	query := `
		INSERT INTO pages (id, book_id, page_number, image_url, ocr_text, ocr_confidence,
		                  detected_lang, ocr_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		page.ID,
		page.BookID,
		page.PageNumber,
		page.ImageURL,
		page.OCRText,
		page.OCRConfidence,
		page.DetectedLang,
		page.OCRStatus,
		page.CreatedAt,
		page.UpdatedAt,
	)

	return err
}

// Update はページを更新する
func (r *pageRepositoryPostgres) Update(ctx context.Context, page *models.Page) error {
	query := `
		UPDATE pages
		SET image_url = $1, ocr_text = $2, ocr_confidence = $3,
		    detected_lang = $4, ocr_status = $5, updated_at = NOW()
		WHERE id = $6
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		page.ImageURL,
		page.OCRText,
		page.OCRConfidence,
		page.DetectedLang,
		page.OCRStatus,
		page.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrBookNotFound // ページが見つからない
	}

	return nil
}

// FindByID はIDでページを取得する
func (r *pageRepositoryPostgres) FindByID(ctx context.Context, id uuid.UUID) (*models.Page, error) {
	query := `
		SELECT id, book_id, page_number, image_url, ocr_text, ocr_confidence,
		       detected_lang, ocr_status, created_at, updated_at
		FROM pages
		WHERE id = $1
	`

	page := &models.Page{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&page.ID,
		&page.BookID,
		&page.PageNumber,
		&page.ImageURL,
		&page.OCRText,
		&page.OCRConfidence,
		&page.DetectedLang,
		&page.OCRStatus,
		&page.CreatedAt,
		&page.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return page, nil
}

// FindByBookID は書籍IDで全ページを取得する
func (r *pageRepositoryPostgres) FindByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.Page, error) {
	query := `
		SELECT id, book_id, page_number, image_url, ocr_text, ocr_confidence,
		       detected_lang, ocr_status, created_at, updated_at
		FROM pages
		WHERE book_id = $1
		ORDER BY page_number ASC
	`

	rows, err := r.db.QueryContext(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []*models.Page
	for rows.Next() {
		page := &models.Page{}
		err := rows.Scan(
			&page.ID,
			&page.BookID,
			&page.PageNumber,
			&page.ImageURL,
			&page.OCRText,
			&page.OCRConfidence,
			&page.DetectedLang,
			&page.OCRStatus,
			&page.CreatedAt,
			&page.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pages, nil
}

// UpdateOCRResult はOCR結果を更新する
func (r *pageRepositoryPostgres) UpdateOCRResult(ctx context.Context, pageID uuid.UUID, ocrText string, confidence float64, detectedLang string) error {
	query := `
		UPDATE pages
		SET ocr_text = $1, ocr_confidence = $2, detected_lang = $3,
		    ocr_status = $4, updated_at = NOW()
		WHERE id = $5
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		ocrText,
		confidence,
		detectedLang,
		models.OCRStatusCompleted,
		pageID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrBookNotFound // ページが見つからない
	}

	return nil
}
