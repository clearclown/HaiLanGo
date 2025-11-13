'use client';

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import type { LearningTimeDataPoint } from '@/lib/types/stats';

interface LearningTimeChartProps {
  data: LearningTimeDataPoint[];
}

export function LearningTimeChart({ data }: LearningTimeChartProps) {
  // Format data for recharts
  const chartData = data.map((point) => ({
    date: new Date(point.date).toLocaleDateString('ja-JP', {
      month: 'short',
      day: 'numeric',
    }),
    minutes: Math.round(point.seconds / 60),
  }));

  if (data.length === 0) {
    return (
      <div className="flex items-center justify-center h-64 text-gray-500">
        データがありません
      </div>
    );
  }

  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={chartData}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="date" />
        <YAxis label={{ value: '分', angle: -90, position: 'insideLeft' }} />
        <Tooltip
          formatter={(value: number) => [`${value}分`, '学習時間']}
          labelStyle={{ color: '#333' }}
        />
        <Legend />
        <Bar dataKey="minutes" fill="#4A90E2" name="学習時間（分）" />
      </BarChart>
    </ResponsiveContainer>
  );
}
