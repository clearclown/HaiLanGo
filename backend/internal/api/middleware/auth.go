package middleware

import (
	"net/http"
	"strings"

	"github.com/clearclown/HaiLanGo/backend/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// AuthRequired は認証必須のミドルウェア
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// まずAuthorizationヘッダーを確認
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// Bearer トークンの抽出
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
				c.Abort()
				return
			}
			tokenString = parts[1]
		} else {
			// Authorizationヘッダーがない場合、クエリパラメータを確認（WebSocket用）
			tokenString = c.Query("token")
			if tokenString == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header or token query parameter required"})
				c.Abort()
				return
			}
		}

		// トークンの検証
		claims, err := jwt.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// ユーザーIDをコンテキストに設定
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
