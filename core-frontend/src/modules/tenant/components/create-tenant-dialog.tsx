'use client';

import React, { useState } from 'react';
import {
  Building2,
  AlertCircle,
  Loader2,
  Layers,
  Database,
  ShieldCheck,
  Lock,
} from 'lucide-react';
import {
  Dialog,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
  DialogClose,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useCreateTenant } from '@/modules/tenant';
import { useI18n } from '@/core/i18n/i18n-context';
import type { Tenant, TenantIsolateMode, CreateTenantInput } from '../types/tenant';

interface CreateTenantDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: (createdTenant: Tenant) => void;
}

export function CreateTenantDialog({ open, onOpenChange, onSuccess }: CreateTenantDialogProps) {
  const { formatError, t } = useI18n();
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [isolateMode, setIsolateMode] = useState<TenantIsolateMode>('shared');
  const [dsn, setDsn] = useState('');
  const [clientErrors, setClientErrors] = useState<Record<string, string>>({});

  const createTenantMutation = useCreateTenant();

  const resetForm = () => {
    setName('');
    setDescription('');
    setIsolateMode('shared');
    setDsn('');
    setClientErrors({});
    createTenantMutation.reset();
  };

  const validate = (): boolean => {
    const errors: Record<string, string> = {};
    if (!name.trim()) {
      errors.name = t('tenant.create.validation.name_required', 'Tenant name is required.');
    } else if (name.trim().length < 3) {
      errors.name = t(
        'tenant.create.validation.name_min_length',
        'Tenant name must be at least 3 characters.'
      );
    }

    if (isolateMode === 'dedicated') {
      if (!dsn.trim()) {
        errors.dsn = t(
          'tenant.create.validation.dsn_required',
          'Database connection DSN is required for dedicated isolation mode.'
        );
      }
    }

    setClientErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;

    const payload: CreateTenantInput = {
      name: name.trim(),
      description: description.trim() || undefined,
      isolate_mode: isolateMode,
    };

    if (isolateMode === 'dedicated') {
      // Both connection pool and SSL mode are strictly system-enforced for security.
      payload.db_config = {
        dsn: dsn.trim(),
        ssl_mode: 'require',
      };
    }

    createTenantMutation.mutate(payload, {
      onSuccess: (newTenant) => {
        resetForm();
        onOpenChange(false);
        if (onSuccess) {
          onSuccess(newTenant);
        }
      },
    });
  };

  const backendError = createTenantMutation.isError
    ? formatError(createTenantMutation.error)
    : null;

  return (
    <Dialog
      open={open}
      onOpenChange={(next) => {
        if (!next) resetForm();
        onOpenChange(next);
      }}
    >
      <DialogClose onClose={() => onOpenChange(false)} />
      <DialogHeader>
        <DialogTitle className="flex items-center gap-2">
          <Building2 className="h-5 w-5 text-neutral-900 stroke-[1.75]" />
          <span>{t('tenant.create.title', 'Provision New Tenant')}</span>
        </DialogTitle>
        <DialogDescription>
          {t(
            'tenant.create.description',
            'Create an isolated organization workspace with dedicated or shared multi-tenancy.'
          )}
        </DialogDescription>
      </DialogHeader>

      <form onSubmit={handleSubmit} className="space-y-4">
        {/* Backend Top Error Banner */}
        {backendError && (
          <div className="rounded-md border border-red-200 bg-red-50/70 p-3 text-xs text-red-900 flex items-start gap-2.5">
            <AlertCircle className="h-4 w-4 text-red-600 shrink-0 mt-0.5" />
            <div>
              <span className="font-semibold text-red-800">
                {backendError.code ? `Error ${backendError.code}: ` : ''}
              </span>
              <span>{backendError.message}</span>
            </div>
          </div>
        )}

        {/* Tenant Name Field */}
        <div className="space-y-1.5">
          <Label htmlFor="tenant-name" required>
            {t('tenant.create.name_label', 'Tenant Name')}
          </Label>
          <Input
            id="tenant-name"
            value={name}
            onChange={(e) => {
              setName(e.target.value);
              if (clientErrors.name) {
                setClientErrors((prev) => {
                  const next = { ...prev };
                  delete next.name;
                  return next;
                });
              }
            }}
            placeholder={t(
              'tenant.create.name_placeholder',
              'e.g. Acme Corporation, Stark Industries'
            )}
            error={Boolean(clientErrors.name || backendError?.fieldErrors?.name)}
          />
          {(clientErrors.name || backendError?.fieldErrors?.name) && (
            <p className="text-[11px] text-red-600">
              {clientErrors.name || backendError?.fieldErrors?.name}
            </p>
          )}
        </div>

        {/* Tenant Description Field */}
        <div className="space-y-1.5">
          <Label htmlFor="tenant-description">
            {t('tenant.create.description_label', 'Description (Optional)')}
          </Label>
          <Input
            id="tenant-description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder={t(
              'tenant.create.description_placeholder',
              'e.g. Primary production workspace'
            )}
          />
        </div>

        {/* Isolation Mode Selection (Flat Cards) */}
        <div className="space-y-1.5 pt-1">
          <Label>{t('tenant.create.isolation_mode_label', 'Database Isolation Mode')}</Label>
          <div className="grid grid-cols-2 gap-3">
            {/* Shared Mode Card */}
            <button
              type="button"
              onClick={() => setIsolateMode('shared')}
              className={`p-3 rounded-lg border text-left transition-colors flex flex-col justify-between ${
                isolateMode === 'shared'
                  ? 'border-neutral-900 bg-neutral-900 text-white'
                  : 'border-neutral-200 bg-white text-neutral-900 hover:bg-neutral-50'
              }`}
            >
              <div className="flex items-center justify-between w-full mb-1">
                <span className="text-xs font-semibold">
                  {t('tenant.isolation.shared', 'Shared Database')}
                </span>
                <Layers
                  className={`h-3.5 w-3.5 ${
                    isolateMode === 'shared' ? 'text-white' : 'text-neutral-500'
                  }`}
                />
              </div>
              <p
                className={`text-[11px] leading-tight ${
                  isolateMode === 'shared' ? 'text-neutral-300' : 'text-neutral-500'
                }`}
              >
                {t(
                  'tenant.isolation.shared_desc',
                  'Shared pool with row-level tenant context isolation.'
                )}
              </p>
            </button>

            {/* Dedicated Mode Card */}
            <button
              type="button"
              onClick={() => setIsolateMode('dedicated')}
              className={`p-3 rounded-lg border text-left transition-colors flex flex-col justify-between ${
                isolateMode === 'dedicated'
                  ? 'border-neutral-900 bg-neutral-900 text-white'
                  : 'border-neutral-200 bg-white text-neutral-900 hover:bg-neutral-50'
              }`}
            >
              <div className="flex items-center justify-between w-full mb-1">
                <span className="text-xs font-semibold">
                  {t('tenant.isolation.dedicated', 'Dedicated DB')}
                </span>
                <Database
                  className={`h-3.5 w-3.5 ${
                    isolateMode === 'dedicated' ? 'text-white' : 'text-neutral-500'
                  }`}
                />
              </div>
              <p
                className={`text-[11px] leading-tight ${
                  isolateMode === 'dedicated' ? 'text-neutral-300' : 'text-neutral-500'
                }`}
              >
                {t(
                  'tenant.isolation.dedicated_desc',
                  'Isolated PostgreSQL instance with custom DSN router.'
                )}
              </p>
            </button>
          </div>
        </div>

        {/* Dedicated DB Configuration Fields */}
        {isolateMode === 'dedicated' && (
          <div className="rounded-lg border border-neutral-200 bg-neutral-50/50 p-3.5 space-y-3 duration-200 animate-in fade-in-0">
            <div className="text-[11px] font-semibold uppercase tracking-wider text-neutral-700 flex items-center gap-1.5">
              <Database className="h-3.5 w-3.5 text-neutral-900" />
              <span>{t('tenant.create.dedicated_params', 'Dedicated DB Parameters')}</span>
            </div>

            {/* Connection DSN */}
            <div className="space-y-1.5">
              <Label htmlFor="tenant-dsn" required>
                {t('tenant.create.dsn_label', 'Connection DSN')}
              </Label>
              <Input
                id="tenant-dsn"
                value={dsn}
                onChange={(e) => {
                  setDsn(e.target.value);
                  if (clientErrors.dsn) {
                    setClientErrors((prev) => {
                      const next = { ...prev };
                      delete next.dsn;
                      return next;
                    });
                  }
                }}
                placeholder={t(
                  'tenant.create.dsn_placeholder',
                  'postgres://user:pass@host:5432/tenant_db'
                )}
                className="font-mono text-xs"
                error={Boolean(clientErrors.dsn || backendError?.fieldErrors?.dsn)}
              />
              {(clientErrors.dsn || backendError?.fieldErrors?.dsn) && (
                <p className="text-[11px] text-red-600">
                  {clientErrors.dsn || backendError?.fieldErrors?.dsn}
                </p>
              )}
            </div>

            <div className="grid grid-cols-2 gap-3">
              {/* Connection Pool: STRICTLY DENIED TO USER / PLATFORM MANAGED */}
              <div className="space-y-1.5">
                <Label className="text-neutral-600">
                  {t('tenant.create.pool_policy_label', 'Connection Pool Policy')}
                </Label>
                <div
                  data-testid="connection-pool-managed"
                  className="rounded-md border border-neutral-200 bg-white p-2.5 text-neutral-700 flex items-center gap-2 text-xs"
                >
                  <ShieldCheck className="h-4 w-4 text-emerald-600 shrink-0" />
                  <div className="leading-tight">
                    <span className="font-semibold text-neutral-900">10 connections</span>
                    <p className="text-[10px] text-neutral-500">
                      {t('tenant.create.pool_policy_desc', 'Platform-managed to prevent exhaustion.')}
                    </p>
                  </div>
                </div>
              </div>

              {/* SSL Mode: STRICTLY DENIED TO USER / PLATFORM ENFORCED (require) */}
              <div className="space-y-1.5">
                <Label className="text-neutral-600">
                  {t('tenant.create.ssl_policy_label', 'SSL Security Mode')}
                </Label>
                <div
                  data-testid="ssl-mode-managed"
                  className="rounded-md border border-neutral-200 bg-white p-2.5 text-neutral-700 flex items-center gap-2 text-xs"
                >
                  <Lock className="h-4 w-4 text-emerald-600 shrink-0" />
                  <div className="leading-tight">
                    <span className="font-semibold text-neutral-900">require</span>
                    <p className="text-[10px] text-neutral-500">
                      {t('tenant.create.ssl_policy_desc', 'TLS encryption strictly required for all dedicated databases.')}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Dialog Actions */}
        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={createTenantMutation.isPending}
          >
            {t('common.cancel', 'Cancel')}
          </Button>
          <Button type="submit" disabled={createTenantMutation.isPending}>
            {createTenantMutation.isPending ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                <span>{t('common.creating', 'Creating...')}</span>
              </>
            ) : (
              t('tenant.create.submit', 'Create Tenant')
            )}
          </Button>
        </DialogFooter>
      </form>
    </Dialog>
  );
}
