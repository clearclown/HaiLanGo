import type { Metadata } from 'next';
import './globals.css';
import { WebSocketProvider } from '@/components/WebSocketProvider';
import { ToastProvider } from '@/components/notifications/ToastContainer';
import { NotificationManager } from '@/components/notifications/NotificationManager';
import { AuthProvider } from '@/components/AuthProvider';

export const metadata: Metadata = {
  title: 'HaiLanGo - AI-Powered Language Learning',
  description: '既存の言語学習本 × AI技術 = 個人に最適化された能動的な学習体験',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <body>
        <ToastProvider>
          <WebSocketProvider>
            <AuthProvider>
              <NotificationManager />
              {children}
            </AuthProvider>
          </WebSocketProvider>
        </ToastProvider>
      </body>
    </html>
  );
}
