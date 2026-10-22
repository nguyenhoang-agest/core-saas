export type TenantIsolateMode = 'shared' | 'dedicated';

export interface TenantDBConfig {
  dsn: string;
  max_conns?: number;
  ssl_mode?: string;
}

export interface Tenant {
  id: string;
  name: string;
  description?: string;
  status: 'active' | 'suspended' | 'deleted' | string;
  isolate_mode: TenantIsolateMode;
  db_config?: TenantDBConfig;
  create_at?: string;
  update_at?: string;
}

export interface CreateTenantInput {
  name: string;
  description?: string;
  isolate_mode?: TenantIsolateMode;
  db_config?: {
    dsn: string;
    max_conns?: number;
    ssl_mode?: string;
  };
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
}

export interface ApiResponse<T> {
  status: 'success' | 'error';
  data?: T;
  message_code?: string;
  error_detail?: FieldErrorDetail[];
}
