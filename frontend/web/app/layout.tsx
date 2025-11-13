import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'HaiLanGo - AI言語学習プラットフォーム',
  description: '既存の言語学習本をAI技術で最大限に活用する革新的な学習プラットフォーム',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  );
}
