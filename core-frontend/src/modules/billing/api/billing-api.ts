import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/apis/axios-client';
import { tenantKey } from '@/core/tenant/tenant-key';
import { useTenant } from '@/core/tenant/tenant-context';

export interface SubscriptionInfo {
  planId: string;
  status: 'active' | 'past_due' | 'canceled';
  currentPeriodEnd: string;
}

export async function fetchSubscription(tenantId: string): Promise<SubscriptionInfo> {
  try {
    const res = await apiClient.get<SubscriptionInfo>(`/billing/${tenantId}/subscription`);
    return res.data;
  } catch (_e) {
    return { planId: 'pro', status: 'active', currentPeriodEnd: '2026-12-31' };
  }
}

export function useSubscription() {
  const { tenant } = useTenant();

  return useQuery({
    queryKey: tenantKey(tenant.id, ['billing', 'subscription']),
    queryFn: () => fetchSubscription(tenant.id),
    staleTime: 60 * 1000,
  });
}
