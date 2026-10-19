import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { parseBackendError } from './api/tenant-api';
import { TenantKpiCards } from './components/tenant-kpi-cards';
import { CreateTenantDialog } from './components/create-tenant-dialog';
import { TenantListTable } from './components/tenant-list-table';
import { TenantProvider } from '@/core/tenant/tenant-context';
import { I18nProvider } from '@/core/i18n/i18n-context';

vi.mock('@/apis/axios-client', () => ({
  apiClient: {
    get: vi.fn().mockImplementation((url: string) => {
      if (url.startsWith('/tenants/')) {
        return Promise.resolve({
          data: {
            status: 'success',
            data: {
              id: '0192e21b-8526-70e2-a044-cf66a2985f47',
              name: 'Acme Global',
              isolate_mode: 'shared',
              status: 'active',
            },
          },
        });
      }
      const search = url.includes('?') ? url.split('?')[1] : '';
      const params = new URLSearchParams(search);
      const page = Number(params.get('page') || 1);

      return Promise.resolve({
        data: {
          status: 'success',
          data: {
            items: [
              {
                id: `t-${page}-1`,
                name: 'Acme Corp',
                description: 'Main production workspace',
                isolate_mode: 'shared',
                status: 'active',
                create_at: '2026-10-01T00:00:00Z',
              },
              {
                id: `t-${page}-2`,
                name: 'Beta Dedicated',
                description: 'Enterprise dedicated DB',
                isolate_mode: 'dedicated',
                status: 'active',
                create_at: '2026-10-02T00:00:00Z',
              },
            ],
            pagination: {
              page,
              page_size: 10,
              total_items: 25,
              total_pages: 3,
            },
          },
        },
      });
    }),
    post: vi.fn().mockResolvedValue({
      data: {
        status: 'success',
        data: {
          id: 'new-id',
          name: 'New Tenant',
          isolate_mode: 'shared',
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

describe('Tenant Module', () => {
  describe('parseBackendError', () => {
    it('maps E031002 (TenantAlreadyExists) to friendly error message in English and Vietnamese', () => {
      const axiosError = {
        isAxiosError: true,
        response: {
          data: {
            status: 'error',
            message_code: 'E031002',
          },
        },
      };

      const resultEn = parseBackendError(axiosError, 'en');
      expect(resultEn.code).toBe('E031002');
      expect(resultEn.message).toBe('A tenant with this name already exists.');

      const resultVi = parseBackendError(axiosError, 'vi');
      expect(resultVi.code).toBe('E031002');
      expect(resultVi.message).toBe('Tên Tenant này đã tồn tại trong hệ thống.');
    });

    it('maps E031004 (MissingDedicatedDBConfig) correctly', () => {
      const axiosError = {
        isAxiosError: true,
        response: {
          data: {
            status: 'error',
            message_code: 'E031004',
          },
        },
      };

      const result = parseBackendError(axiosError);
      expect(result.code).toBe('E031004');
      expect(result.message).toBe(
        'Database DSN is required when dedicated isolation mode is selected.'
      );
    });

    it('maps E041001 (NotFoundTenant) correctly', () => {
      const axiosError = {
        isAxiosError: true,
        response: {
          data: {
            status: 'error',
            message_code: 'E041001',
          },
        },
      };

      const result = parseBackendError(axiosError);
      expect(result.code).toBe('E041001');
      expect(result.message).toBe('Tenant not found with the specified ID.');
    });

    it('maps field-level validation errors (E050000)', () => {
      const axiosError = {
        isAxiosError: true,
        response: {
          data: {
            status: 'error',
            message_code: 'E050000',
            error_detail: [
              {
                field: 'name',
                message_code: 'E050001',
              },
            ],
          },
        },
      };

      const result = parseBackendError(axiosError);
      expect(result.code).toBe('E050000');
      expect(result.fieldErrors?.name).toBe('This field is required.');
    });
  });

  describe('TenantKpiCards', () => {
    it('renders tenant kpi cards in English', () => {
      renderWithProviders(<TenantKpiCards />);
      expect(screen.getByText('Active Context')).toBeInTheDocument();
      expect(screen.getByText('Supported Isolation')).toBeInTheDocument();
      expect(screen.getByText('3-Layer Errors')).toBeInTheDocument();
    });

    it('renders tenant kpi cards in Vietnamese', () => {
      renderWithProviders(<TenantKpiCards />, 'vi');
      expect(screen.getByText('Ngữ Cảnh Hiện Tại')).toBeInTheDocument();
      expect(screen.getByText('Chế Độ Cô Lập')).toBeInTheDocument();
      expect(screen.getByText('Hệ Thống Lỗi 3 Tầng')).toBeInTheDocument();
    });
  });

  describe('CreateTenantDialog', () => {
    it('renders form and validates required tenant name in English', () => {
      renderWithProviders(<CreateTenantDialog open={true} onOpenChange={() => {}} />);

      expect(screen.getByText('Provision New Tenant')).toBeInTheDocument();
      expect(screen.getByText('Shared Database')).toBeInTheDocument();
      expect(screen.getByText('Dedicated DB')).toBeInTheDocument();

      // Submit without name
      const submitBtn = screen.getByRole('button', { name: /Create Tenant/i });
      fireEvent.click(submitBtn);

      expect(screen.getByText('Tenant name is required.')).toBeInTheDocument();
    });

    it('renders form and validates required tenant name in Vietnamese', () => {
      renderWithProviders(<CreateTenantDialog open={true} onOpenChange={() => {}} />, 'vi');

      expect(screen.getByText('Khởi Tạo Tenant Mới')).toBeInTheDocument();
      expect(screen.getByText('Cơ Sở Dữ Liệu Dùng Chung')).toBeInTheDocument();
      expect(screen.getByText('DB Độc Lập (Dedicated)')).toBeInTheDocument();

      // Submit without name
      const submitBtn = screen.getByRole('button', { name: /Tạo Tenant/i });
      fireEvent.click(submitBtn);

      expect(screen.getByText('Tên Tenant là bắt buộc.')).toBeInTheDocument();
    });

    it('unfolds dedicated db inputs when clicking Dedicated DB', () => {
      renderWithProviders(<CreateTenantDialog open={true} onOpenChange={() => {}} />);

      const dedicatedBtn = screen.getByText('Dedicated DB');
      fireEvent.click(dedicatedBtn);

      expect(screen.getByText('Connection DSN')).toBeInTheDocument();
      expect(screen.getByText('Max Connections')).toBeInTheDocument();
    });
  });

  describe('TenantListTable', () => {
    it('renders list table with rows in English', async () => {
      renderWithProviders(<TenantListTable />);

      expect(screen.getByText('Tenant Workspaces')).toBeInTheDocument();
      expect(screen.getByPlaceholderText('Search tenants by name...')).toBeInTheDocument();

      await waitFor(() => {
        expect(screen.getByText('Acme Corp')).toBeInTheDocument();
        expect(screen.getByText('Beta Dedicated')).toBeInTheDocument();
      });

      expect(screen.getByText('Showing 1 - 10 of 25 tenants')).toBeInTheDocument();
    });

    it('renders list table with rows in Vietnamese', async () => {
      renderWithProviders(<TenantListTable />, 'vi');

      expect(screen.getByText('Danh Sách Workspace Tenant')).toBeInTheDocument();
      expect(screen.getByPlaceholderText('Tìm kiếm tenant theo tên...')).toBeInTheDocument();

      await waitFor(() => {
        expect(screen.getByText('Acme Corp')).toBeInTheDocument();
        expect(screen.getByText('Beta Dedicated')).toBeInTheDocument();
      });

      expect(screen.getByText('Hiển thị 1 - 10 trên tổng số 25 tenant')).toBeInTheDocument();
    });

    it('calls onInspectTenant when inspect button is clicked', async () => {
      const inspectSpy = vi.fn();
      renderWithProviders(<TenantListTable onInspectTenant={inspectSpy} />);

      await waitFor(() => {
        expect(screen.getByText('Acme Corp')).toBeInTheDocument();
      });

      const inspectBtns = screen.getAllByRole('button', { name: /Inspect/i });
      fireEvent.click(inspectBtns[0]);

      expect(inspectSpy).toHaveBeenCalledWith(
        expect.objectContaining({
          id: 't-1-1',
          name: 'Acme Corp',
        })
      );
    });

    it('handles pagination next and previous controls correctly', async () => {
      renderWithProviders(<TenantListTable />);

      await waitFor(() => {
        expect(screen.getByText('Acme Corp')).toBeInTheDocument();
      });

      expect(screen.getByText('Showing 1 - 10 of 25 tenants')).toBeInTheDocument();
      expect(screen.getByText('Page 1 of 3')).toBeInTheDocument();

      // On page 1: Prev should be disabled, Next should be enabled
      const prevBtn = screen.getByRole('button', { name: /Previous/i });
      const nextBtn = screen.getByRole('button', { name: /Next/i });

      expect(prevBtn).toBeDisabled();
      expect(nextBtn).toBeEnabled();

      // Click Next -> should advance to page 2
      fireEvent.click(nextBtn);

      await waitFor(() => {
        expect(screen.getByText('Page 2 of 3')).toBeInTheDocument();
      });

      expect(screen.getByText('Showing 11 - 20 of 25 tenants')).toBeInTheDocument();
      expect(prevBtn).toBeEnabled();
      expect(nextBtn).toBeEnabled();

      // Click Prev -> should go back to page 1
      fireEvent.click(prevBtn);

      await waitFor(() => {
        expect(screen.getByText('Page 1 of 3')).toBeInTheDocument();
      });

      expect(prevBtn).toBeDisabled();
      expect(nextBtn).toBeEnabled();
    });
  });
});


