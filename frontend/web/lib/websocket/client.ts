import type {
  Message,
  MessageType,
  MessageHandler,
  WebSocketConfig,
} from './types';

const DEFAULT_CONFIG: Required<WebSocketConfig> = {
  url: process.env.NEXT_PUBLIC_WS_URL ||
       (typeof window !== 'undefined'
         ? `ws://${window.location.hostname}:8080/api/v1/ws`
         : 'ws://localhost:8080/api/v1/ws'),
  reconnectInterval: 1000,
  maxReconnectAttempts: 5,
  heartbeatInterval: 30000,
};

export class WebSocketClient {
  private ws: WebSocket | null = null;
  private config: Required<WebSocketConfig>;
  private listeners: Map<MessageType, Set<MessageHandler>> = new Map();
  private reconnectAttempts = 0;
  private reconnectTimeout: NodeJS.Timeout | null = null;
  private heartbeatInterval: NodeJS.Timeout | null = null;
  private isConnecting = false;
  private isIntentionallyClosed = false;
  private token: string | null = null;

  constructor(config?: WebSocketConfig) {
    this.config = { ...DEFAULT_CONFIG, ...config };
  }

  /**
   * WebSocket接続を確立する
   * @param token JWT認証トークン
   */
  public connect(token: string): void {
    if (this.isConnecting || this.ws?.readyState === WebSocket.OPEN) {
      console.warn('WebSocket is already connecting or connected');
      return;
    }

    this.token = token;
    this.isIntentionallyClosed = false;
    this.isConnecting = true;

    try {
      const wsUrl = `${this.config.url}?token=${encodeURIComponent(token)}`;
      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = this.handleOpen.bind(this);
      this.ws.onmessage = this.handleMessage.bind(this);
      this.ws.onerror = this.handleError.bind(this);
      this.ws.onclose = this.handleClose.bind(this);
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error);
      this.isConnecting = false;
      this.handleReconnect();
    }
  }

  /**
   * WebSocket接続を切断する
   */
  public disconnect(): void {
    this.isIntentionallyClosed = true;
    this.clearReconnectTimeout();
    this.clearHeartbeat();

    if (this.ws) {
      this.ws.close(1000, 'Client disconnected');
      this.ws = null;
    }

    this.reconnectAttempts = 0;
  }

  /**
   * 接続状態を取得する
   */
  public isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }

  /**
   * メッセージタイプごとのリスナーを登録する
   * @param type メッセージタイプ
   * @param handler ハンドラー関数
   */
  public on<T = unknown>(type: MessageType, handler: MessageHandler<T>): void {
    if (!this.listeners.has(type)) {
      this.listeners.set(type, new Set());
    }
    this.listeners.get(type)!.add(handler as MessageHandler);
  }

  /**
   * リスナーを解除する
   * @param type メッセージタイプ
   * @param handler ハンドラー関数
   */
  public off<T = unknown>(type: MessageType, handler: MessageHandler<T>): void {
    const handlers = this.listeners.get(type);
    if (handlers) {
      handlers.delete(handler as MessageHandler);
    }
  }

  /**
   * メッセージを送信する
   * @param message 送信するメッセージ
   */
  public send(message: unknown): void {
    if (!this.isConnected()) {
      console.error('WebSocket is not connected');
      return;
    }

    try {
      this.ws!.send(JSON.stringify(message));
    } catch (error) {
      console.error('Failed to send message:', error);
    }
  }

  private handleOpen(): void {
    console.log('WebSocket connected');
    this.isConnecting = false;
    this.reconnectAttempts = 0;
    this.startHeartbeat();
  }

  private handleMessage(event: MessageEvent): void {
    try {
      const messages = event.data.split('\n').filter((line: string) => line.trim());

      for (const messageData of messages) {
        const message: Message = JSON.parse(messageData);
        this.dispatchMessage(message);
      }
    } catch (error) {
      console.error('Failed to parse WebSocket message:', error);
    }
  }

  private handleError(event: Event): void {
    console.error('WebSocket error:', event);
    this.isConnecting = false;
  }

  private handleClose(event: CloseEvent): void {
    console.log('WebSocket closed:', event.code, event.reason);
    this.isConnecting = false;
    this.clearHeartbeat();

    if (!this.isIntentionallyClosed) {
      this.handleReconnect();
    }
  }

  private handleReconnect(): void {
    if (this.reconnectAttempts >= this.config.maxReconnectAttempts) {
      console.error('Max reconnection attempts reached');
      return;
    }

    this.clearReconnectTimeout();

    const delay = this.config.reconnectInterval * Math.pow(2, this.reconnectAttempts);
    this.reconnectAttempts++;

    console.log(`Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts}/${this.config.maxReconnectAttempts})`);

    this.reconnectTimeout = setTimeout(() => {
      if (this.token && !this.isIntentionallyClosed) {
        this.connect(this.token);
      }
    }, delay);
  }

  private clearReconnectTimeout(): void {
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }
  }

  private startHeartbeat(): void {
    this.clearHeartbeat();

    this.heartbeatInterval = setInterval(() => {
      if (this.isConnected()) {
        // Send ping (WebSocket ping/pong is handled automatically by browser)
        // We can send a custom ping if needed
        this.send({ type: 'ping' });
      }
    }, this.config.heartbeatInterval);
  }

  private clearHeartbeat(): void {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
      this.heartbeatInterval = null;
    }
  }

  private dispatchMessage(message: Message): void {
    const handlers = this.listeners.get(message.type);
    if (handlers && handlers.size > 0) {
      handlers.forEach(handler => {
        try {
          handler(message.payload);
        } catch (error) {
          console.error(`Error in message handler for type ${message.type}:`, error);
        }
      });
    } else {
      console.log('No handlers registered for message type:', message.type);
    }
  }
}

// Singleton instance
let wsClient: WebSocketClient | null = null;

/**
 * WebSocketクライアントのシングルトンインスタンスを取得
 */
export function getWebSocketClient(config?: WebSocketConfig): WebSocketClient {
  if (!wsClient) {
    wsClient = new WebSocketClient(config);
  }
  return wsClient;
}

/**
 * WebSocketクライアントをリセット（主にテスト用）
 */
export function resetWebSocketClient(): void {
  if (wsClient) {
    wsClient.disconnect();
    wsClient = null;
  }
}
