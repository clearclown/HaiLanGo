import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { PageLearning } from './PageLearning';
import * as learningApi from '@/services/learningApi';

// モックの設定
vi.mock('@/services/learningApi');

const mockPage = {
  id: '123',
  bookId: 'book-123',
  pageNumber: 1,
  imageUrl: 'https://example.com/page1.png',
  ocrText: 'Здравствуйте!',
  translation: 'こんにちは！',
  audioUrl: 'https://example.com/audio1.mp3',
  createdAt: '2025-01-01T00:00:00Z',
  updatedAt: '2025-01-01T00:00:00Z',
  isCompleted: false,
};

describe('PageLearning', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should render page content', async () => {
    vi.mocked(learningApi.getPage).mockResolvedValue(mockPage);

    render(<PageLearning bookId="book-123" pageNumber={1} userId="user-123" />);

    await waitFor(() => {
      expect(screen.getByText('Здравствуйте!')).toBeInTheDocument();
    });
  });

  it('should show loading state', () => {
    vi.mocked(learningApi.getPage).mockImplementation(
      () => new Promise(() => {}) // 永遠に解決しないPromise
    );

    render(<PageLearning bookId="book-123" pageNumber={1} userId="user-123" />);

    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });

  it('should show error state', async () => {
    vi.mocked(learningApi.getPage).mockRejectedValue(new Error('Failed to fetch'));

    render(<PageLearning bookId="book-123" pageNumber={1} userId="user-123" />);

    await waitFor(() => {
      expect(screen.getByText(/error/i)).toBeInTheDocument();
    });
  });

  it('should mark page as completed', async () => {
    vi.mocked(learningApi.getPage).mockResolvedValue(mockPage);
    vi.mocked(learningApi.markPageCompleted).mockResolvedValue(undefined);

    render(<PageLearning bookId="book-123" pageNumber={1} userId="user-123" />);

    await waitFor(() => {
      expect(screen.getByText('Здравствуйте!')).toBeInTheDocument();
    });

    const completeButton = screen.getByRole('button', { name: /学習完了/i });
    fireEvent.click(completeButton);

    await waitFor(() => {
      expect(learningApi.markPageCompleted).toHaveBeenCalledWith('book-123', 1, {
        userId: 'user-123',
        studyTime: expect.any(Number),
      });
    });
  });

  it('should navigate to next page', async () => {
    vi.mocked(learningApi.getPage).mockResolvedValue(mockPage);

    const onPageChange = vi.fn();

    render(
      <PageLearning
        bookId="book-123"
        pageNumber={1}
        userId="user-123"
        onPageChange={onPageChange}
      />
    );

    await waitFor(() => {
      expect(screen.getByText('Здравствуйте!')).toBeInTheDocument();
    });

    const nextButton = screen.getByRole('button', { name: /次へ/i });
    fireEvent.click(nextButton);

    expect(onPageChange).toHaveBeenCalledWith(2);
  });

  it('should navigate to previous page', async () => {
    vi.mocked(learningApi.getPage).mockResolvedValue({ ...mockPage, pageNumber: 2 });

    const onPageChange = vi.fn();

    render(
      <PageLearning
        bookId="book-123"
        pageNumber={2}
        userId="user-123"
        onPageChange={onPageChange}
      />
    );

    await waitFor(() => {
      expect(screen.getByText('Здравствуйте!')).toBeInTheDocument();
    });

    const prevButton = screen.getByRole('button', { name: /前へ/i });
    fireEvent.click(prevButton);

    expect(onPageChange).toHaveBeenCalledWith(1);
  });

  it('should disable previous button on first page', async () => {
    vi.mocked(learningApi.getPage).mockResolvedValue(mockPage);

    render(<PageLearning bookId="book-123" pageNumber={1} userId="user-123" />);

    await waitFor(() => {
      expect(screen.getByText('Здравствуйте!')).toBeInTheDocument();
    });

    const prevButton = screen.getByRole('button', { name: /前へ/i });
    expect(prevButton).toBeDisabled();
  });
});
