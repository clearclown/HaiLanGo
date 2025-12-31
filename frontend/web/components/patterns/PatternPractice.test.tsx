import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import type { Pattern, PatternPractice as PatternPracticeType } from '../../lib/types/pattern';
import { PatternPractice } from './PatternPractice';

// Mock Math.random to return predictable values (keeps original order)
const originalRandom = Math.random;
beforeEach(() => {
  vi.useFakeTimers();
  Math.random = vi.fn(() => 0.6); // Returns value that keeps sort order stable
});
afterEach(() => {
  vi.useRealTimers();
  Math.random = originalRandom;
});

const mockPattern: Pattern = {
  id: 'pattern-1',
  book_id: 'book-1',
  type: 'greeting',
  pattern: 'Hello',
  translation: 'こんにちは',
  frequency: 5,
  created_at: '2025-01-01T00:00:00Z',
  updated_at: '2025-01-01T00:00:00Z',
};

const mockPractices: PatternPracticeType[] = [
  {
    id: 'practice-1',
    pattern_id: 'pattern-1',
    question: "How do you say 'Hello' in Japanese?",
    correct_answer: 'こんにちは',
    alternative_answers: ['おはよう', 'こんばんは', 'さようなら'],
    difficulty: 1,
    created_at: '2025-01-01T00:00:00Z',
  },
  {
    id: 'practice-2',
    pattern_id: 'pattern-1',
    question: "Choose the correct response to 'Hello'",
    correct_answer: 'Hello, how are you?',
    alternative_answers: ['Goodbye', 'Thank you', 'Sorry'],
    difficulty: 2,
    created_at: '2025-01-01T00:00:00Z',
  },
];

describe('PatternPractice', () => {
  it('renders pattern information', () => {
    render(<PatternPractice pattern={mockPattern} practices={mockPractices} />);

    expect(screen.getByText('Hello')).toBeInTheDocument();
    // "こんにちは" appears both in pattern info and as answer option
    // Check for at least one occurrence
    expect(screen.getAllByText('こんにちは').length).toBeGreaterThan(0);
  });

  it('displays practice question', () => {
    render(<PatternPractice pattern={mockPattern} practices={mockPractices} />);

    expect(screen.getByText("How do you say 'Hello' in Japanese?")).toBeInTheDocument();
  });

  it('shows all answer options', () => {
    render(<PatternPractice pattern={mockPattern} practices={mockPractices} />);

    // Use getAllByRole to check answer buttons
    const buttons = screen.getAllByRole('button');
    expect(buttons.length).toBe(4);

    // Check that each answer text exists (as button name)
    expect(screen.getByRole('button', { name: 'こんにちは' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'おはよう' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'こんばんは' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'さようなら' })).toBeInTheDocument();
  });

  it('marks correct answer when clicked', async () => {
    const onAnswerSubmit = vi.fn();
    render(
      <PatternPractice
        pattern={mockPattern}
        practices={mockPractices}
        onAnswerSubmit={onAnswerSubmit}
      />
    );

    const correctAnswer = screen.getByRole('button', { name: 'こんにちは' });
    await act(async () => {
      fireEvent.click(correctAnswer);
    });

    expect(onAnswerSubmit).toHaveBeenCalledWith({
      practice_id: 'practice-1',
      answer: 'こんにちは',
      is_correct: true,
    });
  });

  it('marks incorrect answer when clicked', async () => {
    const onAnswerSubmit = vi.fn();
    render(
      <PatternPractice
        pattern={mockPattern}
        practices={mockPractices}
        onAnswerSubmit={onAnswerSubmit}
      />
    );

    const incorrectAnswer = screen.getByRole('button', { name: 'おはよう' });
    await act(async () => {
      fireEvent.click(incorrectAnswer);
    });

    expect(onAnswerSubmit).toHaveBeenCalledWith({
      practice_id: 'practice-1',
      answer: 'おはよう',
      is_correct: false,
    });
  });

  it('moves to next question after answer', async () => {
    render(<PatternPractice pattern={mockPattern} practices={mockPractices} />);

    // Answer first question
    const correctAnswer = screen.getByRole('button', { name: 'こんにちは' });
    await act(async () => {
      fireEvent.click(correctAnswer);
    });

    // Advance timer past the 1500ms setTimeout in the component
    await act(async () => {
      vi.advanceTimersByTime(1600);
    });

    // Check for next question
    expect(screen.getByText("Choose the correct response to 'Hello'")).toBeInTheDocument();
  });

  it('shows completion message when all questions answered', async () => {
    render(<PatternPractice pattern={mockPattern} practices={mockPractices} />);

    // Answer first question (correct answer)
    const firstAnswer = screen.getByRole('button', { name: 'こんにちは' });
    await act(async () => {
      fireEvent.click(firstAnswer);
    });
    await act(async () => {
      vi.advanceTimersByTime(1600);
    });

    // Answer second question (first option available)
    const secondQuestionAnswers = screen.getAllByRole('button');
    await act(async () => {
      fireEvent.click(secondQuestionAnswers[0]);
    });
    await act(async () => {
      vi.advanceTimersByTime(1600);
    });

    // Check completion message
    expect(screen.getByText(/completed/i)).toBeInTheDocument();
  });

  it('displays progress indicator', () => {
    render(<PatternPractice pattern={mockPattern} practices={mockPractices} />);

    // Component shows "Question 1 of 2"
    expect(screen.getByText(/Question 1 of 2/)).toBeInTheDocument();
  });

  it('shows difficulty level', () => {
    render(<PatternPractice pattern={mockPattern} practices={mockPractices} />);

    expect(screen.getByText(/difficulty/i)).toBeInTheDocument();
  });
});
