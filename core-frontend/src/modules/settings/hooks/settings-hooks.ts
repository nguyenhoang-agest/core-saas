import { useQuery } from '@tanstack/react-query';
import { tenantKey } from '@/core/tenant/tenant-key';
import { useTenant } from '@/core/tenant/tenant-context';
import { fetchTenantSettings } from '../api/settings-api';

export function useTenantSettings() {
  const { tenant } = useTenant();

  return useQuery({
    queryKey: tenantKey(tenant.id, ['settings', 'general']),
    queryFn: () => fetchTenantSettings(tenant.id),
    staleTime: 60 * 1000,
  });
}
