/**
 * 教師モードコンポーネント
 */

'use client';

import { useEffect } from 'react';
import { useTeacherMode } from '@/hooks/useTeacherMode';
import type { TeacherModeSettings } from '@/types/teacher-mode';
import { DEFAULT_TEACHER_MODE_SETTINGS } from '@/types/teacher-mode';

/** Props */
export interface TeacherModeProps {
  /** 書籍ID */
  bookId: string;
  /** 教師モード設定 */
  settings?: TeacherModeSettings;
}

/**
 * 教師モードコンポーネント
 */
export function TeacherMode({ bookId, settings = DEFAULT_TEACHER_MODE_SETTINGS as TeacherModeSettings }: TeacherModeProps) {
  const {
    playbackState,
    playlist,
    loading,
    error,
    play,
    pause,
    stop,
    next,
    previous,
  } = useTeacherMode(bookId, settings);

  /**
   * Media Session API設定
   */
  useEffect(() => {
    if (!('mediaSession' in navigator)) {
      return;
    }

    // メタデータ設定
    navigator.mediaSession.metadata = new MediaMetadata({
      title: '教師モード',
      artist: 'HaiLanGo',
      album: bookId,
    });

    // アクションハンドラー設定
    navigator.mediaSession.setActionHandler('play', () => {
      if (playbackState.status === 'paused') {
        play();
      }
    });

    navigator.mediaSession.setActionHandler('pause', () => {
      if (playbackState.status === 'playing') {
        pause();
      }
    });

    navigator.mediaSession.setActionHandler('previoustrack', () => {
      previous();
    });

    navigator.mediaSession.setActionHandler('nexttrack', () => {
      next();
    });

    return () => {
      // クリーンアップ
      navigator.mediaSession.setActionHandler('play', null);
      navigator.mediaSession.setActionHandler('pause', null);
      navigator.mediaSession.setActionHandler('previoustrack', null);
      navigator.mediaSession.setActionHandler('nexttrack', null);
    };
  }, [bookId, playbackState.status, play, pause, previous, next]);

  /**
   * ローディング表示
   */
  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div role="status" className="flex flex-col items-center gap-4">
          <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-500" />
          <p className="text-gray-600">プレイリストを読み込み中...</p>
        </div>
      </div>
    );
  }

  /**
   * エラー表示
   */
  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="bg-red-50 border border-red-200 rounded-lg p-6 max-w-md">
          <h2 className="text-red-800 text-lg font-semibold mb-2">
            エラーが発生しました
          </h2>
          <p className="text-red-600 mb-4">{error.message}</p>
          <button
            type="button"
            onClick={() => window.location.reload()}
            className="bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700"
          >
            再読み込み
          </button>
        </div>
      </div>
    );
  }

  /**
   * プレイリストがない場合
   */
  if (!playlist) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <p className="text-gray-600">プレイリストが見つかりません</p>
      </div>
    );
  }

  const isPlaying = playbackState.status === 'playing';
  const isPaused = playbackState.status === 'paused';
  const isStopped = playbackState.status === 'stopped';
  const isFirstPage = playbackState.currentPage <= 1;
  const isLastPage = playbackState.currentPage >= playlist.pages.length;

  return (
    <div className="flex flex-col h-screen bg-gray-50">
      {/* ヘッダー */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <h1 className="text-xl font-semibold text-gray-900">教師モード</h1>
          <button
            type="button"
            aria-label="設定"
            className="p-2 text-gray-600 hover:bg-gray-100 rounded-full"
          >
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
              />
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
              />
            </svg>
          </button>
        </div>
      </header>

      {/* メインコンテンツ */}
      <main className="flex-1 flex flex-col items-center justify-center p-8">
        {/* ステータス表示 */}
        <div className="mb-8 text-center">
          {isPlaying && (
            <div className="inline-flex items-center gap-2 bg-green-100 text-green-800 px-4 py-2 rounded-full">
              <div className="w-2 h-2 bg-green-600 rounded-full animate-pulse" />
              <span className="font-medium">再生中</span>
            </div>
          )}
          {isPaused && (
            <div className="inline-flex items-center gap-2 bg-yellow-100 text-yellow-800 px-4 py-2 rounded-full">
              <div className="w-2 h-2 bg-yellow-600 rounded-full" />
              <span className="font-medium">一時停止中</span>
            </div>
          )}
        </div>

        {/* ページ表示 */}
        {!isStopped && playbackState.currentPage > 0 && (
          <div className="text-center mb-8">
            <p className="text-4xl font-bold text-gray-900">
              ページ {playbackState.currentPage}
            </p>
            <p className="text-gray-600 mt-2">
              全 {playlist.pages.length} ページ中
            </p>
          </div>
        )}

        {/* 進捗バー */}
        {!isStopped && (
          <div className="w-full max-w-md mb-8">
            <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
              <div
                className="h-full bg-blue-500 transition-all duration-300"
                style={{
                  width: `${(playbackState.currentPage / playlist.pages.length) * 100}%`,
                }}
              />
            </div>
          </div>
        )}

        {/* コントロールボタン */}
        <div className="flex items-center gap-4">
          {/* 前のページ */}
          <button
            type="button"
            onClick={previous}
            disabled={isFirstPage || isStopped}
            aria-label="前のページ"
            className="p-3 rounded-full bg-white shadow-md hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed transition-all"
          >
            <svg
              className="w-6 h-6 text-gray-700"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15 19l-7-7 7-7"
              />
            </svg>
          </button>

          {/* 再生/一時停止ボタン */}
          {isStopped && (
            <button
              type="button"
              onClick={play}
              className="p-6 rounded-full bg-blue-500 text-white shadow-lg hover:shadow-xl hover:bg-blue-600 transition-all"
            >
              <svg
                className="w-8 h-8"
                fill="currentColor"
                viewBox="0 0 24 24"
              >
                <path d="M8 5v14l11-7z" />
              </svg>
              <span className="sr-only">開始</span>
            </button>
          )}

          {isPlaying && (
            <button
              type="button"
              onClick={pause}
              className="p-6 rounded-full bg-blue-500 text-white shadow-lg hover:shadow-xl hover:bg-blue-600 transition-all"
            >
              <svg
                className="w-8 h-8"
                fill="currentColor"
                viewBox="0 0 24 24"
              >
                <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z" />
              </svg>
              <span className="sr-only">一時停止</span>
            </button>
          )}

          {isPaused && (
            <button
              type="button"
              onClick={play}
              className="p-6 rounded-full bg-blue-500 text-white shadow-lg hover:shadow-xl hover:bg-blue-600 transition-all"
            >
              <svg
                className="w-8 h-8"
                fill="currentColor"
                viewBox="0 0 24 24"
              >
                <path d="M8 5v14l11-7z" />
              </svg>
              <span className="sr-only">再開</span>
            </button>
          )}

          {/* 次のページ */}
          <button
            type="button"
            onClick={next}
            disabled={isLastPage || isStopped}
            aria-label="次のページ"
            className="p-3 rounded-full bg-white shadow-md hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed transition-all"
          >
            <svg
              className="w-6 h-6 text-gray-700"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 5l7 7-7 7"
              />
            </svg>
          </button>
        </div>

        {/* 停止ボタン */}
        {!isStopped && (
          <button
            type="button"
            onClick={stop}
            className="mt-8 px-6 py-3 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
          >
            停止
          </button>
        )}
      </main>
    </div>
  );
}
