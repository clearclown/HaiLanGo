import { mockDashboardData } from "@/lib/__mocks__/data"
import { render, screen } from "@testing-library/react"
import { beforeEach, describe, expect, it, vi } from "vitest"
import HomePage from "./page"

// Mock the API module
vi.mock("@/lib/api", () => ({
  fetchDashboard: vi.fn(),
}))

describe("HomePage", () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it("should render welcome message with user name", async () => {
    const { fetchDashboard } = await import("@/lib/api")
    vi.mocked(fetchDashboard).mockResolvedValue(mockDashboardData)

    render(await HomePage())

    expect(screen.getByText(/こんにちは、太郎さん/)).toBeDefined()
    expect(screen.getByText(/今日も頑張りましょう！/)).toBeDefined()
  })

  it("should display today's learning card", async () => {
    const { fetchDashboard } = await import("@/lib/api")
    vi.mocked(fetchDashboard).mockResolvedValue(mockDashboardData)

    render(await HomePage())

    expect(screen.getByText(/今日の学習/)).toBeDefined()
    expect(screen.getByText(/ロシア語入門/)).toBeDefined()
  })

  it("should show learning progress", async () => {
    const { fetchDashboard } = await import("@/lib/api")
    vi.mocked(fetchDashboard).mockResolvedValue(mockDashboardData)

    render(await HomePage())

    // Check for progress text (12/150 pages completed)
    expect(screen.getByText(/12\/150/)).toBeDefined()
  })

  it("should display learning stats", async () => {
    const { fetchDashboard } = await import("@/lib/api")
    vi.mocked(fetchDashboard).mockResolvedValue(mockDashboardData)

    render(await HomePage())

    expect(screen.getByText(/連続学習/)).toBeDefined()
    expect(screen.getByText(/7日/)).toBeDefined()
    expect(screen.getByText(/総学習時間/)).toBeDefined()
  })

  it("should display quick access buttons", async () => {
    const { fetchDashboard } = await import("@/lib/api")
    vi.mocked(fetchDashboard).mockResolvedValue(mockDashboardData)

    render(await HomePage())

    expect(screen.getByText(/マイ本/)).toBeDefined()
    expect(screen.getByText(/復習/)).toBeDefined()
  })
})
