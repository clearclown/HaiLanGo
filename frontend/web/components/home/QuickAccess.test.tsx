import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { QuickAccess } from './QuickAccess';

describe('QuickAccess', () => {
  const mockData = {
    booksCount: 5,
    reviewItemsCount: 12,
  };

  it('should render My Books shortcut', () => {
    render(<QuickAccess data={mockData} />);

    expect(screen.getByText(/マイ本/)).toBeDefined();
  });

  it('should display books count', () => {
    render(<QuickAccess data={mockData} />);

    expect(screen.getByText('5冊')).toBeDefined();
  });

  it('should render Review shortcut', () => {
    render(<QuickAccess data={mockData} />);

    expect(screen.getByText(/復習/)).toBeDefined();
  });

  it('should display review items count', () => {
    render(<QuickAccess data={mockData} />);

    expect(screen.getByText('12項目')).toBeDefined();
  });

  it('should render book icon', () => {
    const { container } = render(<QuickAccess data={mockData} />);

    expect(container.textContent).toContain('📖');
  });

  it('should render target icon', () => {
    const { container } = render(<QuickAccess data={mockData} />);

    expect(container.textContent).toContain('🎯');
  });
});
