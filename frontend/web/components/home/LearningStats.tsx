import type { LearningStats as LearningStatsType } from '@/lib/types';
import { cn } from '@/lib/utils';

interface LearningStatsProps {
  stats: LearningStatsType;
}

function formatLearningTime(seconds: number): string {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);

  if (hours === 0) {
    return `${minutes}分`;
  }

  return `${hours}時間${minutes}分`;
}

export function LearningStats({ stats }: LearningStatsProps) {
  return (
    <div className={cn('rounded-xl border border-border bg-white p-6 shadow-sm')}>
      <div className="mb-4 flex items-center gap-2">
        <span className="text-2xl">📊</span>
        <h2 className="text-xl font-semibold text-text-primary">学習統計</h2>
      </div>

      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <span className="text-text-secondary">連続学習</span>
          <span className="font-semibold text-text-primary">
            {stats.streakDays}日 {stats.streakDays > 0 && '🔥'}
          </span>
        </div>

        <div className="flex items-center justify-between">
          <span className="text-text-secondary">総学習時間</span>
          <span className="font-semibold text-text-primary">
            {formatLearningTime(stats.totalLearningTimeSeconds)}
          </span>
        </div>
      </div>
    </div>
  );
}
