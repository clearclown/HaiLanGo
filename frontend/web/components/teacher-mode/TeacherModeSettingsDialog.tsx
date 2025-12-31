/**
 * 教師モード設定ダイアログ
 */

'use client';

import type {
  AudioQuality,
  PlaybackSpeed,
  RepeatCount,
  TeacherModeSettings,
} from '@/types/teacher-mode';
import { useState } from 'react';

interface TeacherModeSettingsDialogProps {
  /** 表示状態 */
  isOpen: boolean;
  /** 閉じるコールバック */
  onClose: () => void;
  /** 設定 */
  settings: TeacherModeSettings;
  /** 設定変更コールバック */
  onSettingsChange: (settings: TeacherModeSettings) => void;
}

const SPEED_OPTIONS: { value: PlaybackSpeed; label: string }[] = [
  { value: 0.5, label: '0.5x' },
  { value: 0.75, label: '0.75x' },
  { value: 1.0, label: '1.0x' },
  { value: 1.25, label: '1.25x' },
  { value: 1.5, label: '1.5x' },
  { value: 2.0, label: '2.0x' },
];

const REPEAT_OPTIONS: { value: RepeatCount; label: string }[] = [
  { value: 1, label: '1回' },
  { value: 2, label: '2回' },
  { value: 3, label: '3回' },
];

const QUALITY_OPTIONS: { value: AudioQuality; label: string }[] = [
  { value: 'standard', label: '標準' },
  { value: 'premium', label: '高品質' },
];

