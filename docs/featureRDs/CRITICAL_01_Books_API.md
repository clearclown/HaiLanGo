# 🚨 CRITICAL - Books API Implementation

**優先度**: P0 - CRITICAL
**担当者**: Backend Engineer
**見積もり**: 4-6時間
**期限**: 即座
**ブロッカー**: フロントエンドが実装済みで現在失敗中

## 現状の問題

❌ **Books APIが未実装のため、フロントエンドの本棚機能が完全に動作していない**
- フロントエンドは `/api/v1/books` を呼び出すが404エラー
- 本の作成、一覧表示、削除がすべて失敗
- E2Eテストが失敗（20%がこの影響）

## 実装要件

### 1. APIエンドポイント

#### 1.1 本の一覧取得
```
GET /api/v1/books
Headers: Authorization: Bearer {token}
Response 200:
{
  "books": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "title": "ロシア語入門",
      "target_language": "ru",
      "native_language": "ja",
      "reference_language": "ja",
      "cover_image_url": "https://...",
      "total_pages": 150,
      "processed_pages": 45,
      "status": "ready",
      "ocr_status": "completed",
      "created_at": "2025-11-14T10:00:00Z",
      "updated_at": "2025-11-14T12:00:00Z"
    }
  ]
}
```

#### 1.2 本の詳細取得
```
GET /api/v1/books/:id
Headers: Authorization: Bearer {token}
Response 200: (同じBook object)
Response 404: { "error": "Book not found" }
```

#### 1.3 本の作成
```
POST /api/v1/books
Headers: Authorization: Bearer {token}
Content-Type: application/json
Body:
{
  "title": "ロシア語入門",
  "target_language": "ru",
  "native_language": "ja",
  "reference_language": "ja"
}
Response 201:
{
  "book": {
    "id": "uuid",
    "user_id": "uuid",
    "title": "ロシア語入門",
    "target_language": "ru",
    "native_language": "ja",
    "reference_language": "ja",
    "total_pages": 0,
    "processed_pages": 0,
    "status": "uploading",
    "ocr_status": "pending",
    "created_at": "2025-11-14T10:00:00Z",
    "updated_at": "2025-11-14T10:00:00Z"
  }
}
Response 400: { "error": "Invalid request body" }
```

#### 1.4 本の削除
```
DELETE /api/v1/books/:id
Headers: Authorization: Bearer {token}
Response 200: { "success": true }
Response 404: { "error": "Book not found" }
Response 403: { "error": "Forbidden" }
```

### 2. データベーススキーマ

```sql
CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    target_language VARCHAR(10) NOT NULL,
    native_language VARCHAR(10) NOT NULL,
    reference_language VARCHAR(10),
    cover_image_url TEXT,
    total_pages INTEGER DEFAULT 0,
    processed_pages INTEGER DEFAULT 0,
    status VARCHAR(50) DEFAULT 'uploading' CHECK (status IN ('uploading', 'processing', 'ready', 'failed')),
    ocr_status VARCHAR(50) DEFAULT 'pending' CHECK (ocr_status IN ('pending', 'processing', 'completed', 'failed')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
);
```

### 3. 実装コード（handler/books.go）

```go
package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
)

type BooksHandler struct {
	repo repository.BookRepository
}

func NewBooksHandler(repo repository.BookRepository) *BooksHandler {
	return &BooksHandler{repo: repo}
}

// GetBooks godoc
// @Summary Get all books for user
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string][]models.Book
// @Failure 401 {object} map[string]string
// @Router /api/v1/books [get]
func (h *BooksHandler) GetBooks(c *gin.Context) {
	userID := c.GetString("user_id") // ミドルウェアからユーザーIDを取得

	books, err := h.repo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"books": books})
}

// GetBook godoc
// @Summary Get book by ID
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Success 200 {object} models.Book
// @Failure 404 {object} map[string]string
// @Router /api/v1/books/{id} [get]
func (h *BooksHandler) GetBook(c *gin.Context) {
	bookID := c.Param("id")
	userID := c.GetString("user_id")

	book, err := h.repo.FindByID(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// ユーザー所有権チェック
	if book.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// CreateBook godoc
// @Summary Create new book
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param book body models.CreateBookRequest true "Book data"
// @Success 201 {object} map[string]models.Book
// @Failure 400 {object} map[string]string
// @Router /api/v1/books [post]
func (h *BooksHandler) CreateBook(c *gin.Context) {
	var req models.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID := c.GetString("user_id")

	book := &models.Book{
		UserID:            userID,
		Title:             req.Title,
		TargetLanguage:    req.TargetLanguage,
		NativeLanguage:    req.NativeLanguage,
		ReferenceLanguage: req.ReferenceLanguage,
		Status:            "uploading",
		OCRStatus:         "pending",
	}

	if err := h.repo.Create(c.Request.Context(), book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"book": book})
}

// DeleteBook godoc
// @Summary Delete book
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Success 200 {object} map[string]bool
// @Failure 404 {object} map[string]string
// @Router /api/v1/books/{id} [delete]
func (h *BooksHandler) DeleteBook(c *gin.Context) {
	bookID := c.Param("id")
	userID := c.GetString("user_id")

	// 所有権チェック
	book, err := h.repo.FindByID(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if book.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), bookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// RegisterRoutes registers book routes
func (h *BooksHandler) RegisterRoutes(rg *gin.RouterGroup) {
	books := rg.Group("/books")
	books.Use(middleware.AuthRequired()) // 認証必須
	{
		books.GET("", h.GetBooks)
		books.POST("", h.CreateBook)
		books.GET("/:id", h.GetBook)
		books.DELETE("/:id", h.DeleteBook)
	}
}
```

