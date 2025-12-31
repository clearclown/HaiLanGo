package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestIDHeader = "X-Request-ID"
	requestIDKey    = "request_id"
)

// RequestID は各リクエストにRequest IDを付与するミドルウェア
// - 既に X-Request-ID が付与されていれば安全性チェック後に採用
// - 不正/未指定の場合はサーバ側で生成
// - 応答にも X-Request-ID を返す
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := strings.TrimSpace(c.GetHeader(requestIDHeader))
		if !isSafeRequestID(requestID) {
			requestID = uuid.NewString()
		}

		c.Set(requestIDKey, requestID)
		c.Writer.Header().Set(requestIDHeader, requestID)

		c.Next()
	}
}

// GetRequestID はContextからRequest IDを取得する
func GetRequestID(c *gin.Context) string {
	if c == nil {
		return ""
	}
	v, ok := c.Get(requestIDKey)
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

func isSafeRequestID(s string) bool {
	// 安全なログ出力/ヘッダー返却のため、文字種と長さを制限する
	// 許可: 英数字 + - _ . : のみ（一般的なUUID/トレースIDをカバー）
	if s == "" || len(s) > 128 {
		return false
	}
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if (ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '_' || ch == '.' || ch == ':' {
			continue
		}
		return false
	}
	return true
}
