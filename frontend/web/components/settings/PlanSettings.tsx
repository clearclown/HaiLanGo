'use client';

import type { Plan } from '@/types/settings';

interface PlanSettingsProps {
  plan: Plan;
  onUpgrade: () => void;
}

export default function PlanSettings({ plan, onUpgrade }: PlanSettingsProps) {
  const formatDate = (dateString?: string) => {
    if (!dateString) return '';
    return new Date(dateString).toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <h2 className="text-xl font-semibold mb-4">プラン</h2>

      <div className="mb-4">
        <div className="text-lg font-semibold">
          {plan.type === 'free' ? '無料プラン' : 'プレミアムプラン'}
        </div>
        {plan.expiresAt && (
          <div className="text-sm text-text-secondary mt-2">
            有効期限: {formatDate(plan.expiresAt)}
          </div>
        )}
      </div>

      {plan.type === 'free' && (
        <button
          type="button"
          onClick={onUpgrade}
          className="bg-secondary text-white px-6 py-2 rounded-lg hover:bg-opacity-90"
        >
          プレミアムにアップグレード
        </button>
      )}

      {plan.type === 'premium' && (
        <div className="text-sm text-text-secondary">プレミアム機能をご利用いただけます</div>
      )}
    </div>
  );
}
