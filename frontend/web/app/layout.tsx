import type { Metadata } from "next"
import "./globals.css"

export const metadata: Metadata = {
  title: "HaiLanGo - AI-Powered Language Learning",
  description: "AI-powered language learning platform",
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  )
}
