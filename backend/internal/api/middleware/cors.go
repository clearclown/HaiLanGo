package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS はCORSミドルウェア
func CORS() gin.HandlerFunc {
	cfg := newCORSConfigFromEnv()

	if cfg.allowAll {
		log.Println("⚠️  CORS_ALLOWED_ORIGINS is not set (or '*'): allowing all origins. Set CORS_ALLOWED_ORIGINS for production.")
	} else {
		log.Printf("✅ CORS_ALLOWED_ORIGINS configured (%d origins)", len(cfg.allowedOrigins))
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Preflight/actual 共通で返すヘッダー
		c.Writer.Header().Set("Access-Control-Allow-Methods", cfg.allowedMethods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", cfg.allowedHeaders)

		// Credentials を許可する場合（Allow-Origin が '*' の場合は仕様上不可）
		if cfg.allowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Origin がある場合のみ CORS 判定（サーバ間通信等は Origin が付かないことがある）
		if origin != "" {
			if cfg.allowAll {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			} else if cfg.isOriginAllowed(origin) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Add("Vary", "Origin")
			} else if c.Request.Method == http.MethodOptions {
				// preflight は明示的に拒否（デバッグ容易化 & 早期失敗）
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CORS origin not allowed"})
				return
			}
		} else if cfg.allowAll {
			// Origin がないケースは CORS と無関係だが、開発時の互換性のため '*' を返す
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

type corsConfig struct {
	allowAll         bool
	allowedOrigins   map[string]struct{}
	allowCredentials bool
	allowedHeaders   string
	allowedMethods   string
}

func newCORSConfigFromEnv() corsConfig {
	raw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	allowAll := raw == "" || raw == "*"

	allowed := make(map[string]struct{})
	if !allowAll {
		for _, part := range strings.Split(raw, ",") {
			o := strings.TrimSpace(part)
			if o == "" {
				continue
			}
			allowed[o] = struct{}{}
		}
	}

	// 明示指定がなければ「許可オリジン固定時のみ credentials を許可」する
	allowCreds := false
	if v := strings.TrimSpace(os.Getenv("CORS_ALLOW_CREDENTIALS")); v != "" {
		allowCreds = strings.EqualFold(v, "true")
	} else {
		allowCreds = !allowAll
	}

	// "*" + credentials は仕様上無効なので強制的に無効化（ログで気づけるようにする）
	if allowAll && allowCreds {
		log.Println("⚠️  CORS_ALLOW_CREDENTIALS=true with CORS_ALLOWED_ORIGINS='*' is invalid; forcing credentials=false")
		allowCreds = false
	}

	return corsConfig{
		allowAll:         allowAll,
		allowedOrigins:   allowed,
		allowCredentials: allowCreds,
		allowedHeaders:   "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With",
		allowedMethods:   "POST, OPTIONS, GET, PUT, DELETE, PATCH",
	}
}

func (c corsConfig) isOriginAllowed(origin string) bool {
	if c.allowAll {
		return true
	}
	_, ok := c.allowedOrigins[origin]
	return ok
}
