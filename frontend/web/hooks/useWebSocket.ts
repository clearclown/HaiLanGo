'use client';

import { useEffect, useRef, useState, useCallback } from 'react';
import {
  getWebSocketClient,
  type MessageType,
  type MessageHandler,
  type WebSocketConfig,
  type OCRProgressPayload,
  type BookReadyPayload,
  type ReviewReminderPayload,
  type LearningUpdatePayload,
  type NotificationPayload,
  type ErrorPayload,
} from '@/lib/websocket';

interface UseWebSocketOptions {
  config?: WebSocketConfig;
  autoConnect?: boolean;
}

interface UseWebSocketReturn {
  connected: boolean;
  connect: (token: string) => void;
  disconnect: () => void;
  subscribe: <T = unknown>(type: MessageType, handler: MessageHandler<T>) => () => void;
}

/**
 * WebSocket接続を管理するReactフック
 * @param options WebSocketオプション
 * @returns WebSocket接続の状態と制御関数
 */
export function useWebSocket(options: UseWebSocketOptions = {}): UseWebSocketReturn {
  const { config, autoConnect = false } = options;
  const [connected, setConnected] = useState(false);
  const clientRef = useRef(getWebSocketClient(config));
  const checkIntervalRef = useRef<NodeJS.Timeout | null>(null);
  const subscribersRef = useRef<Map<MessageType, Set<MessageHandler>>>(new Map());

  // 接続状態を定期的にチェック
  const startConnectionCheck = useCallback(() => {
    if (checkIntervalRef.current) {
      clearInterval(checkIntervalRef.current);
    }

    checkIntervalRef.current = setInterval(() => {
      const isConnected = clientRef.current.isConnected();
      setConnected(isConnected);
    }, 1000);
  }, []);

  const stopConnectionCheck = useCallback(() => {
    if (checkIntervalRef.current) {
      clearInterval(checkIntervalRef.current);
      checkIntervalRef.current = null;
    }
  }, []);

  // WebSocket接続
  const connect = useCallback((token: string) => {
    try {
      clientRef.current.connect(token);
      startConnectionCheck();
    } catch (error) {
      console.error('Failed to connect WebSocket:', error);
    }
  }, [startConnectionCheck]);

  // WebSocket切断
  const disconnect = useCallback(() => {
    clientRef.current.disconnect();
    stopConnectionCheck();
    setConnected(false);
  }, [stopConnectionCheck]);

  // メッセージハンドラーを購読
  const subscribe = useCallback(<T = unknown>(
    type: MessageType,
    handler: MessageHandler<T>
  ): (() => void) => {
    // 購読者を追跡
    if (!subscribersRef.current.has(type)) {
      subscribersRef.current.set(type, new Set());
    }
    subscribersRef.current.get(type)!.add(handler as MessageHandler);

    // WebSocketクライアントに登録
    clientRef.current.on(type, handler);

    // クリーンアップ関数を返す
    return () => {
      const subscribers = subscribersRef.current.get(type);
      if (subscribers) {
        subscribers.delete(handler as MessageHandler);
      }
      clientRef.current.off(type, handler);
    };
  }, []);

  // 自動接続
  useEffect(() => {
    if (autoConnect && typeof window !== 'undefined') {
      // トークンを取得（localStorage, cookie, etc.から）
      const token = localStorage.getItem('auth_token');
      if (token) {
        connect(token);
      }
    }

    return () => {
      // コンポーネントのアンマウント時にクリーンアップ
      stopConnectionCheck();

      // 全ての購読を解除
      subscribersRef.current.forEach((handlers, type) => {
        handlers.forEach(handler => {
          clientRef.current.off(type, handler);
        });
      });
      subscribersRef.current.clear();
    };
  }, [autoConnect, connect, stopConnectionCheck]);

  return {
    connected,
    connect,
    disconnect,
    subscribe,
  };
}

/**
 * 特定のメッセージタイプを購読する簡易フック
 * @param type メッセージタイプ
 * @param handler ハンドラー関数
 */
export function useWebSocketSubscription<T = unknown>(
  type: MessageType,
  handler: MessageHandler<T>
): void {
  const { subscribe } = useWebSocket();

  useEffect(() => {
    const unsubscribe = subscribe(type, handler);
    return unsubscribe;
  }, [type, handler, subscribe]);
}

/**
 * 複数のメッセージタイプを購読する簡易フック
 * @param subscriptions メッセージタイプとハンドラーのマップ
 */
export function useWebSocketSubscriptions(
  subscriptions: Record<MessageType, MessageHandler>
): void {
  const { subscribe } = useWebSocket();

  useEffect(() => {
    const unsubscribers = Object.entries(subscriptions).map(([type, handler]) =>
      subscribe(type as MessageType, handler)
    );

    return () => {
      unsubscribers.forEach(unsubscribe => unsubscribe());
    };
  }, [subscriptions, subscribe]);
}

// ========================================
// Typed hooks for specific message types
// ========================================

/**
 * OCR進捗通知を購読するフック
 */
export function useOCRProgress(handler: MessageHandler<OCRProgressPayload>): void {
  useWebSocketSubscription('ocr_progress', handler);
}

/**
 * 書籍準備完了通知を購読するフック
 */
export function useBookReady(handler: MessageHandler<BookReadyPayload>): void {
  useWebSocketSubscription('book_ready', handler);
}

/**
 * 復習リマインダー通知を購読するフック
 */
export function useReviewReminder(handler: MessageHandler<ReviewReminderPayload>): void {
  useWebSocketSubscription('review_reminder', handler);
}

/**
 * 学習更新通知を購読するフック
 */
export function useLearningUpdate(handler: MessageHandler<LearningUpdatePayload>): void {
  useWebSocketSubscription('learning_update', handler);
}

/**
 * 一般通知を購読するフック
 */
export function useNotification(handler: MessageHandler<NotificationPayload>): void {
  useWebSocketSubscription('notification', handler);
}

/**
 * エラー通知を購読するフック
 */
export function useErrorNotification(handler: MessageHandler<ErrorPayload>): void {
  useWebSocketSubscription('error', handler);
}
