package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/clearclown/HaiLanGo/internal/service/stats"
	"github.com/google/uuid"
)

// StatsHandler handles statistics-related HTTP requests
type StatsHandler struct {
	service *stats.Service
}

// NewStatsHandler creates a new stats handler
func NewStatsHandler(service *stats.Service) *StatsHandler {
	return &StatsHandler{
		service: service,
	}
}

// GetDashboard handles GET /api/v1/stats/dashboard
func (h *StatsHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (assuming authentication middleware sets this)
	userIDStr := r.Context().Value("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get dashboard stats
	dashboard, err := h.service.GetDashboardStats(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get dashboard stats", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}

// GetLearningTime handles GET /api/v1/stats/learning-time
func (h *StatsHandler) GetLearningTime(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userIDStr := r.Context().Value("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get learning time stats
	learningTime, err := h.service.GetLearningTimeStats(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get learning time stats", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(learningTime)
}

// GetProgress handles GET /api/v1/stats/progress
func (h *StatsHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userIDStr := r.Context().Value("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get progress stats
	progress, err := h.service.GetProgressStats(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get progress stats", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}

// GetStreak handles GET /api/v1/stats/streak
func (h *StatsHandler) GetStreak(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userIDStr := r.Context().Value("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get streak stats
	streak, err := h.service.GetStreakStats(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get streak stats", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(streak)
}

// GetLearningTimeChart handles GET /api/v1/stats/learning-time-chart?days=7
func (h *StatsHandler) GetLearningTimeChart(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userIDStr := r.Context().Value("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get days parameter (default: 7)
	daysStr := r.URL.Query().Get("days")
	days := 7
	if daysStr != "" {
		days, err = strconv.Atoi(daysStr)
		if err != nil || days < 1 || days > 365 {
			http.Error(w, "Invalid days parameter", http.StatusBadRequest)
			return
		}
	}

	// Get learning time chart
	chart, err := h.service.GetLearningTimeChart(r.Context(), userID, days)
	if err != nil {
		http.Error(w, "Failed to get learning time chart", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chart)
}

// GetProgressChart handles GET /api/v1/stats/progress-chart?days=30
func (h *StatsHandler) GetProgressChart(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userIDStr := r.Context().Value("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get days parameter (default: 30)
	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		days, err = strconv.Atoi(daysStr)
		if err != nil || days < 1 || days > 365 {
			http.Error(w, "Invalid days parameter", http.StatusBadRequest)
			return
		}
	}

	// Get progress chart
	chart, err := h.service.GetProgressChart(r.Context(), userID, days)
	if err != nil {
		http.Error(w, "Failed to get progress chart", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chart)
}

// GetWeakWords handles GET /api/v1/stats/weak-words?limit=10
func (h *StatsHandler) GetWeakWords(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userIDStr := r.Context().Value("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get limit parameter (default: 10)
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 100 {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	// Get weak words
	weakWords, err := h.service.GetWeakWords(r.Context(), userID, limit)
	if err != nil {
		http.Error(w, "Failed to get weak words", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"weak_words": weakWords,
	})
}
