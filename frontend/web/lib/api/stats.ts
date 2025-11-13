import axios from 'axios';
import type {
  DashboardStats,
  LearningTimeStats,
  ProgressStats,
  StreakStats,
  LearningTimeDataPoint,
  ProgressDataPoint,
} from '../types/stats';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Create axios instance with default config
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const statsApi = {
  // Get dashboard stats
  getDashboard: async (): Promise<DashboardStats> => {
    const response = await apiClient.get<DashboardStats>('/api/v1/stats/dashboard');
    return response.data;
  },

  // Get learning time stats
  getLearningTime: async (): Promise<LearningTimeStats> => {
    const response = await apiClient.get<LearningTimeStats>('/api/v1/stats/learning-time');
    return response.data;
  },

  // Get progress stats
  getProgress: async (): Promise<ProgressStats> => {
    const response = await apiClient.get<ProgressStats>('/api/v1/stats/progress');
    return response.data;
  },

  // Get streak stats
  getStreak: async (): Promise<StreakStats> => {
    const response = await apiClient.get<StreakStats>('/api/v1/stats/streak');
    return response.data;
  },

  // Get learning time chart
  getLearningTimeChart: async (days = 7): Promise<LearningTimeDataPoint[]> => {
    const response = await apiClient.get<LearningTimeDataPoint[]>(
      `/api/v1/stats/learning-time-chart?days=${days}`,
    );
    return response.data;
  },

  // Get progress chart
  getProgressChart: async (days = 30): Promise<ProgressDataPoint[]> => {
    const response = await apiClient.get<ProgressDataPoint[]>(
      `/api/v1/stats/progress-chart?days=${days}`,
    );
    return response.data;
  },

  // Get weak words
  getWeakWords: async (limit = 10): Promise<string[]> => {
    const response = await apiClient.get<{ weak_words: string[] }>(
      `/api/v1/stats/weak-words?limit=${limit}`,
    );
    return response.data.weak_words;
  },
};
