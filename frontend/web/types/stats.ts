// Stats types matching backend models

export interface DashboardStats {
  learningTimeToday: number;
  learningTimeThisWeek: number;
  totalLearningTime: number;
  currentStreak: number;
  longestStreak: number;
  completedPages: number;
  totalPages: number;
  masteredWords: number;
  masteredPhrases: number;
  completedBooks: number;
  totalBooks: number;
  averagePronunciationScore: number;
}

export interface DailyLearningTime {
  date: string;
  minutes: number;
}

export interface LearningTimeData {
  period: string;
  data: DailyLearningTime[];
  totalMinutes: number;
  averageMinutes: number;
}

export interface TimeSeriesData {
  date: string;
  count: number;
}

export interface ProgressData {
  period: string;
  words: TimeSeriesData[];
  phrases: TimeSeriesData[];
  pages: TimeSeriesData[];
}

export interface WeakItem {
  word?: string;
  phrase?: string;
  language: string;
  attempts: number;
  averageScore: number;
}

export interface WeakPointsData {
  weakWords: WeakItem[];
  weakPhrases: WeakItem[];
}

export interface StatsResponse {
  dashboard: DashboardStats;
  learningTime: LearningTimeData;
  progress: ProgressData;
  weakPoints: WeakPointsData;
}
