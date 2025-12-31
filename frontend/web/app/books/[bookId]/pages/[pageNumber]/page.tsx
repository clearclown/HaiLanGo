'use client';

import { useAuth } from '@/components/AuthProvider';
import { PageLearning } from '@/components/learning/PageLearning';
import { useRouter } from 'next/navigation';

interface PageProps {
  params: {
    bookId: string;
    pageNumber: string;
  };
}

export default function LearningPage({ params }: PageProps) {
  const router = useRouter();
  const { user, isLoading } = useAuth();
  const pageNumber = Number.parseInt(params.pageNumber, 10);

  const handlePageChange = (newPageNumber: number) => {
    router.push(`/books/${params.bookId}/pages/${newPageNumber}`);
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
      </div>
    );
  }

  if (!user) {
    return null; // AuthProvider will redirect to login
  }

  return (
    <PageLearning
      bookId={params.bookId}
      pageNumber={pageNumber}
      userId={user.id}
      onPageChange={handlePageChange}
    />
  );
}
