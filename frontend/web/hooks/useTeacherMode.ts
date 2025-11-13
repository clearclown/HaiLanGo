/**
 * 教師モードカスタムフック
 */

import { useState, useEffect, useCallback, useRef } from 'react';
import { teacherModeApi } from '@/services/teacherModeApi';
import type {
  TeacherModePlaylist,
  PlaybackState,
  TeacherModeSettings,
  DEFAULT_TEACHER_MODE_SETTINGS,
} from '@/types/teacher-mode';

/** useTeacherModeの戻り値 */
export interface UseTeacherModeReturn {
  /** 再生状態 */
  playbackState: PlaybackState;
  /** プレイリスト */
  playlist: TeacherModePlaylist | null;
  /** ローディング中 */
  loading: boolean;
  /** エラー */
  error: Error | null;
  /** 再生開始 */
  play: () => void;
  /** 一時停止 */
  pause: () => void;
  /** 停止 */
  stop: () => void;
  /** 次のページへ */
  next: () => void;
  /** 前のページへ */
  previous: () => void;
  /** 指定ページへシーク */
  seekTo: (pageNumber: number) => void;
}

/**
 * 教師モードフック
 */
export function useTeacherMode(
  bookId: string,
  settings: TeacherModeSettings = DEFAULT_TEACHER_MODE_SETTINGS as TeacherModeSettings,
): UseTeacherModeReturn {
  // 状態管理
  const [playbackState, setPlaybackState] = useState<PlaybackState>({
    status: 'stopped',
    currentPage: 0,
    currentSegmentIndex: 0,
    elapsedTime: 0,
    totalDuration: 0,
  });

  const [playlist, setPlaylist] = useState<TeacherModePlaylist | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  // タイマー管理用のRef
  const timerRef = useRef<NodeJS.Timeout | null>(null);
  const audioRef = useRef<HTMLAudioElement | null>(null);

  /**
   * プレイリストを取得
   */
  useEffect(() => {
    let mounted = true;

    const fetchPlaylist = async () => {
      try {
        setLoading(true);
        setError(null);

        const data = await teacherModeApi.fetchPlaylist(bookId, settings);

        if (mounted) {
          setPlaylist(data);
          setPlaybackState((prev) => ({
            ...prev,
            totalDuration: data.totalDuration,
          }));
        }
      } catch (err) {
        if (mounted) {
          setError(err instanceof Error ? err : new Error('Failed to fetch playlist'));
        }
      } finally {
        if (mounted) {
          setLoading(false);
        }
      }
    };

    fetchPlaylist();

    return () => {
      mounted = false;
    };
  }, [bookId, settings]);

  /**
   * 音声を再生
   */
  const playAudio = useCallback((audioUrl: string, duration: number) => {
    if (audioRef.current) {
      audioRef.current.pause();
    }

    audioRef.current = new Audio(audioUrl);
    audioRef.current.playbackRate = settings.speed;
    audioRef.current.play();

    // セグメント再生完了後に次のセグメントへ
    audioRef.current.onended = () => {
      // タイマーをセット（ページ間隔）
      timerRef.current = setTimeout(() => {
        setPlaybackState((prev) => {
          const currentPageData = playlist?.pages.find(
            (p) => p.pageNumber === prev.currentPage,
          );

          if (!currentPageData) return prev;

          // 次のセグメントがあるか確認
          if (prev.currentSegmentIndex < currentPageData.segments.length - 1) {
            return {
              ...prev,
              currentSegmentIndex: prev.currentSegmentIndex + 1,
            };
          }

          // 次のページがあるか確認
          if (playlist && prev.currentPage < playlist.pages.length) {
            return {
              ...prev,
              currentPage: prev.currentPage + 1,
              currentSegmentIndex: 0,
            };
          }

          // 最後のページなので停止
          return {
            ...prev,
            status: 'stopped',
            currentPage: 0,
            currentSegmentIndex: 0,
          };
        });
      }, settings.pageInterval * 1000);
    };
  }, [playlist, settings.speed, settings.pageInterval]);

  /**
   * 現在のセグメントを再生
   */
  useEffect(() => {
    if (playbackState.status !== 'playing' || !playlist) {
      return;
    }

    const currentPageData = playlist.pages.find(
      (p) => p.pageNumber === playbackState.currentPage,
    );

    if (!currentPageData) {
      return;
    }

    const currentSegment = currentPageData.segments[playbackState.currentSegmentIndex];

    if (currentSegment) {
      playAudio(currentSegment.audioUrl, currentSegment.duration);
    }
  }, [
    playbackState.status,
    playbackState.currentPage,
    playbackState.currentSegmentIndex,
    playlist,
    playAudio,
  ]);

  /**
   * 再生開始
   */
  const play = useCallback(() => {
    if (!playlist || playlist.pages.length === 0) {
      return;
    }

    setPlaybackState((prev) => {
      // 停止状態から開始する場合は最初のページから
      if (prev.status === 'stopped') {
        return {
          ...prev,
          status: 'playing',
          currentPage: 1,
          currentSegmentIndex: 0,
        };
      }

      // 一時停止からの再開
      return {
        ...prev,
        status: 'playing',
      };
    });
  }, [playlist]);

  /**
   * 一時停止
   */
  const pause = useCallback(() => {
    if (audioRef.current) {
      audioRef.current.pause();
    }

    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = null;
    }

    setPlaybackState((prev) => ({
      ...prev,
      status: 'paused',
    }));
  }, []);

  /**
   * 停止
   */
  const stop = useCallback(() => {
    if (audioRef.current) {
      audioRef.current.pause();
      audioRef.current = null;
    }

    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = null;
    }

    setPlaybackState((prev) => ({
      ...prev,
      status: 'stopped',
      currentPage: 0,
      currentSegmentIndex: 0,
      elapsedTime: 0,
    }));
  }, []);

  /**
   * 次のページへ
   */
  const next = useCallback(() => {
    if (!playlist) return;

    setPlaybackState((prev) => {
      // 最後のページかチェック
      if (prev.currentPage >= playlist.pages.length) {
        return {
          ...prev,
          status: 'stopped',
          currentPage: 0,
          currentSegmentIndex: 0,
        };
      }

      return {
        ...prev,
        currentPage: prev.currentPage + 1,
        currentSegmentIndex: 0,
      };
    });
  }, [playlist]);

  /**
   * 前のページへ
   */
  const previous = useCallback(() => {
    setPlaybackState((prev) => {
      // 最初のページかチェック
      if (prev.currentPage <= 1) {
        return prev;
      }

      return {
        ...prev,
        currentPage: prev.currentPage - 1,
        currentSegmentIndex: 0,
      };
    });
  }, []);

  /**
   * 指定ページへシーク
   */
  const seekTo = useCallback(
    (pageNumber: number) => {
      if (!playlist) return;

      if (pageNumber < 1 || pageNumber > playlist.pages.length) {
        return;
      }

      setPlaybackState((prev) => ({
        ...prev,
        currentPage: pageNumber,
        currentSegmentIndex: 0,
      }));
    },
    [playlist],
  );

  /**
   * クリーンアップ
   */
  useEffect(() => {
    return () => {
      if (audioRef.current) {
        audioRef.current.pause();
        audioRef.current = null;
      }

      if (timerRef.current) {
        clearTimeout(timerRef.current);
        timerRef.current = null;
      }
    };
  }, []);

  return {
    playbackState,
    playlist,
    loading,
    error,
    play,
    pause,
    stop,
    next,
    previous,
    seekTo,
  };
}
