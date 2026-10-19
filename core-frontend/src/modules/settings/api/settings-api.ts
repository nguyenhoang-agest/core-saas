import { apiClient } from '@/apis/axios-client';

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
