# 実装指示書: 学習統計ダッシュボード

## 概要
ユーザーの学習進捗を可視化するダッシュボード機能の実装。学習時間、習得単語数、連続学習日数などの統計情報をグラフとカードで表示する。

## 担当範囲
- **フロントエンド**: `frontend/web/app/stats/page.tsx`
- **コンポーネント**: `frontend/web/components/stats/*` (一部実装済み)
- **バックエンドAPI**: すでに実装済み（`/api/v1/stats/*`）

## 前提条件
- Node.js 18+、pnpm がインストール済み
- バックエンドAPI が http://localhost:8080 で起動中
- 既存の stats コンポーネント（Dashboard, LearningTimeChart, ProgressChart）を活用

## 実装ステップ

### Step 1: 型定義の作成

**ファイル**: `frontend/web/types/stats.ts`

```typescript
export interface LearningTimeStats {
  total_seconds: number;
  total_hours: number;
  daily_average: number;
  weekly_average: number;
  monthly_average: number;
}

export interface ProgressStats {
  completed_pages: number;
  mastered_words: number;
  mastered_phrases: number;
  completed_books: number;
}

export interface StreakStats {
  current_streak: number;
  longest_streak: number;
  last_study_date: string;
}

export interface LearningTimeDataPoint {
  date: string;
  seconds: number;
}

export interface ProgressDataPoint {
  date: string;
  words: number;
  phrases: number;
  pages: number;
}

export interface DashboardStats {
  learning_time: LearningTimeStats;
  progress: ProgressStats;
  streak: StreakStats;
  pronunciation_avg: number;
  weak_words: string[];
  learning_time_chart: LearningTimeDataPoint[];
  progress_chart: ProgressDataPoint[];
}
```

### Step 2: API クライアントの拡張

**ファイル**: `frontend/web/lib/api/client.ts`

**追加する内容**:

```typescript
import type { DashboardStats } from '@/types/stats';

stats = {
  getDashboard: async (): Promise<DashboardStats> => {
    return this.fetch<DashboardStats>('/api/v1/stats/dashboard');
  },

  getLearningTime: async (days: number = 7): Promise<LearningTimeDataPoint[]> => {
    return this.fetch<LearningTimeDataPoint[]>(`/api/v1/stats/learning-time?days=${days}`);
  },

  getProgress: async (days: number = 30): Promise<ProgressDataPoint[]> => {
    return this.fetch<ProgressDataPoint[]>(`/api/v1/stats/progress?days=${days}`);
  },
};
```

### Step 3: StatsCard コンポーネントの作成

**ファイル**: `frontend/web/components/stats/StatsCard.tsx`

```typescript
interface StatsCardProps {
  icon: string;
  title: string;
  value: string | number;
  subtitle?: string;
  color?: 'blue' | 'green' | 'yellow' | 'red';
}

export function StatsCard({
  icon,
  title,
  value,
  subtitle,
  color = 'blue',
}: StatsCardProps) {
  const colorClasses = {
    blue: 'text-blue-600 bg-blue-50',
    green: 'text-green-600 bg-green-50',
    yellow: 'text-yellow-600 bg-yellow-50',
    red: 'text-red-600 bg-red-50',
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <div className="flex items-center gap-4">
        <div className={`text-3xl p-3 rounded-lg ${colorClasses[color]}`}>
          {icon}
        </div>
        <div className="flex-1">
          <h3 className="text-sm text-gray-600 mb-1">{title}</h3>
          <p className="text-2xl font-bold">{value}</p>
          {subtitle && (
            <p className="text-xs text-gray-500 mt-1">{subtitle}</p>
          )}
        </div>
      </div>
    </div>
  );
}
```

### Step 4: WeakWordsList コンポーネントの作成

**ファイル**: `frontend/web/components/stats/WeakWordsList.tsx`

