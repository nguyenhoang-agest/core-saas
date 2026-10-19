import * as React from 'react';
import { cn } from '@/lib/utils';

export interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  error?: boolean;
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type, error, ...props }, ref) => {
    return (
      <input
        type={type}
        className={cn(
          'flex h-9 w-full rounded-md border border-neutral-200 bg-white px-3 py-1 text-sm text-neutral-900 transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-neutral-400 focus-visible:outline-none focus-visible:border-neutral-900 focus-visible:ring-1 focus-visible:ring-neutral-900 disabled:cursor-not-allowed disabled:opacity-50 disabled:bg-neutral-50',
          error && 'border-red-500 focus-visible:border-red-600 focus-visible:ring-red-600',
          className
        )}
        ref={ref}
        {...props}
      />
    );
  }
);
Input.displayName = 'Input';

export { Input };
