'use client';

import { BottomNav } from './BottomNav';
import { Header } from './Header';
import { Sidebar } from './Sidebar';

interface AppLayoutProps {
  children: React.ReactNode;
  title?: string;
  showBack?: boolean;
  backHref?: string;
  headerActions?: React.ReactNode;
}

export function AppLayout({ children, title, showBack, backHref, headerActions }: AppLayoutProps) {
  return (
    <div className="min-h-screen bg-gray-50">
      <div className="flex">
        <Sidebar />
        <main className="flex-1 min-h-screen pb-20 lg:pb-0">
          <Header title={title} showBack={showBack} backHref={backHref} actions={headerActions} />
          {children}
        </main>
      </div>
      <BottomNav />
    </div>
  );
}
