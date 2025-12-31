export interface Word {
  id: string;
  user_id: string;
  book_id: string;
  page_number: number;
  text: string;
  meaning: string;
  pronunciation: string;
  part_of_speech: string;
  example: string;
  language: string;
  review_count: number;
  average_score: number;
  mastery: number;
  tags: string[];
  last_reviewed_at: string;
  created_at: string;
  updated_at: string;
}

export interface WordFilter {
  book_id?: string;
  language?: string;
  query?: string;
  tags?: string[];
  min_mastery?: number;
  max_mastery?: number;
  limit?: number;
  offset?: number;
  sort_by?: 'created_at' | 'mastery' | 'review_count';
  sort_order?: 'asc' | 'desc';
}

export interface WordStats {
  total_words: number;
  mastered_words: number;
  average_mastery: number;
  total_reviews: number;
}

export interface AddWordRequest {
  book_id?: string;
  page_number?: number;
  text: string;
  meaning?: string;
  pronunciation?: string;
  part_of_speech?: string;
  example?: string;
  language: string;
  tags?: string[];
}

export interface UpdateWordRequest {
  text?: string;
  meaning?: string;
  pronunciation?: string;
  part_of_speech?: string;
  example?: string;
  tags?: string[];
}

export interface RecordReviewRequest {
  score: number;
}
