package main

import (
	"fmt"
	"log"
	"os"

	"github.com/clearclown/HaiLanGo/internal/api/router"
	"github.com/joho/godotenv"
)

func main() {
	// 環境変数を読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// ストレージパスを取得（環境変数またはデフォルト値）
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}

	// ストレージディレクトリを作成
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	// サーバーポートを取得（環境変数またはデフォルト値）
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	// ルーターをセットアップ
	r := router.SetupRouter(storagePath)

	// サーバーを起動
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Starting HaiLanGo API server on %s", addr)
	log.Printf("Storage path: %s", storagePath)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
