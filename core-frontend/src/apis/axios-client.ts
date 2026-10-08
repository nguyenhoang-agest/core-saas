import axios from 'axios';

/**
 * Extracts the tenant ID from the document cookies in client-side runtime,
 * or returns 'default' as fallback.
 */
export function getClientTenantId(): string {
  if (typeof document !== 'undefined') {
    const match = document.cookie.match(/(?:^|; )x-tenant-id=([^;]*)/);
    if (match && match[1]) {
      return decodeURIComponent(match[1]);
    }
  }
  return 'default';
}

export const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || '/api',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
});

apiClient.interceptors.request.use(
  (config) => {
    const tenantId = getClientTenantId();
    if (tenantId) {
      config.headers['X-Tenant-ID'] = tenantId;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    return Promise.reject(error);
  }
);
