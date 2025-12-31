import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import SettingsPage from './page';

// APIクライアントをモック
const mockGet = vi.fn();
const mockPlanGet = vi.fn();

vi.mock('@/lib/api/client', () => ({
  apiClient: {
    settings: {
      get: () => mockGet(),
      updateProfile: vi.fn().mockResolvedValue({ success: true }),
      updateNotifications: vi.fn().mockResolvedValue({ success: true }),
    },
    plan: {
      get: () => mockPlanGet(),
    },
    auth: {
      logout: vi.fn().mockResolvedValue({}),
    },
  },
}));

describe('SettingsPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockGet.mockResolvedValue({
      profile: {
        id: '1',
        name: '太郎',
        email: 'taro@example.com',
      },
      notifications: {
        learningReminder: true,
        reviewNotification: true,
        emailNotification: false,
      },
      interfaceLanguage: 'ja',
    });
    mockPlanGet.mockResolvedValue({
      type: 'free',
    });
  });

  it('設定画面が表示される', async () => {
    await act(async () => {
      render(<SettingsPage />);
    });

    // Wait for loading to complete and check for heading
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 100));
    });

    // Use heading role to avoid matching navbar link
    expect(screen.getByRole('heading', { name: '設定' })).toBeInTheDocument();
    expect(screen.getByText('アカウント')).toBeInTheDocument();
    expect(screen.getByText('プラン')).toBeInTheDocument();
    expect(screen.getByText('通知設定')).toBeInTheDocument();
  });

  it('ローディング状態を表示する', () => {
    // Make the API never resolve
    mockGet.mockReturnValue(new Promise(() => {}));
    mockPlanGet.mockReturnValue(new Promise(() => {}));

    render(<SettingsPage />);

    expect(screen.getByText('読み込み中...')).toBeInTheDocument();
  });
});
