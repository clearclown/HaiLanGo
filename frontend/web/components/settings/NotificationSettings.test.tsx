import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import NotificationSettings from './NotificationSettings';

describe('NotificationSettings', () => {
  const mockSettings = {
    learningReminder: true,
    reviewNotification: true,
    emailNotification: false,
  };

  const mockOnUpdate = vi.fn();

  it('通知設定が表示される', () => {
    render(<NotificationSettings settings={mockSettings} onUpdate={mockOnUpdate} />);

    expect(screen.getByLabelText('学習リマインダー')).toBeChecked();
    expect(screen.getByLabelText('復習通知')).toBeChecked();
    expect(screen.getByLabelText('メール通知')).not.toBeChecked();
  });

  it('学習リマインダーをオフにできる', async () => {
    const user = userEvent.setup();
    render(<NotificationSettings settings={mockSettings} onUpdate={mockOnUpdate} />);

    const learningReminderSwitch = screen.getByLabelText('学習リマインダー');
    await user.click(learningReminderSwitch);

    await waitFor(() => {
      expect(mockOnUpdate).toHaveBeenCalledWith({
        ...mockSettings,
        learningReminder: false,
      });
    });
  });

  it('復習通知をオフにできる', async () => {
    const user = userEvent.setup();
    render(<NotificationSettings settings={mockSettings} onUpdate={mockOnUpdate} />);

    const reviewNotificationSwitch = screen.getByLabelText('復習通知');
    await user.click(reviewNotificationSwitch);

    await waitFor(() => {
      expect(mockOnUpdate).toHaveBeenCalledWith({
        ...mockSettings,
        reviewNotification: false,
      });
    });
  });

  it('メール通知をオンにできる', async () => {
    const user = userEvent.setup();
    render(<NotificationSettings settings={mockSettings} onUpdate={mockOnUpdate} />);

    const emailNotificationSwitch = screen.getByLabelText('メール通知');
    await user.click(emailNotificationSwitch);

    await waitFor(() => {
      expect(mockOnUpdate).toHaveBeenCalledWith({
        ...mockSettings,
        emailNotification: true,
      });
    });
  });
});
