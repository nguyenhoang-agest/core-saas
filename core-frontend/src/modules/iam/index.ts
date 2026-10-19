import type { ModuleDefinition } from '@/core/modules/types';
import { IamPage } from './pages/iam-page';

export * from './store/iam-store';
export * from './api/iam-api';
export * from './hooks/iam-hooks';
export * from './pages/iam-page';

export const iamModule: ModuleDefinition = {
  id: 'iam',
  name: 'Identity & Access',
  description: 'Manage users, teams, invites, and role-based permissions',
  routes: [
    {
      path: '/iam',
      component: IamPage,
    },
  ],
  menuItems: [
    {
      key: 'iam',
      label: 'IAM',
      path: '/iam',
      order: 20,
      permission: 'iam:view',
    },
  ],
  permissions: ['iam:view', 'iam:users:manage', 'iam:roles:manage'],
  slots: [],
};

export default iamModule;
