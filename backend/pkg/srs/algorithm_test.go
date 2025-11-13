package srs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGetBaseInterval は基本間隔の取得をテスト
func TestGetBaseInterval(t *testing.T) {
	tests := []struct {
		name         string
		reviewCount  int
		wantInterval int
	}{
		{"初回学習", 0, 1},
		{"2回目", 1, 3},
		{"3回目", 2, 7},
		{"4回目", 3, 14},
		{"5回目", 4, 30},
		{"6回目", 5, 60},
		{"7回目以降", 6, 60},
		{"10回目以降", 10, 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetBaseInterval(tt.reviewCount)
			assert.Equal(t, tt.wantInterval, got)
		})
	}
}

// TestAdjustInterval はスコアに基づく間隔調整をテスト
func TestAdjustInterval(t *testing.T) {
	tests := []struct {
		name         string
		baseInterval int
		score        int
		wantInterval int
	}{
		{"スコア85点以上（1.5倍）", 10, 90, 15},
		{"スコア85点ちょうど（1.5倍）", 10, 85, 15},
		{"スコア84点（通常）", 10, 84, 10},
		{"スコア70点（通常）", 10, 70, 10},
		{"スコア69点（半分）", 10, 69, 5},
		{"スコア50点（半分）", 10, 50, 5},
		{"スコア49点（翌日）", 10, 49, 1},
		{"スコア0点（翌日）", 10, 0, 1},
		{"小数点の切り捨て（1.5倍）", 7, 90, 10}, // 7 * 1.5 = 10.5 → 10
		{"小数点の切り捨て（半分）", 7, 60, 3},   // 7 / 2 = 3.5 → 3
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AdjustInterval(tt.baseInterval, tt.score)
			assert.Equal(t, tt.wantInterval, got)
		})
	}
}

// TestCalculateNextReviewDate は次回復習日の計算をテスト
func TestCalculateNextReviewDate(t *testing.T) {
	baseTime := time.Date(2025, 11, 13, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		reviewCount    int
		score          int
		lastReviewDate time.Time
		wantDate       time.Time
	}{
		{
			name:           "初回学習（1日後）",
			reviewCount:    0,
			score:          0,
			lastReviewDate: baseTime,
			wantDate:       baseTime.AddDate(0, 0, 1),
		},
		{
			name:           "2回目・高得点（3日 * 1.5 = 4.5日 → 4日）",
			reviewCount:    1,
			score:          90,
			lastReviewDate: baseTime,
			wantDate:       baseTime.AddDate(0, 0, 4),
		},
		{
			name:           "3回目・普通（7日）",
			reviewCount:    2,
			score:          75,
			lastReviewDate: baseTime,
			wantDate:       baseTime.AddDate(0, 0, 7),
		},
		{
			name:           "4回目・低得点（14日 / 2 = 7日）",
			reviewCount:    3,
			score:          60,
			lastReviewDate: baseTime,
			wantDate:       baseTime.AddDate(0, 0, 7),
		},
		{
			name:           "5回目・最低得点（翌日）",
			reviewCount:    4,
			score:          40,
			lastReviewDate: baseTime,
			wantDate:       baseTime.AddDate(0, 0, 1),
		},
		{
			name:           "6回目以降（60日）",
			reviewCount:    5,
			score:          80,
			lastReviewDate: baseTime,
			wantDate:       baseTime.AddDate(0, 0, 60),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateNextReviewDate(tt.reviewCount, tt.score, tt.lastReviewDate)
			assert.Equal(t, tt.wantDate.Year(), got.Year())
			assert.Equal(t, tt.wantDate.Month(), got.Month())
			assert.Equal(t, tt.wantDate.Day(), got.Day())
		})
	}
}

// TestCalculateNextReviewDateEdgeCases はエッジケースをテスト
func TestCalculateNextReviewDateEdgeCases(t *testing.T) {
	baseTime := time.Date(2025, 11, 13, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		reviewCount    int
		score          int
		lastReviewDate time.Time
		description    string
	}{
		{
			name:           "非常に大きな復習回数",
			reviewCount:    100,
			score:          80,
			lastReviewDate: baseTime,
			description:    "60日間隔を維持",
		},
		{
			name:           "負の復習回数（エラーケース）",
			reviewCount:    -1,
			score:          80,
			lastReviewDate: baseTime,
			description:    "初回として扱う",
		},
		{
			name:           "100点超過",
			reviewCount:    1,
			score:          150,
			lastReviewDate: baseTime,
			description:    "1.5倍として扱う",
		},
		{
			name:           "負のスコア",
			reviewCount:    1,
			score:          -10,
			lastReviewDate: baseTime,
			description:    "翌日として扱う",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateNextReviewDate(tt.reviewCount, tt.score, tt.lastReviewDate)
			assert.NotNil(t, got, tt.description)
			// 結果が過去の日付でないことを確認
			assert.False(t, got.Before(tt.lastReviewDate), "次回復習日は過去の日付であってはならない")
		})
	}
}

// TestShouldReviewToday は今日復習すべきかを判定するテスト
func TestShouldReviewToday(t *testing.T) {
	now := time.Date(2025, 11, 13, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		nextReviewDate time.Time
		want           bool
	}{
		{"過去の日付", now.AddDate(0, 0, -1), true},
		{"今日", now, true},
		{"明日", now.AddDate(0, 0, 1), false},
		{"1週間後", now.AddDate(0, 0, 7), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldReviewToday(tt.nextReviewDate, now)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestGetPriorityLevel は優先度レベルの判定をテスト
func TestGetPriorityLevel(t *testing.T) {
	now := time.Date(2025, 11, 13, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		nextReviewDate time.Time
		wantPriority   string
	}{
		{"過去の日付（緊急）", now.AddDate(0, 0, -1), "urgent"},
		{"今日（緊急）", now, "urgent"},
		{"明日（推奨）", now.AddDate(0, 0, 1), "recommended"},
		{"2日後（余裕あり）", now.AddDate(0, 0, 2), "relaxed"},
		{"1週間後（余裕あり）", now.AddDate(0, 0, 7), "relaxed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPriorityLevel(tt.nextReviewDate, now)
			assert.Equal(t, tt.wantPriority, got)
		})
	}
}

// BenchmarkCalculateNextReviewDate はパフォーマンステスト
func BenchmarkCalculateNextReviewDate(b *testing.B) {
	baseTime := time.Now()
	for i := 0; i < b.N; i++ {
		CalculateNextReviewDate(5, 85, baseTime)
	}
}

// BenchmarkAdjustInterval はパフォーマンステスト
func BenchmarkAdjustInterval(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AdjustInterval(30, 75)
	}
}
