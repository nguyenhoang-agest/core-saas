'use client';

import React, { createContext, useContext, useState, useEffect, useMemo } from 'react';

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
  const [tenant, setTenant] = useState<TenantConfig>(initialTenant);

  // Read persisted tenant from cookie on client mount
  useEffect(() => {
    if (typeof document !== 'undefined') {
      const match = document.cookie.match(/(?:^|; )x-tenant-id=([^;]*)/);
      if (match && match[1]) {
        const cookieTenantId = decodeURIComponent(match[1]);
        if (cookieTenantId) {
          let savedName = cookieTenantId;
          try {
            savedName = localStorage.getItem('core_saas_tenant_name') || cookieTenantId;
          } catch {
            // ignore
          }
          setTenant((prev) => {
            if (prev.id !== cookieTenantId) {
              return {
                ...prev,
                id: cookieTenantId,
                name: savedName,
              };
            }
            return prev;
          });
        }
      }
    }
  }, []);

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
