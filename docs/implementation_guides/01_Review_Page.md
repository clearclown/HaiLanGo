# 実装指示書: 復習ページ (Review Page)

## 概要
間隔反復学習（SRS）に基づいた復習機能の実装。ユーザーが学習した単語・フレーズを最適なタイミングで復習できる機能を提供する。

## 担当範囲
- **フロントエンド**: `frontend/web/app/review/page.tsx`
- **コンポーネント**: `frontend/web/components/review/*`
- **バックエンドAPI**: すでに実装済み（`/api/v1/review/*`）

## 前提条件
- Node.js 18+、pnpm がインストール済み
- バックエンドAPI が http://localhost:8080 で起動中
- TypeScript、React、Next.js の基本知識

## 実装ステップ

### Step 1: 型定義の作成

**ファイル**: `frontend/web/types/review.ts`

```typescript
export type ReviewPriority = 'urgent' | 'recommended' | 'optional';

export interface ReviewItem {
  id: string;
  type: 'word' | 'phrase';
  text: string;
  translation: string;
  language: string;
  mastery_level: number; // 0-100
  last_reviewed: string;
  next_review: string;
  priority: ReviewPriority;
}

export interface ReviewStats {
  urgent_count: number;
  recommended_count: number;
  optional_count: number;
  total_completed_today: number;
  weekly_completion_rate: number;
}

export interface ReviewResult {
  item_id: string;
  score: number; // 0-100
  pronunciation_score?: number;
  completed_at: string;
}
```

### Step 2: API クライアントの拡張

**ファイル**: `frontend/web/lib/api/client.ts`

**追加する内容**:

```typescript
// インポートに追加
import type { ReviewItem, ReviewStats, ReviewResult } from '@/types/review';

// APIClient クラス内に追加
review = {
  getStats: async (): Promise<ReviewStats> => {
    return this.fetch<ReviewStats>('/api/v1/review/stats');
  },

  getItems: async (priority?: 'urgent' | 'recommended' | 'optional'): Promise<{ items: ReviewItem[] }> => {
    const query = priority ? `?priority=${priority}` : '';
    return this.fetch<{ items: ReviewItem[] }>(`/api/v1/review/items${query}`);
  },

  submit: async (result: ReviewResult): Promise<{ success: boolean; next_review: string }> => {
    return this.fetch<{ success: boolean; next_review: string }>('/api/v1/review/submit', {
      method: 'POST',
      body: JSON.stringify(result),
    });
  },
};
```

### Step 3: ReviewCard コンポーネントの作成

**ファイル**: `frontend/web/components/review/ReviewCard.tsx`

```typescript
'use client';

import { useState } from 'react';
import type { ReviewItem, ReviewPriority } from '@/types/review';

interface ReviewCardProps {
  items: ReviewItem[];
  priority: ReviewPriority;
  onStartReview: () => void;
}

export function ReviewCard({ items, priority, onStartReview }: ReviewCardProps) {
  const getPriorityConfig = (priority: ReviewPriority) => {
    switch (priority) {
      case 'urgent':
        return {
          color: 'red',
          bgColor: 'bg-red-50',
          textColor: 'text-red-600',
          borderColor: 'border-red-200',
          icon: '🔴',
          title: '緊急',
          description: '今日中に復習が必要',
        };
      case 'recommended':
        return {
          color: 'yellow',
          bgColor: 'bg-yellow-50',
          textColor: 'text-yellow-600',
          borderColor: 'border-yellow-200',
          icon: '🟡',
          title: '推奨',
          description: '今日復習すると効果的',
        };
      case 'optional':
        return {
          color: 'green',
          bgColor: 'bg-green-50',
          textColor: 'text-green-600',
          borderColor: 'border-green-200',
          icon: '🟢',
          title: '余裕あり',
          description: '明日以降でもOK',
        };
    }
  };

  const config = getPriorityConfig(priority);

  return (
    <div className={`rounded-lg border-2 ${config.borderColor} ${config.bgColor} p-6`}>
      <div className="flex items-center gap-3 mb-3">
        <span className="text-2xl">{config.icon}</span>
        <div>
          <h3 className={`text-lg font-semibold ${config.textColor}`}>
            {config.title} ({items.length}項目)
          </h3>
          <p className="text-sm text-gray-600">{config.description}</p>
        </div>
      </div>

      {items.length > 0 && (
        <button
          type="button"
          onClick={onStartReview}
          className={`w-full mt-4 px-4 py-3 bg-${config.color}-500 text-white rounded-lg hover:bg-${config.color}-600 transition-colors font-medium`}
        >
          復習する
        </button>
      )}
    </div>
  );
}
```

