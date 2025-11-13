package router

import (
	"github.com/clearclown/HaiLanGo/backend/internal/api/handler"
	"github.com/gin-gonic/gin"
)

// SetupRouter はAPIルーターをセットアップ
func SetupRouter(reviewHandler *handler.ReviewHandler) *gin.Engine {
	router := gin.Default()

	// ヘルスチェック
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		// 復習API
		review := v1.Group("/review")
		{
			// 復習項目取得（優先度別）
			review.GET("/items/:user_id", reviewHandler.GetReviewItems)

			// 復習完了
			review.POST("/items/:item_id/complete", reviewHandler.CompleteReview)

			// 統計情報取得
			review.GET("/stats/:user_id", reviewHandler.GetStats)
		}
	}

	return router
}
