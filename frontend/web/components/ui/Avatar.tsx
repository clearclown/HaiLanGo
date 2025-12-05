'use client';

import { cn } from '@/lib/utils';
import { useState } from 'react';

interface AvatarProps extends React.HTMLAttributes<HTMLDivElement> {
  src?: string;
  alt?: string;
  fallback?: string;
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl';
  variant?: 'circle' | 'rounded';
  status?: 'online' | 'offline' | 'busy' | 'away';
  ring?: boolean;
}

export function Avatar({
  src,
  alt = '',
  fallback,
  size = 'md',
  variant = 'circle',
  status,
  ring,
  className,
  ...props
}: AvatarProps) {
  const [imageError, setImageError] = useState(false);

  const sizes = {
    xs: 'w-6 h-6 text-[10px]',
    sm: 'w-8 h-8 text-xs',
    md: 'w-10 h-10 text-sm',
    lg: 'w-12 h-12 text-base',
    xl: 'w-16 h-16 text-lg',
    '2xl': 'w-20 h-20 text-xl',
  };

  const statusSizes = {
    xs: 'w-1.5 h-1.5 border',
    sm: 'w-2 h-2 border',
    md: 'w-2.5 h-2.5 border-2',
    lg: 'w-3 h-3 border-2',
    xl: 'w-3.5 h-3.5 border-2',
    '2xl': 'w-4 h-4 border-2',
  };

  const statusColors = {
    online: 'bg-success',
    offline: 'bg-text-secondary',
    busy: 'bg-error',
    away: 'bg-warning',
  };

  const variants = {
    circle: 'rounded-full',
    rounded: 'rounded-lg',
  };

  const initials =
    fallback ||
    alt
      ?.split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2) ||
    '?';

  const showFallback = !src || imageError;

  return (
    <div
      className={cn(
        'relative inline-flex items-center justify-center',
        'bg-gradient-to-br from-primary/20 to-secondary/20',
        'text-text-primary font-medium overflow-hidden',
        sizes[size],
        variants[variant],
        ring && 'ring-2 ring-white ring-offset-2 ring-offset-background-secondary',
        className
      )}
      {...props}
    >
      {showFallback ? (
        <span className="select-none">{initials}</span>
      ) : (
        <img
          src={src}
          alt={alt}
          className="w-full h-full object-cover"
          onError={() => setImageError(true)}
        />
      )}
      {status && (
        <span
          className={cn(
            'absolute bottom-0 right-0 border-white',
            statusSizes[size],
            statusColors[status],
            variant === 'circle' ? 'rounded-full' : 'rounded-sm'
          )}
        />
      )}
    </div>
  );
}

interface AvatarGroupProps {
  children: React.ReactNode;
  max?: number;
  size?: AvatarProps['size'];
  className?: string;
}

export function AvatarGroup({
  children,
  max = 5,
  size = 'md',
  className,
}: AvatarGroupProps) {
  const childArray = Array.isArray(children) ? children : [children];
  const visibleAvatars = childArray.slice(0, max);
  const remainingCount = childArray.length - max;

  const overlapSizes = {
    xs: '-ml-2',
    sm: '-ml-2.5',
    md: '-ml-3',
    lg: '-ml-4',
    xl: '-ml-5',
    '2xl': '-ml-6',
  };

  return (
    <div className={cn('flex items-center', className)}>
      {visibleAvatars.map((child, index) => (
        <div
          key={index}
          className={cn(index > 0 && overlapSizes[size], 'ring-2 ring-white rounded-full')}
        >
          {child}
        </div>
      ))}
      {remainingCount > 0 && (
        <Avatar
          size={size}
          fallback={`+${remainingCount}`}
          className={cn(overlapSizes[size], 'ring-2 ring-white')}
        />
      )}
    </div>
  );
}
