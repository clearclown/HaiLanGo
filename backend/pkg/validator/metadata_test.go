package validator

import (
	"testing"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
		errType error
	}{
		{
			name:    "正常なタイトル",
			title:   "ロシア語入門",
			wantErr: false,
		},
		{
			name:    "長いタイトル（200文字以内）",
			title:   string(make([]rune, 200)),
			wantErr: false,
		},
		{
			name:    "空のタイトル",
			title:   "",
			wantErr: true,
			errType: ErrEmptyTitle,
		},
		{
			name:    "タイトルが長すぎる",
			title:   string(make([]rune, 201)),
			wantErr: true,
			errType: ErrTitleTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTitle(tt.title)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateLanguageCode(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{name: "日本語", code: "ja", wantErr: false},
		{name: "英語", code: "en", wantErr: false},
		{name: "ロシア語", code: "ru", wantErr: false},
		{name: "ペルシャ語", code: "fa", wantErr: false},
		{name: "ヘブライ語", code: "he", wantErr: false},
		{name: "中国語", code: "zh", wantErr: false},
		{name: "スペイン語", code: "es", wantErr: false},
		{name: "大文字（無効）", code: "JA", wantErr: true},
		{name: "3文字（無効）", code: "jpn", wantErr: true},
		{name: "1文字（無効）", code: "j", wantErr: true},
		{name: "サポートされていない言語", code: "xx", wantErr: true},
		{name: "空文字", code: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLanguageCode(tt.code)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidLanguageCode, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateBookMetadata(t *testing.T) {
	tests := []struct {
		name     string
		metadata models.BookMetadata
		wantErr  bool
		errType  error
	}{
		{
			name: "正常なメタデータ",
			metadata: models.BookMetadata{
				Title:          "ロシア語入門",
				TargetLanguage: "ru",
				NativeLanguage: "ja",
			},
			wantErr: false,
		},
		{
			name: "参照言語も含む",
			metadata: models.BookMetadata{
				Title:             "クルド語入門",
				TargetLanguage:    "tr",
				NativeLanguage:    "ja",
				ReferenceLanguage: "en",
			},
			wantErr: false,
		},
		{
			name: "タイトルが空",
			metadata: models.BookMetadata{
				Title:          "",
				TargetLanguage: "ru",
				NativeLanguage: "ja",
			},
			wantErr: true,
			errType: ErrEmptyTitle,
		},
		{
			name: "学習先言語が無効",
			metadata: models.BookMetadata{
				Title:          "テスト",
				TargetLanguage: "xx",
				NativeLanguage: "ja",
			},
			wantErr: true,
			errType: ErrInvalidLanguageCode,
		},
		{
			name: "母国語が無効",
			metadata: models.BookMetadata{
				Title:          "テスト",
				TargetLanguage: "en",
				NativeLanguage: "invalid",
			},
			wantErr: true,
			errType: ErrInvalidLanguageCode,
		},
		{
			name: "参照言語が無効",
			metadata: models.BookMetadata{
				Title:             "テスト",
				TargetLanguage:    "en",
				NativeLanguage:    "ja",
				ReferenceLanguage: "xxx",
			},
			wantErr: true,
			errType: ErrInvalidLanguageCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBookMetadata(tt.metadata)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsValidLanguageCode(t *testing.T) {
	assert.True(t, IsValidLanguageCode("ja"))
	assert.True(t, IsValidLanguageCode("en"))
	assert.False(t, IsValidLanguageCode("xx"))
	assert.False(t, IsValidLanguageCode("JA"))
}

func TestGetSupportedLanguages(t *testing.T) {
	languages := GetSupportedLanguages()
	assert.NotEmpty(t, languages)
	assert.Contains(t, languages, "ja")
	assert.Contains(t, languages, "en")
	assert.Contains(t, languages, "ru")
}
