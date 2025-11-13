package ocr

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// Mock implementations

type mockPageRepository struct {
	pages map[uuid.UUID]*models.Page
	err   error
}

func newMockPageRepository() *mockPageRepository {
	return &mockPageRepository{
		pages: make(map[uuid.UUID]*models.Page),
	}
}

func (m *mockPageRepository) GetByID(ctx context.Context, pageID uuid.UUID) (*models.Page, error) {
	if m.err != nil {
		return nil, m.err
	}
	page, exists := m.pages[pageID]
	if !exists {
		return nil, nil
	}
	return page, nil
}

func (m *mockPageRepository) Update(ctx context.Context, page *models.Page) error {
	if m.err != nil {
		return m.err
	}
	m.pages[page.ID] = page
	return nil
}

type mockCorrectionRepository struct {
	corrections map[uuid.UUID][]models.OCRTextCorrection
	err         error
}

func newMockCorrectionRepository() *mockCorrectionRepository {
	return &mockCorrectionRepository{
		corrections: make(map[uuid.UUID][]models.OCRTextCorrection),
	}
}

func (m *mockCorrectionRepository) Create(ctx context.Context, correction *models.OCRTextCorrection) error {
	if m.err != nil {
		return m.err
	}
	m.corrections[correction.PageID] = append(m.corrections[correction.PageID], *correction)
	return nil
}

func (m *mockCorrectionRepository) GetByPageID(ctx context.Context, pageID uuid.UUID, limit, offset int) ([]models.OCRTextCorrection, error) {
	if m.err != nil {
		return nil, m.err
	}
	corrections, exists := m.corrections[pageID]
	if !exists {
		return []models.OCRTextCorrection{}, nil
	}

	start := offset
	if start > len(corrections) {
		return []models.OCRTextCorrection{}, nil
	}

	end := start + limit
	if end > len(corrections) {
		end = len(corrections)
	}

	return corrections[start:end], nil
}

func (m *mockCorrectionRepository) CountByPageID(ctx context.Context, pageID uuid.UUID) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	corrections, exists := m.corrections[pageID]
	if !exists {
		return 0, nil
	}
	return len(corrections), nil
}

// Tests

