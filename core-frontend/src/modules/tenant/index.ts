import type { ModuleDefinition } from '@/core/modules/types';
import { TenantAdminPage } from './pages/tenant-admin-page';

export * from './types/tenant';
export * from './api/tenant-api';
export * from './pages/tenant-admin-page';
export * from './components/tenant-kpi-cards';
export * from './components/tenant-inspector';
export * from './components/tenant-list-table';
export * from './components/create-tenant-dialog';

export const tenantModule: ModuleDefinition = {
  id: 'tenant',
  name: 'Tenant Management',
  description: 'Manage tenant workspaces, database isolation modes, and runtime contexts',
  routes: [
    {
      path: '/tenants',
      component: TenantAdminPage,
    },
  ],
  menuItems: [
    {
      key: 'tenants',
      label: 'Tenants',
      path: '/tenants',
      order: 5,
      permission: 'tenant:view',
    },
  ],
  permissions: ['tenant:view', 'tenant:create', 'tenant:update', 'tenant:delete'],
};

export default tenantModule;
