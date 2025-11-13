import { test, expect } from '@playwright/test';

test('homepage has title', async ({ page }) => {
  await page.goto('/');
  await expect(page.getByRole('heading', { name: /HaiLanGo/i })).toBeVisible();
});

test('homepage has description', async ({ page }) => {
  await page.goto('/');
  await expect(page.getByText(/Coming Soon/i)).toBeVisible();
});
