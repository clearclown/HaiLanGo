'use client';

import { TeacherMode } from '@/components/teacher-mode/TeacherMode';

interface PageProps {
  params: {
    bookId: string;
  };
}

export default function TeacherModePage({ params }: PageProps) {
  return <TeacherMode bookId={params.bookId} />;
}
