package repository

import (
	"context"
	"database/sql"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// BookRepository は書籍データアクセスのインターフェース
type BookRepository interface {
	// Create は新しい書籍を作成する
	Create(ctx context.Context, book *models.Book) error

	// GetByID はIDで書籍を取得する
	GetByID(ctx context.Context, id uuid.UUID) (*models.Book, error)

	// GetByUserID はユーザーIDで書籍一覧を取得する
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Book, error)

	// Update は書籍情報を更新する
	Update(ctx context.Context, book *models.Book) error

	// Delete は書籍を削除する
	Delete(ctx context.Context, id uuid.UUID) error

	// UpdateStatus は書籍のステータスを更新する
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.BookStatus) error
}

// BookFileRepository は書籍ファイルデータアクセスのインターフェース
type BookFileRepository interface {
	// Create は新しい書籍ファイルを作成する
	Create(ctx context.Context, file *models.BookFile) error

	// GetByBookID は書籍IDでファイル一覧を取得する
	GetByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.BookFile, error)

	// GetByID はIDでファイルを取得する
	GetByID(ctx context.Context, id uuid.UUID) (*models.BookFile, error)

	// Delete はファイルを削除する
	Delete(ctx context.Context, id uuid.UUID) error
}

// ChunkUploadRepository はチャンクアップロードデータアクセスのインターフェース
type ChunkUploadRepository interface {
	// Create は新しいチャンクアップロードを作成する
	Create(ctx context.Context, chunk *models.ChunkUpload) error

	// GetByID はIDでチャンクアップロードを取得する
	GetByID(ctx context.Context, id uuid.UUID) (*models.ChunkUpload, error)

	// GetByBookID は書籍IDでチャンクアップロードを取得する
	GetByBookID(ctx context.Context, bookID uuid.UUID) (*models.ChunkUpload, error)

	// Update はチャンクアップロード情報を更新する
	Update(ctx context.Context, chunk *models.ChunkUpload) error

	// Delete はチャンクアップロードを削除する
	Delete(ctx context.Context, id uuid.UUID) error

	// IncrementUploadedChunks はアップロード済みチャンク数をインクリメントする
	IncrementUploadedChunks(ctx context.Context, id uuid.UUID) error
}

// メモリ内実装（開発・テスト用）

// InMemoryBookRepository はメモリ内の書籍リポジトリ実装
type InMemoryBookRepository struct {
	books map[uuid.UUID]*models.Book
}

// NewInMemoryBookRepository は新しいInMemoryBookRepositoryを作成する
func NewInMemoryBookRepository() *InMemoryBookRepository {
	return &InMemoryBookRepository{
		books: make(map[uuid.UUID]*models.Book),
	}
}

// Create は書籍を作成する
func (r *InMemoryBookRepository) Create(ctx context.Context, book *models.Book) error {
	r.books[book.ID] = book
	return nil
}

// GetByID はIDで書籍を取得する
func (r *InMemoryBookRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Book, error) {
	book, ok := r.books[id]
	if !ok {
		return nil, ErrBookNotFound
	}
	return book, nil
}

// GetByUserID はユーザーIDで書籍一覧を取得する
func (r *InMemoryBookRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Book, error) {
	var books []*models.Book
	for _, book := range r.books {
		if book.UserID == userID {
			books = append(books, book)
		}
	}
	return books, nil
}

// Update は書籍情報を更新する
func (r *InMemoryBookRepository) Update(ctx context.Context, book *models.Book) error {
	if _, ok := r.books[book.ID]; !ok {
		return ErrBookNotFound
	}
	r.books[book.ID] = book
	return nil
}

// Delete は書籍を削除する
func (r *InMemoryBookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	delete(r.books, id)
	return nil
}

// UpdateStatus は書籍のステータスを更新する
func (r *InMemoryBookRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.BookStatus) error {
	book, ok := r.books[id]
	if !ok {
		return ErrBookNotFound
	}
	book.Status = status
	return nil
}

