import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { TenantProvider, useTenant } from './tenant-context';

let mockSession: any = null;
let mockStatus: 'authenticated' | 'unauthenticated' | 'loading' = 'unauthenticated';

vi.mock('next-auth/react', () => ({
  useSession: () => ({
    data: mockSession,
    status: mockStatus,
  }),
}));

function TenantConsumer() {
  const { tenant } = useTenant();
  return (
    <div>
      <span data-testid="tenant-id">{tenant.id}</span>
      <span data-testid="tenant-name">{tenant.name}</span>
    </div>
  );
}

describe('TenantProvider Session Synchronization', () => {
  beforeEach(() => {
    mockSession = null;
    mockStatus = 'unauthenticated';
    if (typeof document !== 'undefined') {
      document.cookie = 'x-tenant-id=; max-age=0; path=/';
      localStorage.clear();
    }
  });

  it('renders default tenant when unauthenticated', () => {
    render(
      <TenantProvider>
        <TenantConsumer />
      </TenantProvider>
    );

    expect(screen.getByTestId('tenant-id').textContent).toBe('default');
    expect(screen.getByTestId('tenant-name').textContent).toBe('Default Organization');
  });

  it('enforces and locks tenant to owner session tenant, overwriting stale admin tenant', () => {
    // Simulate stale admin tenant in cookie and localStorage
    document.cookie = 'x-tenant-id=default; path=/';
    localStorage.setItem('core_saas_tenant_id', 'default');
    localStorage.setItem('core_saas_tenant_name', 'Platform Host');

    // Tenant owner logs in
    mockStatus = 'authenticated';
    mockSession = {
      user: {
        id: 'user-stark-1',
        name: 'Tony Stark',
        email: 'tony@stark.com',
        role: 'owner',
        tenantId: 'tenant-stark-888',
        tenantName: 'Stark Industries',
        isHost: false,
      },
    };

    render(
      <TenantProvider>
        <TenantConsumer />
      </TenantProvider>
    );

    // Context should immediately reflect owner's workspace
    expect(screen.getByTestId('tenant-id').textContent).toBe('tenant-stark-888');
    expect(screen.getByTestId('tenant-name').textContent).toBe('Stark Industries');

    // Cookie and localStorage should be synced
    expect(document.cookie).toContain('x-tenant-id=tenant-stark-888');
    expect(localStorage.getItem('core_saas_tenant_name')).toBe('Stark Industries');
  });

  it('allows host admin to retain or switch tenant', () => {
    // Host admin switches to a specific tenant
    document.cookie = 'x-tenant-id=tenant-custom-999; path=/';
    localStorage.setItem('core_saas_tenant_name', 'Custom Client Org');

    mockStatus = 'authenticated';
    mockSession = {
      user: {
        id: 'admin-root',
        name: 'Platform Admin',
        email: 'admin@example.com',
        role: 'owner',
        tenantId: 'default',
        tenantName: 'Platform Host',
        isHost: true,
      },
    };

    render(
      <TenantProvider>
        <TenantConsumer />
      </TenantProvider>
    );

    expect(screen.getByTestId('tenant-id').textContent).toBe('tenant-custom-999');
    expect(screen.getByTestId('tenant-name').textContent).toBe('Custom Client Org');
  });
});
