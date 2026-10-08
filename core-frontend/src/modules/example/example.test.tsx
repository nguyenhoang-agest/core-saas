import { render, screen, act } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ExamplePage } from './pages/example-page';
import { ExampleWidget } from './components/example-widget';
import { exampleModule } from './index';
import { TenantProvider } from '@/core/tenant/tenant-context';

vi.mock('next-auth/react', () => ({
  useSession: () => ({
    data: {
      user: {
        permissions: ['example:view', 'example:write'],
      },
    },
    status: 'authenticated',
  }),
}));

vi.mock('@/apis/axios-client', () => ({
  apiClient: {
    get: vi.fn().mockResolvedValue({
      data: [
        { id: '1', title: 'Setup Authentication & RBAC', status: 'completed', createdAt: '2026-10-01' },
      ],
    }),
  },
}));

function renderWithProviders(ui: React.ReactElement) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
    },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <TenantProvider initialTenant={{ id: 'test-tenant', name: 'Test Tenant', theme: { colorPrimary: '#1677ff' } }}>
        {ui}
      </TenantProvider>
    </QueryClientProvider>
  );
}

describe('Example Module Slice', () => {
  it('has a valid ModuleDefinition structure', () => {
    expect(exampleModule.id).toBe('example');
    expect(exampleModule.routes).toHaveLength(1);
    expect(exampleModule.menuItems).toHaveLength(1);
    expect(exampleModule.slots).toHaveLength(1);
    expect(exampleModule.permissions).toContain('example:view');
  });

  it('renders ExampleWidget correctly', () => {
    render(<ExampleWidget />);
    expect(screen.getByText('Example Slice Ready')).toBeInTheDocument();
  });

  it('renders ExamplePage and responds to store actions', async () => {
    renderWithProviders(<ExamplePage />);

    expect(screen.getByText('Example Module Slice')).toBeInTheDocument();
    expect(screen.getByText(/Active Tenant:/)).toBeInTheDocument();

    const incrementBtn = screen.getByRole('button', { name: /Store Action Clicks: 0/i });
    expect(incrementBtn).toBeInTheDocument();

    act(() => {
      incrementBtn.click();
    });

    expect(screen.getByRole('button', { name: /Store Action Clicks: 1/i })).toBeInTheDocument();
  });

  it('shows protected action button when user has permission', () => {
    renderWithProviders(<ExamplePage />);
    expect(screen.getByRole('button', { name: /Protected Mutation Action/i })).toBeInTheDocument();
  });
});
