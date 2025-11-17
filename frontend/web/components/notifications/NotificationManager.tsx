'use client';

import { useToast } from './ToastContainer';
import {
  useOCRProgress,
  useBookReady,
  useNotification,
  useReviewReminder,
  useLearningUpdate,
} from '@/hooks/useWebSocket';
import type {
  OCRProgressPayload,
  BookReadyPayload,
  NotificationPayload,
  ReviewReminderPayload,
  LearningUpdatePayload,
} from '@/lib/websocket/types';

/**
 * NotificationManager
 *
 * WebSocketメッセージをリッスンし、適切なトースト通知を表示するコンポーネント
 * アプリケーションのルートレベルで使用する
 */
export const NotificationManager = () => {
  const toast = useToast();

  // OCR進捗通知
  useOCRProgress((payload: OCRProgressPayload) => {
    const progress = Math.round(payload.progress);
    toast.showInfo(
      'OCR処理中',
      `${payload.processedPages}/${payload.totalPages} ページ処理済み (${progress}%)`
    );
  });

  // 書籍準備完了通知
  useBookReady((payload: BookReadyPayload) => {
    toast.showSuccess(
      '書籍の準備が完了しました！',
      payload.message || `${payload.title} が学習可能になりました`
    );
  });

  // 一般通知
  useNotification((payload: NotificationPayload) => {
    const level = payload.level || 'info';

    switch (level) {
      case 'error':
        toast.showError(payload.title, payload.message);
        break;
      case 'warning':
        toast.showWarning(payload.title, payload.message);
        break;
      case 'success':
        toast.showSuccess(payload.title, payload.message);
        break;
      default:
        toast.showInfo(payload.title, payload.message);
    }
  });

  // 復習リマインダー通知
  useReviewReminder((payload: ReviewReminderPayload) => {
    toast.showInfo(
      '復習の時間です！',
      `${payload.count}個の復習項目があります`
    );
  });

  // 学習進捗更新通知
  useLearningUpdate((payload: LearningUpdatePayload) => {
    toast.showSuccess(
      '学習進捗を更新しました',
      payload.message || '統計情報が更新されました'
    );
  });

  return null; // このコンポーネントはUIを持たない
};
