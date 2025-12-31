import { cn } from '@/lib/utils';

interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: 'default' | 'primary' | 'secondary' | 'success' | 'warning' | 'danger' | 'info';
  size?: 'sm' | 'md' | 'lg';
  dot?: boolean;
  icon?: React.ReactNode;
}

export function Badge({
  className,
  variant = 'default',
  size = 'sm',
  dot,
  icon,
  children,
  ...props
}: BadgeProps) {
  const variants = {
    default: 'bg-background-secondary text-text-secondary',
    primary: 'bg-primary/10 text-primary',
    secondary: 'bg-secondary/10 text-secondary',
    success: 'bg-success/10 text-success',
    warning: 'bg-warning/10 text-warning',
    danger: 'bg-error/10 text-error',
    info: 'bg-info/10 text-info',
  };

  const dotColors = {
    default: 'bg-text-secondary',
    primary: 'bg-primary',
    secondary: 'bg-secondary',
    success: 'bg-success',
    warning: 'bg-warning',
    danger: 'bg-error',
    info: 'bg-info',
  };

  const sizes = {
    sm: 'px-2 py-0.5 text-xs gap-1',
    md: 'px-2.5 py-1 text-sm gap-1.5',
    lg: 'px-3 py-1.5 text-sm gap-2',
  };

  const dotSizes = {
    sm: 'w-1.5 h-1.5',
    md: 'w-2 h-2',
    lg: 'w-2.5 h-2.5',
  };

  const iconSizes = {
    sm: 'w-3 h-3',
    md: 'w-3.5 h-3.5',
    lg: 'w-4 h-4',
  };

  return (
    <span
      className={cn(
        'inline-flex items-center font-medium rounded-full',
        variants[variant],
        sizes[size],
        className
      )}
      {...props}
    >
      {dot && <span className={cn('rounded-full', dotColors[variant], dotSizes[size])} />}
      {icon && <span className={cn('flex-shrink-0', iconSizes[size])}>{icon}</span>}
      {children}
    </span>
  );
}

interface BadgeGroupProps {
  children: React.ReactNode;
  className?: string;
}

export function BadgeGroup({ children, className }: BadgeGroupProps) {
  return <div className={cn('flex flex-wrap gap-2', className)}>{children}</div>;
}
