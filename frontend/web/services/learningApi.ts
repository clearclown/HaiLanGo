import type { PageWithProgress, LearningProgress, MarkPageCompletedRequest } from '@/lib/api/types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export class LearningApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'LearningApiError';
  }
}

/**
 * ページを取得する
 */
export async function getPage(
  bookId: string,
  pageNumber: number
): Promise<PageWithProgress> {
  const response = await fetch(
    `${API_BASE_URL}/api/v1/books/${bookId}/pages/${pageNumber}`
  );

  if (!response.ok) {
    throw new LearningApiError(response.status, 'Failed to fetch page');
  }

  return response.json();
}

/**
 * ページを完了としてマークする
 */
export async function markPageCompleted(
  bookId: string,
  pageNumber: number,
  data: MarkPageCompletedRequest
): Promise<void> {
  const response = await fetch(
    `${API_BASE_URL}/api/v1/books/${bookId}/pages/${pageNumber}/complete`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    }
  );

  if (!response.ok) {
    throw new LearningApiError(response.status, 'Failed to mark page as completed');
  }
}

/**
 * 学習進捗を取得する
 */
export async function getProgress(
  bookId: string,
  userId: string
): Promise<LearningProgress> {
  const response = await fetch(
    `${API_BASE_URL}/api/v1/books/${bookId}/progress?userId=${userId}`
  );

  if (!response.ok) {
    throw new LearningApiError(response.status, 'Failed to fetch progress');
  }

  return response.json();
}
