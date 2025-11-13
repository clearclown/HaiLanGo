package router

import (
	"github.com/clearclown/HaiLanGo/backend/internal/api/handler"
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
}