func TestUpdateOCRText_Success(t *testing.T) {
	pageRepo := newMockPageRepository()
	correctionRepo := newMockCorrectionRepository()
	service := NewEditorService(pageRepo, correctionRepo)

	bookID := uuid.New()
	pageID := uuid.New()
	userID := uuid.New()

	// Setup test page
	pageRepo.pages[pageID] = &models.Page{
		ID:          pageID,
		BookID:      bookID,
		PageNumber:  1,
		OCRText:     "Original text",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	correctedText := "Corrected text"

	// Test
	correction, err := service.UpdateOCRText(context.Background(), bookID, pageID, userID, correctedText)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if correction == nil {
		t.Fatal("Expected correction to be non-nil")
	}

	if correction.OriginalText != "Original text" {
		t.Errorf("Expected original text 'Original text', got: %s", correction.OriginalText)
	}

	if correction.CorrectedText != correctedText {
		t.Errorf("Expected corrected text '%s', got: %s", correctedText, correction.CorrectedText)
	}

	// Verify page was updated
	updatedPage := pageRepo.pages[pageID]
	if updatedPage.CorrectedText == nil || *updatedPage.CorrectedText != correctedText {
		t.Errorf("Expected page corrected text '%s', got: %v", correctedText, updatedPage.CorrectedText)
	}
}

func TestUpdateOCRText_InvalidText(t *testing.T) {
	pageRepo := newMockPageRepository()
	correctionRepo := newMockCorrectionRepository()
	service := NewEditorService(pageRepo, correctionRepo)

	bookID := uuid.New()
	pageID := uuid.New()
	userID := uuid.New()

	testCases := []struct {
		name          string
		correctedText string
		expectedError error
	}{
		{
			name:          "empty text",
			correctedText: "",
			expectedError: ErrInvalidCorrectedText,
		},
		{
			name:          "text too long",
			correctedText: string(make([]rune, 10001)),
			expectedError: ErrTextTooLong,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.UpdateOCRText(context.Background(), bookID, pageID, userID, tc.correctedText)

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("Expected error %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestUpdateOCRText_PageNotFound(t *testing.T) {
	pageRepo := newMockPageRepository()
	correctionRepo := newMockCorrectionRepository()
	service := NewEditorService(pageRepo, correctionRepo)

	bookID := uuid.New()
	pageID := uuid.New()
	userID := uuid.New()

	_, err := service.UpdateOCRText(context.Background(), bookID, pageID, userID, "Corrected text")

	if !errors.Is(err, ErrPageNotFound) {
		t.Errorf("Expected error %v, got: %v", ErrPageNotFound, err)
	}
}

func TestUpdateOCRText_Unauthorized(t *testing.T) {
	pageRepo := newMockPageRepository()
	correctionRepo := newMockCorrectionRepository()
	service := NewEditorService(pageRepo, correctionRepo)

	bookID := uuid.New()
	wrongBookID := uuid.New()
	pageID := uuid.New()
	userID := uuid.New()

	// Setup test page with different book ID
	pageRepo.pages[pageID] = &models.Page{
		ID:         pageID,
		BookID:     bookID,
		PageNumber: 1,
		OCRText:    "Original text",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err := service.UpdateOCRText(context.Background(), wrongBookID, pageID, userID, "Corrected text")

	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("Expected error %v, got: %v", ErrUnauthorized, err)
	}
}

func TestGetCorrectionHistory_Success(t *testing.T) {
	pageRepo := newMockPageRepository()
	correctionRepo := newMockCorrectionRepository()
	service := NewEditorService(pageRepo, correctionRepo)

	bookID := uuid.New()
	pageID := uuid.New()
	userID := uuid.New()

	// Setup test page
	pageRepo.pages[pageID] = &models.Page{
		ID:         pageID,
		BookID:     bookID,
		PageNumber: 1,
		OCRText:    "Original text",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Add some corrections
	for i := 0; i < 5; i++ {
		correctionRepo.corrections[pageID] = append(correctionRepo.corrections[pageID], models.OCRTextCorrection{
			ID:            uuid.New(),
			BookID:        bookID,
			PageID:        pageID,
			OriginalText:  "Original text",
			CorrectedText: "Corrected text",
			UserID:        userID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
	}

	history, err := service.GetCorrectionHistory(context.Background(), bookID, pageID, userID, 10, 0)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if history == nil {
		t.Fatal("Expected history to be non-nil")
	}

	if history.TotalCount != 5 {
		t.Errorf("Expected total count 5, got: %d", history.TotalCount)
	}

	if len(history.Corrections) != 5 {
		t.Errorf("Expected 5 corrections, got: %d", len(history.Corrections))
	}
}

func TestCalculateDiff(t *testing.T) {
	service := NewEditorService(nil, nil)

	testCases := []struct {
		name          string
		original      string
		corrected     string
		expectChanges bool
	}{
		{
			name:          "no changes",
			original:      "Same text",
			corrected:     "Same text",
			expectChanges: false,
		},
		{
			name:          "with changes",
			original:      "Original text",
			corrected:     "Corrected text",
			expectChanges: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			diff := service.CalculateDiff(tc.original, tc.corrected)

			if diff.HasChanges != tc.expectChanges {
				t.Errorf("Expected HasChanges %v, got: %v", tc.expectChanges, diff.HasChanges)
			}

			if diff.Original != tc.original {
				t.Errorf("Expected original '%s', got: '%s'", tc.original, diff.Original)
			}

			if diff.Corrected != tc.corrected {
				t.Errorf("Expected corrected '%s', got: '%s'", tc.corrected, diff.Corrected)
			}
		})
	}
}
