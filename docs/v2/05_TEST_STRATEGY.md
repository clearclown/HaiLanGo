# テスト戦略

## 概要

HaiLanGoのテスト戦略を定義する。
**Testcontainers**を使用して、実際のPostgreSQLとRedisに対するテストを実行する。

---

## テストピラミッド

```
                    ┌─────────┐
                    │  E2E    │  ← Playwright (少数)
                    │  Tests  │
                 ┌──┴─────────┴──┐
                 │  Integration  │  ← Testcontainers (中程度)
                 │    Tests      │
              ┌──┴───────────────┴──┐
              │     Unit Tests      │  ← 多数
              └─────────────────────┘
```

| レベル | 目的 | 実行時間 | カバレッジ目標 |
|--------|------|----------|----------------|
| Unit | 単一関数/メソッドの検証 | ~1秒 | 80%+ |
| Integration | DB/Redis含むサービス検証 | ~10秒 | 主要フロー |
| E2E | ユーザーフロー全体の検証 | ~1分 | クリティカルパス |

---

## Testcontainersセットアップ

### インストール

```bash
go get github.com/testcontainers/testcontainers-go
go get github.com/testcontainers/testcontainers-go/modules/postgres
go get github.com/testcontainers/testcontainers-go/modules/redis
```

### テストヘルパー

```go
// backend/tests/testhelpers/containers.go
package testhelpers

import (
    "context"
    "testing"
    "time"

    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/modules/redis"
    "github.com/testcontainers/testcontainers-go/wait"
)

// TestContainers holds all test containers
type TestContainers struct {
    PostgresContainer *postgres.PostgresContainer
    RedisContainer    *redis.RedisContainer
    PostgresURL       string
    RedisURL          string
}

// SetupTestContainers creates and starts test containers
func SetupTestContainers(t *testing.T) *TestContainers {
    ctx := context.Background()

    // PostgreSQL Container
    pgContainer, err := postgres.Run(ctx,
        "postgres:15-alpine",
        postgres.WithDatabase("hailango_test"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(60*time.Second),
        ),
    )
    if err != nil {
        t.Fatalf("Failed to start PostgreSQL container: %v", err)
    }

    pgURL, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
    if err != nil {
        t.Fatalf("Failed to get PostgreSQL connection string: %v", err)
    }

    // Redis Container
    redisContainer, err := redis.Run(ctx,
        "redis:7-alpine",
        testcontainers.WithWaitStrategy(
            wait.ForLog("Ready to accept connections").
                WithStartupTimeout(30*time.Second),
        ),
    )
    if err != nil {
        t.Fatalf("Failed to start Redis container: %v", err)
    }

    redisURL, err := redisContainer.ConnectionString(ctx)
    if err != nil {
        t.Fatalf("Failed to get Redis connection string: %v", err)
    }

    // Cleanup on test completion
    t.Cleanup(func() {
        if err := pgContainer.Terminate(ctx); err != nil {
            t.Logf("Failed to terminate PostgreSQL container: %v", err)
        }
        if err := redisContainer.Terminate(ctx); err != nil {
            t.Logf("Failed to terminate Redis container: %v", err)
        }
    })

    return &TestContainers{
        PostgresContainer: pgContainer,
        RedisContainer:    redisContainer,
        PostgresURL:       pgURL,
        RedisURL:          redisURL,
    }
}
```

### マイグレーション実行ヘルパー

```go
// backend/tests/testhelpers/migration.go
package testhelpers

import (
    "database/sql"
    "testing"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations executes all migrations on the test database
func RunMigrations(t *testing.T, dbURL string) {
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        t.Fatalf("Failed to open database: %v", err)
    }
    defer db.Close()

    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        t.Fatalf("Failed to create migration driver: %v", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://../../migrations",
        "postgres",
        driver,
    )
    if err != nil {
        t.Fatalf("Failed to create migrate instance: %v", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        t.Fatalf("Failed to run migrations: %v", err)
    }
}
```

---

## 統合テスト例

### リポジトリテスト

