import { test, expect } from '@playwright/test';

test.describe('Upload Page Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/upload');
  });

  test('should display upload page correctly', async ({ page }) => {
    // ステップインジケーターの確認
    const stepLabels = ['メタデータ', 'ファイル選択', 'アップロード', '完了'];

    for (const label of stepLabels) {
      const stepElement = page.getByText(label);
      await expect(stepElement).toBeVisible();
    }

    // メタデータフォームが表示されることを確認
    const heading = page.getByRole('heading', { name: /本の情報を入力/i });
    await expect(heading).toBeVisible();
  });

  test('should have all required metadata form fields', async ({ page }) => {
    // タイトルフィールド
    const titleInput = page.getByLabel(/本のタイトル/i);
    await expect(titleInput).toBeVisible();

    // 学習先言語
    const targetLanguageSelect = page.getByLabel(/学習先言語/i);
    await expect(targetLanguageSelect).toBeVisible();

    // 母国語
    const nativeLanguageSelect = page.getByLabel(/母国語/i);
    await expect(nativeLanguageSelect).toBeVisible();

    // 参照言語（オプション）
    const referenceLanguageSelect = page.getByLabel(/参照言語/i);
    await expect(referenceLanguageSelect).toBeVisible();
  });

  test('should validate required fields', async ({ page }) => {
    // 空のまま次へボタンをクリック
    const nextButton = page.getByRole('button', { name: /次へ/i });
    await nextButton.click();

    // HTML5バリデーションまたはアラートが表示されるはず
    // フィールドが未入力の場合、ページ遷移しないことを確認
    await expect(page).toHaveURL(/.*upload/);
  });

  test('should fill metadata form correctly', async ({ page }) => {
    // タイトル入力
    const titleInput = page.getByLabel(/本のタイトル/i);
    await titleInput.fill('テスト本');

    // 学習先言語選択
    const targetLanguageSelect = page.getByLabel(/学習先言語/i);
    await targetLanguageSelect.selectOption('ru'); // ロシア語

    // 母国語選択（デフォルトで日本語のはず）
    const nativeLanguageSelect = page.getByLabel(/母国語/i);
    const nativeValue = await nativeLanguageSelect.inputValue();
    expect(nativeValue).toBeTruthy();

    // フォームが正しく入力されたことを確認
    const titleValue = await titleInput.inputValue();
    expect(titleValue).toBe('テスト本');
  });

  test('should have cancel button that redirects to books page', async ({ page }) => {
    const cancelButton = page.getByRole('button', { name: /キャンセル/i });
    await expect(cancelButton).toBeVisible();

    await cancelButton.click();
    await expect(page).toHaveURL(/.*books/);
  });

  test('should show progress steps correctly', async ({ page }) => {
    // ステップ1が現在アクティブであることを確認
    const step1 = page.locator('text=メタデータ').locator('..');

    // ステップインジケーターの色やスタイルを確認
    // アクティブなステップは青色、それ以外はグレー
    const step1Classes = await step1.getAttribute('class');
    expect(step1Classes).toContain('blue');
  });
});
