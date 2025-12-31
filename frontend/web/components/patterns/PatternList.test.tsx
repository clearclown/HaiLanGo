import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import type { Pattern } from '../../lib/types/pattern';
import { PatternList } from './PatternList';

// Mock patterns data
const mockPatterns: Pattern[] = [
  {
    id: '1',
    book_id: 'book-1',
    type: 'greeting',
    pattern: 'Hello',
    translation: 'こんにちは',
    frequency: 5,
    created_at: '2025-01-01T00:00:00Z',
    updated_at: '2025-01-01T00:00:00Z',
  },
  {
    id: '2',
    book_id: 'book-1',
    type: 'question',
    pattern: 'How are you?',
    translation: '元気ですか？',
    frequency: 3,
    created_at: '2025-01-01T00:00:00Z',
    updated_at: '2025-01-01T00:00:00Z',
  },
  {
    id: '3',
    book_id: 'book-1',
    type: 'response',
    pattern: 'Thank you',
    translation: 'ありがとう',
    frequency: 4,
    created_at: '2025-01-01T00:00:00Z',
    updated_at: '2025-01-01T00:00:00Z',
  },
];

describe('PatternList', () => {
  it('renders pattern list', () => {
    render(<PatternList patterns={mockPatterns} />);

    expect(screen.getByText('Hello')).toBeInTheDocument();
    expect(screen.getByText('How are you?')).toBeInTheDocument();
    expect(screen.getByText('Thank you')).toBeInTheDocument();
  });

  it('displays pattern translations', () => {
    render(<PatternList patterns={mockPatterns} />);

    expect(screen.getByText('こんにちは')).toBeInTheDocument();
    expect(screen.getByText('元気ですか？')).toBeInTheDocument();
    expect(screen.getByText('ありがとう')).toBeInTheDocument();
  });

  it('displays pattern frequencies', () => {
    render(<PatternList patterns={mockPatterns} />);

    expect(screen.getByText(/5/)).toBeInTheDocument();
    expect(screen.getByText(/3/)).toBeInTheDocument();
    expect(screen.getByText(/4/)).toBeInTheDocument();
  });

  it('filters patterns by type', () => {
    const onFilterChange = vi.fn();
    render(<PatternList patterns={mockPatterns} onFilterChange={onFilterChange} />);

    // The button label is "Greeting" (capitalized)
    const filterButton = screen.getByRole('button', { name: 'Greeting' });
    fireEvent.click(filterButton);

    expect(onFilterChange).toHaveBeenCalledWith('greeting');
  });

  it('calls onPatternClick when pattern is clicked', () => {
    const onPatternClick = vi.fn();
    render(<PatternList patterns={mockPatterns} onPatternClick={onPatternClick} />);

    const patternCard = screen.getByText('Hello').closest('div');
    if (patternCard) {
      fireEvent.click(patternCard);
    }

    expect(onPatternClick).toHaveBeenCalledWith(mockPatterns[0]);
  });

  it('renders empty state when no patterns', () => {
    render(<PatternList patterns={[]} />);

    expect(screen.getByText(/no patterns found/i)).toBeInTheDocument();
  });

  it('sorts patterns by frequency', () => {
    render(<PatternList patterns={mockPatterns} sortBy="frequency" />);

    const patterns = screen.getAllByTestId('pattern-card');
    expect(patterns[0]).toHaveTextContent('Hello');
    expect(patterns[1]).toHaveTextContent('Thank you');
    expect(patterns[2]).toHaveTextContent('How are you?');
  });
});
