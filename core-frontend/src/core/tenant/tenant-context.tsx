'use client';

import React, { createContext, useContext, useState, useEffect, useMemo } from 'react';

export interface TenantTheme {
  colorPrimary: string;
  borderRadius?: number;
  colorSuccess?: string;
  colorWarning?: string;
  colorError?: string;
}

export interface TenantConfig {
  id: string;
  name: string;
  theme: TenantTheme;
}

export const DEFAULT_TENANT_CONFIG: TenantConfig = {
  id: 'default',
  name: 'Default Organization',
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

  const value = useMemo(() => ({ tenant, setTenant }), [tenant]);

  return <TenantContext.Provider value={value}>{children}</TenantContext.Provider>;
}

export function useTenant(): TenantContextValue {
  return useContext(TenantContext);
}
