import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { LearningStats } from './LearningStats';

describe('LearningStats', () => {
  const mockStats = {
    streakDays: 7,
    totalLearningTimeSeconds: 13320, // 3時間42分
    completedPagesCount: 12,
    booksCount: 5,
    reviewItemsCount: 12,
  };

  it('should render section title', () => {
    render(<LearningStats stats={mockStats} />);

    expect(screen.getByText(/学習統計/)).toBeDefined();
  });

  it('should display streak days', () => {
    const { container } = render(<LearningStats stats={mockStats} />);

    expect(screen.getByText(/連続学習/)).toBeDefined();
    expect(container.textContent).toContain('7日');
  });

  it('should show fire emoji for streak', () => {
    const { container } = render(<LearningStats stats={mockStats} />);

    expect(container.textContent).toContain('🔥');
  });

  it('should display total learning time formatted correctly', () => {
    render(<LearningStats stats={mockStats} />);

    expect(screen.getByText(/総学習時間/)).toBeDefined();
    // 13320 seconds = 3 hours 42 minutes
    expect(screen.getByText(/3時間42分/)).toBeDefined();
  });

  it('should format learning time with hours and minutes', () => {
    const statsWithDifferentTime = {
      ...mockStats,
      totalLearningTimeSeconds: 7260, // 2時間1分
    };

    render(<LearningStats stats={statsWithDifferentTime} />);

    expect(screen.getByText(/2時間1分/)).toBeDefined();
  });

  it('should handle zero streak days', () => {
    const statsWithNoStreak = {
      ...mockStats,
      streakDays: 0,
    };

    const { container } = render(<LearningStats stats={statsWithNoStreak} />);

    expect(container.textContent).toContain('0日');
  });
});
