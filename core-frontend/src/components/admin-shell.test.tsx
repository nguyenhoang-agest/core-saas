import React from 'react';
import { render, screen, within } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AdminShell } from './admin-shell';
import { TenantProvider } from '@/core/tenant/tenant-context';
import { I18nProvider } from '@/core/i18n/i18n-context';

let currentPath = '/';

vi.mock('next/navigation', () => ({
  usePathname: () => currentPath,
}));

function renderWithShell(ui: React.ReactNode = <div>Content</div>, locale: 'en' | 'vi' = 'en') {
  return render(
    <I18nProvider initialLocale={locale}>
      <TenantProvider>
        <AdminShell>{ui}</AdminShell>
      </TenantProvider>
    </I18nProvider>
  );
}

describe('AdminShell Breadcrumb Bar', () => {
  beforeEach(() => {
    currentPath = '/';
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

  it('handles other platform modules like /billing and /iam', () => {
    currentPath = '/billing';
    renderWithShell();

    const breadcrumbNav = screen.getByRole('navigation', { name: /Breadcrumb/i });

    const billingItem = within(breadcrumbNav).getByText('Billing');
    expect(billingItem).toHaveAttribute('aria-current', 'page');
  });
});
