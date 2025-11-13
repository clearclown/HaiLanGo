import { expect, test } from '@playwright/test';

test.describe('設定画面', () => {
  test.beforeEach(async ({ page }) => {
    // モックAPIをセットアップ
    await page.route('**/api/v1/settings', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          profile: {
            id: '1',
            name: '太郎',
            email: 'taro@example.com',
          },
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
        body: JSON.stringify({
          type: 'free',
        }),
      });
    });

    await page.goto('/settings');
  });

  test('設定画面が表示される', async ({ page }) => {
    await expect(page.locator('h1:has-text("設定")')).toBeVisible();
    await expect(page.locator('h2:has-text("アカウント")')).toBeVisible();
    await expect(page.locator('h2:has-text("プラン")')).toBeVisible();
    await expect(page.locator('h2:has-text("通知設定")')).toBeVisible();
  });

  test('アカウント情報が表示される', async ({ page }) => {
    await expect(page.locator('input[id="name"]')).toHaveValue('太郎');
    await expect(page.locator('input[id="email"]')).toHaveValue('taro@example.com');
  });

  test('名前を変更できる', async ({ page }) => {
    await page.route('**/api/v1/settings/profile', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ success: true }),
      });
    });

    const nameInput = page.locator('input[id="name"]');
    await nameInput.fill('次郎');

    const saveButton = page.locator('button:has-text("保存")');
    await saveButton.click();

    await expect(nameInput).toHaveValue('次郎');
  });

  test('通知設定を変更できる', async ({ page }) => {
    await page.route('**/api/v1/settings/notifications', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ success: true }),
      });
    });

    const learningReminderSwitch = page.locator('input[id="learningReminder"]');
    await expect(learningReminderSwitch).toBeChecked();
    await learningReminderSwitch.click();
    await expect(learningReminderSwitch).not.toBeChecked();
  });

  test('プランが表示される', async ({ page }) => {
    await expect(page.locator('text=無料プラン')).toBeVisible();
    await expect(page.locator('button:has-text("プレミアムにアップグレード")')).toBeVisible();
  });

  test('ログアウトダイアログが表示される', async ({ page }) => {
    const logoutButton = page.locator('button:has-text("ログアウト")');
    await logoutButton.click();

    await expect(page.locator('text=ログアウトしますか？')).toBeVisible();
  });

  test('ログアウトをキャンセルできる', async ({ page }) => {
    const logoutButton = page.locator('button:has-text("ログアウト")');
    await logoutButton.click();

    const cancelButton = page.locator('button:has-text("キャンセル")');
    await cancelButton.click();

    await expect(page.locator('text=ログアウトしますか？')).not.toBeVisible();
  });
});
