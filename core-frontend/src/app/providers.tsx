'use client';

import React from 'react';
import { PlatformProvider, type PlatformProviderProps } from '@/core/providers/platform-provider';

/**
 * Replaces hard-coded providers with composable PlatformProvider
 * while preserving backwards compatibility with existing root layouts.
 */
export function Providers({ children, ...props }: PlatformProviderProps) {
  return <PlatformProvider {...props}>{children}</PlatformProvider>;
}
