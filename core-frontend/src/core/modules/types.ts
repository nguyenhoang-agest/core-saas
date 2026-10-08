import type { ComponentType } from 'react';

export interface ModuleMenuItem {
  key: string;
  label: string;
  path: string;
  icon?: string;
  order?: number;
  permission?: string;
}

export interface ModuleRoute {
  path: string;
  component: ComponentType<unknown>;
}

export interface SlotContribution<P = Record<string, unknown>> {
  id: string;
  slotName: string;
  order?: number;
  component: ComponentType<P>;
}

export interface ModuleDefinition {
  id: string;
  name: string;
  description?: string;
  routes?: ModuleRoute[];
  menuItems?: ModuleMenuItem[];
  permissions?: string[];
  slots?: SlotContribution[];
}
