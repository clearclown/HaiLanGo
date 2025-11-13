import type { NotificationSettings, Plan, UserProfile, UserSettings } from '@/types/settings';

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
}

export const apiClient = new APIClient();
