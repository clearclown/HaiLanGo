package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders は最低限のセキュリティヘッダーを付与する
// API用途のため、過度に強いポリシー（CSP等）はここでは付与しない
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()

		// MIME sniffing を防止
		h.Set("X-Content-Type-Options", "nosniff")

		// クリックジャッキング対策（API用途では基本的に問題にならないが安全側）
		h.Set("X-Frame-Options", "DENY")

		// Referer からの情報漏洩を抑制
		h.Set("Referrer-Policy", "no-referrer")

		c.Next()
	}
}
