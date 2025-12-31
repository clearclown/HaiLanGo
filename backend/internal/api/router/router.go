package router

import (
	"database/sql"
	"log"

	"github.com/clearclown/HaiLanGo/backend/internal/api/handler"
	"github.com/clearclown/HaiLanGo/backend/internal/api/middleware"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/clearclown/HaiLanGo/backend/internal/service"
	ocrservice "github.com/clearclown/HaiLanGo/backend/internal/service/ocr"
	"github.com/clearclown/HaiLanGo/backend/internal/websocket"
	"github.com/clearclown/HaiLanGo/backend/pkg/cache"
	"github.com/clearclown/HaiLanGo/backend/pkg/ocr"
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
	r := gin.New()

	// ミドルウェアの設定
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger())
	r.Use(gin.Recovery())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimiter())

	// ストレージを初期化
	localStorage := storage.NewLocalStorage(storagePath)
	tempDir := storagePath + "/temp"

	// ========================================
	// リポジトリの初期化（データベース接続に応じてフォールバック）
	// ========================================

	// BookRepository: InMemory fallback対応
	var bookRepo repository.BookRepository
	if err := db.Ping(); err != nil {
		log.Printf("⚠️  データベースPing失敗 (BookRepository: InMemoryを使用): %v", err)
		bookRepo = repository.NewInMemoryBookRepository()
		log.Println("✅ InMemoryBookRepositoryを使用します")
	} else {
		bookRepo = repository.NewBookRepositoryPostgres(db)
	}

	// 他のリポジトリ: InMemory fallback対応
	var reviewRepo repository.ReviewRepository
	var statsRepo repository.StatsRepositoryInterface
	var learningRepo repository.LearningRepositoryInterface
	var ocrRepo repository.OCRRepositoryInterface
	var ttsRepo repository.TTSRepositoryInterface
	var sttRepo repository.STTRepositoryInterface
	var paymentRepo repository.PaymentRepositoryInterface
	var dictionaryRepo repository.DictionaryRepositoryInterface
	var patternRepo repository.PatternRepositoryInterface

	if err := db.Ping(); err != nil {
		log.Println("⚠️  データベース接続失敗 - すべてのリポジトリでInMemory実装を使用します")
		reviewRepo = repository.NewInMemoryReviewRepository()
		statsRepo = repository.NewInMemoryStatsRepository()
		learningRepo = repository.NewInMemoryLearningRepository()
		ocrRepo = repository.NewInMemoryOCRRepository()
		ttsRepo = repository.NewInMemoryTTSRepository()
		sttRepo = repository.NewInMemorySTTRepository()
		paymentRepo = repository.NewInMemoryPaymentRepository()
		dictionaryRepo = repository.NewInMemoryDictionaryRepository()
		patternRepo = repository.NewInMemoryPatternRepository()
	} else {
		reviewRepo = repository.NewReviewRepositoryPostgres(db)
		statsRepo = repository.NewStatsRepository(db)
		learningRepo = repository.NewLearningRepositoryPostgres(db)
		ocrRepo = repository.NewOCRRepositoryPostgres(db)
		ttsRepo = repository.NewTTSRepositoryPostgres(db)
		sttRepo = repository.NewSTTRepositoryPostgres(db)
		paymentRepo = repository.NewPaymentRepositoryPostgres(db)
		dictionaryRepo = repository.NewDictionaryRepositoryPostgres(db)
		patternRepo = repository.NewPatternRepositoryPostgres(db)
	}

	// PageRepository と TeacherModeRepository: InMemory fallback対応
	var pageRepo repository.PageRepository
	var teacherModeRepo repository.TeacherModeRepository
	if err := db.Ping(); err != nil {
		log.Println("⚠️  データベース接続失敗 - PageRepository, TeacherModeRepositoryでInMemory実装を使用します")
		pageRepo = repository.NewInMemoryPageRepository()
		teacherModeRepo = repository.NewInMemoryTeacherModeRepository()
	} else {
		pageRepo = repository.NewPageRepositoryPostgres(db)
		teacherModeRepo = repository.NewTeacherModeRepositoryPostgres(db)
	}

	// ========================================
	// サービスの初期化
	// ========================================
	uploadService := service.NewUploadService(localStorage, bookRepo, tempDir)
	teacherModeService := service.NewTeacherModeService(teacherModeRepo, pageRepo, bookRepo, ttsRepo)

	// OCRサービスの初期化
	ocrClient, err := ocr.NewOCRClient() // 環境変数に基づいて実際のAPIまたはモックを返す
	if err != nil {
		panic("Failed to initialize OCR client: " + err.Error())
	}
	// キャッシュの初期化（Redis優先、フォールバックでInMemory）
	appCache, err := cache.NewCache()
	if err != nil {
		log.Printf("⚠️  キャッシュ初期化失敗（InMemoryを使用）: %v", err)
		appCache = cache.NewInMemoryCache()
	}
	ocrSvc := ocrservice.NewOCRService(ocrClient, appCache)
	ocrSvc.SetPageRepository(pageRepo)

	// Note: Stats/Review handlers work directly with repositories
	// SRS algorithm is embedded in ReviewHandler (service.SM2Algorithm)

	// WebSocketハブを初期化（先に初期化してサービスで使用できるようにする）
	wsHub := websocket.NewHub()
	go wsHub.Run()
	wsHandler := handler.NewWebSocketHandler(wsHub)

	// OCRサービスにWebSocketハブを設定
	ocrSvc.SetWebSocketHub(wsHub)

	// ========================================
	// ハンドラーの初期化
	// ========================================
	uploadHandler := handler.NewUploadHandler(uploadService)
	booksHandler := handler.NewBooksHandler(bookRepo, wsHub)
	reviewHandler := handler.NewReviewHandler(reviewRepo, wsHub)
	statsHandler := handler.NewStatsHandler(statsRepo)
	learningHandler := handler.NewLearningHandler(learningRepo)
	ocrHandler := handler.NewOCRHandler(ocrRepo, ocrSvc, wsHub)
	ttsHandler := handler.NewTTSHandler(ttsRepo)
	sttHandler := handler.NewSTTHandler(sttRepo)
	paymentHandler := handler.NewPaymentHandler(paymentRepo)
	dictionaryHandler := handler.NewDictionaryHandler(dictionaryRepo)
	patternHandler := handler.NewPatternHandler(patternRepo)
	teacherModeHandler := handler.NewTeacherModeHandler(teacherModeService)

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
			ocrHandler.RegisterRoutes(authenticated)

			// TTS API
			ttsHandler.RegisterRoutes(authenticated)

			// STT API
			sttHandler.RegisterRoutes(authenticated)

			// Payment API
			paymentHandler.RegisterRoutes(authenticated)

			// Dictionary API
			dictionaryHandler.RegisterRoutes(authenticated)

			// Pattern API
			patternHandler.RegisterRoutes(authenticated)

			// Teacher Mode API
			teacherModeHandler.RegisterRoutes(authenticated)

			// WebSocket API
			wsHandler.RegisterRoutes(authenticated)

			// WebSocket統計（デバッグ用）
			authenticated.GET("/ws/stats", wsHandler.GetStats)
		}
	}

	return r
}
