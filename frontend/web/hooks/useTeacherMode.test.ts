/**
 * useTeacherMode フックのテスト
 */

import { renderHook, act, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useTeacherMode } from './useTeacherMode';

// モックの設定
const mockFetchPlaylist = vi.fn();
vi.mock('@/services/teacherModeApi', () => ({
  teacherModeApi: {
    fetchPlaylist: () => mockFetchPlaylist(),
  },
}));

describe('useTeacherMode', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.useFakeTimers();

    // デフォルトのモックプレイリスト
    mockFetchPlaylist.mockResolvedValue({
      id: 'playlist-1',
      bookId: 'test-book',
      pages: [
        {
          pageNumber: 1,
          segments: [
            {
              id: 'segment-1',
              type: 'phrase',
              audioUrl: 'http://example.com/audio1.mp3',
              duration: 2000,
              text: 'Hello',
              language: 'en',
            },
          ],
          totalDuration: 2000,
        },
        {
          pageNumber: 2,
          segments: [
            {
              id: 'segment-2',
              type: 'phrase',
              audioUrl: 'http://example.com/audio2.mp3',
              duration: 2000,
              text: 'World',
              language: 'en',
            },
          ],
          totalDuration: 2000,
        },
      ],
      settings: {
        speed: 1.0,
        pageInterval: 5,
        repeatCount: 1,
        audioQuality: 'standard',
        content: {
          includeTranslation: true,
          includeWordExplanation: true,
          includeGrammarExplanation: false,
          includePronunciationPractice: false,
          includeExampleSentences: false,
        },
      },
      totalDuration: 4000,
    });
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('初期状態では停止している', () => {
    const { result } = renderHook(() => useTeacherMode('test-book'));

    expect(result.current.playbackState.status).toBe('stopped');
    expect(result.current.playbackState.currentPage).toBe(0);
  });

  it('play() を呼ぶと再生が開始される', async () => {
    const { result } = renderHook(() => useTeacherMode('test-book'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.play();
    });

    expect(result.current.playbackState.status).toBe('playing');
    expect(result.current.playbackState.currentPage).toBe(1);
  });

  it('pause() を呼ぶと一時停止される', async () => {
    const { result } = renderHook(() => useTeacherMode('test-book'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.play();
    });

    act(() => {
      result.current.pause();
    });

    expect(result.current.playbackState.status).toBe('paused');
  });

  it('stop() を呼ぶと停止される', async () => {
    const { result } = renderHook(() => useTeacherMode('test-book'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.play();
    });

    act(() => {
      result.current.stop();
    });

    expect(result.current.playbackState.status).toBe('stopped');
    expect(result.current.playbackState.currentPage).toBe(0);
  });

  it('next() を呼ぶと次のページに移動する', async () => {
    const { result } = renderHook(() => useTeacherMode('test-book'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.play();
    });

    act(() => {
      result.current.next();
    });

    expect(result.current.playbackState.currentPage).toBe(2);
  });

  it('previous() を呼ぶと前のページに移動する', async () => {
    const { result } = renderHook(() => useTeacherMode('test-book'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.play();
    });

    act(() => {
      result.current.next();
    });

    act(() => {
      result.current.previous();
    });

    expect(result.current.playbackState.currentPage).toBe(1);
  });

  it('最初のページで previous() を呼んでも0ページには移動しない', async () => {
    const { result } = renderHook(() => useTeacherMode('test-book'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.play();
    });

    act(() => {
      result.current.previous();
    });

    expect(result.current.playbackState.currentPage).toBe(1);
  });

  it('最後のページで next() を呼ぶと停止する', async () => {
    const { result } = renderHook(() => useTeacherMode('test-book'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.play();
    });

    // ページ2に移動
    act(() => {
      result.current.next();
    });

    // 最後のページでnext()を呼ぶ
    act(() => {
      result.current.next();
    });

    expect(result.current.playbackState.status).toBe('stopped');
  });

  it('プレイリストの取得に失敗した場合はエラーが設定される', async () => {
    mockFetchPlaylist.mockRejectedValue(new Error('API Error'));

    const { result } = renderHook(() => useTeacherMode('test-book'));

    await waitFor(() => {
      expect(result.current.error).not.toBeNull();
    });

    expect(result.current.error?.message).toBe('API Error');
  });

  it('ページ間隔の後に自動的に次のページに移動する', async () => {
    const { result } = renderHook(() => useTeacherMode('test-book'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.play();
    });

    // ページ間隔（5秒）を進める
    act(() => {
      vi.advanceTimersByTime(7000); // セグメント再生時間(2秒) + ページ間隔(5秒)
    });

    await waitFor(() => {
      expect(result.current.playbackState.currentPage).toBe(2);
    });
  });
});
