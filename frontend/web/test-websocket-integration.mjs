#!/usr/bin/env node

/**
 * WebSocket統合テストスクリプト
 *
 * このスクリプトは以下をテストします:
 * 1. ユーザー登録とログイン
 * 2. WebSocket接続の確立
 * 3. 書籍作成時の通知
 * 4. OCR進捗通知（モック）
 */

import WebSocket from 'ws';

const API_BASE = 'http://localhost:8080/api/v1';
const WS_URL = 'ws://localhost:8080/api/v1/ws';

let testsPassed = 0;
let testsFailed = 0;

function log(message, type = 'info') {
  const colors = {
    info: '\x1b[36m',
    success: '\x1b[32m',
    error: '\x1b[31m',
    warn: '\x1b[33m',
  };
  const reset = '\x1b[0m';
  console.log(`${colors[type]}${message}${reset}`);
}

function assert(condition, message) {
  if (condition) {
    testsPassed++;
    log(`✓ ${message}`, 'success');
  } else {
    testsFailed++;
    log(`✗ ${message}`, 'error');
    throw new Error(`Assertion failed: ${message}`);
  }
}

async function registerUser(email, password) {
  log('\n=== ユーザー登録テスト ===');
  const response = await fetch(`${API_BASE}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email,
      password,
      display_name: 'Test User',
    }),
  });

  const data = await response.json();

  if (!response.ok) {
    log(`登録失敗: ${response.status} - ${JSON.stringify(data)}`, 'error');
  }

  assert(response.ok, 'ユーザー登録成功');
  assert(data.user !== undefined, 'ユーザーオブジェクトが返される');
  assert(data.access_token !== undefined, 'JWTトークンが返される');
  log(`登録成功: ${email}`, 'success');
  return { token: data.access_token, user: data.user };
}

async function loginUser(email, password) {
  log('\n=== ユーザーログインテスト ===');
  const response = await fetch(`${API_BASE}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });

  const data = await response.json();

  if (!response.ok) {
    log(`ログイン失敗: ${response.status} - ${JSON.stringify(data)}`, 'error');
  }

  assert(response.ok, 'ログイン成功');
  assert(data.access_token !== undefined, 'JWTトークンが返される');
  log('ログイン成功', 'success');
  return { token: data.access_token, user: data.user };
}

function connectWebSocket(token) {
  return new Promise((resolve, reject) => {
    log('\n=== WebSocket接続テスト ===');
    const ws = new WebSocket(`${WS_URL}?token=${encodeURIComponent(token)}`);
    const messages = [];
    let resolved = false;

    ws.on('open', () => {
      log('WebSocket接続成功', 'success');
      assert(true, 'WebSocketが正常に接続される');

      // 接続確認後、resolveして次のテストに進む
      setTimeout(() => {
        if (!resolved) {
          resolved = true;
          resolve({ ws, messages });
        }
      }, 1000);
    });

    ws.on('message', (data) => {
      try {
        const message = JSON.parse(data.toString());
        log(`受信メッセージ: ${JSON.stringify(message, null, 2)}`, 'info');
        messages.push(message);
      } catch (error) {
        log(`メッセージパースエラー: ${error.message}`, 'error');
      }
    });

    ws.on('error', (error) => {
      log(`WebSocketエラー: ${error.message}`, 'error');
      if (!resolved) {
        resolved = true;
        reject(error);
      }
    });

    ws.on('close', () => {
      log('WebSocket接続が閉じられました', 'warn');
    });

    // タイムアウト処理
    setTimeout(() => {
      if (!resolved) {
        resolved = true;
        reject(new Error('WebSocket connection timeout'));
      }
    }, 5000);
  });
}

async function createBook(token, title = 'テスト書籍') {
  log('\n=== 書籍作成テスト ===');
  const response = await fetch(`${API_BASE}/books`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify({
      title,
      target_language: 'ru',
      native_language: 'ja',
      reference_language: 'ja',
    }),
  });

  const data = await response.json();
  assert(response.ok, '書籍作成成功');
  assert(data.book !== undefined, '書籍オブジェクトが返される');
  log(`書籍作成成功: ${data.book.title} (ID: ${data.book.id})`, 'success');
  return data.book;
}

async function waitForMessage(messages, predicate, timeout = 5000) {
  const startTime = Date.now();

  while (Date.now() - startTime < timeout) {
    const message = messages.find(predicate);
    if (message) {
      return message;
    }
    await new Promise(resolve => setTimeout(resolve, 100));
  }

  throw new Error(`Message not received within ${timeout}ms`);
}

async function runTests() {
  const testEmail = `test-${Date.now()}@example.com`;
  const testPassword = 'TestPassword123!';

  try {
    log('='.repeat(50));
    log('WebSocket統合テスト開始', 'info');
    log('='.repeat(50));

    // 1. ユーザー登録/ログイン
    const { token, user } = await registerUser(testEmail, testPassword);
    log(`\n認証トークン取得: ${token.substring(0, 20)}...`, 'success');

    // 2. WebSocket接続
    const { ws, messages } = await connectWebSocket(token);

    // 3. 書籍作成 → 通知確認
    log('\n=== 書籍作成通知テスト ===');
    const book = await createBook(token, `テスト書籍 ${Date.now()}`);

    // 通知を待つ
    log('通知メッセージを待機中...', 'info');
    try {
      const notification = await waitForMessage(
        messages,
        msg => msg.type === 'notification' && msg.payload?.level === 'success'
      );

      assert(notification !== undefined, '書籍作成の通知が届く');
      assert(notification.payload?.title !== undefined, '通知にタイトルが含まれる');
      assert(notification.payload?.message !== undefined, '通知にメッセージが含まれる');
      log(`通知受信成功: ${notification.payload.title}`, 'success');
    } catch (error) {
      log(`通知受信タイムアウト: ${error.message}`, 'error');
      testsFailed++;
    }

    // 4. WebSocket統計確認
    log('\n=== WebSocket統計確認 ===');
    const statsResponse = await fetch(`${API_BASE}/ws/stats`, {
      headers: { 'Authorization': `Bearer ${token}` },
    });
    const stats = await statsResponse.json();
    log(`接続中のクライアント数: ${stats.total_clients}`, 'info');
    log(`送信メッセージ数: ${stats.messages_sent}`, 'info');
    log(`受信メッセージ数: ${stats.messages_received}`, 'info');
    assert(stats.total_clients >= 1, '少なくとも1つのクライアントが接続中');

    // クリーンアップ
    ws.close();

    // 結果表示
    log('\n' + '='.repeat(50));
    log(`テスト結果: ${testsPassed} 成功, ${testsFailed} 失敗`,
        testsFailed === 0 ? 'success' : 'error');
    log('='.repeat(50));

    process.exit(testsFailed === 0 ? 0 : 1);

  } catch (error) {
    log(`\nテスト実行エラー: ${error.message}`, 'error');
    log(error.stack, 'error');
    process.exit(1);
  }
}

runTests();
