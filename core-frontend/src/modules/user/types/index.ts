import type { PaginationParams } from '@/types';

export type UserRole = 'owner' | 'admin' | 'member' | 'viewer' | string;
export type UserStatus = 'active' | 'suspended' | 'invited';

export interface User {
  id: string;
  tenant_id: string;
  email: string;
  name: string;
  role: UserRole;
  roles?: string[];
  permissions: string[];
  effective_permissions: string[];
  status: UserStatus;
  phone?: string;
  avatar_url?: string;
  last_login_at?: string;
  create_at: string;
  update_at: string;
}

export interface RoleDefinition {
  role: string;
  name: string;
  description: string;
  default_permissions: string[];
}

export interface RoleItem {
  id: string;
  tenant_id: string;
  slug: string;
  name: string;
  description: string;
  is_system: boolean;
  version: number;
  permissions: string[];
  user_count: number;
  create_at: string;
  update_at: string;
}

export interface PermissionDefinition {
  id: string;
  module: string;
  name: string;
  description: string;
}

export interface PermissionCatalogItem {
  key: string;
  module: string;
  scope: 'tenant' | 'platform';
  i18n_name_key: string;
  i18n_desc_key: string;
}

export interface CreateRoleInput {
  name: string;
  slug: string;
  description?: string;
  permissions: string[];
}

export interface UpdateRoleInput {
  name: string;
  description?: string;
  permissions: string[];
  version: number;
}

export interface CloneRoleInput {
  name: string;
  slug: string;
  description?: string;
}

export interface AssignRolesInput {
  role_ids: string[];
}

export interface CreateUserInput {
  name: string;
  email: string;
  role: UserRole;
  role_ids?: string[];
  permissions?: string[];
  status?: UserStatus;
  phone?: string;
}

export interface UpdateUserInput {
  name?: string;
  role?: UserRole;
  role_ids?: string[];
  permissions?: string[];
  status?: UserStatus;
  phone?: string;
}

export interface ListUsersParams extends PaginationParams {
  keyword?: string;
  role?: string;
  status?: string;
}

export interface FieldErrorDetail {
  field?: string;
  data_error?: unknown;
  message_code?: string;
}

export interface BackendResponseError {
  status: 'error';
  message_code: string;
  error_detail?: FieldErrorDetail[];
  affected_users?: number;
}

export interface ApiResponse<T> {
  status: 'success' | 'error';
  data?: T;
  message_code?: string;
  error_detail?: FieldErrorDetail[];
}
