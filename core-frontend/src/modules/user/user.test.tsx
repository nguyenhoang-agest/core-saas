import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { userModule } from './index';
import { UserPage } from './pages/user-page';
import { parseBackendError, UserApiError } from './api';
import { TenantProvider } from '@/core/tenant/tenant-context';
import { I18nProvider } from '@/core/i18n/i18n-context';

vi.mock('next-auth/react', () => ({
  useSession: vi.fn().mockReturnValue({
    data: {
      user: {
        id: 'admin-1',
        name: 'Super Admin',
        email: 'admin@acme.com',
        permissions: ['*'],
      },
    },
    status: 'authenticated',
  }),
}));

vi.mock('@/apis/axios-client', () => ({
  apiClient: {
    get: vi.fn().mockImplementation((url: string) => {
      if (url.includes('/roles/meta') || url === '/users/roles' || url.endsWith('/roles')) {
        return Promise.resolve({
          data: {
            status: 'success',
            data: [
              { id: 'r-owner', slug: 'owner', role: 'owner', name: 'Owner', description: 'Super admin', is_system: true, version: 1, permissions: ['*'], default_permissions: ['*'], user_count: 1 },
              { id: 'r-admin', slug: 'admin', role: 'admin', name: 'Administrator', description: 'Admin', is_system: true, version: 1, permissions: ['iam:*', 'tenant:*'], default_permissions: ['iam:*', 'tenant:*'], user_count: 1 },
              { id: 'r-member', slug: 'member', role: 'member', name: 'Member', description: 'Member', is_system: true, version: 1, permissions: ['iam:view'], default_permissions: ['iam:view'], user_count: 1 },
              { id: 'r-viewer', slug: 'viewer', role: 'viewer', name: 'Viewer', description: 'Viewer', is_system: true, version: 1, permissions: ['iam:view'], default_permissions: ['iam:view'], user_count: 0 },
            ],
          },
        });
      }
      if (url.includes('/permissions/meta') || url === '/users/permissions' || url.endsWith('/permissions')) {
        return Promise.resolve({
          data: {
            status: 'success',
            data: [
              { id: 'iam:view', key: 'iam:view', module: 'iam', scope: 'tenant', name: 'View IAM & Users', description: 'Browse directory', i18n_name_key: 'iam:view', i18n_desc_key: 'Browse directory' },
              { id: 'iam:users:manage', key: 'iam:users:manage', module: 'iam', scope: 'tenant', name: 'Manage Users', description: 'Invite and edit members', i18n_name_key: 'iam:users:manage', i18n_desc_key: 'Manage members' },
            ],
          },
        });
      }
      // /users query
      return Promise.resolve({
        data: {
          status: 'success',
          data: {
            items: [
              {
                id: 'u-1',
                tenant_id: 'default',
                email: 'alice@acme.com',
                name: 'Alice Admin',
                role: 'admin',
                permissions: [],
                effective_permissions: ['iam:*', 'tenant:*'],
                status: 'active',
                create_at: '2026-10-01T00:00:00Z',
              },
              {
                id: 'u-2',
                tenant_id: 'default',
                email: 'bob@acme.com',
                name: 'Bob Member',
                role: 'member',
                permissions: ['billing:view'],
                effective_permissions: ['iam:view', 'billing:view'],
                status: 'active',
                create_at: '2026-10-02T00:00:00Z',
              },
            ],
            pagination: {
              page: 1,
              page_size: 10,
              total_items: 2,
              total_pages: 1,
            },
          },
        },
      });
    }),
    post: vi.fn().mockResolvedValue({
      data: {
        status: 'success',
        data: {
          id: 'u-new',
          email: 'new@acme.com',
          name: 'New User',
          role: 'member',
          status: 'active',
        },
      },
    }),
  },
}));

function renderWithProviders(ui: React.ReactElement, locale: 'en' | 'vi' = 'en') {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <I18nProvider initialLocale={locale}>
        <TenantProvider>{ui}</TenantProvider>
      </I18nProvider>
    </QueryClientProvider>
  );
}

