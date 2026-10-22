import { useQuery } from '@tanstack/react-query';
import { tenantKey } from '@/core/tenant/tenant-key';
import { useTenant } from '@/core/tenant/tenant-context';
import { fetchExampleItems } from '../api/example-api';

export function useExampleItems() {
  const { tenant } = useTenant();

  return useQuery({
    queryKey: tenantKey(tenant.id, ['example', 'items']),
    queryFn: () => fetchExampleItems(tenant.id),
    staleTime: 30 * 1000,
  });
}
