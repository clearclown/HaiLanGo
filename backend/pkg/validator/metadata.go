package validator

import (
	"errors"
	"regexp"

	"github.com/clearclown/HaiLanGo/internal/models"
)

var (
	// ErrInvalidLanguageCode は無効な言語コードのエラー
	ErrInvalidLanguageCode = errors.New("invalid language code")
	// ErrEmptyTitle はタイトルが空のエラー
	ErrEmptyTitle = errors.New("title cannot be empty")
	// ErrTitleTooLong はタイトルが長すぎるエラー
	ErrTitleTooLong = errors.New("title is too long (max 200 characters)")
)

// ISO 639-1 言語コード（主要な言語）
var validLanguageCodes = map[string]bool{
	"ja": true, // 日本語
	"en": true, // 英語
	"zh": true, // 中国語
	"ru": true, // ロシア語
	"fa": true, // ペルシャ語
	"he": true, // ヘブライ語
	"es": true, // スペイン語
	"fr": true, // フランス語
	"pt": true, // ポルトガル語
	"de": true, // ドイツ語
	"it": true, // イタリア語
	"tr": true, // トルコ語
	"ar": true, // アラビア語
	"ko": true, // 韓国語
	"vi": true, // ベトナム語
	"th": true, // タイ語
	"hi": true, // ヒンディー語
	"bn": true, // ベンガル語
	"pl": true, // ポーランド語
	"uk": true, // ウクライナ語
	"nl": true, // オランダ語
	"sv": true, // スウェーデン語
	"no": true, // ノルウェー語
	"da": true, // デンマーク語
	"fi": true, // フィンランド語
	"el": true, // ギリシャ語
	"cs": true, // チェコ語
	"hu": true, // ハンガリー語
	"ro": true, // ルーマニア語
	"id": true, // インドネシア語
	"ms": true, // マレー語
	"tl": true, // タガログ語
}

const (
	maxTitleLength = 200
)

// ValidateBookMetadata は書籍メタデータを検証する
func ValidateBookMetadata(metadata models.BookMetadata) error {
	// タイトルの検証
	if err := ValidateTitle(metadata.Title); err != nil {
		return err
	}

	// 学習先言語の検証
	if err := ValidateLanguageCode(metadata.TargetLanguage); err != nil {
		return err
	}

	// 母国語の検証
	if err := ValidateLanguageCode(metadata.NativeLanguage); err != nil {
		return err
	}

	// 参照言語の検証（オプション）
	if metadata.ReferenceLanguage != "" {
		if err := ValidateLanguageCode(metadata.ReferenceLanguage); err != nil {
			return err
		}
	}

	return nil
}

// ValidateTitle はタイトルを検証する
func ValidateTitle(title string) error {
	// 空チェック
	if title == "" {
		return ErrEmptyTitle
	}

	// 長さチェック
	if len([]rune(title)) > maxTitleLength {
		return ErrTitleTooLong
	}

	return nil
}

// ValidateLanguageCode は言語コードを検証する（ISO 639-1）
func ValidateLanguageCode(code string) error {
	// 形式チェック（2文字の小文字）
	matched, _ := regexp.MatchString(`^[a-z]{2}$`, code)
	if !matched {
		return ErrInvalidLanguageCode
	}

	// サポートされている言語コードかチェック
	if !validLanguageCodes[code] {
		return ErrInvalidLanguageCode
	}

	return nil
}

// IsValidLanguageCode は言語コードが有効かどうかを返す
func IsValidLanguageCode(code string) bool {
	return validLanguageCodes[code]
}

// GetSupportedLanguages はサポートされている言語コードの一覧を返す
func GetSupportedLanguages() []string {
	codes := make([]string, 0, len(validLanguageCodes))
	for code := range validLanguageCodes {
		codes = append(codes, code)
	}
	return codes
}
