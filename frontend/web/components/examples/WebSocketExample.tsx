'use client';

import { useCallback } from 'react';
import {
  useOCRProgress,
  useBookReady,
  useNotification,
} from '@/hooks/useWebSocket';
import type {
  OCRProgressPayload,
  BookReadyPayload,
  NotificationPayload,
} from '@/lib/websocket';

/**
 * WebSocket通知を受信する例のコンポーネント
 */
export function WebSocketExample() {
  // OCR進捗通知のハンドラー
  const handleOCRProgress = useCallback((payload: OCRProgressPayload) => {
    console.log('OCR Progress:', payload);
    // UI更新ロジックをここに実装
    // 例: トーストメッセージ表示、プログレスバー更新など
  }, []);

  // 書籍準備完了通知のハンドラー
  const handleBookReady = useCallback((payload: BookReadyPayload) => {
    console.log('Book Ready:', payload);
    // UI更新ロジックをここに実装
    // 例: 成功メッセージ表示、ページ遷移など
  }, []);

  // 一般通知のハンドラー
  const handleNotification = useCallback((payload: NotificationPayload) => {
    console.log('Notification:', payload);
    // UI更新ロジックをここに実装
    // 例: トーストメッセージ表示など
  }, []);

  // WebSocketメッセージを購読
  useOCRProgress(handleOCRProgress);
  useBookReady(handleBookReady);
  useNotification(handleNotification);

  return (
    <div className="p-4 border rounded-lg bg-gray-50">
      <h3 className="text-lg font-semibold mb-2">WebSocket通知</h3>
      <p className="text-sm text-gray-600">
        WebSocket経由でリアルタイム通知を受信しています。
      </p>
      <div className="mt-4 text-xs text-gray-500">
        <p>購読中のイベント:</p>
        <ul className="list-disc list-inside">
          <li>OCR進捗更新</li>
          <li>書籍準備完了</li>
          <li>一般通知</li>
        </ul>
      </div>
    </div>
  );
}
