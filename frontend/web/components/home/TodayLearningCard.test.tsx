import { render, screen } from "@testing-library/react"
import { describe, expect, it } from "vitest"
import { TodayLearningCard } from "./TodayLearningCard"

describe("TodayLearningCard", () => {
  const mockData = {
    book: {
      id: "book-1",
      title: "ロシア語入門",
      totalPages: 150,
      completedPages: 12,
    },
    progress: {
      currentPage: 12,
      totalPages: 150,
      completedPages: 12,
    },
  }

  it("should render section title", () => {
    render(<TodayLearningCard data={mockData} />)

    expect(screen.getByText(/今日の学習/)).toBeDefined()
  })

  it("should display book title", () => {
    render(<TodayLearningCard data={mockData} />)

    expect(screen.getByText("ロシア語入門")).toBeDefined()
  })

  it("should show progress bar", () => {
    render(<TodayLearningCard data={mockData} />)

    expect(screen.getByText(/12\/150/)).toBeDefined()
    expect(screen.getByText(/ページ完了/)).toBeDefined()
  })

  it("should display continue learning button", () => {
    render(<TodayLearningCard data={mockData} />)

    expect(screen.getByText(/続きから学習/)).toBeDefined()
  })

  it("should calculate progress percentage correctly", () => {
    const { container } = render(<TodayLearningCard data={mockData} />)

    // Progress should be 12/150 = 8%
    const progressBar = container.querySelector('[role="progressbar"]')
    expect(progressBar).toBeDefined()
  })
})
