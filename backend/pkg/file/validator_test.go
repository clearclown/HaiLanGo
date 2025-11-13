package file

import (
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateFileType はファイルタイプの検証をテストする
func TestValidateFileType(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		fileName    string
		want        bool
	}{
		{
			name:        "PDFファイル - 正常",
			contentType: "application/pdf",
			fileName:    "test.pdf",
			want:        true,
		},
		{
			name:        "PNGファイル - 正常",
			contentType: "image/png",
			fileName:    "test.png",
			want:        true,
		},
		{
			name:        "JPGファイル - 正常",
			contentType: "image/jpeg",
			fileName:    "test.jpg",
			want:        true,
		},
		{
			name:        "JPEGファイル - 正常",
			contentType: "image/jpeg",
			fileName:    "test.jpeg",
			want:        true,
		},
		{
			name:        "HEICファイル - 正常",
			contentType: "image/heic",
			fileName:    "test.heic",
			want:        true,
		},
		{
			name:        "不正なファイル - ZIPファイル",
			contentType: "application/zip",
			fileName:    "test.zip",
			want:        false,
		},
		{
			name:        "不正なファイル - EXEファイル",
			contentType: "application/x-msdownload",
			fileName:    "test.exe",
			want:        false,
		},
		{
			name:        "不正なファイル - テキストファイル",
			contentType: "text/plain",
			fileName:    "test.txt",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateFileType(tt.contentType, tt.fileName)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestValidateFileSize はファイルサイズの検証をテストする
func TestValidateFileSize(t *testing.T) {
	tests := []struct {
		name string
		size int64
		want bool
	}{
		{
			name: "正常サイズ - 1MB",
			size: 1 * 1024 * 1024,
			want: true,
		},
		{
			name: "正常サイズ - 10MB",
			size: 10 * 1024 * 1024,
			want: true,
		},
		{
			name: "正常サイズ - 50MB",
			size: 50 * 1024 * 1024,
			want: true,
		},
		{
			name: "正常サイズ - 100MB（上限）",
			size: 100 * 1024 * 1024,
			want: true,
		},
		{
			name: "サイズ超過 - 101MB",
			size: 101 * 1024 * 1024,
			want: false,
		},
		{
			name: "サイズ超過 - 150MB",
			size: 150 * 1024 * 1024,
			want: false,
		},
		{
			name: "サイズ0 - 無効",
			size: 0,
			want: false,
		},
		{
			name: "負のサイズ - 無効",
			size: -1,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateFileSize(tt.size)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestValidateFile はファイル全体の検証をテストする
func TestValidateFile(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		fileName    string
		fileSize    int64
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "正常なPDFファイル",
			contentType: "application/pdf",
			fileName:    "book.pdf",
			fileSize:    10 * 1024 * 1024,
			wantErr:     false,
		},
		{
			name:        "正常なPNGファイル",
			contentType: "image/png",
			fileName:    "page1.png",
			fileSize:    5 * 1024 * 1024,
			wantErr:     false,
		},
		{
			name:        "ファイルタイプが不正",
			contentType: "application/zip",
			fileName:    "archive.zip",
			fileSize:    10 * 1024 * 1024,
			wantErr:     true,
			errMsg:      "unsupported file type",
		},
		{
			name:        "ファイルサイズが大きすぎる",
			contentType: "application/pdf",
			fileName:    "large.pdf",
			fileSize:    150 * 1024 * 1024,
			wantErr:     true,
			errMsg:      "file size exceeds limit",
		},
		{
			name:        "ファイルサイズが0",
			contentType: "application/pdf",
			fileName:    "empty.pdf",
			fileSize:    0,
			wantErr:     true,
			errMsg:      "file size is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの multipart.FileHeader を作成
			fileHeader := &multipart.FileHeader{
				Filename: tt.fileName,
				Size:     tt.fileSize,
				Header:   make(map[string][]string),
			}
			fileHeader.Header.Set("Content-Type", tt.contentType)

			err := ValidateFile(fileHeader)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGetFileExtension はファイル拡張子の取得をテストする
func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     string
	}{
		{
			name:     "PDF拡張子",
			fileName: "book.pdf",
			want:     ".pdf",
		},
		{
			name:     "PNG拡張子",
			fileName: "image.png",
			want:     ".png",
		},
		{
			name:     "複数のドット",
			fileName: "my.book.name.pdf",
			want:     ".pdf",
		},
		{
			name:     "拡張子なし",
			fileName: "filename",
			want:     "",
		},
		{
			name:     "大文字拡張子",
			fileName: "IMAGE.PNG",
			want:     ".png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFileExtension(tt.fileName)
			assert.Equal(t, tt.want, got)
		})
	}
}
