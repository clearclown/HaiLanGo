package handler

import (
	"net/http"
	"strconv"

	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StatsHandler handles statistics-related HTTP requests
type StatsHandler struct {
	repo repository.StatsRepositoryInterface
}

// NewStatsHandler creates a new stats handler
func NewStatsHandler(repo repository.StatsRepositoryInterface) *StatsHandler {
	return &StatsHandler{
		repo: repo,
	}
}

// GetDashboard handles GET /api/v1/stats/dashboard
// @Summary Get dashboard statistics
// @Description Get overall learning statistics for dashboard
// @Tags stats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.DashboardStatsFlat
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/stats/dashboard [get]
func (h *StatsHandler) GetDashboard(c *gin.Context) {
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

	stats, err := h.repo.GetDashboardStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetLearningTime handles GET /api/v1/stats/learning-time
// @Summary Get learning time data
// @Description Get learning time data for specified period
// @Tags stats
// @Accept json
// @Produce json
// @Param period query string false "Period (day|week|month|year)" default(week)
// @Security BearerAuth
// @Success 200 {object} models.LearningTimeData
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/stats/learning-time [get]
func (h *StatsHandler) GetLearningTime(c *gin.Context) {
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

	period := c.DefaultQuery("period", "week")

	data, err := h.repo.GetLearningTimeData(c.Request.Context(), userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get learning time data"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetProgress handles GET /api/v1/stats/progress
// @Summary Get progress data
// @Description Get learning progress data for specified period
// @Tags stats
// @Accept json
// @Produce json
// @Param period query string false "Period (week|month|year)" default(month)
// @Security BearerAuth
// @Success 200 {object} models.ProgressData
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/stats/progress [get]
func (h *StatsHandler) GetProgress(c *gin.Context) {
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

	period := c.DefaultQuery("period", "month")

	data, err := h.repo.GetProgressData(c.Request.Context(), userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get progress data"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetWeakPoints handles GET /api/v1/stats/weak-points
// @Summary Get weak points analysis
// @Description Get weak points (words/phrases with low scores)
// @Tags stats
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Security BearerAuth
// @Success 200 {object} models.WeakPointsData
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/stats/weak-points [get]
func (h *StatsHandler) GetWeakPoints(c *gin.Context) {
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

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	data, err := h.repo.GetWeakPoints(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get weak points"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// RegisterRoutes registers stats routes
func (h *StatsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	stats := rg.Group("/stats")
	{
		stats.GET("/dashboard", h.GetDashboard)
		stats.GET("/learning-time", h.GetLearningTime)
		stats.GET("/progress", h.GetProgress)
		stats.GET("/weak-points", h.GetWeakPoints)
	}
}
