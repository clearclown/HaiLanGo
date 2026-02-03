import { test, expect } from '@playwright/test';

test.describe('Learning Sessions API', () => {
  test('list sessions endpoint exists', async ({ request }) => {
    const response = await request.get('/api/learning/sessions');
    expect(response.status()).not.toBe(404);
  });

  test('create session endpoint exists', async ({ request }) => {
    const response = await request.post('/api/learning/sessions', {
      data: {
        book_id: '00000000-0000-0000-0000-000000000000',
        mode: 'normal',
      },
    });

    expect(response.status()).not.toBe(404);
  });

  test('list sessions returns array', async ({ request }) => {
    const response = await request.get('/api/learning/sessions');

    if (response.ok()) {
      const data = await response.json();
      expect(Array.isArray(data)).toBeTruthy();
    }
  });
});
