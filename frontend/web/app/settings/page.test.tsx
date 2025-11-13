import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import SettingsPage from './page';

// APIクライアントをモック
vi.mock('@/lib/api/client', () => ({
  apiClient: {
    settings: {
      get: vi.fn().mockResolvedValue({
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
      }),
      updateProfile: vi.fn().mockResolvedValue({ success: true }),
      updateNotifications: vi.fn().mockResolvedValue({ success: true }),
    },
    plan: {
      get: vi.fn().mockResolvedValue({
        type: 'free',
      }),
    },
  },
}));

describe('SettingsPage', () => {
  it('設定画面が表示される', async () => {
    render(<SettingsPage />);

    await waitFor(() => {
      expect(screen.getByText('設定')).toBeInTheDocument();
    });

    expect(screen.getByText('アカウント')).toBeInTheDocument();
    expect(screen.getByText('プラン')).toBeInTheDocument();
    expect(screen.getByText('通知設定')).toBeInTheDocument();
  });

  it('アカウント設定が表示される', async () => {
    render(<SettingsPage />);

    await waitFor(() => {
      expect(screen.getByDisplayValue('太郎')).toBeInTheDocument();
      expect(screen.getByDisplayValue('taro@example.com')).toBeInTheDocument();
    });
  });

  it('通知設定が表示される', async () => {
    render(<SettingsPage />);

    await waitFor(() => {
      expect(screen.getByLabelText('学習リマインダー')).toBeChecked();
      expect(screen.getByLabelText('復習通知')).toBeChecked();
      expect(screen.getByLabelText('メール通知')).not.toBeChecked();
    });
  });

  it('ログアウトボタンが表示される', async () => {
    render(<SettingsPage />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /ログアウト/ })).toBeInTheDocument();
    });
  });

  it('ログアウトができる', async () => {
    const user = userEvent.setup();
    const mockLogout = vi.fn();
    global.localStorage.clear = mockLogout;

    render(<SettingsPage />);

    await waitFor(() => {
      const logoutButton = screen.getByRole('button', { name: /ログアウト/ });
      expect(logoutButton).toBeInTheDocument();
    });

    const logoutButton = screen.getByRole('button', { name: /ログアウト/ });
    await user.click(logoutButton);

    // ログアウト確認ダイアログが表示されることを確認
    await waitFor(() => {
      expect(screen.getByText(/ログアウトしますか/)).toBeInTheDocument();
    });
  });
});
