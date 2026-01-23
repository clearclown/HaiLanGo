import { test as base, expect } from '@playwright/test';

// Mock user data for tests
export const mockUser = {
  id: 'test-user-1',
  email: 'test@example.com',
  display_name: '太郎',
};

// Mock tokens
const mockAccessToken = 'test-access-token-for-e2e';
const mockRefreshToken = 'test-refresh-token-for-e2e';

/**
 * Extended test fixture that handles authentication by setting
 * access_token in localStorage before each test
 */
export const test = base.extend({
  // This hook runs before each test, setting up authentication
  page: async ({ page }, use) => {
    // Set E2E test mode flag before page loads
    await page.addInitScript(() => {
      // biome-ignore lint/suspicious/noExplicitAny: E2E test flag
      (window as any).__E2E_TEST_MODE__ = true;
    });

    // Navigate to a page first to set localStorage (localStorage is domain-bound)
    await page.goto('/login');

    // Set authentication tokens in localStorage
    await page.evaluate(
      ({ token, refreshToken, user }) => {
        localStorage.setItem('access_token', token);
        localStorage.setItem('refresh_token', refreshToken);
        localStorage.setItem('user', JSON.stringify(user));
      },
      {
        token: mockAccessToken,
        refreshToken: mockRefreshToken,
        user: mockUser,
      }
    );

    // Now the page is authenticated for subsequent navigations
    await use(page);
  },
});

export { expect };

/**
 * Helper to setup API mocks for common endpoints
 */
export async function setupAuthenticatedMocks(
  page: typeof test extends typeof base ? Parameters<Parameters<typeof test>[1]>[0]['page'] : never
) {
  // Mock auth validation endpoint if it exists
  await page.route('**/api/v1/auth/validate', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        valid: true,
        user: mockUser,
      }),
    });
  });

  // Mock user profile endpoint
  await page.route('**/api/v1/users/me', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(mockUser),
    });
  });
}
