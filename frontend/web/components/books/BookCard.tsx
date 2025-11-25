import type { Book } from '@/types/book';
import Link from 'next/link';

interface BookCardProps {
  book: Book;
  onDelete?: (bookId: string) => void;
}

export function BookCard({ book, onDelete }: BookCardProps) {
  const progressPercentage = book.total_pages > 0
    ? Math.round((book.processed_pages / book.total_pages) * 100)
    : 0;

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
        return 'text-green-600 bg-green-50';
      case 'processing':
      case 'uploading':
        return 'text-blue-600 bg-blue-50';
      case 'failed':
        return 'text-red-600 bg-red-50';
      default:
        return 'text-gray-600 bg-gray-50';
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow">
      <div className="flex gap-4">
        {/* Book Cover */}
        <div className="flex-shrink-0">
          {book.cover_image_url ? (
            <img
              src={book.cover_image_url}
              alt={book.title}
              className="w-24 h-32 object-cover rounded"
            />
          ) : (
            <div className="w-24 h-32 bg-gray-200 rounded flex items-center justify-center text-gray-400">
              📕
            </div>
          )}
        </div>

        {/* Book Info */}
        <div className="flex-1 min-w-0">
          <h3 className="text-xl font-semibold mb-2 truncate">{book.title}</h3>

          <div className="space-y-1 text-sm text-gray-600 mb-3">
            <p>
              <span className="font-medium">学習言語:</span> {book.target_language}
            </p>
            <p>
              <span className="font-medium">母国語:</span> {book.native_language}
            </p>
            {book.reference_language && (
              <p>
                <span className="font-medium">参照言語:</span> {book.reference_language}
              </p>
            )}
          </div>

          {/* Status Badge */}
          <div className="mb-3">
            <span className={`inline-block px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(book.status)}`}>
              {getStatusText(book.status)}
            </span>
          </div>

          {/* Progress Bar */}
          {book.status === 'ready' && (
            <div className="mb-3">
              <div className="flex justify-between text-sm text-gray-600 mb-1">
                <span>進捗</span>
                <span>{progressPercentage}% ({book.processed_pages}/{book.total_pages}ページ)</span>
              </div>
              <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
                <div
                  className="h-full bg-blue-500"
                  style={{ width: `${progressPercentage}%` }}
                />
              </div>
            </div>
          )}

          {/* Actions */}
          <div className="flex gap-2">
            {book.status === 'ready' && (
              <Link
                href={`/books/${book.id}/pages/1`}
                className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors text-sm"
              >
                学習を続ける
              </Link>
            )}
            <Link
              href={`/books/${book.id}`}
              className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors text-sm"
            >
              詳細
            </Link>
            {onDelete && (
              <button
                type="button"
                onClick={() => onDelete(book.id)}
                className="px-4 py-2 bg-red-50 text-red-600 rounded-lg hover:bg-red-100 transition-colors text-sm"
              >
                削除
              </button>
            )}
          </div>
        </div>
      </div>

      {/* Updated At */}
      <div className="mt-4 text-xs text-gray-500">
        最終更新: {new Date(book.updated_at).toLocaleString('ja-JP')}
      </div>
    </div>
  );
}