describe('User Module (Frontend & RBAC Menu)', () => {
  it('registers correct ModuleDefinition properties', () => {
    expect(userModule.id).toBe('user');
    expect(userModule.name).toBe('Users');
    expect(userModule.routes).toBeDefined();
    expect(userModule.routes).toHaveLength(1);
    expect(userModule.routes![0].path).toBe('/users');
    expect(userModule.menuItems).toBeDefined();
    expect(userModule.menuItems).toHaveLength(1);
    expect(userModule.menuItems![0].key).toBe('users');
    expect(userModule.menuItems![0].path).toBe('/users');
    expect(userModule.menuItems![0].order).toBe(15);
    expect(userModule.permissions).toContain('iam:view');
    expect(userModule.permissions).toContain('iam:users:manage');
  });

  it('parses backend error codes accurately in English and Vietnamese', () => {
    const customErr = new UserApiError('E042001');
    const parsedEn = parseBackendError(customErr, 'en');
    expect(parsedEn.code).toBe('E042001');
    expect(parsedEn.message).toBe('User not found with the specified ID.');

    const parsedVi = parseBackendError(customErr, 'vi');
    expect(parsedVi.code).toBe('E042001');
    expect(parsedVi.message).toBe('Không tìm thấy người dùng với định danh đã chỉ định.');

    // RBAC error codes
    const roleInUseErr = new UserApiError('E032006');
    expect(parseBackendError(roleInUseErr, 'en').message).toBe(
      'Cannot delete role because it is currently assigned to one or more members.'
    );
    expect(parseBackendError(roleInUseErr, 'vi').message).toBe(
      'Không thể xóa vai trò vì hiện đang được gán cho một hoặc nhiều thành viên.'
    );

    const escalationErr = new UserApiError('E032009');
    expect(parseBackendError(escalationErr, 'en').message).toBe(
      'Privilege escalation violation: you cannot grant or assign permissions that you do not possess.'
    );
    expect(parseBackendError(escalationErr, 'vi').message).toBe(
      'Vi phạm leo thang đặc quyền: bạn không thể cấp hoặc gán các quyền hạn mà bạn không sở hữu.'
    );

    const notFoundRoleErr = new UserApiError('E042002');
    expect(parseBackendError(notFoundRoleErr, 'en').message).toBe(
      'Role not found with the specified ID.'
    );
    expect(parseBackendError(notFoundRoleErr, 'vi').message).toBe(
      'Không tìm thấy vai trò với mã định danh đã chỉ định.'
    );
  });

  it('renders UserPage smoke test with header, tabs, members, and responds to tab switches', async () => {
    renderWithProviders(<UserPage />);

    // Header & Tabs
    expect(screen.getByText('User Management')).toBeInTheDocument();
    expect(screen.getByText('Members & Directory')).toBeInTheDocument();
    expect(screen.getByText('Roles & Permissions Matrix')).toBeInTheDocument();

    // Data table items
    await waitFor(() => {
      expect(screen.getByText('Alice Admin')).toBeInTheDocument();
      expect(screen.getByText('Bob Member')).toBeInTheDocument();
    });

    // Switch to Roles & Permissions Matrix tab
    const matrixTab = screen.getByText('Roles & Permissions Matrix');
    fireEvent.click(matrixTab);

    await waitFor(() => {
      expect(screen.getByText('Role-Based Access Control (RBAC) Matrix')).toBeInTheDocument();
      expect(screen.getByText('View IAM & Users')).toBeInTheDocument();
    });
  });

  it('opens create user dialog when Create User button is clicked', async () => {
    renderWithProviders(<UserPage />);

    const createBtns = screen.getAllByText('Create User');
    fireEvent.click(createBtns[0]);

    await waitFor(() => {
      expect(screen.getByText('Full Name')).toBeInTheDocument();
      expect(screen.getByText('Primary Role')).toBeInTheDocument();
    });
  });
});
