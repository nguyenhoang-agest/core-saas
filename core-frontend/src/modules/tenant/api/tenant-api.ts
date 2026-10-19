import axios from 'axios';
import { apiClient } from '@/apis/axios-client';
import { getClientLocale, type Locale } from '@/core/i18n/i18n-context';
import enErrors from '@/locales/en/errors.json';
import viErrors from '@/locales/vi/errors.json';
import type { PaginatedData, PaginationMeta, PaginationParams } from '@/types';
import type { ApiResponse, BackendResponseError, CreateTenantInput, Tenant } from '../types/tenant';

const errorDictionaries: Record<Locale, Record<string, string>> = {
  en: enErrors,
  vi: viErrors,
};

export class TenantApiError extends Error {
  code: string;
  fieldErrors?: Record<string, string>;

  constructor(code: string, message?: string, fieldErrors?: Record<string, string>) {
    super(message || code);
    this.name = 'TenantApiError';
    this.code = code;
    this.fieldErrors = fieldErrors;
  }
}

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

    // Two-tier error resolution: local dictionary -> backend localized message -> localized fallback
    const fallbackMessage = code
      ? (dict.ERR_OPERATION_FAILED || 'Operation failed with error code: {code}').replace(
          '{code}',
          code
        )
      : dict.ERR_UNEXPECTED || 'An unexpected error occurred.';

    const message = (code && dict[code]) || data.message || fallbackMessage;

    const fieldErrors: Record<string, string> = {};
    if (data.error_detail && Array.isArray(data.error_detail)) {
      for (const item of data.error_detail) {
        if (item.field) {
          fieldErrors[item.field] =
            (item.message_code && dict[item.message_code]) ||
            (item as { message?: string }).message ||
            dict.ERR_INVALID_FIELD ||
            'Invalid value';
        }
      }
    }

    return { message, code, fieldErrors };
  }

  if (error instanceof TenantApiError) {
    const message = dict[error.code] || error.message;
    return { message, code: error.code, fieldErrors: error.fieldErrors };
  }

  if (error instanceof Error) {
    const message = dict[error.message] || error.message;
    return { message };
  }

  return {
    message: dict.ERR_NETWORK || 'An unexpected network error occurred.',
  };
}

export async function getTenantById(id: string): Promise<Tenant> {
  const response = await apiClient.get<ApiResponse<Tenant>>(`/tenants/${encodeURIComponent(id)}`);
  if (response.data.data) {
    return response.data.data;
  }
  const code = response.data.message_code || 'E021002';
  const dict = errorDictionaries[getClientLocale()] || errorDictionaries.en;
  throw new TenantApiError(code, dict[code] || dict.E021002);
}

export async function createTenant(input: CreateTenantInput): Promise<Tenant> {
  const response = await apiClient.post<ApiResponse<Tenant>>('/tenants', input);
  if (response.data.data) {
    return response.data.data;
  }
  const code = response.data.message_code || 'E021001';
  const dict = errorDictionaries[getClientLocale()] || errorDictionaries.en;
  throw new TenantApiError(code, dict[code] || dict.E021001);
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
    const p = raw.pagination;
    const resolvedPage = Number(p?.page ?? page);
    const resolvedPageSize = Number(p?.pageSize ?? p?.page_size ?? pageSize);
    const totalItems = Number(p?.totalItems ?? p?.total_items ?? 0);
    const totalPages = Math.max(
      1,
      Number(
        p?.totalPages ??
          p?.total_pages ??
          (totalItems > 0 ? Math.ceil(totalItems / resolvedPageSize) : 1)
      )
    );

    const pagination: PaginationMeta = {
      page: resolvedPage,
      pageSize: resolvedPageSize,
      page_size: resolvedPageSize,
      totalItems,
      total_items: totalItems,
      totalPages,
      total_pages: totalPages,
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
