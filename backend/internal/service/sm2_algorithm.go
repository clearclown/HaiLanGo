package service

import (
	"math"
	"time"
)

// SM2Algorithm (SuperMemo 2) アルゴリズム実装
type SM2Algorithm struct{}

func NewSM2Algorithm() *SM2Algorithm {
	return &SM2Algorithm{}
}

// CalculateNextReview は次の復習日時を計算する
// score: 0-100 (30=思い出せない, 70=少し時間がかかった, 100=完璧)
func (s *SM2Algorithm) CalculateNextReview(
	currentEaseFactor float64,
	currentInterval int,
	score int,
) (nextInterval int, nextEaseFactor float64, nextReview time.Time) {

	// スコアを0-5の品質スケールに変換
	quality := s.scoreToQuality(score)

	// 新しい容易度係数を計算
	nextEaseFactor = currentEaseFactor + (0.1 - (5-float64(quality))*(0.08+(5-float64(quality))*0.02))
	if nextEaseFactor < 1.3 {
		nextEaseFactor = 1.3
	}

	// 次の間隔を計算
	if quality < 3 {
		// 失敗：最初からやり直し
		nextInterval = 1
	} else {
		if currentInterval == 0 {
			nextInterval = 1
		} else if currentInterval == 1 {
			nextInterval = 6
		} else {
			nextInterval = int(math.Round(float64(currentInterval) * nextEaseFactor))
		}
	}

	// 次の復習日時
	nextReview = time.Now().Add(time.Duration(nextInterval) * 24 * time.Hour)

	return nextInterval, nextEaseFactor, nextReview
}

func (s *SM2Algorithm) scoreToQuality(score int) int {
	switch {
	case score >= 90:
		return 5 // 完璧
	case score >= 70:
		return 4 // 正解だが努力が必要
	case score >= 50:
		return 3 // かろうじて正解
	case score >= 30:
		return 2 // 不正解だが覚えていた
	default:
		return 0 // 完全に忘れた
	}
}

// CalculatePriority は復習の優先度を計算する
func (s *SM2Algorithm) CalculatePriority(nextReview time.Time) string {
	now := time.Now()
	hoursUntil := nextReview.Sub(now).Hours()

	if hoursUntil <= 0 {
		return "urgent" // 期限切れ
	} else if hoursUntil <= 24 {
		return "urgent" // 今日中
	} else if hoursUntil <= 48 {
		return "recommended" // 明日まで
	} else {
		return "optional" // 余裕あり
	}
}
