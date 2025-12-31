'use client';

import { apiClient } from '@/lib/api/client';
import type { Book } from '@/types/book';
import Link from 'next/link';
import { useParams, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

export default function BookDetailPage() {
  const params = useParams();
  const router = useRouter();
  const bookId = params.bookId as string;

  const [book, setBook] = useState<Book | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    const fetchBook = async () => {
      try {
        const bookData = await apiClient.books.get(bookId);
        setBook(bookData);
      } catch (err) {
        setError('本の情報を取得できませんでした');
        console.error('Failed to fetch book:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchBook();
  }, [bookId]);

  const handleDelete = async () => {
    if (!confirm('この本を削除してもよろしいですか？この操作は取り消せません。')) {
      return;
    }

    setDeleting(true);
    try {
      await apiClient.books.delete(bookId);
      router.push('/books');
    } catch (err) {
      console.error('Failed to delete book:', err);
      alert('本の削除に失敗しました');
      setDeleting(false);
    }
  };

  const getStatusText = (status: Book['status']) => {
    switch (status) {
      case 'uploading':
        return 'アップロード中';
      case 'processing':
        return 'OCR処理中';
      case 'ready':
        return '学習可能';
      case 'failed':
        return '処理失敗';
      default:
        return status;
    }
  };

  const getStatusColor = (status: Book['status']) => {
    switch (status) {
      case 'ready':
        return 'text-green-600 bg-green-100';
      case 'processing':
      case 'uploading':
        return 'text-blue-600 bg-blue-100';
      case 'failed':
        return 'text-red-600 bg-red-100';
      default:
        return 'text-gray-600 bg-gray-100';
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
      </div>
    );
  }

  if (error || !book) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-600 mb-4">{error || '本が見つかりませんでした'}</p>
          <Link href="/books" className="text-blue-600 hover:underline">
            本棚に戻る
          </Link>
        </div>
      </div>
    );
  }

  const progressPercentage =
    book.total_pages > 0 ? Math.round((book.processed_pages / book.total_pages) * 100) : 0;

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-4xl mx-auto px-4 py-8">
        {/* Back Button */}
        <Link
          href="/books"
          className="inline-flex items-center text-gray-600 hover:text-gray-800 mb-6"
        >
          ← 本棚に戻る
        </Link>

        <div className="bg-white rounded-lg shadow-sm p-8">
          <div className="flex gap-8">
            {/* Book Cover */}
            <div className="flex-shrink-0">
              {book.cover_image_url ? (
                <img
                  src={book.cover_image_url}
                  alt={book.title}
                  className="w-48 h-64 object-cover rounded-lg shadow"
                />
              ) : (
                <div className="w-48 h-64 bg-gray-200 rounded-lg flex items-center justify-center text-6xl text-gray-400">
                  📕
                </div>
              )}
            </div>

            {/* Book Info */}
            <div className="flex-1">
              <h1 className="text-3xl font-bold mb-4">{book.title}</h1>

              <div className="mb-6">
                <span
                  className={`inline-block px-4 py-2 rounded-full text-sm font-medium ${getStatusColor(book.status)}`}
                >
                  {getStatusText(book.status)}
                </span>
              </div>

              <div className="space-y-3 text-gray-600 mb-6">
                <div className="flex">
                  <span className="w-32 font-medium">学習先言語:</span>
                  <span>{book.target_language}</span>
                </div>
                <div className="flex">
                  <span className="w-32 font-medium">母国語:</span>
                  <span>{book.native_language}</span>
                </div>
                {book.reference_language && (
                  <div className="flex">
                    <span className="w-32 font-medium">参照言語:</span>
                    <span>{book.reference_language}</span>
                  </div>
                )}
                <div className="flex">
                  <span className="w-32 font-medium">総ページ数:</span>
                  <span>{book.total_pages}ページ</span>
                </div>
                <div className="flex">
                  <span className="w-32 font-medium">作成日:</span>
                  <span>{new Date(book.created_at).toLocaleDateString('ja-JP')}</span>
                </div>
                <div className="flex">
                  <span className="w-32 font-medium">最終更新:</span>
                  <span>{new Date(book.updated_at).toLocaleString('ja-JP')}</span>
                </div>
              </div>

              {/* Progress */}
              {book.status === 'ready' && book.total_pages > 0 && (
                <div className="mb-6">
                  <div className="flex justify-between text-sm text-gray-600 mb-2">
                    <span>学習進捗</span>
                    <span>
                      {progressPercentage}% ({book.processed_pages}/{book.total_pages}ページ)
                    </span>
                  </div>
                  <div className="h-3 bg-gray-200 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-blue-500 transition-all"
                      style={{ width: `${progressPercentage}%` }}
                    />
                  </div>
                </div>
              )}

              {/* Action Buttons */}
              <div className="flex gap-4">
                {book.status === 'ready' && (
                  <Link
                    href={`/books/${book.id}/pages/1`}
                    className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium"
                  >
                    学習を開始する
                  </Link>
                )}
                <button
                  type="button"
                  onClick={handleDelete}
                  disabled={deleting}
                  className="px-6 py-3 bg-red-100 text-red-600 rounded-lg hover:bg-red-200 transition-colors font-medium disabled:opacity-50"
                >
                  {deleting ? '削除中...' : '本を削除'}
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Pages List (if ready) */}
        {book.status === 'ready' && book.total_pages > 0 && (
          <div className="mt-8 bg-white rounded-lg shadow-sm p-8">
            <h2 className="text-xl font-bold mb-6">ページ一覧</h2>
            <div className="grid grid-cols-5 sm:grid-cols-8 md:grid-cols-10 gap-2">
              {Array.from({ length: book.total_pages }, (_, i) => i + 1).map((pageNum) => (
                <Link
                  key={pageNum}
                  href={`/books/${book.id}/pages/${pageNum}`}
                  className="aspect-square flex items-center justify-center bg-gray-100 hover:bg-blue-100 text-gray-700 hover:text-blue-600 rounded-lg transition-colors text-sm font-medium"
                >
                  {pageNum}
                </Link>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
