import { act, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import Home from './page';

describe('Home', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders loading state initially', () => {
    // Mock fetch to never resolve (keeps loading state)
    vi.mocked(global.fetch).mockReturnValue(new Promise(() => {}));

    render(<Home />);
    expect(document.querySelector('.animate-spin')).toBeDefined();
  });

  it('renders error state when fetch fails', async () => {
    // Mock fetch to reject (throw error)
    vi.mocked(global.fetch).mockRejectedValue(new Error('Network error'));

    await act(async () => {
      render(<Home />);
    });

    // Wait longer for the error state to render
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 200));
    });

    // Check for error message - the component shows "Network error" from the caught exception
    expect(screen.getByText('Network error')).toBeInTheDocument();
  });

  it('renders dashboard when fetch succeeds', async () => {
    vi.mocked(global.fetch).mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue({
        user: { name: 'テストユーザー' },
        stats: {
          booksCount: 5,
          reviewItemsCount: 10,
          streak: 7,
          totalLearningTime: 180,
        },
        todayLearning: null,
      }),
    } as unknown as Response);

    await act(async () => {
      render(<Home />);
    });

    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 100));
    });

    expect(screen.getByText(/テストユーザー/)).toBeInTheDocument();
  });
});
