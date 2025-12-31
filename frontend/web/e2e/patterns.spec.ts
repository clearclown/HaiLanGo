import { expect, test } from './fixtures';

test.describe('Pattern Extraction and Practice', () => {
  test.beforeEach(async ({ page }) => {
    // Mock patterns list API
    await page.route('**/api/v1/books/*/patterns', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          patterns: [
            {
              id: 'pattern-1',
              book_id: 'test-book',
              type: 'greeting',
              pattern: 'Hello',
              translation: 'こんにちは',
              frequency: 5,
              created_at: '2024-01-01T00:00:00Z',
              updated_at: '2024-01-01T00:00:00Z',
            },
            {
              id: 'pattern-2',
              book_id: 'test-book',
              type: 'question',
              pattern: 'How are you?',
              translation: 'お元気ですか？',
              frequency: 3,
              created_at: '2024-01-01T00:00:00Z',
              updated_at: '2024-01-01T00:00:00Z',
            },
          ],
        }),
      });
    });

    // Mock pattern details API
    await page.route('**/api/v1/patterns/*/practice', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          questions: [
            {
              id: 'practice-1',
              pattern_id: 'pattern-1',
              question: 'How do you say "Hello" in Japanese?',
              correct_answer: 'こんにちは',
              alternative_answers: ['さようなら', 'ありがとう', 'おはよう'],
              difficulty: 1,
              created_at: '2024-01-01T00:00:00Z',
            },
            {
              id: 'practice-2',
              pattern_id: 'pattern-1',
              question: 'How do you greet someone in the afternoon?',
              correct_answer: 'こんにちは',
              alternative_answers: ['おはよう', 'こんばんは', 'おやすみ'],
              difficulty: 2,
              created_at: '2024-01-01T00:00:00Z',
            },
          ],
        }),
      });
    });

    // Mock single pattern API
    await page.route(/\/api\/v1\/patterns\/[^/]+$/, async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: 'pattern-1',
          book_id: 'test-book',
          type: 'greeting',
          pattern: 'Hello',
          translation: 'こんにちは',
          frequency: 5,
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        }),
      });
    });

    // Navigate to the patterns page
    await page.goto('/books/test-book/patterns');
    await page.waitForLoadState('networkidle');
  });

  test('should display pattern list', async ({ page }) => {
    // Wait for patterns to load
    await page.waitForSelector('[data-testid="pattern-card"]');

    // Check that patterns are displayed
    const patterns = await page.$$('[data-testid="pattern-card"]');
    expect(patterns.length).toBeGreaterThan(0);

    // Verify pattern content
    await expect(page.locator('text=Hello')).toBeVisible();
    await expect(page.locator('text=こんにちは')).toBeVisible();
  });

  test('should filter patterns by type', async ({ page }) => {
    // Wait for page to load
    await page.waitForSelector('[data-testid="pattern-card"]');

    // Click greeting filter
    await page.click('button:has-text("Greeting")');

    // Verify only greeting patterns are shown
    const greetingBadges = await page.$$('span:has-text("Greeting")');
    expect(greetingBadges.length).toBeGreaterThan(0);
  });

  test('should open pattern practice on click', async ({ page }) => {
    // Wait for patterns to load
    await page.waitForSelector('[data-testid="pattern-card"]');

    // Click on first pattern
    await page.click('[data-testid="pattern-card"]');

    // Verify practice screen is shown
    await expect(page.locator('text=Question 1 of')).toBeVisible();
    await expect(page.locator('text=Difficulty:')).toBeVisible();
  });

  test('should complete practice exercise', async ({ page }) => {
    // Navigate to practice
    await page.goto('/patterns/pattern-1/practice');
    await page.waitForLoadState('networkidle');

    // Wait for question to load
    await page.waitForSelector("button:has-text('こんにちは')");

    // Answer first question
    await page.click("button:has-text('こんにちは')");

    // Wait for feedback (shows "Correct" in green background)
    await expect(page.locator('text=Correct')).toBeVisible();

    // Wait for next question or completion
    await page.waitForTimeout(2000);

    // Continue answering if there are more questions
    const completionText = await page.locator('text=Completed').isVisible();
    if (!completionText) {
      // Answer remaining questions - click on the correct answer
      await page.click("button:has-text('こんにちは')");
    }
  });

  test('should show completion screen after all questions', async ({ page }) => {
    // Navigate to practice
    await page.goto('/patterns/pattern-1/practice');
    await page.waitForLoadState('networkidle');

    // Answer all questions (2 questions in mock)
    // First question
    await page.waitForSelector("button:has-text('こんにちは')");
    await page.click("button:has-text('こんにちは')");
    await page.waitForTimeout(2000);

    // Second question
    await page.waitForSelector("button:has-text('こんにちは')");
    await page.click("button:has-text('こんにちは')");
    await page.waitForTimeout(2000);

    // Verify completion screen (component shows "🎉 Practice Completed!")
    await expect(page.locator('text=Practice Completed!')).toBeVisible({ timeout: 10000 });
    await expect(page.locator('text=out of')).toBeVisible();
  });

  test('should restart practice after completion', async ({ page }) => {
    // Navigate to practice and complete it
    await page.goto('/patterns/pattern-1/practice');
    await page.waitForLoadState('networkidle');

    // Answer all questions to reach completion screen
    // First question
    await page.waitForSelector("button:has-text('こんにちは')");
    await page.click("button:has-text('こんにちは')");
    await page.waitForTimeout(2000);

    // Second question
    await page.waitForSelector("button:has-text('こんにちは')");
    await page.click("button:has-text('こんにちは')");
    await page.waitForTimeout(2000);

    // Wait for completion screen
    await expect(page.locator('text=Practice Completed!')).toBeVisible({ timeout: 10000 });

    // Click practice again button
    await page.click('button:has-text("Practice Again")');

    // Verify back to first question
    await expect(page.locator('text=Question 1 of')).toBeVisible();
  });

  test('should display pattern frequency', async ({ page }) => {
    // Wait for patterns to load
    await page.waitForSelector('[data-testid="pattern-card"]');

    // Verify frequency is displayed
    await expect(page.locator('text=×5')).toBeVisible();
  });

  test('should sort patterns by frequency', async ({ page }) => {
    // Wait for patterns to load
    await page.waitForSelector('[data-testid="pattern-card"]');

    // Get first pattern
    const firstPattern = await page
      .$('[data-testid="pattern-card"]')
      .then((el) => el?.textContent());

    // Verify sorting (highest frequency first)
    expect(firstPattern).toContain('×5');
  });

  test('should show progress bar during practice', async ({ page }) => {
    // Navigate to practice
    await page.goto('/patterns/pattern-1/practice');
    await page.waitForLoadState('networkidle');

    // Wait for progress bar
    await page.waitForSelector('.bg-blue-600');

    // Verify progress bar exists
    const progressBar = await page.$('.bg-blue-600');
    expect(progressBar).not.toBeNull();
  });

  test('should highlight correct and incorrect answers', async ({ page }) => {
    // Navigate to practice
    await page.goto('/patterns/pattern-1/practice');
    await page.waitForLoadState('networkidle');

    // Wait for question
    await page.waitForSelector("button:has-text('さようなら')");

    // Click wrong answer
    await page.click("button:has-text('さようなら')");

    // Wait for highlighting
    await page.waitForTimeout(500);

    // Verify incorrect answer is highlighted (red background)
    const incorrectAnswer = await page.$('.bg-red-100');
    expect(incorrectAnswer).not.toBeNull();

    // Verify correct answer is also shown (green background)
    const correctAnswer = await page.$('.bg-green-100');
    expect(correctAnswer).not.toBeNull();
  });
});
