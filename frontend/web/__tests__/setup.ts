import { expect, afterEach, vi } from 'vitest';
import { cleanup } from '@testing-library/react';
import '@testing-library/jest-dom/vitest';

// テスト後にクリーンアップ
afterEach(() => {
  cleanup();
});

// モックの設定
global.fetch = vi.fn();

// 環境変数の設定
process.env.NEXT_PUBLIC_API_URL = 'http://localhost:8080';
process.env.NEXT_PUBLIC_USE_MOCK_APIS = 'true';
