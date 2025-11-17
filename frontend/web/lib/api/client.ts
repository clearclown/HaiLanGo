import type { NotificationSettings, Plan, UserProfile, UserSettings } from '@/types/settings';
import type { Book, BookMetadata } from '@/types/book';
import type { UploadMetadata } from '@/types/upload';
import type { ReviewItem, ReviewStats, ReviewResult } from '@/types/review';
import type { DashboardStats, LearningTimeData, ProgressData, WeakPointsData } from '@/types/stats';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class APIClient {
  private async fetch<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.statusText}`);
    }

    return response.json();
  }

  settings = {
    get: async (): Promise<UserSettings> => {
      return this.fetch<UserSettings>('/api/v1/settings');
    },

    updateProfile: async (profile: Partial<UserProfile>): Promise<{ success: boolean }> => {
      return this.fetch<{ success: boolean }>('/api/v1/settings/profile', {
        method: 'PUT',
        body: JSON.stringify(profile),
      });
    },

    updatePassword: async (
      currentPassword: string,
      newPassword: string
    ): Promise<{ success: boolean }> => {
      return this.fetch<{ success: boolean }>('/api/v1/settings/password', {
        method: 'PUT',
        body: JSON.stringify({ currentPassword, newPassword }),
      });
    },

    updateNotifications: async (
      notifications: NotificationSettings
    ): Promise<{ success: boolean }> => {
      return this.fetch<{ success: boolean }>('/api/v1/settings/notifications', {
        method: 'PUT',
        body: JSON.stringify(notifications),
      });
    },
  };

  plan = {
    get: async (): Promise<Plan> => {
      return this.fetch<Plan>('/api/v1/plan');
    },

    upgrade: async (): Promise<{ checkoutUrl: string }> => {
      return this.fetch<{ checkoutUrl: string }>('/api/v1/plan/upgrade', {
        method: 'POST',
      });
    },
  };

  auth = {
    logout: async (): Promise<{ success: boolean }> => {
      return this.fetch<{ success: boolean }>('/api/v1/auth/logout', {
        method: 'POST',
      });
    },
  };

  books = {
    list: async (): Promise<{ books: Book[] }> => {
      return this.fetch<{ books: Book[] }>('/api/v1/books');
    },

    get: async (bookId: string): Promise<Book> => {
      return this.fetch<Book>(`/api/v1/books/${bookId}`);
    },

    create: async (metadata: BookMetadata): Promise<{ book: Book }> => {
      return this.fetch<{ book: Book }>('/api/v1/books', {
        method: 'POST',
        body: JSON.stringify(metadata),
      });
    },

    delete: async (bookId: string): Promise<{ success: boolean }> => {
      return this.fetch<{ success: boolean }>(`/api/v1/books/${bookId}`, {
        method: 'DELETE',
      });
    },

    update: async (bookId: string, metadata: Partial<BookMetadata>): Promise<{ success: boolean }> => {
      return this.fetch<{ success: boolean }>(`/api/v1/books/${bookId}`, {
        method: 'PUT',
        body: JSON.stringify(metadata),
      });
    },
  };

  upload = {
    createBook: async (metadata: Omit<UploadMetadata, 'book_id'>): Promise<{ book_id: string }> => {
      return this.fetch<{ book_id: string }>('/api/v1/upload/create', {
        method: 'POST',
        body: JSON.stringify(metadata),
      });
    },

    uploadFile: async (
      bookId: string,
      file: File,
      onProgress?: (progress: number) => void
    ): Promise<{ success: boolean; file_id: string }> => {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('book_id', bookId);

      return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();

        xhr.upload.addEventListener('progress', (e) => {
          if (e.lengthComputable && onProgress) {
            const progress = (e.loaded / e.total) * 100;
            onProgress(progress);
          }
        });

        xhr.addEventListener('load', () => {
          if (xhr.status >= 200 && xhr.status < 300) {
            resolve(JSON.parse(xhr.responseText));
          } else {
            reject(new Error(`Upload failed: ${xhr.statusText}`));
          }
        });

        xhr.addEventListener('error', () => reject(new Error('Upload failed')));

        xhr.open('POST', `${API_BASE_URL}/api/v1/upload/file`);
        xhr.send(formData);
      });
    },

    complete: async (bookId: string): Promise<{ success: boolean }> => {
      return this.fetch<{ success: boolean }>('/api/v1/upload/complete', {
        method: 'POST',
        body: JSON.stringify({ book_id: bookId }),
      });
    },
  };

  review = {
    getStats: async (): Promise<ReviewStats> => {
      return this.fetch<ReviewStats>('/api/v1/review/stats');
    },

    getItems: async (priority?: 'urgent' | 'recommended' | 'optional'): Promise<{ items: ReviewItem[] }> => {
      const query = priority ? `?priority=${priority}` : '';
      return this.fetch<{ items: ReviewItem[] }>(`/api/v1/review/items${query}`);
    },

    submit: async (result: ReviewResult): Promise<{ success: boolean; next_review: string }> => {
      return this.fetch<{ success: boolean; next_review: string }>('/api/v1/review/submit', {
        method: 'POST',
        body: JSON.stringify(result),
      });
    },
  };

  stats = {
    getDashboard: async (): Promise<DashboardStats> => {
      return this.fetch<DashboardStats>('/api/v1/stats/dashboard');
    },

    getLearningTime: async (period: 'week' | 'month' | 'year' = 'week'): Promise<LearningTimeData> => {
      return this.fetch<LearningTimeData>(`/api/v1/stats/learning-time?period=${period}`);
    },

    getProgress: async (period: 'week' | 'month' | 'year' = 'week'): Promise<ProgressData> => {
      return this.fetch<ProgressData>(`/api/v1/stats/progress?period=${period}`);
    },

    getWeakPoints: async (limit = 10): Promise<WeakPointsData> => {
      return this.fetch<WeakPointsData>(`/api/v1/stats/weak-points?limit=${limit}`);
    },
  };
}

export const apiClient = new APIClient();
