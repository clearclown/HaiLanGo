package handler

import (
	"net/http"
	"strconv"

	"github.com/clearclown/HaiLanGo/backend/internal/service/stats"
	"github.com/gin-gonic/gin"
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
func (h *StatsHandler) GetDashboard(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get dashboard stats
	dashboard, err := h.service.GetDashboardStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard stats"})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

// GetLearningTime handles GET /api/v1/stats/learning-time
func (h *StatsHandler) GetLearningTime(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get learning time stats
	learningTime, err := h.service.GetLearningTimeStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get learning time stats"})
		return
	}

	c.JSON(http.StatusOK, learningTime)
}

// GetProgress handles GET /api/v1/stats/progress
func (h *StatsHandler) GetProgress(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get progress stats
	progress, err := h.service.GetProgressStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get progress stats"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetStreak handles GET /api/v1/stats/streak
func (h *StatsHandler) GetStreak(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get streak stats
	streak, err := h.service.GetStreakStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get streak stats"})
		return
	}

	c.JSON(http.StatusOK, streak)
}

// GetLearningTimeChart handles GET /api/v1/stats/learning-time-chart?days=7
func (h *StatsHandler) GetLearningTimeChart(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get days parameter (default: 7)
	daysStr := c.Query("days")
	days := 7
	if daysStr != "" {
		days, err = strconv.Atoi(daysStr)
		if err != nil || days < 1 || days > 365 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid days parameter"})
			return
		}
	}

	// Get learning time chart
	chart, err := h.service.GetLearningTimeChart(c.Request.Context(), userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get learning time chart"})
		return
	}

	c.JSON(http.StatusOK, chart)
}

// GetProgressChart handles GET /api/v1/stats/progress-chart?days=30
func (h *StatsHandler) GetProgressChart(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get days parameter (default: 30)
	daysStr := c.Query("days")
	days := 30
	if daysStr != "" {
		days, err = strconv.Atoi(daysStr)
		if err != nil || days < 1 || days > 365 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid days parameter"})
			return
		}
	}

	// Get progress chart
	chart, err := h.service.GetProgressChart(c.Request.Context(), userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get progress chart"})
		return
	}

	c.JSON(http.StatusOK, chart)
}

// GetWeakWords handles GET /api/v1/stats/weak-words?limit=10
func (h *StatsHandler) GetWeakWords(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get limit parameter (default: 10)
	limitStr := c.Query("limit")
	limit := 10
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
	}

	// Get weak words
	weakWords, err := h.service.GetWeakWords(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get weak words"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"weak_words": weakWords})
}

// RegisterRoutes registers stats routes
func (h *StatsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	stats := rg.Group("/stats")
	{
		stats.GET("/dashboard", h.GetDashboard)
		stats.GET("/learning-time", h.GetLearningTime)
		stats.GET("/progress", h.GetProgress)
		stats.GET("/streak", h.GetStreak)
		stats.GET("/learning-time-chart", h.GetLearningTimeChart)
		stats.GET("/progress-chart", h.GetProgressChart)
		stats.GET("/weak-words", h.GetWeakWords)
	}
}
