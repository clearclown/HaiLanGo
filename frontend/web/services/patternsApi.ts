/**
 * Patterns API Service
 */

import type { Pattern, PatternPractice } from '@/lib/types/pattern';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export class PatternsApiError extends Error {
  constructor(
    public status: number,
    message: string
  ) {
    super(message);
    this.name = 'PatternsApiError';
  }
}

/**
 * Get patterns for a book
 */
export async function getPatterns(bookId: string): Promise<Pattern[]> {
  const response = await fetch(`${API_BASE_URL}/api/v1/books/${bookId}/patterns`);

  if (!response.ok) {
    throw new PatternsApiError(response.status, 'Failed to fetch patterns');
  }

  const data = await response.json();
  return data.patterns;
}

/**
 * Get practice questions for a pattern
 */
export async function getPatternPractice(patternId: string): Promise<PatternPractice[]> {
  const response = await fetch(`${API_BASE_URL}/api/v1/patterns/${patternId}/practice`);

  if (!response.ok) {
    throw new PatternsApiError(response.status, 'Failed to fetch practice questions');
  }

  const data = await response.json();
  return data.questions;
}

/**
 * Submit practice answer
 */
export async function submitPracticeAnswer(
  patternId: string,
  questionId: string,
  answer: string
): Promise<{ correct: boolean; correctAnswer: string }> {
  const response = await fetch(`${API_BASE_URL}/api/v1/patterns/${patternId}/practice/answer`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ questionId, answer }),
  });

  if (!response.ok) {
    throw new PatternsApiError(response.status, 'Failed to submit answer');
  }

  return response.json();
}
