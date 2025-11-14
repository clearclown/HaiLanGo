package router

import (
	"database/sql"

	"github.com/clearclown/HaiLanGo/backend/internal/api/handler"
	"github.com/clearclown/HaiLanGo/backend/internal/api/middleware"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/clearclown/HaiLanGo/backend/internal/service"
	"github.com/clearclown/HaiLanGo/backend/pkg/storage"
	"github.com/gin-gonic/gin"
)

// SetupRouter はAPIルーターをセットアップする
func SetupRouter(
	db *sql.DB,
	authHandler *handler.AuthHandler,
	storagePath string,
) *gin.Engine {
	// Ginエンジンの作成
	r := gin.Default()

	// ミドルウェアの設定
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimiter())

	// ストレージを初期化
	localStorage := storage.NewLocalStorage(storagePath)
	tempDir := storagePath + "/temp"

	// ========================================
	// リポジトリの初期化
	// ========================================
	bookRepo := repository.NewBookRepositoryPostgres(db)     // PostgreSQL実装
	reviewRepo := repository.NewInMemoryReviewRepository()   // InMemory実装
	statsRepo := repository.NewInMemoryStatsRepository()     // InMemory実装
	learningRepo := repository.NewInMemoryLearningRepository() // InMemory実装

	// ========================================
	// サービスの初期化
	// ========================================
	uploadService := service.NewUploadService(localStorage, tempDir)
	// ocrService := service.NewOCRService() // TODO: 実装必要
	// statsService := stats.NewService(statsRepo) // TODO: 実装必要
	// srsService := srs.NewSRSService(reviewRepo) // TODO: 実装必要

	// ========================================
	// ハンドラーの初期化
	// ========================================
	uploadHandler := handler.NewUploadHandler(uploadService)
	booksHandler := handler.NewBooksHandler(bookRepo)
	reviewHandler := handler.NewReviewHandler(reviewRepo)
	statsHandler := handler.NewStatsHandler(statsRepo)
	learningHandler := handler.NewLearningHandler(learningRepo)
	// dictionaryHandler := handler.NewDictionaryHandler() // TODO: 実装必要
	// patternHandler := handler.NewPatternHandler() // TODO: 実装必要
	// ocrHandler := ocr.NewOCRHandler(ocrService) // TODO: 実装必要
	// paymentHandler := payment.NewPaymentHandler() // TODO: 実装必要
	// teacherModeHandler := teachermode.NewTeacherModeHandler() // TODO: 実装必要

	// WebSocketハブを初期化
	// wsHub := websocket.NewHub() // TODO: 実装必要
	// go wsHub.Run()
	// wsHandler := websocket.NewHandler(wsHub) // TODO: 実装必要

	// ========================================
	// ヘルスチェックエンドポイント
	// ========================================
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "HaiLanGo API is running",
			"version": "1.0.0",
		})
	})

	// ========================================
	// API v1グループ
	// ========================================
	v1 := r.Group("/api/v1")
	{
		// 認証エンドポイント（認証不要）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}

		// 以下、認証必須
		authenticated := v1.Group("")
		authenticated.Use(middleware.AuthRequired())
		{
			// Books API
			booksHandler.RegisterRoutes(authenticated)

			// Upload API
			uploadHandler.RegisterRoutes(authenticated)

			// Review API
			reviewHandler.RegisterRoutes(authenticated)

			// Stats API
			statsHandler.RegisterRoutes(authenticated)

			// Learning API
			learningHandler.RegisterRoutes(authenticated)

			// OCR API
			// ocrHandler.RegisterRoutes(authenticated) // TODO: Uncomment when implemented

			// Pattern API
			// patternHandler.RegisterRoutes(authenticated) // TODO: Uncomment when implemented

			// Teacher Mode API
			// teacherModeHandler.RegisterRoutes(authenticated) // TODO: Uncomment when implemented

			// Dictionary API
			// dictionaryHandler.RegisterRoutes(authenticated) // TODO: Uncomment when implemented

			// Payment API
			// paymentHandler.RegisterRoutes(authenticated) // TODO: Uncomment when implemented

			// WebSocket API
			// authenticated.GET("/ws", wsHandler.HandleWebSocket) // TODO: Uncomment when implemented
		}
	}

	return r
}
