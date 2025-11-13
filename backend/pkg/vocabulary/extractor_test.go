package vocabulary

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExtractWords はテキストから単語を抽出するテスト
func TestExtractWords(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		language string
		expected []string
	}{
		{
			name:     "ロシア語テキストの単語抽出",
			text:     "Здравствуйте! Как дела?",
			language: "ru",
			expected: []string{"здравствуйте", "дела"}, // 正規化（小文字）、ストップワード除去
		},
		{
			name:     "英語テキストの単語抽出",
			text:     "Hello, how are you?",
			language: "en",
			expected: []string{"hello", "how", "you"}, // 正規化（小文字）、ストップワード除去（"are"は除外）
		},
		{
			name:     "日本語テキストの単語抽出",
			text:     "こんにちは、元気ですか？",
			language: "ja",
			expected: []string{"こんにちは", "元気"}, // ストップワード除去（"ですか"は除外）
		},
		{
			name:     "中国語テキストの単語抽出",
			text:     "你好，你好吗？",
			language: "zh",
			expected: []string{"你好", "你好吗"}, // ストップワード除去後の結果
		},
		{
			name:     "空のテキスト",
			text:     "",
			language: "en",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words := ExtractWords(tt.text, tt.language)
			assert.ElementsMatch(t, tt.expected, words)
		})
	}
}

// TestNormalizeWord は単語の正規化テスト
func TestNormalizeWord(t *testing.T) {
	tests := []struct {
		name     string
		word     string
		language string
		expected string
	}{
		{
			name:     "英語の大文字を小文字に変換",
			word:     "HELLO",
			language: "en",
			expected: "hello",
		},
		{
			name:     "ロシア語の正規化",
			word:     "Здравствуйте",
			language: "ru",
			expected: "здравствуйте",
		},
		{
			name:     "句読点の除去",
			word:     "hello!",
			language: "en",
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalized := NormalizeWord(tt.word, tt.language)
			assert.Equal(t, tt.expected, normalized)
		})
	}
}

// TestRemoveStopWords はストップワードの除去テスト
func TestRemoveStopWords(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		language string
		expected []string
	}{
		{
			name:     "英語のストップワード除去",
			words:    []string{"the", "hello", "is", "world", "a"},
			language: "en",
			expected: []string{"hello", "world"},
		},
		{
			name:     "ロシア語のストップワード除去",
			words:    []string{"это", "привет", "и", "мир"},
			language: "ru",
			expected: []string{"привет", "мир"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := RemoveStopWords(tt.words, tt.language)
			assert.ElementsMatch(t, tt.expected, filtered)
		})
	}
}

// TestRemoveDuplicates は重複単語の除去テスト
func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		expected []string
	}{
		{
			name:     "重複単語の除去",
			words:    []string{"hello", "world", "hello", "test", "world"},
			expected: []string{"hello", "world", "test"},
		},
		{
			name:     "重複なし",
			words:    []string{"hello", "world", "test"},
			expected: []string{"hello", "world", "test"},
		},
		{
			name:     "空のリスト",
			words:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unique := RemoveDuplicates(tt.words)
			assert.ElementsMatch(t, tt.expected, unique)
		})
	}
}

// TestIsValidWord は有効な単語かどうかのテスト
func TestIsValidWord(t *testing.T) {
	tests := []struct {
		name     string
		word     string
		language string
		expected bool
	}{
		{
			name:     "有効な英語単語",
			word:     "hello",
			language: "en",
			expected: true,
		},
		{
			name:     "短すぎる単語（1文字）",
			word:     "a",
			language: "en",
			expected: false,
		},
		{
			name:     "数字のみ",
			word:     "123",
			language: "en",
			expected: false,
		},
		{
			name:     "記号のみ",
			word:     "!!!",
			language: "en",
			expected: false,
		},
		{
			name:     "有効なロシア語単語",
			word:     "привет",
			language: "ru",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := IsValidWord(tt.word, tt.language)
			assert.Equal(t, tt.expected, valid)
		})
	}
}

// TestExtractWordsWithContext はコンテキスト付きの単語抽出テスト
func TestExtractWordsWithContext(t *testing.T) {
	text := "Hello world! This is a test. Hello again!"
	language := "en"

	wordsWithContext := ExtractWordsWithContext(text, language)

	require.NotEmpty(t, wordsWithContext)

	// "Hello" が2回出現することを確認
	helloCount := 0
	for _, wc := range wordsWithContext {
		if wc.Word == "hello" {
			helloCount++
			assert.NotEmpty(t, wc.Context, "コンテキストが空ではないこと")
		}
	}
	assert.Equal(t, 2, helloCount, "Hello が2回出現すること")
}

// Benchmark tests
func BenchmarkExtractWords(b *testing.B) {
	text := "Здравствуйте! Как дела? Меня зовут Иван. Я изучаю русский язык."
	language := "ru"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExtractWords(text, language)
	}
}

func BenchmarkNormalizeWord(b *testing.B) {
	word := "Здравствуйте!"
	language := "ru"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NormalizeWord(word, language)
	}
}
