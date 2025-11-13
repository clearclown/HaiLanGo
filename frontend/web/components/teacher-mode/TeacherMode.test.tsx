/**
 * 教師モードコンポーネントのテスト
 */

import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { TeacherMode } from './TeacherMode';

// モックの設定
const mockUseTeacherMode = vi.fn();
vi.mock('@/hooks/useTeacherMode', () => ({
  useTeacherMode: () => mockUseTeacherMode(),
}));

describe('TeacherMode', () => {
  beforeEach(() => {
    vi.useFakeTimers();

    // デフォルトのモック戻り値
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

  it('開始ボタンをクリックすると教師モードが開始される', async () => {
    const mockPlay = vi.fn();
    mockUseTeacherMode.mockReturnValue({
      ...mockUseTeacherMode(),
      play: mockPlay,
    });

    render(<TeacherMode bookId="test-book" />);
    fireEvent.click(screen.getByText('開始'));

    await waitFor(() => {
      expect(mockPlay).toHaveBeenCalled();
    });
  });

  it('再生中は一時停止ボタンが表示される', () => {
    mockUseTeacherMode.mockReturnValue({
      ...mockUseTeacherMode(),
      playbackState: {
        status: 'playing',
        currentPage: 1,
        currentSegmentIndex: 0,
        elapsedTime: 1000,
        totalDuration: 10000,
      },
    });

    render(<TeacherMode bookId="test-book" />);
    expect(screen.getByText('一時停止')).toBeInTheDocument();
  });

  it('一時停止ボタンをクリックすると再生が一時停止される', async () => {
    const mockPause = vi.fn();
    mockUseTeacherMode.mockReturnValue({
      ...mockUseTeacherMode(),
      playbackState: {
        status: 'playing',
        currentPage: 1,
        currentSegmentIndex: 0,
        elapsedTime: 1000,
        totalDuration: 10000,
      },
      pause: mockPause,
    });

    render(<TeacherMode bookId="test-book" />);
    fireEvent.click(screen.getByText('一時停止'));

    await waitFor(() => {
      expect(mockPause).toHaveBeenCalled();
    });
  });

  it('一時停止中は再開ボタンが表示される', () => {
    mockUseTeacherMode.mockReturnValue({
      ...mockUseTeacherMode(),
      playbackState: {
        status: 'paused',
        currentPage: 1,
        currentSegmentIndex: 0,
        elapsedTime: 1000,
        totalDuration: 10000,
      },
    });

    render(<TeacherMode bookId="test-book" />);
    expect(screen.getByText('再開')).toBeInTheDocument();
  });

  it('現在のページ番号が表示される', () => {
    mockUseTeacherMode.mockReturnValue({
      ...mockUseTeacherMode(),
      playbackState: {
        status: 'playing',
        currentPage: 12,
        currentSegmentIndex: 0,
        elapsedTime: 1000,
        totalDuration: 10000,
      },
    });

    render(<TeacherMode bookId="test-book" />);
    expect(screen.getByText(/ページ 12/)).toBeInTheDocument();
  });

  it('前のページボタンをクリックすると前のページに移動する', async () => {
    const mockPrevious = vi.fn();
    mockUseTeacherMode.mockReturnValue({
      ...mockUseTeacherMode(),
      playbackState: {
        status: 'playing',
        currentPage: 2,
        currentSegmentIndex: 0,
        elapsedTime: 1000,
        totalDuration: 10000,
      },
      previous: mockPrevious,
    });

    render(<TeacherMode bookId="test-book" />);
    fireEvent.click(screen.getByLabelText('前のページ'));

    await waitFor(() => {
      expect(mockPrevious).toHaveBeenCalled();
    });
  });

  it('次のページボタンをクリックすると次のページに移動する', async () => {
    const mockNext = vi.fn();
    mockUseTeacherMode.mockReturnValue({
      ...mockUseTeacherMode(),
      playbackState: {
        status: 'playing',
        currentPage: 1,
        currentSegmentIndex: 0,
        elapsedTime: 1000,
        totalDuration: 10000,
      },
      next: mockNext,
    });

    render(<TeacherMode bookId="test-book" />);
    fireEvent.click(screen.getByLabelText('次のページ'));

    await waitFor(() => {
      expect(mockNext).toHaveBeenCalled();
    });
  });

  it('停止ボタンをクリックすると教師モードが停止される', async () => {
    const mockStop = vi.fn();
    mockUseTeacherMode.mockReturnValue({
      ...mockUseTeacherMode(),
      playbackState: {
        status: 'playing',
        currentPage: 1,
        currentSegmentIndex: 0,
        elapsedTime: 1000,
        totalDuration: 10000,
      },
      stop: mockStop,
    });

    render(<TeacherMode bookId="test-book" />);
    fireEvent.click(screen.getByText('停止'));

    await waitFor(() => {
      expect(mockStop).toHaveBeenCalled();
    });
  });

  it('ローディング中はスピナーが表示される', () => {
    mockUseTeacherMode.mockReturnValue({
      ...mockUseTeacherMode(),
      loading: true,
    });

    render(<TeacherMode bookId="test-book" />);
    expect(screen.getByRole('status')).toBeInTheDocument();
  });

  it('エラーが発生した場合はエラーメッセージが表示される', () => {
    mockUseTeacherMode.mockReturnValue({
      ...mockUseTeacherMode(),
      error: new Error('テストエラー'),
    });

    render(<TeacherMode bookId="test-book" />);
    expect(screen.getByText(/エラーが発生しました/)).toBeInTheDocument();
  });

  it('Media Session APIが設定される', () => {
    const mockSetActionHandler = vi.fn();
    global.navigator.mediaSession.setActionHandler = mockSetActionHandler;

    render(<TeacherMode bookId="test-book" />);

    expect(mockSetActionHandler).toHaveBeenCalledWith('play', expect.any(Function));
    expect(mockSetActionHandler).toHaveBeenCalledWith('pause', expect.any(Function));
    expect(mockSetActionHandler).toHaveBeenCalledWith('previoustrack', expect.any(Function));
    expect(mockSetActionHandler).toHaveBeenCalledWith('nexttrack', expect.any(Function));
  });
});
