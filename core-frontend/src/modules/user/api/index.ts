import axios from 'axios';
import { apiClient } from '@/apis/axios-client';
import { getClientLocale, type Locale } from '@/core/i18n/i18n-context';
import enErrors from '@/locales/en/errors.json';
import viErrors from '@/locales/vi/errors.json';
import type { PaginatedData, PaginationMeta } from '@/types';
import type {
  ApiResponse,
  AssignRolesInput,
  BackendResponseError,
  CloneRoleInput,
  CreateRoleInput,
  CreateUserInput,
  ListUsersParams,
  PermissionCatalogItem,
  PermissionDefinition,
  RoleDefinition,
  RoleItem,
  UpdateRoleInput,
  UpdateUserInput,
  User,
  UserStatus,
} from '../types';

const errorDictionaries: Record<Locale, Record<string, string>> = {
  en: enErrors,
  vi: viErrors,
};

export class UserApiError extends Error {
  code: string;
  fieldErrors?: Record<string, string>;
  affectedUsers?: number;

  constructor(code: string, message?: string, fieldErrors?: Record<string, string>, affectedUsers?: number) {
    super(message || code);
    this.name = 'UserApiError';
    this.code = code;
    this.fieldErrors = fieldErrors;
    this.affectedUsers = affectedUsers;
  }
}

export function parseBackendError(
  error: unknown,
  overrideLocale?: Locale
): {
  message: string;
  code?: string;
  fieldErrors?: Record<string, string>;
  affectedUsers?: number;
} {
  const currentLocale = overrideLocale || getClientLocale();
  const dict = errorDictionaries[currentLocale] || errorDictionaries.en;

  if (axios.isAxiosError(error) && error.response?.data) {
    const data = error.response.data as BackendResponseError & { message?: string };
    const code = data.message_code;
    const affectedUsers = data.affected_users;

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

    return { message, code, fieldErrors, affectedUsers };
  }

  if (error instanceof UserApiError) {
    const message = dict[error.code] || error.message;
    return { message, code: error.code, fieldErrors: error.fieldErrors, affectedUsers: error.affectedUsers };
  }

  if (error instanceof Error) {
    const message = dict[error.message] || error.message;
    return { message };
  }

  return {
    message: dict.ERR_NETWORK || 'An unexpected network error occurred.',
  };
}

