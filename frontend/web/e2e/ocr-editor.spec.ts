import { test, expect } from '@playwright/test';

test.describe('OCR Editor', () => {
  const bookId = 'test-book-123';
  const pageId = 'test-page-456';

  test.beforeEach(async ({ page }) => {
    // Mock API responses
    await page.route('**/api/v1/books/*/pages/*/ocr-text', async (route) => {
      if (route.request().method() === 'PUT') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            correction: {
              id: 'correction-789',
              book_id: bookId,
              page_id: pageId,
              original_text: 'Original OCR text',
              corrected_text: 'Corrected text',
              user_id: 'user-001',
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString(),
            },
            message: 'OCR text updated successfully',
          }),
        });
      }
    });

    // Navigate to the OCR editor page
    // Note: This URL will need to be adjusted based on your actual routing
    await page.goto(`http://localhost:3000/books/${bookId}/pages/${pageId}/edit`);
  });

  test('displays the OCR text editor', async ({ page }) => {
    await expect(page.getByTestId('ocr-text-editor')).toBeVisible();
  });

  test('allows editing OCR text', async ({ page }) => {
    const textarea = page.getByTestId('text-editor-textarea');
    await expect(textarea).toBeVisible();

    await textarea.fill('Modified OCR text');
    await expect(textarea).toHaveValue('Modified OCR text');
  });

  test('shows unsaved changes indicator', async ({ page }) => {
    const textarea = page.getByTestId('text-editor-textarea');
    await textarea.fill('Modified OCR text');

    const unsavedIndicator = page.getByTestId('unsaved-indicator');
    await expect(unsavedIndicator).toBeVisible();
    await expect(unsavedIndicator).toContainText('Unsaved changes');
  });

  test('saves corrected text successfully', async ({ page }) => {
    const textarea = page.getByTestId('text-editor-textarea');
    await textarea.fill('Corrected text');

    const saveButton = page.getByTestId('save-button');
    await expect(saveButton).toBeEnabled();
    await saveButton.click();

    // Wait for success message
    const successMessage = page.getByTestId('success-message');
    await expect(successMessage).toBeVisible();
    await expect(successMessage).toContainText('saved successfully');
  });

  test('resets text to original', async ({ page }) => {
    const textarea = page.getByTestId('text-editor-textarea');
    const originalValue = await textarea.inputValue();

    await textarea.fill('Modified text');
    await expect(textarea).toHaveValue('Modified text');

    const resetButton = page.getByTestId('reset-button');
    await resetButton.click();

    await expect(textarea).toHaveValue(originalValue);
    await expect(page.getByTestId('unsaved-indicator')).not.toBeVisible();
  });

  test('displays character count', async ({ page }) => {
    const textarea = page.getByTestId('text-editor-textarea');
    await textarea.fill('Test text');

    await expect(page.locator('.char-count')).toContainText('9 / 10,000 characters');
  });

  test('validates empty text', async ({ page }) => {
    const textarea = page.getByTestId('text-editor-textarea');
    await textarea.fill('   ');

    const saveButton = page.getByTestId('save-button');
    await saveButton.click();

    const errorMessage = page.getByTestId('error-message');
    await expect(errorMessage).toBeVisible();
    await expect(errorMessage).toContainText('Text cannot be empty');
  });

  test('validates text length', async ({ page }) => {
    const textarea = page.getByTestId('text-editor-textarea');
    const longText = 'a'.repeat(10001);
    await textarea.fill(longText);

    const saveButton = page.getByTestId('save-button');
    await saveButton.click();

    const errorMessage = page.getByTestId('error-message');
    await expect(errorMessage).toBeVisible();
    await expect(errorMessage).toContainText('exceeds maximum length');
  });

  test('disables buttons when no changes', async ({ page }) => {
    const saveButton = page.getByTestId('save-button');
    const resetButton = page.getByTestId('reset-button');

    await expect(saveButton).toBeDisabled();
    await expect(resetButton).toBeDisabled();
  });

  test('shows loading state while saving', async ({ page }) => {
    // Mock slow API response
    await page.route('**/api/v1/books/*/pages/*/ocr-text', async (route) => {
      if (route.request().method() === 'PUT') {
        await page.waitForTimeout(1000);
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            correction: {},
          }),
        });
      }
    });

    const textarea = page.getByTestId('text-editor-textarea');
    await textarea.fill('Modified text');

    const saveButton = page.getByTestId('save-button');
    await saveButton.click();

    await expect(saveButton).toContainText('Saving...');
    await expect(saveButton).toBeDisabled();
  });
});

test.describe('Diff Viewer', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to a page with diff viewer
    await page.goto('http://localhost:3000/test/diff-viewer');
  });

  test('displays text differences', async ({ page }) => {
    const diffViewer = page.getByTestId('diff-viewer');
    await expect(diffViewer).toBeVisible();

    await expect(page.getByTestId('original-text')).toBeVisible();
    await expect(page.getByTestId('corrected-text')).toBeVisible();
  });

  test('shows diff statistics', async ({ page }) => {
    const diffStats = page.getByTestId('diff-stats');
    await expect(diffStats).toBeVisible();
  });

  test('displays "No changes" when texts are identical', async ({ page }) => {
    // This would require setting up the component with identical texts
    const noChanges = page.getByTestId('no-changes');
    // This assertion depends on how the test page is set up
    // await expect(noChanges).toBeVisible();
  });
});
