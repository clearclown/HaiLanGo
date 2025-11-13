package handler

import (
	"net/http"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/service/srs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReviewHandler は復習APIのハンドラー
type ReviewHandler struct {
	srsService *srs.SRSService
}

// NewReviewHandler は新しいReviewHandlerを作成
func NewReviewHandler(srsService *srs.SRSService) *ReviewHandler {
	return &ReviewHandler{
		srsService: srsService,
	}
}

// GetReviewItemsRequest は復習項目取得のリクエスト
type GetReviewItemsRequest struct {
	UserID string `uri:"user_id" binding:"required,uuid"`
}

// GetReviewItemsResponse は復習項目取得のレスポンス
type GetReviewItemsResponse struct {
	UrgentItems      []ReviewItemResponse `json:"urgent_items"`
	RecommendedItems []ReviewItemResponse `json:"recommended_items"`
	RelaxedItems     []ReviewItemResponse `json:"relaxed_items"`
}

// ReviewItemResponse は復習項目のレスポンス
type ReviewItemResponse struct {
	ID             string     `json:"id"`
	BookID         string     `json:"book_id"`
	PageNumber     int        `json:"page_number"`
	ItemType       string     `json:"item_type"`
	Content        string     `json:"content"`
	Translation    string     `json:"translation"`
	ReviewCount    int        `json:"review_count"`
	LastReviewDate *time.Time `json:"last_review_date"`
	NextReviewDate *time.Time `json:"next_review_date"`
	LastScore      int        `json:"last_score"`
}

// CompleteReviewRequest は復習完了のリクエスト
type CompleteReviewRequest struct {
	ItemID       string `uri:"item_id" binding:"required,uuid"`
	Score        int    `json:"score" binding:"required,min=0,max=100"`
	TimeSpentSec int    `json:"time_spent_sec" binding:"min=0"`
}

// CompleteReviewResponse は復習完了のレスポンス
type CompleteReviewResponse struct {
	Success        bool       `json:"success"`
	NextReviewDate *time.Time `json:"next_review_date"`
}

// GetStatsRequest は統計取得のリクエスト
type GetStatsRequest struct {
	UserID string `uri:"user_id" binding:"required,uuid"`
}

// GetStatsResponse は統計取得のレスポンス
type GetStatsResponse struct {
	TotalReviewItems int     `json:"total_review_items"`
	UrgentItems      int     `json:"urgent_items"`
	RecommendedItems int     `json:"recommended_items"`
	RelaxedItems     int     `json:"relaxed_items"`
	WeeklyReviewCount int    `json:"weekly_review_count"`
	CurrentStreak    int     `json:"current_streak"`
	LongestStreak    int     `json:"longest_streak"`
	AverageScore     float64 `json:"average_score"`
}

// GetReviewItems は復習項目を取得（優先度別）
// GET /api/v1/review/items/:user_id
func (h *ReviewHandler) GetReviewItems(c *gin.Context) {
	var req GetReviewItemsRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	now := time.Now()
	items, err := h.srsService.GetReviewItems(c.Request.Context(), userID, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := &GetReviewItemsResponse{
		UrgentItems:      make([]ReviewItemResponse, 0),
		RecommendedItems: make([]ReviewItemResponse, 0),
		RelaxedItems:     make([]ReviewItemResponse, 0),
	}

	// 緊急項目
	for _, item := range items.UrgentItems {
		response.UrgentItems = append(response.UrgentItems, ReviewItemResponse{
			ID:             item.ID.String(),
			BookID:         item.BookID.String(),
			PageNumber:     item.PageNumber,
			ItemType:       item.ItemType,
			Content:        item.Content,
			Translation:    item.Translation,
			ReviewCount:    item.ReviewCount,
			LastReviewDate: item.LastReviewDate,
			NextReviewDate: item.NextReviewDate,
			LastScore:      item.LastScore,
		})
	}

	// 推奨項目
	for _, item := range items.RecommendedItems {
		response.RecommendedItems = append(response.RecommendedItems, ReviewItemResponse{
			ID:             item.ID.String(),
			BookID:         item.BookID.String(),
			PageNumber:     item.PageNumber,
			ItemType:       item.ItemType,
			Content:        item.Content,
			Translation:    item.Translation,
			ReviewCount:    item.ReviewCount,
			LastReviewDate: item.LastReviewDate,
			NextReviewDate: item.NextReviewDate,
			LastScore:      item.LastScore,
		})
	}

	// 余裕あり項目
	for _, item := range items.RelaxedItems {
		response.RelaxedItems = append(response.RelaxedItems, ReviewItemResponse{
			ID:             item.ID.String(),
			BookID:         item.BookID.String(),
			PageNumber:     item.PageNumber,
			ItemType:       item.ItemType,
			Content:        item.Content,
			Translation:    item.Translation,
			ReviewCount:    item.ReviewCount,
			LastReviewDate: item.LastReviewDate,
			NextReviewDate: item.NextReviewDate,
			LastScore:      item.LastScore,
		})
	}

	c.JSON(http.StatusOK, response)
}

// CompleteReview は復習完了処理
// POST /api/v1/review/items/:item_id/complete
func (h *ReviewHandler) CompleteReview(c *gin.Context) {
	var req CompleteReviewRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemID, err := uuid.Parse(req.ItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	err = h.srsService.CompleteReview(c.Request.Context(), itemID, req.Score, req.TimeSpentSec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新後の復習項目を取得して次回復習日を返す
	// （実装簡略化のため、ここでは成功レスポンスのみ）
	c.JSON(http.StatusOK, &CompleteReviewResponse{
		Success: true,
	})
}

// GetStats は統計情報を取得
// GET /api/v1/review/stats/:user_id
func (h *ReviewHandler) GetStats(c *gin.Context) {
	var req GetStatsRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	now := time.Now()
	stats, err := h.srsService.GetStats(c.Request.Context(), userID, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := &GetStatsResponse{
		TotalReviewItems:  stats.TotalReviewItems,
		UrgentItems:       stats.UrgentItems,
		RecommendedItems:  stats.RecommendedItems,
		RelaxedItems:      stats.RelaxedItems,
		WeeklyReviewCount: stats.WeeklyReviewCount,
		CurrentStreak:     stats.CurrentStreak,
		LongestStreak:     stats.LongestStreak,
		AverageScore:      stats.AverageScore,
	}

	c.JSON(http.StatusOK, response)
}
