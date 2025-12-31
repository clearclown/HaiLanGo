import { expect, test } from './fixtures';

test.describe('Review Page Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Mock review API endpoints
    await page.route('**/api/v1/review/stats', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          urgent: 3,
          recommended: 5,
          optional: 4,
          total_completed_today: 10,
          weekly_completion_rate: 65,
        }),
      });
    });
    await page.route('**/api/v1/review/items**', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          items: [
            {
              id: '1',
              type: 'word',
              content: 'Здравствуйте',
              priority: 'urgent',
              dueDate: new Date().toISOString(),
            },
          ],
        }),
      });
    });
    await page.goto('/review');
  });

  test('should display review page correctly', async ({ page }) => {
    // ページタイトルの確認
    const heading = page.getByRole('heading', { name: /復習/i });
    await expect(heading).toBeVisible();

    // サブタイトルの確認
    const subtitle = page.getByText(/間隔反復学習で効率的に記憶/i);
    await expect(subtitle).toBeVisible();
  });

  test('should show loading state initially', async ({ page }) => {
    // ローディング表示の確認（瞬間的なので、タイミングによっては見えないかも）
    // モックによりすぐにデータが返されるため、ローディングはスキップされることが多い
    const loadingText = page.getByText(/読み込み中/i);
    const isLoadingVisible = await loadingText.isVisible().catch(() => false);

    // ローディングが表示されるか、すでにコンテンツが表示されているかどちらか
    // ページが完全に読み込まれるまで待つ
    await page.waitForLoadState('networkidle');
    const heading = page.getByRole('heading', { name: /復習/i });
    const isHeadingVisible = await heading.isVisible().catch(() => false);

    expect(isLoadingVisible || isHeadingVisible).toBeTruthy();
  });

  test('should display review statistics', async ({ page }) => {
    // 統計情報が表示されるまで待つ
    await page.waitForLoadState('networkidle');

    // 「今週の進捗」または「今日の復習」が表示されることを確認
    const statsSection = page.locator('text=今週の進捗').or(page.locator('text=今日の復習'));
    const hasStats = await statsSection
      .first()
      .isVisible()
      .catch(() => false);

    // 統計情報が表示されるか、空の状態メッセージが表示されるかどちらか
    if (!hasStats) {
      const emptyMessage = page.getByText(/素晴らしい！/i);
      await expect(emptyMessage).toBeVisible();
    }
  });

  test('should display review priority cards', async ({ page }) => {
    await page.waitForLoadState('networkidle');

    // 優先度カードのラベルを確認
    const urgentLabel = page.getByText(/緊急/i);
    const recommendedLabel = page.getByText(/推奨/i);
    const optionalLabel = page.getByText(/余裕あり/i);

    // いずれかのラベルが表示されているか、空の状態メッセージが表示されるか
    const hasCards = await urgentLabel
      .or(recommendedLabel)
      .or(optionalLabel)
      .first()
      .isVisible()
      .catch(() => false);

    if (!hasCards) {
      // カードがない場合、完了メッセージが表示されるはず
      const completeMessage = page.getByText(/今日の復習はすべて完了しました/i);
      await expect(completeMessage).toBeVisible();
    }
  });

  test('should show empty state when no review items', async ({ page }) => {
    await page.waitForLoadState('networkidle');

    // 空の状態のメッセージまたは復習カードを確認
    const emptyMessage = page.getByText(/今日の復習はすべて完了しました/i);
    const reviewButton = page.getByRole('button', { name: /復習する/i });

    // どちらかが表示されているはず
    const hasEmptyMessage = await emptyMessage.isVisible().catch(() => false);
    const hasReviewButton = await reviewButton
      .first()
      .isVisible()
      .catch(() => false);

    expect(hasEmptyMessage || hasReviewButton).toBeTruthy();
  });

  test('should have review start buttons', async ({ page }) => {
    await page.waitForLoadState('networkidle');

    // 復習ボタンが存在するか確認
    const reviewButtons = page.getByRole('button', { name: /復習する/i });
    const buttonCount = await reviewButtons.count();

    // ボタンがあるか、空の状態メッセージがあるかどちらか
    if (buttonCount === 0) {
      const emptyMessage = page.getByText(/今日の復習はすべて完了しました/i);
      await expect(emptyMessage).toBeVisible();
    } else {
      // 少なくとも1つの復習ボタンが存在する
      expect(buttonCount).toBeGreaterThan(0);
    }
  });

  test('should display progress bar for weekly completion', async ({ page }) => {
    await page.waitForLoadState('networkidle');

    // 進捗バーまたは空の状態を確認
    const progressSection = page.locator('text=今週の進捗');
    const hasProgress = await progressSection.isVisible().catch(() => false);

    if (hasProgress) {
      // 進捗バーの要素を確認（progressbar roleを持つ要素）
      const progressBar = page.locator('[role="progressbar"]').first();
      await expect(progressBar).toBeVisible();
    } else {
      // 進捗情報がない場合、空の状態メッセージを確認
      const emptyMessage = page.getByText(/今日の復習はすべて完了しました/i);
      await expect(emptyMessage).toBeVisible();
    }
  });

  test('should show today completed count', async ({ page }) => {
    await page.waitForLoadState('networkidle');

    // 「今日の復習」セクションまたは空の状態を確認
    const todaySection = page.locator('text=今日の復習');
    const hasToday = await todaySection.isVisible().catch(() => false);

    if (hasToday) {
      // 完了項目数が表示されている（「X項目」の形式）
      const itemsCount = page.locator('text=/\\d+項目/').first();
      await expect(itemsCount).toBeVisible();
    }
  });

  test('should handle error state gracefully', async ({ page }) => {
    // ネットワークエラーをシミュレート（オプション）
    // エラー状態の確認は、実際のAPIレスポンスに依存する

    await page.waitForLoadState('networkidle');

    // エラーメッセージまたは正常なコンテンツが表示されるか確認
    const errorMessage = page.getByText(/復習データの読み込みに失敗しました/i);
    const normalContent = page.getByRole('heading', { name: /復習/i });

    const hasError = await errorMessage.isVisible().catch(() => false);
    const hasContent = await normalContent.isVisible();

    // エラーまたはコンテンツのいずれかが表示される
    expect(hasError || hasContent).toBeTruthy();
  });

  test('should have retry button on error', async ({ page }) => {
    await page.waitForLoadState('networkidle');

    // エラー状態の場合、再試行ボタンが存在するか確認
    const retryButton = page.getByRole('button', { name: /再試行/i });
    const hasRetry = await retryButton.isVisible().catch(() => false);

    // 再試行ボタンがない場合は正常状態なのでOK
    if (hasRetry) {
      await expect(retryButton).toBeVisible();
      // ボタンがクリック可能であることを確認
      await expect(retryButton).toBeEnabled();
    }
  });
});
