'use client';

import { AppLayout } from '@/components/layout';
import { ReviewCard } from '@/components/review/ReviewCard';
import { ReviewSession } from '@/components/review/ReviewSession';
import { Button, Card, CardContent, CardHeader, CardTitle, Progress } from '@/components/ui';
import { apiClient } from '@/lib/api/client';
import type { ReviewItem, ReviewStats } from '@/types/review';
import { useEffect, useState } from 'react';

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
      <AppLayout>
        <div className="container-app py-6 lg:py-8 flex items-center justify-center min-h-[60vh]">
          <div className="text-center">
            <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-primary mx-auto mb-4" />
            <p className="text-gray-600">読み込み中...</p>
          </div>
        </div>
      </AppLayout>
    );
  }

  if (error) {
    return (
      <AppLayout>
        <div className="container-app py-6 lg:py-8 flex items-center justify-center min-h-[60vh]">
          <div className="text-center">
            <p className="text-error mb-4">{error}</p>
            <Button variant="primary" onClick={loadReviewData}>
              再試行
            </Button>
          </div>
        </div>
      </AppLayout>
    );
  }

  return (
    <AppLayout>
      <div className="container-app py-6 lg:py-8">
        {/* ヘッダー */}
        <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 mb-6">
          <div>
            <h1 className="text-2xl lg:text-3xl font-bold text-gray-900">復習</h1>
            <p className="text-gray-600 mt-1">間隔反復学習で効率的に記憶</p>
          </div>
          {stats && (
            <div className="text-left sm:text-right">
              <p className="text-sm text-gray-500">今日の復習</p>
              <p className="text-2xl font-bold text-secondary">{stats.total_completed_today}項目</p>
            </div>
          )}
        </div>

        {/* 統計情報 */}
        {stats && (
          <Card className="mb-6">
            <CardHeader>
              <CardTitle className="text-lg">今週の進捗</CardTitle>
            </CardHeader>
            <CardContent>
              <Progress value={stats.weekly_completion_rate} showLabel color="primary" size="md" />
            </CardContent>
          </Card>
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
    </AppLayout>
  );
}
