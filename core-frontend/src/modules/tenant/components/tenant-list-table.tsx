'use client';

import React, { useState, useMemo } from 'react';
import {
  Search,
  Building,
  Database,
  Server,
  Copy,
  Check,
  CheckCircle2,
  ArrowRight,
  ChevronLeft,
  ChevronRight,
  AlertCircle,
  RefreshCw,
  Plus,
  Filter,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useTenants } from '../api/tenant-api';
import { useTenant as useTenantContext } from '@/core/tenant/tenant-context';
import { useI18n } from '@/core/i18n/i18n-context';
import type { Tenant, TenantIsolateMode } from '../types/tenant';

interface TenantListTableProps {
  onInspectTenant?: (tenant: Tenant) => void;
  onProvisionClick?: () => void;
}

export function TenantListTable({
  onInspectTenant,
  onProvisionClick,
}: TenantListTableProps) {
  const { t } = useI18n();
  const { tenant: currentContextTenant, setTenant: setContextTenant } = useTenantContext();

  // Search & Filter state
  const [searchInput, setSearchInput] = useState('');
  const [debouncedKeyword, setDebouncedKeyword] = useState('');
  const [statusFilter, setStatusFilter] = useState<'all' | 'active' | 'suspended'>('all');
  const [isolationFilter, setIsolationFilter] = useState<'all' | 'shared' | 'dedicated'>('all');

  // Pagination state
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  // Interaction feedback states
  const [copiedId, setCopiedId] = useState<string | null>(null);
  const [switchedId, setSwitchedId] = useState<string | null>(null);

  // Fetch tenants with query
  const queryParams = useMemo(() => {
    return {
      page,
      pageSize,
      keyword: debouncedKeyword,
      status: statusFilter !== 'all' ? statusFilter : undefined,
      isolateMode: isolationFilter !== 'all' ? isolationFilter : undefined,
    };
  }, [page, pageSize, debouncedKeyword, statusFilter, isolationFilter]);

  const { data, isLoading, isError, error, refetch, isFetching } = useTenants(queryParams);

  // Apply client-side isolation filter fallback
  const filteredItems = useMemo(() => {
    const rawItems = data?.items || [];
    if (isolationFilter === 'all') return rawItems;
    return rawItems.filter((item) => (item.isolate_mode || 'shared') === isolationFilter);
  }, [data?.items, isolationFilter]);

  const pagination = data?.pagination;
  const totalItems = Number(pagination?.totalItems ?? pagination?.total_items ?? 0);
  const totalPages = Math.max(
    1,
    Number(
      pagination?.totalPages ??
      pagination?.total_pages ??
      (totalItems > 0 ? Math.ceil(totalItems / pageSize) : 1)
    )
  );

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setDebouncedKeyword(searchInput.trim());
    setPage(1);
  };

  const handleClearSearch = () => {
    setSearchInput('');
    setDebouncedKeyword('');
    setPage(1);
  };

  const handleCopyId = (id: string) => {
    navigator.clipboard.writeText(id);
    setCopiedId(id);
    setTimeout(() => setCopiedId(null), 2000);
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

    setSwitchedId(target.id);
    setTimeout(() => setSwitchedId(null), 3000);
  };

  const fromIndex = totalItems === 0 ? 0 : (page - 1) * pageSize + 1;
  const toIndex = Math.min(page * pageSize, totalItems);

  return (
    <Card className="shadow-none">
      <CardHeader className="pb-4">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
          <div>
            <CardTitle className="text-base flex items-center gap-2">
              <Building className="h-4 w-4 text-neutral-900 stroke-[1.75]" />
              <span>{t('tenant.list.title', 'Tenant Workspaces')}</span>
            </CardTitle>
            <CardDescription className="mt-1">
              {t(
                'tenant.list.subtitle',
                'Overview of registered tenants, isolation configurations, and status across the platform.'
              )}
            </CardDescription>
          </div>

          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => refetch()}
              disabled={isFetching}
              className="text-xs gap-1.5 h-8 text-neutral-700"
            >
              <RefreshCw className={`h-3.5 w-3.5 ${isFetching ? 'animate-spin' : ''}`} />
              <span>{isFetching ? t('common.loading', 'Loading...') : t('common.refresh', 'Refresh')}</span>
            </Button>
            {onProvisionClick && (
              <Button size="sm" onClick={onProvisionClick} className="text-xs gap-1.5 h-8">
                <Plus className="h-3.5 w-3.5" />
                <span>{t('tenant.page.provision_btn', 'Provision Tenant')}</span>
              </Button>
            )}
          </div>
        </div>
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Search & Filters Bar */}
        <div className="flex flex-col md:flex-row items-stretch md:items-center justify-between gap-3">
          <form onSubmit={handleSearchSubmit} className="relative flex-1">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-neutral-500 stroke-[1.5]" />
            <Input
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              placeholder={t('tenant.list.search_placeholder', 'Search tenants by name...')}
              className="pl-9 pr-14 text-xs h-9 font-sans"
            />
            {searchInput && (
              <button
                type="button"
                onClick={handleClearSearch}
                className="absolute right-2.5 top-2 text-xs text-neutral-400 hover:text-neutral-700 px-1 py-0.5 rounded"
              >
                ✕
              </button>
            )}
          </form>

          <div className="flex items-center gap-2 shrink-0">
            {/* Status Filter */}
            <div className="flex items-center gap-1.5">
              <Filter className="h-3.5 w-3.5 text-neutral-500 hidden sm:inline" />
              <select
                value={statusFilter}
                onChange={(e) => {
                  setStatusFilter(e.target.value as 'all' | 'active' | 'suspended');
                  setPage(1);
                }}
                className="text-xs rounded-md border border-neutral-200 bg-white px-2.5 py-1.5 text-neutral-800 focus:outline-none focus:ring-1 focus:ring-neutral-900"
                aria-label={t('tenant.list.filter_status', 'Status')}
              >
                <option value="all">{t('tenant.list.all_statuses', 'All Statuses')}</option>
                <option value="active">{t('tenant.list.active', 'Active')}</option>
                <option value="suspended">{t('tenant.list.suspended', 'Suspended')}</option>
              </select>
            </div>

            {/* Isolation Mode Filter */}
            <div>
              <select
                value={isolationFilter}
                onChange={(e) => {
                  setIsolationFilter(e.target.value as 'all' | 'shared' | 'dedicated');
                  setPage(1);
                }}
                className="text-xs rounded-md border border-neutral-200 bg-white px-2.5 py-1.5 text-neutral-800 focus:outline-none focus:ring-1 focus:ring-neutral-900"
                aria-label={t('tenant.list.filter_isolation', 'Isolation')}
              >
                <option value="all">{t('tenant.list.all_isolation', 'All Isolation Modes')}</option>
                <option value="shared">{t('tenant.list.shared', 'Shared')}</option>
                <option value="dedicated">{t('tenant.list.dedicated', 'Dedicated')}</option>
              </select>
            </div>
          </div>
        </div>

        {/* Error Alert */}
        {isError && (
          <div className="rounded-lg border border-red-200 bg-red-50/50 p-4 text-xs text-red-900 flex items-start justify-between gap-3">
            <div className="flex items-start gap-2.5">
              <AlertCircle className="h-4 w-4 text-red-600 mt-0.5 shrink-0" />
              <div>
                <div className="font-semibold">{t('tenant.inspect.failed', 'Failed to load tenants')}</div>
                <div className="text-red-700 mt-0.5">
                  {(error as Error)?.message || 'An error occurred while communicating with the Go backend.'}
                </div>
              </div>
            </div>
            <Button
              size="sm"
              variant="outline"
              className="text-xs bg-white text-red-900 border-red-200 hover:bg-red-50"
              onClick={() => refetch()}
            >
              {t('common.retry', 'Retry')}
            </Button>
          </div>
        )}

        {/* Table Content */}
        <div className="rounded-md border border-neutral-200 overflow-x-auto bg-white">
          <table className="w-full text-left border-collapse text-xs">
            <thead>
              <tr className="border-b border-neutral-200 bg-neutral-50/75 text-[11px] font-semibold text-neutral-600 uppercase tracking-wider">
                <th className="py-3 px-4">{t('tenant.list.col_tenant', 'Tenant')}</th>
                <th className="py-3 px-4">{t('tenant.list.col_id', 'Identifier')}</th>
                <th className="py-3 px-4">{t('tenant.list.col_isolation', 'Data Isolation')}</th>
                <th className="py-3 px-4">{t('tenant.list.col_status', 'Status')}</th>
                <th className="py-3 px-4">{t('tenant.list.col_created', 'Created At')}</th>
                <th className="py-3 px-4 text-right">{t('tenant.list.col_actions', 'Actions')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-neutral-200">
              {/* Loading State */}
              {isLoading && (
                <>
                  {[1, 2, 3, 4].map((idx) => (
                    <tr key={idx} className="animate-pulse">
                      <td className="py-3 px-4">
                        <div className="h-3.5 bg-neutral-200 rounded w-32 mb-1.5" />
                        <div className="h-2.5 bg-neutral-100 rounded w-48" />
                      </td>
                      <td className="py-3 px-4">
                        <div className="h-3 bg-neutral-100 rounded w-24" />
                      </td>
                      <td className="py-3 px-4">
                        <div className="h-4 bg-neutral-100 rounded w-16" />
                      </td>
                      <td className="py-3 px-4">
                        <div className="h-4 bg-neutral-100 rounded w-14" />
                      </td>
                      <td className="py-3 px-4">
                        <div className="h-3 bg-neutral-100 rounded w-20" />
                      </td>
                      <td className="py-3 px-4 text-right">
                        <div className="h-7 bg-neutral-100 rounded w-24 ml-auto" />
                      </td>
                    </tr>
                  ))}
                </>
              )}

              {/* Empty State */}
              {!isLoading && filteredItems.length === 0 && (
                <tr>
                  <td colSpan={6} className="py-12 text-center text-neutral-500">
                    <div className="flex flex-col items-center justify-center max-w-sm mx-auto">
                      <div className="p-2.5 rounded-full bg-neutral-100 text-neutral-600 mb-2">
                        <Building className="h-6 w-6 stroke-[1.5]" />
                      </div>
                      <div className="font-semibold text-neutral-800 text-sm">
                        {t('tenant.list.empty_title', 'No tenants found')}
                      </div>
                      <p className="text-xs text-neutral-500 mt-1 mb-4 leading-relaxed">
                        {debouncedKeyword || statusFilter !== 'all' || isolationFilter !== 'all'
                          ? t(
                              'tenant.list.empty_search',
                              'No tenants match your search and filter criteria. Try adjusting your query.'
                            )
                          : t(
                              'tenant.list.empty_initial',
                              'No tenants provisioned yet. Click below to create your first tenant.'
                            )}
                      </p>
                      {debouncedKeyword || statusFilter !== 'all' || isolationFilter !== 'all' ? (
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => {
                            handleClearSearch();
                            setStatusFilter('all');
                            setIsolationFilter('all');
                          }}
                          className="text-xs"
                        >
                          {t('tenant.list.reset_filters', 'Reset Filters')}
                        </Button>
                      ) : (
                        onProvisionClick && (
                          <Button size="sm" onClick={onProvisionClick} className="text-xs gap-1.5">
                            <Plus className="h-3.5 w-3.5" />
                            <span>{t('tenant.page.provision_btn', 'Provision Tenant')}</span>
                          </Button>
                        )
                      )}
                    </div>
                  </td>
                </tr>
              )}

              {/* Rows */}
              {!isLoading &&
                filteredItems.map((item) => {
                  const isCurrentActive = currentContextTenant.id === item.id;
                  const isJustSwitched = switchedId === item.id;
                  const isDedicated = item.isolate_mode === 'dedicated';

                  return (
                    <tr
                      key={item.id}
                      className={`hover:bg-neutral-50/60 transition-colors ${
                        isCurrentActive ? 'bg-neutral-50/80 font-medium' : ''
                      }`}
                    >
                      {/* Name & Description */}
                      <td className="py-3 px-4">
                        <div className="flex items-start gap-2.5">
                          <div
                            className={`p-1.5 rounded border mt-0.5 shrink-0 ${
                              isCurrentActive
                                ? 'bg-neutral-900 border-neutral-900 text-white'
                                : 'bg-neutral-50 border-neutral-200 text-neutral-600'
                            }`}
                          >
                            <Building className="h-3.5 w-3.5" />
                          </div>
                          <div>
                            <div className="font-semibold text-neutral-900 flex items-center gap-1.5">
                              <span>{item.name}</span>
                              {isCurrentActive && (
                                <span className="inline-block px-1.5 py-0.2 rounded text-[10px] bg-neutral-900 text-white font-normal">
                                  Current
                                </span>
                              )}
                            </div>
                            {item.description && (
                              <p className="text-[11px] text-neutral-500 line-clamp-1 mt-0.5">
                                {item.description}
                              </p>
                            )}
                          </div>
                        </div>
                      </td>

                      {/* UUID with Copy */}
                      <td className="py-3 px-4">
                        <div className="flex items-center gap-1.5 font-mono text-[11px] text-neutral-700">
                          <span className="truncate max-w-[110px]" title={item.id}>
                            {item.id}
                          </span>
                          <button
                            type="button"
                            onClick={() => handleCopyId(item.id)}
                            className="p-1 hover:text-neutral-900 text-neutral-400 rounded transition-colors"
                            title={t('tenant.inspect.copy_uuid', 'Copy UUID')}
                          >
                            {copiedId === item.id ? (
                              <Check className="h-3 w-3 text-emerald-600" />
                            ) : (
                              <Copy className="h-3 w-3" />
                            )}
                          </button>
                        </div>
                      </td>

                      {/* Isolation Mode */}
                      <td className="py-3 px-4">
                        <Badge
                          variant="outline"
                          className={`gap-1.5 text-[11px] font-normal ${
                            isDedicated
                              ? 'border-indigo-200 bg-indigo-50/70 text-indigo-800'
                              : 'border-neutral-200 bg-neutral-50 text-neutral-700'
                          }`}
                        >
                          {isDedicated ? (
                            <Server className="h-3 w-3 text-indigo-600" />
                          ) : (
                            <Database className="h-3 w-3 text-neutral-500" />
                          )}
                          <span className="capitalize">{item.isolate_mode || 'Shared'}</span>
                        </Badge>
                      </td>

                      {/* Status */}
                      <td className="py-3 px-4">
                        <Badge
                          variant="outline"
                          className="gap-1.5 text-[11px] font-normal border-neutral-200 bg-white"
                        >
                          <span
                            className={`h-1.5 w-1.5 rounded-full ${
                              item.status === 'suspended' ? 'bg-amber-500' : 'bg-emerald-500'
                            }`}
                          />
                          <span className="capitalize">{item.status || 'Active'}</span>
                        </Badge>
                      </td>

                      {/* Created At */}
                      <td className="py-3 px-4 text-neutral-500 text-[11px]">
                        {item.create_at ? new Date(item.create_at).toLocaleDateString() : '—'}
                      </td>

                      {/* Actions */}
                      <td className="py-3 px-4 text-right">
                        <div className="flex items-center justify-end gap-1.5">
                          {onInspectTenant && (
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => onInspectTenant(item)}
                              className="text-xs h-7 px-2 text-neutral-700 border-neutral-200 hover:border-neutral-900"
                            >
                              {t('tenant.list.inspect', 'Inspect')}
                            </Button>
                          )}

                          {isCurrentActive ? (
                            <Badge variant="success" className="gap-1 text-[11px] py-1 px-2 h-7">
                              <CheckCircle2 className="h-3 w-3" />
                              <span>{t('tenant.list.active_context', 'Active Context')}</span>
                            </Badge>
                          ) : (
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => handleSwitchTenant(item)}
                              className="text-xs h-7 px-2 gap-1 text-neutral-800 border-neutral-200 hover:border-neutral-900"
                            >
                              {isJustSwitched ? (
                                <>
                                  <Check className="h-3 w-3 text-emerald-600" />
                                  <span>{t('common.switched', 'Switched!')}</span>
                                </>
                              ) : (
                                <>
                                  <span>{t('tenant.list.switch_context', 'Switch')}</span>
                                  <ArrowRight className="h-3 w-3 stroke-[1.5]" />
                                </>
                              )}
                            </Button>
                          )}
                        </div>
                      </td>
                    </tr>
                  );
                })}
            </tbody>
          </table>
        </div>

        {/* Server-Side Pagination Bar */}
        <div className="flex flex-col sm:flex-row items-center justify-between gap-3 pt-2 text-xs text-neutral-600">
          <div className="flex items-center gap-3">
            <span>
              {t(
                'tenant.list.showing_results',
                { from: fromIndex, to: toIndex, total: totalItems },
                `Showing ${fromIndex} - ${toIndex} of ${totalItems} tenants`
              )}
            </span>

            {/* Page Size Selector */}
            <div className="flex items-center gap-1.5">
              <span className="text-[11px] text-neutral-400">|</span>
              <span className="text-[11px]">{t('tenant.list.page_size', 'Per page')}:</span>
              <select
                value={pageSize}
                onChange={(e) => {
                  setPageSize(Number(e.target.value));
                  setPage(1);
                }}
                className="text-xs rounded border border-neutral-200 bg-white px-2 py-0.5 text-neutral-800 focus:outline-none"
              >
                <option value={10}>10</option>
                <option value={20}>20</option>
                <option value={50}>50</option>
              </select>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-[11px] text-neutral-500 mr-1">
              {t('tenant.list.page', { current: page, total: totalPages }, `Page ${page} of ${totalPages}`)}
            </span>
            <Button
              variant="outline"
              size="sm"
              disabled={page <= 1 || isLoading}
              onClick={() => setPage((prev) => Math.max(1, prev - 1))}
              className="h-7 px-2 text-xs gap-1 border-neutral-200"
            >
              <ChevronLeft className="h-3.5 w-3.5" />
              <span className="hidden sm:inline">{t('tenant.list.prev', 'Previous')}</span>
            </Button>
            <Button
              variant="outline"
              size="sm"
              disabled={page >= totalPages || isLoading}
              onClick={() => setPage((prev) => Math.min(totalPages, prev + 1))}
              className="h-7 px-2 text-xs gap-1 border-neutral-200"
            >
              <span className="hidden sm:inline">{t('tenant.list.next', 'Next')}</span>
              <ChevronRight className="h-3.5 w-3.5" />
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
