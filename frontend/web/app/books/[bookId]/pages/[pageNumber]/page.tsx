'use client';

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
  const pageNumber = parseInt(params.pageNumber, 10);

  // TODO: ユーザーIDを認証から取得
  const userId = 'user-123';

  const handlePageChange = (newPageNumber: number) => {
    router.push(`/books/${params.bookId}/pages/${newPageNumber}`);
  };

  return (
    <PageLearning
      bookId={params.bookId}
      pageNumber={pageNumber}
      userId={userId}
      onPageChange={handlePageChange}
    />
  );
}
