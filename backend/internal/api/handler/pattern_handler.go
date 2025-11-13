package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/service/pattern"
	"github.com/google/uuid"
)

// PatternHandler handles pattern-related HTTP requests
type PatternHandler struct {
	extractor *pattern.Extractor
}

// NewPatternHandler creates a new pattern handler
func NewPatternHandler() *PatternHandler {
	return &PatternHandler{
		extractor: pattern.NewExtractor(),
	}
}

// ExtractPatterns handles POST /api/v1/books/{book_id}/patterns/extract
func (h *PatternHandler) ExtractPatterns(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req models.PatternExtractionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Fetch book pages from database
	// For now, this is a placeholder
	pages := []pattern.PageText{} // Would fetch from database

	// Extract patterns
	patterns, err := h.extractor.ExtractPatterns(r.Context(), req.BookID, pages, req.MinFrequency)
	if err != nil {
		http.Error(w, "Failed to extract patterns", http.StatusInternalServerError)
		return
	}

	// Create response
	resp := models.PatternExtractionResponse{
		Patterns:       patterns,
		TotalFound:     len(patterns),
		ProcessedPages: req.PageEnd - req.PageStart + 1,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetPatterns handles GET /api/v1/books/{book_id}/patterns
func (h *PatternHandler) GetPatterns(w http.ResponseWriter, r *http.Request) {
	// Extract book_id from URL
	bookIDStr := r.URL.Query().Get("book_id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	// TODO: Fetch patterns from database
	// For now, return empty list
	patterns := []models.Pattern{}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"patterns": patterns,
		"book_id":  bookID,
	})
}

// GetPatternPractice handles GET /api/v1/patterns/{pattern_id}/practice
func (h *PatternHandler) GetPatternPractice(w http.ResponseWriter, r *http.Request) {
	// Extract pattern_id from URL
	patternIDStr := r.URL.Query().Get("pattern_id")
	patternID, err := uuid.Parse(patternIDStr)
	if err != nil {
		http.Error(w, "Invalid pattern ID", http.StatusBadRequest)
		return
	}

	// Extract count parameter (default: 10)
	countStr := r.URL.Query().Get("count")
	count := 10
	if countStr != "" {
		if c, err := strconv.Atoi(countStr); err == nil {
			count = c
		}
	}

	// TODO: Fetch pattern and generate practice exercises from database
	// For now, return placeholder
	practices := []models.PatternPractice{}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pattern_id": patternID,
		"practices":  practices,
		"count":      count,
	})
}
