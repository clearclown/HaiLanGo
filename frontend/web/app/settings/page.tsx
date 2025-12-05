'use client';

import { AppLayout } from '@/components/layout';
import { Button, Card, CardContent, CardHeader, CardTitle } from '@/components/ui';
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
      <AppLayout>
        <div className="container-app py-6 lg:py-8 flex items-center justify-center min-h-[60vh]">
          <div className="text-center">
            <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-primary mx-auto mb-4" />
            <p className="text-gray-600">読み込み中...</p>
          </div>
        </div>
      </AppLayout>
    );
  }

  return (
    <AppLayout>
      <div className="container-app py-6 lg:py-8">
        <h1 className="text-2xl lg:text-3xl font-bold text-gray-900 mb-6 lg:mb-8">設定</h1>

        <div className="space-y-6">
          {profile && <AccountSettings profile={profile} onUpdate={handleUpdateProfile} />}

          {plan && <PlanSettings plan={plan} onUpgrade={handleUpgrade} />}

          {notifications && (
            <NotificationSettings settings={notifications} onUpdate={handleUpdateNotifications} />
          )}

          <Card>
            <CardHeader>
              <CardTitle>その他</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                <a href="/help" className="block text-primary hover:text-primary-dark transition-colors">
                  ヘルプ・サポート
                </a>
                <a href="/terms" className="block text-primary hover:text-primary-dark transition-colors">
                  利用規約
                </a>
                <a href="/privacy" className="block text-primary hover:text-primary-dark transition-colors">
                  プライバシーポリシー
                </a>
                <button
                  type="button"
                  onClick={() => setShowLogoutDialog(true)}
                  className="text-error hover:text-error-dark transition-colors"
                >
                  ログアウト
                </button>
              </div>
            </CardContent>
          </Card>
        </div>

        {showLogoutDialog && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
            <Card className="w-full max-w-sm animate-fade-in">
              <CardContent className="pt-6">
                <h3 className="text-lg font-semibold mb-4 text-center">ログアウトしますか？</h3>
                <div className="flex gap-3">
                  <Button
                    variant="ghost"
                    onClick={() => setShowLogoutDialog(false)}
                    className="flex-1"
                  >
                    キャンセル
                  </Button>
                  <Button
                    variant="danger"
                    onClick={handleLogout}
                    className="flex-1"
                  >
                    ログアウト
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        )}
      </div>
    </AppLayout>
  );
}
