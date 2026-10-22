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
  User,
  RefreshCw,
  Eye,
  EyeOff,
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
  onSuccess?: (createdTenant: Tenant, owner?: { name: string; email: string }) => void;
}

export function CreateTenantDialog({ open, onOpenChange, onSuccess }: CreateTenantDialogProps) {
  const { formatError, t } = useI18n();
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [isolateMode, setIsolateMode] = useState<TenantIsolateMode>('shared');
  const [dsn, setDsn] = useState('');
  const [createOwner, setCreateOwner] = useState(true);
  const [ownerName, setOwnerName] = useState('');
  const [ownerEmail, setOwnerEmail] = useState('');
  const [ownerPassword, setOwnerPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [clientErrors, setClientErrors] = useState<Record<string, string>>({});

  const createTenantMutation = useCreateTenant();

  const resetForm = () => {
    setName('');
    setDescription('');
    setIsolateMode('shared');
    setDsn('');
    setCreateOwner(true);
    setOwnerName('');
    setOwnerEmail('');
    setOwnerPassword('');
    setShowPassword(false);
    setClientErrors({});
    createTenantMutation.reset();
  };

  const generatePassword = () => {
    const chars = 'abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789!@#$%&*';
    const array = new Uint8Array(12);
    crypto.getRandomValues(array);
    const pwd = Array.from(array, (byte) => chars[byte % chars.length]).join('');
    setOwnerPassword(pwd);
    setShowPassword(true);
    setClientErrors((prev) => {
      const next = { ...prev };
      delete next.ownerPassword;
      return next;
    });
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

    if (createOwner) {
      if (!ownerName.trim()) {
        errors.ownerName = t(
          'tenant.create.validation.owner_name_required',
          'Owner name is required.'
        );
      } else if (ownerName.trim().length < 2) {
        errors.ownerName = t(
          'tenant.create.validation.owner_name_min_length',
          'Owner name must be at least 2 characters.'
        );
      }

      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (!ownerEmail.trim()) {
        errors.ownerEmail = t(
          'tenant.create.validation.owner_email_required',
          'Owner email is required.'
        );
      } else if (!emailRegex.test(ownerEmail.trim())) {
        errors.ownerEmail = t(
          'tenant.create.validation.owner_email_invalid',
          'Please enter a valid owner email address.'
        );
      }

      if (!ownerPassword) {
        errors.ownerPassword = t(
          'tenant.create.validation.owner_password_required',
          'Owner password is required.'
        );
      } else if (ownerPassword.length < 8) {
        errors.ownerPassword = t(
          'tenant.create.validation.owner_password_min',
          'Password must be at least 8 characters.'
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

    if (createOwner) {
      payload.owner = {
        name: ownerName.trim(),
        email: ownerEmail.trim(),
        password: ownerPassword,
      };
    }

    createTenantMutation.mutate(payload, {
      onSuccess: (newTenant) => {
        const ownerInfo = createOwner
          ? { name: ownerName.trim(), email: ownerEmail.trim() }
          : undefined;
        resetForm();
        onOpenChange(false);
        if (onSuccess) {
          onSuccess(newTenant, ownerInfo);
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

      <form onSubmit={handleSubmit} noValidate className="space-y-4">
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

        {/* Workspace Owner Section */}
        <div className="rounded-lg border border-neutral-200 bg-neutral-50/50 p-3.5 space-y-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-1.5">
              <User className="h-4 w-4 text-neutral-900" />
              <span className="text-xs font-semibold text-neutral-900">
                {t('tenant.create.owner_section_title', 'Initial Workspace Owner')}
              </span>
            </div>
            <label className="flex items-center gap-1.5 text-xs text-neutral-600 cursor-pointer select-none">
              <input
                type="checkbox"
                checked={createOwner}
                onChange={(e) => setCreateOwner(e.target.checked)}
                className="rounded border-neutral-300 text-neutral-900 focus:ring-neutral-900 h-3.5 w-3.5"
              />
              <span>{t('tenant.create.provision_owner_toggle', 'Create Owner Account')}</span>
            </label>
          </div>
          <p className="text-[11px] text-neutral-500 leading-tight">
            {t(
              'tenant.create.owner_section_desc',
              'Creates the primary administrator account with role "owner" for this workspace.'
            )}
          </p>

          {createOwner && (
            <div className="space-y-3 pt-1">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                {/* Owner Name */}
                <div className="space-y-1">
                  <Label htmlFor="owner-name" required className="text-xs">
                    {t('tenant.create.owner_name_label', 'Owner Full Name')}
                  </Label>
                  <Input
                    id="owner-name"
                    value={ownerName}
                    onChange={(e) => {
                      setOwnerName(e.target.value);
                      if (clientErrors.ownerName) {
                        setClientErrors((prev) => {
                          const next = { ...prev };
                          delete next.ownerName;
                          return next;
                        });
                      }
                    }}
                    placeholder={t('tenant.create.owner_name_placeholder', 'e.g. Tony Stark')}
                    className="h-8 text-xs"
                    error={Boolean(clientErrors.ownerName || backendError?.fieldErrors?.['owner.name'])}
                  />
                  {(clientErrors.ownerName || backendError?.fieldErrors?.['owner.name']) && (
                    <p className="text-[11px] text-red-600">
                      {clientErrors.ownerName || backendError?.fieldErrors?.['owner.name']}
                    </p>
                  )}
                </div>

                {/* Owner Email */}
                <div className="space-y-1">
                  <Label htmlFor="owner-email" required className="text-xs">
                    {t('tenant.create.owner_email_label', 'Owner Email')}
                  </Label>
                  <Input
                    id="owner-email"
                    type="email"
                    value={ownerEmail}
                    onChange={(e) => {
                      setOwnerEmail(e.target.value);
                      if (clientErrors.ownerEmail) {
                        setClientErrors((prev) => {
                          const next = { ...prev };
                          delete next.ownerEmail;
                          return next;
                        });
                      }
                    }}
                    placeholder={t('tenant.create.owner_email_placeholder', 'e.g. tony@stark.com')}
                    className="h-8 text-xs"
                    error={Boolean(clientErrors.ownerEmail || backendError?.fieldErrors?.['owner.email'])}
                  />
                  {(clientErrors.ownerEmail || backendError?.fieldErrors?.['owner.email']) && (
                    <p className="text-[11px] text-red-600">
                      {clientErrors.ownerEmail || backendError?.fieldErrors?.['owner.email']}
                    </p>
                  )}
                </div>
              </div>

              {/* Owner Password */}
              <div className="space-y-1">
                <div className="flex items-center justify-between">
                  <Label htmlFor="owner-password" required className="text-xs">
                    {t('tenant.create.owner_password_label', 'Initial Password')}
                  </Label>
                  <button
                    type="button"
                    onClick={generatePassword}
                    className="text-[11px] font-medium text-neutral-600 hover:text-neutral-900 flex items-center gap-1 underline underline-offset-2"
                  >
                    <RefreshCw className="h-3 w-3" />
                    <span>{t('tenant.create.generate_password', 'Generate')}</span>
                  </button>
                </div>
                <div className="relative">
                  <Input
                    id="owner-password"
                    type={showPassword ? 'text' : 'password'}
                    value={ownerPassword}
                    onChange={(e) => {
                      setOwnerPassword(e.target.value);
                      if (clientErrors.ownerPassword) {
                        setClientErrors((prev) => {
                          const next = { ...prev };
                          delete next.ownerPassword;
                          return next;
                        });
                      }
                    }}
                    placeholder={t('tenant.create.owner_password_placeholder', 'At least 8 characters')}
                    className="h-8 text-xs pr-8 font-mono"
                    error={Boolean(clientErrors.ownerPassword || backendError?.fieldErrors?.['owner.password'])}
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    tabIndex={-1}
                    className="absolute right-2 top-1/2 -translate-y-1/2 text-neutral-400 hover:text-neutral-700"
                  >
                    {showPassword ? <EyeOff className="h-3.5 w-3.5" /> : <Eye className="h-3.5 w-3.5" />}
                  </button>
                </div>
                {(clientErrors.ownerPassword || backendError?.fieldErrors?.['owner.password']) && (
                  <p className="text-[11px] text-red-600">
                    {clientErrors.ownerPassword || backendError?.fieldErrors?.['owner.password']}
                  </p>
                )}
              </div>
            </div>
          )}
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
