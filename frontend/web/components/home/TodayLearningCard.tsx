import type { Book, LearningProgress } from "@/lib/types"
import { cn } from "@/lib/utils"
import Link from "next/link"

interface TodayLearningCardProps {
  data: {
    book: Book
    progress: LearningProgress
  }
}

export function TodayLearningCard({ data }: TodayLearningCardProps) {
  const { book, progress } = data
  const progressPercentage = (progress.completedPages / progress.totalPages) * 100

  return (
    <div className={cn("rounded-xl border border-border bg-white p-6 shadow-sm")}>
      <div className="mb-4 flex items-center gap-2">
        <span className="text-2xl">📚</span>
        <h2 className="text-xl font-semibold text-text-primary">今日の学習</h2>
      </div>

      <div className="mb-4 space-y-2">
        <h3 className="text-lg font-medium text-text-primary">{book.title}</h3>

        <div className="space-y-2">
          <div className="flex items-center justify-between text-sm text-text-secondary">
            <span>
              {progress.completedPages}/{progress.totalPages} ページ完了
            </span>
            <span>{Math.round(progressPercentage)}%</span>
          </div>

          <div className="h-2 w-full overflow-hidden rounded-full bg-background-secondary">
            <div
              className="h-full bg-secondary transition-all duration-300"
              style={{ width: `${progressPercentage}%` }}
              role="progressbar"
              aria-valuenow={progressPercentage}
              aria-valuemin={0}
              aria-valuemax={100}
              tabIndex={0}
            />
          </div>
        </div>
      </div>

      <Link
        href={`/books/${book.id}/pages/${progress.currentPage}`}
        className={cn(
          "block w-full rounded-lg bg-primary px-4 py-3 text-center font-medium text-white",
          "transition-colors hover:bg-primary/90",
          "focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2"
        )}
      >
        続きから学習
      </Link>
    </div>
  )
}
