package tts

import (
	"context"
	"testing"
)

// BenchmarkGenerateAudio は音声生成のベンチマーク
func BenchmarkGenerateAudio(b *testing.B) {
	service := NewTTSService()
	ctx := context.Background()
	text := "Hello, world! This is a benchmark test."
	lang := "en"
	quality := "standard"
	speed := 1.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GenerateAudio(ctx, text, lang, quality, speed)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGenerateAudioCached はキャッシュされた音声生成のベンチマーク
func BenchmarkGenerateAudioCached(b *testing.B) {
	service := NewTTSService()
	ctx := context.Background()
	text := "Hello, world! This is a benchmark test."
	lang := "en"
	quality := "standard"
	speed := 1.0

	// 事前にキャッシュに保存
	_, err := service.GenerateAudio(ctx, text, lang, quality, speed)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GenerateAudio(ctx, text, lang, quality, speed)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGenerateAudioShortText は短いテキストの音声生成ベンチマーク
func BenchmarkGenerateAudioShortText(b *testing.B) {
	service := NewTTSService()
	ctx := context.Background()
	texts := []string{"Hello", "Hi", "Good morning", "Thank you", "Goodbye"}
	lang := "en"
	quality := "standard"
	speed := 1.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		text := texts[i%len(texts)]
		_, err := service.GenerateAudio(ctx, text, lang, quality, speed)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGenerateAudioLongText は長いテキストの音声生成ベンチマーク
func BenchmarkGenerateAudioLongText(b *testing.B) {
	service := NewTTSService()
	ctx := context.Background()
	// 約300文字のテキスト
	text := ""
	for i := 0; i < 30; i++ {
		text += "This is a long text for benchmarking purposes. "
	}
	lang := "en"
	quality := "standard"
	speed := 1.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GenerateAudio(ctx, text, lang, quality, speed)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBatchGenerate はバッチ生成のベンチマーク
func BenchmarkBatchGenerate(b *testing.B) {
	service := NewTTSService()
	ctx := context.Background()
	texts := []string{"Hello", "Goodbye", "Thank you", "Good morning", "Good night"}
	lang := "en"
	quality := "standard"
	speed := 1.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.BatchGenerate(ctx, texts, lang, quality, speed)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMultipleLanguages は複数言語の音声生成ベンチマーク
func BenchmarkMultipleLanguages(b *testing.B) {
	service := NewTTSService()
	ctx := context.Background()
	text := "Hello, world!"
	languages := []string{"en", "ja", "zh", "ru", "es", "fr"}
	quality := "standard"
	speed := 1.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lang := languages[i%len(languages)]
		_, err := service.GenerateAudio(ctx, text, lang, quality, speed)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDifferentSpeeds は異なる速度の音声生成ベンチマーク
func BenchmarkDifferentSpeeds(b *testing.B) {
	service := NewTTSService()
	ctx := context.Background()
	text := "Hello, world!"
	lang := "en"
	quality := "standard"
	speeds := []float64{0.5, 0.75, 1.0, 1.25, 1.5, 2.0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		speed := speeds[i%len(speeds)]
		_, err := service.GenerateAudio(ctx, text, lang, quality, speed)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConcurrentGeneration は並列音声生成のベンチマーク
func BenchmarkConcurrentGeneration(b *testing.B) {
	service := NewTTSService()
	ctx := context.Background()
	text := "Hello, world!"
	lang := "en"
	quality := "standard"
	speed := 1.0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := service.GenerateAudio(ctx, text, lang, quality, speed)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
