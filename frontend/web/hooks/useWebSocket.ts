"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import type {
  ErrorData,
  LearningUpdateData,
  Notification,
  OCRProgressData,
  TTSProgressData,
} from "@/lib/types/notification";

export interface WebSocketOptions {
  url: string;
  userId: string;
  onOCRProgress?: (data: OCRProgressData) => void;
  onTTSProgress?: (data: TTSProgressData) => void;
  onLearningUpdate?: (data: LearningUpdateData) => void;
  onError?: (data: ErrorData) => void;
  reconnectAttempts?: number;
  reconnectInterval?: number;
}

export interface WebSocketState {
  isConnected: boolean;
  isConnecting: boolean;
  error: Error | null;
}

export function useWebSocket(options: WebSocketOptions) {
  const {
    url,
    userId,
    onOCRProgress,
    onTTSProgress,
    onLearningUpdate,
    onError,
    reconnectAttempts = 5,
    reconnectInterval = 3000,
  } = options;

  const [state, setState] = useState<WebSocketState>({
    isConnected: false,
    isConnecting: false,
    error: null,
  });

  const wsRef = useRef<WebSocket | null>(null);
  const reconnectCountRef = useRef(0);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      return;
    }

    setState((prev) => ({ ...prev, isConnecting: true, error: null }));

    try {
      const wsUrl = `${url}?user_id=${userId}`;
      const ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        console.log("WebSocket connected");
        setState({ isConnected: true, isConnecting: false, error: null });
        reconnectCountRef.current = 0;

        // Send initial ping
        ws.send(
          JSON.stringify({
            type: "ping",
            data: null,
            timestamp: new Date().toISOString(),
          }),
        );
      };

      ws.onmessage = (event) => {
        try {
          const notification = JSON.parse(event.data) as Notification;

          switch (notification.type) {
            case "ocr_progress":
              onOCRProgress?.(notification.data as OCRProgressData);
              break;
            case "tts_progress":
              onTTSProgress?.(notification.data as TTSProgressData);
              break;
            case "learning_update":
              onLearningUpdate?.(notification.data as LearningUpdateData);
              break;
            case "error":
              onError?.(notification.data as ErrorData);
              break;
            case "pong":
              // Handle pong response
              break;
          }
        } catch (error) {
          console.error("Failed to parse WebSocket message:", error);
        }
      };

      ws.onerror = (error) => {
        console.error("WebSocket error:", error);
        setState((prev) => ({
          ...prev,
          error: new Error("WebSocket connection error"),
        }));
      };

      ws.onclose = () => {
        console.log("WebSocket disconnected");
        setState({ isConnected: false, isConnecting: false, error: null });
        wsRef.current = null;

        // Attempt to reconnect
        if (reconnectCountRef.current < reconnectAttempts) {
          reconnectCountRef.current += 1;
          console.log(
            `Attempting to reconnect (${reconnectCountRef.current}/${reconnectAttempts})...`,
          );

          reconnectTimeoutRef.current = setTimeout(() => {
            connect();
          }, reconnectInterval);
        } else {
          setState((prev) => ({
            ...prev,
            error: new Error("Max reconnection attempts reached"),
          }));
        }
      };

      wsRef.current = ws;
    } catch (error) {
      console.error("Failed to create WebSocket:", error);
      setState((prev) => ({
        ...prev,
        isConnecting: false,
        error: error instanceof Error ? error : new Error("Unknown error"),
      }));
    }
  }, [
    url,
    userId,
    onOCRProgress,
    onTTSProgress,
    onLearningUpdate,
    onError,
    reconnectAttempts,
    reconnectInterval,
  ]);

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }

    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }

    reconnectCountRef.current = 0;
    setState({ isConnected: false, isConnecting: false, error: null });
  }, []);

  const sendPing = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(
        JSON.stringify({
          type: "ping",
          data: null,
          timestamp: new Date().toISOString(),
        }),
      );
    }
  }, []);

  useEffect(() => {
    connect();

    // Send ping every 30 seconds to keep connection alive
    const pingInterval = setInterval(() => {
      sendPing();
    }, 30000);

    return () => {
      clearInterval(pingInterval);
      disconnect();
    };
  }, [connect, disconnect, sendPing]);

  return {
    ...state,
    connect,
    disconnect,
    sendPing,
  };
}
