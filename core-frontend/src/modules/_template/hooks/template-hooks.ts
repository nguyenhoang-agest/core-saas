import { useQuery } from '@tanstack/react-query';
import { tenantKey } from '@/core/tenant/tenant-key';
import { useTenant } from '@/core/tenant/tenant-context';
import { fetchTemplateData } from '../api/template-api';

export function useTemplateData() {
  const { tenant } = useTenant();

  return useQuery({
    queryKey: tenantKey(tenant.id, ['__MODULE_NAME__', 'list']),
    queryFn: () => fetchTemplateData(tenant.id),
    staleTime: 60 * 1000,
  });
}
