import { cn } from '@/lib/utils';

interface ProgressProps extends React.HTMLAttributes<HTMLDivElement> {
  value: number;
  max?: number;
  size?: 'sm' | 'md' | 'lg';
  variant?: 'primary' | 'secondary' | 'success' | 'warning' | 'danger' | 'gradient';
  showLabel?: boolean;
  label?: string;
  animated?: boolean;
}

export function Progress({
  value,
  max = 100,
  size = 'md',
  variant = 'primary',
  showLabel = false,
  label,
  animated = false,
  className,
  ...props
}: ProgressProps) {
  const percentage = Math.min(Math.max((value / max) * 100, 0), 100);

  const sizes = {
    sm: 'h-1.5',
    md: 'h-2',
    lg: 'h-3',
  };

  const trackColors = {
    primary: 'bg-primary/15',
    secondary: 'bg-secondary/15',
    success: 'bg-success/15',
    warning: 'bg-warning/15',
    danger: 'bg-error/15',
    gradient: 'bg-gray-100',
  };

  const barColors = {
    primary: 'bg-primary',
    secondary: 'bg-secondary',
    success: 'bg-success',
    warning: 'bg-warning',
    danger: 'bg-error',
    gradient: 'bg-gradient-to-r from-primary to-secondary',
  };

  return (
    <div className={cn('w-full', className)} {...props}>
      {(showLabel || label) && (
        <div className="flex justify-between items-center mb-1.5">
          <span className="text-sm text-text-secondary">{label || 'Progress'}</span>
          <span className="text-sm font-medium text-text-primary tabular-nums">
            {Math.round(percentage)}%
          </span>
        </div>
      )}
      {/* biome-ignore lint/a11y/useFocusableInteractive: progressbar is display-only, not interactive */}
      <div
        className={cn('w-full rounded-full overflow-hidden', sizes[size], trackColors[variant])}
        role="progressbar"
        aria-valuenow={value}
        aria-valuemin={0}
        aria-valuemax={max}
      >
        <div
          className={cn(
            'h-full rounded-full transition-all duration-500 ease-out',
            barColors[variant],
            animated && 'animate-pulse'
          )}
          style={{ width: `${percentage}%` }}
        />
      </div>
    </div>
  );
}

interface ProgressCircleProps extends React.HTMLAttributes<HTMLDivElement> {
  value: number;
  max?: number;
  size?: 'sm' | 'md' | 'lg';
  variant?: 'primary' | 'secondary' | 'success' | 'warning' | 'danger';
  strokeWidth?: number;
  showLabel?: boolean;
}

export function ProgressCircle({
  value,
  max = 100,
  size = 'md',
  variant = 'primary',
  strokeWidth = 4,
  showLabel = true,
  className,
  ...props
}: ProgressCircleProps) {
  const percentage = Math.min(Math.max((value / max) * 100, 0), 100);

  const sizes = {
    sm: { size: 48, fontSize: 'text-xs' },
    md: { size: 72, fontSize: 'text-base' },
    lg: { size: 96, fontSize: 'text-xl' },
  };

  const colors = {
    primary: { stroke: '#4A90E2', track: 'rgba(74, 144, 226, 0.15)' },
    secondary: { stroke: '#50C878', track: 'rgba(80, 200, 120, 0.15)' },
    success: { stroke: '#27AE60', track: 'rgba(39, 174, 96, 0.15)' },
    warning: { stroke: '#F39C12', track: 'rgba(243, 156, 18, 0.15)' },
    danger: { stroke: '#E74C3C', track: 'rgba(231, 76, 60, 0.15)' },
  };

  const { size: sizeValue, fontSize } = sizes[size];
  const { stroke, track } = colors[variant];

  const radius = (sizeValue - strokeWidth) / 2;
  const circumference = radius * 2 * Math.PI;
  const offset = circumference - (percentage / 100) * circumference;

  return (
    // biome-ignore lint/a11y/useFocusableInteractive: progressbar is display-only, not interactive
    <div
      className={cn('relative inline-flex items-center justify-center', className)}
      style={{ width: sizeValue, height: sizeValue }}
      role="progressbar"
      aria-valuenow={value}
      aria-valuemin={0}
      aria-valuemax={max}
      {...props}
    >
      <svg aria-hidden="true" className="transform -rotate-90" width={sizeValue} height={sizeValue}>
        <circle
          cx={sizeValue / 2}
          cy={sizeValue / 2}
          r={radius}
          stroke={track}
          strokeWidth={strokeWidth}
          fill="none"
        />
        <circle
          cx={sizeValue / 2}
          cy={sizeValue / 2}
          r={radius}
          stroke={stroke}
          strokeWidth={strokeWidth}
          fill="none"
          strokeLinecap="round"
          strokeDasharray={circumference}
          strokeDashoffset={offset}
          className="transition-all duration-500 ease-out"
        />
      </svg>
      {showLabel && (
        <span className={cn('absolute font-semibold text-text-primary tabular-nums', fontSize)}>
          {Math.round(percentage)}
        </span>
      )}
    </div>
  );
}
