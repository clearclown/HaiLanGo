package router

import (
	"github.com/clearclown/HaiLanGo/internal/api/handler"
	"github.com/clearclown/HaiLanGo/internal/service"
	"github.com/clearclown/HaiLanGo/pkg/storage"
	"github.com/gin-gonic/gin"
)

// SetupRouter はGinルーターをセットアップする
func SetupRouter(storagePath string) *gin.Engine {
	// Ginエンジンを作成
	router := gin.Default()

	// CORS設定（開発環境用）
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// ストレージを初期化
	localStorage := storage.NewLocalStorage(storagePath)

	// 一時ディレクトリを作成
	tempDir := storagePath + "/temp"

	// サービスを初期化
	uploadService := service.NewUploadService(localStorage, tempDir)

	// ハンドラーを初期化
	uploadHandler := handler.NewUploadHandler(uploadService)

	// APIグループ
	api := router.Group("/api/v1")
	{
		// ヘルスチェック
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"message": "HaiLanGo API is running",
			})
		})

		// アップロードルートを登録
		uploadHandler.RegisterRoutes(api)
	}

	return router
}
