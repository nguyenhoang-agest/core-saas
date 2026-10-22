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
    it('renders form and validates required tenant name in English with isolation mode selection', () => {
      renderWithProviders(<CreateTenantDialog open={true} onOpenChange={() => {}} />);

      expect(screen.getByText('Provision New Tenant')).toBeInTheDocument();
      expect(screen.getByLabelText(/Tenant Name/i)).toBeInTheDocument();
      expect(screen.getByText('Shared Database')).toBeInTheDocument();
      expect(screen.getByText('Dedicated DB')).toBeInTheDocument();

      // Submit without name
      const submitBtn = screen.getByRole('button', { name: /Create Tenant/i });
      fireEvent.click(submitBtn);

      expect(screen.getByText('Tenant name is required.')).toBeInTheDocument();
    });

    it('renders form and validates required tenant name in Vietnamese with isolation mode selection', () => {
      renderWithProviders(<CreateTenantDialog open={true} onOpenChange={() => {}} />, 'vi');

      expect(screen.getByText('Khởi Tạo Tenant Mới')).toBeInTheDocument();
      expect(screen.getByLabelText(/Tên Tenant/i)).toBeInTheDocument();
      expect(screen.getByText('Cơ Sở Dữ Liệu Dùng Chung')).toBeInTheDocument();
      expect(screen.getByText('DB Độc Lập (Dedicated)')).toBeInTheDocument();

      // Submit without name
      const submitBtn = screen.getByRole('button', { name: /Tạo Tenant/i });
      fireEvent.click(submitBtn);

      expect(screen.getByText('Tên Tenant là bắt buộc.')).toBeInTheDocument();
    });

    it('allows selecting Dedicated DB with DSN while strictly denying connection pool and SSL mode choices', () => {
      renderWithProviders(<CreateTenantDialog open={true} onOpenChange={() => {}} />);

      // Switch to Dedicated DB
      const dedicatedBtn = screen.getByText('Dedicated DB');
      fireEvent.click(dedicatedBtn);

      // Connection DSN is available
      expect(screen.getByLabelText(/Connection DSN/i)).toBeInTheDocument();

      // Connection pool choice is strictly DENIED to user (no editable input)
      expect(screen.queryByLabelText(/Max Connections/i)).not.toBeInTheDocument();
      expect(screen.queryByRole('spinbutton', { name: /Max Connections/i })).not.toBeInTheDocument();

      // SSL mode choice is strictly DENIED to user (no editable input)
      expect(screen.queryByLabelText(/^SSL Mode$/i)).not.toBeInTheDocument();

      // Platform managed read-only badges are shown
      expect(screen.getByTestId('connection-pool-managed')).toBeInTheDocument();
      expect(screen.getByText('10 connections')).toBeInTheDocument();

      expect(screen.getByTestId('ssl-mode-managed')).toBeInTheDocument();
      expect(screen.getByText('require')).toBeInTheDocument();
    });

    it('validates owner fields and generates a secure password', async () => {
      renderWithProviders(<CreateTenantDialog open={true} onOpenChange={() => {}} />);

      // Fill valid tenant name
      fireEvent.change(screen.getByLabelText(/Tenant Name/i), {
        target: { value: 'Wayne Enterprises' },
      });

      // Submit without filling owner fields
      const submitBtn = screen.getByRole('button', { name: /Create Tenant/i });
      fireEvent.click(submitBtn);

      expect(screen.getByText('Owner name is required.')).toBeInTheDocument();
      expect(screen.getByText('Owner email is required.')).toBeInTheDocument();
      expect(screen.getByText('Owner password is required.')).toBeInTheDocument();

      // Click "Generate" password button
      const generateBtn = screen.getByRole('button', { name: /Generate/i });
      fireEvent.click(generateBtn);

      const passwordInput = screen.getByLabelText(/Initial Password/i) as HTMLInputElement;
      expect(passwordInput.value.length).toBe(12);

      // Fill owner name and invalid email
      fireEvent.change(screen.getByLabelText(/Owner Full Name/i), {
        target: { value: 'Bruce Wayne' },
      });
      fireEvent.change(screen.getByLabelText(/Owner Email/i), {
        target: { value: 'invalid-email' },
      });

      fireEvent.click(submitBtn);
      expect(screen.getByText('Please enter a valid owner email address.')).toBeInTheDocument();
    });

    it('submits payload with owner details and triggers onSuccess', async () => {
      const onSuccessMock = vi.fn();
      renderWithProviders(<CreateTenantDialog open={true} onOpenChange={() => {}} onSuccess={onSuccessMock} />);

      fireEvent.change(screen.getByLabelText(/Tenant Name/i), {
        target: { value: 'Wayne Enterprises' },
      });
      fireEvent.change(screen.getByLabelText(/Owner Full Name/i), {
        target: { value: 'Bruce Wayne' },
      });
      fireEvent.change(screen.getByLabelText(/Owner Email/i), {
        target: { value: 'bruce@wayne.com' },
      });
      fireEvent.change(screen.getByLabelText(/Initial Password/i), {
        target: { value: 'Batman2026!' },
      });

      const submitBtn = screen.getByRole('button', { name: /Create Tenant/i });
      fireEvent.click(submitBtn);

      await waitFor(() => {
        expect(onSuccessMock).toHaveBeenCalledWith(
          expect.objectContaining({ id: 'new-id' }),
          { name: 'Bruce Wayne', email: 'bruce@wayne.com' }
        );
      });
    });

    it('allows toggling off owner creation to provision headless tenant', async () => {
      const onSuccessMock = vi.fn();
      renderWithProviders(<CreateTenantDialog open={true} onOpenChange={() => {}} onSuccess={onSuccessMock} />);

      fireEvent.change(screen.getByLabelText(/Tenant Name/i), {
        target: { value: 'Headless Corp' },
      });

      // Uncheck owner creation
      const toggleCheckbox = screen.getByRole('checkbox', { name: /Create Owner Account/i });
      fireEvent.click(toggleCheckbox);

      expect(screen.queryByLabelText(/Owner Full Name/i)).not.toBeInTheDocument();

      const submitBtn = screen.getByRole('button', { name: /Create Tenant/i });
      fireEvent.click(submitBtn);

      await waitFor(() => {
        expect(onSuccessMock).toHaveBeenCalledWith(
          expect.objectContaining({ id: 'new-id' }),
          undefined
        );
      });
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


