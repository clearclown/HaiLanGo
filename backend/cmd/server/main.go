package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	host := getEnv("SERVER_HOST", "0.0.0.0")
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

	// JWT署名鍵の初期化（本番では永続化した鍵を読み込むこと）
	privPath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	pubPath := os.Getenv("JWT_PUBLIC_KEY_PATH")
	if privPath != "" || pubPath != "" {
		if privPath == "" || pubPath == "" {
			log.Fatal("JWT_PRIVATE_KEY_PATH と JWT_PUBLIC_KEY_PATH は両方設定してください")
		}
		if err := jwt.LoadRSAKeysFromFiles(privPath, pubPath); err != nil {
			log.Fatalf("JWT鍵読み込みエラー: %v", err)
		}
		log.Println("✅ JWT RSA鍵をファイルから読み込みました")
	} else {
		if err := jwt.GenerateRSAKeys(); err != nil {
			log.Fatalf("RSA鍵生成エラー: %v", err)
		}
		log.Println("⚠️  JWT RSA鍵を起動時に生成しました（開発向け）。本番では JWT_PRIVATE_KEY_PATH/JWT_PUBLIC_KEY_PATH を設定してください。")
	}

	// サービスの初期化
	authService := service.NewAuthService(userRepo)

	// ハンドラーの初期化
	authHandler := handler.NewAuthHandler(authService)

	// ルーターのセットアップ
	r := router.SetupRouter(db, authHandler, storagePath)

	// サーバー起動（タイムアウト設定で Slowloris 等に耐性を付ける）
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("HaiLanGo APIサーバーを起動します: %s", addr)
	log.Printf("ストレージパス: %s", storagePath)

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	// サーバーはゴルーチンで起動し、シグナルでgraceful shutdownする
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("サーバー起動エラー: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 シャットダウンシグナルを受信しました。サーバーを停止します...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("サーバー停止エラー: %v", err)
	}

	log.Println("✅ サーバーを正常に停止しました")
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
