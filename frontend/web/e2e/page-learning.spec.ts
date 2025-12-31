import { expect, test } from './fixtures';

test.describe('Page Learning', () => {
  let completedPages: Set<string> = new Set();

  test.beforeEach(async ({ page }) => {
    completedPages = new Set();

    // Mock learning API endpoints - matches localhost:8080/api/v1/books/*/pages/**
    await page.route('**/api/v1/books/*/pages/**', async (route) => {
      const url = route.request().url();
      const method = route.request().method();

      // Extract bookId and pageNumber from URL
      const match = url.match(/\/books\/([^/]+)\/pages\/(\d+)/);
      if (!match) {
        await route.continue();
        return;
      }

      const bookId = match[1];
      const pageNumber = Number.parseInt(match[2], 10);

      // Handle POST to /complete
      if (method === 'POST' && url.includes('/complete')) {
        completedPages.add(`${bookId}-${pageNumber}`);
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ success: true }),
        });
        return;
      }

      // Handle GET for page data
      const isCompleted = completedPages.has(`${bookId}-${pageNumber}`);
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: `page-${pageNumber}`,
          bookId: bookId,
          pageNumber: pageNumber,
          imageUrl: `/images/page${pageNumber}.jpg`,
          ocrText: pageNumber === 1 ? 'Здравствуйте!' : 'Привет!',
          translation: pageNumber === 1 ? 'こんにちは！' : 'やあ！',
          audioUrl: `/audio/page${pageNumber}.mp3`,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          isCompleted: isCompleted,
        }),
      });
    });

    // テスト用のページに移動
    await page.goto('/books/test-book/pages/1');
    // Wait for page to load
    await page.waitForLoadState('networkidle');
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

  test('should play audio when clicking play button', async ({ page }) => {
    // 再生ボタンをクリック
    const playButton = page.locator('button[aria-label="再生"]');
    await expect(playButton).toBeVisible();
    await playButton.click();

    // 音声プレイヤーの再生ボタンが一時停止ボタンに変わることを確認
    await expect(page.locator('button[aria-label="一時停止"]')).toBeVisible();
  });

  test('should mark page as completed', async ({ page }) => {
    // 完了済みバッジがまだ表示されていないことを確認
    await expect(page.locator('text=完了済み')).not.toBeVisible();

    // 学習完了ボタンをクリック
    await page.click('button:has-text("学習完了")');

    // APIが完了を記録し、ページがリロードされて完了済みバッジが表示されるまで待つ
    await expect(page.locator('text=完了済み')).toBeVisible({ timeout: 10000 });
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
    // 速度ボタン（デフォルト1x）をクリック
    const speedButton = page.locator('button[aria-label="1x"]');
    await expect(speedButton).toBeVisible();
    await speedButton.click();

    // 速度メニューが表示されることを確認（1.5xオプション）
    const speed15xOption = page.locator('button[aria-label="1.5x"]');
    await expect(speed15xOption).toBeVisible();

    // 1.5xを選択
    await speed15xOption.click();

    // 速度が変更されたことを確認（ボタンのaria-labelが変わる）
    await expect(page.locator('button[aria-label="1.5x"]').first()).toBeVisible();
  });

  test('should show progress bar', async ({ page }) => {
    // 進捗バーが表示されることを確認（ヘッダー内の進捗バー）
    const progressBar = page.locator('header .bg-blue-500').first();
    await expect(progressBar).toBeVisible();

    // 進捗バーの幅が設定されていることを確認
    const width = await progressBar.getAttribute('style');
    expect(width).toContain('width');
  });
});
