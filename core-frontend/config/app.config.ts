import type { ModuleDefinition } from '@/core/modules/types';
import { tenantModule } from '@/modules/tenant';
import { exampleModule } from '@/modules/example';
import { iamModule } from '@/modules/iam';
import { billingModule } from '@/modules/billing';
import { settingsModule } from '@/modules/settings';

export interface AppConfig {
  branding: {
    appName: string;
    tagline: string;
    logoUrl?: string;
    defaultTenant: string;
  };
  /**
   * Explicit list of enabled modules.
   * Add/remove module definition objects here to enable/disable features in the kit.
   */
  enabledModules: ModuleDefinition[];
}

export const appConfig: AppConfig = {
  branding: {
    appName: 'Core SaaS Starter Kit',
    tagline: 'Modular production-grade SaaS boilerplate',
    defaultTenant: 'default',
  },
  enabledModules: [
    tenantModule,
    exampleModule,
    iamModule,
    billingModule,
    settingsModule,
  ],
};
