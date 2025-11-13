export type NotificationType =
  | "ocr_progress"
  | "tts_progress"
  | "learning_update"
  | "error"
  | "ping"
  | "pong";

export interface Notification<T = unknown> {
  type: NotificationType;
  data: T;
  timestamp: string;
}

export interface OCRProgressData {
  book_id: string;
  total_pages: number;
  processed_pages: number;
  current_page: number;
  progress: number; // 0-100
  estimated_time_ms: number;
  status: "processing" | "completed" | "failed";
}

export interface TTSProgressData {
  book_id: string;
  page_number: number;
  total_segments: number;
  processed_segments: number;
  progress: number; // 0-100
  status: "processing" | "completed" | "failed";
}

export interface LearningUpdateData {
  user_id: string;
  book_id: string;
  page_number: number;
  completed_pages: number;
  total_pages: number;
  learned_words: number;
  study_time_ms: number;
}

export interface ErrorData {
  code: string;
  message: string;
  details?: string;
}
