'use client';

import React, { createContext, useContext, useEffect } from 'react';
import { useWebSocket } from '@/hooks/useWebSocket';
import type { MessageType, MessageHandler } from '@/lib/websocket';

interface WebSocketContextValue {
  connected: boolean;
  subscribe: <T = unknown>(type: MessageType, handler: MessageHandler<T>) => () => void;
}

const WebSocketContext = createContext<WebSocketContextValue | null>(null);

interface WebSocketProviderProps {
  children: React.ReactNode;
  autoConnect?: boolean;
}

/**
 * WebSocket接続を提供するプロバイダーコンポーネント
 * アプリケーションのルートレベルで使用する
 */
export function WebSocketProvider({ children, autoConnect = true }: WebSocketProviderProps) {
  const { connected, connect, subscribe } = useWebSocket({ autoConnect });

  useEffect(() => {
    if (autoConnect && typeof window !== 'undefined') {
      // トークンを取得（実際の実装に応じて調整）
      const token = localStorage.getItem('auth_token');
      if (token) {
        connect(token);
      }
    }
  }, [autoConnect, connect]);

  return (
    <WebSocketContext.Provider value={{ connected, subscribe }}>
      {children}
    </WebSocketContext.Provider>
  );
}

/**
 * WebSocketコンテキストを使用するフック
 */
export function useWebSocketContext(): WebSocketContextValue {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error('useWebSocketContext must be used within WebSocketProvider');
  }
  return context;
}
