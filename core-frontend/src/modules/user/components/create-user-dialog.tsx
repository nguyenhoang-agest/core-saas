'use client';

import React, { useState } from 'react';
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
import { useI18n } from '@/core/i18n/i18n-context';
import { useSession } from 'next-auth/react';
import type { CustomSession } from '@/core/auth/auth-options';
import { Badge } from '@/components/ui/badge';
import { useCreateUser, usePermissionsMeta, useRoles, useAssignUserRoles } from '../hooks';
import { parseBackendError } from '../api';
import type { User, RoleItem } from '../types';
import { ChevronDown, ChevronRight, Key, Shield, Lock } from 'lucide-react';

interface CreateUserDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: (user: User) => void;
}

export function CreateUserDialog({ open, onOpenChange, onSuccess }: CreateUserDialogProps) {
  const { t } = useI18n();
  const { data: session } = useSession();
  const customSession = session as CustomSession | null;
  const actorPermissions = customSession?.user?.permissions || [];
  const isActorOwner = actorPermissions.includes('*');

  const createUserMutation = useCreateUser();
  const assignRolesMutation = useAssignUserRoles();
  const { data: roles = [], isLoading: rolesLoading } = useRoles();
  const { data: permissionsCatalog = [] } = usePermissionsMeta();

  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [selectedRoleIds, setSelectedRoleIds] = useState<string[]>([]);
  const [phone, setPhone] = useState('');
  const [customAcl, setCustomAcl] = useState<string[]>([]);
  const [aclExpanded, setAclExpanded] = useState(false);

  const [formError, setFormError] = useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

  // Initialize selected role to "member" when roles load
  React.useEffect(() => {
    if (roles.length > 0 && selectedRoleIds.length === 0) {
      const memberRole = roles.find((r) => r.slug === 'member');
      if (memberRole) {
        setSelectedRoleIds([memberRole.id]);
      } else {
        setSelectedRoleIds([roles[0].id]);
      }
    }
  }, [roles, selectedRoleIds.length]);

  const fallbackRoles: RoleItem[] = [
    { id: 'role-owner', tenant_id: '', slug: 'owner', name: 'Owner', description: 'Full super-admin ownership', is_system: true, version: 1, permissions: ['*'], user_count: 0, create_at: '', update_at: '' },
    { id: 'role-admin', tenant_id: '', slug: 'admin', name: 'Admin', description: 'High-privilege administrator', is_system: true, version: 1, permissions: [], user_count: 0, create_at: '', update_at: '' },
    { id: 'role-member', tenant_id: '', slug: 'member', name: 'Member', description: 'Standard operational user', is_system: true, version: 1, permissions: [], user_count: 0, create_at: '', update_at: '' },
    { id: 'role-viewer', tenant_id: '', slug: 'viewer', name: 'Viewer', description: 'Read-only audit & inspection', is_system: true, version: 1, permissions: [], user_count: 0, create_at: '', update_at: '' },
  ];

  const effectiveRoles = roles.length > 0 ? roles : fallbackRoles;

  const resetForm = () => {
    setName('');
    setEmail('');
    const defaultRole = roles.find((r) => r.slug === 'member') || roles[0];
    setSelectedRoleIds(defaultRole ? [defaultRole.id] : []);
    setPhone('');
    setCustomAcl([]);
    setFormError(null);
    setFieldErrors({});
    setAclExpanded(false);
  };

  const handleClose = () => {
    resetForm();
    onOpenChange(false);
  };

  const toggleRole = (roleId: string) => {
    setSelectedRoleIds((prev) => {
      if (prev.includes(roleId)) {
        if (prev.length <= 1) return prev; // At least one role must remain selected
        return prev.filter((id) => id !== roleId);
      }
      return [...prev, roleId];
    });
  };

  const togglePermission = (permId: string) => {
    setCustomAcl((prev) =>
      prev.includes(permId) ? prev.filter((p) => p !== permId) : [...prev, permId]
    );
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError(null);
    setFieldErrors({});

    if (!name.trim()) {
      setFieldErrors((prev) => ({ ...prev, name: t('iam.create.err_name_required', 'Full name is required') }));
      return;
    }
    if (!email.trim() || !email.includes('@')) {
      setFieldErrors((prev) => ({ ...prev, email: t('iam.create.err_email_invalid', 'A valid email address is required') }));
      return;
    }

    // Determine primary system role slug for legacy compatibility
    const selectedRoles = effectiveRoles.filter((r) => selectedRoleIds.includes(r.id));
    const primarySlug = selectedRoles[0]?.slug || 'member';
    const validSystemRoles = ['owner', 'admin', 'member', 'viewer'];
    const primarySystemRole = validSystemRoles.includes(primarySlug) ? primarySlug : 'member';

    try {
      const created = await createUserMutation.mutateAsync({
        name: name.trim(),
        email: email.trim().toLowerCase(),
        role: primarySystemRole,
        phone: phone.trim() || undefined,
        permissions: customAcl.length > 0 ? customAcl : undefined,
      });

      // Assign multiple roles if custom roles exist in DB
      if (roles.length > 0 && selectedRoleIds.length > 0) {
        try {
          await assignRolesMutation.mutateAsync({
            userId: created.id,
            input: { role_ids: selectedRoleIds },
          });
        } catch (assignErr) {
          console.error('Failed to assign multi-roles after user creation:', assignErr);
        }
      }

      handleClose();
      onSuccess(created);
    } catch (err) {
      const parsed = parseBackendError(err);
      setFormError(parsed.message);
      if (parsed.fieldErrors) {
        setFieldErrors(parsed.fieldErrors);
      }
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogClose onClose={handleClose} />
      <div className="max-h-[85vh] overflow-y-auto px-6 py-6">
        <DialogHeader>
          <div className="flex items-center gap-2 mb-1">
            <div className="p-1.5 rounded-md border border-neutral-200 bg-neutral-50">
              <Shield className="h-4 w-4 text-neutral-800" />
            </div>
            <DialogTitle>{t('user.create.title', t('iam.create.title', 'Create User'))}</DialogTitle>
          </div>
          <DialogDescription className="text-xs">
            {t(
              'user.create.description',
              t(
                'iam.create.description',
                'Add a new user to this tenant workspace. Set their primary RBAC role and optionally grant direct ACL permissions.'
              )
            )}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4 pt-4">
          {formError && (
            <div className="rounded-md border border-red-200 bg-red-50 p-3 text-xs text-red-800 animate-in fade-in-0">
              <span className="font-semibold">{t('common.error', 'Error')}:</span> {formError}
            </div>
          )}

          {/* Full Name */}
          <div className="space-y-1">
            <Label htmlFor="create-user-name" className="text-xs">
              {t('iam.create.name_label', 'Full Name')} <span className="text-red-500">*</span>
            </Label>
            <Input
              id="create-user-name"
              placeholder={t('iam.create.name_placeholder')}
              value={name}
              onChange={(e) => {
                setName(e.target.value);
                if (fieldErrors.name) {
                  setFieldErrors((prev) => {
                    const next = { ...prev };
                    delete next.name;
                    return next;
                  });
                }
              }}
              className="text-xs"
            />
            {fieldErrors.name && (
              <p className="text-[11px] text-red-600 mt-1">{fieldErrors.name}</p>
            )}
          </div>

          {/* Email */}
          <div className="space-y-1">
            <Label htmlFor="create-user-email" className="text-xs">
              {t('iam.create.email_label', 'Email Address')} <span className="text-red-500">*</span>
            </Label>
            <Input
              id="create-user-email"
              type="email"
              placeholder={t('iam.create.email_placeholder')}
              value={email}
              onChange={(e) => {
                setEmail(e.target.value);
                if (fieldErrors.email) {
                  setFieldErrors((prev) => {
                    const next = { ...prev };
                    delete next.email;
                    return next;
                  });
                }
              }}
              className="text-xs font-mono"
            />
            {fieldErrors.email && (
              <p className="text-[11px] text-red-600 mt-1">{fieldErrors.email}</p>
            )}
          </div>

          {/* Dynamic Role Selection */}
          <div className="space-y-1.5">
            <div className="flex items-center justify-between">
              <Label className="text-xs font-medium">
                {t('iam.create.role_label', 'Assigned Roles')} <span className="text-red-500">*</span>
              </Label>
              <span className="text-[10px] text-neutral-500 font-mono">
                {selectedRoleIds.length} selected
              </span>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 pt-1 max-h-56 overflow-y-auto pr-1">
              {rolesLoading ? (
                <div className="col-span-2 py-4 text-center text-xs text-neutral-400">
                  {t('common.loading', 'Loading roles...')}
                </div>
              ) : (
                effectiveRoles.map((r) => {
                  const isOwnerRole = r.slug === 'owner';
                  const isOwnerDisabled = isOwnerRole && !isActorOwner;
                  const isSelected = selectedRoleIds.includes(r.id);

                  return (
                    <div
                      key={r.id}
                      onClick={() => {
                        if (!isOwnerDisabled) toggleRole(r.id);
                      }}
                      className={`flex flex-col p-2.5 rounded-md border text-xs transition-colors select-none ${
                        isOwnerDisabled
                          ? 'opacity-50 bg-neutral-100 border-neutral-200 cursor-not-allowed'
                          : isSelected
                          ? 'border-neutral-900 bg-neutral-50/80 ring-1 ring-neutral-900 cursor-pointer'
                          : 'border-neutral-200 bg-white hover:border-neutral-300 cursor-pointer'
                      }`}
                    >
                      <div className="flex items-center justify-between mb-1">
                        <div className="flex items-center gap-1.5 min-w-0">
                          <span className="font-semibold capitalize text-neutral-900 truncate">
                            {r.name}
                          </span>
                          {r.is_system ? (
                            <Badge
                              variant="secondary"
                              className="text-[9px] px-1 py-0 h-4 bg-neutral-100 text-neutral-600 border border-neutral-200"
                            >
                              System
                            </Badge>
                          ) : (
                            <Badge
                              variant="outline"
                              className="text-[9px] px-1 py-0 h-4 bg-indigo-50 text-indigo-700 border-indigo-200"
                            >
                              Custom
                            </Badge>
                          )}
                        </div>
                        <input
                          type="checkbox"
                          checked={isSelected}
                          disabled={isOwnerDisabled}
                          onChange={() => {
                            if (!isOwnerDisabled) toggleRole(r.id);
                          }}
                          className="accent-neutral-900 h-3.5 w-3.5 rounded"
                        />
                      </div>
                      {r.description && (
                        <span className="text-[10px] text-neutral-500 leading-tight line-clamp-1">
                          {r.description}
                        </span>
                      )}
                      <div className="flex items-center justify-between mt-1 text-[10px] text-neutral-400">
                        <span>{r.permissions.length} capabilities</span>
                        {isOwnerDisabled && (
                          <span className="text-amber-600 font-medium flex items-center gap-0.5">
                            <Lock className="h-2.5 w-2.5" /> Requires Owner
                          </span>
                        )}
                      </div>
                    </div>
                  );
                })
              )}
            </div>
            <p className="text-[10px] text-neutral-500">
              Select one or more roles. The first selected role serves as the primary system role.
            </p>
          </div>

          {/* Phone Number */}
          <div className="space-y-1">
            <Label htmlFor="create-user-phone" className="text-xs">
              {t('iam.create.phone_label', 'Phone Number (Optional)')}
            </Label>
            <Input
              id="create-user-phone"
              type="tel"
              value={phone}
              onChange={(e) => setPhone(e.target.value)}
              className="text-xs"
            />
          </div>

          {/* Granular ACL Custom Permissions Accordion */}
          <div className="border border-neutral-200 rounded-md bg-neutral-50/50 p-3 space-y-2">
            <button
              type="button"
              onClick={() => setAclExpanded(!aclExpanded)}
              className="flex items-center justify-between w-full text-xs font-medium text-neutral-800"
            >
              <div className="flex items-center gap-1.5">
                <Key className="h-3.5 w-3.5 text-neutral-600" />
                <span>{t('iam.create.acl_accordion', 'Custom ACL Permissions (Direct Overrides)')}</span>
                {customAcl.length > 0 && (
                  <span className="ml-1.5 px-1.5 py-0.2 rounded-full bg-amber-100 text-amber-800 text-[10px] font-mono">
                    +{customAcl.length}
                  </span>
                )}
              </div>
              {aclExpanded ? (
                <ChevronDown className="h-4 w-4 text-neutral-500" />
              ) : (
                <ChevronRight className="h-4 w-4 text-neutral-500" />
              )}
            </button>

            {aclExpanded && (
              <div className="pt-2 border-t border-neutral-200 space-y-2">
                <p className="text-[11px] text-neutral-500">
                  {t(
                    'iam.create.acl_hint',
                    'Direct ACL permissions augment the base role and are granted explicitly to this user account.'
                  )}
                </p>

                <div className="max-h-44 overflow-y-auto space-y-1.5 pr-1">
                  {permissionsCatalog.map((perm) => {
                    const isChecked = customAcl.includes(perm.id);
                    return (
                      <label
                        key={perm.id}
                        className={`flex items-start gap-2 p-2 rounded border text-xs cursor-pointer transition-colors ${
                          isChecked
                            ? 'border-amber-300 bg-amber-50/60'
                            : 'border-neutral-200 bg-white hover:bg-neutral-50'
                        }`}
                      >
                        <input
                          type="checkbox"
                          checked={isChecked}
                          onChange={() => togglePermission(perm.id)}
                          className="accent-neutral-900 rounded mt-0.5"
                        />
                        <div className="flex-1 min-w-0">
                          <div className="font-mono text-[11px] text-neutral-900 font-medium">
                            {perm.id}
                          </div>
                          <div className="text-[10px] text-neutral-500 truncate">
                            {perm.name}
                          </div>
                        </div>
                      </label>
                    );
                  })}
                </div>
              </div>
            )}
          </div>

          <DialogFooter className="pt-3">
            <Button
              type="button"
              variant="outline"
              onClick={handleClose}
              disabled={createUserMutation.isPending}
              className="text-xs"
            >
              {t('common.cancel', 'Cancel')}
            </Button>
            <Button
              type="submit"
              disabled={createUserMutation.isPending}
              className="text-xs min-w-[100px]"
            >
              {createUserMutation.isPending
                ? t('common.creating', 'Creating...')
                : t('user.create.submit_btn', t('iam.create.submit_btn', 'Create User'))}
            </Button>
          </DialogFooter>
        </form>
      </div>
    </Dialog>
  );
}
