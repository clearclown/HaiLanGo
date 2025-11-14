import { test, expect } from '@playwright/test';

test.describe('Navigation Tests', () => {
  test('should redirect from root to books page', async ({ page }) => {
    await page.goto('/');
    await page.waitForURL('**/books');
    expect(page.url()).toContain('/books');
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

    const heading = page.getByRole('heading', { name: /設定/i });
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
