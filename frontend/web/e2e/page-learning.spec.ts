import { test, expect } from '@playwright/test';

test.describe('Page Learning', () => {
  test.beforeEach(async ({ page }) => {
    // テスト用のページに移動
    await page.goto('/books/test-book/pages/1');
  });

  test('should display page content', async ({ page }) => {
    // ページ番号が表示されることを確認
    await expect(page.locator('text=ページ 1')).toBeVisible();

    // ページ画像が表示されることを確認
    const pageImage = page.locator('img[alt*="ページ"]');
    await expect(pageImage).toBeVisible();

    // テキストが表示されることを確認
    await expect(page.locator('text=Здравствуйте!')).toBeVisible();
  });

  test('should navigate to next page', async ({ page }) => {
    // 次へボタンをクリック
    await page.click('button:has-text("次へ")');

    // ページ2に遷移することを確認
    await expect(page).toHaveURL('/books/test-book/pages/2');
    await expect(page.locator('text=ページ 2')).toBeVisible();
  });

  test('should navigate to previous page', async ({ page }) => {
    // まずページ2に移動
    await page.goto('/books/test-book/pages/2');
    await expect(page.locator('text=ページ 2')).toBeVisible();

    // 前へボタンをクリック
    await page.click('button:has-text("前へ")');

    // ページ1に遷移することを確認
    await expect(page).toHaveURL('/books/test-book/pages/1');
    await expect(page.locator('text=ページ 1')).toBeVisible();
  });

  test('should disable previous button on first page', async ({ page }) => {
    // 前へボタンが無効化されていることを確認
    const prevButton = page.locator('button:has-text("前へ")');
    await expect(prevButton).toBeDisabled();
  });

  test('should play audio when clicking page', async ({ page }) => {
    // ページ画像をクリック
    const pageContent = page.locator('[data-testid="page-content"]');
    await pageContent.click();

    // 音声プレイヤーの再生ボタンが一時停止ボタンに変わることを確認
    await expect(page.locator('button[aria-label="一時停止"]')).toBeVisible();
  });

  test('should mark page as completed', async ({ page }) => {
    // 学習完了ボタンをクリック
    await page.click('button:has-text("学習完了")');

    // 完了済みバッジが表示されることを確認
    await expect(page.locator('text=完了済み')).toBeVisible();
  });

  test('should show loading state', async ({ page }) => {
    // ネットワークを遅くする
    await page.route('**/api/v1/books/**', (route) => {
      setTimeout(() => route.continue(), 2000);
    });

    await page.goto('/books/test-book/pages/1');

    // ローディング状態が表示されることを確認
    await expect(page.locator('text=Loading...')).toBeVisible();
  });

  test('should show error state when API fails', async ({ page }) => {
    // APIエラーをシミュレート
    await page.route('**/api/v1/books/**', (route) => {
      route.fulfill({
        status: 500,
        body: JSON.stringify({ error: 'Internal Server Error' }),
      });
    });

    await page.goto('/books/test-book/pages/1');

    // エラーメッセージが表示されることを確認
    await expect(page.locator('text=/error/i')).toBeVisible();
  });

  test('should change audio speed', async ({ page }) => {
    // 速度ボタンをクリック
    await page.click('button:has-text("1.0x")');

    // 速度メニューが表示されることを確認
    await expect(page.locator('button:has-text("1.5x")')).toBeVisible();

    // 1.5xを選択
    await page.click('button:has-text("1.5x")');

    // 速度が変更されたことを確認
    await expect(page.locator('button:has-text("1.5x")')).toBeVisible();
  });

  test('should show progress bar', async ({ page }) => {
    // 進捗バーが表示されることを確認
    const progressBar = page.locator('.bg-blue-500');
    await expect(progressBar).toBeVisible();

    // 進捗バーの幅が設定されていることを確認
    const width = await progressBar.getAttribute('style');
    expect(width).toContain('width');
  });
});
