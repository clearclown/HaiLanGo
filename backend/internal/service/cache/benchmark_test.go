package cache

import (
	"testing"
	"time"
)

// BenchmarkCacheSet はキャッシュ保存のベンチマーク
func BenchmarkCacheSet(b *testing.B) {
	cache := NewAudioCache()
	audioURL := "https://example.com/audio.mp3"
	ttl := 7 * 24 * time.Hour

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := cache.GenerateKey("text", "en", "standard", float64(i))
		err := cache.Set(key, audioURL, ttl)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCacheGet はキャッシュ取得のベンチマーク
func BenchmarkCacheGet(b *testing.B) {
	cache := NewAudioCache()
	audioURL := "https://example.com/audio.mp3"
	ttl := 7 * 24 * time.Hour

	// 事前に1000個のアイテムをキャッシュに保存
	for i := 0; i < 1000; i++ {
		key := cache.GenerateKey("text", "en", "standard", float64(i))
		cache.Set(key, audioURL, ttl)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := cache.GenerateKey("text", "en", "standard", float64(i%1000))
		_, found := cache.Get(key)
		if !found {
			b.Fatal("cache miss")
		}
	}
}

// BenchmarkCacheGenerateKey はキャッシュキー生成のベンチマーク
func BenchmarkCacheGenerateKey(b *testing.B) {
	cache := NewAudioCache()
	text := "Hello, world! This is a test."
	lang := "en"
	quality := "standard"
	speed := 1.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cache.GenerateKey(text, lang, quality, speed)
	}
}

// BenchmarkCacheHitRate はキャッシュヒット率のベンチマーク
func BenchmarkCacheHitRate(b *testing.B) {
	cache := NewAudioCache()
	audioURL := "https://example.com/audio.mp3"
	ttl := 7 * 24 * time.Hour

	// 事前に100個のアイテムをキャッシュに保存
	for i := 0; i < 100; i++ {
		key := cache.GenerateKey("text", "en", "standard", float64(i))
		cache.Set(key, audioURL, ttl)
	}

	hits := 0
	misses := 0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 70%の確率でキャッシュヒット、30%の確率でキャッシュミス
		var key string
		if i%10 < 7 {
			// キャッシュヒット
			key = cache.GenerateKey("text", "en", "standard", float64(i%100))
		} else {
			// キャッシュミス
			key = cache.GenerateKey("text", "en", "standard", float64(i+1000))
		}

		_, found := cache.Get(key)
		if found {
			hits++
		} else {
			misses++
		}
	}

	b.ReportMetric(float64(hits)/float64(hits+misses)*100, "hit_rate_%")
}

// BenchmarkConcurrentCacheAccess は並列キャッシュアクセスのベンチマーク
func BenchmarkConcurrentCacheAccess(b *testing.B) {
	cache := NewAudioCache()
	audioURL := "https://example.com/audio.mp3"
	ttl := 7 * 24 * time.Hour

	// 事前に100個のアイテムをキャッシュに保存
	for i := 0; i < 100; i++ {
		key := cache.GenerateKey("text", "en", "standard", float64(i))
		cache.Set(key, audioURL, ttl)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := cache.GenerateKey("text", "en", "standard", float64(i%100))
			_, _ = cache.Get(key)
			i++
		}
	})
}
