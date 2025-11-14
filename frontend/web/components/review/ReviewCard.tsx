'use client';

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
          bgColor: 'bg-red-50',
          textColor: 'text-red-600',
          borderColor: 'border-red-200',
          buttonColor: 'bg-red-500 hover:bg-red-600',
          icon: '🔴',
          title: '緊急',
          description: '今日中に復習が必要',
        };
      case 'recommended':
        return {
          bgColor: 'bg-yellow-50',
          textColor: 'text-yellow-600',
          borderColor: 'border-yellow-200',
          buttonColor: 'bg-yellow-500 hover:bg-yellow-600',
          icon: '🟡',
          title: '推奨',
          description: '今日復習すると効果的',
        };
      case 'optional':
        return {
          bgColor: 'bg-green-50',
          textColor: 'text-green-600',
          borderColor: 'border-green-200',
          buttonColor: 'bg-green-500 hover:bg-green-600',
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
          className={`w-full mt-4 px-4 py-3 text-white rounded-lg transition-colors font-medium ${config.buttonColor}`}
        >
          復習する
        </button>
      )}
    </div>
  );
}
