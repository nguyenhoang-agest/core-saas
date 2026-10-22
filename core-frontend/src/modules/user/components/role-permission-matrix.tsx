'use client';

import React, { useMemo } from 'react';
import { Key, Check, Minus, Lock } from 'lucide-react';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { useI18n } from '@/core/i18n/i18n-context';
import { usePermissionsCatalog, useRoles } from '../hooks';
import type { RoleItem, PermissionCatalogItem } from '../types';

export function RolePermissionMatrix() {
  const { t } = useI18n();
  const { data: catalogData = [], isLoading: catalogLoading } = usePermissionsCatalog();
  const { data: rolesData = [], isLoading: rolesLoading } = useRoles();
  const catalog = Array.isArray(catalogData) ? catalogData : [];
  const roles = Array.isArray(rolesData) ? rolesData : [];

  // Check if a role default/assigned permissions contain the given perm
  const hasPermission = (role: RoleItem, permKey: string) => {
    if (role.slug === 'owner') return true;
    return (role.permissions || []).includes(permKey) || (role.permissions || []).includes('*');
  };

  // Group permissions by module
  const groupedPerms = useMemo(() => {
    const map = new Map<string, PermissionCatalogItem[]>();
    for (const p of catalog) {
      const mod = p.module || 'general';
      if (!map.has(mod)) {
        map.set(mod, []);
      }
      map.get(mod)!.push(p);
    }
    return Array.from(map.entries());
  }, [catalog]);

  const isLoading = catalogLoading || rolesLoading;

  return (
    <Card className="shadow-none border border-neutral-200">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-sm font-semibold text-neutral-900 flex items-center gap-2">
              <Key className="h-4 w-4 text-indigo-600" />
              <span>{t('iam.matrix.title', 'Role & Permission Access Matrix')}</span>
            </CardTitle>
            <CardDescription className="text-xs text-neutral-500 mt-1">
              {t(
                'iam.matrix.subtitle',
                'Comprehensive comparative grid of capabilities across all system and custom tenant roles.'
              )}
            </CardDescription>
          </div>
        </div>
      </CardHeader>

      <CardContent className="p-0">
        <div className="overflow-x-auto">
          <table className="w-full text-xs text-left border-collapse min-w-[700px]">
            <thead>
              <tr className="border-y border-neutral-200 bg-neutral-50 text-neutral-600 font-medium">
                <th className="py-2.5 px-4 w-1/3 min-w-[240px]">
                  {t('iam.matrix.col_capability', 'Capability & Module')}
                </th>
                {roles.map((r) => (
                  <th key={r.id} className="py-2.5 px-3 text-center min-w-[120px]">
                    <div className="font-semibold text-neutral-900">{r.name}</div>
                    <div className="mt-0.5">
                      {r.is_system ? (
                        <Badge
                          variant="secondary"
                          className="text-[10px] py-0 px-1.5 h-4 bg-neutral-100 text-neutral-600"
                        >
                          <Lock className="h-2.5 w-2.5 mr-0.5 inline" />
                          {t('iam.roles.badge_system', 'System')}
                        </Badge>
                      ) : (
                        <Badge
                          variant="outline"
                          className="text-[10px] py-0 px-1.5 h-4 bg-indigo-50 text-indigo-700 border-indigo-200"
                        >
                          {t('iam.roles.badge_custom', 'Custom')}
                        </Badge>
                      )}
                    </div>
                  </th>
                ))}
              </tr>
            </thead>

            <tbody className="divide-y divide-neutral-200">
              {isLoading ? (
                <tr>
                  <td colSpan={roles.length + 1} className="py-8 text-center text-neutral-500">
                    {t('common.loading', 'Loading access matrix...')}
                  </td>
                </tr>
              ) : groupedPerms.length === 0 ? (
                <tr>
                  <td colSpan={roles.length + 1} className="py-8 text-center text-neutral-500">
                    {t('iam.matrix.no_capabilities', 'No capabilities registered in catalog.')}
                  </td>
                </tr>
              ) : (
                groupedPerms.map(([moduleName, perms]) => (
                  <React.Fragment key={moduleName}>
                    <tr className="bg-neutral-100/60 font-semibold text-neutral-700">
                      <td colSpan={roles.length + 1} className="py-1.5 px-4 uppercase text-[11px] tracking-wider">
                        {moduleName}
                      </td>
                    </tr>
                    {perms.map((p) => (
                      <tr key={p.key} className="hover:bg-neutral-50/70 transition-colors">
                        <td className="py-2.5 px-4">
                          <div className="font-medium text-neutral-900">
                            {(p as any).name || t(p.i18n_name_key, p.key)}
                          </div>
                          <div className="text-[11px] text-neutral-500 font-mono">
                            {p.key}
                          </div>
                        </td>
                        {roles.map((r) => {
                          const active = hasPermission(r, p.key);
                          return (
                            <td key={r.id} className="py-2.5 px-3 text-center">
                              {active ? (
                                <span className="inline-flex items-center justify-center h-5 w-5 rounded-full bg-emerald-100 text-emerald-700">
                                  <Check className="h-3 w-3 stroke-[2.5]" />
                                </span>
                              ) : (
                                <span className="inline-flex items-center justify-center h-5 w-5 rounded-full bg-neutral-100 text-neutral-400">
                                  <Minus className="h-3 w-3" />
                                </span>
                              )}
                            </td>
                          );
                        })}
                      </tr>
                    ))}
                  </React.Fragment>
                ))
              )}
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  );
}
