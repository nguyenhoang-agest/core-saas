import { useQuery } from '@tanstack/react-query';
import { tenantKey } from '@/core/tenant/tenant-key';
import { useTenant } from '@/core/tenant/tenant-context';
import { fetchIamUsers } from '../api/iam-api';

export function useIamUsers() {
  const { tenant } = useTenant();

  return useQuery({
    queryKey: tenantKey(tenant.id, ['iam', 'users']),
    queryFn: () => fetchIamUsers(tenant.id),
    staleTime: 60 * 1000,
  });
}
