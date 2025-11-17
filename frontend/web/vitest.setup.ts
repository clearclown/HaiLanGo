import '@testing-library/jest-dom';
import { cleanup } from '@testing-library/react';
import { afterEach } from 'vitest';

// テスト後のクリーンアップ
afterEach(() => {
  cleanup();
});

// Media Session API のモック
Object.defineProperty(global.navigator, 'mediaSession', {
  writable: true,
  value: {
    metadata: null,
    playbackState: 'none',
    setActionHandler: () => {},
    setPositionState: () => {},
  },
});
