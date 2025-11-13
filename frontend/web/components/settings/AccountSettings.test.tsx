import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import AccountSettings from './AccountSettings';

describe('AccountSettings', () => {
  const mockProfile = {
    id: '1',
    name: '太郎',
    email: 'taro@example.com',
  };

  const mockOnUpdate = vi.fn().mockResolvedValue(undefined);

  it('プロフィール情報が表示される', () => {
    render(<AccountSettings profile={mockProfile} onUpdate={mockOnUpdate} />);

    expect(screen.getByDisplayValue('太郎')).toBeInTheDocument();
    expect(screen.getByDisplayValue('taro@example.com')).toBeInTheDocument();
  });

  it('名前を変更できる', async () => {
    const user = userEvent.setup();
    render(<AccountSettings profile={mockProfile} onUpdate={mockOnUpdate} />);

    const nameInput = screen.getByLabelText('名前');
    await user.clear(nameInput);
    await user.type(nameInput, '次郎');

    const saveButton = screen.getByRole('button', { name: '保存' });
    await user.click(saveButton);

    await waitFor(() => {
      expect(mockOnUpdate).toHaveBeenCalledWith({
        name: '次郎',
        email: 'taro@example.com',
      });
    });
  });

  it('メールアドレスを変更できる', async () => {
    const user = userEvent.setup();
    render(<AccountSettings profile={mockProfile} onUpdate={mockOnUpdate} />);

    const emailInput = screen.getByLabelText('メールアドレス');
    await user.clear(emailInput);
    await user.type(emailInput, 'jiro@example.com');

    const saveButton = screen.getByRole('button', { name: '保存' });
    await user.click(saveButton);

    await waitFor(() => {
      expect(mockOnUpdate).toHaveBeenCalledWith({
        name: '太郎',
        email: 'jiro@example.com',
      });
    });
  });

  it('バリデーションエラーが表示される', async () => {
    mockOnUpdate.mockClear();
    render(<AccountSettings profile={mockProfile} onUpdate={mockOnUpdate} />);

    const emailInput = screen.getByLabelText('メールアドレス');
    fireEvent.change(emailInput, { target: { value: 'invalid' } });

    const form = emailInput.closest('form');
    if (form) {
      fireEvent.submit(form);
    }

    await waitFor(() => {
      expect(screen.getByText(/有効なメールアドレスを入力してください/)).toBeInTheDocument();
    });
    expect(mockOnUpdate).not.toHaveBeenCalled();
  });
});
