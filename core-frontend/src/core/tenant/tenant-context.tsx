'use client';

import React, { createContext, useContext, useState, useEffect, useMemo } from 'react';
import { useSession } from 'next-auth/react';

function useSafeSession() {
  try {
    return useSession();
  } catch {
    return { data: null, status: 'unauthenticated' as const };
  }
}

export interface TenantTheme {
  colorPrimary: string;
  borderRadius?: number;
  colorSuccess?: string;
  colorWarning?: string;
  colorError?: string;
}

export type TenantIsolateMode = 'shared' | 'schema' | 'dedicated';

export interface TenantConfig {
  id: string;
  name: string;
  theme: TenantTheme;
  isolateMode?: TenantIsolateMode;
}

export const DEFAULT_TENANT_CONFIG: TenantConfig = {
  id: 'default',
  name: 'Default Organization',
  isolateMode: 'shared',
  theme: {
    colorPrimary: '#1677ff',
    borderRadius: 6,
  },
};


interface TenantContextValue {
  tenant: TenantConfig;
  setTenant: (tenant: TenantConfig) => void;
}

const TenantContext = createContext<TenantContextValue>({
  tenant: DEFAULT_TENANT_CONFIG,
  setTenant: () => {},
});

export function TenantProvider({
  children,
  initialTenant = DEFAULT_TENANT_CONFIG,
}: {
  children: React.ReactNode;
  initialTenant?: TenantConfig;
}) {
  const sessionResult = useSafeSession();
  const session = sessionResult?.data;
  const status = sessionResult?.status;

  const [tenant, setTenant] = useState<TenantConfig>(initialTenant);

  // Sync active tenant whenever authenticated session changes
  useEffect(() => {
    if (status === 'unauthenticated') {
      setTenant(DEFAULT_TENANT_CONFIG);
      return;
    }

    if (status !== 'authenticated' || !session?.user) return;

    const userTenantId = (session.user as any).tenantId;
    const userTenantName = (session.user as any).tenantName || userTenantId;
    const isHost = !!(session.user as any).isHost;

    if (!isHost && userTenantId) {
      // Non-host users (such as tenant owners) are strictly bound to their own workspace.
      // Automatically overwrite any leftover admin tenant cookie or localStorage.
      setTenant((prev) => {
        if (prev.id !== userTenantId || prev.name !== userTenantName) {
          return {
            ...prev,
            id: userTenantId,
            name: userTenantName,
          };
        }
        return prev;
      });

      if (typeof document !== 'undefined') {
        document.cookie = `x-tenant-id=${encodeURIComponent(userTenantId)}; path=/; max-age=31536000; SameSite=Lax`;
        try {
          localStorage.setItem('core_saas_tenant_id', userTenantId);
          localStorage.setItem('core_saas_tenant_name', userTenantName);
        } catch {
          // ignore
        }
      }
    } else if (isHost) {
      // Host administrator: read persisted tenant or default to session tenant
      if (typeof document !== 'undefined') {
        const match = document.cookie.match(/(?:^|; )x-tenant-id=([^;]*)/);
        const cookieTenantId = match && match[1] ? decodeURIComponent(match[1]) : '';
        const effectiveId = cookieTenantId || userTenantId || 'default';
        let savedName = effectiveId;
        try {
          savedName = localStorage.getItem('core_saas_tenant_name') || userTenantName || effectiveId;
        } catch {
          // ignore
        }
        setTenant((prev) => {
          if (prev.id !== effectiveId || prev.name !== savedName) {
            return {
              ...prev,
              id: effectiveId,
              name: savedName,
            };
          }
          return prev;
        });
      }
    }
  }, [status, session]);

  const handleSetTenant = (nextTenant: TenantConfig) => {
    setTenant(nextTenant);
    if (typeof document !== 'undefined') {
      document.cookie = `x-tenant-id=${encodeURIComponent(nextTenant.id)}; path=/; max-age=31536000; SameSite=Lax`;
      try {
        localStorage.setItem('core_saas_tenant_id', nextTenant.id);
        localStorage.setItem('core_saas_tenant_name', nextTenant.name);
      } catch {
        // ignore
      }
    }
  };

  // Sync theme to CSS custom properties
  useEffect(() => {
    if (typeof document !== 'undefined') {
      const root = document.documentElement;
      if (tenant.theme?.colorPrimary) {
        root.style.setProperty('--primary-color', tenant.theme.colorPrimary);
      }
      if (tenant.theme?.borderRadius !== undefined) {
        root.style.setProperty('--border-radius', `${tenant.theme.borderRadius}px`);
      }
    }
  }, [tenant]);

  const value = useMemo(() => ({ tenant, setTenant: handleSetTenant }), [tenant]);

  return <TenantContext.Provider value={value}>{children}</TenantContext.Provider>;
}

export function useTenant(): TenantContextValue {
  return useContext(TenantContext);
}
