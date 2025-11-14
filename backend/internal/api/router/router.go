package router

import (
	"context"
	"database/sql"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/api/handler"
	"github.com/clearclown/HaiLanGo/backend/internal/api/learning"
	"github.com/clearclown/HaiLanGo/backend/internal/api/middleware"
	"github.com/clearclown/HaiLanGo/backend/internal/api/ocr"
	"github.com/clearclown/HaiLanGo/backend/internal/api/payment"
	teachermode "github.com/clearclown/HaiLanGo/backend/internal/api/teacher-mode"
	"github.com/clearclown/HaiLanGo/backend/internal/api/websocket"
	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/clearclown/HaiLanGo/backend/internal/service"
	dictionaryService "github.com/clearclown/HaiLanGo/backend/internal/service/dictionary"
	ocrService "github.com/clearclown/HaiLanGo/backend/internal/service/ocr"
	paymentService "github.com/clearclown/HaiLanGo/backend/internal/service/payment"
	statsService "github.com/clearclown/HaiLanGo/backend/internal/service/stats"
	"github.com/clearclown/HaiLanGo/backend/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	bookRepo := repository.NewInMemoryBookRepository()
	reviewRepo := repository.NewInMemoryReviewRepository()

	// ========================================
	// サービスの初期化
	// ========================================
	uploadService := service.NewUploadService(localStorage, tempDir)

	// Create services/handlers - using simple mocks where needed
	ocrEditorService := ocrService.NewEditorService(nil, nil)
	statsServiceMock := &mockStatsService{}
	dictionaryServiceInstance, _ := dictionaryService.NewService()
	paymentServiceMock := &mockPaymentService{}
	learningServiceMock := &mockLearningService{}

	// ========================================
	// ハンドラーの初期化
	// ========================================
	uploadHandler := handler.NewUploadHandler(uploadService)
	booksHandler := handler.NewBooksHandler(bookRepo)
	reviewHandler := handler.NewReviewHandler(reviewRepo)
	statsHandler := handler.NewStatsHandler(statsServiceMock)
	dictionaryHandler := handler.NewDictionaryHandler(dictionaryServiceInstance)
	patternHandler := handler.NewPatternHandler()
	ocrHandler := ocr.NewHandler(ocrEditorService)
	learningHandler := learning.NewHandler(learningServiceMock)
	paymentHandler := payment.NewHandler(paymentServiceMock)
	teacherModeHandler := teachermode.NewHandler()

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

			// Upload API
			uploadHandler.RegisterRoutes(authenticated)

			// Review API
			reviewHandler.RegisterRoutes(authenticated)

			// Stats API
			statsHandler.RegisterRoutes(authenticated)

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

			// Payment API
			paymentHandler.RegisterRoutes(authenticated)

			// WebSocket API
			wsHandler.RegisterRoutes(authenticated)
		}
	}

	return r
}

// mockStatsService is a temporary mock implementation
type mockStatsService struct{}

func (m *mockStatsService) GetDashboardStats(ctx context.Context, userID uuid.UUID) (*statsService.DashboardStats, error) {
	return &statsService.DashboardStats{
		CompletedBooks:   5,
		TotalLearningTime: 7200,
		CurrentStreak:    7,
	}, nil
}

func (m *mockStatsService) GetLearningTimeStats(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	return map[string]int{"today": 3600}, nil
}

func (m *mockStatsService) GetProgressStats(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	return map[string]int{"pages": 45}, nil
}

func (m *mockStatsService) GetStreakStats(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	return map[string]int{"current": 7}, nil
}

func (m *mockStatsService) GetLearningTimeChart(ctx context.Context, userID uuid.UUID, days int) (interface{}, error) {
	return []map[string]interface{}{{"date": "2025-01-01", "time": 3600}}, nil
}

func (m *mockStatsService) GetProgressChart(ctx context.Context, userID uuid.UUID, days int) (interface{}, error) {
	return []map[string]interface{}{{"date": "2025-01-01", "progress": 10}}, nil
}

func (m *mockStatsService) GetWeakWords(ctx context.Context, userID uuid.UUID, limit int) ([]interface{}, error) {
	return []interface{}{}, nil
}

// mockPaymentService is a temporary mock implementation
type mockPaymentService struct{}

type mockSubscription struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	PlanID   uuid.UUID
	Status   string
	CreatedAt time.Time
}

type mockPlan struct {
	ID    uuid.UUID
	Name  string
	Price float64
}

func (m *mockPaymentService) CreateSubscription(ctx context.Context, userID, planID uuid.UUID) (*mockSubscription, error) {
	return &mockSubscription{
		ID:       uuid.New(),
		UserID:   userID,
		PlanID:   planID,
		Status:   "active",
		CreatedAt: time.Now(),
	}, nil
}

func (m *mockPaymentService) ListPlans(ctx context.Context) ([]*mockPlan, error) {
	return []*mockPlan{
		{
			ID:    uuid.New(),
			Name:  "Premium",
			Price: 9.99,
		},
	}, nil
}

func (m *mockPaymentService) GetSubscription(ctx context.Context, id uuid.UUID) (*mockSubscription, error) {
	return &mockSubscription{
		ID:     id,
		Status: "active",
	}, nil
}

func (m *mockPaymentService) CancelSubscription(ctx context.Context, subscriptionID uuid.UUID, cancelAtPeriodEnd bool) error {
	return nil
}

// mockLearningService is a temporary mock implementation
type mockLearningService struct{}

func (m *mockLearningService) GetPage(ctx context.Context, bookID uuid.UUID, pageNumber int) (*models.PageWithProgress, error) {
	return &models.PageWithProgress{
		Page: &models.Page{
			ID:            uuid.New(),
			BookID:        bookID,
			PageNumber:    pageNumber,
			ImageURL:      "https://example.com/page.jpg",
			OCRText:       "Sample OCR text",
			OCRConfidence: 0.95,
			DetectedLang:  "en",
			OCRStatus:     models.OCRStatusCompleted,
		},
		IsCompleted: false,
	}, nil
}

func (m *mockLearningService) MarkPageCompleted(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, studyTime int) error {
	return nil
}

func (m *mockLearningService) GetProgress(ctx context.Context, userID, bookID uuid.UUID) (*models.LearningProgress, error) {
	return &models.LearningProgress{
		BookID:         bookID,
		CompletedPages: 5,
		TotalPages:     100,
	}, nil
}
