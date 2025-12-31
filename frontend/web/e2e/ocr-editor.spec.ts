import { expect, test } from './fixtures';

test.describe('OCR Editor', () => {
  const bookId = 'test-book-123';
  const pageNumber = '1';

  test.beforeEach(async ({ page }) => {
    // Mock page data API
    await page.route(`**/api/v1/books/${bookId}/pages/${pageNumber}`, async (route) => {
      if (route.request().method() === 'GET') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            id: 'page-001',
            book_id: bookId,
            page_number: 1,
            ocr_text: 'Original OCR text from the book page.',
            corrected_text: null,
            image_url: null,
            created_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z',
          }),
        });
      }
    });

    // Mock OCR text update API
    await page.route('**/api/v1/books/*/pages/*/ocr-text', async (route) => {
      if (route.request().method() === 'PUT') {
        const body = JSON.parse(route.request().postData() || '{}');
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            correction: {
              id: 'correction-789',
              book_id: bookId,
              page_id: pageNumber,
              original_text: 'Original OCR text from the book page.',
              corrected_text: body.corrected_text || 'Corrected text',
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
    await page.goto(`/books/${bookId}/pages/${pageNumber}/edit`);
    await page.waitForLoadState('networkidle');
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

    // Original text is loaded first
    await expect(page.locator('.char-count')).toContainText('/ 10,000 characters');

    // Fill with known text and verify count
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
    // Override with slow API response
    await page.route('**/api/v1/books/*/pages/*/ocr-text', async (route) => {
      if (route.request().method() === 'PUT') {
        // Add delay to simulate slow response
        await new Promise((resolve) => setTimeout(resolve, 1500));
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            correction: {
              id: 'correction-789',
              book_id: 'test-book-123',
              page_id: '1',
              original_text: 'Original OCR text from the book page.',
              corrected_text: 'Modified text',
              user_id: 'user-001',
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString(),
            },
          }),
        });
      }
    });

    const textarea = page.getByTestId('text-editor-textarea');
    await textarea.fill('Modified text');

    const saveButton = page.getByTestId('save-button');
    await saveButton.click();

    // Check that it shows saving state
    await expect(saveButton).toContainText('Saving...');
    await expect(saveButton).toBeDisabled();

    // Wait for save to complete
    await expect(saveButton).toContainText('Save Changes', { timeout: 5000 });
  });
});

test.describe('Diff Viewer', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the diff viewer test page
    await page.goto('/test/diff-viewer');
    await page.waitForLoadState('networkidle');
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
    // Click the button to set identical texts
    await page.click('text=Set Identical Texts');

    // Wait for the no-changes indicator to appear
    const noChanges = page.getByTestId('no-changes');
    await expect(noChanges).toBeVisible();
    await expect(noChanges).toContainText('No changes');
  });
});
