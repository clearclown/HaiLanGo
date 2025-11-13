package password

import (
	"errors"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	// MinPasswordLength はパスワードの最小長
	MinPasswordLength = 8
	// BcryptCost はbcryptのコスト（10は標準的な値）
	BcryptCost = 10
)

var (
	// ErrEmptyPassword は空のパスワードエラー
	ErrEmptyPassword = errors.New("パスワードが空です")
	// ErrPasswordTooShort はパスワードが短すぎるエラー
	ErrPasswordTooShort = errors.New("パスワードは最低8文字必要です")
	// ErrPasswordTooWeak はパスワードが弱すぎるエラー
	ErrPasswordTooWeak = errors.New("パスワードは大文字、小文字、数字、記号のうち3種類以上を含む必要があります")
)

// HashPassword はパスワードをbcryptでハッシュ化する
func HashPassword(password string) (string, error) {
	// 空のパスワードチェック
	if password == "" {
		return "", ErrEmptyPassword
	}

	// パスワードをbcryptでハッシュ化
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// VerifyPassword はパスワードとハッシュを比較検証する
func VerifyPassword(password, hash string) bool {
	// 空のパスワードまたはハッシュの場合は失敗
	if password == "" || hash == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePasswordStrength はパスワードの強度を検証する
func ValidatePasswordStrength(password string) error {
	// 長さチェック
	if len(password) < MinPasswordLength {
		return ErrPasswordTooShort
	}

	// パスワード強度チェック：大文字、小文字、数字、記号のうち3種類以上を含む
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// 4種類のうち3種類以上が含まれているかチェック
	typesCount := 0
	if hasUpper {
		typesCount++
	}
	if hasLower {
		typesCount++
	}
	if hasNumber {
		typesCount++
	}
	if hasSpecial {
		typesCount++
	}

	if typesCount < 3 {
		return ErrPasswordTooWeak
	}

	return nil
}

// IsValidEmail はメールアドレスの形式を検証する
func IsValidEmail(email string) bool {
	// シンプルなメールアドレスの正規表現
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
