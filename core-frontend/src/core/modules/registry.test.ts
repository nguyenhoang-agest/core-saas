import { describe, it, expect, beforeEach } from 'vitest';
import React from 'react';
import { ModuleRegistry } from './registry';
import type { ModuleDefinition } from './types';

describe('ModuleRegistry', () => {
  let registry: ModuleRegistry;

  const DummyComponent: React.FC = () => React.createElement('div', null, 'dummy');

  beforeEach(() => {
    registry = new ModuleRegistry();
  });

  it('registers a module and retrieves it by id', () => {
    const mod: ModuleDefinition = {
      id: 'billing',
      name: 'Billing Module',
    };
    registry.register(mod);

    expect(registry.get('billing')).toEqual(mod);
    expect(registry.getAll()).toHaveLength(1);
  });

  it('rejects duplicate module id with an error', () => {
    const mod1: ModuleDefinition = { id: 'iam', name: 'IAM v1' };
    const mod2: ModuleDefinition = { id: 'iam', name: 'IAM v2' };

    registry.register(mod1);
    expect(() => registry.register(mod2)).toThrowError(
      /Duplicate module id detected: "iam"/
    );
  });

  it('rejects duplicate slot contribution id with an error', () => {
    const mod1: ModuleDefinition = {
      id: 'mod1',
      name: 'Module 1',
      slots: [
        {
          id: 'user-profile-menu-item',
          slotName: 'header-user-menu',
          component: DummyComponent,
        },
      ],
    };
    const mod2: ModuleDefinition = {
      id: 'mod2',
      name: 'Module 2',
      slots: [
        {
          id: 'user-profile-menu-item',
          slotName: 'header-user-menu',
          component: DummyComponent,
        },
      ],
    };

    registry.register(mod1);
    expect(() => registry.register(mod2)).toThrowError(
      /Duplicate slot contribution id detected: "user-profile-menu-item"/
    );
  });

  it('orders slot contributions strictly by order ascending', () => {
    const mod1: ModuleDefinition = {
      id: 'mod1',
      name: 'Module 1',
      slots: [
        {
          id: 'slot-c',
          slotName: 'header-actions',
          order: 30,
          component: DummyComponent,
        },
        {
          id: 'slot-a',
          slotName: 'header-actions',
          order: 10,
          component: DummyComponent,
        },
      ],
    };
    const mod2: ModuleDefinition = {
      id: 'mod2',
      name: 'Module 2',
      slots: [
        {
          id: 'slot-b',
          slotName: 'header-actions',
          order: 20,
          component: DummyComponent,
        },
      ],
    };

    registry.register(mod1);
    registry.register(mod2);

    const contributions = registry.getSlotContributions('header-actions');
    expect(contributions.map((c) => c.id)).toEqual(['slot-a', 'slot-b', 'slot-c']);
  });

  it('orders menu items by order ascending', () => {
    const mod1: ModuleDefinition = {
      id: 'mod1',
      name: 'Module 1',
      menuItems: [
        { key: 'settings', label: 'Settings', path: '/settings', order: 50 },
        { key: 'dashboard', label: 'Dashboard', path: '/', order: 10 },
      ],
    };
    const mod2: ModuleDefinition = {
      id: 'mod2',
      name: 'Module 2',
      menuItems: [
        { key: 'iam', label: 'IAM', path: '/iam', order: 20 },
      ],
    };

    registry.register(mod1);
    registry.register(mod2);

    const menu = registry.getMenuItems();
    expect(menu.map((m) => m.key)).toEqual(['dashboard', 'iam', 'settings']);
  });
});
