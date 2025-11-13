import React, { useState, useEffect } from 'react';
import { usePageLearning } from '@/hooks/usePageLearning';
import { AudioPlayer } from './AudioPlayer';

interface PageLearningProps {
  bookId: string;
  pageNumber: number;
  userId: string;
  onPageChange?: (newPageNumber: number) => void;
}

/**
 * ページバイページ学習コンポーネント
 */
export const PageLearning: React.FC<PageLearningProps> = ({
  bookId,
  pageNumber,
  userId,
  onPageChange,
}) => {
  const { page, loading, error, markCompleted } = usePageLearning({
    bookId,
    pageNumber,
    userId,
  });
  const [studyStartTime, setStudyStartTime] = useState<number>(Date.now());

  useEffect(() => {
    // ページが変わるたびに学習開始時刻をリセット
    setStudyStartTime(Date.now());
  }, [pageNumber]);

  const handleMarkCompleted = async () => {
    const studyTime = Math.floor((Date.now() - studyStartTime) / 1000);
    await markCompleted(studyTime);
  };

  const handlePrevPage = () => {
    if (pageNumber > 1 && onPageChange) {
      onPageChange(pageNumber - 1);
    }
  };

  const handleNextPage = () => {
    if (onPageChange) {
      onPageChange(pageNumber + 1);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-lg">Loading...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-red-500">
          <p>Error: {error.message}</p>
        </div>
      </div>
    );
  }

  if (!page) {
    return null;
  }

  return (
    <div className="flex flex-col h-screen bg-gray-50">
      {/* ヘッダー */}
      <header className="bg-white shadow-sm p-4">
        <div className="flex items-center justify-between">
          <h1 className="text-xl font-bold">ページ {page.pageNumber}</h1>
          {page.isCompleted && (
            <span className="px-3 py-1 bg-green-100 text-green-800 rounded-full text-sm">
              完了済み
            </span>
          )}
        </div>
        {/* 進捗バー */}
        <div className="mt-2 h-1 bg-gray-200 rounded">
          <div
            className="h-full bg-blue-500 rounded"
            style={{ width: `${(page.pageNumber / 150) * 100}%` }}
          />
        </div>
      </header>

      {/* メインコンテンツ */}
      <main className="flex-1 overflow-auto p-4">
        {/* ページ画像 */}
        <div className="mb-4">
          <img
            src={page.imageUrl}
            alt={`ページ ${page.pageNumber}`}
            className="w-full max-w-2xl mx-auto rounded-lg shadow-md"
            data-testid="page-content"
          />
        </div>

        {/* OCRテキスト */}
        <div className="max-w-2xl mx-auto mb-4 p-4 bg-white rounded-lg shadow">
          <h2 className="text-lg font-semibold mb-2">テキスト</h2>
          <p className="text-2xl mb-4">{page.ocrText}</p>
          <p className="text-gray-600">{page.translation}</p>
        </div>

        {/* 音声プレイヤー */}
        <div className="max-w-2xl mx-auto mb-4">
          <AudioPlayer audioUrl={page.audioUrl} />
        </div>
      </main>

      {/* フッター */}
      <footer className="bg-white shadow-lg p-4">
        <div className="max-w-2xl mx-auto flex justify-between items-center">
          <button
            onClick={handlePrevPage}
            disabled={pageNumber === 1}
            className="px-6 py-2 bg-gray-200 text-gray-700 rounded-lg disabled:opacity-50"
            aria-label="前へ"
          >
            ← 前へ
          </button>

          <button
            onClick={handleMarkCompleted}
            className="px-6 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600"
            aria-label="学習完了"
          >
            学習完了 ✓
          </button>

          <button
            onClick={handleNextPage}
            className="px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
            aria-label="次へ"
          >
            次へ →
          </button>
        </div>
      </footer>
    </div>
  );
};
