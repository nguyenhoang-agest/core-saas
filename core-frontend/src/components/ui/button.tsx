import * as React from 'react';
import { cn } from '@/lib/utils';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'default' | 'outline' | 'secondary' | 'ghost' | 'link' | 'destructive';
  size?: 'default' | 'sm' | 'lg' | 'icon';
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'default', size = 'default', ...props }, ref) => {
    const baseStyles =
      'inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-neutral-900 disabled:pointer-events-none disabled:opacity-50';

    const variants = {
      default: 'bg-neutral-900 text-white hover:bg-neutral-800 active:bg-neutral-950 shadow-none',
      outline:
        'border border-neutral-200 bg-white text-neutral-900 hover:bg-neutral-50 hover:border-neutral-300 shadow-none',
      secondary: 'bg-neutral-100 text-neutral-900 hover:bg-neutral-200 shadow-none',
      ghost: 'text-neutral-700 hover:bg-neutral-100 hover:text-neutral-900',
      link: 'text-neutral-900 underline-offset-4 hover:underline',
      destructive: 'bg-red-600 text-white hover:bg-red-700 shadow-none',
    };

    const sizes = {
      default: 'h-9 px-4 py-2',
      sm: 'h-8 rounded-md px-3 text-xs',
      lg: 'h-10 rounded-md px-6 text-sm',
      icon: 'h-9 w-9 p-0',
    };

    return (
      <button
        ref={ref}
        className={cn(baseStyles, variants[variant], sizes[size], className)}
        {...props}
      />
    );
  }
);
Button.displayName = 'Button';

export { Button };