export function TeacherModeSettingsDialog({
  isOpen,
  onClose,
  settings,
  onSettingsChange,
}: TeacherModeSettingsDialogProps) {
  const [localSettings, setLocalSettings] = useState<TeacherModeSettings>(settings);

  if (!isOpen) {
    return null;
  }

  const handleSave = () => {
    onSettingsChange(localSettings);
    onClose();
  };

  const handleContentToggle = (key: keyof TeacherModeSettings['content']) => {
    setLocalSettings((prev) => ({
      ...prev,
      content: {
        ...prev.content,
        [key]: !prev.content[key],
      },
    }));
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* オーバーレイ */}
      <div
        className="absolute inset-0 bg-black/50"
        onClick={onClose}
        onKeyDown={(e) => e.key === 'Escape' && onClose()}
        role="button"
        tabIndex={0}
        aria-label="閉じる"
      />

      {/* ダイアログ */}
      <dialog
        open
        className="relative bg-white rounded-xl shadow-2xl max-w-md w-full mx-4 max-h-[90vh] overflow-y-auto m-0"
        aria-labelledby="settings-title"
      >
        {/* ヘッダー */}
        <div className="sticky top-0 bg-white border-b px-6 py-4 flex items-center justify-between">
          <h2 id="settings-title" className="text-xl font-semibold text-gray-900">
            教師モード設定
          </h2>
          <button
            type="button"
            onClick={onClose}
            className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-full transition-colors"
            aria-label="閉じる"
          >
            <svg
              aria-hidden="true"
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        {/* コンテンツ */}
        <div className="p-6 space-y-6">
          {/* 再生速度 */}
          <fieldset>
            <legend className="block text-sm font-medium text-gray-700 mb-2">再生速度</legend>
            <div className="flex flex-wrap gap-2">
              {SPEED_OPTIONS.map((option) => (
                <button
                  key={option.value}
                  type="button"
                  onClick={() => setLocalSettings((prev) => ({ ...prev, speed: option.value }))}
                  className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                    localSettings.speed === option.value
                      ? 'bg-blue-500 text-white'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  }`}
                >
                  {option.label}
                </button>
              ))}
            </div>
          </fieldset>

          {/* ページ間隔 */}
          <div>
            <label htmlFor="page-interval" className="block text-sm font-medium text-gray-700 mb-2">
              ページ間隔: {localSettings.pageInterval}秒
            </label>
            <input
              id="page-interval"
              type="range"
              min={0}
              max={30}
              value={localSettings.pageInterval}
              onChange={(e) =>
                setLocalSettings((prev) => ({
                  ...prev,
                  pageInterval: Number(e.target.value),
                }))
              }
              className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-blue-500"
            />
            <div className="flex justify-between text-xs text-gray-500 mt-1">
              <span>0秒</span>
              <span>30秒</span>
            </div>
          </div>

          {/* リピート回数 */}
          <fieldset>
            <legend className="block text-sm font-medium text-gray-700 mb-2">リピート回数</legend>
            <div className="flex gap-2">
              {REPEAT_OPTIONS.map((option) => (
                <button
                  key={option.value}
                  type="button"
                  onClick={() =>
                    setLocalSettings((prev) => ({ ...prev, repeatCount: option.value }))
                  }
                  className={`flex-1 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                    localSettings.repeatCount === option.value
                      ? 'bg-blue-500 text-white'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  }`}
                >
                  {option.label}
                </button>
              ))}
            </div>
          </fieldset>

          {/* 音質 */}
          <fieldset>
            <legend className="block text-sm font-medium text-gray-700 mb-2">音質</legend>
            <div className="flex gap-2">
              {QUALITY_OPTIONS.map((option) => (
                <button
                  key={option.value}
                  type="button"
                  onClick={() =>
                    setLocalSettings((prev) => ({ ...prev, audioQuality: option.value }))
                  }
                  className={`flex-1 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                    localSettings.audioQuality === option.value
                      ? 'bg-blue-500 text-white'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  }`}
                >
                  {option.label}
                  {option.value === 'premium' && <span className="ml-1 text-xs">👑</span>}
                </button>
              ))}
            </div>
          </fieldset>

          {/* 学習内容 */}
          <fieldset>
            <legend className="block text-sm font-medium text-gray-700 mb-3">学習内容</legend>
            <div className="space-y-3">
              <label className="flex items-center gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  checked={localSettings.content.includeTranslation}
                  onChange={() => handleContentToggle('includeTranslation')}
                  className="w-5 h-5 rounded border-gray-300 text-blue-500 focus:ring-blue-500"
                />
                <span className="text-gray-700">母国語訳を含む</span>
              </label>

              <label className="flex items-center gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  checked={localSettings.content.includeWordExplanation}
                  onChange={() => handleContentToggle('includeWordExplanation')}
                  className="w-5 h-5 rounded border-gray-300 text-blue-500 focus:ring-blue-500"
                />
                <span className="text-gray-700">単語解説を含む</span>
              </label>

              <label className="flex items-center gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  checked={localSettings.content.includeGrammarExplanation}
                  onChange={() => handleContentToggle('includeGrammarExplanation')}
                  className="w-5 h-5 rounded border-gray-300 text-blue-500 focus:ring-blue-500"
                />
                <span className="text-gray-700">文法解説を含む</span>
              </label>

              <label className="flex items-center gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  checked={localSettings.content.includePronunciationPractice}
                  onChange={() => handleContentToggle('includePronunciationPractice')}
                  className="w-5 h-5 rounded border-gray-300 text-blue-500 focus:ring-blue-500"
                />
                <span className="text-gray-700">発音練習を含む</span>
              </label>

              <label className="flex items-center gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  checked={localSettings.content.includeExampleSentences}
                  onChange={() => handleContentToggle('includeExampleSentences')}
                  className="w-5 h-5 rounded border-gray-300 text-blue-500 focus:ring-blue-500"
                />
                <span className="text-gray-700">例文を含む</span>
              </label>
            </div>
          </fieldset>
        </div>

        {/* フッター */}
        <div className="sticky bottom-0 bg-white border-t px-6 py-4 flex gap-3">
          <button
            type="button"
            onClick={onClose}
            className="flex-1 px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors font-medium"
          >
            キャンセル
          </button>
          <button
            type="button"
            onClick={handleSave}
            className="flex-1 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors font-medium"
          >
            保存
          </button>
        </div>
      </dialog>
    </div>
  );
}
