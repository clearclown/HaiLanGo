package middleware

import (
	"os"
	"strings"
)

// IsOriginAllowed は CORS_ALLOWED_ORIGINS に基づいてOriginを許可するか判定する
// - origin が空の場合は true（ブラウザ以外のクライアント等）
// - CORS_ALLOWED_ORIGINS が未設定 or "*" の場合は true（開発向け）
func IsOriginAllowed(origin string) bool {
	if origin == "" {
		return true
	}

	raw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if raw == "" || raw == "*" {
		return true
	}

	for _, part := range strings.Split(raw, ",") {
		if strings.TrimSpace(part) == origin {
			return true
		}
	}

	return false
}
