import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import PlanSettings from './PlanSettings';

describe('PlanSettings', () => {
  it('無料プランが表示される', () => {
    const mockPlan = { type: 'free' as const };
    const mockOnUpgrade = vi.fn();

    render(<PlanSettings plan={mockPlan} onUpgrade={mockOnUpgrade} />);

    expect(screen.getByText('無料プラン')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /プレミアムにアップグレード/ })).toBeInTheDocument();
  });

  it('プレミアムプランが表示される', () => {
    const mockPlan = {
      type: 'premium' as const,
      expiresAt: '2025-12-31T00:00:00Z',
    };
    const mockOnUpgrade = vi.fn();

    render(<PlanSettings plan={mockPlan} onUpgrade={mockOnUpgrade} />);

    expect(screen.getByText('プレミアムプラン')).toBeInTheDocument();
    expect(screen.getByText(/有効期限/)).toBeInTheDocument();
  });

  it('アップグレードボタンをクリックできる', async () => {
    const user = userEvent.setup();
    const mockPlan = { type: 'free' as const };
    const mockOnUpgrade = vi.fn();

    render(<PlanSettings plan={mockPlan} onUpgrade={mockOnUpgrade} />);

    const upgradeButton = screen.getByRole('button', { name: /プレミアムにアップグレード/ });
    await user.click(upgradeButton);

    expect(mockOnUpgrade).toHaveBeenCalled();
  });
});
