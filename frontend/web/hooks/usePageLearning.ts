import { useState, useEffect } from 'react';
import type { PageWithProgress } from '@/lib/api/types';
import { getPage, markPageCompleted } from '@/services/learningApi';

interface UsePageLearningProps {
  bookId: string;
  pageNumber: number;
  userId: string;
}

interface UsePageLearningReturn {
  page: PageWithProgress | null;
  loading: boolean;
  error: Error | null;
  markCompleted: (studyTime: number) => Promise<void>;
  refetch: () => Promise<void>;
}

/**
 * ページ学習のカスタムフック
 */
export function usePageLearning({
  bookId,
  pageNumber,
  userId,
}: UsePageLearningProps): UsePageLearningReturn {
  const [page, setPage] = useState<PageWithProgress | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const fetchPage = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await getPage(bookId, pageNumber);
      setPage(data);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPage();
  }, [bookId, pageNumber]);

  const markCompleted = async (studyTime: number) => {
    try {
      await markPageCompleted(bookId, pageNumber, { userId, studyTime });
      // ページ情報を再取得
      await fetchPage();
    } catch (err) {
      throw err instanceof Error ? err : new Error('Failed to mark page as completed');
    }
  };

  return {
    page,
    loading,
    error,
    markCompleted,
    refetch: fetchPage,
  };
}
