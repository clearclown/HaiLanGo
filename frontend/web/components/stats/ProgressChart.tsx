'use client';

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import type { ProgressDataPoint } from '@/lib/types/stats';

interface ProgressChartProps {
  data: ProgressDataPoint[];
}

export function ProgressChart({ data }: ProgressChartProps) {
  // Format data for recharts
  const chartData = data.map((point) => ({
    date: new Date(point.date).toLocaleDateString('ja-JP', {
      month: 'short',
      day: 'numeric',
    }),
    words: point.words,
    phrases: point.phrases,
    pages: point.pages,
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
      <LineChart data={chartData}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="date" />
        <YAxis />
        <Tooltip labelStyle={{ color: '#333' }} />
        <Legend />
        <Line
          type="monotone"
          dataKey="words"
          stroke="#9333EA"
          name="単語"
          strokeWidth={2}
        />
        <Line
          type="monotone"
          dataKey="phrases"
          stroke="#10B981"
          name="フレーズ"
          strokeWidth={2}
        />
        <Line
          type="monotone"
          dataKey="pages"
          stroke="#3B82F6"
          name="ページ"
          strokeWidth={2}
        />
      </LineChart>
    </ResponsiveContainer>
  );
}
