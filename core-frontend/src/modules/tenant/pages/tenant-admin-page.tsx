'use client';

import React, { useState } from 'react';
import {
  Plus,
  Building2,
  CheckCircle2,
  Info,
  ShieldCheck,
  Server,
  Database,
  ArrowRight,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { useI18n } from '@/core/i18n/i18n-context';
import { TenantKpiCards } from '../components/tenant-kpi-cards';
import { TenantInspector } from '../components/tenant-inspector';
import { TenantListTable } from '../components/tenant-list-table';
import { CreateTenantDialog } from '../components/create-tenant-dialog';
import type { Tenant } from '../types/tenant';

export function TenantAdminPage() {
  const { t } = useI18n();
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [newlyCreatedTenant, setNewlyCreatedTenant] = useState<Tenant | null>(null);
  const [selectedTenantId, setSelectedTenantId] = useState<string | null>(null);

  const handleTenantCreated = (tenant: Tenant) => {
    setNewlyCreatedTenant(tenant);
    setSelectedTenantId(tenant.id);
  };

  const handleInspectTenant = (tenant: Tenant) => {
    setSelectedTenantId(tenant.id);
    const el = document.getElementById('tenant-inspector-section');
    if (el) {
      el.scrollIntoView({ behavior: 'smooth' });
    }
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between pb-5 border-b border-neutral-200 gap-4">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <div className="p-1.5 rounded-md border border-neutral-200 bg-white shadow-none">
              <Building2 className="h-5 w-5 text-neutral-900 stroke-[1.75]" />
            </div>
            <h1 className="text-xl font-bold tracking-tight text-neutral-900">
              {t('tenant.page.title', 'Tenant Management')}
            </h1>
          </div>
          <p className="text-xs text-neutral-500 max-w-2xl">
            {t(
              'tenant.page.description',
              'Provision, inspect, and switch active tenant workspaces. Built on top of the hardened Go Fiber backend with dynamic database routing and 3-layer typed error handling.'
            )}
          </p>
        </div>

        {/* Primary Action Button */}
        <Button
          onClick={() => setCreateDialogOpen(true)}
          className="gap-1.5 self-start sm:self-auto shrink-0"
        >
          <Plus className="h-4 w-4 stroke-[2]" />
          <span>{t('tenant.page.provision_btn', 'Provision Tenant')}</span>
        </Button>
      </div>

      {/* Success Notification Banner after Creation */}
      {newlyCreatedTenant && (
        <div className="rounded-lg border border-emerald-200 bg-emerald-50/60 p-4 text-xs text-emerald-900 flex items-start justify-between gap-3 duration-200 animate-in fade-in-0">
          <div className="flex items-start gap-2.5">
            <CheckCircle2 className="h-4 w-4 text-emerald-600 mt-0.5 shrink-0" />
            <div>
              <div className="font-semibold text-emerald-800">
                {t(
                  'tenant.page.created_success',
                  { name: newlyCreatedTenant.name },
                  `Tenant "${newlyCreatedTenant.name}" successfully created!`
                )}
              </div>
              <div className="text-emerald-700 font-mono text-[11px] mt-0.5">
                {t(
                  'tenant.page.created_details',
                  { id: newlyCreatedTenant.id, mode: newlyCreatedTenant.isolate_mode },
                  `ID: ${newlyCreatedTenant.id} • Mode: ${newlyCreatedTenant.isolate_mode}`
                )}
              </div>
            </div>
          </div>
          <Button
            size="sm"
            variant="outline"
            className="text-xs bg-white text-emerald-900 border-emerald-200 hover:bg-emerald-50"
            onClick={() => setNewlyCreatedTenant(null)}
          >
            {t('common.dismiss', 'Dismiss')}
          </Button>
        </div>
      )}

      {/* KPI Cards */}
      <TenantKpiCards />

      {/* Tenant Workspaces List Table */}
      <TenantListTable
        onInspectTenant={handleInspectTenant}
        onProvisionClick={() => setCreateDialogOpen(true)}
      />

      {/* Main Tenant Inspector Section */}
      <div id="tenant-inspector-section">
        <TenantInspector
          key={selectedTenantId || newlyCreatedTenant?.id || 'default'}
          initialTenantId={selectedTenantId || newlyCreatedTenant?.id}
        />
      </div>

      {/* Architecture Guidance Card (Outlined Minimalist) */}
      <Card>
        <CardHeader>
          <CardTitle className="text-sm flex items-center gap-2">
            <Info className="h-4 w-4 text-neutral-900 stroke-[1.75]" />
            <span>
              {t(
                'tenant.architecture.title',
                'Architecture & Data Isolation Guarantees'
              )}
            </span>
          </CardTitle>
          <CardDescription className="text-xs">
            {t(
              'tenant.architecture.description',
              'How the Core SaaS platform handles multi-tenancy and database routing behind the scenes.'
            )}
          </CardDescription>
        </CardHeader>
        <CardContent className="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs">
          <div className="rounded border border-neutral-200 p-3.5 bg-neutral-50/30 space-y-1.5">
            <div className="font-semibold text-neutral-900 flex items-center gap-1.5">
              <Database className="h-3.5 w-3.5 text-neutral-700" />
              <span>
                {t('tenant.architecture.shared_title', 'Shared Multi-Tenant (Default)')}
              </span>
            </div>
            <p className="text-neutral-600 leading-relaxed text-[11px]">
              {t(
                'tenant.architecture.shared_desc',
                'Tenants share a unified PostgreSQL database pool. Every database operation is automatically scoped with the active tenant ID extracted from HTTP header'
              )}{' '}
              <code className="text-neutral-900 font-mono text-[10px] bg-neutral-200/60 px-1 py-0.5 rounded">
                X-Tenant-ID
              </code>
              .
            </p>
          </div>

          <div className="rounded border border-neutral-200 p-3.5 bg-neutral-50/30 space-y-1.5">
            <div className="font-semibold text-neutral-900 flex items-center gap-1.5">
              <Server className="h-3.5 w-3.5 text-neutral-700" />
              <span>
                {t('tenant.architecture.dedicated_title', 'Dedicated Database Mode')}
              </span>
            </div>
            <p className="text-neutral-600 leading-relaxed text-[11px]">
              {t(
                'tenant.architecture.dedicated_desc',
                "The dynamic router resolves the tenant record and acquires an isolated connection pool using the tenant's custom DSN configuration, with automated idle pool reaping."
              )}
            </p>
          </div>
        </CardContent>
      </Card>

      {/* Create Tenant Modal */}
      <CreateTenantDialog
        open={createDialogOpen}
        onOpenChange={setCreateDialogOpen}
        onSuccess={handleTenantCreated}
      />
    </div>
  );
}