### Step 4: ReviewSession コンポーネントの作成

**ファイル**: `frontend/web/components/review/ReviewSession.tsx`

```typescript
'use client';

import { useState } from 'react';
import type { ReviewItem } from '@/types/review';
import { apiClient } from '@/lib/api/client';

interface ReviewSessionProps {
  items: ReviewItem[];
  onComplete: () => void;
  onCancel: () => void;
}

export function ReviewSession({ items, onComplete, onCancel }: ReviewSessionProps) {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [showTranslation, setShowTranslation] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const currentItem = items[currentIndex];
  const progress = ((currentIndex + 1) / items.length) * 100;

  const handleScore = async (score: number) => {
    setIsSubmitting(true);
    try {
      await apiClient.review.submit({
        item_id: currentItem.id,
        score,
        completed_at: new Date().toISOString(),
      });

      if (currentIndex < items.length - 1) {
        setCurrentIndex(currentIndex + 1);
        setShowTranslation(false);
      } else {
        onComplete();
      }
    } catch (error) {
      console.error('Failed to submit review:', error);
      alert('復習の送信に失敗しました');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-8 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold">復習セッション</h2>
          <button
            type="button"
            onClick={onCancel}
            className="text-gray-500 hover:text-gray-700"
          >
            ✕
          </button>
        </div>

        {/* Progress Bar */}
        <div className="mb-6">
          <div className="flex justify-between text-sm text-gray-600 mb-2">
            <span>{currentIndex + 1} / {items.length}</span>
            <span>{Math.round(progress)}%</span>
          </div>
          <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
            <div
              className="h-full bg-blue-500 transition-all duration-300"
              style={{ width: `${progress}%` }}
            />
          </div>
        </div>

        {/* Question */}
        <div className="mb-8">
          <div className="text-center mb-4">
            <span className="text-sm text-gray-500 uppercase">
              {currentItem.type === 'word' ? '単語' : 'フレーズ'}
            </span>
          </div>
          <div className="text-4xl font-bold text-center mb-4">
            {currentItem.text}
          </div>

          {/* Show Translation Button */}
          {!showTranslation && (
            <button
              type="button"
              onClick={() => setShowTranslation(true)}
              className="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200"
            >
              翻訳を表示
            </button>
          )}

          {/* Translation */}
          {showTranslation && (
            <div className="text-2xl text-center text-gray-700 mb-6">
              {currentItem.translation}
            </div>
          )}
        </div>

        {/* Score Buttons */}
        {showTranslation && (
          <div className="space-y-3">
            <h3 className="font-semibold text-center mb-4">どれくらい覚えていましたか？</h3>
            <button
              type="button"
              onClick={() => handleScore(100)}
              disabled={isSubmitting}
              className="w-full px-4 py-3 bg-green-500 text-white rounded-lg hover:bg-green-600 disabled:opacity-50"
            >
              🟢 完璧に覚えていた
            </button>
            <button
              type="button"
              onClick={() => handleScore(70)}
              disabled={isSubmitting}
              className="w-full px-4 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50"
            >
              🔵 少し時間がかかった
            </button>
            <button
              type="button"
              onClick={() => handleScore(30)}
              disabled={isSubmitting}
              className="w-full px-4 py-3 bg-red-500 text-white rounded-lg hover:bg-red-600 disabled:opacity-50"
            >
              🔴 思い出せなかった
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
```

### Step 5: Review ページの実装

**ファイル**: `frontend/web/app/review/page.tsx`

