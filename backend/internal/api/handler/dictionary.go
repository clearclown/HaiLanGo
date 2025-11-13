package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/clearclown/HaiLanGo/backend/internal/service/dictionary"
	pkgDict "github.com/clearclown/HaiLanGo/backend/pkg/dictionary"
)

// DictionaryHandler handles dictionary-related HTTP requests
type DictionaryHandler struct {
	service *dictionary.Service
}

// NewDictionaryHandler creates a new dictionary handler
func NewDictionaryHandler(service *dictionary.Service) *DictionaryHandler {
	return &DictionaryHandler{
		service: service,
	}
}

// LookupWord handles GET /api/v1/dictionary/words/{word}
func (h *DictionaryHandler) LookupWord(w http.ResponseWriter, r *http.Request) {
	// Extract word from URL path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 6 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	word := parts[5]

	// Get language from query parameter (default: en)
	language := r.URL.Query().Get("language")
	if language == "" {
		language = "en"
	}

	// Lookup word
	entry, err := h.service.LookupWord(r.Context(), word, language)
	if err != nil {
		if err == pkgDict.ErrWordNotFound {
			http.Error(w, "Word not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entry); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// LookupWordDetails handles GET /api/v1/dictionary/words/{word}/details
func (h *DictionaryHandler) LookupWordDetails(w http.ResponseWriter, r *http.Request) {
	// Extract word from URL path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 6 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	word := parts[5]

	// Get language from query parameter (default: en)
	language := r.URL.Query().Get("language")
	if language == "" {
		language = "en"
	}

	// Lookup word details
	entry, err := h.service.LookupWordDetails(r.Context(), word, language)
	if err != nil {
		if err == pkgDict.ErrWordNotFound {
			http.Error(w, "Word not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entry); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
