import { test, expect } from '@playwright/test';

test.describe('Health Checks', () => {
  test('health endpoint returns healthy status', async ({ request }) => {
    const response = await request.get('/health');
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.status).toBe('healthy');
    expect(data.app).toBe('HaiLanGo');
  });

  test('root endpoint returns API info', async ({ request }) => {
    const response = await request.get('/');
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.app).toBe('HaiLanGo');
    expect(data.endpoints).toBeDefined();
  });
});
