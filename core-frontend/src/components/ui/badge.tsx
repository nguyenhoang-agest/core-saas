import * as React from 'react';
import { cn } from '@/lib/utils';

export interface BadgeProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: 'outline' | 'default' | 'secondary' | 'success';
}

function Badge({ className, variant = 'outline', ...props }: BadgeProps) {
  const variants = {
    outline: 'border border-neutral-200 bg-neutral-50 text-neutral-800',
    default: 'bg-neutral-900 text-white border-transparent',
    secondary: 'bg-neutral-100 text-neutral-800 border-transparent',
    success: 'border border-emerald-200 bg-emerald-50 text-emerald-800',
  };

  return (
    <div
      className={cn(
        'inline-flex items-center rounded-md px-2 py-0.5 text-xs font-medium tracking-wide transition-colors',
        variants[variant],
        className
      )}
      {...props}
    />
  );
}

export { Badge };
