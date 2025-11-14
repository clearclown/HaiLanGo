# 🚨 CRITICAL - Router Integration (All Handlers)

**優先度**: P0 - CRITICAL
**担当者**: Backend Engineer
**見積もり**: 2-3時間
**期限**: 即座
**ブロッカー**: すべてのAPI機能がルーターに未登録

## 現状の問題

❌ **多くのハンドラーが実装されているのにルーターに登録されていない**

**現在のrouter.go:**
- ✅ Auth（認証）- 登録済み
- ✅ Upload（アップロード）- 登録済み
- ❌ **Books** - 未登録
- ❌ **Review** - 未登録
- ❌ **OCR** - 未登録
- ❌ **Learning** - 未登録
- ❌ **Pattern** - 未登録
- ❌ **Teacher Mode** - 未登録
- ❌ **Payment** - 未登録
- ❌ **WebSocket** - 未登録
- ❌ **Stats** - 未登録
- ❌ **Dictionary** - 未登録

**これは基本的な統合作業の怠慢。即座に修正せよ。**

## 実装要件

### 1. 既存ハンドラーの確認

```bash
# 確認されたハンドラーファイル
backend/internal/api/handler/auth.go             # ✅ 登録済み
backend/internal/api/handler/upload.go           # ✅ 登録済み
backend/internal/api/handler/stats.go            # ❌ 未登録
backend/internal/api/handler/dictionary.go       # ❌ 未登録
backend/internal/api/handler/pattern_handler.go  # ❌ 未登録
backend/internal/api/handler/review_handler.go   # ❌ 未登録
backend/internal/api/ocr/handler.go              # ❌ 未登録
backend/internal/api/learning/handler.go         # ❌ 未登録
backend/internal/api/payment/handler.go          # ❌ 未登録
backend/internal/api/teacher-mode/handler.go     # ❌ 未登録
backend/internal/api/websocket/handler.go        # ❌ 未登録
```

### 2. router.go の完全版実装

```go
package router

import (
	"database/sql"

	"github.com/clearclown/HaiLanGo/backend/internal/api/handler"
	"github.com/clearclown/HaiLanGo/backend/internal/api/learning"
	"github.com/clearclown/HaiLanGo/backend/internal/api/middleware"
	"github.com/clearclown/HaiLanGo/backend/internal/api/ocr"
	"github.com/clearclown/HaiLanGo/backend/internal/api/payment"
	"github.com/clearclown/HaiLanGo/backend/internal/api/teachermode"
	"github.com/clearclown/HaiLanGo/backend/internal/api/websocket"
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
	bookRepo := repository.NewBookRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	statsRepo := repository.NewStatsRepository(db)
	// 他のリポジトリも必要に応じて追加

	// ========================================
	// サービスの初期化
	// ========================================
	uploadService := service.NewUploadService(localStorage, tempDir)
	ocrService := service.NewOCRService()
	// 他のサービスも必要に応じて追加

	// ========================================
	// ハンドラーの初期化
	// ========================================
	uploadHandler := handler.NewUploadHandler(uploadService)
	booksHandler := handler.NewBooksHandler(bookRepo)
	reviewHandler := handler.NewReviewHandler(reviewRepo)
	statsHandler := handler.NewStatsHandler(statsRepo)
	dictionaryHandler := handler.NewDictionaryHandler()
	patternHandler := handler.NewPatternHandler()
	ocrHandler := ocr.NewOCRHandler(ocrService)
	learningHandler := learning.NewLearningHandler()
	paymentHandler := payment.NewPaymentHandler()
	teacherModeHandler := teachermode.NewTeacherModeHandler()

	// WebSocketハブを初期化
	wsHub := websocket.NewHub()
	go wsHub.Run()
	wsHandler := websocket.NewHandler(wsHub)

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

			// Review API
			reviewHandler.RegisterRoutes(authenticated)

			// Stats API
			statsHandler.RegisterRoutes(authenticated)

			// Upload API
			uploadHandler.RegisterRoutes(authenticated)

			// OCR API
			ocrHandler.RegisterRoutes(authenticated)

			// Learning API
			learningHandler.RegisterRoutes(authenticated)

			// Pattern API
			patternHandler.RegisterRoutes(authenticated)

			// Teacher Mode API
			teacherModeHandler.RegisterRoutes(authenticated)

			// Dictionary API
			dictionaryHandler.RegisterRoutes(authenticated)

			// Payment API（一部エンドポイントは認証不要の場合あり）
			paymentHandler.RegisterRoutes(authenticated)

			// WebSocket API
			authenticated.GET("/ws", wsHandler.HandleWebSocket)
		}
	}

	return r
}
```

### 3. 各ハンドラーの RegisterRoutes 実装確認

**すべてのハンドラーは `RegisterRoutes(rg *gin.RouterGroup)` メソッドを実装すること。**

#### Books Handler
```go
func (h *BooksHandler) RegisterRoutes(rg *gin.RouterGroup) {
	books := rg.Group("/books")
	{
		books.GET("", h.GetBooks)
		books.POST("", h.CreateBook)
		books.GET("/:id", h.GetBook)
		books.DELETE("/:id", h.DeleteBook)
	}
}
```

#### Review Handler
```go
func (h *ReviewHandler) RegisterRoutes(rg *gin.RouterGroup) {
	review := rg.Group("/review")
	{
		review.GET("/stats", h.GetStats)
		review.GET("/items", h.GetItems)
		review.POST("/submit", h.SubmitReview)
	}
}
```

#### Stats Handler
```go
func (h *StatsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	stats := rg.Group("/stats")
	{
		stats.GET("", h.GetStats)
		stats.GET("/learning", h.GetLearningStats)
		stats.GET("/weak-words", h.GetWeakWords)
	}
}
```

