/**
 * 教師モード E2Eテスト
 */

import { expect, test } from './fixtures';

test.describe('教師モード', () => {
  test.beforeEach(async ({ page }) => {
    // Mock teacher-mode API endpoints
    await page.route('**/api/v1/books/*/teacher-mode/generate', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          playlistId: 'playlist-1',
          totalPages: 5,
          estimatedDuration: 300,
          pages: [
            {
              pageNumber: 1,
              segments: [{ type: 'text', audioUrl: '/audio/p1s1.mp3', duration: 5, text: 'Hello' }],
            },
            {
              pageNumber: 2,
              segments: [{ type: 'text', audioUrl: '/audio/p2s1.mp3', duration: 5, text: 'World' }],
            },
            {
              pageNumber: 3,
              segments: [{ type: 'text', audioUrl: '/audio/p3s1.mp3', duration: 5, text: 'Test' }],
            },
            {
              pageNumber: 4,
              segments: [{ type: 'text', audioUrl: '/audio/p4s1.mp3', duration: 5, text: 'Page' }],
            },
            {
              pageNumber: 5,
              segments: [{ type: 'text', audioUrl: '/audio/p5s1.mp3', duration: 5, text: 'Five' }],
            },
          ],
        }),
      });
    });

    await page.goto('/books/test-book/teacher-mode');
    await page.waitForLoadState('networkidle');
  });

  test('教師モードの開始', async ({ page }) => {
    // 開始ボタン（再生アイコン）をクリック
    const playButton = page.locator('button:has(.sr-only:text("開始"))');
    await expect(playButton).toBeVisible();
    await playButton.click();

    // 再生中の表示を確認
    await expect(page.locator('text=再生中')).toBeVisible();
  });

  test('教師モードの一時停止と再開', async ({ page }) => {
    // 教師モードを開始
    const playButton = page.locator('button:has(.sr-only:text("開始"))');
    await playButton.click();
    await expect(page.locator('text=再生中')).toBeVisible();

    // 一時停止ボタンをクリック
    const pauseButton = page.locator('button:has(.sr-only:text("一時停止"))');
    await pauseButton.click();
    await expect(page.locator('text=一時停止中')).toBeVisible();

    // 再開ボタンをクリック
    const resumeButton = page.locator('button:has(.sr-only:text("再開"))');
    await resumeButton.click();
    await expect(page.locator('text=再生中')).toBeVisible();
  });

  test('ページナビゲーション', async ({ page }) => {
    // 教師モードを開始
    const playButton = page.locator('button:has(.sr-only:text("開始"))');
    await playButton.click();
    await expect(page.locator('text=ページ 1')).toBeVisible();

    // 次のページへ
    await page.click('[aria-label="次のページ"]');
    await expect(page.locator('text=ページ 2')).toBeVisible();

    // 前のページへ
    await page.click('[aria-label="前のページ"]');
    await expect(page.locator('text=ページ 1')).toBeVisible();
  });

  test('教師モードの停止', async ({ page }) => {
    // 教師モードを開始
    const playButton = page.locator('button:has(.sr-only:text("開始"))');
    await playButton.click();
    await expect(page.locator('text=再生中')).toBeVisible();

    // 停止ボタンをクリック（赤い停止ボタンを正確に選択）
    const stopButton = page.getByRole('button', { name: '停止', exact: true });
    await stopButton.click();

    // 停止ボタンが消えて、開始ボタンのみが表示される
    await expect(stopButton).not.toBeVisible({ timeout: 10000 });
  });

  test('設定画面の表示と変更', async ({ page }) => {
    // 設定ボタンをクリック
    await page.click('[aria-label="設定"]');

    // 設定ダイアログが表示される
    await expect(page.locator('text=教師モード設定')).toBeVisible();

    // 再生速度を変更
    await page.click('text=1.5x');

    // 保存ボタンをクリック
    await page.click('text=保存');

    // ダイアログが閉じる
    await expect(page.locator('text=教師モード設定')).not.toBeVisible();
  });

  test.skip('自動ページ遷移', async ({ page }) => {
    // TODO: 音声再生が必要なテストはスキップ（音声ファイルのモック複雑）
    // 教師モードを開始
    const playButton = page.locator('button:has(.sr-only:text("開始"))');
    await playButton.click();
    await expect(page.locator('text=ページ 1')).toBeVisible();

    // ページ間隔（デフォルト5秒）+ セグメント再生時間を待つ
    await page.waitForTimeout(7000);

    // 次のページに自動遷移する
    await expect(page.locator('text=ページ 2')).toBeVisible();
  });

  test('バックグラウンド再生の確認', async ({ page }) => {
    // 教師モードを開始
    const playButton = page.locator('button:has(.sr-only:text("開始"))');
    await playButton.click();
    await expect(page.locator('text=再生中')).toBeVisible();

    // Media Session APIが設定されているか確認
    const mediaSessionMetadata = await page.evaluate(() => {
      return navigator.mediaSession?.metadata?.title;
    });

    expect(mediaSessionMetadata).toBe('教師モード');
  });

  test('エラーハンドリング', async ({ page }) => {
    // ネットワークエラーをシミュレート
    await page.route('**/api/v1/books/*/teacher-mode/**', (route) => {
      route.abort('failed');
    });

    await page.goto('/books/test-book/teacher-mode');

    // エラーメッセージが表示される
    await expect(page.locator('text=エラーが発生しました')).toBeVisible();
  });

  test('最後のページでの動作', async ({ page }) => {
    // 教師モードを開始
    const playButton = page.locator('button:has(.sr-only:text("開始"))');
    await playButton.click();
    await expect(page.locator('text=ページ 1')).toBeVisible();

    // 最後のページまで移動（5ページ）
    for (let i = 0; i < 4; i++) {
      await page.click('[aria-label="次のページ"]');
      await page.waitForTimeout(300);
    }

    // 最後のページで次へボタンが無効化される
    await expect(page.locator('[aria-label="次のページ"]')).toBeDisabled();
  });

  test('最初のページでの動作', async ({ page }) => {
    // 教師モードを開始
    const playButton = page.locator('button:has(.sr-only:text("開始"))');
    await playButton.click();

    // 最初のページで前へボタンが無効化される
    await expect(page.locator('[aria-label="前のページ"]')).toBeDisabled();
  });

  test('ローディング状態の表示', async ({ page }) => {
    // 遅いレスポンスをシミュレート - 新しいルートを設定する前にナビゲート前に設定
    await page.route('**/api/v1/books/*/teacher-mode/generate', async (route) => {
      await new Promise((resolve) => setTimeout(resolve, 2000));
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          playlistId: 'playlist-1',
          totalPages: 1,
          estimatedDuration: 10,
          pages: [
            {
              pageNumber: 1,
              segments: [{ type: 'text', audioUrl: '/audio/1.mp3', duration: 5, text: 'Test' }],
            },
          ],
        }),
      });
    });

    await page.goto('/books/test-book/teacher-mode');

    // ローディングスピナーが表示される（output要素は暗黙的にstatus roleを持つ）
    await expect(page.locator('output[aria-live="polite"]')).toBeVisible();
  });
});
