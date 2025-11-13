package main

import (
<<<<<<< HEAD
=======
	"database/sql"
	"fmt"
>>>>>>> origin/main
	"log"
	"os"

	"github.com/clearclown/HaiLanGo/backend/internal/api/handler"
	"github.com/clearclown/HaiLanGo/backend/internal/api/router"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
<<<<<<< HEAD
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
=======
	"github.com/clearclown/HaiLanGo/backend/internal/service"
	"github.com/clearclown/HaiLanGo/backend/pkg/jwt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// 環境変数を読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// 環境変数の読み込み
	port := getEnv("BACKEND_PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgresql://HaiLanGo:password@localhost:5432/HaiLanGo_dev?sslmode=disable")
	storagePath := getEnv("STORAGE_PATH", "./storage")

	// ストレージディレクトリを作成
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		log.Fatalf("ストレージディレクトリ作成エラー: %v", err)
	}

	// データベース接続
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}
	defer db.Close()

	// データベース接続テスト
	if err := db.Ping(); err != nil {
		log.Fatalf("データベースPingエラー: %v", err)
	}

	log.Println("データベースに接続しました")

	// RSA鍵ペアの生成（本番環境では事前に生成した鍵を読み込むこと）
	if err := jwt.GenerateRSAKeys(); err != nil {
		log.Fatalf("RSA鍵生成エラー: %v", err)
	}

	log.Println("RSA鍵ペアを生成しました")

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(db)

	// サービスの初期化
	authService := service.NewAuthService(userRepo)

	// ハンドラーの初期化
	authHandler := handler.NewAuthHandler(authService)

	// ルーターのセットアップ
	r := router.SetupRouter(authHandler, storagePath)

	// サーバー起動
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("HaiLanGo APIサーバーを起動します: %s", addr)
	log.Printf("ストレージパス: %s", storagePath)

	if err := r.Run(addr); err != nil {
		log.Fatalf("サーバー起動エラー: %v", err)
	}
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
>>>>>>> origin/main
