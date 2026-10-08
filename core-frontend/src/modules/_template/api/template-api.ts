import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/apis/axios-client';
import { tenantKey } from '@/core/tenant/tenant-key';
import { useTenant } from '@/core/tenant/tenant-context';

export interface TemplateData {
  id: string;
  name: string;
}

export async function fetchTemplateData(tenantId: string): Promise<TemplateData[]> {
  try {
    const response = await apiClient.get<TemplateData[]>('/__MODULE_NAME__');
    return response.data;
  } catch (_error) {
    return [];
  }
}

export function useTemplateData() {
  const { tenant } = useTenant();

  return useQuery({
    queryKey: tenantKey(tenant.id, ['__MODULE_NAME__', 'list']),
    queryFn: () => fetchTemplateData(tenant.id),
    staleTime: 60 * 1000,
  });
}
