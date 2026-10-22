'use client';

import React, { useState, useEffect, useMemo } from 'react';
import { Shield, AlertCircle, Check, Copy, Save, X, Lock } from 'lucide-react';
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
import { useCreateRole, useUpdateRole, useCloneRole, usePermissionsCatalog } from '../hooks';
import { parseBackendError } from '../api';
import type { RoleItem, PermissionCatalogItem } from '../types';

interface RoleEditorDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  mode: 'create' | 'edit' | 'clone';
  role?: RoleItem | null;
  actorPermissions?: string[];
}

export function RoleEditorDialog({
  open,
  onOpenChange,
  mode,
  role,
  actorPermissions = ['*'],
}: RoleEditorDialogProps) {
  const { t } = useI18n();
  const { data: catalog = [] } = usePermissionsCatalog();

  const createMutation = useCreateRole();
  const updateMutation = useUpdateRole();
  const cloneMutation = useCloneRole();

  const [name, setName] = useState('');
  const [slug, setSlug] = useState('');
  const [description, setDescription] = useState('');
  const [selectedPerms, setSelectedPerms] = useState<string[]>([]);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  // Group permissions by module
  const groupedPermissions = useMemo(() => {
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

  const isActorOwner = actorPermissions.includes('*');

  // Reset form when opened or role changes
  useEffect(() => {
    if (!open) {
      setErrorMessage(null);
      return;
    }

    if (mode === 'create') {
      setName('');
      setSlug('');
      setDescription('');
      setSelectedPerms([]);
    } else if (mode === 'edit' && role) {
      setName(role.name);
      setSlug(role.slug);
      setDescription(role.description || '');
      setSelectedPerms([...role.permissions]);
    } else if (mode === 'clone' && role) {
      setName(`${role.name} (Copy)`);
      setSlug(`${role.slug}-copy`);
      setDescription(role.description || '');
      setSelectedPerms([...role.permissions]);
    }
    setErrorMessage(null);
  }, [open, mode, role]);

  // Auto-slugify name during creation if slug hasn't been manually diverged
  const handleNameChange = (val: string) => {
    setName(val);
    if (mode === 'create') {
      const generated = val
        .toLowerCase()
        .replace(/[^a-z0-9]+/g, '-')
        .replace(/^-+|-+$/g, '');
      setSlug(generated);
    }
  };

  const isPermAllowedForActor = (key: string) => {
    if (isActorOwner) return true;
    return actorPermissions.includes(key);
  };

  const handleTogglePerm = (key: string) => {
    if (!isPermAllowedForActor(key)) return;

    setSelectedPerms((prev) =>
      prev.includes(key) ? prev.filter((p) => p !== key) : [...prev, key]
    );
  };

  const handleToggleModule = (modulePerms: PermissionCatalogItem[]) => {
    // Only toggle permissions the actor possesses
    const allowedKeys = modulePerms.map((p) => p.key).filter(isPermAllowedForActor);
    const allSelected = allowedKeys.every((k) => selectedPerms.includes(k));

    if (allSelected) {
      setSelectedPerms((prev) => prev.filter((k) => !allowedKeys.includes(k)));
    } else {
      setSelectedPerms((prev) => Array.from(new Set([...prev, ...allowedKeys])));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMessage(null);

    if (!name.trim()) {
      setErrorMessage(t('iam.roles.error_name_required', 'Role name is required.'));
      return;
    }

    if (!slug.trim()) {
      setErrorMessage(t('iam.roles.error_slug_required', 'Role identifier is required.'));
      return;
    }

    if (selectedPerms.length === 0) {
      setErrorMessage(t('iam.roles.error_perms_required', 'Please select at least one permission capability.'));
      return;
    }

    try {
      if (mode === 'create') {
        await createMutation.mutateAsync({
          name: name.trim(),
          slug: slug.trim(),
          description: description.trim(),
          permissions: selectedPerms,
        });
      } else if (mode === 'edit' && role) {
        await updateMutation.mutateAsync({
          id: role.id,
          input: {
            name: name.trim(),
            description: description.trim(),
            permissions: selectedPerms,
            version: role.version,
          },
        });
      } else if (mode === 'clone' && role) {
        await cloneMutation.mutateAsync({
          id: role.id,
          input: {
            name: name.trim(),
            slug: slug.trim(),
            description: description.trim(),
          },
        });
      }
      onOpenChange(false);
    } catch (err) {
      const parsed = parseBackendError(err);
      setErrorMessage(parsed.message);
    }
  };

  const isSubmitting =
    createMutation.isPending || updateMutation.isPending || cloneMutation.isPending;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogClose onClose={() => onOpenChange(false)} />
      <div className="max-h-[85vh] flex flex-col">
        <form onSubmit={handleSubmit} className="flex flex-col h-full overflow-hidden">
          <DialogHeader className="p-2 pb-2">
            <DialogTitle className="text-base font-semibold flex items-center gap-2">
              <Shield className="h-5 w-5 text-indigo-600" />
              <span>
                {mode === 'create' && t('iam.roles.create_title', 'Create Custom Role')}
                {mode === 'edit' && t('iam.roles.edit_title', 'Edit Role Permissions')}
                {mode === 'clone' && t('iam.roles.clone_title', 'Clone Role')}
              </span>
            </DialogTitle>
            <DialogDescription className="text-xs text-neutral-500">
              {t(
                'iam.roles.dialog_desc',
                'Configure fine-grained permissions for this organizational role. Custom roles store explicit capabilities.'
              )}
            </DialogDescription>
          </DialogHeader>

          {errorMessage && (
            <div className="mx-6 my-2 p-3 bg-red-50 border border-red-200 rounded-md flex items-start gap-2 text-xs text-red-700">
              <AlertCircle className="h-4 w-4 shrink-0 text-red-600 mt-0.5" />
              <span>{errorMessage}</span>
            </div>
          )}

          <div className="flex-1 overflow-y-auto px-6 py-2 space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label htmlFor="role-name" className="text-xs font-medium text-neutral-700">
                  {t('iam.roles.name', 'Role Name')} *
                </Label>
                <Input
                  id="role-name"
                  value={name}
                  onChange={(e) => handleNameChange(e.target.value)}
                  placeholder="e.g. Compliance Auditor"
                  className="mt-1 h-9 text-xs"
                  required
                />
              </div>

              <div>
                <Label htmlFor="role-slug" className="text-xs font-medium text-neutral-700">
                  {t('iam.roles.slug', 'Role Identifier (Slug)')} *
                </Label>
                <Input
                  id="role-slug"
                  value={slug}
                  onChange={(e) => setSlug(e.target.value)}
                  placeholder="e.g. compliance-auditor"
                  disabled={mode === 'edit'}
                  className="mt-1 h-9 text-xs font-mono disabled:bg-neutral-100"
                  required
                />
              </div>
            </div>

            <div>
              <Label htmlFor="role-desc" className="text-xs font-medium text-neutral-700">
                {t('iam.roles.description', 'Description')}
              </Label>
              <Input
                id="role-desc"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Briefly describe what members with this role can do"
                className="mt-1 h-9 text-xs"
              />
            </div>

            <div className="pt-2">
              <div className="flex items-center justify-between pb-2 border-b border-neutral-200">
                <span className="text-xs font-semibold text-neutral-900 uppercase tracking-wider">
                  {t('iam.roles.capabilities_title', 'Role Capabilities & Modules')}
                </span>
                <span className="text-xs text-neutral-500">
                  {selectedPerms.length} {t('iam.roles.selected_count', 'selected')}
                </span>
              </div>

              <div className="mt-3 space-y-4">
                {groupedPermissions.map(([moduleName, perms]) => {
                  const allowedModuleKeys = perms
                    .map((p) => p.key)
                    .filter(isPermAllowedForActor);
                  const allSelected =
                    allowedModuleKeys.length > 0 &&
                    allowedModuleKeys.every((k) => selectedPerms.includes(k));

                  return (
                    <div
                      key={moduleName}
                      className="border border-neutral-200 rounded-lg p-3 bg-neutral-50/50"
                    >
                      <div className="flex items-center justify-between pb-2 border-b border-neutral-100">
                        <span className="text-xs font-bold text-neutral-800 uppercase tracking-wide">
                          {moduleName}
                        </span>
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => handleToggleModule(perms)}
                          className="h-6 px-2 text-[11px] text-indigo-600 hover:text-indigo-800"
                        >
                          {allSelected
                            ? t('common.deselect_all', 'Deselect All')
                            : t('common.select_all', 'Select All')}
                        </Button>
                      </div>

                      <div className="grid grid-cols-1 md:grid-cols-2 gap-2 mt-2">
                        {perms.map((p) => {
                          const allowed = isPermAllowedForActor(p.key);
                          const isChecked = selectedPerms.includes(p.key);

                          return (
                            <label
                              key={p.key}
                              className={`flex items-start gap-2 p-2 rounded border text-xs cursor-pointer transition-colors ${
                                isChecked
                                  ? 'bg-indigo-50 border-indigo-200 text-indigo-900'
                                  : 'bg-white border-neutral-200 text-neutral-700 hover:bg-neutral-50'
                              } ${!allowed ? 'opacity-50 cursor-not-allowed bg-neutral-100' : ''}`}
                            >
                              <input
                                type="checkbox"
                                checked={isChecked}
                                disabled={!allowed}
                                onChange={() => handleTogglePerm(p.key)}
                                className="mt-0.5 rounded border-neutral-300 text-indigo-600 focus:ring-indigo-500"
                              />
                              <div className="flex-1">
                                <div className="flex items-center gap-1.5 font-medium">
                                  <span>{t(p.i18n_name_key, p.key)}</span>
                                  {!allowed && (
                                    <span title={t('iam.roles.subset_rule_tooltip', 'You do not have permission to grant this capability.')}>
                                      <Lock className="h-3 w-3 text-neutral-400" />
                                    </span>
                                  )}
                                </div>
                                <div className="text-[11px] text-neutral-500 mt-0.5">
                                  {t(p.i18n_desc_key, p.key)}
                                </div>
                              </div>
                            </label>
                          );
                        })}
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          </div>

          <DialogFooter className="p-4 bg-neutral-50 border-t border-neutral-200 flex items-center justify-end gap-2">
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => onOpenChange(false)}
              disabled={isSubmitting}
            >
              {t('common.cancel', 'Cancel')}
            </Button>
            <Button
              type="submit"
              size="sm"
              disabled={isSubmitting}
              className="bg-indigo-600 hover:bg-indigo-700 text-white flex items-center gap-1.5"
            >
              {isSubmitting ? (
                <span>{t('common.saving', 'Saving...')}</span>
              ) : (
                <>
                  <Save className="h-4 w-4" />
                  <span>
                    {mode === 'create' && t('iam.roles.btn_create', 'Create Role')}
                    {mode === 'edit' && t('iam.roles.btn_save', 'Save Changes')}
                    {mode === 'clone' && t('iam.roles.btn_clone', 'Clone Role')}
                  </span>
                </>
              )}
            </Button>
          </DialogFooter>
        </form>
      </div>
    </Dialog>
  );
}
