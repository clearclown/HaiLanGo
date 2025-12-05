'use client';

import { cn } from '@/lib/utils';
import { forwardRef, useState } from 'react';

export interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  hint?: string;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
  variant?: 'default' | 'filled';
  inputSize?: 'sm' | 'md' | 'lg';
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  (
    {
      className,
      label,
      error,
      hint,
      leftIcon,
      rightIcon,
      variant = 'default',
      inputSize = 'md',
      id,
      type,
      ...props
    },
    ref
  ) => {
    const inputId = id || label?.toLowerCase().replace(/\s+/g, '-');
    const [showPassword, setShowPassword] = useState(false);

    const isPassword = type === 'password';
    const inputType = isPassword && showPassword ? 'text' : type;

    const sizes = {
      sm: 'h-9 px-3 text-sm',
      md: 'h-12 px-4 text-base',
      lg: 'h-14 px-5 text-lg',
    };

    const variants = {
      default: cn(
        'bg-white border border-border',
        'focus:border-primary focus:ring-2 focus:ring-primary/20',
        error && 'border-error focus:border-error focus:ring-error/20'
      ),
      filled: cn(
        'bg-background-secondary border border-transparent',
        'focus:bg-white focus:border-primary focus:ring-2 focus:ring-primary/20',
        error && 'border-error focus:border-error focus:ring-error/20'
      ),
    };

    const PasswordToggle = () => (
      <button
        type="button"
        onClick={() => setShowPassword(!showPassword)}
        className="text-text-secondary hover:text-text-primary transition-colors focus:outline-none"
        tabIndex={-1}
      >
        {showPassword ? (
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
            />
          </svg>
        ) : (
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
            />
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
            />
          </svg>
        )}
      </button>
    );

    return (
      <div className="w-full">
        {label && (
          <label
            htmlFor={inputId}
            className="block text-sm font-medium text-text-primary mb-2"
          >
            {label}
          </label>
        )}
        <div className="relative">
          {leftIcon && (
            <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-text-secondary">
              <span className="w-5 h-5">{leftIcon}</span>
            </div>
          )}
          <input
            ref={ref}
            id={inputId}
            type={inputType}
            className={cn(
              'block w-full rounded-lg transition-all duration-200',
              'text-text-primary placeholder:text-text-muted',
              'focus:outline-none',
              'disabled:bg-background-secondary disabled:text-text-secondary disabled:cursor-not-allowed',
              sizes[inputSize],
              variants[variant],
              leftIcon && 'pl-12',
              (rightIcon || isPassword) && 'pr-12',
              className
            )}
            {...props}
          />
          {(rightIcon || isPassword) && (
            <div className="absolute inset-y-0 right-0 pr-4 flex items-center">
              {isPassword ? <PasswordToggle /> : rightIcon}
            </div>
          )}
        </div>
        {error && (
          <p className="mt-2 text-sm text-error flex items-center gap-1.5">
            <svg className="w-4 h-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            {error}
          </p>
        )}
        {hint && !error && (
          <p className="mt-2 text-sm text-text-secondary">{hint}</p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';

export interface TextareaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
  error?: string;
  hint?: string;
  variant?: 'default' | 'filled';
}

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ className, label, error, hint, variant = 'default', id, ...props }, ref) => {
    const inputId = id || label?.toLowerCase().replace(/\s+/g, '-');

    const variants = {
      default: cn(
        'bg-white border border-border',
        'focus:border-primary focus:ring-2 focus:ring-primary/20',
        error && 'border-error focus:border-error focus:ring-error/20'
      ),
      filled: cn(
        'bg-background-secondary border border-transparent',
        'focus:bg-white focus:border-primary focus:ring-2 focus:ring-primary/20',
        error && 'border-error focus:border-error focus:ring-error/20'
      ),
    };

    return (
      <div className="w-full">
        {label && (
          <label
            htmlFor={inputId}
            className="block text-sm font-medium text-text-primary mb-2"
          >
            {label}
          </label>
        )}
        <textarea
          ref={ref}
          id={inputId}
          className={cn(
            'block w-full px-4 py-3 rounded-lg transition-all duration-200',
            'text-text-primary placeholder:text-text-muted',
            'focus:outline-none resize-none',
            'disabled:bg-background-secondary disabled:text-text-secondary disabled:cursor-not-allowed',
            variants[variant],
            className
          )}
          {...props}
        />
        {error && (
          <p className="mt-2 text-sm text-error flex items-center gap-1.5">
            <svg className="w-4 h-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            {error}
          </p>
        )}
        {hint && !error && (
          <p className="mt-2 text-sm text-text-secondary">{hint}</p>
        )}
      </div>
    );
  }
);

Textarea.displayName = 'Textarea';
