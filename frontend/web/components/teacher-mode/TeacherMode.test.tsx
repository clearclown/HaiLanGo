/**
 * 教師モードコンポーネントのテスト
 */

import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { TeacherMode } from './TeacherMode';

// MediaMetadataのグローバルモック
class MockMediaMetadata {
  title: string;
  artist: string;
  album: string;
  constructor(init: { title?: string; artist?: string; album?: string }) {
    this.title = init.title || '';
    this.artist = init.artist || '';
    this.album = init.album || '';
  }
}
(global as { MediaMetadata?: typeof MockMediaMetadata }).MediaMetadata = MockMediaMetadata;

// モックの設定
const mockUseTeacherMode = vi.fn();
vi.mock('@/hooks/useTeacherMode', () => ({
  useTeacherMode: () => mockUseTeacherMode(),
}));

// モックプレイリスト
const mockPlaylist = {
  id: 'playlist-1',
  bookId: 'test-book',
  pages: [
    { pageNumber: 1, segments: [], totalDuration: 5000 },
    { pageNumber: 2, segments: [], totalDuration: 5000 },
    { pageNumber: 3, segments: [], totalDuration: 5000 },
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
  totalDuration: 15000,
};

describe('TeacherMode', () => {
  beforeEach(() => {
    vi.useFakeTimers();

    // Media Session APIのモック
    Object.defineProperty(global.navigator, 'mediaSession', {
      value: {
        metadata: null,
        setActionHandler: vi.fn(),
      },
      writable: true,
      configurable: true,
    });

    // デフォルトのモック戻り値（プレイリストあり）
    mockUseTeacherMode.mockReturnValue({
      playbackState: {
        status: 'stopped',
        currentPage: 0,
        currentSegmentIndex: 0,
        elapsedTime: 0,
        totalDuration: 0,
      },
      playlist: mockPlaylist,
      loading: false,
      error: null,
      play: vi.fn(),
      pause: vi.fn(),
      stop: vi.fn(),
      next: vi.fn(),
      previous: vi.fn(),
      seekTo: vi.fn(),
    });
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.clearAllMocks();
  });

  it('教師モードの開始ボタンが表示される', () => {
    render(<TeacherMode bookId="test-book" />);
    expect(screen.getByText('開始')).toBeInTheDocument();
  });

  it('開始ボタンをクリックすると教師モードが開始される', () => {
    const mockPlay = vi.fn();
    mockUseTeacherMode.mockReturnValue({
      playbackState: {
        status: 'stopped',
        currentPage: 0,
        currentSegmentIndex: 0,
        elapsedTime: 0,
        totalDuration: 0,
      },
      playlist: mockPlaylist,
      loading: false,
      error: null,
      play: mockPlay,
      pause: vi.fn(),
      stop: vi.fn(),
      next: vi.fn(),
      previous: vi.fn(),
      seekTo: vi.fn(),
    });

    render(<TeacherMode bookId="test-book" />);
    fireEvent.click(screen.getByText('開始'));

    expect(mockPlay).toHaveBeenCalled();
  });

  it('再生中は一時停止ボタンが表示される', () => {
    mockUseTeacherMode.mockReturnValue({
      playbackState: {
        status: 'playing',
        currentPage: 1,
        currentSegmentIndex: 0,
        elapsedTime: 1000,
        totalDuration: 10000,
      },
      playlist: mockPlaylist,
      loading: false,
      error: null,
      play: vi.fn(),
      pause: vi.fn(),
      stop: vi.fn(),
      next: vi.fn(),
      previous: vi.fn(),
      seekTo: vi.fn(),
    });

    render(<TeacherMode bookId="test-book" />);
    expect(screen.getByText('一時停止')).toBeInTheDocument();
  });

  it('一時停止ボタンをクリックすると再生が一時停止される', () => {
    const mockPause = vi.fn();
    mockUseTeacherMode.mockReturnValue({
      playbackState: {
        status: 'playing',
        currentPage: 1,
        currentSegmentIndex: 0,
        elapsedTime: 1000,
        totalDuration: 10000,
      },
      playlist: mockPlaylist,
      loading: false,
      error: null,
      play: vi.fn(),
      pause: mockPause,
      stop: vi.fn(),
      next: vi.fn(),
      previous: vi.fn(),
      seekTo: vi.fn(),
    });

    render(<TeacherMode bookId="test-book" />);
    fireEvent.click(screen.getByText('一時停止'));

    expect(mockPause).toHaveBeenCalled();
  });

  it('一時停止中は再開ボタンが表示される', () => {
    mockUseTeacherMode.mockReturnValue({
      playbackState: {
        status: 'paused',
        currentPage: 1,
        currentSegmentIndex: 0,
        elapsedTime: 1000,
        totalDuration: 10000,
      },
      playlist: mockPlaylist,
      loading: false,
      error: null,
      play: vi.fn(),
      pause: vi.fn(),
      stop: vi.fn(),
      next: vi.fn(),
      previous: vi.fn(),
      seekTo: vi.fn(),
    });

    render(<TeacherMode bookId="test-book" />);
    expect(screen.getByText('再開')).toBeInTheDocument();
  });

  it('現在のページ番号が表示される', () => {
    mockUseTeacherMode.mockReturnValue({
      playbackState: {
        status: 'playing',
        currentPage: 2,
        currentSegmentIndex: 0,
        elapsedTime: 1000,
        totalDuration: 10000,
      },
      playlist: mockPlaylist,
      loading: false,
      error: null,
      play: vi.fn(),
      pause: vi.fn(),
      stop: vi.fn(),
      next: vi.fn(),
      previous: vi.fn(),
      seekTo: vi.fn(),
    });

    render(<TeacherMode bookId="test-book" />);
    expect(screen.getByText(/ページ 2/)).toBeInTheDocument();
  });

  it('前のページボタンをクリックすると前のページに移動する', () => {
    const mockPrevious = vi.fn();
    mockUseTeacherMode.mockReturnValue({
      playbackState: {
        status: 'playing',
        currentPage: 2,
        currentSegmentIndex: 0,
        elapsedTime: 1000,
        totalDuration: 10000,
      },
      playlist: mockPlaylist,
      loading: false,
      error: null,
      play: vi.fn(),
      pause: vi.fn(),
      stop: vi.fn(),
      next: vi.fn(),
      previous: mockPrevious,
      seekTo: vi.fn(),
    });

    render(<TeacherMode bookId="test-book" />);
    fireEvent.click(screen.getByLabelText('前のページ'));

    expect(mockPrevious).toHaveBeenCalled();
  });

  it('次のページボタンをクリックすると次のページに移動する', () => {
    const mockNext = vi.fn();
    mockUseTeacherMode.mockReturnValue({
      playbackState: {
        status: 'playing',
        currentPage: 1,
        currentSegmentIndex: 0,
        elapsedTime: 1000,
        totalDuration: 10000,
      },
      playlist: mockPlaylist,
      loading: false,
      error: null,
      play: vi.fn(),
      pause: vi.fn(),
      stop: vi.fn(),
      next: mockNext,
      previous: vi.fn(),
      seekTo: vi.fn(),
    });

    render(<TeacherMode bookId="test-book" />);
    fireEvent.click(screen.getByLabelText('次のページ'));

    expect(mockNext).toHaveBeenCalled();
  });

  it('停止ボタンをクリックすると教師モードが停止される', () => {
    const mockStop = vi.fn();
    mockUseTeacherMode.mockReturnValue({
      playbackState: {
        status: 'playing',
        currentPage: 1,
        currentSegmentIndex: 0,
        elapsedTime: 1000,
        totalDuration: 10000,
      },
      playlist: mockPlaylist,
      loading: false,
      error: null,
      play: vi.fn(),
      pause: vi.fn(),
      stop: mockStop,
      next: vi.fn(),
      previous: vi.fn(),
      seekTo: vi.fn(),
    });

    render(<TeacherMode bookId="test-book" />);
    fireEvent.click(screen.getByText('停止'));

    expect(mockStop).toHaveBeenCalled();
  });

  it('ローディング中はスピナーが表示される', () => {
    mockUseTeacherMode.mockReturnValue({
      playbackState: {
        status: 'stopped',
        currentPage: 0,
        currentSegmentIndex: 0,
        elapsedTime: 0,
        totalDuration: 0,
      },
      playlist: null,
      loading: true,
      error: null,
      play: vi.fn(),
      pause: vi.fn(),
      stop: vi.fn(),
      next: vi.fn(),
      previous: vi.fn(),
      seekTo: vi.fn(),
    });

    render(<TeacherMode bookId="test-book" />);
    expect(screen.getByRole('status')).toBeInTheDocument();
  });

  it('エラーが発生した場合はエラーメッセージが表示される', () => {
    mockUseTeacherMode.mockReturnValue({
      playbackState: {
        status: 'stopped',
        currentPage: 0,
        currentSegmentIndex: 0,
        elapsedTime: 0,
        totalDuration: 0,
      },
      playlist: null,
      loading: false,
      error: new Error('テストエラー'),
      play: vi.fn(),
      pause: vi.fn(),
      stop: vi.fn(),
      next: vi.fn(),
      previous: vi.fn(),
      seekTo: vi.fn(),
    });

    render(<TeacherMode bookId="test-book" />);
    expect(screen.getByText(/エラーが発生しました/)).toBeInTheDocument();
  });

  it('Media Session APIが設定される', () => {
    const mockSetActionHandler = vi.fn();
    Object.defineProperty(global.navigator, 'mediaSession', {
      value: {
        metadata: null,
        setActionHandler: mockSetActionHandler,
      },
      writable: true,
      configurable: true,
    });

    render(<TeacherMode bookId="test-book" />);

    expect(mockSetActionHandler).toHaveBeenCalledWith('play', expect.any(Function));
    expect(mockSetActionHandler).toHaveBeenCalledWith('pause', expect.any(Function));
    expect(mockSetActionHandler).toHaveBeenCalledWith('previoustrack', expect.any(Function));
    expect(mockSetActionHandler).toHaveBeenCalledWith('nexttrack', expect.any(Function));
  });
});
