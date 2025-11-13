package stt

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

// CalculateAccuracyScore は正確性スコアを計算する（0-100点）
func CalculateAccuracyScore(expected, recognized string) int {
	// 大文字小文字を無視して比較
	expectedLower := strings.ToLower(strings.TrimSpace(expected))
	recognizedLower := strings.ToLower(strings.TrimSpace(recognized))

	// 完全一致の場合
	if expectedLower == recognizedLower {
		return 100
	}

	// Levenshtein距離を使用して類似度を計算
	distance := levenshteinDistance(expectedLower, recognizedLower)
	maxLen := math.Max(float64(len(expectedLower)), float64(len(recognizedLower)))

	if maxLen == 0 {
		return 0
	}

	// 類似度を0-100のスコアに変換
	similarity := 1.0 - (float64(distance) / maxLen)
	score := int(similarity * 100)

	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}

	return score
}

// CalculateFluencyScore は流暢性スコアを計算する（0-100点）
func CalculateFluencyScore(words []models.WordInfo, duration float64) int {
	if len(words) == 0 || duration == 0 {
		return 0
	}

	// 単語あたりの平均時間を計算（秒）
	avgTimePerWord := duration / float64(len(words))

	// 理想的な単語あたりの時間（秒）- 自然な会話速度
	// 英語では約0.5秒/単語が自然
	idealTime := 0.5

	// 間隔の安定性を計算
	var gaps []float64
	for i := 0; i < len(words)-1; i++ {
		gap := words[i+1].StartTime - words[i].EndTime
		gaps = append(gaps, gap)
	}

	// 間隔の標準偏差を計算
	gapVariance := calculateVariance(gaps)

	// スコア計算
	// 1. 速度スコア（理想的な速度に近いほど高い）
	speedDiff := math.Abs(avgTimePerWord - idealTime)
	speedScore := math.Max(0, 100-speedDiff*100)

	// 2. 安定性スコア（間隔が安定しているほど高い）
	stabilityScore := math.Max(0, 100-gapVariance*200)

	// 総合スコア（速度と安定性の平均）
	totalScore := (speedScore + stabilityScore) / 2

	return int(totalScore)
}

// CalculatePronunciationScore は発音スコアを計算する（0-100点）
func CalculatePronunciationScore(expectedWords, recognizedWords []models.WordInfo) int {
	if len(expectedWords) == 0 {
		return 0
	}

	totalScore := 0
	matchCount := 0

	// 各単語の発音スコアを計算
	for i := 0; i < len(expectedWords) && i < len(recognizedWords); i++ {
		expectedWord := strings.ToLower(expectedWords[i].Word)
		recognizedWord := strings.ToLower(recognizedWords[i].Word)

		wordScore := CalculateAccuracyScore(expectedWord, recognizedWord)
		totalScore += wordScore
		matchCount++
	}

	if matchCount == 0 {
		return 0
	}

	return totalScore / matchCount
}

// GenerateFeedback はスコアに基づいてフィードバックを生成する
func GenerateFeedback(score *models.PronunciationScore) *models.Feedback {
	feedback := &models.Feedback{
		PositivePoints: []string{},
		Improvements:   []string{},
		SpecificAdvice: []string{},
	}

	// レベルとメッセージを決定
	if score.TotalScore >= ScoreExcellentThreshold {
		feedback.Level = FeedbackLevelExcellent
		feedback.Message = "🎉 素晴らしい！完璧に近い発音です。"
		feedback.PositivePoints = append(feedback.PositivePoints,
			"発音が非常に明瞭です",
			"イントネーションが自然です",
			"流暢に話せています",
		)
	} else if score.TotalScore >= ScoreGoodThreshold {
		feedback.Level = FeedbackLevelGood
		feedback.Message = "👍 良好です！もう少しで完璧です。"
		feedback.PositivePoints = append(feedback.PositivePoints,
			"基本的な発音は正確です",
			"理解しやすい発音です",
		)
		feedback.Improvements = append(feedback.Improvements,
			"いくつかの単語の発音を改善できます",
		)
	} else if score.TotalScore >= ScoreFairThreshold {
		feedback.Level = FeedbackLevelFair
		feedback.Message = "💪 頑張りましょう！改善の余地があります。"
		feedback.Improvements = append(feedback.Improvements,
			"発音の正確性を向上させましょう",
			"単語の区切りを意識しましょう",
		)
	} else {
		feedback.Level = FeedbackLevelPoor
		feedback.Message = "📚 練習を重ねましょう。"
		feedback.Improvements = append(feedback.Improvements,
			"基本的な発音から練習しましょう",
			"ゆっくり丁寧に発音しましょう",
		)
	}

	// 具体的なアドバイスを生成
	if score.AccuracyScore < 80 {
		feedback.SpecificAdvice = append(feedback.SpecificAdvice,
			"正確な発音を意識してください",
		)
	}

	if score.FluencyScore < 70 {
		feedback.SpecificAdvice = append(feedback.SpecificAdvice,
			"自然なリズムで話すように心がけてください",
		)
	}

	if score.PronuncScore < 75 {
		feedback.SpecificAdvice = append(feedback.SpecificAdvice,
			"個々の音素をはっきりと発音しましょう",
		)
	}

	// 単語レベルの改善点を追加
	for _, wordScore := range score.WordScores {
		if !wordScore.IsCorrect && wordScore.Score < 70 {
			advice := fmt.Sprintf("「%s」の発音を練習してください（認識結果: %s）",
				wordScore.ExpectedWord, wordScore.RecognizedWord)
			feedback.SpecificAdvice = append(feedback.SpecificAdvice, advice)
		}
	}

	return feedback
}

// levenshteinDistance はLevenshtein距離を計算する
func levenshteinDistance(s1, s2 string) int {
	r1 := []rune(s1)
	r2 := []rune(s2)

	len1 := len(r1)
	len2 := len(r2)

	// 動的計画法でLevenshtein距離を計算
	matrix := make([][]int, len1+1)
	for i := range matrix {
		matrix[i] = make([]int, len2+1)
		matrix[i][0] = i
	}

	for j := 0; j <= len2; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if !unicode.IsSpace(r1[i-1]) && !unicode.IsSpace(r2[j-1]) && r1[i-1] != r2[j-1] {
				cost = 1
			}

			matrix[i][j] = min3(
				matrix[i-1][j]+1,      // 削除
				matrix[i][j-1]+1,      // 挿入
				matrix[i-1][j-1]+cost, // 置換
			)
		}
	}

	return matrix[len1][len2]
}

// min3 は3つの整数の最小値を返す
func min3(a, b, c int) int {
	min := a
	if b < min {
		min = b
	}
	if c < min {
		min = c
	}
	return min
}

// calculateVariance は分散を計算する
func calculateVariance(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	// 平均を計算
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	// 分散を計算
	variance := 0.0
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}

	return variance / float64(len(values))
}
