'use client';

import { cn } from '@/lib/utils';
import { forwardRef } from 'react';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger' | 'success';
  size?: 'sm' | 'md' | 'lg';
  isLoading?: boolean;
  fullWidth?: boolean;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      className,
      variant = 'primary',
      size = 'md',
      isLoading,
      disabled,
      fullWidth,
      leftIcon,
      rightIcon,
      children,
      ...props
    },
    ref
  ) => {
    const baseStyles = cn(
      'inline-flex items-center justify-center font-medium',
      'rounded-lg transition-all duration-200',
      'focus:outline-none focus:ring-2 focus:ring-offset-2',
      'disabled:opacity-50 disabled:cursor-not-allowed',
      'active:scale-[0.98]'
    );

    const variants = {
      primary: cn(
        'bg-primary text-white',
        'hover:bg-primary-dark',
        'focus:ring-primary/50',
        'shadow-[0_2px_4px_rgba(74,144,226,0.3)]',
        'hover:shadow-[0_4px_8px_rgba(74,144,226,0.4)]'
      ),
      secondary: cn(
        'bg-white text-primary',
        'border border-primary',
        'hover:bg-primary/5',
        'focus:ring-primary/50'
      ),
      outline: cn(
        'bg-transparent text-text-primary',
        'border border-border',
        'hover:bg-background-secondary hover:border-border-dark',
        'focus:ring-primary/30'
      ),
      ghost: cn(
        'bg-transparent text-text-secondary',
        'hover:bg-background-secondary hover:text-text-primary',
        'focus:ring-gray-300'
      ),
      danger: cn(
        'bg-error text-white',
        'hover:bg-error-dark',
        'focus:ring-error/50',
        'shadow-[0_2px_4px_rgba(231,76,60,0.3)]'
      ),
      success: cn(
        'bg-success text-white',
        'hover:bg-success-dark',
        'focus:ring-success/50',
        'shadow-[0_2px_4px_rgba(39,174,96,0.3)]'
      ),
    };

    const sizes = {
      sm: 'h-9 px-3 text-sm gap-1.5',
      md: 'h-12 px-5 text-base gap-2',
      lg: 'h-14 px-7 text-lg gap-2.5',
    };

    const iconSizes = {
      sm: 'w-4 h-4',
      md: 'w-5 h-5',
      lg: 'w-6 h-6',
    };

    const LoadingSpinner = () => (
      <svg
        className={cn('animate-spin', iconSizes[size])}
        fill="none"
        viewBox="0 0 24 24"
      >
        <circle
          className="opacity-25"
          cx="12"
          cy="12"
          r="10"
          stroke="currentColor"
          strokeWidth="4"
        />
        <path
          className="opacity-75"
          fill="currentColor"
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
        />
      </svg>
    );

    return (
      <button
        ref={ref}
        className={cn(
          baseStyles,
          variants[variant],
          sizes[size],
          fullWidth && 'w-full',
          className
        )}
        disabled={disabled || isLoading}
        {...props}
      >
        {isLoading ? (
          <>
            <LoadingSpinner />
            <span>Loading...</span>
          </>
        ) : (
          <>
            {leftIcon && (
              <span className={cn('flex-shrink-0', iconSizes[size])}>
                {leftIcon}
              </span>
            )}
            {children}
            {rightIcon && (
              <span className={cn('flex-shrink-0', iconSizes[size])}>
                {rightIcon}
              </span>
            )}
          </>
        )}
      </button>
    );
  }
);

Button.displayName = 'Button';
