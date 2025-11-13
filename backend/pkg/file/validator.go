package file

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"
)

const (
	// MaxFileSize はアップロード可能な最大ファイルサイズ（100MB）
	MaxFileSize = 100 * 1024 * 1024 // 100MB in bytes
)

var (
	// ErrUnsupportedFileType はサポートされていないファイルタイプのエラー
	ErrUnsupportedFileType = errors.New("unsupported file type")
	// ErrFileSizeExceedsLimit はファイルサイズが制限を超えているエラー
	ErrFileSizeExceedsLimit = errors.New("file size exceeds limit")
	// ErrInvalidFileSize はファイルサイズが無効なエラー
	ErrInvalidFileSize = errors.New("file size is invalid")
)

// AllowedContentTypes はサポートされているMIMEタイプのリスト
var AllowedContentTypes = map[string]bool{
	"application/pdf": true,
	"image/png":       true,
	"image/jpeg":      true,
	"image/jpg":       true,
	"image/heic":      true,
}

// AllowedExtensions はサポートされているファイル拡張子のリスト
var AllowedExtensions = map[string]bool{
	".pdf":  true,
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".heic": true,
}

// ValidateFileType はファイルタイプが許可されているかを検証する
func ValidateFileType(contentType, fileName string) bool {
	// Content-Typeをチェック
	if AllowedContentTypes[contentType] {
		return true
	}

	// 拡張子をチェック（フォールバック）
	ext := GetFileExtension(fileName)
	return AllowedExtensions[ext]
}

// ValidateFileSize はファイルサイズが制限内かを検証する
func ValidateFileSize(size int64) bool {
	return size > 0 && size <= MaxFileSize
}

// ValidateFile はファイルヘッダー全体を検証する
func ValidateFile(fileHeader *multipart.FileHeader) error {
	// ファイルサイズの検証
	if fileHeader.Size <= 0 {
		return ErrInvalidFileSize
	}

	if fileHeader.Size > MaxFileSize {
		return ErrFileSizeExceedsLimit
	}

	// Content-Typeを取得
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		// Content-Typeがない場合は拡張子から判断
		ext := GetFileExtension(fileHeader.Filename)
		contentType = getContentTypeFromExtension(ext)
	}

	// ファイルタイプの検証
	if !ValidateFileType(contentType, fileHeader.Filename) {
		return ErrUnsupportedFileType
	}

	return nil
}

// GetFileExtension はファイル名から拡張子を取得する（小文字化）
func GetFileExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	return strings.ToLower(ext)
}

// getContentTypeFromExtension は拡張子からContent-Typeを取得する
func getContentTypeFromExtension(ext string) string {
	contentTypeMap := map[string]string{
		".pdf":  "application/pdf",
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".heic": "image/heic",
	}
	return contentTypeMap[ext]
}
