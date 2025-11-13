package router

import (
	"github.com/clearclown/HaiLanGo/backend/internal/api/handler"
<<<<<<< HEAD
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
=======
	"github.com/clearclown/HaiLanGo/backend/internal/api/middleware"
	"github.com/clearclown/HaiLanGo/backend/internal/service"
	"github.com/clearclown/HaiLanGo/backend/pkg/storage"
	"github.com/gin-gonic/gin"
)

// SetupRouter はAPIルーターをセットアップする
func SetupRouter(authHandler *handler.AuthHandler, storagePath string) *gin.Engine {
	// Ginエンジンの作成
	r := gin.Default()

	// ミドルウェアの設定
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimiter())

	// ストレージを初期化
	localStorage := storage.NewLocalStorage(storagePath)

	// 一時ディレクトリを作成
	tempDir := storagePath + "/temp"

	// サービスを初期化
	uploadService := service.NewUploadService(localStorage, tempDir)

	// ハンドラーを初期化
	uploadHandler := handler.NewUploadHandler(uploadService)

	// ヘルスチェックエンドポイント
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "HaiLanGo API is running",
		})
	})

	// API v1グループ
	v1 := r.Group("/api/v1")
	{
		// 認証エンドポイント
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}

		// アップロードルートを登録
		uploadHandler.RegisterRoutes(v1)
	}

	return r
>>>>>>> origin/main
}
