package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter はレート制限ミドルウェア
// 簡易的なインメモリ実装（本番環境ではRedisを使用すること）
func RateLimiter() gin.HandlerFunc {
	type client struct {
		count      int
		lastAccess time.Time
	}

	var (
		clients = make(map[string]*client)
		mu      sync.Mutex
	)

	// 設定
	maxRequests := 100 // 1分あたりの最大リクエスト数
	window := time.Minute

	// クリーンアップゴルーチン
	go func() {
		for {
			time.Sleep(window)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastAccess) > window {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		defer mu.Unlock()

		if _, exists := clients[ip]; !exists {
			clients[ip] = &client{
				count:      0,
				lastAccess: time.Now(),
			}
		}

		cli := clients[ip]

		// ウィンドウの更新
		if time.Since(cli.lastAccess) > window {
			cli.count = 0
			cli.lastAccess = time.Now()
		}

		// レート制限チェック
		if cli.count >= maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "リクエスト数が多すぎます。しばらく待ってから再度お試しください。",
			})
			c.Abort()
			return
		}

		cli.count++
		c.Next()
	}
}