#### OCR Handler
```go
func (h *OCRHandler) RegisterRoutes(rg *gin.RouterGroup) {
	ocr := rg.Group("/ocr")
	{
		ocr.POST("/process", h.ProcessImage)
		ocr.GET("/result/:id", h.GetResult)
		ocr.POST("/edit", h.EditResult)
	}
}
```

#### Learning Handler
```go
func (h *LearningHandler) RegisterRoutes(rg *gin.RouterGroup) {
	learning := rg.Group("/learning")
	{
		learning.GET("/page", h.GetPage)
		learning.POST("/audio", h.GenerateAudio)
		learning.POST("/pronunciation", h.EvaluatePronunciation)
	}
}
```

#### Pattern Handler
```go
func (h *PatternHandler) RegisterRoutes(rg *gin.RouterGroup) {
	patterns := rg.Group("/patterns")
	{
		patterns.GET("/:bookId", h.GetPatterns)
		patterns.POST("/extract", h.ExtractPatterns)
	}
}
```

#### Teacher Mode Handler
```go
func (h *TeacherModeHandler) RegisterRoutes(rg *gin.RouterGroup) {
	teacher := rg.Group("/teacher-mode")
	{
		teacher.POST("/generate", h.GeneratePlaylist)
		teacher.GET("/playlist/:id", h.GetPlaylist)
		teacher.POST("/download-package", h.CreateDownloadPackage)
	}
}
```

#### Dictionary Handler
```go
func (h *DictionaryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	dictionary := rg.Group("/dictionary")
	{
		dictionary.GET("/lookup/:word", h.Lookup)
		dictionary.GET("/examples/:word", h.GetExamples)
	}
}
```

#### Payment Handler
```go
func (h *PaymentHandler) RegisterRoutes(rg *gin.RouterGroup) {
	payment := rg.Group("/payment")
	{
		payment.POST("/create-checkout", h.CreateCheckoutSession)
		payment.GET("/success", h.HandleSuccess)
		payment.GET("/cancel", h.HandleCancel)
	}

	// Webhook（認証不要）は別途登録
	// rg.POST("/webhook/stripe", h.HandleStripeWebhook)
}
```

### 4. main.go の更新

```go
package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/clearclown/HaiLanGo/backend/internal/api/handler"
	"github.com/clearclown/HaiLanGo/backend/internal/api/router"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/clearclown/HaiLanGo/backend/pkg/database"

	_ "github.com/lib/pq"
)

func main() {
	// 環境変数から設定を読み込み
	dbURL := os.Getenv("DATABASE_URL")
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}

	// データベース接続
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// マイグレーション実行
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 認証ハンドラーの初期化
	userRepo := repository.NewUserRepository(db)
	authHandler := handler.NewAuthHandler(userRepo)

	// ルーターのセットアップ
	r := router.SetupRouter(db, authHandler, storagePath)

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
```

### 5. 不足しているハンドラーの作成

以下のハンドラーファイルが存在しない場合は作成すること：

- `handler/books.go` - CRITICAL_01_Books_API.md参照
- `handler/stats.go` - 既に存在するが実装が不完全な可能性あり
- `handler/dictionary.go` - 既に存在するが実装が不完全な可能性あり

### 6. テスト方法

```bash
# サーバー起動
cd backend
go run cmd/server/main.go

# 各エンドポイントの確認
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/books -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/review/stats -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/stats -H "Authorization: Bearer {token}"
# ... その他すべてのエンドポイント
```

## 完了条件（Definition of Done）

- [ ] `router/router.go` が完全に更新されている
- [ ] すべてのハンドラーが `RegisterRoutes` メソッドを実装している
- [ ] すべてのリポジトリが初期化されている
- [ ] すべてのエンドポイントが `/health` と同様に応答する
- [ ] `curl` ですべてのエンドポイントにアクセスできる（認証トークンあり）
- [ ] フロントエンドからのAPI呼び出しが成功する
- [ ] ログに404エラーが出ない

## 検証チェックリスト

### 起動確認
- [ ] サーバーが起動する（エラーなし）
- [ ] `/health` エンドポイントが応答する

### 各APIエンドポイント確認
- [ ] `GET /api/v1/books` - 200または401
- [ ] `POST /api/v1/books` - 201または400/401
- [ ] `GET /api/v1/review/stats` - 200または401
- [ ] `GET /api/v1/review/items` - 200または401
- [ ] `POST /api/v1/review/submit` - 200または400/401
- [ ] `GET /api/v1/stats` - 200または401
- [ ] `POST /api/v1/ocr/process` - 200または400/401
- [ ] `GET /api/v1/learning/page` - 200または404/401
- [ ] `GET /api/v1/patterns/:bookId` - 200または404/401
- [ ] `POST /api/v1/teacher-mode/generate` - 200または400/401
- [ ] `GET /api/v1/dictionary/lookup/:word` - 200または404/401
- [ ] `POST /api/v1/payment/create-checkout` - 200または400/401
- [ ] `GET /api/v1/ws` - WebSocket接続成功

### フロントエンド統合確認
- [ ] Books ページが動作する
- [ ] Review ページが動作する
- [ ] Upload ページが動作する
- [ ] Settings ページが動作する

## 注意事項

**❌ 絶対にやってはいけないこと:**
- ルートを登録せずにハンドラーだけ実装する
- パニックを引き起こすコードを書く
- 認証ミドルウェアを忘れる

**✅ 必ず守ること:**
- すべてのハンドラーをルーターに登録
- 適切なミドルウェアを適用
- エラーハンドリングを実装
- ログを適切に出力

## 期限

**本日中に完了させること。これ以上の遅延は許されない。**