```typescript
'use client';

import { useEffect, useState } from 'react';
import { apiClient } from '@/lib/api/client';
import type { ReviewItem, ReviewStats } from '@/types/review';
import { ReviewCard } from '@/components/review/ReviewCard';
import { ReviewSession } from '@/components/review/ReviewSession';

export default function ReviewPage() {
  const [stats, setStats] = useState<ReviewStats | null>(null);
  const [urgentItems, setUrgentItems] = useState<ReviewItem[]>([]);
  const [recommendedItems, setRecommendedItems] = useState<ReviewItem[]>([]);
  const [optionalItems, setOptionalItems] = useState<ReviewItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [activeSession, setActiveSession] = useState<ReviewItem[] | null>(null);

  useEffect(() => {
    loadReviewData();
  }, []);

  const loadReviewData = async () => {
    try {
      setIsLoading(true);
      const [statsData, urgentData, recommendedData, optionalData] = await Promise.all([
        apiClient.review.getStats(),
        apiClient.review.getItems('urgent'),
        apiClient.review.getItems('recommended'),
        apiClient.review.getItems('optional'),
      ]);

      setStats(statsData);
      setUrgentItems(urgentData.items);
      setRecommendedItems(recommendedData.items);
      setOptionalItems(optionalData.items);
    } catch (error) {
      console.error('Failed to load review data:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleStartSession = (items: ReviewItem[]) => {
    setActiveSession(items);
  };

  const handleCompleteSession = () => {
    setActiveSession(null);
    loadReviewData(); // Reload data after completing session
  };

  const handleCancelSession = () => {
    setActiveSession(null);
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background-secondary flex items-center justify-center">
        <div className="text-gray-600">読み込み中...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background-secondary">
      <div className="max-w-4xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="flex items-center gap-3 mb-8">
          <h1 className="text-3xl font-bold">復習</h1>
          {stats && stats.total_completed_today > 0 && (
            <span className="text-2xl">🔥</span>
          )}
        </div>

        {/* Stats */}
        {stats && (
          <div className="grid grid-cols-2 gap-4 mb-8">
            <div className="bg-white rounded-lg p-4">
              <h3 className="text-sm text-gray-600 mb-1">今日の復習</h3>
              <p className="text-2xl font-bold">{stats.total_completed_today}項目</p>
            </div>
            <div className="bg-white rounded-lg p-4">
              <h3 className="text-sm text-gray-600 mb-1">今週の達成率</h3>
              <p className="text-2xl font-bold">{stats.weekly_completion_rate}%</p>
            </div>
          </div>
        )}

        {/* Review Cards */}
        <div className="space-y-4">
          <ReviewCard
            items={urgentItems}
            priority="urgent"
            onStartReview={() => handleStartSession(urgentItems)}
          />
          <ReviewCard
            items={recommendedItems}
            priority="recommended"
            onStartReview={() => handleStartSession(recommendedItems)}
          />
          <ReviewCard
            items={optionalItems}
            priority="optional"
            onStartReview={() => handleStartSession(optionalItems)}
          />
        </div>

        {/* No Reviews Message */}
        {urgentItems.length === 0 && recommendedItems.length === 0 && optionalItems.length === 0 && (
          <div className="text-center py-12">
            <div className="text-6xl mb-4">🎉</div>
            <h3 className="text-xl font-semibold mb-2">素晴らしい！</h3>
            <p className="text-gray-600">今日の復習はすべて完了しました</p>
            <a
              href="/books"
              className="inline-block mt-6 px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
            >
              新しいページを学習する
            </a>
          </div>
        )}
      </div>

      {/* Review Session */}
      {activeSession && (
        <ReviewSession
          items={activeSession}
          onComplete={handleCompleteSession}
          onCancel={handleCancelSession}
        />
      )}
    </div>
  );
}
```

## テスト方法

1. **開発サーバー起動**:
   ```bash
   cd frontend/web
   pnpm run dev
   ```

2. **ブラウザで確認**: http://localhost:3000/review

3. **確認項目**:
   - [ ] 復習統計が表示される
   - [ ] 緊急・推奨・余裕ありの3つのカードが表示される
   - [ ] 「復習する」ボタンをクリックするとセッションが開始される
   - [ ] フラッシュカード形式で復習できる
   - [ ] スコアを選択すると次の項目に進む
   - [ ] すべて完了すると統計が更新される

## 完了条件

- [ ] 型定義ファイルが作成されている
- [ ] API クライアントが拡張されている
- [ ] ReviewCard コンポーネントが動作する
- [ ] ReviewSession コンポーネントが動作する
- [ ] Review ページが正しくレンダリングされる
- [ ] 復習セッションが最後まで完了できる
- [ ] エラー処理が適切に実装されている

## トラブルシューティング

### APIエラーが発生する場合
- バックエンドが起動しているか確認: `curl http://localhost:8080/health`
- ブラウザの開発者ツールでNetwork タブを確認

### スタイルが崩れる場合
- Tailwind CSS が正しくビルドされているか確認
- `pnpm run dev` を再起動

### TypeScript エラーが発生する場合
- 型定義ファイルが正しくインポートされているか確認
- `pnpm run type-check` でエラーを確認

## 参考資料

- [間隔反復学習アルゴリズム](../../docs/featureRDs/8_間隔反復学習SRS.md)
- [UI/UX設計書](../../docs/ui_ux_design_document.md)
- [Next.js公式ドキュメント](https://nextjs.org/docs)
