/**
 * 教師モードの型定義
 */

/** 音声セグメントのタイプ */
export type AudioSegmentType = 'phrase' | 'translation' | 'explanation' | 'pause';

/** 再生状態 */
export type PlaybackStatus = 'stopped' | 'playing' | 'paused';

/** 音声品質 */
export type AudioQuality = 'standard' | 'premium';

/** 再生速度 */
export type PlaybackSpeed = 0.5 | 0.75 | 1.0 | 1.25 | 1.5 | 2.0;

/** リピート回数 */
export type RepeatCount = 1 | 2 | 3;

/** 教師モード設定 */
export interface TeacherModeSettings {
  /** 再生速度 */
  speed: PlaybackSpeed;
  /** ページ間隔（秒） */
  pageInterval: number;
  /** リピート回数 */
  repeatCount: RepeatCount;
  /** 音質 */
  audioQuality: AudioQuality;
  /** 学習内容 */
  content: {
    /** 母国語訳を含む */
    includeTranslation: boolean;
    /** 単語解説を含む */
    includeWordExplanation: boolean;
    /** 文法解説を含む */
    includeGrammarExplanation: boolean;
    /** 発音練習を含む */
    includePronunciationPractice: boolean;
    /** 例文を含む */
    includeExampleSentences: boolean;
  };
}

/** 音声セグメント */
export interface AudioSegment {
  /** ID */
  id: string;
  /** タイプ */
  type: AudioSegmentType;
  /** 音声URL */
  audioUrl: string;
  /** 長さ（ミリ秒） */
  duration: number;
  /** テキスト */
  text: string;
  /** 言語 */
  language: string;
}

/** ページ音声 */
export interface PageAudio {
  /** ページ番号 */
  pageNumber: number;
  /** セグメント */
  segments: AudioSegment[];
  /** 総長さ（ミリ秒） */
  totalDuration: number;
}

/** 教師モードプレイリスト */
export interface TeacherModePlaylist {
  /** プレイリストID */
  id: string;
  /** 書籍ID */
  bookId: string;
  /** ページ */
  pages: PageAudio[];
  /** 設定 */
  settings: TeacherModeSettings;
  /** 総長さ（ミリ秒） */
  totalDuration: number;
}

/** 再生状態 */
export interface PlaybackState {
  /** 状態 */
  status: PlaybackStatus;
  /** 現在のページ */
  currentPage: number;
  /** 現在のセグメントインデックス */
  currentSegmentIndex: number;
  /** 経過時間（ミリ秒） */
  elapsedTime: number;
  /** 総長さ（ミリ秒） */
  totalDuration: number;
}

/** デフォルト設定 */
export const DEFAULT_TEACHER_MODE_SETTINGS: TeacherModeSettings = {
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
};
