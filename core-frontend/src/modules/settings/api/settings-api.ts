import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/apis/axios-client';
import { tenantKey } from '@/core/tenant/tenant-key';
import { useTenant } from '@/core/tenant/tenant-context';

export interface TenantSettings {
  name: string;
  subdomain: string;
  themeColor: string;
}

export async function fetchTenantSettings(tenantId: string): Promise<TenantSettings> {
  try {
    const res = await apiClient.get<TenantSettings>(`/settings/${tenantId}`);
    return res.data;
  } catch (_e) {
    return { name: 'Default Organization', subdomain: tenantId, themeColor: '#1677ff' };
  }
}

export function useTenantSettings() {
  const { tenant } = useTenant();

  return useQuery({
    queryKey: tenantKey(tenant.id, ['settings', 'general']),
    queryFn: () => fetchTenantSettings(tenant.id),
    staleTime: 60 * 1000,
  });
}
