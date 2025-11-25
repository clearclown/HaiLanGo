import type { NotificationSettings, Plan, UserProfile, UserSettings } from '@/types/settings';
import type { Book, BookMetadata } from '@/types/book';
import type { UploadMetadata } from '@/types/upload';
import type { ReviewItem, ReviewStats, ReviewResult } from '@/types/review';
import type { DashboardStats, LearningTimeData, ProgressData, WeakPointsData } from '@/types/stats';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Get auth token from localStorage (client-side only)
function getAuthToken(): string | null {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem('access_token');
}

class APIClient {
  private async fetch<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const token = getAuthToken();
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options?.headers,
    };

    // Add Authorization header if token exists
    if (token) {
      (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers,
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
      // Use the books endpoint to create a book
      const response = await this.fetch<{ book: { id: string } }>('/api/v1/books', {
        method: 'POST',
        body: JSON.stringify({
          title: metadata.title,
          target_language: metadata.target_language,
          native_language: metadata.native_language,
          reference_language: metadata.reference_language || '',
        }),
      });
      // Return in the expected format
      return { book_id: response.book.id };
    },

    uploadFile: async (
      bookId: string,
      file: File,
      onProgress?: (progress: number) => void
    ): Promise<{ success: boolean; file_id: string }> => {
      const formData = new FormData();
      formData.append('files', file);

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
            try {
              const response = JSON.parse(xhr.responseText);
              resolve({ success: true, file_id: response.file_id || bookId });
            } catch {
              resolve({ success: true, file_id: bookId });
            }
          } else {
            reject(new Error(`Upload failed: ${xhr.statusText}`));
          }
        });

        xhr.addEventListener('error', () => reject(new Error('Upload failed')));

        xhr.open('POST', `${API_BASE_URL}/api/v1/upload/books/${bookId}/files`);
        // Add Authorization header if token exists
        const token = getAuthToken();
        if (token) {
          xhr.setRequestHeader('Authorization', `Bearer ${token}`);
        }
        xhr.send(formData);
      });
    },

    complete: async (bookId: string): Promise<{ success: boolean }> => {
      // The backend doesn't have a complete endpoint, just return success
      return Promise.resolve({ success: true });
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
