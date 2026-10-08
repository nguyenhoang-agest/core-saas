import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/apis/axios-client';
import { tenantKey } from '@/core/tenant/tenant-key';
import { useTenant } from '@/core/tenant/tenant-context';

export interface ExampleItem {
  id: string;
  title: string;
  status: 'active' | 'pending' | 'completed';
  createdAt: string;
}

export const MOCK_EXAMPLE_ITEMS: ExampleItem[] = [
  { id: '1', title: 'Setup Authentication & RBAC', status: 'completed', createdAt: '2026-10-01' },
  { id: '2', title: 'Configure Multi-tenant Subdomain Resolution', status: 'completed', createdAt: '2026-10-02' },
  { id: '3', title: 'Add Custom Domain Ingestion', status: 'pending', createdAt: '2026-10-03' },
  { id: '4', title: 'Audit Tenant Resource Isolation', status: 'active', createdAt: '2026-10-04' },
];

/**
 * Fetches example items scoped to the active tenant.
 */
export async function fetchExampleItems(tenantId: string): Promise<ExampleItem[]> {
  try {
    const response = await apiClient.get<ExampleItem[]>('/example/items');
    return response.data;
  } catch (_error) {
    // In starter kit mode without a running backend, gracefully fallback to mock data
    return MOCK_EXAMPLE_ITEMS;
  }
}

export function useExampleItems() {
  const { tenant } = useTenant();

  return useQuery({
    queryKey: tenantKey(tenant.id, ['example', 'items']),
    queryFn: () => fetchExampleItems(tenant.id),
    staleTime: 30 * 1000,
  });
}
