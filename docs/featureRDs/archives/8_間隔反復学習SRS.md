# 機能実装: 間隔反復学習（SRS）

## 要件

### 機能要件

#### 間隔反復学習アルゴリズム（Spaced Repetition System）
1. **基本スケジュール**
   - 初回: 1日後
   - 2回目: 3日後
   - 3回目: 7日後
   - 4回目: 14日後
   - 5回目: 30日後
   - 6回目以降: 60日後

2. **調整ルール**
   - **スコア85点以上**: 次の間隔を1.5倍に延長
   - **スコア70-84点**: 通常の間隔
   - **スコア50-69点**: 間隔を半分に短縮
   - **スコア50点未満**: 翌日に再度復習

3. **復習項目の管理**
   - フレーズ単位で管理
   - 単語単位で管理（オプション）
   - 学習履歴の記録
   - 次回復習日の計算

4. **復習画面**
   - 緊急（今日中に復習が必要）
   - 推奨（今日復習すると効果的）
   - 余裕あり（明日以降でもOK）

### 非機能要件

- **パフォーマンス**:
  - 復習項目の取得: 100ms以内
  - スケジュール計算: 10ms以内
  - 大量データ（1000項目以上）への対応

- **セキュリティ**:
  - ユーザー認証必須
  - 学習データのプライバシー保護

- **拡張性**:
  - 将来的なアルゴリズム改善への対応
  - カスタムスケジュール設定（オプション）

- **エラーハンドリング**:
  - データ不整合の検出
  - 復習日の再計算
  - ユーザーへのエラー通知

## 実装指示

### Step 1: テスト設計

以下の順でテストを作成：

1. **ユニットテスト（関数/メソッドレベル）**
   - 間隔計算アルゴリズム
   - スコアに基づく調整
   - 次回復習日の計算
   - 復習項目のフィルタリング

2. **統合テスト（モジュール間）**
   - 復習項目の取得
   - 復習完了後の更新
   - スケジュールの再計算
   - 学習履歴の記録

3. **エッジケースのテスト**
   - 初回学習項目
   - 長期未復習項目
   - スコア0点の項目
   - 大量データ（1000項目以上）

テストファイルは `backend/internal/service/srs/srs_test.go` に配置。

テストは **実装前に実行してすべて失敗すること** を確認。

#### テスト例（Go）

```go
// backend/internal/service/srs/srs_test.go
package srs

import (
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestCalculateNextReviewDate(t *testing.T) {
    now := time.Now()

    // 初回学習
    nextDate := CalculateNextReviewDate(0, 0, now)
    expectedDate := now.AddDate(0, 0, 1)
    assert.Equal(t, expectedDate.Year(), nextDate.Year())
    assert.Equal(t, expectedDate.Month(), nextDate.Month())
    assert.Equal(t, expectedDate.Day(), nextDate.Day())

    // 2回目（スコア85点以上）
    nextDate = CalculateNextReviewDate(1, 90, now)
    expectedDate = now.AddDate(0, 0, 3*3/2) // 3日 * 1.5倍
    assert.Equal(t, expectedDate.Day(), nextDate.Day())

    // スコア50点未満
    nextDate = CalculateNextReviewDate(5, 40, now)
    expectedDate = now.AddDate(0, 0, 1) // 翌日
    assert.Equal(t, expectedDate.Day(), nextDate.Day())
}

func TestAdjustInterval(t *testing.T) {
    // 間隔調整のテスト
    baseInterval := 7 // 7日

    // スコア85点以上
    adjusted := AdjustInterval(baseInterval, 90)
    assert.Equal(t, 10, adjusted) // 7 * 1.5 = 10.5 → 10

    // スコア70-84点
    adjusted = AdjustInterval(baseInterval, 75)
    assert.Equal(t, 7, adjusted)

    // スコア50-69点
    adjusted = AdjustInterval(baseInterval, 60)
    assert.Equal(t, 3, adjusted) // 7 / 2 = 3.5 → 3
}

func TestGetReviewItems(t *testing.T) {
    // 復習項目取得のテスト
    ctx := context.Background()
    userID := "test-user"

    items, err := GetReviewItems(ctx, userID, time.Now())
    require.NoError(t, err)

    // 緊急項目の確認
    urgentItems := FilterUrgentItems(items, time.Now())
    assert.NotNil(t, urgentItems)
}

func TestUpdateReviewHistory(t *testing.T) {
    // 復習履歴更新のテスト
    ctx := context.Background()
    itemID := "test-item"
    score := 85

    err := UpdateReviewHistory(ctx, itemID, score)
    require.NoError(t, err)

    // 次回復習日の確認
    item, err := GetReviewItem(ctx, itemID)
    require.NoError(t, err)
    assert.NotNil(t, item.NextReviewDate)
}
```

