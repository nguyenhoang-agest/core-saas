'use client';

import React from 'react';
import { Users, UserCheck, Shield, Key } from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { useI18n } from '@/core/i18n/i18n-context';
import type { User } from '../types';

interface UserKpiCardsProps {
  users?: User[];
  totalCount?: number;
  isLoading?: boolean;
}

export function UserKpiCards({ users = [], totalCount = 0, isLoading = false }: UserKpiCardsProps) {
  const { t } = useI18n();

  const activeCount = users.filter((u) => u.status === 'active').length;
  const adminCount = users.filter((u) => u.role === 'admin' || u.role === 'owner').length;
  const aclCount = users.filter((u) => u.permissions && u.permissions.length > 0).length;

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {/* Total Users */}
      <Card className="shadow-none">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-xs font-medium text-neutral-500">
            {t('iam.kpi.total_users', 'Total Members')}
          </CardTitle>
          <div className="p-1.5 rounded-md border border-neutral-200 bg-neutral-50">
            <Users className="h-4 w-4 text-neutral-700" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold tracking-tight text-neutral-900">
            {isLoading ? '...' : (totalCount || users.length)}
          </div>
          <p className="text-[11px] text-neutral-500 mt-1">
            {t('iam.kpi.total_users_sub', 'Active & invited directory')}
          </p>
        </CardContent>
      </Card>

      {/* Active Users */}
      <Card className="shadow-none">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-xs font-medium text-neutral-500">
            {t('iam.kpi.active_users', 'Active Users')}
          </CardTitle>
          <div className="p-1.5 rounded-md border border-neutral-200 bg-emerald-50">
            <UserCheck className="h-4 w-4 text-emerald-600" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold tracking-tight text-neutral-900">
            {isLoading ? '...' : activeCount}
          </div>
          <p className="text-[11px] text-neutral-500 mt-1">
            {t('iam.kpi.active_users_sub', 'Operational members')}
          </p>
        </CardContent>
      </Card>

      {/* Admin / Owners */}
      <Card className="shadow-none">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-xs font-medium text-neutral-500">
            {t('iam.kpi.admin_users', 'Administrators')}
          </CardTitle>
          <div className="p-1.5 rounded-md border border-neutral-200 bg-indigo-50">
            <Shield className="h-4 w-4 text-indigo-600" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold tracking-tight text-neutral-900">
            {isLoading ? '...' : adminCount}
          </div>
          <p className="text-[11px] text-neutral-500 mt-1">
            {t('iam.kpi.admin_users_sub', 'Owner & Admin roles')}
          </p>
        </CardContent>
      </Card>

      {/* Direct ACL Overrides */}
      <Card className="shadow-none">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-xs font-medium text-neutral-500">
            {t('iam.kpi.custom_acl', 'ACL Overrides')}
          </CardTitle>
          <div className="p-1.5 rounded-md border border-neutral-200 bg-amber-50">
            <Key className="h-4 w-4 text-amber-600" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold tracking-tight text-neutral-900">
            {isLoading ? '...' : aclCount}
          </div>
          <p className="text-[11px] text-neutral-500 mt-1">
            {t('iam.kpi.custom_acl_sub', 'Users with direct permissions')}
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
