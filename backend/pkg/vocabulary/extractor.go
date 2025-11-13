package vocabulary

import (
	"regexp"
	"strings"
	"unicode"
)

// WordWithContext はコンテキスト付きの単語
type WordWithContext struct {
	Word    string
	Context string // 単語が出現した文脈
}

// ExtractWords はテキストから単語を抽出する
func ExtractWords(text string, language string) []string {
	// 1. テキストを単語に分割
	words := tokenize(text, language)

	// 2. 正規化（小文字化など）
	normalized := make([]string, 0, len(words))
	for _, word := range words {
		normalizedWord := NormalizeWord(word, language)
		if normalizedWord != "" {
			normalized = append(normalized, normalizedWord)
		}
	}

	// 3. ストップワード除去
	filtered := RemoveStopWords(normalized, language)

	// 4. 有効な単語のみをフィルタ
	valid := make([]string, 0, len(filtered))
	for _, word := range filtered {
		if IsValidWord(word, language) {
			valid = append(valid, word)
		}
	}

	// 5. 重複除去
	unique := RemoveDuplicates(valid)

	return unique
}

// tokenize はテキストを単語に分割する
func tokenize(text string, language string) []string {
	switch language {
	case "ja", "zh": // 日本語・中国語の場合
		return tokenizeCJK(text)
	default: // その他の言語
		return tokenizeDefault(text)
	}
}

// tokenizeDefault はデフォルトのトークン化（スペース区切り）
func tokenizeDefault(text string) []string {
	// 句読点を削除してスペースで分割
	re := regexp.MustCompile(`[^\p{L}\s]+`)
	cleaned := re.ReplaceAllString(text, " ")
	words := strings.Fields(cleaned)
	return words
}

// tokenizeCJK は中国語・日本語のトークン化
func tokenizeCJK(text string) []string {
	// 簡易的な実装: 文字ごとに分割し、連続する同じ文字種をグループ化
	words := make([]string, 0)
	var currentWord strings.Builder
	var lastType int // 0: その他, 1: ひらがな, 2: カタカナ, 3: 漢字

	for _, r := range text {
		var currentType int
		if isHiragana(r) {
			currentType = 1
		} else if isKatakana(r) {
			currentType = 2
		} else if isKanji(r) {
			currentType = 3
		} else {
			currentType = 0
		}

		if currentType == 0 {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
			lastType = 0
			continue
		}

		if lastType != 0 && lastType != currentType {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		}

		currentWord.WriteRune(r)
		lastType = currentType
	}

	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	return words
}

// isHiragana はひらがなかどうか判定
func isHiragana(r rune) bool {
	return r >= '\u3040' && r <= '\u309F'
}

// isKatakana はカタカナかどうか判定
func isKatakana(r rune) bool {
	return r >= '\u30A0' && r <= '\u30FF'
}

// isKanji は漢字かどうか判定
func isKanji(r rune) bool {
	return (r >= '\u4E00' && r <= '\u9FFF') || // CJK統合漢字
		(r >= '\u3400' && r <= '\u4DBF') // CJK統合漢字拡張A
}

// NormalizeWord は単語を正規化する
func NormalizeWord(word string, language string) string {
	// 小文字に変換
	normalized := strings.ToLower(word)

	// 句読点を除去
	re := regexp.MustCompile(`[^\p{L}]+`)
	normalized = re.ReplaceAllString(normalized, "")

	return normalized
}

// RemoveStopWords はストップワードを除去する
func RemoveStopWords(words []string, language string) []string {
	stopWords := getStopWords(language)
	filtered := make([]string, 0, len(words))

	for _, word := range words {
		if !contains(stopWords, word) {
			filtered = append(filtered, word)
		}
	}

	return filtered
}

// getStopWords は言語ごとのストップワードリストを返す
func getStopWords(language string) []string {
	stopWordMap := map[string][]string{
		"en": {"the", "a", "an", "and", "or", "but", "in", "on", "at", "to", "for", "of", "with", "is", "are", "was", "were", "be", "been", "being"},
		"ru": {"и", "в", "не", "на", "я", "что", "он", "с", "как", "а", "то", "это", "она", "по", "но", "они", "мы"},
		"ja": {"は", "が", "を", "に", "の", "と", "で", "や", "も", "から", "まで", "より", "ですか", "です", "ます"},
		"zh": {"的", "了", "在", "是", "我", "有", "和", "人", "这", "中", "大", "为", "上", "个", "国", "一"},
	}

	if stopWords, ok := stopWordMap[language]; ok {
		return stopWords
	}
	return []string{}
}

// RemoveDuplicates は重複を除去する
func RemoveDuplicates(words []string) []string {
	seen := make(map[string]bool)
	unique := make([]string, 0, len(words))

	for _, word := range words {
		if !seen[word] {
			seen[word] = true
			unique = append(unique, word)
		}
	}

	return unique
}

// IsValidWord は有効な単語かどうか判定する
func IsValidWord(word string, language string) bool {
	// 空文字はNG
	if word == "" {
		return false
	}

	// 短すぎる単語はNG（1文字）
	if len([]rune(word)) < 2 {
		return false
	}

	// 数字のみはNG
	if isNumeric(word) {
		return false
	}

	// 記号のみはNG
	if isPunctuation(word) {
		return false
	}

	return true
}

// isNumeric は数字のみかどうか判定
func isNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// isPunctuation は記号のみかどうか判定
func isPunctuation(s string) bool {
	for _, r := range s {
		if !unicode.IsPunct(r) {
			return false
		}
	}
	return true
}

// contains はスライスに要素が含まれるか判定
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ExtractWordsWithContext はコンテキスト付きで単語を抽出する
func ExtractWordsWithContext(text string, language string) []WordWithContext {
	sentences := splitIntoSentences(text)
	wordsWithContext := make([]WordWithContext, 0)

	for _, sentence := range sentences {
		words := ExtractWords(sentence, language)
		for _, word := range words {
			wordsWithContext = append(wordsWithContext, WordWithContext{
				Word:    word,
				Context: sentence,
			})
		}
	}

	return wordsWithContext
}

// splitIntoSentences はテキストを文に分割する
func splitIntoSentences(text string) []string {
	// 簡易的な実装: 。！？.!?で分割
	re := regexp.MustCompile(`[。！？.!?]+`)
	sentences := re.Split(text, -1)

	// 空の文を除去
	filtered := make([]string, 0, len(sentences))
	for _, sentence := range sentences {
		trimmed := strings.TrimSpace(sentence)
		if trimmed != "" {
			filtered = append(filtered, trimmed)
		}
	}

	return filtered
}