### 4. Repository Interface（repository/book.go）

```go
package repository

import (
	"context"
	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

type BookRepository interface {
	Create(ctx context.Context, book *models.Book) error
	FindByID(ctx context.Context, id string) (*models.Book, error)
	FindByUserID(ctx context.Context, userID string) ([]*models.Book, error)
	Update(ctx context.Context, book *models.Book) error
	Delete(ctx context.Context, id string) error
}

type bookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) BookRepository {
	return &bookRepository{db: db}
}

// 実装省略（標準的なCRUD操作）
```

### 5. Router Integration（router/router.go）

```go
// SetupRouter 内に追加
bookRepo := repository.NewBookRepository(db)
booksHandler := handler.NewBooksHandler(bookRepo)

// Books API
booksHandler.RegisterRoutes(v1)
```

### 6. テストケース

```go
func TestBooksHandler_GetBooks(t *testing.T) {
	// モックDBをセットアップ
	// ユーザーIDを設定
	// GetBooks を呼び出し
	// レスポンスが正しいことを検証
}

func TestBooksHandler_CreateBook(t *testing.T) {
	// 有効なリクエストボディで CreateBook を呼び出し
	// 201 Created が返ることを検証
	// DBに保存されていることを検証
}

func TestBooksHandler_DeleteBook_Forbidden(t *testing.T) {
	// 他のユーザーの本を削除しようとする
	// 403 Forbidden が返ることを検証
}
```

## 完了条件（Definition of Done）

- [ ] `handler/books.go` ファイルが作成され、すべてのメソッドが実装されている
- [ ] `repository/book.go` が実装されている
- [ ] `router/router.go` にルートが登録されている
- [ ] データベーススキーマが適用されている
- [ ] すべてのエンドポイントが動作する（Postmanでテスト）
- [ ] ユニットテストが書かれ、すべてパスする
- [ ] フロントエンドの Books ページが正常に動作する
- [ ] E2Eテストが成功する（books.spec.ts）

## 検証方法

### 1. Postman / cURL テスト
```bash
# 本の作成
curl -X POST http://localhost:8080/api/v1/books \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Book","target_language":"ru","native_language":"ja"}'

# 本の一覧取得
curl -X GET http://localhost:8080/api/v1/books \
  -H "Authorization: Bearer {token}"

# 本の削除
curl -X DELETE http://localhost:8080/api/v1/books/{id} \
  -H "Authorization: Bearer {token}"
```

### 2. フロントエンド動作確認
- http://localhost:3000/books にアクセス
- 「本を追加」ボタンが動作する
- 本のリストが表示される
- 削除ボタンが動作する

### 3. E2Eテスト実行
```bash
cd frontend/web
pnpm playwright test books.spec.ts
# すべてのテストがパスすること
```

## 注意事項

**❌ 絶対にやってはいけないこと:**
- ハードコードされた値を使用する
- エラーハンドリングを省略する
- ユーザー認証チェックを省略する
- テストを書かない

**✅ 必ず守ること:**
- すべてのエラーケースをハンドル
- ログを適切に出力
- トランザクションを使用（DB操作）
- 入力値のバリデーション

## 質問・不明点

不明点がある場合は即座にPMに確認すること。**推測で実装するな。**
