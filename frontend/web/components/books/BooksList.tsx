'use client';

import { useEffect, useState } from 'react';
import { apiClient } from '@/lib/api/client';
import type { Book } from '@/types/book';
import { BookCard } from './BookCard';

export function BooksList() {
  const [books, setBooks] = useState<Book[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [bookToDelete, setBookToDelete] = useState<string | null>(null);

  useEffect(() => {
    loadBooks();
  }, []);

  const loadBooks = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await apiClient.books.list();
      setBooks(response.books);
    } catch (err) {
      console.error('Failed to load books:', err);
      setError('本の読み込みに失敗しました');
    } finally {
      setIsLoading(false);
    }
  };

  const handleDeleteClick = (bookId: string) => {
    setBookToDelete(bookId);
    setShowDeleteDialog(true);
  };

  const handleDeleteConfirm = async () => {
    if (!bookToDelete) return;

    try {
      await apiClient.books.delete(bookToDelete);
      setBooks(books.filter(book => book.id !== bookToDelete));
      setShowDeleteDialog(false);
      setBookToDelete(null);
    } catch (err) {
      console.error('Failed to delete book:', err);
      alert('本の削除に失敗しました');
    }
  };

  const handleDeleteCancel = () => {
    setShowDeleteDialog(false);
    setBookToDelete(null);
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-gray-600">読み込み中...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <div className="text-red-600 mb-4">{error}</div>
        <button
          type="button"
          onClick={loadBooks}
          className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
        >
          再試行
        </button>
      </div>
    );
  }

  if (books.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <div className="text-6xl mb-4">📚</div>
        <h3 className="text-xl font-semibold mb-2">まだ本がありません</h3>
        <p className="text-gray-600 mb-6">新しい本を追加して学習を始めましょう！</p>
        <a
          href="/upload"
          className="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
        >
          本を追加
        </a>
      </div>
    );
  }

  return (
    <>
      <div className="space-y-4">
        {books.map(book => (
          <BookCard
            key={book.id}
            book={book}
            onDelete={handleDeleteClick}
          />
        ))}
      </div>

      {/* Delete Confirmation Dialog */}
      {showDeleteDialog && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-sm mx-4">
            <h3 className="text-lg font-semibold mb-4">本を削除しますか？</h3>
            <p className="text-gray-600 mb-6">
              学習記録も削除されます。この操作は取り消せません。
            </p>
            <div className="flex gap-4">
              <button
                type="button"
                onClick={handleDeleteCancel}
                className="flex-1 bg-gray-200 text-gray-700 px-4 py-2 rounded-lg hover:bg-gray-300"
              >
                キャンセル
              </button>
              <button
                type="button"
                onClick={handleDeleteConfirm}
                className="flex-1 bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700"
              >
                削除
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
