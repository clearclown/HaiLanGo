import type { DashboardData } from "../types"

export const mockDashboardData: DashboardData = {
  user: {
    id: "user-1",
    name: "太郎",
    email: "taro@example.com",
  },
  todayLearning: {
    book: {
      id: "book-1",
      title: "ロシア語入門",
      totalPages: 150,
      completedPages: 12,
      lastStudiedAt: "2025-11-13T08:00:00Z",
    },
    progress: {
      currentPage: 12,
      totalPages: 150,
      completedPages: 12,
    },
  },
  stats: {
    streakDays: 7,
    totalLearningTimeSeconds: 13320, // 3時間42分
    completedPagesCount: 12,
    booksCount: 5,
    reviewItemsCount: 12,
  },
}
