'use client';

import { useState, useEffect } from 'react';
import { apiClient } from '@/lib/api/client';
import type { DashboardStats, LearningTimeData, ProgressData, WeakPointsData } from '@/types/stats';

export default function StatsPage() {
  const [dashboard, setDashboard] = useState<DashboardStats | null>(null);
  const [learningTime, setLearningTime] = useState<LearningTimeData | null>(null);
  const [progress, setProgress] = useState<ProgressData | null>(null);
  const [weakPoints, setWeakPoints] = useState<WeakPointsData | null>(null);
  const [period, setPeriod] = useState<'week' | 'month' | 'year'>('week');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        setLoading(true);
        setError(null);

        const [dashboardData, learningTimeData, progressData, weakPointsData] = await Promise.all([
          apiClient.stats.getDashboard(),
          apiClient.stats.getLearningTime(period),
          apiClient.stats.getProgress(period),
          apiClient.stats.getWeakPoints(10),
        ]);

        setDashboard(dashboardData);
        setLearningTime(learningTimeData);
        setProgress(progressData);
        setWeakPoints(weakPointsData);
      } catch (err) {
        console.error('Failed to fetch stats:', err);
        setError('統計データの取得に失敗しました');
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, [period]);

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/4 mb-8"></div>
          <div className="space-y-4">
            <div className="h-64 bg-gray-200 rounded"></div>
            <div className="h-48 bg-gray-200 rounded"></div>
            <div className="h-48 bg-gray-200 rounded"></div>
          </div>
        </div>
      </div>
    );
  }

  if (error || !dashboard || !learningTime || !progress) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
          <p className="text-red-800 mb-4">{error || 'データの読み込みに失敗しました'}</p>
          <button
            onClick={() => window.location.reload()}
            className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
          >
            再読み込み
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">学習統計</h1>
        <p className="text-gray-600">あなたの学習状況を確認しましょう</p>
      </div>

      {/* Period Selector */}
      <div className="mb-6 flex gap-2">
        {(['week', 'month', 'year'] as const).map((p) => (
          <button
            key={p}
            onClick={() => setPeriod(p)}
            className={`px-4 py-2 rounded-lg font-medium transition-colors ${
              period === p
                ? 'bg-blue-600 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            }`}
          >
            {p === 'week' ? '週' : p === 'month' ? '月' : '年'}
          </button>
        ))}
      </div>

      {/* Dashboard Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <StatsCard
          title="現在のストリーク"
          value={`${dashboard.currentStreak}日`}
          icon="🔥"
          color="orange"
        />
        <StatsCard
          title="最長ストリーク"
          value={`${dashboard.longestStreak}日`}
          icon="⭐"
          color="yellow"
        />
        <StatsCard
          title="完了ページ"
          value={`${dashboard.completedPages}`}
          subtitle={`全${dashboard.totalPages}ページ`}
          icon="📄"
          color="blue"
        />
        <StatsCard
          title="習得単語数"
          value={`${dashboard.masteredWords}`}
          icon="📚"
          color="green"
        />
      </div>

      {/* Learning Time Chart */}
      <div className="bg-white rounded-lg shadow-md p-6 mb-8">
        <h2 className="text-xl font-bold mb-4">今週の学習時間</h2>
        <div className="flex items-end justify-between h-64 gap-2">
          {learningTime.data.length > 0 ? (
            learningTime.data.map((item, index) => {
              const maxMinutes = Math.max(...learningTime.data.map((d) => d.minutes), 1);
              const heightPercent = (item.minutes / maxMinutes) * 100;

              return (
                <div key={index} className="flex-1 flex flex-col items-center">
                  <div
                    className="w-full bg-blue-600 rounded-t transition-all hover:bg-blue-700"
                    style={{ height: `${heightPercent}%`, minHeight: item.minutes > 0 ? '4px' : '0' }}
                    title={`${item.minutes}分`}
                  />
                  <div className="text-xs text-gray-600 mt-2">{item.date.split('-')[2]}</div>
                  <div className="text-xs text-gray-500">{item.minutes}分</div>
                </div>
              );
            })
          ) : (
            <div className="w-full flex items-center justify-center h-full text-gray-400">
              データがありません
            </div>
          )}
        </div>
        <div className="mt-4 text-sm text-gray-600">
          <p>総学習時間: {learningTime.totalMinutes}分</p>
          <p>平均学習時間: {learningTime.averageMinutes.toFixed(1)}分/日</p>
        </div>
      </div>

      {/* Progress Overview */}
      <div className="bg-white rounded-lg shadow-md p-6 mb-8">
        <h2 className="text-xl font-bold mb-4">進捗状況</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <ProgressItem
            title="単語"
            data={progress.words}
            color="blue"
          />
          <ProgressItem
            title="フレーズ"
            data={progress.phrases}
            color="green"
          />
          <ProgressItem
            title="ページ"
            data={progress.pages}
            color="purple"
          />
        </div>
      </div>

      {/* Weak Points */}
      {weakPoints && (weakPoints.weakWords.length > 0 || weakPoints.weakPhrases.length > 0) && (
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-bold mb-4">苦手な項目</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {weakPoints.weakWords.length > 0 && (
              <div>
                <h3 className="font-semibold mb-3 text-gray-700">単語</h3>
                <div className="space-y-2">
                  {weakPoints.weakWords.map((item, index) => (
                    <WeakPointItem
                      key={index}
                      text={item.word || ''}
                      language={item.language}
                      attempts={item.attempts}
                      averageScore={item.averageScore}
                    />
                  ))}
                </div>
              </div>
            )}
            {weakPoints.weakPhrases.length > 0 && (
              <div>
                <h3 className="font-semibold mb-3 text-gray-700">フレーズ</h3>
                <div className="space-y-2">
                  {weakPoints.weakPhrases.map((item, index) => (
                    <WeakPointItem
                      key={index}
                      text={item.phrase || ''}
                      language={item.language}
                      attempts={item.attempts}
                      averageScore={item.averageScore}
                    />
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

// Helper Components

interface StatsCardProps {
  title: string;
  value: string;
  subtitle?: string;
  icon: string;
  color: 'blue' | 'green' | 'orange' | 'yellow' | 'purple';
}

function StatsCard({ title, value, subtitle, icon, color }: StatsCardProps) {
  const colorClasses = {
    blue: 'bg-blue-50 border-blue-200 text-blue-600',
    green: 'bg-green-50 border-green-200 text-green-600',
    orange: 'bg-orange-50 border-orange-200 text-orange-600',
    yellow: 'bg-yellow-50 border-yellow-200 text-yellow-600',
    purple: 'bg-purple-50 border-purple-200 text-purple-600',
  };

  return (
    <div className={`rounded-lg border-2 p-4 ${colorClasses[color]}`}>
      <div className="flex items-center justify-between mb-2">
        <span className="text-2xl">{icon}</span>
        <h3 className="text-sm font-medium text-gray-600">{title}</h3>
      </div>
      <p className="text-2xl font-bold mb-1">{value}</p>
      {subtitle && <p className="text-xs text-gray-500">{subtitle}</p>}
    </div>
  );
}

interface ProgressItemProps {
  title: string;
  data: Array<{ date: string; count: number }>;
  color: 'blue' | 'green' | 'purple';
}

function ProgressItem({ title, data, color }: ProgressItemProps) {
  const colorClasses = {
    blue: 'text-blue-600',
    green: 'text-green-600',
    purple: 'text-purple-600',
  };

  const total = data.reduce((sum, item) => sum + item.count, 0);
  const latest = data.length > 0 ? data[data.length - 1].count : 0;

  return (
    <div>
      <h3 className={`font-semibold mb-2 ${colorClasses[color]}`}>{title}</h3>
      <div className="space-y-2">
        <div className="flex justify-between items-center">
          <span className="text-sm text-gray-600">今週の合計:</span>
          <span className="font-bold text-lg">{total}</span>
        </div>
        <div className="flex justify-between items-center">
          <span className="text-sm text-gray-600">最新:</span>
          <span className="font-semibold">{latest}</span>
        </div>
      </div>
    </div>
  );
}

interface WeakPointItemProps {
  text: string;
  language: string;
  attempts: number;
  averageScore: number;
}

function WeakPointItem({ text, language, attempts, averageScore }: WeakPointItemProps) {
  const scoreColor = averageScore >= 70 ? 'text-yellow-600' : 'text-red-600';

  return (
    <div className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
      <div className="flex-1">
        <p className="font-medium text-gray-900">{text}</p>
        <p className="text-xs text-gray-500">{language}</p>
      </div>
      <div className="text-right">
        <p className={`font-bold ${scoreColor}`}>{averageScore.toFixed(0)}点</p>
        <p className="text-xs text-gray-500">{attempts}回</p>
      </div>
    </div>
  );
}
