'use client';

import { Avatar } from '@/components/ui';
import type { Route } from 'next';
import Link from 'next/link';

interface HeaderProps {
  title?: string;
  showBack?: boolean;
  backHref?: string;
  actions?: React.ReactNode;
}

export function Header({ title, showBack, backHref = '/', actions }: HeaderProps) {
  return (
    <header className="lg:hidden sticky top-0 z-40 bg-white border-b border-gray-200">
      <div className="flex items-center justify-between h-14 px-4">
        <div className="flex items-center gap-3">
          {showBack ? (
            <Link
              href={backHref as Route}
              className="p-2 -ml-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
              aria-label="戻る"
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
                  d="M15 19l-7-7 7-7"
                />
              </svg>
            </Link>
          ) : (
            <Link href="/" className="flex items-center gap-2">
              <div className="w-8 h-8 bg-gradient-to-br from-primary to-secondary rounded-lg flex items-center justify-center">
                <span className="text-white font-bold text-sm">H</span>
              </div>
              <span className="font-bold text-gray-900">HaiLanGo</span>
            </Link>
          )}
          {title && <h1 className="font-semibold text-gray-900">{title}</h1>}
        </div>

        <div className="flex items-center gap-2">
          {actions}
          {!showBack && (
            <>
              <button
                type="button"
                className="p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
                aria-label="通知"
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
                    d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
                  />
                </svg>
              </button>
              <Link href="/settings">
                <Avatar size="sm" fallback="U" />
              </Link>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
