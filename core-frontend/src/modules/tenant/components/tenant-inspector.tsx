'use client';

import React, { useState } from 'react';
import {
  Search,
  Copy,
  Check,
  Building,
  Database,
  Calendar,
  Clock,
  ArrowRight,
  AlertCircle,
  Eye,
  EyeOff,
  CheckCircle2,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useTenant as useTenantQuery } from '../hooks/tenant-hooks';
import { useTenant as useTenantContext } from '@/core/tenant/tenant-context';
import { useI18n } from '@/core/i18n/i18n-context';
import type { Tenant } from '../types/tenant';

interface TenantInspectorProps {
  initialTenantId?: string;
  onTenantCreated?: (tenant: Tenant) => void;
}

export function TenantInspector({ initialTenantId }: TenantInspectorProps) {
  const { formatError, t } = useI18n();
  const [searchId, setSearchId] = useState<string>(initialTenantId || '');
  const [activeQueryId, setActiveQueryId] = useState<string>(initialTenantId || '');
  const [copied, setCopied] = useState<boolean>(false);
  const [showDsn, setShowDsn] = useState<boolean>(false);
  const [switchedSuccess, setSwitchedSuccess] = useState<boolean>(false);

  const { tenant: currentContextTenant, setTenant: setContextTenant } = useTenantContext();

  const {
    data: tenant,
    isLoading,
    isError,
    error,
    refetch,
  } = useTenantQuery(activeQueryId, Boolean(activeQueryId));

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = searchId.trim();
    if (trimmed) {
      setActiveQueryId(trimmed);
    }
  };

  const handleCopyId = (id: string) => {
    navigator.clipboard.writeText(id);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleSwitchTenant = (target: Tenant) => {
    setContextTenant({
      id: target.id,
      name: target.name,
      isolateMode: target.isolate_mode,
      theme: {
        colorPrimary: '#0f172a',
        borderRadius: 8,
      },
    });

    if (typeof document !== 'undefined') {
      document.cookie = `x-tenant-id=${encodeURIComponent(target.id)}; path=/; max-age=31536000; SameSite=Lax`;
    }

    setSwitchedSuccess(true);
    setTimeout(() => setSwitchedSuccess(false), 3000);
  };

  const errorDetail = isError ? formatError(error) : null;

  return (
    <Card className="mb-6">
      <CardHeader>
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
          <div>
            <CardTitle className="text-base flex items-center gap-2">
              <Building className="h-4 w-4 text-neutral-900 stroke-[1.75]" />
              <span>{t('tenant.inspect.title', 'Tenant Lookup & Inspector')}</span>
            </CardTitle>
            <CardDescription className="mt-1">
              {t(
                'tenant.inspect.description',
                'Query and inspect tenant state directly from the Go backend via'
              )}{' '}
              <code className="text-neutral-900 font-mono text-xs bg-neutral-100 px-1 py-0.5 rounded border border-neutral-200">
                GET /tenants/:id
              </code>
            </CardDescription>
          </div>

          {currentContextTenant.id && (
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                setSearchId(currentContextTenant.id);
                setActiveQueryId(currentContextTenant.id);
              }}
              className="text-xs text-neutral-600 self-start sm:self-auto"
            >
              {t('tenant.inspect.inspect_active', 'Inspect Active Context')} (
              {currentContextTenant.id})
            </Button>
          )}
        </div>
      </CardHeader>

      <CardContent>
        {/* Search Bar */}
        <form onSubmit={handleSearch} className="flex gap-2 mb-6">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-neutral-600 stroke-[1.5]" />
            <Input
              value={searchId}
              onChange={(e) => setSearchId(e.target.value)}
              placeholder={t(
                'tenant.inspect.search_placeholder',
                'Enter Tenant UUID (e.g. 0192e21b-8526-70e2-a044-cf66a2985f47)...'
              )}
              className="pl-9 font-mono text-xs"
            />
          </div>
          <Button type="submit" disabled={!searchId.trim() || isLoading} size="default">
            {isLoading
              ? t('tenant.inspect.inspecting_btn', 'Inspecting...')
              : t('tenant.inspect.inspect_btn', 'Inspect')}
          </Button>
        </form>

        {/* Loading State */}
        {isLoading && (
          <div className="py-12 flex flex-col items-center justify-center text-neutral-600 gap-2">
            <div className="h-5 w-5 animate-spin rounded-full border-2 border-neutral-300 border-t-neutral-900" />
            <span className="text-xs">
              {t('tenant.inspect.fetching', 'Fetching tenant details from backend...')}
            </span>
          </div>
        )}

        {/* Error State */}
        {isError && errorDetail && (
          <div className="rounded-lg border border-red-200 bg-red-50/50 p-4 mb-4 text-xs text-red-900 flex items-start gap-3">
            <AlertCircle className="h-4 w-4 text-red-600 shrink-0 mt-0.5 stroke-[1.75]" />
            <div>
              <div className="font-semibold text-red-800">
                {errorDetail.code
                  ? `Error ${errorDetail.code}`
                  : t('tenant.inspect.failed', 'Inspection Failed')}
              </div>
              <div className="text-red-700 mt-0.5">{errorDetail.message}</div>
            </div>
          </div>
        )}

        {/* Result Detail View */}
        {tenant && !isLoading && (
          <div className="rounded-lg border border-neutral-200 bg-neutral-50/30 p-5 space-y-5">
            {/* Top row: Name, Status, Actions */}
            <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-4 border-b border-neutral-200 gap-3">
              <div>
                <div className="flex items-center gap-2.5">
                  <h3 className="text-lg font-bold text-neutral-900 tracking-tight">
                    {tenant.name}
                  </h3>
                  <Badge variant="outline" className="gap-1.5 bg-white">
                    <span className="h-1.5 w-1.5 rounded-full bg-emerald-500" />
                    <span className="capitalize">
                      {tenant.status || t('common.active', 'Active')}
                    </span>
                  </Badge>
                  <Badge
                    variant="outline"
                    className="uppercase text-[10px] tracking-wider font-semibold border-neutral-300 bg-white"
                  >
                    {tenant.isolate_mode || 'Shared'}
                  </Badge>
                </div>
                {tenant.description && (
                  <p className="text-xs text-neutral-500 mt-1">{tenant.description}</p>
                )}
              </div>

              {/* Set as Active Context Button */}
              <div className="flex items-center gap-2">
                {currentContextTenant.id === tenant.id ? (
                  <Badge variant="success" className="gap-1 text-xs py-1">
                    <CheckCircle2 className="h-3 w-3" />
                    <span>{t('tenant.inspect.active_context', 'Active Context')}</span>
                  </Badge>
                ) : (
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => handleSwitchTenant(tenant)}
                    className="text-xs gap-1.5 border-neutral-300 hover:border-neutral-900"
                  >
                    {switchedSuccess ? (
                      <>
                        <Check className="h-3.5 w-3.5 text-emerald-600" />
                        <span>{t('tenant.inspect.switched', 'Switched!')}</span>
                      </>
                    ) : (
                      <>
                        <span>{t('tenant.inspect.switch_context', 'Switch Context')}</span>
                        <ArrowRight className="h-3.5 w-3.5 stroke-[1.5]" />
                      </>
                    )}
                  </Button>
                )}
              </div>
            </div>

            {/* Metadata Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs">
              <div className="space-y-1">
                <span className="text-[10px] font-semibold uppercase tracking-wider text-neutral-600">
                  {t('tenant.inspect.id_label', 'Tenant Identifier (UUID)')}
                </span>
                <div className="flex items-center gap-2 bg-white px-2.5 py-1.5 rounded border border-neutral-200 font-mono text-[11px] text-neutral-800">
                  <span className="truncate">{tenant.id}</span>
                  <button
                    type="button"
                    onClick={() => handleCopyId(tenant.id)}
                    className="ml-auto p-1 hover:text-neutral-900 text-neutral-600 rounded transition-colors"
                    title={t('tenant.inspect.copy_uuid', 'Copy UUID')}
                  >
                    {copied ? (
                      <Check className="h-3 w-3 text-emerald-600" />
                    ) : (
                      <Copy className="h-3 w-3" />
                    )}
                  </button>
                </div>
              </div>

              <div className="space-y-1">
                <span className="text-[10px] font-semibold uppercase tracking-wider text-neutral-600">
                  {t('tenant.inspect.isolation_mode', 'Isolation Mode')}
                </span>
                <div className="bg-white px-2.5 py-1.5 rounded border border-neutral-200 text-neutral-800 flex items-center justify-between">
                  <span className="font-medium capitalize">{tenant.isolate_mode}</span>
                  <span className="text-[10px] text-neutral-600">
                    {tenant.isolate_mode === 'dedicated'
                      ? t('tenant.inspect.separate_db', 'Separate Database Connection')
                      : t('tenant.inspect.shared_db', 'Multi-tenant Shared DB')}
                  </span>
                </div>
              </div>

              {tenant.create_at && (
                <div className="space-y-1">
                  <span className="text-[10px] font-semibold uppercase tracking-wider text-neutral-600">
                    {t('tenant.inspect.created_at', 'Created At')}
                  </span>
                  <div className="bg-white px-2.5 py-1.5 rounded border border-neutral-200 text-neutral-700 flex items-center gap-2">
                    <Calendar className="h-3.5 w-3.5 text-neutral-600 stroke-[1.5]" />
                    <span>{new Date(tenant.create_at).toLocaleString()}</span>
                  </div>
                </div>
              )}

              {tenant.update_at && (
                <div className="space-y-1">
                  <span className="text-[10px] font-semibold uppercase tracking-wider text-neutral-600">
                    {t('tenant.inspect.updated_at', 'Last Updated')}
                  </span>
                  <div className="bg-white px-2.5 py-1.5 rounded border border-neutral-200 text-neutral-700 flex items-center gap-2">
                    <Clock className="h-3.5 w-3.5 text-neutral-600 stroke-[1.5]" />
                    <span>{new Date(tenant.update_at).toLocaleString()}</span>
                  </div>
                </div>
              )}
            </div>

            {/* Dedicated Database Parameters */}
            {tenant.isolate_mode === 'dedicated' && tenant.db_config && (
              <div className="mt-4 pt-4 border-t border-neutral-200">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-[11px] font-semibold uppercase tracking-wider text-neutral-700 flex items-center gap-1.5">
                    <Database className="h-3.5 w-3.5 text-neutral-900" />
                    <span>
                      {t('tenant.inspect.dedicated_config', 'Dedicated Database Configuration')}
                    </span>
                  </span>
                  <button
                    type="button"
                    onClick={() => setShowDsn(!showDsn)}
                    className="flex items-center gap-1 text-[11px] text-neutral-600 hover:text-neutral-900"
                  >
                    {showDsn ? <EyeOff className="h-3 w-3" /> : <Eye className="h-3 w-3" />}
                    <span>
                      {showDsn
                        ? t('tenant.inspect.mask_dsn', 'Mask DSN')
                        : t('tenant.inspect.reveal_dsn', 'Reveal DSN')}
                    </span>
                  </button>
                </div>

                <div className="rounded border border-neutral-200 bg-white p-3 font-mono text-xs space-y-2">
                  <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-1 text-neutral-800">
                    <span className="text-neutral-600 font-sans text-[11px]">
                      {t('tenant.inspect.connection_dsn', 'Connection DSN:')}
                    </span>
                    <span className="text-[11px] bg-neutral-50 px-2 py-0.5 rounded border border-neutral-200">
                      {showDsn
                        ? tenant.db_config.dsn
                        : tenant.db_config.dsn.replace(/:([^@]+)@/, ':••••••@')}
                    </span>
                  </div>
                  <div className="flex items-center justify-between text-neutral-800 pt-1 border-t border-neutral-100">
                    <span className="text-neutral-600 font-sans text-[11px]">
                      {t('tenant.inspect.max_conns', 'Max Connections:')}
                    </span>
                    <span className="font-medium">{tenant.db_config.max_conns || 10}</span>
                  </div>
                  {tenant.db_config.ssl_mode && (
                    <div className="flex items-center justify-between text-neutral-800 pt-1 border-t border-neutral-100">
                      <span className="text-neutral-600 font-sans text-[11px]">
                        {t('tenant.inspect.ssl_mode', 'SSL Mode:')}
                      </span>
                      <span className="font-medium">{tenant.db_config.ssl_mode}</span>
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
