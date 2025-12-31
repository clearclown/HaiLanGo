package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger はRequest IDと合わせてリクエストの概要をログ出力する
// - クエリ文字列（token等が混入しやすい）は出力しない
// - Authorization等のヘッダーも出力しない
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		method := sanitizeLogValue(c.Request.Method)
		path := sanitizeLogValue(c.Request.URL.Path)
		clientIP := sanitizeLogValue(c.ClientIP())

		c.Next()

		status := c.Writer.Status()
		latency := time.Since(start)

		requestID := sanitizeLogValue(GetRequestID(c))

		userID := ""
		if v, exists := c.Get("user_id"); exists {
			if s, ok := v.(string); ok {
				userID = sanitizeLogValue(s)
			}
		}

		errCount := len(c.Errors)
		errMsg := ""
		if errCount > 0 {
			// gin.Errors は内部用メッセージが含まれる可能性があるため、最小限に抑える
			errMsg = sanitizeLogValue(c.Errors.String())
		}

		// ログレベル風のプレフィックス（標準logのみで運用できるようにする）
		prefix := "INFO"
		if status >= 500 {
			prefix = "ERROR"
		} else if status >= 400 {
			prefix = "WARN"
		}

		if userID != "" {
			log.Printf("%s request_id=%s status=%d method=%s path=%s latency=%s ip=%s user_id=%s errors=%d err=%s",
				prefix, requestID, status, method, path, latency, clientIP, userID, errCount, errMsg)
			return
		}

		log.Printf("%s request_id=%s status=%d method=%s path=%s latency=%s ip=%s errors=%d err=%s",
			prefix, requestID, status, method, path, latency, clientIP, errCount, errMsg)
	}
}

func sanitizeLogValue(s string) string {
	// ログ改ざん（改行注入）を避ける
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}
