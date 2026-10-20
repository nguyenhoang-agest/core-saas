'use client';

import React, { useState } from 'react';
import {
  Shield,
  Lock,
  Plus,
  Edit2,
  Trash2,
  Copy,
  Users,
  AlertTriangle,
  Key,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  Dialog,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
  DialogClose,
} from '@/components/ui/dialog';
import { useI18n } from '@/core/i18n/i18n-context';
import { useRoles, useDeleteRole } from '../hooks';
import { parseBackendError } from '../api';
import { RoleEditorDialog } from './role-editor-dialog';
import type { RoleItem } from '../types';

interface RoleListTableProps {
  actorPermissions?: string[];
}

export function RoleListTable({ actorPermissions = ['*'] }: RoleListTableProps) {
  const { t } = useI18n();
  const { data: rolesData = [], isLoading } = useRoles();
  const roles = Array.isArray(rolesData) ? rolesData : [];
  const deleteMutation = useDeleteRole();

  // Role Editor Dialog state
  const [editorOpen, setEditorOpen] = useState(false);
  const [editorMode, setEditorMode] = useState<'create' | 'edit' | 'clone'>('create');
  const [targetRole, setTargetRole] = useState<RoleItem | null>(null);

  // Deletion Modal state
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [roleToDelete, setRoleToDelete] = useState<RoleItem | null>(null);
  const [deleteErrorMessage, setDeleteErrorMessage] = useState<string | null>(null);

  const handleOpenCreate = () => {
    setEditorMode('create');
    setTargetRole(null);
    setEditorOpen(true);
  };

  const handleOpenEdit = (role: RoleItem) => {
    setEditorMode('edit');
    setTargetRole(role);
    setEditorOpen(true);
  };

  const handleOpenClone = (role: RoleItem) => {
    setEditorMode('clone');
    setTargetRole(role);
    setEditorOpen(true);
  };

  const handleOpenDelete = (role: RoleItem) => {
    setRoleToDelete(role);
    setDeleteErrorMessage(null);
    setDeleteModalOpen(true);
  };

  const handleConfirmDelete = async () => {
    if (!roleToDelete) return;
    setDeleteErrorMessage(null);

    try {
      await deleteMutation.mutateAsync(roleToDelete.id);
      setDeleteModalOpen(false);
      setRoleToDelete(null);
    } catch (err) {
      const parsed = parseBackendError(err);
      setDeleteErrorMessage(parsed.message);
    }
  };

  return (
    <>
      <Card className="shadow-none border border-neutral-200">
        <CardHeader className="pb-3 flex flex-row items-center justify-between">
          <div>
            <CardTitle className="text-sm font-semibold text-neutral-900 flex items-center gap-2">
              <Shield className="h-4 w-4 text-indigo-600" />
              <span>{t('iam.roles.table_title', 'Roles & Access Control')}</span>
            </CardTitle>
            <CardDescription className="text-xs text-neutral-500 mt-1">
              {t(
                'iam.roles.table_desc',
                'Manage system built-in roles and customized tenant capabilities.'
              )}
            </CardDescription>
          </div>

          <Button
            size="sm"
            onClick={handleOpenCreate}
            className="h-8 bg-indigo-600 hover:bg-indigo-700 text-white text-xs flex items-center gap-1.5"
          >
            <Plus className="h-3.5 w-3.5" />
            <span>{t('iam.roles.btn_create_role', 'Create Role')}</span>
          </Button>
        </CardHeader>

        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full text-xs text-left border-collapse">
              <thead>
                <tr className="border-y border-neutral-200 bg-neutral-50/75 text-neutral-600 font-medium">
                  <th className="py-2.5 px-4">{t('iam.roles.col_name', 'Role Name')}</th>
                  <th className="py-2.5 px-4">{t('iam.roles.col_type', 'Type')}</th>
                  <th className="py-2.5 px-4">{t('iam.roles.col_capabilities', 'Capabilities')}</th>
                  <th className="py-2.5 px-4">{t('iam.roles.col_members', 'Members')}</th>
                  <th className="py-2.5 px-4 text-right">{t('common.actions', 'Actions')}</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-neutral-200">
                {isLoading ? (
                  <tr>
                    <td colSpan={5} className="py-8 text-center text-neutral-500">
                      {t('common.loading', 'Loading roles...')}
                    </td>
                  </tr>
                ) : roles.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="py-8 text-center text-neutral-500">
                      {t('iam.roles.no_roles_found', 'No roles configured.')}
                    </td>
                  </tr>
                ) : (
                  roles.map((r) => {
                    const isSystem = r.is_system;

                    return (
                      <tr key={r.id} className="hover:bg-neutral-50/60 transition-colors">
                        <td className="py-3 px-4">
                          <div className="font-semibold text-neutral-900">{r.name}</div>
                          <div className="text-[11px] font-mono text-neutral-400 mt-0.5">
                            {r.slug}
                          </div>
                          {r.description && (
                            <div className="text-[11px] text-neutral-500 mt-1 line-clamp-1">
                              {r.description}
                            </div>
                          )}
                        </td>

                        <td className="py-3 px-4">
                          {isSystem ? (
                            <Badge
                              variant="secondary"
                              className="bg-neutral-100 text-neutral-700 border border-neutral-200 flex items-center gap-1 w-fit"
                            >
                              <Lock className="h-3 w-3 text-neutral-500" />
                              <span>{t('iam.roles.badge_system', 'System')}</span>
                            </Badge>
                          ) : (
                            <Badge
                              variant="outline"
                              className="bg-indigo-50 text-indigo-700 border-indigo-200 w-fit"
                            >
                              {t('iam.roles.badge_custom', 'Custom')}
                            </Badge>
                          )}
                        </td>

                        <td className="py-3 px-4">
                          <div className="flex items-center gap-1.5 text-neutral-700">
                            <Key className="h-3.5 w-3.5 text-neutral-400" />
                            <span className="font-medium">
                              {r.permissions.length} {t('iam.roles.permissions_count', 'permissions')}
                            </span>
                          </div>
                        </td>

                        <td className="py-3 px-4">
                          <div className="flex items-center gap-1.5 text-neutral-700">
                            <Users className="h-3.5 w-3.5 text-neutral-400" />
                            <span>
                              {r.user_count} {t('iam.roles.members_count', 'members')}
                            </span>
                          </div>
                        </td>

                        <td className="py-3 px-4 text-right">
                          <div className="flex items-center justify-end gap-1">
                            {!isSystem && (
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleOpenEdit(r)}
                                className="h-7 w-7 p-0 text-neutral-600 hover:text-indigo-600"
                                title={t('common.edit', 'Edit')}
                              >
                                <Edit2 className="h-3.5 w-3.5" />
                              </Button>
                            )}

                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleOpenClone(r)}
                              className="h-7 w-7 p-0 text-neutral-600 hover:text-indigo-600"
                              title={t('common.clone', 'Clone')}
                            >
                              <Copy className="h-3.5 w-3.5" />
                            </Button>

                            {!isSystem && (
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleOpenDelete(r)}
                                className="h-7 w-7 p-0 text-neutral-600 hover:text-red-600"
                                title={t('common.delete', 'Delete')}
                              >
                                <Trash2 className="h-3.5 w-3.5" />
                              </Button>
                            )}
                          </div>
                        </td>
                      </tr>
                    );
                  })
                )}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      {/* Role Editor Modal */}
      <RoleEditorDialog
        open={editorOpen}
        onOpenChange={setEditorOpen}
        mode={editorMode}
        role={targetRole}
        actorPermissions={actorPermissions}
      />

      {/* Role Deletion Modal */}
      <Dialog open={deleteModalOpen} onOpenChange={setDeleteModalOpen}>
        <DialogClose onClose={() => setDeleteModalOpen(false)} />
        <div>
          <DialogHeader>
            <DialogTitle className="text-sm font-semibold flex items-center gap-2 text-red-600">
              <AlertTriangle className="h-4 w-4" />
              <span>{t('iam.roles.delete_title', 'Delete Role')}</span>
            </DialogTitle>
            <DialogDescription className="text-xs text-neutral-600 mt-2">
              {roleToDelete && roleToDelete.user_count > 0 ? (
                <div className="p-3 bg-amber-50 border border-amber-200 rounded-md text-amber-800 space-y-2">
                  <div className="font-semibold text-xs">
                    {t(
                      'iam.roles.delete_blocked_title',
                      'Cannot Delete Role Currently In Use'
                    )}
                  </div>
                  <div className="text-[11px] leading-relaxed">
                    {t(
                      'iam.roles.delete_blocked_desc',
                      `This role is currently assigned to ${roleToDelete.user_count} active member(s). To protect security and data integrity, you must reassign all members to another role before this role can be deleted.`
                    )}
                  </div>
                </div>
              ) : (
                <span>
                  {t(
                    'iam.roles.delete_confirm',
                    `Are you sure you want to permanently delete the custom role "${roleToDelete?.name}"? This action cannot be undone.`
                  )}
                </span>
              )}
            </DialogDescription>
          </DialogHeader>

          {deleteErrorMessage && (
            <div className="p-2.5 bg-red-50 border border-red-200 rounded text-xs text-red-700">
              {deleteErrorMessage}
            </div>
          )}

          <DialogFooter className="gap-2 mt-2">
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => setDeleteModalOpen(false)}
            >
              {roleToDelete && roleToDelete.user_count > 0
                ? t('common.close', 'Close')
                : t('common.cancel', 'Cancel')}
            </Button>

            {roleToDelete && roleToDelete.user_count === 0 && (
              <Button
                type="button"
                size="sm"
                variant="destructive"
                disabled={deleteMutation.isPending}
                onClick={handleConfirmDelete}
              >
                {deleteMutation.isPending
                  ? t('common.deleting', 'Deleting...')
                  : t('common.confirm_delete', 'Delete Role')}
              </Button>
            )}
          </DialogFooter>
        </div>
      </Dialog>
    </>
  );
}
