export type ReviewPriority = 'urgent' | 'recommended' | 'optional';

export interface ReviewItem {
  id: string;
  type: 'word' | 'phrase';
  text: string;
  translation: string;
  language: string;
  mastery_level: number; // 0-100
  last_reviewed: string;
  next_review: string;
  priority: ReviewPriority;
}

export interface ReviewStats {
  urgent_count: number;
  recommended_count: number;
  optional_count: number;
  total_completed_today: number;
  weekly_completion_rate: number;
}

export interface ReviewResult {
  item_id: string;
  score: number; // 0-100
  pronunciation_score?: number;
  completed_at: string;
}
