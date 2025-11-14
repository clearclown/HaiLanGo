'use client';

import { useState, useEffect } from 'react';
import { ReviewCard } from '@/components/review/ReviewCard';
import { ReviewSession } from '@/components/review/ReviewSession';
import type { ReviewItem, ReviewStats } from '@/types/review';
import { apiClient } from '@/lib/api/client';

export default function ReviewPage() {
  const [stats, setStats] = useState<ReviewStats | null>(null);
  const [urgentItems, setUrgentItems] = useState<ReviewItem[]>([]);
  const [recommendedItems, setRecommendedItems] = useState<ReviewItem[]>([]);
  const [optionalItems, setOptionalItems] = useState<ReviewItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeSession, setActiveSession] = useState<{
    items: ReviewItem[];
    priority: 'urgent' | 'recommended' | 'optional';
  } | null>(null);

  useEffect(() => {
    loadReviewData();
  }, []);

  const loadReviewData = async () => {
    try {
      setIsLoading(true);
      setError(null);

      // 統計情報を取得
      const statsData = await apiClient.review.getStats();
      setStats(statsData);

      // 各優先度のアイテムを取得
      const [urgent, recommended, optional] = await Promise.all([
        apiClient.review.getItems('urgent'),
        apiClient.review.getItems('recommended'),
        apiClient.review.getItems('optional'),
      ]);

      setUrgentItems(urgent.items);
      setRecommendedItems(recommended.items);
      setOptionalItems(optional.items);
    } catch (err) {
      console.error('Failed to load review data:', err);
      setError('復習データの読み込みに失敗しました');
    } finally {
      setIsLoading(false);
    }
  };

  const handleStartReview = (
    items: ReviewItem[],
    priority: 'urgent' | 'recommended' | 'optional'
  ) => {
    setActiveSession({ items, priority });
  };

  const handleCompleteSession = async () => {
    setActiveSession(null);
    // セッション完了後、データを再読み込み
    await loadReviewData();
  };

  const handleCancelSession = () => {
    setActiveSession(null);
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background-secondary flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-600">読み込み中...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-background-secondary flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-500 mb-4">{error}</p>
          <button
            type="button"
            onClick={loadReviewData}
            className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
          >
            再試行
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background-secondary">
      <div className="max-w-4xl mx-auto px-4 py-8">
        {/* ヘッダー */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold">復習</h1>
            <p className="text-gray-600 mt-1">間隔反復学習で効率的に記憶</p>
          </div>
          {stats && (
            <div className="text-right">
              <p className="text-sm text-gray-500">今日の復習</p>
              <p className="text-2xl font-bold text-green-500">
                {stats.total_completed_today}項目
              </p>
            </div>
          )}
        </div>

        {/* 統計情報 */}
        {stats && (
          <div className="bg-white rounded-lg p-6 mb-6 shadow-sm">
            <h2 className="text-lg font-semibold mb-4">今週の進捗</h2>
            <div className="flex items-center gap-4">
              <div className="flex-1">
                <div className="h-4 bg-gray-200 rounded-full overflow-hidden">
                  <div
                    className="h-full bg-blue-500 transition-all duration-300"
                    style={{ width: `${stats.weekly_completion_rate}%` }}
                  />
                </div>
              </div>
              <span className="text-lg font-semibold">
                {Math.round(stats.weekly_completion_rate)}%
              </span>
            </div>
          </div>
        )}

        {/* 復習カード */}
        <div className="space-y-4">
          <ReviewCard
            items={urgentItems}
            priority="urgent"
            onStartReview={() => handleStartReview(urgentItems, 'urgent')}
          />
          <ReviewCard
            items={recommendedItems}
            priority="recommended"
            onStartReview={() => handleStartReview(recommendedItems, 'recommended')}
          />
          <ReviewCard
            items={optionalItems}
            priority="optional"
            onStartReview={() => handleStartReview(optionalItems, 'optional')}
          />
        </div>

        {/* すべて完了した場合 */}
        {urgentItems.length === 0 &&
          recommendedItems.length === 0 &&
          optionalItems.length === 0 && (
            <div className="text-center py-12">
              <div className="text-6xl mb-4">🎉</div>
              <h2 className="text-2xl font-bold mb-2">素晴らしい！</h2>
              <p className="text-gray-600">
                今日の復習はすべて完了しました
                <br />
                新しいページを学習して、語彙を増やしましょう
              </p>
            </div>
          )}
      </div>

      {/* 復習セッション */}
      {activeSession && (
        <ReviewSession
          items={activeSession.items}
          onComplete={handleCompleteSession}
          onCancel={handleCancelSession}
        />
      )}
    </div>
  );
}
