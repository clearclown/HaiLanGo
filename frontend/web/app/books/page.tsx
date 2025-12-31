import { BooksList } from '@/components/books/BooksList';
import { AppLayout } from '@/components/layout';
import { Button, Input } from '@/components/ui';
import Link from 'next/link';

export default function BooksPage() {
  return (
    <AppLayout>
      <div className="container-app py-6 lg:py-8">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 mb-6">
          <div>
            <h1 className="text-2xl lg:text-3xl font-bold text-gray-900">マイ本</h1>
            <p className="text-gray-600 mt-1">あなたの学習教材</p>
          </div>
          <Link href="/upload">
            <Button variant="primary" size="md" className="w-full sm:w-auto">
              <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 4v16m8-8H4"
                />
              </svg>
              本を追加
            </Button>
          </Link>
        </div>

        {/* Search Bar */}
        <div className="mb-6">
          <Input
            type="text"
            placeholder="本を検索..."
            leftIcon={
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
            }
          />
        </div>

        {/* Books List */}
        <BooksList />
      </div>
    </AppLayout>
  );
}
