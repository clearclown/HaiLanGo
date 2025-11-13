/**
 * 教師モード E2Eテスト
 */

import { test, expect } from '@playwright/test';

test.describe('教師モード', () => {
  test.beforeEach(async ({ page }) => {
    // モックAPIサーバーを使用
    await page.goto('/?useMocks=true');
  });

  test('教師モードの開始', async ({ page }) => {
    await page.goto('/books/test-book/teacher-mode');

    // 開始ボタンをクリック
    await page.click('text=開始');

    // 再生中の表示を確認
    await expect(page.locator('text=再生中')).toBeVisible();
  });

  test('教師モードの一時停止と再開', async ({ page }) => {
    await page.goto('/books/test-book/teacher-mode');

    // 教師モードを開始
    await page.click('text=開始');
    await expect(page.locator('text=再生中')).toBeVisible();

    // 一時停止
    await page.click('text=一時停止');
    await expect(page.locator('text=一時停止中')).toBeVisible();

    // 再開
    await page.click('text=再開');
    await expect(page.locator('text=再生中')).toBeVisible();
  });

  test('ページナビゲーション', async ({ page }) => {
    await page.goto('/books/test-book/teacher-mode');

    // 教師モードを開始
    await page.click('text=開始');
    await expect(page.locator('text=ページ 1')).toBeVisible();

    // 次のページへ
    await page.click('[aria-label="次のページ"]');
    await expect(page.locator('text=ページ 2')).toBeVisible();

    // 前のページへ
    await page.click('[aria-label="前のページ"]');
    await expect(page.locator('text=ページ 1')).toBeVisible();
  });

  test('教師モードの停止', async ({ page }) => {
    await page.goto('/books/test-book/teacher-mode');

    // 教師モードを開始
    await page.click('text=開始');
    await expect(page.locator('text=再生中')).toBeVisible();

    // 停止
    await page.click('text=停止');
    await expect(page.locator('text=開始')).toBeVisible();
  });

  test('設定画面の表示と変更', async ({ page }) => {
    await page.goto('/books/test-book/teacher-mode');

    // 設定ボタンをクリック
    await page.click('[aria-label="設定"]');

    // 設定ダイアログが表示される
    await expect(page.locator('text=教師モード設定')).toBeVisible();

    // 再生速度を変更
    await page.click('text=1.5x');

    // ページ間隔を変更
    await page.fill('[name="pageInterval"]', '10');

    // 保存
    await page.click('text=保存');

    // 設定が反映される
    await expect(page.locator('text=教師モード設定')).not.toBeVisible();
  });

  test('自動ページ遷移', async ({ page }) => {
    await page.goto('/books/test-book/teacher-mode');

    // 教師モードを開始
    await page.click('text=開始');
    await expect(page.locator('text=ページ 1')).toBeVisible();

    // ページ間隔（デフォルト5秒）+ セグメント再生時間を待つ
    await page.waitForTimeout(7000);

    // 次のページに自動遷移する
    await expect(page.locator('text=ページ 2')).toBeVisible();
  });

  test('バックグラウンド再生の確認', async ({ page, context }) => {
    await page.goto('/books/test-book/teacher-mode');

    // 教師モードを開始
    await page.click('text=開始');
    await expect(page.locator('text=再生中')).toBeVisible();

    // Media Session APIが設定されているか確認
    const mediaSessionMetadata = await page.evaluate(() => {
      return navigator.mediaSession?.metadata?.title;
    });

    expect(mediaSessionMetadata).toBeTruthy();
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
    await page.goto('/books/test-book/teacher-mode');

    // 教師モードを開始
    await page.click('text=開始');

    // 最後のページまで移動
    while (await page.locator('[aria-label="次のページ"]').isEnabled()) {
      await page.click('[aria-label="次のページ"]');
      await page.waitForTimeout(500);
    }

    // 最後のページで次へボタンが無効化される
    await expect(page.locator('[aria-label="次のページ"]')).toBeDisabled();
  });

  test('最初のページでの動作', async ({ page }) => {
    await page.goto('/books/test-book/teacher-mode');

    // 教師モードを開始
    await page.click('text=開始');

    // 最初のページで前へボタンが無効化される
    await expect(page.locator('[aria-label="前のページ"]')).toBeDisabled();
  });

  test('ローディング状態の表示', async ({ page }) => {
    // 遅いレスポンスをシミュレート
    await page.route('**/api/v1/books/*/teacher-mode/**', async (route) => {
      await page.waitForTimeout(2000);
      await route.continue();
    });

    await page.goto('/books/test-book/teacher-mode');

    // ローディングスピナーが表示される
    await expect(page.locator('[role="status"]')).toBeVisible();
  });
});
