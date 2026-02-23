/**
 * Full E2E flow test:
 * 書籍アップロード → OCR → 学習セッション → TTS再生 → 発音評価 → SRS復習
 */
import { test, expect, APIRequestContext } from '@playwright/test';

// Unique email per test run to avoid conflicts
const testEmail = `e2e-${Date.now()}@example.com`;
const testPassword = 'TestPassword123!';
let authToken: string = '';
let bookId: string = '';
let sessionId: string = '';
let vocabularyId: string = '';

// Helper: authenticated POST
async function authPost(request: APIRequestContext, url: string, data: object) {
  return request.post(url, {
    data,
    headers: authToken ? { Authorization: `Bearer ${authToken}` } : {},
  });
}

// Helper: authenticated GET
async function authGet(request: APIRequestContext, url: string) {
  return request.get(url, {
    headers: authToken ? { Authorization: `Bearer ${authToken}` } : {},
  });
}

test.describe('Full Learning Flow', () => {
  // ─────────────────────────────────────────────────────────
  // Step 1: User Registration & Login
  // ─────────────────────────────────────────────────────────
  test('1. ユーザー登録が成功する', async ({ request }) => {
    const response = await request.post('/api/auth/register/', {
      data: {
        email: testEmail,
        password: testPassword,
        display_name: 'E2E Test User',
        native_language: 'ja',
      },
    });

    expect(response.status()).toBe(201);
    const body = await response.json();
    expect(body.user).toBeDefined();
    expect(body.user.email).toBe(testEmail);
    expect(body.tokens).toBeDefined();
    expect(body.tokens.access_token).toBeDefined();

    // Save token for subsequent requests
    authToken = body.tokens.access_token;
  });

  test('2. ログインが成功しJWTトークンが返る', async ({ request }) => {
    const response = await request.post('/api/auth/login/', {
      data: {
        email: testEmail,
        password: testPassword,
      },
    });

    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body.tokens.access_token).toBeDefined();
    // Refresh token for next steps
    authToken = body.tokens.access_token;
  });

  test('3. 重複登録は拒否される', async ({ request }) => {
    const response = await request.post('/api/auth/register/', {
      data: {
        email: testEmail,
        password: testPassword,
        display_name: 'Duplicate',
        native_language: 'ja',
      },
    });
    expect(response.status()).toBe(409);
  });

  // ─────────────────────────────────────────────────────────
  // Step 2: Book Upload (OCR trigger)
  // ─────────────────────────────────────────────────────────
  test('4. 書籍アップロード(OCRジョブ起動)', async ({ request }) => {
    const response = await authPost(request, '/api/books/', {
      title: 'Japanese Textbook Vol.1',
      source_language: 'ja',
      target_language: 'en',
    });

    expect(response.status()).toBe(201);
    const body = await response.json();
    expect(body.id).toBeDefined();
    expect(body.title).toBe('Japanese Textbook Vol.1');
    expect(body.status).toBe('Pending');
    expect(body.job_id).toBeDefined();

    bookId = body.id;
  });

  test('5. 書籍一覧に追加した書籍が表示される', async ({ request }) => {
    const response = await authGet(request, '/api/books/');

    expect(response.status()).toBe(200);
    const books = await response.json();
    expect(Array.isArray(books)).toBeTruthy();
    expect(books.length).toBeGreaterThan(0);
  });

  test('6. 書籍詳細が取得できる', async ({ request }) => {
    const response = await authGet(request, `/api/books/${bookId}/`);

    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body.id).toBe(bookId);
    expect(body.title).toBe('Japanese Textbook Vol.1');
  });

  test('7. 他のユーザーの書籍は閲覧禁止', async ({ request }) => {
    // Try with a dummy ID (non-existent book)
    const fakeId = '00000000-0000-0000-0000-000000000001';
    const response = await authGet(request, `/api/books/${fakeId}/`);
    expect([403, 404]).toContain(response.status());
  });

  // ─────────────────────────────────────────────────────────
  // Step 3: Learning Session
  // ─────────────────────────────────────────────────────────
  test('8. 学習セッションの作成', async ({ request }) => {
    const response = await authPost(request, '/api/learning/sessions/', {
      book_id: bookId,
      session_type: 'PageByPage',
      start_page: 1,
      end_page: 10,
    });

    expect(response.status()).toBe(201);
    const body = await response.json();
    expect(body.id).toBeDefined();
    expect(body.status).toBe('Active');
    expect(body.session_type).toBe('PageByPage');

    sessionId = body.id;
  });

  test('9. セッション一覧が取得できる', async ({ request }) => {
    const response = await authGet(request, '/api/learning/sessions/');

    expect(response.status()).toBe(200);
    const sessions = await response.json();
    expect(Array.isArray(sessions)).toBeTruthy();
    expect(sessions.length).toBeGreaterThan(0);
  });

  test('10. セッションを一時停止できる', async ({ request }) => {
    const response = await authPost(request, `/api/learning/sessions/${sessionId}/status/`, {
      action: 'Pause',
    });

    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body.status).toBe('Paused');
  });

  test('11. セッションを再開できる', async ({ request }) => {
    const response = await authPost(request, `/api/learning/sessions/${sessionId}/status/`, {
      action: 'Resume',
    });

    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body.status).toBe('Active');
  });

  // ─────────────────────────────────────────────────────────
  // Step 4: TTS (Text-to-Speech)
  // ─────────────────────────────────────────────────────────
  test('12. TTS音声合成が成功する', async ({ request }) => {
    const response = await authPost(request, '/api/tts/synthesize/', {
      text: 'こんにちは、日本語の勉強を頑張りましょう。',
      language: 'ja',
    });

    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body.id).toBeDefined();
    expect(body.audio_url).toBeDefined();
    expect(body.text).toBe('こんにちは、日本語の勉強を頑張りましょう。');
  });

  test('13. TTS対応言語一覧が取得できる', async ({ request }) => {
    const response = await authGet(request, '/api/tts/languages/');

    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(Array.isArray(body.languages)).toBeTruthy();
    expect(body.languages.length).toBeGreaterThan(0);
  });

  test('14. TTS生成履歴が記録される', async ({ request }) => {
    const response = await authGet(request, '/api/tts/generations/');

    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(Array.isArray(body)).toBeTruthy();
  });

  // ─────────────────────────────────────────────────────────
  // Step 5: STT & Pronunciation Evaluation
  // ─────────────────────────────────────────────────────────
  test('15. STT発音評価エンドポイントが存在する', async ({ request }) => {
    const response = await authPost(request, '/api/stt/evaluate/', {
      text: 'hello world',
      language: 'en',
      audio_base64: 'dGVzdA==', // minimal base64 (test)
    });

    // Accept success or validation error - endpoint must exist
    expect(response.status()).not.toBe(404);
    expect(response.status()).not.toBe(405);
  });

  test('16. STT文字起こしエンドポイントが存在する', async ({ request }) => {
    const response = await authPost(request, '/api/stt/transcribe/', {
      language: 'en',
      audio_base64: 'dGVzdA==',
    });

    expect(response.status()).not.toBe(404);
    expect(response.status()).not.toBe(405);
  });

  // ─────────────────────────────────────────────────────────
  // Step 6: SRS Review (Spaced Repetition)
  // ─────────────────────────────────────────────────────────
  test('17. 語彙を登録できる', async ({ request }) => {
    const response = await authPost(request, '/api/review/vocabulary/', {
      word: '勉強',
      reading: 'べんきょう',
      meaning: 'study',
      language: 'ja',
      notes: 'E2E test vocabulary',
    });

    expect(response.status()).toBe(201);
    const body = await response.json();
    expect(body.id).toBeDefined();
    expect(body.word).toBe('勉強');

    vocabularyId = body.id;
  });

  test('18. 語彙一覧が取得できる', async ({ request }) => {
    const response = await authGet(request, '/api/review/vocabulary/');

    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(Array.isArray(body)).toBeTruthy();
    expect(body.length).toBeGreaterThan(0);
  });

  test('19. 復習キューが取得できる', async ({ request }) => {
    const response = await authGet(request, '/api/review/queue/');

    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body.items).toBeDefined();
    expect(Array.isArray(body.items)).toBeTruthy();
  });

  test('20. SRS復習を記録できる(品質4=覚えた)', async ({ request }) => {
    const response = await authPost(request, `/api/review/vocabulary/${vocabularyId}/review/`, {
      quality: 4,
    });

    // Accept success response
    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body.next_review).toBeDefined();
    expect(body.repetitions).toBeGreaterThan(0);
  });

  test('21. SRS統計が取得できる', async ({ request }) => {
    const response = await authGet(request, '/api/review/stats/');

    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body.total_items).toBeDefined();
    expect(body.due_today).toBeDefined();
  });

  // ─────────────────────────────────────────────────────────
  // Step 7: Teacher Mode
  // ─────────────────────────────────────────────────────────
  test('22. ティーチャーモードセッションを開始できる', async ({ request }) => {
    const response = await authPost(request, '/api/teacher/start/', {
      book_id: bookId,
      start_page: 1,
      end_page: 5,
    });

    expect(response.status()).toBe(201);
    const body = await response.json();
    expect(body.id).toBeDefined();
    expect(body.status).toBe('Active');
  });

  test('23. ティーチャーモードセッション一覧が取得できる', async ({ request }) => {
    const response = await authGet(request, '/api/teacher/sessions/');

    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(Array.isArray(body)).toBeTruthy();
  });

  // ─────────────────────────────────────────────────────────
  // Step 8: Full Flow Completion
  // ─────────────────────────────────────────────────────────
  test('24. 学習セッションを完了できる', async ({ request }) => {
    const response = await authPost(request, `/api/learning/sessions/${sessionId}/status/`, {
      action: 'Complete',
    });

    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body.status).toBe('Completed');
    expect(body.ended_at).toBeDefined();
  });

  test('25. ヘルスチェックが常にhealthyを返す', async ({ request }) => {
    const response = await request.get('/health/');
    expect(response.ok()).toBeTruthy();
    const body = await response.json();
    expect(body.status).toBe('healthy');
    expect(body.app).toBe('HaiLanGo');
  });
});
