package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/clearclown/HaiLanGo/backend/internal/api/handler"
	"github.com/clearclown/HaiLanGo/backend/internal/api/router"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
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
	var userRepo repository.UserRepository
	if err := db.Ping(); err != nil {
		log.Printf("⚠️  データベースPing失敗 (InMemoryリポジトリを使用): %v", err)
		// InMemoryリポジトリを使用
		userRepo = repository.NewInMemoryUserRepository()
		log.Println("✅ InMemoryUserRepositoryを使用します")
	} else {
		log.Println("✅ データベースに接続しました")
		// PostgreSQLリポジトリを使用
		userRepo = repository.NewUserRepository(db)
	}

	// RSA鍵ペアの生成（本番環境では事前に生成した鍵を読み込むこと）
	if err := jwt.GenerateRSAKeys(); err != nil {
		log.Fatalf("RSA鍵生成エラー: %v", err)
	}

	log.Println("RSA鍵ペアを生成しました")

	// サービスの初期化
	authService := service.NewAuthService(userRepo)

	// ハンドラーの初期化
	authHandler := handler.NewAuthHandler(authService)

	// ルーターのセットアップ
	r := router.SetupRouter(db, authHandler, storagePath)

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
