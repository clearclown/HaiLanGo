import { cn } from '@/lib/utils';
import { forwardRef } from 'react';

interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: 'default' | 'bordered' | 'elevated' | 'interactive' | 'highlight';
  padding?: 'none' | 'sm' | 'md' | 'lg';
}

export const Card = forwardRef<HTMLDivElement, CardProps>(
  ({ className, variant = 'default', padding = 'md', children, ...props }, ref) => {
    const variants = {
      default: cn('bg-white rounded-xl', 'shadow-[0_2px_8px_rgba(0,0,0,0.06)]'),
      bordered: cn(
        'bg-white rounded-xl',
        'border border-border',
        'hover:border-border-dark transition-colors'
      ),
      elevated: cn('bg-white rounded-xl', 'shadow-[0_4px_16px_rgba(0,0,0,0.1)]'),
      interactive: cn(
        'bg-white rounded-xl',
        'shadow-[0_2px_8px_rgba(0,0,0,0.06)]',
        'hover:shadow-[0_4px_16px_rgba(0,0,0,0.1)]',
        'hover:translate-y-[-2px]',
        'transition-all duration-200 cursor-pointer'
      ),
      highlight: cn(
        'bg-gradient-to-br from-primary/5 to-secondary/5',
        'rounded-xl border border-primary/20'
      ),
    };

    const paddings = {
      none: '',
      sm: 'p-4',
      md: 'p-6',
      lg: 'p-8',
    };

    return (
      <div ref={ref} className={cn(variants[variant], paddings[padding], className)} {...props}>
        {children}
      </div>
    );
  }
);

Card.displayName = 'Card';

interface CardHeaderProps extends React.HTMLAttributes<HTMLDivElement> {
  action?: React.ReactNode;
}

export function CardHeader({ className, action, children, ...props }: CardHeaderProps) {
  return (
    <div className={cn('flex items-start justify-between mb-4', className)} {...props}>
      <div className="flex-1">{children}</div>
      {action && <div className="ml-4 flex-shrink-0">{action}</div>}
    </div>
  );
}

interface CardTitleProps extends React.HTMLAttributes<HTMLHeadingElement> {
  as?: 'h1' | 'h2' | 'h3' | 'h4';
}

export function CardTitle({ className, as: Tag = 'h3', children, ...props }: CardTitleProps) {
  return (
    <Tag
      className={cn(
        'font-semibold text-text-primary',
        Tag === 'h1' && 'text-2xl',
        Tag === 'h2' && 'text-xl',
        Tag === 'h3' && 'text-lg',
        Tag === 'h4' && 'text-base',
        className
      )}
      {...props}
    >
      {children}
    </Tag>
  );
}

interface CardDescriptionProps extends React.HTMLAttributes<HTMLParagraphElement> {}

export function CardDescription({ className, children, ...props }: CardDescriptionProps) {
  return (
    <p className={cn('text-sm text-text-secondary mt-1 leading-relaxed', className)} {...props}>
      {children}
    </p>
  );
}

interface CardContentProps extends React.HTMLAttributes<HTMLDivElement> {}

export function CardContent({ className, children, ...props }: CardContentProps) {
  return (
    <div className={cn('', className)} {...props}>
      {children}
    </div>
  );
}

interface CardFooterProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: 'default' | 'bordered' | 'actions';
}

export function CardFooter({
  className,
  variant = 'default',
  children,
  ...props
}: CardFooterProps) {
  const variants = {
    default: 'mt-4',
    bordered: 'mt-4 pt-4 border-t border-border-light',
    actions: 'mt-4 pt-4 border-t border-border-light flex items-center gap-3',
  };

  return (
    <div className={cn(variants[variant], className)} {...props}>
      {children}
    </div>
  );
}

interface CardImageProps extends React.ImgHTMLAttributes<HTMLImageElement> {
  aspectRatio?: 'auto' | 'video' | 'square';
  alt: string;
}

export function CardImage({ className, aspectRatio = 'auto', alt, ...props }: CardImageProps) {
  const aspectRatios = {
    auto: '',
    video: 'aspect-video',
    square: 'aspect-square',
  };

  return (
    <div className={cn('overflow-hidden rounded-t-xl -mx-6 -mt-6 mb-4', aspectRatios[aspectRatio])}>
      <img className={cn('w-full h-full object-cover', className)} {...props} alt={alt} />
    </div>
  );
}
