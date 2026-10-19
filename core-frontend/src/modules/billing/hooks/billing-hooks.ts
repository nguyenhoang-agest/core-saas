import { useQuery } from '@tanstack/react-query';
import { tenantKey } from '@/core/tenant/tenant-key';
import { useTenant } from '@/core/tenant/tenant-context';
import { fetchSubscription } from '../api/billing-api';

export function useSubscription() {
  const { tenant } = useTenant();

  return useQuery({
    queryKey: tenantKey(tenant.id, ['billing', 'subscription']),
    queryFn: () => fetchSubscription(tenant.id),
    staleTime: 60 * 1000,
  });
}