### Step 2: 実装

- テストが通るように実装
- コードコメントは日本語で
- 型安全性を最優先
- エラーハンドリングを必ず含める

#### 実装ファイル構造

```
backend/
├── internal/
│   ├── service/
│   │   └── srs/
│   │       ├── srs.go               # SRSアルゴリズム
│   │       ├── scheduler.go         # スケジューラー
│   │       └── srs_test.go
│   └── repository/
│       ├── review_item.go           # 復習項目DB操作
│       └── review_item_test.go
└── pkg/
    └── srs/
        ├── algorithm.go             # SRSアルゴリズム実装
        └── algorithm_test.go
```

#### 主要な実装ポイント

1. **SRSアルゴリズム**
   ```go
   // pkg/srs/algorithm.go
   func CalculateNextReviewDate(reviewCount int, score int, lastReviewDate time.Time) time.Time {
       // 基本間隔の取得
       baseInterval := GetBaseInterval(reviewCount)

       // スコアに基づく調整
       adjustedInterval := AdjustInterval(baseInterval, score)

       // 次回復習日の計算
       return lastReviewDate.AddDate(0, 0, adjustedInterval)
   }

   func AdjustInterval(baseInterval int, score int) int {
       if score >= 85 {
           return int(float64(baseInterval) * 1.5)
       } else if score >= 70 {
           return baseInterval
       } else if score >= 50 {
           return baseInterval / 2
       } else {
           return 1 // 翌日
       }
   }
   ```

2. **復習項目取得**
   ```go
   // internal/service/srs/srs.go
   func (s *SRSService) GetReviewItems(ctx context.Context, userID string, now time.Time) ([]*ReviewItem, error) {
       // 1. ユーザーの全復習項目取得
       // 2. 次回復習日でフィルタリング
       // 3. 優先度でソート
       // 4. 返却
   }

   func (s *SRSService) UpdateReview(ctx context.Context, itemID string, score int) error {
       // 1. 復習履歴記録
       // 2. 次回復習日計算
       // 3. DB更新
   }
   ```

### Step 3: リファクタリング

- DRY原則の適用
- パフォーマンス最適化（インデックス最適化）
- コードの可読性向上

### Step 4: ドキュメント

- README.md への機能説明追加
- APIドキュメント
- SRSアルゴリズムの説明

#### APIエンドポイント

```
GET  /api/v1/review/items
POST /api/v1/review/items/{item_id}/complete
GET  /api/v1/review/stats
```

## 制約事項

- 既存の `backend/internal/models/` は変更しない（新規追加のみ）
- `backend/internal/service/srs/` 配下のみ編集可能
- 依存関係の追加は `go.mod` のみ

## 完了条件

- [ ] すべてのテストが通る
- [ ] lintエラーがない
- [ ] タイプエラーがない
- [ ] ドキュメントが更新されている
- [ ] SRSアルゴリズムの動作確認
- [ ] 大量データでのパフォーマンス確認

## 追加のテスト要件

### セキュリティテスト
- [ ] ユーザー認証の動作確認
- [ ] 学習データのプライバシー保護

### パフォーマンステスト
- [ ] 復習項目取得時間（100ms以内）
- [ ] スケジュール計算時間（10ms以内）
- [ ] 大量データ（1000項目以上）での動作確認

### E2Eテスト
- [ ] 復習画面の表示
- [ ] 復習完了後の更新
- [ ] 次回復習日の計算
