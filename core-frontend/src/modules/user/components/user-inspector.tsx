'use client';

import React, { useState } from 'react';
import {
  User as UserIcon,
  Copy,
  Check,
  Shield,
  Key,
  Calendar,
  X,
  Mail,
  Building,
  CheckCircle2,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useI18n } from '@/core/i18n/i18n-context';
import type { User } from '../types';

interface UserInspectorProps {
  user: User | null;
  onClose?: () => void;
}

export function UserInspector({ user, onClose }: UserInspectorProps) {
  const { t } = useI18n();
  const [copiedId, setCopiedId] = useState(false);

  if (!user) return null;

  const handleCopyId = () => {
    navigator.clipboard.writeText(user.id);
    setCopiedId(true);
    setTimeout(() => setCopiedId(false), 2000);
  };

  const isOwner = user.role === 'owner';
  const customAcl = user.permissions || [];
  const effectivePerms = user.effective_permissions || [];

  return (
    <Card className="border border-neutral-300 shadow-sm transition-all duration-200">
      <CardHeader className="flex flex-row items-start justify-between pb-3">
        <div className="flex items-center gap-3">
          <div className="h-10 w-10 rounded-full border border-neutral-300 bg-neutral-100 flex items-center justify-center text-neutral-800 font-semibold text-sm">
            {user.avatar_url ? (
              // eslint-disable-next-line @next/next/no-img-element
              <img
                src={user.avatar_url}
                alt={user.name}
                className="h-full w-full rounded-full object-cover"
              />
            ) : (
              user.name.slice(0, 2).toUpperCase()
            )}
          </div>
          <div>
            <div className="flex items-center gap-2">
              <CardTitle className="text-base font-semibold text-neutral-900">
                {user.name}
              </CardTitle>
              <Badge
                variant={
                  user.status === 'active'
                    ? 'success'
                    : user.status === 'invited'
                    ? 'outline'
                    : 'secondary'
                }
                className="text-[10px] uppercase font-mono tracking-wider"
              >
                {t(`iam.statuses.${user.status}`, user.status)}
              </Badge>
            </div>
            <CardDescription className="text-xs text-neutral-500 flex items-center gap-1.5 mt-0.5">
              <Mail className="h-3.5 w-3.5 text-neutral-400" />
              <span>{user.email}</span>
            </CardDescription>
          </div>
        </div>

        {onClose && (
          <Button
            size="sm"
            variant="ghost"
            onClick={onClose}
            className="h-8 w-8 p-0 text-neutral-400 hover:text-neutral-700"
            aria-label="Close user inspector"
          >
            <X className="h-4 w-4" />
          </Button>
        )}
      </CardHeader>

      <CardContent className="space-y-4 pt-1">
        {/* Core Metadata Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-3 bg-neutral-50 p-3 rounded-md border border-neutral-200 text-xs">
          <div>
            <span className="text-[10px] text-neutral-500 font-medium uppercase tracking-wider block">
              User ID
            </span>
            <div className="flex items-center gap-1.5 mt-0.5">
              <span className="font-mono text-[11px] text-neutral-800 truncate">{user.id}</span>
              <button
                type="button"
                onClick={handleCopyId}
                className="text-neutral-400 hover:text-neutral-700 p-0.5"
                title={t('common.copy', 'Copy')}
              >
                {copiedId ? (
                  <Check className="h-3 w-3 text-emerald-600" />
                ) : (
                  <Copy className="h-3 w-3" />
                )}
              </button>
            </div>
          </div>

          <div>
            <span className="text-[10px] text-neutral-500 font-medium uppercase tracking-wider block">
              Role
            </span>
            <div className="flex items-center gap-1 mt-0.5 font-medium text-neutral-900 capitalize">
              <Shield className="h-3.5 w-3.5 text-neutral-600" />
              <span>{t(`iam.roles.${user.role}`, user.role)}</span>
            </div>
          </div>

          <div>
            <span className="text-[10px] text-neutral-500 font-medium uppercase tracking-wider block">
              Tenant ID
            </span>
            <div className="flex items-center gap-1 mt-0.5 font-mono text-[11px] text-neutral-700">
              <Building className="h-3.5 w-3.5 text-neutral-400" />
              <span>{user.tenant_id}</span>
            </div>
          </div>

          <div>
            <span className="text-[10px] text-neutral-500 font-medium uppercase tracking-wider block">
              Joined Date
            </span>
            <div className="flex items-center gap-1 mt-0.5 text-[11px] text-neutral-700">
              <Calendar className="h-3.5 w-3.5 text-neutral-400" />
              <span>{new Date(user.create_at).toLocaleDateString()}</span>
            </div>
          </div>
        </div>

        {/* Effective Permissions Section */}
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-neutral-900 flex items-center gap-1.5">
              <Key className="h-3.5 w-3.5 text-neutral-600" />
              <span>Effective Permissions (RBAC + ACL)</span>
            </span>
            <span className="text-[10px] text-neutral-500 font-mono">
              {isOwner ? 'Universal (*)' : `${effectivePerms.length} computed`}
            </span>
          </div>

          {isOwner ? (
            <div className="rounded border border-indigo-200 bg-indigo-50/50 p-2.5 text-xs text-indigo-950 flex items-start gap-2">
              <CheckCircle2 className="h-4 w-4 text-indigo-600 mt-0.5 shrink-0" />
              <div>
                <span className="font-semibold">Super-Admin Wildcard (*)</span>: As workspace
                Owner, this account possesses unconstrained access to all present and future system
                modules and actions.
              </div>
            </div>
          ) : (
            <div className="flex flex-wrap gap-1.5 max-h-36 overflow-y-auto p-2 border border-neutral-200 rounded-md bg-white">
              {effectivePerms.length > 0 ? (
                effectivePerms.map((perm) => (
                  <Badge
                    key={perm}
                    variant="outline"
                    className={`text-[10px] font-mono px-2 py-0.5 ${
                      customAcl.includes(perm)
                        ? 'border-amber-300 bg-amber-50 text-amber-900 font-medium'
                        : 'border-neutral-200 bg-neutral-50 text-neutral-700'
                    }`}
                  >
                    {perm}
                    {customAcl.includes(perm) && ' (custom ACL)'}
                  </Badge>
                ))
              ) : (
                <span className="text-xs text-neutral-400 italic">No permissions assigned.</span>
              )}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
