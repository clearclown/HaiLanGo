import { cn } from "@/lib/utils"

interface WelcomeCardProps {
  userName: string
}

export function WelcomeCard({ userName }: WelcomeCardProps) {
  return (
    <div className={cn("mb-6 space-y-2")}>
      <h1 className="text-2xl font-bold text-text-primary">👋 こんにちは、{userName}さん</h1>
      <p className="text-text-secondary">今日も頑張りましょう！</p>
    </div>
  )
}
