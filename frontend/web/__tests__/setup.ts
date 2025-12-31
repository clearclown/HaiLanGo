import { cleanup } from '@testing-library/react';
import { afterEach, vi } from 'vitest';
import '@testing-library/jest-dom/vitest';

// Mock Next.js navigation globally
vi.mock('next/navigation', () => ({
  usePathname: () => '/',
  useRouter: () => ({
    push: vi.fn(),
    replace: vi.fn(),
    refresh: vi.fn(),
    back: vi.fn(),
    forward: vi.fn(),
    prefetch: vi.fn(),
  }),
  useSearchParams: () => new URLSearchParams(),
}));

// テスト後にクリーンアップ
afterEach(() => {
  cleanup();
});

// モックの設定
global.fetch = vi.fn();

// 環境変数の設定
process.env.NEXT_PUBLIC_API_URL = 'http://localhost:8080';
process.env.NEXT_PUBLIC_USE_MOCK_APIS = 'true';
