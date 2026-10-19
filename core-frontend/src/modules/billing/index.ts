import type { ModuleDefinition } from '@/core/modules/types';
import { BillingPage } from './pages/billing-page';

export * from './store/billing-store';
export * from './api/billing-api';
export * from './hooks/billing-hooks';
export * from './pages/billing-page';

export const billingModule: ModuleDefinition = {
  id: 'billing',
  name: 'Billing',
  description: 'Manage subscription plans, invoices, and payment methods',
  routes: [
    {
      path: '/billing',
      component: BillingPage,
    },
  ],
  menuItems: [
    {
      key: 'billing',
      label: 'Billing',
      path: '/billing',
      order: 30,
      permission: 'billing:view',
    },
  ],
  permissions: ['billing:view', 'billing:plans:update', 'billing:invoices:view'],
  slots: [],
};

export default billingModule;
