import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import React from 'react';
import { billingModule } from './index';
import { BillingPage } from './pages/billing-page';
import { TenantProvider } from '@/core/tenant/tenant-context';

describe('Billing Module (Stub)', () => {
  it('has a valid ModuleDefinition structure', () => {
    expect(billingModule.id).toBe('billing');
    expect(billingModule.routes).toHaveLength(1);
    expect(billingModule.menuItems).toHaveLength(1);
    expect(billingModule.permissions).toContain('billing:view');
  });

  it('renders BillingPage smoke test without crashing', () => {
    render(
      <TenantProvider>
        <BillingPage />
      </TenantProvider>
    );
    expect(screen.getByText('Billing & Subscriptions')).toBeInTheDocument();
    expect(screen.getByText('Billing Module Stub')).toBeInTheDocument();
  });
});
