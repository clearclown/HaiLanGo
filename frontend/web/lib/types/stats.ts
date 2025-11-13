// Statistics types (matching backend models)

export interface LearningTimeStats {
  total_seconds: number;
  total_hours: number;
  daily_average: number;
  weekly_average: number;
  monthly_average: number;
}

export interface ProgressStats {
  completed_pages: number;
  mastered_words: number;
  mastered_phrases: number;
  completed_books: number;
}

export interface StreakStats {
  current_streak: number;
  longest_streak: number;
  last_study_date: string;
}

export interface LearningTimeDataPoint {
  date: string;
  seconds: number;
}

export interface ProgressDataPoint {
  date: string;
  words: number;
  phrases: number;
  pages: number;
}

export interface DashboardStats {
  learning_time: LearningTimeStats;
  progress: ProgressStats;
  streak: StreakStats;
  pronunciation_avg: number;
  weak_words: string[];
  learning_time_chart: LearningTimeDataPoint[];
  progress_chart: ProgressDataPoint[];
}
