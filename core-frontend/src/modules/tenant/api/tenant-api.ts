import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axios from 'axios';
import { apiClient } from '@/apis/axios-client';
import { getClientLocale, type Locale } from '@/core/i18n/i18n-context';
import enErrors from '@/locales/en/errors.json';
import viErrors from '@/locales/vi/errors.json';
import type { PaginatedData, PaginationParams } from '@/types';
import type {
  Tenant,
  CreateTenantInput,
  ApiResponse,
  BackendResponseError,
} from '../types/tenant';

const errorDictionaries: Record<Locale, Record<string, string>> = {
  en: enErrors,
  vi: viErrors,
};

export function parseBackendError(
  error: unknown,
  overrideLocale?: Locale
): {
  message: string;
  code?: string;
  fieldErrors?: Record<string, string>;
} {
  const currentLocale = overrideLocale || getClientLocale();
  const dict = errorDictionaries[currentLocale] || errorDictionaries.en;

  if (axios.isAxiosError(error) && error.response?.data) {
    const data = error.response.data as BackendResponseError & { message?: string };
    const code = data.message_code;

    // Two-tier error resolution: local dictionary -> backend localized message -> fallback
    const message =
      (code && dict[code]) ||
      data.message ||
      (code ? `Operation failed with error code: ${code}` : 'An unexpected error occurred.');

    const fieldErrors: Record<string, string> = {};
    if (data.error_detail && Array.isArray(data.error_detail)) {
      for (const item of data.error_detail) {
        if (item.field) {
          const itemMsg =
            (item.message_code && dict[item.message_code]) ||
            (item as { message?: string }).message ||
            'Invalid value';
          fieldErrors[item.field] = itemMsg;
        }
      }
    }

    return { message, code, fieldErrors };
  }

  if (error instanceof Error) {
    return { message: error.message };
  }

  return {
    message:
      currentLocale === 'vi'
        ? 'Đã xảy ra lỗi mạng hoặc lỗi không xác định.'
        : 'An unexpected network error occurred.',
  };
}

export async function getTenantById(id: string): Promise<Tenant> {
  const response = await apiClient.get<ApiResponse<Tenant>>(`/tenants/${encodeURIComponent(id)}`);
  if (response.data.data) {
    return response.data.data;
  }
  throw new Error('Tenant data missing in response');
}

export async function createTenant(input: CreateTenantInput): Promise<Tenant> {
  const response = await apiClient.post<ApiResponse<Tenant>>('/tenants', input);
  if (response.data.data) {
    return response.data.data;
  }
  throw new Error('Created tenant data missing in response');
}

export interface ListTenantsParams extends PaginationParams {
  keyword?: string;
  status?: string;
  isolateMode?: string;
  isolate_mode?: string;
}

export async function getTenants(params: ListTenantsParams = {}): Promise<PaginatedData<Tenant>> {
  const query = new URLSearchParams();
  const page = params.page ?? 1;
  const pageSize = params.pageSize ?? params.page_size ?? 10;
  query.set('page', String(page));
  query.set('page_size', String(pageSize));

  if (params.keyword?.trim()) query.set('keyword', params.keyword.trim());
  if (params.status && params.status !== 'all') query.set('status', params.status);
  const isolateMode = params.isolateMode ?? params.isolate_mode;
  if (isolateMode && isolateMode !== 'all') query.set('isolate_mode', isolateMode);
  if (params.sortBy ?? params.sort_by) query.set('sort_by', (params.sortBy ?? params.sort_by)!);
  if (params.order) query.set('order', params.order);

  const qs = query.toString();
  const url = qs ? `/tenants?${qs}` : '/tenants';
  const response = await apiClient.get<ApiResponse<PaginatedData<Tenant>>>(url);

  if (response.data.data) {
    const raw = response.data.data;
    const rawMeta = (raw.pagination || {}) as Record<string, unknown>;

    const currentPage = Number(rawMeta.page ?? rawMeta.page_index ?? page);
    const resolvedPageSize = Number(rawMeta.page_size ?? rawMeta.pageSize ?? pageSize);
    const totalItems = Number(rawMeta.total_items ?? rawMeta.totalItems ?? rawMeta.total ?? 0);
    const totalPages = Number(
      rawMeta.total_pages ??
      rawMeta.totalPages ??
      (totalItems > 0 ? Math.ceil(totalItems / resolvedPageSize) : 1)
    );

    const pagination: PaginationMeta = {
      page: currentPage,
      pageSize: resolvedPageSize,
      page_size: resolvedPageSize,
      totalItems,
      total_items: totalItems,
      totalPages: Math.max(1, totalPages),
      total_pages: Math.max(1, totalPages),
    };

    return {
      items: raw.items || [],
      pagination,
    };
  }

  return {
    items: [],
    pagination: {
      page,
      pageSize,
      page_size: pageSize,
      totalItems: 0,
      total_items: 0,
      totalPages: 1,
      total_pages: 1,
    },
  };
}

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

