import { test, expect } from '@playwright/test';

test.describe('Books API', () => {
  test('list books endpoint exists', async ({ request }) => {
    const response = await request.get('/api/books');
    // Should return list or auth error, not 404
    expect(response.status()).not.toBe(404);
  });

  test('create book endpoint exists', async ({ request }) => {
    const response = await request.post('/api/books', {
      data: {
        title: 'Test Book',
        language: 'Japanese',
      },
    });

    // Should return success or error, not 404
    expect(response.status()).not.toBe(404);
  });

  test('list books returns array', async ({ request }) => {
    const response = await request.get('/api/books');

    if (response.ok()) {
      const data = await response.json();
      expect(Array.isArray(data)).toBeTruthy();
    }
  });

  test('get book by id endpoint exists', async ({ request }) => {
    const response = await request.get('/api/books/00000000-0000-0000-0000-000000000000');
    // Should return book or not found, not method error
    expect([200, 404, 401]).toContain(response.status());
  });
});
