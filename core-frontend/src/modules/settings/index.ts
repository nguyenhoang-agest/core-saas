import type { ModuleDefinition } from '@/core/modules/types';
import { SettingsPage } from './pages/settings-page';

export * from './store/settings-store';
export * from './api/settings-api';
export * from './hooks/settings-hooks';
export * from './pages/settings-page';

export const settingsModule: ModuleDefinition = {
  id: 'settings',
  name: 'Settings',
  description: 'Manage tenant details, branding customization, and preferences',
  routes: [
    {
      path: '/settings',
      component: SettingsPage,
    },
  ],
  menuItems: [
    {
      key: 'settings',
      label: 'Settings',
      path: '/settings',
      order: 40,
      permission: 'settings:view',
    },
  ],
  permissions: ['settings:view', 'settings:branding:update'],
  slots: [],
};

export default settingsModule;
