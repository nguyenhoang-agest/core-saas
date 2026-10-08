import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import React from 'react';
import { iamModule } from './index';
import { IamPage } from './pages/iam-page';
import { TenantProvider } from '@/core/tenant/tenant-context';

describe('IAM Module (Stub)', () => {
  it('has a valid ModuleDefinition structure', () => {
    expect(iamModule.id).toBe('iam');
    expect(iamModule.routes).toHaveLength(1);
    expect(iamModule.menuItems).toHaveLength(1);
    expect(iamModule.permissions).toContain('iam:view');
  });

  it('renders IamPage smoke test without crashing', () => {
    render(
      <TenantProvider>
        <IamPage />
      </TenantProvider>
    );
    expect(screen.getByText('Identity & Access Management (IAM)')).toBeInTheDocument();
    expect(screen.getByText('IAM Module Stub')).toBeInTheDocument();
  });
});
