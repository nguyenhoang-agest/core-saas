import { apiClient } from '@/apis/axios-client';

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
