package main

import (
	"log"
	"os"

	"github.com/clearclown/HaiLanGo/internal/api/learning"
	learningService "github.com/clearclown/HaiLanGo/internal/service/learning"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 環境変数の読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// モックリポジトリの作成（本番ではPostgreSQLリポジトリを使用）
	repo := NewMockRepository()

	// サービスの作成
	service := learningService.NewService(repo)

	// ハンドラーの作成
	handler := learning.NewHandler(service)

	// Ginルーターの設定
	router := gin.Default()

	// CORSの設定
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

	// APIルートの設定
	api := router.Group("/api/v1")
	{
		books := api.Group("/books")
		{
			books.GET("/:bookId/pages/:pageNumber", handler.GetPage)
			books.POST("/:bookId/pages/:pageNumber/complete", handler.MarkPageCompleted)
			books.GET("/:bookId/progress", handler.GetProgress)
		}
	}

	// サーバーの起動
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