```go
// backend/internal/repository/postgres/book_repository_test.go
package postgres_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "hailango/internal/models"
    "hailango/internal/repository/postgres"
    "hailango/tests/testhelpers"
)

func TestBookRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    // Setup containers
    containers := testhelpers.SetupTestContainers(t)
    testhelpers.RunMigrations(t, containers.PostgresURL)

    // Create repository
    repo, err := postgres.NewBookRepository(containers.PostgresURL)
    require.NoError(t, err)

    ctx := context.Background()

    t.Run("Create and Get Book", func(t *testing.T) {
        // Create user first
        userRepo, _ := postgres.NewUserRepository(containers.PostgresURL)
        user := &models.User{
            Email:        "test@example.com",
            PasswordHash: "hash",
            DisplayName:  "Test User",
        }
        err := userRepo.Create(ctx, user)
        require.NoError(t, err)

        // Create book
        book := &models.Book{
            UserID:         user.ID,
            Title:          "ロシア語入門",
            TargetLanguage: "ru",
            NativeLanguage: "ja",
            Status:         "processing",
        }

        err = repo.Create(ctx, book)
        require.NoError(t, err)
        assert.NotEmpty(t, book.ID)

        // Get book
        found, err := repo.GetByID(ctx, book.ID)
        require.NoError(t, err)
        assert.Equal(t, book.Title, found.Title)
        assert.Equal(t, book.TargetLanguage, found.TargetLanguage)
    })

    t.Run("List Books by User", func(t *testing.T) {
        // Create multiple books for a user
        userRepo, _ := postgres.NewUserRepository(containers.PostgresURL)
        user := &models.User{
            Email:        "list-test@example.com",
            PasswordHash: "hash",
            DisplayName:  "List Test User",
        }
        userRepo.Create(ctx, user)

        for i := 0; i < 5; i++ {
            book := &models.Book{
                UserID:         user.ID,
                Title:          fmt.Sprintf("Book %d", i),
                TargetLanguage: "ru",
                NativeLanguage: "ja",
                Status:         "ready",
            }
            repo.Create(ctx, book)
        }

        // List books
        books, total, err := repo.ListByUserID(ctx, user.ID, 1, 10)
        require.NoError(t, err)
        assert.Equal(t, int64(5), total)
        assert.Len(t, books, 5)
    })
}
```

### サービステスト

```go
// backend/internal/service/book_service_test.go
package service_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "hailango/internal/service"
    "hailango/internal/repository/postgres"
    "hailango/tests/testhelpers"
)

func TestBookService_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    containers := testhelpers.SetupTestContainers(t)
    testhelpers.RunMigrations(t, containers.PostgresURL)

    // Setup repositories
    bookRepo, _ := postgres.NewBookRepository(containers.PostgresURL)
    userRepo, _ := postgres.NewUserRepository(containers.PostgresURL)

    // Create service
    bookService := service.NewBookService(bookRepo, userRepo)

    ctx := context.Background()

    t.Run("CreateBook with validation", func(t *testing.T) {
        // First create a user
        user := testhelpers.CreateTestUser(t, userRepo)

        // Test valid book creation
        book, err := bookService.CreateBook(ctx, user.ID, &service.CreateBookInput{
            Title:          "Valid Book",
            TargetLanguage: "ru",
            NativeLanguage: "ja",
        })
        require.NoError(t, err)
        assert.NotEmpty(t, book.ID)

        // Test invalid language
        _, err = bookService.CreateBook(ctx, user.ID, &service.CreateBookInput{
            Title:          "Invalid Book",
            TargetLanguage: "invalid",
            NativeLanguage: "ja",
        })
        assert.Error(t, err)
    })
}
```

### ハンドラーテスト

