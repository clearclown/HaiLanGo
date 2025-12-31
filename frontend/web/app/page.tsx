'use client';

import { LearningStats } from '@/components/home/LearningStats';
import { QuickAccess } from '@/components/home/QuickAccess';
import { TodayLearningCard } from '@/components/home/TodayLearningCard';
import { WelcomeCard } from '@/components/home/WelcomeCard';
import { AppLayout } from '@/components/layout';
import type { DashboardData } from '@/lib/types';
import { useEffect, useState } from 'react';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export default function HomePage() {
  const [data, setData] = useState<DashboardData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchData() {
      try {
        const token = localStorage.getItem('access_token');
        const response = await fetch(`${API_BASE_URL}/api/v1/home/dashboard`, {
          headers: {
            Authorization: token ? `Bearer ${token}` : '',
          },
        });

        if (!response.ok) {
          throw new Error('Failed to fetch dashboard data');
        }

        const dashboardData = await response.json();
        setData(dashboardData);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An error occurred');
      } finally {
        setLoading(false);
      }
    }

    fetchData();
  }, []);

  if (loading) {
    return (
      <AppLayout>
        <div className="container-app py-6 lg:py-8 flex items-center justify-center min-h-[50vh]">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
        </div>
      </AppLayout>
    );
  }

  if (error || !data) {
    return (
      <AppLayout>
        <div className="container-app py-6 lg:py-8">
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
            {error || 'Failed to load dashboard data'}
          </div>
        </div>
      </AppLayout>
    );
  }

  return (
    <AppLayout>
      <div className="container-app py-6 lg:py-8">
        <WelcomeCard userName={data.user.name} />

        <div className="space-y-6">
          {data.todayLearning && <TodayLearningCard data={data.todayLearning} />}

          <QuickAccess
            data={{
              booksCount: data.stats.booksCount,
              reviewItemsCount: data.stats.reviewItemsCount,
            }}
          />

          <LearningStats stats={data.stats} />
        </div>
      </div>
    </AppLayout>
  );
}
