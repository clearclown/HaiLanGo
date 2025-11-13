package main

import (
	"log"
	"os"

	"github.com/clearclown/HaiLanGo/backend/internal/api/handler"
	"github.com/clearclown/HaiLanGo/backend/internal/api/router"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/clearclown/HaiLanGo/backend/internal/service/srs"
	"github.com/joho/godotenv"
)

func main() {
	// .envファイルを読み込み（存在する場合）
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// ポート設定
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	// モックモードの確認
	useMocks := os.Getenv("USE_MOCK_APIS") == "true"

	log.Println("=================================")
	log.Println("HaiLanGo Backend Server")
	log.Println("=================================")
	log.Printf("Port: %s\n", port)
	log.Printf("Mock Mode: %v\n", useMocks)
	log.Println("=================================")

	// モックリポジトリを使用（実際のDB実装は後で追加）
	reviewItemRepo := repository.NewMockReviewItemRepository()

	// サービス層を初期化
	srsService := srs.NewSRSService(reviewItemRepo)

	// ハンドラーを初期化
	reviewHandler := handler.NewReviewHandler(srsService)

	// ルーターをセットアップ
	r := router.SetupRouter(reviewHandler)

	// サーバー起動
	log.Printf("Server starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
