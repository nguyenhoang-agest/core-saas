import React from 'react';
import { Database, ShieldCheck, Server, Sparkles } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';
import { useTenant } from '@/core/tenant/tenant-context';
import { useI18n } from '@/core/i18n/i18n-context';

export function TenantKpiCards() {
  const { tenant } = useTenant();
  const { t } = useI18n();

  return (
    <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
      {/* KPI 1: Active Tenant Context */}
      <Card className="hover:border-neutral-300 transition-colors">
        <CardContent className="p-4">
          <div className="flex items-center justify-between mb-2">
            <span className="text-[11px] font-semibold uppercase tracking-wider text-neutral-600">
              {t('tenant.kpi.active_context', 'Active Context')}
            </span>
            <span className="flex h-2 w-2 rounded-full bg-emerald-500 ring-2 ring-emerald-100" />
          </div>
          <div className="text-base font-bold text-neutral-900 truncate">
            {tenant.name || t('tenant.kpi.default_org', 'Default Organization')}
          </div>
          <div className="text-[11px] font-mono text-neutral-600 mt-1 truncate">
            id: {tenant.id}
          </div>
        </CardContent>
      </Card>

      {/* KPI 2: Isolation Modes */}
      <Card className="hover:border-neutral-300 transition-colors">
        <CardContent className="p-4">
          <div className="flex items-center justify-between mb-2">
            <span className="text-[11px] font-semibold uppercase tracking-wider text-neutral-600">
              {t('tenant.kpi.supported_isolation', 'Supported Isolation')}
            </span>
            <Database className="h-4 w-4 text-neutral-600 stroke-[1.5]" />
          </div>
          <div className="text-base font-bold text-neutral-900">
            {t('tenant.kpi.shared_dedicated', 'Shared & Dedicated')}
          </div>
          <div className="text-[11px] text-neutral-600 mt-1">
            {t('tenant.kpi.dynamic_routing', 'Dynamic connection pool routing')}
          </div>
        </CardContent>
      </Card>

      {/* KPI 3: Error Architecture */}
      <Card className="hover:border-neutral-300 transition-colors">
        <CardContent className="p-4">
          <div className="flex items-center justify-between mb-2">
            <span className="text-[11px] font-semibold uppercase tracking-wider text-neutral-600">
              {t('tenant.kpi.api_architecture', 'API Architecture')}
            </span>
            <ShieldCheck className="h-4 w-4 text-neutral-600 stroke-[1.5]" />
          </div>
          <div className="text-base font-bold text-neutral-900">
            {t('tenant.kpi.three_layer_errors', '3-Layer Errors')}
          </div>
          <div className="text-[11px] text-neutral-600 mt-1">
            {t('tenant.kpi.type_safe_codes', 'Type-safe domain codes (E031001+)')}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
