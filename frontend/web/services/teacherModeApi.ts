/**
 * 教師モードAPIサービス
 */

import type {
  TeacherModePlaylist,
  TeacherModeSettings,
  PlaybackState,
} from '@/types/teacher-mode';

/** APIベースURL */
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

/** APIレスポンスエラー */
export class ApiError extends Error {
  constructor(
    message: string,
    public statusCode: number,
    public response?: unknown,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

/**
 * APIリクエストを実行する
 */
async function fetchApi<T>(
  endpoint: string,
  options: RequestInit = {},
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  };

  // 認証トークンがあれば追加
  const token = localStorage.getItem('auth_token');
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  try {
    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new ApiError(
        errorData.message || `HTTP Error ${response.status}`,
        response.status,
        errorData,
      );
    }

    return await response.json();
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }
    throw new ApiError(
      error instanceof Error ? error.message : 'Network error',
      0,
    );
  }
}

/**
 * プレイリスト生成APIレスポンス
 */
interface GeneratePlaylistResponse {
  playlistId: string;
  totalPages: number;
  estimatedDuration: number;
  pages: Array<{
    pageNumber: number;
    segments: Array<{
      type: string;
      audioUrl: string;
      duration: number;
      text: string;
    }>;
  }>;
}

/**
 * ダウンロードパッケージAPIレスポンス
 */
interface GenerateDownloadPackageResponse {
  packageId: string;
  downloadUrl: string;
  totalSize: number;
  expiresAt: string;
}

/**
 * 教師モードAPIクライアント
 */
export const teacherModeApi = {
  /**
   * プレイリストを取得
   */
  async fetchPlaylist(
    bookId: string,
    settings: TeacherModeSettings,
  ): Promise<TeacherModePlaylist> {
    const response = await fetchApi<GeneratePlaylistResponse>(
      `/api/v1/books/${bookId}/teacher-mode/generate`,
      {
        method: 'POST',
        body: JSON.stringify({
          settings,
          pageRange: {
            start: 1,
            end: 999,
          },
        }),
      },
    );

    // レスポンスを内部形式に変換
    return {
      id: response.playlistId,
      bookId,
      pages: response.pages.map((page) => ({
        pageNumber: page.pageNumber,
        segments: page.segments.map((seg, idx) => ({
          id: `segment-${page.pageNumber}-${idx}`,
          type: seg.type as any,
          audioUrl: seg.audioUrl,
          duration: seg.duration,
          text: seg.text,
          language: 'en', // TODO: 実際の言語を返すようにバックエンド修正
        })),
        totalDuration: page.segments.reduce((sum, seg) => sum + seg.duration, 0),
      })),
      settings,
      totalDuration: response.estimatedDuration,
    };
  },

  /**
   * オフラインダウンロードパッケージを生成
   */
  async generateDownloadPackage(
    bookId: string,
    settings: TeacherModeSettings,
  ): Promise<GenerateDownloadPackageResponse> {
    return fetchApi<GenerateDownloadPackageResponse>(
      `/api/v1/books/${bookId}/teacher-mode/download-package`,
      {
        method: 'POST',
        body: JSON.stringify({ settings }),
      },
    );
  },

  /**
   * 再生状態を保存
   */
  async updatePlaybackState(
    bookId: string,
    state: Partial<PlaybackState>,
  ): Promise<void> {
    await fetchApi(`/api/v1/books/${bookId}/teacher-mode/playback-state`, {
      method: 'PUT',
      body: JSON.stringify({
        currentPage: state.currentPage,
        currentSegmentIndex: state.currentSegmentIndex,
        elapsedTime: state.elapsedTime,
      }),
    });
  },

  /**
   * 保存された再生状態を取得
   */
  async getPlaybackState(bookId: string): Promise<PlaybackState | null> {
    try {
      return await fetchApi<PlaybackState>(
        `/api/v1/books/${bookId}/teacher-mode/playback-state`,
      );
    } catch (error) {
      if (error instanceof ApiError && error.statusCode === 404) {
        return null;
      }
      throw error;
    }
  },
};
