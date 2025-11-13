'use client';

import type { NotificationSettings as NotificationSettingsType } from '@/types/settings';

interface NotificationSettingsProps {
  settings: NotificationSettingsType;
  onUpdate: (settings: NotificationSettingsType) => Promise<void>;
}

export default function NotificationSettings({ settings, onUpdate }: NotificationSettingsProps) {
  const handleToggle = async (key: keyof NotificationSettingsType) => {
    const newSettings = {
      ...settings,
      [key]: !settings[key],
    };
    await onUpdate(newSettings);
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <h2 className="text-xl font-semibold mb-4">通知設定</h2>

      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <label htmlFor="learningReminder" className="text-sm font-medium">
            学習リマインダー
          </label>
          <input
            id="learningReminder"
            type="checkbox"
            checked={settings.learningReminder}
            onChange={() => handleToggle('learningReminder')}
            className="w-12 h-6 appearance-none bg-gray-300 rounded-full relative cursor-pointer transition-colors checked:bg-primary"
            style={{
              WebkitAppearance: 'none',
            }}
          />
        </div>

        <div className="flex items-center justify-between">
          <label htmlFor="reviewNotification" className="text-sm font-medium">
            復習通知
          </label>
          <input
            id="reviewNotification"
            type="checkbox"
            checked={settings.reviewNotification}
            onChange={() => handleToggle('reviewNotification')}
            className="w-12 h-6 appearance-none bg-gray-300 rounded-full relative cursor-pointer transition-colors checked:bg-primary"
            style={{
              WebkitAppearance: 'none',
            }}
          />
        </div>

        <div className="flex items-center justify-between">
          <label htmlFor="emailNotification" className="text-sm font-medium">
            メール通知
          </label>
          <input
            id="emailNotification"
            type="checkbox"
            checked={settings.emailNotification}
            onChange={() => handleToggle('emailNotification')}
            className="w-12 h-6 appearance-none bg-gray-300 rounded-full relative cursor-pointer transition-colors checked:bg-primary"
            style={{
              WebkitAppearance: 'none',
            }}
          />
        </div>
      </div>
    </div>
  );
}
