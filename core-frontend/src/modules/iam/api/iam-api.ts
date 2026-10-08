import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/apis/axios-client';
import { tenantKey } from '@/core/tenant/tenant-key';
import { useTenant } from '@/core/tenant/tenant-context';

export interface IamUser {
  id: string;
  email: string;
  name: string;
  role: string;
}

export async function fetchIamUsers(tenantId: string): Promise<IamUser[]> {
  try {
    const res = await apiClient.get<IamUser[]>('/iam/users');
    return res.data;
  } catch (_e) {
    return [];
  }
}

export function useIamUsers() {
  const { tenant } = useTenant();

  return useQuery({
    queryKey: tenantKey(tenant.id, ['iam', 'users']),
    queryFn: () => fetchIamUsers(tenant.id),
    staleTime: 60 * 1000,
  });
}
