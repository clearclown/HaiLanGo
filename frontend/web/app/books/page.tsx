import { BooksList } from '@/components/books/BooksList';
import Link from 'next/link';

export default function BooksPage() {
  return (
    <div className="min-h-screen bg-background-secondary">
      <div className="max-w-6xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold">マイ本</h1>
            <p className="text-gray-600 mt-1">あなたの学習教材</p>
          </div>
          <Link
            href="/upload"
            className="flex items-center gap-2 px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
          >
            <span>＋</span>
            <span>本を追加</span>
          </Link>
        </div>

        {/* Search Bar */}
        <div className="mb-6">
          <input
            type="text"
            placeholder="本を検索..."
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        {/* Books List */}
        <BooksList />
      </div>
    </div>
  );
}
