package ocr

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

var (
	// ErrInvalidCorrectedText is returned when the corrected text is invalid
	ErrInvalidCorrectedText = errors.New("corrected text is invalid")
	// ErrPageNotFound is returned when the page is not found
	ErrPageNotFound = errors.New("page not found")
	// ErrUnauthorized is returned when the user is not authorized
	ErrUnauthorized = errors.New("unauthorized access")
	// ErrTextTooLong is returned when the text exceeds the maximum length
	ErrTextTooLong = errors.New("text exceeds maximum length of 10000 characters")
)

// PageRepository defines methods for page data access
type PageRepository interface {
	GetByID(ctx context.Context, pageID uuid.UUID) (*models.Page, error)
	Update(ctx context.Context, page *models.Page) error
}

// CorrectionRepository defines methods for correction history data access
type CorrectionRepository interface {
	Create(ctx context.Context, correction *models.OCRTextCorrection) error
	GetByPageID(ctx context.Context, pageID uuid.UUID, limit, offset int) ([]models.OCRTextCorrection, error)
	CountByPageID(ctx context.Context, pageID uuid.UUID) (int, error)
}

// EditorService handles OCR text correction operations
type EditorService struct {
	pageRepo       PageRepository
	correctionRepo CorrectionRepository
}

// NewEditorService creates a new EditorService
func NewEditorService(pageRepo PageRepository, correctionRepo CorrectionRepository) *EditorService {
	return &EditorService{
		pageRepo:       pageRepo,
		correctionRepo: correctionRepo,
	}
}

// UpdateOCRText updates the OCR text for a page with manual corrections
func (s *EditorService) UpdateOCRText(
	ctx context.Context,
	bookID, pageID, userID uuid.UUID,
	correctedText string,
) (*models.OCRTextCorrection, error) {
	// Validate corrected text
	if err := s.validateCorrectedText(correctedText); err != nil {
		return nil, err
	}

	// Get the page
	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get page: %w", err)
	}
	if page == nil {
		return nil, ErrPageNotFound
	}

	// Verify the page belongs to the specified book
	if page.BookID != bookID {
		return nil, ErrUnauthorized
	}

	// Create correction record
	correction := &models.OCRTextCorrection{
		ID:            uuid.New(),
		BookID:        bookID,
		PageID:        pageID,
		OriginalText:  page.OCRText,
		CorrectedText: correctedText,
		UserID:        userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save correction to history
	if err := s.correctionRepo.Create(ctx, correction); err != nil {
		return nil, fmt.Errorf("failed to save correction: %w", err)
	}

	// Update the page with corrected text
	// Note: We store corrections separately, but also update the OCRText
	page.OCRText = correctedText
	page.UpdatedAt = time.Now()
	if err := s.pageRepo.Update(ctx, page); err != nil {
		return nil, fmt.Errorf("failed to update page: %w", err)
	}

	return correction, nil
}

// GetCorrectionHistory retrieves the correction history for a page
func (s *EditorService) GetCorrectionHistory(
	ctx context.Context,
	bookID, pageID, userID uuid.UUID,
	limit, offset int,
) (*models.OCRCorrectionHistory, error) {
	// Get the page to verify access
	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get page: %w", err)
	}
	if page == nil {
		return nil, ErrPageNotFound
	}

	// Verify the page belongs to the specified book
	if page.BookID != bookID {
		return nil, ErrUnauthorized
	}

	// Get corrections
	corrections, err := s.correctionRepo.GetByPageID(ctx, pageID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get corrections: %w", err)
	}

	// Get total count
	totalCount, err := s.correctionRepo.CountByPageID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("failed to count corrections: %w", err)
	}

	return &models.OCRCorrectionHistory{
		PageID:      pageID,
		Corrections: corrections,
		TotalCount:  totalCount,
	}, nil
}

// validateCorrectedText validates the corrected text
func (s *EditorService) validateCorrectedText(text string) error {
	if text == "" {
		return fmt.Errorf("%w: text cannot be empty", ErrInvalidCorrectedText)
	}

	// Check character count (not byte count)
	if utf8.RuneCountInString(text) > 10000 {
		return ErrTextTooLong
	}

	return nil
}

// CalculateDiff calculates the difference between original and corrected text
func (s *EditorService) CalculateDiff(original, corrected string) *TextDiff {
	return &TextDiff{
		Original:  original,
		Corrected: corrected,
		HasChanges: original != corrected,
	}
}

// TextDiff represents the difference between two texts
type TextDiff struct {
	Original   string `json:"original"`
	Corrected  string `json:"corrected"`
	HasChanges bool   `json:"has_changes"`
}