```go
// backend/internal/api/handler/book_handler_test.go
package handler_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "hailango/internal/api/handler"
    "hailango/internal/api/router"
    "hailango/tests/testhelpers"
)

func TestBookHandler_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    containers := testhelpers.SetupTestContainers(t)
    testhelpers.RunMigrations(t, containers.PostgresURL)

    // Setup app with test containers
    app := testhelpers.SetupTestApp(t, containers)

    t.Run("POST /api/v1/books", func(t *testing.T) {
        // Create test user and get token
        token := testhelpers.CreateUserAndGetToken(t, app)

        body := map[string]interface{}{
            "title":           "Test Book",
            "target_language": "ru",
            "native_language": "ja",
        }
        bodyBytes, _ := json.Marshal(body)

        req := httptest.NewRequest("POST", "/api/v1/books", bytes.NewReader(bodyBytes))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+token)

        rec := httptest.NewRecorder()
        app.Router.ServeHTTP(rec, req)

        assert.Equal(t, http.StatusCreated, rec.Code)

        var response map[string]interface{}
        json.Unmarshal(rec.Body.Bytes(), &response)

        assert.True(t, response["success"].(bool))
        data := response["data"].(map[string]interface{})
        assert.NotEmpty(t, data["id"])
        assert.Equal(t, "Test Book", data["title"])
    })

    t.Run("GET /api/v1/books", func(t *testing.T) {
        token := testhelpers.CreateUserAndGetToken(t, app)

        // Create some books first
        for i := 0; i < 3; i++ {
            testhelpers.CreateBookForUser(t, app, token)
        }

        req := httptest.NewRequest("GET", "/api/v1/books", nil)
        req.Header.Set("Authorization", "Bearer "+token)

        rec := httptest.NewRecorder()
        app.Router.ServeHTTP(rec, req)

        assert.Equal(t, http.StatusOK, rec.Code)

        var response map[string]interface{}
        json.Unmarshal(rec.Body.Bytes(), &response)

        data := response["data"].(map[string]interface{})
        books := data["books"].([]interface{})
        assert.Len(t, books, 3)
    })
}
```

---

## フロントエンドテスト

### Vitest 単体テスト

```typescript
// frontend/web/lib/api/__tests__/books.test.ts
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BooksAPI } from '../books';

describe('BooksAPI', () => {
  let api: BooksAPI;
  let mockFetch: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    mockFetch = vi.fn();
    global.fetch = mockFetch;
    api = new BooksAPI('http://localhost:8080');
  });

  describe('getBooks', () => {
    it('should fetch books with pagination', async () => {
      const mockResponse = {
        success: true,
        data: {
          books: [{ id: '1', title: 'Test Book' }],
          pagination: { page: 1, limit: 20, total: 1 },
        },
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      });

      const result = await api.getBooks({ page: 1, limit: 20 });

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/books?page=1&limit=20',
        expect.any(Object)
      );
      expect(result.books).toHaveLength(1);
    });
  });

  describe('createBook', () => {
    it('should create a new book', async () => {
      const mockResponse = {
        success: true,
        data: { id: '1', title: 'New Book' },
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      });

      const result = await api.createBook({
        title: 'New Book',
        targetLanguage: 'ru',
        nativeLanguage: 'ja',
      });

      expect(result.id).toBe('1');
    });
  });
});
```

### Playwright E2Eテスト

```typescript
// frontend/web/e2e/book-upload.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Book Upload Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Login first
    await page.goto('/login');
    await page.fill('[data-testid="email"]', 'test@example.com');
    await page.fill('[data-testid="password"]', 'password123');
    await page.click('[data-testid="login-button"]');
    await page.waitForURL('/');
  });

  test('should upload a new book', async ({ page }) => {
    // Navigate to upload page
    await page.click('[data-testid="add-book-button"]');
    await expect(page).toHaveURL('/upload');

    // Fill book details
    await page.fill('[data-testid="book-title"]', 'Test Russian Book');
    await page.selectOption('[data-testid="target-language"]', 'ru');
    await page.selectOption('[data-testid="native-language"]', 'ja');

    // Upload file
    const fileInput = page.locator('[data-testid="file-input"]');
    await fileInput.setInputFiles('tests/fixtures/sample.pdf');

    // Submit
    await page.click('[data-testid="upload-button"]');

    // Wait for upload completion
    await expect(page.locator('[data-testid="upload-progress"]')).toBeVisible();
    await expect(page.locator('[data-testid="upload-complete"]')).toBeVisible({
      timeout: 30000,
    });

    // Verify book appears in list
    await page.goto('/');
    await expect(page.locator('text=Test Russian Book')).toBeVisible();
  });

  test('should show OCR progress', async ({ page }) => {
    // Upload a book (abbreviated)
    // ...

    // Check OCR progress via WebSocket
    await expect(page.locator('[data-testid="ocr-progress"]')).toBeVisible();
    await expect(page.locator('[data-testid="ocr-status"]')).toContainText(
      'Processing'
    );
  });
});
```

