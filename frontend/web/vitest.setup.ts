import '@testing-library/jest-dom';
import { cleanup } from '@testing-library/react';
import { afterEach } from 'vitest';

// テスト後のクリーンアップ
afterEach(() => {
  cleanup();
});

// Media Session API のモック
global.navigator.mediaSession = {
  metadata: null,
  playbackState: 'none',
  setActionHandler: () => {},
  setPositionState: () => {},
} as any;
