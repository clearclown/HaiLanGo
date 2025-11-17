package handler

import (
	"net/http"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/clearclown/HaiLanGo/backend/internal/service"
	"github.com/clearclown/HaiLanGo/backend/internal/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReviewHandler struct {
	repo    repository.ReviewRepository
	srsAlgo *service.SM2Algorithm
	wsHub   *websocket.Hub
}

func NewReviewHandler(repo repository.ReviewRepository, wsHub *websocket.Hub) *ReviewHandler {
	return &ReviewHandler{
		repo:    repo,
		srsAlgo: service.NewSM2Algorithm(),
		wsHub:   wsHub,
	}
}

// GetStats godoc
// @Summary Get review statistics
// @Tags review
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.ReviewStats
// @Router /api/v1/review/stats [get]
func (h *ReviewHandler) GetStats(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := userIDStr.(string)

	// 今日の開始時刻
	todayStart := time.Now().Truncate(24 * time.Hour)

	// すべての復習アイテムを取得
	items, err := h.repo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch review items"})
		return
	}

	stats := models.ReviewStats{
		UrgentCount:      0,
		RecommendedCount: 0,
		OptionalCount:    0,
	}

	// 優先度別にカウント
	for _, item := range items {
		priority := h.srsAlgo.CalculatePriority(item.NextReview)
		switch priority {
		case "urgent":
			stats.UrgentCount++
		case "recommended":
			stats.RecommendedCount++
		case "optional":
			stats.OptionalCount++
		}
	}

	// 今日完了した復習数を取得
	stats.TotalCompletedToday, err = h.repo.CountCompletedToday(c.Request.Context(), userID, todayStart)
	if err != nil {
		stats.TotalCompletedToday = 0
	}

	// 今週の完了率を計算
	weekStart := todayStart.Add(-7 * 24 * time.Hour)
	weeklyCompleted, _ := h.repo.CountCompletedSince(c.Request.Context(), userID, weekStart)
	weeklyTarget := len(items) * 7 // 1日1回 × 7日
	if weeklyTarget > 0 {
		stats.WeeklyCompletionRate = float64(weeklyCompleted) / float64(weeklyTarget) * 100
	}

	// WebSocket通知: 緊急の復習がある場合に通知
	if stats.UrgentCount > 0 && h.wsHub != nil {
		userUUID, err := uuid.Parse(userID)
		if err == nil {
			// ReviewReminderMessageを送信
			wsReviewItems := []websocket.ReviewItem{}
			for _, item := range items {
				priority := h.srsAlgo.CalculatePriority(item.NextReview)
				if priority == "urgent" {
					// UUIDに変換
					itemUUID, err := uuid.Parse(item.ID)
					if err != nil {
						continue
					}
					wsItem := websocket.ReviewItem{
						ID:          itemUUID,
						Content:     item.Text,
						Translation: item.Translation,
						DueDate:     item.NextReview,
						Priority:    priority,
					}
					wsReviewItems = append(wsReviewItems, wsItem)
				}
			}

			message, err := websocket.NewReviewReminderMessage(
				stats.UrgentCount,
				wsReviewItems,
			)
			if err == nil {
				h.wsHub.SendToUser(userUUID, message)
			}
		}
	}

	c.JSON(http.StatusOK, stats)
}

// GetItems godoc
// @Summary Get review items by priority
// @Tags review
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param priority query string false "Priority filter (urgent, recommended, optional)"
// @Success 200 {object} map[string][]models.ReviewItem
// @Router /api/v1/review/items [get]
func (h *ReviewHandler) GetItems(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := userIDStr.(string)
	priorityFilter := c.Query("priority")

	items, err := h.repo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch review items"})
		return
	}

	// 優先度でフィルタリング
	var filteredItems []*models.ReviewItem
	for _, item := range items {
		priority := h.srsAlgo.CalculatePriority(item.NextReview)
		item.Priority = priority

		if priorityFilter == "" || priority == priorityFilter {
			filteredItems = append(filteredItems, item)
		}
	}

	if filteredItems == nil {
		filteredItems = []*models.ReviewItem{}
	}

	c.JSON(http.StatusOK, gin.H{"items": filteredItems})
}

// SubmitReview godoc
// @Summary Submit review result
// @Tags review
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param result body models.ReviewResult true "Review result"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/review/submit [post]
func (h *ReviewHandler) SubmitReview(c *gin.Context) {
	var result models.ReviewResult
	if err := c.ShouldBindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := userIDStr.(string)

	// 復習アイテムを取得
	item, err := h.repo.FindByID(c.Request.Context(), result.ItemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review item not found"})
		return
	}

	// 所有権チェック
	if item.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// SRSアルゴリズムで次の復習日時を計算
	nextInterval, nextEaseFactor, nextReview := h.srsAlgo.CalculateNextReview(
		item.EaseFactor,
		item.IntervalDays,
		result.Score,
	)

	// 習熟度を更新
	newMasteryLevel := item.MasteryLevel
	if result.Score >= 70 {
		newMasteryLevel += 10
		if newMasteryLevel > 100 {
			newMasteryLevel = 100
		}
	} else if result.Score < 50 {
		newMasteryLevel -= 5
		if newMasteryLevel < 0 {
			newMasteryLevel = 0
		}
	}

	// アイテムを更新
	item.MasteryLevel = newMasteryLevel
	item.IntervalDays = nextInterval
	item.EaseFactor = nextEaseFactor
	item.LastReviewed = time.Now()
	item.NextReview = nextReview
	item.ReviewCount++

	if err := h.repo.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review item"})
		return
	}

	// 履歴を保存
	history := &models.ReviewHistory{
		ReviewItemID: item.ID,
		UserID:       userID,
		Score:        result.Score,
		ReviewedAt:   time.Now(),
	}

	if err := h.repo.SaveHistory(c.Request.Context(), history); err != nil {
		// エラーログは出すが、レスポンスは成功を返す
		c.JSON(http.StatusOK, gin.H{
			"success":     true,
			"next_review": nextReview.Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"next_review": nextReview.Format(time.RFC3339),
	})
}

// RegisterRoutes registers review routes
func (h *ReviewHandler) RegisterRoutes(rg *gin.RouterGroup) {
	review := rg.Group("/review")
	{
		review.GET("/stats", h.GetStats)
		review.GET("/items", h.GetItems)
		review.POST("/submit", h.SubmitReview)
	}
}
