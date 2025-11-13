export interface User {
  id: string
  name: string
  email: string
}

export interface Book {
  id: string
  title: string
  coverImage?: string
  totalPages: number
  completedPages: number
  lastStudiedAt?: string
}

export interface LearningProgress {
  currentPage: number
  totalPages: number
  completedPages: number
}

export interface LearningStats {
  streakDays: number
  totalLearningTimeSeconds: number
  completedPagesCount: number
  booksCount: number
  reviewItemsCount: number
}

export interface DashboardData {
  user: User
  todayLearning?: {
    book: Book
    progress: LearningProgress
  }
  stats: LearningStats
}
