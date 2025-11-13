'use client';

import AccountSettings from '@/components/settings/AccountSettings';
import NotificationSettings from '@/components/settings/NotificationSettings';
import PlanSettings from '@/components/settings/PlanSettings';
import { apiClient } from '@/lib/api/client';
import type {
  NotificationSettings as NotificationSettingsType,
  UserProfile,
} from '@/types/settings';
import { useEffect, useState } from 'react';

export default function SettingsPage() {
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [notifications, setNotifications] = useState<NotificationSettingsType | null>(null);
  const [plan, setPlan] = useState<{ type: 'free' | 'premium'; expiresAt?: string } | null>(null);
  const [showLogoutDialog, setShowLogoutDialog] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const loadSettings = async () => {
      try {
        const [settingsData, planData] = await Promise.all([
          apiClient.settings.get(),
          apiClient.plan.get(),
        ]);

        setProfile(settingsData.profile);
        setNotifications(settingsData.notifications);
        setPlan(planData);
      } catch (error) {
        console.error('設定の読み込みに失敗しました', error);
      } finally {
        setIsLoading(false);
      }
    };

    loadSettings();
  }, []);

  const handleUpdateProfile = async (updatedProfile: Partial<UserProfile>) => {
    try {
      await apiClient.settings.updateProfile(updatedProfile);
      setProfile((prev) => (prev ? { ...prev, ...updatedProfile } : null));
    } catch (error) {
      console.error('プロフィールの更新に失敗しました', error);
      throw error;
    }
  };

  const handleUpdateNotifications = async (updatedNotifications: NotificationSettingsType) => {
    try {
      await apiClient.settings.updateNotifications(updatedNotifications);
      setNotifications(updatedNotifications);
    } catch (error) {
      console.error('通知設定の更新に失敗しました', error);
      throw error;
    }
  };

  const handleUpgrade = async () => {
    try {
      const { checkoutUrl } = await apiClient.plan.upgrade();
      window.location.href = checkoutUrl;
    } catch (error) {
      console.error('アップグレードに失敗しました', error);
    }
  };

  const handleLogout = async () => {
    try {
      await apiClient.auth.logout();
      localStorage.clear();
      window.location.href = '/login';
    } catch (error) {
      console.error('ログアウトに失敗しました', error);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background-secondary flex items-center justify-center">
        <div className="text-text-secondary">読み込み中...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background-secondary">
      <div className="max-w-4xl mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-8">設定</h1>

        <div className="space-y-6">
          {profile && <AccountSettings profile={profile} onUpdate={handleUpdateProfile} />}

          {plan && <PlanSettings plan={plan} onUpgrade={handleUpgrade} />}

          {notifications && (
            <NotificationSettings settings={notifications} onUpdate={handleUpdateNotifications} />
          )}

          <div className="bg-white rounded-lg shadow-sm p-6">
            <h2 className="text-xl font-semibold mb-4">その他</h2>
            <div className="space-y-2">
              <a href="/help" className="block text-primary hover:underline">
                ヘルプ・サポート
              </a>
              <a href="/terms" className="block text-primary hover:underline">
                利用規約
              </a>
              <a href="/privacy" className="block text-primary hover:underline">
                プライバシーポリシー
              </a>
              <button
                type="button"
                onClick={() => setShowLogoutDialog(true)}
                className="text-error hover:underline"
              >
                ログアウト
              </button>
            </div>
          </div>
        </div>

        {showLogoutDialog && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
            <div className="bg-white rounded-lg p-6 max-w-sm">
              <h3 className="text-lg font-semibold mb-4">ログアウトしますか？</h3>
              <div className="flex gap-4">
                <button
                  type="button"
                  onClick={() => setShowLogoutDialog(false)}
                  className="flex-1 bg-gray-200 text-text-primary px-4 py-2 rounded-lg"
                >
                  キャンセル
                </button>
                <button
                  type="button"
                  onClick={handleLogout}
                  className="flex-1 bg-error text-white px-4 py-2 rounded-lg"
                >
                  ログアウト
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
