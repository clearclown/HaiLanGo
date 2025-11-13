package srs

import (
	"time"
)

// GetBaseInterval は復習回数に基づく基本間隔（日数）を返す
// 初回: 1日後、2回目: 3日後、3回目: 7日後、4回目: 14日後、5回目: 30日後、6回目以降: 60日後
func GetBaseInterval(reviewCount int) int {
	// 負の復習回数は初回として扱う
	if reviewCount < 0 {
		reviewCount = 0
	}

	intervals := []int{1, 3, 7, 14, 30, 60}

	if reviewCount >= len(intervals) {
		return intervals[len(intervals)-1] // 6回目以降は60日
	}

	return intervals[reviewCount]
}

// AdjustInterval はスコアに基づいて間隔を調整する
// スコア85点以上: 1.5倍、70-84点: 通常、50-69点: 半分、50点未満: 翌日
func AdjustInterval(baseInterval int, score int) int {
	if score >= 85 {
		// 1.5倍に延長（小数点切り捨て）
		return int(float64(baseInterval) * 1.5)
	} else if score >= 70 {
		// 通常の間隔
		return baseInterval
	} else if score >= 50 {
		// 半分に短縮（小数点切り捨て）
		return baseInterval / 2
	} else {
		// 翌日に復習
		return 1
	}
}

// CalculateNextReviewDate は次回復習日を計算する
func CalculateNextReviewDate(reviewCount int, score int, lastReviewDate time.Time) time.Time {
	baseInterval := GetBaseInterval(reviewCount)
	adjustedInterval := AdjustInterval(baseInterval, score)

	return lastReviewDate.AddDate(0, 0, adjustedInterval)
}

// ShouldReviewToday は今日復習すべきかを判定する
func ShouldReviewToday(nextReviewDate time.Time, now time.Time) bool {
	// 次回復習日が今日以前なら復習すべき
	return nextReviewDate.Before(now) || isSameDay(nextReviewDate, now)
}

// GetPriorityLevel は次回復習日に基づいて優先度レベルを返す
// urgent: 今日中に復習が必要、recommended: 今日復習すると効果的、relaxed: 明日以降でもOK
func GetPriorityLevel(nextReviewDate time.Time, now time.Time) string {
	daysDiff := daysBetween(now, nextReviewDate)

	if daysDiff < 0 {
		// 過去の日付 = 緊急
		return "urgent"
	} else if daysDiff == 0 {
		// 今日 = 緊急
		return "urgent"
	} else if daysDiff == 1 {
		// 明日 = 推奨
		return "recommended"
	} else {
		// 2日後以降 = 余裕あり
		return "relaxed"
	}
}

// isSameDay は2つの時刻が同じ日かを判定する
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// daysBetween は2つの時刻の間の日数を返す
func daysBetween(from, to time.Time) int {
	// 時刻をリセットして日付のみで比較
	fromDate := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	toDate := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, to.Location())

	duration := toDate.Sub(fromDate)
	days := int(duration.Hours() / 24)

	return days
}
