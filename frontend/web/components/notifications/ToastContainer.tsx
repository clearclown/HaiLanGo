'use client';

import { createContext, useContext, useState, useCallback, useEffect } from 'react';
import type { ReactNode } from 'react';
import Toast from './Toast';
import type { ToastProps } from './Toast';

interface ToastContextType {
  addToast: (toast: Omit<ToastProps, 'id' | 'onClose'>) => void;
  removeToast: (id: string) => void;
  showInfo: (title: string, message?: string) => void;
  showSuccess: (title: string, message?: string) => void;
  showWarning: (title: string, message?: string) => void;
  showError: (title: string, message?: string) => void;
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export const useToast = () => {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within ToastProvider');
  }
  return context;
};

interface ToastProviderProps {
  children: ReactNode;
}

export const ToastProvider = ({ children }: ToastProviderProps) => {
  const [toasts, setToasts] = useState<ToastProps[]>([]);

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((toast) => toast.id !== id));
  }, []);

  const addToast = useCallback(
    (toast: Omit<ToastProps, 'id' | 'onClose'>) => {
      const id = `toast-${Date.now()}-${Math.random()}`;
      const newToast: ToastProps = {
        ...toast,
        id,
        onClose: () => removeToast(id),
      };
      setToasts((prev) => [...prev, newToast]);
    },
    [removeToast]
  );

  const showInfo = useCallback(
    (title: string, message?: string) => {
      addToast({ type: 'info', title, message });
    },
    [addToast]
  );

  const showSuccess = useCallback(
    (title: string, message?: string) => {
      addToast({ type: 'success', title, message });
    },
    [addToast]
  );

  const showWarning = useCallback(
    (title: string, message?: string) => {
      addToast({ type: 'warning', title, message });
    },
    [addToast]
  );

  const showError = useCallback(
    (title: string, message?: string) => {
      addToast({ type: 'error', title, message });
    },
    [addToast]
  );

  // Expose toast methods to window for E2E testing
  useEffect(() => {
    if (typeof window !== 'undefined' && process.env.NODE_ENV !== 'production') {
      (window as any).__TEST_TOAST__ = {
        showInfo,
        showSuccess,
        showWarning,
        showError,
        addToast,
        removeToast,
      };
    }
    return () => {
      if (typeof window !== 'undefined') {
        delete (window as any).__TEST_TOAST__;
      }
    };
  }, [showInfo, showSuccess, showWarning, showError, addToast, removeToast]);

  return (
    <ToastContext.Provider
      value={{
        addToast,
        removeToast,
        showInfo,
        showSuccess,
        showWarning,
        showError,
      }}
    >
      {children}
      {/* Toast Container - 画面右上に固定表示 */}
      <div className="fixed top-4 right-4 z-50 space-y-2">
        {toasts.map((toast) => (
          <Toast key={toast.id} {...toast} />
        ))}
      </div>
    </ToastContext.Provider>
  );
};
