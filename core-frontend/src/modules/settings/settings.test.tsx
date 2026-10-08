import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import React from 'react';
import { settingsModule } from './index';
import { SettingsPage } from './pages/settings-page';
import { TenantProvider } from '@/core/tenant/tenant-context';

describe('Settings Module (Stub)', () => {
  it('has a valid ModuleDefinition structure', () => {
    expect(settingsModule.id).toBe('settings');
    expect(settingsModule.routes).toHaveLength(1);
    expect(settingsModule.menuItems).toHaveLength(1);
    expect(settingsModule.permissions).toContain('settings:view');
  });

  it('renders SettingsPage smoke test without crashing', () => {
    render(
      <TenantProvider>
        <SettingsPage />
      </TenantProvider>
    );
    expect(screen.getByText('Tenant Settings')).toBeInTheDocument();
    expect(screen.getByText('Settings Module Stub')).toBeInTheDocument();
  });
});
