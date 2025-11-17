import { LearningStats } from "@/components/home/LearningStats"
import { QuickAccess } from "@/components/home/QuickAccess"
import { TodayLearningCard } from "@/components/home/TodayLearningCard"
import { WelcomeCard } from "@/components/home/WelcomeCard"
import { fetchDashboard } from "@/lib/api"

// Force dynamic rendering (no static generation at build time)
export const dynamic = 'force-dynamic'

export default async function HomePage() {
  const data = await fetchDashboard()

  return (
    <div className="min-h-screen bg-background-secondary">
      <div className="mx-auto max-w-4xl px-4 py-8">
        <WelcomeCard userName={data.user.name} />

        <div className="space-y-6">
          {data.todayLearning && <TodayLearningCard data={data.todayLearning} />}

          <QuickAccess
            data={{
              booksCount: data.stats.booksCount,
              reviewItemsCount: data.stats.reviewItemsCount,
            }}
          />

          <LearningStats stats={data.stats} />
        </div>
      </div>
    </div>
  )
}