```typescript
interface WeakWordsListProps {
  words: string[];
}

export function WeakWordsList({ words }: WeakWordsListProps) {
  if (words.length === 0) {
    return (
      <div className="bg-white rounded-lg shadow-sm p-6">
        <h2 className="text-xl font-semibold mb-4">🎯 苦手な単語</h2>
        <p className="text-gray-600 text-center py-8">
          素晴らしい！苦手な単語はありません
        </p>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <h2 className="text-xl font-semibold mb-4">🎯 苦手な単語</h2>
      <p className="text-sm text-gray-600 mb-4">
        習熟度が低い順に表示しています
      </p>
      <div className="space-y-2">
        {words.map((word, index) => (
          <div
            key={index}
            className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
          >
            <span className="font-mono text-gray-500">{index + 1}</span>
            <span className="font-medium">{word}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
```

### Step 5: Stats ページの実装

**ファイル**: `frontend/web/app/stats/page.tsx`

```typescript
'use client';

import { useEffect, useState } from 'react';
import { apiClient } from '@/lib/api/client';
import type { DashboardStats } from '@/types/stats';
import { StatsCard } from '@/components/stats/StatsCard';
import { LearningTimeChart } from '@/components/stats/LearningTimeChart';
import { ProgressChart } from '@/components/stats/ProgressChart';
import { WeakWordsList } from '@/components/stats/WeakWordsList';

export default function StatsPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [timeRange, setTimeRange] = useState<7 | 30 | 90>(7);

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const data = await apiClient.stats.getDashboard();
      setStats(data);
    } catch (err) {
      console.error('Failed to load stats:', err);
      setError('統計情報の読み込みに失敗しました');
    } finally {
      setIsLoading(false);
    }
  };

  const formatTime = (seconds: number): string => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    if (hours > 0) {
      return `${hours}時間${minutes}分`;
    }
    return `${minutes}分`;
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background-secondary flex items-center justify-center">
        <div className="text-gray-600">読み込み中...</div>
      </div>
    );
  }

  if (error || !stats) {
    return (
      <div className="min-h-screen bg-background-secondary flex flex-col items-center justify-center">
        <div className="text-red-600 mb-4">{error || 'データがありません'}</div>
        <button
          type="button"
          onClick={loadStats}
          className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
        >
          再試行
        </button>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background-secondary">
      <div className="max-w-6xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="flex items-center gap-3 mb-8">
          <h1 className="text-3xl font-bold">📊 学習統計</h1>
          {stats.streak.current_streak > 0 && (
            <span className="text-2xl">🔥</span>
          )}
        </div>

        {/* Quick Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <StatsCard
            icon="⏱️"
            title="総学習時間"
            value={`${Math.round(stats.learning_time.total_hours)}時間`}
            subtitle={`1日平均: ${formatTime(stats.learning_time.daily_average)}`}
            color="blue"
          />

          <StatsCard
            icon="🔥"
            title="連続学習"
            value={`${stats.streak.current_streak}日`}
            subtitle={`最長: ${stats.streak.longest_streak}日`}
            color="red"
          />

          <StatsCard
            icon="📚"
            title="完了ページ"
            value={stats.progress.completed_pages}
            subtitle={`${stats.progress.completed_books}冊完了`}
            color="green"
          />

          <StatsCard
            icon="✨"
            title="習得単語"
            value={stats.progress.mastered_words}
            subtitle={`フレーズ: ${stats.progress.mastered_phrases}個`}
            color="yellow"
          />
        </div>

        {/* Pronunciation Score */}
        {stats.pronunciation_avg > 0 && (
          <div className="bg-white rounded-lg shadow-sm p-6 mb-8">
            <h2 className="text-xl font-semibold mb-4">🎤 発音スコア</h2>
            <div className="flex items-center gap-4">
              <div className="flex-1">
                <div className="flex justify-between text-sm text-gray-600 mb-2">
                  <span>平均スコア</span>
                  <span>{Math.round(stats.pronunciation_avg)}点</span>
                </div>
                <div className="h-4 bg-gray-200 rounded-full overflow-hidden">
                  <div
                    className="h-full bg-green-500"
                    style={{ width: `${stats.pronunciation_avg}%` }}
                  />
                </div>
              </div>
              <div className="text-4xl font-bold text-green-600">
                {Math.round(stats.pronunciation_avg)}
              </div>
            </div>
          </div>
        )}

        {/* Time Range Selector */}
        <div className="flex justify-end mb-4">
          <div className="bg-white rounded-lg shadow-sm p-1 inline-flex">
            {[7, 30, 90].map((days) => (
              <button
                key={days}
                type="button"
                onClick={() => setTimeRange(days as 7 | 30 | 90)}
                className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                  timeRange === days
                    ? 'bg-blue-500 text-white'
                    : 'text-gray-600 hover:bg-gray-100'
                }`}
              >
                {days}日
              </button>
            ))}
          </div>
        </div>

        {/* Charts */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {/* Learning Time Chart */}
          <div className="bg-white rounded-lg shadow-sm p-6">
            <h2 className="text-xl font-semibold mb-4">📈 学習時間の推移</h2>
            <LearningTimeChart
              data={stats.learning_time_chart.filter((_, i) =>
                i >= stats.learning_time_chart.length - timeRange
              )}
            />
          </div>

          {/* Progress Chart */}
          <div className="bg-white rounded-lg shadow-sm p-6">
            <h2 className="text-xl font-semibold mb-4">📊 習得の推移</h2>
            <ProgressChart
              data={stats.progress_chart.filter((_, i) =>
                i >= stats.progress_chart.length - timeRange
              )}
            />
          </div>
        </div>

        {/* Weak Words */}
        <WeakWordsList words={stats.weak_words} />

        {/* Study Insights */}
        <div className="bg-white rounded-lg shadow-sm p-6 mt-8">
          <h2 className="text-xl font-semibold mb-4">💡 学習のヒント</h2>
          <div className="space-y-3">
            {stats.streak.current_streak === 0 && (
              <div className="p-4 bg-blue-50 rounded-lg">
                <p className="text-blue-800">
                  💪 今日から学習を始めて、連続記録を作りましょう！
                </p>
              </div>
            )}

            {stats.learning_time.daily_average < 600 && (
              <div className="p-4 bg-yellow-50 rounded-lg">
                <p className="text-yellow-800">
                  ⏰ 1日10分以上の学習を目標にしましょう。継続が大切です！
                </p>
              </div>
            )}

            {stats.weak_words.length > 0 && (
              <div className="p-4 bg-green-50 rounded-lg">
                <p className="text-green-800">
                  🎯 苦手な単語を復習して、習熟度を上げましょう！
                </p>
              </div>
            )}

            {stats.pronunciation_avg > 0 && stats.pronunciation_avg < 70 && (
              <div className="p-4 bg-purple-50 rounded-lg">
                <p className="text-purple-800">
                  🎤 発音練習を増やして、スコアアップを目指しましょう！
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
```

## テスト方法

1. ブラウザで http://localhost:3000/stats にアクセス

2. **確認項目**:
   - [ ] 統計カードが4つ表示される
   - [ ] 学習時間グラフが表示される
   - [ ] 進捗グラフが表示される
   - [ ] 苦手な単語リストが表示される
   - [ ] 時間範囲（7日/30日/90日）を切り替えられる
   - [ ] 学習のヒントが表示される

## 完了条件

- [ ] 型定義が作成されている
- [ ] API クライアントが拡張されている
- [ ] StatsCard コンポーネントが動作する
- [ ] WeakWordsList コンポーネントが動作する
- [ ] Stats ページが正しくレンダリングされる
- [ ] グラフが正しく表示される
- [ ] エラーハンドリングが適切に実装されている

## 参考資料

- [学習統計ダッシュボードRD](../../docs/featureRDs/10_学習統計ダッシュボード.md)
- [UI/UX設計書](../../docs/ui_ux_design_document.md)
