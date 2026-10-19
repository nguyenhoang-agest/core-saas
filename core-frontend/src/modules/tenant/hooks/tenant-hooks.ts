import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { createTenant, getTenantById, getTenants, type ListTenantsParams } from '@/modules/tenant';

export const tenantQueryKeys = {
  all: ['tenants'] as const,
  list: (params: ListTenantsParams) => ['tenants', 'list', params] as const,
  detail: (id: string) => ['tenants', 'detail', id] as const,
};

export function useTenants(params: ListTenantsParams = {}) {
  return useQuery({
    queryKey: tenantQueryKeys.list(params),
    queryFn: () => getTenants(params),
  });
}

export function useTenant(id: string | null | undefined, enabled = true) {
  return useQuery({
    queryKey: tenantQueryKeys.detail(id || ''),
    queryFn: () => getTenantById(id!),
    enabled: Boolean(id && id.trim().length > 0 && enabled),
    retry: false,
  });
}

export function useCreateTenant() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createTenant,
    onSuccess: (newTenant) => {
      queryClient.setQueryData(tenantQueryKeys.detail(newTenant.id), newTenant);
      queryClient.invalidateQueries({ queryKey: tenantQueryKeys.all });
    },
  });
}
