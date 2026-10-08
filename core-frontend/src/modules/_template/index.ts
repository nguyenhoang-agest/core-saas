import type { ModuleDefinition } from '@/core/modules/types';
import { TemplatePage } from './pages/template-page';

export * from './store/template-store';
export * from './api/template-api';
export * from './pages/template-page';

export const templateModule: ModuleDefinition = {
  id: '__MODULE_NAME__',
  name: '__MODULE_TITLE__',
  description: 'Scaffolded starter module for __MODULE_TITLE__',
  routes: [
    {
      path: '/__MODULE_NAME__',
      component: TemplatePage,
    },
  ],
  menuItems: [
    {
      key: '__MODULE_NAME__',
      label: '__MODULE_TITLE__',
      path: '/__MODULE_NAME__',
      order: 50,
      permission: '__MODULE_NAME__:view',
    },
  ],
  permissions: ['__MODULE_NAME__:view', '__MODULE_NAME__:manage'],
  slots: [],
};

export default templateModule;
