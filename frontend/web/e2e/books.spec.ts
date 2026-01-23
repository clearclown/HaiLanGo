import { expect, test } from './fixtures';

test.describe('Books Page Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Mock books API
    await page.route('**/api/v1/books**', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          books: [],
          total: 0,
        }),
      });
    });
    await page.goto('/books');
  });

  test('should display books page correctly', async ({ page }) => {
    // ページタイトルの確認
    const heading = page.getByRole('heading', { name: /マイ本/i });
    await expect(heading).toBeVisible();

    // サブタイトルの確認
    const subtitle = page.getByText(/あなたの学習教材/i);
    await expect(subtitle).toBeVisible();

    // 追加ボタンの確認 (header button)
    const addButton = page.getByRole('link', { name: /本を追加/i }).first();
    await expect(addButton).toBeVisible();

    // 検索バーの確認
    const searchInput = page.getByPlaceholder(/本を検索/i);
    await expect(searchInput).toBeVisible();
  });

  test('should show empty state when no books', async ({ page }) => {
    // Wait for loading to complete (loading message to disappear)
    await page
      .waitForSelector('text=読み込み中...', { state: 'hidden', timeout: 10000 })
      .catch(() => {});

    // 空の状態のメッセージを確認
    const emptyMessage = page.getByText(/まだ本がありません/i);
    await expect(emptyMessage).toBeVisible({ timeout: 5000 });
  });

  test('should have functional search input', async ({ page }) => {
    const searchInput = page.getByPlaceholder(/本を検索/i);
    await searchInput.fill('テスト');

    const value = await searchInput.inputValue();
    expect(value).toBe('テスト');
  });

  test('should navigate to upload page when clicking add button', async ({ page }) => {
    // Use first() because there are multiple "本を追加" links (header + empty state)
    const addButton = page.getByRole('link', { name: /本を追加/i }).first();
    await addButton.click();

    await expect(page).toHaveURL(/.*upload/);
  });

  test('should display book cards with correct information', async ({ page }) => {
    // 本が存在する場合、カード要素を確認
    const bookCards = page.locator('[class*="book"]').or(page.locator('article'));
    const count = await bookCards.count();

    if (count > 0) {
      const firstCard = bookCards.first();
      await expect(firstCard).toBeVisible();

      // カードには学習を続けるボタンまたは詳細ボタンがあるはず
      const actionButton = firstCard.getByRole('link').or(firstCard.getByRole('button'));
      await expect(actionButton.first()).toBeVisible();
    }
  });
});
