'use client';

import React, { useState, useEffect } from 'react';
import { SessionProvider } from 'next-auth/react';
import type { Session } from 'next-auth';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AntdRegistry } from '@ant-design/nextjs-registry';
import { ConfigProvider, theme } from 'antd';
import { TenantProvider, useTenant, type TenantConfig } from '../tenant/tenant-context';
import { OverrideProvider, type ComponentMap } from '../overrides/override-context';
import { moduleRegistry } from '../modules/registry';
import type { ModuleDefinition } from '../modules/types';

export interface PlatformProviderProps {
  children: React.ReactNode;
  session?: Session | null;
  initialTenant?: TenantConfig;
  overrides?: ComponentMap;
  modules?: ModuleDefinition[];
}

function InnerAntdThemeConfig({ children }: { children: React.ReactNode }) {
  const { tenant } = useTenant();

  return (
    <ConfigProvider
      theme={{
        algorithm: theme.defaultAlgorithm,
        token: {
          colorPrimary: tenant.theme?.colorPrimary || '#1677ff',
          borderRadius: tenant.theme?.borderRadius ?? 6,
        },
      }}
    >
      {children}
    </ConfigProvider>
  );
}

export function PlatformProvider({
  children,
  session,
  initialTenant,
  overrides,
  modules,
}: PlatformProviderProps) {
  // Synchronously register enabled modules if provided
  if (modules && modules.length > 0) {
    for (const mod of modules) {
      if (!moduleRegistry.get(mod.id)) {
        moduleRegistry.register(mod);
      }
    }
  }

  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 60 * 1000,
            refetchOnWindowFocus: false,
          },
        },
      })
  );

  return (
    <SessionProvider session={session}>
      <QueryClientProvider client={queryClient}>
        <TenantProvider initialTenant={initialTenant}>
          <AntdRegistry>
            <InnerAntdThemeConfig>
              <OverrideProvider overrides={overrides}>{children}</OverrideProvider>
            </InnerAntdThemeConfig>
          </AntdRegistry>
        </TenantProvider>
      </QueryClientProvider>
    </SessionProvider>
  );
}
