'use client';

import React from 'react';
import { useSession } from 'next-auth/react';
import type { CustomSession } from './auth-options';

export interface CanProps {
  /**
   * The permission code required to display the child content.
   * Wildcard '*' grants access to all features.
   */
  permission: string;
  /**
   * Fallback content to display when permission check fails.
   */
  fallback?: React.ReactNode;
  /**
   * Protected children nodes.
   */
  children: React.ReactNode;
}

/**
 * `<Can>` is a client-side utility strictly for showing/hiding UI elements
 * based on session permissions.
 *
 * IMPORTANT SECURITY NOTICE:
 * Client-side permission checks are purely a user experience convenience.
 * Attackers can inspect client code, manipulate memory, or bypass UI elements.
 * Real, authoritative authorization MUST ALWAYS be enforced on the backend / API services.
 */
export function Can({ permission, fallback = null, children }: CanProps) {
  const { data: session } = useSession();
  const customSession = session as CustomSession | null;
  const userPermissions = customSession?.user?.permissions || [];

  const hasPermission =
    userPermissions.includes('*') || userPermissions.includes(permission);

  if (!hasPermission) {
    return <>{fallback}</>;
  }

  return <>{children}</>;
}
