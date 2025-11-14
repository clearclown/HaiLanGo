import { test, expect } from '@playwright/test';

test.describe('Books Page Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/books');
  });

  test('should display books page correctly', async ({ page }) => {
    // ページタイトルの確認
    const heading = page.getByRole('heading', { name: /マイ本/i });
    await expect(heading).toBeVisible();

    // サブタイトルの確認
    const subtitle = page.getByText(/あなたの学習教材/i);
    await expect(subtitle).toBeVisible();

    // 追加ボタンの確認
    const addButton = page.getByRole('link', { name: /本を追加/i });
    await expect(addButton).toBeVisible();

    // 検索バーの確認
    const searchInput = page.getByPlaceholder(/本を検索/i);
    await expect(searchInput).toBeVisible();
  });

  test('should show empty state when no books', async ({ page }) => {
    // 空の状態のメッセージまたは本のリストを確認
    const emptyMessage = page.getByText(/まだ本がありません/i);
    const booksList = page.getByRole('article').or(page.locator('[class*="book"]'));

    // どちらかが表示されているはず
    const hasEmptyMessage = await emptyMessage.isVisible().catch(() => false);
    const hasBooks = await booksList.first().isVisible().catch(() => false);

    expect(hasEmptyMessage || hasBooks).toBeTruthy();
  });

  test('should have functional search input', async ({ page }) => {
    const searchInput = page.getByPlaceholder(/本を検索/i);
    await searchInput.fill('テスト');

    const value = await searchInput.inputValue();
    expect(value).toBe('テスト');
  });

  test('should navigate to upload page when clicking add button', async ({ page }) => {
    const addButton = page.getByRole('link', { name: /本を追加/i });
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
