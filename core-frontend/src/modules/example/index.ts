import type { ModuleDefinition } from '@/core/modules/types';
import { ExamplePage } from './pages/example-page';
import { ExampleWidget } from './components/example-widget';

export * from './store/example-store';
export * from './api/example-api';
export * from './hooks/example-hooks';
export * from './pages/example-page';
export * from './components/example-widget';

export const exampleModule: ModuleDefinition = {
  id: 'example',
  name: 'Example Slice',
  description: 'Thin vertical slice showcasing routes, store, API, and slots',
  routes: [
    {
      path: '/example',
      component: ExamplePage,
    },
  ],
  menuItems: [
    {
      key: 'example',
      label: 'Example Slice',
      path: '/example',
      order: 10,
      permission: 'example:view',
    },
  ],
  permissions: ['example:view', 'example:write'],
  slots: [
    {
      id: 'example-header-widget',
      slotName: 'header-actions',
      order: 10,
      component: ExampleWidget,
    },
  ],
};

export default exampleModule;
