'use client';

import React, { useState, useEffect, useMemo } from 'react';
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
import { useUpdateUser, usePermissionsMeta, useRolesMeta, useRoles, useAssignUserRoles } from '../hooks';
import { parseBackendError } from '../api';
import type { User, UserRole, UserStatus, RoleItem } from '../types';
import { Edit2, Key, Check, Shield, Lock } from 'lucide-react';

interface EditUserDialogProps {
  user: User | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: (user: User) => void;
}

export function EditUserDialog({ user, open, onOpenChange, onSuccess }: EditUserDialogProps) {
  const { t } = useI18n();
  const { data: session } = useSession();
  const customSession = session as CustomSession | null;
  const actorPermissions = customSession?.user?.permissions || [];
  const isActorOwner = actorPermissions.includes('*');

  const updateMutation = useUpdateUser();
  const assignRolesMutation = useAssignUserRoles();
  const { data: roles = [], isLoading: rolesLoading } = useRoles();
  const { data: permissionsCatalog = [] } = usePermissionsMeta();
  const { data: rolesMeta = [] } = useRolesMeta();

  const [name, setName] = useState('');
  const [selectedRoleIds, setSelectedRoleIds] = useState<string[]>([]);
  const [status, setStatus] = useState<UserStatus>('active');
  const [phone, setPhone] = useState('');
  const [customAcl, setCustomAcl] = useState<string[]>([]);

  const [formError, setFormError] = useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

  const fallbackRoles: RoleItem[] = [
    { id: 'role-owner', tenant_id: '', slug: 'owner', name: 'Owner', description: 'Full super-admin ownership', is_system: true, version: 1, permissions: ['*'], user_count: 0, create_at: '', update_at: '' },
    { id: 'role-admin', tenant_id: '', slug: 'admin', name: 'Admin', description: 'High-privilege administrator', is_system: true, version: 1, permissions: [], user_count: 0, create_at: '', update_at: '' },
    { id: 'role-member', tenant_id: '', slug: 'member', name: 'Member', description: 'Standard operational user', is_system: true, version: 1, permissions: [], user_count: 0, create_at: '', update_at: '' },
    { id: 'role-viewer', tenant_id: '', slug: 'viewer', name: 'Viewer', description: 'Read-only audit & inspection', is_system: true, version: 1, permissions: [], user_count: 0, create_at: '', update_at: '' },
  ];

  const effectiveRoles = roles.length > 0 ? roles : fallbackRoles;

  useEffect(() => {
    if (user) {
      setName(user.name);
      setStatus(user.status);
      setPhone(user.phone || '');
      setCustomAcl(user.permissions || []);
      setFormError(null);
      setFieldErrors({});

      // Populate selected role IDs based on user's current roles
      const matchingIds: string[] = [];
      if (effectiveRoles.length > 0) {
        if (user.roles && user.roles.length > 0) {
          for (const ur of user.roles) {
            const found = effectiveRoles.find((r) => r.id === ur || r.slug === ur);
            if (found && !matchingIds.includes(found.id)) {
              matchingIds.push(found.id);
            }
          }
        }
        if (matchingIds.length === 0 && user.role) {
          const found = effectiveRoles.find((r) => r.slug === user.role || r.id === user.role);
          if (found) {
            matchingIds.push(found.id);
          }
        }
        if (matchingIds.length > 0) {
          setSelectedRoleIds(matchingIds);
        } else if (effectiveRoles.length > 0) {
          const memberRole = effectiveRoles.find((r) => r.slug === 'member') || effectiveRoles[0];
          setSelectedRoleIds([memberRole.id]);
        }
      }
    }
  }, [user, roles]);

  const handleClose = () => {
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

  // Find default permissions for currently selected primary role
  const roleDefaultPerms = useMemo(() => {
    const selectedRoles = effectiveRoles.filter((r) => selectedRoleIds.includes(r.id));
    const primary = selectedRoles[0];
    if (!primary) return [];
    if (primary.slug === 'owner') return ['*'];
    return primary.permissions || [];
  }, [selectedRoleIds, effectiveRoles]);

  const togglePermission = (permId: string) => {
    setCustomAcl((prev) =>
      prev.includes(permId) ? prev.filter((p) => p !== permId) : [...prev, permId]
    );
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user) return;

    setFormError(null);
    setFieldErrors({});

    if (!name.trim()) {
      setFieldErrors({ name: t('iam.create.err_name_required', 'Full name is required') });
      return;
    }

    const selectedRoles = effectiveRoles.filter((r) => selectedRoleIds.includes(r.id));
    const primarySlug = selectedRoles[0]?.slug || 'member';
    const validSystemRoles = ['owner', 'admin', 'member', 'viewer'];
    const primarySystemRole = validSystemRoles.includes(primarySlug) ? primarySlug : 'member';

    try {
      const updated = await updateMutation.mutateAsync({
        id: user.id,
        input: {
          name: name.trim(),
          role: primarySystemRole,
          status,
          phone: phone.trim() || undefined,
          permissions: customAcl,
        },
      });

      if (roles.length > 0 && selectedRoleIds.length > 0) {
        try {
          await assignRolesMutation.mutateAsync({
            userId: user.id,
            input: { role_ids: selectedRoleIds },
          });
        } catch (assignErr) {
          console.error('Failed to assign multi-roles after user update:', assignErr);
        }
      }

      handleClose();
      onSuccess(updated);
    } catch (err) {
      const parsed = parseBackendError(err);
      setFormError(parsed.message);
      if (parsed.fieldErrors) {
        setFieldErrors(parsed.fieldErrors);
      }
    }
  };

  if (!user) return null;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogClose onClose={handleClose} />
      <div className="max-h-[85vh] overflow-y-auto px-6 py-6">
        <DialogHeader>
          <div className="flex items-center gap-2 mb-1">
            <div className="p-1.5 rounded-md border border-neutral-200 bg-neutral-50">
              <Edit2 className="h-4 w-4 text-neutral-800" />
            </div>
            <DialogTitle>
              {t('iam.edit.title', { name: user.name }, 'Edit Member: {name}')}
            </DialogTitle>
          </div>
          <DialogDescription className="text-xs">
            {t(
              'iam.edit.description',
              'Update member identity, assign functional roles, and modify granular ACL overrides.'
            )}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4 pt-4">
          {formError && (
            <div className="rounded-md border border-red-200 bg-red-50 p-3 text-xs text-red-800 animate-in fade-in-0">
              <span className="font-semibold">{t('common.error', 'Error')}:</span> {formError}
            </div>
          )}

          {/* Email (Readonly) */}
          <div className="space-y-1">
            <Label className="text-xs text-neutral-500">
              {t('iam.create.email_label', 'Email Address')} (Readonly)
            </Label>
            <Input
              value={user.email}
              disabled
              className="text-xs font-mono bg-neutral-100 text-neutral-600 cursor-not-allowed"
            />
          </div>

          {/* Name */}
          <div className="space-y-1">
            <Label htmlFor="edit-user-name" className="text-xs">
              {t('iam.create.name_label', 'Full Name')} <span className="text-red-500">*</span>
            </Label>
            <Input
              id="edit-user-name"
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

          {/* Dynamic Role Selection */}
          <div className="space-y-1.5">
            <div className="flex items-center justify-between">
              <Label className="text-xs font-medium">
                {t('iam.create.role_label', 'Assigned Roles')} <span className="text-red-500">*</span>
              </Label>
              <span className="text-[10px] text-neutral-500 font-mono">
                {selectedRoleIds.length} assigned
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
                  // Non-owner actors cannot grant owner role, nor revoke it if they are not owner
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
              Users inherit permissions from all assigned roles. Primary role corresponds to the first selected role.
            </p>
          </div>

          {/* Status */}
          <div className="space-y-1">
            <Label className="text-xs">{t('iam.list.filter_status', 'Status')}</Label>
            <div className="grid grid-cols-3 gap-2 pt-1">
              {(['active', 'suspended', 'invited'] as UserStatus[]).map((st) => (
                <label
                  key={st}
                  className={`flex items-center justify-between p-2 rounded-md border text-xs cursor-pointer ${
                    status === st
                      ? 'border-neutral-900 bg-neutral-50 font-semibold'
                      : 'border-neutral-200'
                  }`}
                >
                  <span className="capitalize">{t(`iam.statuses.${st}`, st)}</span>
                  <input
                    type="radio"
                    name="edit-user-status"
                    value={st}
                    checked={status === st}
                    onChange={() => setStatus(st)}
                    className="accent-neutral-900 h-3 w-3"
                  />
                </label>
              ))}
            </div>
          </div>

          {/* Phone */}
          <div className="space-y-1">
            <Label htmlFor="edit-user-phone" className="text-xs">
              {t('iam.create.phone_label', 'Phone Number (Optional)')}
            </Label>
            <Input
              id="edit-user-phone"
              type="tel"
              value={phone}
              onChange={(e) => setPhone(e.target.value)}
              className="text-xs"
            />
          </div>

          {/* Granular ACL Custom Permissions */}
          <div className="border border-neutral-200 rounded-md bg-neutral-50/50 p-3 space-y-2">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-neutral-900 flex items-center gap-1.5">
                <Key className="h-3.5 w-3.5 text-neutral-600" />
                <span>{t('iam.edit.acl_title', 'Granular ACL Permissions (Direct Overrides)')}</span>
              </span>
              <span className="text-[10px] text-neutral-500 font-mono">
                {customAcl.length} direct grants
              </span>
            </div>
            <p className="text-[11px] text-neutral-500">
              {t(
                'iam.edit.acl_desc',
                'Toggle individual permissions granted specifically to this user beyond their role default.'
              )}
            </p>

            <div className="max-h-48 overflow-y-auto space-y-1.5 pr-1 pt-1">
              {permissionsCatalog.map((perm) => {
                const isDirect = customAcl.includes(perm.id);
                const isRoleDefault =
                  roleDefaultPerms.includes('*') ||
                  roleDefaultPerms.includes(perm.id) ||
                  roleDefaultPerms.some((p) => p.endsWith(':*') && perm.id.startsWith(p.slice(0, -1)));

                return (
                  <label
                    key={perm.id}
                    className={`flex items-start gap-2 p-2 rounded border text-xs cursor-pointer transition-colors ${
                      isDirect
                        ? 'border-amber-300 bg-amber-50/60'
                        : isRoleDefault
                        ? 'border-indigo-200 bg-indigo-50/40'
                        : 'border-neutral-200 bg-white hover:bg-neutral-50'
                    }`}
                  >
                    <input
                      type="checkbox"
                      checked={isDirect}
                      onChange={() => togglePermission(perm.id)}
                      className="accent-neutral-900 rounded mt-0.5"
                    />
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center justify-between">
                        <span className="font-mono text-[11px] text-neutral-900 font-medium">
                          {perm.id}
                        </span>
                        {isRoleDefault && (
                          <span className="text-[10px] text-indigo-700 bg-indigo-100/70 px-1.5 py-0.2 rounded font-sans">
                            Base Role
                          </span>
                        )}
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

          <DialogFooter className="pt-3">
            <Button
              type="button"
              variant="outline"
              onClick={handleClose}
              disabled={updateMutation.isPending}
              className="text-xs"
            >
              {t('common.cancel', 'Cancel')}
            </Button>
            <Button
              type="submit"
              disabled={updateMutation.isPending}
              className="text-xs min-w-[100px]"
            >
              {updateMutation.isPending
                ? t('common.loading', 'Saving...')
                : t('common.save', 'Save Changes')}
            </Button>
          </DialogFooter>
        </form>
      </div>
    </Dialog>
  );
}
