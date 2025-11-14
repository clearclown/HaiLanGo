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
		// Authorizationヘッダーを取得
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Bearer トークンの抽出
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

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
