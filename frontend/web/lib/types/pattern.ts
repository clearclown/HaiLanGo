export type PatternType =
  | "greeting"
  | "question"
  | "response"
  | "request"
  | "confirmation"
  | "other";

export interface Pattern {
  id: string;
  book_id: string;
  type: PatternType;
  pattern: string;
  translation: string;
  frequency: number;
  created_at: string;
  updated_at: string;
}

export interface PatternExample {
  id: string;
  pattern_id: string;
  page_number: number;
  original_text: string;
  translated_text: string;
  context: string;
  created_at: string;
}

export interface PatternPractice {
  id: string;
  pattern_id: string;
  question: string;
  correct_answer: string;
  alternative_answers: string[];
  difficulty: number;
  created_at: string;
}

export interface PatternProgress {
  id: string;
  user_id: string;
  pattern_id: string;
  mastery_level: number;
  practice_count: number;
  correct_count: number;
  last_practiced_at?: string;
  created_at: string;
  updated_at: string;
}

export interface PatternExtractionRequest {
  book_id: string;
  page_start: number;
  page_end: number;
  min_frequency: number;
}

export interface PatternExtractionResponse {
  patterns: Pattern[];
  total_found: number;
  processed_pages: number;
  duration: number;
}