// ErrBookNotFound は書籍が見つからないエラー
var ErrBookNotFound = &RepositoryError{
	Code:    "BOOK_NOT_FOUND",
	Message: "book not found",
}

// RepositoryError はリポジトリエラー
type RepositoryError struct {
	Code    string
	Message string
}

func (e *RepositoryError) Error() string {
	return e.Message
}

// PostgreSQL実装

// bookRepositoryPostgres はPostgreSQLベースの書籍リポジトリ実装
type bookRepositoryPostgres struct {
	db *sql.DB
}

// NewBookRepositoryPostgres は新しいPostgreSQL実装のBookRepositoryを作成する
func NewBookRepositoryPostgres(db *sql.DB) BookRepository {
	return &bookRepositoryPostgres{db: db}
}

// Create は書籍を作成する
func (r *bookRepositoryPostgres) Create(ctx context.Context, book *models.Book) error {
	query := `
		INSERT INTO books (id, user_id, title, target_language, native_language, reference_language,
		                  cover_image_url, total_pages, processed_pages, status, ocr_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		book.ID,
		book.UserID,
		book.Title,
		book.TargetLanguage,
		book.NativeLanguage,
		book.ReferenceLanguage,
		book.CoverImageURL,
		book.TotalPages,
		book.ProcessedPages,
		book.Status,
		book.OCRStatus,
		book.CreatedAt,
		book.UpdatedAt,
	)

	return err
}

// GetByID はIDで書籍を取得する
func (r *bookRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (*models.Book, error) {
	query := `
		SELECT id, user_id, title, target_language, native_language, reference_language,
		       cover_image_url, total_pages, processed_pages, status, ocr_status, created_at, updated_at
		FROM books
		WHERE id = $1
	`

	book := &models.Book{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.UserID,
		&book.Title,
		&book.TargetLanguage,
		&book.NativeLanguage,
		&book.ReferenceLanguage,
		&book.CoverImageURL,
		&book.TotalPages,
		&book.ProcessedPages,
		&book.Status,
		&book.OCRStatus,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrBookNotFound
		}
		return nil, err
	}

	return book, nil
}

// GetByUserID はユーザーIDで書籍一覧を取得する
func (r *bookRepositoryPostgres) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Book, error) {
	query := `
		SELECT id, user_id, title, target_language, native_language, reference_language,
		       cover_image_url, total_pages, processed_pages, status, ocr_status, created_at, updated_at
		FROM books
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		book := &models.Book{}
		err := rows.Scan(
			&book.ID,
			&book.UserID,
			&book.Title,
			&book.TargetLanguage,
			&book.NativeLanguage,
			&book.ReferenceLanguage,
			&book.CoverImageURL,
			&book.TotalPages,
			&book.ProcessedPages,
			&book.Status,
			&book.OCRStatus,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

// Update は書籍情報を更新する
func (r *bookRepositoryPostgres) Update(ctx context.Context, book *models.Book) error {
	query := `
		UPDATE books
		SET title = $1, target_language = $2, native_language = $3, reference_language = $4,
		    cover_image_url = $5, total_pages = $6, processed_pages = $7, status = $8, ocr_status = $9,
		    updated_at = NOW()
		WHERE id = $10
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		book.Title,
		book.TargetLanguage,
		book.NativeLanguage,
		book.ReferenceLanguage,
		book.CoverImageURL,
		book.TotalPages,
		book.ProcessedPages,
		book.Status,
		book.OCRStatus,
		book.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrBookNotFound
	}

	return nil
}

// Delete は書籍を削除する
func (r *bookRepositoryPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM books WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrBookNotFound
	}

	return nil
}

// UpdateStatus は書籍のステータスを更新する
func (r *bookRepositoryPostgres) UpdateStatus(ctx context.Context, id uuid.UUID, status models.BookStatus) error {
	query := `
		UPDATE books
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrBookNotFound
	}

	return nil
}
