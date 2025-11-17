// WebSocket message types (matching backend)
export type MessageType =
  | 'ocr_progress'
  | 'book_ready'
  | 'review_reminder'
  | 'learning_update'
  | 'notification'
  | 'error'
  | 'connection_established';

export type NotificationLevel = 'info' | 'success' | 'warning' | 'error';

// Base message structure
export interface Message {
  type: MessageType;
  payload: unknown;
  timestamp: string;
}

// Specific payload types
export interface OCRProgressPayload {
  bookId: string;
  totalPages: number;
  processedPages: number;
  progress: number;
  currentPage?: number;
  status: string;
  message?: string;
}

export interface BookReadyPayload {
  bookId: string;
  title: string;
  totalPages: number;
  message: string;
}

export interface ReviewItem {
  id: string;
  content: string;
  translation?: string;
  dueDate: string;
  priority: 'urgent' | 'recommended' | 'optional';
}

export interface ReviewReminderPayload {
  count: number;
  items: ReviewItem[];
  message: string;
}

export interface LearningStats {
  totalTime: number; // seconds
  completedPages: number;
  masteredWords: number;
  pronunciationScore?: number;
  streakDays: number;
  lastStudiedAt?: string;
}

export interface LearningUpdatePayload {
  sessionId: string;
  bookId?: string;
  stats: LearningStats;
  message?: string;
}

export interface NotificationPayload {
  title: string;
  message: string;
  level: NotificationLevel;
  action?: {
    label: string;
    url: string;
  };
}

export interface ErrorPayload {
  code: string;
  message: string;
  details?: string;
}

export interface ConnectionEstablishedPayload {
  userId: string;
  message: string;
  timestamp: string;
}

// Typed message handlers
export type MessageHandler<T = unknown> = (payload: T) => void;

export interface WebSocketConfig {
  url?: string;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
  heartbeatInterval?: number;
}
