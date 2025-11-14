'use client';

import { useState } from 'react';
import type { ReviewItem } from '@/types/review';
import { apiClient } from '@/lib/api/client';

interface ReviewSessionProps {
  items: ReviewItem[];
  onComplete: () => void;
  onCancel: () => void;
}

export function ReviewSession({ items, onComplete, onCancel }: ReviewSessionProps) {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [showTranslation, setShowTranslation] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const currentItem = items[currentIndex];
  const progress = ((currentIndex + 1) / items.length) * 100;

  const handleScore = async (score: number) => {
    setIsSubmitting(true);
    try {
      await apiClient.review.submit({
        item_id: currentItem.id,
        score,
        completed_at: new Date().toISOString(),
      });

      if (currentIndex < items.length - 1) {
        setCurrentIndex(currentIndex + 1);
        setShowTranslation(false);
      } else {
        onComplete();
      }
    } catch (error) {
      console.error('Failed to submit review:', error);
      alert('復習の送信に失敗しました');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-8 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold">復習セッション</h2>
          <button
            type="button"
            onClick={onCancel}
            className="text-gray-500 hover:text-gray-700"
          >
            ✕
          </button>
        </div>

        {/* Progress Bar */}
        <div className="mb-6">
          <div className="flex justify-between text-sm text-gray-600 mb-2">
            <span>{currentIndex + 1} / {items.length}</span>
            <span>{Math.round(progress)}%</span>
          </div>
          <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
            <div
              className="h-full bg-blue-500 transition-all duration-300"
              style={{ width: `${progress}%` }}
            />
          </div>
        </div>

        {/* Question */}
        <div className="mb-8">
          <div className="text-center mb-4">
            <span className="text-sm text-gray-500 uppercase">
              {currentItem.type === 'word' ? '単語' : 'フレーズ'}
            </span>
          </div>
          <div className="text-4xl font-bold text-center mb-4">
            {currentItem.text}
          </div>

          {/* Show Translation Button */}
          {!showTranslation && (
            <button
              type="button"
              onClick={() => setShowTranslation(true)}
              className="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200"
            >
              翻訳を表示
            </button>
          )}

          {/* Translation */}
          {showTranslation && (
            <div className="text-2xl text-center text-gray-700 mb-6">
              {currentItem.translation}
            </div>
          )}
        </div>

        {/* Score Buttons */}
        {showTranslation && (
          <div className="space-y-3">
            <h3 className="font-semibold text-center mb-4">どれくらい覚えていましたか？</h3>
            <button
              type="button"
              onClick={() => handleScore(100)}
              disabled={isSubmitting}
              className="w-full px-4 py-3 bg-green-500 text-white rounded-lg hover:bg-green-600 disabled:opacity-50"
            >
              🟢 完璧に覚えていた
            </button>
            <button
              type="button"
              onClick={() => handleScore(70)}
              disabled={isSubmitting}
              className="w-full px-4 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50"
            >
              🔵 少し時間がかかった
            </button>
            <button
              type="button"
              onClick={() => handleScore(30)}
              disabled={isSubmitting}
              className="w-full px-4 py-3 bg-red-500 text-white rounded-lg hover:bg-red-600 disabled:opacity-50"
            >
              🔴 思い出せなかった
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
