package repository

import (
	"context"

	"github.com/clearclown/HaiLanGo/internal/models"
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
