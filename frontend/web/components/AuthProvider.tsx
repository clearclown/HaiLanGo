'use client';

import { usePathname, useRouter } from 'next/navigation';
import { createContext, useCallback, useContext, useEffect, useState } from 'react';

// Pages that don't require authentication
const PUBLIC_PATHS = ['/login', '/register', '/forgot-password'];

// Auth context type
interface AuthContextType {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: User | null;
  login: (token: string, user?: User) => void;
  logout: () => void;
  refreshAuth: () => void;
}

interface User {
  id: string;
  email: string;
  name?: string;
}

// Create context with default values
const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  isLoading: true,
  user: null,
  login: () => {},
  logout: () => {},
  refreshAuth: () => {},
});

// Custom hook to use auth context
export function useAuth(): AuthContextType {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const [isLoading, setIsLoading] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState<User | null>(null);

  const checkAuth = useCallback(() => {
    // Check for token in localStorage
    const token = localStorage.getItem('access_token');
    const storedUser = localStorage.getItem('user');

    if (token) {
      setIsAuthenticated(true);
      if (storedUser) {
        try {
          setUser(JSON.parse(storedUser));
        } catch {
          setUser(null);
        }
      }
    } else {
      setIsAuthenticated(false);
      setUser(null);
    }

    setIsLoading(false);
  }, []);

  const login = useCallback((token: string, userData?: User) => {
    localStorage.setItem('access_token', token);
    if (userData) {
      localStorage.setItem('user', JSON.stringify(userData));
      setUser(userData);
    }
    setIsAuthenticated(true);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('user');
    setIsAuthenticated(false);
    setUser(null);
    router.push('/login');
  }, [router]);

  const refreshAuth = useCallback(() => {
    checkAuth();
  }, [checkAuth]);

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  useEffect(() => {
    // Check if we're on a public path
    const isPublicPath = PUBLIC_PATHS.some((path) => pathname?.startsWith(path));

    if (!isLoading) {
      if (!isAuthenticated && !isPublicPath) {
        // No token and not on public path - redirect to login
        router.push('/login');
      } else if (isAuthenticated && pathname === '/login') {
        // Already logged in but on login page - redirect to home
        router.push('/');
      }
    }
  }, [pathname, router, isAuthenticated, isLoading]);

  // Show loading while checking auth on protected pages
  if (isLoading && !PUBLIC_PATHS.some((path) => pathname?.startsWith(path))) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
      </div>
    );
  }

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        isLoading,
        user,
        login,
        logout,
        refreshAuth,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}