export async function getUsers(params: ListUsersParams = {}): Promise<PaginatedData<User>> {
  const query = new URLSearchParams();
  const page = params.page ?? 1;
  const pageSize = params.pageSize ?? params.page_size ?? 10;
  query.set('page', String(page));
  query.set('page_size', String(pageSize));

  if (params.keyword?.trim()) query.set('keyword', params.keyword.trim());
  if (params.role && params.role !== 'all') query.set('role', params.role);
  if (params.status && params.status !== 'all') query.set('status', params.status);
  if (params.sortBy) query.set('sortBy', params.sortBy);
  if (params.order) query.set('order', params.order);

  const url = `/users?${query.toString()}`;
  const response = await apiClient.get<ApiResponse<PaginatedData<User>>>(url);

  if (response.data?.data) {
    const raw = response.data.data;
    const rawPagination = (raw as unknown as { pagination?: Record<string, unknown> }).pagination;

    const currentPage = Number(
      rawPagination?.page ??
      (rawPagination as unknown as { current_page?: unknown })?.current_page ??
      page
    );

    const resolvedPageSize = Number(
      rawPagination?.pageSize ??
      (rawPagination as unknown as { page_size?: unknown })?.page_size ??
      pageSize
    );

    const totalItems = Number(
      rawPagination?.totalItems ??
      (rawPagination as unknown as { total_items?: unknown })?.total_items ??
      raw.items?.length ??
      0
    );

    const totalPages = Number(
      rawPagination?.totalPages ??
      (rawPagination as unknown as { total_pages?: unknown })?.total_pages ??
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

export async function getUserById(id: string): Promise<User> {
  const response = await apiClient.get<ApiResponse<User>>(`/users/${encodeURIComponent(id)}`);
  if (response.data.data) {
    return response.data.data;
  }
  const code = response.data.message_code || 'E042001';
  const dict = errorDictionaries[getClientLocale()] || errorDictionaries.en;
  throw new UserApiError(code, dict[code] || dict.E042001);
}

export async function createUser(input: CreateUserInput): Promise<User> {
  const response = await apiClient.post<ApiResponse<User>>('/users', input);
  if (response.data.data) {
    return response.data.data;
  }
  const code = response.data.message_code || 'E022001';
  const dict = errorDictionaries[getClientLocale()] || errorDictionaries.en;
  throw new UserApiError(code, dict[code] || dict.E022001);
}

export async function updateUser(id: string, input: UpdateUserInput): Promise<User> {
  const response = await apiClient.put<ApiResponse<User>>(`/users/${encodeURIComponent(id)}`, input);
  if (response.data.data) {
    return response.data.data;
  }
  const code = response.data.message_code || 'E022003';
  const dict = errorDictionaries[getClientLocale()] || errorDictionaries.en;
  throw new UserApiError(code, dict[code] || dict.E022003);
}

export async function updateUserStatus(id: string, status: UserStatus): Promise<User> {
  const response = await apiClient.patch<ApiResponse<User>>(
    `/users/${encodeURIComponent(id)}/status`,
    { status }
  );
  if (response.data.data) {
    return response.data.data;
  }
  const code = response.data.message_code || 'E022003';
  const dict = errorDictionaries[getClientLocale()] || errorDictionaries.en;
  throw new UserApiError(code, dict[code] || dict.E022003);
}

export async function deleteUser(id: string): Promise<void> {
  await apiClient.delete(`/users/${encodeURIComponent(id)}`);
}

// Legacy metadata helpers
export async function getRolesMeta(): Promise<RoleDefinition[]> {
  try {
    const response = await apiClient.get<ApiResponse<RoleDefinition[]>>('/users/roles/meta');
    return response.data.data || [];
  } catch {
    return [];
  }
}

export async function getPermissionsMeta(): Promise<PermissionDefinition[]> {
  try {
    const response = await apiClient.get<ApiResponse<PermissionDefinition[]>>(
      '/users/permissions/meta'
    );
    return response.data.data || [];
  } catch {
    return [];
  }
}

// Dynamic RBAC API endpoints
export async function getRoles(): Promise<RoleItem[]> {
  const response = await apiClient.get<ApiResponse<RoleItem[]>>('/users/roles');
  return response.data.data || [];
}

export async function getRoleById(id: string): Promise<RoleItem> {
  const response = await apiClient.get<ApiResponse<RoleItem>>(`/users/roles/${encodeURIComponent(id)}`);
  if (response.data.data) {
    return response.data.data;
  }
  const code = response.data.message_code || 'E042002';
  throw new UserApiError(code);
}

export async function createRole(input: CreateRoleInput): Promise<RoleItem> {
  const response = await apiClient.post<ApiResponse<RoleItem>>('/users/roles', input);
  if (response.data.data) {
    return response.data.data;
  }
  const code = response.data.message_code || 'E032001';
  throw new UserApiError(code);
}

export async function updateRole(id: string, input: UpdateRoleInput): Promise<RoleItem> {
  const response = await apiClient.put<ApiResponse<RoleItem>>(`/users/roles/${encodeURIComponent(id)}`, input);
  if (response.data.data) {
    return response.data.data;
  }
  const code = response.data.message_code || 'E032002';
  throw new UserApiError(code);
}

export async function deleteRole(id: string): Promise<void> {
  try {
    await apiClient.delete(`/users/roles/${encodeURIComponent(id)}`);
  } catch (err: unknown) {
    if (axios.isAxiosError(err) && err.response?.data) {
      const data = err.response.data as BackendResponseError;
      throw new UserApiError(
        data.message_code || 'E032006',
        undefined,
        undefined,
        data.affected_users
      );
    }
    throw err;
  }
}

export async function cloneRole(id: string, input: CloneRoleInput): Promise<RoleItem> {
  const response = await apiClient.post<ApiResponse<RoleItem>>(`/users/roles/${encodeURIComponent(id)}/clone`, input);
  if (response.data.data) {
    return response.data.data;
  }
  const code = response.data.message_code || 'E032001';
  throw new UserApiError(code);
}

export async function assignUserRoles(userId: string, input: AssignRolesInput): Promise<void> {
  await apiClient.put(`/users/${encodeURIComponent(userId)}/roles`, input);
}

export async function getPermissionsCatalog(): Promise<PermissionCatalogItem[]> {
  const response = await apiClient.get<ApiResponse<PermissionCatalogItem[]>>('/users/permissions');
  return response.data.data || [];
}
