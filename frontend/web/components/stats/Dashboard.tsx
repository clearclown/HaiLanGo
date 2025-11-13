'use client';

import { useQuery } from '@tanstack/react-query';
import { statsApi } from '@/lib/api/stats';
import { LearningTimeChart } from './LearningTimeChart';
import { ProgressChart } from './ProgressChart';

export function Dashboard() {
  const { data: dashboard, isLoading, error } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: statsApi.getDashboard,
    refetchInterval: 60000, // Refetch every minute
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-lg">読み込み中...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-lg text-red-500">
          エラーが発生しました: {error.message}
        </div>
      </div>
    );
  }

  if (!dashboard) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-lg">データがありません</div>
      </div>
    );
  }

  const { learning_time, progress, streak, pronunciation_avg, weak_words, learning_time_chart, progress_chart } = dashboard;

  // Convert seconds to hours and minutes
  const hours = Math.floor(learning_time.total_seconds / 3600);
  const minutes = Math.floor((learning_time.total_seconds % 3600) / 60);

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">学習統計ダッシュボード</h1>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {/* Learning Time */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-gray-500 text-sm font-medium mb-2">総学習時間</h3>
          <p className="text-3xl font-bold text-blue-600">
            {hours}時間{minutes}分
          </p>
          <p className="text-sm text-gray-600 mt-2">
            1日平均: {Math.floor(learning_time.daily_average / 60)}分
          </p>
        </div>

        {/* Streak */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-gray-500 text-sm font-medium mb-2">連続学習記録</h3>
          <p className="text-3xl font-bold text-orange-600 flex items-center">
            {streak.current_streak}日
            <span className="ml-2 text-2xl">🔥</span>
          </p>
          <p className="text-sm text-gray-600 mt-2">
            最長記録: {streak.longest_streak}日
          </p>
        </div>

        {/* Progress */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-gray-500 text-sm font-medium mb-2">学習進捗</h3>
          <p className="text-3xl font-bold text-green-600">
            {progress.completed_pages}
          </p>
          <p className="text-sm text-gray-600 mt-2">
            完了ページ数
          </p>
        </div>

        {/* Mastered Words */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-gray-500 text-sm font-medium mb-2">習得単語数</h3>
          <p className="text-3xl font-bold text-purple-600">
            {progress.mastered_words}語
          </p>
          <p className="text-sm text-gray-600 mt-2">
            フレーズ: {progress.mastered_phrases}個
          </p>
        </div>
      </div>

      {/* Additional Stats */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        {/* Pronunciation Average */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold mb-4">発音スコア平均</h3>
          <div className="flex items-center">
            <div className="w-full bg-gray-200 rounded-full h-4">
              <div
                className="bg-blue-600 h-4 rounded-full transition-all"
                style={{ width: `${pronunciation_avg}%` }}
              />
            </div>
            <span className="ml-4 text-2xl font-bold text-blue-600">
              {pronunciation_avg.toFixed(1)}
            </span>
          </div>
        </div>

        {/* Completed Books */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold mb-4">完了書籍数</h3>
          <p className="text-4xl font-bold text-indigo-600">
            {progress.completed_books}冊
          </p>
        </div>
      </div>

      {/* Weak Words */}
      {weak_words.length > 0 && (
        <div className="bg-white rounded-lg shadow p-6 mb-8">
          <h3 className="text-lg font-semibold mb-4">苦手な単語</h3>
          <div className="flex flex-wrap gap-2">
            {weak_words.map((word, index) => (
              <span
                key={index}
                className="px-3 py-1 bg-red-100 text-red-700 rounded-full text-sm"
              >
                {word}
              </span>
            ))}
          </div>
        </div>
      )}

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Learning Time Chart */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold mb-4">学習時間の推移（過去7日間）</h3>
          <LearningTimeChart data={learning_time_chart} />
        </div>

        {/* Progress Chart */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold mb-4">学習進捗の推移（過去30日間）</h3>
          <ProgressChart data={progress_chart} />
        </div>
      </div>
    </div>
  );
}