---

## テスト実行

### バックエンド

```bash
# 全テスト実行
go test ./...

# 統合テスト込み（Testcontainers使用）
go test ./... -v

# 短いテストのみ（Testcontainersスキップ）
go test ./... -short

# カバレッジ
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# 特定のパッケージ
go test ./internal/repository/postgres/... -v

# 特定のテスト
go test ./internal/repository/postgres/... -run TestBookRepository -v
```

### フロントエンド

```bash
# 単体テスト
pnpm test

# ウォッチモード
pnpm test:watch

# カバレッジ
pnpm test:coverage

# E2Eテスト
pnpm test:e2e

# E2Eテスト（UIモード）
pnpm test:e2e:ui
```

---

## CI/CD設定

### GitHub Actions

```yaml
# .github/workflows/test.yml
name: Tests

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  backend-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Install dependencies
        run: |
          cd backend
          go mod download

      - name: Run unit tests
        run: |
          cd backend
          go test ./... -short -v

      - name: Run integration tests
        run: |
          cd backend
          go test ./... -v -coverprofile=coverage.out

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./backend/coverage.out

  frontend-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup pnpm
        uses: pnpm/action-setup@v2
        with:
          version: 8

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'pnpm'
          cache-dependency-path: frontend/web/pnpm-lock.yaml

      - name: Install dependencies
        run: |
          cd frontend/web
          pnpm install

      - name: Run tests
        run: |
          cd frontend/web
          pnpm test:coverage

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./frontend/web/coverage/lcov.info

  e2e-test:
    runs-on: ubuntu-latest
    needs: [backend-test, frontend-test]
    steps:
      - uses: actions/checkout@v4

      - name: Setup
        # ... setup steps

      - name: Start services
        run: |
          docker-compose -f docker-compose.test.yml up -d
          sleep 10

      - name: Run E2E tests
        run: |
          cd frontend/web
          pnpm test:e2e

      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: playwright-report
          path: frontend/web/playwright-report
```

---

## テストデータ管理

### Fixtures

```go
// backend/tests/fixtures/fixtures.go
package fixtures

import (
    "hailango/internal/models"
)

func SampleBook() *models.Book {
    return &models.Book{
        Title:          "ロシア語入門",
        TargetLanguage: "ru",
        NativeLanguage: "ja",
        Status:         "ready",
        TotalPages:     150,
    }
}

func SampleUser() *models.User {
    return &models.User{
        Email:          "test@example.com",
        PasswordHash:   "$2a$10$...", // bcrypt hash of "password123"
        DisplayName:    "Test User",
        NativeLanguage: "ja",
    }
}
```

### Factory

```go
// backend/tests/factories/factories.go
package factories

import (
    "context"
    "testing"

    "hailango/internal/models"
    "hailango/internal/repository"
)

type Factory struct {
    userRepo repository.UserRepository
    bookRepo repository.BookRepository
}

func (f *Factory) CreateUser(t *testing.T, overrides ...func(*models.User)) *models.User {
    user := &models.User{
        Email:        fmt.Sprintf("user-%s@example.com", uuid.New().String()),
        PasswordHash: "hash",
        DisplayName:  "Test User",
    }

    for _, override := range overrides {
        override(user)
    }

    err := f.userRepo.Create(context.Background(), user)
    if err != nil {
        t.Fatalf("Failed to create user: %v", err)
    }

    return user
}

func (f *Factory) CreateBook(t *testing.T, userID string, overrides ...func(*models.Book)) *models.Book {
    book := &models.Book{
        UserID:         userID,
        Title:          "Test Book",
        TargetLanguage: "ru",
        NativeLanguage: "ja",
        Status:         "ready",
    }

    for _, override := range overrides {
        override(book)
    }

    err := f.bookRepo.Create(context.Background(), book)
    if err != nil {
        t.Fatalf("Failed to create book: %v", err)
    }

    return book
}
```

---

## 次のドキュメント

- [06_IMPLEMENTATION_ROADMAP.md](./06_IMPLEMENTATION_ROADMAP.md) - 実装ロードマップ
