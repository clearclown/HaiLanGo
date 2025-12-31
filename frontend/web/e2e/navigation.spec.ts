import { expect, test } from './fixtures';

test.describe('Navigation Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Mock common API endpoints
    await page.route('**/api/v1/home/dashboard', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          user: { id: '1', name: '太郎', email: 'test@example.com' },
          todayLearning: null,
          stats: {
            streakDays: 0,
            totalLearningTimeSeconds: 0,
            completedPagesCount: 0,
            booksCount: 0,
            reviewItemsCount: 0,
          },
        }),
      });
    });
    await page.route('**/api/v1/books**', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ books: [], total: 0 }),
      });
    });
    await page.route('**/api/v1/settings', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          profile: { id: '1', name: '太郎', email: 'test@example.com' },
          notifications: {
            learningReminder: true,
            reviewNotification: true,
            emailNotification: false,
          },
          interfaceLanguage: 'ja',
        }),
      });
    });
    await page.route('**/api/v1/plan', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ type: 'free' }),
      });
    });
    await page.route('**/api/v1/review/**', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ items: [], stats: { urgent: 0, recommended: 0, optional: 0 } }),
      });
    });
  });

  test('should display home page at root', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveTitle(/HaiLanGo/i);
    // Home page should show welcome message
    await expect(page.getByText(/こんにちは/)).toBeVisible();
  });

  test('should have functional navigation links', async ({ page }) => {
    await page.goto('/books');

    // ページが正しくロードされることを確認
    await expect(page).toHaveTitle(/HaiLanGo/i);

    // ヘッダーまたはナビゲーション要素が存在することを確認
    const heading = page.getByRole('heading', { name: /マイ本/i });
    await expect(heading).toBeVisible();
  });

  test('should navigate to upload page', async ({ page }) => {
    await page.goto('/books');

    // 「本を追加」リンクをクリック
    const uploadLink = page.getByRole('link', { name: /本を追加/i });
    if (await uploadLink.isVisible()) {
      await uploadLink.click();
      await expect(page).toHaveURL(/.*upload/);
    }
  });

  test('should navigate to settings page', async ({ page }) => {
    await page.goto('/settings');
    await expect(page).toHaveURL(/.*settings/);

    const heading = page.getByRole('heading', { name: '設定', exact: true });
    await expect(heading).toBeVisible();
  });

  test('should navigate to review page', async ({ page }) => {
    await page.goto('/review');
    await expect(page).toHaveURL(/.*review/);

    // 復習ページの要素を確認
    const heading = page.getByRole('heading', { name: /復習/i });
    await expect(heading).toBeVisible();
  });
});
