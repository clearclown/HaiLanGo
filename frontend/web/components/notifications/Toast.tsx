'use client';

import { useEffect, useState } from 'react';

export type ToastType = 'info' | 'success' | 'warning' | 'error';

export interface ToastProps {
  id: string;
  type: ToastType;
  title: string;
  message?: string;
  duration?: number;
  onClose?: () => void;
}

const Toast = ({ type, title, message, duration = 5000, onClose }: ToastProps) => {
  const [isVisible, setIsVisible] = useState(true);

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsVisible(false);
      setTimeout(() => onClose?.(), 300); // アニメーション完了後にonClose
    }, duration);

    return () => clearTimeout(timer);
  }, [duration, onClose]);

  const handleClose = () => {
    setIsVisible(false);
    setTimeout(() => onClose?.(), 300);
  };

  const typeStyles = {
    info: 'bg-blue-50 border-blue-200 text-blue-900',
    success: 'bg-green-50 border-green-200 text-green-900',
    warning: 'bg-yellow-50 border-yellow-200 text-yellow-900',
    error: 'bg-red-50 border-red-200 text-red-900',
  };

  const iconStyles = {
    info: '🔵',
    success: '✅',
    warning: '⚠️',
    error: '❌',
  };

  return (
    <div
      role="alert"
      aria-live="assertive"
      data-testid="toast-notification"
      data-toast-type={type}
      className={`
        ${typeStyles[type]}
        border rounded-lg p-4 shadow-lg transition-all duration-300 ease-in-out min-w-[300px] max-w-[500px]
        ${isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-2'}
      `}
    >
      <div className="flex items-start justify-between gap-3">
        <div className="flex items-start gap-3 flex-1">
          <span className="text-xl">{iconStyles[type]}</span>
          <div className="flex-1">
            <h4 className="font-semibold text-sm">{title}</h4>
            {message && <p className="text-sm mt-1 opacity-80">{message}</p>}
          </div>
        </div>
        <button
          onClick={handleClose}
          className="text-gray-500 hover:text-gray-700 transition-colors"
          aria-label="Close"
        >
          <svg
            className="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      </div>
    </div>
  );
};

export default Toast;
