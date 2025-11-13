import type { DashboardData, LearningStats } from "./types"

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"

export async function fetchDashboard(): Promise<DashboardData> {
  const response = await fetch(`${API_BASE_URL}/api/v1/home/dashboard`, {
    credentials: "include",
  })

  if (!response.ok) {
    throw new Error("Failed to fetch dashboard data")
  }

  return response.json()
}

export async function fetchStats(): Promise<LearningStats> {
  const response = await fetch(`${API_BASE_URL}/api/v1/home/stats`, {
    credentials: "include",
  })

  if (!response.ok) {
    throw new Error("Failed to fetch stats")
  }

  return response.json()
}
