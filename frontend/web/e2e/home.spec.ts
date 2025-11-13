import { expect, test } from "@playwright/test"

test.describe("Home Page", () => {
  test.beforeEach(async ({ page }) => {
    // Mock the API responses
    await page.route("**/api/v1/home/dashboard", async (route) => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
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
            totalLearningTimeSeconds: 13320,
            completedPagesCount: 12,
            booksCount: 5,
            reviewItemsCount: 12,
          },
        }),
      })
    })

    await page.goto("/")
  })

  test("should display home page", async ({ page }) => {
    await expect(page).toHaveTitle(/HaiLanGo/)
  })

  test("should show welcome message", async ({ page }) => {
    await expect(page.getByText(/こんにちは、太郎さん/)).toBeVisible()
  })

  test("should display today's learning card", async ({ page }) => {
    await expect(page.getByText(/今日の学習/)).toBeVisible()
    await expect(page.getByText("ロシア語入門")).toBeVisible()
  })

  test("should show learning progress", async ({ page }) => {
    await expect(page.getByText(/12\/150/)).toBeVisible()
  })

  test("should have continue learning button", async ({ page }) => {
    const button = page.getByRole("button", { name: /続きから学習/ })
    await expect(button).toBeVisible()
  })

  test("should display quick access section", async ({ page }) => {
    await expect(page.getByText(/マイ本/)).toBeVisible()
    await expect(page.getByText(/復習/)).toBeVisible()
  })

  test("should show learning stats", async ({ page }) => {
    await expect(page.getByText(/学習統計/)).toBeVisible()
    await expect(page.getByText(/連続学習/)).toBeVisible()
    await expect(page.getByText("7日")).toBeVisible()
  })

  test("should navigate to book list when clicking My Books", async ({ page }) => {
    await page.click("text=マイ本")
    // Should navigate to /books (we'll implement this later)
    await expect(page).toHaveURL(/\/books/)
  })
})
