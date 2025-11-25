'use client';

import { useEffect, useState } from 'react';
import { usePathname, useRouter } from 'next/navigation';

// Pages that don't require authentication
const PUBLIC_PATHS = ['/login', '/register', '/forgot-password'];

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const [isChecking, setIsChecking] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    const checkAuth = () => {
      // Check if we're on a public path
      const isPublicPath = PUBLIC_PATHS.some(path => pathname?.startsWith(path));

      // Check for token in localStorage
      const token = localStorage.getItem('access_token');

      if (!token && !isPublicPath) {
        // No token and not on public path - redirect to login
        router.push('/login');
        return;
      }

      if (token && pathname === '/login') {
        // Already logged in but on login page - redirect to home
        router.push('/');
        return;
      }

      setIsAuthenticated(!!token);
      setIsChecking(false);
    };

    checkAuth();
  }, [pathname, router]);

  // Show loading while checking auth on protected pages
  if (isChecking && !PUBLIC_PATHS.some(path => pathname?.startsWith(path))) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return <>{children}</>;
}
