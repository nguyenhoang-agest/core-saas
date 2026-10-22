import React from 'react';
import { render, screen, within } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AdminShell } from './admin-shell';
import { TenantProvider } from '@/core/tenant/tenant-context';
import { I18nProvider } from '@/core/i18n/i18n-context';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

let currentPath = '/';
const mockReplace = vi.fn();

vi.mock('next/navigation', () => ({
  usePathname: () => currentPath,
  useRouter: () => ({
    replace: mockReplace,
    push: vi.fn(),
  }),
}));

let mockIsHost = false;
let mockStatus: 'authenticated' | 'unauthenticated' | 'loading' = 'authenticated';

vi.mock('next-auth/react', () => ({
  useSession: () => ({
    data:
      mockStatus === 'authenticated'
        ? {
            user: {
              id: 'user-1',
              name: 'Demo Admin',
              email: 'admin@example.com',
              isHost: mockIsHost,
            },
          }
        : null,
    status: mockStatus,
  }),
  signOut: vi.fn(),
  signIn: vi.fn(),
}));

function renderWithShell(ui: React.ReactNode = <div>Content</div>, locale: 'en' | 'vi' = 'en') {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <I18nProvider initialLocale={locale}>
        <TenantProvider>
          <AdminShell>{ui}</AdminShell>
        </TenantProvider>
      </I18nProvider>
    </QueryClientProvider>
  );
}

describe('AdminShell Breadcrumb Bar', () => {
  beforeEach(() => {
    currentPath = '/';
    mockIsHost = false;
    mockStatus = 'authenticated';
    mockReplace.mockClear();
  });

  it('renders root breadcrumb on dashboard with clickable Admin link', () => {
    currentPath = '/';
    renderWithShell();

    const breadcrumbNav = screen.getByRole('navigation', { name: /Breadcrumb/i });
    expect(breadcrumbNav).toBeInTheDocument();

    const adminLink = within(breadcrumbNav).getByRole('link', { name: /Admin/i });
    expect(adminLink).toHaveAttribute('href', '/');

    const activeItem = within(breadcrumbNav).getByText('Dashboard');
    expect(activeItem).toHaveAttribute('aria-current', 'page');
  });

  it('renders breadcrumb on /tenants with parent link to root', () => {
    currentPath = '/tenants';
    renderWithShell();

    const breadcrumbNav = screen.getByRole('navigation', { name: /Breadcrumb/i });

    const adminLink = within(breadcrumbNav).getByRole('link', { name: /Admin/i });
    expect(adminLink).toHaveAttribute('href', '/');

    const activeItem = within(breadcrumbNav).getByText('Tenants');
    expect(activeItem).toHaveAttribute('aria-current', 'page');
  });

  it('renders nested breadcrumbs on /settings/general', () => {
    currentPath = '/settings/general';
    renderWithShell();

    const breadcrumbNav = screen.getByRole('navigation', { name: /Breadcrumb/i });

    const adminLink = within(breadcrumbNav).getByRole('link', { name: /Admin/i });
    expect(adminLink).toHaveAttribute('href', '/');

    const settingsLink = within(breadcrumbNav).getByRole('link', { name: /Settings/i });
    expect(settingsLink).toHaveAttribute('href', '/settings');

    const activeItem = within(breadcrumbNav).getByText('General');
    expect(activeItem).toHaveAttribute('aria-current', 'page');
  });

  it('renders localized breadcrumbs in Vietnamese', () => {
    currentPath = '/tenants';
    renderWithShell(undefined, 'vi');

    const breadcrumbNav = screen.getByRole('navigation', { name: /Breadcrumb/i });

    const adminLink = within(breadcrumbNav).getByRole('link', { name: /Quản trị/i });
    expect(adminLink).toHaveAttribute('href', '/');

    const activeItem = within(breadcrumbNav).getByText('Không gian Tenant');
    expect(activeItem).toHaveAttribute('aria-current', 'page');
  });

  it('handles platform module /billing breadcrumb', () => {
    currentPath = '/billing';
    renderWithShell();

    const breadcrumbNav = screen.getByRole('navigation', { name: /Breadcrumb/i });
    const billingItem = within(breadcrumbNav).getByText('Billing');
    expect(billingItem).toHaveAttribute('aria-current', 'page');
  });

  it('handles platform module /users breadcrumb', () => {
    currentPath = '/users';
    renderWithShell();

    const breadcrumbNav = screen.getByRole('navigation', { name: /Breadcrumb/i });
    const usersItem = within(breadcrumbNav).getByText('Users');
    expect(usersItem).toHaveAttribute('aria-current', 'page');
  });

  it('bypasses admin shell chrome on /login route', () => {
    currentPath = '/login';
    renderWithShell(<div data-testid="login-content">Login Form</div>);

    expect(screen.getByTestId('login-content')).toBeInTheDocument();
    expect(screen.queryByRole('navigation', { name: /Breadcrumb/i })).not.toBeInTheDocument();
  });

  it('renders user profile and sign out button in header', () => {
    currentPath = '/';
    renderWithShell();

    expect(screen.getByText('Demo Admin')).toBeInTheDocument();
    expect(screen.getByTitle('Sign Out')).toBeInTheDocument();
  });

  it('hides /tenants link and tenant switcher for non-host tenant user', () => {
    mockIsHost = false;
    renderWithShell();

    const aside = screen.getByRole('complementary');
    expect(within(aside).queryByRole('link', { name: /Tenants/i })).not.toBeInTheDocument();
    expect(screen.queryByTitle('Switch active organization / tenant')).not.toBeInTheDocument();
  });

  it('shows /tenants link and tenant switcher for host administrator', () => {
    mockIsHost = true;
    renderWithShell();

    const aside = screen.getByRole('complementary');
    expect(within(aside).getByRole('link', { name: /Tenants/i })).toBeInTheDocument();
    expect(screen.getByTitle('Switch active organization / tenant')).toBeInTheDocument();
  });

  it('restricts unauthenticated guest to Overview only in navigation and shows guest banner', () => {
    mockStatus = 'unauthenticated';
    currentPath = '/';
    renderWithShell();

    const aside = screen.getByRole('complementary');
    // Overview link must be present
    expect(within(aside).getByRole('link', { name: /Overview/i })).toBeInTheDocument();

    // All management pages must be hidden
    expect(within(aside).queryByRole('link', { name: /Tenants/i })).not.toBeInTheDocument();
    expect(within(aside).queryByRole('link', { name: /Users/i })).not.toBeInTheDocument();
    expect(within(aside).queryByRole('link', { name: /Settings/i })).not.toBeInTheDocument();

    // Guest mode banner and Sign In button shown in footer
    expect(within(aside).getByText('Browsing as Guest')).toBeInTheDocument();
    expect(within(aside).getByRole('link', { name: /Sign In/i })).toBeInTheDocument();
  });

  it('triggers redirect to /login when unauthenticated guest accesses protected route', () => {
    mockStatus = 'unauthenticated';
    currentPath = '/users';
    renderWithShell(<div data-testid="protected-content">Secret Users Content</div>);

    // Protected content should NOT be rendered
    expect(screen.queryByTestId('protected-content')).not.toBeInTheDocument();

    // Redirecting message rendered
    expect(
      screen.getByText('Sign in required. Redirecting to login...')
    ).toBeInTheDocument();

    // Router redirect triggered
    expect(mockReplace).toHaveBeenCalledWith('/login?callbackUrl=%2Fusers');
  });
});
