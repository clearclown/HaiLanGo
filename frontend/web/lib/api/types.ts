// API型定義

export interface Page {
  id: string;
  bookId: string;
  pageNumber: number;
  imageUrl: string;
  ocrText: string;
  translation: string;
  audioUrl: string;
  createdAt: string;
  updatedAt: string;
}

export interface PageWithProgress extends Page {
  isCompleted: boolean;
  completedAt?: string;
}

export interface LearningProgress {
  bookId: string;
  totalPages: number;
  completedPages: number;
  progress: number;
  totalStudyTime: number;
}

export interface MarkPageCompletedRequest {
  userId: string;
  studyTime: number;
}
