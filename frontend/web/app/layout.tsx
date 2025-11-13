export const metadata = {
  title: 'HaiLanGo - AI Language Learning',
  description: 'AI-Powered Language Learning Platform',
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
