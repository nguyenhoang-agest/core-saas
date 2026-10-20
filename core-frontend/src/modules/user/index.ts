import type { ModuleDefinition } from '@/core/modules/types';
import { UserPage } from './pages/user-page';

export * from './types';
export * from './api';
export * from './hooks';
export * from './components';
export * from './store/user-store';
export * from './pages/user-page';

export const userModule: ModuleDefinition = {
  id: 'user',
  name: 'Users',
  description: 'Multi-tenant user directory, RBAC roles, and granular ACL permissions',
  routes: [
    {
      path: '/users',
      component: UserPage,
    },
  ],
  menuItems: [
    {
      key: 'users',
      label: 'Users',
      path: '/users',
      order: 15,
      permission: 'iam:view',
    },
  ],
  permissions: ['iam:view', 'iam:users:manage', 'iam:roles:manage'],
  slots: [],
};

export default userModule;
