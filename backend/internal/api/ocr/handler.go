package ocr

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/service/ocr"
	"github.com/google/uuid"
)

// Handler handles OCR-related HTTP requests
type Handler struct {
	editorService *ocr.EditorService
}

// NewHandler creates a new OCR Handler
func NewHandler(editorService *ocr.EditorService) *Handler {
	return &Handler{
		editorService: editorService,
	}
}

// UpdateOCRText handles PUT /api/v1/books/{book_id}/pages/{page_id}/ocr-text
func (h *Handler) UpdateOCRText(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract path parameters
	bookID, err := uuid.Parse(getPathParam(r, "book_id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	pageID, err := uuid.Parse(getPathParam(r, "page_id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid page ID")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		respondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	// Parse request body
	var req models.UpdateOCRTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Update OCR text
	correction, err := h.editorService.UpdateOCRText(ctx, bookID, pageID, userID, req.CorrectedText)
	if err != nil {
		switch err {
		case ocr.ErrPageNotFound:
			respondError(w, http.StatusNotFound, "page not found")
		case ocr.ErrUnauthorized:
			respondError(w, http.StatusForbidden, "access denied")
		case ocr.ErrInvalidCorrectedText, ocr.ErrTextTooLong:
			respondError(w, http.StatusBadRequest, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "failed to update OCR text")
		}
		return
	}

	// Respond with success
	respondJSON(w, http.StatusOK, models.UpdateOCRTextResponse{
		Success:    true,
		Correction: *correction,
		Message:    "OCR text updated successfully",
	})
}

// GetOCRHistory handles GET /api/v1/books/{book_id}/pages/{page_id}/ocr-history
func (h *Handler) GetOCRHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract path parameters
	bookID, err := uuid.Parse(getPathParam(r, "book_id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	pageID, err := uuid.Parse(getPathParam(r, "page_id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid page ID")
		return
	}

	// Get user ID from context
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		respondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	// Parse query parameters
	limit := 10
	offset := 0
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get correction history
	history, err := h.editorService.GetCorrectionHistory(ctx, bookID, pageID, userID, limit, offset)
	if err != nil {
		switch err {
		case ocr.ErrPageNotFound:
			respondError(w, http.StatusNotFound, "page not found")
		case ocr.ErrUnauthorized:
			respondError(w, http.StatusForbidden, "access denied")
		default:
			respondError(w, http.StatusInternalServerError, "failed to get correction history")
		}
		return
	}

	respondJSON(w, http.StatusOK, history)
}

// Helper functions

func getPathParam(r *http.Request, key string) string {
	// This would typically use a router like gorilla/mux or chi
	// For now, we'll use a simple implementation
	return r.URL.Query().Get(key)
}

func getUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	// In a real implementation, this would extract the user ID from the context
	// set by an authentication middleware
	userIDValue := ctx.Value("user_id")
	if userIDValue == nil {
		return uuid.Nil, false
	}

	switch v := userIDValue.(type) {
	case uuid.UUID:
		return v, true
	case string:
		if userID, err := uuid.Parse(v); err == nil {
			return userID, true
		}
	}

	return uuid.Nil, false
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{
		"error": message,
	})
}
