import { cn } from '@/lib/utils';
import Link from 'next/link';

interface QuickAccessProps {
  data: {
    booksCount: number;
    reviewItemsCount: number;
  };
}

export function QuickAccess({ data }: QuickAccessProps) {
  return (
    <div className="grid grid-cols-2 gap-4">
      <Link
        href="/books"
        className={cn(
          'flex flex-col items-center justify-center rounded-xl border border-border',
          'bg-white p-6 shadow-sm transition-all',
          'hover:border-primary hover:shadow-md',
          'focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2'
        )}
      >
        <span className="mb-2 text-3xl">📖</span>
        <h3 className="mb-1 text-lg font-medium text-text-primary">マイ本</h3>
        <p className="text-sm text-text-secondary">{data.booksCount}冊</p>
      </Link>

      <Link
        href="/review"
        className={cn(
          'flex flex-col items-center justify-center rounded-xl border border-border',
          'bg-white p-6 shadow-sm transition-all',
          'hover:border-primary hover:shadow-md',
          'focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2'
        )}
      >
        <span className="mb-2 text-3xl">🎯</span>
        <h3 className="mb-1 text-lg font-medium text-text-primary">復習</h3>
        <p className="text-sm text-text-secondary">{data.reviewItemsCount}項目</p>
      </Link>
    </div>
  );
}
