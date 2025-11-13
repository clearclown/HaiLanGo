package storage

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLocalStorage_SaveFile はローカルストレージへのファイル保存をテストする
func TestLocalStorage_SaveFile(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	storage := NewLocalStorage(tempDir)
	ctx := context.Background()

	tests := []struct {
		name     string
		userID   uuid.UUID
		bookID   uuid.UUID
		fileName string
		data     []byte
		wantErr  bool
	}{
		{
			name:     "正常にPDFファイルを保存",
			userID:   uuid.New(),
			bookID:   uuid.New(),
			fileName: "book.pdf",
			data:     []byte("PDF content"),
			wantErr:  false,
		},
		{
			name:     "正常にPNGファイルを保存",
			userID:   uuid.New(),
			bookID:   uuid.New(),
			fileName: "page1.png",
			data:     []byte("PNG content"),
			wantErr:  false,
		},
		{
			name:     "空のデータを保存",
			userID:   uuid.New(),
			bookID:   uuid.New(),
			fileName: "empty.pdf",
			data:     []byte{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := storage.SaveFile(ctx, tt.userID, tt.bookID, tt.fileName, bytes.NewReader(tt.data))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, path)

			// ファイルが実際に保存されているか確認
			fullPath := filepath.Join(tempDir, path)
			savedData, err := os.ReadFile(fullPath)
			require.NoError(t, err)
			assert.Equal(t, tt.data, savedData)
		})
	}
}

// TestLocalStorage_GetFile はローカルストレージからのファイル取得をテストする
func TestLocalStorage_GetFile(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewLocalStorage(tempDir)
	ctx := context.Background()

	// テストデータを保存
	userID := uuid.New()
	bookID := uuid.New()
	fileName := "test.pdf"
	testData := []byte("test content")

	path, err := storage.SaveFile(ctx, userID, bookID, fileName, bytes.NewReader(testData))
	require.NoError(t, err)

	// ファイルを取得
	reader, err := storage.GetFile(ctx, path)
	require.NoError(t, err)
	defer reader.Close()

	// 内容を確認
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)
	require.NoError(t, err)
	assert.Equal(t, testData, buf.Bytes())
}

// TestLocalStorage_DeleteFile はローカルストレージからのファイル削除をテストする
func TestLocalStorage_DeleteFile(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewLocalStorage(tempDir)
	ctx := context.Background()

	// テストデータを保存
	userID := uuid.New()
	bookID := uuid.New()
	fileName := "delete_test.pdf"
	testData := []byte("delete me")

	path, err := storage.SaveFile(ctx, userID, bookID, fileName, bytes.NewReader(testData))
	require.NoError(t, err)

	// ファイルが存在することを確認
	fullPath := filepath.Join(tempDir, path)
	_, err = os.Stat(fullPath)
	require.NoError(t, err)

	// ファイルを削除
	err = storage.DeleteFile(ctx, path)
	require.NoError(t, err)

	// ファイルが削除されたことを確認
	_, err = os.Stat(fullPath)
	assert.True(t, os.IsNotExist(err))
}

// TestLocalStorage_FileExists はファイルの存在チェックをテストする
func TestLocalStorage_FileExists(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewLocalStorage(tempDir)
	ctx := context.Background()

	// テストデータを保存
	userID := uuid.New()
	bookID := uuid.New()
	fileName := "exists_test.pdf"
	testData := []byte("exists")

	path, err := storage.SaveFile(ctx, userID, bookID, fileName, bytes.NewReader(testData))
	require.NoError(t, err)

	// ファイルが存在することを確認
	exists, err := storage.FileExists(ctx, path)
	require.NoError(t, err)
	assert.True(t, exists)

	// 存在しないファイルをチェック
	exists, err = storage.FileExists(ctx, "non/existent/file.pdf")
	require.NoError(t, err)
	assert.False(t, exists)
}

// TestGenerateStoragePath はストレージパスの生成をテストする
func TestGenerateStoragePath(t *testing.T) {
	userID := uuid.New()
	bookID := uuid.New()

	tests := []struct {
		name     string
		fileName string
	}{
		{
			name:     "PDFファイル",
			fileName: "book.pdf",
		},
		{
			name:     "PNGファイル",
			fileName: "page1.png",
		},
		{
			name:     "スペースを含むファイル名",
			fileName: "my book name.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := GenerateStoragePath(userID, bookID, tt.fileName)

			// パスの形式を確認
			assert.NotEmpty(t, path)
			assert.Contains(t, path, userID.String())
			assert.Contains(t, path, bookID.String())

			// 拡張子が保持されているか確認
			ext := filepath.Ext(tt.fileName)
			assert.Contains(t, path, ext)
		})
	}
}

// TestSanitizeFileName はファイル名のサニタイズをテストする
func TestSanitizeFileName(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     string
	}{
		{
			name:     "正常なファイル名",
			fileName: "book.pdf",
			want:     "book.pdf",
		},
		{
			name:     "スペースをアンダースコアに変換",
			fileName: "my book.pdf",
			want:     "my_book.pdf",
		},
		{
			name:     "特殊文字を削除",
			fileName: "book#123!.pdf",
			want:     "book123.pdf",
		},
		{
			name:     "日本語ファイル名",
			fileName: "本.pdf",
			want:     "本.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeFileName(tt.fileName)
			assert.Equal(t, tt.want, got)
		})
	}
}
